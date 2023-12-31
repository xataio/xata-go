// SPDX-License-Identifier: Apache-2.0

// This file was auto-generated by Fern from our API Definition.

package api

type BranchMetadata struct {
	Branch     *BranchName `json:"branch,omitempty"`
	Labels     *[]string   `json:"labels,omitempty"`
	Repository *string     `json:"repository,omitempty"`
	Stage      *string     `json:"stage,omitempty"`
}
