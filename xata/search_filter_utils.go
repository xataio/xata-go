package xata

import "github.com/xataio/xata-go/xata/internal/fern-workspace/generated/go"

type SortOrder api.SortOrder

const (
	SortOrderAsc SortOrder = iota + 1
	SortOrderDesc
	SortOrderRandom
)

type FilterExpression api.FilterExpression

type PageConfig api.PageConfig

type QueryTableRequestConsistency uint8

const (
	ConsistencyStrong QueryTableRequestConsistency = iota + 1
	ConsistencyEventual
)

type PrefixExpression uint8

const (
	PrefixExpressionPhrase PrefixExpression = iota + 1
	PrefixExpressionDisabled
)

type SearchBranchRequestTablesItem api.SearchBranchRequestTablesItem

type HighlightExpression api.HighlightExpression

type SearchPageConfig api.SearchPageConfig

type BoosterExpression api.BoosterExpression

func NewSortExpressionFromStringList(value []string) *api.SortExpression {
	return api.NewSortExpressionFromStringList(value)
}

func NewSortExpressionFromStringSortOrderMap(value map[string]SortOrder) *api.SortExpression {
	valGen := make(map[string]api.SortOrder, len(value))
	for k, v := range value {
		valGen[k] = (api.SortOrder)(v)
	}

	return api.NewSortExpressionFromStringSortOrderMap(valGen)
}

func NewSortExpressionFromStringSortOrderMapList(value []map[string]SortOrder) *api.SortExpression {
	valGen := make([]map[string]api.SortOrder, len(value))
	for i, vs := range value {
		for k, v := range vs {
			valGen[i] = make(map[string]api.SortOrder, len(vs))
			valGen[i][k] = (api.SortOrder)(v)
		}
	}

	return api.NewSortExpressionFromStringSortOrderMapList(valGen)
}

func NewFilterListFromFilterExpression(value *FilterExpression) *api.FilterList {
	return api.NewFilterListFromFilterExpression((*api.FilterExpression)(value))
}

func NewFilterListFromFilterExpressionList(value []*FilterExpression) *api.FilterList {
	var valGen []*api.FilterExpression
	for _, v := range value {
		valGen = append(valGen, (*api.FilterExpression)(v))
	}
	return api.NewFilterListFromFilterExpressionList(valGen)
}

func NewSearchBranchRequestTablesItemFromString(value string) *SearchBranchRequestTablesItem {
	return (*SearchBranchRequestTablesItem)(api.NewSearchBranchRequestTablesItemFromString(value))
}

type TargetExpressionItem api.TargetExpressionItem

type TargetExpression []*TargetExpressionItem

func NewTargetExpression(columnName string) *TargetExpressionItem {
	return (*TargetExpressionItem)(api.NewTargetExpressionItemFromString(columnName))
}

type TargetExpressionItemColumn api.TargetExpressionItemColumn

func NewTargetExpressionWithColumnObject(colObj TargetExpressionItemColumn) *TargetExpressionItem {
	colObjGen := &api.TargetExpressionItemColumn{
		Column: colObj.Column,
		Weight: colObj.Weight,
	}
	return (*TargetExpressionItem)(api.NewTargetExpressionItemFromTargetExpressionItemColumn(colObjGen))
}

type ValueBooster struct {
	Column          string
	Value           *api.ValueBoosterValue
	Factor          float64
	IfMatchesFilter *FilterExpression
}

func NewValueBoosterValueFromString(value string) *api.ValueBoosterValue {
	return api.NewValueBoosterValueFromString(value)
}

type BoosterExpressionValueBooster struct {
	ValueBooster *ValueBooster
}

func NewBoosterExpressionFromBoosterExpressionValueBooster(value *BoosterExpressionValueBooster) *BoosterExpression {
	genBoosterExpVal := &api.BoosterExpressionValueBooster{ValueBooster: &api.ValueBooster{
		Column:          value.ValueBooster.Column,
		Value:           value.ValueBooster.Value,
		Factor:          value.ValueBooster.Factor,
		IfMatchesFilter: (*api.FilterExpression)(value.ValueBooster.IfMatchesFilter),
	}}
	return (*BoosterExpression)(api.NewBoosterExpressionFromBoosterExpressionValueBooster(genBoosterExpVal))
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
	genValue := &api.BoosterExpressionNumericBooster{
		NumericBooster: &api.NumericBooster{
			Column:          value.NumericBooster.Column,
			Factor:          value.NumericBooster.Factor,
			Modifier:        (*api.NumericBoosterModifier)(value.NumericBooster.Modifier),
			IfMatchesFilter: (*api.FilterExpression)(value.NumericBooster.IfMatchesFilter),
		},
	}
	return (*BoosterExpression)(api.NewBoosterExpressionFromBoosterExpressionNumericBooster(genValue))
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
	genVal := &api.BoosterExpressionDateBooster{
		DateBooster: &api.DateBooster{
			Column:          value.DateBooster.Column,
			Origin:          value.DateBooster.Origin,
			Scale:           value.DateBooster.Scale,
			Decay:           value.DateBooster.Decay,
			Factor:          value.DateBooster.Factor,
			IfMatchesFilter: (*api.FilterExpression)(value.DateBooster.IfMatchesFilter),
		},
	}
	return (*BoosterExpression)(api.NewBoosterExpressionFromBoosterExpressionDateBooster(genVal))
}

type AskTableRequestSearch struct {
	Fuzziness *int
	Target    TargetExpression
	Prefix    *PrefixExpression
	Filter    *FilterExpression
	Boosters  []*BoosterExpression
}

type AskTableRequestSearchType uint8

const (
	AskTableRequestSearchTypeKeyword AskTableRequestSearchType = iota + 1
	AskTableRequestSearchTypeVector
)

type AskTableRequestVectorSearch struct {
	// The column to use for vector search. It must be of type `vector`.
	Column string
	// The column containing the text for vector search. Must be of type `text`.
	ContentColumn string
	Filter        *FilterExpression
}

type SummarizeTableRequestConsistency uint8

//func NewAggExpressionFromAggExpressionCount(value *AggExpressionCount) *api.AggExpression {
//
//}
