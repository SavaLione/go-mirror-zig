package config

import (
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
