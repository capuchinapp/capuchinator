.EXPORT_ALL_VARIABLES:
GOBIN = $(shell pwd)/bin

.PHONY: init
init: tools

.PHONY: deps
deps:
	@go mod tidy

.PHONY: audit
audit: tools
	@export PATH="$(shell pwd)/bin:$(PATH)"; govulncheck ./...

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: test
test:
	@go test -race -failfast -count=1 ./...

.PHONY: run
run:
	@go run -race ./main.go

.PHONY: tools
tools: deps
	@go install golang.org/x/vuln/cmd/govulncheck@latest
