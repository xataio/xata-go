// SPDX-License-Identifier: Apache-2.0

// This file was auto-generated by Fern from our API Definition.

package api

// Insert operation
type TransactionInsertOp struct {
	// If set, the call will return the requested fields as part of the response.
	Columns *[]string `json:"columns,omitempty"`
	// createOnly is used to change how Xata acts when an explicit ID is set in the `record` key.
	//
	// If `createOnly` is set to `true`, Xata will only attempt to insert the record. If there's a conflict, Xata
	// will cancel the transaction.
	//
	// If `createOnly` is set to `false`, Xata will attempt to insert the record. If there's no
	// conflict, the record is inserted. If there is a conflict, Xata will replace the record.
	CreateOnly *bool `json:"createOnly,omitempty"`
	// The version of the record you expect to be overwriting. Only valid with an
	// explicit ID is also set in the `record` key.
	IfVersion *int `json:"ifVersion,omitempty"`
	// The record to insert. The `id` field is optional; when specified, it will be used as the ID for the record.
	Record map[string]any `json:"record,omitempty"`
	// The table name
	Table string `json:"table"`
}
