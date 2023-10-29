package xata_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xataio/xata-go/xata"
	xatagencore "github.com/xataio/xata-go/xata/internal/fern-workspace/generated/go/core"
)

func TestNewRecordsClient(t *testing.T) {
	t.Run("should construct a new client", func(t *testing.T) {
		got, err := xata.NewRecordsClient(
			xata.WithBaseURL("https://www.example.com"),
			xata.WithAPIKey("my-api-token"),
		)
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})
}

var errTestCasesWorkspace = []struct {
	name       string
	statusCode int
	apiErr     *xatagencore.APIError
}{
	{
		name:       "should return bad request error",
		statusCode: http.StatusBadRequest,
		apiErr:     xatagencore.NewAPIError(http.StatusBadRequest, testErrBody),
	},
	{
		name:       "should return authentication error",
		statusCode: http.StatusUnauthorized,
		apiErr:     xatagencore.NewAPIError(http.StatusUnauthorized, testErrBody),
	},
	{
		name:       "should return not-found error",
		statusCode: http.StatusNotFound,
		apiErr:     xatagencore.NewAPIError(http.StatusNotFound, testErrBody),
	},
	{
		name:       "should return server error",
		statusCode: http.StatusServiceUnavailable,
		apiErr:     xatagencore.NewAPIError(http.StatusServiceUnavailable, testErrBody),
	},
	{
		name:       "should handle undocumented (not in the api specs) error",
		statusCode: http.StatusConflict,
		apiErr:     xatagencore.NewAPIError(http.StatusConflict, testErrBody),
	},
}

func Test_recordsClient_Get(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		want       *xata.Record
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should get a record successfully",
			want: &xata.Record{
				RecordMeta: xata.RecordMeta{Id: "some-id"},
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

			cli, err := xata.NewRecordsClient(xata.WithBaseURL(testSrv.URL), xata.WithAPIKey("test-key"))
			assert.NoError(err)
			assert.NotNil(cli)

			got, err := cli.Get(context.TODO(), xata.GetRecordRequest{
				RecordRequest: xata.RecordRequest{
					DatabaseName: xata.String("test-db"),
					BranchName:   xata.String("main"),
					TableName:    "test-table",
				},
				RecordID: "test-id",
				Columns:  []string{"test-column"},
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
				assert.Equal(tt.want.Id, got.Id)
				assert.NoError(err)
			}
		})
	}
}

func Test_recordsClient_Insert(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		want       *xata.Record
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should insert a record successfully",
			want: &xata.Record{
				RecordMeta: xata.RecordMeta{Id: "some-id"},
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

			cli, err := xata.NewRecordsClient(xata.WithBaseURL(testSrv.URL), xata.WithAPIKey("test-key"))
			assert.NoError(err)
			assert.NotNil(cli)

			got, err := cli.Insert(context.TODO(), xata.InsertRecordRequest{
				RecordRequest: xata.RecordRequest{
					DatabaseName: xata.String("test-db"),
					BranchName:   xata.String("main"),
					TableName:    "test-table",
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
				assert.NoError(err)
				assert.Equal(tt.want.Id, got.Id)
			}
		})
	}
}

func Test_recordsClient_InsertWithID(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		want       *xata.Record
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should insert a record with ID successfully",
			want: &xata.Record{
				RecordMeta: xata.RecordMeta{Id: "some-id"},
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

			cli, err := xata.NewRecordsClient(xata.WithBaseURL(testSrv.URL), xata.WithAPIKey("test-key"))
			assert.NoError(err)
			assert.NotNil(cli)

			got, err := cli.InsertWithID(context.TODO(), xata.InsertRecordWithIDRequest{
				RecordRequest: xata.RecordRequest{
					DatabaseName: xata.String("test-db"),
					BranchName:   xata.String("main"),
					TableName:    "test-table",
				},
				RecordID:  "test-id",
				IfVersion: xata.Int(29),
				Columns:   []string{"test-column"},
				Body:      map[string]*xata.DataInputRecordValue{"test": xata.ValueFromBoolean(true)},
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
				assert.Equal(tt.want.Id, got.Id)
			}
		})
	}
}

func Test_recordsClient_Update(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		want       *xata.Record
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should update a record successfully",
			want: &xata.Record{
				RecordMeta: xata.RecordMeta{Id: "some-id"},
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
			testSrv := testService(t, http.MethodPatch, "/db", tt.statusCode, tt.apiErr != nil, tt.want)

			cli, err := xata.NewRecordsClient(xata.WithBaseURL(testSrv.URL), xata.WithAPIKey("test-key"))
			assert.NoError(err)
			assert.NotNil(cli)

			got, err := cli.Update(context.TODO(), xata.UpdateRecordRequest{
				RecordRequest: xata.RecordRequest{
					DatabaseName: xata.String("test-db"),
					BranchName:   xata.String("main"),
					TableName:    "test-table",
				},
				RecordID: "test-id",
				Columns:  []string{"test-column"},
				Body:     map[string]*xata.DataInputRecordValue{"test": xata.ValueFromBoolean(true)},
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
				assert.Equal(tt.want.Id, got.Id)
			}
		})
	}
}
