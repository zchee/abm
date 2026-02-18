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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/go-cmp/cmp"
)

func TestParseECDSAPrivateKeyFromPEM(t *testing.T) {
	ctx := t.Context()
	if err := ctx.Err(); err != nil {
		t.Fatalf("context error: %v", err)
	}

	p256Key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate P-256 key: %v", err)
	}
	p256SEC1, err := x509.MarshalECPrivateKey(p256Key)
	if err != nil {
		t.Fatalf("marshal P-256 EC key: %v", err)
	}
	p256SEC1PEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: p256SEC1,
	})
	p256PKCS8, err := x509.MarshalPKCS8PrivateKey(p256Key)
	if err != nil {
		t.Fatalf("marshal P-256 PKCS8 key: %v", err)
	}
	p256PKCS8PEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: p256PKCS8,
	})

	p384Key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatalf("generate P-384 key: %v", err)
	}
	p384SEC1, err := x509.MarshalECPrivateKey(p384Key)
	if err != nil {
		t.Fatalf("marshal P-384 EC key: %v", err)
	}
	p384SEC1PEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: p384SEC1,
	})

	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate RSA key: %v", err)
	}
	rsaPKCS8, err := x509.MarshalPKCS8PrivateKey(rsaKey)
	if err != nil {
		t.Fatalf("marshal RSA PKCS8 key: %v", err)
	}
	rsaPKCS8PEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: rsaPKCS8,
	})

	tests := map[string]struct {
		pemBytes []byte
		wantErr  bool
	}{
		"success: sec1 EC key": {
			pemBytes: p256SEC1PEM,
		},
		"success: pkcs8 EC key": {
			pemBytes: p256PKCS8PEM,
		},
		"error: wrong curve": {
			pemBytes: p384SEC1PEM,
			wantErr:  true,
		},
		"error: non-EC key": {
			pemBytes: rsaPKCS8PEM,
			wantErr:  true,
		},
		"error: invalid pem": {
			pemBytes: []byte("not-a-pem"),
			wantErr:  true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := t.Context()
			if err := ctx.Err(); err != nil {
				t.Fatalf("context error: %v", err)
			}

			key, err := parseECDSAPrivateKeyFromPEM(tt.pemBytes)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseECDSAPrivateKeyFromPEM error mismatch: err=%v wantErr=%v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			if key == nil {
				t.Fatal("parseECDSAPrivateKeyFromPEM returned nil key without error")
			}
			if diff := cmp.Diff(elliptic.P256().Params().Name, key.Curve.Params().Name); diff != "" {
				t.Fatalf("curve mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNewAssertion(t *testing.T) {
	ctx := t.Context()
	if err := ctx.Err(); err != nil {
		t.Fatalf("context error: %v", err)
	}

	clientID := "BUSINESSAPI.9703f56c-10ce-4876-8f59-e78e5e23a152"
	keyID := "d136aa66-0c3b-4bd4-9892-c20e8db024ab"

	p256Key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate P-256 key: %v", err)
	}
	p256SEC1, err := x509.MarshalECPrivateKey(p256Key)
	if err != nil {
		t.Fatalf("marshal P-256 EC key: %v", err)
	}
	p256PEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: p256SEC1,
	})

	privateKeyPath := filepath.Join(t.TempDir(), "private-key.pem")
	if err := os.WriteFile(privateKeyPath, p256PEM, 0o600); err != nil {
		t.Fatalf("write private key: %v", err)
	}

	tokenString, err := NewAssertion(ctx, clientID, keyID, privateKeyPath)
	if err != nil {
		t.Fatalf("NewAssertion returned error: %v", err)
	}
	if tokenString == "" {
		t.Fatal("NewAssertion returned empty token")
	}

	parsedToken, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodES256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
		}
		return &p256Key.PublicKey, nil
	})
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}
	if !parsedToken.Valid {
		t.Fatal("parsed token is invalid")
	}

	kid, ok := parsedToken.Header["kid"].(string)
	if !ok {
		t.Fatalf("token header kid missing or not a string: %#v", parsedToken.Header["kid"])
	}
	if diff := cmp.Diff(keyID, kid); diff != "" {
		t.Fatalf("kid mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(jwt.SigningMethodES256.Alg(), parsedToken.Method.Alg()); diff != "" {
		t.Fatalf("alg mismatch (-want +got):\n%s", diff)
	}

	claims, ok := parsedToken.Claims.(*jwt.RegisteredClaims)
	if !ok {
		t.Fatalf("unexpected claims type: %T", parsedToken.Claims)
	}
	if diff := cmp.Diff(clientID, claims.Issuer); diff != "" {
		t.Fatalf("issuer mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(clientID, claims.Subject); diff != "" {
		t.Fatalf("subject mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(jwt.ClaimStrings{Audience}, claims.Audience); diff != "" {
		t.Fatalf("audience mismatch (-want +got):\n%s", diff)
	}
	if claims.IssuedAt == nil || claims.ExpiresAt == nil {
		t.Fatalf("missing issuedAt or expiresAt: issuedAt=%v expiresAt=%v", claims.IssuedAt, claims.ExpiresAt)
	}
	if diff := cmp.Diff(180*24*time.Hour, claims.ExpiresAt.Time.Sub(claims.IssuedAt.Time)); diff != "" {
		t.Fatalf("expiration window mismatch (-want +got):\n%s", diff)
	}
	if claims.ID == "" {
		t.Fatalf("missing jti claim")
	}

	now := time.Now().UTC()
	if claims.IssuedAt.Time.After(now.Add(2 * time.Second)) {
		t.Fatalf("issuedAt is in the future: issuedAt=%v now=%v", claims.IssuedAt.Time, now)
	}
	if claims.ExpiresAt.Time.Before(now) {
		t.Fatalf("expiresAt is in the past: expiresAt=%v now=%v", claims.ExpiresAt.Time, now)
	}
}

func TestNewAssertionErrors(t *testing.T) {
	ctx := t.Context()
	if err := ctx.Err(); err != nil {
		t.Fatalf("context error: %v", err)
	}

	clientID := "BUSINESSAPI.9703f56c-10ce-4876-8f59-e78e5e23a152"
	keyID := "d136aa66-0c3b-4bd4-9892-c20e8db024ab"

	p384Key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatalf("generate P-384 key: %v", err)
	}
	p384SEC1, err := x509.MarshalECPrivateKey(p384Key)
	if err != nil {
		t.Fatalf("marshal P-384 EC key: %v", err)
	}
	p384PEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: p384SEC1,
	})

	tests := map[string]struct {
		privateKeyPath string
		writeKey       []byte
		wantErr        bool
	}{
		"error: missing key file": {
			privateKeyPath: filepath.Join(t.TempDir(), "missing.pem"),
			wantErr:        true,
		},
		"error: wrong curve": {
			privateKeyPath: filepath.Join(t.TempDir(), "p384.pem"),
			writeKey:       p384PEM,
			wantErr:        true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := t.Context()
			if err := ctx.Err(); err != nil {
				t.Fatalf("context error: %v", err)
			}

			if len(tt.writeKey) > 0 {
				if err := os.WriteFile(tt.privateKeyPath, tt.writeKey, 0o600); err != nil {
					t.Fatalf("write key: %v", err)
				}
			}

			_, err := NewAssertion(ctx, clientID, keyID, tt.privateKeyPath)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewAssertion error mismatch: err=%v wantErr=%v", err, tt.wantErr)
			}
		})
	}
}

