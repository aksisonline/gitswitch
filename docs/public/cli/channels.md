---
title: Versions & Release Channels
description: Check your version, upgrade through Homebrew or the install script, and switch between stable and beta builds
---
## What am I on?

```bash
gitswitch --version    # just the number — instant, works offline
gitswitch version      # number, channel, how it's managed, and whether there's an update
```

```
gitswitch v0.5.0
stable channel · managed by Homebrew
Update available: v0.6.0 — run: brew upgrade gitswitch
```

`-v` and `--v` work too. `gitswitch version --json` prints the same facts for scripts. The update check uses the GitHub Releases API and is cached for 24 hours, so it costs nothing on a normal run.

## One rule: whoever installed it, updates it

gitswitch never replaces a file that a package manager owns. How it updates depends on how you installed it:

| Installed with | Updates with | What `gitswitch upgrade` does |
|---|---|---|
| **Homebrew** | `brew` | Runs `brew update && brew upgrade aksisonline/tap/gitswitch` for you |
| **The install script** (`curl … \| bash`) | gitswitch itself | Downloads the new binary from get.gitswitch.dev |

Why so strict? If gitswitch overwrote a Homebrew-installed binary, Homebrew would still believe the old version is installed — the next `brew upgrade` or `brew cleanup` then fights whatever is actually on disk. Keeping every update inside `brew` avoids that entirely.

Not sure which one you have? `gitswitch version` says so on its second line.

## Upgrading

```bash
gitswitch upgrade            # update now
gitswitch upgrade --check    # only say whether an update exists
gitswitch upgrade --dry-run  # show the exact command that would run, change nothing
```

On Homebrew you can also just run `brew upgrade gitswitch` — it is the same thing.

In the TUI, a pending update shows in the status line at the bottom. Press `u` to open the update prompt, which names the command it will run before you confirm. If you answer "not now", the prompt won't reopen by itself for that release — the status-line reminder stays.

## Two channels

| Channel | What it is | Homebrew formula |
|---|---|---|
| **stable** | Normal releases. This is you, unless you chose otherwise. | `aksisonline/tap/gitswitch` |
| **beta** (canary) | Pre-release builds — new features first, rough edges included. | `aksisonline/tap/gitswitch-beta` |

```bash
gitswitch beta      # move to beta
gitswitch stable    # move back
```

Both show what they are about to do and ask first (`--yes` skips the question, `--dry-run` only shows the plan), and both leave `~/.config/gitswitch/` completely alone — your accounts, pins, and history survive the trip in either direction.

**On Homebrew**, switching channel swaps one formula for the other. The two can't be installed together (they provide the same `gitswitch` command), so gitswitch runs:

```bash
brew update
brew uninstall gitswitch
brew install aksisonline/tap/gitswitch-beta
```

Homebrew stays in charge afterwards — beta updates arrive with `brew upgrade gitswitch-beta`. If installing the new formula fails, the old one is put back, so a failed switch never leaves you without gitswitch. You can run those three commands yourself instead; `gitswitch beta` is only a shortcut for them.

**With the install script**, the binary is replaced with the latest build of the other channel.

Right after a stable release there may be no beta that is newer than it. `gitswitch beta` says so and leaves you on stable rather than moving you backwards.

### Update notifications stay in your lane

Canary builds are only offered newer canary builds; stable builds are only offered stable releases. Moving between channels is something you do deliberately, never something an update prompt does to you.

## "What's New"

After an upgrade that includes new features, gitswitch shows a one-time **What's New** screen on next launch — the release notes, in the app, so you don't have to go looking for a changelog.

It appears when the installed version's MINOR or MAJOR number goes up, after any upgrade you start from gitswitch itself, and for every channel switch. A brand-new install skips it — there is no "old version" to compare with.

## "Shell integration updated — run `gitswitch shell`"

Occasionally a release improves the shell hook, prompt, or completion script. When that happens gitswitch tells you once, and `gitswitch shell` replaces the gitswitch block in your rc file in place.

It only nags when your installed shell setup is genuinely out of date — not on every version bump.

## Next

- **[Installation](/docs/get-started/installation)** — first-time install and uninstall
- **[CLI Reference](/docs/cli/commands)** — every command
