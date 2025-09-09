# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
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
