// SPDX-License-Identifier: Apache-2.0

package xata_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xataio/xata-go/xata"

	xatagenworkspace "github.com/xataio/xata-go/xata/internal/fern-workspace/generated/go"
	xatagencore "github.com/xataio/xata-go/xata/internal/fern-workspace/generated/go/core"
)

func TestNewSearchAndFilterClient(t *testing.T) {
	t.Run("should construct a new client", func(t *testing.T) {
		got, err := xata.NewSearchAndFilterClient(
			xata.WithBaseURL("https://www.example.com"),
			xata.WithAPIKey("my-api-token"),
		)
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})
}

func Test_searchAndFilterCli_Query(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		want       *xatagenworkspace.QueryTableResponse
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should query",
			want: &xatagenworkspace.QueryTableResponse{
				Records: []*xatagenworkspace.Record{
					{
						"test-key": "test-value",
					},
				},
				Meta: &xatagenworkspace.RecordsMetadata{
					Page: &xatagenworkspace.RecordsMetadataPage{
						Cursor: "some-cursor",
						More:   false,
						Size:   10,
					},
				},
			},
			statusCode: http.StatusOK,
		},
	}

	for _, eTC := range errTestCasesWorkspace {
		tests = append(tests, tc{
			name:       eTC.name,
			statusCode: eTC.statusCode,
			apiErr:     eTC.apiErr,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testSrv := testService(t, http.MethodPost, "/db", tt.statusCode, tt.apiErr != nil, tt.want)

			cli, err := xata.NewSearchAndFilterClient(
				xata.WithBaseURL(testSrv.URL),
				xata.WithAPIKey("test-key"),
			)
			assert.NoError(err)
			assert.NotNil(cli)

			got, err := cli.Query(context.TODO(), xata.QueryTableRequest{
				BranchRequestOptional: xata.BranchRequestOptional{
					DatabaseName: xata.String("some-db"),
				},
				TableName: "some-table",
				Payload: xata.QueryTableRequestPayload{
					Columns:     []string{"column-name"},
					Consistency: xata.ConsistencyEventual,
					Sort: xata.NewSortExpressionFromStringSortOrderMapList(
						[]map[string]xata.SortOrder{
							{
								"column-name": xata.SortOrderDesc,
							},
						}),
				},
			})

			if tt.apiErr != nil {
				errAPI := tt.apiErr.Unwrap()
				if errAPI == nil {
					t.Fatal("expected error but got nil")
				}
				assert.ErrorAs(err, &errAPI)
				assert.Equal(err.Error(), tt.apiErr.Error())
				assert.Nil(got)
			} else {
				assert.Equal(tt.want.Records[0], got.Records[0])
				assert.NoError(err)
			}
		})
	}
}

