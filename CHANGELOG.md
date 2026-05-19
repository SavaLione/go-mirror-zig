# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.6] - 2026-05-19
### Added
- Added test suites for the index, middleware, redirect, and configuration parser handlers.
- Added unit tests for Zig internal logic in `internal/zig/zig.go`.
- Added a deployment subsection to the README for using Nginx as a reverse proxy.
- Added a direct link to the Zig Programming Language website in the README header.
- Added functionality to fetch and parse the upstream `index.json` for release processing.
- Added a function to calculate the total cacheable size of all upstream Zig artifacts.
- Added the `-show-possible-size` flag to display upstream artifact statistics (total size, release counts, and median sizes) before exiting.
- Added intelligent periodic cache cleanup that cross-references the local cache with the upstream `index.json` to preserve active `master` builds while removing stale ones.
- Added aggregate logging for the cleanup task, reporting the total number of files removed and space reclaimed (in bytes).
- Added new targets to the build script:
    - `linux/riscv64`
    - `linux/386`
    - `linux/ppc64le`
    - `linux/loong64`
    - `windows/arm64`
    - `freebsd/amd64`
    - `freebsd/arm64`
    - `netbsd/amd64`
    - `netbsd/arm64`
    - `openbsd/amd64`
    - `openbsd/arm64`


### Fixed
- Fixed inconsistent behavior of temporary files across platforms (Windows vs. Linux file locking). On Linux it is allowed to delete an opened file, while on Windows it is not allowed
- Added a warning message when ACME Terms of Service are not explicitly accepted.
- Removed an unused version string from the `config` package.
- Corrected the documentation regarding the default interval for cleaning up cached dev builds.
- Fixed grammar in the changelog.

### Changed
- Replaced GitHub links with Codeberg links for the official Zig main repository.
- The HTTP to HTTPS redirect now omits the `443` port if it's not required
- Reduced the interval in seconds to clean up cached dev builds (from 1 day to 2 hours), thanks to the intelligent periodic cache cleanup feature.
- Updated table driven tests for artifact regex matching.
- Improved the configuration parser tests.
- Updated readme.

## [1.1.0] - 2026-03-14
### Added
- Added `-show-index-page` and `-index-page` flags to allow users to customize or disable the root index page.
- Added `-clear-builds-interval` flag to automatically clean up stale Zig development artifacts from the cache.
- Added background task to periodically remove cached dev builds based on the configured interval.
- Added a new `zig` internal package to centralize artifact validation and parsing logic.
- Added server version and system architecture information to the index page.

### Security
- Updated `golang.org/x/crypto` library (`v0.41.0` -> `v0.49.0`) (thanks dependabot for that)

### Fixed
- Fixed minimal required Go version in the README.md.

### Changed
- Changed project's Go version (`1.25.0` -> `1.26.1`)

## [1.0.0] - 2025-09-09
### Added
- Initial project structure and basic mirroring functionality.
- Added README.md with project description and usage instructions.
- Added an internal package for parsing command-line flags.
- Comments to the Config package.
- Added support for serving content over HTTPS (TLS).
- Added command-line flags for configuring ports, TLS certificates, and HTTP-to-HTTPS redirection.
- Added optional redirection of all HTTP traffic to HTTPS.
- Implemented caching for official Zig download paths (`/builds/` and `/download/*/*`)
- Added ACME support for automatic TLS certificate acquisition from Let's Encrypt.
- Added a 'Features' section to README.md.
- Added new sections in the readme: Getting Started, Examples, Configuration, Deployment, Contributing, Licenses and Acknowledgements.
- Added version information flag support, updated readme.
- Added a make script for building the application.
- Added a build script for cross-compiling releases.
- Compiled releases are not tracked by git.
- Added 'From a Precompiled Binary' section in readme.

### Fixed
- Fixed newline convention CRLF -> LF
- The application now correctly exits if the server fails to start.
- Improved server lifecycle management for graceful startup and shutdown.
- Ensured atomic writes to the cache to prevent serving partially downloaded files.
- Fixed grammar in the changelog.

### Changed
- Renamed the internal `flags` package to `config` for clarity.
- Improved error handling and robustness of the cache handler.
- Handlers log more information about requests and server tasks.
- Refactored logging to reduce noise and added a source field (e.g., 'cache' or 'upstream') to request logs.
- Changed the cache directory layout to mirror the official Zig download structure (`/builds/` and `/download/*/*`).
- Changed the appearance of the index page.
- Changed the appearance of the readme.

### Removed
- Removed notice about the project being a draft.
