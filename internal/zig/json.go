package zig

import "encoding/json"

// Convert index.json (the file with all Zig releases) to the ZigReleases structure
// It omits all artifacts without a tarball
func (r *Release) UnmarshalJSON(data []byte) error {
	type Alias Release
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// See: https://stackoverflow.com/a/53352532
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	r.Platforms = make(map[string]Artifact)

	for key, val := range raw {
		switch key {
		case "version", "date", "docs", "stdDocs", "notes":
			continue
		default:
			var artifact Artifact
			if err := json.Unmarshal(val, &artifact); err == nil {
				// Omitting artifacts without a tarball
				if artifact.Tarball != "" {
					r.Platforms[key] = artifact
				}
			}
		}

	}
	return nil
}