func Test_searchAndFilterCli_SearchBranch(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		want       *xatagenworkspace.SearchBranchResponse
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should search branch",
			want: &xatagenworkspace.SearchBranchResponse{
				Records: []*xatagenworkspace.Record{
					{
						"test-key": "test-value",
					},
				},
				Warning:    nil,
				TotalCount: 1,
			},
			statusCode: http.StatusOK,
		},
	}

	for _, eTC := range errTestCasesWorkspace {
		tests = append(tests, tc{
			name:       eTC.name,
			statusCode: eTC.statusCode,
			apiErr:     eTC.apiErr,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testSrv := testService(t, http.MethodPost, "/db", tt.statusCode, tt.apiErr != nil, tt.want)

			cli, err := xata.NewSearchAndFilterClient(
				xata.WithBaseURL(testSrv.URL),
				xata.WithAPIKey("test-key"),
			)
			assert.NoError(err)
			assert.NotNil(cli)

			got, err := cli.SearchBranch(context.TODO(), xata.SearchBranchRequest{
				BranchRequestOptional: xata.BranchRequestOptional{
					DatabaseName: xata.String("db-name"),
				},
				Payload: xata.SearchBranchRequestPayload{
					Tables: []*xata.SearchBranchRequestTablesItem{
						xata.NewSearchBranchRequestTablesItemFromString("table-name"),
					},
					Query:     "test",
					Fuzziness: xata.Int(0),
					Prefix:    nil,
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

			if tt.apiErr != nil {
				errAPI := tt.apiErr.Unwrap()
				if errAPI == nil {
					t.Fatal("expected error but got nil")
				}
				assert.ErrorAs(err, &errAPI)
				assert.Equal(err.Error(), tt.apiErr.Error())
				assert.Nil(got)
			} else {
				assert.Equal(tt.want.Records[0], got.Records[0])
				assert.NoError(err)
			}
		})
	}
}

func Test_searchAndFilterCli_SearchTable(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		want       *xatagenworkspace.SearchTableResponse
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should search table",
			want: &xatagenworkspace.SearchTableResponse{
				Records: []*xatagenworkspace.Record{
					{
						"key": "value",
					},
				},
				Warning:    xata.String("warning"),
				TotalCount: 1,
			},
			statusCode: http.StatusOK,
		},
	}

	for _, eTC := range errTestCasesWorkspace {
		tests = append(tests, tc{
			name:       eTC.name,
			statusCode: eTC.statusCode,
			apiErr:     eTC.apiErr,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testSrv := testService(t, http.MethodPost, "/db", tt.statusCode, tt.apiErr != nil, tt.want)

			cli, err := xata.NewSearchAndFilterClient(
				xata.WithBaseURL(testSrv.URL),
				xata.WithAPIKey("test-key"),
			)
			assert.NoError(err)
			assert.NotNil(cli)

			prefix := xata.PrefixExpressionDisabled
			got, err := cli.SearchTable(context.TODO(), xata.SearchTableRequest{
				BranchRequestOptional: xata.BranchRequestOptional{
					DatabaseName: xata.String("db"),
					BranchName:   xata.String("branch"),
				},
				TableName: "table",
				Payload: xata.SearchTableRequestPayload{
					Query:     "test",
					Fuzziness: xata.Int(0),
					Target: []*xata.TargetExpressionItem{
						xata.NewTargetExpression("column"),
						xata.NewTargetExpressionWithColumnObject(xata.TargetExpressionItemColumn{
							Column: "columnt",
							Weight: xata.Float64(2),
						}),
					},
					Prefix: &prefix,
					Filter: &xata.FilterExpression{
						All: xata.NewFilterListFromFilterExpression(&xata.FilterExpression{
							Exists: xata.String("column"),
						}),
						Any: xata.NewFilterListFromFilterExpressionList([]*xata.FilterExpression{
							{
								Exists: xata.String("column"),
							},
						}),
					},
					Highlight: &xata.HighlightExpression{
						Enabled:    xata.Bool(true),
						EncodeHtml: xata.Bool(true),
					},
					Boosters: []*xata.BoosterExpression{
						xata.NewBoosterExpressionFromBoosterExpressionValueBooster(&xata.BoosterExpressionValueBooster{
							ValueBooster: &xata.ValueBooster{
								Column: "column",
								Value:  xata.NewValueBoosterValueFromString("test"),
								Factor: 1,
								IfMatchesFilter: &xata.FilterExpression{
									All: xata.NewFilterListFromFilterExpression(&xata.FilterExpression{
										Exists: xata.String("column"),
									}),
									Any: xata.NewFilterListFromFilterExpressionList([]*xata.FilterExpression{
										{
											Exists: xata.String("column"),
										},
									}),
								},
							},
						}),
						xata.NewBoosterExpressionFromBoosterExpressionNumericBooster(&xata.BoosterExpressionNumericBooster{
							NumericBooster: &xata.NumericBooster{
								Column:   "column",
								Factor:   2,
								Modifier: xata.Uint8(2),
								IfMatchesFilter: &xata.FilterExpression{
									All: xata.NewFilterListFromFilterExpression(&xata.FilterExpression{
										Exists: xata.String("column"),
									}),
									Any: xata.NewFilterListFromFilterExpressionList([]*xata.FilterExpression{
										{
											Exists: xata.String("bool-olumn"),
										},
									}),
								},
							},
						}),
						xata.NewBoosterExpressionFromBoosterExpressionDateBooster(&xata.BoosterExpressionDateBooster{
							DateBooster: &xata.DateBooster{
								Column: "column",
								Origin: xata.String("2023-01-02T15:04:05Z"),
								Scale:  "1d",
								Decay:  1,
								Factor: xata.Float64(2),
								IfMatchesFilter: &xata.FilterExpression{
									All: xata.NewFilterListFromFilterExpression(&xata.FilterExpression{
										Exists: xata.String("column"),
									}),
									Any: xata.NewFilterListFromFilterExpressionList([]*xata.FilterExpression{
										{
											Exists: xata.String("column"),
										},
									}),
								},
							},
						}),
					},
					Page: &xata.SearchPageConfig{
						Size:   xata.Int(2),
						Offset: xata.Int(0),
					},
				},
			})

			if tt.apiErr != nil {
				errAPI := tt.apiErr.Unwrap()
				if errAPI == nil {
					t.Fatal("expected error but got nil")
				}
				assert.ErrorAs(err, &errAPI)
				assert.Equal(err.Error(), tt.apiErr.Error())
				assert.Nil(got)
			} else {
				assert.Equal(tt.want.Records[0], got.Records[0])
				assert.NoError(err)
			}
		})
	}
}

