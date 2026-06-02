---
title: Installation
description: Install gitswitch on macOS and Linux
---

gitswitch runs on macOS and Linux. Choose any of the methods below.

## Homebrew

```bash
brew install aksisonline/tap/gitswitch
```

Update with:

```bash
brew upgrade gitswitch
```

> **Tip:** Homebrew is the recommended method — it handles PATH, updates, and uninstallation cleanly.

## curl installer

Works on macOS and Linux:

```bash
curl -fsSL https://raw.githubusercontent.com/aksisonline/gitswitch/main/.github/install.sh | bash
```

The script downloads the appropriate binary for your OS and architecture to `/usr/local/bin`.

To upgrade later:

```bash
gitswitch upgrade
```

## Go install

Requires Go 1.21 or later:

```bash
go install github.com/aksisonline/gitswitch@latest
```

## Build from source

```bash
git clone https://github.com/aksisonline/gitswitch
cd gitswitch
make install
```

## Verify

```bash
gitswitch version
```

Expected output:

```
gitswitch v0.1.11
Already on latest version.
```

If a newer release is available:

```
gitswitch v0.1.11
New version available: v0.1.12
Run: gitswitch upgrade
```

## Prerequisites

- `git` must be installed and on `PATH` — gitswitch checks this at startup and exits with an error if not found
- `gh` (GitHub CLI) is optional — required only if you use the `--gh-user` flag on profiles

## Shell integration

After installing the binary, set up the shell prompt and identity nudges:

```bash
gitswitch install
```

See [Shell Integration](/docs/features/shell) for details.

## Next steps

- [Quick Start](/docs/quick-start)
- [Shell Integration](/docs/features/shell)
- [CLI Reference](/docs/cli/commands)
