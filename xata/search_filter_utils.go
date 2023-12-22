// SPDX-License-Identifier: Apache-2.0

package xata

import (
	xatagenworkspace "github.com/xataio/xata-go/xata/internal/fern-workspace/generated/go"
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
	ConsistencyStrong QueryTableRequestConsistency = iota + 1
	ConsistencyEventual
)

type PrefixExpression uint8

const (
	PrefixExpressionPhrase PrefixExpression = iota + 1
	PrefixExpressionDisabled
)

type SearchBranchRequestTablesItem xatagenworkspace.SearchBranchRequestTablesItem

type HighlightExpression xatagenworkspace.HighlightExpression

type SearchPageConfig xatagenworkspace.SearchPageConfig

type BoosterExpression xatagenworkspace.BoosterExpression

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

func NewSearchBranchRequestTablesItemFromString(value string) *SearchBranchRequestTablesItem {
	return (*SearchBranchRequestTablesItem)(xatagenworkspace.NewSearchBranchRequestTablesItemFromString(value))
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
	DateBooster *DateBooster
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

type AggExpression *xatagenworkspace.AggExpression

type AggExpressionMap = map[string]AggExpression

type AggExpressionCount xatagenworkspace.AggExpressionCount

type CountAggFilter struct {
	Filter FilterExpression
}

func CountByFilter(value CountAggFilter) *xatagenworkspace.CountAgg {
	return xatagenworkspace.NewCountAggFromCountAggFilter(&xatagenworkspace.CountAggFilter{
		Filter: (*xatagenworkspace.FilterExpression)(&value.Filter),
	})
}

func CountAll() *xatagenworkspace.CountAgg {
	return xatagenworkspace.NewCountAggWithStringLiteral()
}

func NewCountAggExpression(value AggExpressionCount) AggExpression {
	return xatagenworkspace.NewAggExpressionFromAggExpressionCount((*xatagenworkspace.AggExpressionCount)(&value))
}

func NewSumAggExpression(column string) *xatagenworkspace.AggExpression {
	return xatagenworkspace.NewAggExpressionFromAggExpressionSum(&xatagenworkspace.AggExpressionSum{
		Sum: &xatagenworkspace.SumAgg{Column: column},
	})
}

func NewMaxAggExpression(column string) *xatagenworkspace.AggExpression {
	return xatagenworkspace.NewAggExpressionFromAggExpressionMax(&xatagenworkspace.AggExpressionMax{
		Max: &xatagenworkspace.MaxAgg{Column: column},
	})
}

func NewMinAggExpression(column string) *xatagenworkspace.AggExpression {
	return xatagenworkspace.NewAggExpressionFromAggExpressionMin(&xatagenworkspace.AggExpressionMin{
		Min: &xatagenworkspace.MinAgg{Column: column},
	})
}

func NewAverageAggExpression(column string) *xatagenworkspace.AggExpression {
	return xatagenworkspace.NewAggExpressionFromAggExpressionAverage(&xatagenworkspace.AggExpressionAverage{
		Average: &xatagenworkspace.AverageAgg{Column: column},
	})
}

type UniqueCountAgg struct {
	Column             string
	PrecisionThreshold *int
}

func NewUniqueCountAggExpression(value UniqueCountAgg) *xatagenworkspace.AggExpression {
	return xatagenworkspace.NewAggExpressionFromAggExpressionUniqueCount(&xatagenworkspace.AggExpressionUniqueCount{
		UniqueCount: &xatagenworkspace.UniqueCountAgg{
			Column:             value.Column,
			PrecisionThreshold: value.PrecisionThreshold,
		},
	})
}

type DateHistogramAggCalendarInterval uint8

const (
	DateHistogramAggCalendarIntervalMinute DateHistogramAggCalendarInterval = iota + 1
	DateHistogramAggCalendarIntervalHour
	DateHistogramAggCalendarIntervalDay
	DateHistogramAggCalendarIntervalWeek
	DateHistogramAggCalendarIntervalMonth
	DateHistogramAggCalendarIntervalQuarter
	DateHistogramAggCalendarIntervalYear
)

// Split data into buckets by a datetime column. Accepts sub-aggregations for each bucket.
type DateHistogramAgg struct {
	// The column to use for bucketing. Must be of type datetime.
	Column string
	// The fixed interval to use when bucketing.
	// It is formatted as number + units, for example: `5d`, `20m`, `10s`.
	Interval *string
	// The calendar-aware interval to use when bucketing. Possible values are: `minute`,
	// `hour`, `day`, `week`, `month`, `quarter`, `year`.
	CalendarInterval *DateHistogramAggCalendarInterval `json:"calendarInterval,omitempty"`
	// The timezone to use for bucketing. By default, UTC is assumed.
	// The accepted format is as an ISO 8601 UTC offset. For example: `+01:00` or
	// `-08:00`.
	Timezone *string
}

func NewDateHistogramAggExpression(value DateHistogramAgg) *xatagenworkspace.AggExpression {
	return xatagenworkspace.NewAggExpressionFromAggExpressionDateHistogram(&xatagenworkspace.AggExpressionDateHistogram{
		DateHistogram: &xatagenworkspace.DateHistogramAgg{
			Column:           value.Column,
			Interval:         value.Interval,
			CalendarInterval: (*xatagenworkspace.DateHistogramAggCalendarInterval)(value.CalendarInterval),
			Timezone:         value.Timezone,
		},
	})
}

// Split data into buckets by the unique values in a column. Accepts sub-aggregations for each bucket.
// The top values as ordered by the number of records (`$count`) are returned.
type TopValuesAgg struct {
	// The column to use for bucketing. Accepted types are `string`, `email`, `int`, `float`, or `bool`.
	Column string
	// The maximum number of unique values to return.
	Size *int
	// Sub Aggregations
	Aggs *xatagenworkspace.AggExpressionMap `json:"aggs,omitempty"`
}

func NewTopValuesAggExpression(value TopValuesAgg) *xatagenworkspace.AggExpression {
	return xatagenworkspace.NewAggExpressionFromAggExpressionTopValues(&xatagenworkspace.AggExpressionTopValues{
		TopValues: &xatagenworkspace.TopValuesAgg{
			Column: value.Column,
			Size:   value.Size,
			Aggs:   (*xatagenworkspace.AggExpressionMap)(value.Aggs),
		},
	})
}

// Split data into buckets by dynamic numeric ranges. Accepts sub-aggregations for each bucket.
type NumericHistogramAgg struct {
	// The column to use for bucketing. Must be of numeric type.
	Column string
	// The numeric interval to use for bucketing. The resulting buckets will be ranges
	// with this value as size.
	Interval float64
	// By default the bucket keys start with 0 and then continue in `interval` steps. The bucket
	// boundaries can be shifted by using the offset option. For example, if the `interval` is 100,
	// but you prefer the bucket boundaries to be `[50, 150), [150, 250), etc.`, you can set `offset`
	// to 50.
	Offset *float64
}

func NewNumericHistogramAggExpression(value NumericHistogramAgg) *xatagenworkspace.AggExpression {
	return xatagenworkspace.NewAggExpressionFromAggExpressionNumericHistogram(&xatagenworkspace.AggExpressionNumericHistogram{NumericHistogram: &xatagenworkspace.NumericHistogramAgg{
		Column:   value.Column,
		Offset:   value.Offset,
		Interval: value.Interval,
	}})
}
