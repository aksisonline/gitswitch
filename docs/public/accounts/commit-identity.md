---
title: Commit Identity
description: What a switch actually writes, and how git decides who you are
---

## Switching

```bash
gitswitch work          # or: gitswitch switch work, or enter in the TUI
```

```
✓ Switched to 'work' — Alice Smith <alice@company.com>
```

That writes, in order:

```bash
git config --global user.name    "Alice Smith"
git config --global user.email   "alice@company.com"
git config --global user.signingkey "<key>"    # if the account has one
git config --global gpg.format   ssh           # only when that key is an SSH key
git config --global core.sshCommand "ssh -i <path> -o IdentitiesOnly=yes"   # if it has an SSH key
gh auth switch --user alice-corp               # if it has a GitHub username
```

Two things worth knowing:

- **Accounts without a key clear that setting.** Switching to an account with no signing key unsets `user.signingkey` *and* `gpg.format`; no SSH key unsets `core.sshCommand`. So one account's setup never leaks into the next one.
- **The `gh` step is best-effort.** If `gh` isn't installed or that account isn't logged in, you get a warning and the git config still switches correctly:

```
warning: gh auth switch --user alice-corp: exec: "gh": executable file not found in $PATH
✓ Switched to 'work' — Alice Smith <alice@company.com>
```

## Who am I right now?

```bash
gitswitch current
```

```
work — Alice Smith <alice@company.com>  (pinned to this repo)
HTTPS credential helper: active
```

The parenthetical is the important part — see [Scopes](/docs/accounts/scopes). For scripts:

```bash
gitswitch current --json
```

Or ask git directly:

```bash
git config --global user.name
git config --global user.email
```

## How git decides

Narrowest scope wins:

1. `GIT_CONFIG_*` environment variables — this terminal only
2. The repo's own `.git/config` — where [pins](/docs/routing/identity-awareness#pin-a-repo) live
3. `~/.gitconfig` — where a plain `gitswitch <name>` writes
4. `/etc/gitconfig` — system-wide

So a pinned repo ignores your global identity entirely. That's the point of pinning.

## Check what a commit says

```bash
git log -1 --format="%an <%ae>"
```

## Fixing a wrong one

Not pushed yet:

```bash
gitswitch work
git commit --amend --reset-author --no-edit
```

More than one commit, or already pushed:

```bash
gitswitch reauthor 3 --to work --push
```

One command instead of a hand-rolled interactive rebase. Full reference in [AI Coding Agents](/docs/ai/ai-agents) — it's the same tool whether a human or an agent made the mess.

> Once pushed, attribution is permanent unless you rewrite history. Cheaper to prevent: [pin the repo](/docs/routing/identity-awareness#pin-a-repo), or let the [shell nudge](/docs/routing/shell#the-nudge) catch you on the way in.

## Removing a per-repo identity

If a repo has its own identity you no longer want:

```bash
gitswitch unpin
```

That clears `user.name`, `user.email`, `user.signingkey`, `gpg.format`, and `core.sshCommand` from the repo's local config, and it falls back to your global identity. (The by-hand equivalent is `git config --local --unset user.email`, five times over.)

## Next

- **[Scopes](/docs/accounts/scopes)** — global vs. repo vs. terminal
- **[SSH Keys](/docs/accounts/ssh-keys)** — forcing the right key
- **[Commit Signing](/docs/accounts/gpg)** — GPG or SSH signing per account
