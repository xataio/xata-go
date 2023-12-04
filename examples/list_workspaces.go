// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/xataio/xata-go/xata"
)

func main() {
	workspaces, err := xata.NewWorkspacesClient()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := workspaces.List(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", *resp.Workspaces[0])
}
