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

package abm

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/oauth2"
)

func testClientForServer(t *testing.T, server *httptest.Server) *Client {
	t.Helper()

	ctx := t.Context()
	if err := ctx.Err(); err != nil {
		t.Fatalf("context error: %v", err)
	}

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "test-token"})
	client, err := NewClientWithBaseURL(server.Client(), tokenSource, server.URL)
	if err != nil {
		t.Fatalf("NewClientWithBaseURL returned error: %v", err)
	}

	return client
}

func TestNewClientWithBaseURL(t *testing.T) {
	ctx := t.Context()
	if err := ctx.Err(); err != nil {
		t.Fatalf("context error: %v", err)
	}

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "token"})

	tests := map[string]struct {
		httpClient   *http.Client
		tokenSource  oauth2.TokenSource
		baseURL      string
		wantErr      bool
		wantBaseHost string
	}{
		"success: default base url": {
			httpClient:   http.DefaultClient,
			tokenSource:  tokenSource,
			baseURL:      "",
			wantBaseHost: "api-business.apple.com",
		},
		"success: custom base url": {
			httpClient:   http.DefaultClient,
			tokenSource:  tokenSource,
			baseURL:      "https://example.test/abm",
			wantBaseHost: "example.test",
		},
		"error: missing token source": {
			httpClient: http.DefaultClient,
			baseURL:    DefaultAPIBaseURL,
			wantErr:    true,
		},
		"error: invalid base url": {
			httpClient:  http.DefaultClient,
			tokenSource: tokenSource,
			baseURL:     "://bad-url",
			wantErr:     true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := t.Context()
			if err := ctx.Err(); err != nil {
				t.Fatalf("context error: %v", err)
			}

			client, err := NewClientWithBaseURL(tt.httpClient, tt.tokenSource, tt.baseURL)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewClientWithBaseURL error mismatch: err=%v wantErr=%v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			if client == nil {
				t.Fatal("NewClientWithBaseURL returned nil client without error")
			}
			if diff := cmp.Diff(tt.wantBaseHost, client.baseURL.Host); diff != "" {
				t.Fatalf("base url host mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	ctx := t.Context()
	if err := ctx.Err(); err != nil {
		t.Fatalf("context error: %v", err)
	}

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "token"})
	client, err := NewClient(nil, tokenSource)
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}
	if diff := cmp.Diff("api-business.apple.com", client.baseURL.Host); diff != "" {
		t.Fatalf("base url host mismatch (-want +got):\n%s", diff)
	}
}

