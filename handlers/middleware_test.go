package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMiddleware(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	middleware := Middleware(nextHandler)
	rr := httptest.NewRecorder()

	req := httptest.NewRequest("POST", "/", strings.NewReader(strings.Repeat("a", 1024)))
	middleware.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("got %v, want %v", rr.Code, http.StatusOK)
	}
}

func TestGetRemoteIp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		remoteAddr string
		xForwarded string
		want       string
	}{
		{"Direct connection", "1.2.3.4:5678", "", "1.2.3.4"},
		{"Forwarded", "127.0.0.1:80", "200.200.200.200, 1.1.1.1", "200.200.200.200"},
		{"Malformed remote addr", "invalid-addr", "", "invalid-addr"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := &http.Request{
				Header:     http.Header{"X-Forwarded-For": []string{tt.xForwarded}},
				RemoteAddr: tt.remoteAddr,
			}
			if got := GetRemoteIP(*req); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}

}

func TestGetSource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		url  string
		want string
	}{
		{"With source", "/?source=health-check", "health-check"},
		{"Empty source", "/", ""},
		{"Multiple params", "/?other=1&source=zig-vscode", "zig-vscode"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest("GET", tt.url, nil)
			if got := GetSource(*req); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
