# Makefile for lab2-terrylneal

GO = go
MIGRATE = migrate
DB_URL ?= postgresql://lab2_user:lab2pass@localhost/lab2_terrylneal?sslmode=disable
PORT ?= 4000

# Colors
GREEN = \033[0;32m
RED = \033[0;31m
YELLOW = \033[1;33m
NC = \033[0m 

.PHONY: help run build test migrate-up migrate-down migrate-force clean db-create db-drop

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

env-check: ## Check if required environment variables are set
	@if [ -z "$(DB_URL)" ]; then \
		echo "$(RED)Error: DB_URL is not set$(NC)"; \
		echo "Please set it in .envrc or export it manually"; \
		exit 1; \
	fi
	@echo "$(GREEN)Environment variables OK$(NC)"

run: env-check ## Run the app
	@echo "$(GREEN)Starting server on port $(PORT)...$(NC)"
	$(GO) run cmd/api/main.go

build: env-check ## Build the app
	@echo "$(GREEN)Building application...$(NC)"
	$(GO) build -o bin/api cmd/api/main.go
	@echo "$(GREEN)Build complete! Binary saved to bin/api$(NC)"

test: ## Run tests
	@echo "$(GREEN)Running tests...$(NC)"
	$(GO) test -v ./...

test-json: ## Run tests with JSON output
	$(GO) test -json ./...

cover: ## Run tests with coverage
	$(GO) test -cover ./...

migrate-create: ## Create a new migration
	@if [ -z "$(NAME)" ]; then \
		echo "$(RED)Error: NAME is required. Usage: make migrate-create NAME=your_migration_name$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)Creating new migration: $(NAME)...$(NC)"
	$(MIGRATE) create -ext sql -dir ./migrations -seq $(NAME)

migrate-up: env-check ## Run all up migrations
	@echo "$(GREEN)Running migrations up...$(NC)"
	$(MIGRATE) -path ./migrations -database "$(DB_URL)" up

migrate-up-1: env-check ## Run the next migration
	@echo "$(GREEN)Running next migration...$(NC)"
	$(MIGRATE) -path ./migrations -database "$(DB_URL)" up 1

migrate-down: env-check ## Rollback the last migration
	@echo "$(RED)Rolling back last migration...$(NC)"
	$(MIGRATE) -path ./migrations -database "$(DB_URL)" down 1

migrate-down-all: env-check ## Rollback all migrations
	@echo "$(RED)Rolling back all migrations...$(NC)"
	$(MIGRATE) -path ./migrations -database "$(DB_URL)" down -all

migrate-force: env-check ## Force set migration version
	@if [ -z "$(VERSION)" ]; then \
		echo "$(RED)Error: VERSION is required. Usage: make migrate-force VERSION=1$(NC)"; \
		exit 1; \
	fi
	@echo "$(RED)Forcing migration version to $(VERSION)...$(NC)"
	$(MIGRATE) -path ./migrations -database "$(DB_URL)" force $(VERSION)

migrate-version: env-check ## Show current migration version
	@echo "$(YELLOW)Current migration version:$(NC)"
	$(MIGRATE) -path ./migrations -database "$(DB_URL)" version

migrate-status: env-check ## Show migration status
	@echo "$(YELLOW)Migration status:$(NC)"
	$(MIGRATE) -path ./migrations -database "$(DB_URL)" version

db-create: ## Create the database and user (requires postgres superuser)
	@echo "$(GREEN)Creating database and user...$(NC)"
	sudo -u postgres psql -c "CREATE DATABASE lab2_terrylneal;"
	sudo -u postgres psql -c "CREATE USER lab2_user WITH PASSWORD 'lab2pass';"
	sudo -u postgres psql -c "ALTER DATABASE lab2_terrylneal OWNER TO lab2_user;"
	sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE lab2_terrylneal TO lab2_user;"
	sudo -u postgres psql -d lab2_terrylneal -c "GRANT ALL ON SCHEMA public TO lab2_user;"

db-drop: ## Drop the database and user (careful!)
	@echo "$(RED)Dropping database and user...$(NC)"
	sudo -u postgres psql -c "DROP DATABASE IF EXISTS lab2_terrylneal;"
	sudo -u postgres psql -c "DROP USER IF EXISTS lab2_user;"

db-reset: db-drop db-create migrate-up ## Reset database and run migrations
	@echo "$(GREEN)Database reset complete!$(NC)"

db-connect: ## Connect to the database using psql
	@echo "$(GREEN)Connecting to database...$(NC)"
	@command -v psql >/dev/null 2>&1 || { echo "$(RED)psql is not installed$(NC)" >&2; exit 1; }
	psql "$(DB_URL)"

db-migrate-status: env-check ## Show detailed migration status
	@echo "$(YELLOW)Checking migration status...$(NC)"
	@if [ -d "./migrations" ]; then \
		echo "$(GREEN)Migrations directory exists$(NC)"; \
		echo "$(YELLOW)Migration files:$(NC)"; \
		ls -la ./migrations/*.sql 2>/dev/null || echo "No migration files found"; \
	else \
		echo "$(RED)Migrations directory not found$(NC)"; \
	fi
	@echo "$(YELLOW)Current database version:$(NC)"
	-$(MIGRATE) -path ./migrations -database "$(DB_URL)" version 2>/dev/null || echo "No migrations applied yet"

deps: ## Download dependencies
	$(GO) mod download
	$(GO) mod tidy

clean: ## Clean build artifacts
	@echo "$(GREEN)Cleaning...$(NC)"
	rm -rf bin/
	$(GO) clean

fmt: ## Format code
	$(GO) fmt ./...

lint: ## Run linter (requires golangci-lint)
	@which golangci-lint > /dev/null && golangci-lint run || echo "golangci-lint not installed"

dev: env-check ## Run with live reload (requires air)
	@which air > /dev/null && air || echo "air not installed, using regular run"
	$(MAKE) run

curl-get: ## Test GET endpoints with curl
	@echo "$(GREEN)Testing / endpoint...$(NC)"
	curl -s http://localhost:$(PORT)/
	@echo "\n$(GREEN)Testing /about endpoint...$(NC)"
	curl -s http://localhost:$(PORT)/about
	@echo "\n$(GREEN)Testing /api/info endpoint...$(NC)"
	curl -s http://localhost:$(PORT)/api/info | jq '.' 2>/dev/null || curl -s http://localhost:$(PORT)/api/info
	@echo "\n$(GREEN)Testing /api/users endpoint...$(NC)"
	curl -s http://localhost:$(PORT)/api/users | jq '.' 2>/dev/null || curl -s http://localhost:$(PORT)/api/users

curl-post: ## Test POST endpoint with curl
	@echo "$(GREEN)Testing POST /api/users/create endpoint...$(NC)"
	curl -X POST http://localhost:$(PORT)/api/users/create \
		-H "Content-Type: application/json" \
		-d '{"username":"testuser","email":"test@example.com"}' \
		| jq '.' 2>/dev/null || curl -X POST http://localhost:$(PORT)/api/users/create \
		-H "Content-Type: application/json" \
		-d '{"username":"testuser","email":"test@example.com"}'

.PHONY: install-tools
install-tools: ## Install development tools
	@echo "$(GREEN)Installing development tools...$(NC)"
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/air-verse/air@latest
	go install github.com/go-delve/delve/cmd/dlv@latest
	@echo "$(GREEN)Tools installed!$(NC)"

# Show current configuration
show-config: ## Show current configuration (without exposing passwords)
	@echo "$(YELLOW)Current configuration:$(NC)"
	@echo "PORT: $(PORT)"
	@echo "DB_URL: $(DB_URL)"
	@echo "$(YELLOW)Migration commands available:$(NC)"
	@echo "  make migrate-create NAME=<name>  - Create new migration"
	@echo "  make migrate-up                  - Run all pending migrations"
	@echo "  make migrate-up-1                - Run next migration"
	@echo "  make migrate-down                - Rollback last migration"
	@echo "  make migrate-status              - Check migration status"
	@echo "  make db-migrate-status           - Detailed migration status"

# Default target
all: deps build

# Quick migration check
.PHONY: migration-check
migration-check: ## Quick check if migrations need to be run
	@echo "$(YELLOW)Checking if migrations are needed...$(NC)"
	@$(MAKE) -s env-check
	@if [ -d "./migrations" ] && [ "$$(ls -A ./migrations/*.sql 2>/dev/null)" ]; then \
		echo "$(GREEN)Migrations directory has files$(NC)"; \
		CURRENT_VERSION=$$($(MIGRATE) -path ./migrations -database "$(DB_URL)" version 2>&1 | grep -o '^[0-9]\+' || echo "0"); \
		LATEST_MIGRATION=$$(ls -1 ./migrations/*.up.sql 2>/dev/null | sort -n | tail -1 | grep -o '[0-9]\+' || echo "0"); \
		if [ "$$CURRENT_VERSION" -lt "$$LATEST_MIGRATION" ]; then \
			echo "$(YELLOW)Migrations pending (current: $$CURRENT_VERSION, latest: $$LATEST_MIGRATION)$(NC)"; \
			echo "Run 'make migrate-up' to apply them"; \
		else \
			echo "$(GREEN)Database is up to date (version: $$CURRENT_VERSION)$(NC)"; \
		fi \
	else \
		echo "$(YELLOW)No migration files found$(NC)"; \
	fi