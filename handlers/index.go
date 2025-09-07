package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
)

func RootHandler(t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		if err := t.ExecuteTemplate(w, "index.html", nil); err != nil {
			slog.Error(
				"failed to execute index template",
				"error", err,
				"remote_ip", GetRemoteIP(*r),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
