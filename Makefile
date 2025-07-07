# Build variables
BINARY_NAME=block-flow-server
BUILD_DIR=bin
MAIN_PATH=cmd/server

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOLINT=golangci-lint

# Build targets
.PHONY: all build clean test deps fmt lint run dev help install-frontend build-frontend dev-frontend

all: deps fmt lint test build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(MAIN_PATH)

# Build for production with optimizations
build-prod:
	@echo "Building $(BINARY_NAME) for production..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME) ./$(MAIN_PATH)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

# Run tests with coverage report
test-coverage: test
	@echo "Generating coverage report..."
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) -s -w .

# Lint code
lint:
	@echo "Linting code..."
	$(GOLINT) run

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Development mode with file watching (requires air)
dev:
	@echo "Starting development server..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Installing air for hot reload..."; \
		$(GOGET) -u github.com/cosmtrek/air; \
		air; \
	fi

# Install development tools
dev-tools:
	@echo "Installing development tools..."
	$(GOGET) -u github.com/cosmtrek/air
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint

# Initialize the project
init:
	@echo "Initializing project..."
	$(GOMOD) init block-flow
	make deps
	make dev-tools

# Frontend commands
install-frontend:
	@echo "Installing frontend dependencies..."
	@cd web && npm install

build-frontend:
	@echo "Building frontend..."
	@cd web && npm run build

dev-frontend:
	@echo "Starting frontend development server..."
	@cd web && npm start

# Full development setup
setup: deps install-frontend dev-tools
	@echo "Project setup complete!"

# Build everything
build-all: build build-frontend
	@echo "Building complete application..."

# Help
help:
	@echo "Available commands:"
	@echo "  build        Build the application"
	@echo "  build-prod   Build for production"
	@echo "  clean        Clean build artifacts"
	@echo "  test         Run tests"
	@echo "  test-coverage Run tests with coverage"
	@echo "  deps         Download dependencies"
	@echo "  fmt          Format code"
	@echo "  lint         Lint code"
	@echo "  run          Build and run the application"
	@echo "  dev          Run in development mode with hot reload"
	@echo "  dev-tools    Install development tools"
	@echo "  init         Initialize the project"
	@echo "  help         Show this help message"
