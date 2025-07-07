# Block-Flow API Documentation

## Overview

The Block-Flow API provides RESTful endpoints for managing flows and blocks in the visual programming platform.

**Base URL:** `http://localhost:8080/api/v1`

## Authentication

Currently, no authentication is required. This will be added in future versions.

## Error Handling

All endpoints return appropriate HTTP status codes:

- `200 OK` - Successful operation
- `201 Created` - Resource created successfully
- `204 No Content` - Successful operation with no response body
- `400 Bad Request` - Invalid request data
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

Error responses include a JSON object with an error message:
```json
{
  "error": "Description of the error"
}
```

## Endpoints

### Health Check

#### GET /health

Check the health status of the API.

**Response:**
```json
{
  "status": "ok"
}
```

### Flows

#### GET /flows

List all flows.

**Response:**
```json
[
  {
    "id": "flow-123",
    "name": "My Flow",
    "description": "A sample flow",
    "nodes": [...],
    "connections": [...],
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z",
    "version": "1.0.0",
    "active": false
  }
]
```

#### POST /flows

Create a new flow.

**Request Body:**
```json
{
  "name": "My New Flow",
  "description": "Description of the flow",
  "nodes": [],
  "connections": [],
  "properties": {}
}
```

**Response:** `201 Created`
```json
{
  "id": "generated-flow-id",
  "name": "My New Flow",
  "description": "Description of the flow",
  "nodes": [],
  "connections": [],
  "properties": {},
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z",
  "version": "1.0.0",
  "active": false
}
```

#### GET /flows/{id}

Get a specific flow by ID.

**Parameters:**
- `id` (string) - Flow ID

**Response:**
```json
{
  "id": "flow-123",
  "name": "My Flow",
  "description": "A sample flow",
  "nodes": [...],
  "connections": [...],
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z",
  "version": "1.0.0",
  "active": false
}
```

#### PUT /flows/{id}

Update an existing flow.

**Parameters:**
- `id` (string) - Flow ID

**Request Body:**
```json
{
  "name": "Updated Flow Name",
  "description": "Updated description",
  "nodes": [...],
  "connections": [...]
}
```

**Response:**
```json
{
  "id": "flow-123",
  "name": "Updated Flow Name",
  "description": "Updated description",
  "nodes": [...],
  "connections": [...],
  "updated_at": "2025-01-01T00:00:00Z"
}
```

#### DELETE /flows/{id}

Delete a flow.

**Parameters:**
- `id` (string) - Flow ID

**Response:** `204 No Content`

#### POST /flows/{id}/start

Start execution of a flow.

**Parameters:**
- `id` (string) - Flow ID

**Response:**
```json
{
  "status": "started"
}
```

#### POST /flows/{id}/stop

Stop execution of a flow.

**Parameters:**
- `id` (string) - Flow ID

**Response:**
```json
{
  "status": "stopped"
}
```

#### POST /flows/{id}/trigger

Manually trigger a flow with optional input data.

**Parameters:**
- `id` (string) - Flow ID

**Request Body (optional):**
```json
{
  "payload": "any data",
  "topic": "optional topic",
  "headers": {
    "key": "value"
  }
}
```

**Response:**
```json
{
  "status": "triggered"
}
```

#### GET /flows/{id}/status

Get the execution status of a flow.

**Parameters:**
- `id` (string) - Flow ID

**Response:**
```json
{
  "id": "execution-123",
  "flow_id": "flow-123",
  "status": "running",
  "started_at": "2025-01-01T00:00:00Z",
  "nodes": {
    "node-1": {
      "node_id": "node-1",
      "status": "success",
      "executed_at": "2025-01-01T00:00:01Z",
      "duration": 150000000,
      "input_count": 1,
      "output_count": 1
    }
  }
}
```

### Blocks

#### GET /blocks

List all available block types.

