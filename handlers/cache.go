package handlers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

var filenameRegex = regexp.MustCompile(`^zig(?:|-bootstrap|-[a-zA-Z0-9_]+-[a-zA-Z0-9_]+)-(\d+\.\d+\.\d+(?:-dev\.\d+\+[0-9a-f]+)?)\.(?:tar\.xz|zip)(?:\.minisig)?$`)

// Cache holds the dependencies for the cache handler, making it more testable and organized.
type Cache struct {
	upstreamHost string
	cacheDir     string
	client       *http.Client // Use a custom client for timeouts.
	fileLocks    sync.Map     // Safely stores locks for in-flight downloads.
}

// NewCache creates a new Cache handler dependency object.
func NewCache(upstreamHost, cacheDir string) *Cache {
	return &Cache{
		upstreamHost: upstreamHost,
		cacheDir:     cacheDir,
		client: &http.Client{
			Timeout: 30 * time.Minute, // Timeout for the entire download.
			Transport: &http.Transport{
				IdleConnTimeout: 90 * time.Second,
			},
		},
		fileLocks: sync.Map{},
	}
}

// Handler returns the http.HandlerFunc for caching.
func (c *Cache) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := filepath.Base(r.URL.Path)

		logger := slog.With(
			"remote_ip", GetRemoteIP(*r),
			"path", r.URL.Path,
			"filename", filename,
			"source", GetSource(*r),
		)

		// Validate filename.
		matches := filenameRegex.FindStringSubmatch(filename)
		if len(matches) < 2 {
			logger.Warn("invalid filename format")
			http.Error(w, "Invalid filename format", http.StatusBadRequest)
			return
		}

		// Serve the file if it is cached.
		fileFullPath := filepath.Join(c.cacheDir, filename)
		if fileExists(fileFullPath) {
			serveFile(w, r, fileFullPath, logger)
			return
		}

		// Lock and download. The file is not in the cache.
		// The file needs to be downloaded. Lock to prevent multiple concurrent
		// downloads for the same file.
		mu, _ := c.fileLocks.LoadOrStore(filename, &sync.Mutex{})
		fileMutex := mu.(*sync.Mutex)

		fileMutex.Lock()
		defer func() {
			fileMutex.Unlock()
			c.fileLocks.Delete(filename)
		}()

		// Double-check if another request downloaded the file while we were waiting for the lock.
		if fileExists(fileFullPath) {
			logger.Info("file was cached by another request while waiting for lock")
			serveFile(w, r, fileFullPath, logger)
			return
		}

		// Fetch from upstream
		logger.Info("file not in cache, starting download")
		if err := c.fetchAndCacheFile(r.Context(), logger, filename, matches[1]); err != nil {
			if errors.Is(err, errUpstreamNotFound) {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			} else if errors.Is(err, errUpstreamUnavailable) {
				http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		// Serve the newly cached file.
		serveFile(w, r, fileFullPath, logger)
	}
}

// Custom error types for clearer error handling.
var (
	errUpstreamNotFound    = errors.New("file not found on upstream")
	errUpstreamUnavailable = errors.New("upstream server returned non-OK status")
)

func (c *Cache) fetchAndCacheFile(ctx context.Context, logger *slog.Logger, filename, version string) error {
	// Determine upstream URL
	var sourceURL string
	if strings.Contains(version, "-dev") {
		sourceURL = fmt.Sprintf("%s/builds/%s", c.upstreamHost, filename)
	} else {
		sourceURL = fmt.Sprintf("%s/download/%s/%s", c.upstreamHost, version, filename)
	}

	logger = logger.With("source_url", sourceURL)
	logger.Info("fetching file from upstream")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sourceURL, nil)
	if err != nil {
		logger.Error("failed to create upstream request", "error", err)
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		logger.Error("failed to fetch file from upstream", "error", err)
		return errUpstreamUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		logger.Warn("file not found on upstream")
		return errUpstreamNotFound
	}
	if resp.StatusCode != http.StatusOK {
		logger.Error("upstream server returned non-OK status", "status_code", resp.StatusCode)
		return errUpstreamUnavailable
	}

	// Create a temporary file to avoid serving a partially downloaded file.
	tmpFile, err := os.CreateTemp(c.cacheDir, filename+".*.tmp")
	if err != nil {
		logger.Error("failed to create temporary file", "error", err)
		return err
	}
	defer func() {
		// tmpFile.Close()
		os.Remove(tmpFile.Name())
	}()

	// Stream the download to the temp file.
	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		logger.Error("failed to write to temporary file", "temp_file", tmpFile.Name(), "error", err)
		return err
	}
	tmpFile.Close()

	// Atomically move the file to its final destination.
	finalPath := filepath.Join(c.cacheDir, filename)
	if err := os.Rename(tmpFile.Name(), finalPath); err != nil {
		logger.Error("failed to rename temporary file", "from", tmpFile.Name(), "to", finalPath, "error", err)
		return err
	}

	logger.Info("successfully downloaded and cached file")
	return nil
}

func serveFile(w http.ResponseWriter, r *http.Request, path string, logger *slog.Logger) {
	logger.Info("serving file from cache")
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, path)
}

// fileExists checks if a file exists and is not a directory.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
