package handlers

import (
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

		http.Redirect(w, r, targetURL, http.StatusMovedPermanently)
	}
}
