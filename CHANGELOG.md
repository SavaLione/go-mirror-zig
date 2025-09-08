# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- Added draft working version
- Added the project description (README.md)
- An internal package for parsing command-line flags.
- Comments to the Config package.
- Now it is possible to start a secure server (TLS, HTTPS).
- New config parameters that are responsible for HTTP/HTTPS ports, TLS cert/key files, redirect behavior.
- Ability to redirect incoming requests from http-port (default 80) to tls-port (default 443).
- The server can act as a caching upstream mirror. Now it allows serving files
  on `/builds/*` and `/download/*/*` paths. Can be useful if you want to daisy
  chain server.
- ACME challenges support.

### Fixed
- Fixed newline convention CRLF -> LF
- The application now correctly exits if the server fails to start.
- Servers (http, tls, redirect) start and stop correctly
- Files that have to be stored in the cache are downloaded safely right now.

### Changed
- The package Flags renamed to Config.
- Cache handler is safer.
- Handlers log more information about requests and server tasks.
- Reduced amount of logging information. Middleware and redirect handlers don't
  log any data anymore. Added 'source' field to the cache handler's logger.
- Cache directory layout for cached files. Previously all files were put in a
  root cache directory (that is set by `--cache-dir` flag), but now they follow
  the official Zig download layout. Dev builds are now stored in `/builds/` and
  release builds are stored in `/downloads/*/`.
