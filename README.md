# OOMKill Simulator

A lightweight Go service for simulating Out-of-Memory (OOM) conditions and server crashes.

## Quick Start

```bash
go build -o oomkill-simulator main.go
./oomkill-simulator
```

## API Endpoints

- `GET /` - Health check
- `POST /trigger` - Start OOM simulation (allocates memory in 500MB chunks until crash)
- `GET /crash` - Force immediate server exit

## Environment Variables

- `SHOULD_CRASH` - Reserved for future crash control features (currently logged but not enforced)

## Log Levels

The service uses structured logging with the following prefixes:
- `[STARTUP]` - Server initialization events
- `[INFO]` - General information (health checks, requests)
- `[WARN]` - Warnings (invalid requests)
- `[ACTION]` - Active operations (OOM simulation, crash triggers)
- `[MEMORY]` - Memory allocation progress
- `[CRITICAL]` - Critical events (crashes)
- `[ERROR]` - Error conditions

## Usage Examples

```bash
# Start the server
./oomkill-simulator

# Trigger OOM simulation
curl -X POST http://localhost:9090/trigger

# Force crash
curl http://localhost:9090/crash

# Health check
curl http://localhost:9090/
```

Server runs on port `9090`. Use for testing container resilience and monitoring systems.
