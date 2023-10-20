package xata

import (
	"context"
	"fmt"

	xatagenworkspace "github.com/xataio/xata-go/xata/internal/fern-workspace/generated/go"
	xatagenclient "github.com/xataio/xata-go/xata/internal/fern-workspace/generated/go/core"
)

type SortOrder xatagenworkspace.SortOrder

const (
	SortOrderAsc SortOrder = iota + 1
	SortOrderDesc
	SortOrderRandom
)

type FilterExpression xatagenworkspace.FilterExpression

type PageConfig xatagenworkspace.PageConfig

type QueryTableRequestConsistency uint8

const (
	QueryTableRequestConsistencyStrong QueryTableRequestConsistency = iota + 1
	QueryTableRequestConsistencyEventual
)

type QueryTableRequestPayload struct {
	Filter      *FilterExpression
	Sort        *xatagenworkspace.SortExpression
	Page        *PageConfig
	Columns     []string
	Consistency QueryTableRequestConsistency
}

type QueryTableRequest struct {
	TableRequest
	Payload QueryTableRequestPayload
}

type PrefixExpression uint8

const (
	PrefixExpressionPhrase PrefixExpression = iota + 1
	PrefixExpressionDisabled
)

type SearchBranchRequestTablesItem xatagenworkspace.SearchBranchRequestTablesItem

type HighlightExpression xatagenworkspace.HighlightExpression

type SearchPageConfig xatagenworkspace.SearchPageConfig

type SearchBranchRequestPayload struct {
	Tables    []*SearchBranchRequestTablesItem
	Query     string
	Fuzziness *int
	Prefix    *PrefixExpression
	Highlight *HighlightExpression
	Page      *SearchPageConfig
}

type SearchBranchRequest struct {
	TableRequest
	Payload SearchBranchRequestPayload
}

type SearchAndFilterClient interface {
	Query(ctx context.Context, request QueryTableRequest) (*xatagenworkspace.QueryTableResponse, error)
	SearchBranch(ctx context.Context, request SearchBranchRequest) (*xatagenworkspace.SearchBranchResponse, error)
	// SearchTable(ctx context.Context, dbBranchName DbBranchName, tableName TableName, request *SearchTableRequest) (*SearchTableResponse, error)
	// VectorSearchTable(ctx context.Context, dbBranchName DbBranchName, tableName TableName, request *VectorSearchTableRequest) (*VectorSearchTableResponse, error)
	// AskTable(ctx context.Context, dbBranchName DbBranchName, tableName TableName, request *AskTableRequest) (*AskTableResponse, error)
	// AskTableSession(ctx context.Context, dbBranchName DbBranchName, tableName TableName, sessionId string, request *AskTableSessionRequest) (*AskTableSessionResponse, error)
	// SummarizeTable(ctx context.Context, dbBranchName DbBranchName, tableName TableName, request *SummarizeTableRequest) (*SummarizeTableResponse, error)
	// AggregateTable(ctx context.Context, dbBranchName DbBranchName, tableName TableName, request *AggregateTableRequest) (*AggregateTableResponse, error)
}

type searchAndFilterCli struct {
	generated  xatagenworkspace.SearchAndFilterClient
	dbName     string
	branchName string
}

func (s searchAndFilterCli) dbBranchName(request TableRequest) (string, error) {
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

func NewSortExpressionFromStringList(value []string) *xatagenworkspace.SortExpression {
	return xatagenworkspace.NewSortExpressionFromStringList(value)
}

func NewSortExpressionFromStringSortOrderMap(value map[string]SortOrder) *xatagenworkspace.SortExpression {
	valGen := make(map[string]xatagenworkspace.SortOrder, len(value))
	for k, v := range value {
		valGen[k] = (xatagenworkspace.SortOrder)(v)
	}

	return xatagenworkspace.NewSortExpressionFromStringSortOrderMap(valGen)
}

func NewSortExpressionFromStringSortOrderMapList(value []map[string]SortOrder) *xatagenworkspace.SortExpression {
	valGen := make([]map[string]xatagenworkspace.SortOrder, len(value))
	for i, vs := range value {
		for k, v := range vs {
			valGen[i] = make(map[string]xatagenworkspace.SortOrder, len(vs))
			valGen[i][k] = (xatagenworkspace.SortOrder)(v)
		}
	}

	return xatagenworkspace.NewSortExpressionFromStringSortOrderMapList(valGen)
}

func NewFilterListFromFilterExpression(value *FilterExpression) *xatagenworkspace.FilterList {
	return xatagenworkspace.NewFilterListFromFilterExpression((*xatagenworkspace.FilterExpression)(value))
}

func NewFilterListFromFilterExpressionList(value []*FilterExpression) *xatagenworkspace.FilterList {
	var valGen []*xatagenworkspace.FilterExpression
	for _, v := range value {
		valGen = append(valGen, (*xatagenworkspace.FilterExpression)(v))
	}
	return xatagenworkspace.NewFilterListFromFilterExpressionList(valGen)
}

func (s searchAndFilterCli) Query(ctx context.Context, request QueryTableRequest) (*xatagenworkspace.QueryTableResponse, error) {
	dbBranchName, err := s.dbBranchName(request.TableRequest)
	if err != nil {
		return nil, err
	}

	return s.generated.QueryTable(ctx, dbBranchName, request.TableName, &xatagenworkspace.QueryTableRequest{
		Filter:      (*xatagenworkspace.FilterExpression)(request.Payload.Filter),
		Sort:        request.Payload.Sort,
		Page:        (*xatagenworkspace.PageConfig)(request.Payload.Page),
		Columns:     &request.Payload.Columns,
		Consistency: (*xatagenworkspace.QueryTableRequestConsistency)(&request.Payload.Consistency),
	})
}

func NewSearchBranchRequestTablesItemFromString(value string) *SearchBranchRequestTablesItem {
	return (*SearchBranchRequestTablesItem)(xatagenworkspace.NewSearchBranchRequestTablesItemFromString(value))
}

func (s searchAndFilterCli) SearchBranch(ctx context.Context, request SearchBranchRequest) (*xatagenworkspace.SearchBranchResponse, error) {
	dbBranchName, err := s.dbBranchName(request.TableRequest)
	if err != nil {
		return nil, err
	}

	var tables []*xatagenworkspace.SearchBranchRequestTablesItem
	if len(request.Payload.Tables) != 0 {
		for _, t := range request.Payload.Tables {
			tables = append(tables, (*xatagenworkspace.SearchBranchRequestTablesItem)(t))
		}
	}

	return s.generated.SearchBranch(ctx, dbBranchName, &xatagenworkspace.SearchBranchRequest{
		Tables:    &tables,
		Query:     request.Payload.Query,
		Fuzziness: request.Payload.Fuzziness,
		Prefix:    (*xatagenworkspace.PrefixExpression)(request.Payload.Prefix),
		Highlight: (*xatagenworkspace.HighlightExpression)(request.Payload.Highlight),
		Page:      (*xatagenworkspace.SearchPageConfig)(request.Payload.Page),
	})
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
