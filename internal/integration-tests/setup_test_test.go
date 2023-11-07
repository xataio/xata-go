// SPDX-License-Identifier: Apache-2.0

package integrationtests

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_setup_cleanup(t *testing.T) {
	cfg, err := setupDatabase()
	if err != nil {
		t.Fatalf("unable to setup db: %v", err)
	}
	assert.NotNil(t, cfg)

	err = setupTableWithColumns(context.Background(), cfg)
	if err != nil {
		t.Fatalf("unable to setup table: %v", err)
	}

	t.Logf("%#v", cfg)

	t.Cleanup(func() {
		err = cleanup(cfg)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func Test_cleanupIntegrationWorkspaces(t *testing.T) {
	if _, found := os.LookupEnv("CLEAN_UP_INTEGRATION_WORKSPACES"); !found {
		t.Skip("skipping integration workspaces cleanup")
	}
	err := cleanAllWorkspaces()
	if err != nil {
		t.Fatal(err)
	}
}
