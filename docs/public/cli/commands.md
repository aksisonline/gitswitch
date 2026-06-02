---
title: CLI Commands
description: Complete reference for every gitswitch command and flag
---

This page is the exhaustive reference for all `gitswitch` commands, flags, and exit codes.

---

## `gitswitch` — open the TUI

```bash
gitswitch
```

Opens the interactive TUI. On first run with no existing profiles, imports current `git config --global user.name` and `git config --global user.email` as a profile named `default`.

**TUI keybindings**

| Key | Action |
|-----|--------|
| `↑` / `↓` or `k` / `j` | Move cursor |
| `Enter` | Switch to selected profile |
| `a` | Open add-profile form |
| `e` | Open edit-profile form |
| `ctrl+d` (in edit form) | Delete profile |
| `c` | Cycle color theme (12 themes) |
| `?` | Show CLI quick-reference screen |
| `u` | Upgrade (shown only when a newer version is available) |
| `q` / `ctrl+c` | Quit |

---

## `gitswitch <nickname>` — quick switch

```bash
gitswitch work
```

Switches to the named profile immediately and exits. Equivalent to `gitswitch switch <nickname>`. Useful in scripts and shell aliases.

**Output:**

```
✓ Switched to 'work' — Alice Smith <alice@company.com>
```

Exits non-zero if the nickname does not exist.

---

## `gitswitch add` — add a profile

```bash
gitswitch add <nickname> <user-name> <email> [flags]
```

**Arguments**

| Argument | Description |
|----------|-------------|
| `nickname` | Short label for the profile (e.g. `work`, `oss`). Not written to git config. Must be unique. |
| `user-name` | Value for `git config user.name`. Quote names with spaces. |
| `email` | Value for `git config user.email`. |

**Flags**

| Flag | Description |
|------|-------------|
| `--sign-key <key>` | GPG key ID. Sets `git config user.signingkey` on switch. |
| `--ssh-key <path>` | Path to SSH private key. Sets `core.sshCommand` to `ssh -i <path> -o IdentitiesOnly=yes` on switch. |
| `--gh-user <username>` | GitHub CLI username. Runs `gh auth switch --user <username>` on switch. Best-effort — fails gracefully if `gh` is not installed or the account is not logged in. |

**Examples**

```bash
# Git identity only
gitswitch add personal "Alice Smith" alice@gmail.com

# With SSH key
gitswitch add work "Alice Smith" alice@company.com \
  --ssh-key ~/.ssh/id_work

# Full profile
gitswitch add corp "Alice Smith" alice@corp.com \
  --sign-key ABCD1234EF567890 \
  --ssh-key ~/.ssh/id_corp \
  --gh-user alice-corp
```

---

## `gitswitch switch` — switch by name

```bash
gitswitch switch <nickname>
```

Identical to `gitswitch <nickname>`. Both forms apply changes in the same order.

---

## `gitswitch list` — list profiles

```bash
gitswitch list
```

Prints all saved profiles. A `✓` marks the active one.

```
✓  personal        Alice Smith <alice@gmail.com>
   work            Alice Smith <alice@company.com>
   corp            Alice Smith <alice@corp.com>
```

---

## `gitswitch current` — show active profile

```bash
gitswitch current [flags]
```

Prints the nickname, name, and email of the currently active profile.

```
work — Alice Smith <alice@company.com>
```

Prints `No active profile` if none has been applied yet.

**Flags**

| Flag | Description |
|------|-------------|
| `--short` | Outputs `nickname\temail` tab-separated. Used by the Starship prompt block. |
| `--prompt` | Outputs `nickname\tcolor` tab-separated, where color is the ANSI 256-color index for the current theme's primary color. Used by shell prompt functions. |

---

## `gitswitch remove` — delete a profile

```bash
gitswitch remove <nickname>
```

Permanently removes the profile from `~/.config/gitswitch/profiles.json`. Does not revert any git config that was applied when it was last active.

---

## `gitswitch init` — import current git config

```bash
gitswitch init
```

Reads `git config --global user.name` and `git config --global user.email` and saves them as a profile named `default`. Runs automatically on first launch if no profiles exist.

---

## `gitswitch install` — set up shell integration

```bash
gitswitch install [flags]
```

Runs an interactive wizard that installs:

