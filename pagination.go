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
	"io"
	"iter"
	"net/http"
	"net/url"
	"strings"
)

// maxPages is the maximum number of pages the iterator will fetch before stopping,
// matching the ABM API hard limit of 1000 pages.
const maxPages = 1000

// PageDecoderFunc is a function that decodes a paginated API response payload into type T and returns the next link.
type PageDecoderFunc[T any] func(payload []byte) (T, string, error)

// PageIterator iterates paginated API responses from the given baseURL using the provided HTTP client and decoder function.
func PageIterator[T any](ctx context.Context, client *http.Client, decoder PageDecoderFunc[T], baseURL string) iter.Seq2[T, error] {
	var zero T

	return func(yield func(T, error) bool) {
		if err := ctx.Err(); err != nil {
			yield(zero, err)
			return
		}

		nextURL := baseURL
		for page := 0; nextURL != ""; page++ {
			if err := ctx.Err(); err != nil {
				yield(zero, err)
				return
			}

			if page >= maxPages {
				yield(zero, fmt.Errorf("pagination exceeded %d pages", maxPages))
				return
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, nextURL, http.NoBody)
			if err != nil {
				yield(zero, fmt.Errorf("build paginated request: %w", err))
				return
			}

			resp, err := client.Do(req)
			if err != nil {
				yield(zero, fmt.Errorf("paginated request: %w", err))
				return
			}

			payload, readErr := io.ReadAll(resp.Body)
			resp.Body.Close()
			if readErr != nil {
				yield(zero, fmt.Errorf("read response: %w", readErr))
				return
			}
			if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
				yield(zero, fmt.Errorf("request failed: status=%s body=%s", resp.Status, strings.TrimSpace(string(payload))))
				return
			}

			data, nextLink, err := decoder(payload)
			if err != nil {
				yield(zero, err)
				return
			}

			if !yield(data, nil) {
				return
			}

			nextURL, err = resolveNextURL(req.URL, nextLink)
			if err != nil {
				yield(zero, err)
				return
			}
		}
	}
}

func resolveNextURL(baseURL *url.URL, next string) (string, error) {
	if next == "" {
		return "", nil
	}

	parsed, err := url.Parse(next)
	if err != nil {
		return "", fmt.Errorf("parse next links url: %w", err)
	}

	if parsed.IsAbs() {
		return parsed.String(), nil
	}

	return baseURL.ResolveReference(parsed).String(), nil
}
