# Contributing to MayR Labs CLI

Thank you for considering contributing to MayR Labs CLI! This document outlines the process and guidelines for contributing.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Coding Guidelines](#coding-guidelines)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)

## Code of Conduct

This project adheres to a code of conduct. By participating, you are expected to uphold this code. Please be respectful and constructive in all interactions.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/mayrlabs-go.git`
3. Add upstream remote: `git remote add upstream https://github.com/MayR-Labs/mayrlabs-go.git`
4. Create a new branch: `git checkout -b feature/your-feature-name`

## Development Setup

### Prerequisites

- Go 1.17 or higher
- Git

### Install Dependencies

```bash
cd mayrlabs-go
go mod download
```

### Build the Project

```bash
go build -o mayrlabs main.go
```

### Run Tests

```bash
go test ./... -v
```

## How to Contribute

### Reporting Bugs

- Check if the bug has already been reported in Issues
- If not, create a new issue with:
  - Clear title and description
  - Steps to reproduce
  - Expected vs actual behavior
  - Your environment (OS, Go version, etc.)

### Suggesting Features

- Open an issue with the "enhancement" label
- Clearly describe the feature and its use case
- Explain why this feature would be useful

### Code Contributions

1. **Find or create an issue** for what you want to work on
2. **Comment on the issue** to let others know you're working on it
3. **Write your code** following our coding guidelines
4. **Add tests** for any new functionality
5. **Update documentation** if needed
6. **Submit a pull request**

## Coding Guidelines

### Go Style

- Follow standard Go conventions and idioms
- Run `gofmt` on your code before committing
- Use meaningful variable and function names
- Add comments for exported functions and complex logic
- Keep functions small and focused

### Code Structure

```
.
├── cmd/                    # CLI commands (Cobra setup)
├── internal/
│   ├── commands/          # Command implementations
│   ├── utils/             # Utility functions
│   └── config/            # Configuration handling
├── tests/                 # Integration tests
└── main.go                # Application entry point
```

### Adding a New Command

1. Create a new file in `internal/commands/` (e.g., `mycommand.go`)
2. Implement your command using Cobra:

```go
package commands

import (
    "github.com/spf13/cobra"
)

var MyCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Brief description",
    Long:  "Detailed description",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}
```

3. Add the command to `cmd/root.go`:

```go
func init() {
    rootCmd.AddCommand(commands.MyCmd)
}
```

4. Write tests in `internal/commands/mycommand_test.go`

### Error Handling

- Always handle errors explicitly
- Return errors rather than panicking
- Use `fmt.Errorf` with `%w` for error wrapping

### Testing

- Write unit tests for all new functionality
- Aim for >80% code coverage
- Use table-driven tests where appropriate
- Test edge cases and error scenarios

Example:

```go
func TestMyFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        // Test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := MyFunction(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("MyFunction() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("MyFunction() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Submitting Changes

### Before Submitting

- [ ] Code follows project style guidelines
- [ ] All tests pass: `go test ./...`
- [ ] Code passes linting: `go vet ./...`
- [ ] Added/updated tests for changes
- [ ] Updated documentation if needed
- [ ] Commit messages are clear and descriptive

### Pull Request Process

1. Update your branch with the latest upstream:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. Push your changes:
   ```bash
   git push origin feature/your-feature-name
   ```

3. Create a Pull Request:
   - Use a clear, descriptive title
   - Reference any related issues
   - Describe what changes you made and why
   - Include screenshots for UI changes (if applicable)

4. Address review feedback:
   - Make requested changes
   - Push updates to the same branch
   - Respond to comments

### Commit Message Guidelines

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters
- Reference issues and pull requests when relevant

Example:
```
Add support for generating TypeScript config

- Implement TypeScript configuration generation
- Add tests for TypeScript config
- Update documentation

Fixes #123
```

## Questions?

If you have questions, feel free to:
- Open an issue with the "question" label
- Reach out to the maintainers

Thank you for contributing! 🎉