func TestClient_ABMOperationsSuccess(t *testing.T) {
	ctx := t.Context()
	if err := ctx.Err(); err != nil {
		t.Fatalf("context error: %v", err)
	}

	requestTemplate := OrgDeviceActivityCreateRequest{
		Data: OrgDeviceActivityCreateRequestData{
			Attributes: OrgDeviceActivityCreateRequestDataAttributes{
				ActivityType: OrgDeviceActivityTypeAssignDevices,
			},
			Relationships: OrgDeviceActivityCreateRequestDataRelationships{
				Devices: OrgDeviceActivityCreateRequestDataRelationshipsDevices{
					Data: []OrgDeviceActivityCreateRequestDataRelationshipsDevicesData{
						{
							ID:   "device-1",
							Type: "orgDevices",
						},
					},
				},
				MdmServer: OrgDeviceActivityCreateRequestDataRelationshipsMdmServer{
					Data: OrgDeviceActivityCreateRequestDataRelationshipsMdmServerData{
						ID:   "mdm-1",
						Type: "mdmServers",
					},
				},
			},
			Type: "orgDeviceActivities",
		},
	}

	tests := map[string]struct {
		method       string
		path         string
		query        url.Values
		statusCode   int
		responseBody string
		expectBody   *OrgDeviceActivityCreateRequest
		invoke       func(ctx context.Context, client *Client) error
	}{
		"success: get org devices": {
			method: http.MethodGet,
			path:   "/v1/orgDevices",
			query: url.Values{
				"fields[orgDevices]": []string{"partNumber,serialNumber"},
				"limit":              []string{"2"},
			},
			statusCode:   http.StatusOK,
			responseBody: `{"data":[{"id":"device-1","type":"orgDevices","attributes":{"partNumber":"PART-001"}}],"links":{"self":"https://api-business.apple.com/v1/orgDevices"}}`,
			invoke: func(ctx context.Context, client *Client) error {
				resp, err := client.GetOrgDevices(ctx, &GetOrgDevicesOptions{
					Fields: []string{"partNumber", " serialNumber ", ""},
					Limit:  2,
				})
				if err != nil {
					return err
				}
				if len(resp.Data) != 1 {
					return fmt.Errorf("unexpected data length: %d", len(resp.Data))
				}
				if diff := cmp.Diff("PART-001", resp.Data[0].Attributes.PartNumber); diff != "" {
					return fmt.Errorf("part number mismatch (-want +got):\n%s", diff)
				}
				return nil
			},
		},
		"success: get org device": {
			method:       http.MethodGet,
			path:         "/v1/orgDevices/device-1",
			query:        url.Values{"fields[orgDevices]": []string{"partNumber"}},
			statusCode:   http.StatusOK,
			responseBody: `{"data":{"id":"device-1","type":"orgDevices","attributes":{"partNumber":"PART-001"}},"links":{"self":"https://api-business.apple.com/v1/orgDevices/device-1"}}`,
			invoke: func(ctx context.Context, client *Client) error {
				resp, err := client.GetOrgDevice(ctx, "device-1", &GetOrgDeviceOptions{Fields: []string{"partNumber"}})
				if err != nil {
					return err
				}
				if diff := cmp.Diff("device-1", resp.Data.ID); diff != "" {
					return fmt.Errorf("org device id mismatch (-want +got):\n%s", diff)
				}
				return nil
			},
		},
		"success: get org device apple care coverage": {
			method: http.MethodGet,
			path:   "/v1/orgDevices/device-1/appleCareCoverage",
			query: url.Values{
				"fields[appleCareCoverage]": []string{"status"},
				"limit":                     []string{"1"},
			},
			statusCode:   http.StatusOK,
			responseBody: `{"data":[{"id":"coverage-1","type":"appleCareCoverage","attributes":{"status":"ACTIVE"}}],"links":{"self":"https://api-business.apple.com/v1/orgDevices/device-1/appleCareCoverage"}}`,
			invoke: func(ctx context.Context, client *Client) error {
				resp, err := client.GetOrgDeviceAppleCareCoverage(ctx, "device-1", &GetOrgDeviceAppleCareCoverageOptions{
					Fields: []string{"status"},
					Limit:  1,
				})
				if err != nil {
					return err
				}
				if len(resp.Data) != 1 {
					return fmt.Errorf("unexpected data length: %d", len(resp.Data))
				}
				if diff := cmp.Diff(AppleCareCoverageStatusActive, resp.Data[0].Attributes.Status); diff != "" {
					return fmt.Errorf("status mismatch (-want +got):\n%s", diff)
				}
				return nil
			},
		},
		"success: get mdm servers": {
			method: http.MethodGet,
			path:   "/v1/mdmServers",
			query: url.Values{
				"fields[mdmServers]": []string{"serverName"},
				"limit":              []string{"1"},
			},
			statusCode:   http.StatusOK,
			responseBody: `{"data":[{"id":"mdm-1","type":"mdmServers","attributes":{"serverName":"Primary MDM"}}],"links":{"self":"https://api-business.apple.com/v1/mdmServers"}}`,
			invoke: func(ctx context.Context, client *Client) error {
				resp, err := client.GetMdmServers(ctx, &GetMdmServersOptions{
					Fields: []string{"serverName"},
					Limit:  1,
				})
				if err != nil {
					return err
				}
				if len(resp.Data) != 1 {
					return fmt.Errorf("unexpected data length: %d", len(resp.Data))
				}
				if diff := cmp.Diff("mdm-1", resp.Data[0].ID); diff != "" {
					return fmt.Errorf("mdm id mismatch (-want +got):\n%s", diff)
				}
				return nil
			},
		},
		"success: get mdm server device linkages": {
			method:       http.MethodGet,
			path:         "/v1/mdmServers/mdm-1/relationships/devices",
			query:        url.Values{"limit": []string{"2"}},
			statusCode:   http.StatusOK,
			responseBody: `{"data":[{"id":"device-1","type":"orgDevices"},{"id":"device-2","type":"orgDevices"}],"links":{"self":"https://api-business.apple.com/v1/mdmServers/mdm-1/relationships/devices"}}`,
			invoke: func(ctx context.Context, client *Client) error {
				resp, err := client.GetMdmServerDeviceLinkages(ctx, "mdm-1", &GetMdmServerDeviceLinkagesOptions{Limit: 2})
				if err != nil {
					return err
				}
				if len(resp.Data) != 2 {
					return fmt.Errorf("unexpected data length: %d", len(resp.Data))
				}
				if diff := cmp.Diff("device-2", resp.Data[1].ID); diff != "" {
					return fmt.Errorf("device id mismatch (-want +got):\n%s", diff)
				}
				return nil
			},
		},
		"success: get org device assigned server linkage": {
			method:       http.MethodGet,
			path:         "/v1/orgDevices/device-1/relationships/assignedServer",
			query:        url.Values{},
			statusCode:   http.StatusOK,
			responseBody: `{"data":{"id":"mdm-1","type":"mdmServers"},"links":{"self":"https://api-business.apple.com/v1/orgDevices/device-1/relationships/assignedServer"}}`,
			invoke: func(ctx context.Context, client *Client) error {
				resp, err := client.GetOrgDeviceAssignedServerLinkage(ctx, "device-1")
				if err != nil {
					return err
				}
				if diff := cmp.Diff("mdm-1", resp.Data.ID); diff != "" {
					return fmt.Errorf("assigned server id mismatch (-want +got):\n%s", diff)
				}
				return nil
			},
		},
		"success: get org device assigned server": {
			method:       http.MethodGet,
			path:         "/v1/orgDevices/device-1/assignedServer",
			query:        url.Values{"fields[mdmServers]": []string{"serverName"}},
			statusCode:   http.StatusOK,
			responseBody: `{"data":{"id":"mdm-1","type":"mdmServers","attributes":{"serverName":"Primary MDM"}},"links":{"self":"https://api-business.apple.com/v1/orgDevices/device-1/assignedServer"}}`,
			invoke: func(ctx context.Context, client *Client) error {
				resp, err := client.GetOrgDeviceAssignedServer(ctx, "device-1", &GetOrgDeviceAssignedServerOptions{Fields: []string{"serverName"}})
				if err != nil {
					return err
				}
				if diff := cmp.Diff("Primary MDM", resp.Data.Attributes.ServerName); diff != "" {
					return fmt.Errorf("server name mismatch (-want +got):\n%s", diff)
				}
				return nil
			},
		},
		"success: create org device activity": {
			method:       http.MethodPost,
			path:         "/v1/orgDeviceActivities",
			query:        url.Values{},
			statusCode:   http.StatusCreated,
			responseBody: `{"data":{"id":"activity-1","type":"orgDeviceActivities"},"links":{"self":"https://api-business.apple.com/v1/orgDeviceActivities/activity-1"}}`,
			expectBody:   &requestTemplate,
			invoke: func(ctx context.Context, client *Client) error {
				resp, err := client.CreateOrgDeviceActivity(ctx, requestTemplate)
				if err != nil {
					return err
				}
				if diff := cmp.Diff("activity-1", resp.Data.ID); diff != "" {
					return fmt.Errorf("activity id mismatch (-want +got):\n%s", diff)
				}
				return nil
			},
		},
		"success: get org device activity": {
			method:       http.MethodGet,
			path:         "/v1/orgDeviceActivities/activity-1",
			query:        url.Values{"fields[orgDeviceActivities]": []string{"status"}},
			statusCode:   http.StatusOK,
			responseBody: `{"data":{"id":"activity-1","type":"orgDeviceActivities","attributes":{"status":"COMPLETED"}},"links":{"self":"https://api-business.apple.com/v1/orgDeviceActivities/activity-1"}}`,
			invoke: func(ctx context.Context, client *Client) error {
				resp, err := client.GetOrgDeviceActivity(ctx, "activity-1", &GetOrgDeviceActivityOptions{Fields: []string{"status"}})
				if err != nil {
					return err
				}
				if diff := cmp.Diff("COMPLETED", resp.Data.Attributes.Status); diff != "" {
					return fmt.Errorf("activity status mismatch (-want +got):\n%s", diff)
				}
				return nil
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := t.Context()
			if err := ctx.Err(); err != nil {
				t.Fatalf("context error: %v", err)
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if diff := cmp.Diff(tt.method, r.Method); diff != "" {
					t.Fatalf("method mismatch (-want +got):\n%s", diff)
				}
				if diff := cmp.Diff(tt.path, r.URL.Path); diff != "" {
					t.Fatalf("path mismatch (-want +got):\n%s", diff)
				}
				if diff := cmp.Diff(tt.query, r.URL.Query()); diff != "" {
					t.Fatalf("query mismatch (-want +got):\n%s", diff)
				}
				if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
					t.Fatalf("authorization header mismatch: got=%q want=%q", got, "Bearer test-token")
				}

				if tt.expectBody != nil {
					payload, err := io.ReadAll(r.Body)
					if err != nil {
						t.Fatalf("read request body: %v", err)
					}
					var gotReq OrgDeviceActivityCreateRequest
					if err := json.Unmarshal(payload, &gotReq); err != nil {
						t.Fatalf("unmarshal request body: %v", err)
					}
					if diff := cmp.Diff(*tt.expectBody, gotReq); diff != "" {
						t.Fatalf("request body mismatch (-want +got):\n%s", diff)
					}
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				fmt.Fprint(w, tt.responseBody)
			}))
			t.Cleanup(server.Close)

			client := testClientForServer(t, server)
			if err := tt.invoke(ctx, client); err != nil {
				t.Fatalf("invoke returned error: %v", err)
			}
		})
	}
}

