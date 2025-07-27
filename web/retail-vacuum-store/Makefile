# Vacuum Store Makefile

.PHONY: dev build test clean run help

# Default target
help:
	@echo "🌪️ Vacuum Store Commands:"
	@echo "  make dev    - Start development server"
	@echo "  make build  - Build production binary"
	@echo "  make test   - Run tests and checks"
	@echo "  make clean  - Clean build artifacts"
	@echo "  make run    - Run production binary"
	@echo "  make deps   - Install dependencies"

# Development server
dev:
	@./scripts/dev.sh

# Build production binary
build:
	@./scripts/build.sh

# Run tests
test:
	@./scripts/test.sh

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf build/
	@go clean

# Run production binary (assumes it's already built)
run:
	@echo "🌪️ Running production build..."
	@cd build && ./vacuum-store

# Install/update dependencies
deps:
	@echo "📦 Installing dependencies..."
	@go mod tidy
	@go mod download

# Initialize database (for future use)
init-db:
	@echo "🗄️ Initializing database..."
	@mkdir -p db/
	@echo "Database directory created"

# Format code
fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...

# Lint code (requires golangci-lint)
lint:
	@echo "🔍 Linting code..."
	@golangci-lint run -v

# Quick development setup
setup: deps init-db
	@echo "✅ Development environment ready!"
	@echo "Run 'make dev' to start the server"
