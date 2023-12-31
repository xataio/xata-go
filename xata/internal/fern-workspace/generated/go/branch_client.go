// SPDX-License-Identifier: Apache-2.0

// This file was auto-generated by Fern from our API Definition.

package api

import (
	bytes "bytes"
	context "context"
	json "encoding/json"
	errors "errors"
	fmt "fmt"
	core "github.com/xataio/xata-go/xata/internal/fern-workspace/generated/go/core"
	io "io"
	http "net/http"
	url "net/url"
)

type BranchClient interface {
	GetBranchDetails(ctx context.Context, dbBranchName DbBranchName) (*DbBranch, error)
	CreateBranch(ctx context.Context, dbBranchName DbBranchName, request *CreateBranchRequest) (*CreateBranchResponse, error)
	DeleteBranch(ctx context.Context, dbBranchName DbBranchName) (*DeleteBranchResponse, error)
	GetBranchMetadata(ctx context.Context, dbBranchName DbBranchName) (*BranchMetadata, error)
	UpdateBranchMetadata(ctx context.Context, dbBranchName DbBranchName, request *BranchMetadata) error
	ApplyMigration(ctx context.Context, dbBranchName DbBranchName, request []map[string]any) error
	PgRollStatus(ctx context.Context, dbBranchName DbBranchName) (*PgRollStatusResponse, error)
	GetBranchStats(ctx context.Context, dbBranchName DbBranchName) (*GetBranchStatsResponse, error)
	GetBranchList(ctx context.Context, dbName DbName) (*ListBranchesResponse, error)
	GetGitBranchesMapping(ctx context.Context, dbName DbName) (*ListGitBranchesResponse, error)
	AddGitBranchesEntry(ctx context.Context, dbName DbName, request *AddGitBranchesEntryRequest) (*AddGitBranchesEntryResponse, error)
	RemoveGitBranchesEntry(ctx context.Context, dbName DbName, request *RemoveGitBranchesEntryRequest) error
	ResolveBranch(ctx context.Context, dbName DbName, request *ResolveBranchRequest) (*ResolveBranchResponse, error)
}

func NewBranchClient(opts ...core.ClientOption) BranchClient {
	options := core.NewClientOptions()
	for _, opt := range opts {
		opt(options)
	}
	return &branchClient{
		baseURL:    options.BaseURL,
		httpClient: options.HTTPClient,
		header:     options.ToHeader(),
	}
}

type branchClient struct {
	baseURL    string
	httpClient core.HTTPClient
	header     http.Header
}

// The DBBranchName matches the pattern `{db_name}:{branch_name}`.
func (b *branchClient) GetBranchDetails(ctx context.Context, dbBranchName DbBranchName) (*DbBranch, error) {
	baseURL := "/"
	if b.baseURL != "" {
		baseURL = b.baseURL
	}
	endpointURL := fmt.Sprintf(baseURL+"/"+"db/%v", dbBranchName)

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

	var response *DbBranch
	if err := core.DoRequest(
		ctx,
		b.httpClient,
		endpointURL,
		http.MethodGet,
		nil,
		&response,
		false,
		b.header,
		errorDecoder,
	); err != nil {
		return response, err
	}
	return response, nil
}

// The DBBranchName matches the pattern `{db_name}:{branch_name}`.
func (b *branchClient) CreateBranch(ctx context.Context, dbBranchName DbBranchName, request *CreateBranchRequest) (*CreateBranchResponse, error) {
	baseURL := "/"
	if b.baseURL != "" {
		baseURL = b.baseURL
	}
	endpointURL := fmt.Sprintf(baseURL+"/"+"db/%v", dbBranchName)

	queryParams := make(url.Values)
	if request.From != nil {
		queryParams.Add("from", fmt.Sprintf("%v", *request.From))
	}
	if len(queryParams) > 0 {
		endpointURL += "?" + queryParams.Encode()
	}

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

	var response *CreateBranchResponse
	if err := core.DoRequest(
		ctx,
		b.httpClient,
		endpointURL,
		http.MethodPut,
		request,
		&response,
		false,
		b.header,
		errorDecoder,
	); err != nil {
		return response, err
	}
	return response, nil
}

// Delete the branch in the database and all its resources
//
// The DBBranchName matches the pattern `{db_name}:{branch_name}`.
func (b *branchClient) DeleteBranch(ctx context.Context, dbBranchName DbBranchName) (*DeleteBranchResponse, error) {
	baseURL := "/"
	if b.baseURL != "" {
		baseURL = b.baseURL
	}
	endpointURL := fmt.Sprintf(baseURL+"/"+"db/%v", dbBranchName)

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
		case 409:
			value := new(ConflictError)
			value.APIError = apiError
			if err := decoder.Decode(value); err != nil {
				return err
			}
			return value
		}
		return apiError
	}

	var response *DeleteBranchResponse
	if err := core.DoRequest(
		ctx,
		b.httpClient,
		endpointURL,
		http.MethodDelete,
		nil,
		&response,
		false,
		b.header,
		errorDecoder,
	); err != nil {
		return response, err
	}
	return response, nil
}

