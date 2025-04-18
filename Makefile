APP_NAME=team_dev_api
PORT=8080

run:
	go run ./cmd/server/main.go

fmt:
	go fmt ./...
	goimports -w .

GOLANGCI_LINT := $(shell which golangci-lint || echo "$(shell go env GOPATH)/bin/golangci-lint")

lint:
	@if [ ! -x "$(GOLANGCI_LINT)" ]; then \
		echo "golangci-lint not found, installing..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin latest; \
	fi
	PATH="$(shell go env GOPATH)/bin:$$PATH" $(GOLANGCI_LINT) run ./...

swag:
	swag init -g cmd/server/main.go -o ./docs

.PHONY: run fmt lint swag
