package pinghandlers

import (
	"fmt"
	"net/http"
)

type Pong struct {
	Status string `json:"status"`
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	// Return plaintext 'pong' response
	fmt.Fprintln(w, "pong")
}
