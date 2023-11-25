// SPDX-License-Identifier: Apache-2.0

// This file was auto-generated by Fern from our API Definition.

package api

type SummarizeTableRequestPage struct {
	// The number of records returned by summarize. If the amount of data you have exceeds this, or you have
	// more complex reporting requirements, we recommend that you use the aggregate endpoint instead.
	Size *int `json:"size,omitempty"`
}
