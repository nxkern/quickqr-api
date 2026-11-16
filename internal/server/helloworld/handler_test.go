package helloworld

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	ep := &MyHandler{}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	ep.HelloWorld(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("error status code %d", rr.Code)
	}

	waited := "hello world"

	if rr.Body.String() != waited {
		t.Errorf("error get data %s, get %s", waited, rr.Body.String())
	}
}
