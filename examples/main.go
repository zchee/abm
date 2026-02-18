// Copyright 2026 The abm Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"golang.org/x/oauth2"

	"github.com/zchee/abm"
)

const (
	endpointGetOrgDevices                     = "get-org-devices"
	endpointGetOrgDevice                      = "get-org-device"
	endpointGetOrgDeviceAppleCareCoverage     = "get-org-device-applecare-coverage"
	endpointGetMDMServers                     = "get-mdm-servers"
	endpointGetMDMServerDeviceLinkages        = "get-mdm-server-device-linkages"
	endpointGetOrgDeviceAssignedServerLinkage = "get-org-device-assigned-server-linkage"
	endpointGetOrgDeviceAssignedServer        = "get-org-device-assigned-server"
	endpointCreateOrgDeviceActivity           = "create-org-device-activity"
	endpointGetOrgDeviceActivity              = "get-org-device-activity"
	endpointFetchOrgDevicePartNumbers         = "fetch-org-device-part-numbers"
)

var (
	clientID       string
	keyID          string
	privateKeyPath string
	apiBaseURL     string
	endpoint       string

	orgDeviceID         string
	mdmServerID         string
	orgDeviceActivityID string
	fieldsArg           string
	limit               int

	activityType      string
	activityDeviceIDs string
)

func init() {
	flag.StringVar(&clientID, "client-id", "", "ABM client ID")
	flag.StringVar(&keyID, "key-id", "", "ABM key ID")
	flag.StringVar(&privateKeyPath, "private-key", "", "path to private key file")
	flag.StringVar(&apiBaseURL, "api-base-url", "", "optional ABM API base URL override")
	flag.StringVar(&endpoint, "endpoint", endpointGetOrgDevices, "endpoint to call")

	flag.StringVar(&orgDeviceID, "org-device-id", "", "organization device ID")
	flag.StringVar(&mdmServerID, "mdm-server-id", "", "MDM server ID")
	flag.StringVar(&orgDeviceActivityID, "org-device-activity-id", "", "organization device activity ID")
	flag.StringVar(&fieldsArg, "fields", "", "comma-separated fields parameter")
	flag.IntVar(&limit, "limit", 0, "page size limit (0 means API default)")

	flag.StringVar(&activityType, "activity-type", string(abm.OrgDeviceActivityTypeAssignDevices), "activity type for create-org-device-activity")
	flag.StringVar(&activityDeviceIDs, "activity-device-ids", "", "comma-separated org device IDs for create-org-device-activity")

	flag.Usage = usage
}

func main() {
	flag.Parse()

	if clientID == "" {
		log.Fatal("-client-id flag is required")
	}
	if keyID == "" {
		log.Fatal("-key-id flag is required")
	}
	if privateKeyPath == "" {
		log.Fatal("-private-key flag is required")
	}

	ctx := context.Background()

	assertion, err := abm.NewAssertion(ctx, clientID, keyID, privateKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	ts, err := abm.NewTokenSource(ctx, nil, clientID, assertion, "")
	if err != nil {
		log.Fatal(err)
	}

	client, err := newABMClient(ts)
	if err != nil {
		log.Fatal(err)
	}

	response, err := runEndpoint(ctx, client)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.MarshalWrite(os.Stdout, response, jsontext.WithIndent("  ")); err != nil {
		log.Fatal(err)
	}

	_, _ = fmt.Fprintln(os.Stdout)
}

func usage() {
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	_, _ = fmt.Fprintln(flag.CommandLine.Output(), "")
	_, _ = fmt.Fprintln(flag.CommandLine.Output(), "Supported endpoint values:")
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "  - %s\n", endpointGetOrgDevices)
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "  - %s\n", endpointGetOrgDevice)
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "  - %s\n", endpointGetOrgDeviceAppleCareCoverage)
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "  - %s\n", endpointGetMDMServers)
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "  - %s\n", endpointGetMDMServerDeviceLinkages)
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "  - %s\n", endpointGetOrgDeviceAssignedServerLinkage)
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "  - %s\n", endpointGetOrgDeviceAssignedServer)
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "  - %s\n", endpointCreateOrgDeviceActivity)
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "  - %s\n", endpointGetOrgDeviceActivity)
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "  - %s\n", endpointFetchOrgDevicePartNumbers)
}

func newABMClient(tokenSource oauth2.TokenSource) (*abm.Client, error) {
	if apiBaseURL == "" {
		return abm.NewClient(nil, tokenSource)
	}

	return abm.NewClientWithBaseURL(nil, tokenSource, apiBaseURL)
}

