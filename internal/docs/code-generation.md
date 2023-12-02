# Generating from OpenAPI spec

This SDK is a wrapper around auto generated Go code from Xata [OpenAPI specs](https://xata.io/docs/rest-api/contexts#openapi-specifications). [Fern](https://github.com/fern-api/fern) is used for code generation.

The process is automated with the following Make targets:

Download the latest server OpenAPI specs
```shell
make download-openapi-specs
```

Before running following commands, ensure that you are logged in to fern:
```shell
make fern-login
```
Note that the code generation automation supports fern version [0.15.10](https://github.com/fern-api/fern/releases/tag/0.15.10).
If there is an issue with code generation, ensure that the installed fern is in the supported version.

Generate code for `core` scope
```shell
make generate-core-code
```

Generate code for `workspace` scope
```shell
make generate-workspace-code
```

Code generation requires some updates in the API specs and auto-generated code for various reasons.
For more information, see [this PR](https://github.com/xataio/xata-go/pull/26#issue-1989477775) and the issues it resolves.

Add the license headers in case a new file was created:
```shell
make check-license-header
```
