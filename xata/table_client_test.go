package xata_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xataio/xata-go/xata"
)

func TestNewTableClient(t *testing.T) {
	t.Run("should construct a new client", func(t *testing.T) {
		got, err := xata.NewTableClient(
			xata.WithBaseURL("https://www.example.com"),
			xata.WithAPIKey("my-api-token"),
		)
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})
}
