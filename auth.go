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
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	// Audience is the Apple Business Manager OAuth2 audience.
	Audience = "https://account.apple.com/auth/oauth2/v2/token"

	// TokenURL is the Apple Business Manager OAuth2 token URL.
	TokenURL = "https://account.apple.com/auth/oauth2/token"

	// ClientAssertionURI is the client assertion type for JWT bearer.
	ClientAssertionURI = "urn:ietf:params:oauth:client-assertion-type:jwt-bearer"

	// ScopeBusinessAPI is the Apple Business Manager Business API scope.
	ScopeBusinessAPI = "business.api"
)

// NewAssertion creates a signed client assertion for Apple Business Manager (ABM).
func NewAssertion(ctx context.Context, clientID, keyID, privateKeyPath string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return "", fmt.Errorf("read private key: %w", err)
	}

	ecKey, err := parseECDSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", fmt.Errorf("parse private key: %w", err)
	}

	issuedAt := time.Now().UTC()
	expiresAt := issuedAt.Add(180 * 24 * time.Hour) // 180 days
	claims := jwt.RegisteredClaims{
		Issuer:    clientID,
		Subject:   clientID,
		Audience:  jwt.ClaimStrings{Audience},
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ID:        uuid.NewString(),
	}
	token := &jwt.Token{
		Header: map[string]any{
			"typ": "JWT",
			"alg": jwt.SigningMethodES256.Alg(),
			"kid": keyID,
		},
		Claims: claims,
		Method: jwt.SigningMethodES256,
	}

	signed, err := token.SignedString(ecKey)
	if err != nil {
		return "", fmt.Errorf("sign client assertion: %w", err)
	}

	return signed, nil
}

// parseECDSAPrivateKeyFromPEM parses an ECDSA P-256 private key from PEM-encoded bytes.
// ABM private keys are stored in PKCS#8 DER format but may carry either the
// "EC PRIVATE KEY" or "PRIVATE KEY" PEM block label, so both are handled via
// x509.ParsePKCS8PrivateKey rather than x509.ParseECPrivateKey (which expects
// the SEC 1 / RFC 5915 encoding used by the "EC PRIVATE KEY" label in OpenSSL).
func parseECDSAPrivateKeyFromPEM(pemBytes []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("missing PEM block")
	}

	switch block.Type {
	case "EC PRIVATE KEY", "PRIVATE KEY":
		parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parse %q private key: %w", block.Type, err)
		}

		key, ok := parsed.(*ecdsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("unexpected private key type: %T", parsed)
		}

		if key.Curve.Params().Name != elliptic.P256().Params().Name {
			return nil, fmt.Errorf("unexpected elliptic curve: %s", key.Curve.Params().Name)
		}

		return key, nil

	default:
		return nil, fmt.Errorf("unsupported PEM block type: %q", block.Type)
	}
}

type clientCredentialsTokenSource struct {
	ctx    context.Context
	config clientcredentials.Config
}

var _ oauth2.TokenSource = (*clientCredentialsTokenSource)(nil)

// NewTokenSource returns a token source for Apple Business Manager using a JWT client assertion.
func NewTokenSource(ctx context.Context, httpClient *http.Client, clientID, clientAssertion, scope string) (oauth2.TokenSource, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if clientID == "" {
		return nil, fmt.Errorf("client ID is required")
	}
	if clientAssertion == "" {
		return nil, fmt.Errorf("client assertion is required")
	}
	if scope == "" {
		scope = ScopeBusinessAPI
	}
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	tokenCtx := context.WithValue(ctx, oauth2.HTTPClient, httpClient)

	params := url.Values{}
	params.Set("client_assertion_type", ClientAssertionURI)
	params.Set("client_assertion", clientAssertion)

	config := clientcredentials.Config{
		ClientID:       clientID,
		TokenURL:       TokenURL,
		Scopes:         []string{scope},
		EndpointParams: params,
		AuthStyle:      oauth2.AuthStyleInParams,
	}
	src := &clientCredentialsTokenSource{
		ctx:    tokenCtx,
		config: config,
	}

	return oauth2.ReuseTokenSource(nil, src), nil
}

// Token implements [oauth2.TokenSource].
func (ts *clientCredentialsTokenSource) Token() (*oauth2.Token, error) {
	if err := ts.ctx.Err(); err != nil {
		return nil, err
	}

	token, err := ts.config.Token(ts.ctx)
	if err != nil {
		return nil, fmt.Errorf("token request: %w", err)
	}

	return token, nil
}
