package main

import (
	"github.com/nxkern/quick-qr-api/internal/middlewares"
	"github.com/nxkern/quick-qr-api/internal/server"
)

func main() {
	rateLimiter := middlewares.NewIPRateLimiter(5, 10)
	srv := &server.MyServer{
		Addr:      ":8080",
		Handler:   nil,
		IpLimiter: rateLimiter,
	}
	srv.HandleFunc()
	srv.Start()
}
