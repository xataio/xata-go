# How to add a new client or method

Let's assume we want to add a client for interacting with the workspace-related operations, such as creation, 
modification, or deletion, as documented in the API.

First, we need to ensure that the auto-generated code from the API definition exists. 
For code generation see [here](code-generation.md). 
The auto-generated code for the workspaces client is located in the 
`xata/internal/fern-core/generated/go/workspaces_client.go` file. 
`fern-core` is for the codes generated from the Core API. 
`fern-workspace` is for the Workspace API.

The files in the `xata/internal` folder are not accessible from the SDK. 
For exposing any auto-generated code we need to create wrapper clients and their methods in the `xata` folder manually. 
Each client lives in a dedicated file, i.e., `xata/workspaces_client.go.`
The following code snippet shows how to create a `List` method on the workspaces client that actually points to 
the auto-generated workspaces client, and it's `GetWorkspacesList` method.
```go
type WorkspacesClient interface {
List(ctx context.Context) (*xatagencore.GetWorkspacesListResponse, error)
...
}

type workspaceCli struct {
generated   xatagencore.WorkspacesClient
workspaceID string
}

func (w workspaceCli) List(ctx context.Context) (*xatagencore.GetWorkspacesListResponse, error) {
return w.generated.GetWorkspacesList(ctx)
}
```
In some cases, the auto-generated code might not be quite useful due to the code generation capability limitations. 
In these cases, the manually generated wrapper code has to be expanded for a good user experience. 
See the [records client](../../xata/records_client.go) as an example of such a scenario. 
This might require modifying the auto-generated code, such as defining a new model or updating a method signature. 
To be able to streamline the code generation process we need to incorporate these modifications via the 
[`code_gen.go`](../../xata/internal/code-gen/code_gen.go) script that applies defined alterations on the 
API definitions or the auto-generated Go code.