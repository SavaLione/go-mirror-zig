package config

import (
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"net"
	"strconv"
)

// Config holds configuration values, populated from command-line flags.
type Config struct {
	CacheDir        string
	UpstreamURL     string
	HTTPPort        int
	TLSPort         int
	ListenAddress   string
	EnableTLS       bool
	RedirectToHTTPS bool

	// KeyPair is the loaded TLS certificate and key. It is only populated if EnableTLS is true.
	KeyPair tls.Certificate

	// unexported fields for parsing flags before validation.
	tlsCertFile string
	tlsKeyFile  string
}

// ParseConfig defines and parses command-line flags, validates them, and returns a populated Config struct.
func ParseConfig() (Config, error) {
	var c Config

	flag.StringVar(&c.CacheDir, "cache-dir", "./", "Path to the directory where downloaded content will be cached.")
	flag.StringVar(&c.UpstreamURL, "upstream-url", "https://ziglang.org", "The URL of the upstream server to mirror/proxy.")
	flag.IntVar(&c.HTTPPort, "http-port", 80, "The port for the plain HTTP listener.")
	flag.IntVar(&c.TLSPort, "tls-port", 443, "The port for the secure TLS (HTTPS) listener.")
	flag.StringVar(&c.ListenAddress, "listen-address", "", "The IP address to listen on. If empty, listens on all available interfaces.")
	flag.BoolVar(&c.EnableTLS, "enable-tls", false, "Enable the TLS (HTTPS) server. Requires -tls-cert-file and -tls-key-file.")
	flag.StringVar(&c.tlsCertFile, "tls-cert-file", "", "Path to the TLS certificate file.")
	flag.StringVar(&c.tlsKeyFile, "tls-key-file", "", "Path to the TLS private key file.")
	flag.BoolVar(&c.RedirectToHTTPS, "redirect-to-https", false, "Enable automatic redirection of HTTP requests to HTTPS. Requires -enable-tls.")

	flag.Parse()

	if c.EnableTLS {
		if c.tlsKeyFile == "" || c.tlsCertFile == "" {
			return c, errors.New("to enable TLS, both -tls-cert-file and -tls-key-file flags must be provided")
		}

		keyPair, err := tls.LoadX509KeyPair(c.tlsCertFile, c.tlsKeyFile)
		if err != nil {
			return c, fmt.Errorf("failed to load TLS key pair: %w", err)
		}
		c.KeyPair = keyPair
	}

	if c.RedirectToHTTPS && !c.EnableTLS {
		return c, errors.New("-redirect-to-https requires -enable-tls to be set")
	}

	return c, nil
}

// HTTPAddress returns the full address for the HTTP server.
func (c Config) HTTPAddress() string {
	return net.JoinHostPort(c.ListenAddress, strconv.Itoa(c.HTTPPort))
}

// HTTPSAddress returns the full address for the HTTPS server.
func (c Config) HTTPSAddress() string {
	return net.JoinHostPort(c.ListenAddress, strconv.Itoa(c.TLSPort))
}
