.PHONY: build
build: ## Build osdd CLI with VERSION (Usage: make build VERSION=vX.Y.Z)
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION not set. Usage: make build VERSION=vX.Y.Z"; \
		exit 1; \
	fi
	./build.sh

.PHONY: build-dev
build-dev: ## Build osdd CLI for development (version=dev)
	mkdir -p out && go build -o out/osdd ./cmd/osdd

.PHONY: clean
clean: ## Remove built binaries
	rm -f ./out/

.PHONY: test
test: ## Run tests
	go test ./...

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'
