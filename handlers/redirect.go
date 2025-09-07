package handlers

import (
	"log/slog"
	"net"
	"net/http"
	"strconv"
)

func RedirectHandler(port int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// net.SplitHostPort to safely extract the hostname without the port.
		host, _, err := net.SplitHostPort(r.Host)
		if err != nil {
			host = r.Host
		}

		targetURL := "https://" + host + ":" + strconv.Itoa(port) + r.URL.RequestURI()

		slog.Info(
			"redirecting http to https",
			"remote_ip", GetRemoteIP(*r),
			"method", r.Method,
			"host", r.Host,
			"path", r.URL.Path,
			"target", targetURL,
		)

		http.Redirect(w, r, targetURL, http.StatusMovedPermanently)
	}
}
