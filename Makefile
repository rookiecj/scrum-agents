.PHONY: build test lint fmt clean install run version help

help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build both backend and frontend
	$(MAKE) -C backend build
	$(MAKE) -C frontend build

test: ## Run all tests
	$(MAKE) -C backend test
	$(MAKE) -C frontend test

lint: ## Run linters on both codebases
	$(MAKE) -C backend lint
	$(MAKE) -C frontend lint

fmt: ## Format backend Go files
	$(MAKE) -C backend fmt

install: ## Install frontend dependencies
	$(MAKE) -C frontend install

run: ## Run the backend server
	$(MAKE) -C backend run

clean: ## Remove all build artifacts
	$(MAKE) -C backend clean
	$(MAKE) -C frontend clean

version: ## Show versions
	@echo "Backend:  $$(cat backend/VERSION)"
	@echo "Frontend: $$(cat frontend/VERSION)"
