package healthcheckhandlers

import (
	"encoding/json"
	"net/http"
	"time"
)

type Health struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

func Healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Health{
		Status: "healthy",
		Time:   time.Now().UTC().Format(time.RFC3339),
	})
}
