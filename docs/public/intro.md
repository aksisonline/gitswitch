---
title: Introduction
description: What gitswitch is and what it does
---

gitswitch manages multiple git identities — name, email, SSH key, GPG signing key, and GitHub CLI account — and switches between them with a single command.

## The problem

Every commit is attributed to whoever `git config --global user.name` and `git config --global user.email` point to. If you work across personal, work, and client repos on the same machine, you're constantly at risk of attributing commits to the wrong identity. Wrong attribution is permanent after a push.

SSH keys compound the problem: your agent may silently authenticate with whichever key works, not necessarily the one you intended.

## What gitswitch does

```bash
gitswitch              # Open interactive TUI to browse and switch profiles
gitswitch work         # Immediately switch to the 'work' profile
gitswitch install      # Set up shell prompt, identity nudges, and tab completion
```

When you switch to a profile, gitswitch runs in this order:

1. `git config --global user.name "<name>"`
2. `git config --global user.email "<email>"`
3. `git config --global user.signingkey "<key>"` — if the profile has a GPG key
4. `git config --global core.sshCommand "ssh -i <path> -o IdentitiesOnly=yes"` — if the profile has an SSH key
5. `gh auth switch --user <username>` — if the profile has a GitHub username (best-effort)

## How it compares to `gh auth switch`

`gh auth switch` changes which GitHub API token the `gh` CLI uses. It does not touch `user.name`, `user.email`, SSH keys, or GPG keys.

| | `gh auth switch` | `gitswitch` |
|---|---|---|
| Changes commit author identity | No | Yes |
| Switches SSH keys | No | Yes |
| Switches GPG signing keys | No | Yes |
| Switches GitHub CLI account | Yes | Yes (optional, via `--gh-user`) |
| Works with GitLab / Bitbucket | Yes | Yes |
| Interactive TUI | No | Yes |

## Core features

- **Commit identity** — switches `user.name` and `user.email` globally
- **SSH key routing** — forces `IdentitiesOnly=yes` so no accidental agent fallback
- **GPG signing** — switches `user.signingkey` alongside your identity
- **GitHub CLI sync** — optionally calls `gh auth switch` on every profile switch
- **Identity awareness** — learns which identity you use per-repo and nudges you when something looks wrong
- **Shell integration** — prompt segments, automatic nudges on `cd`, tab completion
- **12 color themes** — Default, Ocean, Sunset, Forest, Mono, Rose, Arctic, Gold, Violet, Ember, Matrix, Steel

## Storage

| File | Contents |
|------|----------|
| `~/.config/gitswitch/profiles.json` | All profile definitions |
| `~/.config/gitswitch/config.json` | UI preferences (color theme) |
| `~/.config/gitswitch/history.json` | Per-repo identity usage counts and pins |

## Next steps

- [Install gitswitch](/docs/installation)
- [Quick Start](/docs/quick-start)
- [CLI Reference](/docs/cli/commands)
