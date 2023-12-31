// SPDX-License-Identifier: Apache-2.0

// This file was auto-generated by Fern from our API Definition.

package api

import (
	core "github.com/xataio/xata-go/xata/internal/fern-core/generated/go/core"
	http "net/http"
)

type Client interface {
	Users() UsersClient
	Authentication() AuthenticationClient
	OAuth() OAuthClient
	Workspaces() WorkspacesClient
	Invites() InvitesClient
	Databases() DatabasesClient
}

func NewClient(opts ...core.ClientOption) Client {
	options := core.NewClientOptions()
	for _, opt := range opts {
		opt(options)
	}
	return &client{
		baseURL:              options.BaseURL,
		httpClient:           options.HTTPClient,
		header:               options.ToHeader(),
		usersClient:          NewUsersClient(opts...),
		authenticationClient: NewAuthenticationClient(opts...),
		oAuthClient:          NewOAuthClient(opts...),
		workspacesClient:     NewWorkspacesClient(opts...),
		invitesClient:        NewInvitesClient(opts...),
		databasesClient:      NewDatabasesClient(opts...),
	}
}

type client struct {
	baseURL              string
	httpClient           core.HTTPClient
	header               http.Header
	usersClient          UsersClient
	authenticationClient AuthenticationClient
	oAuthClient          OAuthClient
	workspacesClient     WorkspacesClient
	invitesClient        InvitesClient
	databasesClient      DatabasesClient
}

func (c *client) Users() UsersClient {
	return c.usersClient
}

func (c *client) Authentication() AuthenticationClient {
	return c.authenticationClient
}

func (c *client) OAuth() OAuthClient {
	return c.oAuthClient
}

func (c *client) Workspaces() WorkspacesClient {
	return c.workspacesClient
}

func (c *client) Invites() InvitesClient {
	return c.invitesClient
}

func (c *client) Databases() DatabasesClient {
	return c.databasesClient
}
