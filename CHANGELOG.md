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

### Fixed
- Fixed newline convention CRLF -> LF
- The application now correctly exits if the server fails to start.
- Servers (http, tls, redirect) start and stop correctly

### Changed
- The package Flags renamed to Config.
