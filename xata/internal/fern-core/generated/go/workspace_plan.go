// SPDX-License-Identifier: Apache-2.0

// This file was auto-generated by Fern from our API Definition.

package api

import (
	json "encoding/json"
	fmt "fmt"
	strconv "strconv"
)

type WorkspacePlan uint8

const (
	WorkspacePlanFree WorkspacePlan = iota + 1
	WorkspacePlanPro
)

func (w WorkspacePlan) String() string {
	switch w {
	default:
		return strconv.Itoa(int(w))
	case WorkspacePlanFree:
		return "free"
	case WorkspacePlanPro:
		return "pro"
	}
}

func (w WorkspacePlan) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", w.String())), nil
}

func (w *WorkspacePlan) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	switch raw {
	case "free":
		value := WorkspacePlanFree
		*w = value
	case "pro":
		value := WorkspacePlanPro
		*w = value
	}
	return nil
}
