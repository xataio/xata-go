// SPDX-License-Identifier: Apache-2.0

// This file was auto-generated by Fern from our API Definition.

package api

type GetWorkspacesListResponseWorkspacesItem struct {
	Id   WorkspaceId   `json:"id"`
	Name string        `json:"name"`
	Slug string        `json:"slug"`
	Role Role          `json:"role,omitempty"`
	Plan WorkspacePlan `json:"plan,omitempty"`
}
