# abm

Go client library for the Apple Business Manager API.

## Features

- JWT client assertion generation (ES256) and OAuth2 token source creation.
- Typed client methods for all currently documented Apple Business Manager REST operations:
  - GetOrgDevices
  - GetOrgDevice
  - GetOrgDeviceAppleCareCoverage
  - GetMdmServers
  - GetMdmServerDeviceLinkages
  - GetOrgDeviceAssignedServerLinkage
  - GetOrgDeviceAssignedServer
  - CreateOrgDeviceActivity
  - GetOrgDeviceActivity
- Structured request/response models for ABM resources.
- Structured API error decoding (APIError + ErrorResponse).
- Backward-compatible FetchOrgDevicePartNumbers helper.

## Installation

```bash
go get github.com/zchee/abm
```

## Quick Start

```go
package main

import (
	"context"
	"log"
	"log/slog"

	"github.com/zchee/abm"
)

func main() {
	ctx := context.Background()

	assertion, err := abm.NewAssertion(ctx, "<client-id>", "<key-id>", "/path/to/private-key.pem or $(cat /path/to/private-key.pem)")
	if err != nil {
		log.Fatal(err)
	}

	tokenSource, err := abm.NewTokenSource(ctx, nil, "<client-id>", assertion, "")
	if err != nil {
		log.Fatal(err)
	}

	client, err := abm.NewClient(nil, tokenSource)
	if err != nil {
		log.Fatal(err)
	}

	orgDevices, err := client.GetOrgDevices(ctx, &abm.GetOrgDevicesOptions{
		Fields: []string{"partNumber", "serialNumber"},
		Limit:  100,
	})
	if err != nil {
		log.Fatal(err)
	}

	slog.Info(ctx, "devices fetched", "length", len(orgDevices.Data))
}
```

## Endpoint Coverage

| Method | Path | Client Method |
| --- | --- | --- |
| GET | /v1/orgDevices | GetOrgDevices |
| GET | /v1/orgDevices/{id} | GetOrgDevice |
| GET | /v1/orgDevices/{id}/appleCareCoverage | GetOrgDeviceAppleCareCoverage |
| GET | /v1/mdmServers | GetMdmServers |
| GET | /v1/mdmServers/{id}/relationships/devices | GetMdmServerDeviceLinkages |
| GET | /v1/orgDevices/{id}/relationships/assignedServer | GetOrgDeviceAssignedServerLinkage |
| GET | /v1/orgDevices/{id}/assignedServer | GetOrgDeviceAssignedServer |
| POST | /v1/orgDeviceActivities | CreateOrgDeviceActivity |
| GET | /v1/orgDeviceActivities/{id} | GetOrgDeviceActivity |

## References

- Apple Business Manager API: https://developer.apple.com/documentation/applebusinessmanagerapi
