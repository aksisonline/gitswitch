---
title: Account Fields
description: What an account contains, what each field does, and the file it lives in
---

An account (also called a profile) is a nickname plus the git settings that go with it.

| Field | YAML key | Required | Example |
|---|---|---|---|
| Nickname | `nickname` | yes | `work` |
| Full name | `user_name` | yes | `Alice Smith` |
| Email | `email` | yes | `alice@company.com` |
| SSH key path | `ssh_key` | no | `~/.ssh/id_work` |
| Signing key | `sign_key` | no | `ABCD1234EF567890` or `~/.ssh/id_work.pub` |
| GitHub username | `gh_user` | no | `alice-corp` |
| Keychain reference | `token_ref` | auto | `gitswitch:work:github.com` |

## Field by field

### `nickname`

Your label for the account. No spaces, lowercase reads best, must be unique. It's what you type — `gitswitch work`, `gitswitch pin work` — and it is **never** written to git config.

### `user_name` and `email`

Become `user.name` and `user.email` on switch, and appear on every commit you make. The email needs to be verified on the git host if you want commits linked to your account there.

### `ssh_key`

Path to the **private** key (not the `.pub`). On switch:

```bash
git config --global core.sshCommand "ssh -i <expanded-path> -o IdentitiesOnly=yes"
```

`~` and `$HOME` are expanded when applied. If an account has **no** SSH key, switching to it *unsets* `core.sshCommand` — so the previous account's key can't linger. See [SSH Keys](/docs/features/ssh-keys).

### `sign_key`

Either a GPG key ID or an SSH key — gitswitch tells them apart by the value:

- Hex key ID → sets `user.signingkey`, clears `gpg.format` (OpenPGP)
- A path or an inline `ssh-…` key → sets `gpg.format=ssh` and `user.signingkey`

No signing key clears **both**, so a stale `gpg.format=ssh` never breaks the next account. See [Commit Signing](/docs/features/gpg).

### `gh_user`

Your GitHub username — not your email, not your display name. On switch it runs `gh auth switch --user <username>` (best-effort), and it's what [HTTPS routing](/docs/features/https) and [Session Isolation](/docs/features/session-isolation) use to find the right token. See [GitHub CLI Sync](/docs/features/github-sync).

### `token_ref`

Set for you by `gitswitch login`. It's a *reference* to a keychain entry, never the token itself. Editing an account keeps it, so fixing a typo won't log you out.

## Creating

```bash
gitswitch login                             # let GitHub fill it in
gitswitch add work "Alice Smith" alice@company.com --gh-user alice-corp   # or type it
gitswitch                                    # or press `a` in the UI
```

## Editing and deleting

In the UI — `e` on the account, `ctrl+d` inside the form to delete. Or:

```bash
gitswitch remove work
```

## On disk

`~/.config/gitswitch/config.yaml`:

```yaml
version: 2
profiles:
  - nickname: personal
    user_name: Alice Smith
    email: alice@gmail.com
    active: false
  - nickname: work
    user_name: Alice Smith
    email: alice@company.com
    ssh_key: ~/.ssh/id_work
    sign_key: ABCD1234EF567890
    gh_user: alice-corp
    token_ref: gitswitch:work:github.com
    active: true
```

Empty fields are omitted. `active: true` marks the last account you switched to. Written atomically with mode `0600`.

You can hand-edit it — the Settings tab even opens it in your `$EDITOR` and reloads when you're done — but a YAML syntax error means gitswitch can load *nothing*, so prefer the CLI or UI for anything routine.

> **Upgrading from an older version?** gitswitch migrates `profiles.json` → `config.yaml` on first run and leaves `profiles.json.v1.bak` behind as a safety net. Nothing to do.

## Backup

```bash
cp -r ~/.config/gitswitch ~/backups/gitswitch-$(date +%Y%m%d)
```

Note that tokens are in your OS keychain, not in that folder — a restore on a new machine means re-running `gitswitch login`, which takes about ten seconds per account.

## Next

- **[Configuration](/docs/cli/config)** — every file and environment variable
- **[Commands](/docs/cli/commands)** — the full CLI
- **[Connecting Accounts](/docs/features/accounts)** — `gitswitch login`
