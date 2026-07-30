---
title: Commands
description: Every gitswitch command and flag
---

`gitswitch --help` groups commands the same way this page does.

| | |
|---|---|
| **[Identity](#identity)** | `add` `switch` `list` `remove` `current` `init` |
| **[Session Isolation](#session-isolation)** | `pin` `unpin` `install` `uninstall` `doctor` |
| **[Channels](#channels)** | `version` `upgrade` `beta` `stable` |
| **[Extras](#extras)** | `login` `claude` `reauthor` `setup` `pacman` |

Plus [`gitswitch`](#gitswitch--the-interactive-ui) with no arguments, and [`gitswitch <nickname>`](#gitswitch-nickname--quick-switch) as a shortcut.

---

## `gitswitch` — the interactive UI

```bash
gitswitch
```

Three tabs — **Accounts**, **Utilities**, **Settings** — reachable with `1` `2` `3`, `tab` / `shift+tab`, or a mouse click.

| Key | |
|---|---|
| `↑` `↓` / `k` `j` | Move |
| `enter` | Switch to this account (globally) |
| `p` | Pin/unpin this account to the current repo — only shown inside a repo |
| `a` | Add an account |
| `e` | Edit — `ctrl+d` inside the form deletes |
| `v` | Toggle the second column between email and GitHub username |
| `c` | Cycle color theme (12 of them) |
| `u` | Upgrade — only appears when a newer version exists |
| `?` | CLI cheat sheet |
| `1` `2` `3` / `tab` | Switch tab |
| `q` / `ctrl+c` | Quit |

The mouse works throughout: hover moves focus, clicks select and toggle, the scroll wheel moves through lists.

**Utilities tab** — three toggles: Shell Integration, Session Isolation, HTTPS Credential Helper.
**Settings tab** — theme, config file location (opens it in `$EDITOR`), and the `gs` shell alias (rename with `e`, toggle with `enter`).

First launch with no accounts runs a short onboarding wizard, which scans your `gh` logins and `~/.ssh/` keys and offers to import them.

---

## `gitswitch <nickname>` — quick switch

```bash
gitswitch work
```

```
✓ Switched to 'work' — Alice Smith <alice@company.com>
```

Identical to `gitswitch switch work`. Exits non-zero if the nickname doesn't exist.

---

# Identity

## `add`

```bash
gitswitch add <nickname> <user-name> <email> [flags]
```

| Argument | |
|---|---|
| `nickname` | Short label — `work`, `oss`, `client-a`. Must be unique. Never written to git config. |
| `user-name` | Value for `user.name`. Quote it if it has spaces. |
| `email` | Value for `user.email`. |

| Flag | |
|---|---|
| `--ssh-key <path>` | Private key path. Sets `core.sshCommand` to `ssh -i <path> -o IdentitiesOnly=yes` on switch. |
| `--sign-key <key>` | GPG key ID **or** an SSH key. Sets `user.signingkey`, plus `gpg.format=ssh` for SSH keys. |
| `--gh-user <username>` | GitHub username. Runs `gh auth switch --user <username>` on switch (best-effort). |

```bash
gitswitch add personal "Alice Smith" alice@gmail.com

gitswitch add work "Alice Smith" alice@company.com \
  --ssh-key  ~/.ssh/id_work \
  --sign-key ABCD1234EF567890 \
  --gh-user  alice-corp
```

Prefer not typing all that? [`gitswitch login`](#login) fills it in from GitHub.

## `switch`

```bash
gitswitch switch <nickname>
```

Same as `gitswitch <nickname>`. See [Commit Identity](/docs/features/commit-identity) for exactly what gets written.

## `list`

```bash
gitswitch list [--json]
```

```
✓  personal        Alice Smith <alice@gmail.com>
   work            Alice Smith <alice@company.com>
```

`✓` marks your active global account.

## `current`

```bash
gitswitch current [--json | --short | --prompt]
```

```
work — Alice Smith <alice@company.com>  (pinned to this repo)
HTTPS credential helper: active
```

The suffix names the [scope](/docs/concepts/scopes): nothing for global, `(pinned to this repo)`, or `(this terminal's session)`. If a pin exists but Session Isolation is off, it says so.

Prints `No active profile` if nothing has been applied yet.

| Flag | |
|---|---|
| `--json` | `{nickname, user_name, email, scope, gh_user, ssh_key, credential_helper_active}` |
| `--short` | `nickname<TAB>email`, with the scope marker on the nickname (`work●`). For Starship. |
| `--prompt` | `nickname<TAB>color<TAB>marker` — theme color as an ANSI 256 index. For prompt functions. |

## `remove`

```bash
gitswitch remove <nickname>
```

Deletes the account. Does **not** revert git config that was applied while it was active, and doesn't revoke anything on GitHub. Any history entries under that nickname stay behind harmlessly.

## `init`

```bash
gitswitch init
```

Imports your current `--global` `user.name` / `user.email` as an account named `default`. Runs automatically on first use if you have no accounts.

---

# Session Isolation

## `pin`

```bash
gitswitch pin [nickname]
```

Writes the account's full identity into **this repo's** local git config — `user.name`, `user.email`, `user.signingkey`/`gpg.format`, `core.sshCommand`. Your global identity is untouched.

With no nickname, it adopts the repo's existing local `user.email` and matches it to a stored account.

Turns [Session Isolation](/docs/features/session-isolation) on automatically if it was off, since a pin can't take effect without it. Must be run inside a git repo.

## `unpin`

```bash
gitswitch unpin
```

Removes those keys from the repo's local config; the repo falls back to your global identity. Must be run inside a git repo.

## `install`

```bash
gitswitch install [flags]
```

Interactive wizard with three steps, each skippable:

1. **Shell integration** — prompt segment, nudge on `cd`, tab completion, `gs` alias
2. **HTTPS credential routing** — right token per repo on push
3. **Session Isolation** — right `gh` account per repo, and working pins

| Flag | |
|---|---|
| `--shell <zsh\|bash\|fish>` | Override shell detection. Also skips the wizard. |
| `--yes` / `-y` | Accept all defaults, no prompts. For scripts and CI. |
| `--https` | Register the HTTPS credential helper. Default `true`; prompted interactively when omitted. |

Safe to re-run — everything lives between `# gitswitch` markers and gets replaced in place, never duplicated. Re-run it after an upgrade to pick up an improved hook, or to repair credential-helper ordering.

## `uninstall`

```bash
gitswitch uninstall [--shell <zsh|bash|fish>]
```

Removes the shell integration block (or oh-my-zsh plugin), unregisters the HTTPS credential helper, and removes the `gh` wrapper — whichever were installed. Your accounts are left alone. Reload your shell to finish.

## `doctor`

```bash
gitswitch doctor [--json]
```

```
  ✓  git 2.50.1
  ✓  gh  2.95.0
  ✓  HTTPS pushes routed by gitswitch
  ✓  Session Isolation active (bare `gh` commands resolve per-repo)
```

The HTTPS check is the interesting one — git asks credential helpers in config order and takes the first answer, so another tool's helper registered ahead of gitswitch means gitswitch never gets asked:

```
  ✗  HTTPS pushes answered by another helper before gitswitch:
       credential.https://github.com.helper → !/opt/homebrew/bin/gh auth git-credential
       pushes may use the wrong account — run: gitswitch install
```

`gitswitch install` fixes the order without removing the other helper. See [HTTPS Push Routing](/docs/features/https).

---

# Channels

## `version`

```bash
gitswitch version
```

Prints your version and checks for a newer one (cached 24h).

## `upgrade`

```bash
gitswitch upgrade
```

Installs the latest release. On a Homebrew install it tells you to run `brew upgrade gitswitch` instead.

## `beta` / `stable`

```bash
gitswitch beta      # canary pre-releases
gitswitch stable    # back to normal
```

Both confirm first and never touch `~/.config/gitswitch/`. See [Versions & Release Channels](/docs/cli/channels).

---

# Extras

## `login`

```bash
gitswitch login [--profile <nickname>] [--host <hostname>] [--client-id <id>]
```

GitHub device flow in your browser. Creates an account with your name, email, and username filled in, and stores the token in your OS keychain. Re-running it on an existing account refreshes the token and keeps your SSH/signing keys. See [Connecting Accounts](/docs/features/accounts).

## `claude`

```bash
gitswitch claude [--scope user|project]
```

Installs the gitswitch skill into Claude Code, embedded in the binary. `user` → `~/.claude/skills/` (default), `project` → `.claude/skills/`.

## `reauthor`

```bash
gitswitch reauthor <base> --to <nickname> [--from <email>] [--push] [--yes]
```

Rewrites author and committer on commits between `<base>` and `HEAD`. `<base>` is a commit-ish or a bare number meaning "the last N commits".

| Flag | |
|---|---|
| `--to <nickname>` | Required. The account to attribute to. |
| `--from <email>` | Only rewrite commits currently authored by this email. |
| `--push` | Force-push (`--force-with-lease`) afterwards. |
| `--yes` / `-y` | No confirmations — for scripts and agents. |

This rewrites history. Full context in [AI Coding Agents](/docs/features/ai-agents).

## `setup`

```bash
gitswitch setup [--agent]
```

Checks `git` and `gh` and tells you what to do next. `--agent` emits a JSON manifest (version, account count, git/gh state) for AI agents.

## `pacman`

```bash
gitswitch pacman
```

Toggles [arcade mode](/docs/features/arcade) on or off for every future launch. No, we won't explain it here.

## `record` / `recommend`

Used by the shell hook. You rarely run these by hand.

```bash
gitswitch record    [--path <dir>]    # count the active account for this repo
gitswitch recommend [--path <dir>]    # print the recommended account, if any
```

`recommend` exits `0` and prints `nickname<TAB>name<TAB>email` when there's a recommendation, `1` silently when there isn't. See [Pins & Identity Awareness](/docs/features/identity-awareness).

---

## Exit codes

| | |
|---|---|
| `0` | Success |
| `1` | Something went wrong — account not found, git missing, not in a repo |

`gitswitch recommend` exiting `1` is normal and expected: it means "nothing to suggest".

## Where things are stored

| | |
|---|---|
| `~/.config/gitswitch/config.yaml` | Your accounts |
| `~/.config/gitswitch/config.json` | UI preferences — theme, alias, arcade high score |
| `~/.config/gitswitch/history.json` | Per-repo usage counts and pin records |
| OS keychain | Tokens from `gitswitch login` |

Full detail in [Configuration](/docs/cli/config).
