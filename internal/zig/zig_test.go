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
		// https://github.com/ziglang/www.ziglang.org/blob/main/check-mirrors/main.zig
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
		// https://github.com/ziglang/www.ziglang.org/blob/main/check-mirrors/main.zig
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
