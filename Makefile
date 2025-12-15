# B.DEV CLI - Build Automation
# Usage: make [target]

.PHONY: all build run test clean install release fmt lint help

# Variables
VERSION := 3.0.0
BINARY := bdev
BUILD_DIR := ./build
LDFLAGS := -ldflags "-s -w -X github.com/badie/bdev/internal/cmd/version.Version=$(VERSION)"

# Default target
all: clean build

# Build binary
build:
	@echo "Building $(BINARY) v$(VERSION)..."
	@go build $(LDFLAGS) -o $(BINARY).exe ./cmd/bdev
	@echo "Build complete: $(BINARY).exe"

# Run without building
run:
	@go run ./cmd/bdev

# Run with arguments
run-args:
	@go run ./cmd/bdev $(ARGS)

# Run tests
test:
	@echo "Running tests..."
	@go test -v -race ./...

# Run tests with coverage
test-coverage:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Benchmark
bench:
	@go test -bench=. -benchmem ./...

# Clean build artifacts
clean:
	@if exist $(BINARY).exe del $(BINARY).exe
	@if exist coverage.out del coverage.out
	@if exist coverage.html del coverage.html
	@echo "Cleaned"

# Install to GOPATH
install: build
	@copy $(BINARY).exe $(GOPATH)\bin\
	@echo "Installed to $(GOPATH)\bin\"

# Format code
fmt:
	@go fmt ./...
	@echo "Formatted"

# Get dependencies
deps:
	@go mod download
	@go mod tidy
	@echo "Dependencies updated"

# Lint code (requires golangci-lint)
lint:
	@golangci-lint run ./...

# Cross-compile for all platforms
release:
	@echo "Building releases..."
	@if not exist $(BUILD_DIR) mkdir $(BUILD_DIR)
	@set GOOS=windows&& set GOARCH=amd64&& go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-windows-amd64.exe ./cmd/bdev
	@set GOOS=linux&& set GOARCH=amd64&& go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-linux-amd64 ./cmd/bdev
	@set GOOS=darwin&& set GOARCH=amd64&& go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-darwin-amd64 ./cmd/bdev
	@set GOOS=darwin&& set GOARCH=arm64&& go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-darwin-arm64 ./cmd/bdev
	@echo "Releases built in $(BUILD_DIR)/"

# Help
help:
	@echo B.DEV CLI Makefile
	@echo.
	@echo Usage: make [target]
	@echo.
	@echo Targets:
	@echo   build          Build the binary
	@echo   run            Run without building
	@echo   test           Run tests
	@echo   test-coverage  Run tests with coverage
	@echo   clean          Clean build artifacts
	@echo   install        Install to GOPATH
	@echo   release        Cross-compile for all platforms
	@echo   fmt            Format code
	@echo   deps           Download dependencies
	@echo   lint           Lint code
	@echo   help           Show this help
