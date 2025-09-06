package main

import (
	"context"
	"embed"
	"errors"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/savalione/go-mirror-zig/handlers"
)

type Flags struct {
	CacheDir    *string
	UpstreamURL *string
	Port        *int
	IP          *string
}

func (f *Flags) init() {
	f.CacheDir = flag.String("cache-dir", "./", "Directory to store cache")
	f.UpstreamURL = flag.String("upstream-url", "https://ziglang.org", "Zig upstream mirror")
	f.Port = flag.Int("port", 8080, "Port to listen on")
	f.IP = flag.String("ip", "", "IP to listen on")
}

func (f Flags) address() string {
	return *f.IP + ":" + strconv.Itoa(*f.Port)
}

//go:embed templates/*
var content embed.FS

func main() {
	// Signals
	shutdownCtx, shutdownCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer shutdownCancel()

	go func() {
		<-shutdownCtx.Done()

		<-time.After(10 * time.Second)
		slog.Info("Hard shutdown")
		os.Exit(0)
	}()

	tmpl, err := template.ParseFS(content, "templates/*.html")
	if err != nil {
		slog.Error("Error parsing templates", "error", err)
	}

	// Flags
	var flags Flags
	flags.init()
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.RootHandler(tmpl))
	mux.HandleFunc("/{file}", handlers.CacheHandler(*flags.UpstreamURL, *flags.CacheDir))
	mux.HandleFunc("/zig/{file}", handlers.CacheHandler(*flags.UpstreamURL, *flags.CacheDir))

	srv := &http.Server{
		Addr:              flags.address(),
		Handler:           handlers.Middleware(mux),
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	go func() {
		slog.Info("Starting server")
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Couldn't start the server", "error", err)
		}
	}()

	<-shutdownCtx.Done()

	shutdownTimeoutCtx, shutdownTimeoutCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownTimeoutCancel()

	slog.Info("Stopping the server")
	if err := srv.Shutdown(shutdownTimeoutCtx); err != nil {
		slog.Error("Couldn't shutdown the server", "error", err)
	}

	slog.Info("Server stopped")
}
