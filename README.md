# byoc-test

This contains example user apps for choreo byoc feature

## Go Application

A lightweight HTTP service with health checks, proxy capabilities, and server monitoring.

### Getting Started

```bash
cd go
go run main.go
```

The server will start on port `9090`.

### API Endpoints

#### `GET /`
Returns the active status of the server.

**Response:**
```json
{"active": true}
```

#### `GET /healthz/`
Health check endpoint.

**Response:**
```json
{"healthy": true}
```

#### `GET /hello/`
Simple greeting endpoint.

**Query Parameters:**
- `name` (string): Name to greet

**Example:**
```bash
curl "http://localhost:9090/hello/?name=World"
```

**Response:**
```
Hello World
```

#### `GET /status/`
Returns comprehensive server metrics including memory usage, runtime information, and uptime.

**Response:**
```json
{
  "server": {
    "status": "running",
    "uptime": "2h15m30s",
    "startTime": "2025-11-11T10:30:00Z"
  },
  "memory": {
    "allocatedMB": 12,
    "totalAllocMB": 25,
    "sysMB": 50,
    "numGC": 5
  },
  "runtime": {
    "numGoroutines": 10,
    "goVersion": "go1.21.0",
    "numCPU": 8
  },
  "endpoints": ["/", "/healthz/", "/hello/", "/status/", "/proxy/"]
}
```

**Example:**
```bash
curl http://localhost:9090/status/
```

#### `GET /config/`
Returns the current application configuration, showing all active settings.

**Response:**
```json
{
  "server": {
    "port": 9090,
    "logLevel": "info"
  },
  "proxy": {
    "defaultHost": "http://postman-echo.com",
    "defaultPath": "get?foo1=bar1&foo2=bar2",
    "requestTimeoutSeconds": 30
  },
  "features": {
    "enableStatusEndpoint": true,
    "enableConfigEndpoint": true
  },
  "configSource": {
    "filePathEnvVar": "/config/app-config.json",
    "loadedFromFile": true
  }
}
```

**Example:**
```bash
curl http://localhost:9090/config/
```

#### `GET /crash/`
Terminates the server immediately to simulate a container crash. Useful for testing container orchestration, restart policies, and crash recovery.

**Response:**
```json
{"message": "Server crashing now..."}
```

**Example:**
```bash
curl http://localhost:9090/crash/
```

**Note:** The server will exit with code 1 immediately after responding. This is intended for testing purposes only.

#### `POST /proxy/`
Proxies requests to external APIs.

**Request Body:**
```json
{
  "host": "http://postman-echo.com",
  "path": "get?foo1=bar1&foo2=bar2"
}
```

**Example:**
```bash
curl -X POST http://localhost:9090/proxy/ \
  -H "Content-Type: application/json" \
  -d '{
    "host": "http://postman-echo.com",
    "path": "get?foo1=bar1&foo2=bar2"
  }'
```

### Configuration

The application uses a **simplified configuration system** with only **3 environment variables** and **1 file mount**, keeping the Choreo UI clean and minimal.

#### Configuration Priority (highest to lowest):
1. **Environment Variables** (highest priority - only PORT and LOG)
2. **Configuration File** (if `CONFIG` is set)
3. **Default Values** (hardcoded defaults)

#### Environment Variables (Limited to 3)

The application supports **only 3 environment variables** (max 8 characters each):

- `PORT`: Server port (default: `9090`)
- `LOG`: Logging level - `info`, `debug`, `error` (default: `info`)
- `CONFIG`: Path to configuration file (default: none)

**Example:**
```bash
export PORT=8080
export LOG=debug
export CONFIG=/config/app-config.json
go run main.go
```

#### Configuration File (File Mount)

**All other settings** must be configured through a JSON file mount. This includes proxy settings, timeouts, and feature flags.

**Example `app-config.json`:**
```json
{
  "server": {
    "port": 9090,
    "logLevel": "info"
  },
  "proxy": {
    "defaultHost": "http://postman-echo.com",
    "defaultPath": "get?foo1=bar1&foo2=bar2",
    "requestTimeoutSeconds": 30
  },
  "features": {
    "enableStatusEndpoint": true,
    "enableConfigEndpoint": true
  }
}
```

**To use a config file:**
```bash
export CONFIG=/config/app-config.json
go run main.go
```

#### Using with Choreo

For the cleanest Choreo UI experience, use:

**Environment Variables** (3 only):
- `PORT`: `8080` (or your desired port)
- `LOG`: `info` (or `debug` for verbose logging)
- `CONFIG`: `/config/app-config.json` (path to your file mount)

**File Mount** (1 only):
- Mount Path: `/config/app-config.json`
- Content: Your JSON configuration (as shown above)

This approach keeps the Choreo configuration UI minimal while maintaining full configurability through the file mount

### Building

```bash
cd go
go build -o server main.go
./server
```

### Docker

```bash
cd go
docker build -t byoc-test .
docker run -p 9090:9090 byoc-test
```

### OpenAPI Specification

The API is documented using OpenAPI 3.0. See `go/openapi.yaml` for the complete specification.
