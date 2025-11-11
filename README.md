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

### Environment Variables

- `HOST`: Default host for proxy requests (default: `http://postman-echo.com`)
- `PATH`: Default path for proxy requests (default: `get?foo1=bar1&foo2=bar2`)

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
