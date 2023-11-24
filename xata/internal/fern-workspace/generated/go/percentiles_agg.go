// This file was auto-generated by Fern from our API Definition.

package api

// Calculate given percentiles of the numeric values in a particular column.
type PercentilesAgg struct {
	// The column on which to compute the average. Must be a numeric type.
	Column      string    `json:"column"`
	Percentiles []float64 `json:"percentiles,omitempty"`
}
