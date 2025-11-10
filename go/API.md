# API Documentation

## OOM Kill Simulator API

### Base URL
```
http://localhost:9090
```

### Endpoints

#### GET /
**Description**: Health check and basic information about the simulator.

**Response**:
```
OOMKill Simulator
Send a POST request to /trigger to simulate OOMKill
```

#### POST /trigger
**Description**: Triggers the OOM (Out of Memory) simulation process.

**Method**: POST

**Response**:
```
OOMKill simulation triggered
```

**Notes**: 
- This endpoint starts a background goroutine that continuously allocates memory
- The simulation will continue until the system runs out of memory
- Each allocation chunk is 500MB by default

#### GET /crash
**Description**: Immediately crashes the server for testing purposes.

**Response**: Server exits with code 1

**Logs**: "crashing server..." message before exit

### Simulation Process

1. **Initialization**: Sets up memory allocation variables
2. **Chunk Allocation**: Allocates 500MB chunks in a loop
3. **Memory Filling**: Fills each chunk with non-zero data
4. **Progress Logging**: Logs allocation progress and chunk completion
5. **Continue**: Repeats until OOM kill occurs

### Log Messages

The simulator provides detailed logging:
- Server startup messages
- Simulation initialization
- Chunk size configuration
- Memory allocation progress
- Memory filling progress
- Total chunks allocated
- Request method logging

### Usage Example


## Advanced Usage

### Monitoring the Simulation

You can monitor the simulation progress by watching the server logs:

```bash
# Run with verbose logging
go run main.go 2>&1 | tee simulation.log

# In another terminal, monitor memory usage
watch -n 1 'ps aux | grep main | grep -v grep'
```


### Testing Scenarios

1. **Basic OOM Test**: Trigger simulation and wait for OOM kill
2. **Memory Limit Test**: Run with Docker memory constraints
3. **Load Test**: Multiple concurrent trigger requests
4. **Monitoring Test**: Track memory allocation patterns

### Response Codes

- `200`: Successful request
- `405`: Method not allowed (GET on /trigger)
- `500`: Server error (rare, usually before crash)
- `504`: Server is not incoming 
- `201`: Success please check

### Error Handling

The simulator handles various error conditions:
- Invalid HTTP methods
- Server overload scenarios
- Memory allocation failures



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
    
    # Clean up any test logs
    if [ -f "test_results.log" ]; then
        rm "test_results.log"
        print_status "Removed test logs"
    fi
}

# Run tests
run_tests() {
    print_status "Running automated tests..."
    if [ -f "test.sh" ]; then
        ./test.sh
    else
        print_error "Test script not found: test.sh"
        exit 1
    fi
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
    "test")
        run_tests
        ;;
    "check")
        check_requirements
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
    
    # Clean up any test logs
    if [ -f "test_results.log" ]; then
        rm "test_results.log"
        print_status "Removed test logs"
    fi
}

# Run tests
run_tests() {
    print_status "Running automated tests..."
    if [ -f "test.sh" ]; then
        ./test.sh
    else
        print_error "Test script not found: test.sh"
        exit 1
    fi
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
    "test")
        run_tests
        ;;
    "check")
        check_requirements
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
    
    # Clean up any test logs
    if [ -f "test_results.log" ]; then
        rm "test_results.log"
        print_status "Removed test logs"
    fi
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
    
    # Clean up any test logs
    if [ -f "test_results.log" ]; then
        rm "test_results.log"
        print_status "Removed test logs"
    fi
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

# Function to check if server is running
check_server() {
    echo "Checking if server is running..." | tee -a $LOG_FILE
    if curl -s "$SERVER_URL" > /dev/null 2>&1; then
        echo "✓ Server is running" | tee -a $LOG_FILE
        return 0
    else
        echo "✗ Server is not running" | tee -a $LOG_FILE
        return 1
    fi
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


