package xata_test

import (
	"context"
	"net/http"
	"testing"

	xatagenworkspace "github.com/xataio/xata-go/xata/internal/fern-workspace/generated/go"
	xatagencore "github.com/xataio/xata-go/xata/internal/fern-workspace/generated/go/core"

	"github.com/stretchr/testify/assert"
	"github.com/xataio/xata-go/xata"
)

func TestNewTableClient(t *testing.T) {
	t.Run("should construct a new client", func(t *testing.T) {
		got, err := xata.NewTableClient(
			xata.WithBaseURL("https://www.example.com"),
			xata.WithAPIKey("my-api-token"),
		)
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})
}

func Test_tableClient_Create(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		request    xata.TableRequest
		want       *xatagenworkspace.CreateTableResponse
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should create a table",
			request: xata.TableRequest{
				DatabaseName: xata.String("db-name"),
				TableName:    "table-name",
			},
			want: &xatagenworkspace.CreateTableResponse{
				BranchName: "main",
				TableName:  "table-name",
				Status:     xatagenworkspace.MigrationStatusCompleted,
			},
			statusCode: http.StatusCreated,
		},
	}

	for _, eTC := range errTestCasesWorkspace {
		tests = append(tests, tc{
			name: eTC.name,
			request: xata.TableRequest{
				DatabaseName: xata.String("db-name"),
				TableName:    "table-name",
			},
			statusCode: eTC.statusCode,
			apiErr:     eTC.apiErr,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testSrv := testService(t, http.MethodPut, "/db", tt.statusCode, tt.apiErr != nil, tt.want)

			cli, err := xata.NewTableClient(xata.WithBaseURL(testSrv.URL), xata.WithAPIKey("test-key"))
			assert.NoError(err)
			assert.NotNil(cli)

			got, err := cli.Create(context.TODO(), tt.request)

			if tt.apiErr != nil {
				errAPI := tt.apiErr.Unwrap()
				if errAPI == nil {
					t.Fatal("expected error but got nil")
				}
				assert.ErrorAs(err, &errAPI)
				assert.Equal(err.Error(), tt.apiErr.Error())
				assert.Nil(got)
			} else {
				assert.NoError(err)
				assert.Equal(tt.want, got)
			}
		})
	}
}
