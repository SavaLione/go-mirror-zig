package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
	"runtime"
)

func RootHandler(t *template.Template, version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		info := struct {
			Version      string
			Architecture string
		}{
			Version:      version,
			Architecture: runtime.GOARCH,
		}

		if version == "" {
			info.Version = "unknown"
		}

		if err := t.ExecuteTemplate(w, "index.html", info); err != nil {
			slog.Error(
				"failed to execute index template",
				"error", err,
				"remote_ip", GetRemoteIP(*r),
			)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
