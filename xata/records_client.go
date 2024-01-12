// SPDX-License-Identifier: Apache-2.0

package xata

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	xatagenworkspace "github.com/xataio/xata-go/xata/internal/fern-workspace/generated/go"
	xatagenclient "github.com/xataio/xata-go/xata/internal/fern-workspace/generated/go/core"
)

type RecordRequest struct {
	DatabaseName *string
	BranchName   *string
	TableName    string
}

type InsertRecordRequest struct {
	RecordRequest
	Columns []string
	Body    map[string]*DataInputRecordValue
}

type TransactionRequest struct {
	RecordRequest
	Operations []TransactionOperation
}

type BulkInsertRecordRequest struct {
	RecordRequest
	Columns []string
	Records []map[string]*DataInputRecordValue
}

type InsertRecordWithIDRequest struct {
	RecordRequest
	RecordID   string
	CreateOnly *bool
	IfVersion  *int
	Columns    []string
	Body       map[string]*DataInputRecordValue
}

type UpdateRecordRequest struct {
	RecordRequest
	RecordID  string
	IfVersion *int
	Columns   []string
	Body      map[string]*DataInputRecordValue
}

type UpsertRecordRequest UpdateRecordRequest

type GetRecordRequest struct {
	RecordRequest
	RecordID string
	Columns  []string
}

type DeleteRecordRequest struct {
	RecordRequest
	RecordID string
}

type RecordMeta struct {
	Id   string                           `json:"id"` // nolint
	Xata *xatagenworkspace.RecordMetaXata `json:"xata,omitempty"`
}

type Record struct {
	RecordMeta
	Data map[string]interface{}
}

type BulkRecords struct {
	RecordIDs []string
	Records   []Record
}

type RecordsClient interface {
	Transaction(ctx context.Context, request TransactionRequest) (*xatagenworkspace.TransactionSuccess, error)
	Insert(ctx context.Context, request InsertRecordRequest) (*Record, error)
	BulkInsert(ctx context.Context, request BulkInsertRecordRequest) ([]*Record, error)
	Update(ctx context.Context, request UpdateRecordRequest) (*Record, error)
	Upsert(ctx context.Context, request UpsertRecordRequest) (*Record, error)
	InsertWithID(ctx context.Context, request InsertRecordWithIDRequest) (*Record, error)
	Get(ctx context.Context, request GetRecordRequest) (*Record, error)
	Delete(ctx context.Context, request DeleteRecordRequest) error
}

type DataInputRecordValue xatagenworkspace.DataInputRecordValue

func ValueFromString(value string) *DataInputRecordValue {
	return (*DataInputRecordValue)(xatagenworkspace.NewDataInputRecordValueFromString(value))
}

func ValueFromBoolean(value bool) *DataInputRecordValue {
	return (*DataInputRecordValue)(xatagenworkspace.NewDataInputRecordValueFromBoolean(value))
}

func ValueFromDateTime(value time.Time) *DataInputRecordValue {
	return (*DataInputRecordValue)(xatagenworkspace.NewDataInputRecordValueFromDateTime(value))
}

func ValueFromDouble(value float64) *DataInputRecordValue {
	return (*DataInputRecordValue)(xatagenworkspace.NewDataInputRecordValueFromDouble(value))
}

func ValueFromInteger(value int) *DataInputRecordValue {
	return (*DataInputRecordValue)(xatagenworkspace.NewDataInputRecordValueFromDouble(float64(value)))
}

func ValueFromStringList(value []string) *DataInputRecordValue {
	return (*DataInputRecordValue)(xatagenworkspace.NewDataInputRecordValueFromStringList(value))
}

func ValueFromDoubleList(value []float64) *DataInputRecordValue {
	return (*DataInputRecordValue)(xatagenworkspace.NewDataInputRecordValueFromDoubleList(value))
}

type InputFileEntry xatagenworkspace.InputFileEntry

type InputFileArray []*InputFileEntry

