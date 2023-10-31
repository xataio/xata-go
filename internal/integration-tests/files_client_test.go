package integrationtests

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
	"github.com/xataio/xata-go/xata"
)

func Test_filesClient(t *testing.T) {
	cfg, err := setupDatabase()
	if err != nil {
		t.Fatalf("unable to setup database: %v", err)
	}

	ctx := context.TODO()
	err = setupTableWithColumns(ctx, cfg)
	if err != nil {
		t.Fatalf("unable to setup table: %v", err)
	}

	t.Cleanup(func() {
		err = cleanup(cfg)
		if err != nil {
			t.Fatalf("unable to cleanup test setup: %v", err)
		}
	})

	filesCli, err := xata.NewFilesClient(xata.WithAPIKey(cfg.apiKey),
		xata.WithBaseURL(fmt.Sprintf(
			"https://%s.%s.xata.sh",
			cfg.wsID,
			cfg.region,
		)),
		xata.WithHTTPClient(retryablehttp.NewClient().StandardClient()),
	)
	if err != nil {
		t.Fatalf("unable to construct files cli: %v", err)
	}

	recordsCli, err := xata.NewRecordsClient(
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

	insertRecordRequest := generateInsertRecordRequest(cfg.databaseName, cfg.tableName)

	record, err := recordsCli.Insert(ctx, insertRecordRequest)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, record)

	t.Run("delete a file", func(t *testing.T) {
		delRes, err := filesCli.Delete(ctx, xata.DeleteFileRequest{
			BranchRequestOptional: xata.BranchRequestOptional{
				DatabaseName: xata.String(cfg.databaseName),
			},
			TableName:  cfg.tableName,
			RecordId:   record.Id,
			ColumnName: fileColumn,
		})
		assert.NoError(t, err)
		assert.Equal(t, testFileName, delRes.Name)
	})
}