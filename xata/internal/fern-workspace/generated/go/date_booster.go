// This file was auto-generated by Fern from our API Definition.

package api

// Boost records based on the value of a datetime column. It is configured via "origin", "scale", and "decay". The further away from the "origin",
// the more the score is decayed. The decay function uses an exponential function. For example if origin is "now", and scale is 10 days and decay is 0.5, it
// should be interpreted as: a record with a date 10 days before/after origin will be boosted 2 times less than a record with the date at origin.
// The result of the exponential function is a boost between 0 and 1. The "factor" allows you to control how impactful this boost is, by multiplying it with a given value.
type DateBooster struct {
	// The column in which to look for the value.
	Column string `json:"column"`
	// The datetime (formatted as RFC3339) from where to apply the score decay function. The maximum boost will be applied for records with values at this time.
	// If it is not specified, the current date and time is used.
	Origin *string `json:"origin,omitempty"`
	// The duration at which distance from origin the score is decayed with factor, using an exponential function. It is formatted as number + units, for example: `5d`, `20m`, `10s`.
	Scale string `json:"scale"`
	// The decay factor to expect at "scale" distance from the "origin".
	Decay float64 `json:"decay"`
	// The factor with which to multiply the added boost.
	Factor          *float64          `json:"factor,omitempty"`
	IfMatchesFilter *FilterExpression `json:"ifMatchesFilter,omitempty"`
}
