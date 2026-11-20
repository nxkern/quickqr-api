# quick-qr-api

a lightweight, high-performance go api for generating qr codes, featuring built-in rate limiting and health monitoring.

> **notice:** this open-source and public api includes opentelemetry integration (`internal/server/otel.go`) to collect traces, metrics, and logs for observability and performance monitoring.

## features

- **qr code generation:** encode text or urls into png qr codes.
- **rate limiting:** per-client ip rate limiting middleware with cloudflare and reverse proxy header support.
- **health check:** endpoint to verify api status and service availability.
- **opentelemetry:** built-in sdk setup for metrics, tracing, and logging observability.

## endpoints

### post /api/generateqr
generates a qr code png image from the json request payload.

request body:
```json
{
  "text": "https://example.com",
  "size": 256
}
```

### get /api/health
returns the status of the service and its internal checks.

response:
```json
{
  "status": "ok",
  "checks": {
    "qrgenerator": "ok"
  }
}
```

### get /
returns a hello world greeting message.

## getting started

### prerequisites
- go 1.27+

### running the api
```bash
go run main.go
```

### running tests
```bash
go test ./...
```
