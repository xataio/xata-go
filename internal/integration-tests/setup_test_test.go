package integrationtests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_setup_cleanup(t *testing.T) {
	cfg, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, cfg)

	t.Logf("%#v", cfg)

	t.Cleanup(func() {
		err = cleanup(cfg)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func Test_cleanupIntegrationWorkspaces(t *testing.T) {
	err := cleanAllWorkspaces()
	if err != nil {
		t.Fatal(err)
	}
}
