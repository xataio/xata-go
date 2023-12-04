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

func TestNewFilesClient(t *testing.T) {
	t.Run("should construct a new client", func(t *testing.T) {
		got, err := xata.NewFilesClient(
			xata.WithBaseURL("https://www.example.com"),
			xata.WithAPIKey("my-api-token"),
		)
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})
}

func Test_filesClient_Delete(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		want       *xatagenworkspace.FileResponse
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should delete a file successfully",
			want: &xatagenworkspace.FileResponse{
				Name: "test-file.txt",
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
			testSrv := testService(t, http.MethodDelete, "/db", tt.statusCode, tt.apiErr != nil, tt.want)

			cli, err := xata.NewFilesClient(xata.WithBaseURL(testSrv.URL), xata.WithAPIKey("test-key"))
			assert.NoError(err)
			assert.NotNil(cli)

			got, err := cli.Delete(context.TODO(), xata.DeleteFileRequest{
				BranchRequestOptional: xata.BranchRequestOptional{
					DatabaseName: xata.String("my-db"),
				},
				TableName:  "my-table",
				RecordID:   "my-id",
				ColumnName: "file-column",
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
				assert.Equal(tt.want.Name, got.Name)
				assert.NoError(err)
			}
		})
	}
}

func Test_filesClient_Put(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		want       *xatagenworkspace.FileResponse
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should put a file successfully",
			want: &xatagenworkspace.FileResponse{
				MediaType: "text/plain",
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
			testSrv := testService(t, http.MethodPut, "/db", tt.statusCode, tt.apiErr != nil, tt.want)

			cli, err := xata.NewFilesClient(xata.WithBaseURL(testSrv.URL), xata.WithAPIKey("test-key"))
			assert.NoError(err)
			assert.NotNil(cli)

			got, err := cli.Put(context.TODO(), xata.PutFileRequest{
				BranchRequestOptional: xata.BranchRequestOptional{
					DatabaseName: xata.String("my-db"),
				},
				TableName:  "my-table",
				RecordID:   "my-id",
				ColumnName: "file-column",
				Data:       []byte(`hola`),
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
				assert.Equal(tt.want.MediaType, got.MediaType)
				assert.NoError(err)
			}
		})
	}
}

func Test_filesClient_Get(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name:       "should get a file successfully",
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
			testSrv := testService(t, http.MethodGet, "/db", tt.statusCode, tt.apiErr != nil, nil)

			cli, err := xata.NewFilesClient(xata.WithBaseURL(testSrv.URL), xata.WithAPIKey("test-key"))
			assert.NoError(err)
			assert.NotNil(cli)

			got, err := cli.Get(context.TODO(), xata.GetFileRequest{
				BranchRequestOptional: xata.BranchRequestOptional{
					DatabaseName: xata.String("my-db"),
				},
				TableName:  "my-table",
				RecordID:   "my-id",
				ColumnName: "file-column",
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
				assert.NotNil(got)
				assert.NoError(err)
			}
		})
	}
}

func Test_filesClient_GetItem(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name:       "should get a file item successfully",
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
			testSrv := testService(t, http.MethodGet, "/db", tt.statusCode, tt.apiErr != nil, nil)

			cli, err := xata.NewFilesClient(xata.WithBaseURL(testSrv.URL), xata.WithAPIKey("test-key"))
			assert.NoError(err)
			assert.NotNil(cli)

			got, err := cli.GetItem(context.TODO(), xata.GetFileItemRequest{
				BranchRequestOptional: xata.BranchRequestOptional{
					DatabaseName: xata.String("my-db"),
				},
				TableName:  "my-table",
				RecordID:   "my-id",
				ColumnName: "file-column",
				FileID:     "some-id",
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
				assert.NotNil(got)
				assert.NoError(err)
			}
		})
	}
}

func Test_filesClient_PutItem(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		want       *xatagenworkspace.FileResponse
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should put a file item successfully",
			want: &xatagenworkspace.FileResponse{
				MediaType: "text/plain",
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
			testSrv := testService(t, http.MethodPut, "/db", tt.statusCode, tt.apiErr != nil, tt.want)

			cli, err := xata.NewFilesClient(xata.WithBaseURL(testSrv.URL), xata.WithAPIKey("test-key"))
			assert.NoError(err)
			assert.NotNil(cli)

			got, err := cli.PutItem(context.TODO(), xata.PutFileItemRequest{
				BranchRequestOptional: xata.BranchRequestOptional{
					DatabaseName: xata.String("my-db"),
				},
				TableName:  "my-table",
				RecordID:   "my-id",
				ColumnName: "file-column",
				FileID:     "some-id",
				Data:       []byte(`hola`),
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
				assert.Equal(tt.want.MediaType, got.MediaType)
				assert.NoError(err)
			}
		})
	}
}

func Test_filesClient_DeleteItem(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		want       *xatagenworkspace.FileResponse
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should delete a file item successfully",
			want: &xatagenworkspace.FileResponse{
				Name: "test-file.txt",
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
			testSrv := testService(t, http.MethodDelete, "/db", tt.statusCode, tt.apiErr != nil, tt.want)

			cli, err := xata.NewFilesClient(xata.WithBaseURL(testSrv.URL), xata.WithAPIKey("test-key"))
			assert.NoError(err)
			assert.NotNil(cli)

			got, err := cli.DeleteItem(context.TODO(), xata.DeleteFileItemRequest{
				BranchRequestOptional: xata.BranchRequestOptional{
					DatabaseName: xata.String("my-db"),
				},
				TableName:  "my-table",
				RecordID:   "my-id",
				ColumnName: "file-column",
				FileID:     "some-id",
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
				assert.Equal(tt.want.Name, got.Name)
				assert.NoError(err)
			}
		})
	}
}
