package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nxkern/quick-qr-api/internal/server/helloworld"
	"go.opentelemetry.io/otel"
)

func TestRateLimit(t *testing.T) {
	rateLimiter := NewIPRateLimiter(5, 10)
	meter := otel.Meter("test")
	counter, _ := meter.Int64Counter("http.requests")

	hello := &helloworld.MyHandler{}

	handlerWithMiddleware := PerClientRateLimit(rateLimiter, http.HandlerFunc(hello.HelloWorld), counter)
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
