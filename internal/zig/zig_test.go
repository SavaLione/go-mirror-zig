package zig

import (
	"slices"
	"testing"
)

func TestIsZigArtifact(t *testing.T) {
	t.Parallel()

	var zigArtifacts = []struct {
		in       string
		expected bool
	}{
		{"", false},
		{"zig", false},
		{"zig-0.14.1.tarxz", false},
		{"zig-0.14.1.tar .xz", false},

		// See:
		// https://codeberg.org/ziglang/ziglang.org/src/branch/main/check-mirrors/main.zig
		{"zig-arm-linux-0.16.0.tar.xz", true},
		{"zig-riscv64-linux-0.16.0.tar.xz", true},
		{"zig-loongarch64-linux-0.16.0.tar.xz", true},
		{"zig-powerpc64le-freebsd-0.16.0.tar.xz", true},
		{"zig-aarch64-netbsd-0.16.0.tar.xz", true},
		{"zig-arm-openbsd-0.16.0.tar.xz", true},

		{"zig-0.14.1.tar.xz", true},
		{"zig-x86_64-windows-0.14.1.zip", true},
		{"zig-aarch64-macos-0.14.1.tar.xz", true},
		{"zig-x86_64-linux-0.14.1.tar.xz", true},
		{"zig-aarch64-linux-0.14.1.tar.xz", true},
		{"zig-armv7a-linux-0.14.1.tar.xz", true},
		{"zig-riscv64-linux-0.14.1.tar.xz", true},
		{"zig-powerpc64le-linux-0.14.1.tar.xz", true},
		{"zig-x86-linux-0.14.1.tar.xz", true},
		{"zig-loongarch64-linux-0.14.1.tar.xz", true},
		{"zig-s390x-linux-0.14.1.tar.xz", true},
		{"zig-0.10.1.tar.xz", true},
		{"zig-bootstrap-0.10.1.tar.xz", true},
		{"zig-linux-i386-0.10.1.tar.xz", true},
		{"zig-macos-aarch64-0.10.1.tar.xz", true},
		{"zig-windows-x86_64-0.10.1.zip", true},
		{"zig-0.7.1.tar.xz", true},
		{"zig-linux-x86_64-0.7.1.tar.xz", true},
		{"zig-0.6.0.tar.xz", true},
		{"zig-linux-x86_64-0.6.0.tar.xz", true},

		// With .minisig
		{"zig-arm-linux-0.16.0.tar.xz.minisig", true},
		{"zig-riscv64-linux-0.16.0.tar.xz.minisig", true},
		{"zig-loongarch64-linux-0.16.0.tar.xz.minisig", true},
		{"zig-powerpc64le-freebsd-0.16.0.tar.xz.minisig", true},
		{"zig-aarch64-netbsd-0.16.0.tar.xz.minisig", true},
		{"zig-arm-openbsd-0.16.0.tar.xz.minisig", true},

		{"zig-0.14.1.tar.xz.minisig", true},
		{"zig-x86_64-windows-0.14.1.zip.minisig", true},
		{"zig-aarch64-macos-0.14.1.tar.xz.minisig", true},
		{"zig-x86_64-linux-0.14.1.tar.xz.minisig", true},
		{"zig-aarch64-linux-0.14.1.tar.xz.minisig", true},
		{"zig-armv7a-linux-0.14.1.tar.xz.minisig", true},
		{"zig-riscv64-linux-0.14.1.tar.xz.minisig", true},
		{"zig-powerpc64le-linux-0.14.1.tar.xz.minisig", true},
		{"zig-x86-linux-0.14.1.tar.xz.minisig", true},
		{"zig-loongarch64-linux-0.14.1.tar.xz.minisig", true},
		{"zig-s390x-linux-0.14.1.tar.xz.minisig", true},
		{"zig-0.10.1.tar.xz.minisig", true},
		{"zig-bootstrap-0.10.1.tar.xz.minisig", true},
		{"zig-linux-i386-0.10.1.tar.xz.minisig", true},
		{"zig-macos-aarch64-0.10.1.tar.xz.minisig", true},
		{"zig-windows-x86_64-0.10.1.zip.minisig", true},
		{"zig-0.7.1.tar.xz.minisig", true},
		{"zig-linux-x86_64-0.7.1.tar.xz.minisig", true},
		{"zig-0.6.0.tar.xz.minisig", true},
		{"zig-linux-x86_64-0.6.0.tar.xz.minisig", true},
	}

	for _, tt := range zigArtifacts {
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()
			check := IsZigArtifact(tt.in)
			if check != tt.expected {
				t.Errorf("got %v, want %v", check, tt.expected)
			}
		})
	}
}