func TestClient_APIError(t *testing.T) {
	ctx := t.Context()
	if err := ctx.Err(); err != nil {
		t.Fatalf("context error: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"errors":[{"code":"NOT_FOUND","detail":"device not found","status":"404","title":"Not Found"}]}`)
	}))
	t.Cleanup(server.Close)

	client := testClientForServer(t, server)
	_, err := client.GetOrgDevice(ctx, "missing-device", nil)
	if err == nil {
		t.Fatal("expected GetOrgDevice to return API error")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got: %T", err)
	}
	if diff := cmp.Diff(http.StatusNotFound, apiErr.StatusCode); diff != "" {
		t.Fatalf("status code mismatch (-want +got):\n%s", diff)
	}
	if len(apiErr.Response.Errors) != 1 {
		t.Fatalf("unexpected errors length: %d", len(apiErr.Response.Errors))
	}
	if diff := cmp.Diff("NOT_FOUND", apiErr.Response.Errors[0].Code); diff != "" {
		t.Fatalf("error code mismatch (-want +got):\n%s", diff)
	}
}

func TestClient_ParameterValidation(t *testing.T) {
	ctx := t.Context()
	if err := ctx.Err(); err != nil {
		t.Fatalf("context error: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[],"links":{"self":"https://api-business.apple.com/v1/orgDevices"}}`)
	}))
	t.Cleanup(server.Close)

	client := testClientForServer(t, server)

	tests := map[string]struct {
		invoke  func() error
		wantErr bool
	}{
		"error: missing org device id": {
			invoke: func() error {
				_, err := client.GetOrgDevice(ctx, "", nil)
				return err
			},
			wantErr: true,
		},
		"error: missing mdm server id": {
			invoke: func() error {
				_, err := client.GetMdmServerDeviceLinkages(ctx, "  ", nil)
				return err
			},
			wantErr: true,
		},
		"error: missing org device activity id": {
			invoke: func() error {
				_, err := client.GetOrgDeviceActivity(ctx, "", nil)
				return err
			},
			wantErr: true,
		},
		"error: negative limit": {
			invoke: func() error {
				_, err := client.GetOrgDevices(ctx, &GetOrgDevicesOptions{Limit: -1})
				return err
			},
			wantErr: true,
		},
		"error: too large limit": {
			invoke: func() error {
				_, err := client.GetMdmServers(ctx, &GetMdmServersOptions{Limit: 1001})
				return err
			},
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := t.Context()
			if err := ctx.Err(); err != nil {
				t.Fatalf("context error: %v", err)
			}

			err := tt.invoke()
			if (err != nil) != tt.wantErr {
				t.Fatalf("invoke error mismatch: err=%v wantErr=%v", err, tt.wantErr)
			}
		})
	}
}
