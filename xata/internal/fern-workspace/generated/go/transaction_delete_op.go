// SPDX-License-Identifier: Apache-2.0

// This file was auto-generated by Fern from our API Definition.

package api

// A delete operation. The transaction will continue if no record matches the ID by default. To override this behaviour, set failIfMissing to true.
type TransactionDeleteOp struct {
	// If set, the call will return the requested fields as part of the response.
	Columns *[]string `json:"columns,omitempty"`
	// If true, the transaction will fail when the record doesn't exist.
	FailIfMissing *bool    `json:"failIfMissing,omitempty"`
	Id            RecordId `json:"id"`
	// The table name
	Table string `json:"table"`
}
