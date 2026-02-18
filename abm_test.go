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
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/oauth2"
)

func newTLSServerHTTPClient(server *httptest.Server) (*http.Client, error) {
	if server == nil {
		return nil, fmt.Errorf("server is required")
	}

	baseTransport, ok := server.Client().Transport.(*http.Transport)
	if !ok {
		return nil, fmt.Errorf("unexpected transport type: %T", server.Client().Transport)
	}

	transport := baseTransport.Clone()
	tlsConfig := transport.TLSClientConfig
	if tlsConfig == nil {
		tlsConfig = &tls.Config{}
	}
	tlsConfig = tlsConfig.Clone()
	tlsConfig.InsecureSkipVerify = true
	transport.TLSClientConfig = tlsConfig
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		dialer := &net.Dialer{}
		return dialer.DialContext(ctx, network, server.Listener.Addr().String())
	}

	return &http.Client{Transport: transport}, nil
}

func TestClient_FetchOrgDevicePartNumbersCanceledContext(t *testing.T) {
	ctx := t.Context()
	if err := ctx.Err(); err != nil {
		t.Fatalf("context error: %v", err)
	}

	canceledCtx, cancel := context.WithCancel(ctx)
	cancel()

	client := &Client{}
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "token"})
	_, err := client.FetchOrgDevicePartNumbers(canceledCtx, http.DefaultClient, tokenSource)
	if err == nil {
		t.Fatal("expected error for canceled context")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_FetchOrgDevicePartNumbersMissingTokenSource(t *testing.T) {
	ctx := t.Context()
	if err := ctx.Err(); err != nil {
		t.Fatalf("context error: %v", err)
	}

	client := &Client{}
	_, err := client.FetchOrgDevicePartNumbers(ctx, http.DefaultClient, nil)
	if err == nil {
		t.Fatal("expected error for missing token source")
	}
}

func TestClient_FetchOrgDevicePartNumbersSuccess(t *testing.T) {
	ctx := t.Context()
	if err := ctx.Err(); err != nil {
		t.Fatalf("context error: %v", err)
	}

	tests := map[string]struct {
		want         []string
		wantRequests int32
	}{
		"success: two pages": {
			want:         []string{"PART-001", "PART-002"},
			wantRequests: 2,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := t.Context()
			if err := ctx.Err(); err != nil {
				t.Fatalf("context error: %v", err)
			}

			var requestCount int32
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&requestCount, 1)

				if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
					w.WriteHeader(http.StatusUnauthorized)
					fmt.Fprintf(w, `{"error":"unauthorized","authorization":%q}`, got)
					return
				}

				switch r.URL.RawQuery {
				case "":
					w.Header().Set("Content-Type", "application/json")
					fmt.Fprint(w, `{"data":[{"attributes":{"partNumber":"PART-001"}}],"links":{"next":"/v1/orgDevices?page=2"}}`)
				case "page=2":
					w.Header().Set("Content-Type", "application/json")
					fmt.Fprint(w, `{"data":[{"attributes":{"partNumber":"PART-002"}}],"links":{"next":""}}`)
				default:
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, `{"error":"unexpected query: %s"}`, r.URL.RawQuery)
				}
			}))
			t.Cleanup(server.Close)

			httpClient, err := newTLSServerHTTPClient(server)
			if err != nil {
				t.Fatalf("newTLSServerHTTPClient returned error: %v", err)
			}
			tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "test-token"})
			client := &Client{}

			got, err := client.FetchOrgDevicePartNumbers(ctx, httpClient, tokenSource)
			if err != nil {
				t.Fatalf("FetchOrgDevicePartNumbers returned error: %v", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("part numbers mismatch (-want +got):\n%s", diff)
			}
			if count := atomic.LoadInt32(&requestCount); count != tt.wantRequests {
				t.Fatalf("unexpected request count: got=%d want=%d", count, tt.wantRequests)
			}
		})
	}
}
