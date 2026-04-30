# CLI Reference

Full reference for every `gitswitch` command and flag.

---

## `gitswitch` — open the TUI

```bash
gitswitch
```

Opens the interactive terminal UI. On first run, auto-imports your existing `git config` as a `default` profile.

**TUI key bindings**

| Key | Action |
|-----|--------|
| `↑` / `↓` or `k` / `j` | Move cursor |
| `enter` | Switch to selected profile |
| `a` | Open add-profile form |
| `e` | Open edit-profile form |
| `ctrl+d` (in edit form) | Delete profile |
| `?` | Show CLI quick-reference screen |
| `u` | Upgrade to latest version (shown only when update available) |
| `q` / `ctrl+c` | Quit |

---

## `gitswitch <nickname>` — quick switch

```bash
gitswitch work
```

Switches to the named profile immediately and exits. No UI. Useful in scripts or shell aliases.

Exits with a non-zero status if the nickname doesn't exist.

---

## `gitswitch add` — add a profile

```bash
gitswitch add <nickname> <user-name> <email> [flags]
```

**Arguments**

| Argument | Description |
|----------|-------------|
| `nickname` | Short label used to identify the profile (e.g. `work`, `oss`). Not written to git config. |
| `user-name` | Value for `git config user.name`. Use quotes for names with spaces. |
| `email` | Value for `git config user.email`. |

**Flags**

| Flag | Description |
|------|-------------|
| `--sign-key <key>` | GPG key ID. Sets `git config user.signingkey`. |
| `--ssh-key <path>` | Path to SSH private key (e.g. `~/.ssh/id_work`). Sets `core.sshCommand` to force that key and disable SSH agent fallback. |
| `--gh-user <username>` | GitHub CLI username. Runs `gh auth switch --user <username>` on every switch. Fails gracefully if `gh` is not installed. |

**Examples**

```bash
# Minimal — just git identity
gitswitch add personal "Alice Smith" alice@gmail.com

# With SSH key for a work repo host
gitswitch add work "Alice Smith" alice@company.com --ssh-key ~/.ssh/id_work

# Full setup: GPG signing + SSH + GitHub CLI
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

Same as `gitswitch <nickname>` (the positional quick-switch form). Both forms apply the same changes in the same order.

---

## `gitswitch list` — list profiles

```bash
gitswitch list
```

Prints all saved profiles. A `✓` marks the active one.

```
✓  personal        alice@gmail.com
   work            alice@company.com
   oss             alice@example.com
```

---

## `gitswitch current` — show active profile

```bash
gitswitch current [--short]
```

Prints the nickname, name, and email of the currently active profile.

```
work — Alice Smith <alice@company.com>
```

**Flags**

| Flag | Description |
|------|-------------|
| `--short` | Output `nickname\temail` tab-separated. Used by Starship prompt integration. |

Returns `No active profile` if none has been applied yet (omitted in `--short` mode).

---

## `gitswitch remove` — delete a profile

```bash
gitswitch remove <nickname>
```

Permanently removes the profile from storage. Does **not** revert any git config that was applied when it was last switched to.

---

## `gitswitch init` — import current git config

```bash
gitswitch init
```

Reads `git config --global user.name` and `git config --global user.email` and saves them as a profile named `default`. Useful for bootstrapping from an existing setup.

Also run automatically on first launch if no profiles exist yet.

---

## `gitswitch install` — set up shell integration

```bash
gitswitch install [--shell zsh|bash|fish]
```

Installs three shell features at once:

- **Prompt segment** — shows active git identity in the prompt when inside a repo
- **Identity nudge** — when you enter a new git repo, nudges if a different identity is usually used there
- **Tab completion** — completes commands and profile nicknames

Auto-detects your prompt framework:

| Framework | Strategy |
|-----------|----------|
| Starship | Appends `[custom.gitswitch]` to `~/.config/starship.toml` |
| oh-my-zsh | Creates `~/.oh-my-zsh/custom/plugins/gitswitch/gitswitch.plugin.zsh` |
| Powerlevel10k | Drops segment function to rc file, prints manual step |
| Raw zsh/bash/fish | Appends prompt + nudge + completion snippet to rc file |

Idempotent — uses a `# gitswitch shell integration` marker to skip if already installed.

**Flags**

| Flag | Description |
|------|-------------|
| `--shell <shell>` | Override shell detection. Values: `zsh`, `bash`, `fish`. |

---

## `gitswitch pin` — pin an identity to a repo

```bash
gitswitch pin <nickname>
```

Permanently marks a profile as the recommended identity for the current repo. The pin takes priority over all usage-count data — `gitswitch recommend` will always return the pinned identity regardless of history.

Must be run from inside a git repo. Validates that the nickname exists.

---

## `gitswitch unpin` — remove a repo pin

```bash
gitswitch unpin
```

Clears the pinned identity for the current repo. The auto-recommender falls back to usage counts.

Must be run from inside a git repo.

---

## `gitswitch record` — log current identity for this repo

```bash
gitswitch record [--path <dir>]
```

Records the currently active profile against the current repo's remote URL (or path if no remote). Called automatically by the shell nudge hook each time you enter a new repo — you rarely need to run this manually.

**Flags**

| Flag | Description |
|------|-------------|
| `--path <dir>` | Directory to record for. Defaults to current working directory. |

---

## `gitswitch recommend` — check what would be recommended

```bash
gitswitch recommend [--path <dir>]
```

Checks usage history (and any pin) for the current repo and prints the recommended identity if one exists.

- **Exits 0** and prints `nickname\tname\temail` if a recommendation is warranted
- **Exits 1** silently if already on the right identity, no history, or threshold not met

Used internally by the shell nudge hook. Useful for scripting or debugging.

Recommendation threshold (auto-learned, no pin): top identity has **≥3 entries** and **≥60% share** and differs from the current identity.

**Flags**

| Flag | Description |
|------|-------------|
| `--path <dir>` | Directory to check. Defaults to current working directory. |

---

## `gitswitch claude` — install the Claude Code skill

```bash
gitswitch claude [--scope user|project]
```

Installs the git-switcher skill into Claude Code. The skill teaches Claude to detect and fix git identity problems automatically.

The SKILL.md is embedded in the binary — no network call required, always matches your installed version.

**Flags**

| Flag | Description |
|------|-------------|
| `--scope user` | Install to `~/.claude/skills/git-switcher/` (default — all projects) |
| `--scope project` | Install to `.claude/skills/git-switcher/` (this project only) |

After installing, reload Claude Code or open a new session to activate.

---

## `gitswitch version` — show version info

```bash
gitswitch version
```

Prints the installed version and checks for a newer release (via a 24-hour cache).

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

Downloads and runs the install script for the latest release. Bypasses the version cache. Exits early if already on the latest version.

---

## What switching does

Every switch (CLI or TUI) applies changes in this order:

1. `git config --global user.name "<name>"`
2. `git config --global user.email "<email>"`
3. `git config --global user.signingkey "<key>"` — if the profile has a GPG key set
4. `git config --global core.sshCommand "ssh -i <path> -o IdentitiesOnly=yes"` — if the profile has an SSH key set
5. `gh auth switch --user <username>` — if the profile has a GitHub username set; warning only on failure

Step 5 is best-effort. If `gh` is not installed, or the account isn't logged in to `gh`, the git config changes (steps 1–4) are still applied.

---

## Storage

| File | Contents |
|------|----------|
| `~/.config/gitswitch/profiles.json` | All profiles |
| `~/.config/gitswitch/config.json` | UI preferences (color theme) |
| `~/.config/gitswitch/history.json` | Per-repo identity usage counts and pins |
