// This file was auto-generated by Fern from our API Definition.

package api

import (
	json "encoding/json"
	fmt "fmt"
	strconv "strconv"
)

type MigrationStatus uint8

const (
	MigrationStatusCompleted MigrationStatus = iota + 1
	MigrationStatusPending
	MigrationStatusFailed
)

func (m MigrationStatus) String() string {
	switch m {
	default:
		return strconv.Itoa(int(m))
	case MigrationStatusCompleted:
		return "completed"
	case MigrationStatusPending:
		return "pending"
	case MigrationStatusFailed:
		return "failed"
	}
}

func (m MigrationStatus) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", m.String())), nil
}

func (m *MigrationStatus) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	switch raw {
	case "completed":
		value := MigrationStatusCompleted
		*m = value
	case "pending":
		value := MigrationStatusPending
		*m = value
	case "failed":
		value := MigrationStatusFailed
		*m = value
	}
	return nil
}
