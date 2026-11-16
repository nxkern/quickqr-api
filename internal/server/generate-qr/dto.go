// Package generateqr provides HTTP handlers and services for QR code generation.
package generateqr

// Qr defines the request payload for QR code generation.
type Qr struct {
	Text string `json:"text"`
	Size int    `json:"size"`
}
