.PHONY: help build test run-env run-agent clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build both binaries
	go build -o bin/bandit-env ./cmd/bandit-env
	go build -o bin/agent ./cmd/agent

test: ## Run all tests
	go test ./... -v

run-env: build ## Start bandit environment server
	./bin/bandit-env

run-agent: build ## Run agent experiment (env must be running)
	./bin/agent

clean: ## Remove build artifacts
	rm -rf bin/
