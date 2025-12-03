.PHONY: help build test migrate-up migrate-down migrate-create run-inventory run-issuer run-adapter docker-build clean

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Database migrations
MIGRATE_CMD = migrate -path backend/migrations -database "postgres://cmp_user:cmp_pass@localhost:5432/cmp_db?sslmode=disable" -verbose

migrate-up: ## Apply database migrations
	@which migrate > /dev/null || (echo "Install golang-migrate: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest" && exit 1)
	$(MIGRATE_CMD) up

migrate-down: ## Rollback last migration
	@which migrate > /dev/null || (echo "Install golang-migrate: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest" && exit 1)
	$(MIGRATE_CMD) down 1

migrate-create: ## Create a new migration (usage: make migrate-create NAME=add_table)
	@which migrate > /dev/null || (echo "Install golang-migrate: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest" && exit 1)
	@test -n "$(NAME)" || (echo "Usage: make migrate-create NAME=migration_name" && exit 1)
	migrate create -ext sql -dir backend/migrations -seq $(NAME)

# Testing
test: ## Run unit tests
	cd backend && go test -v -race -coverprofile=coverage.out ./...

test-integration: ## Run integration tests
	cd tests/integration && go test -v -tags=integration ./...

test-coverage: ## Generate test coverage report
	cd backend && go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html

# Building
build: ## Build all backend services
	cd backend && go build -o bin/inventory-service ./cmd/inventory-service
	cd backend && go build -o bin/issuer-service ./cmd/issuer-service
	cd backend && go build -o bin/adapter-service ./cmd/adapter-service

build-agent: ## Build Linux agent
	cd agents/linux-agent && go build -o bin/agent ./cmd/agent

# Running services (dev mode)
run-inventory: ## Run inventory service locally
	cd backend && go run ./cmd/inventory-service -port=8081

run-issuer: ## Run issuer service locally
	cd backend && go run ./cmd/issuer-service -port=8082

run-adapter: ## Run adapter service locally
	cd backend && go run ./cmd/adapter-service -port=8083

# Docker
docker-build: ## Build Docker images
	docker-compose -f deploy/docker-compose.yml build

# Cleanup
clean: ## Remove build artifacts
	rm -rf backend/bin
	rm -rf agents/linux-agent/bin
	rm -rf backend/coverage.out backend/coverage.html
