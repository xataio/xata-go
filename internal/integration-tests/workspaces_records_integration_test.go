package integrationtests

import (
	"context"
	"fmt"
	"os"
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
	apiKey, found := os.LookupEnv("XATA_API_KEY")
	if !found {
		t.Skipf("%s not found in env vars", "XATA_API_KEY")
	}

	cfg, err := setup()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err = cleanup(cfg)
		if err != nil {
			t.Fatal(err)
		}
	})

	ctx := context.TODO()
	recordsCli, err := xata.NewRecordsClient(
		xata.WithAPIKey(apiKey),
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
		record, err = recordsCli.Get(ctx, getRecordRequest)
		assert.NoError(t, err)
		assert.NotNil(t, record)
		assert.Equal(t, insertRecordRequest.Body[emailColumn].String, record.Data[emailColumn])
		assert.Equal(t, insertRecordRequest.Body[boolColumn].Boolean, record.Data[boolColumn])
		assert.Equal(t, insertRecordRequest.Body[stringColumn].String, record.Data[stringColumn])
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
			textColumn:     xata.ValueFromString("test-for-text-column"),
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
