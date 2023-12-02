// SPDX-License-Identifier: Apache-2.0

package xata

import (
	"context"
	"fmt"

	xatagenworkspace "github.com/xataio/xata-go/xata/internal/fern-workspace/generated/go"
	xatagenclient "github.com/xataio/xata-go/xata/internal/fern-workspace/generated/go/core"
)

type BranchRequest struct {
	DatabaseName *string
	BranchName   string
}

type BranchMetadataWS xatagenworkspace.BranchMetadata

type CreateBranchRequestPayload struct {
	CreateBranchRequestFrom *string
	Metadata                *BranchMetadataWS
}

type CreateBranchRequest struct {
	BranchName   string
	DatabaseName *string
	From         *string
	Payload      *CreateBranchRequestPayload
}

type BranchClient interface {
	List(ctx context.Context, dbName string) (*xatagenworkspace.ListBranchesResponse, error)
	GetDetails(ctx context.Context, request BranchRequest) (*xatagenworkspace.DbBranch, error)
	Create(ctx context.Context, request CreateBranchRequest) (*xatagenworkspace.CreateBranchResponse, error)
	Delete(ctx context.Context, request BranchRequest) (*xatagenworkspace.DeleteBranchResponse, error)
}

type branchCli struct {
	generated  xatagenworkspace.BranchClient
	dbName     string
	branchName string
}

func (b branchCli) dbBranchName(dbName *string, branchName string) (string, error) {
	if dbName == nil {
		if b.dbName == "" {
			return "", fmt.Errorf("database name cannot be empty")
		}
		dbName = String(b.dbName)
	}

	if branchName == "" {
		return "", fmt.Errorf("branch name cannot be empty")
	}

	return fmt.Sprintf("%s:%s", *dbName, branchName), nil
}

// List lists all available branches.
// https://xata.io/docs/api-reference/dbs/db_name#list-branches
func (b branchCli) List(ctx context.Context, dbName string) (*xatagenworkspace.ListBranchesResponse, error) {
	return b.generated.GetBranchList(ctx, dbName)
}

// GetDetails gets branch schema and metadata.
// https://xata.io/docs/api-reference/db/db_branch_name#get-branch-schema-and-metadata
func (b branchCli) GetDetails(ctx context.Context, request BranchRequest) (*xatagenworkspace.DbBranch, error) {
	dbBranchName, err := b.dbBranchName(request.DatabaseName, request.BranchName)
	if err != nil {
		return nil, err
	}

	return b.generated.GetBranchDetails(ctx, dbBranchName)
}

// Create creates a database branch.
// https://xata.io/docs/api-reference/db/db_branch_name#create-database-branch
func (b branchCli) Create(ctx context.Context, request CreateBranchRequest) (*xatagenworkspace.CreateBranchResponse, error) {
	dbBranchName, err := b.dbBranchName(request.DatabaseName, request.BranchName)
	if err != nil {
		return nil, err
	}

	var payloadFrom *string
	if request.Payload != nil && request.Payload.CreateBranchRequestFrom != nil {
		payloadFrom = request.Payload.CreateBranchRequestFrom
	}

	var payloadMetadata *xatagenworkspace.BranchMetadata
	if request.Payload != nil && request.Payload.Metadata != nil {
		payloadMetadata = (*xatagenworkspace.BranchMetadata)(request.Payload.Metadata)
	}

	req := &xatagenworkspace.CreateBranchRequest{
		From:                    request.From,
		CreateBranchRequestFrom: payloadFrom,
		Metadata:                payloadMetadata,
	}
	return b.generated.CreateBranch(ctx, dbBranchName, req)
}

// Delete deletes a database branch.
// https://xata.io/docs/api-reference/db/db_branch_name#delete-database-branch
func (b branchCli) Delete(ctx context.Context, request BranchRequest) (*xatagenworkspace.DeleteBranchResponse, error) {
	dbBranchName, err := b.dbBranchName(request.DatabaseName, request.BranchName)
	if err != nil {
		return nil, err
	}

	return b.generated.DeleteBranch(ctx, dbBranchName)
}

// NewBranchClient constructs a new client to interact with database branches.
func NewBranchClient(opts ...ClientOption) (BranchClient, error) {
	cliOpts, dbCfg, err := consolidateClientOptionsForWorkspace(opts...)
	if err != nil {
		return nil, err
	}

	return branchCli{
		generated: xatagenworkspace.NewBranchClient(
			func(options *xatagenclient.ClientOptions) {
				options.HTTPClient = cliOpts.HTTPClient
				options.BaseURL = cliOpts.BaseURL
				options.Bearer = cliOpts.Bearer
			}),
		dbName:     dbCfg.dbName,
		branchName: dbCfg.branchName,
	}, nil
}
