# MayR Labs CLI - Architecture Documentation

This document provides an in-depth look at the architecture, design decisions, and internal structure of the MayR Labs CLI.

## Table of Contents

1. [Overview](#overview)
2. [Project Structure](#project-structure)
3. [Core Architecture](#core-architecture)
4. [Command System](#command-system)
5. [Module Organization](#module-organization)
6. [Design Patterns](#design-patterns)
7. [Data Flow](#data-flow)
8. [External Dependencies](#external-dependencies)
9. [Security Considerations](#security-considerations)
10. [Testing Strategy](#testing-strategy)
11. [Build & Release](#build--release)
12. [Future Architecture Considerations](#future-architecture-considerations)

---

## Overview

MayR Labs CLI is a command-line tool built with Go that provides a unified interface for common development tasks. The architecture follows Go best practices and emphasizes:

- **Modularity**: Clear separation of concerns
- **Extensibility**: Easy to add new commands
- **Testability**: Comprehensive test coverage
- **Cross-platform**: Works on macOS, Linux, and Windows
- **Single Binary**: No external dependencies at runtime

### Technology Stack

- **Language**: Go 1.21+
- **CLI Framework**: [Cobra](https://github.com/spf13/cobra)
- **Interactive Prompts**: [Survey](https://github.com/AlecAivazis/survey)
- **AI Integration**: [Google Generative AI Go SDK](https://github.com/google/generative-ai-go)
- **Build Tool**: Go toolchain + Makefile

---

## Project Structure

```
mayrlabs-go/
├── main.go                 # Application entry point
├── cmd/                    # Command registration
│   └── root.go            # Root command and command tree
├── internal/              # Private application code
│   ├── commands/          # Command implementations
│   │   ├── general.go    # General utility commands
│   │   ├── system.go     # System-level commands
│   │   ├── git.go        # Git operations
│   │   ├── env.go        # Environment file management
│   │   ├── changelog.go  # Changelog management
│   │   ├── flutter.go    # Flutter-specific commands
│   │   ├── php.go        # PHP-specific commands
│   │   ├── js.go         # JavaScript-specific commands
│   │   ├── ai.go         # AI integration commands
│   │   ├── session.go    # Session management
│   │   ├── alias.go      # Alias management
│   │   ├── base64.go     # Encoding/decoding
│   │   ├── browser.go    # Browser operations
│   │   ├── ci.go         # CI/CD generation
│   │   ├── dice.go       # Random dice rolling
│   │   ├── format.go     # Code formatting
│   │   ├── license.go    # License generation
│   │   ├── quote.go      # Motivational quotes
│   │   ├── utils.go      # Utility functions
│   │   └── version.go    # Version information
│   └── utils/            # Shared utilities
│       └── ai_test.go    # AI utility tests
├── examples/             # Usage examples
│   └── README.md
├── .github/              # GitHub configuration
│   └── workflows/        # CI/CD workflows
│       ├── ci.yml        # Continuous integration
│       └── release.yml   # Release automation
├── go.mod                # Go module definition
├── go.sum                # Dependency checksums
├── Makefile              # Build automation
├── README.md             # User documentation
├── API.md                # API reference
├── ARCHITECTURE.md       # This file
├── CONTRIBUTING.md       # Contribution guidelines
├── DEVELOPMENT.md        # Developer guide
├── CHANGELOG.md          # Version history
├── LICENSE               # MIT License
└── install.sh           # Installation script
```

---

## Core Architecture

### Application Flow

```
main.go
    ↓
cmd/root.go (Cobra root command)
    ↓
Command Registration (init())
    ↓
Command Execution (internal/commands/)
    ↓
Utility Functions (internal/utils/)
    ↓
External Services (filesystem, network, AI API)
```

### Entry Point

**File**: `main.go`

```go
func main() {
    if err := cmd.Execute(); err != nil {
        os.Exit(1)
    }
}
```

The entry point is minimal, delegating all work to the `cmd` package.

### Root Command

**File**: `cmd/root.go`

The root command:
- Defines the CLI name and description
- Registers all subcommands in `init()`
- Handles the help display when no command is provided
- Uses Cobra's command tree structure

```go
var rootCmd = &cobra.Command{
    Use:   "mayrlabs",
    Short: "🧰 MayR Labs CLI - Streamline your development workflow",
    // ...
}

func init() {
    // Register all commands
    rootCmd.AddCommand(commands.UUIDCmd)
    rootCmd.AddCommand(commands.PasswordCmd)
    // ... more commands
}
```

---

## Command System

### Command Structure

Each command follows the Cobra command pattern:

```go
var ExampleCmd = &cobra.Command{
    Use:   "example [args]",
    Short: "Short description",
    Long:  "Long description with details",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Command implementation
        return nil
    },
}

func init() {
    // Register flags
    ExampleCmd.Flags().StringP("flag", "f", "default", "flag description")
}
```

### Command Categories

Commands are organized into logical groups:

1. **General Commands** (`general.go`, `utils.go`)
   - UUID/ULID generation
   - Password generation
   - Hashing
   - Random number generation

2. **System Commands** (`system.go`, `alias.go`)
   - DNS cache clearing
   - System upgrades
   - Alias management

3. **Development Tools** (`ci.go`, `format.go`, `license.go`)
   - CI/CD generation
   - Code formatting
   - License creation
   - Editor config

4. **Language-Specific** (`flutter.go`, `php.go`, `js.go`)
   - Flutter build scripts
   - PHP code quality tools
   - JavaScript Prettier setup

5. **Version Control** (`git.go`)
   - Branch management
   - Stale branch pruning

6. **Project Management** (`env.go`, `changelog.go`)
   - Environment file management
   - Changelog maintenance

7. **AI Features** (`ai.go`)
   - AI query interface
   - File analysis
   - API key management

8. **Session Management** (`session.go`)
   - Development sessions
   - Encrypted sessions
   - Session history

### Interactive vs Non-Interactive

Commands support both modes:

**Non-Interactive**: All parameters provided via flags
```bash
mayrlabs add-license --type mit --author "John" --year 2025
```

**Interactive**: Missing parameters trigger prompts
```bash
mayrlabs add-license
# Prompts for: type, author, year
```

Implementation uses Survey library:

```go
import "github.com/AlecAivazis/survey/v2"

func askForInput() (string, error) {
    var result string
    prompt := &survey.Input{
        Message: "Enter value:",
    }
    return result, survey.AskOne(prompt, &result)
}
```

---

## Module Organization

### Internal Package

The `internal/` directory ensures code is not importable by external projects, following Go best practices.

#### Commands Module (`internal/commands/`)

Each file typically contains:
- Command definition (Cobra command)
- Command implementation (RunE function)
- Helper functions specific to that command
- Flag definitions

**Example**: `uuid.go`
```go
package commands

import (
    "fmt"
    "github.com/google/uuid"
    "github.com/spf13/cobra"
)

var UUIDCmd = &cobra.Command{
    Use:   "uuid",
    Short: "Generate UUID v4",
    RunE: func(cmd *cobra.Command, args []string) error {
        id := uuid.New()
        fmt.Println(id.String())
        return nil
    },
}
```

#### Utils Module (`internal/utils/`)

Shared utilities used across commands:
- File operations
- Configuration management
- API key storage
- Common prompts

---

## Design Patterns

### 1. Command Pattern

Each command encapsulates a specific action, following the Command Pattern. Commands are self-contained and can be composed.

### 2. Factory Pattern

Commands are created and registered in the `init()` function, acting as a factory for command instances.

### 3. Strategy Pattern

Different formatters (Go, JavaScript, Python) implement a common interface, allowing selection at runtime.

### 4. Template Method

Many commands follow a template:
1. Parse arguments/flags
2. Validate input
3. Execute action
4. Handle errors
5. Provide feedback

### 5. Builder Pattern

Complex configurations (CI/CD files, editor configs) use a builder-like approach to construct output.

---

## Data Flow

### Simple Command Flow

```
User Input → Cobra Parsing → Command.RunE → Action → Output
```

### Interactive Command Flow

```
User Input → Cobra Parsing → Missing Params? → Survey Prompts
    ↓
Input Collected → Validation → Action → Output
```

### AI Command Flow

```
User Query → API Key Check → Gemini API Call → Response Processing → Output
```

### Session Command Flow

```
Session Start → User Interactions → AI Queries → Note Taking
    ↓
Session End → Save to File → Display Summary
```

---

## External Dependencies

### Core Dependencies

1. **Cobra** (`github.com/spf13/cobra`)
   - Purpose: CLI framework
   - Usage: Command definition and parsing
   - Why: Industry standard, excellent documentation

2. **Survey** (`github.com/AlecAivazis/survey/v2`)
   - Purpose: Interactive prompts
   - Usage: User input collection
   - Why: Rich terminal UI, great UX

3. **Google Generative AI** (`github.com/google/generative-ai-go`)
   - Purpose: AI integration
   - Usage: Gemini API access
   - Why: Official Google SDK, reliable

4. **UUID** (`github.com/google/uuid`)
   - Purpose: UUID generation
   - Usage: UUID v4 creation
   - Why: Standard implementation

5. **ULID** (`github.com/oklog/ulid/v2`)
   - Purpose: ULID generation
   - Usage: Sortable identifier creation
   - Why: Standard implementation

6. **Clipboard** (`github.com/atotto/clipboard`)
   - Purpose: Clipboard operations
   - Usage: Copy output to clipboard
   - Why: Cross-platform support

### Dependency Management

Dependencies are managed via Go modules:
- `go.mod`: Direct dependencies
- `go.sum`: Checksums for reproducible builds
- Minimal dependency tree to reduce attack surface

---

## Security Considerations

### 1. Input Validation

All user inputs are validated before processing:
```go
if len(input) == 0 {
    return fmt.Errorf("input cannot be empty")
}
```

### 2. Secure Password Generation

Uses `crypto/rand` for cryptographically secure random generation:
```go
import "crypto/rand"

func generatePassword(length int) (string, error) {
    bytes := make([]byte, length)
    _, err := rand.Read(bytes)
    // ...
}
```

### 3. API Key Storage

API keys are stored in plain text in `~/.mayrlabs/gemini-api-key` with:
- File permissions: 0600 (read/write owner only)
- Location: User home directory
- Warning: Users should protect their home directory

**Improvement Opportunity**: Use OS keychain/keyring for more secure storage.

### 4. Encrypted Sessions

Secure sessions use AES-256-GCM encryption:
- User-provided password as key
- Random nonce per session
- Authenticated encryption

### 5. Command Execution

When executing system commands (e.g., `git`, `npm`), the CLI:
- Uses `exec.Command` with explicit paths when possible
- Validates command output
- Handles errors appropriately
- Doesn't execute arbitrary user input as shell commands

### 6. File Operations

File operations:
- Check existence before overwriting (unless `--force`)
- Validate paths to prevent directory traversal
- Use appropriate file permissions
- Handle errors gracefully

---

## Testing Strategy

### Unit Tests

Located alongside implementation files with `_test.go` suffix.

**Coverage Goals:**
- Core utilities: >80%
- Command logic: >50%
- Overall: >50%

**Testing Approach:**
```go
func TestGeneratePassword(t *testing.T) {
    tests := []struct {
        name   string
        length int
        want   int
    }{
        {"16 chars", 16, 16},
        {"32 chars", 32, 32},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := generatePassword(tt.length)
            if err != nil {
                t.Errorf("unexpected error: %v", err)
            }
            if len(result) != tt.want {
                t.Errorf("got %d, want %d", len(result), tt.want)
            }
        })
    }
}
```

### Table-Driven Tests

Most tests use table-driven approach for comprehensive coverage:
```go
tests := []struct {
    name    string
    input   string
    want    string
    wantErr bool
}{
    // Test cases
}
```

### Integration Tests

Currently minimal; focus on unit tests. Future improvement opportunity.

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
make coverage

# Run with race detection
go test -race ./...
```

---

## Build & Release

### Build Process

**Makefile** defines common tasks:

```makefile
build:
    go build -o mayrlabs main.go

test:
    go test -v -race ./...

lint:
    golangci-lint run

build-all:
    # Multi-platform builds
    GOOS=linux GOARCH=amd64 go build -o dist/mayrlabs-linux-amd64
    GOOS=darwin GOARCH=arm64 go build -o dist/mayrlabs-darwin-arm64
    # ... more platforms
```

### CI/CD Pipeline

**GitHub Actions** (`.github/workflows/ci.yml`):

1. **On Push/PR**:
   - Checkout code
   - Set up Go
   - Install dependencies
   - Run tests with coverage
   - Run linters
   - Build for multiple platforms

2. **On Release Tag** (`.github/workflows/release.yml`):
   - Build binaries for all platforms:
     - Linux: amd64, arm64
     - macOS: amd64 (Intel), arm64 (Apple Silicon)
     - Windows: amd64
   - Generate SHA256 checksums
   - Create GitHub release
   - Upload binaries as release assets
   - Generate release notes

### Version Management

Version is defined in `internal/commands/version.go`:

```go
var Version = "1.0.0"
```

During build, it can be injected via ldflags:
```bash
go build -ldflags "-X github.com/MayR-Labs/mayrlabs-go/internal/commands.Version=1.2.3"
```

### Release Process

1. Update version in `version.go`
2. Update `CHANGELOG.md`
3. Commit changes
4. Create and push Git tag: `git tag v1.0.0 && git push origin v1.0.0`
5. GitHub Actions automatically builds and releases

---

## Future Architecture Considerations

### 1. Plugin System

**Goal**: Allow third-party command extensions

**Approach**:
- Define plugin interface
- Load plugins from `~/.mayrlabs/plugins/`
- Use Go's plugin package or external binaries

### 2. Configuration File

**Goal**: User-customizable defaults

**Approach**:
- YAML configuration at `~/.mayrlabs/config.yaml`
- Override command defaults
- Store preferences

**Example**:
```yaml
defaults:
  license:
    type: mit
    author: "John Doe"
  ci:
    vcs: github
```

### 3. Remote Command Execution

**Goal**: Execute commands on remote servers

**Approach**:
- SSH integration
- Remote session management
- Secure credential storage

### 4. Web UI

**Goal**: Optional web interface for certain features

**Approach**:
- Embedded web server
- Session visualization
- AI chat interface

### 5. Shell Completion Enhancement

**Goal**: Rich shell completions for all commands

**Approach**:
- Use Cobra's completion system
- Dynamic completions for file paths, branches, etc.
- Install scripts for all shells

### 6. Package Manager Integration

**Goal**: Distribute via package managers

**Approach**:
- Homebrew formula
- APT/YUM repositories
- Chocolatey package
- Scoop manifest

### 7. API Server Mode

**Goal**: Run as a background service with HTTP API

**Approach**:
- Daemon mode
- RESTful API
- WebSocket support for sessions

### 8. Enhanced Security

**Goal**: Improve secret management

**Approach**:
- OS keychain integration (macOS Keychain, Windows Credential Manager, Linux Secret Service)
- Optional GPG encryption for configs
- Audit logging

### 9. Telemetry

**Goal**: Optional usage analytics

**Approach**:
- Opt-in telemetry
- Privacy-respecting
- Help prioritize features

### 10. Internationalization

**Goal**: Support multiple languages

**Approach**:
- i18n framework
- Language files
- Locale detection

---

## Code Quality & Standards

### Linting

Uses `golangci-lint` with configuration in `.golangci.yml`:

```yaml
linters:
  enable:
    - gofmt
    - govet
    - staticcheck
    - errcheck
    - gosimple
    - ineffassign
```

### Code Style

- Follow Go idioms and conventions
- Use `gofmt` for formatting
- Meaningful variable names
- Comments for exported functions
- Error handling over panic
- Interfaces over concrete types where appropriate

### Documentation

- Code comments for all exported functions
- README for users
- API.md for reference
- ARCHITECTURE.md (this document) for developers
- DEVELOPMENT.md for setup and contribution

---

## Performance Considerations

### Optimization Strategies

1. **Lazy Loading**: Commands only load dependencies when executed
2. **Minimal Allocations**: Reuse buffers where possible
3. **Concurrent Operations**: Use goroutines for independent tasks
4. **Caching**: Cache expensive operations (e.g., API calls)

### Binary Size

Current binary size: ~15-20 MB (including dependencies)

**Reduction strategies**:
- Strip debug info: `-ldflags "-s -w"`
- Use UPX compression (optional)
- Minimize dependencies

### Startup Time

Target: <100ms for simple commands

**Achieved by**:
- Minimal initialization
- Lazy command registration
- No global state initialization

---

## Error Handling

### Error Strategy

1. **Return errors, don't panic**
   ```go
   if err != nil {
       return fmt.Errorf("operation failed: %w", err)
   }
   ```

2. **Wrap errors with context**
   ```go
   return fmt.Errorf("failed to read file %s: %w", filename, err)
   ```

3. **User-friendly messages**
   ```go
   fmt.Fprintf(os.Stderr, "Error: %v\n", err)
   ```

4. **Exit codes**
   - 0: Success
   - 1: Error

---

## Logging

Currently uses simple `fmt` package for output.

**Future Enhancement**: Structured logging with levels (debug, info, warn, error)

---

## Cross-Platform Compatibility

### OS-Specific Code

Uses build tags for platform-specific implementations:

```go
// +build darwin

package commands

func clearDNSCache() error {
    // macOS implementation
}
```

### File Paths

Always use `filepath` package:
```go
import "path/filepath"

path := filepath.Join(home, ".mayrlabs", "config")
```

### Command Execution

Platform detection for system commands:
```go
import "runtime"

if runtime.GOOS == "windows" {
    // Windows command
} else {
    // Unix command
}
```

---

## Contribution Architecture

### Adding a New Command

1. Create new file in `internal/commands/` (e.g., `newcommand.go`)
2. Define Cobra command:
   ```go
   var NewCmd = &cobra.Command{
       Use:   "new",
       Short: "Description",
       RunE:  runNew,
   }
   
   func runNew(cmd *cobra.Command, args []string) error {
       // Implementation
       return nil
   }
   ```
3. Register in `cmd/root.go`:
   ```go
   rootCmd.AddCommand(commands.NewCmd)
   ```
4. Add tests in `newcommand_test.go`
5. Update documentation

### Code Review Checklist

- [ ] Tests pass
- [ ] Linter passes
- [ ] Documentation updated
- [ ] Error handling appropriate
- [ ] Cross-platform compatibility checked
- [ ] No hardcoded paths
- [ ] Follows existing patterns

---

## Debugging

### Debug Mode

Set environment variable:
```bash
export MAYRLABS_DEBUG=1
mayrlabs command
```

### Profiling

```bash
go build -o mayrlabs main.go
./mayrlabs command
go tool pprof mayrlabs cpu.prof
```

---

## Resources

### Internal Documentation
- [README.md](README.md) - User guide
- [API.md](API.md) - API reference
- [DEVELOPMENT.md](DEVELOPMENT.md) - Developer setup
- [CONTRIBUTING.md](CONTRIBUTING.md) - Contribution guidelines

### External Resources
- [Cobra Documentation](https://cobra.dev/)
- [Go Best Practices](https://go.dev/doc/effective_go)
- [Go Module Reference](https://go.dev/ref/mod)

---

**Last Updated**: 2025-10-27 (v1.0.0)

**Maintainers**: MayR Labs Team

**License**: MIT
