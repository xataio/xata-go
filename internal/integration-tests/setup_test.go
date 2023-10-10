package integrationtests

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/xataio/xata-go/xata"
)

type config struct {
	wsID         string
	tableName    string
	databaseName string
	region       string
}

func setup() (*config, error) {
	ctx := context.Background()
	apiKey, found := os.LookupEnv("XATA_API_KEY")
	if !found {
		return nil, fmt.Errorf("%s not found in env vars", "XATA_API_KEY")
	}

	testID := testIdentifier()

	cfg := &config{}

	workspaceCli, err := xata.NewWorkspacesClient(
		xata.WithAPIKey(apiKey),
		xata.WithHTTPClient(retryablehttp.NewClient().StandardClient()),
	)
	if err != nil {
		return nil, err
	}

	ws, err := workspaceCli.Create(ctx, &xata.WorkspaceMeta{Name: "ws" + testID})
	if err != nil {
		return nil, err
	}

	cfg.wsID = ws.Id

	databaseCli, err := xata.NewDatabasesClient(
		xata.WithAPIKey(apiKey),
		xata.WithHTTPClient(retryablehttp.NewClient().StandardClient()),
	)
	if err != nil {
		return nil, err
	}

	listRegionsResponse, err := databaseCli.GetRegionsWithWorkspaceID(ctx, ws.Id)
	if err != nil {
		return nil, err
	}

	cfg.region = listRegionsResponse.Regions[0].Id
	db, err := databaseCli.Create(ctx, xata.CreateDatabaseRequest{
		DatabaseName: "db" + testID,
		WorkspaceID:  xata.String(ws.Id),
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

	tableCli, err := xata.NewTableClient(
		xata.WithAPIKey(apiKey),
		xata.WithBaseURL(fmt.Sprintf(
			"https://%s.%s.xata.sh",
			ws.Id,
			cfg.region,
		)),
		xata.WithHTTPClient(retryablehttp.NewClient().StandardClient()),
	)
	if err != nil {
		return nil, err
	}

	createTableResponse, err := tableCli.Create(ctx, xata.TableRequest{
		DatabaseName: xata.String(db.DatabaseName),
		TableName:    "table-integration-test_" + testID,
	})
	if err != nil {
		return nil, err
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
		return nil, err
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
		return nil, err
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
		return nil, err
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
		return nil, err
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
		return nil, err
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
		return nil, err
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
		return nil, err
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
		return nil, err
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
		return nil, err
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
		return nil, err
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
		return nil, err
	}

	return cfg, nil
}

func cleanup(cfg *config) error {
	ctx := context.Background()
	apiKey, found := os.LookupEnv("XATA_API_KEY")
	if !found {
		return fmt.Errorf("%s not found in env vars", "XATA_API_KEY")
	}

	tableCli, err := xata.NewTableClient(
		xata.WithAPIKey(apiKey),
		xata.WithBaseURL(fmt.Sprintf(
			"https://%s.%s.xata.sh",
			cfg.wsID,
			cfg.region,
		)),
		xata.WithHTTPClient(retryablehttp.NewClient().StandardClient()),
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

	databaseCli, err := xata.NewDatabasesClient(xata.WithAPIKey(apiKey), xata.WithHTTPClient(retryablehttp.NewClient().StandardClient()))
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

	workspaceCli, err := xata.NewWorkspacesClient(xata.WithAPIKey(apiKey), xata.WithHTTPClient(retryablehttp.NewClient().StandardClient()))
	if err != nil {
		return err
	}

	err = workspaceCli.Delete(ctx, cfg.wsID)
	if err != nil {
		return err
	}

	return nil
}

func testIdentifier() string {
	currentTime := time.Now()

	// Print the time
	return fmt.Sprintf(
		"-integration-test-%d-%d-%d_%d_%d_%d",
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
