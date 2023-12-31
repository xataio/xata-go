// SPDX-License-Identifier: Apache-2.0

// This file was auto-generated by Fern from our API Definition.

package api

type CompareBranchSchemasResponse struct {
	Edits  *SchemaEditScript `json:"edits,omitempty"`
	Source *Schema           `json:"source,omitempty"`
	Target *Schema           `json:"target,omitempty"`
}
