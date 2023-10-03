// This file was auto-generated by Fern from our API Definition.

package api

import (
	json "encoding/json"
)

// InsertRecordRequest is an in-lined request used by the InsertRecord endpoint.
type InsertRecordRequest struct {
	// Column filters
	Columns []*string                        `json:"-"`
	Body    map[string]*DataInputRecordValue `json:"-"`
}

func (i *InsertRecordRequest) UnmarshalJSON(data []byte) error {
	var body map[string]*DataInputRecordValue
	if err := json.Unmarshal(data, &body); err != nil {
		return err
	}
	i.Body = body
	return nil
}

func (i *InsertRecordRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Body)
}
