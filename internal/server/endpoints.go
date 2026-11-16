package server

import (
	"fmt"
	"net/http"

	generateqr "codeberg.org/Nxus/quick-qr-api/internal/server/generate-qr"
	"codeberg.org/Nxus/quick-qr-api/internal/server/helloworld"
)


type Server interface {
	Start() error
	HandleEndpoints()
}

type MyServer struct {
	Addr    string
	Handler http.Handler
}

func (s *MyServer) Start() error {
	fmt.Printf("Server listening in http://localhost%s\n", s.Addr)
	return http.ListenAndServe(s.Addr, s.Handler)
}

func (s *MyServer) HandleFunc() {
	helloW := helloworld.MyHandler{}
	http.HandleFunc("/", helloW.HelloWorld)
	qr := generateqr.MyHandler{}
	http.HandleFunc("/api/generateqr", qr.GenerateQr)
}



