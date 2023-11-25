// SPDX-License-Identifier: Apache-2.0

// This file was auto-generated by Fern from our API Definition.

package api

// SearchBranchRequest is an in-lined request used by the SearchBranch endpoint.
type SearchBranchRequest struct {
	Fuzziness *FuzzinessExpression `json:"fuzziness,omitempty"`
	Highlight *HighlightExpression `json:"highlight,omitempty"`
	Page      *SearchPageConfig    `json:"page,omitempty"`
	Prefix    *PrefixExpression    `json:"prefix,omitempty"`
	// The query string.
	Query string `json:"query"`
	// An array with the tables in which to search. By default, all tables are included. Optionally, filters can be included that apply to each table.
	Tables *[]*SearchBranchRequestTablesItem `json:"tables,omitempty"`
}
