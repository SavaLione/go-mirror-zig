package main

import (
	"context"
	"crypto/tls"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/savalione/go-mirror-zig/handlers"
	"github.com/savalione/go-mirror-zig/internal/config"
	"golang.org/x/crypto/acme/autocert"
)

//go:embed templates/*
var content embed.FS

func main() {
	if err := run(); err != nil {
		slog.Error("server exited with an error", "error", err)
		os.Exit(1)
	}
	slog.Info("server exited gracefully")
}

func run() error {
	// Graceful shutdown.
	shutdownCtx, shutdownCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer shutdownCancel()

	// WaitGroup to wait for all goroutines to finish.
	var wg sync.WaitGroup

	// Parsing templates that are in cmd/templates
	tmpl, err := template.ParseFS(content, "templates/*.html")
	if err != nil {
		return fmt.Errorf("error parsing templates: %w", err)
	}

	cfg, err := config.ParseConfig()
	if err != nil {
		return fmt.Errorf("error parsing configuration: %w", err)
	}

	// HTTP and HTTPS Handler setup
	mux := http.NewServeMux()
	cache := handlers.NewCache(cfg.UpstreamURL, cfg.CacheDir)

	mux.HandleFunc("/", handlers.RootHandler(tmpl))
	mux.HandleFunc("/{file}", cache.Handler())
	mux.HandleFunc("/zig/{file}", cache.Handler())
	mux.HandleFunc("/builds/{file}", cache.Handler())
	mux.HandleFunc("/download/", cache.Handler())
	mainHandler := handlers.Middleware(mux)

	var servers []*http.Server

	if cfg.EnableTLS {
		tlsConfig := &tls.Config{
			Certificates:             []tls.Certificate{cfg.KeyPair},
			PreferServerCipherSuites: true,
			MinVersion:               tls.VersionTLS13,
			CurvePreferences:         []tls.CurveID{tls.CurveP256, tls.X25519},
		}

		httpsServer := &http.Server{
			Addr:         cfg.HTTPSAddress(),
			Handler:      mainHandler,
			TLSConfig:    tlsConfig,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		}

		servers = append(servers, httpsServer)
	} else if cfg.ACME {
		acmeManager := &autocert.Manager{
			Cache:      autocert.DirCache(cfg.ACMECache),
			Prompt:     cfg.AcceptTOS,
			Email:      cfg.ACMEEmail,
			HostPolicy: autocert.HostWhitelist(cfg.ACMEHost),
		}

		acmeServer := &http.Server{
			Addr:         cfg.HTTPSAddress(),
			Handler:      mainHandler,
			TLSConfig:    acmeManager.TLSConfig(),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		}

		servers = append(servers, acmeServer)
	} else {
		httpServer := &http.Server{
			Addr:         cfg.HTTPAddress(),
			Handler:      mainHandler,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		}
		servers = append(servers, httpServer)
	}

	if cfg.RedirectToHTTPS {
		redirectServer := &http.Server{
			Addr:         cfg.HTTPAddress(),
			Handler:      handlers.RedirectHandler(cfg.TLSPort),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		}
		servers = append(servers, redirectServer)
	}

	for _, srv := range servers {
		wg.Add(1)
		go startServer(shutdownCtx, &wg, srv)
	}

	<-shutdownCtx.Done()
	slog.Info("shutdown signal received. waiting for servers to close.")

	wg.Wait()

	// A hard shutdown failsafe
	go func() {
		time.Sleep(10 * time.Second)
		slog.Warn("hard shutdown initiated")
		os.Exit(1)
	}()

	return nil
}

func startServer(ctx context.Context, wg *sync.WaitGroup, srv *http.Server) {
	defer wg.Done()

	isTLS := srv.TLSConfig != nil
	serverType := "HTTP"
	if isTLS {
		serverType = "HTTPS"
	}

	go func() {
		slog.Info(fmt.Sprintf("starting %s server", serverType), "addr", srv.Addr)
		var err error
		if isTLS {
			err = srv.ListenAndServeTLS("", "")
		} else {
			err = srv.ListenAndServe()
		}

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error(fmt.Sprintf("%s server failed to start", serverType), "addr", srv.Addr, "error", err)
		}
	}()

	<-ctx.Done()

	slog.Info(fmt.Sprintf("shutting down %s server", serverType), "addr", srv.Addr)

	// A context with a timeout for the shutdown
	shutdownTimeoutCtx, shutdownTimeoutCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownTimeoutCancel()

	if err := srv.Shutdown(shutdownTimeoutCtx); err != nil {
		slog.Error(fmt.Sprintf("failed to shut down %s server gracefully", serverType), "addr", srv.Addr, "error", err)
	}
}
