// SPDX-License-Identifier: Apache-2.0

// This file was auto-generated by Fern from our API Definition.

package api

import (
	json "encoding/json"
	fmt "fmt"
)

// Branch schema migration operations.
type MigrationOp struct {
	typeName          string
	MigrationTableOp  *MigrationTableOp
	MigrationColumnOp *MigrationColumnOp
}

func NewMigrationOpFromMigrationTableOp(value *MigrationTableOp) *MigrationOp {
	return &MigrationOp{typeName: "migrationTableOp", MigrationTableOp: value}
}

func NewMigrationOpFromMigrationColumnOp(value *MigrationColumnOp) *MigrationOp {
	return &MigrationOp{typeName: "migrationColumnOp", MigrationColumnOp: value}
}

func (m *MigrationOp) UnmarshalJSON(data []byte) error {
	valueMigrationTableOp := new(MigrationTableOp)
	if err := json.Unmarshal(data, &valueMigrationTableOp); err == nil {
		m.typeName = "migrationTableOp"
		m.MigrationTableOp = valueMigrationTableOp
		return nil
	}
	valueMigrationColumnOp := new(MigrationColumnOp)
	if err := json.Unmarshal(data, &valueMigrationColumnOp); err == nil {
		m.typeName = "migrationColumnOp"
		m.MigrationColumnOp = valueMigrationColumnOp
		return nil
	}
	return fmt.Errorf("%s cannot be deserialized as a %T", data, m)
}

func (m MigrationOp) MarshalJSON() ([]byte, error) {
	switch m.typeName {
	default:
		return nil, fmt.Errorf("invalid type %s in %T", m.typeName, m)
	case "migrationTableOp":
		return json.Marshal(m.MigrationTableOp)
	case "migrationColumnOp":
		return json.Marshal(m.MigrationColumnOp)
	}
}

type MigrationOpVisitor interface {
	VisitMigrationTableOp(*MigrationTableOp) error
	VisitMigrationColumnOp(*MigrationColumnOp) error
}

func (m *MigrationOp) Accept(v MigrationOpVisitor) error {
	switch m.typeName {
	default:
		return fmt.Errorf("invalid type %s in %T", m.typeName, m)
	case "migrationTableOp":
		return v.VisitMigrationTableOp(m.MigrationTableOp)
	case "migrationColumnOp":
		return v.VisitMigrationColumnOp(m.MigrationColumnOp)
	}
}
