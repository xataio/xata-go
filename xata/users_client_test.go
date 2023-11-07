// SPDX-License-Identifier: Apache-2.0

package xata_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xataio/xata-go/xata"
	xatagen "github.com/xataio/xata-go/xata/internal/fern-core/generated/go"
	xatagencore "github.com/xataio/xata-go/xata/internal/fern-core/generated/go/core"
)

func TestNewUsersClient(t *testing.T) {
	t.Run("should construct a new users client", func(t *testing.T) {
		got, err := xata.NewUsersClient(xata.WithAPIKey("my-api-token"))
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})
}

func Test_usersCli_Get(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		want       *xatagen.UserWithId
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should get current user",
			want: &xatagen.UserWithId{
				Email:    "email@test.com",
				Fullname: "name lastname",
				Image:    "some-image",
				Id:       "some-id",
			},
			statusCode: http.StatusOK,
		},
	}

	for _, eTC := range errTestCasesCore {
		tests = append(tests, tc{
			name:       eTC.name,
			statusCode: eTC.statusCode,
			apiErr:     eTC.apiErr,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testSrv := testService(t, http.MethodGet, "/user", tt.statusCode, tt.apiErr != nil, tt.want)

			cli, err := xata.NewUsersClient(xata.WithBaseURL(testSrv.URL), xata.WithAPIKey("test-key"))
			if err != nil {
				t.Fatal(err)
			}

			user, err := cli.Get(context.TODO())

			if tt.apiErr != nil {
				errAPI := tt.apiErr.Unwrap()
				if errAPI == nil {
					t.Fatal("expected error but got nil")
				}
				assert.ErrorAs(err, &errAPI)
				assert.Equal(err.Error(), tt.apiErr.Error())
			} else {
				assert.NoError(err)
				assert.Equal(tt.want, user)
			}
		})
	}
}
