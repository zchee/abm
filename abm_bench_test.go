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
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/oauth2"
)

func BenchmarkDecodeOrgDevices(b *testing.B) {
	ctx := b.Context()
	if err := ctx.Err(); err != nil {
		b.Fatalf("context error: %v", err)
	}

	payloadSizes := map[string]int{
		"small_25":   25,
		"medium_200": 200,
		"large_1000": 1000,
	}

	for name, deviceCount := range payloadSizes {
		b.Run(name, func(b *testing.B) {
			ctx := b.Context()
			if err := ctx.Err(); err != nil {
				b.Fatalf("context error: %v", err)
			}

			payload := buildOrgDevicesPayload(deviceCount, "/v1/orgDevices?page=next")
			wantCount := deviceCount

			b.ReportAllocs()
			b.ResetTimer()

			for b.Loop() {
				partNumbers, next, err := decodeOrgDevices(payload)
				if err != nil {
					b.Fatalf("decodeOrgDevices returned error: %v", err)
				}
				if got := len(partNumbers); got != wantCount {
					b.Fatalf("part numbers length mismatch: got=%d want=%d", got, wantCount)
				}
				if next != "/v1/orgDevices?page=next" {
					b.Fatalf("next link mismatch: got=%q want=%q", next, "/v1/orgDevices?page=next")
				}
			}
		})
	}
}

func BenchmarkClientFetchOrgDevicePartNumbers(b *testing.B) {
	ctx := b.Context()
	if err := ctx.Err(); err != nil {
		b.Fatalf("context error: %v", err)
	}

	const (
		pageSize  = 100
		pageCount = 8
	)
	wantTotal := pageSize * pageCount

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer bench-token" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, `{"error":"unauthorized","authorization":%q}`, got)
			return
		}

		pageNumber := 1
		if page := r.URL.Query().Get("page"); page != "" {
			parsed, err := strconv.Atoi(page)
			if err != nil || parsed < 1 || parsed > pageCount {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `{"error":"invalid page","page":%q}`, page)
				return
			}
			pageNumber = parsed
		}

		nextLink := ""
		if pageNumber < pageCount {
			nextLink = fmt.Sprintf("/v1/orgDevices?page=%d", pageNumber+1)
		}

		w.Header().Set("Content-Type", "application/json")
		payload := buildOrgDevicesPageJSON(pageNumber, pageSize, nextLink)
		if _, err := w.Write(payload); err != nil {
			b.Fatalf("write response payload: %v", err)
		}
	}))
	b.Cleanup(server.Close)

	httpClient, err := newTLSServerHTTPClient(server)
	if err != nil {
		b.Fatalf("newTLSServerHTTPClient returned error: %v", err)
	}

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "bench-token"})
	client, err := NewClientWithBaseURL(httpClient, tokenSource, server.URL)
	if err != nil {
		b.Fatalf("NewClientWithBaseURL returned error: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		partNumbers, err := client.FetchOrgDevicePartNumbers(ctx)
		if err != nil {
			b.Fatalf("FetchOrgDevicePartNumbers returned error: %v", err)
		}
		if got := len(partNumbers); got != wantTotal {
			b.Fatalf("part numbers length mismatch: got=%d want=%d", got, wantTotal)
		}
	}
}

func buildOrgDevicesPayload(deviceCount int, nextLink string) []byte {
	return buildOrgDevicesPageJSON(1, deviceCount, nextLink)
}

func buildOrgDevicesPageJSON(pageNumber, pageSize int, nextLink string) []byte {
	var builder strings.Builder

	builder.Grow(pageSize * 640)
	builder.WriteString(`{"data":[`)
	for i := range pageSize {
		if i > 0 {
			builder.WriteByte(',')
		}
		partNumber := fmt.Sprintf("PART-%04d-%05d", pageNumber, i+1)
		fmt.Fprintf(&builder, `{"id":"device-%d-%d","type":"orgDevices","attributes":{"partNumber":"%s","status":"ASSIGNED","productFamily":"iPhone","deviceModel":"iPhone 15 Pro","orderNumber":"ORDER-%04d","orderDateTime":"2026-01-02T03:04:05Z","updatedDateTime":"2026-01-03T04:05:06Z","addedToOrgDateTime":"2026-01-04T05:06:07Z","serialNumber":"SER-%d-%d","productType":"iPhone16,2","deviceCapacity":"256GB","purchaseSourceType":"APPLE","purchaseSourceId":"PS-%04d","imei":["123456789012345","123456789012346"],"meid":["A1000000000001"],"wifiMacAddress":"00:11:22:33:44:55","bluetoothMacAddress":"66:77:88:99:AA:BB","ethernetMacAddress":["CC:DD:EE:FF:00:11"]},"links":{"self":"/v1/orgDevices/%d"},"relationships":{"assignedServer":{"links":{"related":"/v1/servers/1"}}}}`, pageNumber, i+1, partNumber, pageNumber, pageNumber, i+1, pageNumber, pageNumber*10000+i+1)
	}
	builder.WriteString(`],"links":{"next":"`)
	builder.WriteString(nextLink)
	builder.WriteString(`"}}`)

	return []byte(builder.String())
}
