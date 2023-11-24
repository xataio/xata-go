// SPDX-License-Identifier: Apache-2.0

package xata

import (
	"context"
	"fmt"

	xatagenworkspace "github.com/xataio/xata-go/xata/internal/fern-workspace/generated/go"
	xatagenclient "github.com/xataio/xata-go/xata/internal/fern-workspace/generated/go/core"
)

type SearchAndFilterClient interface {
	Query(ctx context.Context, request QueryTableRequest) (*xatagenworkspace.QueryTableResponse, error)
	SearchBranch(ctx context.Context, request SearchBranchRequest) (*xatagenworkspace.SearchBranchResponse, error)
	SearchTable(ctx context.Context, request SearchTableRequest) (*xatagenworkspace.SearchTableResponse, error)
	VectorSearch(ctx context.Context, request VectorSearchTableRequest) (*xatagenworkspace.VectorSearchTableResponse, error)
	Ask(ctx context.Context, request AskTableRequest) (*xatagenworkspace.AskTableResponse, error)
	AskFollowUp(ctx context.Context, request AskFollowUpRequest) (*xatagenworkspace.AskTableSessionResponse, error)
	Summarize(ctx context.Context, request SummarizeTableRequest) (*xatagenworkspace.SummarizeTableResponse, error)
	Aggregate(ctx context.Context, request AggregateTableRequest) (*xatagenworkspace.AggregateTableResponse, error)
}

type BranchRequestOptional struct {
	DatabaseName *string
	BranchName   *string
}

type searchAndFilterCli struct {
	generated  xatagenworkspace.SearchAndFilterClient
	dbName     string
	branchName string
}

