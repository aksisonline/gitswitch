---
title: Connecting Accounts
description: gitswitch login, the browser flow, and where your tokens live
---
An "account" (or profile) in gitswitch is a named bundle: a nickname, a git name and email, and optionally an SSH key, a signing key, and a GitHub username. The fastest way to create one is to let GitHub fill it in for you.

## `gitswitch login`

```bash
gitswitch login
```

```
  ┌──────────────────────────────────────────┐
  │  gitswitch · Log in with GitHub          │
  └──────────────────────────────────────────┘

  (gh auth login runs here — browser, device code, or token, whichever gh offers)

  ✓  Logged in as alice-corp (github.com)
  ✓  Profile "alice-corp" created
```

gitswitch delegates credential storage entirely to GitHub CLI — you approve in your browser, gitswitch never sees your password, and the token lives in `gh`'s secure store, not anywhere gitswitch manages.

Your name, email, and GitHub username come back from the API and become the profile. The nickname defaults to your GitHub username.

**Flags**

| Flag | |
|---|---|
| `--profile <nickname>` | Pick the nickname yourself instead of using the GitHub username |
| `--host <hostname>` | GitHub Enterprise Server, e.g. `--host github.acme.com`. Defaults to `github.com` |

Run it once per account:

```bash
gitswitch login --profile personal
gitswitch login --profile work
```

The first profile you create becomes your active identity automatically.

### Logging in again

Safe. Running `gitswitch login` on an account you already have refreshes the token and updates the name/email, and **keeps** the SSH key and signing key you configured. It won't wipe your setup to start over.

## Adding an account by hand

No browser, no GitHub, works for GitLab/Bitbucket/anything:

```bash
gitswitch add <nickname> <name> <email> [flags]
```

```bash
# the minimum
gitswitch add personal "Alice Smith" alice@gmail.com

# the works
gitswitch add work "Alice Smith" alice@company.com \
  --ssh-key  ~/.ssh/id_work \
  --sign-key ABCD1234EF567890 \
  --gh-user  alice-corp
```

| Flag | |
|---|---|
| `--ssh-key <path>` | Private key to force for this account — see [SSH Keys](/docs/accounts/ssh-keys) |
| `--sign-key <key>` | GPG key ID *or* an SSH key for SSH signing — see [Commit Signing](/docs/accounts/gpg) |
| `--gh-user <username>` | GitHub username, for `gh` and push routing — see [GitHub Account Sync](/docs/routing/github-sync) |

## Where your tokens live

gitswitch stores nothing. GitHub CLI (`gh`) owns all credential storage — tokens live in your OS keychain, never in a file, never in plaintext. gitswitch just reads back your logged-in GitHub identity via `gh` after you've authenticated.

| Platform | Store |
|---|---|
| macOS | Keychain (local only — never synced to iCloud) |
| Linux | Secret Service: GNOME Keyring, KeePassXC, anything libsecret-compatible |
| Windows | Credential Manager |

If you revoke gitswitch entirely but keep `gh`, your tokens stay safe in `gh` — it manages them independently.

## Editing and removing

Editing lives in the UI:

```bash
gitswitch
# select the account, press e
# ctrl+d inside the form deletes it
```

Saving an edit keeps your existing GitHub login — fixing a typo in your name won't log you out.

From the CLI:

```bash
gitswitch remove work
```

That deletes the profile. It does **not** rewind any git config that was applied while it was active, and it doesn't revoke anything on GitHub's side.

## Importing what you already have

First launch runs a short wizard that scans your `gh` logins and the private keys in `~/.ssh/`, and offers to import them as accounts. To pull in your current global git config as a profile named `default` at any time:

```bash
gitswitch init
```

## Next

- **[Commit Identity](/docs/accounts/commit-identity)** — what switching writes
- **[Profiles reference](/docs/cli/profiles)** — every field, and the storage format
- **[Quick Start](/docs/get-started/quick-start)** — connect, pin, done
