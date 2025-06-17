# Contributing to MCP Observability

This is a simple REST API that serves telemetry data using Go and Gin framework.

## Prerequisites

- Go 1.21 or later
- Git

## Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd mcp-observability
```

2. Install dependencies:
```bash
go mod download
```

3. Run the server:
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Endpoint

### Get Spans

Retrieves spans based on query parameters:
- **All spans:** No query parameters
- **By span ID:** `?spanId=...`
- **By trace ID:** `?traceId=...`

#### Get all spans
```bash
curl -X GET http://localhost:8080/spans | jq
```

#### Get a span by spanId
```bash
curl -X GET "http://localhost:8080/spans?spanId=span1" | jq
```

#### Get all spans for a traceId
```bash
curl -X GET "http://localhost:8080/spans?traceId=trace123" | jq
```

#### Advanced Query Examples

You can filter spans by any field using query parameters. Use dot notation for nested fields and attributes.

- **By name:**
```bash
curl -X GET "http://localhost:8080/spans?name=Database%20Query" | jq
```

- **By kind:**
```bash
curl -X GET "http://localhost:8080/spans?kind=CLIENT" | jq
```

- **By status code:**
```bash
curl -X GET "http://localhost:8080/spans?status.code=0" | jq
```

- **By status message:**
```bash
curl -X GET "http://localhost:8080/spans?status.message=Timeout" | jq
```

- **By attribute (e.g., http.method):**
```bash
curl -X GET "http://localhost:8080/spans?attributes.http.method=GET" | jq
```

- **By boolean attribute (e.g., cache.hit):**
```bash
curl -X GET "http://localhost:8080/spans?attributes.cache.hit=true" | jq
```

- **By custom attribute (e.g., job.name):**
```bash
curl -X GET "http://localhost:8080/spans?attributes.job.name=send-email" | jq
```

- **Combining multiple filters:**
```bash
curl -X GET "http://localhost:8080/spans?kind=INTERNAL&attributes.job.status=success" | jq
```

#### Example responses

**Single span:**
```json
{
  "traceId": "trace123",
  "spanId": "span1",
  "name": "HTTP GET /api/users",
  "kind": "SERVER",
  "startTime": "2024-03-20T10:00:00Z",
  "endTime": "2024-03-20T10:00:01Z",
  "attributes": {
    "http.method": "GET",
    "http.url": "/api/users",
    "http.status_code": 200
  },
  "status": {
    "code": 0,
    "message": "OK"
  },
  "events": [
    {
      "time": "2024-03-20T10:00:00.5Z",
      "name": "processing",
      "attributes": {
        "processing.time": "500ms"
      }
    }
  ],
  "links": []
}
```

**Multiple spans (trace):**
```json
[
  {
    "traceId": "trace123",
    "spanId": "span1",
    "name": "HTTP GET /api/users",
    "kind": "SERVER",
    "startTime": "2024-03-20T10:00:00Z",
    "endTime": "2024-03-20T10:00:01Z",
    "attributes": {
      "http.method": "GET",
      "http.url": "/api/users",
      "http.status_code": 200
    },
    "status": {
      "code": 0,
      "message": "OK"
    },
    "events": [
      {
        "time": "2024-03-20T10:00:00.5Z",
        "name": "processing",
        "attributes": {
          "processing.time": "500ms"
        }
      }
    ],
    "links": []
  },
  {
    "traceId": "trace123",
    "spanId": "span2",
    "name": "Database Query",
    "kind": "CLIENT",
    "startTime": "2024-03-20T10:00:00.1Z",
    "endTime": "2024-03-20T10:00:00.8Z",
    "attributes": {
      "db.system": "postgresql",
      "db.operation": "SELECT",
      "db.statement": "SELECT * FROM users"
    },
    "status": {
      "code": 0,
      "message": "OK"
    },
    "events": [],
    "links": [
      {
        "traceId": "trace123",
        "spanId": "span1",
        "attributes": {
          "link.type": "child_of"
        }
      }
    ]
  }
]
```

## Error Responses

The API returns appropriate HTTP status codes:

- 200: Success
- 404: Resource not found
- 500: Internal server error

Error responses include a message explaining the error:

```json
{
  "error": "Span not found"
}
``` 