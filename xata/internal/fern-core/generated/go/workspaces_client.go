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

type WorkspacesClient interface {
	GetWorkspacesList(ctx context.Context) (*GetWorkspacesListResponse, error)
	CreateWorkspace(ctx context.Context, request *WorkspaceMeta) (*Workspace, error)
	GetWorkspace(ctx context.Context, workspaceId WorkspaceId) (*Workspace, error)
	UpdateWorkspace(ctx context.Context, workspaceId WorkspaceId, request *WorkspaceMeta) (*Workspace, error)
	DeleteWorkspace(ctx context.Context, workspaceId WorkspaceId) error
	GetWorkspaceMembersList(ctx context.Context, workspaceId WorkspaceId) (*WorkspaceMembers, error)
	UpdateWorkspaceMemberRole(ctx context.Context, workspaceId WorkspaceId, userId UserId, request *UpdateWorkspaceMemberRoleRequest) error
	RemoveWorkspaceMember(ctx context.Context, workspaceId WorkspaceId, userId UserId) error
}

func NewWorkspacesClient(opts ...core.ClientOption) WorkspacesClient {
	options := core.NewClientOptions()
	for _, opt := range opts {
		opt(options)
	}
	return &workspacesClient{
		baseURL:    options.BaseURL,
		httpClient: options.HTTPClient,
		header:     options.ToHeader(),
	}
}

type workspacesClient struct {
	baseURL    string
	httpClient core.HTTPClient
	header     http.Header
}

// Retrieve the list of workspaces the user belongs to
func (w *workspacesClient) GetWorkspacesList(ctx context.Context) (*GetWorkspacesListResponse, error) {
	baseURL := "/"
	if w.baseURL != "" {
		baseURL = w.baseURL
	}
	endpointURL := baseURL + "/" + "workspaces"

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

	var response *GetWorkspacesListResponse
	if err := core.DoRequest(
		ctx,
		w.httpClient,
		endpointURL,
		http.MethodGet,
		nil,
		&response,
		false,
		w.header,
		errorDecoder,
	); err != nil {
		return response, err
	}
	return response, nil
}

// Creates a new workspace with the user requesting it as its single owner.
func (w *workspacesClient) CreateWorkspace(ctx context.Context, request *WorkspaceMeta) (*Workspace, error) {
	baseURL := "/"
	if w.baseURL != "" {
		baseURL = w.baseURL
	}
	endpointURL := baseURL + "/" + "workspaces"

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

	var response *Workspace
	if err := core.DoRequest(
		ctx,
		w.httpClient,
		endpointURL,
		http.MethodPost,
		request,
		&response,
		false,
		w.header,
		errorDecoder,
	); err != nil {
		return response, err
	}
	return response, nil
}

// Retrieve workspace info from a workspace ID
//
// Workspace ID
func (w *workspacesClient) GetWorkspace(ctx context.Context, workspaceId WorkspaceId) (*Workspace, error) {
	baseURL := "/"
	if w.baseURL != "" {
		baseURL = w.baseURL
	}
	endpointURL := fmt.Sprintf(baseURL+"/"+"workspaces/%v", workspaceId)

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
		case 403:
			value := new(ForbiddenError)
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

	var response *Workspace
	if err := core.DoRequest(
		ctx,
		w.httpClient,
		endpointURL,
		http.MethodGet,
		nil,
		&response,
		false,
		w.header,
		errorDecoder,
	); err != nil {
		return response, err
	}
	return response, nil
}

// Update workspace info
//
// Workspace ID
func (w *workspacesClient) UpdateWorkspace(ctx context.Context, workspaceId WorkspaceId, request *WorkspaceMeta) (*Workspace, error) {
	baseURL := "/"
	if w.baseURL != "" {
		baseURL = w.baseURL
	}
	endpointURL := fmt.Sprintf(baseURL+"/"+"workspaces/%v", workspaceId)

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
		case 403:
			value := new(ForbiddenError)
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

	var response *Workspace
	if err := core.DoRequest(
		ctx,
		w.httpClient,
		endpointURL,
		http.MethodPut,
		request,
		&response,
		false,
		w.header,
		errorDecoder,
	); err != nil {
		return response, err
	}
	return response, nil
}

// Delete the workspace with the provided ID
//
// Workspace ID
func (w *workspacesClient) DeleteWorkspace(ctx context.Context, workspaceId WorkspaceId) error {
	baseURL := "/"
	if w.baseURL != "" {
		baseURL = w.baseURL
	}
	endpointURL := fmt.Sprintf(baseURL+"/"+"workspaces/%v", workspaceId)

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
		case 403:
			value := new(ForbiddenError)
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
		w.httpClient,
		endpointURL,
		http.MethodDelete,
		nil,
		nil,
		false,
		w.header,
		errorDecoder,
	); err != nil {
		return err
	}
	return nil
}

// Retrieve the list of members of the given workspace
//
// Workspace ID
func (w *workspacesClient) GetWorkspaceMembersList(ctx context.Context, workspaceId WorkspaceId) (*WorkspaceMembers, error) {
	baseURL := "/"
	if w.baseURL != "" {
		baseURL = w.baseURL
	}
	endpointURL := fmt.Sprintf(baseURL+"/"+"workspaces/%v/members", workspaceId)

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
		case 403:
			value := new(ForbiddenError)
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

	var response *WorkspaceMembers
	if err := core.DoRequest(
		ctx,
		w.httpClient,
		endpointURL,
		http.MethodGet,
		nil,
		&response,
		false,
		w.header,
		errorDecoder,
	); err != nil {
		return response, err
	}
	return response, nil
}

// Update a workspace member role. Workspaces must always have at least one owner, so this operation will fail if trying to remove owner role from the last owner in the workspace.
//
// Workspace ID
// UserID
func (w *workspacesClient) UpdateWorkspaceMemberRole(ctx context.Context, workspaceId WorkspaceId, userId UserId, request *UpdateWorkspaceMemberRoleRequest) error {
	baseURL := "/"
	if w.baseURL != "" {
		baseURL = w.baseURL
	}
	endpointURL := fmt.Sprintf(baseURL+"/"+"workspaces/%v/members/%v", workspaceId, userId)

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
		case 403:
			value := new(ForbiddenError)
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
		w.httpClient,
		endpointURL,
		http.MethodPut,
		request,
		nil,
		false,
		w.header,
		errorDecoder,
	); err != nil {
		return err
	}
	return nil
}

// Remove the member from the workspace
//
// Workspace ID
// UserID
func (w *workspacesClient) RemoveWorkspaceMember(ctx context.Context, workspaceId WorkspaceId, userId UserId) error {
	baseURL := "/"
	if w.baseURL != "" {
		baseURL = w.baseURL
	}
	endpointURL := fmt.Sprintf(baseURL+"/"+"workspaces/%v/members/%v", workspaceId, userId)

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
		case 403:
			value := new(ForbiddenError)
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
		w.httpClient,
		endpointURL,
		http.MethodDelete,
		nil,
		nil,
		false,
		w.header,
		errorDecoder,
	); err != nil {
		return err
	}
	return nil
}
