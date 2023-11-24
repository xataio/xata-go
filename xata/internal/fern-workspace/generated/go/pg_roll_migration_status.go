// This file was auto-generated by Fern from our API Definition.

package api

import (
	json "encoding/json"
	fmt "fmt"
	strconv "strconv"
)

type PgRollMigrationStatus uint8

const (
	PgRollMigrationStatusNoMigrations PgRollMigrationStatus = iota + 1
	PgRollMigrationStatusInProgress
	PgRollMigrationStatusComplete
)

func (p PgRollMigrationStatus) String() string {
	switch p {
	default:
		return strconv.Itoa(int(p))
	case PgRollMigrationStatusNoMigrations:
		return "no migrations"
	case PgRollMigrationStatusInProgress:
		return "in progress"
	case PgRollMigrationStatusComplete:
		return "complete"
	}
}

func (p PgRollMigrationStatus) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", p.String())), nil
}

func (p *PgRollMigrationStatus) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	switch raw {
	case "no migrations":
		value := PgRollMigrationStatusNoMigrations
		*p = value
	case "in progress":
		value := PgRollMigrationStatusInProgress
		*p = value
	case "complete":
		value := PgRollMigrationStatusComplete
		*p = value
	}
	return nil
}
