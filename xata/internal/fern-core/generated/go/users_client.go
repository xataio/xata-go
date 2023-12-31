// SPDX-License-Identifier: Apache-2.0

// This file was auto-generated by Fern from our API Definition.

package api

import (
	bytes "bytes"
	context "context"
	json "encoding/json"
	errors "errors"
	core "github.com/xataio/xata-go/xata/internal/fern-core/generated/go/core"
	io "io"
	http "net/http"
)

type UsersClient interface {
	GetUser(ctx context.Context) (*UserWithId, error)
	UpdateUser(ctx context.Context, request *User) (*UserWithId, error)
	DeleteUser(ctx context.Context) error
}

func NewUsersClient(opts ...core.ClientOption) UsersClient {
	options := core.NewClientOptions()
	for _, opt := range opts {
		opt(options)
	}
	return &usersClient{
		baseURL:    options.BaseURL,
		httpClient: options.HTTPClient,
		header:     options.ToHeader(),
	}
}

type usersClient struct {
	baseURL    string
	httpClient core.HTTPClient
	header     http.Header
}

// Return details of the user making the request
func (u *usersClient) GetUser(ctx context.Context) (*UserWithId, error) {
	baseURL := "/"
	if u.baseURL != "" {
		baseURL = u.baseURL
	}
	endpointURL := baseURL + "/" + "user"

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

	var response *UserWithId
	if err := core.DoRequest(
		ctx,
		u.httpClient,
		endpointURL,
		http.MethodGet,
		nil,
		&response,
		false,
		u.header,
		errorDecoder,
	); err != nil {
		return response, err
	}
	return response, nil
}

// Update user info
func (u *usersClient) UpdateUser(ctx context.Context, request *User) (*UserWithId, error) {
	baseURL := "/"
	if u.baseURL != "" {
		baseURL = u.baseURL
	}
	endpointURL := baseURL + "/" + "user"

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

	var response *UserWithId
	if err := core.DoRequest(
		ctx,
		u.httpClient,
		endpointURL,
		http.MethodPut,
		request,
		&response,
		false,
		u.header,
		errorDecoder,
	); err != nil {
		return response, err
	}
	return response, nil
}

// Delete the user making the request
func (u *usersClient) DeleteUser(ctx context.Context) error {
	baseURL := "/"
	if u.baseURL != "" {
		baseURL = u.baseURL
	}
	endpointURL := baseURL + "/" + "user"

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
		u.httpClient,
		endpointURL,
		http.MethodDelete,
		nil,
		nil,
		false,
		u.header,
		errorDecoder,
	); err != nil {
		return err
	}
	return nil
}
