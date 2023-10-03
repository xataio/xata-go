# Golang SDK for Xata

Simple Golang client for xata.io databases. Currently work in progress.

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

## Code generation with [Fern](https://github.com/fern-api/fern)
```shell
cd xata/internal
mkdir fern-sql
cd fern-sql
fern fern init --openapi https://xata.io/api/openapi\?scope\=sql
fern add fern-go-sdk
# delete typescript from fern/api/generators.yaml
# update the import path for go in fern/api/generators.yaml
# importPath: github.com/github-user/xata-go/xata/internal/fern-sql/generated/go
fern generate --log-level debug --local
```
