// SPDX-License-Identifier: Apache-2.0

package xata_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	generatedwrapper "github.com/xataio/xata-go/xata"
)

func TestClientOptions(t *testing.T) {
	t.Run("WithAPIKey", func(t *testing.T) {
		c := &generatedwrapper.ClientOptions{}
		apiToken := "my-token"
		generatedwrapper.WithAPIKey("my-token")(c)
		assert.Equal(t, apiToken, c.Bearer)
	})
	t.Run("WithHTTPClient", func(t *testing.T) {
		c := &generatedwrapper.ClientOptions{}
		cli := &http.Client{}
		generatedwrapper.WithHTTPClient(cli)(c)
		assert.Equal(t, cli, c.HTTPClient)
	})
	t.Run("WithWorkspaceID", func(t *testing.T) {
		c := &generatedwrapper.ClientOptions{}
		workspaceID := "workspace-123"
		generatedwrapper.WithWorkspaceID(workspaceID)(c)
		assert.Equal(t, workspaceID, c.WorkspaceID)
	})
	t.Run("WithBranch", func(t *testing.T) {
		c := &generatedwrapper.ClientOptions{}
		branch := "branch-123"
		generatedwrapper.WithBranch(branch)(c)
		assert.Equal(t, branch, c.Branch)
	})
	t.Run("WithRegion", func(t *testing.T) {
		c := &generatedwrapper.ClientOptions{}
		region := "region-123"
		generatedwrapper.WithRegion(region)(c)
		assert.Equal(t, region, c.Region)
	})
}
