// This file was auto-generated by Fern from our API Definition.

package api

import (
	bytes "bytes"
	context "context"
	json "encoding/json"
	errors "errors"
	fmt "fmt"
	io "io"
	http "net/http"

	core "github.com/omerdemirok/xata-go/xata/internal/fern-core/generated/go/core"
)

type InvitesClient interface {
	InviteWorkspaceMember(ctx context.Context, workspaceId WorkspaceId, request *InviteWorkspaceMemberRequest) (*WorkspaceInvite, error)
	UpdateWorkspaceMemberInvite(ctx context.Context, workspaceId WorkspaceId, inviteId InviteId, request *UpdateWorkspaceMemberInviteRequest) (*WorkspaceInvite, error)
	CancelWorkspaceMemberInvite(ctx context.Context, workspaceId WorkspaceId, inviteId InviteId) error
	AcceptWorkspaceMemberInvite(ctx context.Context, workspaceId WorkspaceId, inviteKey InviteKey) error
	ResendWorkspaceMemberInvite(ctx context.Context, workspaceId WorkspaceId, inviteId InviteId) error
}

func NewInvitesClient(opts ...core.ClientOption) InvitesClient {
	options := core.NewClientOptions()
	for _, opt := range opts {
		opt(options)
	}
	return &invitesClient{
		baseURL:    options.BaseURL,
		httpClient: options.HTTPClient,
		header:     options.ToHeader(),
	}
}

type invitesClient struct {
	baseURL    string
	httpClient core.HTTPClient
	header     http.Header
}

// Invite some user to join the workspace with the given role
//
// Workspace ID
func (i *invitesClient) InviteWorkspaceMember(ctx context.Context, workspaceId WorkspaceId, request *InviteWorkspaceMemberRequest) (*WorkspaceInvite, error) {
	baseURL := "/"
	if i.baseURL != "" {
		baseURL = i.baseURL
	}
	endpointURL := fmt.Sprintf(baseURL+"/"+"workspaces/%v/invites", workspaceId)

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

	var response *WorkspaceInvite
	if err := core.DoRequest(
		ctx,
		i.httpClient,
		endpointURL,
		http.MethodPost,
		request,
		&response,
		false,
		i.header,
		errorDecoder,
	); err != nil {
		return response, err
	}
	return response, nil
}

// This operation provides a way to update an existing invite. Updates are performed in-place; they do not change the invite link, the expiry time, nor do they re-notify the recipient of the invite.
//
// Workspace ID
// Invite identifier
func (i *invitesClient) UpdateWorkspaceMemberInvite(ctx context.Context, workspaceId WorkspaceId, inviteId InviteId, request *UpdateWorkspaceMemberInviteRequest) (*WorkspaceInvite, error) {
	baseURL := "/"
	if i.baseURL != "" {
		baseURL = i.baseURL
	}
	endpointURL := fmt.Sprintf(baseURL+"/"+"workspaces/%v/invites/%v", workspaceId, inviteId)

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
		case 422:
			value := new(UnprocessableEntityError)
			value.APIError = apiError
			if err := decoder.Decode(value); err != nil {
				return err
			}
			return value
		}
		return apiError
	}

	var response *WorkspaceInvite
	if err := core.DoRequest(
		ctx,
		i.httpClient,
		endpointURL,
		http.MethodPatch,
		request,
		&response,
		false,
		i.header,
		errorDecoder,
	); err != nil {
		return response, err
	}
	return response, nil
}

// This operation provides a way to cancel invites by deleting them. Already accepted invites cannot be deleted.
//
// Workspace ID
// Invite identifier
func (i *invitesClient) CancelWorkspaceMemberInvite(ctx context.Context, workspaceId WorkspaceId, inviteId InviteId) error {
	baseURL := "/"
	if i.baseURL != "" {
		baseURL = i.baseURL
	}
	endpointURL := fmt.Sprintf(baseURL+"/"+"workspaces/%v/invites/%v", workspaceId, inviteId)

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
		i.httpClient,
		endpointURL,
		http.MethodDelete,
		nil,
		nil,
		false,
		i.header,
		errorDecoder,
	); err != nil {
		return err
	}
	return nil
}

// Accept the invitation to join a workspace. If the operation succeeds the user will be a member of the workspace
//
// Workspace ID
// Invite Key (secret) for the invited user
func (i *invitesClient) AcceptWorkspaceMemberInvite(ctx context.Context, workspaceId WorkspaceId, inviteKey InviteKey) error {
	baseURL := "/"
	if i.baseURL != "" {
		baseURL = i.baseURL
	}
	endpointURL := fmt.Sprintf(baseURL+"/"+"workspaces/%v/invites/%v/accept", workspaceId, inviteKey)

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
		i.httpClient,
		endpointURL,
		http.MethodPost,
		nil,
		nil,
		false,
		i.header,
		errorDecoder,
	); err != nil {
		return err
	}
	return nil
}

// This operation provides a way to resend an Invite notification. Invite notifications can only be sent for Invites not yet accepted.
//
// Workspace ID
// Invite identifier
func (i *invitesClient) ResendWorkspaceMemberInvite(ctx context.Context, workspaceId WorkspaceId, inviteId InviteId) error {
	baseURL := "/"
	if i.baseURL != "" {
		baseURL = i.baseURL
	}
	endpointURL := fmt.Sprintf(baseURL+"/"+"workspaces/%v/invites/%v/resend", workspaceId, inviteId)

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
		i.httpClient,
		endpointURL,
		http.MethodPost,
		nil,
		nil,
		false,
		i.header,
		errorDecoder,
	); err != nil {
		return err
	}
	return nil
}
