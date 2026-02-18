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
	"fmt"
	"net/http"

	"github.com/go-json-experiment/json"
	"golang.org/x/oauth2"
)

// Client represents an Apple Business Manager (ABM) API client.
type Client struct {
	hc *http.Client
}

// FetchOrgDevicePartNumbers returns the orgDevices part numbers for Apple Business Manager.
func (c *Client) FetchOrgDevicePartNumbers(ctx context.Context, httpClient *http.Client, tokenSource oauth2.TokenSource) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, ctx.Err()
	}

	if tokenSource == nil {
		return nil, fmt.Errorf("token source is required")
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseTransport := httpClient.Transport
	if baseTransport == nil {
		baseTransport = http.DefaultTransport
	}
	client := *httpClient
	client.Transport = &oauth2.Transport{
		Base:   baseTransport,
		Source: tokenSource,
	}

	partNumbers := make([]string, 0, 64)

	for pagePartNumbers, err := range PageIterator(ctx, &client, decodeOrgDevices, "https://api-business.apple.com/v1/orgDevices") {
		if err != nil {
			return nil, err
		}
		partNumbers = append(partNumbers, pagePartNumbers...)
	}

	return partNumbers, nil
}

func decodeOrgDevices(payload []byte) ([]string, string, error) {
	var response struct {
		Data []struct {
			Attributes struct {
				PartNumber string `json:"partNumber"`
			} `json:"attributes"`
		} `json:"data"`
		Links struct {
			Next string `json:"next"`
		} `json:"links"`
	}
	if err := json.Unmarshal(payload, &response); err != nil {
		return nil, "", fmt.Errorf("decode org devices response: %w", err)
	}

	partNumbers := make([]string, len(response.Data))
	for i, device := range response.Data {
		partNumbers[i] = device.Attributes.PartNumber
	}

	return partNumbers, response.Links.Next, nil
}
