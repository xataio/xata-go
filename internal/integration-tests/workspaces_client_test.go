// SPDX-License-Identifier: Apache-2.0

package integrationtests

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
	"github.com/xataio/xata-go/xata"
)

func Test_workspacesClient(t *testing.T) {
	apiKey, found := os.LookupEnv(xata.EnvXataAPIKey)
	if !found {
		t.Skipf("%s not found in env vars", xata.EnvXataAPIKey)
	}

	t.Run("should create, get, list, update and delete workspace", func(t *testing.T) {
		workspaceCli, err := xata.NewWorkspacesClient(
			xata.WithAPIKey(apiKey),
			xata.WithHTTPClient(retryablehttp.NewClient().StandardClient()),
		)
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.Background()

		ws, err := workspaceCli.Create(ctx, &xata.WorkspaceMeta{Name: "test-integration", Slug: xata.String("test_integration")})
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, ws)

		workspace, err := workspaceCli.GetWithWorkspaceID(ctx, ws.Id)
		assert.Nil(t, err)
		assert.Equal(t, ws.Id, workspace.Id)

		resp, err := workspaceCli.List(context.Background())
		assert.Nil(t, err)
		var wsIDs []string
		for _, ws := range resp.Workspaces {
			wsIDs = append(wsIDs, ws.Id)
		}
		assert.Contains(t, wsIDs, workspace.Id)

		updatedWSName := "updated-name"
		updated, err := workspaceCli.Update(ctx, xata.UpdateWorkspaceRequest{
			Payload:     &xata.WorkspaceMeta{Name: updatedWSName},
			WorkspaceID: xata.String(workspace.Id),
		})
		assert.Nil(t, err)
		assert.Equal(t, updatedWSName, updated.Name)

		err = workspaceCli.Delete(context.Background(), ws.Id)
		assert.NoError(t, err)
	})
}
