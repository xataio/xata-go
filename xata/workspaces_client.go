package xata

import (
	"context"

	xatagencore "github.com/xataio/xata-go/xata/internal/fern-core/generated/go"
	xatagenclient "github.com/xataio/xata-go/xata/internal/fern-core/generated/go/core"
)

type WorkspaceMeta xatagencore.WorkspaceMeta

type WorkspacesClient interface {
	List(ctx context.Context) (*xatagencore.GetWorkspacesListResponse, error)
	Create(ctx context.Context, request *WorkspaceMeta) (*xatagencore.Workspace, error)
	Delete(ctx context.Context, workspaceID string) error
}

type workspaceCli struct {
	generated xatagencore.WorkspacesClient
}

func (w workspaceCli) List(ctx context.Context) (*xatagencore.GetWorkspacesListResponse, error) {
	return w.generated.GetWorkspacesList(ctx)
}

func (w workspaceCli) Create(ctx context.Context, request *WorkspaceMeta) (*xatagencore.Workspace, error) {
	return w.generated.CreateWorkspace(ctx, &xatagencore.WorkspaceMeta{
		Name: request.Name,
		Slug: request.Slug,
	})
}

func (w workspaceCli) Delete(ctx context.Context, workspaceID string) error {
	return w.generated.DeleteWorkspace(ctx, workspaceID)
}

func NewWorkspacesClient(opts ...ClientOption) (WorkspacesClient, error) {
	cliOpts, err := consolidateClientOptionsForCore(opts...)
	if err != nil {
		return nil, err
	}

	return workspaceCli{generated: xatagencore.NewWorkspacesClient(
		func(options *xatagenclient.ClientOptions) {
			options.HTTPClient = cliOpts.HTTPClient
			options.BaseURL = cliOpts.BaseURL
			options.Bearer = cliOpts.Bearer
		})}, nil
}
