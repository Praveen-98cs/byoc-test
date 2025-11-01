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

