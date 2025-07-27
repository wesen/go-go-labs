#!/bin/bash

# Test script for vacuum store
set -e

echo "🌪️ Running Vacuum Store Tests"

# Go to project directory
cd "$(dirname "$0")/.."

# Install dependencies
echo "Installing dependencies..."
go mod tidy

# Run tests
echo "Running unit tests..."
go test ./... -v

# Run race condition tests
echo "Running race condition tests..."
go test ./... -race

# Check formatting
echo "Checking Go formatting..."
if ! gofmt -l . | grep -q '^$'; then
    echo "Code not formatted. Run: go fmt ./..."
    gofmt -l .
    exit 1
fi

# Run vet
echo "Running go vet..."
go vet ./...

echo "All tests passed! ✅"
