---
title: Commands
description: Every gitswitch command and flag
---
`gitswitch --help` groups commands the same way this page does.

| | |
|---|---|
| **[Identity](#identity)** | `login` `add` `switch` `list` `current` `remove` `init` |
| **[Setup & Session Isolation](#session-isolation)** | `shell` `doctor` `pin` `unpin` `guard` `xcode` `uninstall` |
| **[Updates & Channels](#channels)** | `version` `upgrade` `beta` `stable` |
| **[Extras](#extras)** | `reauthor` `skills` `pacman` `completion` `help` |

Plus [`gitswitch`](#gitswitch--the-interactive-ui) with no arguments, and [`gitswitch <nickname>`](#gitswitch-nickname--quick-switch) as a shortcut.

## Help, version, and flags that work everywhere

```bash
gitswitch --help              # every command, grouped
gitswitch <command> --help    # that command's flags and examples (-h works too)
gitswitch --version           # the installed version (-v and --v work too)
```

Flags mean the same thing on every command that has them:

| Flag | Meaning | On |
|---|---|---|
| `--json` | Machine-readable output | `list` `current` `version` `doctor` `guard status` `xcode status` |
| `-y`, `--yes` | Don't ask, just do it | `shell` `doctor` `reauthor` `beta` `stable` |
| `--dry-run` | Show what would happen, change nothing | `upgrade` `beta` `stable` `xcode sync` |
| `-q`, `--quiet` | Print nothing on success | `switch` |

Short names: `ls` = `list`, `rm` / `delete` = `remove`, `use` = `switch`, `whoami` = `current`, `update` = `upgrade`.

When a command is called wrongly — a missing argument, an unknown flag, a value that isn't allowed — gitswitch says what it expected, points you at that command's `--help`, and exits `2`. A mistyped command or nickname gets a "did you mean…?".

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
| `u` | Open the update prompt — when the status line shows a newer version |
| `?` | CLI cheat sheet |
| `1` `2` `3` / `tab` | Switch tab |
| `g` / `G` | Jump to the first / last account |
| `q` / `ctrl+c` | Quit |

The mouse works throughout: hover moves focus, clicks select and toggle, the scroll wheel moves through lists.

**Utilities tab** — five toggles: Shell Integration, Session Isolation, HTTPS Credential Helper, Auto-pin Repeat Repos, Pre-commit Guard.
**Settings tab** — theme (`←` `→`), config file location (opens it in `$EDITOR`), and the `gs` shell alias (rename with `e`, toggle with `enter`).

The line under each tab's content is the status line: it shows what gitswitch is doing (with a spinner), the result of the last action, or a pending update.

**Any window size.** The layout tightens as the terminal gets smaller — boxes become rows, blank lines go, long lists scroll — rather than running off the screen. The smallest window it draws in is 44 columns by 12 rows; below that it shows a short notice (and `gitswitch <nickname>` still works). Set `GITSWITCH_REDUCE_MOTION=1` to turn off the decorative motion (the sliding tab marker, the launch shimmer, the staged setup scan); progress spinners stay.

First launch on a machine with no profiles yet checks for `git`/`gh` and offers to install whichever's missing, then runs a short setup wizard: it scans your git config, `gh` logins and `~/.ssh/` keys, lets you pick what to import, and offers to add more accounts.

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

Same as `gitswitch <nickname>`. See [Commit Identity](/docs/accounts/commit-identity) for exactly what gets written.

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

The suffix names the [scope](/docs/accounts/scopes): nothing for global, `(pinned to this repo)`, or `(this terminal's session)`. If a pin exists but Session Isolation is off, it says so.

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

Turns [Session Isolation](/docs/routing/session-isolation) on automatically if it was off, since a pin can't take effect without it. Must be run inside a git repo.

## `unpin`

```bash
gitswitch unpin
```

Removes those keys from the repo's local config; the repo falls back to your global identity. Must be run inside a git repo.

## `guard`

Refuses a commit whose author doesn't match the right gitswitch profile *before* it enters history. Installs a global pre-commit hook (`git config --global core.hooksPath`) that chains into any existing global hooks path and each repo's own hooks, so Husky, lefthook and `pre-commit` keep working.

```bash
gitswitch guard enable    # install the hook
gitswitch guard disable   # remove it (restores a previous core.hooksPath)
gitswitch guard status    # hook state + whether this repo would pass (--json supported)
gitswitch guard check     # what the hook runs
```

A commit is blocked when there is no `user.email`, the email matches no profile, the repo is pinned to a different profile, or (for unpinned repos) gitswitch has learned a different profile for it. With no profiles yet, everything is allowed.

## `xcode` (early)

Keeps Xcode's own GitHub accounts (**Xcode → Settings → Accounts**) in step with gitswitch. Xcode stores tokens in a keychain group only Apple's processes can write, so each account is still added to Xcode once — gitswitch makes that one paste and then manages which account Xcode *picks*. Commit authorship needs nothing: Xcode already reads git config.

```bash
gitswitch xcode status             # Xcode's accounts vs your profiles (--json supported)
gitswitch xcode add [nickname]     # copy the profile's gh token to the clipboard (--pat: open GitHub's token page)
gitswitch xcode sync [nickname]    # make that profile's account Xcode's default (--dry-run to preview)
gitswitch xcode bind [nickname]    # embed the username in this repo's https remotes
gitswitch xcode unbind
gitswitch xcode enable|disable     # run sync on every switch, bind/unbind on pin/unpin
```

Xcode matches a remote to an account by a username in the URL, then the repo owner, and otherwise uses the first account listed. `sync` reorders that list; `bind` pins a single repo regardless of order. SSH remotes are never rewritten. Xcode reads its accounts at launch, so restart it after a change.

## `shell`

```bash
gitswitch shell [flags]
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
       pushes may use the wrong account — run: gitswitch shell
```

`gitswitch shell` fixes the order without removing the other helper. See [HTTPS Push Routing](/docs/routing/https).

---

# Channels

## `version`

```bash
gitswitch version [--json]
```

Prints your version, your channel, how gitswitch is managed (Homebrew or the install script), and whether a newer release exists (cached 24h). For just the number, with no network access, use `gitswitch --version`.

## `upgrade`

```bash
gitswitch upgrade [--check] [--dry-run]
```

Updates to the latest release on your channel. On a Homebrew install it runs `brew update && brew upgrade` for you, so Homebrew stays in charge; on a script install it downloads the new binary. `--check` only reports; `--dry-run` shows the command without running it.

## `beta` / `stable`

```bash
gitswitch beta   [--yes] [--dry-run]    # pre-release builds
gitswitch stable [--yes] [--dry-run]    # back to normal
```

Both show their plan and confirm first, and never touch `~/.config/gitswitch/`. On Homebrew they swap the `gitswitch` and `gitswitch-beta` formulae. See [Versions & Release Channels](/docs/cli/channels).

---

# Extras

## `login`

```bash
gitswitch login [--profile <nickname>] [--host <hostname>]
```

Delegates to `gh auth login` (the real GitHub CLI's own login command) to authenticate with GitHub, then reads the resulting account back via `gh`. Creates an account with your name, email, and username filled in — no custom OAuth app or separate registration step needed since `gh` owns the account the moment it authenticates. Re-running it on an existing account refreshes your login and keeps your SSH/signing keys. See [Connecting Accounts](/docs/accounts/accounts).

## `skills`

```bash
gitswitch skills [--scope user|project] [--offline] [--check] [--update]
```

Installs or updates the gitswitch skill, embedded in the binary. Tries [skills.sh](https://www.skills.sh) first (works with Claude Code, Cursor, Codex, and more); falls back to installing directly for Claude Code and the shared `.agents/skills` convention when `--offline` is passed or `npx` isn't available. `user` → `~/.claude/skills/`, `~/.agents/skills/` (default), `project` → `.claude/skills/`, `.agents/skills/`.

Use `--check` to verify whether installed skill files match the current binary, or `--update` to rewrite all installed skill files. `gitswitch upgrade` updates installed skills automatically.

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

This rewrites history. Full context in [AI Coding Agents](/docs/ai/ai-agents).

## `pacman`

```bash
gitswitch pacman
```

Toggles [arcade mode](/docs/extras/arcade) on or off for every future launch. No, we won't explain it here.

## `completion`

```bash
gitswitch completion zsh    # or bash, fish, powershell
```

Prints a tab-completion script for your shell. [`gitswitch shell`](#shell) already loads it for you; use this only if you manage completions yourself. It completes commands, flags, flag values, and your own account nicknames.

## `record` / `recommend`

Plumbing for the shell hook — hidden from `--help` and tab completion, but still callable.

```bash
gitswitch record    [--path <dir>]    # count the active account for this repo
gitswitch recommend [--path <dir>]    # print the recommended account, if any
```

`recommend` exits `0` and prints `nickname<TAB>name<TAB>email` when there's a recommendation, `1` silently when there isn't. See [Pins & Identity Awareness](/docs/routing/identity-awareness).

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

Full detail in [Configuration](/docs/cli/config).
