package xata

import (
	"context"
	"fmt"

	xatagenworkspace "github.com/xataio/xata-go/xata/internal/fern-workspace/generated/go"
	xatagenclient "github.com/xataio/xata-go/xata/internal/fern-workspace/generated/go/core"
)

type TableClient interface {
	Create(ctx context.Context, request TableRequest) (*xatagenworkspace.CreateTableResponse, error)
	Delete(ctx context.Context, request TableRequest) (*xatagenworkspace.DeleteTableResponse, error)
	AddColumn(ctx context.Context, request AddColumnRequest) (*xatagenworkspace.AddTableColumnResponse, error)
	DeleteColumn(ctx context.Context, request DeleteColumnRequest) (*xatagenworkspace.DeleteColumnResponse, error)
}

type tableClient struct {
	generated  xatagenworkspace.TableClient
	dbName     string
	branchName string
}

func (t tableClient) dbBranchName(request TableRequest) string {
	if request.DatabaseName == nil {
		request.DatabaseName = String(t.dbName)
	}

	if request.BranchName == nil {
		request.BranchName = String(t.branchName)
	}

	return fmt.Sprintf("%s:%s", *request.DatabaseName, *request.BranchName)
}

type TableRequest struct {
	DatabaseName *string
	BranchName   *string
	TableName    string
}

func (t tableClient) Create(ctx context.Context, request TableRequest) (*xatagenworkspace.CreateTableResponse, error) {
	return t.generated.CreateTable(ctx, t.dbBranchName(request), request.TableName)
}

func (t tableClient) Delete(ctx context.Context, request TableRequest) (*xatagenworkspace.DeleteTableResponse, error) {
	return t.generated.DeleteTable(ctx, t.dbBranchName(request), request.TableName)
}

type ColumnType xatagenworkspace.ColumnType

type ColumnLink xatagenworkspace.ColumnLink

type ColumnVector xatagenworkspace.ColumnVector

type ColumnFile xatagenworkspace.ColumnFile

const (
	ColumnTypeBool ColumnType = iota + 1
	ColumnTypeInt
	ColumnTypeFloat
	ColumnTypeString
	ColumnTypeText
	ColumnTypeEmail
	ColumnTypeMultiple
	ColumnTypeLink
	ColumnTypeObject
	ColumnTypeDatetime
	ColumnTypeVector
	ColumnTypeFile
	ColumnTypeFileMap
	ColumnTypeJSON
)

type Column struct {
	Name         string        `json:"name"`
	Type         ColumnType    `json:"type,omitempty"`
	Link         *ColumnLink   `json:"link,omitempty"`
	Vector       *ColumnVector `json:"vector,omitempty"`
	File         *ColumnFile   `json:"file,omitempty"`
	FileMap      *ColumnFile   `json:"fileMap,omitempty"`
	NotNull      *bool         `json:"notNull,omitempty"`
	DefaultValue *string       `json:"defaultValue,omitempty"`
	Unique       *bool         `json:"unique,omitempty"`
	Columns      *[]*Column    `json:"columns,omitempty"`
}

type AddColumnRequest struct {
	TableRequest
	Column *Column
}

func (t tableClient) AddColumn(ctx context.Context, request AddColumnRequest) (*xatagenworkspace.AddTableColumnResponse, error) {
	return t.generated.AddTableColumn(ctx, t.dbBranchName(request.TableRequest), request.TableName, copyColumn(*request.Column))
}

func copyColumn(in Column) *xatagenworkspace.Column {
	return &xatagenworkspace.Column{
		Name:         in.Name,
		Type:         (xatagenworkspace.ColumnType)(in.Type),
		Link:         (*xatagenworkspace.ColumnLink)(in.Link),
		Vector:       (*xatagenworkspace.ColumnVector)(in.Vector),
		File:         (*xatagenworkspace.ColumnFile)(in.File),
		FileMap:      (*xatagenworkspace.ColumnFile)(in.FileMap),
		NotNull:      in.NotNull,
		DefaultValue: in.DefaultValue,
		Unique:       in.Unique,
		// Columns: &[]*xatagenworkspace.Column{}, TODO learn usage
	}
}

type DeleteColumnRequest struct {
	TableRequest
	ColumnName string
}

func (t tableClient) DeleteColumn(ctx context.Context, request DeleteColumnRequest) (*xatagenworkspace.DeleteColumnResponse, error) {
	return t.generated.DeleteColumn(ctx, t.dbBranchName(request.TableRequest), request.TableName, request.ColumnName)
}

func NewTableClient(opts ...ClientOption) (TableClient, error) {
	cliOpts, dbCfg, err := consolidateClientOptionsForWorkspace(opts...)
	if err != nil {
		return nil, err
	}

	return tableClient{
			generated: xatagenworkspace.NewTableClient(
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