func Test_searchAndFilterCli_VectorSearch(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		want       *xatagenworkspace.VectorSearchTableResponse
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should vector search",
			want: &xatagenworkspace.VectorSearchTableResponse{
				Records: []*xatagenworkspace.Record{
					{
						"key": "value",
					},
				},
				Warning:    xata.String("warning"),
				TotalCount: 1,
			},
			statusCode: http.StatusOK,
		},
	}

	for _, eTC := range errTestCasesWorkspace {
		tests = append(tests, tc{
			name:       eTC.name,
			statusCode: eTC.statusCode,
			apiErr:     eTC.apiErr,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testSrv := testService(t, http.MethodPost, "/db", tt.statusCode, tt.apiErr != nil, tt.want)

			cli, err := xata.NewSearchAndFilterClient(
				xata.WithBaseURL(testSrv.URL),
				xata.WithAPIKey("test-key"),
			)
			assert.NoError(err)
			assert.NotNil(cli)

			got, err := cli.VectorSearch(context.TODO(), xata.VectorSearchTableRequest{
				BranchRequestOptional: xata.BranchRequestOptional{
					DatabaseName: xata.String("db"),
					BranchName:   xata.String("branch"),
				},
				TableName: "table",
				Payload: xata.VectorSearchTableRequestPayload{
					QueryVector:        []float64{10, 2},
					Column:             "vectorColumn",
					SimilarityFunction: xata.String("cosineSimilarity"),
					Size:               xata.Int(2),
					Filter: &xata.FilterExpression{
						All: xata.NewFilterListFromFilterExpression(&xata.FilterExpression{
							Exists: xata.String("vectorColumn"),
						}),
					},
				},
			})

			if tt.apiErr != nil {
				errAPI := tt.apiErr.Unwrap()
				if errAPI == nil {
					t.Fatal("expected error but got nil")
				}
				assert.ErrorAs(err, &errAPI)
				assert.Equal(err.Error(), tt.apiErr.Error())
				assert.Nil(got)
			} else {
				assert.Equal(tt.want.Records[0], got.Records[0])
				assert.NoError(err)
			}
		})
	}
}

