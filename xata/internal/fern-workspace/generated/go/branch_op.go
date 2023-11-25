// SPDX-License-Identifier: Apache-2.0

// This file was auto-generated by Fern from our API Definition.

package api

type BranchOp struct {
	CreatedAt  DateTime        `json:"createdAt"`
	Id         string          `json:"id"`
	Message    *string         `json:"message,omitempty"`
	Migration  *Commit         `json:"migration,omitempty"`
	ModifiedAt *DateTime       `json:"modifiedAt,omitempty"`
	ParentId   *string         `json:"parentID,omitempty"`
	Status     MigrationStatus `json:"status,omitempty"`
	Title      *string         `json:"title,omitempty"`
}
