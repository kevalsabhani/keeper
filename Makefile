# --- Colors and Formatting ---
GREEN  := $(shell printf "\033[32m")
YELLOW := $(shell printf "\033[33m")
BLUE   := $(shell printf "\033[34m")
CYAN   := $(shell printf "\033[36m")
RESET  := $(shell printf "\033[0m")

.PHONY: help build run dev test test-coverage lint fmt clean docker-build docker-up docker-down migrate-up migrate-down deps

help: ## Show this help message
	@echo "$(YELLOW)Usage:$(RESET) make $(GREEN)<target>$(RESET)"
	@echo ""
	@echo "$(YELLOW)Targets:$(RESET)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-20s$(RESET) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

BINARY_NAME=api
GO_FILES=$(shell find . -name "*.go")
MIGRATION_PATH=migrations

build: ## Compile the API binary
	@echo "$(CYAN)🚀 Building $(BINARY_NAME)...$(RESET)"
	@mkdir -p bin
	@CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/$(BINARY_NAME) cmd/api/main.go
	@echo "$(GREEN)✅ Build complete!$(RESET)"

run: build ## Build and run the API
	@echo "$(CYAN)🏃‍♂️ Starting $(BINARY_NAME)...$(RESET)"
	@./bin/$(BINARY_NAME)

dev: ## Run the API in development mode
	@echo "$(YELLOW)🔄 Starting in development mode...$(RESET)"
	@go run cmd/api/main.go

test: ## Run tests
	@echo "$(CYAN)🧪 Running tests...$(RESET)"
	@go test -v -race -coverprofile=coverage.out ./...
	@echo "$(GREEN)✅ Tests complete!$(RESET)"

test-coverage: test ## View test coverage in browser
	@echo "$(CYAN)📊 Generating coverage report...$(RESET)"
	@go tool cover -html=coverage.out

lint: ## Run linter
	@echo "$(CYAN)🧹 Running linter...$(RESET)"
	@golangci-lint run
	@echo "$(GREEN)✅ Linting complete!$(RESET)"

fmt: ## Format source code
	@echo "$(CYAN)📝 Formatting code...$(RESET)"
	@goimports -w .
	@echo "$(GREEN)✅ Code formatted!$(RESET)"

clean: ## Clean build files and coverage reports
	@echo "$(YELLOW)🗑️  Cleaning up...$(RESET)"
	@rm -rf bin/ coverage.out
	@echo "$(GREEN)✅ Clean complete!$(RESET)"

docker-build: ## Build Docker image
	@echo "$(CYAN)🐳 Building Docker image...$(RESET)"
	@docker build -t keeper-api:latest .
	@echo "$(GREEN)✅ Docker image built!$(RESET)"

docker-up: ## Start services using docker-compose
	@echo "$(CYAN)🐳 Starting Docker services...$(RESET)"
	@docker-compose up
	@echo "$(GREEN)✅ Services started!$(RESET)"

docker-down: ## Stop services using docker-compose
	@echo "$(YELLOW)🐳 Stopping Docker services...$(RESET)"
	@docker-compose down
	@echo "$(GREEN)✅ Services stopped!$(RESET)"

migrate-up: ## Apply database migrations
	@echo "$(CYAN)🆙 Applying database migrations...$(RESET)"
	@migrate -path ${MIGRATION_PATH} -database "${DB_URL}" -verbose up
	@echo "$(GREEN)✅ Migrations applied!$(RESET)"

migrate-down: ## Rollback database migrations
	@echo "$(YELLOW)⏬ Rolling back database migrations...$(RESET)"
	@migrate -path ${MIGRATION_PATH} -database "${DB_URL}" -verbose down
	@echo "$(GREEN)✅ Rollback complete!$(RESET)"

deps: ## Download dependencies
	@echo "$(CYAN)📦 Downloading dependencies...$(RESET)"
	@go mod tidy
	@go mod verify
	@echo "$(GREEN)✅ Dependencies ready!$(RESET)"