# Repository Guidelines

This document provides comprehensive guidelines for contributors working on the `abm` (Apple Business Manager) Go library.

## Project Overview

This library provides a Go client for the Apple Business Manager (ABM) API, enabling authentication via JWT client assertions and retrieval of organizational device data through paginated API endpoints.

## Project Structure & Module Organization

```
abm/
├── abm.go              # Core client and device fetching logic
├── auth.go             # OAuth2 authentication and JWT token generation
├── auth_test.go        # Comprehensive authentication tests
├── pagination.go       # Generic pagination iterator using Go 1.26 iterators
├── types.go            # ABM API response types and constants
├── doc.go              # Package documentation
├── examples/           # Usage examples
│   └── main.go
├── .github/            # GitHub configuration
│   ├── dependabot.yaml
│   └── renovate.json5
├── go.mod              # Module dependencies
└── go.sum              # Dependency checksums
```

### Key Components

- **Authentication**: JWT-based OAuth2 client credentials flow with ECDSA P-256 signing
- **API Client**: HTTP client wrapper with token management
- **Pagination**: Generic iterator pattern using Go 1.26's `iter.Seq2`
- **Types**: Strongly-typed structs for ABM API resources

## Build, Test, and Development Commands

### Building

```bash
# Build the library (verify compilation)
go build ./...

# Build the example
go build ./examples
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with race detector
go test -race ./...

# Run tests with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out  # View coverage in browser
```

### Code Quality

```bash
# Format code
gofmt -s -w .

# Run static analysis
go vet ./...

# Run linter (if golangci-lint is installed)
golangci-lint run
```

### Dependency Management

```bash
# Update dependencies
go get -u ./...
go mod tidy

# Verify dependencies
go mod verify

# View dependency graph
go mod graph
```

## Coding Style & Naming Conventions

