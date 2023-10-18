package xata

import (
	"context"
	"fmt"

	xatagenworkspace "github.com/xataio/xata-go/xata/internal/fern-workspace/generated/go"
	xatagenclient "github.com/xataio/xata-go/xata/internal/fern-workspace/generated/go/core"
)

type SearchRequest struct {
	DatabaseName *string
	BranchName   *string
	TableName    string
}

type FilterExpression xatagenworkspace.FilterExpression

type SortExpression xatagenworkspace.SortExpression

type PageConfig xatagenworkspace.PageConfig

type QueryTableRequestConsistency uint8

const (
	QueryTableRequestConsistencyStrong QueryTableRequestConsistency = iota + 1
	QueryTableRequestConsistencyEventual
)

type QueryTableRequestPayload struct {
	Filter      *FilterExpression
	Sort        *SortExpression
	Page        *PageConfig
	Columns     *[]string
	Consistency QueryTableRequestConsistency
}

type QueryTableRequest struct {
	SearchRequest
	Payload QueryTableRequestPayload
}

type SearchAndFilterClient interface {
	Query(ctx context.Context, request QueryTableRequest) (*xatagenworkspace.QueryTableResponse, error)
	//SearchBranch(ctx context.Context, dbBranchName DbBranchName, request *SearchBranchRequest) (*SearchBranchResponse, error)
	//SearchTable(ctx context.Context, dbBranchName DbBranchName, tableName TableName, request *SearchTableRequest) (*SearchTableResponse, error)
	//VectorSearchTable(ctx context.Context, dbBranchName DbBranchName, tableName TableName, request *VectorSearchTableRequest) (*VectorSearchTableResponse, error)
	//AskTable(ctx context.Context, dbBranchName DbBranchName, tableName TableName, request *AskTableRequest) (*AskTableResponse, error)
	//AskTableSession(ctx context.Context, dbBranchName DbBranchName, tableName TableName, sessionId string, request *AskTableSessionRequest) (*AskTableSessionResponse, error)
	//SummarizeTable(ctx context.Context, dbBranchName DbBranchName, tableName TableName, request *SummarizeTableRequest) (*SummarizeTableResponse, error)
	//AggregateTable(ctx context.Context, dbBranchName DbBranchName, tableName TableName, request *AggregateTableRequest) (*AggregateTableResponse, error)
}

type searchAndFilterCli struct {
	generated  xatagenworkspace.SearchAndFilterClient
	dbName     string
	branchName string
}

func (s searchAndFilterCli) dbBranchName(request SearchRequest) (string, error) {
	if request.DatabaseName == nil {
		if s.dbName == "" {
			return "", fmt.Errorf("database name cannot be empty")
		}
		request.DatabaseName = String(s.dbName)
	}

	if request.BranchName == nil {
		if s.branchName == "" {
			return "", fmt.Errorf("branch name cannot be empty")
		}
		request.BranchName = String(s.branchName)
	}

	return fmt.Sprintf("%s:%s", *request.DatabaseName, *request.BranchName), nil
}

func (s searchAndFilterCli) Query(ctx context.Context, request QueryTableRequest) (*xatagenworkspace.QueryTableResponse, error) {
	dbBranchName, err := s.dbBranchName(request.SearchRequest)
	if err != nil {
		return nil, err
	}

	payload := &xatagenworkspace.QueryTableRequest{
		Filter:      (*xatagenworkspace.FilterExpression)(request.Payload.Filter),
		Sort:        (*xatagenworkspace.SortExpression)(request.Payload.Sort),
		Page:        (*xatagenworkspace.PageConfig)(request.Payload.Page),
		Columns:     request.Payload.Columns,
		Consistency: (*xatagenworkspace.QueryTableRequestConsistency)(&request.Payload.Consistency),
	}
	return s.generated.QueryTable(ctx, dbBranchName, request.TableName, payload)
}

func NewSearchAndFilterClient(opts ...ClientOption) (SearchAndFilterClient, error) {
	cliOpts, dbCfg, err := consolidateClientOptionsForWorkspace(opts...)
	if err != nil {
		return nil, err
	}

	return searchAndFilterCli{
			generated: xatagenworkspace.NewSearchAndFilterClient(
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
