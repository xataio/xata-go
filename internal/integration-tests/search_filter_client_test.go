package integrationtests

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
	"github.com/xataio/xata-go/xata"
	"testing"
)

func Test_searchAndFilterClient(t *testing.T) {
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

	searchFilterCli, err := xata.NewSearchAndFilterClient(xata.WithAPIKey(cfg.apiKey),
		xata.WithBaseURL(fmt.Sprintf(
			"https://%s.%s.xata.sh",
			cfg.wsID,
			cfg.region,
		)),
		xata.WithHTTPClient(retryablehttp.NewClient().StandardClient()),
	)
	if err != nil {
		t.Fatalf("unable to construct s&f cli: %v", err)
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
	assert.NoError(t, err)
	assert.NotNil(t, record)

	columns := []string{stringColumn}
	queryTableResponse, err := searchFilterCli.Query(ctx, xata.QueryTableRequest{
		SearchRequest: xata.SearchRequest{
			DatabaseName: xata.String(cfg.databaseName),
			TableName:    cfg.tableName,
		},
		Payload: xata.QueryTableRequestPayload{
			Columns: &columns,
			//Page: &xata.PageConfig{
			//	After: xata.String("cursor"),
			//},
			Consistency: xata.QueryTableRequestConsistencyStrong,
		},
	})
	assert.NoError(t, err)
	fmt.Println(queryTableResponse.Meta.Page)
	fmt.Println(queryTableResponse.Records[0])
}
