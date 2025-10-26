# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.2.0] - 2025-10-26

### Added

- **New Commands:**
  - `mayrlabs version` - Display the version of the CLI
  - `mayrlabs visit` - Open the MayR Labs website (https://mayrlabs.com) in your browser
  - `mayrlabs github` - Open the GitHub repository in your browser

- **Installation Improvements:**
  - Added `install.sh` script for easy installation via curl or wget
  - Installation script automatically detects OS and architecture
  - Script supports Linux (AMD64/ARM64), macOS (Intel/Apple Silicon), and Windows

- **Release Workflow:**
  - Enhanced release workflow with version injection in binaries
  - Automatic generation of SHA256 checksums for all release artifacts
  - Comprehensive release notes with installation instructions
  - Support for multiple platforms: Linux (AMD64/ARM64), macOS (Intel/Apple Silicon), Windows (AMD64)

- **Documentation:**
  - Added CHANGELOG.md to track version history
  - Updated README.md with installation script usage
  - Added instructions for configuring $GOPATH/bin in PATH
  - Improved installation documentation with multiple methods

### Changed

- Updated release workflow to use Go 1.24 and enhanced build process
- Improved binary naming and installation process
- Enhanced release notes format with detailed installation instructions

### Fixed

- Fixed binary name issue when using `go install` (documented workaround)

## [v0.1.0] - Initial Release

### Added

- Initial release of MayR Labs CLI
- Core commands for development workflow automation:
  - UUID generation
  - Password generation
  - Hash generation (MD5, SHA1, SHA256)
  - Keystore creation
  - DNS cache clearing
  - CI/CD workflow generation
  - Code formatting
  - License file creation
  - Editor config generation
  - Motivational quotes
- Git operations (prune stale branches)
- Environment file management (validate, update-example, arrange)
- Changelog management (create, record)
- Flutter build scripts generation
- PHP code quality tools (cs-fix, lint)
- JavaScript Prettier setup and formatting
- Cross-platform support (macOS, Linux, Windows)
- Built with Go and Cobra framework

[v0.2.0]: https://github.com/MayR-Labs/mayrlabs-go/releases/tag/v0.2.0
[v0.1.0]: https://github.com/MayR-Labs/mayrlabs-go/releases/tag/v0.1.0
