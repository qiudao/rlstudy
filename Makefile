.PHONY: help build test run-env run-agent run stop clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build both binaries
	go build -o bin/bandit-env ./cmd/bandit-env
	go build -o bin/agent ./cmd/agent

test: ## Run all tests
	go test ./... -v

run-env: build ## Start bandit environment server (foreground)
	./bin/bandit-env

run-agent: build ## Run agent experiment (env must be running)
	./bin/agent

run: build ## Start env daemon, run agent, stop env
	@./bin/bandit-env -daemon & sleep 0.3
	@./bin/agent; ret=$$?; $(MAKE) -s stop; exit $$ret

stop: ## Stop env daemon
	@if [ -f .bandit-env.pid ]; then kill $$(cat .bandit-env.pid) 2>/dev/null; rm -f .bandit-env.pid /tmp/bandit-env.sock; echo "bandit-env stopped"; else echo "no daemon running"; fi

clean: ## Remove build artifacts
	rm -rf bin/
