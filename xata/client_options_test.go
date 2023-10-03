package xata_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	generatedwrapper "github.com/xataio/xata-go/xata"
)

func TestWithAPIToken(t *testing.T) {
	t.Run("should use the provided API key in client options", func(t *testing.T) {
		c := &generatedwrapper.ClientOptions{}
		apiToken := "my-token"
		generatedwrapper.WithAPIKey("my-token")(c)

		assert.Equal(t, apiToken, c.Bearer)
	})
}

func TestWithHTTPClient(t *testing.T) {
	t.Run("should use the provided HTTP client in client options", func(t *testing.T) {
		c := &generatedwrapper.ClientOptions{}
		cli := &http.Client{}
		generatedwrapper.WithHTTPClient(cli)(c)

		assert.Equal(t, cli, c.HTTPClient)
	})
}
