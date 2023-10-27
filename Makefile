.PHONY: lint gofmt test run integration-test

gofmt:
	docker compose -f internal/build/docker-compose.yaml up gofmt

lint: gofmt
	docker compose -f internal/build/docker-compose.yaml up lint

smoke-test:
	@echo "Running test app for sanity check"
	@cd internal/smoke-tests && go run . && cd ../..

test:
	@echo "Running unit tests"
	@go test -count=1 -cover -race ./...

integration-test:
	@echo "Running integration test"
	@go test -v --tags=integration -count=1 -cover -race ./internal/integration-tests

download-openapi-specs:
	@echo "Downloading openapi specs"
	@cd internal/docs && go run . && cd ../..