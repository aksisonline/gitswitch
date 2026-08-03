---
title: Versions & Release Channels
description: Check your version, upgrade, and switch between stable and canary builds
---

## What am I on?

```bash
gitswitch version
```

```
gitswitch v0.3.0
Already on latest version.
```

```
gitswitch v0.2.9
New version available: v0.3.0
Run: gitswitch upgrade
```

The check hits the GitHub Releases API and is cached for 24 hours, so it costs nothing on a normal run.

## Upgrading

```bash
gitswitch upgrade
```

Downloads and installs the latest release. Installed via Homebrew? It notices and tells you the right command instead of fighting your package manager:

```
gitswitch was installed via Homebrew.
Run: brew upgrade gitswitch
```

You can also press `u` in the TUI — the key only appears when there's actually something newer.

## Two channels

| Channel | What it is |
|---|---|
| **stable** | Normal releases. This is you, unless you chose otherwise. |
| **canary** | Pre-release builds — new features first, rough edges included. |

```bash
gitswitch beta      # move to canary
gitswitch stable    # move back
```

Both confirm before doing anything, and both leave `~/.config/gitswitch/` completely alone — your accounts, pins, and history survive the trip in either direction.

If you're on Homebrew, `gitswitch beta` has to remove the formula first (canary builds aren't published to the tap) and installs the binary to `/usr/local/bin`. It tells you that up front and asks before proceeding. Going back to stable, `gitswitch stable` will suggest re-installing via Homebrew if you prefer.

Restart your shell afterwards so the new binary is picked up.

### Update notifications stay in your lane

Canary builds are only offered newer canary builds; stable builds are only offered stable releases. Moving between channels is something you do deliberately, never something an update prompt does to you.

## "What's New"

After an upgrade that includes new features, gitswitch shows a one-time **What's New** screen on next launch — the release notes, in the app, so you don't have to go looking for a changelog.

It appears when the installed version's MINOR or MAJOR number goes up, and for every canary switch (canary versions only differ in their last number, so gitswitch schedules the screen explicitly for those).

## "Shell integration updated — run `gitswitch install`"

Occasionally a release improves the shell hook, prompt, or completion script. When that happens gitswitch tells you once, and `gitswitch install` replaces the gitswitch block in your rc file in place.

It only nags when your installed shell setup is genuinely out of date — not on every version bump.

## Next

- **[Installation](/docs/installation)** — first-time install and uninstall
- **[CLI Reference](/docs/cli/commands)** — every command
