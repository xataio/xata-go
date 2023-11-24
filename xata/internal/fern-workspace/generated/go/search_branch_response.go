// This file was auto-generated by Fern from our API Definition.

package api

type SearchBranchResponse struct {
	Records []*Record `json:"records,omitempty"`
	// The total count of records matched. It will be accurately returned up to 10000 records.
	TotalCount int     `json:"totalCount"`
	Warning    *string `json:"warning,omitempty"`
}
