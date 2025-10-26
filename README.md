# 🧰 MayR Labs CLI

**MayR Labs CLI** is a lightweight, cross-platform command-line tool built with Go to streamline common development, configuration, and automation tasks across projects.

It aims to provide developers with a unified interface for generating configs, formatting code, managing environments, handling CI/CD, and keeping project structure consistent — all in seconds.

---

## 🚀 Features

- Consistent project setup and formatting across multiple languages.
- Fast cross-platform DNS cache clearing.
- Auto-generation of CI/CD configuration files for major languages and VCS providers.
- Easy license creation and editor configuration.
- Smart environment file management and changelog handling.
- Integrated support for Flutter, PHP, JavaScript, and Go workflows.

---

## 📦 Installation

### Using Installation Script (Recommended)

The easiest way to install mayrlabs is using our installation script:

```bash
# Using curl
curl -sSL https://raw.githubusercontent.com/MayR-Labs/mayrlabs-go/main/install.sh | bash

# Or using wget
wget -qO- https://raw.githubusercontent.com/MayR-Labs/mayrlabs-go/main/install.sh | bash
```

The script automatically:
- Detects your operating system and architecture
- Downloads the latest release
- Installs the binary to an appropriate location
- Verifies the installation

### Using Go Install

If you have Go 1.17+ installed:

```bash
go install github.com/MayR-Labs/mayrlabs-go@latest
```

**Note:** This will install the binary as `mayrlabs-go` (not `mayrlabs`) in your `$GOPATH/bin` directory. You can create an alias or symlink:

```bash
# Create an alias (add to your ~/.bashrc or ~/.zshrc)
alias mayrlabs='mayrlabs-go'

# Or create a symlink
ln -s $(which mayrlabs-go) $(dirname $(which mayrlabs-go))/mayrlabs
```

**Ensure `$GOPATH/bin` is in your PATH:**

Add this to your shell configuration file (`~/.bashrc`, `~/.zshrc`, or `~/.profile`):

```bash
export PATH="$GOPATH/bin:$PATH"
# or if GOPATH is not set
export PATH="$HOME/go/bin:$PATH"
```

Then reload your shell configuration:

```bash
source ~/.bashrc  # or ~/.zshrc, ~/.profile depending on your shell
```

### Download Pre-built Binaries

