package zig

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchAllReleases(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/download/index.json", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		w.Write([]byte(`{
			"master": {
				"version": "0.17.0-dev.305+bdfbf432d",
				"date": "2026-05-15",
    			"docs": "https://ziglang.org/documentation/master/",
    			"stdDocs": "https://ziglang.org/documentation/master/std/",
				"src": {
				  "tarball": "https://ziglang.org/builds/zig-0.17.0-dev.305+bdfbf432d.tar.xz",
				  "shasum": "f4e02500223c65225cb98651d32f744ee1f8939db05c4718598f1d89eddaa5dd",
				  "size": "22530120"
				}
			},
            "0.16.0": {
              "version": "0.16.0",
              "date": "2026-04-13",
              "docs": "https://ziglang.org/documentation/0.16.0/",
              "stdDocs": "https://ziglang.org/documentation/0.16.0/std/",
              "notes": "https://ziglang.org/download/0.16.0/release-notes.html",
              "src": {
                "tarball": "https://ziglang.org/download/0.16.0/zig-0.16.0.tar.xz",
                "shasum": "43186959edc87d5c7a1be7b7d2a25efffd22ce5807c7af99067f86f99641bfdf",
                "size": "22503260"
              },
              "bootstrap": {
                "tarball": "https://ziglang.org/download/0.16.0/zig-bootstrap-0.16.0.tar.xz",
                "shasum": "2a8266a4205772ef40838c8cbdf14875855a515ff3adf89b49c2d2ae93613d10",
                "size": "55245980"
              }
			}
		}`))
	})

	mux.HandleFunc("/download/bad-status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`Server Error`))
	})

	mux.HandleFunc("/download/invalid-json.json", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "master": INVALID JSON ]}`))
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	tests := []struct {
		name        string
		url         string
		ctx         func() context.Context
		expectError bool
	}{
		{
			name:        "successful fetch",
			url:         ts.URL + "/download/index.json",
			ctx:         func() context.Context { return context.Background() },
			expectError: false,
		},

		{
			name:        "server returns non 200 status",
			url:         ts.URL + "/download/bad-status",
			ctx:         func() context.Context { return context.Background() },
			expectError: true,
		},
		{
			name:        "server returns an invalid json",
			url:         ts.URL + "/download/invalid-json.json",
			ctx:         func() context.Context { return context.Background() },
			expectError: true,
		},
		{
			name:        "invalid url format",
			url:         "://0.0.0.0",
			ctx:         func() context.Context { return context.Background() },
			expectError: true,
		},
		{
			name: "canceled context",
			url:  ts.URL + "/download/index.json",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := tt.ctx()
			zr, err := FetchAllReleases(ctx, tt.url)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected an error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("did not expect an error, but got: %v", err)
			}

			if len(zr) == 0 {
				t.Fatalf("expected parsed ZigReleases, got empty or nil")
			}
		})
	}
}
