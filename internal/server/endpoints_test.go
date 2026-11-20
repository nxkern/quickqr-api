package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nxkern/quick-qr-api/internal/middlewares"
	generateqr "github.com/nxkern/quick-qr-api/internal/server/generate-qr"
)

func TestServer(t *testing.T) {
	srv := &MyServer{
		IpLimiter: middlewares.NewIPRateLimiter(5, 10),
	}

	srv.HandleEndpoints()

	ts := httptest.NewServer(srv.Handler)
	defer ts.Close()

	respHealthCheck, err := http.Get(fmt.Sprintf("%s/api/health", ts.URL))
	if err != nil {
		t.Errorf("error request healthcheck: %v", err)
	}
	defer respHealthCheck.Body.Close()
	if respHealthCheck.StatusCode != http.StatusOK {
		t.Errorf("error status code: %v", respHealthCheck.StatusCode)
	}

	bodyGenerateQr := &generateqr.Qr{
		Text: "testing ci",
		Size: 256,
	}

	bodyJson, err := json.Marshal(bodyGenerateQr)

	if err != nil {
		t.Errorf("error marshal json: %v", err)
	}

	respGenerateQr, err := http.Post(fmt.Sprintf("%s/api/generateqr", ts.URL),"application/json", bytes.NewBuffer(bodyJson))

	if err != nil {
		t.Errorf("error request generate: %v", err)
	}

	defer respGenerateQr.Body.Close()
	if respGenerateQr.StatusCode != http.StatusOK {
		t.Errorf("error status code: %v", respGenerateQr.StatusCode)
	}
}
