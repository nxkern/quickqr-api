package main

import "codeberg.org/Nxus/quick-qr-api/internal/server"

func main() {
	srv := &server.MyServer{
		Addr:    ":8080",
		Handler: nil,
	}
	srv.HandleFunc()
	srv.Start()
}
