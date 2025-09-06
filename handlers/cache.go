package handlers

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var filenameRegex = regexp.MustCompile(`^zig(?:|-bootstrap|-[a-zA-Z0-9_]+-[a-zA-Z0-9_]+)-(\d+\.\d+\.\d+(?:-dev\.\d+\+[0-9a-f]+)?)\.(?:tar\.xz|zip)(?:\.minisig)?$`)

func CacheHandler(upstreamHost string, cacheDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := filepath.Base(r.URL.Path)

		// Validate filename and extract version
		matches := filenameRegex.FindStringSubmatch(filename)
		if len(matches) < 2 {
			http.Error(w, "Invalid filename format", http.StatusNotFound)
			return
		}
		version := matches[1]

		// Serve file if it is cached
		fileFullPath := filepath.Join(cacheDir, filename)
		if fileExists(fileFullPath) {
			slog.Info("serving cached file", "ip", GetRemoteIP(*r), "source", GetSource(*r), "file", filename)
			w.Header().Set("Content-Type", "application/octet-stream")
			http.ServeFile(w, r, fileFullPath)
			return
		}

		// The file isn't in the cache
		// Trying to download it

		// Determine upstream URL
		var sourceURL string
		if strings.Contains(version, "-dev") {
			sourceURL = fmt.Sprintf("%s/builds/%s", upstreamHost, filename)
		} else {
			sourceURL = fmt.Sprintf("%s/download/%s/%s", upstreamHost, version, filename)
		}

		// Download the file
		slog.Info("file is not in cache, trying to cache it", "file", filename, "source_url", sourceURL)

		localFilepath := filepath.Join(cacheDir, filename)

		resp, err := http.Get(sourceURL)
		if err != nil {
			slog.Warn("error fetching file from url", "source_url", sourceURL, "error", err)
			http.Error(w, "Bad Gateway", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			slog.Warn("File not found on upstream", "file", filename, "source_url", sourceURL, "upstream", upstreamHost)
			s := "File not found on upstream: " + sourceURL
			http.Error(w, s, http.StatusNotFound)
			return
		}

		if resp.StatusCode != http.StatusOK {
			slog.Warn("Upstream server returned error", "status_code", resp.StatusCode, "source_url", sourceURL, "upstream", upstreamHost)
			http.Error(w, "Bad Gateway", http.StatusBadGateway)
			return
		}

		// Create a temporary file to download to
		tmpfile, err := os.CreateTemp(cacheDir, filename)
		if err != nil {
			slog.Warn("Error creating temp file", "temp_file", tmpfile.Name(), "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer os.Remove(tmpfile.Name())

		// Stream the download to the temp file
		_, err = io.Copy(tmpfile, resp.Body)
		if err != nil {
			slog.Warn("Error writing to temp file", "temp_file", tmpfile.Name(), "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		tmpfile.Close()

		// Move the file to its final destination
		if err := os.Rename(tmpfile.Name(), localFilepath); err != nil {
			slog.Warn("Error renaming temp file", "temp_file", tmpfile.Name(), "local_filepath", localFilepath, "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		slog.Info("Successfully downloaded and cached", "filename", filename)

		// Serve the file
		slog.Info("serving cached file", "ip", GetRemoteIP(*r), "source", GetSource(*r), "file", filename)
		w.Header().Set("Content-Type", "application/octet-stream")
		http.ServeFile(w, r, fileFullPath)
	}
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !errors.Is(err, os.ErrNotExist)
}
