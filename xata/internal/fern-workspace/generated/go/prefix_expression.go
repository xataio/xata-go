// SPDX-License-Identifier: Apache-2.0

// This file was auto-generated by Fern from our API Definition.

package api

import (
	json "encoding/json"
	fmt "fmt"
	strconv "strconv"
)

// If the prefix type is set to "disabled" (the default), the search only matches full words. If the prefix type is set to "phrase", the search will return results that match prefixes of the search phrase.
type PrefixExpression uint8

const (
	PrefixExpressionPhrase PrefixExpression = iota + 1
	PrefixExpressionDisabled
)

func (p PrefixExpression) String() string {
	switch p {
	default:
		return strconv.Itoa(int(p))
	case PrefixExpressionPhrase:
		return "phrase"
	case PrefixExpressionDisabled:
		return "disabled"
	}
}

func (p PrefixExpression) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", p.String())), nil
}

func (p *PrefixExpression) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	switch raw {
	case "phrase":
		value := PrefixExpressionPhrase
		*p = value
	case "disabled":
		value := PrefixExpressionDisabled
		*p = value
	}
	return nil
}
