package xata

import (
	"context"
	"fmt"

	xatagenworkspace "github.com/xataio/xata-go/xata/internal/fern-workspace/generated/go"
	xatagenclient "github.com/xataio/xata-go/xata/internal/fern-workspace/generated/go/core"
)

type FilesClient interface {
	GetItem(ctx context.Context, request GetFileItemRequest) (*xatagenworkspace.GetFileResponse, error)
	PutItem(ctx context.Context, request PutFileItemRequest) (*xatagenworkspace.FileResponse, error)
	// DeleteFileItem(ctx context.Context, dbBranchName DbBranchName, tableName TableName, recordId RecordId, columnName ColumnName, fileId FileItemId) (*FileResponse, error)
	Get(ctx context.Context, request GetFileRequest) (*xatagenworkspace.GetFileResponse, error)
	Put(ctx context.Context, request PutFileRequest) (*xatagenworkspace.FileResponse, error)
	Delete(ctx context.Context, request DeleteFileRequest) (*xatagenworkspace.FileResponse, error)
}

type filesClient struct {
	generated  xatagenworkspace.FilesClient
	dbName     string
	branchName string
}

func (f filesClient) dbBranchName(request BranchRequestOptional) (string, error) {
	if request.DatabaseName == nil {
		if f.dbName == "" {
			return "", fmt.Errorf("database name cannot be empty")
		}
		request.DatabaseName = String(f.dbName)
	}

	if request.BranchName == nil {
		if f.branchName == "" {
			return "", fmt.Errorf("branch name cannot be empty")
		}
		request.BranchName = String(f.branchName)
	}

	return fmt.Sprintf("%s:%s", *request.DatabaseName, *request.BranchName), nil
}

type DeleteFileRequest struct {
	BranchRequestOptional
	TableName  string
	RecordId   string
	ColumnName string
}

func (f filesClient) Delete(ctx context.Context, request DeleteFileRequest) (*xatagenworkspace.FileResponse, error) {
	dbBranchName, err := f.dbBranchName(request.BranchRequestOptional)
	if err != nil {
		return nil, err
	}

	return f.generated.DeleteFile(ctx, dbBranchName, request.TableName, request.RecordId, request.ColumnName)
}

type PutFileRequest struct {
	BranchRequestOptional
	ContentType *string
	TableName   string
	RecordId    string
	ColumnName  string
	Data        []byte
}

func (f filesClient) Put(ctx context.Context, request PutFileRequest) (*xatagenworkspace.FileResponse, error) {
	dbBranchName, err := f.dbBranchName(request.BranchRequestOptional)
	if err != nil {
		return nil, err
	}

	contentType := "application/octet-stream"
	if request.ContentType != nil && *request.ContentType != "" {
		contentType = *request.ContentType
	}

	f.generated.SetContentTypeHeader(contentType)

	return f.generated.PutFile(ctx, dbBranchName, request.TableName, request.RecordId, request.ColumnName, request.Data)
}

type GetFileRequest struct {
	BranchRequestOptional
	TableName  string
	RecordId   string
	ColumnName string
}

func (f filesClient) Get(ctx context.Context, request GetFileRequest) (*xatagenworkspace.GetFileResponse, error) {
	dbBranchName, err := f.dbBranchName(request.BranchRequestOptional)
	if err != nil {
		return nil, err
	}

	return f.generated.GetFile(ctx, dbBranchName, request.TableName, request.RecordId, request.ColumnName)
}

type GetFileItemRequest struct {
	BranchRequestOptional
	TableName  string
	RecordId   string
	ColumnName string
	FileID     string
}

func (f filesClient) GetItem(ctx context.Context, request GetFileItemRequest) (*xatagenworkspace.GetFileResponse, error) {
	dbBranchName, err := f.dbBranchName(request.BranchRequestOptional)
	if err != nil {
		return nil, err
	}

	return f.generated.GetFileItem(ctx, dbBranchName, request.TableName, request.RecordId, request.ColumnName, request.FileID)
}

type PutFileItemRequest struct {
	BranchRequestOptional
	ContentType *string
	TableName   string
	RecordId    string
	ColumnName  string
	FileID      string
	Data        []byte
}

func (f filesClient) PutItem(ctx context.Context, request PutFileItemRequest) (*xatagenworkspace.FileResponse, error) {
	dbBranchName, err := f.dbBranchName(request.BranchRequestOptional)
	if err != nil {
		return nil, err
	}

	contentType := "application/octet-stream"
	if request.ContentType != nil && *request.ContentType != "" {
		contentType = *request.ContentType
	}

	f.generated.SetContentTypeHeader(contentType)

	return f.generated.PutFileItem(ctx, dbBranchName, request.TableName, request.RecordId, request.ColumnName, request.FileID, request.Data)
}

func NewFilesClient(opts ...ClientOption) (FilesClient, error) {
	cliOpts, dbCfg, err := consolidateClientOptionsForWorkspace(opts...)
	if err != nil {
		return nil, err
	}

	return filesClient{
			generated: xatagenworkspace.NewFilesClient(
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
