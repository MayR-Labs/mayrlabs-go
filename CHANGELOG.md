# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.4.0] - 2025-10-26

### Added

- **Survey Library Integration:**
  - Integrated `github.com/AlecAivazis/survey/v2` for enhanced interactive prompts
  - All interactive commands now use survey fields for better user experience
  - Interactive modes automatically prompt for flags like `--copy`, `--force`, `--wip`

- **Enhanced Existing Commands:**
  - `add-license`: Added optional `--author-url` flag to include author URL in licenses
  - `add-license`: Interactive mode now asks for year with current year as default
  - `editor-config`: Now fully interactive, prompts for language if not specified
  - `hash-file`: Added interactive mode when no arguments provided
  - `uuid`: Added `--copy` flag to copy UUID to clipboard
  - `ulid`: Added `--copy` flag to copy ULID to clipboard
  - `env validate`: Now also checks for keys in `.env.example` that are missing in `.env`
  - `changelog record`: Added separator line (----) between changelog versions
  - `changelog create`: Interactive mode asks for force flag if file exists
  - `changelog record`: Interactive mode asks for wip flag

- **New Commands:**
  - `mayrlabs roll-dice [n]` - Roll n dice (1-100) and display results with total and average
  - `mayrlabs alias [name]` - Create a permanent shell alias for mayrlabs command
  - `mayrlabs upgrade` - Upgrade mayrlabs to the latest version
  - `mayrlabs base64 encode/decode [string]` - Encode or decode base64 strings with `--copy` flag
  - `mayrlabs base64-file [path]` - Encode a file to base64 with `--copy` flag
  - `mayrlabs base64-decode-to-file [string]` - Decode base64 string and write to file

### Changed

- `create-keystore`: Now uses PKCS12 format instead of JKS (industry standard)
- `create-keystore`: Better error handling when filename is not provided
- All interactive prompts upgraded to use survey library for better UX
- Interactive modes now prompt for additional flags when not specified via command line

### Fixed

- Fixed `changelog create` - interactive mode now properly handles force flag
- Fixed `changelog record` - interactive mode now properly handles wip flag
- Fixed `create-keystore` - removed JKS deprecation warning by switching to PKCS12
- Fixed `create-keystore` - proper error message when filename is empty

----
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

----
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

[v0.4.0]: https://github.com/MayR-Labs/mayrlabs-go/releases/tag/v0.4.0
[v0.2.0]: https://github.com/MayR-Labs/mayrlabs-go/releases/tag/v0.2.0
[v0.1.0]: https://github.com/MayR-Labs/mayrlabs-go/releases/tag/v0.1.0
