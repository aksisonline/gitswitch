---
name: gitswitch
description: >
  Switch git identity and GitHub account using the gitswitch tool when the user is
  committing or pushing as the wrong account. Use this skill whenever the user mentions:
  their commits showing the wrong name/email, pushing to the wrong GitHub account,
  needing to switch git profiles, wanting to check which git identity is active,
  getting attribution errors on commits, working across multiple GitHub accounts
  (personal vs work, multiple clients), setting up or troubleshooting Xcode git accounts,
  enforcing commit identities with a pre-commit guard, or any time the phrase "wrong account"
  or "wrong identity" appears in a git/GitHub context. Don't wait for the user to say
  "gitswitch" — if the problem smells like a git identity mismatch, invoke this skill.
---

# gitswitch skill

`gitswitch` is a terminal tool (installed via Homebrew or curl) that manages multiple local git identities. It switches `git config user.name`, `git config user.email`, the SSH key, and (optionally) the `gh` CLI account — all at once, with one command.

## When to use it

Reach for `gitswitch` whenever:
- A commit landed with the wrong author name or email
- The user is about to push to a work/client repo but is still on their personal identity
- The user asks "which git account am I on?"
- The user needs to add a new profile for a new job, client, or side project
- The user is confused about why `gh` and `git` show different accounts (they're independent)
- The user wants to prevent commits from ever landing with the wrong author (`gitswitch guard`)
- Xcode or IDE source control is pushing or committing under the wrong GitHub account

## Step 1 — Diagnose

Always check what's currently active before doing anything:

```bash
gitswitch current     # shows active profile name + email
gitswitch list        # shows all profiles, ✓ marks the active one
gitswitch doctor      # checks git, gh, credential helper, and guard status
```

If `gitswitch` is not installed, say so clearly and offer the install instructions:
```bash
# Homebrew (recommended):
brew install aksisonline/tap/gitswitch

# Or, without Homebrew (curl one-liner):
curl -fsSL https://get.gitswitch.dev | bash
```

## Step 2 — Switch to the right profile

If the right profile already exists:

```bash
gitswitch <nickname>
# e.g. gitswitch work
```

That's it. One command applies `user.name`, `user.email`, the SSH key, runs `gh auth switch` if a GitHub username is stored in the profile, and updates Xcode preferences if Xcode auto-sync is enabled.

If the user needs to pick from a list, show them `gitswitch list` output and ask which one.

## Step 3 — Create a profile if it doesn't exist yet

```bash
gitswitch add <nickname> "<Full Name>" <email> [flags]

# Common flags:
#   --ssh-key ~/.ssh/id_work       force a specific SSH key
#   --gh-user <github-handle>      also switch gh CLI account
#   --sign-key <key>               signed commits: GPG key ID, or an SSH
#                                  key path (~/.ssh/id_ed25519.pub) — gitswitch
#                                  sets gpg.format=ssh automatically

# Examples:
gitswitch add work "Alice Smith" alice@company.com --ssh-key ~/.ssh/id_work --gh-user alice-corp
gitswitch add personal "Alice Smith" alice@gmail.com --ssh-key ~/.ssh/id_personal --gh-user alice
gitswitch add clienta "Alice Smith" alice@clienta.com --ssh-key ~/.ssh/id_clienta
```

After adding, switch immediately:
```bash
gitswitch work
```

## Step 4 — Verify

```bash
gitswitch current     # confirm the new identity is active
git config --global user.email   # double-check git sees it
```

If the user already made commits with the wrong identity, don't hand-run `git rebase`/`git commit --amend` yourself — use `gitswitch reauthor` (below). It's one command instead of a multi-step git dance, and it's the tool the user expects to own this.

---

## Pre-commit identity guard — `gitswitch guard`

To prevent wrong-author commits from *ever* entering git history, enable the global pre-commit guard:

```bash
gitswitch guard enable    # install global pre-commit enforcement hook
gitswitch guard disable   # uninstall hook and restore previous core.hooksPath
gitswitch guard status    # inspect hook status and test current repo
```

### How the guard works
- Configures `git config --global core.hooksPath ~/.config/gitswitch/hooks`.
- Chains seamlessly into any pre-existing global hooks and each repo's `.git/hooks/pre-commit` (plus passthrough shims for all client hooks so Husky, lefthook, and `pre-commit` keep working).
- **Refuses commits before they are created** when:
  1. No `user.email` is configured in git.
  2. The effective commit email matches no stored gitswitch profile.
  3. The repository is pinned (`gitswitch pin`) to a profile whose email differs from the commit email.
  4. The repository is unpinned (`ScopeGlobal`) and learned repo history (`history.Recommend`) expects a different profile.
- If zero profiles exist (fresh install), commits are permitted so onboarding is never blocked.

---

## Xcode integration — `gitswitch xcode`

Xcode manages its own GitHub accounts in **Xcode → Settings → Accounts** and stores tokens in a private, Apple-only Keychain access group. While third-party tools cannot write tokens directly, `gitswitch xcode` solves account mismatch and synchronization:

```bash
gitswitch xcode status             # see Xcode accounts vs gitswitch profiles
gitswitch xcode add [nickname]     # copies profile's gh token to clipboard for 1-paste setup
                                   # (--pat opens GitHub's pre-filled token page instead)
gitswitch xcode sync [nickname]    # promotes profile to index 0 (Xcode's default account)
gitswitch xcode bind [nickname]    # binds repo's https remotes (https://user@github.com/...)
gitswitch xcode unbind             # removes embedded username from https remotes
gitswitch xcode enable             # auto-sync on 'switch' and auto-bind on 'pin'
gitswitch xcode disable            # turn off automatic syncing
```

### Why this matters for Xcode
- **Default Account Bias:** Xcode uses the **first** account in its accounts list for organization repos and Swift Package Manager (SPM) dependency resolutions. `gitswitch xcode sync` (or auto-sync on `switch`) ensures Xcode's default always matches your active profile.
- **Per-Project Pinning:** `gitswitch xcode bind` embeds `<gh-user>@github.com` into the repository's HTTPS remote. Xcode's `IDESourceControlAccountManager` matches the username in the URL first, ensuring that project always pushes and fetches using the pinned account.
- **One-Paste Onboarding:** `gitswitch xcode add <nickname>` extracts the token from `gh auth token` and puts it on your clipboard so adding an account in Xcode Settings is simply: `+` → `GitHub` → Account name → ⌘V.

*(Note: Xcode reads its accounts on application launch, so restart Xcode after syncing accounts).*

---

## Editor & IDE extensions

gitswitch provides lightweight status bar widgets and identity pickers for modern editors:

- **VS Code Extension** (`extensions/vscode`):
  - Displays current profile in the status bar (`●` when repo is pinned, `[session]` when overridden).
  - Status bar tooltip reports pre-commit guard status.
  - Commands: `Switch Profile`, `Pin Profile to This Repo`, `Unpin This Repo`, `Toggle Pre-commit Guard`.
- **JetBrains Plugin** (`extensions/jetbrains`):
  - Status bar widget for IntelliJ IDEA, PyCharm, WebStorm, Android Studio, etc.
  - Tools menu actions: `Switch Profile`, `Pin Profile to This Repo`, `Unpin This Repo`, `Toggle Pre-commit Guard`.

---

## Fixing already-made commits — `gitswitch reauthor`

When commits were authored under the wrong profile (before the user switched, or before a profile existed), rewrite them in one step instead of scripting `git rebase`/`git commit --amend` yourself:

```bash
gitswitch reauthor <base> --to <nickname> [--from <old-email>] [--push] [--yes]
```

- `<base>` — a commit-ish (`HEAD~3`, a SHA) or a bare number N meaning "the last N commits". Range rewritten is `(base, HEAD]`.
- `--to <nickname>` — required; the profile whose name/email gets applied.
- `--from <old-email>` — optional filter; only commits currently authored by this email are rewritten, others in range are left alone. Use this when the range mixes commits from multiple identities.
- `--push` — force-pushes (`--force-with-lease`) after rewriting. Without it, the rewrite is local only.
- `--yes` — skips the interactive y/N prompts. **Ask the user before passing `--yes`, especially with `--push`** — rewriting pushed history is destructive to anyone else who has pulled the branch. Confirm scope (which commits, which identity) in plain language first, then run it.

Example — the two most recent commits were made as the old identity before switching to `work`:
```bash
gitswitch reauthor 2 --to work --push
```

Requires a clean working tree (it rebases). If a merge conflict interrupts the rebase, gitswitch leaves the repo mid-rebase and tells you to run `git rebase --abort` to bail out — do that and report back to the user rather than trying to resolve conflicts on their behalf.

---

## Identity awareness — shell integration

`gitswitch shell` sets up shell integration, HTTPS routing, and guard hooks:

```bash
gitswitch shell
```

- **Prompt segment** — shows current git identity in your shell prompt
- **Identity nudge** — when you `cd` into a repo, suggests the identity you usually use there
- **Tab completion** — completes commands and profile nicknames
- **Pre-commit guard** — optional prompt to turn on wrong-author commit enforcement

Detects your prompt framework automatically (Starship, oh-my-zsh, Powerlevel10k, or raw zsh/bash/fish) and installs the right way for each. Idempotent — safe to run multiple times.

After running, reload your shell:
```bash
source ~/.zshrc   # or open a new terminal
```

---

## Pinning an identity to a repo

When one repo should *always* use a given identity, pin it instead of switching globally:

```bash
gitswitch pin work      # write 'work' into this repo's local git config
gitswitch pin           # adopt the identity the repo already has configured
gitswitch unpin         # remove it, fall back to the global identity
```

A pin is a **local** identity: gitswitch writes `user.name`, `user.email`,
`user.signingkey`/`gpg.format`, and `core.sshCommand` to the repo's own
`.git/config`, so commits here are correct no matter which profile is active
globally. It also switches the gh CLI to the matching account, and the HTTPS
credential helper serves this repo's pushes from it. If Xcode auto-sync is enabled,
it binds the repo's remotes for Xcode as well.

Prefer this over `gitswitch switch` whenever the user's problem is "this one repo
keeps committing as the wrong person" — it fixes that repo permanently without
disturbing the others.

If a repo was configured by hand before gitswitch (`git config --local
user.email …`), bare `gitswitch pin` adopts that identity: it matches the email
to a stored profile and fills in the rest.

`gitswitch current` names the scope, so you can tell these apart without reading
git config yourself:

```
work — Alice W <alice@work.com>  (pinned to this repo)     # repo-local
work — Alice W <alice@work.com>  (this terminal's session) # GIT_CONFIG_* env vars
personal — Alice <alice@personal.dev>                      # global identity
```

Check this in Step 1 before switching: if the repo is already pinned, a global
`switch` will not change what commits here use, and switching is the wrong fix.

---

## When HTTPS pushes use the wrong account

If commits are attributed correctly but a **push** is rejected with something like
`Permission to owner/repo denied to <other-account>`, the identity is fine and the
*credential* is wrong. Diagnose with:

```bash
gitswitch doctor
```

Git asks credential helpers in config order and takes the first answer, so another
tool's helper (usually `gh auth setup-git`'s) can be asked before gitswitch. doctor
names it; `gitswitch shell` repairs the order without removing the other helper.
Re-running `gh auth setup-git` or an interactive `gh auth login` undoes the repair,
so if it recurs, that is why.

