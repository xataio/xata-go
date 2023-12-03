// SPDX-License-Identifier: Apache-2.0

package xata

import (
	"context"
	"log"

	xatagencore "github.com/xataio/xata-go/xata/internal/fern-core/generated/go"
	xatagenclient "github.com/xataio/xata-go/xata/internal/fern-core/generated/go/core"
)

type UI xatagencore.CreateDatabaseRequestUi

type BranchMetadata xatagencore.BranchMetadata

type CreateDatabaseRequest struct {
	DatabaseName   string
	WorkspaceID    *string
	BranchName     *string
	Region         *string
	UI             *UI
	BranchMetaData *BranchMetadata
}

type DeleteDatabaseRequest struct {
	DatabaseName string
	WorkspaceID  *string
}

type RenameDatabaseRequest struct {
	DatabaseName string
	NewName      string
	WorkspaceID  *string
}

type DatabasesClient interface {
	Create(ctx context.Context, request CreateDatabaseRequest) (*xatagencore.CreateDatabaseResponse, error)
	Delete(ctx context.Context, request DeleteDatabaseRequest) (*xatagencore.DeleteDatabaseResponse, error)
	GetRegions(ctx context.Context) (*xatagencore.ListRegionsResponse, error)
	GetRegionsWithWorkspaceID(ctx context.Context, workspaceID string) (*xatagencore.ListRegionsResponse, error)
	List(ctx context.Context) (*xatagencore.ListDatabasesResponse, error)
	ListWithWorkspaceID(ctx context.Context, workspaceID string) (*xatagencore.ListDatabasesResponse, error)
	Rename(ctx context.Context, request RenameDatabaseRequest) (*xatagencore.DatabaseMetadata, error)
}

type databaseCli struct {
	generated   xatagencore.DatabasesClient
	WorkspaceID string
	BranchName  string
	Region      string
}

// Create creates a database.
// https://xata.io/docs/api-reference/workspaces/workspace_id/dbs/db_name#create-database
func (d databaseCli) Create(ctx context.Context, request CreateDatabaseRequest) (*xatagencore.CreateDatabaseResponse, error) {
	var workspaceID string
	if request.WorkspaceID == nil {
		workspaceID = d.WorkspaceID
	} else {
		workspaceID = *request.WorkspaceID
	}

	var branchName string
	if request.BranchName == nil {
		branchName = d.BranchName
	} else {
		branchName = *request.BranchName
	}

	var region string
	if request.Region == nil {
		region = d.Region
	} else {
		region = *request.Region
	}

	return d.generated.CreateDatabase(ctx, workspaceID, request.DatabaseName, &xatagencore.CreateDatabaseRequest{
		BranchName: String(branchName),
		Region:     region,
		Ui:         (*xatagencore.CreateDatabaseRequestUi)(request.UI),
		Metadata:   (*xatagencore.BranchMetadata)(request.BranchMetaData),
	})
}

// Delete deletes a database.
// https://xata.io/docs/api-reference/workspaces/workspace_id/dbs/db_name#delete-database
func (d databaseCli) Delete(ctx context.Context, request DeleteDatabaseRequest) (*xatagencore.DeleteDatabaseResponse, error) {
	var workspaceID string
	if request.WorkspaceID == nil {
		workspaceID = d.WorkspaceID
	} else {
		workspaceID = *request.WorkspaceID
	}

	return d.generated.DeleteDatabase(ctx, workspaceID, request.DatabaseName)
}

// GetRegions lists available regions.
// https://xata.io/docs/api-reference/workspaces/workspace_id/regions#list-available-regions
func (d databaseCli) GetRegions(ctx context.Context) (*xatagencore.ListRegionsResponse, error) {
	return d.generated.ListRegions(ctx, d.WorkspaceID)
}

// GetRegionsWithWorkspaceID lists available regions for a given workspace ID.
// https://xata.io/docs/api-reference/workspaces/workspace_id/regions#list-available-regions
func (d databaseCli) GetRegionsWithWorkspaceID(ctx context.Context, workspaceID string) (*xatagencore.ListRegionsResponse, error) {
	return d.generated.ListRegions(ctx, workspaceID)
}

// List lists databases for the default workspace.
// https://xata.io/docs/api-reference/workspaces/workspace_id/dbs#list-databases
func (d databaseCli) List(ctx context.Context) (*xatagencore.ListDatabasesResponse, error) {
	return d.generated.GetDatabaseList(ctx, d.WorkspaceID)
}

// ListWithWorkspaceID lists databases for a given workspace ID.
// https://xata.io/docs/api-reference/workspaces/workspace_id/dbs#list-databases
func (d databaseCli) ListWithWorkspaceID(ctx context.Context, workspaceID string) (*xatagencore.ListDatabasesResponse, error) {
	return d.generated.GetDatabaseList(ctx, workspaceID)
}

// Rename renames a database.
// https://xata.io/docs/api-reference/workspaces/workspace_id/dbs/db_name/rename#rename-database
func (d databaseCli) Rename(ctx context.Context, request RenameDatabaseRequest) (*xatagencore.DatabaseMetadata, error) {
	wsID := d.WorkspaceID
	if request.WorkspaceID != nil && *request.WorkspaceID != "" {
		wsID = *request.WorkspaceID
	}

	return d.generated.RenameDatabase(
		ctx,
		wsID,
		request.DatabaseName,
		&xatagencore.RenameDatabaseRequest{NewName: request.NewName},
	)
}

// NewDatabasesClient constructs a client for interacting databases.
func NewDatabasesClient(opts ...ClientOption) (DatabasesClient, error) {
	cliOpts, err := consolidateClientOptionsForCore(opts...)
	if err != nil {
		return nil, err
	}

	dbCfg, err := loadDatabaseConfig()
	if err != nil {
		// No err, because the config values can be provided by the users.
		log.Println(err)
	}

	return databaseCli{
		generated: xatagencore.NewDatabasesClient(
			func(options *xatagenclient.ClientOptions) {
				options.HTTPClient = cliOpts.HTTPClient
				options.BaseURL = cliOpts.BaseURL
				options.Bearer = cliOpts.Bearer
			}),
		WorkspaceID: dbCfg.workspaceID,
		Region:      dbCfg.region,
		BranchName:  dbCfg.branchName,
	}, nil
}
