package zig

import (
	"regexp"
	"strconv"
)

// Matches standard Zig distribution filenames, capturing the version string.
var filenameRegex = regexp.MustCompile(`^zig(?:|-bootstrap|-[a-zA-Z0-9_]+-[a-zA-Z0-9_]+)-(\d+\.\d+\.\d+(?:-dev\.\d+\+[0-9a-f]+)?)\.(?:tar\.xz|zip)(?:\.minisig)?$`)

// Checks whether the provided string follow the official Zig artifact naming convention.
// It returns true if the string matches the expected pattern, false otherwise.
func IsZigArtifact(s string) bool {
	if len(filenameRegex.FindStringSubmatch(s)) < 2 {
		return false
	}
	return true
}

// Extracts the submatches from a Zig artifact filename.
// Returns nil if the provided string does not follow the official Zig artifact naming convention.
func ArtifactSubmatches(s string) []string {

	matches := filenameRegex.FindStringSubmatch(s)
	if len(matches) < 2 {
		return nil
	}

	return matches
}

// A Zig artifact from index.json
type Artifact struct {
	Tarball string
	Shasum  string
	Size    string
}

// A Zig release from index.json
type Release struct {
	Version   string              `json:"version"`
	Date      string              `json:"date"`
	Docs      string              `json:"docs"`
	StdDocs   string              `json:"stdDocs"`
	Notes     string              `json:"notes"`
	Platforms map[string]Artifact `json:"-"`
}

// All releases and artifacts (including the dev branch one)
type ZigReleases map[string]Release

// Total size of all Zig artifacts
func TotalSize(zr ZigReleases) int {
	total := 0
	for _, release := range zr {
		for _, artifact := range release.Platforms {
			i, err := strconv.Atoi(artifact.Size)
			if err != nil {
				continue
			}
			total += i
		}
	}
	return total
}
