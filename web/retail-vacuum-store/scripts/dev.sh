#!/bin/bash

# Development script for vacuum store
set -e

echo "🌪️ Starting Vacuum Store Development Server"

# Kill any existing server on port 8080
if lsof -ti:8080 >/dev/null 2>&1; then
    echo "Killing existing server on port 8080..."
    lsof -ti:8080 | xargs kill -9
fi

# Go to project directory
cd "$(dirname "$0")/.."

# Install dependencies
echo "Installing dependencies..."
go mod tidy

# Run the server with development settings
echo "Starting server on http://localhost:8080"
go run ./cmd/server --log-level=debug --port=8080
