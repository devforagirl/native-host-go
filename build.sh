#!/bin/bash

# Exit on error
set -e

# Create dist directory
mkdir -p dist

echo "Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -o dist/flowmeter-host-win.exe main.go

echo "Building for macOS (arm64)..."
GOOS=darwin GOARCH=arm64 go build -o dist/flowmeter-host-macos main.go

echo "Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -o dist/flowmeter-host-linux main.go

echo "Build complete. Binaries are in dist/"
