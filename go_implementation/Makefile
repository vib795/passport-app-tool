.PHONY: build clean run install test help

# Default target
.DEFAULT_GOAL := help

# Binary name
BINARY_NAME=passport-tracker

# Build directory
BUILD_DIR=build

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags="-s -w"

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	$(GOBUILD) -o $(BINARY_NAME) -v

build-all: ## Build for all platforms
	@mkdir -p $(BUILD_DIR)
	@echo "Building for Linux..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux
	@echo "Building for macOS (AMD64)..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-macos-amd64
	@echo "Building for macOS (ARM64)..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-macos-arm64
	@echo "Building for Windows..."
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME).exe
	@echo "Done! Binaries are in $(BUILD_DIR)/"

clean: ## Remove build artifacts
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)

run: ## Run the application (requires -r and -d flags)
	$(GOCMD) run main.go $(ARGS)

install: ## Install dependencies
	$(GOMOD) download
	$(GOMOD) tidy

test: ## Run tests
	$(GOTEST) -v ./...

deps: ## Display dependencies
	$(GOMOD) graph

verify: ## Verify dependencies
	$(GOMOD) verify

fmt: ## Format code
	$(GOCMD) fmt ./...

vet: ## Run go vet
	$(GOCMD) vet ./...

# Example usage targets
example: ## Run with example data
	$(GOCMD) run main.go -r 25-2055287305 -d 12/09/1994

# Development helpers
dev: install fmt vet build ## Prepare for development (install, format, vet, build)

release: clean install test build-all ## Create release builds
