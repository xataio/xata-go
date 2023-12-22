// SPDX-License-Identifier: Apache-2.0

package integrationtests

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/xataio/xata-go/xata"
)

type config struct {
	apiKey              string
	wsID                string
	tableName           string
	databaseName        string
	region              string
	testID              string
	httpCli             *http.Client
	workspaceCliBaseURL string
}

func setupDatabase() (*config, error) {
	ctx := context.Background()
	apiKey, found := os.LookupEnv("XATA_API_KEY")
	if !found {
		return nil, fmt.Errorf("%s not found in env vars", "XATA_API_KEY")
	}
	// require workspace ID to come from the env var
	// instead of creating new workspace on each client
	wsID, found := os.LookupEnv("XATA_WORKSPACE_ID")
	if !found {
		return nil, fmt.Errorf("%s not found in env vars", "XATA_WORKSPACE_ID")
	}

	testID := testIdentifier()
	cfg := &config{
		apiKey:  apiKey,
		testID:  testID,
		wsID:    wsID,
		httpCli: retryablehttp.NewClient().StandardClient(),
	}

	databaseCli, err := xata.NewDatabasesClient(
		xata.WithAPIKey(cfg.apiKey),
		xata.WithHTTPClient(cfg.httpCli),
	)
	if err != nil {
		return nil, err
	}

	listRegionsResponse, err := databaseCli.GetRegionsWithWorkspaceID(ctx, cfg.wsID)
	if err != nil {
		return nil, err
	}

	cfg.region = listRegionsResponse.Regions[0].Id

	cfg.workspaceCliBaseURL = fmt.Sprintf(
		"https://%s.%s.xata.sh",
		cfg.wsID,
		cfg.region,
	)

	db, err := databaseCli.Create(ctx, xata.CreateDatabaseRequest{
		DatabaseName: "db" + cfg.testID,
		WorkspaceID:  xata.String(cfg.wsID),
		Region:       &cfg.region,
		UI:           &xata.UI{Color: xata.String("RED")},
		BranchMetaData: &xata.BranchMetadata{
			Repository: xata.String("github.com/my/repository"),
			Branch:     xata.String("feature-branch"),
			Stage:      xata.String("testing"),
			Labels:     &[]string{"development"},
		},
	})
	if err != nil {
		return nil, err
	}

	cfg.databaseName = db.DatabaseName

	return cfg, nil
}