Download the latest release for your platform from the [Releases page](https://github.com/MayR-Labs/mayrlabs-go/releases):

**Linux (AMD64):**
```bash
curl -L https://github.com/MayR-Labs/mayrlabs-go/releases/latest/download/mayrlabs-linux-amd64 -o mayrlabs
chmod +x mayrlabs
sudo mv mayrlabs /usr/local/bin/
```

**macOS (Apple Silicon):**
```bash
curl -L https://github.com/MayR-Labs/mayrlabs-go/releases/latest/download/mayrlabs-darwin-arm64 -o mayrlabs
chmod +x mayrlabs
sudo mv mayrlabs /usr/local/bin/
```

**Windows:**
Download `mayrlabs-windows-amd64.exe` from the releases page and add it to your PATH.

### Build from Source

Clone and build manually:

```bash
git clone https://github.com/MayR-Labs/mayrlabs-go.git
cd mayrlabs-go
make build
# Or: go build -o mayrlabs main.go

# Optionally install to system
sudo mv mayrlabs /usr/local/bin/
```

### Verify Installation

```bash
mayrlabs --version
mayrlabs --help
```

---

## 🧭 Usage

### Basic Usage

```bash
mayrlabs [command] [flags]
```

### Get Help

```bash
# General help
mayrlabs --help

# Help for specific command
mayrlabs [command] --help

# Example
mayrlabs ci --help
```

### Interactive Mode

Most commands support interactive mode if you omit required parameters:

```bash
# Will prompt you for language and VCS
mayrlabs ci

# Will prompt for language
mayrlabs format
```

---

## 📚 Quick Start Examples

### Setup a New Project

```bash
# Generate editor configuration
mayrlabs editor-config go

# Add MIT license
mayrlabs add-license --type mit --author "Your Name" --year 2025

# Create CI/CD workflow
mayrlabs ci --lang go --vcs github

# Create changelog
mayrlabs changelog create
```

### Generate Utilities

```bash
# Generate a UUID
mayrlabs uuid
# Output: 550e8400-e29b-41d4-a716-446655440000

# Generate a secure password (32 characters)
mayrlabs password 32
# Output: Xy9$mK2#pL8vN@qR5wT7uZ1aB3cD6eF0

# Hash a string
mayrlabs hash "my-secret" --algorithm sha256
# Output: my-secret (sha256): 8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92
```

### Environment Management

```bash
# Validate your .env file
mayrlabs env validate

# Update .env.example from .env
mayrlabs env update-example

# Organize .env by prefix (APP_*, DB_*, etc.)
mayrlabs env arrange
```

### Code Formatting

```bash
# Format Go code
mayrlabs format go

# Format JavaScript code  
mayrlabs format javascript

# Format Python code
mayrlabs format python
```

### Git Operations

```bash
# Remove stale local branches
mayrlabs git prune-stale
```

---

## ⚙️ Commands

### 🧱 General Commands

| Command                             | Description                                                            |
| ----------------------------------- | ---------------------------------------------------------------------- |
| `mayrlabs uuid`                     | Generate UUID v4                                                       |
| `mayrlabs password [length]`        | Generate a random password of specified length (default: 16)           |
| `mayrlabs hash [string]`            | Generate hash of a string using md5, sha1, or sha256                   |
| `mayrlabs create-keystore`          | Create a new keystore interactively                                    |
| `mayrlabs dns-clear`                | Clear the DNS cache (choose macOS, Linux, or Windows)                  |
| `mayrlabs ci`                       | Generate CI/CD workflow YAML for your language and VCS                 |
| `mayrlabs format [language]`        | Format project files for a given language (interactive if omitted)     |
| `mayrlabs add-license`              | Create a LICENSE file based on selected type, year, and author         |
| `mayrlabs editor-config [language]` | Generate `.editorconfig` for a specific language (supports `--force`)  |
| `mayrlabs quote`                    | Display a random motivational quote for developers                     |
| `mayrlabs version`                  | Display the version of mayrlabs CLI                                    |
| `mayrlabs visit`                    | Open the MayR Labs website in your browser                             |
| `mayrlabs github`                   | Open the mayrlabs-go GitHub repository in your browser                 |

---

### 🌿 Git Commands

| Command                    | Description                                        |
| -------------------------- | -------------------------------------------------- |
| `mayrlabs git`             | Show available Git-related commands                |
| `mayrlabs git prune-stale` | Delete all local branches not found on the remote  |

---

### 🔐 ENV Commands

| Command                                | Description                                                                           |
| -------------------------------------- | ------------------------------------------------------------------------------------- |
| `mayrlabs env`                         | List available environment commands                                                   |
| `mayrlabs env update-example [source]` | Sync `.env.example` with `.env` or `.env.staging`. Creates `.env.example` if missing |
| `mayrlabs env validate`                | Check for missing keys, invalid values, or duplicated variables                       |
| `mayrlabs env arrange [file]`          | Sort and group environment keys by prefix (e.g., `APP_*`, `DB_*`)                     |

---

### 📝 CHANGELOG Commands

| Command                                         | Description                                           |
| ----------------------------------------------- | ----------------------------------------------------- |
| `mayrlabs changelog`                            | Display all changelog commands                        |
| `mayrlabs changelog create [--force]`           | Create or overwrite `CHANGELOG.md`                    |
| `mayrlabs changelog record [version] [summary]` | Add a new entry to `CHANGELOG.md` (supports `--wip`)  |

---

### 🐦 Flutter Commands

| Command                           | Description                                                         |
| --------------------------------- | ------------------------------------------------------------------- |
| `mayrlabs flutter`                | List Flutter-related commands                                       |
| `mayrlabs flutter create-scripts` | Add useful build scripts to `scripts/` (IPA, APK, AppBundle, etc.) |

---

### 🐘 PHP Commands

| Command               | Description                   |
| --------------------- | ----------------------------- |
| `mayrlabs php`        | List PHP commands             |
| `mayrlabs php cs-fix` | Run PHP CodeSniffer/CS-Fixer  |
| `mayrlabs php lint`   | Lint PHP files                |

---

### ⚡ JavaScript Commands

| Command                      | Description                                             |
| ---------------------------- | ------------------------------------------------------- |
| `mayrlabs js`                | List JavaScript commands                                |
| `mayrlabs js setup-prettier` | Install and configure Prettier with `.prettierrc.yaml`  |
| `mayrlabs js pretty`         | Run Prettier on the project                             |

---

## 🧠 Practical Examples

### Example 1: Setup a Go Project

```bash
# Initialize project files
mayrlabs editor-config go
mayrlabs add-license --type mit --author "Your Name" --year 2025
mayrlabs ci --lang go --vcs github
mayrlabs changelog create

# Format code
mayrlabs format go
```

### Example 2: Setup a JavaScript Project

```bash
# Configure project
mayrlabs editor-config javascript
mayrlabs js setup-prettier
mayrlabs add-license --type mit --author "Your Team" --year 2025

# Generate CI
mayrlabs ci --lang javascript --vcs github

# Format code
mayrlabs js pretty
```

### Example 3: Manage Environment Files

```bash
# After adding variables to .env
mayrlabs env validate          # Check for errors
mayrlabs env arrange           # Organize by prefix
mayrlabs env update-example    # Sync .env.example
```

### Example 4: Flutter Project Setup

```bash
# Generate build scripts
mayrlabs flutter create-scripts

# This creates:
# - scripts/build-apk.sh
# - scripts/build-apk-release.sh
# - scripts/build-appbundle.sh
# - scripts/build-ios.sh
# - scripts/build-ipa.sh
# - scripts/clean.sh
# - scripts/run-tests.sh

# Run a script
./scripts/build-apk-release.sh
```

### Example 5: Generate CI for Different Platforms

```bash
# GitHub Actions for Go
mayrlabs ci --lang go --vcs github

# GitLab CI for Python
mayrlabs ci --lang python --vcs gitlab

# GitHub Actions for Flutter
mayrlabs ci --lang flutter --vcs github
```

---

## 🏗️ Tech Stack

- **Language:** Go
- **Framework:** Cobra
- **Build:** Single static binary
- **Platforms:** macOS, Linux, Windows

---

## 🤝 Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🔗 Links

- **GitHub:** [https://github.com/MayR-Labs/mayrlabs-go](https://github.com/MayR-Labs/mayrlabs-go)
- **Issues:** [Report a Bug](https://github.com/MayR-Labs/mayrlabs-go/issues)
- **Documentation:** [Examples](examples/README.md) | [Development Guide](DEVELOPMENT.md)

---

## 💡 Need Help?

- Check the [examples directory](examples/) for more usage scenarios
- Run `mayrlabs [command] --help` for command-specific help
- Open an issue on GitHub for bug reports or feature requests

---

**Built with ❤️ by MayR Labs**
mayrlabs --help
mayrlabs <command> --help
```

---

## ⚙️ Commands

### 🧱 General Commands

| Command                             | Description                                                            |
| ----------------------------------- | ---------------------------------------------------------------------- |
| `mayrlabs create-keystore`          | Create a new keystore interactively.                                   |
| `mayrlabs dns-clear`                | Clear the DNS cache (choose macOS, Linux, or Windows).                 |
| `mayrlabs ci`                       | Generate CI/CD workflow YAML for your language and VCS.                |
| `mayrlabs format [language]`        | Format project files for a given language (interactive if omitted).    |
| `mayrlabs add-license`              | Create a LICENSE file based on selected type, year, and author.        |
| `mayrlabs editor-config [language]` | Generate `.editorconfig` for a specific language (supports `--force`). |
| `mayrlabs hash [string]`            | Generate hash of a string using md5, sha1, or sha256.                  |
| `mayrlabs uuid`                     | Generate UUID v4.                                                      |
| `mayrlabs password [length]`        | Generate a random password of specified length (default: 16).          |
| `mayrlabs version`                  | Display the version of mayrlabs CLI.                                   |
| `mayrlabs visit`                    | Open the MayR Labs website in your browser.                            |
| `mayrlabs github`                   | Open the mayrlabs-go GitHub repository in your browser.                |

---

### 🌿 Git Commands

| Command                    | Description                                        |
| -------------------------- | -------------------------------------------------- |
| `mayrlabs git`             | Show available Git-related commands.               |
| `mayrlabs git prune-stale` | Delete all local branches not found on the remote. |

---

### 🔐 ENV Commands

| Command                                | Description                                                                           |
| -------------------------------------- | ------------------------------------------------------------------------------------- |
| `mayrlabs env`                         | List available environment commands.                                                  |
| `mayrlabs env update-example [source]` | Sync `.env.example` with `.env` or `.env.staging`. Creates `.env.example` if missing. |
| `mayrlabs env validate`                | Check for missing keys, invalid values, or duplicated variables.                      |
| `mayrlabs env arrange [file]`          | Sort and group environment keys by prefix (e.g., `APP_*`, `DB_*`).                    |

---

### 📝 CHANGELOG Commands

| Command                                         | Description                                           |
| ----------------------------------------------- | ----------------------------------------------------- |
| `mayrlabs changelog`                            | Display all changelog commands.                       |
| `mayrlabs changelog create [--force]`           | Create or overwrite `CHANGELOG.md`.                   |
| `mayrlabs changelog record [version] [summary]` | Add a new entry to `CHANGELOG.md` (supports `--wip`). |

---

### 🐦 Flutter Commands

| Command                           | Description                                                         |
| --------------------------------- | ------------------------------------------------------------------- |
| `mayrlabs flutter`                | List Flutter-related commands.                                      |
| `mayrlabs flutter create-scripts` | Add useful build scripts to `scripts/` (IPA, APK, AppBundle, etc.). |

---

### 🐘 PHP Commands

| Command               | Description                   |
| --------------------- | ----------------------------- |
| `mayrlabs php`        | List PHP commands.            |
| `mayrlabs php cs-fix` | Run PHP CodeSniffer/CS-Fixer. |
| `mayrlabs php lint`   | Lint PHP files.               |

---

### ⚡ JavaScript Commands

| Command                      | Description                                             |
| ---------------------------- | ------------------------------------------------------- |
| `mayrlabs js`                | List JavaScript commands.                               |
| `mayrlabs js setup-prettier` | Install and configure Prettier with `.prettierrc.yaml`. |
| `mayrlabs js pretty`         | Run Prettier on the project.                            |

---

## 🪄 Bonus (Fun / Easter Eggs)

| Command           | Description                                         |
| ----------------- | --------------------------------------------------- |
| `mayrlabs quote`  | Display a random motivational quote for developers. |
| `mayrlabs visit`  | Open the MayR Labs website in your browser.         |
| `mayrlabs github` | Open the GitHub repository in your browser.         |

---

## 🧠 Examples

```bash
# Create a CI file for a Flutter project using GitHub Actions
mayrlabs ci --lang flutter --vcs github

# Sync .env.example with .env
mayrlabs env update-example

# Format Go code
mayrlabs format go

# Add MIT license for MayR Labs, 2025
mayrlabs add-license --type mit --author "MayR Labs" --year 2025

# Generate default scripts for a Flutter project
mayrlabs flutter create-scripts
```

---

## 🏗️ Tech Stack

- **Language:** Go
- **Framework:** Cobra
- **Build:** Single static binary
- **Platforms:** macOS, Linux, Windows

---
