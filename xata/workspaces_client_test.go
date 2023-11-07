// SPDX-License-Identifier: Apache-2.0

package xata_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"

	"github.com/xataio/xata-go/xata"
	xatagen "github.com/xataio/xata-go/xata/internal/fern-core/generated/go"
	xatagencore "github.com/xataio/xata-go/xata/internal/fern-core/generated/go/core"

	"github.com/stretchr/testify/assert"
)

type errBody struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func (e errBody) Error() string {
	response, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	return string(response)
}

var testErrBody = errBody{
	ID:      "test-err-id",
	Message: "test-err-message",
}

func TestNewWorkspacesClient(t *testing.T) {
	t.Run("should construct a new workspace client", func(t *testing.T) {
		got, err := xata.NewWorkspacesClient(xata.WithAPIKey("my-api-token"))
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})
}

var errTestCasesCore = []struct {
	name       string
	statusCode int
	apiErr     *xatagencore.APIError
}{
	{
		name:       "should return bad request error",
		statusCode: http.StatusBadRequest,
		apiErr:     xatagencore.NewAPIError(http.StatusBadRequest, testErrBody),
	},
	{
		name:       "should return authentication error",
		statusCode: http.StatusUnauthorized,
		apiErr:     xatagencore.NewAPIError(http.StatusUnauthorized, testErrBody),
	},
	{
		name:       "should return not-found error",
		statusCode: http.StatusNotFound,
		apiErr:     xatagencore.NewAPIError(http.StatusNotFound, testErrBody),
	},
	{
		name:       "should return server error",
		statusCode: http.StatusServiceUnavailable,
		apiErr:     xatagencore.NewAPIError(http.StatusServiceUnavailable, testErrBody),
	},
	{
		name:       "should handle undocumented (not in the api specs) error",
		statusCode: http.StatusConflict,
		apiErr:     xatagencore.NewAPIError(http.StatusConflict, testErrBody),
	},
}

func Test_workspacesClient_List(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		want       *xatagen.GetWorkspacesListResponse
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should list workspaces successfully for owner",
			want: &xatagen.GetWorkspacesListResponse{
				Workspaces: []*xatagen.GetWorkspacesListResponseWorkspacesItem{
					{
						Id:   "some-id",
						Name: "test-workspace-dbName",
						Slug: "test-workspace-slug",
						Role: xatagen.RoleOwner,
					},
				},
			},
			statusCode: http.StatusOK,
		},
		{
			name: "should list workspaces successfully for maintainer",
			want: &xatagen.GetWorkspacesListResponse{
				Workspaces: []*xatagen.GetWorkspacesListResponseWorkspacesItem{
					{
						Id:   "some-id",
						Name: "test-workspace-dbName",
						Slug: "test-workspace-slug",
						Role: xatagen.RoleMaintainer,
					},
				},
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
			testSrv := testService(t, http.MethodGet, "/workspaces", tt.statusCode, tt.apiErr != nil, tt.want)

			wsCli, err := xata.NewWorkspacesClient(xata.WithBaseURL(testSrv.URL), xata.WithAPIKey("test-key"))
			if err != nil {
				t.Fatal(err)
			}

			got, err := wsCli.List(context.TODO())

			if tt.apiErr != nil {
				errAPI := tt.apiErr.Unwrap()
				if errAPI == nil {
					t.Fatal("expected error but got nil")
				}
				assert.ErrorAs(err, &errAPI)
				assert.Equal(err.Error(), tt.apiErr.Error())
				assert.Nil(got)
			} else {
				assert.Equal(tt.want, got)
				assert.NoError(err)
			}
		})
	}
}

func Test_workspacesClient_Create(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		request    *xata.WorkspaceMeta
		want       *xatagen.Workspace
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should create workspace successfully",
			request: &xata.WorkspaceMeta{
				Name: "my-workspace",
				Slug: xata.String("my_workspace"),
			},
			want: &xatagen.Workspace{
				Name:        "my-workspace",
				Slug:        xata.String("my_workspace"),
				Id:          uuid.NewString(),
				MemberCount: 0,
				Plan:        xatagen.WorkspacePlanFree,
			},
			statusCode: http.StatusCreated,
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
			testSrv := testService(t, http.MethodPost, "/workspaces", tt.statusCode, tt.apiErr != nil, tt.want)

			wsCli, err := xata.NewWorkspacesClient(xata.WithBaseURL(testSrv.URL), xata.WithAPIKey("test-key"))
			if err != nil {
				t.Fatal(err)
			}

			got, err := wsCli.Create(context.TODO(), &xata.WorkspaceMeta{})

			if tt.apiErr != nil {
				errAPI := tt.apiErr.Unwrap()
				if errAPI == nil {
					t.Fatal("expected error but got nil")
				}
				assert.ErrorAs(err, &errAPI)
				assert.Equal(err.Error(), tt.apiErr.Error())
				assert.Nil(got)
			} else {
				assert.Equal(tt.want, got)
				assert.NoError(err)
			}
		})
	}
}

