package integrationtests

import (
	"context"
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

//func Test_cleanupIntegrationWorkspaces(t *testing.T) {
//	err := cleanAllWorkspaces()
//	if err != nil {
//		t.Fatal(err)
//	}
//}
