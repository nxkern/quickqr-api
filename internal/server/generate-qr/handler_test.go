package generateqr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/liyue201/goqr"
)

func TestHandlerGenerateQr(t *testing.T) {
	ep := MyHandler{}
	body := &Qr{
		Text: "testing",
		Size: 256,
	}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		t.Errorf("error marshal json %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/generateqr", bytes.NewBuffer(bodyJson))
	rr := httptest.NewRecorder()
	ep.GenerateQr(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("error status code: %v", rr.Code)
	}
	var img []byte
	if err := json.NewDecoder(rr.Body).Decode(&img); err != nil {
		t.Errorf("error decode json %v", err)
	}
	if err := helper(img, body.Text); err != nil {
		t.Errorf("error data %v", err)
	}
}

func helper(png []byte, text string) error {
	img, _, err := image.Decode(bytes.NewReader(png))
	if err != nil {
		return fmt.Errorf("imgae Decode error: %v", err)
	}

	qrCodes, err := goqr.Recognize(img)
	if err != nil {
		return fmt.Errorf("failed recognize %v", err)
	}
	if len(qrCodes) == 1 {

		if text != string(qrCodes[0].Payload) {
			return fmt.Errorf("incorrect content")
		}
	}
	return nil
}
