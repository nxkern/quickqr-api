// Package main is the entry point for the quick-qr-api server application.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nxkern/quick-qr-api/internal/middlewares"
	"github.com/nxkern/quick-qr-api/internal/server"
)

var watchSignals = []os.Signal{os.Interrupt, syscall.SIGUSR1}

func main() {
	rateLimiter := middlewares.NewIPRateLimiter(5, 10)
	ctx, stop := signal.NotifyContext(context.Background(), watchSignals...)
	defer stop()

	srv := &server.MyServer{
		Addr:      ":8080",
		Handler:   nil,
		IpLimiter: rateLimiter,
	}
	srv.HandleEndpoints()
	if err := srv.Start(ctx); err != nil {
		log.Fatalf("internal server error: %v", err)
	}
	log.Println("Shutdown server...")
}
