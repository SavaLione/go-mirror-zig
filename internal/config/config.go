package config

import (
	"flag"
	"net"
	"strconv"
)

// Config holds configuration values, populated from command-line flags.
type Config struct {
	CacheDir      string
	UpstreamURL   string
	HTTPPort      int
	ListenAddress string
}

// ParseConfig defines and parses command-line flags, validates them, and returns a populated Config struct.
func ParseConfig() (Config, error) {
	var c Config

	flag.StringVar(&c.CacheDir, "cache-dir", "./", "Path to the directory where downloaded content will be cached.")
	flag.StringVar(&c.UpstreamURL, "upstream-url", "https://ziglang.org", "The URL of the upstream server to mirror/proxy.")
	flag.IntVar(&c.HTTPPort, "http-port", 80, "The port for the plain HTTP listener.")
	flag.StringVar(&c.ListenAddress, "listen-address", "", "The IP address to listen on. If empty, listens on all available interfaces.")

	flag.Parse()

	return c, nil
}

// HTTPAddress returns the full address for the HTTP server.
func (c Config) HTTPAddress() string {
	return net.JoinHostPort(c.ListenAddress, strconv.Itoa(c.HTTPPort))
}
