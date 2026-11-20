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
	srv.Start(ctx)
	<-ctx.Done()
	log.Println("Shutdown server...")

}
