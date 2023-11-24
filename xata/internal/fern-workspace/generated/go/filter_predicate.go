// This file was auto-generated by Fern from our API Definition.

package api

import (
	json "encoding/json"
	fmt "fmt"
)

type FilterPredicate struct {
	typeName               string
	FilterValue            *FilterValue
	FilterPredicateList    []*FilterPredicate
	FilterPredicateOp      *FilterPredicateOp
	FilterPredicateRangeOp *FilterPredicateRangeOp
}

func NewFilterPredicateFromFilterValue(value *FilterValue) *FilterPredicate {
	return &FilterPredicate{typeName: "filterValue", FilterValue: value}
}

func NewFilterPredicateFromFilterPredicateList(value []*FilterPredicate) *FilterPredicate {
	return &FilterPredicate{typeName: "filterPredicateList", FilterPredicateList: value}
}

func NewFilterPredicateFromFilterPredicateOp(value *FilterPredicateOp) *FilterPredicate {
	return &FilterPredicate{typeName: "filterPredicateOp", FilterPredicateOp: value}
}

func NewFilterPredicateFromFilterPredicateRangeOp(value *FilterPredicateRangeOp) *FilterPredicate {
	return &FilterPredicate{typeName: "filterPredicateRangeOp", FilterPredicateRangeOp: value}
}

func (f *FilterPredicate) UnmarshalJSON(data []byte) error {
	valueFilterValue := new(FilterValue)
	if err := json.Unmarshal(data, &valueFilterValue); err == nil {
		f.typeName = "filterValue"
		f.FilterValue = valueFilterValue
		return nil
	}
	var valueFilterPredicateList []*FilterPredicate
	if err := json.Unmarshal(data, &valueFilterPredicateList); err == nil {
		f.typeName = "filterPredicateList"
		f.FilterPredicateList = valueFilterPredicateList
		return nil
	}
	valueFilterPredicateOp := new(FilterPredicateOp)
	if err := json.Unmarshal(data, &valueFilterPredicateOp); err == nil {
		f.typeName = "filterPredicateOp"
		f.FilterPredicateOp = valueFilterPredicateOp
		return nil
	}
	valueFilterPredicateRangeOp := new(FilterPredicateRangeOp)
	if err := json.Unmarshal(data, &valueFilterPredicateRangeOp); err == nil {
		f.typeName = "filterPredicateRangeOp"
		f.FilterPredicateRangeOp = valueFilterPredicateRangeOp
		return nil
	}
	return fmt.Errorf("%s cannot be deserialized as a %T", data, f)
}

func (f FilterPredicate) MarshalJSON() ([]byte, error) {
	switch f.typeName {
	default:
		return nil, fmt.Errorf("invalid type %s in %T", f.typeName, f)
	case "filterValue":
		return json.Marshal(f.FilterValue)
	case "filterPredicateList":
		return json.Marshal(f.FilterPredicateList)
	case "filterPredicateOp":
		return json.Marshal(f.FilterPredicateOp)
	case "filterPredicateRangeOp":
		return json.Marshal(f.FilterPredicateRangeOp)
	}
}

type FilterPredicateVisitor interface {
	VisitFilterValue(*FilterValue) error
	VisitFilterPredicateList([]*FilterPredicate) error
	VisitFilterPredicateOp(*FilterPredicateOp) error
	VisitFilterPredicateRangeOp(*FilterPredicateRangeOp) error
}

func (f *FilterPredicate) Accept(v FilterPredicateVisitor) error {
	switch f.typeName {
	default:
		return fmt.Errorf("invalid type %s in %T", f.typeName, f)
	case "filterValue":
		return v.VisitFilterValue(f.FilterValue)
	case "filterPredicateList":
		return v.VisitFilterPredicateList(f.FilterPredicateList)
	case "filterPredicateOp":
		return v.VisitFilterPredicateOp(f.FilterPredicateOp)
	case "filterPredicateRangeOp":
		return v.VisitFilterPredicateRangeOp(f.FilterPredicateRangeOp)
	}
}
