// SPDX-License-Identifier: Apache-2.0

// This file was auto-generated by Fern from our API Definition.

package api

type TableMigration struct {
	ModifiedColumns *[]*ColumnMigration `json:"modifiedColumns,omitempty"`
	NewColumnOrder  []string            `json:"newColumnOrder,omitempty"`
	NewColumns      *map[string]*Column `json:"newColumns,omitempty"`
	RemovedColumns  *[]string           `json:"removedColumns,omitempty"`
}