This project strictly follows the [Google Go Style Guide](https://google.github.io/styleguide/go/).

### General Principles

- **Use Go 1.26+**: This project requires Go 1.26 or higher and uses modern features like range-over-func iterators
- **Use `any` instead of `interface{}`**
- **Use generics** where appropriate (see `PageIterator` for example)
- **Format with `gofmt -s`**: Always use simplified formatting
- **Godoc comments must end with a period**

### Naming Conventions

- **Exported identifiers**: Use `PascalCase` (e.g., `Client`, `NewAssertion`)
- **Unexported identifiers**: Use `camelCase` (e.g., `decodeOrgDevices`, `parseECDSAPrivateKeyFromPEM`)
- **Constants**: Use `PascalCase` with descriptive prefixes (e.g., `StatusAssigned`, `ProductFamilyIPhone`)
- **Acronyms**: Keep uppercase in names (e.g., `ABM`, `JWT`, `ID`, `URL`)

### JSON Struct Tags

- **Primitive types**: Use non-pointer with `omitempty`
  ```go
  Color string `json:"color,omitempty"`
  ```
- **Struct types**: Use pointer with `omitzero`
  ```go
  Links *PagedDocumentLinks `json:"links,omitzero"`
  ```
- **Slices**: Use `omitempty` (slices of primitives or complex types)
  ```go
  IMEI []string `json:"imei,omitempty"`
  ```

### Code Organization

- **Expand struct fields** when it improves readability:
  ```go
  // GOOD
  HTTPOptions: genai.HTTPOptions{
      Headers: http.Header{
          "User-Agent": []string{version.UserAgent("genai")},
      },
  }

  // BAD (too compressed)
  HTTPOptions: genai.HTTPOptions{Headers: http.Header{"User-Agent": []string{version.UserAgent("genai")}}}
  ```

### Third-Party Packages

- **Use `github.com/go-json-experiment/json`** instead of `encoding/json`
- **Use `github.com/golang-jwt/jwt/v5`** for JWT operations
- **Use `golang.org/x/oauth2`** for OAuth2 flows

## Testing Guidelines

### Testing Framework

- **Use `testing` package**: Standard library testing only
- **Assertions**: Use `github.com/google/go-cmp/cmp` for comparisons
- **No mocking frameworks**: Make actual API calls when testing authentication
- **Use `t.Context()`**: Always use `t.Context()` instead of `context.Background()`

### Test Structure

All tests must follow this pattern:

```go
func TestFeatureName(t *testing.T) {
    ctx := t.Context()
    if err := ctx.Err(); err != nil {
        t.Fatalf("context error: %v", err)
    }

    tests := map[string]struct {
        input    string
        expected string
        wantErr  bool
    }{
        "success: basic case": {
            input:    "hello",
            expected: "HELLO",
        },
        "error: empty input": {
            input:   "",
            wantErr: true,
        },
    }

    for name, tt := range tests {
        t.Run(name, func(t *testing.T) {
            ctx := t.Context()
            if err := ctx.Err(); err != nil {
                t.Fatalf("context error: %v", err)
            }

            // Test logic here
            // Do NOT use: tt := tt (copying variable is unneeded in modern Go)
        })
    }
}
```

### Test Naming Conventions

Test case names must follow the pattern: `"<status>: <description>"`

- **Success cases**: `"success: basic case"`, `"success: with pagination"`
- **Error cases**: `"error: missing parameter"`, `"error: invalid format"`

### Test Coverage

- **All public functions must have tests**
- **Cover error paths**: Test both success and failure cases
- **Test edge cases**: Empty inputs, nil values, context cancellation
- **Use table-driven tests**: Map-based test structure (see example above)

### Assertions

Use `github.com/google/go-cmp/cmp` for all comparisons:

```go
if diff := cmp.Diff(want, got); diff != "" {
    t.Fatalf("mismatch (-want +got):\n%s", diff)
}
```

### Context Handling

Every test must check context errors:

```go
ctx := t.Context()
if err := ctx.Err(); err != nil {
    t.Fatalf("context error: %v", err)
}
```

### HTTP Testing

Use `httptest` for HTTP server mocking:

```go
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Handler logic
}))
t.Cleanup(server.Close)
```

## Commit & Pull Request Guidelines

### Commit Message Format

Based on the repository's commit history, follow this format:

```
<scope>: <description>

Examples:
- abm: initial implements
- go.mod: init module
- github: add .github directory
- auth: add JWT token generation
- pagination: implement generic iterator
- types: add OrgDevice response types
```

### Commit Message Rules

- **Scope prefix**: Use the affected component as prefix (e.g., `abm:`, `auth:`, `go.mod:`)
- **Lowercase description**: Keep descriptions lowercase and concise
- **Imperative mood**: Use imperative present tense ("add", not "added" or "adds")
- **No period at end**: Omit trailing periods in the subject line

### GPG Signing

All commits must be signed:

```bash
git commit --gpg-sign --signoff -m "scope: description"
```

### Pull Request Requirements

1. **Follow the PR template** (`.github/PULL_REQUEST_TEMPLATE.md`)
2. **All tests must pass**: Run `go test ./...` before submitting
3. **Code must be formatted**: Run `gofmt -s -w .`
4. **Update documentation**: If adding new public APIs, update godoc comments
5. **Add tests for new code**: Maintain or improve test coverage
6. **Link related issues**: Reference any related GitHub issues

### Code Review Checklist

Before submitting a PR, verify:

- [ ] All tests pass (`go test ./...`)
- [ ] Code is formatted (`gofmt -s -w .`)
- [ ] No linter warnings (`go vet ./...`)
- [ ] Godoc comments are complete and end with periods
- [ ] Test coverage is maintained or improved
- [ ] Commit messages follow the format
- [ ] No unnecessary dependencies added
- [ ] Error handling is comprehensive
- [ ] Context cancellation is handled properly

## Security & Configuration

### API Credentials

**Never commit sensitive data:**

- Private keys (`.pem` files)
- Client IDs or secrets
- Access tokens
- API credentials

### ECDSA Key Requirements

The library requires ECDSA P-256 (ES256) private keys in PEM format:

- **Supported formats**: `EC PRIVATE KEY` or `PRIVATE KEY` (PKCS8)
- **Required curve**: P-256 (secp256r1)
- **File permissions**: Recommend `0600` for private key files

### OAuth2 Configuration

JWT assertions must include:

- **Issuer (`iss`)**: Client ID
- **Subject (`sub`)**: Client ID
- **Audience (`aud`)**: `https://account.apple.com/auth/oauth2/v2/token`
- **Expiration (`exp`)**: Maximum 180 days from issuance
- **JWT ID (`jti`)**: Unique identifier (UUID recommended)

## Architecture Overview

### Authentication Flow

1. Load ECDSA P-256 private key from PEM file
2. Generate JWT with required claims and sign with private key
3. Create OAuth2 token source using client credentials flow with JWT assertion
4. Token source automatically refreshes tokens as needed

### Pagination Pattern

The library uses Go 1.26's range-over-func iterators for pagination:

```go
for page, err := range PageIterator(ctx, client, decoder, baseURL) {
    if err != nil {
        return err
    }
    // Process page data
}
```

**Benefits:**
- Automatic pagination handling
- Early termination with `break`
- Context cancellation support
- Generic implementation for any paginated endpoint

### Error Handling

- **Fail fast**: Return errors immediately for invalid inputs
- **Context-aware**: Check `ctx.Err()` at function entry and in loops
- **Wrapped errors**: Use `fmt.Errorf` with `%w` for error wrapping
- **Descriptive messages**: Include context in error messages

## Dependencies

Current dependencies (see `go.mod`):

- `github.com/go-json-experiment/json` - Modern JSON library
- `github.com/golang-jwt/jwt/v5` - JWT token generation and parsing
- `github.com/google/go-cmp` - Test assertions
- `github.com/google/uuid` - UUID generation for JWT IDs
- `golang.org/x/oauth2` - OAuth2 client credentials flow

### Updating Dependencies

Dependabot is configured to automatically update dependencies daily at 11:00 Asia/Tokyo time.

Manual updates:
```bash
go get -u ./...
go mod tidy
go test ./...  # Verify updates don't break tests
```

## License

This project is licensed under the Apache License 2.0. See the `LICENSE` file for details.

All source files must include the Apache 2.0 license header with SPDX identifier:

```go
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
```

## Additional Resources

- [Apple Business Manager API Documentation](https://developer.apple.com/documentation/applebusinessmanagerapi)
- [Google Go Style Guide](https://google.github.io/styleguide/go/)
- [Go 1.26 Release Notes](https://go.dev/doc/go1.26)
- [OAuth 2.0 RFC 6749](https://datatracker.ietf.org/doc/html/rfc6749)
- [JWT RFC 7519](https://datatracker.ietf.org/doc/html/rfc7519)
