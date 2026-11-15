package endpoints

import (
	"io"
	"net/http"
)

type Endpoints interface {
	HelloWorld()
} 

type MyEndpoints struct {}


func (s *MyEndpoints) HelloWorld() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "holis")
	})
}


