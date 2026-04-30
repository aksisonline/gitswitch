---
name: git-switcher
description: >
  Switch git identity and GitHub account using the gitswitch tool when the user is
  committing or pushing as the wrong account. Use this skill whenever the user mentions:
  their commits showing the wrong name/email, pushing to the wrong GitHub account,
  needing to switch git profiles, wanting to check which git identity is active,
  getting attribution errors on commits, working across multiple GitHub accounts
  (personal vs work, multiple clients), or any time the phrase "wrong account" or
  "wrong identity" appears in a git/GitHub context. Don't wait for the user to say
  "gitswitch" — if the problem smells like a git identity mismatch, invoke this skill.
---

# git-switcher skill

`gitswitch` is a terminal tool (installed via Homebrew or curl) that manages multiple local git identities. It switches `git config user.name`, `git config user.email`, the SSH key, and (optionally) the `gh` CLI account — all at once, with one command.

## When to use it

Reach for `gitswitch` whenever:
- A commit landed with the wrong author name or email
- The user is about to push to a work/client repo but is still on their personal identity
- The user asks "which git account am I on?"
- The user needs to add a new profile for a new job, client, or side project
- The user is confused about why `gh` and `git` show different accounts (they're independent)

## Step 1 — Diagnose

Always check what's currently active before doing anything:

```bash
gitswitch current     # shows active profile name + email
gitswitch list        # shows all profiles, ✓ marks the active one
```

If `gitswitch` is not installed, say so clearly and offer the install instructions:
```bash
# Homebrew (recommended):
brew tap aksisonline/gitswitch
brew install gitswitch

# Or, without Homebrew (curl one-liner):
curl -fsSL https://raw.githubusercontent.com/aksisonline/gitswitch/main/.github/install.sh | bash
```

## Step 2 — Switch to the right profile

If the right profile already exists:

```bash
gitswitch <nickname>
# e.g. gitswitch work
```

That's it. One command applies `user.name`, `user.email`, the SSH key, and runs `gh auth switch` if a GitHub username is stored in the profile.

If the user needs to pick from a list, show them `gitswitch list` output and ask which one.

## Step 3 — Create a profile if it doesn't exist yet

```bash
gitswitch add <nickname> "<Full Name>" <email> [flags]

# Common flags:
#   --ssh-key ~/.ssh/id_work       force a specific SSH key
#   --gh-user <github-handle>      also switch gh CLI account
#   --sign-key <GPG-key-ID>        for signed commits

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

If the user had already made commits with the wrong identity on a branch, let them know that `gitswitch` fixes *future* commits — amending history is a separate step (and usually unnecessary unless they haven't pushed yet).

---

## Identity awareness — shell integration

`gitswitch install` sets up three things at once:

```bash
gitswitch install
```

- **Prompt segment** — shows current git identity in your shell prompt
- **Identity nudge** — when you `cd` into a repo, suggests the identity you usually use there
- **Tab completion** — completes commands and profile nicknames

Detects your prompt framework automatically (Starship, oh-my-zsh, Powerlevel10k, or raw zsh/bash/fish) and installs the right way for each. Idempotent — safe to run multiple times.

After running, reload your shell:
```bash
source ~/.zshrc   # or open a new terminal
```

---

## Pinning an identity to a repo

If you always want a specific identity recommended for a repo regardless of usage history:

```bash
gitswitch pin work      # always recommend 'work' in this repo
gitswitch unpin         # remove the pin, fall back to auto-recommendation
```

The pin takes permanent priority over the auto-learned counts. `unpin` reverts to the usage-based recommender.

---

## How the auto-recommender works

`gitswitch` tracks which identity is active each time you enter a repo (via the shell cd hook). When the top identity has **≥3 entries** and **≥60% share** and differs from your current identity, it nudges:

```
gitswitch: this repo usually uses work <alice@company.com> — switch? [y/N]
```

One keypress. Defaults to N if you ignore it.

To manually record the current identity for the current repo:
```bash
gitswitch record
```

To check what would be recommended right now:
```bash
gitswitch recommend     # exits 0 + prints nickname/name/email if nudge warranted
                        # exits 1 silently if already on the right identity or no data
```

---

## Common scenarios

### "My commits on GitHub show my personal email at my work repo"
```bash
gitswitch list           # find the work profile name
gitswitch work           # switch to it
```
Then for future commits they'll be fine. If they need to fix the last commit:
```bash
git commit --amend --reset-author --no-edit
```

### "I have two GitHub accounts and keep using the wrong one"
```bash
gitswitch add personal "Name" personal@email.com --ssh-key ~/.ssh/id_personal --gh-user my-handle
gitswitch add work     "Name" work@email.com     --ssh-key ~/.ssh/id_work     --gh-user my-corp-handle
# then before each session:
gitswitch personal   # or: gitswitch work
```

### "I want gitswitch to always remind me to use the right account in a repo"
```bash
gitswitch install        # sets up the shell nudge hook
gitswitch pin work       # optional: force 'work' regardless of history
```

### "Is this tool the same as `gh auth switch`?"
No — they solve different problems:
- `gh auth switch` changes which account the **`gh` CLI** (API tokens) uses
- `gitswitch` changes your **commit identity** (`user.name` / `user.email`) and SSH key

You often need both. `gitswitch` calls `gh auth switch` automatically if you stored a `--gh-user` in the profile.

### "I want to see my profiles in a UI"
```bash
gitswitch     # opens the interactive TUI, use ↑↓ to navigate, Enter to switch
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
| `gitswitch pin <name>` | Pin an identity to the current repo permanently |
| `gitswitch unpin` | Remove pin, fall back to usage-based recommendation |
| `gitswitch record` | Log current identity for this repo (called by shell hooks) |
| `gitswitch recommend` | Print recommended identity if threshold met |
| `gitswitch install` | Set up shell prompt segment, nudge hook, and tab completion |
| `gitswitch claude` | Install this skill into Claude Code (~/.claude/skills) |
| `gitswitch claude --scope project` | Install skill for this project only (.claude/skills) |
| `gitswitch init` | Re-import current `git config` as `default` |
| `gitswitch version` | Show version, check for updates |
| `gitswitch upgrade` | Upgrade to latest release |

Profiles stored at `~/.config/gitswitch/profiles.json`.
Usage history stored at `~/.config/gitswitch/history.json`.
