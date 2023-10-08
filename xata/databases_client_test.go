package xata_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xataio/xata-go/xata"
	xatagen "github.com/xataio/xata-go/xata/internal/fern-core/generated/go"
	xatagencore "github.com/xataio/xata-go/xata/internal/fern-core/generated/go/core"
)

func TestNewDatabasesClient(t *testing.T) {
	t.Run("should construct a new client", func(t *testing.T) {
		got, err := xata.NewRecordsClient(
			xata.WithBaseURL("https://www.example.com"),
			xata.WithAPIKey("my-api-token"),
		)
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})
}

func Test_databaseCli_Create(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		request    xata.CreateDatabaseRequest
		want       *xatagen.CreateDatabaseResponse
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should create db",
			request: xata.CreateDatabaseRequest{
				DatabaseName: "db-name",
			},
			want: &xatagen.CreateDatabaseResponse{
				DatabaseName: "db-name",
				BranchName:   xata.String("main"),
				Status:       xatagen.MigrationStatusCompleted,
			},
			statusCode: http.StatusCreated,
		},
	}

	for _, eTC := range errTestCasesCore {
		tests = append(tests, tc{
			name:       eTC.name,
			statusCode: eTC.statusCode,
			apiErr:     eTC.apiErr,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testSrv := testService(t, http.MethodPut, "/workspaces", tt.statusCode, tt.apiErr != nil, tt.want)

			cli, err := xata.NewDatabasesClient(xata.WithBaseURL(testSrv.URL), xata.WithAPIKey("test-key"))
			if err != nil {
				t.Fatal(err)
			}

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
				assert.Equal(tt.want, got)
				assert.NoError(err)
			}
		})
	}
}

func Test_databaseCli_Delete(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		request    xata.DeleteDatabaseRequest
		want       *xatagen.DeleteDatabaseResponse
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should delete db",
			request: xata.DeleteDatabaseRequest{
				DatabaseName: "db-name",
			},
			want: &xatagen.DeleteDatabaseResponse{
				Status: xatagen.MigrationStatusCompleted,
			},
			statusCode: http.StatusOK,
		},
	}

	for _, eTC := range errTestCasesCore {
		tests = append(tests, tc{
			name:       eTC.name,
			statusCode: eTC.statusCode,
			apiErr:     eTC.apiErr,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testSrv := testService(t, http.MethodDelete, "/workspaces", tt.statusCode, tt.apiErr != nil, tt.want)

			cli, err := xata.NewDatabasesClient(xata.WithBaseURL(testSrv.URL), xata.WithAPIKey("test-key"))
			if err != nil {
				t.Fatal(err)
			}

			got, err := cli.Delete(context.TODO(), tt.request)

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