func ValueFromInputFileArray(value InputFileArray) *DataInputRecordValue {
	var xValue xatagenworkspace.InputFileArray
	for _, a := range value {
		xValue = append(xValue, (*xatagenworkspace.InputFileEntry)(a))
	}
	return (*DataInputRecordValue)(xatagenworkspace.NewDataInputRecordValueFromInputFileArray(xValue))
}

func ValueFromInputFile(value InputFile) *DataInputRecordValue {
	v := value
	return (*DataInputRecordValue)(xatagenworkspace.NewDataInputRecordValueFromInputFile((*xatagenworkspace.InputFile)(&v)))
}

type InputFile xatagenworkspace.InputFile

/*
	type InputFile struct {
		// Base64 encoded content
		Base64Content *string `json:"base64Content,omitempty"`
		// Enable public access to the file
		EnablePublicUrl *bool   `json:"enablePublicUrl,omitempty"`
		MediaType       *string `json:"mediaType,omitempty"`
		Name            string  `json:"name"`
		// Time to live for signed URLs
		SignedUrlTimeout *int `json:"signedUrlTimeout,omitempty"`
		// Time to live for upload URLs
		UploadUrlTimeout *int `json:"uploadUrlTimeout,omitempty"`
	}
*/
type recordsClient struct {
	generated  xatagenworkspace.RecordsClient
	dbName     string
	branchName string
}

// Insert inserts a record.
// https://xata.io/docs/api-reference/db/db_branch_name/tables/table_name/data#insert-record
func (r recordsClient) Insert(ctx context.Context, request InsertRecordRequest) (*Record, error) {
	recGen := &xatagenworkspace.InsertRecordRequest{
		Columns: constructColumns(request.Columns),
		Body:    make(map[string]*xatagenworkspace.DataInputRecordValue),
	}

	for k, v := range request.Body {
		recGen.Body[k] = (*xatagenworkspace.DataInputRecordValue)(v)
	}

	dbBranchName, err := r.dbBranchName(request.RecordRequest)
	if err != nil {
		return nil, err
	}

	record, err := r.generated.InsertRecord(ctx, dbBranchName, request.TableName, recGen)
	if err != nil {
		return nil, err
	}

	respRec, err := constructRecord(*record)
	if err != nil {
		return nil, err
	}

	return respRec, nil
}

// BulkInsert bulk inserts records.
// https://xata.io/docs/api-reference/db/db_branch_name/tables/table_name/bulk#bulk-insert-records
func (r recordsClient) BulkInsert(ctx context.Context, request BulkInsertRecordRequest) ([]*Record, error) {
	recGen := &xatagenworkspace.BulkInsertTableRecordsRequest{
		Columns: constructColumns(request.Columns),
	}

	for _, record := range request.Records {
		dataInput := make(map[string]*xatagenworkspace.DataInputRecordValue, len(record))
		for col, val := range record {
			dataInput[col] = (*xatagenworkspace.DataInputRecordValue)(val)
		}
		recGen.Records = append(recGen.Records, dataInput)
	}

	dbBranchName, err := r.dbBranchName(request.RecordRequest)
	if err != nil {
		return nil, err
	}

	records, err := r.generated.BulkInsertTableRecords(ctx, dbBranchName, request.TableName, recGen)
	if err != nil {
		return nil, err
	}

	return constructBulkRecords(*records)
}

// InsertWithID inserts a record with ID.
// https://xata.io/docs/api-reference/db/db_branch_name/tables/table_name/data/record_id#insert-record-with-id
func (r recordsClient) InsertWithID(ctx context.Context, request InsertRecordWithIDRequest) (*Record, error) {
	recGen := &xatagenworkspace.InsertRecordWithIdRequest{
		CreateOnly: request.CreateOnly,
		IfVersion:  request.IfVersion,
		Columns:    constructColumns(request.Columns),
		Body:       make(map[string]*xatagenworkspace.DataInputRecordValue),
	}

	for k, v := range request.Body {
		recGen.Body[k] = (*xatagenworkspace.DataInputRecordValue)(v)
	}

	dbBranchName, err := r.dbBranchName(request.RecordRequest)
	if err != nil {
		return nil, err
	}

	record, err := r.generated.InsertRecordWithId(ctx, dbBranchName, request.TableName, request.RecordID, recGen)
	if err != nil {
		return nil, err
	}

	respRec, err := constructRecord(*record)
	if err != nil {
		return nil, err
	}

	return respRec, nil
}

