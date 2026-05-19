package config

import (
	"flag"
	"testing"
)

func TestHTTPAddress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		in       Config
		expected string
	}{
		{
			name: "Localhost",
			in: Config{
				ListenAddress: "127.0.0.1",
				HTTPPort:      80,
			},
			expected: "127.0.0.1:80",
		},
		{
			name: "All IPs",
			in: Config{
				ListenAddress: "0.0.0.0",
				HTTPPort:      8080,
			},
			expected: "0.0.0.0:8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.in.HTTPAddress(); got != tt.expected {
				t.Errorf("got %v, want %v", got, tt.expected)
			}
		})
	}

}

func TestHTTPSAddress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		in       Config
		expected string
	}{
		{
			name: "Localhost",
			in: Config{
				ListenAddress: "127.0.0.1",
				TLSPort:       443,
			},
			expected: "127.0.0.1:443",
		},
		{
			name: "All IPs",
			in: Config{
				ListenAddress: "0.0.0.0",
				TLSPort:       443,
			},
			expected: "0.0.0.0:443",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.in.HTTPSAddress(); got != tt.expected {
				t.Errorf("got %v, want %v", got, tt.expected)
			}
		})
	}

}

func TestAcceptTOS(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		in             Config
		inTOSurl       string
		expectedTOS    bool
		expectedTOSMsg string
	}{
		{
			name:           "TOS accepted",
			in:             Config{ACMEAcceptTOS: true},
			inTOSurl:       "lets-encrypt",
			expectedTOS:    true,
			expectedTOSMsg: "Accepting ACME Terms of Service at: lets-encrypt",
		},
		{
			name:           "TOS isn't accepted",
			in:             Config{ACMEAcceptTOS: false},
			expectedTOS:    false,
			expectedTOSMsg: "Terms of Service are not accepted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.in.AcceptTOS(tt.inTOSurl)

			if tt.in.ACMEAcceptTOS != tt.expectedTOS {
				t.Errorf("got %v, want %v", tt.in.ACMEAcceptTOS, tt.expectedTOS)
			}

			if got := tt.in.acmeTOSURL; got != tt.expectedTOSMsg {
				t.Errorf("got %v, want %v", got, tt.expectedTOSMsg)
			}
		})
	}

}

func TestParseConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		args      []string
		wantError bool
	}{
		{"defaults", []string{}, false},
		{"TLS and ACME together", []string{"-enable-tls", "-acme"}, true},
		{"TLS enabled without files", []string{"-enable-tls"}, true},
		{"TLS enabled with missing key", []string{"-enable-tls", "-tls-cert-file", "cert.pem"}, true},
		{"TLS enabled with non existent files", []string{"-enable-tls", "-tls-cert-file", "cert.pem", "-tls-key-file", "key.pem"}, true},
		{"Redirect without TLS or ACME", []string{"-redirect-to-https"}, true},
		{"Redirect with ACME (valid setup)", []string{
			"-acme",
			"-acme-accept-tos",
			"-acme-cache", "/tmp/certs",
			"-acme-email", "someone@example.com",
			"-acme-host", "example.com",
			"-redirect-to-https",
		}, false},
		{"ACME missing TOS", []string{"-acme", "-acme-cache", "/tmp", "-acme-email", "someone@example.com", "-acme-host", "example.com"}, true},
		{"ACME missing cache", []string{"-acme", "-acme-accept-tos", "-acme-email", "someone@example.com", "-acme-host", "example.com"}, true},
		{"ACME missing email", []string{"-acme", "-acme-accept-tos", "-acme-cache", "/tmp", "-acme-host", "example.com"}, true},
		{"ACME missing host", []string{"-acme", "-acme-accept-tos", "-acme-cache", "/tmp", "-acme-email", "someone@example.com"}, true},
		{"ACME full valid config", []string{
			"-acme",
			"-acme-accept-tos",
			"-acme-cache", "/tmp/certs",
			"-acme-email", "someone@example.com",
			"-acme-host", "example.com",
		}, false},
		{"Negative clear builds interval", []string{"-clear-builds-interval", "-5"}, true},
		{"Zero clear builds interval (valid)", []string{"-clear-builds-interval", "0"}, false},
		{"Invalid flag name", []string{"-not-a-flag", "value"}, true},
		{"Invalid port type", []string{"-http-port", "not-a-number"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := ParseConfig(tt.args, flag.ContinueOnError)
			if (err != nil) != tt.wantError {
				t.Errorf("got error %v, want error %v", err, tt.wantError)
				return
			}
		})
	}
}
