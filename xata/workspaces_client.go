package xata

import (
	"context"
	"log"

	xatagencore "github.com/xataio/xata-go/xata/internal/fern-core/generated/go"
	xatagenclient "github.com/xataio/xata-go/xata/internal/fern-core/generated/go/core"
)

type WorkspaceMeta xatagencore.WorkspaceMeta

type WorkspacesClient interface {
	List(ctx context.Context) (*xatagencore.GetWorkspacesListResponse, error)
	Create(ctx context.Context, request *WorkspaceMeta) (*xatagencore.Workspace, error)
	Delete(ctx context.Context, workspaceID string) error
	Get(ctx context.Context) (*xatagencore.Workspace, error)
	GetWithWorkspaceID(ctx context.Context, workspaceID string) (*xatagencore.Workspace, error)
}

type workspaceCli struct {
	generated   xatagencore.WorkspacesClient
	workspaceID string
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

func (w workspaceCli) Get(ctx context.Context) (*xatagencore.Workspace, error) {
	return w.generated.GetWorkspace(ctx, w.workspaceID)
}

func (w workspaceCli) GetWithWorkspaceID(ctx context.Context, workspaceID string) (*xatagencore.Workspace, error) {
	return w.generated.GetWorkspace(ctx, workspaceID)
}

func NewWorkspacesClient(opts ...ClientOption) (WorkspacesClient, error) {
	cliOpts, err := consolidateClientOptionsForCore(opts...)
	if err != nil {
		return nil, err
	}

	dbCfg, err := loadDatabaseConfig()
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
