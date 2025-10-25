# Development Guide

This guide provides detailed information for developers who want to contribute to or understand the MayR Labs CLI codebase.

## Table of Contents

- [Project Architecture](#project-architecture)
- [Command Structure](#command-structure)
- [Adding New Commands](#adding-new-commands)
- [Testing Strategy](#testing-strategy)
- [Development Workflow](#development-workflow)
- [Common Tasks](#common-tasks)

## Project Architecture

### Directory Structure

```
.
├── .github/
│   ├── agents/              # Copilot agent configuration
│   └── workflows/           # CI/CD workflows
├── cmd/                     # CLI entry point and root command
│   └── root.go             # Cobra root command setup
├── internal/                # Private application code
│   ├── commands/           # Command implementations
│   │   ├── changelog.go    # Changelog commands
│   │   ├── ci.go           # CI/CD generation
│   │   ├── env.go          # Environment file management
│   │   ├── flutter.go      # Flutter-specific commands
│   │   ├── format.go       # Code formatting
│   │   ├── general.go      # General utility commands
│   │   ├── git.go          # Git operations
│   │   ├── js.go           # JavaScript commands
│   │   ├── license.go      # License generation
│   │   ├── php.go          # PHP commands
│   │   ├── quote.go        # Motivational quotes
│   │   └── utils.go        # Utility commands (hash, uuid, password)
│   └── utils/              # Shared utility functions
│       └── helpers.go      # Common helper functions
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
├── Makefile                # Build automation
└── README.md               # User documentation
```

### Key Dependencies

- **cobra** (`github.com/spf13/cobra`): CLI framework for building commands
- **uuid** (`github.com/google/uuid`): UUID generation

## Command Structure

All commands follow the Cobra command pattern:

```go
var MyCmd = &cobra.Command{
    Use:   "command-name [args]",
    Short: "Brief description",
    Long:  "Detailed description of what this command does",
    Args:  cobra.ExactArgs(1),  // Argument validation
    RunE: func(cmd *cobra.Command, args []string) error {
        // Command implementation
        return nil
    },
}

func init() {
    // Add flags
    MyCmd.Flags().StringP("flag", "f", "default", "Flag description")
}
```

### Command Categories

1. **General Commands**: Standalone utilities (hash, uuid, password, etc.)
2. **Group Commands**: Parent commands with subcommands (git, env, changelog, etc.)

## Adding New Commands

### Step 1: Create the Command File

Create a new file in `internal/commands/`:

```go
// internal/commands/myfeature.go
package commands

import (
    "fmt"
    "github.com/spf13/cobra"
)

var MyFeatureCmd = &cobra.Command{
    Use:   "myfeature",
    Short: "Short description",
    Long:  "Long description",
    RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Println("My feature works!")
        return nil
    },
}
```

### Step 2: Register the Command

Add it to `cmd/root.go`:

```go
func init() {
    // ... existing commands ...
    rootCmd.AddCommand(commands.MyFeatureCmd)
}
```

### Step 3: Add Tests

Create `internal/commands/myfeature_test.go`:

```go
package commands

import (
    "testing"
)

func TestMyFeature(t *testing.T) {
    // Test implementation
}
```

### Step 4: Update Documentation

Update the README.md with the new command documentation.

## Testing Strategy

### Unit Tests

- Test individual functions in isolation
- Use table-driven tests for multiple scenarios
- Mock external dependencies when needed

Example:

```go
func TestHashString(t *testing.T) {
    tests := []struct {
        name      string
        input     string
        algorithm string
        want      string
        wantErr   bool
    }{
        {"MD5", "test", "md5", "098f6bcd4621d373cade4e832627b4f6", false},
        {"SHA256", "test", "sha256", "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := HashString(tt.input, tt.algorithm)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
            }
            if got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Integration Tests

Integration tests should go in the `tests/` directory and test command execution end-to-end.

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
make coverage

# Run specific package tests
go test -v ./internal/commands/
```

## Development Workflow

### 1. Setup

```bash
# Clone the repository
git clone https://github.com/MayR-Labs/mayrlabs-go.git
cd mayrlabs-go

# Install dependencies
make deps
```

### 2. Make Changes

```bash
# Create a feature branch
git checkout -b feature/my-feature

# Make your changes
# ...

# Format code
make fmt

# Run linter
make lint

# Run tests
make test
```

### 3. Build and Test

```bash
# Build the binary
make build

# Test the binary
./mayrlabs --help
./mayrlabs your-command
```

### 4. Submit Changes

```bash
# Commit your changes
git add .
git commit -m "Add feature: description"

# Push to your fork
git push origin feature/my-feature

# Create a pull request on GitHub
```

## Common Tasks

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Install to GOPATH
make install
```

### Testing

```bash
# Run tests
make test

# Run tests with race detection
go test -race ./...

# Run specific test
go test -v -run TestHashString ./internal/utils/
```

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Run go vet
make vet
```

### Debugging

```bash
# Build with debug symbols
go build -gcflags="all=-N -l" -o mayrlabs main.go

# Use delve for debugging
dlv debug main.go -- your-command args
```

### Adding Dependencies

```bash
# Add a new dependency
go get github.com/some/package

# Update dependencies
go get -u ./...

# Tidy up go.mod
go mod tidy
```

## Code Style Guidelines

### Naming Conventions

- **Commands**: Use verb-noun format (e.g., `CreateKeystore`, `ClearDNS`)
- **Functions**: Use camelCase, exported functions start with uppercase
- **Variables**: Use descriptive names, avoid abbreviations unless common
- **Files**: Use snake_case for filenames (e.g., `editor_config.go`)

### Error Handling

```go
// Good: Return errors
func DoSomething() error {
    if err := someOperation(); err != nil {
        return fmt.Errorf("failed to do something: %w", err)
    }
    return nil
}

// Bad: Panic
func DoSomething() {
    if err := someOperation(); err != nil {
        panic(err)  // Don't do this
    }
}
```

### Comments

```go
// Good: Document exported functions
// HashString generates a hash of the input string using the specified algorithm.
// Supported algorithms are: md5, sha1, sha256.
func HashString(input, algorithm string) (string, error) {
    // ...
}

// Good: Explain complex logic
// Calculate the checksum using a rolling hash algorithm
// to detect duplicate content efficiently
checksum := calculateRollingHash(data)
```

## Performance Considerations

- Avoid unnecessary allocations in hot paths
- Use `strings.Builder` for string concatenation in loops
- Close file handles and resources properly
- Use buffered I/O for large files

## Security Best Practices

- Never log sensitive information (passwords, tokens, keys)
- Validate all user inputs
- Use secure random generation for passwords/tokens
- Be careful with file operations (path traversal attacks)
- Don't execute arbitrary shell commands without validation

## Resources

- [Go Documentation](https://golang.org/doc/)
- [Cobra Documentation](https://cobra.dev/)
- [Go Testing](https://golang.org/pkg/testing/)
- [Effective Go](https://golang.org/doc/effective_go.html)

## Getting Help

- Check existing issues on GitHub
- Ask questions in discussions
- Contact maintainers

---

Happy coding! 🚀