// Update updates a record.
// https://xata.io/docs/api-reference/db/db_branch_name/tables/table_name/data/record_id#update-record-with-id
func (r recordsClient) Update(ctx context.Context, request UpdateRecordRequest) (*Record, error) {
	recGen := &xatagenworkspace.UpdateRecordWithIdRequest{
		IfVersion: request.IfVersion,
		Columns:   constructColumns(request.Columns),
		Body:      make(map[string]*xatagenworkspace.DataInputRecordValue),
	}

	for k, v := range request.Body {
		recGen.Body[k] = (*xatagenworkspace.DataInputRecordValue)(v)
	}

	dbBranchName, err := r.dbBranchName(request.RecordRequest)
	if err != nil {
		return nil, err
	}

	record, err := r.generated.UpdateRecordWithId(ctx, dbBranchName, request.TableName, request.RecordID, recGen)
	if err != nil {
		return nil, err
	}

	respRec, err := constructRecord(*record)
	if err != nil {
		return nil, err
	}

	return respRec, nil
}

// Upsert inserts or updates a record.
// https://xata.io/docs/api-reference/db/db_branch_name/tables/table_name/data/record_id#upsert-record-with-id
func (r recordsClient) Upsert(ctx context.Context, request UpsertRecordRequest) (*Record, error) {
	recGen := &xatagenworkspace.UpdateRecordWithIdRequest{
		IfVersion: request.IfVersion,
		Columns:   constructColumns(request.Columns),
		Body:      make(map[string]*xatagenworkspace.DataInputRecordValue),
	}

	for k, v := range request.Body {
		recGen.Body[k] = (*xatagenworkspace.DataInputRecordValue)(v)
	}

	dbBranchName, err := r.dbBranchName(request.RecordRequest)
	if err != nil {
		return nil, err
	}

	record, err := r.generated.UpdateRecordWithId(ctx, dbBranchName, request.TableName, request.RecordID, recGen)
	if err != nil {
		return nil, err
	}

	respRec, err := constructRecord(*record)
	if err != nil {
		return nil, err
	}

	return respRec, nil
}

// Get gets a record by its ID.
// https://xata.io/docs/api-reference/db/db_branch_name/tables/table_name/data/record_id#get-record-by-id
func (r recordsClient) Get(ctx context.Context, request GetRecordRequest) (*Record, error) {
	getRecReq := &xatagenworkspace.GetRecordRequest{
		Columns: constructColumns(request.Columns),
	}

	dbBranchName, err := r.dbBranchName(request.RecordRequest)
	if err != nil {
		return nil, err
	}

	record, err := r.generated.GetRecord(
		ctx,
		dbBranchName,
		request.TableName,
		request.RecordID,
		getRecReq,
	)
	if err != nil {
		return nil, err
	}

	respRec, err := constructRecord(*record)
	if err != nil {
		return nil, err
	}

	return respRec, nil
}

type TransactionOperation *xatagenworkspace.TransactionOperation

type TransactionInsertOp xatagenworkspace.TransactionInsertOp

func NewInsertTransaction(value TransactionInsertOp) TransactionOperation {
	return xatagenworkspace.NewTransactionOperationFromTransactionOperationInsert(&xatagenworkspace.TransactionOperationInsert{
		Insert: (*xatagenworkspace.TransactionInsertOp)(&value),
	})
}

type TransactionGetOp xatagenworkspace.TransactionGetOp

func NewGetTransaction(value TransactionGetOp) TransactionOperation {
	return xatagenworkspace.NewTransactionOperationFromTransactionOperationGet(&xatagenworkspace.TransactionOperationGet{
		Get: (*xatagenworkspace.TransactionGetOp)(&value),
	})
}

type TransactionUpdateOp xatagenworkspace.TransactionUpdateOp

