---
title: Profiles
description: Profile structure, field constraints, and storage format
---

This page covers what a profile contains, how each field is used, and what the storage format looks like.

## Profile structure

A profile is a named set of git identity fields:

| Field | JSON key | Required | Example |
|-------|----------|----------|---------|
| Nickname | `nickname` | Yes | `work` |
| Full name | `user_name` | Yes | `Alice Smith` |
| Email | `email` | Yes | `alice@company.com` |
| SSH key path | `ssh_key` | No | `~/.ssh/id_work` |
| GPG key ID | `sign_key` | No | `ABCD1234EF567890` |
| GitHub username | `gh_user` | No | `alice-work` |

The `nickname` is used to identify the profile in CLI commands. It is never written to git config.

## Create a profile

### CLI

```bash
# Minimal
gitswitch add personal "Alice Smith" alice@gmail.com

# With all optional fields
gitswitch add work "Alice Smith" alice@company.com \
  --ssh-key ~/.ssh/id_work \
  --sign-key ABCD1234EF567890 \
  --gh-user alice-work
```

### TUI

```bash
gitswitch
# Press 'a' to open the add-profile form
# Tab between fields; Enter to save
```

## Field details

### Nickname

- No spaces. Lowercase recommended.
- Must be unique across all profiles.
- Used in: `gitswitch <nickname>`, `gitswitch switch <nickname>`, `gitswitch pin <nickname>`, `gitswitch remove <nickname>`
- Examples: `work`, `personal`, `client-a`, `oss`

### Full name (`user_name`)

- Becomes `git config user.name` on switch.
- Quote in the CLI if it contains spaces: `"Alice Smith"`
- Appears on every commit as the author name.

### Email

- Becomes `git config user.email` on switch.
- Appears on every commit as the author email.
- Must match the email registered with the git host for commit verification.

### SSH key (`ssh_key`)

- Path to the SSH **private** key (not the `.pub` file).
- On switch: sets `git config --global core.sshCommand "ssh -i <expanded-path> -o IdentitiesOnly=yes"`
- Tilde (`~`) and `$HOME` expand to your home directory.
- If omitted, SSH key routing is not changed when switching.

### GPG key (`sign_key`)

- 16-character GPG key ID from `gpg --list-secret-keys --keyid-format LONG`.
- On switch: sets `git config --global user.signingkey <key>`
- If omitted, the signing key is not changed when switching.
- See [GPG Signing](/docs/features/gpg) for key setup.

### GitHub username (`gh_user`)

- Your GitHub username (not email, not display name).
- On switch: runs `gh auth switch --user <username>` (best-effort).
- If omitted, the `gh` CLI account is not changed when switching.
- See [GitHub Account Sync](/docs/features/github-sync) for setup.

## Storage format

Profiles are stored as a JSON array in `~/.config/gitswitch/profiles.json`:

```json
[
  {
    "nickname": "personal",
    "user_name": "Alice Smith",
    "email": "alice@gmail.com",
    "active": false
  },
  {
    "nickname": "work",
    "user_name": "Alice Smith",
    "email": "alice@company.com",
    "ssh_key": "~/.ssh/id_work",
    "sign_key": "ABCD1234EF567890",
    "gh_user": "alice-work",
    "active": true
  }
]
```

Fields with empty values are omitted. The `active: true` field marks the last-switched profile.

You can edit this file manually, but prefer the CLI — a JSON syntax error will prevent gitswitch from loading any profiles.

## Edit a profile

Editing is done via the TUI:

```bash
gitswitch
# Select profile, press 'e' to open the edit form
# ctrl+d in the edit form to delete
```

## Remove a profile

```bash
gitswitch remove work
```

```
Profile 'work' removed
```

Removal is immediate and permanent. Any history entries for this nickname remain in `history.json` but are harmless.

## List profiles

```bash
gitswitch list
```

```
✓  personal        Alice Smith <alice@gmail.com>
   work            Alice Smith <alice@company.com>
```

## Backup and restore

All gitswitch data is in `~/.config/gitswitch/`:

```bash
# Backup
cp -r ~/.config/gitswitch ~/backups/gitswitch-$(date +%Y%m%d)

# Restore
cp -r ~/backups/gitswitch-20240115 ~/.config/gitswitch
```

## Next steps

- [Configuration](/docs/cli/config)
- [CLI Reference](/docs/cli/commands)
- [Quick Start](/docs/quick-start)
