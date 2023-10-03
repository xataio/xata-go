package xata

import (
	"net/http"
)

// httpClient is an interface for a subset of the *http.Client.
type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type ClientOptions struct {
	BaseURL    string
	HTTPClient httpClient
	HTTPHeader http.Header
	Bearer     string
}

type ClientOption func(*ClientOptions)

// WithAPIKey enables passing API key as a parameter.
// If not provided API key will be looked up in env vars.
func WithAPIKey(token string) func(options *ClientOptions) {
	return func(options *ClientOptions) {
		options.Bearer = token
	}
}

// WithHTTPClient enables passing an HTTP client as a parameter.
// If not provided, http.DefaultClient will be used.
func WithHTTPClient(client httpClient) func(options *ClientOptions) {
	return func(options *ClientOptions) {
		options.HTTPClient = client
	}
}

func WithBaseURL(baseURL string) func(options *ClientOptions) {
	return func(options *ClientOptions) {
		options.BaseURL = baseURL
	}
}
