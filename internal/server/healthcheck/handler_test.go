package healthcheck

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealhCheck(t *testing.T) {
	ep := MyHandler{}
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rr := httptest.NewRecorder()
	ep.HealthCheck(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("error status code: %d", rr.Code)
	}
}
