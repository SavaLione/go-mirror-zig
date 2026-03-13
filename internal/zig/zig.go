package zig

import "regexp"

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
