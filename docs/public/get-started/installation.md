---
title: Installation
description: Install gitswitch on macOS or Linux
---

gitswitch is a single Go binary. macOS and Linux, Intel and Apple Silicon/ARM.

## Homebrew

```bash
brew install aksisonline/tap/gitswitch
```

Recommended — it handles `PATH`, updates, and uninstall cleanly.

```bash
brew upgrade gitswitch   # later
```

## curl

```bash
curl -fsSL https://raw.githubusercontent.com/aksisonline/gitswitch/main/.github/install.sh | bash
```

Downloads the right binary for your OS and architecture into `/usr/local/bin`. Upgrade later with `gitswitch upgrade`.

## Go

```bash
go install github.com/aksisonline/gitswitch@latest
```

Requires Go 1.21+.

## From source

```bash
git clone https://github.com/aksisonline/gitswitch
cd gitswitch
make install
```

## Check it worked

```bash
gitswitch version
```

```
gitswitch v0.3.0
Already on latest version.
```

## What you need installed

| | |
|---|---|
| **`git`** | Required. Missing? The first time you run bare `gitswitch` on a machine with no profiles yet, it offers to install it for you — no separate step needed. |
| **`gh`** ([GitHub CLI](https://cli.github.com)) | Optional, but recommended. gitswitch works fine without it — but [HTTPS push routing](/docs/routing/https) and [Session Isolation](/docs/routing/session-isolation) get their tokens from `gh`, so those two features install but stay inert until `gh` is set up. That first-run check offers to install it too. |

Check both at any time:

```bash
gitswitch doctor
```

## Run it

```bash
gitswitch
```

That's the whole "then set it up" step — first run on a machine with no profiles checks for git/gh and offers to install what's missing, then opens a short wizard that gets you a working account.

Want shell integration too (prompt segment, nudges, HTTPS push routing, Session Isolation)? That's a separate, optional step:

```bash
gitswitch shell
```

Everything is opt-out, defaults are sensible, and it's safe to re-run. See [Shell Integration](/docs/routing/shell).

Then head to the [Quick Start](/docs/get-started/quick-start).

## Staying up to date

```bash
gitswitch version    # what am I on? is there anything newer?
gitswitch upgrade    # get it
```

gitswitch also tells you inside the TUI when a new release lands, and shows a "What's New" screen after a feature release. Want the pre-release builds? See [Release Channels](/docs/cli/channels).

## Uninstall

```bash
gitswitch uninstall            # remove shell integration, HTTPS routing, gh wrapper
brew uninstall gitswitch       # then remove the binary
rm -rf ~/.config/gitswitch     # and your profiles, if you're sure
```
