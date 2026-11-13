---
title: Installation
description: Install gitswitch on macOS or Linux
---
gitswitch is a single Go binary. macOS and Linux, Intel and Apple Silicon/ARM.

## Homebrew

```bash
brew install aksisonline/tap/gitswitch
```

Recommended — it handles `PATH`, updates, and uninstall cleanly, and gitswitch respects it: every later update goes through `brew`, never around it.

```bash
brew upgrade gitswitch   # later
```

Want pre-release builds? The beta channel is a second formula in the same tap — see [Release Channels](/docs/cli/channels).

## curl

```bash
curl -fsSL https://get.gitswitch.dev | bash
```

Downloads the right binary for your OS and architecture into `/usr/local/bin`. Upgrade later with `gitswitch upgrade`.

Pick one install method and stay with it. If Homebrew already manages gitswitch, the script stops and tells you to use `brew upgrade` instead of installing a second copy that would shadow the first.

## Go

```bash
go install github.com/aksisonline/gitswitch/cmd/gitswitch@latest
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
gitswitch --version
```

```
gitswitch v0.5.0
```

`gitswitch --help` lists every command; `gitswitch <command> --help` shows that command's flags and examples.

## What you need installed

| | |
|---|---|
| **`git`** | Required. gitswitch installs it for you on first run if it's missing. |
| **`gh`** ([GitHub CLI](https://cli.github.com)) | Optional, but recommended — [HTTPS push routing](/docs/routing/https) and [Session Isolation](/docs/routing/session-isolation) need it. Same first-run offer. |

Check both at any time:

```bash
gitswitch doctor
```

## Run it

```bash
gitswitch
```

That's it — one command handles both cases:

- **Nothing installed** — installs `git`/`gh` itself, then you log in with GitHub once.
- **Already set up** — imports your git config, `gh` accounts, and SSH keys, matching a `gh` account to your git config by verified email so you don't get duplicates.

Want shell integration too (prompt segment, nudges, HTTPS push routing, Session Isolation)? That's a separate, optional step:

```bash
gitswitch shell
```

Everything is opt-out, defaults are sensible, and it's safe to re-run. See [Shell Integration](/docs/routing/shell).

Then head to the [Quick Start](/docs/get-started/quick-start).

## Staying up to date

```bash
gitswitch version    # what am I on? is there anything newer?
gitswitch upgrade    # get it — through Homebrew if that's how you installed
```

gitswitch also tells you inside the TUI when a new release lands (press `u`), and shows a "What's New" screen after a feature release. Want the pre-release builds? See [Release Channels](/docs/cli/channels).

## Uninstall

```bash
gitswitch uninstall            # remove shell integration, HTTPS routing, gh wrapper
brew uninstall gitswitch       # then remove the binary
rm -rf ~/.config/gitswitch     # and your profiles, if you're sure
```
