// Package server provides HTTP server configuration, routing, and lifecycle management.
package server

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/nxkern/quick-qr-api/internal/middlewares"
	generateqr "github.com/nxkern/quick-qr-api/internal/server/generate-qr"
	"github.com/nxkern/quick-qr-api/internal/server/healthcheck"
	"github.com/nxkern/quick-qr-api/internal/server/helloworld"
)

// Middleware defines a function signature for wrapping an http.Handler.
type Middleware func(http.Handler) http.Handler

// Server defines the interface for managing an HTTP server lifecycle and route setup.
type Server interface {
	Start(ctx context.Context) error
	HandleEndpoints()
}

// MyServer implements the Server interface for managing the HTTP application server.
type MyServer struct {
	Addr      string
	Handler   http.Handler
	IpLimiter *middlewares.IPRateLimiter
}

// Start listens for incoming requests and gracefully shuts down when ctx is done.
func (s *MyServer) Start(ctx context.Context) error {
	srv := http.Server{
		Addr:    s.Addr,
		Handler: s.Handler,
	}
	errChan := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Go(func() {
		fmt.Printf("Server listening in http://localhost%s\n", s.Addr)
		if err := srv.ListenAndServe(); err != nil {
			errChan <- err
		}
	})
	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()

		return srv.Shutdown(shutdownCtx)
	}
}

// HandleEndpoints initializes routes and registers middleware handlers.
func (s *MyServer) HandleEndpoints() {
	mux := http.NewServeMux()

	helloW := helloworld.MyHandler{}
	helloHandle := Chain(http.HandlerFunc(helloW.HelloWorld), func(h http.Handler) http.Handler {
		return middlewares.PerClientRateLimit(s.IpLimiter, h)
	})
	mux.Handle("/", helloHandle)

	qr := generateqr.MyHandler{}
	generateQrHandle := Chain(http.HandlerFunc(qr.GenerateQr), func(h http.Handler) http.Handler {
		return middlewares.PerClientRateLimit(s.IpLimiter, h)
	})
	mux.Handle("/api/generateqr", generateQrHandle)

	health := healthcheck.MyHandler{}
	healthHandle := Chain(http.HandlerFunc(health.HealthCheck), func(h http.Handler) http.Handler {
		return middlewares.PerClientRateLimit(s.IpLimiter, h)
	})
	mux.Handle("/api/health", healthHandle)
	s.Handler = mux
}

// Chain applies multiple Middleware decorators to an http.Handler.
func Chain(handle http.Handler, middle ...Middleware) http.Handler {
	for i := range slices.Backward(middle) {
		handle = middle[i](handle)
	}
	return handle
}
