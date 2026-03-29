#!/bin/bash

# Setup script for Go Inventory System
# Author: nielwyn

set -e

echo "========================================="
echo "Go Inventory System - Setup Script"
echo "========================================="
echo ""

# Check if Go is installed
echo "Checking Go installation..."
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed"
    echo "Please install Go 1.21 or higher from https://golang.org/dl/"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
echo "✅ Go is installed: $GO_VERSION"
echo ""

# Check if PostgreSQL is accessible
echo "Checking PostgreSQL..."
if command -v psql &> /dev/null; then
    echo "✅ PostgreSQL client is installed"
else
    echo "⚠️  Warning: PostgreSQL client not found"
    echo "Please ensure PostgreSQL is installed and accessible"
fi
echo ""

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "Creating .env file from .env.example..."
    cp .env.example .env
    echo "✅ .env file created"
    echo "⚠️  Please edit .env file and update the configuration values"
else
    echo "✅ .env file already exists"
fi
echo ""

# Install Go dependencies
echo "Installing Go dependencies..."
go mod download
go mod verify
echo "✅ Dependencies installed successfully"
echo ""

# Build the application
echo "Building the application..."
go build -o bin/api ./cmd/api
echo "✅ Application built successfully"
echo ""

echo "========================================="
echo "Setup completed successfully! 🎉"
echo "========================================="
echo ""
echo "Next steps:"
echo "1. Edit .env file with your configuration"
echo "2. Ensure PostgreSQL is running"
echo "3. Run 'make run' or './bin/api' to start the server"
echo "4. Or run 'make docker-run' to start with Docker Compose"
echo ""
