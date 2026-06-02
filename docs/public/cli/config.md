---
title: Configuration
description: Config files, storage layout, and environment variables
---

This page covers the three files gitswitch reads and writes, and the environment variable that controls the config directory.

## Config directory

All gitswitch data lives in:

```
~/.config/gitswitch/
├── profiles.json    # Profile definitions
├── config.json      # UI preferences (color theme)
└── history.json     # Per-repo usage counts and pins
```

## profiles.json

Stores all profile definitions. Written by `gitswitch add`, `gitswitch remove`, and the TUI editor.

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

Fields with empty values are omitted. `"active": true` marks the last-switched profile.

## config.json

Stores UI preferences. Written when you cycle themes in the TUI.

```json
{
  "color_theme": 3
}
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `color_theme` | integer | `0` | Active theme index (0–11). Changed with `c` in the TUI. |

**Theme index to name**

| Index | Name | Index | Name |
|-------|------|-------|------|
| 0 | Default | 6 | Arctic |
| 1 | Ocean | 7 | Gold |
| 2 | Sunset | 8 | Violet |
| 3 | Forest | 9 | Ember |
| 4 | Mono | 10 | Matrix |
| 5 | Rose | 11 | Steel |

To change theme without opening the TUI:

```bash
# Edit config.json directly
echo '{"color_theme": 3}' > ~/.config/gitswitch/config.json
```

## history.json

Stores per-repo usage counts and pins. Written by `gitswitch record`, `gitswitch pin`, and `gitswitch unpin`.

```json
{
  "repos": {
    "git@github.com:company/api.git": {
      "identities": {
        "work": 14,
        "personal": 1
      },
      "last_used": "work",
      "pinned": "work"
    },
    "git@github.com:alice/blog.git": {
      "identities": {
        "personal": 8
      },
      "last_used": "personal"
    }
  }
}
```

The repo key is the git remote URL (`git remote get-url origin`), or the absolute repo root path if no remote exists.

To clear all learning history:

```bash
rm ~/.config/gitswitch/history.json
```

## What gitswitch writes to ~/.gitconfig

On every profile switch, these keys in `~/.gitconfig` are updated:

| Key | Set when |
|-----|---------|
| `user.name` | Always |
| `user.email` | Always |
| `user.signingkey` | Profile has `sign_key` |
| `core.sshCommand` | Profile has `ssh_key` |

Inspect the current values:

```bash
git config --global user.name
git config --global user.email
git config --global user.signingkey
git config --global core.sshCommand
```

## Environment variables

| Variable | Description |
|----------|-------------|
| `GITSWITCH_CONFIG_DIR` | Override the config directory. Default: `~/.config/gitswitch` |

Example:

```bash
export GITSWITCH_CONFIG_DIR=~/.gitswitch
gitswitch  # reads/writes ~/.gitswitch instead
```

Useful for testing or if you want to keep gitswitch config outside `~/.config`.

## File permissions

gitswitch writes config files with mode `0600` (owner read/write only). The directory is created with `0755`. If you restore from backup:

```bash
chmod 700 ~/.config/gitswitch
chmod 600 ~/.config/gitswitch/*.json
```

## Troubleshooting

**`Error: config directory not found`**

Recreate it:

```bash
mkdir -p ~/.config/gitswitch
gitswitch init
```

**`Error: failed to read config: permission denied`**

Fix permissions:

```bash
chmod 600 ~/.config/gitswitch/profiles.json
```

**Corrupted history.json**

gitswitch automatically backs up a corrupted `history.json` to `history.json.bak` and starts fresh. You can inspect or restore the backup:

```bash
cat ~/.config/gitswitch/history.json.bak
cp ~/.config/gitswitch/history.json.bak ~/.config/gitswitch/history.json
```

**Profiles not loading after manual edit**

Validate JSON before saving:

```bash
cat ~/.config/gitswitch/profiles.json | python3 -m json.tool
# or
cat ~/.config/gitswitch/profiles.json | jq .
```

## Next steps

- [Profiles](/docs/cli/profiles)
- [CLI Reference](/docs/cli/commands)
- [Identity Awareness](/docs/features/identity-awareness)
