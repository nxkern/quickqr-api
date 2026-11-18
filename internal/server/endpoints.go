package server

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/nxkern/quick-qr-api/internal/middlewares"
	generateqr "github.com/nxkern/quick-qr-api/internal/server/generate-qr"
	"github.com/nxkern/quick-qr-api/internal/server/healthcheck"
	"github.com/nxkern/quick-qr-api/internal/server/helloworld"
)

type Middleware func(http.Handler) http.Handler

type Server interface {
	Start() error
	HandleEndpoints()
}

type MyServer struct {
	Addr      string
	Handler   http.Handler
	IpLimiter *middlewares.IPRateLimiter
}

func (s *MyServer) Start() error {
	fmt.Printf("Server listening in http://localhost%s\n", s.Addr)
	return http.ListenAndServe(s.Addr, s.Handler)
}

func (s *MyServer) HandleFunc() {
	helloW := helloworld.MyHandler{}

	helloHandle := Chain(http.HandlerFunc(helloW.HelloWorld), func(h http.Handler) http.Handler {
		return middlewares.PerClientRateLimit(s.IpLimiter, h)
	})
	http.Handle("/", helloHandle)

	qr := generateqr.MyHandler{}
	generateQrHandle := Chain(http.HandlerFunc(qr.GenerateQr), func(h http.Handler) http.Handler {
		return middlewares.PerClientRateLimit(s.IpLimiter, h)
	})
	http.Handle("/api/generateqr", generateQrHandle)

	health := healthcheck.MyHandler{}
	healthHandle := Chain(http.HandlerFunc(health.HealthCheck), func(h http.Handler) http.Handler {
		return middlewares.PerClientRateLimit(s.IpLimiter, h)
	})
	http.Handle("/api/health", healthHandle)
}

func Chain(handle http.Handler, middle ...Middleware) http.Handler {
	for i := range slices.Backward(middle) {
		handle = middle[i](handle)
	}
	return handle
}
