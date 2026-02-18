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

package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"

	"github.com/zchee/abm"
)

var (
	clientID       string
	keyID          string
	privateKeyPath string
)

func init() {
	flag.StringVar(&clientID, "client-id", "", "ABM client id")
	flag.StringVar(&keyID, "key-id", "", "ABM key id")
	flag.StringVar(&privateKeyPath, "private-key", "", "path to private-key filepath")
}

func main() {
	flag.Parse()

	if clientID == "" {
		log.Fatal("-client-id flag is required")
	}
	if keyID == "" {
		log.Fatal("-key-id flag is required")
	}
	if privateKeyPath == "" {
		log.Fatal("-private-key flag is required")
	}

	ctx := context.Background()

	assertion, err := abm.NewAssertion(ctx, clientID, keyID, privateKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	ts, err := abm.NewTokenSource(ctx, nil, clientID, assertion, "")
	if err != nil {
		log.Fatal(err)
	}

	client, err := abm.NewClient(nil, ts)
	if err != nil {
		log.Fatal(err)
	}

	out, err := client.GetOrgDevices(ctx, &abm.GetOrgDevicesOptions{
		Fields: []string{
			"partNumber",
			"serialNumber",
		},
		Limit: 100,
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := json.MarshalWrite(os.Stdout, out, jsontext.WithIndent("  ")); err != nil {
		log.Fatal(err)
	}
}
