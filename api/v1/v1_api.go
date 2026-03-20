package v1

import (
	"net/http"

	healthcheckhandlers "github.com/redjax/serverbeacon/internal/handlers/healthcheckHandlers"
	pinghandlers "github.com/redjax/serverbeacon/internal/handlers/pingHandlers"
	"github.com/redjax/serverbeacon/internal/middleware"
)

func New() http.Handler {
	mux := http.NewServeMux()

	// API version
	v1 := http.NewServeMux()

	v1.HandleFunc("GET /ping", pinghandlers.Ping)
	v1.HandleFunc("GET /health", healthcheckhandlers.Healthcheck)

	mux.Handle("/v1/", http.StripPrefix("/v1", middleware.Logging(v1)))

	return mux
}
