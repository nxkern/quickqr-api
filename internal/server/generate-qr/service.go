package generateqr

import (
	"fmt"

	"github.com/skip2/go-qrcode"
)

// QR defines the service interface for generating QR codes.
type QR interface {
	GenerateQr(png *[]byte) error
}

// MyQR implements the QR interface for QR code generation.
type MyQR struct{}

// GenerateQr generates QR code and return the image
func (qr *MyQR) GenerateQr(png *[]byte, text string, size int) error {
	var err error

	*png, err = qrcode.Encode(text, qrcode.Medium, size)
	if err != nil {
		return fmt.Errorf("[/api/generateqr] error enconde text: %v", err)
	}
	return nil
}
