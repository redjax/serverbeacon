package middleware

import (
	"log"
	"net"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, status: http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.wroteHeader = true
	rw.ResponseWriter.WriteHeader(code)
}

func realIP(r *http.Request) string {
	// Respect X-Forwarded-For / X-Real-IP for if the app is behind a proxy
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}

	if xr := r.Header.Get("X-Real-IP"); xr != "" {
		return xr
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return host
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := wrapResponseWriter(w)

		next.ServeHTTP(wrapped, r)

		ua := r.UserAgent()
		ip := realIP(r)
		dur := time.Since(start)

		log.Printf(
			"%s \"%s %s\" %d %v UA=%q",
			ip,
			r.Method,
			r.URL.Path,
			wrapped.status,
			dur,
			ua,
		)
	})
}
