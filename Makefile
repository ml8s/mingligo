GO ?= go
IMAGE_REGISTRY ?= mingligo-local
IMAGE_TAG ?= local
GOLANGCI_LINT ?= golangci-lint
NODE_BIN := $(HOME)/.nvm/versions/node/v$(shell cat .nvmrc)/bin
NPM_FROM_PATH := $(shell command -v npm 2>/dev/null)
NPM_FALLBACK := $(NODE_BIN)/npm
NPM ?= $(if $(NPM_FROM_PATH),$(NPM_FROM_PATH),$(NPM_FALLBACK))
NPM_BIN_DIR := $(if $(NPM_FROM_PATH),$(dir $(NPM_FROM_PATH)),$(NODE_BIN))
export PATH := $(NPM_BIN_DIR):$(PATH)
GIT_COMMIT ?= $(shell git rev-parse HEAD 2>/dev/null || echo unknown)
GIT_SHA ?= $(shell git rev-parse --short=12 HEAD 2>/dev/null || echo unknown)
BUILD_TIME ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
# A stale host GOROOT must never override the toolchain selected from go.mod.
unexport GOROOT

.PHONY: help build run test vet fmt fmt-check lint security-go check gate docs-lint image-mcp runimage smoke-mcp boundary

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

build: ## Build the mingligo-mcp example binary
	$(GO) build -trimpath -o bin/mingligo-mcp ./examples/mingligo-mcp

run: ## Run the reference MCP server for quick validation
	$(GO) run ./examples/mingligo-mcp --addr :8081

test: ## Run all tests
	$(GO) test -count=1 ./...

boundary: ## Reject fact-core imports of MCP adapters or other internals
	@! grep -RIn --include='*.go' 'github.com/ml8s/mingligo/examples\|github.com/ml8s/mingligo/internal' pkg

vet: ## Run go vet
	$(GO) vet ./...

fmt: ## Format Go sources
	$(GO) fmt ./...

fmt-check: ## Reject unformatted Go sources without rewriting them
	@files="$$(gofmt -l .)"; \
	if [ -n "$$files" ]; then \
		echo "unformatted Go sources:"; \
		echo "$$files"; \
		exit 1; \
	fi

lint: ## Run golangci-lint
	$(GOLANGCI_LINT) run ./...

security-go: ## Run the Go vulnerability scanner
	$(GO) run golang.org/x/vuln/cmd/govulncheck@latest ./...

check: docs-lint fmt-check lint vet boundary ## Static checks: docs, format, lint, vet, and fact-core boundary

docs-lint: ## Run Markdown style checks
	$(NPM) run lint:docs

gate: check test build ## Full source gate: static checks, tests, and binary build

image-mcp: ## Build the mingligo-mcp image
	docker build -f examples/mingligo-mcp/deploy/Dockerfile \
		--build-arg GIT_SHA=$(GIT_SHA) \
		--build-arg COMMIT=$(GIT_COMMIT) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		-t $(IMAGE_REGISTRY)/mingligo-mcp:$(IMAGE_TAG) .

runimage: image-mcp ## Run the mingligo-mcp image for quick validation
	docker run --rm -p 8081:8081 \
		--env MINGLIGO_SERVER_NAME=local/mingligo-mcp \
		$(IMAGE_REGISTRY)/mingligo-mcp:$(IMAGE_TAG)

smoke-mcp: image-mcp ## Smoke the mingligo-mcp image
	IMAGE_TAG=$(IMAGE_TAG) IMAGE_REGISTRY=$(IMAGE_REGISTRY) \
		examples/mingligo-mcp/deploy/smoke-mcp.sh
