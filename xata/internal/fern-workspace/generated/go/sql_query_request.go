// This file was auto-generated by Fern from our API Definition.

package api

// SqlQueryRequest is an in-lined request used by the Query endpoint.
type SqlQueryRequest struct {
	// The consistency level for this request.
	Consistency *SqlQueryRequestConsistency `json:"consistency,omitempty"`
	// The query parameter list.
	Params *[]any `json:"params,omitempty"`
	// The SQL statement.
	Statement string `json:"statement"`
}
