package xata

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	xatagenclient "github.com/omerdemirok/xata-go/xata/internal/fern-workspace/generated/go/core"

	xatagenworkspace "github.com/omerdemirok/xata-go/xata/internal/fern-workspace/generated/go"
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

type InputFile xatagenworkspace.InputFile

type recordsClient struct {
	generated  xatagenworkspace.RecordsClient
	dbName     string
	branchName string
}

func (r recordsClient) Insert(ctx context.Context, request InsertRecordRequest) (*Record, error) {
	insRecReq := &xatagenworkspace.InsertRecordRequest{
		Columns: constructColumns(request.Columns),
		Body:    make(map[string]*xatagenworkspace.DataInputRecordValue),
	}

	for k, v := range request.Body {
		insRecReq.Body[k] = (*xatagenworkspace.DataInputRecordValue)(v)
	}

	record, err := r.generated.InsertRecord(ctx, r.dbBranchName(request.RecordRequest), request.TableName, insRecReq)
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

	record, err := r.generated.GetRecord(
		ctx,
		r.dbBranchName(request.RecordRequest),
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

func (r recordsClient) dbBranchName(request RecordRequest) string {
	if request.DatabaseName == nil {
		request.DatabaseName = String(r.dbName)
	}

	if request.BranchName == nil {
		request.BranchName = String(r.branchName)
	}

	return fmt.Sprintf("%s:%s", *request.DatabaseName, *request.BranchName)
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
	cfg, err := loadConfig(configFileName)
	if err != nil {
		return nil, fmt.Errorf("unable to load config: %w", err)
	}

	dbCfg, err := parseDatabaseURL(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %w", err)
	}

	defaultOpts := &ClientOptions{
		BaseURL: fmt.Sprintf(
			"https://%s.%s.%s",
			dbCfg.workspaceID,
			dbCfg.region,
			dbCfg.domainWorkspace,
		),
		HTTPClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(defaultOpts)
	}

	if defaultOpts.Bearer == "" {
		apiKey, err := getAPIKey()
		if err != nil {
			log.Fatal(err)
		}
		defaultOpts.Bearer = apiKey
	}

	return recordsClient{
			generated: xatagenworkspace.NewRecordsClient(
				func(options *xatagenclient.ClientOptions) {
					options.HTTPClient = defaultOpts.HTTPClient
					options.BaseURL = defaultOpts.BaseURL
					options.Bearer = defaultOpts.Bearer
				}),
			dbName:     dbCfg.dbName,
			branchName: dbCfg.branchName,
		},
		nil
}
