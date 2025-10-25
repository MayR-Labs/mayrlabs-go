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

## 🧩 Installation

```bash
go install github.com/mayrlabs/mayrlabs-cli@latest
```

or clone and build manually:

```bash
git clone https://github.com/mayrlabs/mayrlabs-cli.git
cd mayrlabs-cli
go build -o mayrlabs
```

---

## 🧭 Usage

```bash
mayrlabs [command] [flags]
```

Run interactive mode (prompts for actions):

```bash
mayrlabs
```

Get help:

```bash
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

| Command          | Description                                         |
| ---------------- | --------------------------------------------------- |
| `mayrlabs quote` | Display a random motivational quote for developers. |

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
