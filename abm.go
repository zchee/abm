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

	"github.com/go-json-experiment/json"
)

// FetchOrgDevicePartNumbers returns all org-device part numbers for the organization,
// automatically following pagination until all pages are consumed.
func (c *Client) FetchOrgDevicePartNumbers(ctx context.Context) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	baseURL, err := c.buildURL(orgDevicesPath, nil)
	if err != nil {
		return nil, err
	}

	partNumbers := make([]string, 0, 64)

	for pagePartNumbers, err := range PageIterator(ctx, c.httpClient, decodeOrgDevices, baseURL) {
		if err != nil {
			return nil, err
		}
		partNumbers = append(partNumbers, pagePartNumbers...)
	}

	return partNumbers, nil
}

// FetchAllOrgDevices returns all org devices for the organization,
// automatically following pagination until all pages are consumed.
// It also returns the total device count from the first page's metadata.
func (c *Client) FetchAllOrgDevices(ctx context.Context) ([]OrgDevice, int, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	baseURL, err := c.buildURL(orgDevicesPath, nil)
	if err != nil {
		return nil, 0, err
	}

	var all []OrgDevice
	var total int
	firstPage := true

	for pageResult, err := range PageIterator(ctx, c.httpClient, decodeOrgDevicesPage, baseURL) {
		if err != nil {
			return all, total, err
		}
		if firstPage {
			total = pageResult.total
			firstPage = false
		}
		all = append(all, pageResult.devices...)
	}

	return all, total, nil
}

type orgDevicesPageResult struct {
	devices []OrgDevice
	total   int
}

func decodeOrgDevicesPage(payload []byte) (orgDevicesPageResult, string, error) {
	var response OrgDevicesResponse
	if err := json.Unmarshal(payload, &response); err != nil {
		return orgDevicesPageResult{}, "", fmt.Errorf("decode org devices response: %w", err)
	}

	result := orgDevicesPageResult{
		devices: response.Data,
	}
	if response.Meta != nil {
		result.total = response.Meta.Paging.Total
	}

	return result, response.Links.Next, nil
}

func decodeOrgDevices(payload []byte) ([]string, string, error) {
	var response OrgDevicesResponse
	if err := json.Unmarshal(payload, &response); err != nil {
		return nil, "", fmt.Errorf("decode org devices response: %w", err)
	}

	partNumbers := make([]string, len(response.Data))
	for i, device := range response.Data {
		if device.Attributes != nil {
			partNumbers[i] = device.Attributes.PartNumber
		}
	}

	return partNumbers, response.Links.Next, nil
}
