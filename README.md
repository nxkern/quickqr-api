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
  "text": "https://example.com", // content to encode in the qr
  "size": 256                   // dimension in pixels for the generated image
}
```

response:
- **content-type:** image/png
- **body:** raw binary data of the generated png image.

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

## usage

this is a standard http api that can be consumed by any client, framework, or programming language (e.g., fetch, axios, postman, or native http libraries).

### generate a qr code
example using `curl`:
```bash
curl -X POST https://quickqr-api.nxbin.xyz/api/generateqr \
   -H "Content-Type: application/json" \
   -d '{"text": "https://example.com", "size": 256}'
```

### check health status
example using `curl`:
```bash
curl https://quickqr-api.nxbin.xyz/api/health
```
