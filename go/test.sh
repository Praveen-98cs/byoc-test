#!/bin/bash

# OOM Kill Simulator Test Script
# This script tests the various endpoints of the OOM simulator

set -e



# Test root endpoint
test_root_endpoint() {
    echo "" | tee -a $LOG_FILE
    echo "Testing root endpoint..." | tee -a $LOG_FILE
    
    response=$(curl -s "$SERVER_URL/")
    if [[ $response == *"OOMKill Simulator"* ]]; then
        echo "✓ Root endpoint test passed" | tee -a $LOG_FILE
    else
        echo "✗ Root endpoint test failed" | tee -a $LOG_FILE
        echo "Response: $response" | tee -a $LOG_FILE
    fi
}

# Test trigger endpoint (without actually triggering)
test_trigger_endpoint_method() {
    echo "" | tee -a $LOG_FILE
    echo "Testing trigger endpoint with GET method (should fail)..." | tee -a $LOG_FILE
    
    response=$(curl -s "$SERVER_URL/trigger")
    if [[ $response == *"Method not allowed"* ]]; then
        echo "✓ GET method correctly rejected" | tee -a $LOG_FILE
    else
        echo "✗ GET method test failed" | tee -a $LOG_FILE
        echo "Response: $response" | tee -a $LOG_FILE
    fi
}

# Test POST to trigger (warning: this will actually start OOM simulation)
test_trigger_endpoint_warning() {
    echo "" | tee -a $LOG_FILE
    echo "⚠️  WARNING: The next test would trigger actual OOM simulation" | tee -a $LOG_FILE
    echo "   Skipping POST /trigger test to avoid system impact" | tee -a $LOG_FILE
    echo "   To test manually: curl -X POST $SERVER_URL/trigger" | tee -a $LOG_FILE
}



# Test server response time
test_response_time() {
    echo "" | tee -a $LOG_FILE
    echo "Testing server response time..." | tee -a $LOG_FILE
    
    response_time=$(curl -o /dev/null -s -w "%{time_total}" "$SERVER_URL/")
    echo "Response time: ${response_time}s" | tee -a $LOG_FILE
    
    # Check if response time is reasonable (less than 1 second)
    if (( $(echo "$response_time < 1.0" | bc -l) )); then
        echo "✓ Response time is acceptable" | tee -a $LOG_FILE
    else
        echo "⚠️  Response time is slow: ${response_time}s" | tee -a $LOG_FILE
    fi
}

# Main test execution
main() {
    echo "Starting automated tests..." | tee -a $LOG_FILE
    
    if check_server; then
        test_root_endpoint
        test_trigger_endpoint_method
        test_trigger_endpoint_warning
        test_invalid_endpoints
        test_response_time
    else
        echo "Cannot run tests - server is not running" | tee -a $LOG_FILE
        echo "Start the server with: go run main.go" | tee -a $LOG_FILE
        exit 1
    fi
    
    echo "" | tee -a $LOG_FILE
    echo "Tests completed at $(date)" | tee -a $LOG_FILE
    echo "Results saved to: $LOG_FILE" | tee -a $LOG_FILE
}

# Run tests
main

    fi
    
    if docker images | grep -q "$PROJECT_NAME"; then
        docker rmi "$DOCKER_IMAGE" 2>/dev/null || true
        print_status "Removed Docker image: $DOCKER_IMAGE"
    fi
    
    # Clean up any test logs
    if [ -f "test_results.log" ]; then
        rm "test_results.log"
        print_status "Removed test logs"
    fi
}


# Run with Docker
run_docker() {
    print_status "Running with Docker..."
    docker run -p 9090:9090 "$DOCKER_IMAGE"
}

# Check requirements
check_requirements() {
    print_status "Checking requirements..."
    if ! command -v go &> /dev/null; then               
        print_error "Go is not installed. Please install Go to proceed."
        exit 1
    fi  

    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed. Please install Docker to proceed."
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