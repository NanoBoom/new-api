FRONTEND_DIR = ./web/default
FRONTEND_CLASSIC_DIR = ./web/classic
BACKEND_DIR = .

# Docker image config — override via environment, e.g. `make IMAGE_NAME=your-registry/new-api push`
IMAGE_NAME ?= docker.echojoy.cn/gateway/new-api
IMAGE_TAG  ?= $(shell git rev-parse --short HEAD)
PLATFORM   ?= linux/amd64
BUILDER    ?= new-api-builder

.PHONY: all build-frontend build-frontend-classic build-all-frontends start-backend dev dev-api dev-web dev-web-classic buildx-init build push release

help: ## Prints help for targets with comments
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

all: build-all-frontends start-backend

build-frontend:
	@echo "Building default frontend..."
	@cd $(FRONTEND_DIR) && bun install && DISABLE_ESLINT_PLUGIN='true' VITE_REACT_APP_VERSION=$$(cat ../../VERSION) bun run build

build-frontend-classic:
	@echo "Building classic frontend..."
	@cd $(FRONTEND_CLASSIC_DIR) && bun install && VITE_REACT_APP_VERSION=$$(cat ../../VERSION) bun run build

build-all-frontends: build-frontend build-frontend-classic

start-backend: ## start server
	@echo "Starting backend dev server..."
	@cd $(BACKEND_DIR) && go run main.go &

dev-api:
	@echo "Starting backend services (docker)..."
	@docker compose -f docker-compose.dev.yml up -d

dev-web:
	@echo "Starting frontend dev server..."
	@cd $(FRONTEND_DIR) && bun install && bun run dev

dev-web-classic:
	@echo "Starting classic frontend dev server..."
	@cd $(FRONTEND_CLASSIC_DIR) && bun install && bun run dev

dev: dev-api dev-web

buildx-init: ## Ensure a buildx builder exists (needed for cross-arch on macOS)
	@docker buildx inspect $(BUILDER) >/dev/null 2>&1 || docker buildx create --name $(BUILDER) --use
	@docker buildx use $(BUILDER)

build: buildx-init ## Build linux/amd64 docker image locally (loaded into local docker)
	@echo "Building $(IMAGE_NAME):$(IMAGE_TAG) for $(PLATFORM)..."
	@docker buildx build \
		--platform $(PLATFORM) \
		--build-arg VERSION=$(IMAGE_TAG) \
		-t $(IMAGE_NAME):$(IMAGE_TAG) \
		-t $(IMAGE_NAME):latest \
		--load .

push: buildx-init ## Build linux/amd64 image and push to registry
	@echo "Building & pushing $(IMAGE_NAME):$(IMAGE_TAG) for $(PLATFORM)..."
	@docker buildx build \
		--platform $(PLATFORM) \
		--build-arg VERSION=$(IMAGE_TAG) \
		-t $(IMAGE_NAME):$(IMAGE_TAG) \
		-t $(IMAGE_NAME):latest \
		--push .

release: push ## Alias for push
