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

```bash
# Start the server
go run main.go

# Trigger OOM simulation
curl -X POST http://localhost:9090/trigger

# Check server status
curl http://localhost:9090/

# Force crash (for testing)
curl http://localhost:9090/crash
```