package zig

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Fetch the file with all Zig releases
// Such file is usually located here: https://ziglang.org/download/index.json
func FetchAllReleases(ctx context.Context, url string) (ZigReleases, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var zr ZigReleases
	if err := json.Unmarshal(body, &zr); err != nil {
		return nil, err
	}

	return zr, nil
}
