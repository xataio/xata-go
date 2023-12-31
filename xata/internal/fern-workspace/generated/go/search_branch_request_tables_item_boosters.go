// SPDX-License-Identifier: Apache-2.0

// This file was auto-generated by Fern from our API Definition.

package api

type SearchBranchRequestTablesItemBoosters struct {
	Boosters *[]*BoosterExpression `json:"boosters,omitempty"`
	Filter   *FilterExpression     `json:"filter,omitempty"`
	// The name of the table.
	Table  string            `json:"table"`
	Target *TargetExpression `json:"target,omitempty"`
}
