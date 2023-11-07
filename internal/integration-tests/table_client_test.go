// SPDX-License-Identifier: Apache-2.0

package integrationtests

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/xataio/xata-go/xata"
)

func Test_tableClient(t *testing.T) {
	apiKey, found := os.LookupEnv("XATA_API_KEY")
	if !found {
		t.Skipf("%s not found in env vars", "XATA_API_KEY")
	}

	ctx := context.Background()
	t.Run("should create/delete, get schema and columns of a table and add/delete column", func(t *testing.T) {
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

		regionID := listRegionsResponse.Regions[0].Id

		db, err := databaseCli.Create(ctx, xata.CreateDatabaseRequest{
			DatabaseName: "db_name_" + testID,
			WorkspaceID:  xata.String(ws.Id),
			Region:       &regionID,
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

		t.Cleanup(func() {
			_, err = databaseCli.Delete(ctx, xata.DeleteDatabaseRequest{
				WorkspaceID:  xata.String(ws.Id),
				DatabaseName: db.DatabaseName,
			})
			if err != nil {
				t.Fatal(err)
			}
		})

		tableCli, err := xata.NewTableClient(
			xata.WithAPIKey(apiKey),
			xata.WithBaseURL(fmt.Sprintf(
				"https://%s.%s.xata.sh",
				ws.Id,
				regionID,
			)),
			xata.WithHTTPClient(retryablehttp.NewClient().StandardClient()),
		)
		if err != nil {
			t.Fatal(err)
		}

		createTableResponse, err := tableCli.Create(ctx, xata.TableRequest{
			DatabaseName: xata.String(db.DatabaseName),
			TableName:    "table-integration-test_" + testID,
		})
		if err != nil {
			t.Fatal(err)
		}

		_, err = tableCli.AddColumn(ctx, xata.AddColumnRequest{
			TableRequest: xata.TableRequest{
				TableName:    createTableResponse.TableName,
				DatabaseName: xata.String(db.DatabaseName),
			},
			Column: &xata.Column{
				Name:         stringColumn,
				Type:         xata.ColumnTypeString,
				NotNull:      xata.Bool(true),
				DefaultValue: xata.String("defaultValue"),
				Unique:       xata.Bool(false),
			},
		})
		if err != nil {
			t.Fatal(err)
		}

		schema, err := tableCli.GetSchema(ctx, xata.TableRequest{
			TableName:    createTableResponse.TableName,
			DatabaseName: xata.String(db.DatabaseName),
		})
		assert.NoError(t, err)
		assert.Equal(t, stringColumn, schema.Columns[0].Name)

		columns, err := tableCli.GetColumns(ctx, xata.TableRequest{
			TableName:    createTableResponse.TableName,
			DatabaseName: xata.String(db.DatabaseName),
		})
		assert.NoError(t, err)
		assert.Equal(t, stringColumn, columns.Columns[0].Name)

		_, err = tableCli.DeleteColumn(ctx, xata.DeleteColumnRequest{
			TableRequest: xata.TableRequest{
				TableName:    createTableResponse.TableName,
				DatabaseName: xata.String(db.DatabaseName),
			},
			ColumnName: stringColumn,
		})
		if err != nil {
			t.Fatal(err)
		}
	})
}
