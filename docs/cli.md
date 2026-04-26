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

Same as `gitswitch <nickname>` (the positional quick-switch form). Both forms apply the same changes in the same order. Use whichever fits your workflow.

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
gitswitch current
```

Prints the nickname, name, and email of the currently active profile.

```
work — Alice Smith <alice@company.com>
```

Returns `No active profile` if none has been applied yet.

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

This is also run automatically on first launch if no profiles exist yet.

---

## `gitswitch version` — show version info

```bash
gitswitch version
```

Prints the installed version and checks for a newer release (via a 24-hour cache — no network call if recently checked).

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

Downloads and runs the install script for the latest release. Always fetches from GitHub — bypasses the version cache. Safe to run at any time; exits early if already on the latest version.

```
Checking for updates...
Upgrading v0.1.11 → v0.1.12...
✓ Upgrade complete. Restart gitswitch to use the new version.
```

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

## Profile storage

Profiles are stored as JSON at:

```
~/.config/gitswitch/profiles.json
```