func Test_searchAndFilterCli_Ask(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		want       *xatagenworkspace.AskTableResponse
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should ask a question",
			want: &xatagenworkspace.AskTableResponse{
				Answer:    "id1",
				SessionId: "id2",
			},
			statusCode: http.StatusOK,
		},
	}

	for _, eTC := range errTestCasesWorkspace {
		tests = append(tests, tc{
			name:       eTC.name,
			statusCode: eTC.statusCode,
			apiErr:     eTC.apiErr,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testSrv := testService(t, http.MethodPost, "/db", tt.statusCode, tt.apiErr != nil, tt.want)

			cli, err := xata.NewSearchAndFilterClient(
				xata.WithBaseURL(testSrv.URL),
				xata.WithAPIKey("test-key"),
			)
			assert.NoError(err)
			assert.NotNil(cli)

			keyword := xata.AskTableRequestSearchTypeKeyword
			prefix := xata.PrefixExpressionPhrase
			got, err := cli.Ask(context.TODO(), xata.AskTableRequest{
				BranchRequestOptional: xata.BranchRequestOptional{
					DatabaseName: xata.String("db"),
					BranchName:   xata.String("branch"),
				},
				TableName: "table",
				Payload: xata.AskTableRequestPayload{
					Question:   "how are you",
					SearchType: &keyword,
					Search: &xata.AskTableRequestSearch{
						Target: []*xata.TargetExpressionItem{
							xata.NewTargetExpression("columnName"),
						},
						Prefix: &prefix,
						Filter: &xata.FilterExpression{
							All: xata.NewFilterListFromFilterExpression(&xata.FilterExpression{
								Exists: xata.String("column"),
							}),
							Any: xata.NewFilterListFromFilterExpressionList([]*xata.FilterExpression{
								{
									Exists: xata.String("column"),
								},
							}),
						},
						Boosters: []*xata.BoosterExpression{
							xata.NewBoosterExpressionFromBoosterExpressionValueBooster(&xata.BoosterExpressionValueBooster{
								ValueBooster: &xata.ValueBooster{
									Column: "column",
									Value:  xata.NewValueBoosterValueFromString("test"),
									Factor: 1,
									IfMatchesFilter: &xata.FilterExpression{
										All: xata.NewFilterListFromFilterExpression(&xata.FilterExpression{
											Exists: xata.String("column"),
										}),
										Any: xata.NewFilterListFromFilterExpressionList([]*xata.FilterExpression{
											{
												Exists: xata.String("column"),
											},
										}),
									},
								},
							}),
							xata.NewBoosterExpressionFromBoosterExpressionNumericBooster(&xata.BoosterExpressionNumericBooster{
								NumericBooster: &xata.NumericBooster{
									Column:   "column",
									Factor:   2,
									Modifier: xata.Uint8(2),
									IfMatchesFilter: &xata.FilterExpression{
										All: xata.NewFilterListFromFilterExpression(&xata.FilterExpression{
											Exists: xata.String("column"),
										}),
										Any: xata.NewFilterListFromFilterExpressionList([]*xata.FilterExpression{
											{
												Exists: xata.String("bool-olumn"),
											},
										}),
									},
								},
							}),
							xata.NewBoosterExpressionFromBoosterExpressionDateBooster(&xata.BoosterExpressionDateBooster{
								DateBooster: &xata.DateBooster{
									Column: "column",
									Origin: xata.String("2023-01-02T15:04:05Z"),
									Scale:  "1d",
									Decay:  1,
									Factor: xata.Float64(2),
									IfMatchesFilter: &xata.FilterExpression{
										All: xata.NewFilterListFromFilterExpression(&xata.FilterExpression{
											Exists: xata.String("column"),
										}),
										Any: xata.NewFilterListFromFilterExpressionList([]*xata.FilterExpression{
											{
												Exists: xata.String("column"),
											},
										}),
									},
								},
							}),
						},
					},
				},
			})

			if tt.apiErr != nil {
				errAPI := tt.apiErr.Unwrap()
				if errAPI == nil {
					t.Fatal("expected error but got nil")
				}
				assert.ErrorAs(err, &errAPI)
				assert.Equal(err.Error(), tt.apiErr.Error())
				assert.Nil(got)
			} else {
				assert.Equal(tt.want.Answer, got.Answer)
				assert.NoError(err)
			}
		})
	}
}

func Test_searchAndFilterCli_AskFollowUp(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		want       *xatagenworkspace.AskTableSessionResponse
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should ask a follow up question",
			want: &xatagenworkspace.AskTableSessionResponse{
				Answer: "Sure, here are more info...",
			},
			statusCode: http.StatusOK,
		},
	}

	for _, eTC := range errTestCasesWorkspace {
		tests = append(tests, tc{
			name:       eTC.name,
			statusCode: eTC.statusCode,
			apiErr:     eTC.apiErr,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testSrv := testService(t, http.MethodPost, "/db", tt.statusCode, tt.apiErr != nil, tt.want)

			cli, err := xata.NewSearchAndFilterClient(
				xata.WithBaseURL(testSrv.URL),
				xata.WithAPIKey("test-key"),
			)
			assert.NoError(err)
			assert.NotNil(cli)

			got, err := cli.AskFollowUp(context.TODO(), xata.AskFollowUpRequest{
				BranchRequestOptional: xata.BranchRequestOptional{
					DatabaseName: xata.String("db"),
					BranchName:   xata.String("branch"),
				},
				TableName: "table",
				Question:  "give me more info on...",
			})

			if tt.apiErr != nil {
				errAPI := tt.apiErr.Unwrap()
				if errAPI == nil {
					t.Fatal("expected error but got nil")
				}
				assert.ErrorAs(err, &errAPI)
				assert.Equal(err.Error(), tt.apiErr.Error())
				assert.Nil(got)
			} else {
				assert.Equal(tt.want.Answer, got.Answer)
				assert.NoError(err)
			}
		})
	}
}

