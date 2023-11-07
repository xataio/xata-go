// SPDX-License-Identifier: Apache-2.0

package xata_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xataio/xata-go/xata"
	xatagenworkspace "github.com/xataio/xata-go/xata/internal/fern-workspace/generated/go"

	xatagencore "github.com/xataio/xata-go/xata/internal/fern-workspace/generated/go/core"
)

func TestNewBranchClient(t *testing.T) {
	t.Run("should construct a new client", func(t *testing.T) {
		got, err := xata.NewBranchClient(
			xata.WithBaseURL("https://www.example.com"),
			xata.WithAPIKey("my-api-token"),
		)
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})
}

func Test_branchCli_Create(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		want       *xatagenworkspace.CreateBranchResponse
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should create a new branch successfully",
			want: &xatagenworkspace.CreateBranchResponse{
				DatabaseName: "db-name",
				BranchName:   "new-branch",
				Status:       xatagenworkspace.MigrationStatusCompleted,
			},
			statusCode: http.StatusCreated,
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
			testSrv := testService(t, http.MethodPut, "/db", tt.statusCode, tt.apiErr != nil, tt.want)

			cli, err := xata.NewBranchClient(
				xata.WithBaseURL(testSrv.URL),
				xata.WithAPIKey("test-key"),
			)
			assert.NoError(err)
			assert.NotNil(cli)

			got, err := cli.Create(context.TODO(), xata.CreateBranchRequest{
				DatabaseName: xata.String("my-db"),
				BranchName:   "new-branch",
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
				assert.NoError(err)
				assert.Equal(tt.want.BranchName, got.BranchName)
			}
		})
	}
}

func Test_branchCli_GetDetails(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		want       *xatagenworkspace.DbBranch
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should create a new branch successfully",
			want: &xatagenworkspace.DbBranch{
				DatabaseName:    "my-db",
				BranchName:      "my-branch",
				CreatedAt:       time.Now(),
				Id:              "some-id",
				Version:         1,
				LastMigrationId: "id",
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
			testSrv := testService(t, http.MethodGet, "/db", tt.statusCode, tt.apiErr != nil, tt.want)

			cli, err := xata.NewBranchClient(
				xata.WithBaseURL(testSrv.URL),
				xata.WithAPIKey("test-key"),
			)
			assert.NoError(err)
			assert.NotNil(cli)

			got, err := cli.GetDetails(context.TODO(), xata.BranchRequest{
				DatabaseName: xata.String("my-db"),
				BranchName:   "my-branch",
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
				assert.NoError(err)
				assert.Equal(tt.want.BranchName, got.BranchName)
			}
		})
	}
}

func Test_branchCli_List(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		want       *xatagenworkspace.ListBranchesResponse
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should create a new branch successfully",
			want: &xatagenworkspace.ListBranchesResponse{
				DatabaseName: "my-db",
				Branches: []*xatagenworkspace.Branch{
					{
						Name:      "my-branch",
						CreatedAt: time.Now(),
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
			testSrv := testService(t, http.MethodGet, "/db", tt.statusCode, tt.apiErr != nil, tt.want)

			cli, err := xata.NewBranchClient(
				xata.WithBaseURL(testSrv.URL),
				xata.WithAPIKey("test-key"),
			)
			assert.NoError(err)
			assert.NotNil(cli)

			got, err := cli.List(context.TODO(), "my-db")

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
				assert.Equal(tt.want.Branches[0].Name, got.Branches[0].Name)
				assert.Equal(tt.want.DatabaseName, got.DatabaseName)
			}
		})
	}
}

func Test_branchCli_Delete(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		want       *xatagenworkspace.DeleteBranchResponse
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name:       "should create a new branch successfully",
			want:       &xatagenworkspace.DeleteBranchResponse{Status: xatagenworkspace.MigrationStatusCompleted},
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
			testSrv := testService(t, http.MethodDelete, "/db", tt.statusCode, tt.apiErr != nil, tt.want)

			cli, err := xata.NewBranchClient(
				xata.WithBaseURL(testSrv.URL),
				xata.WithAPIKey("test-key"),
			)
			assert.NoError(err)
			assert.NotNil(cli)

			got, err := cli.Delete(context.TODO(), xata.BranchRequest{
				DatabaseName: xata.String("my-db"),
				BranchName:   "my-branch",
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
				assert.NoError(err)
				assert.Equal(tt.want.Status, got.Status)
			}
		})
	}
}
