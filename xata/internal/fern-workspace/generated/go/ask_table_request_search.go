// This file was auto-generated by Fern from our API Definition.

package api

type AskTableRequestSearch struct {
	Boosters  *[]*BoosterExpression `json:"boosters,omitempty"`
	Filter    *FilterExpression     `json:"filter,omitempty"`
	Fuzziness *FuzzinessExpression  `json:"fuzziness,omitempty"`
	Prefix    *PrefixExpression     `json:"prefix,omitempty"`
	Target    *TargetExpression     `json:"target,omitempty"`
}
