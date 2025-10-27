# MayR Labs CLI - API Reference

This document provides a comprehensive reference for all commands available in the MayR Labs CLI.

## Table of Contents

1. [General Commands](#general-commands)
2. [System Commands](#system-commands)
3. [Encoding/Decoding Commands](#encodingdecoding-commands)
4. [Git Commands](#git-commands)
5. [Environment Commands](#environment-commands)
6. [Changelog Commands](#changelog-commands)
7. [Flutter Commands](#flutter-commands)
8. [PHP Commands](#php-commands)
9. [JavaScript Commands](#javascript-commands)
10. [AI Commands](#ai-commands)
11. [Session Management](#session-management)
12. [Alias Management](#alias-management)

---

## General Commands

### `mayrlabs uuid`

Generate a UUID v4 (Universally Unique Identifier).

**Usage:**
```bash
mayrlabs uuid [--copy]
```

**Flags:**
- `--copy` - Copy the generated UUID to clipboard

**Examples:**
```bash
# Generate a UUID
mayrlabs uuid
# Output: 550e8400-e29b-41d4-a716-446655440000

# Generate and copy to clipboard
mayrlabs uuid --copy
```

---

### `mayrlabs ulid`

Generate a ULID (Universally Unique Lexicographically Sortable Identifier).

**Usage:**
```bash
mayrlabs ulid [--copy]
```

**Flags:**
- `--copy` - Copy the generated ULID to clipboard

**Examples:**
```bash
# Generate a ULID
mayrlabs ulid
# Output: 01ARZ3NDEKTSV4RRFFQ69G5FAV

# Generate and copy to clipboard
mayrlabs ulid --copy
```

---

### `mayrlabs password`

Generate a cryptographically secure random password.

**Usage:**
```bash
mayrlabs password [length]
```

**Arguments:**
- `length` - Password length (default: 16, interactive if omitted)

**Examples:**
```bash
# Generate 16-character password
mayrlabs password

# Generate 32-character password
mayrlabs password 32

# Interactive mode
mayrlabs password
# Will prompt for length
```

---

### `mayrlabs random-int`

Generate a random integer within a specified range.

**Usage:**
```bash
mayrlabs random-int
```

**Interactive Mode:**
- Prompts for minimum value
- Prompts for maximum value

**Examples:**
```bash
mayrlabs random-int
# Enter min value: 1
# Enter max value: 100
# Output: 42
```

---

### `mayrlabs hash`

Generate a hash of a string using various algorithms.

**Usage:**
```bash
mayrlabs hash [string] [--algorithm ALGORITHM]
```

**Arguments:**
- `string` - String to hash (interactive if omitted)

**Flags:**
- `--algorithm` - Hash algorithm: `md5`, `sha1`, `sha256` (default: sha256)

**Examples:**
```bash
# Generate SHA256 hash
mayrlabs hash "my-secret"

# Generate MD5 hash
mayrlabs hash "my-secret" --algorithm md5

# Generate SHA1 hash
mayrlabs hash "password123" --algorithm sha1
```

---

### `mayrlabs hash-file`

Generate a hash of a file's contents.

**Usage:**
```bash
mayrlabs hash-file [file] [--algorithm ALGORITHM]
```

**Arguments:**
- `file` - Path to file (interactive if omitted)

**Flags:**
- `--algorithm` - Hash algorithm: `md5`, `sha1`, `sha256` (default: sha256)

**Examples:**
```bash
# Hash a file with SHA256
mayrlabs hash-file document.txt

# Hash with MD5
mayrlabs hash-file image.jpg --algorithm md5
```

---

### `mayrlabs create-keystore`

Create a new Android keystore in PKCS12 format.

**Usage:**
```bash
mayrlabs create-keystore
```

**Interactive Prompts:**
- Key alias
- Keystore filename
- Validity period (days)
- Key password
- Distinguished name details (name, organization, city, state, country)

**Examples:**
```bash
mayrlabs create-keystore
# Prompts for all required information
# Creates: myapp.keystore
```

---

### `mayrlabs dns-clear`

Clear the DNS cache on your system.

**Usage:**
```bash
mayrlabs dns-clear
```

**Supported Platforms:**
- macOS: Runs `sudo dscacheutil -flushcache && sudo killall -HUP mDNSResponder`
- Linux: Runs `sudo systemd-resolve --flush-caches`
- Windows: Runs `ipconfig /flushdns`

**Examples:**
```bash
mayrlabs dns-clear
# Select your OS
# Runs the appropriate command
```

---

### `mayrlabs ci`

Generate CI/CD workflow configuration files.

**Usage:**
```bash
mayrlabs ci [--lang LANGUAGE] [--vcs VCS]
```

**Flags:**
- `--lang` - Programming language: `go`, `javascript`, `python`, `php`, `flutter`, `dart`
- `--vcs` - Version control system: `github`, `gitlab`, `circleci`

**Examples:**
```bash
# Generate GitHub Actions for Go
mayrlabs ci --lang go --vcs github

# Generate GitLab CI for Python
mayrlabs ci --lang python --vcs gitlab

# Interactive mode
mayrlabs ci
```

**Output Files:**
- GitHub: `.github/workflows/ci.yml`
- GitLab: `.gitlab-ci.yml`
- CircleCI: `.circleci/config.yml`

---

### `mayrlabs format`

Format code for various programming languages.

**Usage:**
```bash
mayrlabs format [language]
```

**Supported Languages:**
- `go` - Runs `gofmt` and `go fmt`
- `javascript` - Runs Prettier
- `python` - Runs Black and autopep8
- `php` - Runs PHP-CS-Fixer
- `dart` - Runs `dart format`

**Examples:**
```bash
# Format Go code
mayrlabs format go

# Format JavaScript
mayrlabs format javascript

# Interactive mode
mayrlabs format
```

---

### `mayrlabs add-license`

Create a LICENSE file for your project.

**Usage:**
```bash
mayrlabs add-license [--type TYPE] [--author AUTHOR] [--year YEAR] [--author-url URL] [--force]
```

**Flags:**
- `--type` - License type: `mit`, `apache2`, `gpl3`, `bsd3`
- `--author` - Copyright holder name
- `--year` - Copyright year (default: current year)
- `--author-url` - Optional URL for the author
- `--force` - Overwrite existing LICENSE file

**Examples:**
```bash
# Create MIT license
mayrlabs add-license --type mit --author "John Doe" --year 2025

# With author URL
mayrlabs add-license --type mit --author "ACME Corp" --author-url "https://acme.com"

# Interactive mode
mayrlabs add-license
```

---

### `mayrlabs editor-config`

Generate an `.editorconfig` file for consistent code formatting.

**Usage:**
```bash
mayrlabs editor-config [language]
```

**Supported Languages:**
- `go`
- `javascript`
- `python`
- `php`
- `dart`
- `general`

**Examples:**
```bash
# Generate for Go project
mayrlabs editor-config go

# Generate for JavaScript project
mayrlabs editor-config javascript

# Interactive mode
mayrlabs editor-config
```

---

### `mayrlabs roll-dice`

Roll virtual dice and get results.

**Usage:**
```bash
mayrlabs roll-dice [n]
```

**Arguments:**
- `n` - Number of dice to roll (1-100)

**Examples:**
```bash
# Roll 3 dice
mayrlabs roll-dice 3
# Output: Results: [4 1 3], Total: 8, Average: 2.67

# Interactive mode
mayrlabs roll-dice
```

---

### `mayrlabs quote`

Display a random motivational quote for developers.

**Usage:**
```bash
mayrlabs quote
```

**Examples:**
```bash
mayrlabs quote
# Output: "First, solve the problem. Then, write the code." - John Johnson
```

---

### `mayrlabs version`

Display the current version of MayR Labs CLI.

**Usage:**
```bash
mayrlabs version
```

**Examples:**
```bash
mayrlabs version
# Output: mayrlabs version 1.0.0
```

---

### `mayrlabs visit`

Open the MayR Labs website in your default browser.

**Usage:**
```bash
mayrlabs visit
```

**URL:** https://mayrlabs.com

---

### `mayrlabs github`

Open the mayrlabs-go GitHub repository in your default browser.

**Usage:**
```bash
mayrlabs github
```

**URL:** https://github.com/MayR-Labs/mayrlabs-go

---

## System Commands

### `mayrlabs upgrade`

Upgrade MayR Labs CLI to the latest version.

**Usage:**
```bash
mayrlabs upgrade
```

**Behavior:**
- Downloads the latest release from GitHub
- Detects your OS and architecture automatically
- Replaces the current binary
- Requires appropriate permissions

**Examples:**
```bash
mayrlabs upgrade
```

---

## Encoding/Decoding Commands

### `mayrlabs base64`

Encode or decode strings using Base64.

**Usage:**
```bash
mayrlabs base64 <encode|decode> [string] [--copy]
```

**Arguments:**
- `encode|decode` - Operation to perform
- `string` - String to encode/decode

**Flags:**
- `--copy` - Copy result to clipboard

**Examples:**
```bash
# Encode a string
mayrlabs base64 encode "Hello World"
# Output: SGVsbG8gV29ybGQ=

# Decode a string
mayrlabs base64 decode "SGVsbG8gV29ybGQ="
# Output: Hello World

# Encode and copy
mayrlabs base64 encode "secret" --copy
```

---

### `mayrlabs base64-file`

Encode a file to Base64.

**Usage:**
```bash
mayrlabs base64-file [path] [--copy]
```

**Arguments:**
- `path` - Path to file

**Flags:**
- `--copy` - Copy result to clipboard

**Examples:**
```bash
# Encode a file
mayrlabs base64-file document.txt

# Encode and copy
mayrlabs base64-file image.png --copy
```

---

### `mayrlabs base64-decode-to-file`

Decode a Base64 string and write to a file.

**Usage:**
```bash
mayrlabs base64-decode-to-file [string] [output-file]
```

**Arguments:**
- `string` - Base64 encoded string
- `output-file` - Path where decoded content should be written

**Examples:**
```bash
mayrlabs base64-decode-to-file "SGVsbG8gV29ybGQ=" output.txt
```

---

## Git Commands

### `mayrlabs git`

Show available Git-related commands.

**Usage:**
```bash
mayrlabs git
```

---

### `mayrlabs git prune-stale`

Delete all local branches that no longer exist on the remote.

**Usage:**
```bash
mayrlabs git prune-stale
```

**Behavior:**
- Fetches latest remote information
- Identifies branches deleted on remote
- Prompts for confirmation
- Deletes local copies of removed branches

**Examples:**
```bash
mayrlabs git prune-stale
```

---

## Environment Commands

### `mayrlabs env`

List available environment file management commands.

**Usage:**
```bash
mayrlabs env
```

---

### `mayrlabs env update-example`

Synchronize `.env.example` with `.env` or another source file.

**Usage:**
```bash
mayrlabs env update-example [source]
```

**Arguments:**
- `source` - Source file (default: `.env`)

**Behavior:**
- Reads all keys from source file
- Writes keys with empty values to `.env.example`
- Preserves formatting and comments
- Creates `.env.example` if it doesn't exist

**Examples:**
```bash
# Update from .env
mayrlabs env update-example

# Update from .env.staging
mayrlabs env update-example .env.staging
```

---

### `mayrlabs env validate`

Validate environment files for errors and inconsistencies.

**Usage:**
```bash
mayrlabs env validate
```

**Checks:**
- Missing values
- Duplicate keys
- Invalid syntax
- Keys in `.env.example` but missing in `.env`
- Keys in `.env` but missing in `.env.example`

**Examples:**
```bash
mayrlabs env validate
```

---

### `mayrlabs env arrange`

Organize environment variables by prefix.

**Usage:**
```bash
mayrlabs env arrange [file]
```

**Arguments:**
- `file` - Environment file to arrange (default: `.env`)

**Behavior:**
- Groups variables by prefix (e.g., `APP_*`, `DB_*`, `AWS_*`)
- Sorts keys alphabetically within groups
- Preserves comments
- Creates backup of original file

**Examples:**
```bash
# Arrange .env
mayrlabs env arrange

# Arrange .env.staging
mayrlabs env arrange .env.staging
```

---

## Changelog Commands

### `mayrlabs changelog`

Display all changelog management commands.

**Usage:**
```bash
mayrlabs changelog
```

---

### `mayrlabs changelog create`

Create a new `CHANGELOG.md` file.

**Usage:**
```bash
mayrlabs changelog create [--force]
```

**Flags:**
- `--force` - Overwrite existing CHANGELOG.md

**Examples:**
```bash
# Create new changelog
mayrlabs changelog create

# Force overwrite
mayrlabs changelog create --force
```

---

### `mayrlabs changelog record`

Add a new version entry to the changelog.

**Usage:**
```bash
mayrlabs changelog record [version] [summary] [--wip]
```

**Arguments:**
- `version` - Version number (e.g., `v1.2.0`)
- `summary` - Brief description of changes

**Flags:**
- `--wip` - Mark as Work In Progress

**Examples:**
```bash
# Add a release entry
mayrlabs changelog record v1.2.0 "Added user authentication"

# Add a WIP entry
mayrlabs changelog record v1.3.0 "Working on API improvements" --wip
```

---

## Flutter Commands

### `mayrlabs flutter`

List Flutter-related commands.

**Usage:**
```bash
mayrlabs flutter
```

---

### `mayrlabs flutter create-scripts`

Generate build scripts for Flutter projects.

**Usage:**
```bash
mayrlabs flutter create-scripts
```

**Generated Scripts:**
- `scripts/build-apk.sh` - Build debug APK
- `scripts/build-apk-release.sh` - Build release APK
- `scripts/build-appbundle.sh` - Build App Bundle
- `scripts/build-ios.sh` - Build iOS app
- `scripts/build-ipa.sh` - Build IPA
- `scripts/clean.sh` - Clean build artifacts
- `scripts/run-tests.sh` - Run Flutter tests

**Examples:**
```bash
mayrlabs flutter create-scripts

# Then use:
./scripts/build-apk-release.sh
```

---

## PHP Commands

### `mayrlabs php`

List PHP-related commands.

**Usage:**
```bash
mayrlabs php
```

---

### `mayrlabs php cs-fix`

Run PHP CodeSniffer or PHP-CS-Fixer to fix code style issues.

**Usage:**
```bash
mayrlabs php cs-fix
```

**Examples:**
```bash
mayrlabs php cs-fix
```

---

### `mayrlabs php lint`

Lint PHP files for syntax errors.

**Usage:**
```bash
mayrlabs php lint
```

**Examples:**
```bash
mayrlabs php lint
```

---

## JavaScript Commands

### `mayrlabs js`

List JavaScript-related commands.

**Usage:**
```bash
mayrlabs js
```

---

### `mayrlabs js setup-prettier`

Install and configure Prettier for JavaScript/TypeScript projects.

**Usage:**
```bash
mayrlabs js setup-prettier
```

**Actions:**
- Installs Prettier via npm
- Creates `.prettierrc.yaml` configuration
- Creates `.prettierignore` file

**Examples:**
```bash
mayrlabs js setup-prettier
```

---

### `mayrlabs js pretty`

Format JavaScript/TypeScript code with Prettier.

**Usage:**
```bash
mayrlabs js pretty
```

**Examples:**
```bash
mayrlabs js pretty
```

---

## AI Commands

### `mayrlabs ai-setup`

Configure Gemini API key for AI features.

**Usage:**
```bash
mayrlabs ai-setup
```

**Behavior:**
- Prompts for Gemini API key (get from https://aistudio.google.com/app/apikey)
- Validates the key by making a test request
- Stores the key securely in `~/.mayrlabs/gemini-api-key`

**Examples:**
```bash
mayrlabs ai-setup
# Enter your Gemini API key: [paste key]
# ✅ API key validated and stored successfully
```

---

### `mayrlabs ai`

Query the AI using Google's Gemini model.

**Usage:**
```bash
mayrlabs ai [query...]
```

**Arguments:**
- `query` - Question or prompt (supports multi-word queries)

**Behavior:**
- If query provided: Sends to AI immediately
- If no query: Opens multiline input form
- Uses `gemini-2.0-flash-exp` model
- Maintains conversation context within session

**Examples:**
```bash
# Single query
mayrlabs ai "How do I implement a binary search tree in Go?"

# Multi-line input
mayrlabs ai
# Opens interactive prompt for longer queries
```

---

### `mayrlabs ai-file`

Send a file's content to the AI for analysis or review.

**Usage:**
```bash
mayrlabs ai-file [path]
```

**Arguments:**
- `path` - Path to text-based file

**Behavior:**
- Reads file content
- Sends to AI with context
- Returns AI analysis/suggestions

**Examples:**
```bash
# Review a source file
mayrlabs ai-file src/main.go

# Analyze a configuration file
mayrlabs ai-file nginx.conf
```

---

### `mayrlabs ai-alias`

Create a permanent shell alias for the `mayrlabs ai` command.

**Usage:**
```bash
mayrlabs ai-alias [name]
```

**Arguments:**
- `name` - Alias name

**Examples:**
```bash
# Create an 'ask' alias
mayrlabs ai-alias ask

# Now use:
ask "What is dependency injection?"
```

---

### `mayrlabs ai-clear`

Remove the stored Gemini API key.

**Usage:**
```bash
mayrlabs ai-clear
```

**Examples:**
```bash
mayrlabs ai-clear
# API key cleared successfully
```

---

## Session Management

### `mayrlabs session-start`

Start an interactive development session.

**Usage:**
```bash
mayrlabs session-start [summary]
```

**Arguments:**
- `summary` - Brief description of session (interactive if omitted)

**Features:**
- Records all AI interactions
- Saves notes and commands
- Tracks session duration
- Stores in `~/.mayrlabs/sessions/`

**Examples:**
```bash
# Start with summary
mayrlabs session-start "Implementing user authentication"

# Interactive mode
mayrlabs session-start
```

---

### `mayrlabs sessions`

List and manage all sessions.

**Usage:**
```bash
mayrlabs sessions
```

**Behavior:**
- Lists all session files
- Shows summary and timestamp
- Allows opening/viewing sessions

**Examples:**
```bash
mayrlabs sessions
```

---

### `mayrlabs session-clear`

Clear all sessions with PIN confirmation.

**Usage:**
```bash
mayrlabs session-clear
```

**Behavior:**
- Requires PIN confirmation
- Deletes all session files
- Cannot be undone

**Examples:**
```bash
mayrlabs session-clear
# Enter PIN to confirm: [enter PIN]
```

---

### `mayrlabs session-prune`

Delete sessions older than specified days.

**Usage:**
```bash
mayrlabs session-prune [days]
```

**Arguments:**
- `days` - Age threshold in days

**Examples:**
```bash
# Delete sessions older than 30 days
mayrlabs session-prune 30

# Delete sessions older than 7 days
mayrlabs session-prune 7
```

---

### `mayrlabs secure-session-start`

Start an encrypted interactive development session.

**Usage:**
```bash
mayrlabs secure-session-start [summary]
```

**Arguments:**
- `summary` - Brief description of session

**Features:**
- End-to-end encryption
- Password-protected
- Secure storage
- Same features as regular sessions

**Examples:**
```bash
# Start encrypted session
mayrlabs secure-session-start "Security vulnerability fix"
# Enter encryption password: [password]
```

---

### `mayrlabs secure-sessions`

List and manage encrypted secure sessions.

**Usage:**
```bash
mayrlabs secure-sessions
```

**Behavior:**
- Lists all encrypted sessions
- Requires password to open
- Shows summary and timestamp

**Examples:**
```bash
mayrlabs secure-sessions
```

---

## Alias Management

### `mayrlabs alias`

Create a permanent shell alias for the mayrlabs command.

**Usage:**
```bash
mayrlabs alias [name]
```

**Arguments:**
- `name` - Alias name

**Behavior:**
- Creates alias in shell configuration file
- Supports bash, zsh, fish
- Updates current session
- Persists across restarts

**Examples:**
```bash
# Create 'ml' alias
mayrlabs alias ml

# Now use:
ml version
ml uuid
```

---

### `mayrlabs alias-list`

List all mayrlabs aliases.

**Usage:**
```bash
mayrlabs alias-list
```

**Examples:**
```bash
mayrlabs alias-list
# Output:
# ml -> mayrlabs
# ask -> mayrlabs ai
```

---

### `mayrlabs alias-clear`

Clear all mayrlabs aliases with confirmation.

**Usage:**
```bash
mayrlabs alias-clear
```

**Behavior:**
- Prompts for confirmation
- Removes all aliases from shell configuration
- Reloads shell configuration

**Examples:**
```bash
mayrlabs alias-clear
# Are you sure? [y/N]: y
# All aliases cleared
```

---

## Exit Codes

The CLI uses standard exit codes:

- `0` - Success
- `1` - General error
- `2` - Misuse of shell command (invalid arguments)

---

## Environment Variables

### `MAYRLABS_HOME`

Override the default configuration directory (default: `~/.mayrlabs`).

```bash
export MAYRLABS_HOME=/custom/path
```

---

## Configuration Files

### Location

All configuration files are stored in `~/.mayrlabs/` (or `$MAYRLABS_HOME`).

### Files

- `gemini-api-key` - Stored Gemini API key (plain text)
- `sessions/` - Regular session files (Markdown format)
- `secure-sessions/` - Encrypted session files
- `aliases.conf` - Shell alias configurations

---

## Tips & Best Practices

1. **Use `--copy` flags** when you need to paste output elsewhere
2. **Interactive mode** is great for commands with many options
3. **Create aliases** for frequently used commands
4. **Use sessions** to track complex development work
5. **Secure sessions** for sensitive or confidential work
6. **Regular pruning** keeps your session storage clean
7. **AI features** work best with clear, specific questions

---

## Getting Help

For any command, use the `--help` flag:

```bash
mayrlabs --help
mayrlabs [command] --help
```

**Report Issues:** https://github.com/MayR-Labs/mayrlabs-go/issues

**Documentation:** https://github.com/MayR-Labs/mayrlabs-go

---

**Last Updated:** 2025-10-27 (v1.0.0)
