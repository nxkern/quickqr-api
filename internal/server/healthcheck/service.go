package healthcheck

import generateqr "codeberg.org/Nxus/quick-qr-api/internal/server/generate-qr"

// HealthService defines operations for performing health checks on services.
type HealthService interface {
	Check(health *Health) 
}

// MyHealthService implements the HealthService interface.
type MyHealthService struct {}

// Check verifies system dependencies and updates health status accordingly.
func (h *MyHealthService) Check(health *Health) {
		// QR generator
	var (
		qr = generateqr.MyQR{}
		png []byte
		test = "health"
		size = 256
	)

	if err := qr.GenerateQr(&png, test, size); err != nil {
		health.Status = "error"
		health.Checks.QRGenerator = "down"
	}

}