// The DBBranchName matches the pattern `{db_name}:{branch_name}`.
func (b *branchClient) GetBranchMetadata(ctx context.Context, dbBranchName DbBranchName) (*BranchMetadata, error) {
	baseURL := "/"
	if b.baseURL != "" {
		baseURL = b.baseURL
	}
	endpointURL := fmt.Sprintf(baseURL+"/"+"db/%v/metadata", dbBranchName)

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

	var response *BranchMetadata
	if err := core.DoRequest(
		ctx,
		b.httpClient,
		endpointURL,
		http.MethodGet,
		nil,
		&response,
		false,
		b.header,
		errorDecoder,
	); err != nil {
		return response, err
	}
	return response, nil
}

// Update the branch metadata
//
// The DBBranchName matches the pattern `{db_name}:{branch_name}`.
func (b *branchClient) UpdateBranchMetadata(ctx context.Context, dbBranchName DbBranchName, request *BranchMetadata) error {
	baseURL := "/"
	if b.baseURL != "" {
		baseURL = b.baseURL
	}
	endpointURL := fmt.Sprintf(baseURL+"/"+"db/%v/metadata", dbBranchName)

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
		b.httpClient,
		endpointURL,
		http.MethodPut,
		request,
		nil,
		false,
		b.header,
		errorDecoder,
	); err != nil {
		return err
	}
	return nil
}

// Applies a pgroll migration to the specified database.
//
// The DBBranchName matches the pattern `{db_name}:{branch_name}`.
func (b *branchClient) ApplyMigration(ctx context.Context, dbBranchName DbBranchName, request []map[string]any) error {
	baseURL := "/"
	if b.baseURL != "" {
		baseURL = b.baseURL
	}
	endpointURL := fmt.Sprintf(baseURL+"/"+"db/%v/pgroll/apply", dbBranchName)

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
		b.httpClient,
		endpointURL,
		http.MethodPost,
		request,
		nil,
		false,
		b.header,
		errorDecoder,
	); err != nil {
		return err
	}
	return nil
}

// The DBBranchName matches the pattern `{db_name}:{branch_name}`.
func (b *branchClient) PgRollStatus(ctx context.Context, dbBranchName DbBranchName) (*PgRollStatusResponse, error) {
	baseURL := "/"
	if b.baseURL != "" {
		baseURL = b.baseURL
	}
	endpointURL := fmt.Sprintf(baseURL+"/"+"db/%v/pgroll/status", dbBranchName)

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

	var response *PgRollStatusResponse
	if err := core.DoRequest(
		ctx,
		b.httpClient,
		endpointURL,
		http.MethodGet,
		nil,
		&response,
		false,
		b.header,
		errorDecoder,
	); err != nil {
		return response, err
	}
	return response, nil
}

// Get branch usage metrics.
//
// The DBBranchName matches the pattern `{db_name}:{branch_name}`.
func (b *branchClient) GetBranchStats(ctx context.Context, dbBranchName DbBranchName) (*GetBranchStatsResponse, error) {
	baseURL := "/"
	if b.baseURL != "" {
		baseURL = b.baseURL
	}
	endpointURL := fmt.Sprintf(baseURL+"/"+"db/%v/stats", dbBranchName)

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

	var response *GetBranchStatsResponse
	if err := core.DoRequest(
		ctx,
		b.httpClient,
		endpointURL,
		http.MethodGet,
		nil,
		&response,
		false,
		b.header,
		errorDecoder,
	); err != nil {
		return response, err
	}
	return response, nil
}

// List all available Branches
//
// The Database Name
func (b *branchClient) GetBranchList(ctx context.Context, dbName DbName) (*ListBranchesResponse, error) {
	baseURL := "/"
	if b.baseURL != "" {
		baseURL = b.baseURL
	}
	endpointURL := fmt.Sprintf(baseURL+"/"+"dbs/%v", dbName)

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

	var response *ListBranchesResponse
	if err := core.DoRequest(
		ctx,
		b.httpClient,
		endpointURL,
		http.MethodGet,
		nil,
		&response,
		false,
		b.header,
		errorDecoder,
	); err != nil {
		return response, err
	}
	return response, nil
}

// Lists all the git branches in the mapping, and their associated Xata branches.
//
// Example response:
//
// ```json
//
//	{
//	  "mappings": [
//	      {
//	        "gitBranch": "main",
//	        "xataBranch": "main"
//	      },
//	      {
//	        "gitBranch": "gitBranch1",
//	        "xataBranch": "xataBranch1"
//	      }
//	      {
//	        "gitBranch": "xataBranch2",
//	        "xataBranch": "xataBranch2"
//	      }
//	  ]
//	}
//
// ```
//
// The Database Name
func (b *branchClient) GetGitBranchesMapping(ctx context.Context, dbName DbName) (*ListGitBranchesResponse, error) {
	baseURL := "/"
	if b.baseURL != "" {
		baseURL = b.baseURL
	}
	endpointURL := fmt.Sprintf(baseURL+"/"+"dbs/%v/gitBranches", dbName)

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
		}
		return apiError
	}

	var response *ListGitBranchesResponse
	if err := core.DoRequest(
		ctx,
		b.httpClient,
		endpointURL,
		http.MethodGet,
		nil,
		&response,
		false,
		b.header,
		errorDecoder,
	); err != nil {
		return response, err
	}
	return response, nil
}

