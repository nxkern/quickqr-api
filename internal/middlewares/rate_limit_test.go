package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"codeberg.org/Nxus/quick-qr-api/internal/server/helloworld"
)

func TestRateLimit(t *testing.T) {
	rateLimiter := NewIPRateLimiter(5, 10)
	hello := &helloworld.MyHandler{}

	handlerWithMiddleware := PerClientRateLimit(rateLimiter, http.HandlerFunc(hello.HelloWorld))
	for i := range 15 {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		rr := httptest.NewRecorder()
		handlerWithMiddleware.ServeHTTP(rr, req)
		if i < 10 {
			if rr.Code != http.StatusOK {
				t.Errorf("error code %d", rr.Code)
			}
		} else {
			if rr.Code != http.StatusTooManyRequests {
				t.Errorf("error code %d in request: %d", rr.Code, i)
			}
		}
	}
}
