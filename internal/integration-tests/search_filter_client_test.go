// SPDX-License-Identifier: Apache-2.0

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
			BranchRequestOptional: xata.BranchRequestOptional{
				DatabaseName: xata.String(cfg.databaseName),
			},
			TableName: cfg.tableName,
			Payload: xata.QueryTableRequestPayload{
				Columns:     []string{stringColumn},
				Consistency: xata.ConsistencyStrong,
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
			BranchRequestOptional: xata.BranchRequestOptional{
				DatabaseName: xata.String(cfg.databaseName),
			},
			TableName: cfg.tableName,
			Payload: xata.QueryTableRequestPayload{
				Columns:     []string{stringColumn},
				Consistency: xata.ConsistencyStrong,
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
			BranchRequestOptional: xata.BranchRequestOptional{
				DatabaseName: xata.String(cfg.databaseName),
			},
			TableName: cfg.tableName,
			Payload: xata.QueryTableRequestPayload{
				Columns:     []string{stringColumn},
				Consistency: xata.ConsistencyEventual,
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
				BranchRequestOptional: xata.BranchRequestOptional{
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

	t.Run("free text search in branch with default values", func(t *testing.T) {
		// Query can return 0 records if this case runs too fast, hence the Eventually usage
		assert.Eventually(t, func() bool {
			searchBranchResponse, err := searchFilterCli.SearchBranch(ctx, xata.SearchBranchRequest{
				BranchRequestOptional: xata.BranchRequestOptional{
					DatabaseName: xata.String(cfg.databaseName),
				},
				Payload: xata.SearchBranchRequestPayload{
					Query: "test",
				},
			})
			assert.NoError(t, err)
			return searchBranchResponse.TotalCount == 1
		}, time.Second*10, time.Second)
	})

	t.Run("free text search in table with default values", func(t *testing.T) {
		assert.Eventually(t, func() bool {
			searchTableResponse, err := searchFilterCli.SearchTable(ctx, xata.SearchTableRequest{
				BranchRequestOptional: xata.BranchRequestOptional{
					DatabaseName: xata.String(cfg.databaseName),
				},
				TableName: cfg.tableName,
				Payload: xata.SearchTableRequestPayload{
					Query: "test",
				},
			})
			assert.NoError(t, err)
			return searchTableResponse.TotalCount == 1
		}, time.Second*10, time.Second)
	})

	// TODO: This test takes too long - returns 503 sometimes, investigate
	//t.Run("free text search in table with all params - returns no match", func(t *testing.T) {
	//	prefix := xata.PrefixExpressionDisabled
	//	_, err = searchFilterCli.SearchTable(ctx, xata.SearchTableRequest{
	//		BranchRequestOptional: xata.BranchRequestOptional{
	//			DatabaseName: xata.String(cfg.databaseName),
	//		},
	//		TableName: cfg.tableName,
	//		Payload: xata.SearchTableRequestPayload{
	//			Query:     "test",
	//			Fuzziness: xata.Int(0),
	//			Target: []*xata.TargetExpressionItem{
	//				xata.NewTargetExpression(stringColumn),
	//				xata.NewTargetExpressionWithColumnObject(xata.TargetExpressionItemColumn{
	//					Column: stringColumn,
	//					Weight: xata.Float64(2),
	//				}),
	//			},
	//			Prefix: &prefix,
	//			Filter: &xata.FilterExpression{
	//				All: xata.NewFilterListFromFilterExpression(&xata.FilterExpression{
	//					Exists: xata.String(stringColumn),
	//				}),
	//				Any: xata.NewFilterListFromFilterExpressionList([]*xata.FilterExpression{
	//					{
	//						Exists: xata.String(boolColumn),
	//					},
	//				}),
	//			},
	//			Highlight: &xata.HighlightExpression{
	//				Enabled:    xata.Bool(true),
	//				EncodeHtml: xata.Bool(true),
	//			},
	//			Boosters: []*xata.BoosterExpression{
	//				xata.NewBoosterExpressionFromBoosterExpressionValueBooster(&xata.BoosterExpressionValueBooster{
	//					ValueBooster: &xata.ValueBooster{
	//						Column: stringColumn,
	//						Value:  xata.NewValueBoosterValueFromString("test"),
	//						Factor: 1,
	//						IfMatchesFilter: &xata.FilterExpression{
	//							All: xata.NewFilterListFromFilterExpression(&xata.FilterExpression{
	//								Exists: xata.String(stringColumn),
	//							}),
	//							Any: xata.NewFilterListFromFilterExpressionList([]*xata.FilterExpression{
	//								{
	//									Exists: xata.String(boolColumn),
	//								},
	//							}),
	//						},
	//					},
	//				}),
	//				xata.NewBoosterExpressionFromBoosterExpressionNumericBooster(&xata.BoosterExpressionNumericBooster{
	//					NumericBooster: &xata.NumericBooster{
	//						Column:   stringColumn,
	//						Factor:   2,
	//						Modifier: xata.Uint8(2),
	//						IfMatchesFilter: &xata.FilterExpression{
	//							All: xata.NewFilterListFromFilterExpression(&xata.FilterExpression{
	//								Exists: xata.String(stringColumn),
	//							}),
	//							Any: xata.NewFilterListFromFilterExpressionList([]*xata.FilterExpression{
	//								{
	//									Exists: xata.String(boolColumn),
	//								},
	//							}),
	//						},
	//					},
	//				}),
	//				xata.NewBoosterExpressionFromBoosterExpressionDateBooster(&xata.BoosterExpressionDateBooster{
	//					DateBooster: &xata.DateBooster{
	//						Column: stringColumn,
	//						Origin: xata.String("2023-01-02T15:04:05Z"),
	//						Scale:  "1d",
	//						Decay:  1,
	//						Factor: xata.Float64(2),
	//						IfMatchesFilter: &xata.FilterExpression{
	//							All: xata.NewFilterListFromFilterExpression(&xata.FilterExpression{
	//								Exists: xata.String(stringColumn),
	//							}),
	//							Any: xata.NewFilterListFromFilterExpressionList([]*xata.FilterExpression{
	//								{
	//									Exists: xata.String(boolColumn),
	//								},
	//							}),
	//						},
	//					},
	//				}),
	//			},
	//			//Page: &xata.SearchPageConfig{
	//			//	Size:   xata.Int(2),
	//			//	Offset: xata.Int(0),
	//			//},
	//		},
	//	})
	//	assert.NoError(t, err)
	//})

	t.Run("vector search", func(t *testing.T) {
		searchVectorResp, err := searchFilterCli.VectorSearch(ctx, xata.VectorSearchTableRequest{
			BranchRequestOptional: xata.BranchRequestOptional{
				DatabaseName: xata.String(cfg.databaseName),
			},
			TableName: cfg.tableName,
			Payload: xata.VectorSearchTableRequestPayload{
				QueryVector:        []float64{10, 2},
				Column:             vectorColumn,
				SimilarityFunction: xata.String("cosineSimilarity"),
				Size:               xata.Int(2),
				Filter: &xata.FilterExpression{
					All: xata.NewFilterListFromFilterExpression(&xata.FilterExpression{
						Exists: xata.String(vectorColumn),
					}),
				},
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, searchVectorResp.TotalCount)
	})

	t.Run("ask question to table and ask a follow up one", func(t *testing.T) {
		var sessionID string
		assert.Eventually(t, func() bool {
			keyword := xata.AskTableRequestSearchTypeKeyword
			askResp, err := searchFilterCli.Ask(ctx, xata.AskTableRequest{
				BranchRequestOptional: xata.BranchRequestOptional{
					DatabaseName: xata.String(cfg.databaseName),
				},
				TableName: cfg.tableName,
				Payload: xata.AskTableRequestPayload{
					Question:   "What is atom?",
					SearchType: &keyword,
				},
			})

			if askResp != nil && askResp.SessionId != "" {
				sessionID = askResp.SessionId
			}

			return assert.NoError(t, err) && assert.NotEmpty(t, askResp.Answer) && assert.NotEmpty(t, askResp.SessionId)
		}, time.Second*10, time.Second)

		askFollowUpRes, err := searchFilterCli.AskFollowUp(ctx, xata.AskFollowUpRequest{
			BranchRequestOptional: xata.BranchRequestOptional{
				DatabaseName: xata.String(cfg.databaseName),
			},
			TableName: cfg.tableName,
			SessionID: sessionID,
			Question:  "Give me more info about atom, please.",
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, askFollowUpRes)
	})

	t.Run("summarize table", func(t *testing.T) {
		sumTableRes, err := searchFilterCli.Summarize(
			ctx,
			xata.SummarizeTableRequest{
				BranchRequestOptional: xata.BranchRequestOptional{
					DatabaseName: xata.String(cfg.databaseName),
				},
				TableName: cfg.tableName,
				Payload: xata.SummarizeTableRequestPayload{
					Columns: []string{integerColumn},
					Summaries: map[string]map[string]any{
						"count_integerCol": {
							"count": integerColumn,
						},
						"min_integerCol": {
							"min": integerColumn,
						},
						"max_integerCol": {
							"max": integerColumn,
						},
						"sum_integerCol": {
							"sum": integerColumn,
						},
						"average_integerCol": {
							"average": integerColumn,
						},
					},
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, 10.0, sumTableRes.Summaries[0]["max_integerCol"])
		assert.Equal(t, 10.0, sumTableRes.Summaries[0]["min_integerCol"])
		assert.Equal(t, 10.0, sumTableRes.Summaries[0]["sum_integerCol"])
		assert.Equal(t, 10.0, sumTableRes.Summaries[0]["average_integerCol"])
		assert.Equal(t, 1.0, sumTableRes.Summaries[0]["count_integerCol"])
	})

	t.Run("aggregate table count with filter", func(t *testing.T) {
		assert.Eventually(t, func() bool {
			aggTableRes, err := searchFilterCli.Aggregate(ctx, xata.AggregateTableRequest{
				BranchRequestOptional: xata.BranchRequestOptional{
					DatabaseName: xata.String(cfg.databaseName),
				},
				TableName: cfg.tableName,
				Payload: xata.AggregateTableRequestPayload{
					Aggregations: xata.AggExpressionMap{
						"filteredCount": xata.NewCountAggExpression(xata.AggExpressionCount{
							Count: xata.CountByFilter(xata.CountAggFilter{
								Filter: xata.FilterExpression{
									All: xata.NewFilterListFromFilterExpression(&xata.FilterExpression{
										Exists: xata.String(stringColumn),
									}),
								},
							}),
						}),
					},
				},
			})
			assert.NoError(t, err)
			return int(*(*aggTableRes.Aggs)["filteredCount"].DoubleOptional) == 1
		}, time.Second*10, time.Second)
	})

	t.Run("aggregate table count all", func(t *testing.T) {
		assert.Eventually(t, func() bool {
			aggTableRes, err := searchFilterCli.Aggregate(ctx, xata.AggregateTableRequest{
				BranchRequestOptional: xata.BranchRequestOptional{
					DatabaseName: xata.String(cfg.databaseName),
				},
				TableName: cfg.tableName,
				Payload: xata.AggregateTableRequestPayload{
					Aggregations: xata.AggExpressionMap{
						"totalCount": xata.NewCountAggExpression(xata.AggExpressionCount{
							Count: xata.CountAll(),
						}),
					},
				},
			})
			assert.NoError(t, err)
			return int(*(*aggTableRes.Aggs)["totalCount"].DoubleOptional) == 1
		}, time.Second*10, time.Second)
	})

	t.Run("aggregate table sum", func(t *testing.T) {
		assert.Eventually(t, func() bool {
			aggTableRes, err := searchFilterCli.Aggregate(ctx, xata.AggregateTableRequest{
				BranchRequestOptional: xata.BranchRequestOptional{
					DatabaseName: xata.String(cfg.databaseName),
				},
				TableName: cfg.tableName,
				Payload: xata.AggregateTableRequestPayload{
					Aggregations: xata.AggExpressionMap{
						"sum": xata.NewSumAggExpression(integerColumn),
					},
				},
			})
			assert.NoError(t, err)
			return int(*(*aggTableRes.Aggs)["sum"].DoubleOptional) == 10
		}, time.Second*10, time.Second)
	})

	t.Run("aggregate table max", func(t *testing.T) {
		assert.Eventually(t, func() bool {
			aggTableRes, err := searchFilterCli.Aggregate(ctx, xata.AggregateTableRequest{
				BranchRequestOptional: xata.BranchRequestOptional{
					DatabaseName: xata.String(cfg.databaseName),
				},
				TableName: cfg.tableName,
				Payload: xata.AggregateTableRequestPayload{
					Aggregations: xata.AggExpressionMap{
						"max": xata.NewMaxAggExpression(integerColumn),
					},
				},
			})
			assert.NoError(t, err)
			if (*aggTableRes.Aggs)["max"] != nil {
				return int(*(*aggTableRes.Aggs)["max"].DoubleOptional) == 10
			}
			return false
		}, time.Second*10, time.Second)
	})

	t.Run("aggregate table min", func(t *testing.T) {
		assert.Eventually(t, func() bool {
			aggTableRes, err := searchFilterCli.Aggregate(ctx, xata.AggregateTableRequest{
				BranchRequestOptional: xata.BranchRequestOptional{
					DatabaseName: xata.String(cfg.databaseName),
				},
				TableName: cfg.tableName,
				Payload: xata.AggregateTableRequestPayload{
					Aggregations: xata.AggExpressionMap{
						"min": xata.NewMinAggExpression(integerColumn),
					},
				},
			})
			assert.NoError(t, err)
			if (*aggTableRes.Aggs)["min"] != nil {
				return int(*(*aggTableRes.Aggs)["min"].DoubleOptional) == 10
			}
			return false
		}, time.Second*10, time.Second)
	})

	t.Run("aggregate table average", func(t *testing.T) {
		assert.Eventually(t, func() bool {
			aggTableRes, err := searchFilterCli.Aggregate(ctx, xata.AggregateTableRequest{
				BranchRequestOptional: xata.BranchRequestOptional{
					DatabaseName: xata.String(cfg.databaseName),
				},
				TableName: cfg.tableName,
				Payload: xata.AggregateTableRequestPayload{
					Aggregations: xata.AggExpressionMap{
						"average": xata.NewAverageAggExpression(integerColumn),
					},
				},
			})
			assert.NoError(t, err)
			if (*aggTableRes.Aggs)["average"] != nil {
				return int(*(*aggTableRes.Aggs)["average"].DoubleOptional) == 10
			}
			return false
		}, time.Second*10, time.Second)
	})

	t.Run("aggregate table unique", func(t *testing.T) {
		assert.Eventually(t, func() bool {
			aggTableRes, err := searchFilterCli.Aggregate(ctx, xata.AggregateTableRequest{
				BranchRequestOptional: xata.BranchRequestOptional{
					DatabaseName: xata.String(cfg.databaseName),
				},
				TableName: cfg.tableName,
				Payload: xata.AggregateTableRequestPayload{
					Aggregations: xata.AggExpressionMap{
						"unique": xata.NewUniqueCountAggExpression(xata.UniqueCountAgg{
							Column:             stringColumn,
							PrecisionThreshold: xata.Int(1),
						}),
					},
				},
			})
			assert.NoError(t, err)
			if (*aggTableRes.Aggs)["unique"] != nil {
				return int(*(*aggTableRes.Aggs)["unique"].DoubleOptional) == 1
			}
			return false
		}, time.Second*10, time.Second)
	})

	t.Run("aggregate table date histogram", func(t *testing.T) {
		assert.Eventually(t, func() bool {
			aggTableRes, err := searchFilterCli.Aggregate(ctx, xata.AggregateTableRequest{
				BranchRequestOptional: xata.BranchRequestOptional{
					DatabaseName: xata.String(cfg.databaseName),
				},
				TableName: cfg.tableName,
				Payload: xata.AggregateTableRequestPayload{
					Aggregations: xata.AggExpressionMap{
						"histogram": xata.NewDateHistogramAggExpression(xata.DateHistogramAgg{
							Column:           dateTimeColumn,
							Interval:         xata.String("1d"),
							CalendarInterval: nil,
							Timezone:         nil,
						}),
					},
				},
			})
			assert.NoError(t, err)
			if (*aggTableRes.Aggs)["histogram"] != nil &&
				(*aggTableRes.Aggs)["histogram"].AggResponseValues != nil {
				return len((*aggTableRes.Aggs)["histogram"].AggResponseValues.Values) > 0
			}
			return false
		}, time.Second*10, time.Second)
	})

	t.Run("aggregate table date histogram", func(t *testing.T) {
		assert.Eventually(t, func() bool {
			aggTableRes, err := searchFilterCli.Aggregate(ctx, xata.AggregateTableRequest{
				BranchRequestOptional: xata.BranchRequestOptional{
					DatabaseName: xata.String(cfg.databaseName),
				},
				TableName: cfg.tableName,
				Payload: xata.AggregateTableRequestPayload{
					Aggregations: xata.AggExpressionMap{
						"histogram": xata.NewDateHistogramAggExpression(xata.DateHistogramAgg{
							Column:           dateTimeColumn,
							Interval:         xata.String("1d"),
							CalendarInterval: nil,
							Timezone:         nil,
						}),
					},
				},
			})
			assert.NoError(t, err)
			if (*aggTableRes.Aggs)["histogram"] != nil && (*aggTableRes.Aggs)["histogram"].AggResponseValues != nil {
				return len((*aggTableRes.Aggs)["histogram"].AggResponseValues.Values) > 0
			}
			return false
		}, time.Second*10, time.Second)
	})

	t.Run("aggregate table top values", func(t *testing.T) {
		assert.Eventually(t, func() bool {
			aggTableRes, err := searchFilterCli.Aggregate(ctx, xata.AggregateTableRequest{
				BranchRequestOptional: xata.BranchRequestOptional{
					DatabaseName: xata.String(cfg.databaseName),
				},
				TableName: cfg.tableName,
				Payload: xata.AggregateTableRequestPayload{
					Aggregations: xata.AggExpressionMap{
						"top_values": xata.NewTopValuesAggExpression(xata.TopValuesAgg{
							Column: stringColumn,
							Size:   nil,
						}),
					},
				},
			})
			assert.NoError(t, err)
			if (*aggTableRes.Aggs)["top_values"] != nil && (*aggTableRes.Aggs)["top_values"].AggResponseValues != nil {
				return len((*aggTableRes.Aggs)["top_values"].AggResponseValues.Values) > 0
			}
			return false
		}, time.Second*10, time.Second)
	})

	t.Run("aggregate table numeric histogram", func(t *testing.T) {
		assert.Eventually(t, func() bool {
			aggTableRes, err := searchFilterCli.Aggregate(ctx, xata.AggregateTableRequest{
				BranchRequestOptional: xata.BranchRequestOptional{
					DatabaseName: xata.String(cfg.databaseName),
				},
				TableName: cfg.tableName,
				Payload: xata.AggregateTableRequestPayload{
					Aggregations: xata.AggExpressionMap{
						"num_histogram": xata.NewNumericHistogramAggExpression(xata.NumericHistogramAgg{
							Column:   integerColumn,
							Interval: 1.0,
							Offset:   nil,
						}),
					},
				},
			})
			assert.NoError(t, err)
			if (*aggTableRes.Aggs)["num_histogram"] != nil && (*aggTableRes.Aggs)["num_histogram"].AggResponseValues != nil {
				return len((*aggTableRes.Aggs)["num_histogram"].AggResponseValues.Values) > 0
			}
			return false
		}, time.Second*10, time.Second)
	})
}
