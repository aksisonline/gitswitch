---
title: Configuration
description: Every file gitswitch reads or writes, and the one environment variable
---

## The config directory

```
~/.config/gitswitch/
├── config.yaml           # your accounts
├── config.json           # UI preferences
├── history.json          # per-repo usage counts and pins
├── hook-version          # which shell hook version is installed
├── whats-new-pending.txt # a queued "What's New" screen
├── profiles.json.v1.bak  # backup from the pre-YAML migration, if you upgraded
└── history.json.bak      # backup, only if history was ever corrupted
```

The path is fixed at `~/.config/gitswitch/`. Files are written with mode `0600`, the directory with `0755`.

## `config.yaml` — your accounts

```yaml
version: 2
profiles:
  - nickname: work
    user_name: Alice Smith
    email: alice@company.com
    ssh_key: ~/.ssh/id_work
    gh_user: alice-corp
    token_ref: gitswitch:work:github.com
    active: true
```

Written by `add`, `login`, `remove`, and the editor in the UI. Full field reference: [Account Fields](/docs/cli/profiles).

Open it from the UI's **Settings** tab → *Config Location* — it launches your `$EDITOR` and reloads your accounts as soon as you save.

## `config.json` — preferences

```json
{
  "color_theme": 3,
  "show_username": false,
  "shell_alias": "gs",
  "shell_alias_disabled": false,
  "arcade_hi_score": 7400,
  "arcade_mode": false
}
```

| Field | Default | |
|---|---|---|
| `color_theme` | `0` | Theme index, 0–11. `c` in the UI cycles it. |
| `show_username` | `false` | Show GitHub usernames instead of emails in the list. `v` toggles it. |
| `shell_alias` | `"gs"` | The short alias installed with shell integration. |
| `shell_alias_disabled` | `false` | Turn the alias off. |
| `arcade_hi_score` | — | Your [arcade mode](/docs/features/arcade) high score. Yes, really. |
| `arcade_mode` | `false` | Whether `gitswitch pacman` left arcade mode on. |
| `shell_enabled` | — | Internal: shell integration state. |
| `splash_seen_020` | — | Internal: whether the one-time upgrade splash has been shown. |

**Themes by index**

| | | | |
|---|---|---|---|
| 0 Default | 3 Forest | 6 Arctic | 9 Ember |
| 1 Ocean | 4 Mono | 7 Gold | 10 Matrix |
| 2 Sunset | 5 Rose | 8 Violet | 11 Steel |

```bash
echo '{"color_theme": 10}' > ~/.config/gitswitch/config.json    # Matrix, without opening the UI
```

## `history.json` — what it learned

```json
{
  "repos": {
    "git@github.com:company/api.git": {
      "identities": { "work": 14, "personal": 1 },
      "last_used": "work",
      "pinned": "work"
    }
  }
}
```

The key is `git remote get-url origin`, falling back to the absolute repo root path when there's no remote. Written by `record`, `pin`, and `unpin`.

```bash
rm ~/.config/gitswitch/history.json    # forget everything it learned
```

If this file is ever corrupted, gitswitch copies it to `history.json.bak` and starts fresh rather than refusing to run.

## Where tokens live

Not in any of these files. `gitswitch login` puts tokens in your OS keychain — macOS Keychain (local, never iCloud-synced), Windows Credential Manager, or the Linux Secret Service — filed under `gitswitch:<nickname>:<host>`. The account records only that reference.

## Environment variables

| | |
|---|---|
| `GITSWITCH_SECRETS_BACKEND=none` | Don't touch the keychain at all. gitswitch stores no tokens; everything else works normally. |

## What gitswitch writes to your git config

On a global switch (`~/.gitconfig`):

| Key | When |
|---|---|
| `user.name` | always |
| `user.email` | always |
| `user.signingkey` | account has a signing key — cleared otherwise |
| `gpg.format` | signing key is an SSH key — cleared otherwise |
| `core.sshCommand` | account has an SSH key — cleared otherwise |

On a pin (`<repo>/.git/config`), the same five keys, scoped to that repo. `gitswitch unpin` removes them.

Inspect any time:

```bash
git config --global --get-regexp '^(user|gpg|core\.sshCommand)'
git config --local  --get-regexp '^(user|gpg|core\.sshCommand)'
```

## What gitswitch writes to your shell

Only inside marked blocks, so nothing else in your rc file is ever touched:

| | |
|---|---|
| `# gitswitch shell integration` block | Prompt function, `cd` nudge hook, completion, `gs` alias |
| `# gitswitch gh wrapper` block | The [Session Isolation](/docs/features/session-isolation) `gh` function — separate, so it can be toggled independently |
| `~/.config/starship.toml` | A `[custom.gitswitch]` block, on Starship |
| `~/.oh-my-zsh/custom/plugins/gitswitch/` | A plugin file, on oh-my-zsh |
| `credential.helper` entries in `~/.gitconfig` | [HTTPS routing](/docs/features/https) — added in front of existing helpers, never replacing them |

`gitswitch uninstall` removes all of it.

## Fixing a broken config

```bash
gitswitch doctor                        # start here
cat ~/.config/gitswitch/config.yaml     # valid YAML?
chmod 700 ~/.config/gitswitch && chmod 600 ~/.config/gitswitch/*   # after a restore
```

More in [Troubleshooting](/docs/troubleshooting).

## Next

- **[Account Fields](/docs/cli/profiles)**
- **[Commands](/docs/cli/commands)**
