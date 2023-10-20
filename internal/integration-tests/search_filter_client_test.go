package integrationtests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
	"github.com/xataio/xata-go/xata"
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

	t.Run("filter and sort via string list", func(t *testing.T) {
		queryTableResponse, err := searchFilterCli.Query(ctx, xata.QueryTableRequest{
			TableRequest: xata.TableRequest{
				DatabaseName: xata.String(cfg.databaseName),
				TableName:    cfg.tableName,
			},
			Payload: xata.QueryTableRequestPayload{
				Columns:     []string{stringColumn},
				Consistency: xata.QueryTableRequestConsistencyStrong,
				Sort:        xata.NewSortExpressionFromStringList([]string{stringColumn}),
				Filter: &xata.FilterExpression{
					Exists: xata.String(stringColumn),
				},
			},
		})
		assert.NoError(t, err)
		assert.False(t, queryTableResponse.Meta.Page.More)
		assert.NotEmpty(t, (*queryTableResponse.Records[0])[stringColumn])
		assert.Empty(t, (*queryTableResponse.Records[0])[boolColumn])
	})

	t.Run("filter and sort via string sort order map", func(t *testing.T) {
		queryTableResponse, err := searchFilterCli.Query(ctx, xata.QueryTableRequest{
			TableRequest: xata.TableRequest{
				DatabaseName: xata.String(cfg.databaseName),
				TableName:    cfg.tableName,
			},
			Payload: xata.QueryTableRequestPayload{
				Columns:     []string{stringColumn},
				Consistency: xata.QueryTableRequestConsistencyStrong,
				Sort: xata.NewSortExpressionFromStringSortOrderMap(map[string]xata.SortOrder{
					stringColumn: xata.SortOrderAsc,
				}),
				Filter: &xata.FilterExpression{
					All: xata.NewFilterListFromFilterExpression(&xata.FilterExpression{
						Exists: xata.String(stringColumn),
					}),
					Any: xata.NewFilterListFromFilterExpressionList([]*xata.FilterExpression{
						{
							Exists: xata.String(boolColumn),
						},
					}),
				},
			},
		})
		assert.NoError(t, err)
		assert.False(t, queryTableResponse.Meta.Page.More)
		assert.NotEmpty(t, (*queryTableResponse.Records[0])[stringColumn])
		assert.Empty(t, (*queryTableResponse.Records[0])[boolColumn])
	})

	t.Run("sort via string sort order map list", func(t *testing.T) {
		queryTableResponse, err := searchFilterCli.Query(ctx, xata.QueryTableRequest{
			TableRequest: xata.TableRequest{
				DatabaseName: xata.String(cfg.databaseName),
				TableName:    cfg.tableName,
			},
			Payload: xata.QueryTableRequestPayload{
				Columns:     []string{stringColumn},
				Consistency: xata.QueryTableRequestConsistencyEventual,
				Sort: xata.NewSortExpressionFromStringSortOrderMapList(
					[]map[string]xata.SortOrder{
						{
							stringColumn: xata.SortOrderDesc,
						},
					}),
			},
		})
		assert.NoError(t, err)
		assert.False(t, queryTableResponse.Meta.Page.More)
		assert.NotEmpty(t, (*queryTableResponse.Records[0])[stringColumn])
		assert.Empty(t, (*queryTableResponse.Records[0])[boolColumn])
	})

	t.Run("free text search in branch", func(t *testing.T) {
		// Query can return 0 records if this case runs too fast, hence the Eventually usage
		assert.Eventually(t, func() bool {
			pref := xata.PrefixExpressionDisabled
			searchBranchResponse, err := searchFilterCli.SearchBranch(ctx, xata.SearchBranchRequest{
				TableRequest: xata.TableRequest{
					DatabaseName: xata.String(cfg.databaseName),
				},
				Payload: xata.SearchBranchRequestPayload{
					Tables: []*xata.SearchBranchRequestTablesItem{
						xata.NewSearchBranchRequestTablesItemFromString(cfg.tableName),
					},
					Query:     "test",
					Fuzziness: xata.Int(0),
					Prefix:    &pref,
					Highlight: &xata.HighlightExpression{
						Enabled:    xata.Bool(true),
						EncodeHtml: xata.Bool(true),
					},
					Page: &xata.SearchPageConfig{
						Size:   xata.Int(10),
						Offset: xata.Int(0),
					},
				},
			})
			assert.NoError(t, err)
			return searchBranchResponse.TotalCount == 1
		}, time.Second*10, time.Second)
	})
}
