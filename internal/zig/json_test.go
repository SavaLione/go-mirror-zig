package zig

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestZigReleases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		inputJSON     string
		expectError   bool
		expected      ZigReleases
		expectedTotal int
	}{
		{
			name:          "empty json object",
			inputJSON:     `{}`,
			expectError:   false,
			expected:      ZigReleases{},
			expectedTotal: 0,
		},
		{
			name:        "invalid json",
			inputJSON:   `{ not valid json }`,
			expectError: true,
			expected:    nil,
		},
		{
			name: "dev branch, no artifacts",
			inputJSON: `{
  				"master": {
  				  "version": "0.17.0-dev.305+bdfbf432d",
  				  "date": "2026-05-15",
  				  "docs": "https://ziglang.org/documentation/master/",
  				  "stdDocs": "https://ziglang.org/documentation/master/std/"
  				}
			}`,
			expectError: false,
			expected: ZigReleases{
				"master": {
					Version:   "0.17.0-dev.305+bdfbf432d",
					Date:      "2026-05-15",
					Docs:      "https://ziglang.org/documentation/master/",
					StdDocs:   "https://ziglang.org/documentation/master/std/",
					Platforms: map[string]Artifact{},
				},
			},
		},
		{
			name: "dev branch, single artifact",
			inputJSON: `{
  				"master": {
  				  "version": "0.17.0-dev.305+bdfbf432d",
  				  "date": "2026-05-15",
  				  "docs": "https://ziglang.org/documentation/master/",
  				  "stdDocs": "https://ziglang.org/documentation/master/std/",
  				  "src": {
  				    "tarball": "https://ziglang.org/builds/zig-0.17.0-dev.305+bdfbf432d.tar.xz",
  				    "shasum": "f4e02500223c65225cb98651d32f744ee1f8939db05c4718598f1d89eddaa5dd",
  				    "size": "22530120"
  				  }
  				}
			}`,
			expectError: false,
			expected: ZigReleases{
				"master": {
					Version: "0.17.0-dev.305+bdfbf432d",
					Date:    "2026-05-15",
					Docs:    "https://ziglang.org/documentation/master/",
					StdDocs: "https://ziglang.org/documentation/master/std/",
					Platforms: map[string]Artifact{
						"src": {
							Tarball: "https://ziglang.org/builds/zig-0.17.0-dev.305+bdfbf432d.tar.xz",
							Shasum:  "f4e02500223c65225cb98651d32f744ee1f8939db05c4718598f1d89eddaa5dd",
							Size:    "22530120",
						},
					},
				},
			},
			expectedTotal: 22530120,
		},
		{
			name: "dev branch, multiple artifacts",
			inputJSON: `{
  				"master": {
  				  "version": "0.17.0-dev.305+bdfbf432d",
  				  "date": "2026-05-15",
  				  "docs": "https://ziglang.org/documentation/master/",
  				  "stdDocs": "https://ziglang.org/documentation/master/std/",
  				  "src": {
  				    "tarball": "https://ziglang.org/builds/zig-0.17.0-dev.305+bdfbf432d.tar.xz",
  				    "shasum": "f4e02500223c65225cb98651d32f744ee1f8939db05c4718598f1d89eddaa5dd",
  				    "size": "22530120"
  				  },
    			  "bootstrap": {
    			    "tarball": "https://ziglang.org/builds/zig-bootstrap-0.17.0-dev.305+bdfbf432d.tar.xz",
    			    "shasum": "f458a74d58561e185b69527a80b7f9aeeb721d6f824e293286a48c6fdb047f52",
    			    "size": "56523612"
    			  }
  				}
			}`,
			expectError: false,
			expected: ZigReleases{
				"master": {
					Version: "0.17.0-dev.305+bdfbf432d",
					Date:    "2026-05-15",
					Docs:    "https://ziglang.org/documentation/master/",
					StdDocs: "https://ziglang.org/documentation/master/std/",
					Platforms: map[string]Artifact{
						"src": {
							Tarball: "https://ziglang.org/builds/zig-0.17.0-dev.305+bdfbf432d.tar.xz",
							Shasum:  "f4e02500223c65225cb98651d32f744ee1f8939db05c4718598f1d89eddaa5dd",
							Size:    "22530120",
						},
						"bootstrap": {
							Tarball: "https://ziglang.org/builds/zig-bootstrap-0.17.0-dev.305+bdfbf432d.tar.xz",
							Shasum:  "f458a74d58561e185b69527a80b7f9aeeb721d6f824e293286a48c6fdb047f52",
							Size:    "56523612",
						},
					},
				},
			},
			expectedTotal: 22530120 + 56523612,
		},
		{
			name: "dev branch, omit artifacts without tarballs",
			inputJSON: `{
  				"master": {
  				  "version": "0.17.0-dev.305+bdfbf432d",
  				  "date": "2026-05-15",
  				  "docs": "https://ziglang.org/documentation/master/",
  				  "stdDocs": "https://ziglang.org/documentation/master/std/",
  				  "src": {
  				    "shasum": "f4e02500223c65225cb98651d32f744ee1f8939db05c4718598f1d89eddaa5dd",
  				    "size": "22530120"
  				  },
    			  "bootstrap": {
    			    "shasum": "f458a74d58561e185b69527a80b7f9aeeb721d6f824e293286a48c6fdb047f52",
    			    "size": "56523612"
    			  }
  				}
			}`,
			expectError: false,
			expected: ZigReleases{
				"master": {
					Version:   "0.17.0-dev.305+bdfbf432d",
					Date:      "2026-05-15",
					Docs:      "https://ziglang.org/documentation/master/",
					StdDocs:   "https://ziglang.org/documentation/master/std/",
					Platforms: map[string]Artifact{},
				},
			},
		},
		{
			name: "dev branch, single artifact, wrong size",
			inputJSON: `{
  				"master": {
  				  "version": "0.17.0-dev.305+bdfbf432d",
  				  "date": "2026-05-15",
  				  "docs": "https://ziglang.org/documentation/master/",
  				  "stdDocs": "https://ziglang.org/documentation/master/std/",
  				  "src": {
  				    "tarball": "https://ziglang.org/builds/zig-0.17.0-dev.305+bdfbf432d.tar.xz",
  				    "shasum": "f4e02500223c65225cb98651d32f744ee1f8939db05c4718598f1d89eddaa5dd",
  				    "size": "22530120"
  				  },
    			  "bootstrap": {
    			    "tarball": "https://ziglang.org/builds/zig-bootstrap-0.17.0-dev.305+bdfbf432d.tar.xz",
    			    "shasum": "f458a74d58561e185b69527a80b7f9aeeb721d6f824e293286a48c6fdb047f52",
    			    "size": "it-is-not-a-number"
    			  }
  				}
			}`,
			expectError: false,
			expected: ZigReleases{
				"master": {
					Version: "0.17.0-dev.305+bdfbf432d",
					Date:    "2026-05-15",
					Docs:    "https://ziglang.org/documentation/master/",
					StdDocs: "https://ziglang.org/documentation/master/std/",
					Platforms: map[string]Artifact{
						"src": {
							Tarball: "https://ziglang.org/builds/zig-0.17.0-dev.305+bdfbf432d.tar.xz",
							Shasum:  "f4e02500223c65225cb98651d32f744ee1f8939db05c4718598f1d89eddaa5dd",
							Size:    "22530120",
						},
						"bootstrap": {
							Tarball: "https://ziglang.org/builds/zig-bootstrap-0.17.0-dev.305+bdfbf432d.tar.xz",
							Shasum:  "f458a74d58561e185b69527a80b7f9aeeb721d6f824e293286a48c6fdb047f52",
							Size:    "it-is-not-a-number",
						},
					},
				},
			},
			expectedTotal: 22530120,
		},

		{
			name: "dev branch and a release, multiple artifacts",
			inputJSON: `{
  				"master": {
  				  "version": "0.17.0-dev.305+bdfbf432d",
  				  "date": "2026-05-15",
  				  "docs": "https://ziglang.org/documentation/master/",
  				  "stdDocs": "https://ziglang.org/documentation/master/std/",
  				  "src": {
  				    "tarball": "https://ziglang.org/builds/zig-0.17.0-dev.305+bdfbf432d.tar.xz",
  				    "shasum": "f4e02500223c65225cb98651d32f744ee1f8939db05c4718598f1d89eddaa5dd",
  				    "size": "22530120"
  				  },
    			  "bootstrap": {
    			    "tarball": "https://ziglang.org/builds/zig-bootstrap-0.17.0-dev.305+bdfbf432d.tar.xz",
    			    "shasum": "f458a74d58561e185b69527a80b7f9aeeb721d6f824e293286a48c6fdb047f52",
    			    "size": "56523612"
    			  }
  				},
                "0.16.0": {
                  "version": "0.16.0",
                  "date": "2026-04-13",
                  "docs": "https://ziglang.org/documentation/0.16.0/",
                  "stdDocs": "https://ziglang.org/documentation/0.16.0/std/",
                  "notes": "https://ziglang.org/download/0.16.0/release-notes.html",
                  "src": {
                    "tarball": "https://ziglang.org/download/0.16.0/zig-0.16.0.tar.xz",
                    "shasum": "43186959edc87d5c7a1be7b7d2a25efffd22ce5807c7af99067f86f99641bfdf",
                    "size": "22503260"
                  },
                  "bootstrap": {
                    "tarball": "https://ziglang.org/download/0.16.0/zig-bootstrap-0.16.0.tar.xz",
                    "shasum": "2a8266a4205772ef40838c8cbdf14875855a515ff3adf89b49c2d2ae93613d10",
                    "size": "55245980"
                  }
                }
			}`,
			expectError: false,
			expected: ZigReleases{
				"master": {
					Version: "0.17.0-dev.305+bdfbf432d",
					Date:    "2026-05-15",
					Docs:    "https://ziglang.org/documentation/master/",
					StdDocs: "https://ziglang.org/documentation/master/std/",
					Platforms: map[string]Artifact{
						"src": {
							Tarball: "https://ziglang.org/builds/zig-0.17.0-dev.305+bdfbf432d.tar.xz",
							Shasum:  "f4e02500223c65225cb98651d32f744ee1f8939db05c4718598f1d89eddaa5dd",
							Size:    "22530120",
						},
						"bootstrap": {
							Tarball: "https://ziglang.org/builds/zig-bootstrap-0.17.0-dev.305+bdfbf432d.tar.xz",
							Shasum:  "f458a74d58561e185b69527a80b7f9aeeb721d6f824e293286a48c6fdb047f52",
							Size:    "56523612",
						},
					},
				},
				"0.16.0": {
					Version: "0.16.0",
					Date:    "2026-04-13",
					Docs:    "https://ziglang.org/documentation/0.16.0/",
					StdDocs: "https://ziglang.org/documentation/0.16.0/std/",
					Notes:   "https://ziglang.org/download/0.16.0/release-notes.html",
					Platforms: map[string]Artifact{
						"src": {
							Tarball: "https://ziglang.org/download/0.16.0/zig-0.16.0.tar.xz",
							Shasum:  "43186959edc87d5c7a1be7b7d2a25efffd22ce5807c7af99067f86f99641bfdf",
							Size:    "22503260",
						},
						"bootstrap": {
							Tarball: "https://ziglang.org/download/0.16.0/zig-bootstrap-0.16.0.tar.xz",
							Shasum:  "2a8266a4205772ef40838c8cbdf14875855a515ff3adf89b49c2d2ae93613d10",
							Size:    "55245980",
						},
					},
				},
			},
			expectedTotal: 22530120 + 56523612 + 22503260 + 55245980,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var zr ZigReleases
			err := json.Unmarshal([]byte(tt.inputJSON), &zr)

			if (err != nil) != tt.expectError {
				t.Fatalf("expected error: %v, got: %v", tt.expectError, err)
			}

			if tt.expectError {
				return
			}

			if !reflect.DeepEqual(zr, tt.expected) {
				t.Errorf("\nExpected: %#v\nGot:      %#v", tt.expected, zr)
			}

			actualTotal := TotalSize(zr)
			if actualTotal != tt.expectedTotal {
				t.Errorf("TotalSize: expected %d, got %d", tt.expectedTotal, actualTotal)
			}
		})
	}
}
