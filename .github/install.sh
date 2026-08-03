#!/bin/bash
# gitswitch installer — https://get.gitswitch.dev
#
# Humans / scripts:  curl -fsSL https://get.gitswitch.dev | bash
# Agents:            fetch this URL as text — this file is both the
#                     installer and its own instructions. The code below
#                     just installs the binary; read on for what to do
#                     with it once it's there.
#
# ── After this script finishes ──────────────────────────────────────────
#
#   1. gitswitch doctor --json
#      Reports git/gh presence, HTTPS routing, Session Isolation as JSON.
#      Resolve anything it flags before moving on.
#
#   2. gitswitch login
#      GitHub's OAuth device flow: it prints a URL and a short code. A
#      human has to open that URL in a browser and approve it — an agent
#      can't complete this step on its own. Show the user the URL and
#      code, ask them to approve it, then wait; the command blocks until
#      they do (or it times out). Don't try to script around it.
#
#   3. gitswitch doctor --json   (again, to confirm)
#      "profiles" should now be >= 1 and "gh": {"installed": true}. If
#      shell integration ran too (gitswitch shell), "https.routed_by_gitswitch"
#      and "gh_wrapper_installed" should both be true.
#
# Full CLI reference: https://gitswitch.dev/docs/cli/commands
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

# Use version from arg, or resolve latest from GitHub API
VERSION="${1:-}"
if [ -z "$VERSION" ]; then
  VERSION=$(curl -fsSL "https://api.github.com/repos/aksisonline/gitswitch/releases/latest" \
    -H "Accept: application/vnd.github+json" \
    | grep '"tag_name"' | head -1 | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/')
  if [ -z "$VERSION" ]; then
    echo "Could not determine latest version. Pass a version explicitly: install.sh v0.2.1"
    exit 1
  fi
fi

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

# Create gs short alias (symlink → gitswitch)
GS_PATH="$INSTALL_DIR/gs"
if [ -e "$GS_PATH" ] && [ ! -L "$GS_PATH" ]; then
  echo "⚠  $GS_PATH already exists and is not a symlink — skipping gs alias"
  echo "   (possible conflict with another tool; gs will not be created)"
else
  if [ ! -w "$INSTALL_DIR" ]; then
    sudo ln -sf "$INSTALL_DIR/gitswitch" "$GS_PATH"
  else
    ln -sf "$INSTALL_DIR/gitswitch" "$GS_PATH"
  fi
  echo "✓ gs alias created at $GS_PATH"
fi

echo ""
echo "Run  gs  to get started — first run sets up git/gh and your account automatically."
