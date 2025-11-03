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