func TestNewAssertionCanceledContext(t *testing.T) {
	ctx := t.Context()
	if err := ctx.Err(); err != nil {
		t.Fatalf("context error: %v", err)
	}

	canceledCtx, cancel := context.WithCancel(ctx)
	cancel()

	_, err := NewAssertion(canceledCtx, "client-id", "key-id", "unused")
	if err == nil {
		t.Fatal("expected error for canceled context")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewClientCredentialsTokenSource(t *testing.T) {
	ctx := t.Context()
	if err := ctx.Err(); err != nil {
		t.Fatalf("context error: %v", err)
	}

	tests := map[string]struct {
		clientID        string
		clientAssertion string
		scope           string
		wantErr         bool
	}{
		"success: default scope": {
			clientID:        "client-id",
			clientAssertion: "assertion",
		},
		"success: custom scope": {
			clientID:        "client-id",
			clientAssertion: "assertion",
			scope:           "custom.scope",
		},
		"error: missing client ID": {
			clientAssertion: "assertion",
			wantErr:         true,
		},
		"error: missing client assertion": {
			clientID: "client-id",
			wantErr:  true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := t.Context()
			if err := ctx.Err(); err != nil {
				t.Fatalf("context error: %v", err)
			}

			source, err := NewTokenSource(ctx, http.DefaultClient, tt.clientID, tt.clientAssertion, tt.scope)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewClientCredentialsTokenSource error mismatch: err=%v wantErr=%v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if source == nil {
				t.Fatal("NewClientCredentialsTokenSource returned nil token source")
			}
		})
	}
}

func TestNewClientCredentialsTokenSourceCanceledContext(t *testing.T) {
	ctx := t.Context()
	if err := ctx.Err(); err != nil {
		t.Fatalf("context error: %v", err)
	}

	canceledCtx, cancel := context.WithCancel(ctx)
	cancel()

	_, err := newTokenSource(canceledCtx, http.DefaultClient, "client-id", "assertion", ScopeBusinessAPI)
	if err == nil {
		t.Fatal("expected error for canceled context")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClientCredentialsTokenSourceFormBody(t *testing.T) {
	ctx := t.Context()
	if err := ctx.Err(); err != nil {
		t.Fatalf("context error: %v", err)
	}

	type tokenRequest struct {
		method      string
		query       url.Values
		body        string
		contentType string
	}

	requestCh := make(chan tokenRequest, 1)
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		requestCh <- tokenRequest{
			method:      r.Method,
			query:       r.URL.Query(),
			body:        string(body),
			contentType: r.Header.Get("Content-Type"),
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"access_token":"abc123","token_type":"Bearer","expires_in":3600}`)
	}))
	t.Cleanup(server.Close)

	httpClient, err := newTLSServerHTTPClient(server)
	if err != nil {
		t.Fatalf("newTLSServerHTTPClient returned error: %v", err)
	}

	source, err := newTokenSource(ctx, httpClient, "client-id", "assertion", "business.api")
	if err != nil {
		t.Fatalf("newClientCredentialsTokenSource returned error: %v", err)
	}

	token, err := source.Token()
	if err != nil {
		t.Fatalf("Token returned error: %v", err)
	}
	if token == nil || token.AccessToken == "" {
		t.Fatalf("Token returned empty access token: %#v", token)
	}

	select {
	case req := <-requestCh:
		if diff := cmp.Diff(http.MethodPost, req.method); diff != "" {
			t.Fatalf("method mismatch (-want +got):\n%s", diff)
		}
		if req.body == "" {
			t.Fatal("expected non-empty request body")
		}
		if !strings.HasPrefix(req.contentType, "application/x-www-form-urlencoded") {
			t.Fatalf("unexpected content-type header: %q", req.contentType)
		}

		form, err := url.ParseQuery(req.body)
		if err != nil {
			t.Fatalf("parse form body: %v", err)
		}

		tests := map[string]struct {
			key   string
			value string
		}{
			"grant_type": {
				key:   "grant_type",
				value: "client_credentials",
			},
			"client_id": {
				key:   "client_id",
				value: "client-id",
			},
			"client_assertion_type": {
				key:   "client_assertion_type",
				value: ClientAssertionURI,
			},
			"client_assertion": {
				key:   "client_assertion",
				value: "assertion",
			},
			"scope": {
				key:   "scope",
				value: "business.api",
			},
		}

		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				value := form.Get(tt.key)
				if diff := cmp.Diff(tt.value, value); diff != "" {
					t.Fatalf("query param mismatch (-want +got):\n%s", diff)
				}
			})
		}

		if _, ok := form["client_secret"]; ok {
			t.Fatalf("unexpected client_secret form param: %v", form["client_secret"])
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for token request")
	}
}

func TestDecodeOrgDevices(t *testing.T) {
	ctx := t.Context()
	if err := ctx.Err(); err != nil {
		t.Fatalf("context error: %v", err)
	}

	tests := map[string]struct {
		payload  string
		want     []string
		wantNext string
		wantErr  bool
	}{
		"success: multiple devices with next": {
			payload:  `{"data":[{"attributes":{"partNumber":"PART-001"}},{"attributes":{"partNumber":""}},{"attributes":{}}],"links":{"next":"/v1/orgDevices?page=2"}}`,
			want:     []string{"PART-001", "", ""},
			wantNext: "/v1/orgDevices?page=2",
		},
		"success: empty data": {
			payload: `{"data":[],"links":{}}`,
			want:    []string{},
		},
		"error: invalid json": {
			payload: `{"data":[`,
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := t.Context()
			if err := ctx.Err(); err != nil {
				t.Fatalf("context error: %v", err)
			}

			got, next, err := decodeOrgDevices([]byte(tt.payload))
			if (err != nil) != tt.wantErr {
				t.Fatalf("decodeOrgDevices error mismatch: err=%v wantErr=%v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("part numbers mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantNext, next); diff != "" {
				t.Fatalf("next link mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestOrgDevicePartNumberPagesPagination(t *testing.T) {
	ctx := t.Context()
	if err := ctx.Err(); err != nil {
		t.Fatalf("context error: %v", err)
	}

	tests := map[string]struct {
		stopAfter    int
		want         []string
		wantRequests int32
	}{
		"success: two pages": {
			want:         []string{"PART-001", "PART-002"},
			wantRequests: 2,
		},
		"success: early stop": {
			stopAfter:    1,
			want:         []string{"PART-001"},
			wantRequests: 1,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := t.Context()
			if err := ctx.Err(); err != nil {
				t.Fatalf("context error: %v", err)
			}

			var requestCount int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&requestCount, 1)

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

			orgDevicesURL := server.URL + "/v1/orgDevices"

			var got []string
			pageCount := 0
			for page, err := range PageIterator(ctx, server.Client(), decodeOrgDevices, orgDevicesURL) {
				if err != nil {
					t.Fatalf("orgDevicePartNumberPages returned error: %v", err)
				}
				got = append(got, page...)
				pageCount++
				if tt.stopAfter > 0 && pageCount >= tt.stopAfter {
					break
				}
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
