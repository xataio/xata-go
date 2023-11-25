// SPDX-License-Identifier: Apache-2.0

// This file was auto-generated by Fern from our API Definition.

package api

import (
	bytes "bytes"
	context "context"
	json "encoding/json"
	errors "errors"
	fmt "fmt"
	core "github.com/xataio/xata-go/xata/internal/fern-core/generated/go/core"
	io "io"
	http "net/http"
)

type AuthenticationClient interface {
	GetUserApiKeys(ctx context.Context) (*GetUserApiKeysResponse, error)
	CreateUserApiKey(ctx context.Context, keyName ApiKeyName) (*CreateUserApiKeyResponse, error)
	DeleteUserApiKey(ctx context.Context, keyName ApiKeyName) error
}

func NewAuthenticationClient(opts ...core.ClientOption) AuthenticationClient {
	options := core.NewClientOptions()
	for _, opt := range opts {
		opt(options)
	}
	return &authenticationClient{
		baseURL:    options.BaseURL,
		httpClient: options.HTTPClient,
		header:     options.ToHeader(),
	}
}

type authenticationClient struct {
	baseURL    string
	httpClient core.HTTPClient
	header     http.Header
}

// Retrieve a list of existing user API keys
func (a *authenticationClient) GetUserApiKeys(ctx context.Context) (*GetUserApiKeysResponse, error) {
	baseURL := "/"
	if a.baseURL != "" {
		baseURL = a.baseURL
	}
	endpointURL := baseURL + "/" + "user/keys"

	errorDecoder := func(statusCode int, body io.Reader) error {
		raw, err := io.ReadAll(body)
		if err != nil {
			return err
		}
		apiError := core.NewAPIError(statusCode, errors.New(string(raw)))
		decoder := json.NewDecoder(bytes.NewReader(raw))
		switch statusCode {
		case 400:
			value := new(BadRequestError)
			value.APIError = apiError
			if err := decoder.Decode(value); err != nil {
				return err
			}
			return value
		case 401:
			value := new(UnauthorizedError)
			value.APIError = apiError
			if err := decoder.Decode(value); err != nil {
				return err
			}
			return value
		case 404:
			value := new(NotFoundError)
			value.APIError = apiError
			if err := decoder.Decode(value); err != nil {
				return err
			}
			return value
		}
		return apiError
	}

	var response *GetUserApiKeysResponse
	if err := core.DoRequest(
		ctx,
		a.httpClient,
		endpointURL,
		http.MethodGet,
		nil,
		&response,
		false,
		a.header,
		errorDecoder,
	); err != nil {
		return response, err
	}
	return response, nil
}

// Create and return new API key
//
// API Key name
func (a *authenticationClient) CreateUserApiKey(ctx context.Context, keyName ApiKeyName) (*CreateUserApiKeyResponse, error) {
	baseURL := "/"
	if a.baseURL != "" {
		baseURL = a.baseURL
	}
	endpointURL := fmt.Sprintf(baseURL+"/"+"user/keys/%v", keyName)

	errorDecoder := func(statusCode int, body io.Reader) error {
		raw, err := io.ReadAll(body)
		if err != nil {
			return err
		}
		apiError := core.NewAPIError(statusCode, errors.New(string(raw)))
		decoder := json.NewDecoder(bytes.NewReader(raw))
		switch statusCode {
		case 400:
			value := new(BadRequestError)
			value.APIError = apiError
			if err := decoder.Decode(value); err != nil {
				return err
			}
			return value
		case 401:
			value := new(UnauthorizedError)
			value.APIError = apiError
			if err := decoder.Decode(value); err != nil {
				return err
			}
			return value
		case 404:
			value := new(NotFoundError)
			value.APIError = apiError
			if err := decoder.Decode(value); err != nil {
				return err
			}
			return value
		}
		return apiError
	}

	var response *CreateUserApiKeyResponse
	if err := core.DoRequest(
		ctx,
		a.httpClient,
		endpointURL,
		http.MethodPost,
		nil,
		&response,
		false,
		a.header,
		errorDecoder,
	); err != nil {
		return response, err
	}
	return response, nil
}

// Delete an existing API key
//
// API Key name
func (a *authenticationClient) DeleteUserApiKey(ctx context.Context, keyName ApiKeyName) error {
	baseURL := "/"
	if a.baseURL != "" {
		baseURL = a.baseURL
	}
	endpointURL := fmt.Sprintf(baseURL+"/"+"user/keys/%v", keyName)

	errorDecoder := func(statusCode int, body io.Reader) error {
		raw, err := io.ReadAll(body)
		if err != nil {
			return err
		}
		apiError := core.NewAPIError(statusCode, errors.New(string(raw)))
		decoder := json.NewDecoder(bytes.NewReader(raw))
		switch statusCode {
		case 400:
			value := new(BadRequestError)
			value.APIError = apiError
			if err := decoder.Decode(value); err != nil {
				return err
			}
			return value
		case 401:
			value := new(UnauthorizedError)
			value.APIError = apiError
			if err := decoder.Decode(value); err != nil {
				return err
			}
			return value
		case 404:
			value := new(NotFoundError)
			value.APIError = apiError
			if err := decoder.Decode(value); err != nil {
				return err
			}
			return value
		}
		return apiError
	}

	if err := core.DoRequest(
		ctx,
		a.httpClient,
		endpointURL,
		http.MethodDelete,
		nil,
		nil,
		false,
		a.header,
		errorDecoder,
	); err != nil {
		return err
	}
	return nil
}