Do not "fix" this by running `gh auth switch` unless the user asks — that changes
their active GitHub account machine-wide, which may be in use by another terminal
or agent. To push once with the right token without touching global state:

```bash
git -c credential.helper= -c credential.helper='!gitswitch credential' push
```

---

## All commands

| Command | What it does |
|---------|-------------|
| `gitswitch` | Open interactive TUI |
| `gitswitch <name>` | Quick switch (no UI) |
| `gitswitch list` | List all profiles |
| `gitswitch current` | Show active profile |
| `gitswitch current --short` | Output `nickname\temail` (used by Starship prompt) |
| `gitswitch add …` | Add a new profile |
| `gitswitch remove <name>` | Delete a profile |
| `gitswitch pin [name]` | Pin an identity to this repo's local git config (no name = adopt existing) |
| `gitswitch unpin` | Remove the repo's local identity, fall back to global |
| `gitswitch guard [enable\|disable\|status\|check]` | Manage pre-commit identity guard hook |
| `gitswitch xcode [status\|add\|sync\|bind\|unbind\|enable\|disable]` | Sync & bind identities with Xcode's source control accounts |
| `gitswitch reauthor <base> --to <name>` | Rewrite author identity on already-made commits |
| `gitswitch record` | Log current identity for this repo (called by shell hooks) |
| `gitswitch recommend` | Print recommended identity if threshold met |
| `gitswitch doctor` | Verify git, gh, credential helper, and guard installation |
| `gitswitch shell` | Set up shell prompt segment, nudge hook, and tab completion |
| `gitswitch skills` | Install this skill for your AI agent (skills.sh first, offline fallback) |
| `gitswitch init` | Re-import current `git config` as `default` |
| `gitswitch version` | Show version, check for updates |
| `gitswitch upgrade` | Upgrade to latest release |

Profiles stored at `~/.config/gitswitch/config.yaml` (or migrated from `profiles.json`).
Usage history stored at `~/.config/gitswitch/history.json`.
Xcode auto-sync marker at `~/.config/gitswitch/xcode-sync`.
Guard hook stored at `~/.config/gitswitch/hooks/pre-commit`.
