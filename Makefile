GO ?= go
IMAGE_REGISTRY ?= liki-local
IMAGE_TAG ?= local
GIT_COMMIT ?= $(shell git rev-parse HEAD 2>/dev/null || echo unknown)
GIT_SHA ?= $(shell git rev-parse --short=12 HEAD 2>/dev/null || echo unknown)
BUILD_TIME ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

.PHONY: help build run test vet fmt fmt-check image-mcp smoke-mcp

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

build: ## Build the mingligo-mcp binary
	$(GO) build -trimpath -o bin/mingligo-mcp ./cmd/mingligo-mcp

run: ## Run the reference MCP server for quick validation
	$(GO) run ./cmd/mingligo-mcp --addr :8081

test: ## Run all tests
	$(GO) test -count=1 ./...

vet: ## Run go vet
	$(GO) vet ./...

fmt: ## Format Go sources
	$(GO) fmt ./...

fmt-check: ## Reject unformatted Go sources
	@test -z "$$($(GO) fmt ./...)" || { $(GO) fmt -l .; exit 1; }

image-mcp: ## Build the mingligo-mcp image
	docker build -f dev/Dockerfile \
		--build-arg GIT_SHA=$(GIT_SHA) \
		--build-arg COMMIT=$(GIT_COMMIT) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		-t $(IMAGE_REGISTRY)/mingligo-mcp:$(IMAGE_TAG) .

smoke-mcp: image-mcp ## Smoke the mingligo-mcp image
	IMAGE_TAG=$(IMAGE_TAG) IMAGE_REGISTRY=$(IMAGE_REGISTRY) dev/smoke-mcp.sh
