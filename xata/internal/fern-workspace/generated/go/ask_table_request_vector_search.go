// SPDX-License-Identifier: Apache-2.0

// This file was auto-generated by Fern from our API Definition.

package api

type AskTableRequestVectorSearch struct {
	// The column to use for vector search. It must be of type `vector`.
	Column string `json:"column"`
	// The column containing the text for vector search. Must be of type `text`.
	ContentColumn string            `json:"contentColumn"`
	Filter        *FilterExpression `json:"filter,omitempty"`
}
