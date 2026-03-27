package handler

import (
	"net/http"
	"time"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, ApiResponse{
		Ok: true,
		Data: map[string]string{
			"status":    "ok",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		},
	})
}
