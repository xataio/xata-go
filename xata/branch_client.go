package xata

import (
	"context"
	"fmt"
	xatagenworkspace "github.com/xataio/xata-go/xata/internal/fern-workspace/generated/go"
	xatagenclient "github.com/xataio/xata-go/xata/internal/fern-workspace/generated/go/core"
)

type BranchRequest struct {
	DatabaseName *string
	BranchName   *string
}

type CreateBranchRequestPayload xatagenworkspace.CreateBranchRequest

type CreateBranchRequest struct {
	BranchName     string
	DatabaseName   *string
	From           *string
	BranchMetadata BranchMetadata
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

func (b branchCli) dbBranchName(request CreateBranchRequest) (string, error) {
	if request.DatabaseName == nil {
		if b.dbName == "" {
			return "", fmt.Errorf("database name cannot be empty")
		}
		request.DatabaseName = String(b.dbName)
	}

	if request.BranchName == "" {
		return "", fmt.Errorf("branch name cannot be empty")
	}

	return fmt.Sprintf("%s:%s", *request.DatabaseName, request.BranchName), nil
}

func (b branchCli) List(ctx context.Context, dbName string) (*xatagenworkspace.ListBranchesResponse, error) {
	return b.generated.GetBranchList(ctx, dbName)
}

func (b branchCli) GetDetails(ctx context.Context, request BranchRequest) (*xatagenworkspace.DbBranch, error) {
	//TODO implement me
	panic("implement me")
}

func (b branchCli) Create(ctx context.Context, request CreateBranchRequest) (*xatagenworkspace.CreateBranchResponse, error) {
	dbBranchName, err := b.dbBranchName(request)
	if err != nil {
		return nil, err
	}

	req := &xatagenworkspace.CreateBranchRequest{
		From:                    request.From,
		CreateBranchRequestFrom: request.From,
		Metadata:                (*xatagenworkspace.BranchMetadata)(&request.BranchMetadata),
	}
	return b.generated.CreateBranch(ctx, dbBranchName, req)
}

func (b branchCli) Delete(ctx context.Context, request BranchRequest) (*xatagenworkspace.DeleteBranchResponse, error) {
	//TODO implement me
	panic("implement me")
}

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
