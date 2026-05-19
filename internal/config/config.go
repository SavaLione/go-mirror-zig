package config

import (
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"strconv"
)

// Config holds configuration values, populated from command-line flags.
type Config struct {
	CacheDir         string
	UpstreamURL      string
	HTTPPort         int
	TLSPort          int
	ListenAddress    string
	EnableTLS        bool
	RedirectToHTTPS  bool
	ShowVersion      bool
	ShowPossibleSize bool
	ShowIndexPage    bool
	IndexPage        string
	ClearBuilds      int

	ACME          bool
	ACMEDirectory string
	ACMEAcceptTOS bool
	ACMECache     string
	ACMEEmail     string
	ACMEHost      string
	acmeTOSURL    string

	// KeyPair is the loaded TLS certificate and key. It is only populated if EnableTLS is true.
	KeyPair tls.Certificate

	// unexported fields for parsing flags before validation.
	tlsCertFile string
	tlsKeyFile  string
}

// ParseConfig defines and parses command-line flags, validates them, and returns a populated Config struct.
func ParseConfig(args []string, errorHandling flag.ErrorHandling) (Config, error) {
	var c Config

	fs := flag.NewFlagSet("go-mirror-zig", errorHandling)

	if errorHandling == flag.ContinueOnError {
		fs.SetOutput(io.Discard) // suppress console text on tests
	}

	fs.StringVar(&c.CacheDir, "cache-dir", "./", "Path to the directory where downloaded content will be cached.")
	fs.StringVar(&c.UpstreamURL, "upstream-url", "https://ziglang.org", "The URL of the upstream server to mirror/proxy.")
	fs.IntVar(&c.HTTPPort, "http-port", 80, "The port for the plain HTTP listener.")
	fs.IntVar(&c.TLSPort, "tls-port", 443, "The port for the secure TLS (HTTPS) listener.")
	fs.StringVar(&c.ListenAddress, "listen-address", "", "The IP address to listen on. If empty, listens on all available interfaces.")
	fs.BoolVar(&c.EnableTLS, "enable-tls", false, "Enable the TLS (HTTPS) server. Requires -tls-cert-file and -tls-key-file.")
	fs.StringVar(&c.tlsCertFile, "tls-cert-file", "", "Path to the TLS certificate file.")
	fs.StringVar(&c.tlsKeyFile, "tls-key-file", "", "Path to the TLS private key file.")
	fs.BoolVar(&c.RedirectToHTTPS, "redirect-to-https", false, "Enable automatic redirection of HTTP requests to HTTPS. Requires -enable-tls or -acme.")
	fs.BoolVar(&c.ShowVersion, "version", false, "Print version information and exit.")
	fs.BoolVar(&c.ShowPossibleSize, "show-possible-size", false, "Print estimation stats of all cacheable upstream artifacts (size, release counts) and exit.")
	fs.BoolVar(&c.ShowIndexPage, "show-index-page", true, "Whether to serve a custom index page at the root (/). Set to false to disable.")
	fs.StringVar(&c.IndexPage, "index-page", "", "Path to a directory containing static files to serve as the root index. If empty, uses the default built-in index page.")
	fs.IntVar(&c.ClearBuilds, "clear-builds-interval", 7200, "Interval in seconds to clean up cached dev builds. Set to 0 to disable.")

	fs.BoolVar(&c.ACME, "acme", false, "Obtain TLS certificates using the ACME challenge.")
	fs.StringVar(&c.ACMEDirectory, "acme-directory", "https://acme-v02.api.letsencrypt.org/directory", "ACME directory URL.")
	fs.BoolVar(&c.ACMEAcceptTOS, "acme-accept-tos", false, "Accept the ACME provider's Terms of Service.")
	fs.StringVar(&c.ACMECache, "acme-cache", "", "Directory for storing obtained certificates.")
	fs.StringVar(&c.ACMEEmail, "acme-email", "", "Email address for ACME registration and recovery notices.")
	fs.StringVar(&c.ACMEHost, "acme-host", "", "The hostname (domain name) for which to obtain the ACME certificate.")

	err := fs.Parse(args)
	if err != nil {
		return c, err
	}

	if c.EnableTLS && c.ACME {
		return c, errors.New("cannot use both -enable-tls (manual certificates) and -acme (automatic certificates) at the same time")
	}

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

	// if !c.EnableTLS && (c.tlsCertFile != "" || c.tlsKeyFile != "") {
	// 	slog.Warn("-tls-cert-file or -tls-key-file provided without -enable-tls. These flags will be ignored.")
	// }

	if c.RedirectToHTTPS && !(c.EnableTLS || c.ACME) {
		return c, errors.New("-redirect-to-https requires -enable-tls or -acme to be set")
	}

	if c.ACME {
		if !c.ACMEAcceptTOS {
			return c, errors.New("using ACME requires accepting the Terms of Service; please provide the -acme-accept-tos flag")
		}

		if c.ACMECache == "" {
			return c, errors.New("-acme-cache must be set to a directory for storing certificates")
		}

		if c.ACMEEmail == "" {
			return c, errors.New("the -acme-email flag must be provided")
		}

		if c.ACMEHost == "" {
			return c, errors.New("the -acme-host flag must be provided")
		}

		if c.EnableTLS {
			return c, errors.New("manual TLS settings (-tls-cert-file, -tls-key-file) cannot be used with ACME")
		}
	}

	if c.ClearBuilds < 0 {
		return c, errors.New("the -clear-builds-interval flag can't be negative")
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

func (c *Config) AcceptTOS(tosURL string) bool {
	if c.ACMEAcceptTOS {
		c.acmeTOSURL = "Accepting ACME Terms of Service at: " + tosURL
	} else {
		c.acmeTOSURL = "Terms of Service are not accepted"
	}
	return c.ACMEAcceptTOS
}
