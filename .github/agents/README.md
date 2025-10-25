# MayR Labs CLI - Copilot Agent Instructions

## Project Overview
MayR Labs CLI is a lightweight, cross-platform command-line tool built with Go to streamline common development, configuration, and automation tasks across projects.

## Tech Stack
- **Language:** Go 1.17+
- **CLI Framework:** Cobra
- **Build:** Single static binary
- **Platforms:** macOS, Linux, Windows

## Project Structure
```
.
├── cmd/                    # CLI commands (Cobra)
│   └── root.go            # Root command and app initialization
├── internal/              # Private application code
│   ├── commands/         # Command implementations
│   ├── utils/            # Utility functions
│   └── config/           # Configuration handling
├── pkg/                   # Public libraries (if any)
├── tests/                 # Integration tests
├── .github/
│   └── workflows/        # CI/CD workflows
└── main.go               # Application entry point
```

## Development Guidelines

### Building and Testing
```bash
# Build the project
go build -o mayrlabs main.go

# Run tests
go test ./... -v

# Run tests with coverage
go test ./... -cover -coverprofile=coverage.out

# Run linter
go vet ./...
golangci-lint run
```

### Adding New Commands
1. Create command file in `internal/commands/`
2. Implement command logic with proper error handling
3. Add command to root command in `cmd/root.go`
4. Write unit tests in corresponding `_test.go` file
5. Update documentation if needed

### Code Style
- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Add comments for exported functions and complex logic
- Keep functions small and focused
- Handle errors explicitly, don't ignore them

### Testing Requirements
- Write unit tests for all public functions
- Aim for >80% code coverage
- Include edge cases and error scenarios
- Use table-driven tests where appropriate

### CI/CD
- All PRs must pass CI checks (build, test, lint)
- Code must be formatted with `gofmt`
- All tests must pass
- No security vulnerabilities

## Common Tasks

### Working with Cobra Commands
Commands are organized hierarchically. Use `cobra-cli` or create manually:
```go
var myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Brief description",
    Long:  "Detailed description",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}
```

### Cross-Platform Considerations
- Use `filepath.Join()` for paths
- Use `os.PathSeparator` when needed
- Test on multiple platforms if modifying platform-specific code
- Use build tags for platform-specific code: `// +build windows`

## Dependencies
- `github.com/spf13/cobra` - CLI framework
- `github.com/spf13/viper` - Configuration management
- Additional dependencies as needed for specific features

## Security Considerations
- Never commit sensitive data (API keys, passwords, etc.)
- Validate all user inputs
- Use secure random generation for passwords/tokens
- Be careful with file operations (path traversal attacks)

## Documentation
- Keep README.md updated with new commands
- Add inline comments for complex logic
- Document all exported functions and types
- Include usage examples in help text

## Getting Help
- Check existing issues on GitHub
- Refer to Cobra documentation: https://cobra.dev/
- Follow Go best practices: https://golang.org/doc/effective_go.html
