// SPDX-License-Identifier: Apache-2.0

package integrationtests

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
	"github.com/xataio/xata-go/xata"
)

func Test_databasesClient(t *testing.T) {
	apiKey, found := os.LookupEnv(xata.EnvXataAPIKey)
	if !found {
		t.Skipf("%s not found in env vars", xata.EnvXataAPIKey)
	}

	ctx := context.Background()

	t.Run("should create, list, rename and delete database and list regions", func(t *testing.T) {
		httpCli := retryablehttp.NewClient().StandardClient()

		workspaceCli, err := xata.NewWorkspacesClient(
			xata.WithAPIKey(apiKey),
			xata.WithHTTPClient(httpCli),
		)
		if err != nil {
			t.Fatal(err)
		}

		testID := testIdentifier()

		ws, err := workspaceCli.Create(
			context.Background(),
			&xata.WorkspaceMeta{Name: "ws_name_" + testID},
		)
		if err != nil {
			t.Fatal(err)
		}

		t.Cleanup(func() {
			err := workspaceCli.Delete(ctx, ws.Id)
			if err != nil {
				t.Fatal(err)
			}
		})

		databaseCli, err := xata.NewDatabasesClient(
			xata.WithAPIKey(apiKey),
			xata.WithHTTPClient(httpCli),
		)
		if err != nil {
			t.Fatal(err)
		}

		listRegionsResponse, err := databaseCli.GetRegionsWithWorkspaceID(ctx, ws.Id)
		if err != nil {
			t.Fatal(err)
		}

		db, err := databaseCli.Create(ctx, xata.CreateDatabaseRequest{
			DatabaseName: "db_name_" + testID,
			WorkspaceID:  xata.String(ws.Id),
			Region:       &listRegionsResponse.Regions[0].Id,
			UI:           &xata.UI{Color: xata.String("RED")},
			BranchMetaData: &xata.BranchMetadata{
				Repository: xata.String("github.com/my/repository"),
				Branch:     xata.String("feature-branch"),
				Stage:      xata.String("testing"),
				Labels:     &[]string{"development"},
			},
		})
		if err != nil {
			t.Fatal(err)
		}

		newDBName := "test-db-name-updated"
		dbMeta, err := databaseCli.Rename(ctx, xata.RenameDatabaseRequest{
			DatabaseName: db.DatabaseName,
			NewName:      newDBName,
			WorkspaceID:  xata.String(ws.Id),
		})
		assert.NoError(t, err)
		assert.Equal(t, newDBName, dbMeta.Name)

		listResponse, err := databaseCli.ListWithWorkspaceID(ctx, ws.Id)
		assert.NoError(t, err)
		assert.Equal(t, len(listResponse.Databases), 1)

		_, err = databaseCli.Delete(ctx, xata.DeleteDatabaseRequest{
			WorkspaceID:  xata.String(ws.Id),
			DatabaseName: newDBName,
		})
		assert.NoError(t, err)
	})
}
