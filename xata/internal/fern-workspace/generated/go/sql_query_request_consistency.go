// This file was auto-generated by Fern from our API Definition.

package api

import (
	json "encoding/json"
	fmt "fmt"
	strconv "strconv"
)

// The consistency level for this request.
type SqlQueryRequestConsistency uint8

const (
	SqlQueryRequestConsistencyStrong SqlQueryRequestConsistency = iota + 1
	SqlQueryRequestConsistencyEventual
)

func (s SqlQueryRequestConsistency) String() string {
	switch s {
	default:
		return strconv.Itoa(int(s))
	case SqlQueryRequestConsistencyStrong:
		return "strong"
	case SqlQueryRequestConsistencyEventual:
		return "eventual"
	}
}

func (s SqlQueryRequestConsistency) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", s.String())), nil
}

func (s *SqlQueryRequestConsistency) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	switch raw {
	case "strong":
		value := SqlQueryRequestConsistencyStrong
		*s = value
	case "eventual":
		value := SqlQueryRequestConsistencyEventual
		*s = value
	}
	return nil
}
