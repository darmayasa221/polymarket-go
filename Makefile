# Go Base Framework Makefile
# ===========================

# Variables
APP_NAME := go-base-framework
MAIN_PATH := ./cmd/server
BINARY_PATH := ./bin/server
GO := go
GOFLAGS := -v
LDFLAGS := -s -w

# Database migration variables
MIGRATE_PATH := ./db/migrations
SEED_PATH := ./db/seeds
DB ?= sqlite
SQLITE_URL := sqlite3://./data/local.db
POSTGRES_URL := $(DATABASE_URL)
MYSQL_URL := $(DATABASE_URL)

# Colors for output
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

.PHONY: all build build-prod run dev test test-cover test-short test-race lint fmt vet clean help
.PHONY: docker-build docker-up docker-down
.PHONY: migrate-up migrate-down migrate-create migrate-status seed
.PHONY: mock install-tools hooks deps verify tidy check

# Default target
all: help

## help: Display this help message
help:
	@echo "$(GREEN)$(APP_NAME) - Available Commands$(NC)"
	@echo ""
	@echo "$(YELLOW)Development:$(NC)"
	@echo "  make dev           - Run with hot reload (air)"
	@echo "  make run           - Run the application"
	@echo "  make build         - Build the application (development)"
	@echo "  make build-prod    - Build optimized production binary"
	@echo "  make clean         - Clean build artifacts"
	@echo ""
	@echo "$(YELLOW)Testing:$(NC)"
	@echo "  make test          - Run all tests"
	@echo "  make test-cover    - Run tests with coverage"
	@echo "  make test-short    - Run short tests only"
	@echo "  make test-race     - Run tests with race detector"
	@echo ""
	@echo "$(YELLOW)Code Quality:$(NC)"
	@echo "  make lint          - Run golangci-lint"
	@echo "  make fmt           - Format code with gofmt and goimports"
	@echo "  make vet           - Run go vet"
	@echo "  make check         - Run all quality checks (fmt, vet, lint, test)"
	@echo ""
	@echo "$(YELLOW)Database:$(NC)"
	@echo "  make migrate-up db=<db>         - Run migrations (db: sqlite, postgres, mysql, supabase)"
	@echo "  make migrate-down db=<db>       - Rollback one migration"
	@echo "  make migrate-status db=<db>     - Show migration status"
	@echo "  make migrate-create db=<db> name=<name> - Create new migration"
	@echo "  make seed db=<db>               - Run all seeds for database"
	@echo "  make seed db=<db> name=<name>   - Run specific seed (e.g., name=users)"
	@echo ""
	@echo "$(YELLOW)Docker:$(NC)"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-up     - Start Docker containers"
	@echo "  make docker-down   - Stop Docker containers"
	@echo ""
	@echo "$(YELLOW)Tools:$(NC)"
	@echo "  make install-tools - Install development tools"
	@echo "  make hooks         - Install Git hooks (lefthook)"
	@echo "  make mock          - Generate mocks using mockery"
	@echo "  make deps          - Download dependencies"
	@echo "  make tidy          - Tidy go modules"
	@echo "  make verify        - Verify dependencies"
	@echo ""

## build: Build the application (development)
build:
	@echo "$(GREEN)Building $(APP_NAME) (development)...$(NC)"
	@mkdir -p bin
	$(GO) build $(GOFLAGS) -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "$(GREEN)Build complete: $(BINARY_PATH)$(NC)"

## build-prod: Build optimized production binary
build-prod:
	@echo "$(GREEN)Building $(APP_NAME) (production)...$(NC)"
	@mkdir -p bin
	CGO_ENABLED=0 $(GO) build -ldflags="$(LDFLAGS)" -trimpath -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "$(GREEN)Production build complete: $(BINARY_PATH)$(NC)"
	@ls -lh $(BINARY_PATH)

## run: Run the application
run:
	@echo "$(GREEN)Running $(APP_NAME)...$(NC)"
	$(GO) run $(MAIN_PATH)

## dev: Run with hot reload using air
dev:
	@echo "$(GREEN)Starting development server with hot reload...$(NC)"
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "$(RED)air is not installed. Run 'make install-tools' first$(NC)"; \
		exit 1; \
	fi

