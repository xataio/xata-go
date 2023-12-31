// SPDX-License-Identifier: Apache-2.0

// This file was auto-generated by Fern from our API Definition.

package api

import (
	json "encoding/json"
	fmt "fmt"
)

type TargetExpressionItem struct {
	typeName                   string
	String                     string
	TargetExpressionItemColumn *TargetExpressionItemColumn
}

func NewTargetExpressionItemFromString(value string) *TargetExpressionItem {
	return &TargetExpressionItem{typeName: "string", String: value}
}

func NewTargetExpressionItemFromTargetExpressionItemColumn(value *TargetExpressionItemColumn) *TargetExpressionItem {
	return &TargetExpressionItem{typeName: "targetExpressionItemColumn", TargetExpressionItemColumn: value}
}

func (t *TargetExpressionItem) UnmarshalJSON(data []byte) error {
	var valueString string
	if err := json.Unmarshal(data, &valueString); err == nil {
		t.typeName = "string"
		t.String = valueString
		return nil
	}
	valueTargetExpressionItemColumn := new(TargetExpressionItemColumn)
	if err := json.Unmarshal(data, &valueTargetExpressionItemColumn); err == nil {
		t.typeName = "targetExpressionItemColumn"
		t.TargetExpressionItemColumn = valueTargetExpressionItemColumn
		return nil
	}
	return fmt.Errorf("%s cannot be deserialized as a %T", data, t)
}

func (t TargetExpressionItem) MarshalJSON() ([]byte, error) {
	switch t.typeName {
	default:
		return nil, fmt.Errorf("invalid type %s in %T", t.typeName, t)
	case "string":
		return json.Marshal(t.String)
	case "targetExpressionItemColumn":
		return json.Marshal(t.TargetExpressionItemColumn)
	}
}

type TargetExpressionItemVisitor interface {
	VisitString(string) error
	VisitTargetExpressionItemColumn(*TargetExpressionItemColumn) error
}

func (t *TargetExpressionItem) Accept(v TargetExpressionItemVisitor) error {
	switch t.typeName {
	default:
		return fmt.Errorf("invalid type %s in %T", t.typeName, t)
	case "string":
		return v.VisitString(t.String)
	case "targetExpressionItemColumn":
		return v.VisitTargetExpressionItemColumn(t.TargetExpressionItemColumn)
	}
}
