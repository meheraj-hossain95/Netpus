#!/bin/bash

echo "========================================"
echo "Building Octopus Network Monitor"
echo "========================================"
echo

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed or not in PATH"
    echo "Please install Go 1.22 or higher from https://go.dev/dl/"
    exit 1
fi

# Check if Wails is installed
if ! command -v wails &> /dev/null; then
    echo "Error: Wails CLI is not installed"
    echo "Installing Wails CLI..."
    go install github.com/wailsapp/wails/v2/cmd/wails@latest
    if [ $? -ne 0 ]; then
        echo "Failed to install Wails CLI"
        exit 1
    fi
    echo "Wails CLI installed successfully"
    echo
    
    # Add Go bin to PATH if not already there
    export PATH="$PATH:$(go env GOPATH)/bin"
fi

echo "Downloading dependencies..."
go mod download
if [ $? -ne 0 ]; then
    echo "Failed to download dependencies"
    exit 1
fi
echo

echo "Building application..."
wails build
if [ $? -ne 0 ]; then
    echo
    echo "Build failed!"
    exit 1
fi

echo
echo "========================================"
echo "Build successful!"
echo "========================================"
echo
echo "Executable location: build/bin/octopus"
echo
echo "To install shortcuts, run:"
echo "  ./build/bin/octopus --install"
echo
