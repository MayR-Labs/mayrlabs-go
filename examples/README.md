# Examples

This directory contains example configurations and usage scenarios for MayR Labs CLI.

## Command Examples

### General Commands

#### Generate UUID
```bash
mayrlabs uuid
# Output: 550e8400-e29b-41d4-a716-446655440000
```

#### Generate Password
```bash
# Default 16 characters
mayrlabs password

# Custom length
mayrlabs password 32
```

#### Hash String
```bash
# Default SHA256
mayrlabs hash "my-secret-string"

# Specify algorithm
mayrlabs hash "my-string" --algorithm md5
mayrlabs hash "my-string" --algorithm sha1
mayrlabs hash "my-string" --algorithm sha256
```

### License Management

#### Create MIT License
```bash
mayrlabs add-license --type mit --author "Your Name" --year 2025
```

#### Available License Types
- `mit` - MIT License
- `apache2` - Apache License 2.0
- `gpl3` - GNU General Public License v3
- `bsd3` - BSD 3-Clause License

### Editor Configuration

#### Generate .editorconfig
```bash
# For Go projects
mayrlabs editor-config go

# For JavaScript/TypeScript projects
mayrlabs editor-config javascript

# For Python projects
mayrlabs editor-config python
```

### CI/CD Configuration

#### Generate GitHub Actions Workflow
```bash
mayrlabs ci --lang go --vcs github
mayrlabs ci --lang javascript --vcs github
mayrlabs ci --lang flutter --vcs github
```

#### Generate GitLab CI Configuration
```bash
mayrlabs ci --lang go --vcs gitlab
```

### Code Formatting

#### Format Code
```bash
# Interactive mode (prompts for language)
mayrlabs format

# Direct language specification
mayrlabs format go
mayrlabs format javascript
mayrlabs format python
```

### Environment File Management

#### Update .env.example from .env
```bash
mayrlabs env update-example
```

#### Update from specific source
```bash
mayrlabs env update-example .env.staging
```

#### Validate .env file
```bash
mayrlabs env validate
```

#### Arrange .env file by prefix
```bash
mayrlabs env arrange
mayrlabs env arrange .env.staging
```

### Changelog Management

#### Create CHANGELOG.md
```bash
mayrlabs changelog create
```

#### Record New Version
```bash
mayrlabs changelog record v1.0.0 "Initial release"
```

#### Mark as Work in Progress
```bash
mayrlabs changelog record v1.1.0 "New features" --wip
```

### Git Operations

#### Prune Stale Branches
```bash
mayrlabs git prune-stale
```

### Flutter Commands

#### Create Build Scripts
```bash
mayrlabs flutter create-scripts
```

This creates several scripts in the `scripts/` directory:
- `build-apk.sh` - Build debug APK
- `build-apk-release.sh` - Build release APK
- `build-appbundle.sh` - Build App Bundle
- `build-ios.sh` - Build iOS app
- `build-ipa.sh` - Build IPA
- `clean.sh` - Clean Flutter build
- `run-tests.sh` - Run Flutter tests

### PHP Commands

#### Fix Code Style
```bash
mayrlabs php cs-fix
```

#### Lint PHP Files
```bash
mayrlabs php lint
```

### JavaScript Commands

#### Setup Prettier
```bash
mayrlabs js setup-prettier
```

#### Format with Prettier
```bash
mayrlabs js pretty
```

### Bonus Commands

#### Get Motivational Quote
```bash
mayrlabs quote
```

### DNS Management

#### Clear DNS Cache
```bash
mayrlabs dns-clear
```

Automatically detects your OS and runs the appropriate command:
- macOS: `sudo dscacheutil -flushcache`
- Linux: `sudo systemd-resolve --flush-caches`
- Windows: `ipconfig /flushdns`

### Keystore Creation

#### Create Android Keystore
```bash
mayrlabs create-keystore
```

Interactive prompts for:
- Key alias
- Keystore filename
- Validity period

## Workflow Examples

### Setting Up a New Go Project

```bash
# Generate .editorconfig
mayrlabs editor-config go

# Create LICENSE
mayrlabs add-license --type mit --author "Your Name" --year 2025

# Generate CI/CD configuration
mayrlabs ci --lang go --vcs github

# Create CHANGELOG
mayrlabs changelog create

# Format code
mayrlabs format go
```

### Setting Up a New JavaScript Project

```bash
# Generate .editorconfig
mayrlabs editor-config javascript

# Setup Prettier
mayrlabs js setup-prettier

# Create LICENSE
mayrlabs add-license --type mit --author "Your Name" --year 2025

# Generate CI/CD configuration
mayrlabs ci --lang javascript --vcs github

# Format code
mayrlabs js pretty
```

### Managing Environment Files

```bash
# After adding new variables to .env
mayrlabs env update-example

# Validate the .env file
mayrlabs env validate

# Organize variables by prefix
mayrlabs env arrange
```

### Flutter Project Setup

```bash
# Generate build scripts
mayrlabs flutter create-scripts

# Generate CI/CD
mayrlabs ci --lang flutter --vcs github

# Create .editorconfig
mayrlabs editor-config dart

# Create LICENSE
mayrlabs add-license --type mit --author "Your Name" --year 2025
```

## Integration with Other Tools

### Git Hooks

You can use MayR Labs CLI in Git hooks. For example, create `.git/hooks/pre-commit`:

```bash
#!/bin/bash
# Format code before commit
mayrlabs format go
```

### CI/CD Integration

The generated CI/CD files are ready to use and include:
- Dependency installation
- Testing
- Building
- Linting

### Makefile Integration

You can add MayR Labs CLI commands to your Makefile:

```makefile
.PHONY: setup

setup:
	mayrlabs editor-config go
	mayrlabs add-license --type mit --author "Team" --year 2025
	mayrlabs ci --lang go --vcs github
```

## Tips and Tricks

1. **Interactive Mode**: Most commands support interactive mode if you omit required parameters
2. **Force Overwrite**: Many file-generating commands support `--force` flag
3. **Help**: Use `--help` with any command for detailed information
4. **Chaining**: Combine with shell operators for powerful workflows

```bash
# Generate multiple files at once
mayrlabs editor-config go && \
mayrlabs add-license --type mit --author "Me" --year 2025 && \
mayrlabs changelog create
```

## Common Use Cases

### 1. Quick Project Initialization
```bash
mayrlabs add-license --type mit --author "Your Name" --year 2025
mayrlabs editor-config go
mayrlabs ci --lang go --vcs github
mayrlabs changelog create
```

### 2. Environment Management
```bash
# Daily workflow
mayrlabs env validate
mayrlabs env arrange
mayrlabs env update-example
```

### 3. Code Maintenance
```bash
mayrlabs format go
mayrlabs git prune-stale
```

### 4. Security
```bash
# Generate secure password
mayrlabs password 32

# Generate hash for verification
mayrlabs hash "important-data" --algorithm sha256
```

## Need Help?

Run any command with `--help` for detailed usage:

```bash
mayrlabs --help
mayrlabs ci --help
mayrlabs env --help
```
