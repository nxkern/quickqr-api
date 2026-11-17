package main

import (
	"codeberg.org/Nxus/quick-qr-api/internal/middlewares"
	"codeberg.org/Nxus/quick-qr-api/internal/server"
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