func runEndpoint(ctx context.Context, client *abm.Client) (any, error) {
	fields := splitCommaList(fieldsArg)

	switch endpoint {
	case endpointGetOrgDevices:
		return client.GetOrgDevices(ctx, &abm.GetOrgDevicesOptions{
			Fields: fields,
			Limit:  limit,
		})
	case endpointGetOrgDevice:
		if orgDeviceID == "" {
			return nil, fmt.Errorf("-org-device-id is required for %s", endpointGetOrgDevice)
		}
		return client.GetOrgDevice(ctx, orgDeviceID, &abm.GetOrgDeviceOptions{Fields: fields})
	case endpointGetOrgDeviceAppleCareCoverage:
		if orgDeviceID == "" {
			return nil, fmt.Errorf("-org-device-id is required for %s", endpointGetOrgDeviceAppleCareCoverage)
		}
		return client.GetOrgDeviceAppleCareCoverage(ctx, orgDeviceID, &abm.GetOrgDeviceAppleCareCoverageOptions{
			Fields: fields,
			Limit:  limit,
		})
	case endpointGetMDMServers:
		return client.GetMDMServers(ctx, &abm.GetMDMServersOptions{
			Fields: fields,
			Limit:  limit,
		})
	case endpointGetMDMServerDeviceLinkages:
		if mdmServerID == "" {
			return nil, fmt.Errorf("-mdm-server-id is required for %s", endpointGetMDMServerDeviceLinkages)
		}
		return client.GetMDMServerDeviceLinkages(ctx, mdmServerID, &abm.GetMDMServerDeviceLinkagesOptions{Limit: limit})
	case endpointGetOrgDeviceAssignedServerLinkage:
		if orgDeviceID == "" {
			return nil, fmt.Errorf("-org-device-id is required for %s", endpointGetOrgDeviceAssignedServerLinkage)
		}
		return client.GetOrgDeviceAssignedServerLinkage(ctx, orgDeviceID)
	case endpointGetOrgDeviceAssignedServer:
		if orgDeviceID == "" {
			return nil, fmt.Errorf("-org-device-id is required for %s", endpointGetOrgDeviceAssignedServer)
		}
		return client.GetOrgDeviceAssignedServer(ctx, orgDeviceID, &abm.GetOrgDeviceAssignedServerOptions{Fields: fields})
	case endpointCreateOrgDeviceActivity:
		if mdmServerID == "" {
			return nil, fmt.Errorf("-mdm-server-id is required for %s", endpointCreateOrgDeviceActivity)
		}
		deviceIDs := splitCommaList(activityDeviceIDs)
		if len(deviceIDs) == 0 {
			return nil, fmt.Errorf("-activity-device-ids is required for %s", endpointCreateOrgDeviceActivity)
		}
		request := abm.OrgDeviceActivityCreateRequest{
			Data: abm.OrgDeviceActivityCreateRequestData{
				Type: "orgDeviceActivities",
				Attributes: abm.OrgDeviceActivityCreateRequestDataAttributes{
					ActivityType: abm.OrgDeviceActivityType(activityType),
				},
				Relationships: abm.OrgDeviceActivityCreateRequestDataRelationships{
					Devices: abm.OrgDeviceActivityCreateRequestDataRelationshipsDevices{
						Data: toActivityDeviceRelationships(deviceIDs),
					},
					MDMServer: abm.OrgDeviceActivityCreateRequestDataRelationshipsMDMServer{
						Data: abm.OrgDeviceActivityCreateRequestDataRelationshipsMDMServerData{
							ID:   mdmServerID,
							Type: "mdmServers",
						},
					},
				},
			},
		}
		return client.CreateOrgDeviceActivity(ctx, request)
	case endpointGetOrgDeviceActivity:
		if orgDeviceActivityID == "" {
			return nil, fmt.Errorf("-org-device-activity-id is required for %s", endpointGetOrgDeviceActivity)
		}
		return client.GetOrgDeviceActivity(ctx, orgDeviceActivityID, &abm.GetOrgDeviceActivityOptions{Fields: fields})
	case endpointFetchOrgDevicePartNumbers:
		return client.FetchOrgDevicePartNumbers(ctx)
	default:
		return nil, fmt.Errorf("unsupported -endpoint value %q", endpoint)
	}
}

func splitCommaList(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		out = append(out, trimmed)
	}

	return out
}

func toActivityDeviceRelationships(deviceIDs []string) []abm.OrgDeviceActivityCreateRequestDataRelationshipsDevicesData {
	data := make([]abm.OrgDeviceActivityCreateRequestDataRelationshipsDevicesData, len(deviceIDs))
	for i, id := range deviceIDs {
		data[i] = abm.OrgDeviceActivityCreateRequestDataRelationshipsDevicesData{
			ID:   id,
			Type: "orgDevices",
		}
	}

	return data
}
