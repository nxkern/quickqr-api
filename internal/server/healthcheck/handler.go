package healthcheck

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Handler handles health check HTTP requests.
type Handler interface {
	HealthCheck(w http.ResponseWriter, r *http.Request)
}

// MyHandler implements the Handler interface.
type MyHandler struct{}

// HealthCheck executes health checks and responds with current status JSON.
func (h *MyHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println("[/api/health] method no alllowed")
		http.Error(w, "method no alllowed", http.StatusBadRequest)
		return
	}
	fmt.Println("fmt healthcheck pass")
	log.Println("log healthcheck pass")
	healthService := MyHealthService{}
	health := &Health{
		Status: "ok",
		Checks: ChecksList{
			QRGenerator: "ok",
		},
	}
	healthService.Check(health)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(health); err != nil {
		log.Printf("[/api/health] error encode struct: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
