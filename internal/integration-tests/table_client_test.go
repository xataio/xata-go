// SPDX-License-Identifier: Apache-2.0

package integrationtests

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/xataio/xata-go/xata"
)

func Test_tableClient(t *testing.T) {
	cfg, err := setupDatabase()
	if err != nil {
		t.Fatalf("unable to setup database: %v", err)
	}

	t.Cleanup(func() {
		err = cleanup(cfg)
		if err != nil {
			t.Fatalf("unable to cleanup test setup: %v", err)
		}
	})

	databaseCli, err := xata.NewDatabasesClient(
		xata.WithAPIKey(cfg.apiKey),
		xata.WithBaseURL(fmt.Sprintf(
			"https://%s.%s.xata.sh",
			cfg.wsID,
			cfg.region,
		)),
		xata.WithHTTPClient(retryablehttp.NewClient().StandardClient()),
	)
	if err != nil {
		t.Fatal(err)
	}

	tableCli, err := xata.NewTableClient(
		xata.WithAPIKey(cfg.apiKey),
		xata.WithBaseURL(fmt.Sprintf(
			"https://%s.%s.xata.sh",
			cfg.wsID,
			cfg.region,
		)),
		xata.WithHTTPClient(retryablehttp.NewClient().StandardClient()),
	)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	t.Run("should create/delete, get schema and columns of a table and add/delete column", func(t *testing.T) {
		testID := testIdentifier()
		db, err := databaseCli.Create(ctx, xata.CreateDatabaseRequest{
			DatabaseName: "db_name_" + testID,
			WorkspaceID:  xata.String(cfg.wsID),
		})
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

		t.Cleanup(func() {
			_, err = databaseCli.Delete(ctx, xata.DeleteDatabaseRequest{
				WorkspaceID:  xata.String(cfg.wsID),
				DatabaseName: db.DatabaseName,
			})
			if err != nil {
				t.Fatal(err)
			}
		})
	})
}
