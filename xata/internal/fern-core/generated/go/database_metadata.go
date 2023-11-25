// SPDX-License-Identifier: Apache-2.0

// This file was auto-generated by Fern from our API Definition.

package api

// Metadata of databases
type DatabaseMetadata struct {
	// The machine-readable name of a database
	Name string `json:"name"`
	// Region where this database is hosted
	Region string `json:"region"`
	// The time this database was created
	CreatedAt DateTime `json:"createdAt"`
	// Metadata about the database for display in Xata user interfaces
	Ui *DatabaseMetadataUi `json:"ui,omitempty"`
}
