---
title: GitHub Account Sync
description: How gitswitch integrates with the GitHub CLI
---

This page covers the `--gh-user` flag — what it does when you switch profiles and how it behaves when `gh` is unavailable.

## What it does

When you switch to a profile that has a GitHub username, gitswitch runs:

```bash
gh auth switch --user alice-work
```

This changes which GitHub account `gh` CLI commands operate under — creating PRs, listing issues, fetching notifications, and so on.

## Add a profile with a GitHub username

```bash
gitswitch add work "Alice Smith" alice@company.com \
  --gh-user alice-work
```

The value passed to `--gh-user` must match a username that is already authenticated in `gh`:

```bash
gh auth status
# Shows all authenticated accounts
```

If the account is not yet authenticated:

```bash
gh auth login
# Follow the prompts to authenticate the account
```

## Behavior when `gh` is unavailable

The GitHub switch is best-effort. If `gh` is not installed, or the specified account is not authenticated, gitswitch prints a warning but still completes the profile switch:

```
warning: gh auth switch --user alice-work: exec: "gh": executable file not found in $PATH
✓ Switched to 'work' — Alice Smith <alice@company.com>
```

Steps 1–4 (git config writes) always execute. Step 5 (`gh auth switch`) is advisory only.

## Profiles without GitHub sync

If you don't pass `--gh-user`, switching that profile leaves the active `gh` account unchanged:

```bash
gitswitch add internal "Alice Smith" alice@internal.corp
# No --gh-user — gh account is not touched when switching to 'internal'
```

## Using gh CLI with multiple accounts

```bash
# Switch to personal profile (includes --gh-user alice)
gitswitch personal
gh repo list
# Lists repos owned by alice

# Switch to work profile (includes --gh-user alice-work)
gitswitch work
gh pr create --title "Fix" --body "..."
# Creates PR as alice-work
```

## Verify active GitHub account

```bash
gh auth status
```

```
github.com
  ✓ Logged in to github.com account alice-work (keyring)
  - Active account: true
```

## Troubleshooting

**`gh auth switch` fails silently**

Run it manually to see the error:

```bash
gh auth switch --user alice-work
```

Common causes:
- Account not authenticated (`gh auth login`)
- Username typo in the profile (check with `gh auth status`)
- `gh` not installed

**`gh` account and git identity are out of sync**

This is expected when you use profiles without `--gh-user`. The two are independent:

```bash
gh auth status          # shows GitHub API account
gitswitch current       # shows git commit identity
```

To resync manually:

```bash
gh auth switch --user alice-work
```

## Next steps

- [Identity Awareness](/docs/features/identity-awareness)
- [Shell Integration](/docs/features/shell)
- [CLI Reference](/docs/cli/commands)
