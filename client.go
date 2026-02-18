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
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/go-json-experiment/json"
	"golang.org/x/oauth2"
)

const (
	// DefaultAPIBaseURL is the default Apple Business Manager API base URL.
	DefaultAPIBaseURL = "https://api-business.apple.com/"

	maxPageLimit = 1000
)

const (
	orgDevicesPath         = "v1/orgDevices"
	orgDeviceActivitiesURL = "v1/orgDeviceActivities"
	mdmServersPath         = "v1/mdmServers"
)

// Client represents an Apple Business Manager (ABM) API client.
// The embedded HTTP client is already wrapped with an OAuth2 transport and
// must not be shared with other callers after construction.
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client // authorized via oauth2.Transport
}

// APIError contains API-level error details returned from Apple Business Manager.
type APIError struct {
	StatusCode int
	Status     string
	Response   ErrorResponse
	Body       string
}

func (e *APIError) Error() string {
	if len(e.Response.Errors) > 0 {
		errItem := e.Response.Errors[0]
		if errItem.Code != "" || errItem.Detail != "" {
			return fmt.Sprintf("abm api error: status=%d code=%q detail=%q", e.StatusCode, errItem.Code, errItem.Detail)
		}
	}

	if e.Body == "" {
		return fmt.Sprintf("abm api error: status=%d", e.StatusCode)
	}

	return fmt.Sprintf("abm api error: status=%d body=%q", e.StatusCode, e.Body)
}

// GetOrgDevicesOptions contains optional query parameters for GetOrgDevices.
type GetOrgDevicesOptions struct {
	Fields []string
	Limit  int
}

// GetOrgDeviceOptions contains optional query parameters for GetOrgDevice.
type GetOrgDeviceOptions struct {
	Fields []string
}

// GetOrgDeviceAppleCareCoverageOptions contains optional query parameters for GetOrgDeviceAppleCareCoverage.
type GetOrgDeviceAppleCareCoverageOptions struct {
	Fields []string
	Limit  int
}

// GetMDMServersOptions contains optional query parameters for [Client.GetMDMServers].
type GetMDMServersOptions struct {
	Fields []string
	Limit  int
}

// GetMDMServerDeviceLinkagesOptions contains optional query parameters for [Client.GetMDMServerDeviceLinkages].
type GetMDMServerDeviceLinkagesOptions struct {
	Limit int
}

// GetOrgDeviceAssignedServerOptions contains optional query parameters for [Client.GetOrgDeviceAssignedServer].
type GetOrgDeviceAssignedServerOptions struct {
	Fields []string
}

// GetOrgDeviceActivityOptions contains optional query parameters for [Client.GetOrgDeviceActivity].
type GetOrgDeviceActivityOptions struct {
	Fields []string
}

// NewClient returns an authenticated ABM client using the default API base URL.
func NewClient(httpClient *http.Client, tokenSource oauth2.TokenSource) (*Client, error) {
	return NewClientWithBaseURL(httpClient, tokenSource, DefaultAPIBaseURL)
}