func Test_searchAndFilterCli_Summarize(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		want       *xatagenworkspace.SummarizeTableResponse
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should summarize a table",
			want: &xatagenworkspace.SummarizeTableResponse{Summaries: []map[string]any{
				{
					"k": "v",
				},
			}},
			statusCode: http.StatusOK,
		},
	}

	for _, eTC := range errTestCasesWorkspace {
		tests = append(tests, tc{
			name:       eTC.name,
			statusCode: eTC.statusCode,
			apiErr:     eTC.apiErr,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testSrv := testService(t, http.MethodPost, "/db", tt.statusCode, tt.apiErr != nil, tt.want)

			cli, err := xata.NewSearchAndFilterClient(
				xata.WithBaseURL(testSrv.URL),
				xata.WithAPIKey("test-key"),
			)
			assert.NoError(err)
			assert.NotNil(cli)

			got, err := cli.Summarize(
				context.TODO(),
				xata.SummarizeTableRequest{
					BranchRequestOptional: xata.BranchRequestOptional{
						DatabaseName: xata.String("db"),
						BranchName:   xata.String("branch"),
					},
					TableName: "table",
					Payload: xata.SummarizeTableRequestPayload{
						Columns:   []string{"col"},
						Summaries: map[string]map[string]any{"count_col": {"count": "col"}},
					},
				},
			)

			if tt.apiErr != nil {
				errAPI := tt.apiErr.Unwrap()
				if errAPI == nil {
					t.Fatal("expected error but got nil")
				}
				assert.ErrorAs(err, &errAPI)
				assert.Equal(err.Error(), tt.apiErr.Error())
				assert.Nil(got)
			} else {
				assert.Equal(tt.want, got)
				assert.NoError(err)
			}
		})
	}
}

// Agg: Date Histogram
func Test_searchAndFilterCli_Aggregate_DateHistogram(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		want       *xatagenworkspace.AggregateTableResponse
		statusCode int
		apiErr     *xatagencore.APIError
	}

	aggRes := map[string]*xatagenworkspace.AggResponse{
		"test": xatagenworkspace.NewAggResponseFromDoubleOptional(xata.Float64(2)),
	}

	tests := []tc{
		{
			name:       "should aggregate a table",
			want:       &xatagenworkspace.AggregateTableResponse{Aggs: &aggRes},
			statusCode: http.StatusOK,
		},
	}

	for _, eTC := range errTestCasesWorkspace {
		tests = append(tests, tc{
			name:       eTC.name,
			statusCode: eTC.statusCode,
			apiErr:     eTC.apiErr,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testSrv := testService(t, http.MethodPost, "/db", tt.statusCode, tt.apiErr != nil, tt.want)

			cli, err := xata.NewSearchAndFilterClient(
				xata.WithBaseURL(testSrv.URL),
				xata.WithAPIKey("test-key"),
			)
			assert.NoError(err)
			assert.NotNil(cli)

			got, err := cli.Aggregate(
				context.TODO(),
				xata.AggregateTableRequest{
					BranchRequestOptional: xata.BranchRequestOptional{
						DatabaseName: xata.String("my-db"),
					},
					TableName: "my-table",
					Payload: xata.AggregateTableRequestPayload{
						Aggregations: xata.AggExpressionMap{
							"histogram": xata.NewDateHistogramAggExpression(xata.DateHistogramAgg{
								Column:           "my-column",
								Interval:         xata.String("1d"),
								CalendarInterval: nil,
								Timezone:         nil,
							}),
						},
					},
				},
			)

			if tt.apiErr != nil {
				errAPI := tt.apiErr.Unwrap()
				if errAPI == nil {
					t.Fatal("expected error but got nil")
				}
				assert.ErrorAs(err, &errAPI)
				assert.Equal(err.Error(), tt.apiErr.Error())
				assert.Nil(got)
			} else {
				assert.Equal(tt.want, got)
				assert.NoError(err)
			}
		})
	}
}