func Test_workspacesClient_Delete(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name:       "should delete workspace successfully",
			statusCode: http.StatusNoContent,
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
			testSrv := testService(t, http.MethodDelete, "/workspaces", tt.statusCode, tt.apiErr != nil, nil)

			wsCli, err := xata.NewWorkspacesClient(xata.WithBaseURL(testSrv.URL), xata.WithAPIKey("test-key"))
			if err != nil {
				t.Fatal(err)
			}

			err = wsCli.Delete(context.TODO(), uuid.NewString())

			if tt.apiErr != nil {
				errAPI := tt.apiErr.Unwrap()
				if errAPI == nil {
					t.Fatal("expected error but got nil")
				}
				assert.ErrorAs(err, &errAPI)
				assert.Equal(err.Error(), tt.apiErr.Error())
			} else {
				assert.NoError(err)
			}
		})
	}
}

func Test_workspacesClient_Get(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		want       *xatagen.Workspace
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should get workspace",
			want: &xatagen.Workspace{
				Name:        "ws-name",
				Slug:        xata.String("slug"),
				Id:          "some-id",
				MemberCount: 2,
				Plan:        xatagen.WorkspacePlanFree,
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
			testSrv := testService(t, http.MethodGet, "/workspaces", tt.statusCode, tt.apiErr != nil, tt.want)

			wsCli, err := xata.NewWorkspacesClient(xata.WithBaseURL(testSrv.URL), xata.WithAPIKey("test-key"))
			if err != nil {
				t.Fatal(err)
			}

			got, err := wsCli.GetWithWorkspaceID(context.TODO(), "some-id")

			if tt.apiErr != nil {
				errAPI := tt.apiErr.Unwrap()
				if errAPI == nil {
					t.Fatal("expected error but got nil")
				}
				assert.ErrorAs(err, &errAPI)
				assert.Equal(err.Error(), tt.apiErr.Error())
				assert.Nil(got)
			} else {
				assert.Equal(tt.want, got)
				assert.NoError(err)
			}

			got, err = wsCli.Get(context.TODO())

			if tt.apiErr != nil {
				errAPI := tt.apiErr.Unwrap()
				if errAPI == nil {
					t.Fatal("expected error but got nil")
				}
				assert.ErrorAs(err, &errAPI)
				assert.Equal(err.Error(), tt.apiErr.Error())
				assert.Nil(got)
			} else {
				assert.Equal(tt.want, got)
				assert.NoError(err)
			}
		})
	}
}

func Test_workspacesClient_Update(t *testing.T) {
	assert := assert.New(t)

	type tc struct {
		name       string
		request    xata.UpdateWorkspaceRequest
		want       *xatagen.Workspace
		statusCode int
		apiErr     *xatagencore.APIError
	}

	tests := []tc{
		{
			name: "should update ws",
			request: xata.UpdateWorkspaceRequest{
				Payload:     nil,
				WorkspaceID: nil,
			},
			want: &xatagen.Workspace{
				Name:        "my-workspace",
				Slug:        xata.String("my_workspace"),
				Id:          uuid.NewString(),
				MemberCount: 0,
				Plan:        xatagen.WorkspacePlanFree,
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
			testSrv := testService(t, http.MethodPut, "/workspaces", tt.statusCode, tt.apiErr != nil, tt.want)

			wsCli, err := xata.NewWorkspacesClient(xata.WithBaseURL(testSrv.URL), xata.WithAPIKey("test-key"))
			if err != nil {
				t.Fatal(err)
			}

			got, err := wsCli.Update(context.TODO(), tt.request)

			if tt.apiErr != nil {
				errAPI := tt.apiErr.Unwrap()
				if errAPI == nil {
					t.Fatal("expected error but got nil")
				}
				assert.ErrorAs(err, &errAPI)
				assert.Equal(err.Error(), tt.apiErr.Error())
				assert.Nil(got)
			} else {
				assert.Equal(tt.want, got)
				assert.NoError(err)
			}
		})
	}
}
