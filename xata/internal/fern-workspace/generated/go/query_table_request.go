// SPDX-License-Identifier: Apache-2.0

// This file was auto-generated by Fern from our API Definition.

package api

// QueryTableRequest is an in-lined request used by the QueryTable endpoint.
type QueryTableRequest struct {
	Columns *QueryColumnsProjection `json:"columns,omitempty"`
	// The consistency level for this request.
	Consistency *QueryTableRequestConsistency `json:"consistency,omitempty"`
	Filter      *FilterExpression             `json:"filter,omitempty"`
	Page        *PageConfig                   `json:"page,omitempty"`
	Sort        *SortExpression               `json:"sort,omitempty"`
}