- Shell prompt segment (shows active identity in the prompt when inside a git repo)
- Identity nudge hook (triggers `gitswitch recommend` on `cd` into a repo)
- Tab completion
- HTTPS credential helper (optional — routes HTTPS git operations through the active profile's `gh` account)

**Shell framework detection**

| Framework | Strategy |
|-----------|----------|
| Starship | Appends `[custom.gitswitch]` block to `~/.config/starship.toml` |
| oh-my-zsh | Creates plugin at `~/.oh-my-zsh/custom/plugins/gitswitch/gitswitch.plugin.zsh` |
| Powerlevel10k | Drops segment function to rc file; prints manual step |
| Raw zsh / bash / fish | Appends prompt + nudge + completion snippet to rc file |

Idempotent — uses a `# gitswitch shell integration` marker to skip re-installation.

**Flags**

| Flag | Description |
|------|-------------|
| `--shell <shell>` | Override shell detection. Values: `zsh`, `bash`, `fish`. Also skips the interactive wizard. |
| `--yes` / `-y` | Accept all defaults without prompts. For scripts and CI. |
| `--https` | Register the HTTPS credential helper (default: `true` when using `--yes`). |

---

## `gitswitch uninstall` — remove shell integration

```bash
gitswitch uninstall [flags]
```

Removes the shell integration block written by `gitswitch install` from the rc file (or oh-my-zsh plugin directory) and unregisters the HTTPS credential helper if installed.

**Flags**

| Flag | Description |
|------|-------------|
| `--shell <shell>` | Target shell. Values: `zsh`, `bash`, `fish`. Default: auto-detect. |

---

## `gitswitch pin` — pin an identity to a repo

```bash
gitswitch pin <nickname>
```

Marks the given profile as the permanent recommended identity for the current repo. The pin is stored in `~/.config/gitswitch/history.json` under the repo's remote URL key — no files are written to the repo itself.

Must be run from inside a git repo. Validates that the nickname exists.

---

## `gitswitch unpin` — remove a repo pin

```bash
gitswitch unpin
```

Clears the pinned identity for the current repo. `gitswitch recommend` falls back to usage-count-based recommendations.

Must be run from inside a git repo.

---

## `gitswitch record` — record identity for a repo

```bash
gitswitch record [flags]
```

Increments the usage counter for the currently active profile under the current repo's key. Called automatically by the shell nudge hook on every `cd` into a repo — you rarely need to run this manually.

**Flags**

| Flag | Description |
|------|-------------|
| `--path <dir>` | Directory to record for. Default: current working directory. |

---

## `gitswitch recommend` — print recommendation for current repo

```bash
gitswitch recommend [flags]
```

Checks usage history and any pin for the current repo and prints the recommended identity.

**Exit behavior**

- Exits `0` and prints `nickname\tname\temail` when a recommendation is warranted
- Exits `1` silently when already on the right identity, no history exists, or the threshold is not met

**Recommendation threshold (no pin):** top identity has ≥ 3 entries and ≥ 60% share and differs from the current identity.

**Flags**

| Flag | Description |
|------|-------------|
| `--path <dir>` | Directory to check. Default: current working directory. |

---

## `gitswitch claude` — install the Claude Code skill

```bash
gitswitch claude [flags]
```

Installs the embedded gitswitch skill into Claude Code. The SKILL.md is embedded in the binary — no network request.

**Flags**

| Flag | Description |
|------|-------------|
| `--scope user` | Install to `~/.claude/skills/gitswitch/` — active for all projects (default) |
| `--scope project` | Install to `.claude/skills/gitswitch/` — active for this project only |

After installing, reload Claude Code or open a new session.

---

## `gitswitch version` — show version info

```bash
gitswitch version
```

Prints the installed version and checks for a newer release (cached for 24 hours).

```
gitswitch v0.1.11
New version available: v0.1.12
Run: gitswitch upgrade
```

---

## `gitswitch upgrade` — upgrade to latest

```bash
gitswitch upgrade
```

Downloads and runs the install script for the latest release. Exits early if already on the latest version.

If installed via Homebrew, prints the correct command instead:

```
gitswitch was installed via Homebrew.
Run: brew upgrade gitswitch
```

---

## `gitswitch pacman` — arcade mode

```bash
gitswitch pacman
```

Opens the TUI with the arcade intro animation and double-border style. All normal TUI functionality works — this is a cosmetic variant.

---

## What switching does

Every switch (CLI or TUI) applies changes in this order:

1. `git config --global user.name "<name>"`
2. `git config --global user.email "<email>"`
3. `git config --global user.signingkey "<key>"` — only if the profile has a GPG key
4. `git config --global core.sshCommand "ssh -i <path> -o IdentitiesOnly=yes"` — only if the profile has an SSH key
5. `gh auth switch --user <username>` — only if the profile has a GitHub username; warning only on failure

---

## Exit codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | General error (profile not found, git not installed, etc.) |

`gitswitch recommend` exits `1` silently when no recommendation is available — this is expected behavior used by the shell hook.

---

## Storage files

| File | Contents |
|------|----------|
| `~/.config/gitswitch/profiles.json` | All profiles |
| `~/.config/gitswitch/config.json` | UI preferences (color theme index) |
| `~/.config/gitswitch/history.json` | Per-repo identity usage counts and pins |
