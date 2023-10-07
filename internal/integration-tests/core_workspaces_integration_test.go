package integrationtests

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/xataio/xata-go/xata"
)

func Test_workspacesClient_List(t *testing.T) {
	apiKey, found := os.LookupEnv("XATA_API_KEY")
	if !found {
		t.Skipf("%s not found in env vars", "XATA_API_KEY")
	}

	t.Run("should list the workspaces for the given API key", func(t *testing.T) {
		workspaceCli, err := xata.NewWorkspacesClient(xata.WithAPIKey(apiKey))
		if err != nil {
			t.Fatal(err)
		}

		resp, err := workspaceCli.List(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})
}

func Test_workspacesClient_Create_Delete(t *testing.T) {
	apiKey, found := os.LookupEnv("XATA_API_KEY")
	if !found {
		t.Skipf("%s not found in env vars", "XATA_API_KEY")
	}

	t.Run("should create and delete a workspace", func(t *testing.T) {
		workspaceCli, err := xata.NewWorkspacesClient(xata.WithAPIKey(apiKey))
		if err != nil {
			t.Fatal(err)
		}

		ws, err := workspaceCli.Create(context.Background(), &xata.WorkspaceMeta{Name: "test-integration", Slug: xata.String("test_integration")})
		if err != nil {
			t.Fatal(err)
		}

		assert.NoError(t, err)
		assert.NotNil(t, ws)

		err = workspaceCli.Delete(context.Background(), ws.Id)
		assert.NoError(t, err)
	})
}
