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
		request    xata.QueryTableRequest
		want       *xatagenworkspace.QueryTableResponse
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should get a record successfully",
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
				TableRequest: xata.TableRequest{
					DatabaseName: xata.String("some-db"),
					TableName:    "table-name",
				},
				Payload: xata.QueryTableRequestPayload{
					Columns:     []string{"column-name"},
					Consistency: xata.QueryTableRequestConsistencyEventual,
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
		request    xata.SearchBranchRequest
		want       *xatagenworkspace.SearchBranchResponse
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should get a record successfully",
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
				TableRequest: xata.TableRequest{
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