// NewClientWithBaseURL returns an authenticated ABM client using the provided API base URL.
func NewClientWithBaseURL(httpClient *http.Client, tokenSource oauth2.TokenSource, baseURL string) (*Client, error) {
	if tokenSource == nil {
		return nil, fmt.Errorf("token source is required")
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	resolvedBaseURL, err := parseBaseURL(baseURL)
	if err != nil {
		return nil, err
	}

	baseTransport := httpClient.Transport
	if baseTransport == nil {
		baseTransport = http.DefaultTransport
	}

	authorizedClient := *httpClient
	authorizedClient.Transport = &oauth2.Transport{
		Base:   baseTransport,
		Source: tokenSource,
	}

	return &Client{
		baseURL:    resolvedBaseURL,
		httpClient: &authorizedClient,
	}, nil
}

// GetOrgDevices gets a list of organization devices.
func (c *Client) GetOrgDevices(ctx context.Context, options *GetOrgDevicesOptions) (*OrgDevicesResponse, error) {
	var fields []string
	var limit int
	if options != nil {
		fields = options.Fields
		limit = options.Limit
	}

	query, err := buildFieldsAndLimitQuery("fields[orgDevices]", fields, limit)
	if err != nil {
		return nil, err
	}

	var response OrgDevicesResponse
	if err := c.doJSONRequest(ctx, http.MethodGet, orgDevicesPath, query, nil, &response, http.StatusOK); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetOrgDevice gets information for a single organization device.
func (c *Client) GetOrgDevice(ctx context.Context, orgDeviceID string, options *GetOrgDeviceOptions) (*OrgDeviceResponse, error) {
	escapedID, err := validateAndEscapeID("org device ID", orgDeviceID)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	if options != nil {
		setFieldsQuery(query, "fields[orgDevices]", options.Fields)
	}

	var response OrgDeviceResponse
	if err := c.doJSONRequest(ctx, http.MethodGet, joinPath(orgDevicesPath, escapedID), query, nil, &response, http.StatusOK); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetOrgDeviceAppleCareCoverage gets AppleCare coverage information for a single organization device.
func (c *Client) GetOrgDeviceAppleCareCoverage(ctx context.Context, orgDeviceID string, options *GetOrgDeviceAppleCareCoverageOptions) (*AppleCareCoverageResponse, error) {
	escapedID, err := validateAndEscapeID("org device ID", orgDeviceID)
	if err != nil {
		return nil, err
	}

	var fields []string
	var limit int
	if options != nil {
		fields = options.Fields
		limit = options.Limit
	}

	query, err := buildFieldsAndLimitQuery("fields[appleCareCoverage]", fields, limit)
	if err != nil {
		return nil, err
	}

	var response AppleCareCoverageResponse
	path := joinPath(orgDevicesPath, escapedID, "appleCareCoverage")
	if err := c.doJSONRequest(ctx, http.MethodGet, path, query, nil, &response, http.StatusOK); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetMDMServers gets a list of device management services.
func (c *Client) GetMDMServers(ctx context.Context, options *GetMDMServersOptions) (*MDMServersResponse, error) {
	var fields []string
	var limit int
	if options != nil {
		fields = options.Fields
		limit = options.Limit
	}

	query, err := buildFieldsAndLimitQuery("fields[mdmServers]", fields, limit)
	if err != nil {
		return nil, err
	}

	var response MDMServersResponse
	if err := c.doJSONRequest(ctx, http.MethodGet, mdmServersPath, query, nil, &response, http.StatusOK); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetMDMServerDeviceLinkages gets all org-device serial IDs linked to a device management service.
func (c *Client) GetMDMServerDeviceLinkages(ctx context.Context, mdmServerID string, options *GetMDMServerDeviceLinkagesOptions) (*MDMServerDevicesLinkagesResponse, error) {
	escapedID, err := validateAndEscapeID("mdm server ID", mdmServerID)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	if options != nil {
		if err := setLimitQuery(query, options.Limit); err != nil {
			return nil, err
		}
	}

	var response MDMServerDevicesLinkagesResponse
	path := joinPath(mdmServersPath, escapedID, "relationships", "devices")
	if err := c.doJSONRequest(ctx, http.MethodGet, path, query, nil, &response, http.StatusOK); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetOrgDeviceAssignedServerLinkage gets assigned device-management service ID linkage for a device.
func (c *Client) GetOrgDeviceAssignedServerLinkage(ctx context.Context, orgDeviceID string) (*OrgDeviceAssignedServerLinkageResponse, error) {
	escapedID, err := validateAndEscapeID("org device ID", orgDeviceID)
	if err != nil {
		return nil, err
	}

	var response OrgDeviceAssignedServerLinkageResponse
	path := joinPath(orgDevicesPath, escapedID, "relationships", "assignedServer")
	if err := c.doJSONRequest(ctx, http.MethodGet, path, nil, nil, &response, http.StatusOK); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetOrgDeviceAssignedServer gets assigned device-management service information for a device.
func (c *Client) GetOrgDeviceAssignedServer(ctx context.Context, orgDeviceID string, options *GetOrgDeviceAssignedServerOptions) (*MDMServerResponse, error) {
	escapedID, err := validateAndEscapeID("org device ID", orgDeviceID)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	if options != nil {
		setFieldsQuery(query, "fields[mdmServers]", options.Fields)
	}

	var response MDMServerResponse
	path := joinPath(orgDevicesPath, escapedID, "assignedServer")
	if err := c.doJSONRequest(ctx, http.MethodGet, path, query, nil, &response, http.StatusOK); err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateOrgDeviceActivity creates an org-device activity that assigns or unassigns devices.
func (c *Client) CreateOrgDeviceActivity(ctx context.Context, request OrgDeviceActivityCreateRequest) (*OrgDeviceActivityResponse, error) {
	var response OrgDeviceActivityResponse
	if err := c.doJSONRequest(ctx, http.MethodPost, orgDeviceActivitiesURL, nil, request, &response, http.StatusCreated); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetOrgDeviceActivity gets organization device activity information.
func (c *Client) GetOrgDeviceActivity(ctx context.Context, orgDeviceActivityID string, options *GetOrgDeviceActivityOptions) (*OrgDeviceActivityResponse, error) {
	escapedID, err := validateAndEscapeID("org device activity ID", orgDeviceActivityID)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	if options != nil {
		setFieldsQuery(query, "fields[orgDeviceActivities]", options.Fields)
	}

	var response OrgDeviceActivityResponse
	if err := c.doJSONRequest(ctx, http.MethodGet, joinPath(orgDeviceActivitiesURL, escapedID), query, nil, &response, http.StatusOK); err != nil {
		return nil, err
	}

	return &response, nil
}

func buildFieldsAndLimitQuery(fieldKey string, fields []string, limit int) (url.Values, error) {
	query := url.Values{}
	setFieldsQuery(query, fieldKey, fields)
	if err := setLimitQuery(query, limit); err != nil {
		return nil, err
	}

	return query, nil
}

func setFieldsQuery(query url.Values, key string, fields []string) {
	if len(fields) == 0 {
		return
	}

	parts := make([]string, 0, len(fields))
	for _, field := range fields {
		trimmed := strings.TrimSpace(field)
		if trimmed == "" {
			continue
		}
		parts = append(parts, trimmed)
	}
	if len(parts) == 0 {
		return
	}

	query.Set(key, strings.Join(parts, ","))
}

func setLimitQuery(query url.Values, limit int) error {
	if limit == 0 {
		return nil
	}
	if limit < 0 {
		return fmt.Errorf("limit must be >= 0: %d", limit)
	}
	if limit > maxPageLimit {
		return fmt.Errorf("limit must be <= %d: %d", maxPageLimit, limit)
	}

	query.Set("limit", strconv.Itoa(limit))
	return nil
}

func validateAndEscapeID(name, id string) (string, error) {
	trimmed := strings.TrimSpace(id)
	if trimmed == "" {
		return "", fmt.Errorf("%s is required", name)
	}

	return url.PathEscape(trimmed), nil
}

func joinPath(parts ...string) string {
	filtered := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.Trim(part, "/")
		if trimmed == "" {
			continue
		}
		filtered = append(filtered, trimmed)
	}

	return strings.Join(filtered, "/")
}

func parseBaseURL(rawBaseURL string) (*url.URL, error) {
	if rawBaseURL == "" {
		rawBaseURL = DefaultAPIBaseURL
	}

	parsed, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse base URL: %w", err)
	}
	if !parsed.IsAbs() {
		return nil, fmt.Errorf("base URL must be absolute: %q", rawBaseURL)
	}
	if parsed.Host == "" {
		return nil, fmt.Errorf("base URL host is required")
	}
	if !strings.HasSuffix(parsed.Path, "/") {
		parsed.Path += "/"
	}

	return parsed, nil
}

func (c *Client) buildURL(path string, query url.Values) (string, error) {
	base := *c.baseURL // copy to avoid mutations

	relative, err := url.Parse(strings.TrimPrefix(path, "/"))
	if err != nil {
		return "", fmt.Errorf("parse request path: %w", err)
	}

	resolved := base.ResolveReference(relative)
	if len(query) > 0 {
		resolved.RawQuery = query.Encode()
	}

	return resolved.String(), nil
}

func statusAllowed(statusCode int, expectedStatusCodes []int) bool {
	return slices.Contains(expectedStatusCodes, statusCode)
}

func decodeAPIError(resp *http.Response, payload []byte) error {
	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Body:       strings.TrimSpace(string(payload)),
	}

	if len(payload) == 0 {
		return apiErr
	}

	var errResp ErrorResponse
	if err := json.Unmarshal(payload, &errResp); err == nil {
		apiErr.Response = errResp
	}

	return apiErr
}

func (c *Client) doJSONRequest(ctx context.Context, method, path string, query url.Values, requestBody, responseBody any, expectedStatusCodes ...int) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if len(expectedStatusCodes) == 0 {
		expectedStatusCodes = []int{http.StatusOK}
	}

	requestURL, err := c.buildURL(path, query)
	if err != nil {
		return err
	}

	var body []byte
	if requestBody != nil {
		body, err = json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("encode request body: %w", err)
		}
	}

	requestReader := io.Reader(http.NoBody)
	if len(body) > 0 {
		requestReader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, requestURL, requestReader)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if !statusAllowed(resp.StatusCode, expectedStatusCodes) {
		return decodeAPIError(resp, payload)
	}

	if responseBody == nil || len(payload) == 0 {
		return nil
	}

	if err := json.Unmarshal(payload, responseBody); err != nil {
		return fmt.Errorf("decode response body: %w", err)
	}

	return nil
}
