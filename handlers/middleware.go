package handlers

import (
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"
)

// A custom response writer to capture the status code.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Limit request body size
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(recorder, r)

		slog.Info("request handled",
			"remote_ip", GetRemoteIP(*r),
			"method", r.Method,
			"path", r.URL.Path,
			"status", recorder.status,
			"duration", time.Since(start),
			"user_agent", r.UserAgent(),
		)
	})
}

func GetRemoteIP(r http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return strings.Split(forwarded, ",")[0]
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func GetSource(r http.Request) string {
	if source := r.URL.Query().Get("source"); source != "" {
		return source
	}
	return ""
}
