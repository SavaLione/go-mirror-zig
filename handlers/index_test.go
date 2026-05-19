package handlers

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRootHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		uri            string
		version        string
		expectedStatus int
		expectedBody   string
		template       *template.Template
	}{
		{
			name:           "Template and version check",
			uri:            "/",
			version:        "1.2.3",
			expectedStatus: http.StatusOK,
			expectedBody:   "<b>version 1.2.3</b>",
			template:       template.Must(template.New("index.html").Parse("<b>version {{ .Version }}</b>")),
		},
		{
			name:           "404 on wrong path",
			uri:            "/should-not-work/",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404 page not found",
			template:       template.Must(template.New("index.html").Parse("")),
		},
		{
			name:           "Empty version defaults to unknown",
			uri:            "/",
			version:        "",
			expectedStatus: http.StatusOK,
			expectedBody:   "<b>version unknown</b>",
			template:       template.Must(template.New("index.html").Parse("<b>version {{ .Version }}</b>")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest("GET", tt.uri, nil)
			rr := httptest.NewRecorder()

			handler := RootHandler(tt.template, tt.version)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("got status %v, want %v", status, tt.expectedStatus)
			}

			if body := strings.TrimSpace(rr.Body.String()); body != tt.expectedBody {
				t.Errorf("got %v, want %v", body, tt.expectedBody)
			}
		})
	}
}
