---
title: Introduction
description: gitswitch keeps every commit and every push on the right account — automatically
---
**Every commit goes out as the right you. Automatically.**

New to git entirely? Skip the tutorials on SSH keys and `.gitconfig`. Install gitswitch and run it — the first time, on a machine with no profiles yet, it checks for git and the GitHub CLI and offers to install whichever's missing, then walks you through connecting a GitHub account. You're set up correctly without reading a single guide, and you don't need a second account for any of this to be worth it. See [gitswitch for git beginners](/docs/guides/beginners).

gitswitch keeps track of the git accounts you use — personal, work, each client — and makes sure the right one is used in the right repo. Not just your name and email: your SSH key, your signing key, and which GitHub account `git push` and `gh` talk to.

```bash
gitswitch            # first run sets up git/gh and connects your account automatically
gitswitch pin work   # this repo always uses that account
```

That's it. From then on you can forget about it.

## The problem

Git has exactly one global name and email:

```bash
git config --global user.name   "Alice Smith"
git config --global user.email  "alice@gmail.com"
```

Every commit you make anywhere on your machine uses those two values. So the moment you have more than one account, you're one distracted `git commit` away from a work repo full of commits stamped with your personal email — and once it's pushed, that's permanent.

It gets worse in the corners:

- **SSH keys** — your agent offers every key it has, so a push can authenticate as an account you didn't mean to use.
- **The GitHub CLI** — `gh` remembers one active account for your whole machine, so `gh pr create` in one terminal can be acting as a different account than you think.
- **AI coding agents** — Claude Code and friends run `git commit` for you, using whatever identity happens to be active.

## What gitswitch does about it

| | |
|---|---|
| **Switches everything at once** | name, email, SSH key, signing key, GitHub account — one command |
| **Pins repos** | this repo always uses this account, no thinking required |
| **Routes pushes** | HTTPS pushes use the right account's token, per repo |
| **Isolates terminals** | two terminals, two GitHub accounts, no fighting over `gh` |
| **Learns your habits** | notices which account you use where, and nudges you if something looks off |
| **Fixes past mistakes** | rewrites the author on commits you already made |
| **Stays local** | no servers, no telemetry, no account — tokens live in your OS keychain |

## Is this just `gh auth switch`?

No — they solve different halves of the problem. `gh auth switch` changes which GitHub *API* account the `gh` command uses. It doesn't touch what your commits actually say.

| | `gh auth switch` | gitswitch |
|---|---|---|
| Commit name and email | ✗ | ✓ |
| SSH key | ✗ | ✓ |
| Signing key | ✗ | ✓ |
| GitHub CLI account | ✓ (one, globally) | ✓ (per repo, per terminal) |
| Per-repo pinning | ✗ | ✓ |
| GitLab, Bitbucket, self-hosted | ✗ | ✓ — it's just git config |
| Interactive UI | ✗ | ✓ |

gitswitch happily uses `gh` under the hood — and if it's not installed yet, the first-run check offers to install it for you. Think of it as the layer above.

## What it is *not*

gitswitch is a router, not a vault. Your commits stay in git, your keys stay in `~/.ssh/`, your tokens stay in your OS keychain. gitswitch decides which one applies right now — it never proxies your traffic, phones home, or copies your credentials anywhere.

## Where to go next

- **[Install it](/docs/get-started/installation)** — Homebrew, curl, or Go
- **[Quick Start](/docs/get-started/quick-start)** — first account connected in about two minutes
- **[Scopes](/docs/accounts/scopes)** — global vs. repo vs. terminal, and how to tell which one is winning
- **[Multi-account GitHub](/docs/guides/multi-github)** — the most common setup, start to finish
- **[AI Coding Agents](/docs/ai/ai-agents)** — stop your agent committing as the wrong person
- **[CLI Reference](/docs/cli/commands)** — every command and flag

There is also a hidden arcade mode. We're not going to tell you what it does. ([Fine, here.](/docs/extras/arcade))
