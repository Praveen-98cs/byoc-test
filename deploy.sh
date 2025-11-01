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

# Build Docker image
build_docker() {
    print_status "Building Docker image..."
    docker build -t "$DOCKER_IMAGE" ./go/
    print_status "Docker image built successfully: $DOCKER_IMAGE"
}

# Run locally
run_local() {
    print_status "Running locally..."
    if [ -f "$BINARY_NAME" ]; then
        ./"$BINARY_NAME"
    else
        print_error "Binary not found. Run build first."
        exit 1
    fi
}

# Run with Docker
run_docker() {
    print_status "Running with Docker..."
    docker run -p 9090:9090 "$DOCKER_IMAGE"
}

# Clean up
cleanup() {
    print_status "Cleaning up..."
    if [ -f "$BINARY_NAME" ]; then
        rm "$BINARY_NAME"
        print_status "Removed binary: $BINARY_NAME"
    fi
    
    if docker images | grep -q "$PROJECT_NAME"; then
        docker rmi "$DOCKER_IMAGE" 2>/dev/null || true
        print_status "Removed Docker image: $DOCKER_IMAGE"
    fi
}

# Show usage
show_usage() {
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  build      - Build Go binary"
    echo "  docker     - Build Docker image"
    echo "  run        - Run locally (requires build first)"
    echo "  run-docker - Run with Docker (requires docker build first)"
    echo "  clean      - Clean up binaries and images"
    echo "  help       - Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 build && $0 run"
    echo "  $0 docker && $0 run-docker"
}

# Main logic
case "${1:-help}" in
    "build")
        build_binary
        ;;
    "docker")
        build_docker
        ;;
    "run")
        run_local
        ;;
    "run-docker")
        run_docker
        ;;
    "clean")
        cleanup
        ;;
    "help"|"")
        show_usage
        ;;
    *)
        print_error "Unknown command: $1"
        show_usage
        exit 1
        ;;
esac