// Adds an entry to the mapping of git branches to Xata branches. The git branch and the Xata branch must be present in the body of the request. If the Xata branch doesn't exist, a 400 error is returned.
//
// If the git branch is already present in the mapping, the old entry is overwritten, and a warning message is included in the response. If the git branch is added and didn't exist before, the response code is 204. If the git branch existed and it was overwritten, the response code is 201.
//
// Example request:
//
// ```json
// // POST https://tutorial-ng7s8c.xata.sh/dbs/demo/gitBranches
//
//	{
//	  "gitBranch": "fix/bug123",
//	  "xataBranch": "fix_bug"
//	}
//
// ```
//
// The Database Name
func (b *branchClient) AddGitBranchesEntry(ctx context.Context, dbName DbName, request *AddGitBranchesEntryRequest) (*AddGitBranchesEntryResponse, error) {
	baseURL := "/"
	if b.baseURL != "" {
		baseURL = b.baseURL
	}
	endpointURL := fmt.Sprintf(baseURL+"/"+"dbs/%v/gitBranches", dbName)

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
		}
		return apiError
	}

	var response *AddGitBranchesEntryResponse
	if err := core.DoRequest(
		ctx,
		b.httpClient,
		endpointURL,
		http.MethodPost,
		request,
		&response,
		false,
		b.header,
		errorDecoder,
	); err != nil {
		return response, err
	}
	return response, nil
}

// Removes an entry from the mapping of git branches to Xata branches. The name of the git branch must be passed as a query parameter. If the git branch is not found, the endpoint returns a 404 status code.
//
// Example request:
//
// ```json
// // DELETE https://tutorial-ng7s8c.xata.sh/dbs/demo/gitBranches?gitBranch=fix%2Fbug123
// ```
//
// The Database Name
func (b *branchClient) RemoveGitBranchesEntry(ctx context.Context, dbName DbName, request *RemoveGitBranchesEntryRequest) error {
	baseURL := "/"
	if b.baseURL != "" {
		baseURL = b.baseURL
	}
	endpointURL := fmt.Sprintf(baseURL+"/"+"dbs/%v/gitBranches", dbName)

	queryParams := make(url.Values)
	queryParams.Add("gitBranch", fmt.Sprintf("%v", request.GitBranch))
	if len(queryParams) > 0 {
		endpointURL += "?" + queryParams.Encode()
	}

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
		b.httpClient,
		endpointURL,
		http.MethodDelete,
		request,
		nil,
		false,
		b.header,
		errorDecoder,
	); err != nil {
		return err
	}
	return nil
}

// In order to resolve the database branch, the following algorithm is used:
// * if the `gitBranch` was provided and is found in the [git branches mapping](/api-reference/dbs/db_name/gitBranches), the associated Xata branch is returned
// * else, if a Xata branch with the exact same name as `gitBranch` exists, return it
// * else, if `fallbackBranch` is provided and a branch with that name exists, return it
// * else, return the default branch of the DB (`main` or the first branch)
//
// Example call:
//
// ```json
// // GET https://tutorial-ng7s8c.xata.sh/dbs/demo/dbs/demo/resolveBranch?gitBranch=test&fallbackBranch=tsg
// ```
//
// Example response:
//
// ```json
//
//	{
//	  "branch": "main",
//	  "reason": {
//	    "code": "DEFAULT_BRANCH",
//	    "message": "Default branch for this database (main)"
//	  }
//	}
//
// ```
//
// The Database Name
func (b *branchClient) ResolveBranch(ctx context.Context, dbName DbName, request *ResolveBranchRequest) (*ResolveBranchResponse, error) {
	baseURL := "/"
	if b.baseURL != "" {
		baseURL = b.baseURL
	}
	endpointURL := fmt.Sprintf(baseURL+"/"+"dbs/%v/resolveBranch", dbName)

	queryParams := make(url.Values)
	if request.GitBranch != nil {
		queryParams.Add("gitBranch", fmt.Sprintf("%v", *request.GitBranch))
	}
	if request.FallbackBranch != nil {
		queryParams.Add("fallbackBranch", fmt.Sprintf("%v", *request.FallbackBranch))
	}
	if len(queryParams) > 0 {
		endpointURL += "?" + queryParams.Encode()
	}

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
		}
		return apiError
	}

	var response *ResolveBranchResponse
	if err := core.DoRequest(
		ctx,
		b.httpClient,
		endpointURL,
		http.MethodGet,
		request,
		&response,
		false,
		b.header,
		errorDecoder,
	); err != nil {
		return response, err
	}
	return response, nil
}
