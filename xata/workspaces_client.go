package xata

import (
	"context"
	"fmt"
	"log"
	"net/http"

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

func NewWorkspacesClient(opts ...ClientOption) WorkspacesClient {
	defaultOpts := &ClientOptions{
		BaseURL:    fmt.Sprintf("https://%s", defaultControlPlaneDomain),
		HTTPClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(defaultOpts)
	}

	if defaultOpts.Bearer == "" {
		apiKey, err := getAPIKey()
		if err != nil {
			log.Fatal(err)
		}
		defaultOpts.Bearer = apiKey
	}

	return workspaceCli{generated: xatagencore.NewWorkspacesClient(
		func(options *xatagenclient.ClientOptions) {
			options.HTTPClient = defaultOpts.HTTPClient
			options.BaseURL = defaultOpts.BaseURL
			options.Bearer = defaultOpts.Bearer
		})}
}
