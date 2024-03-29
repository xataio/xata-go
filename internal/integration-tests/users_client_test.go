// SPDX-License-Identifier: Apache-2.0

package integrationtests

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
	"github.com/xataio/xata-go/xata"
)

func Test_usersClient(t *testing.T) {
	apiKey, found := os.LookupEnv(xata.EnvXataAPIKey)
	if !found {
		t.Skipf("%s not found in env vars", xata.EnvXataAPIKey)
	}

	t.Run("should get the current user", func(t *testing.T) {
		usersCli, err := xata.NewUsersClient(
			xata.WithAPIKey(apiKey),
			xata.WithHTTPClient(retryablehttp.NewClient().StandardClient()),
		)
		if err != nil {
			log.Fatal(err)
		}

		user, err := usersCli.Get(context.TODO())
		assert.NoError(t, err)
		assert.NotNil(t, user)
	})
}
