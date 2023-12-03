// SPDX-License-Identifier: Apache-2.0

package xata

import (
	"fmt"
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

func consolidateClientOptionsForCore(opts ...ClientOption) (*ClientOptions, error) {
	cliOpts := &ClientOptions{}

	for _, opt := range opts {
		opt(cliOpts)
	}

	if cliOpts.HTTPClient == nil {
		cliOpts.HTTPClient = http.DefaultClient
	}

	if cliOpts.BaseURL == "" {
		cliOpts.BaseURL = fmt.Sprintf("https://%s", defaultControlPlaneDomain)
	}

	if cliOpts.Bearer == "" {
		apiKey, err := getAPIKey()
		if err != nil {
			return nil, err
		}
		cliOpts.Bearer = apiKey
	}

	return cliOpts, nil
}

func consolidateClientOptionsForWorkspace(opts ...ClientOption) (*ClientOptions, *databaseConfig, error) {
	cliOpts := &ClientOptions{}

	for _, opt := range opts {
		opt(cliOpts)
	}

	if cliOpts.HTTPClient == nil {
		cliOpts.HTTPClient = http.DefaultClient
	}

	dbCfg, err := loadDatabaseConfig()
	if err != nil && cliOpts.BaseURL == "" {
		return nil, nil, err
	}

	if cliOpts.BaseURL == "" {
		cliOpts.BaseURL = fmt.Sprintf(
			"https://%s.%s.%s",
			dbCfg.workspaceID,
			dbCfg.region,
			dbCfg.domainWorkspace,
		)
	}

	if cliOpts.Bearer == "" {
		apiKey, err := getAPIKey()
		if err != nil {
			return nil, nil, err
		}
		cliOpts.Bearer = apiKey
	}

	return cliOpts, &dbCfg, nil
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

// WithBaseURL enables passing the base URL.
func WithBaseURL(baseURL string) func(options *ClientOptions) {
	return func(options *ClientOptions) {
		options.BaseURL = baseURL
	}
}
