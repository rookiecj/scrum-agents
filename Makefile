.PHONY: build build-backend build-frontend test test-backend test-frontend run lint lint-backend lint-frontend clean help

# Default target
help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: build-backend build-frontend ## Build both backend and frontend

build-backend: ## Build Go backend
	cd backend && go build ./...

build-frontend: ## Build frontend
	cd frontend && npm run build

test: test-backend test-frontend ## Run all tests with coverage

test-backend: ## Run Go backend tests with coverage
	cd backend && go test ./... -v -cover

test-frontend: ## Run frontend tests
	cd frontend && npm test

run: ## Run the Go backend server
	cd backend && go run ./cmd/server

lint: lint-backend lint-frontend ## Run linters on both codebases

lint-backend: ## Run go vet on backend
	cd backend && go vet ./...

lint-frontend: ## Run TypeScript type check on frontend
	cd frontend && npm run lint

clean: ## Remove build artifacts
	rm -rf backend/cmd/server/server
	rm -rf frontend/dist
