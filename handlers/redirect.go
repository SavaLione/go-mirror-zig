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

		portStr := ""
		if port != 443 {
			portStr = ":" + strconv.Itoa(port)
		}

		targetURL := "https://" + host + portStr + r.URL.RequestURI()

		http.Redirect(w, r, targetURL, http.StatusMovedPermanently)
	}
}
