#!/bin/bash

# OOM Kill Simulator Test Script
# This script tests the various endpoints of the OOM simulator

set -e

SERVER_URL="http://localhost:9090"
LOG_FILE="test_results.log"

echo "=== OOM Kill Simulator Test Script ===" | tee $LOG_FILE
echo "Starting tests at $(date)" | tee -a $LOG_FILE
echo "" | tee -a $LOG_FILE

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

# Main test execution
main() {
    echo "Starting automated tests..." | tee -a $LOG_FILE
    
    if check_server; then
        test_root_endpoint
        test_trigger_endpoint_method
        test_trigger_endpoint_warning
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