func NewUpdateTransaction(value TransactionUpdateOp) TransactionOperation {
	return xatagenworkspace.NewTransactionOperationFromTransactionOperationUpdate(&xatagenworkspace.TransactionOperationUpdate{
		Update: (*xatagenworkspace.TransactionUpdateOp)(&value),
	})
}

type TransactionDeleteOp xatagenworkspace.TransactionDeleteOp

func NewDeleteTransaction(value TransactionDeleteOp) TransactionOperation {
	return xatagenworkspace.NewTransactionOperationFromTransactionOperationDelete(&xatagenworkspace.TransactionOperationDelete{
		Delete: (*xatagenworkspace.TransactionDeleteOp)(&value),
	})
}

// Transaction executes a transaction on a branch.
// https://xata.io/docs/api-reference/db/db_branch_name/transaction#execute-a-transaction-on-a-branch
func (r recordsClient) Transaction(ctx context.Context, request TransactionRequest) (*xatagenworkspace.TransactionSuccess, error) {
	dbBranchName, err := r.dbBranchName(request.RecordRequest)
	if err != nil {
		return nil, err
	}

	var operationsGen []*xatagenworkspace.TransactionOperation
	for _, op := range request.Operations {
		operationsGen = append(operationsGen, op)
	}

	return r.generated.BranchTransaction(ctx, dbBranchName, &xatagenworkspace.BranchTransactionRequest{
		Operations: operationsGen,
	})
}

// Delete deletes a record from a table.
// https://xata.io/docs/api-reference/db/db_branch_name/tables/table_name/data/record_id#delete-record-from-table
func (r recordsClient) Delete(ctx context.Context, request DeleteRecordRequest) error {
	dbBranchName, err := r.dbBranchName(request.RecordRequest)
	if err != nil {
		return err
	}

	return r.generated.DeleteRecord(ctx, dbBranchName, request.TableName, request.RecordID)
}

func (r recordsClient) dbBranchName(request RecordRequest) (string, error) {
	if request.DatabaseName == nil {
		if r.dbName == "" {
			return "", fmt.Errorf("database name cannot be empty")
		}
		request.DatabaseName = String(r.dbName)
	}

	if request.BranchName == nil {
		if r.branchName == "" {
			return "", fmt.Errorf("branch name cannot be empty")
		}
		request.BranchName = String(r.branchName)
	}

	return fmt.Sprintf("%s:%s", *request.DatabaseName, *request.BranchName), nil
}

func constructColumns(columns []string) []*string {
	if len(columns) == 0 {
		return nil
	}

	return []*string{String(strings.Join(columns, ","))}
}

func constructBulkRecords(in xatagenworkspace.BulkInsertTableRecordsResponse) ([]*Record, error) {
	var records []*Record

	// response has a key `records` that holds the records
	for _, rec := range in["records"] {
		record, err := constructRecord(rec)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, nil
}

func constructRecord(in map[string]interface{}) (*Record, error) {
	rawResponse, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	var meta RecordMeta
	err = json.Unmarshal(rawResponse, &meta)
	if err != nil {
		return nil, err
	}

	record := Record{
		RecordMeta: meta,
		Data:       make(map[string]interface{}),
	}

	for k, v := range in {
		if k == "id" {
			continue
		}

		if k == "xata" {
			continue
		}

		record.Data[k] = v
	}

	return &record, nil
}

// NewRecordsClient constructs a client for interacting with records.
func NewRecordsClient(opts ...ClientOption) (RecordsClient, error) {
	cliOpts, dbCfg, err := consolidateClientOptionsForWorkspace(opts...)
	if err != nil {
		return nil, err
	}

	return recordsClient{
			generated: xatagenworkspace.NewRecordsClient(
				func(options *xatagenclient.ClientOptions) {
					options.HTTPClient = cliOpts.HTTPClient
					options.BaseURL = cliOpts.BaseURL
					options.Bearer = cliOpts.Bearer
				}),
			dbName:     dbCfg.dbName,
			branchName: dbCfg.branchName,
		},
		nil
}
