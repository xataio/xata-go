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
				RecordId:   "my-id",
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
