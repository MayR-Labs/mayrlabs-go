#!/bin/bash

# MayR Labs CLI Installation Script
# Usage: curl -sSL https://raw.githubusercontent.com/MayR-Labs/mayrlabs-go/main/install.sh | bash
# or: wget -qO- https://raw.githubusercontent.com/MayR-Labs/mayrlabs-go/main/install.sh | bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Convert architecture names
case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

# Map OS names
case "$OS" in
    linux)
        OS="linux"
        ;;
    darwin)
        OS="darwin"
        ;;
    mingw*|msys*|cygwin*)
        OS="windows"
        ;;
    *)
        echo -e "${RED}Unsupported operating system: $OS${NC}"
        exit 1
        ;;
esac

# Set binary name
BINARY_NAME="mayrlabs"
if [ "$OS" = "windows" ]; then
    BINARY_NAME="mayrlabs.exe"
fi

# GitHub repository details
REPO="MayR-Labs/mayrlabs-go"
RELEASE_URL="https://api.github.com/repos/$REPO/releases/latest"

echo -e "${GREEN}🧰 MayR Labs CLI Installer${NC}"
echo ""
echo "Detected system: $OS-$ARCH"
echo ""

# Get latest release version
echo "Fetching latest release..."
LATEST_VERSION=$(curl -s "$RELEASE_URL" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    echo -e "${RED}Failed to fetch latest release version${NC}"
    exit 1
fi

echo "Latest version: $LATEST_VERSION"

# Construct download URL
DOWNLOAD_FILE="mayrlabs-${OS}-${ARCH}"
if [ "$OS" = "windows" ]; then
    DOWNLOAD_FILE="${DOWNLOAD_FILE}.exe"
fi

DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_VERSION/$DOWNLOAD_FILE"

echo ""
echo "Downloading $DOWNLOAD_FILE..."

# Download the binary
TMP_DIR=$(mktemp -d)
TMP_FILE="$TMP_DIR/$BINARY_NAME"

if command -v curl >/dev/null 2>&1; then
    curl -L -o "$TMP_FILE" "$DOWNLOAD_URL"
elif command -v wget >/dev/null 2>&1; then
    wget -O "$TMP_FILE" "$DOWNLOAD_URL"
else
    echo -e "${RED}Error: curl or wget is required but not found${NC}"
    exit 1
fi

# Make executable
chmod +x "$TMP_FILE"

# Determine installation directory
INSTALL_DIR=""
if [ "$OS" = "windows" ]; then
    # On Windows, install to user's local bin
    INSTALL_DIR="$HOME/bin"
    mkdir -p "$INSTALL_DIR"
else
    # Try to install to /usr/local/bin (requires sudo)
    if [ -w "/usr/local/bin" ]; then
        INSTALL_DIR="/usr/local/bin"
    else
        # Fallback to user's local bin
        INSTALL_DIR="$HOME/.local/bin"
        mkdir -p "$INSTALL_DIR"
    fi
fi

# Install the binary
echo ""
echo "Installing to $INSTALL_DIR..."

if [ "$INSTALL_DIR" = "/usr/local/bin" ] && [ ! -w "$INSTALL_DIR" ]; then
    # Need sudo for /usr/local/bin
    echo -e "${YELLOW}Installing to $INSTALL_DIR requires sudo privileges${NC}"
    sudo mv "$TMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
else
    mv "$TMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
fi

# Clean up
rm -rf "$TMP_DIR"

echo ""
echo -e "${GREEN}✅ MayR Labs CLI installed successfully!${NC}"
echo ""

# Check if installation directory is in PATH
if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
    echo -e "${YELLOW}⚠️  Warning: $INSTALL_DIR is not in your PATH${NC}"
    echo ""
    echo "Add it to your PATH by adding this line to your shell configuration file:"
    echo ""
    
    if [ "$OS" = "darwin" ] || [ "$OS" = "linux" ]; then
        SHELL_NAME=$(basename "$SHELL")
        case "$SHELL_NAME" in
            bash)
                CONFIG_FILE="~/.bashrc"
                ;;
            zsh)
                CONFIG_FILE="~/.zshrc"
                ;;
            fish)
                CONFIG_FILE="~/.config/fish/config.fish"
                echo "  set -gx PATH $INSTALL_DIR \$PATH"
                ;;
            *)
                CONFIG_FILE="~/.profile"
                ;;
        esac
        
        if [ "$SHELL_NAME" != "fish" ]; then
            echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
        fi
        
        echo ""
        echo "Then run:"
        echo "  source $CONFIG_FILE"
    fi
    echo ""
fi

# Verify installation
echo "Verifying installation..."
if command -v mayrlabs >/dev/null 2>&1; then
    VERSION=$(mayrlabs version 2>&1 || echo "unknown")
    echo -e "${GREEN}$VERSION${NC}"
    echo ""
    echo "Run 'mayrlabs --help' to get started!"
else
    echo -e "${YELLOW}mayrlabs command not found. Please add $INSTALL_DIR to your PATH${NC}"
fi

echo ""
echo "📚 Documentation: https://github.com/$REPO"
echo "🌐 Website: https://mayrlabs.com"
