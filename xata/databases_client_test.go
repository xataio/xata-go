// SPDX-License-Identifier: Apache-2.0

package xata_test

import (
	"context"
	"net/http"
	"testing"
	"time"

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

func Test_databaseCli_GetRegions(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		want       *xatagen.ListRegionsResponse
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should get workspaces",
			want: &xatagen.ListRegionsResponse{
				Regions: []*xatagen.Region{
					{
						Id:   "region-id",
						Name: "region-name",
					},
				},
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
			testSrv := testService(t, http.MethodGet, "/workspaces", tt.statusCode, tt.apiErr != nil, tt.want)

			cli, err := xata.NewDatabasesClient(xata.WithBaseURL(testSrv.URL), xata.WithAPIKey("test-key"))
			if err != nil {
				t.Fatal(err)
			}

			got, err := cli.GetRegions(context.TODO())

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

func Test_databaseCli_GetRegionsWithWorkspaceID(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		want       *xatagen.ListRegionsResponse
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should get workspaces",
			want: &xatagen.ListRegionsResponse{
				Regions: []*xatagen.Region{
					{
						Id:   "region-id",
						Name: "region-name",
					},
				},
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
			testSrv := testService(t, http.MethodGet, "/workspaces", tt.statusCode, tt.apiErr != nil, tt.want)

			cli, err := xata.NewDatabasesClient(xata.WithBaseURL(testSrv.URL), xata.WithAPIKey("test-key"))
			if err != nil {
				t.Fatal(err)
			}

			got, err := cli.GetRegionsWithWorkspaceID(context.TODO(), "workspace-id")

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

func Test_databaseCli_List(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		want       *xatagen.ListDatabasesResponse
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should list dbs",
			want: &xatagen.ListDatabasesResponse{Databases: []*xatagen.DatabaseMetadata{
				{
					Name:      "db-name",
					Region:    "region",
					CreatedAt: time.Now(),
				},
			}},
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
			testSrv := testService(t, http.MethodGet, "/workspaces", tt.statusCode, tt.apiErr != nil, tt.want)

			cli, err := xata.NewDatabasesClient(xata.WithBaseURL(testSrv.URL), xata.WithAPIKey("test-key"))
			if err != nil {
				t.Fatal(err)
			}

			got, err := cli.List(context.TODO())

			if tt.apiErr != nil {
				errAPI := tt.apiErr.Unwrap()
				if errAPI == nil {
					t.Fatal("expected error but got nil")
				}
				assert.ErrorAs(err, &errAPI)
				assert.Equal(err.Error(), tt.apiErr.Error())
				assert.Nil(got)
			} else {
				assert.Equal(tt.want.Databases[0].Name, got.Databases[0].Name)
				assert.NoError(err)
			}
		})
	}
}

func Test_databaseCli_ListWithWorkspaceID(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		want       *xatagen.ListDatabasesResponse
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should list dbs with workspace ID",
			want: &xatagen.ListDatabasesResponse{Databases: []*xatagen.DatabaseMetadata{
				{
					Name:      "db-name",
					Region:    "region",
					CreatedAt: time.Now(),
				},
			}},
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
			testSrv := testService(t, http.MethodGet, "/workspaces", tt.statusCode, tt.apiErr != nil, tt.want)

			cli, err := xata.NewDatabasesClient(xata.WithBaseURL(testSrv.URL), xata.WithAPIKey("test-key"))
			if err != nil {
				t.Fatal(err)
			}

			got, err := cli.ListWithWorkspaceID(context.TODO(), "ws-id")

			if tt.apiErr != nil {
				errAPI := tt.apiErr.Unwrap()
				if errAPI == nil {
					t.Fatal("expected error but got nil")
				}
				assert.ErrorAs(err, &errAPI)
				assert.Equal(err.Error(), tt.apiErr.Error())
				assert.Nil(got)
			} else {
				assert.Equal(tt.want.Databases[0].Name, got.Databases[0].Name)
				assert.NoError(err)
			}
		})
	}
}

func Test_databaseCli_Rename(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		request    xata.RenameDatabaseRequest
		want       *xatagen.DatabaseMetadata
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should update db name",
			request: xata.RenameDatabaseRequest{
				DatabaseName: "old-name",
				NewName:      "new-name",
			},
			want: &xatagen.DatabaseMetadata{
				Name:      "new-name",
				Region:    "region",
				CreatedAt: time.Now(),
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
			testSrv := testService(t, http.MethodPost, "/workspaces", tt.statusCode, tt.apiErr != nil, tt.want)

			cli, err := xata.NewDatabasesClient(xata.WithBaseURL(testSrv.URL), xata.WithAPIKey("test-key"))
			if err != nil {
				t.Fatal(err)
			}

			got, err := cli.Rename(context.TODO(), tt.request)

			if tt.apiErr != nil {
				errAPI := tt.apiErr.Unwrap()
				if errAPI == nil {
					t.Fatal("expected error but got nil")
				}
				assert.ErrorAs(err, &errAPI)
				assert.Equal(err.Error(), tt.apiErr.Error())
				assert.Nil(got)
			} else {
				assert.Equal(tt.want.Name, got.Name)
				assert.NoError(err)
			}
		})
	}
}
