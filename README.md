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

The application supports configuration through both **environment variables** and **file mounts**, making it easy to deploy on platforms like Choreo.

#### Configuration Priority (highest to lowest):
1. **Environment Variables** (highest priority)
2. **Configuration File** (if `CONFIG_FILE_PATH` is set)
3. **Default Values** (hardcoded defaults)

#### Environment Variables

You can configure the application using these environment variables:

- `CONFIG_FILE_PATH`: Path to configuration file (e.g., `/config/app-config.json`)
- `SERVER_PORT`: Server port (default: `9090`)
- `LOG_LEVEL`: Logging level - `info`, `debug`, `error` (default: `info`)
- `PROXY_DEFAULT_HOST`: Default host for proxy requests (default: `http://postman-echo.com`)
- `PROXY_DEFAULT_PATH`: Default path for proxy requests (default: `get?foo1=bar1&foo2=bar2`)
- `REQUEST_TIMEOUT`: Request timeout in seconds (default: `30`)
- `ENABLE_STATUS_ENDPOINT`: Enable/disable status endpoint - `true`/`false` (default: `true`)
- `ENABLE_CONFIG_ENDPOINT`: Enable/disable config endpoint - `true`/`false` (default: `true`)

**Example using environment variables:**
```bash
export SERVER_PORT=8080
export LOG_LEVEL=debug
export PROXY_DEFAULT_HOST=https://api.example.com
go run main.go
```

#### Configuration File (File Mount)

You can provide a JSON configuration file and set its path via the `CONFIG_FILE_PATH` environment variable.

**Example `config.json`:**
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

**Example using file mount:**
```bash
export CONFIG_FILE_PATH=/config/app-config.json
go run main.go
```

#### Using with Choreo

**Option 1: Environment Variables Only**
- In Choreo UI, add key-value pairs as environment variables
- Set individual values like `SERVER_PORT=8080`, `LOG_LEVEL=debug`

**Option 2: File Mount Only**
- In Choreo UI, create a file mount at `/config/app-config.json`
- Add your JSON configuration as the file content
- Set environment variable: `CONFIG_FILE_PATH=/config/app-config.json`

**Option 3: Both (Recommended)**
- Use file mount for base configuration
- Use environment variables for overrides and secrets
- Environment variables take precedence over file configuration

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
