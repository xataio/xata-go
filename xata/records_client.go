package xata

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	xatagenclient "github.com/xataio/xata-go/xata/internal/fern-workspace/generated/go/core"

	xatagenworkspace "github.com/xataio/xata-go/xata/internal/fern-workspace/generated/go"
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

type GetRecordRequest struct {
	RecordRequest
	RecordID string
	Columns  []string
}

type RecordMeta struct {
	Id   string                           `json:"id"` // nolint
	Xata *xatagenworkspace.RecordMetaXata `json:"xata,omitempty"`
}

type Record struct {
	RecordMeta
	Data map[string]interface{}
}

type RecordsClient interface {
	Insert(ctx context.Context, request InsertRecordRequest) (*Record, error)
	Update(ctx context.Context, request UpdateRecordRequest) (*Record, error)
	InsertWithID(ctx context.Context, request InsertRecordWithIDRequest) (*Record, error)
	Get(ctx context.Context, request GetRecordRequest) (*Record, error)
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

type recordsClient struct {
	generated  xatagenworkspace.RecordsClient
	dbName     string
	branchName string
}

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

func constructRecord(response map[string]interface{}) (*Record, error) {
	rawResponse, err := json.Marshal(response)
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

	for k, v := range response {
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
