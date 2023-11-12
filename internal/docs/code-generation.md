This SDK is a wrapper around an auto generated Go code from the OpenAPI specs.
[Fern](https://github.com/fern-api/fern) is used for code generation.

The process is automated with the following Make targets:

Download the latest API specs
```shell
make download-openapi-specs
```

Generate code for CORE scope
```shell
make generate-core-cod
```

Generate code for WORKSPACE scope
```shell
make generate-workspace-cod
```

Code generation requires some updates in the API specs and auto-generated code for various reasons.
For more information, see [this PR](https://github.com/xataio/xata-go/pull/26#issue-1989477775) and the issues it resolves.
