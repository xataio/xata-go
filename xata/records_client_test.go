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

func Test_recordsClient_Upsert(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		want       *xata.Record
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should upsert a record successfully",
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

			got, err := cli.Upsert(context.TODO(), xata.UpsertRecordRequest{
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

func Test_recordsClient_BulkInsert(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		want       []*xata.Record
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should bulk insert successfully",
			want: []*xata.Record{
				{
					RecordMeta: xata.RecordMeta{Id: "some-id"},
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
			testSrv := testService(
				t,
				http.MethodPost,
				"/db",
				tt.statusCode,
				tt.apiErr != nil,
				xatagenworkspace.BulkInsertTableRecordsResponse{"records": []map[string]interface{}{{"id": "some-id"}}},
			)

			cli, err := xata.NewRecordsClient(xata.WithBaseURL(testSrv.URL), xata.WithAPIKey("test-key"))
			assert.NoError(err)
			assert.NotNil(cli)

			got, err := cli.BulkInsert(context.TODO(), xata.BulkInsertRecordRequest{
				RecordRequest: xata.RecordRequest{
					DatabaseName: xata.String("test-db"),
					BranchName:   xata.String("main"),
					TableName:    "test-table",
				},
				Columns: []string{"test-column"},
				Records: []map[string]*xata.DataInputRecordValue{
					{"test": xata.ValueFromBoolean(true)},
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
				assert.Equal(1, len(got))
				assert.Equal(tt.want[0].Id, got[0].Id)
			}
		})
	}
}

func Test_recordsClient_Transaction(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		want       *xatagenworkspace.TransactionSuccess
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name:       "should send a transaction successfully",
			want:       &xatagenworkspace.TransactionSuccess{},
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
			testSrv := testService(
				t,
				http.MethodPost,
				"/db",
				tt.statusCode,
				tt.apiErr != nil,
				tt.want,
			)

			cli, err := xata.NewRecordsClient(xata.WithBaseURL(testSrv.URL), xata.WithAPIKey("test-key"))
			assert.NoError(err)
			assert.NotNil(cli)

			got, err := cli.Transaction(context.TODO(), xata.TransactionRequest{
				RecordRequest: xata.RecordRequest{
					DatabaseName: xata.String("test-db"),
					BranchName:   xata.String("main"),
					TableName:    "test-table",
				},
				Operations: []xata.TransactionOperation{
					xata.NewDeleteTransaction(xata.TransactionDeleteOp{
						Table: "test-table",
						Id:    "some-id",
					}),
					xata.NewGetTransaction(xata.TransactionGetOp{
						Table: "test-table",
						Id:    "some-id",
					}),
					xata.NewUpdateTransaction(xata.TransactionUpdateOp{
						Table:  "test-table",
						Id:     "some-id",
						Fields: map[string]any{"test": "value"},
					}),
					xata.NewInsertTransaction(xata.TransactionInsertOp{
						Table:  "test-table",
						Record: map[string]any{"test": "value"},
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
				assert.NoError(err)
				assert.Equal(tt.want, got)
			}
		})
	}
}

func Test_recordsClient_Delete(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name:       "should delete a record",
			statusCode: http.StatusNoContent,
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
			testSrv := testService(
				t,
				http.MethodDelete,
				"/db",
				tt.statusCode,
				tt.apiErr != nil,
				nil,
			)

			cli, err := xata.NewRecordsClient(xata.WithBaseURL(testSrv.URL), xata.WithAPIKey("test-key"))
			assert.NoError(err)
			assert.NotNil(cli)

			err = cli.Delete(context.TODO(), xata.DeleteRecordRequest{
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
			} else {
				assert.NoError(err)
			}
		})
	}
}
