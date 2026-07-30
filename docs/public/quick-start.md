---
title: Quick Start
description: From zero to "I never think about this again" in about two minutes
---

Three commands. That's the whole thing.

```bash
gitswitch install    # 1. set up
gitswitch login      # 2. connect an account
gitswitch pin work   # 3. claim a repo
```

Below is what each one actually does, in case you like knowing.

## 1. Set up

```bash
gitswitch install
```

A short wizard with three steps, each one explaining itself and each one skippable:

| Step | What you get |
|---|---|
| **Shell integration** | Your prompt shows which account is active, you get nudged if a repo looks wrong, and tab-completion works. |
| **HTTPS credential routing** | `git push` over HTTPS uses the right account's token instead of whatever your keychain hands over. |
| **Session Isolation** | Each terminal's `gh` commands resolve the account for the repo *you're in* — and it's what makes repo pins work. |

Then reload your shell:

```bash
source ~/.zshrc    # zsh
source ~/.bashrc   # bash
# or just open a new terminal
```

> Everything gitswitch writes lives between `# gitswitch` markers in your rc file. Re-running `install` replaces that block in place — it never appends a second copy.

## 2. Connect an account

```bash
gitswitch login
```

Opens a GitHub login in your browser (device flow — you paste a short code, GitHub does the rest):

```
  Open this URL in your browser:

    https://github.com/login/device

  Then enter the code:

    A1B2-C3D4

  Waiting for authorization...

  ✓  Logged in as alice-corp (github.com)
  ✓  Profile "alice-corp" created
  ✓  Token stored in keychain
```

Your name, email, and GitHub username are filled in for you. Run it again for each account you have:

```bash
gitswitch login                      # personal
gitswitch login --profile work       # work, with a nicer nickname
gitswitch login --host github.acme.com   # GitHub Enterprise
```

No SSH keys to generate, nothing to paste into a settings page. If you'd rather bring your own keys, see [SSH Keys](/docs/features/ssh-keys) — or type a profile in by hand:

```bash
gitswitch add work "Alice Smith" alice@company.com --gh-user alice-corp
```

See them all:

```bash
gitswitch list
```

```
✓  personal        Alice Smith <alice@gmail.com>
   work            Alice Smith <alice@company.com>
```

The `✓` is your current global identity.

## 3. Claim your repos

```bash
cd ~/work/api
gitswitch pin work
```

```
✓ Pinned 'work' to this repo — Alice Smith <alice@company.com> (local git config; global identity unchanged)
```

Every commit in that repo now uses `work`. Every push uses `work`'s token. Every `gh pr create` acts as `work`. Your global identity — and every other repo — is left exactly as it was.

Do it once per repo and you are finished thinking about this:

```bash
cd ~/personal/blog     && gitswitch pin personal
cd ~/clients/acme/api  && gitswitch pin acme
```

Changed your mind? `gitswitch unpin`.

> **Already configured a repo by hand years ago?** Run `gitswitch pin` with no name. gitswitch reads the repo's existing `user.email`, matches it to one of your accounts, and fills in the rest.

## Switching by hand

Pins cover the repos you own. For everything else:

```bash
gitswitch work        # switch global identity, right now
gitswitch current     # who am I in here?
```

```
work — Alice Smith <alice@company.com>  (pinned to this repo)
```

That parenthetical matters — it tells you *why* you're that person. See [Scopes](/docs/concepts/scopes).

## The interactive UI

```bash
gitswitch
```

Three tabs — **Accounts**, **Utilities**, **Settings** — switched with `1` `2` `3`, `tab`, or a mouse click. Yes, the mouse works.

| Key | Does |
|---|---|
| `↑` `↓` / `k` `j` | Move |
| `enter` | Switch to this account (globally) |
| `p` | Pin this account to the current repo — press again to unpin |
| `a` | Add an account |
| `e` | Edit — `ctrl+d` inside the form deletes |
| `v` | Show GitHub usernames instead of emails |
| `c` | Next color theme (there are 12) |
| `u` | Upgrade — only appears when there's something to upgrade to |
| `?` | CLI cheat sheet |
| `q` | Out |

First launch with no accounts runs a short onboarding wizard — it scans your `gh` logins and `~/.ssh/` keys and offers to import what it finds.

## Where next

- **[Scopes](/docs/concepts/scopes)** — global, repo, and terminal identities, and who wins
- **[Multi-account GitHub](/docs/scenarios/multi-github)** — the full two-account walkthrough
- **[Identity Awareness](/docs/features/identity-awareness)** — the nudges, and how they're learned
- **[Troubleshooting](/docs/troubleshooting)** — when a push still goes out wrong
