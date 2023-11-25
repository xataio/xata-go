# How to add a new client or method

The SDK is modelling the namespaces from Xata's [API reference](https://xata.io/docs/api-reference/user), and creates a new client for each namespace. Example given, the `user` APIs listed in the [reference](https://xata.io/docs/api-reference/user) would be modelled in a `UserClient`.

Let's assume we want to add a new client for interacting with the workspace-related operations, such as creation, 
modification, or deletion, as documented in the [API reference](https://xata.io/docs/api-reference/workspaces). 
This new client would contain the mentioned operations, modelled in a `WorkspaceClient`, see [workspaces_client.go](https://github.com/xataio/xata-go/blob/main/xata/workspaces_client.go) as a reference point.

First, we need to ensure that the generated code from the API definition exists, for code generation consult [this guide](code-generation.md). 
The generated code for the workspaces client is located in the `xata/internal/fern-core/generated/go/workspaces_client.go` file. 
* `fern-core` is for the codes generated from the Core API
* `fern-workspace` is for the Workspace API

The files in the `xata/internal` folder are not accessible from the SDK. 
For exposing any generated code we need to create wrapper clients and their methods in the `xata` package manually. 
Each client lives in a dedicated file, i.e., `xata/workspaces_client.go`.
The wrappers have the purpose to create a leaner method names and allow for a level of complexity abstraction.
The following code snippet shows how to create a `List` method on the workspace client that points to the generated workspace client, and it's `GetWorkspacesList` method.

```go
type WorkspacesClient interface {
    List(ctx context.Context) (*xatagencore.GetWorkspacesListResponse, error)
}

type workspaceCli struct {
    generated   xatagencore.WorkspacesClient
    workspaceID string
}

func (w workspaceCli) List(ctx context.Context) (*xatagencore.GetWorkspacesListResponse, error) {
    return w.generated.GetWorkspacesList(ctx)
}
```

In some cases, the generated code might not be idomatic, due to limitations in code generation. 
In these cases, the manually generated wrapper code has to be expanded for an improved user experience. 
See the [RecordsClient](../../xata/records_client.go) as an example of such a scenario. 
This might require modifying the generated code by hand, such as defining a new model or updating a method signature. 
To be able to streamline the code generation process we need to incorporate these modifications via the 
[`code_gen.go`](../../xata/internal/code-gen/code_gen.go) script that applies defined alterations on the 
API definitions or the generated Go code.
