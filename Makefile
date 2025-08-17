# Makefile for Chatmate
# Provides common development tasks and build automation

.PHONY: help build test clean install release dev-setup doc lint security format check all

# Default target
.DEFAULT_GOAL := help

# Project configuration
BINARY_NAME := chatmate
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DIR := dist
INSTALL_DIR := /usr/local/bin

# Colors for output
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

# Helper to print colored output
define log
	@echo -e "$(BLUE)▶$(NC) $(1)"
endef

define success
	@echo -e "$(GREEN)✓$(NC) $(1)"
endef

define warning
	@echo -e "$(YELLOW)⚠$(NC) $(1)"
endef

define error
	@echo -e "$(RED)✗$(NC) $(1)"
endef

help: ## Show this help message
	@echo -e "$(BLUE)Chatmate Development Tasks$(NC)"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo ""
	@echo -e "$(GREEN)Common Tasks:$(NC)"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ { printf "  $(BLUE)%-15s$(NC) %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo ""
	@echo -e "$(GREEN)Examples:$(NC)"
	@echo "  make build              # Build for current platform"
	@echo "  make test               # Run full test suite"
	@echo "  make release VERSION=v1.2.3   # Create release"
	@echo "  make install            # Install to $(INSTALL_DIR)"
	@echo "  make dev-setup          # Set up development environment"
	@echo ""

build: ## Build binary for current platform
	$(call log,"Building $(BINARY_NAME) for current platform")
	@./scripts/build.sh --release
	$(call success,"Build completed")

build-all: ## Build binaries for all platforms
	$(call log,"Building $(BINARY_NAME) for all platforms")
	@./scripts/build.sh --platforms --release
	$(call success,"Cross-platform build completed")

test: ## Run comprehensive test suite
	$(call log,"Running comprehensive Go test suite")
	@go test ./... -cover
	$(call success,"All tests passed")

test-unit: ## Run unit tests only
	$(call log,"Running unit tests")
	@go test ./cmd ./internal/manager ./pkg/...
	$(call success,"Unit tests completed")

test-integration: ## Run integration tests only
	$(call log,"Running integration tests")
	@go test ./test/integration/... -v
	$(call success,"Integration tests completed")

test-coverage: ## Run tests with detailed coverage report
	$(call log,"Generating test coverage report")
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	$(call success,"Coverage report generated: coverage.html")

clean: ## Clean build artifacts and caches
	$(call log,"Cleaning build artifacts")
	@rm -f $(BINARY_NAME)*
	@rm -rf $(BUILD_DIR)
	@go clean -cache -testcache -modcache
	$(call success,"Clean completed")

install: build ## Install binary to system
	$(call log,"Installing $(BINARY_NAME) to $(INSTALL_DIR)")
	@sudo cp ./$(BINARY_NAME) $(INSTALL_DIR)/
	@sudo chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	$(call success,"Installed to $(INSTALL_DIR)/$(BINARY_NAME)")

uninstall: ## Remove binary from system
	$(call log,"Uninstalling $(BINARY_NAME) from $(INSTALL_DIR)")
	@sudo rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	$(call success,"Uninstalled $(BINARY_NAME)")

dev-setup: ## Set up development environment
	$(call log,"Setting up development environment")
	@go mod download
	@go mod tidy
	@chmod +x scripts/*.sh
	@chmod +x run-tests.sh
	$(call success,"Development environment ready")

format: ## Format Go code
	$(call log,"Formatting Go code")
	@go fmt ./...
	$(call success,"Code formatting completed")

lint: ## Run code linters
	$(call log,"Running linters")
	@go vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		$(warning,"golangci-lint not installed, skipping advanced linting"); \
	fi
	$(call success,"Linting completed")

security: ## Run security scan
	$(call log,"Running security scan")
	@./$(BINARY_NAME) status --security-scan 2>/dev/null || echo "Building binary first..."
	@make build >/dev/null 2>&1
	@./$(BINARY_NAME) status --security-scan
	$(call success,"Security scan completed")

doc: ## Generate documentation
	$(call log,"Generating documentation")
	@go doc -all ./... > docs/API.md 2>/dev/null || true
	@./scripts/generate-man-pages.go 2>/dev/null || echo "Man pages already generated"
	$(call success,"Documentation generated")

completions: ## Install shell completions
	$(call log,"Installing shell completions")
	@./scripts/install-completions.sh
	$(call success,"Shell completions installed")

release: ## Create a new release (specify VERSION=vX.Y.Z)
ifndef VERSION
	$(error Please specify VERSION, e.g., make release VERSION=v1.2.3)
endif
	$(call log,"Creating release $(VERSION)")
	@./scripts/release.sh $(VERSION)
	$(call success,"Release $(VERSION) created")

release-patch: ## Create a patch release
	$(call log,"Creating patch release")
	@./scripts/release.sh patch
	$(call success,"Patch release created")

release-minor: ## Create a minor release  
	$(call log,"Creating minor release")
	@./scripts/release.sh minor
	$(call success,"Minor release created")

release-major: ## Create a major release
	$(call log,"Creating major release")
	@./scripts/release.sh major
	$(call success,"Major release created")

dry-run: ## Show what a patch release would do
	$(call log,"Dry run for patch release")
	@./scripts/release.sh --dry-run patch

check: lint test security ## Run all quality checks
	$(call success,"All quality checks passed")

all: clean format lint test build ## Full development cycle
	$(call success,"Full development cycle completed")

# Advanced targets
benchmark: ## Run performance benchmarks
	$(call log,"Running benchmarks")
	@go test -bench=. -benchmem ./...
	$(call success,"Benchmarks completed")

coverage: ## Generate test coverage report
	$(call log,"Generating coverage report")
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"
	$(call success,"Coverage report generated")

profile: ## Profile the application
	$(call log,"Profiling application")
	@go test -cpuprofile=cpu.prof -memprofile=mem.prof ./...
	@echo "Profiles: cpu.prof, mem.prof"
	$(call success,"Profiling completed")

deps: ## Update dependencies
	$(call log,"Updating dependencies")
	@go get -u ./...
	@go mod tidy
	$(call success,"Dependencies updated")

# Development shortcuts
dev: format lint test ## Quick development check
	$(call success,"Development checks passed")

quick: ## Quick build and test
	@make build test-unit
	$(call success,"Quick check completed")

# Platform-specific builds
linux: ## Build for Linux
	@GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64
	$(call success,"Linux build completed")

windows: ## Build for Windows  
	@GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows-amd64.exe
	$(call success,"Windows build completed")

macos: ## Build for macOS
	@GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin-amd64
	@GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME)-darwin-arm64  
	$(call success,"macOS builds completed")

# Utility targets
version: ## Show current version
	@echo "Version: $(VERSION)"
	@echo "Git commit: $(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"
	@echo "Build date: $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")"

size: ## Show binary sizes
	$(call log,"Binary sizes")
	@ls -lh $(BINARY_NAME)* 2>/dev/null | awk '{print $$9 "\t" $$5}' || echo "No binaries found"
	@if [ -d "$(BUILD_DIR)" ]; then \
		echo "Distribution binaries:"; \
		ls -lh $(BUILD_DIR)/* | awk '{print $$9 "\t" $$5}'; \
	fi

# Maintenance targets  
reset: clean ## Reset to clean state
	$(call warning,"Resetting to clean state")
	@git clean -fd
	@git reset --hard
	$(call success,"Reset completed")

# Help for make targets
list: ## List all make targets
	@$(MAKE) -pRrq -f $(firstword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'
