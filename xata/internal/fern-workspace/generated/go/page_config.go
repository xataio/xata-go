// SPDX-License-Identifier: Apache-2.0

// This file was auto-generated by Fern from our API Definition.

package api

// Pagination settings.
type PageConfig struct {
	// Query the next page that follow the cursor.
	After *string `json:"after,omitempty"`
	// Query the previous page before the cursor.
	Before *string `json:"before,omitempty"`
	// Query the last page from the cursor.
	End *string `json:"end,omitempty"`
	// Use offset to skip entries. To skip pages set offset to a multiple of size.
	Offset *int `json:"offset,omitempty"`
	// Set page size. If the size is missing it is read from the cursor. If no cursor is given Xata will choose the default page size.
	Size *int `json:"size,omitempty"`
	// Query the first page from the cursor.
	Start *string `json:"start,omitempty"`
}
