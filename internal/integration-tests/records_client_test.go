package integrationtests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/stretchr/testify/assert"
	"github.com/xataio/xata-go/xata"
)

const (
	stringColumn   = "user-name"
	boolColumn     = "active"
	textColumn     = "text-column"
	emailColumn    = "email"
	dateTimeColumn = "date-of-birth"
	integerColumn  = "integer-column"
	floatColumn    = "float-column"
	fileColumn     = "file-column"
	jsonColumn     = "json-column"
	vectorColumn   = "vector-column" // it is important to set a vector dimension on the UI: 2
	multipleColumn = "multiple-column"
	testFileName   = "file-name.txt"
)

func Test_recordsClient_Insert_Get(t *testing.T) {
	cfg, err := setupDatabase()
	if err != nil {
		t.Fatalf("unable to setup database: %v", err)
	}

	ctx := context.TODO()
	err = setupTableWithColumns(ctx, cfg)
	if err != nil {
		t.Fatalf("unable to setup table: %v", err)
	}

	t.Cleanup(func() {
		err = cleanup(cfg)
		if err != nil {
			t.Fatalf("unable to cleanup test setup: %v", err)
		}
	})

	recordsCli, err := xata.NewRecordsClient(
		xata.WithAPIKey(cfg.apiKey),
		xata.WithBaseURL(fmt.Sprintf(
			"https://%s.%s.xata.sh",
			cfg.wsID,
			cfg.region,
		)),
		xata.WithHTTPClient(retryablehttp.NewClient().StandardClient()),
	)
	if err != nil {
		t.Fatal(err)
	}

	databaseName := cfg.databaseName
	tableName := cfg.tableName
	t.Run("should create a record", func(t *testing.T) {
		insertRecordRequest := generateInsertRecordRequest(databaseName, tableName)

		record, err := recordsCli.Insert(ctx, insertRecordRequest)
		assert.NoError(t, err)
		assert.NotNil(t, record)
		assert.Equal(t, insertRecordRequest.Body[emailColumn].String, record.Data[emailColumn])
		assert.Equal(t, insertRecordRequest.Body[boolColumn].Boolean, record.Data[boolColumn])
		assert.Equal(t, insertRecordRequest.Body[stringColumn].String, record.Data[stringColumn])
		assert.Equal(t, insertRecordRequest.Body[textColumn].String, record.Data[textColumn])
		assert.Equal(t, insertRecordRequest.Body[integerColumn].Double, record.Data[integerColumn])
		assert.Equal(t, insertRecordRequest.Body[floatColumn].Double, record.Data[floatColumn])
		assert.Equal(t, insertRecordRequest.Body[fileColumn].InputFile.Name, record.Data[fileColumn].(map[string]interface{})["name"])
		assert.ElementsMatch(t, insertRecordRequest.Body[vectorColumn].DoubleList, record.Data[vectorColumn])
		assert.ElementsMatch(t, insertRecordRequest.Body[multipleColumn].StringList, record.Data[multipleColumn])
		assert.Equal(t, insertRecordRequest.Body[jsonColumn].String, record.Data[jsonColumn])
	})

	t.Run("should bulk insert records", func(t *testing.T) {
		insertRecordRequest := generateInsertRecordRequest(databaseName, tableName)

		// 400: {"errors":[{
		//"status":400,"message":"column [file-column]: file upload not permitted in transaction"},
		delete(insertRecordRequest.Body, fileColumn)

		records, err := recordsCli.BulkInsert(ctx, xata.BulkInsertRecordRequest{
			RecordRequest: insertRecordRequest.RecordRequest,
			Columns:       insertRecordRequest.Columns,
			Records: []map[string]*xata.DataInputRecordValue{
				insertRecordRequest.Body,
				insertRecordRequest.Body,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, records)
		assert.Equal(t, 2, len(records))

		for _, record := range records {
			assert.Equal(t, insertRecordRequest.Body[emailColumn].String, record.Data[emailColumn])
			assert.Equal(t, insertRecordRequest.Body[boolColumn].Boolean, record.Data[boolColumn])
			assert.Equal(t, insertRecordRequest.Body[stringColumn].String, record.Data[stringColumn])
			assert.Equal(t, insertRecordRequest.Body[textColumn].String, record.Data[textColumn])
			assert.Equal(t, insertRecordRequest.Body[integerColumn].Double, record.Data[integerColumn])
			assert.Equal(t, insertRecordRequest.Body[floatColumn].Double, record.Data[floatColumn])
			assert.ElementsMatch(t, insertRecordRequest.Body[vectorColumn].DoubleList, record.Data[vectorColumn])
			assert.ElementsMatch(t, insertRecordRequest.Body[multipleColumn].StringList, record.Data[multipleColumn])
			assert.Equal(t, insertRecordRequest.Body[jsonColumn].String, record.Data[jsonColumn])
		}
	})

	t.Run("should create a record with ID and update/upsert it", func(t *testing.T) {
		providedRecordID := "random-string-for-ID"
		insertRecordRequest := generateInsertRecordRequest(databaseName, tableName)

		record, err := recordsCli.InsertWithID(ctx, xata.InsertRecordWithIDRequest{
			RecordRequest: insertRecordRequest.RecordRequest,
			RecordID:      providedRecordID,
			CreateOnly:    xata.Bool(false),
			IfVersion:     xata.Int(1),
			Columns:       insertRecordRequest.Columns,
			Body:          insertRecordRequest.Body,
		})
		assert.NoError(t, err)
		assert.NotNil(t, record)
		assert.Equal(t, insertRecordRequest.Body[emailColumn].String, record.Data[emailColumn])
		assert.Equal(t, insertRecordRequest.Body[boolColumn].Boolean, record.Data[boolColumn])
		assert.Equal(t, insertRecordRequest.Body[stringColumn].String, record.Data[stringColumn])
		assert.Equal(t, insertRecordRequest.Body[textColumn].String, record.Data[textColumn])
		assert.Equal(t, insertRecordRequest.Body[integerColumn].Double, record.Data[integerColumn])
		assert.Equal(t, insertRecordRequest.Body[floatColumn].Double, record.Data[floatColumn])
		assert.Equal(t, insertRecordRequest.Body[fileColumn].InputFile.Name, record.Data[fileColumn].(map[string]interface{})["name"])
		assert.ElementsMatch(t, insertRecordRequest.Body[vectorColumn].DoubleList, record.Data[vectorColumn])
		assert.ElementsMatch(t, insertRecordRequest.Body[multipleColumn].StringList, record.Data[multipleColumn])
		assert.Equal(t, insertRecordRequest.Body[jsonColumn].String, record.Data[jsonColumn])

		record, err = recordsCli.Update(ctx, xata.UpdateRecordRequest{
			RecordRequest: insertRecordRequest.RecordRequest,
			RecordID:      providedRecordID,
			Columns:       insertRecordRequest.Columns,
			Body:          insertRecordRequest.Body,
		})
		assert.NoError(t, err)
		assert.NotNil(t, record)
		assert.Equal(t, insertRecordRequest.Body[emailColumn].String, record.Data[emailColumn])
		assert.Equal(t, insertRecordRequest.Body[boolColumn].Boolean, record.Data[boolColumn])
		assert.Equal(t, insertRecordRequest.Body[stringColumn].String, record.Data[stringColumn])
		assert.Equal(t, insertRecordRequest.Body[textColumn].String, record.Data[textColumn])
		assert.Equal(t, insertRecordRequest.Body[integerColumn].Double, record.Data[integerColumn])
		assert.Equal(t, insertRecordRequest.Body[floatColumn].Double, record.Data[floatColumn])
		assert.Equal(t, insertRecordRequest.Body[fileColumn].InputFile.Name, record.Data[fileColumn].(map[string]interface{})["name"])
		assert.ElementsMatch(t, insertRecordRequest.Body[vectorColumn].DoubleList, record.Data[vectorColumn])
		assert.ElementsMatch(t, insertRecordRequest.Body[multipleColumn].StringList, record.Data[multipleColumn])
		assert.Equal(t, insertRecordRequest.Body[jsonColumn].String, record.Data[jsonColumn])

		record, err = recordsCli.Upsert(ctx, xata.UpsertRecordRequest{
			RecordRequest: insertRecordRequest.RecordRequest,
			RecordID:      providedRecordID,
			Columns:       insertRecordRequest.Columns,
			Body:          insertRecordRequest.Body,
		})
		assert.NoError(t, err)
		assert.NotNil(t, record)
		assert.Equal(t, insertRecordRequest.Body[emailColumn].String, record.Data[emailColumn])
		assert.Equal(t, insertRecordRequest.Body[boolColumn].Boolean, record.Data[boolColumn])
		assert.Equal(t, insertRecordRequest.Body[stringColumn].String, record.Data[stringColumn])
		assert.Equal(t, insertRecordRequest.Body[textColumn].String, record.Data[textColumn])
		assert.Equal(t, insertRecordRequest.Body[integerColumn].Double, record.Data[integerColumn])
		assert.Equal(t, insertRecordRequest.Body[floatColumn].Double, record.Data[floatColumn])
		assert.Equal(t, insertRecordRequest.Body[fileColumn].InputFile.Name, record.Data[fileColumn].(map[string]interface{})["name"])
		assert.ElementsMatch(t, insertRecordRequest.Body[vectorColumn].DoubleList, record.Data[vectorColumn])
		assert.ElementsMatch(t, insertRecordRequest.Body[multipleColumn].StringList, record.Data[multipleColumn])
		assert.Equal(t, insertRecordRequest.Body[jsonColumn].String, record.Data[jsonColumn])
	})

	t.Run("should get a record", func(t *testing.T) {
		// first, create a record
		insertRecordRequest := generateInsertRecordRequest(databaseName, tableName)

		record, err := recordsCli.Insert(ctx, insertRecordRequest)
		assert.NoError(t, err)
		assert.NotNil(t, record)

		// retrieve the record
		getRecordRequest := xata.GetRecordRequest{
			RecordRequest: xata.RecordRequest{
				DatabaseName: xata.String(databaseName),
				TableName:    tableName,
			},
			RecordID: record.RecordMeta.Id,
		}
		recordRetrieved, err := recordsCli.Get(ctx, getRecordRequest)
		assert.NoError(t, err)
		assert.NotNil(t, recordRetrieved)
		assert.Equal(t, insertRecordRequest.Body[emailColumn].String, recordRetrieved.Data[emailColumn])
		assert.Equal(t, insertRecordRequest.Body[boolColumn].Boolean, recordRetrieved.Data[boolColumn])
		assert.Equal(t, insertRecordRequest.Body[stringColumn].String, recordRetrieved.Data[stringColumn])
	})

	t.Run("should get a record with get transaction and columns in query", func(t *testing.T) {
		// first, create a record
		insertRecordRequest := generateInsertRecordRequest(databaseName, tableName)

		record, err := recordsCli.Insert(ctx, insertRecordRequest)
		assert.NoError(t, err)
		assert.NotNil(t, record)

		// retrieve the record
		columns := []string{stringColumn, emailColumn, boolColumn}
		transactionRes, err := recordsCli.Transaction(ctx, xata.TransactionRequest{
			RecordRequest: xata.RecordRequest{
				DatabaseName: xata.String(databaseName),
				TableName:    tableName,
			},
			Operations: []xata.TransactionOperation{
				xata.NewGetTransaction(xata.TransactionGetOp{
					Table:   tableName,
					Id:      record.Id,
					Columns: &columns,
				}),
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, transactionRes)
		assert.Equal(t, insertRecordRequest.Body[emailColumn].String, (*transactionRes.Results[0].Columns)[emailColumn])
		assert.Equal(t, insertRecordRequest.Body[boolColumn].Boolean, (*transactionRes.Results[0].Columns)[boolColumn])
		assert.Equal(t, insertRecordRequest.Body[stringColumn].String, (*transactionRes.Results[0].Columns)[stringColumn])
	})

	t.Run("should get a record with get transaction", func(t *testing.T) {
		// first, create a record
		insertRecordRequest := generateInsertRecordRequest(databaseName, tableName)

		record, err := recordsCli.Insert(ctx, insertRecordRequest)
		assert.NoError(t, err)
		assert.NotNil(t, record)

		// retrieve the record
		transactionRes, err := recordsCli.Transaction(ctx, xata.TransactionRequest{
			RecordRequest: xata.RecordRequest{
				DatabaseName: xata.String(databaseName),
				TableName:    tableName,
			},
			Operations: []xata.TransactionOperation{
				xata.NewGetTransaction(xata.TransactionGetOp{
					Table: tableName,
					Id:    record.Id,
				}),
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, transactionRes)
		assert.Equal(t, record.Id, (*transactionRes.Results[0].Columns)["id"])
		assert.NotNil(t, (*transactionRes.Results[0].Columns)["xata"])
	})

	t.Run("should insert a record with insert transaction", func(t *testing.T) {
		stringVal := "test-from-insert-transaction"
		columns := []string{stringColumn, emailColumn, boolColumn}
		transactionRes, err := recordsCli.Transaction(ctx, xata.TransactionRequest{
			RecordRequest: xata.RecordRequest{
				DatabaseName: xata.String(databaseName),
				TableName:    tableName,
			},
			Operations: []xata.TransactionOperation{
				xata.NewInsertTransaction(xata.TransactionInsertOp{
					Table:   tableName,
					Record:  map[string]any{stringColumn: stringVal},
					Columns: &columns,
				}),
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, transactionRes)
		assert.Equal(t, stringVal, (*transactionRes.Results[0].Columns)[stringColumn])
	})

	t.Run("should update a record with update transaction", func(t *testing.T) {
		// insert a record first
		stringVal := "test-from-insert-transaction"
		transactionRes, err := recordsCli.Transaction(ctx, xata.TransactionRequest{
			RecordRequest: xata.RecordRequest{
				DatabaseName: xata.String(databaseName),
				TableName:    tableName,
			},
			Operations: []xata.TransactionOperation{
				xata.NewInsertTransaction(xata.TransactionInsertOp{
					Table:  tableName,
					Record: map[string]any{stringColumn: stringVal},
				}),
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, transactionRes)

		// then update it
		updatedStrValue := "this-is-updated"
		columns := []string{stringColumn}
		transactionResUpdate, err := recordsCli.Transaction(ctx, xata.TransactionRequest{
			RecordRequest: xata.RecordRequest{
				DatabaseName: xata.String(databaseName),
				TableName:    tableName,
			},
			Operations: []xata.TransactionOperation{
				xata.NewUpdateTransaction(xata.TransactionUpdateOp{
					Table:   tableName,
					Id:      transactionRes.Results[0].Id,
					Fields:  map[string]any{stringColumn: updatedStrValue},
					Columns: &columns,
				}),
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, transactionResUpdate)
		assert.Equal(t, updatedStrValue, (*transactionResUpdate.Results[0].Columns)[stringColumn])
		assert.Equal(t, transactionRes.Results[0].Id, transactionResUpdate.Results[0].Id)
		assert.Equal(t, 1, transactionResUpdate.Results[0].Rows)
	})

	t.Run("should delete a record with delete transaction", func(t *testing.T) {
		// insert a record first
		stringVal := "test-from-insert-transaction"
		transactionRes, err := recordsCli.Transaction(ctx, xata.TransactionRequest{
			RecordRequest: xata.RecordRequest{
				DatabaseName: xata.String(databaseName),
				TableName:    tableName,
			},
			Operations: []xata.TransactionOperation{
				xata.NewInsertTransaction(xata.TransactionInsertOp{
					Table:  tableName,
					Record: map[string]any{stringColumn: stringVal},
				}),
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, transactionRes)

		// then delete it
		transactionResDel, err := recordsCli.Transaction(ctx, xata.TransactionRequest{
			RecordRequest: xata.RecordRequest{
				DatabaseName: xata.String(databaseName),
				TableName:    tableName,
			},
			Operations: []xata.TransactionOperation{
				xata.NewDeleteTransaction(xata.TransactionDeleteOp{
					Table: tableName,
					Id:    transactionRes.Results[0].Id,
				}),
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, transactionResDel)
		assert.Equal(t, 1, transactionResDel.Results[0].Rows)
	})

	t.Run("should get a record with filtering by columns", func(t *testing.T) {
		// first, create a record
		insertRecordRequest := generateInsertRecordRequest(databaseName, tableName)

		record, err := recordsCli.Insert(ctx, insertRecordRequest)
		assert.NoError(t, err)
		assert.NotNil(t, record)

		// retrieve the record
		getRecordRequest := xata.GetRecordRequest{
			RecordRequest: xata.RecordRequest{
				DatabaseName: xata.String(databaseName),
				TableName:    tableName,
			},
			RecordID: record.RecordMeta.Id,
			Columns:  []string{stringColumn},
		}
		record, err = recordsCli.Get(ctx, getRecordRequest)
		assert.NoError(t, err)
		assert.NotNil(t, record)
		assert.Equal(t, insertRecordRequest.Body[stringColumn].String, record.Data[stringColumn])
		assert.Nil(t, record.Data[emailColumn]) // filtered out from the response
		assert.Nil(t, record.Data[boolColumn])  // filtered out from the response
	})

	t.Run("should fail to create a record when provided a non existing column name", func(t *testing.T) {
		req := xata.InsertRecordRequest{
			RecordRequest: xata.RecordRequest{
				DatabaseName: xata.String(databaseName),
				TableName:    tableName,
			},
			Body: map[string]*xata.DataInputRecordValue{
				"made-up-column-name": xata.ValueFromString("test-value-from-SDK-integration-" + time.Now().String()),
			},
		}

		recordResp, err := recordsCli.Insert(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, recordResp)
	})
}

func generateInsertRecordRequest(databaseName, tableName string) xata.InsertRecordRequest {
	return xata.InsertRecordRequest{
		RecordRequest: xata.RecordRequest{
			DatabaseName: xata.String(databaseName),
			TableName:    tableName,
		},
		Columns: []string{
			emailColumn,
			boolColumn,
			dateTimeColumn,
			stringColumn,
			textColumn,
			integerColumn,
			floatColumn,
			fileColumn,
			jsonColumn,
			vectorColumn,
			multipleColumn,
		},
		Body: map[string]*xata.DataInputRecordValue{
			stringColumn:   xata.ValueFromString("test-value-from-SDK-integration-" + time.Now().String()),
			emailColumn:    xata.ValueFromString("test-value-from-SDK-integration@test.com"),
			boolColumn:     xata.ValueFromBoolean(true),
			dateTimeColumn: xata.ValueFromDateTime(time.Now()),
			textColumn:     xata.ValueFromString(textContent),
			integerColumn:  xata.ValueFromInteger(10),
			floatColumn:    xata.ValueFromDouble(10.3),
			fileColumn: xata.ValueFromInputFile(xata.InputFile{
				Name:          testFileName,
				Base64Content: xata.String("ZmlsZSBjb250ZW50"), // file content
			}),
			vectorColumn:   xata.ValueFromDoubleList([]float64{10.3, 20.2}),
			multipleColumn: xata.ValueFromStringList([]string{"hello", "world"}),
			jsonColumn:     xata.ValueFromString(`{"key":"value"}`),
		},
	}
}

var textContent = "An atom is a particle that consists of a nucleus of protons and neutrons surrounded by an electromagnetically-bound cloud of electrons. The atom is the basic particle of the chemical elements, and the chemical elements are distinguished from each other by the number of protons that are in their atoms. For example, any atom that contains 11 protons is sodium, and any atom that contains 29 protons is copper. The number of neutrons defines the isotope of the element.\n\nAtoms are extremely small, typically around 100 picometers across. A human hair is about a million carbon atoms wide. This is smaller than the shortest wavelength of visible light, which means humans cannot see atoms with conventional microscopes. Atoms are so small that accurately predicting their behavior using classical physics is not possible due to quantum effects.\n\nMore than 99.94% of an atom's mass is in the nucleus. Each proton has a positive electric charge, while each electron has a negative charge, and the neutrons, if any are present, have no electric charge. If the numbers of protons and electrons are equal, as they normally are, then the atom is electrically neutral. If an atom has more electrons than protons, then it has an overall negative charge, and is called a negative ion (or anion). Conversely, if it has more protons than electrons, it has a positive charge, and is called a positive ion (or cation)."
