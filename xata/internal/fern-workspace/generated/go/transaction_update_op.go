// SPDX-License-Identifier: Apache-2.0

// This file was auto-generated by Fern from our API Definition.

package api

// Update operation
type TransactionUpdateOp struct {
	// If set, the call will return the requested fields as part of the response.
	Columns *[]string `json:"columns,omitempty"`
	// The fields of the record you'd like to update
	Fields map[string]any `json:"fields,omitempty"`
	Id     RecordId       `json:"id"`
	// The version of the record you expect to be updating
	IfVersion *int `json:"ifVersion,omitempty"`
	// The table name
	Table string `json:"table"`
	// Xata will insert this record if it cannot be found.
	Upsert *bool `json:"upsert,omitempty"`
}
