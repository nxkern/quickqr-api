// Package healthcheck provides health check types and HTTP handlers for service status monitoring.
package healthcheck

// Health represents the health status response payload.
type Health struct {
	Status string     `json:"status"`
	Checks ChecksList `json:"checks"`
}

// ChecksList details status of individual system dependencies.
type ChecksList struct {
	QRGenerator string `json:"qrgenerator"`
}
