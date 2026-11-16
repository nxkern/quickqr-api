package generateqr

import (
	"encoding/json"
	"log"
	"net/http"
)

// Handler handles QR code HTTP requests.
type Handler interface {
	GenerateQr(w http.ResponseWriter, r *http.Request)
}

// MyHandler implements the Handler interface for QR generation.
type MyHandler struct{}

// GenerateQr calls MyQR.GenerateQr and sends the image in the response
func (h *MyHandler) GenerateQr(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusBadRequest)
		return
	}
	var (
		png []byte
		err error
	)
	services := MyQR{}
	qr := Qr{}
	err = json.NewDecoder(r.Body).Decode(&qr)
	if err != nil {
		log.Printf("[/api/generateqr] error json decoder %v", err)
		http.Error(w, "error incorrect body", http.StatusConflict)
		return
	}
	err = services.GenerateQr(&png, qr.Text, qr.Size)
	if err != nil {
		log.Printf("[api/generateqr] error %v", err)
		http.Error(w, "error encode text", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	err = json.NewEncoder(w).Encode(png)
	if err != nil {
		log.Printf("[/api/generateqr] error encode writer %v", err)
		http.Error(w, "error encode qr", http.StatusInternalServerError)
		return
	}
}
