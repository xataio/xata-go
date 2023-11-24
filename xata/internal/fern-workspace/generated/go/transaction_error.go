// This file was auto-generated by Fern from our API Definition.

package api

// An error message from a failing transaction operation
type TransactionError struct {
	// The index of the failing operation
	Index int `json:"index"`
	// The error message
	Message string `json:"message"`
}
