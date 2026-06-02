---
title: Commit Identity Switching
description: How gitswitch switches git commit attribution
---

This page covers what switching does to your git config and how git resolves identity when both local and global config exist.

## What switching sets

Every profile switch — whether via `gitswitch <nickname>`, `gitswitch switch <nickname>`, or the TUI — runs these commands in order:

```bash
git config --global user.name "<name>"
git config --global user.email "<email>"
git config --global user.signingkey "<key>"     # only if profile has sign-key
git config --global core.sshCommand "ssh -i <path> -o IdentitiesOnly=yes"  # only if profile has ssh-key
```

These are `--global` writes to `~/.gitconfig`. They apply to every repo on your machine until you switch again.

## Quick switch

```bash
gitswitch work
```

```
✓ Switched to 'work' — Alice Smith <alice@company.com>
```

## Check what's active

```bash
gitswitch current
```

```
work — Alice Smith <alice@company.com>
```

Or inspect git config directly:

```bash
git config --global user.name
git config --global user.email
```

## Git's identity resolution order

Git checks identity in this order, stopping at the first match:

1. Repo-local config (`git config --local user.email`) — set per-repo, never touched by gitswitch
2. Global config (`~/.gitconfig`) — gitswitch writes here
3. System config (`/etc/gitconfig`)

If a repo has a local override, gitswitch's global switch has no effect on that repo.

## Verify commits are attributed correctly

Check the author of a recent commit:

```bash
git log -1 --format="%an <%ae>"
```

Before pushing, if you suspect the wrong identity was used:

```bash
# Check what identity was active
gitswitch current

# If wrong, switch and amend the last commit before pushing
gitswitch work
git commit --amend --reset-author --no-edit
```

> **Tip:** Once pushed, commit attribution is permanent. Use [shell integration](/docs/features/shell) to get nudged when entering repos that usually use a different identity.

## Remove a local override

If commits in a repo are using a per-repo identity you set manually:

```bash
git config --local --unset user.email
git config --local --unset user.name
```

After this, the global config (managed by gitswitch) takes effect.

## Next steps

- [SSH Keys](/docs/features/ssh-keys)
- [GPG Signing](/docs/features/gpg)
- [Identity Awareness](/docs/features/identity-awareness)
