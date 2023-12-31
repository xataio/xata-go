.PHONY: lint gofmt test run integration-test

gofmt: ## run gofmt
	docker compose -f internal/build/docker-compose.yaml up gofmt

lint: gofmt ## run linter after gofmt
	docker compose -f internal/build/docker-compose.yaml up lint

smoke-test: ## smoke tests
	@echo "Running test app for sanity check"
	@cd internal/smoke-tests && go run . && cd ../..

test: ## run unit tests
	TEST_DIRECTORY=./xata go run gotest.tools/gotestsum@latest --format testname

integration-test: ## run integration tests
	TEST_DIRECTORY=./internal/integration-tests go run gotest.tools/gotestsum@latest --format testname

download-openapi-specs: ## download openapi specs
	@echo "Downloading openapi specs"
	@cd internal/docs && go run . && cd ../..

clean-workspaces: ## cleanup
	@echo "Cleaning integration test workspaces"
	CLEAN_UP_INTEGRATION_WORKSPACES=true go test -v -run Test_cleanupIntegrationWorkspaces ./...

check-license-header: ## Check if all *.go files have a license header
	curl -s https://raw.githubusercontent.com/lluissm/license-header-checker/master/install.sh | bash
	./bin/license-header-checker -a -r .github/license-header.txt . go

help: ## Display help
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
#------------- <https://suva.sh/posts/well-documented-makefiles> --------------
.PHONY: fern-login ## login fern code generator
fern-login:
	fern login

.PHONY: generate-core-code ## generate code for the scope core
generate-core-code:
	go run xata/internal/code-gen/code_gen.go -scope=core

.PHONY: generate-workspace-code ## generate code for the scope workspace
generate-workspace-code:
	go run xata/internal/code-gen/code_gen.go -scope=workspace