package helloworld

import (
	"io"
	"log"
	"net/http"
)

type Handler interface {
	HelloWorld(w http.ResponseWriter, r *http.Request)
}

type MyHandler struct{}

func (h *MyHandler) HelloWorld(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusBadRequest)
		return
	}
	if _, err := io.WriteString(w, "hello world"); err != nil {
		log.Printf("[api/helloworld] error: %v\n", err)
		return
	}
}
