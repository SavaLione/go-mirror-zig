package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRedirectHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		hostHeader   string
		uri          string
		port         int
		wantLocation string
	}{
		{
			name:         "Simple hostname",
			hostHeader:   "example.com",
			uri:          "/path/to/resource",
			port:         8443,
			wantLocation: "https://example.com:8443/path/to/resource",
		},
		{
			name:         "Hostname with existing port",
			hostHeader:   "example.com:80",
			uri:          "/",
			port:         8443,
			wantLocation: "https://example.com:8443/",
		},
		{
			name:         "IPv6 address",
			hostHeader:   "[2001:db8::1]",
			uri:          "/test",
			port:         8443,
			wantLocation: "https://[2001:db8::1]:8443/test",
		},
		{
			name:         "Hostname with a zig artifact",
			hostHeader:   "example.com",
			uri:          "/zig-riscv64-linux-0.14.1.tar.xz",
			port:         8443,
			wantLocation: "https://example.com:8443/zig-riscv64-linux-0.14.1.tar.xz",
		},
		{
			name:         "Redirect to 443 (omit port)",
			hostHeader:   "example.com:8080",
			uri:          "/zig-riscv64-linux-0.14.1.tar.xz.minisig",
			port:         443,
			wantLocation: "https://example.com/zig-riscv64-linux-0.14.1.tar.xz.minisig",
		},
		{
			name:         "Redirect to 443 (omit port)",
			hostHeader:   "example.com:80",
			uri:          "/zig-riscv64-linux-0.14.1.tar.xz.minisig",
			port:         443,
			wantLocation: "https://example.com/zig-riscv64-linux-0.14.1.tar.xz.minisig",
		},
		{
			name:         "Redirect to custom port (include port)",
			hostHeader:   "example.com",
			uri:          "/zig-riscv64-linux-0.14.1.tar.xz.minisig",
			port:         8443,
			wantLocation: "https://example.com:8443/zig-riscv64-linux-0.14.1.tar.xz.minisig",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest("GET", tt.uri, nil)
			req.Host = tt.hostHeader

			rr := httptest.NewRecorder()
			handler := RedirectHandler(tt.port)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusMovedPermanently {
				t.Errorf("got status %v, want %v", status, http.StatusMovedPermanently)
			}

			gotLocation := rr.Header().Get("Location")
			if gotLocation != tt.wantLocation {
				t.Errorf("got Location %v, want %v", gotLocation, tt.wantLocation)
			}
		})
	}

}
