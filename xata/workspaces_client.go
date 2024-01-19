// SPDX-License-Identifier: Apache-2.0

package xata

import (
	"context"
	"log"

	xatagencore "github.com/xataio/xata-go/xata/internal/fern-core/generated/go"
	xatagenclient "github.com/xataio/xata-go/xata/internal/fern-core/generated/go/core"
)

type WorkspaceMeta xatagencore.WorkspaceMeta

type UpdateWorkspaceRequest struct {
	Payload     *WorkspaceMeta
	WorkspaceID *string
}

type WorkspacesClient interface {
	List(ctx context.Context) (*xatagencore.GetWorkspacesListResponse, error)
	Create(ctx context.Context, request *WorkspaceMeta) (*xatagencore.Workspace, error)
	Delete(ctx context.Context, workspaceID string) error
	Get(ctx context.Context) (*xatagencore.Workspace, error)
	GetWithWorkspaceID(ctx context.Context, workspaceID string) (*xatagencore.Workspace, error)
	Update(ctx context.Context, request UpdateWorkspaceRequest) (*xatagencore.Workspace, error)
}

type workspaceCli struct {
	generated   xatagencore.WorkspacesClient
	workspaceID string
}

// List retrieves the list of workspaces the user belongs to.
// https://xata.io/docs/api-reference/workspaces#get-list-of-workspaces
func (w workspaceCli) List(ctx context.Context) (*xatagencore.GetWorkspacesListResponse, error) {
	return w.generated.GetWorkspacesList(ctx)
}

// Create creates a new workspace with the user requesting it as its single owner.
// https://xata.io/docs/api-reference/workspaces#create-a-new-workspace
func (w workspaceCli) Create(ctx context.Context, request *WorkspaceMeta) (*xatagencore.Workspace, error) {
	return w.generated.CreateWorkspace(ctx, (*xatagencore.WorkspaceMeta)(request))
}

// Delete deletes the workspace with the provided ID.
// https://xata.io/docs/api-reference/workspaces/workspace_id#delete-an-existing-workspace
func (w workspaceCli) Delete(ctx context.Context, workspaceID string) error {
	return w.generated.DeleteWorkspace(ctx, workspaceID)
}

// Get retrieves workspace information for the default workspace.
// https://xata.io/docs/api-reference/workspaces/workspace_id#get-an-existing-workspace
func (w workspaceCli) Get(ctx context.Context) (*xatagencore.Workspace, error) {
	return w.generated.GetWorkspace(ctx, w.workspaceID)
}

// GetWithWorkspaceID retrieves workspace information for the given ID.
// https://xata.io/docs/api-reference/workspaces/workspace_id#get-an-existing-workspace
func (w workspaceCli) GetWithWorkspaceID(ctx context.Context, workspaceID string) (*xatagencore.Workspace, error) {
	return w.generated.GetWorkspace(ctx, workspaceID)
}

// Update updates workspace information.
// https://xata.io/docs/api-reference/workspaces/workspace_id#update-an-existing-workspace
func (w workspaceCli) Update(ctx context.Context, request UpdateWorkspaceRequest) (*xatagencore.Workspace, error) {
	workspaceID := w.workspaceID
	if request.WorkspaceID != nil && *request.WorkspaceID != "" {
		workspaceID = *request.WorkspaceID
	}

	return w.generated.UpdateWorkspace(ctx, workspaceID, (*xatagencore.WorkspaceMeta)(request.Payload))
}

// NewWorkspacesClient constructs a client for interacting with workspaces.
func NewWorkspacesClient(opts ...ClientOption) (WorkspacesClient, error) {
	cliOpts, err := consolidateClientOptionsForCore(opts...)
	if err != nil {
		return nil, err
	}

	dbCfg, err := loadDatabaseConfig(cliOpts)
	if err != nil {
		// No err, because the config values can be provided by the users.
		log.Println(err)
	}

	return workspaceCli{
		generated: xatagencore.NewWorkspacesClient(
			func(options *xatagenclient.ClientOptions) {
				options.HTTPClient = cliOpts.HTTPClient
				options.BaseURL = cliOpts.BaseURL
				options.Bearer = cliOpts.Bearer
			}),
		workspaceID: dbCfg.workspaceID,
	}, nil
}
