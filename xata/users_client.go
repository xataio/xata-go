// SPDX-License-Identifier: Apache-2.0

package xata

import (
	"context"

	xatagencore "github.com/xataio/xata-go/xata/internal/fern-core/generated/go"
	xatagenclient "github.com/xataio/xata-go/xata/internal/fern-core/generated/go/core"
)

type UsersClient interface {
	Get(ctx context.Context) (*xatagencore.UserWithId, error)
}

type usersCli struct {
	generated xatagencore.UsersClient
}

// Get returns details of the user making the request.
// https://xata.io/docs/api-reference/user#get-user-details
func (u usersCli) Get(ctx context.Context) (*xatagencore.UserWithId, error) {
	return u.generated.GetUser(ctx)
}

// NewUsersClient constructs a client for interacting users.
func NewUsersClient(opts ...ClientOption) (UsersClient, error) {
	cliOpts, err := consolidateClientOptionsForCore(opts...)
	if err != nil {
		return nil, err
	}

	return usersCli{
		generated: xatagencore.NewUsersClient(
			func(options *xatagenclient.ClientOptions) {
				options.HTTPClient = cliOpts.HTTPClient
				options.BaseURL = cliOpts.BaseURL
				options.Bearer = cliOpts.Bearer
			}),
	}, nil
}
