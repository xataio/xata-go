<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="./assets/logo_dark.svg">
    <source media="(prefers-color-scheme: light)" srcset="./assets/logo_light.svg">
    <img width="400" alt="Xata" src="./assets/logo_dark.svg">
  </picture>
</p>

# Golang SDK for Xata

Simple Golang client for xata.io databases.

Xata is a Serverless Database that is as easy to use as a spreadsheet, has the
data integrity of PostgresSQL, and the search and analytics functionality of
Elasticsearch.

To install, run:

Assuming that the API key is set as an env var: `XATA_API_KEY=api-key-value`
```Go
package main

import (
	"context"
	"fmt"
	"log"
	
	"github.com/xataio/xata-go/xata"
)

func main() {
	workspaceCli, err := xata.NewWorkspacesClient() 
	if err != nil {
		log.Fatal(err)
	}

	resp, err := workspaceCli.List(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", *resp.Workspaces[0])
	// Output: api.GetWorkspacesListResponseWorkspacesItem{ID:"Alice-s-workspace-abcd", Name:"Alice's workspace", Slug:"Alice-s-workspace", Role:0x1}

	item := *resp.Workspaces[0]
	fmt.Printf("%s\n", item.Role.String())
	// Output: owner
}
```

The API key can also be provided as a parameter to the client constructor:
```Go
workspaceCli, err := xata.NewWorkspacesClient(xata.WithAPIKey("my-api-key"))
```

To learn more about Xata, visit [xata.io](https://xata.io).

- API Reference: https://xata.io/docs/rest-api/contexts#openapi-specifications

## Development

### Requirements

- Go 1.21.0+
- Docker
- Make
- [fern](https://docs.buildwithfern.com/overview/cli/cli) (only if auto code generation is needed)

### Tests

```shell
make test
```

```shell
make integration-test
```

### Linting

```shell
make lint
```
