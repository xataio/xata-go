// SPDX-License-Identifier: Apache-2.0

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
	DeleteItem(ctx context.Context, request DeleteFileItemRequest) (*xatagenworkspace.FileResponse, error)
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
	RecordID   string
	ColumnName string
}

// Delete removes the content from a file column.
// https://xata.io/docs/api-reference/db/db_branch_name/tables/table_name/data/record_id/column/column_name/file#remove-the-content-from-a-file-column
func (f filesClient) Delete(ctx context.Context, request DeleteFileRequest) (*xatagenworkspace.FileResponse, error) {
	dbBranchName, err := f.dbBranchName(request.BranchRequestOptional)
	if err != nil {
		return nil, err
	}

	return f.generated.DeleteFile(ctx, dbBranchName, request.TableName, request.RecordID, request.ColumnName)
}

type PutFileRequest struct {
	BranchRequestOptional
	ContentType *string
	TableName   string
	RecordID    string
	ColumnName  string
	Data        []byte
}

// Put uploads content to a file column.
// https://xata.io/docs/api-reference/db/db_branch_name/tables/table_name/data/record_id/column/column_name/file#upload-content-to-a-file-column
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

	return f.generated.PutFile(ctx, dbBranchName, request.TableName, request.RecordID, request.ColumnName, request.Data)
}

type GetFileRequest struct {
	BranchRequestOptional
	TableName  string
	RecordID   string
	ColumnName string
}

// Get downloads content from a file column.
// https://xata.io/docs/api-reference/db/db_branch_name/tables/table_name/data/record_id/column/column_name/file#download-content-from-a-file-column
func (f filesClient) Get(ctx context.Context, request GetFileRequest) (*xatagenworkspace.GetFileResponse, error) {
	dbBranchName, err := f.dbBranchName(request.BranchRequestOptional)
	if err != nil {
		return nil, err
	}

	return f.generated.GetFile(ctx, dbBranchName, request.TableName, request.RecordID, request.ColumnName)
}

type GetFileItemRequest struct {
	BranchRequestOptional
	TableName  string
	RecordID   string
	ColumnName string
	FileID     string
}

// GetItem downloads content from a file item in a file array column.
// https://xata.io/docs/api-reference/db/db_branch_name/tables/table_name/data/record_id/column/column_name/file/file_id#download-content-from-a-file-item-in-a-file-array-column
func (f filesClient) GetItem(ctx context.Context, request GetFileItemRequest) (*xatagenworkspace.GetFileResponse, error) {
	dbBranchName, err := f.dbBranchName(request.BranchRequestOptional)
	if err != nil {
		return nil, err
	}

	return f.generated.GetFileItem(ctx, dbBranchName, request.TableName, request.RecordID, request.ColumnName, request.FileID)
}

type PutFileItemRequest struct {
	BranchRequestOptional
	ContentType *string
	TableName   string
	RecordID    string
	ColumnName  string
	FileID      string
	Data        []byte
}

// PutItem uploads or updates the content of a file item in a file array column.
// https://xata.io/docs/api-reference/db/db_branch_name/tables/table_name/data/record_id/column/column_name/file/file_id#upload-or-update-the-content-of-a-file-item-in-a-file-array-column
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

	return f.generated.PutFileItem(ctx, dbBranchName, request.TableName, request.RecordID, request.ColumnName, request.FileID, request.Data)
}

type DeleteFileItemRequest struct {
	BranchRequestOptional
	TableName  string
	RecordID   string
	ColumnName string
	FileID     string
}

// DeleteItem deletes an item from a file array.
// https://xata.io/docs/api-reference/db/db_branch_name/tables/table_name/data/record_id/column/column_name/file/file_id#delete-an-item-from-a-file-array
func (f filesClient) DeleteItem(ctx context.Context, request DeleteFileItemRequest) (*xatagenworkspace.FileResponse, error) {
	dbBranchName, err := f.dbBranchName(request.BranchRequestOptional)
	if err != nil {
		return nil, err
	}

	return f.generated.DeleteFileItem(ctx, dbBranchName, request.TableName, request.RecordID, request.ColumnName, request.FileID)
}

// NewFilesClient constructs a client for interacting files.
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