func TestArtifactSubmatches(t *testing.T) {
	t.Parallel()

	var zigArtifacts = []struct {
		in       string
		expected []string
	}{
		{"", nil},
		{"zig", nil},
		{"zig-0.14.1.tarxz", nil},
		{"zig-0.14.1.tar .xz", nil},

		// See:
		// https://codeberg.org/ziglang/ziglang.org/src/branch/main/check-mirrors/main.zig
		{"zig-arm-linux-0.16.0.tar.xz", []string{"zig-arm-linux-0.16.0.tar.xz", "0.16.0"}},
		{"zig-riscv64-linux-0.16.0.tar.xz", []string{"zig-riscv64-linux-0.16.0.tar.xz", "0.16.0"}},
		{"zig-loongarch64-linux-0.16.0.tar.xz", []string{"zig-loongarch64-linux-0.16.0.tar.xz", "0.16.0"}},
		{"zig-powerpc64le-freebsd-0.16.0.tar.xz", []string{"zig-powerpc64le-freebsd-0.16.0.tar.xz", "0.16.0"}},
		{"zig-aarch64-netbsd-0.16.0.tar.xz", []string{"zig-aarch64-netbsd-0.16.0.tar.xz", "0.16.0"}},
		{"zig-arm-openbsd-0.16.0.tar.xz", []string{"zig-arm-openbsd-0.16.0.tar.xz", "0.16.0"}},

		{"zig-0.14.1.tar.xz", []string{"zig-0.14.1.tar.xz", "0.14.1"}},
		{"zig-x86_64-windows-0.14.1.zip", []string{"zig-x86_64-windows-0.14.1.zip", "0.14.1"}},
		{"zig-aarch64-macos-0.14.1.tar.xz", []string{"zig-aarch64-macos-0.14.1.tar.xz", "0.14.1"}},
		{"zig-x86_64-linux-0.14.1.tar.xz", []string{"zig-x86_64-linux-0.14.1.tar.xz", "0.14.1"}},
		{"zig-aarch64-linux-0.14.1.tar.xz", []string{"zig-aarch64-linux-0.14.1.tar.xz", "0.14.1"}},
		{"zig-armv7a-linux-0.14.1.tar.xz", []string{"zig-armv7a-linux-0.14.1.tar.xz", "0.14.1"}},
		{"zig-riscv64-linux-0.14.1.tar.xz", []string{"zig-riscv64-linux-0.14.1.tar.xz", "0.14.1"}},
		{"zig-powerpc64le-linux-0.14.1.tar.xz", []string{"zig-powerpc64le-linux-0.14.1.tar.xz", "0.14.1"}},
		{"zig-x86-linux-0.14.1.tar.xz", []string{"zig-x86-linux-0.14.1.tar.xz", "0.14.1"}},
		{"zig-loongarch64-linux-0.14.1.tar.xz", []string{"zig-loongarch64-linux-0.14.1.tar.xz", "0.14.1"}},
		{"zig-s390x-linux-0.14.1.tar.xz", []string{"zig-s390x-linux-0.14.1.tar.xz", "0.14.1"}},
		{"zig-0.10.1.tar.xz", []string{"zig-0.10.1.tar.xz", "0.10.1"}},
		{"zig-bootstrap-0.10.1.tar.xz", []string{"zig-bootstrap-0.10.1.tar.xz", "0.10.1"}},
		{"zig-linux-i386-0.10.1.tar.xz", []string{"zig-linux-i386-0.10.1.tar.xz", "0.10.1"}},
		{"zig-macos-aarch64-0.10.1.tar.xz", []string{"zig-macos-aarch64-0.10.1.tar.xz", "0.10.1"}},
		{"zig-windows-x86_64-0.10.1.zip", []string{"zig-windows-x86_64-0.10.1.zip", "0.10.1"}},
		{"zig-0.7.1.tar.xz", []string{"zig-0.7.1.tar.xz", "0.7.1"}},
		{"zig-linux-x86_64-0.7.1.tar.xz", []string{"zig-linux-x86_64-0.7.1.tar.xz", "0.7.1"}},
		{"zig-0.6.0.tar.xz", []string{"zig-0.6.0.tar.xz", "0.6.0"}},
		{"zig-linux-x86_64-0.6.0.tar.xz", []string{"zig-linux-x86_64-0.6.0.tar.xz", "0.6.0"}},

		// With .minisig
		{"zig-arm-linux-0.16.0.tar.xz.minisig", []string{"zig-arm-linux-0.16.0.tar.xz.minisig", "0.16.0"}},
		{"zig-riscv64-linux-0.16.0.tar.xz.minisig", []string{"zig-riscv64-linux-0.16.0.tar.xz.minisig", "0.16.0"}},
		{"zig-loongarch64-linux-0.16.0.tar.xz.minisig", []string{"zig-loongarch64-linux-0.16.0.tar.xz.minisig", "0.16.0"}},
		{"zig-powerpc64le-freebsd-0.16.0.tar.xz.minisig", []string{"zig-powerpc64le-freebsd-0.16.0.tar.xz.minisig", "0.16.0"}},
		{"zig-aarch64-netbsd-0.16.0.tar.xz.minisig", []string{"zig-aarch64-netbsd-0.16.0.tar.xz.minisig", "0.16.0"}},
		{"zig-arm-openbsd-0.16.0.tar.xz.minisig", []string{"zig-arm-openbsd-0.16.0.tar.xz.minisig", "0.16.0"}},

		{"zig-0.14.1.tar.xz.minisig", []string{"zig-0.14.1.tar.xz.minisig", "0.14.1"}},
		{"zig-x86_64-windows-0.14.1.zip.minisig", []string{"zig-x86_64-windows-0.14.1.zip.minisig", "0.14.1"}},
		{"zig-aarch64-macos-0.14.1.tar.xz.minisig", []string{"zig-aarch64-macos-0.14.1.tar.xz.minisig", "0.14.1"}},
		{"zig-x86_64-linux-0.14.1.tar.xz.minisig", []string{"zig-x86_64-linux-0.14.1.tar.xz.minisig", "0.14.1"}},
		{"zig-aarch64-linux-0.14.1.tar.xz.minisig", []string{"zig-aarch64-linux-0.14.1.tar.xz.minisig", "0.14.1"}},
		{"zig-armv7a-linux-0.14.1.tar.xz.minisig", []string{"zig-armv7a-linux-0.14.1.tar.xz.minisig", "0.14.1"}},
		{"zig-riscv64-linux-0.14.1.tar.xz.minisig", []string{"zig-riscv64-linux-0.14.1.tar.xz.minisig", "0.14.1"}},
		{"zig-powerpc64le-linux-0.14.1.tar.xz.minisig", []string{"zig-powerpc64le-linux-0.14.1.tar.xz.minisig", "0.14.1"}},
		{"zig-x86-linux-0.14.1.tar.xz.minisig", []string{"zig-x86-linux-0.14.1.tar.xz.minisig", "0.14.1"}},
		{"zig-loongarch64-linux-0.14.1.tar.xz.minisig", []string{"zig-loongarch64-linux-0.14.1.tar.xz.minisig", "0.14.1"}},
		{"zig-s390x-linux-0.14.1.tar.xz.minisig", []string{"zig-s390x-linux-0.14.1.tar.xz.minisig", "0.14.1"}},
		{"zig-0.10.1.tar.xz.minisig", []string{"zig-0.10.1.tar.xz.minisig", "0.10.1"}},
		{"zig-bootstrap-0.10.1.tar.xz.minisig", []string{"zig-bootstrap-0.10.1.tar.xz.minisig", "0.10.1"}},
		{"zig-linux-i386-0.10.1.tar.xz.minisig", []string{"zig-linux-i386-0.10.1.tar.xz.minisig", "0.10.1"}},
		{"zig-macos-aarch64-0.10.1.tar.xz.minisig", []string{"zig-macos-aarch64-0.10.1.tar.xz.minisig", "0.10.1"}},
		{"zig-windows-x86_64-0.10.1.zip.minisig", []string{"zig-windows-x86_64-0.10.1.zip.minisig", "0.10.1"}},
		{"zig-0.7.1.tar.xz.minisig", []string{"zig-0.7.1.tar.xz.minisig", "0.7.1"}},
		{"zig-linux-x86_64-0.7.1.tar.xz.minisig", []string{"zig-linux-x86_64-0.7.1.tar.xz.minisig", "0.7.1"}},
		{"zig-0.6.0.tar.xz.minisig", []string{"zig-0.6.0.tar.xz.minisig", "0.6.0"}},
		{"zig-linux-x86_64-0.6.0.tar.xz.minisig", []string{"zig-linux-x86_64-0.6.0.tar.xz.minisig", "0.6.0"}},
	}

	for _, tt := range zigArtifacts {
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()
			check := ArtifactSubmatches(tt.in)

			if !slices.Equal(check, tt.expected) {
				t.Errorf("got %v, want %v", check, tt.expected)
			}
		})
	}
}

func TestTotalSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    ZigReleases
		expected int
	}{
		{
			name:     "no releases",
			input:    nil,
			expected: 0,
		},
		{
			name: "no artifacts",
			input: ZigReleases{
				"master": Release{
					Version:   "0.17.0-dev.305+bdfbf432d",
					Platforms: map[string]Artifact{},
				},
			},
			expected: 0,
		},
		{
			name: "valid sizes across multiple releases and platforms",
			input: ZigReleases{
				"master": Release{
					Platforms: map[string]Artifact{
						"src":       {Size: "22530120"},
						"bootstrap": {Size: "56523612"},
					},
				},
				"0.16.0": Release{
					Platforms: map[string]Artifact{
						"src":       {Size: "22503260"},
						"bootstrap": {Size: "55245980"},
					},
				},
			},
			expected: 22530120 + 56523612 + 22503260 + 55245980,
		},

		{
			name: "mixed valid and invalid size strings",
			input: ZigReleases{
				"master": Release{
					Platforms: map[string]Artifact{
						"src":           {Size: "22530120"},
						"bootstrap":     {Size: "not-valid"},
						"riscv64-linux": {Size: "57027896"},
					},
				},
			},
			expected: 22530120 + 57027896,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := TotalSize(tt.input)
			if actual != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, actual)
			}
		})
	}
}
