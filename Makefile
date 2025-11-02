.PHONY: build run test clean help

# Binary name
BINARY_NAME=pokedexcli

# Build directory
BUILD_DIR=.

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/pokedexcli

run: build ## Build and run the application
	./$(BUILD_DIR)/$(BINARY_NAME)

test: ## Run tests
	$(GOTEST) -v ./...

test-coverage: ## Run tests with coverage
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	$(GOCLEAN)
	rm -f $(BUILD_DIR)/$(BINARY_NAME)
	rm -f coverage.out coverage.html

fmt: ## Format code
	$(GOCMD) fmt ./...

vet: ## Run go vet
	$(GOCMD) vet ./...

lint: fmt vet ## Run linters

mod-tidy: ## Tidy go modules
	$(GOMOD) tidy

mod-download: ## Download dependencies
	$(GOMOD) download

.DEFAULT_GOAL := help

