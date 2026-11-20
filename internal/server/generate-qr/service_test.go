package generateqr

import (
	"bytes"
	"image"
	"testing"

	"github.com/liyue201/goqr"
)

func TestGenerareQr(t *testing.T) {
	qr := MyQR{}
	var (
		png  []byte
		size = 256
		text = "testing"
	)

	err := qr.GenerateQr(&png, text, size)
	if err != nil {
		t.Errorf("error generateqr %v", err)
	}
	img, _, err := image.Decode(bytes.NewReader(png))
	if err != nil {
		t.Errorf("imgae Decode error: %v", err)
	}

	qrCodes, err := goqr.Recognize(img)
	if err != nil {
		t.Errorf("failed recognize %v", err)
	}
	if len(qrCodes) == 1 {
		if text != string(qrCodes[0].Payload) {
			t.Errorf("incorrect content")
		}
	}
}
