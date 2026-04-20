#!/bin/bash
set -e

# Detect OS and architecture
OS=$(uname -s)
ARCH=$(uname -m)

case "$OS" in
  Darwin)
    case "$ARCH" in
      arm64) BINARY="gitswitch-darwin-arm64" ;;
      x86_64) BINARY="gitswitch-darwin-amd64" ;;
      *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
    esac
    ;;
  Linux)
    case "$ARCH" in
      aarch64) BINARY="gitswitch-linux-arm64" ;;
      x86_64) BINARY="gitswitch-linux-amd64" ;;
      *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
    esac
    ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac

# Use version from arg or latest
VERSION="${1:-LATEST_VERSION}"

# Download URL
RELEASE_URL="https://github.com/aksisonline/gitswitch/releases/download/$VERSION/$BINARY"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

echo "Installing gitswitch ($VERSION) for $OS/$ARCH..."

# Create temp directory
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

# Download binary
echo "Downloading $BINARY..."
if ! curl -fsSL "$RELEASE_URL" -o "$TEMP_DIR/gitswitch"; then
  echo "Failed to download gitswitch. Check version or internet connection."
  exit 1
fi

# Make executable
chmod +x "$TEMP_DIR/gitswitch"

# Verify it works
if ! "$TEMP_DIR/gitswitch" --help > /dev/null 2>&1; then
  echo "Downloaded binary failed sanity check."
  exit 1
fi

# Install
if [ ! -w "$INSTALL_DIR" ]; then
  echo "sudo required to install to $INSTALL_DIR"
  sudo mv "$TEMP_DIR/gitswitch" "$INSTALL_DIR/gitswitch"
else
  mv "$TEMP_DIR/gitswitch" "$INSTALL_DIR/gitswitch"
fi

echo "✓ gitswitch installed to $INSTALL_DIR/gitswitch"
gitswitch --help | head -1
