VERSION := $(shell cat VERSION)

.PHONY: build run lint fmt test clean version help

help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the server binary with version
	go build -ldflags "-X main.Version=$(VERSION)" ./cmd/server

run: ## Run the server with version
	go run -ldflags "-X main.Version=$(VERSION)" ./cmd/server

lint: ## Run go vet
	go vet ./...

fmt: ## Run gofmt on all Go files
	gofmt -w .

test: ## Run tests with coverage
	go test ./... -v -cover

clean: ## Remove build artifacts
	rm -f server

version: ## Show current version
	@echo $(VERSION)