**Response:**
```json
[
  {
    "type": "inject",
    "name": "Inject",
    "description": "Manually trigger the flow",
    "category": "input",
    "version": "1.0.0",
    "author": "Block-Flow",
    "icon": "play",
    "color": "#4CAF50"
  },
  {
    "type": "debug",
    "name": "Debug",
    "description": "Output debug information",
    "category": "output",
    "version": "1.0.0",
    "author": "Block-Flow",
    "icon": "bug",
    "color": "#FF9800"
  }
]
```

#### GET /blocks/{type}

Get detailed information about a specific block type.

**Parameters:**
- `type` (string) - Block type

**Response:**
```json
{
  "type": "inject",
  "name": "Inject",
  "description": "Manually trigger the flow",
  "category": "input",
  "version": "1.0.0",
  "author": "Block-Flow",
  "icon": "play",
  "color": "#4CAF50"
}
```

## WebSocket API

### Connection

Connect to the WebSocket endpoint for real-time updates:

**Endpoint:** `ws://localhost:8080/api/v1/ws`

### Message Format

All WebSocket messages use JSON format:

```json
{
  "type": "message_type",
  "data": {
    "key": "value"
  },
  "timestamp": "2025-01-01T00:00:00Z"
}
```

### Message Types

#### Flow Execution Updates

Receive real-time updates about flow execution:

```json
{
  "type": "flow_execution",
  "data": {
    "flow_id": "flow-123",
    "execution_id": "exec-456",
    "status": "running",
    "node_id": "node-1",
    "event": "node_started"
  }
}
```

#### Debug Messages

Receive debug output from flows:

```json
{
  "type": "debug",
  "data": {
    "flow_id": "flow-123",
    "node_id": "debug-1",
    "message": "Hello, World!",
    "timestamp": "2025-01-01T00:00:00Z"
  }
}
```

## Flow JSON Format

### Flow Structure

```json
{
  "id": "string",
  "name": "string",
  "description": "string (optional)",
  "nodes": [
    {
      "id": "string",
      "type": "string",
      "name": "string",
      "x": 0,
      "y": 0,
      "properties": {},
      "inputs": 1,
      "outputs": 1,
      "wires": [[]]
    }
  ],
  "connections": [
    {
      "id": "string",
      "source": "string",
      "source_port": 0,
      "target": "string",
      "target_port": 0,
      "label": "string (optional)"
    }
  ],
  "properties": {},
  "created_at": "ISO8601 timestamp",
  "updated_at": "ISO8601 timestamp",
  "version": "string",
  "active": false
}
```

### Node Types

#### Inject Node
```json
{
  "type": "inject",
  "properties": {
    "topic": "string",
    "payload": "any",
    "repeat": "interval (optional)",
    "once": true
  }
}
```

#### Debug Node
```json
{
  "type": "debug",
  "properties": {
    "console": true,
    "complete": "payload",
    "target": "debug"
  }
}
```

#### Function Node
```json
{
  "type": "function",
  "properties": {
    "func": "JavaScript code string",
    "outputs": 1
  }
}
```

## Examples

### Creating a Simple Flow

```bash
curl -X POST http://localhost:8080/api/v1/flows \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Hello World Flow",
    "nodes": [
      {
        "id": "inject-1",
        "type": "inject",
        "name": "Start",
        "x": 100,
        "y": 100,
        "properties": {
          "payload": "Hello, World!"
        },
        "outputs": 1,
        "wires": [["debug-1"]]
      },
      {
        "id": "debug-1",
        "type": "debug",
        "name": "Output",
        "x": 300,
        "y": 100,
        "properties": {
          "console": true
        },
        "inputs": 1
      }
    ]
  }'
```

### Starting a Flow

```bash
curl -X POST http://localhost:8080/api/v1/flows/your-flow-id/start
```

### Triggering a Flow

```bash
curl -X POST http://localhost:8080/api/v1/flows/your-flow-id/trigger \
  -H "Content-Type: application/json" \
  -d '{
    "payload": "Custom message",
    "topic": "test"
  }'
```
