VERSION := $(shell cat VERSION)

.PHONY: install dev build lint test test-watch test-e2e preview clean help

help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install: ## Install dependencies
	npm install

dev: ## Start development server
	npm run dev

build: ## Build for production (type check + vite build)
	npm run build

lint: ## Run TypeScript type check
	npm run lint

test: ## Run unit tests
	npm test

test-watch: ## Run tests in watch mode
	npm run test:watch

test-e2e: ## Run end-to-end tests with Playwright
	npm run test:e2e

preview: ## Preview production build
	npm run preview

clean: ## Remove build artifacts
	rm -rf dist node_modules/.vite

version: ## Show current version
	@echo $(VERSION)
