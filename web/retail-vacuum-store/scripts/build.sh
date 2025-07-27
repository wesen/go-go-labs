#!/bin/bash

# Build script for vacuum store
set -e

echo "🌪️ Building Vacuum Store"

# Go to project directory
cd "$(dirname "$0")/.."

# Clean previous builds
echo "Cleaning previous builds..."
rm -rf build/
mkdir -p build/

# Install dependencies
echo "Installing dependencies..."
go mod tidy

# Build the application
echo "Building application..."
go build -o build/vacuum-store ./cmd/server

# Copy static assets
echo "Copying static assets..."
cp -r static/ build/static/
cp config.yaml build/

# Copy database directory structure
mkdir -p build/db/

echo "Build complete! Binary: build/vacuum-store"
echo "To run: cd build && ./vacuum-store"
