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

func Test_tableClient_Delete(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		request    xata.TableRequest
		want       *xatagenworkspace.DeleteTableResponse
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should delete a table",
			request: xata.TableRequest{
				DatabaseName: xata.String("db-name"),
				TableName:    "table-name",
			},
			want:       &xatagenworkspace.DeleteTableResponse{Status: xatagenworkspace.MigrationStatusCompleted},
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
			testSrv := testService(t, http.MethodDelete, "/db", tt.statusCode, tt.apiErr != nil, tt.want)

			cli, err := xata.NewTableClient(xata.WithBaseURL(testSrv.URL), xata.WithAPIKey("test-key"))
			assert.NoError(err)
			assert.NotNil(cli)

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
				assert.NoError(err)
				assert.Equal(tt.want, got)
			}
		})
	}
}

func Test_tableClient_AddColumn(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		request    xata.AddColumnRequest
		want       *xatagenworkspace.AddTableColumnResponse
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should create a table",
			request: xata.AddColumnRequest{
				TableRequest: xata.TableRequest{
					DatabaseName: xata.String("db-name"),
					TableName:    "table-name",
				},
				Column: &xata.Column{
					Name:         "user-name",
					Type:         xata.ColumnTypeString,
					NotNull:      xata.Bool(true),
					DefaultValue: xata.String("defaultValue"),
					Unique:       xata.Bool(false),
				},
			},
			want: &xatagenworkspace.AddTableColumnResponse{
				MigrationId:       "mig-id",
				ParentMigrationId: "parent-mig-id",
				Status:            xatagenworkspace.MigrationStatusCompleted,
			},
			statusCode: http.StatusOK,
		},
	}

	for _, eTC := range errTestCasesWorkspace {
		tests = append(tests, tc{
			name: eTC.name,
			request: xata.AddColumnRequest{
				TableRequest: xata.TableRequest{
					DatabaseName: xata.String("db-name"),
					TableName:    "table-name",
				},
				Column: &xata.Column{
					Name:         "user-name",
					Type:         xata.ColumnTypeString,
					NotNull:      xata.Bool(true),
					DefaultValue: xata.String("defaultValue"),
					Unique:       xata.Bool(false),
				},
			},
			statusCode: eTC.statusCode,
			apiErr:     eTC.apiErr,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testSrv := testService(t, http.MethodPost, "/db", tt.statusCode, tt.apiErr != nil, tt.want)

			cli, err := xata.NewTableClient(xata.WithBaseURL(testSrv.URL), xata.WithAPIKey("test-key"))
			assert.NoError(err)
			assert.NotNil(cli)

			got, err := cli.AddColumn(context.TODO(), tt.request)

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

func Test_tableClient_DeleteColumn(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		request    xata.DeleteColumnRequest
		want       *xatagenworkspace.DeleteColumnResponse
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should create a table",
			request: xata.DeleteColumnRequest{
				TableRequest: xata.TableRequest{
					DatabaseName: xata.String("db-name"),
					TableName:    "table-name",
				},
				ColumnName: "col-name",
			},
			want: &xatagenworkspace.DeleteColumnResponse{
				MigrationId:       "mig-id",
				ParentMigrationId: "par-mig-id",
				Status:            xatagenworkspace.MigrationStatusCompleted,
			},
			statusCode: http.StatusOK,
		},
	}

	for _, eTC := range errTestCasesWorkspace {
		tests = append(tests, tc{
			name: eTC.name,
			request: xata.DeleteColumnRequest{
				TableRequest: xata.TableRequest{
					DatabaseName: xata.String("db-name"),
					TableName:    "table-name",
				},
				ColumnName: "col-name",
			},
			statusCode: eTC.statusCode,
			apiErr:     eTC.apiErr,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testSrv := testService(t, http.MethodDelete, "/db", tt.statusCode, tt.apiErr != nil, tt.want)

			cli, err := xata.NewTableClient(xata.WithBaseURL(testSrv.URL), xata.WithAPIKey("test-key"))
			assert.NoError(err)
			assert.NotNil(cli)

			got, err := cli.DeleteColumn(context.TODO(), tt.request)

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
