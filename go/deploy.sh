#!/bin/bash

# OOM Kill Simulator Deployment Helper
# This script helps with building and deploying the OOM simulator

set -e

PROJECT_NAME="oom-simulator"
BINARY_NAME="oom-simulator"
DOCKER_IMAGE="$PROJECT_NAME:latest"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Build Go binary
build_binary() {
    print_status "Building Go binary..."
    cd go/
    go mod tidy
    go build -o "../$BINARY_NAME" main.go
    cd ..
    print_status "Binary built successfully: $BINARY_NAME"
}

