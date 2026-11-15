package server

import (
	"fmt"
	"net/http"

	api "codeberg.org/Nxus/quick-qr-api/internal/server/endpoints"
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
	endpoints :=  api.MyEndpoints{}
	endpoints.HelloWorld()
}



