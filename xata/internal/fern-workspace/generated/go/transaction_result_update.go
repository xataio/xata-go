// SPDX-License-Identifier: Apache-2.0

// This file was auto-generated by Fern from our API Definition.

package api

import (
	json "encoding/json"
)

// A result from an update operation.
type TransactionResultUpdate struct {
	Columns *TransactionResultColumns `json:"columns,omitempty"`
	Id      RecordId                  `json:"id"`
	// The number of updated rows
	Rows      int `json:"rows"`
	operation string
}

func (t *TransactionResultUpdate) Operation() string {
	return t.operation
}

func (t *TransactionResultUpdate) UnmarshalJSON(data []byte) error {
	type unmarshaler TransactionResultUpdate
	var value unmarshaler
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*t = TransactionResultUpdate(value)
	t.operation = "update"
	return nil
}

func (t *TransactionResultUpdate) MarshalJSON() ([]byte, error) {
	type embed TransactionResultUpdate
	var marshaler = struct {
		embed
		Operation string `json:"operation"`
	}{
		embed:     embed(*t),
		Operation: "update",
	}
	return json.Marshal(marshaler)
}