## test: Run all tests
test:
	@echo "$(GREEN)Running tests...$(NC)"
	$(GO) test ./... -v

## test-short: Run short tests only
test-short:
	@echo "$(GREEN)Running short tests...$(NC)"
	$(GO) test ./... -v -short

## test-cover: Run tests with coverage
test-cover:
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	@mkdir -p coverage
	$(GO) test ./... -v -coverprofile=coverage/coverage.out -covermode=atomic
	$(GO) tool cover -html=coverage/coverage.out -o coverage/coverage.html
	@echo "$(GREEN)Coverage report: coverage/coverage.html$(NC)"

## lint: Run golangci-lint
lint:
	@echo "$(GREEN)Running linter...$(NC)"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run --fix; \
	else \
		echo "$(RED)golangci-lint is not installed. Run 'make install-tools' first$(NC)"; \
		exit 1; \
	fi

## fmt: Format code
fmt:
	@echo "$(GREEN)Formatting code...$(NC)"
	$(GO) fmt ./...
	@if command -v goimports > /dev/null; then \
		goimports -w .; \
	fi

## vet: Run go vet
vet:
	@echo "$(GREEN)Running go vet...$(NC)"
	$(GO) vet ./...

## clean: Clean build artifacts
clean:
	@echo "$(GREEN)Cleaning...$(NC)"
	@rm -rf bin/
	@rm -rf coverage/
	@rm -rf tmp/
	@rm -rf data/*.db
	@echo "$(GREEN)Clean complete$(NC)"

## install-tools: Install development tools
install-tools:
	@echo "$(GREEN)Installing development tools...$(NC)"
	$(GO) install github.com/air-verse/air@latest
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install golang.org/x/tools/cmd/goimports@latest
	$(GO) install github.com/vektra/mockery/v2@latest
	$(GO) install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	$(GO) install github.com/evilmartians/lefthook@latest
	@echo "$(YELLOW)Note: For Supabase, install CLI manually: https://supabase.com/docs/guides/cli$(NC)"
	@echo "$(GREEN)Tools installed successfully$(NC)"
	@lefthook install
	@echo "$(GREEN)Git hooks installed$(NC)"

## hooks: Install Git hooks using lefthook
hooks:
	@echo "$(GREEN)Installing Git hooks with lefthook...$(NC)"
	@if command -v lefthook > /dev/null; then \
		lefthook install; \
	else \
		echo "$(RED)lefthook is not installed. Run 'make install-tools' first$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)Git hooks installed successfully$(NC)"

## mock: Generate mocks using mockery
mock:
	@echo "$(GREEN)Generating mocks...$(NC)"
	@if command -v mockery > /dev/null; then \
		mockery --all --dir=./internal/domains --output=./tests/mocks --outpkg=mocks; \
		mockery --all --dir=./internal/applications/security --output=./tests/mocks --outpkg=mocks; \
	else \
		echo "$(RED)mockery is not installed. Run 'make install-tools' first$(NC)"; \
		exit 1; \
	fi

# Helper function to get database URL
define get_db_url
$(if $(filter sqlite,$(1)),$(SQLITE_URL),\
$(if $(filter postgres,$(1)),$(POSTGRES_URL),\
$(if $(filter mysql,$(1)),$(MYSQL_URL),\
$(error Unknown database: $(1). Use: sqlite, postgres, or mysql))))
endef

## migrate-up: Run database migrations (usage: make migrate-up db=sqlite)
migrate-up:
ifeq ($(DB),supabase)
	@if ! command -v supabase > /dev/null; then \
		echo "$(RED)supabase CLI is not installed. Run 'make install-tools' first$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)Running supabase migrations...$(NC)"
	@supabase db push
else
	@if ! command -v migrate > /dev/null; then \
		echo "$(RED)migrate is not installed. Run 'make install-tools' first$(NC)"; \
		exit 1; \
	fi
	@if [ ! -d "$(MIGRATE_PATH)/$(DB)" ]; then \
		echo "$(RED)No migrations found for $(DB) at $(MIGRATE_PATH)/$(DB)$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)Running $(DB) migrations...$(NC)"
	@migrate -path $(MIGRATE_PATH)/$(DB) -database "$(call get_db_url,$(DB))" up
endif

## migrate-down: Rollback database migrations (usage: make migrate-down db=sqlite)
migrate-down:
ifeq ($(DB),supabase)
	@if ! command -v supabase > /dev/null; then \
		echo "$(RED)supabase CLI is not installed. Run 'make install-tools' first$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)Rolling back supabase migration...$(NC)"
	@supabase db reset
else
	@if ! command -v migrate > /dev/null; then \
		echo "$(RED)migrate is not installed. Run 'make install-tools' first$(NC)"; \
		exit 1; \
	fi
	@if [ ! -d "$(MIGRATE_PATH)/$(DB)" ]; then \
		echo "$(RED)No migrations found for $(DB) at $(MIGRATE_PATH)/$(DB)$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)Rolling back $(DB) migration...$(NC)"
	@migrate -path $(MIGRATE_PATH)/$(DB) -database "$(call get_db_url,$(DB))" down 1
endif

## migrate-status: Show migration status (usage: make migrate-status db=sqlite)
migrate-status:
ifeq ($(DB),supabase)
	@if ! command -v supabase > /dev/null; then \
		echo "$(RED)supabase CLI is not installed. Run 'make install-tools' first$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)Supabase migration status:$(NC)"
	@supabase db diff
else
	@if ! command -v migrate > /dev/null; then \
		echo "$(RED)migrate is not installed. Run 'make install-tools' first$(NC)"; \
		exit 1; \
	fi
	@if [ ! -d "$(MIGRATE_PATH)/$(DB)" ]; then \
		echo "$(RED)No migrations found for $(DB) at $(MIGRATE_PATH)/$(DB)$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)Migration status for $(DB):$(NC)"
	@migrate -path $(MIGRATE_PATH)/$(DB) -database "$(call get_db_url,$(DB))" version
endif

## migrate-create: Create a new migration file (usage: make migrate-create db=sqlite name=add_users)
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "$(RED)Please provide a migration name: make migrate-create db=<db> name=<name>$(NC)"; \
		exit 1; \
	fi
ifeq ($(DB),supabase)
	@if ! command -v supabase > /dev/null; then \
		echo "$(RED)supabase CLI is not installed. Run 'make install-tools' first$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)Creating supabase migration: $(name)$(NC)"
	@supabase migration new $(name)
else
	@if ! command -v migrate > /dev/null; then \
		echo "$(RED)migrate is not installed. Run 'make install-tools' first$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)Creating $(DB) migration: $(name)$(NC)"
	@mkdir -p $(MIGRATE_PATH)/$(DB)
	@migrate create -ext sql -dir $(MIGRATE_PATH)/$(DB) -seq $(name)
endif

## seed: Run database seeds (usage: make seed db=sqlite OR make seed db=sqlite name=users)
seed:
	@if [ ! -d "$(SEED_PATH)/$(DB)" ]; then \
		echo "$(RED)No seeds found for $(DB) at $(SEED_PATH)/$(DB)$(NC)"; \
		exit 1; \
	fi
ifeq ($(DB),supabase)
	@if ! command -v supabase > /dev/null; then \
		echo "$(RED)supabase CLI is not installed$(NC)"; \
		exit 1; \
	fi
ifdef name
	@echo "$(GREEN)Running supabase seed: $(name)$(NC)"
	@file=$$(ls $(SEED_PATH)/$(DB)/*_$(name).sql 2>/dev/null | head -1); \
	if [ -n "$$file" ]; then \
		echo "$(YELLOW)Seeding: $$(basename $$file)$(NC)"; \
		supabase db execute --file $$file; \
	else \
		echo "$(RED)Seed not found: *_$(name).sql$(NC)"; \
		exit 1; \
	fi
else
	@echo "$(GREEN)Running all supabase seeds...$(NC)"
	@for file in $$(ls $(SEED_PATH)/$(DB)/*.sql 2>/dev/null | sort); do \
		echo "$(YELLOW)Seeding: $$(basename $$file)$(NC)"; \
		supabase db execute --file $$file; \
	done
endif
else ifeq ($(DB),sqlite)
ifdef name
	@echo "$(GREEN)Running sqlite seed: $(name)$(NC)"
	@file=$$(ls $(SEED_PATH)/$(DB)/*_$(name).sql 2>/dev/null | head -1); \
	if [ -n "$$file" ]; then \
		echo "$(YELLOW)Seeding: $$(basename $$file)$(NC)"; \
		sqlite3 ./data/local.db < $$file; \
	else \
		echo "$(RED)Seed not found: *_$(name).sql$(NC)"; \
		exit 1; \
	fi
else
	@echo "$(GREEN)Running all sqlite seeds...$(NC)"
	@for file in $$(ls $(SEED_PATH)/$(DB)/*.sql 2>/dev/null | sort); do \
		echo "$(YELLOW)Seeding: $$(basename $$file)$(NC)"; \
		sqlite3 ./data/local.db < $$file; \
	done
endif
else ifeq ($(DB),postgres)
ifdef name
	@echo "$(GREEN)Running postgres seed: $(name)$(NC)"
	@file=$$(ls $(SEED_PATH)/$(DB)/*_$(name).sql 2>/dev/null | head -1); \
	if [ -n "$$file" ]; then \
		echo "$(YELLOW)Seeding: $$(basename $$file)$(NC)"; \
		psql $(DATABASE_URL) -f $$file; \
	else \
		echo "$(RED)Seed not found: *_$(name).sql$(NC)"; \
		exit 1; \
	fi
else
	@echo "$(GREEN)Running all postgres seeds...$(NC)"
	@for file in $$(ls $(SEED_PATH)/$(DB)/*.sql 2>/dev/null | sort); do \
		echo "$(YELLOW)Seeding: $$(basename $$file)$(NC)"; \
		psql $(DATABASE_URL) -f $$file; \
	done
endif
else ifeq ($(DB),mysql)
ifdef name
	@echo "$(GREEN)Running mysql seed: $(name)$(NC)"
	@file=$$(ls $(SEED_PATH)/$(DB)/*_$(name).sql 2>/dev/null | head -1); \
	if [ -n "$$file" ]; then \
		echo "$(YELLOW)Seeding: $$(basename $$file)$(NC)"; \
		mysql < $$file; \
	else \
		echo "$(RED)Seed not found: *_$(name).sql$(NC)"; \
		exit 1; \
	fi
else
	@echo "$(GREEN)Running all mysql seeds...$(NC)"
	@for file in $$(ls $(SEED_PATH)/$(DB)/*.sql 2>/dev/null | sort); do \
		echo "$(YELLOW)Seeding: $$(basename $$file)$(NC)"; \
		mysql < $$file; \
	done
endif
else
	@echo "$(RED)Unknown database: $(DB). Use: sqlite, postgres, mysql, or supabase$(NC)"
	@exit 1
endif

## docker-build: Build Docker image
docker-build:
	@echo "$(GREEN)Building Docker image...$(NC)"
	docker build -t $(APP_NAME):latest .

## docker-up: Start Docker containers
docker-up:
	@echo "$(GREEN)Starting Docker containers...$(NC)"
	docker-compose up -d

## docker-down: Stop Docker containers
docker-down:
	@echo "$(GREEN)Stopping Docker containers...$(NC)"
	docker-compose down

## deps: Download dependencies
deps:
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	$(GO) mod download

## verify: Verify dependencies
verify:
	@echo "$(GREEN)Verifying dependencies...$(NC)"
	$(GO) mod verify

## tidy: Tidy go modules
tidy:
	@echo "$(GREEN)Tidying modules...$(NC)"
	$(GO) mod tidy

## test-race: Run tests with race detector
test-race:
	@echo "$(GREEN)Running tests with race detector...$(NC)"
	$(GO) test ./... -v -race

## check: Run all quality checks
check: fmt vet lint test
	@echo "$(GREEN)All checks passed$(NC)"