func setupTableWithColumns(ctx context.Context, cfg *config) error {
	tableCli, err := xata.NewTableClient(
		xata.WithAPIKey(cfg.apiKey),
		xata.WithBaseURL(cfg.workspaceCliBaseURL),
		xata.WithHTTPClient(cfg.httpCli),
	)
	if err != nil {
		return err
	}

	createTableResponse, err := tableCli.Create(ctx, xata.TableRequest{
		DatabaseName: xata.String(cfg.databaseName),
		TableName:    "table-integration-test_" + cfg.testID,
	})
	if err != nil {
		return err
	}

	cfg.tableName = createTableResponse.TableName

	_, err = tableCli.AddColumn(ctx, xata.AddColumnRequest{
		TableRequest: xata.TableRequest{
			TableName:    cfg.tableName,
			DatabaseName: xata.String(cfg.databaseName),
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
		return err
	}

	_, err = tableCli.AddColumn(ctx, xata.AddColumnRequest{
		TableRequest: xata.TableRequest{
			TableName:    cfg.tableName,
			DatabaseName: xata.String(cfg.databaseName),
		},
		Column: &xata.Column{
			Name:         boolColumn,
			Type:         xata.ColumnTypeBool,
			NotNull:      xata.Bool(true),
			DefaultValue: xata.String("false"),
			Unique:       xata.Bool(false),
		},
	})
	if err != nil {
		return err
	}

	_, err = tableCli.AddColumn(ctx, xata.AddColumnRequest{
		TableRequest: xata.TableRequest{
			TableName:    cfg.tableName,
			DatabaseName: xata.String(cfg.databaseName),
		},
		Column: &xata.Column{
			Name:         textColumn,
			Type:         xata.ColumnTypeText,
			NotNull:      xata.Bool(true),
			DefaultValue: xata.String("defaultValue"),
			Unique:       xata.Bool(false),
		},
	})
	if err != nil {
		return err
	}

	_, err = tableCli.AddColumn(ctx, xata.AddColumnRequest{
		TableRequest: xata.TableRequest{
			TableName:    cfg.tableName,
			DatabaseName: xata.String(cfg.databaseName),
		},
		Column: &xata.Column{
			Name:    emailColumn,
			Type:    xata.ColumnTypeEmail,
			NotNull: xata.Bool(false),
			Unique:  xata.Bool(false),
		},
	})
	if err != nil {
		return err
	}

	_, err = tableCli.AddColumn(ctx, xata.AddColumnRequest{
		TableRequest: xata.TableRequest{
			TableName:    cfg.tableName,
			DatabaseName: xata.String(cfg.databaseName),
		},
		Column: &xata.Column{
			Name:    dateTimeColumn,
			Type:    xata.ColumnTypeDatetime,
			NotNull: xata.Bool(false),
			Unique:  xata.Bool(false),
		},
	})
	if err != nil {
		return err
	}

	_, err = tableCli.AddColumn(ctx, xata.AddColumnRequest{
		TableRequest: xata.TableRequest{
			TableName:    cfg.tableName,
			DatabaseName: xata.String(cfg.databaseName),
		},
		Column: &xata.Column{
			Name:    integerColumn,
			Type:    xata.ColumnTypeInt,
			NotNull: xata.Bool(false),
			Unique:  xata.Bool(false),
		},
	})
	if err != nil {
		return err
	}

	_, err = tableCli.AddColumn(ctx, xata.AddColumnRequest{
		TableRequest: xata.TableRequest{
			TableName:    cfg.tableName,
			DatabaseName: xata.String(cfg.databaseName),
		},
		Column: &xata.Column{
			Name:    floatColumn,
			Type:    xata.ColumnTypeFloat,
			NotNull: xata.Bool(false),
			Unique:  xata.Bool(false),
		},
	})
	if err != nil {
		return err
	}

	_, err = tableCli.AddColumn(ctx, xata.AddColumnRequest{
		TableRequest: xata.TableRequest{
			TableName:    cfg.tableName,
			DatabaseName: xata.String(cfg.databaseName),
		},
		Column: &xata.Column{
			Name: fileColumn,
			Type: xata.ColumnTypeFile,
			File: &xata.ColumnFile{
				DefaultPublicAccess: xata.Bool(true),
			},
			NotNull: xata.Bool(false),
			Unique:  xata.Bool(false),
		},
	})
	if err != nil {
		return err
	}

	_, err = tableCli.AddColumn(ctx, xata.AddColumnRequest{
		TableRequest: xata.TableRequest{
			TableName:    cfg.tableName,
			DatabaseName: xata.String(cfg.databaseName),
		},
		Column: &xata.Column{
			Name: fileArrayColumn,
			Type: xata.ColumnTypeFileMap,
			File: &xata.ColumnFile{
				DefaultPublicAccess: xata.Bool(true),
			},
			NotNull: xata.Bool(false),
			Unique:  xata.Bool(false),
		},
	})
	if err != nil {
		return err
	}

	_, err = tableCli.AddColumn(ctx, xata.AddColumnRequest{
		TableRequest: xata.TableRequest{
			TableName:    cfg.tableName,
			DatabaseName: xata.String(cfg.databaseName),
		},
		Column: &xata.Column{
			Name:    jsonColumn,
			Type:    xata.ColumnTypeJSON,
			NotNull: xata.Bool(false),
			Unique:  xata.Bool(false),
		},
	})
	if err != nil {
		return err
	}

	_, err = tableCli.AddColumn(ctx, xata.AddColumnRequest{
		TableRequest: xata.TableRequest{
			TableName:    cfg.tableName,
			DatabaseName: xata.String(cfg.databaseName),
		},
		Column: &xata.Column{
			Name: vectorColumn,
			Type: xata.ColumnTypeVector,
			Vector: &xata.ColumnVector{
				Dimension: 2,
			},
			NotNull: xata.Bool(false),
			Unique:  xata.Bool(false),
		},
	})
	if err != nil {
		return err
	}

	_, err = tableCli.AddColumn(ctx, xata.AddColumnRequest{
		TableRequest: xata.TableRequest{
			TableName:    cfg.tableName,
			DatabaseName: xata.String(cfg.databaseName),
		},
		Column: &xata.Column{
			Name:    multipleColumn,
			Type:    xata.ColumnTypeMultiple,
			NotNull: xata.Bool(false),
			Unique:  xata.Bool(false),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func cleanup(cfg *config) error {
	ctx := context.Background()

	if cfg.tableName != "" {
		tableCli, err := xata.NewTableClient(
			xata.WithAPIKey(cfg.apiKey),
			xata.WithBaseURL(cfg.workspaceCliBaseURL),
			xata.WithHTTPClient(cfg.httpCli),
		)
		if err != nil {
			return err
		}

		_, err = tableCli.Delete(ctx, xata.TableRequest{
			DatabaseName: xata.String(cfg.databaseName),
			TableName:    cfg.tableName,
		})
		if err != nil {
			return err
		}
	}

	if cfg.databaseName != "" {
		databaseCli, err := xata.NewDatabasesClient(
			xata.WithAPIKey(cfg.apiKey),
			xata.WithHTTPClient(cfg.httpCli),
		)
		if err != nil {
			return err
		}

		_, err = databaseCli.Delete(ctx, xata.DeleteDatabaseRequest{
			WorkspaceID:  xata.String(cfg.wsID),
			DatabaseName: cfg.databaseName,
		})
		if err != nil {
			return err
		}
	}
	if cfg.wsID != "" {
		workspaceCli, err := xata.NewWorkspacesClient(
			xata.WithAPIKey(cfg.apiKey),
			xata.WithHTTPClient(cfg.httpCli),
		)
		if err != nil {
			return err
		}

		err = workspaceCli.Delete(ctx, cfg.wsID)
		if err != nil {
			return err
		}
	}

	return nil
}

func testIdentifier() string {
	currentTime := time.Now()

	// Print the time
	return fmt.Sprintf(
		"integration-test_%d-%d-%d_%d_%d_%d",
		currentTime.Year(),
		currentTime.Month(),
		currentTime.Day(),
		currentTime.Hour(),
		currentTime.Minute(),
		currentTime.Second(),
	)
}

func cleanAllWorkspaces() error {
	ctx := context.Background()
	apiKey, found := os.LookupEnv("XATA_API_KEY")
	if !found {
		return fmt.Errorf("%s not found in env vars", "XATA_API_KEY")
	}

	workspaceCli, err := xata.NewWorkspacesClient(xata.WithAPIKey(apiKey), xata.WithHTTPClient(retryablehttp.NewClient().StandardClient()))
	if err != nil {
		return err
	}

	listResponse, err := workspaceCli.List(ctx)
	if err != nil {
		return err
	}

	for _, ws := range listResponse.Workspaces {
		if strings.Contains(ws.Name, "integration-test") {
			err = workspaceCli.Delete(ctx, ws.Id)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