func (s searchAndFilterCli) dbBranchName(request BranchRequestOptional) (string, error) {
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

type QueryTableRequestPayload struct {
	Filter      *FilterExpression
	Sort        *xatagenworkspace.SortExpression
	Page        *PageConfig
	Columns     []string
	Consistency QueryTableRequestConsistency
}

type QueryTableRequest struct {
	BranchRequestOptional
	TableName string
	Payload   QueryTableRequestPayload
}

func (s searchAndFilterCli) Query(ctx context.Context, request QueryTableRequest) (*xatagenworkspace.QueryTableResponse, error) {
	dbBranchName, err := s.dbBranchName(request.BranchRequestOptional)
	if err != nil {
		return nil, err
	}

	return s.generated.QueryTable(ctx, dbBranchName, request.TableName, &xatagenworkspace.QueryTableRequest{
		Filter:  (*xatagenworkspace.FilterExpression)(request.Payload.Filter),
		Sort:    request.Payload.Sort,
		Page:    (*xatagenworkspace.PageConfig)(request.Payload.Page),
		Columns: &request.Payload.Columns,
		//Consistency: (*xatagenworkspace.QueryTableRequestConsistency)(&request.Payload.Consistency),
	})
}

type SearchBranchRequestPayload struct {
	Tables    []*SearchBranchRequestTablesItem
	Query     string
	Fuzziness *int
	Prefix    *PrefixExpression
	Highlight *HighlightExpression
	Page      *SearchPageConfig
}

type SearchBranchRequest struct {
	BranchRequestOptional
	Payload SearchBranchRequestPayload
}

func (s searchAndFilterCli) SearchBranch(ctx context.Context, request SearchBranchRequest) (*xatagenworkspace.SearchBranchResponse, error) {
	dbBranchName, err := s.dbBranchName(request.BranchRequestOptional)
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

type SearchTableRequestPayload struct {
	Query     string
	Fuzziness *int
	Target    TargetExpression
	Prefix    *PrefixExpression
	Filter    *FilterExpression
	Highlight *HighlightExpression
	Boosters  []*BoosterExpression
	Page      *SearchPageConfig
}

type SearchTableRequest struct {
	BranchRequestOptional
	TableName string
	Payload   SearchTableRequestPayload
}

func (s searchAndFilterCli) SearchTable(ctx context.Context, request SearchTableRequest) (*xatagenworkspace.SearchTableResponse, error) {
	dbBranchName, err := s.dbBranchName(request.BranchRequestOptional)
	if err != nil {
		return nil, err
	}

	var boostersGen []*xatagenworkspace.BoosterExpression
	if len(request.Payload.Boosters) > 0 {
		for _, b := range request.Payload.Boosters {
			boostersGen = append(boostersGen, (*xatagenworkspace.BoosterExpression)(b))
		}
	}

	var targetExpGen []*xatagenworkspace.TargetExpressionItem
	if len(request.Payload.Target) > 0 {
		for _, e := range request.Payload.Target {
			targetExpGen = append(targetExpGen, (*xatagenworkspace.TargetExpressionItem)(e))
		}
	}

	return s.generated.SearchTable(ctx, dbBranchName, request.TableName, &xatagenworkspace.SearchTableRequest{
		Query:     request.Payload.Query,
		Fuzziness: request.Payload.Fuzziness,
		Target:    &targetExpGen,
		Prefix:    (*xatagenworkspace.PrefixExpression)(request.Payload.Prefix),
		Filter:    (*xatagenworkspace.FilterExpression)(request.Payload.Filter),
		Highlight: (*xatagenworkspace.HighlightExpression)(request.Payload.Highlight),
		Boosters:  &boostersGen,
		Page:      (*xatagenworkspace.SearchPageConfig)(request.Payload.Page),
	})
}

type VectorSearchTableRequestPayload struct {
	QueryVector        []float64
	Column             string
	SimilarityFunction *string
	Size               *int
	Filter             *FilterExpression
}

type VectorSearchTableRequest struct {
	BranchRequestOptional
	TableName string
	Payload   VectorSearchTableRequestPayload
}

func (s searchAndFilterCli) VectorSearch(ctx context.Context, request VectorSearchTableRequest) (*xatagenworkspace.VectorSearchTableResponse, error) {
	dbBranchName, err := s.dbBranchName(request.BranchRequestOptional)
	if err != nil {
		return nil, err
	}

	return s.generated.VectorSearchTable(ctx, dbBranchName, request.TableName, &xatagenworkspace.VectorSearchTableRequest{
		QueryVector:        request.Payload.QueryVector,
		Column:             request.Payload.Column,
		SimilarityFunction: request.Payload.SimilarityFunction,
		Size:               request.Payload.Size,
		Filter:             (*xatagenworkspace.FilterExpression)(request.Payload.Filter),
	})
}

type AskTableRequestPayload struct {
	// The question you'd like to ask.
	Question string
	// The type of search to use. If set to `keyword` (the default), the search can be configured by passing
	// a `search` object with the following fields. For more details about each, see the Search endpoint documentation.
	// All fields are optional.
	//   - fuzziness  - typo tolerance
	//   - target - columns to search into, and weights.
	//   - prefix - prefix search type.
	//   - filter - pre-filter before searching.
	//   - boosters - control relevancy.
	//
	// If set to `vector`, a `vectorSearch` object must be passed, with the following parameters. For more details, see the Vector
	// Search endpoint documentation. The `column` and `contentColumn` parameters are required.
	//   - column - the vector column containing the embeddings.
	//   - contentColumn - the column that contains the text from which the embeddings where computed.
	//   - filter - pre-filter before searching.
	SearchType   *AskTableRequestSearchType
	Search       *AskTableRequestSearch
	VectorSearch *AskTableRequestVectorSearch
	Rules        *[]string
}

type AskTableRequest struct {
	BranchRequestOptional
	TableName string
	Payload   AskTableRequestPayload
}

func (s searchAndFilterCli) Ask(ctx context.Context, request AskTableRequest) (*xatagenworkspace.AskTableResponse, error) {
	dbBranchName, err := s.dbBranchName(request.BranchRequestOptional)
	if err != nil {
		return nil, err
	}

	var targetExpGen []*xatagenworkspace.TargetExpressionItem
	var searchGen *xatagenworkspace.AskTableRequestSearch
	if request.Payload.Search != nil {
		searchGen = &xatagenworkspace.AskTableRequestSearch{
			Fuzziness: request.Payload.Search.Fuzziness,
			Target:    &targetExpGen,
		}

		if len(request.Payload.Search.Target) > 0 {
			for _, e := range request.Payload.Search.Target {
				targetExpGen = append(targetExpGen, (*xatagenworkspace.TargetExpressionItem)(e))
			}
		}
	}

	var vectorSearchGen *xatagenworkspace.AskTableRequestVectorSearch
	if request.Payload.VectorSearch != nil {
		vectorSearchGen = &xatagenworkspace.AskTableRequestVectorSearch{
			Column:        request.Payload.VectorSearch.Column,
			ContentColumn: request.Payload.VectorSearch.ContentColumn,
			Filter:        (*xatagenworkspace.FilterExpression)(request.Payload.VectorSearch.Filter),
		}
	}

	return s.generated.AskTable(ctx, dbBranchName, request.TableName, &xatagenworkspace.AskTableRequest{
		Question:     request.Payload.Question,
		SearchType:   (*xatagenworkspace.AskTableRequestSearchType)(request.Payload.SearchType),
		Search:       searchGen,
		VectorSearch: vectorSearchGen,
		Rules:        request.Payload.Rules,
	})
}

type AskFollowUpRequest struct {
	BranchRequestOptional
	TableName string
	SessionID string
	Question  string
}

func (s searchAndFilterCli) AskFollowUp(ctx context.Context, request AskFollowUpRequest) (*xatagenworkspace.AskTableSessionResponse, error) {
	dbBranchName, err := s.dbBranchName(request.BranchRequestOptional)
	if err != nil {
		return nil, err
	}

	return s.generated.AskTableSession(
		ctx,
		dbBranchName,
		request.TableName,
		request.SessionID,
		&xatagenworkspace.AskTableSessionRequest{Message: String(request.Question)},
	)
}

type SummarizeTableRequestPayload struct {
	Filter          *FilterExpression
	Columns         []string
	Summaries       map[string]map[string]any
	Sort            *xatagenworkspace.SortExpression
	SummariesFilter *FilterExpression
	// The consistency level for this request.
	//Consistency  *SummarizeTableRequestConsistency // https://github.com/xataio/xata-go/pull/37#issue-2009859238
	NumberOfPage *int
}

type SummarizeTableRequest struct {
	BranchRequestOptional
	TableName string
	Payload   SummarizeTableRequestPayload
}

func (s searchAndFilterCli) Summarize(ctx context.Context, request SummarizeTableRequest) (*xatagenworkspace.SummarizeTableResponse, error) {
	dbBranchName, err := s.dbBranchName(request.BranchRequestOptional)
	if err != nil {
		return nil, err
	}

	var sumExpList xatagenworkspace.SummaryExpressionList
	if len(request.Payload.Summaries) > 0 {
		sumExpList = make(xatagenworkspace.SummaryExpressionList, len(request.Payload.Summaries))
		for k, v := range request.Payload.Summaries {
			if len(v) > 0 {
				sumExp := make(xatagenworkspace.SummaryExpression, len(v))
				for k1, v1 := range v {
					sumExp[k1] = v1
				}
				sumExpList[k] = sumExp
			}
		}
	}

	return s.generated.SummarizeTable(ctx, dbBranchName, request.TableName, &xatagenworkspace.SummarizeTableRequest{
		Filter:          (*xatagenworkspace.FilterExpression)(request.Payload.Filter),
		Columns:         &request.Payload.Columns,
		Summaries:       &sumExpList,
		Sort:            request.Payload.Sort,
		SummariesFilter: (*xatagenworkspace.FilterExpression)(request.Payload.SummariesFilter),
		//Consistency:     (*xatagenworkspace.SummarizeTableRequestConsistency)(request.Payload.Consistency),
		Page: &xatagenworkspace.SummarizeTableRequestPage{
			Size: request.Payload.NumberOfPage,
		},
	})
}

type AggregateTableRequestPayload struct {
	Filter       *FilterExpression
	Aggregations AggExpressionMap
}

type AggregateTableRequest struct {
	BranchRequestOptional
	TableName string
	Payload   AggregateTableRequestPayload
}

func (s searchAndFilterCli) Aggregate(ctx context.Context, request AggregateTableRequest) (*xatagenworkspace.AggregateTableResponse, error) {
	dbBranchName, err := s.dbBranchName(request.BranchRequestOptional)
	if err != nil {
		return nil, err
	}

	var aggsGen xatagenworkspace.AggExpressionMap
	if len(request.Payload.Aggregations) > 0 {
		aggsGen = make(xatagenworkspace.AggExpressionMap, len(request.Payload.Aggregations))
		for k, v := range request.Payload.Aggregations {
			aggsGen[k] = v
		}
	}

	return s.generated.AggregateTable(ctx, dbBranchName, request.TableName, &xatagenworkspace.AggregateTableRequest{
		Filter: (*xatagenworkspace.FilterExpression)(request.Payload.Filter),
		Aggs:   &aggsGen,
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
