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
	BranchRequestOptional
	TableName string
	Payload   QueryTableRequestPayload
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

type BoosterExpression xatagenworkspace.BoosterExpression

type BranchRequestOptional struct {
	DatabaseName *string
	BranchName   *string
}

type SearchBranchRequest struct {
	BranchRequestOptional
	Payload SearchBranchRequestPayload
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

type SearchAndFilterClient interface {
	Query(ctx context.Context, request QueryTableRequest) (*xatagenworkspace.QueryTableResponse, error)
	SearchBranch(ctx context.Context, request SearchBranchRequest) (*xatagenworkspace.SearchBranchResponse, error)
	SearchTable(ctx context.Context, request SearchTableRequest) (*xatagenworkspace.SearchTableResponse, error)
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
	dbBranchName, err := s.dbBranchName(request.BranchRequestOptional)
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

type TargetExpressionItem xatagenworkspace.TargetExpressionItem

type TargetExpression []*TargetExpressionItem

func NewTargetExpression(columnName string) *TargetExpressionItem {
	return (*TargetExpressionItem)(xatagenworkspace.NewTargetExpressionItemFromString(columnName))
}

type TargetExpressionItemColumn xatagenworkspace.TargetExpressionItemColumn

func NewTargetExpressionWithColumnObject(colObj TargetExpressionItemColumn) *TargetExpressionItem {
	colObjGen := &xatagenworkspace.TargetExpressionItemColumn{
		Column: colObj.Column,
		Weight: colObj.Weight,
	}
	return (*TargetExpressionItem)(xatagenworkspace.NewTargetExpressionItemFromTargetExpressionItemColumn(colObjGen))
}

type ValueBooster struct {
	Column          string
	Value           *xatagenworkspace.ValueBoosterValue
	Factor          float64
	IfMatchesFilter *FilterExpression
}

func NewValueBoosterValueFromString(value string) *xatagenworkspace.ValueBoosterValue {
	return xatagenworkspace.NewValueBoosterValueFromString(value)
}

type BoosterExpressionValueBooster struct {
	ValueBooster *ValueBooster
}

func NewBoosterExpressionFromBoosterExpressionValueBooster(value *BoosterExpressionValueBooster) *BoosterExpression {
	genBoosterExpVal := &xatagenworkspace.BoosterExpressionValueBooster{ValueBooster: &xatagenworkspace.ValueBooster{
		Column:          value.ValueBooster.Column,
		Value:           value.ValueBooster.Value,
		Factor:          value.ValueBooster.Factor,
		IfMatchesFilter: (*xatagenworkspace.FilterExpression)(value.ValueBooster.IfMatchesFilter),
	}}
	return (*BoosterExpression)(xatagenworkspace.NewBoosterExpressionFromBoosterExpressionValueBooster(genBoosterExpVal))
}

type NumericBooster struct {
	// The column in which to look for the value.
	Column string
	// The factor with which to multiply the value of the column before adding it to the item score.
	Factor float64
	// Modifier to be applied to the column value, before being multiplied with the factor. The possible values are:
	//   - none (default).
	//   - log: common logarithm (base 10)
	//   - log1p: add 1 then take the common logarithm. This ensures that the value is positive if the
	//     value is between 0 and 1.
	//   - ln: natural logarithm (base e)
	//   - ln1p: add 1 then take the natural logarithm. This ensures that the value is positive if the
	//     value is between 0 and 1.
	//   - square: raise the value to the power of two.
	//   - sqrt: take the square root of the value.
	//   - reciprocal: reciprocate the value (if the value is `x`, the reciprocal is `1/x`).
	Modifier        *uint8
	IfMatchesFilter *FilterExpression
}

type BoosterExpressionNumericBooster struct {
	NumericBooster *NumericBooster
}

func NewBoosterExpressionFromBoosterExpressionNumericBooster(value *BoosterExpressionNumericBooster) *BoosterExpression {
	genValue := &xatagenworkspace.BoosterExpressionNumericBooster{
		NumericBooster: &xatagenworkspace.NumericBooster{
			Column:          value.NumericBooster.Column,
			Factor:          value.NumericBooster.Factor,
			Modifier:        (*xatagenworkspace.NumericBoosterModifier)(value.NumericBooster.Modifier),
			IfMatchesFilter: (*xatagenworkspace.FilterExpression)(value.NumericBooster.IfMatchesFilter),
		},
	}
	return (*BoosterExpression)(xatagenworkspace.NewBoosterExpressionFromBoosterExpressionNumericBooster(genValue))
}

// DateBooster records based on the value of a datetime column. It is configured via "origin", "scale", and "decay". The further away from the "origin",
// the more the score is decayed. The decay function uses an exponential function. For example if origin is "now", and scale is 10 days and decay is 0.5, it
// should be interpreted as: a record with a date 10 days before/after origin will be boosted 2 times less than a record with the date at origin.
// The result of the exponential function is a boost between 0 and 1. The "factor" allows you to control how impactful this boost is, by multiplying it with a given value.
type DateBooster struct {
	// The column in which to look for the value.
	Column string
	// The datetime (formatted as RFC3339) from where to apply the score decay function. The maximum boost will be applied for records with values at this time.
	// If it is not specified, the current date and time is used.
	Origin *string
	// The duration at which distance from origin the score is decayed with factor, using an exponential function. It is formatted as number + units, for example: `5d`, `20m`, `10s`.
	Scale string
	// The decay factor to expect at "scale" distance from the "origin".
	Decay float64
	// The factor with which to multiply the added boost.
	Factor          *float64
	IfMatchesFilter *FilterExpression
}

type BoosterExpressionDateBooster struct {
	DateBooster *DateBooster `json:"dateBooster,omitempty"`
}

func NewBoosterExpressionFromBoosterExpressionDateBooster(value *BoosterExpressionDateBooster) *BoosterExpression {
	genVal := &xatagenworkspace.BoosterExpressionDateBooster{
		DateBooster: &xatagenworkspace.DateBooster{
			Column:          value.DateBooster.Column,
			Origin:          value.DateBooster.Origin,
			Scale:           value.DateBooster.Scale,
			Decay:           value.DateBooster.Decay,
			Factor:          value.DateBooster.Factor,
			IfMatchesFilter: (*xatagenworkspace.FilterExpression)(value.DateBooster.IfMatchesFilter),
		},
	}
	return (*BoosterExpression)(xatagenworkspace.NewBoosterExpressionFromBoosterExpressionDateBooster(genVal))
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
