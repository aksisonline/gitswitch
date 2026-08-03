---
title: Session Isolation
description: Two terminals, two GitHub accounts, no fighting
---

## The problem

`gh` remembers exactly one "active" account for your whole machine. So this happens:

```bash
# terminal 1 — work repo
gitswitch work          # gh's active account is now alice-corp

# terminal 2 — personal repo, same moment
gh pr create            # ...opens the PR as alice-corp
```

Nothing warned you. Nothing was misconfigured. `gh` simply doesn't have a concept of "this terminal" or "this repo".

## What Session Isolation does

It installs a small `gh` shell function that wraps the real `gh`. Before each call it asks gitswitch which account applies to the directory you're in, gets that account's token, and passes it as `GH_TOKEN` for that one command:

```bash
cd ~/work/api      && gh pr create      # acts as alice-corp
cd ~/personal/blog && gh issue list     # acts as alice
```

gh's global active-account file is never read or written, so any number of terminals can work as different accounts simultaneously without stepping on each other.

It's also the thing that makes [repo pins](/docs/routing/identity-awareness#pin-a-repo) real — a pin only takes effect while Session Isolation is on, because this is the machinery that keeps a repo's account separate from the rest of your machine.

## Turn it on

It's a step in `gitswitch shell`, on by default. Later:

```bash
gitswitch          # Utilities tab → toggle Session Isolation
```

Then reload your shell — it's a shell function, so terminals opened before you flipped the toggle won't have it:

```bash
source ~/.zshrc    # zsh
source ~/.bashrc   # bash
```

Pinning a repo turns it on for you if it was off, so `gitswitch pin work` just works whether or not you've ever visited the Utilities tab.

Check it:

```bash
gitswitch doctor
# ✓  Session Isolation active (bare `gh` commands resolve per-repo)
```

> **Needs `gh`.** The wrapper gets its tokens from `gh auth token --user <account>`, so the account has to be logged in to `gh`. Without that, the wrapper installs but falls straight through to plain `gh` — no change in behavior, nothing broken.

## How the account is chosen

Exactly the same resolution [HTTPS push routing](/docs/routing/https) uses, so `git push` and `gh` can never disagree about who you are:

1. Repo pin or terminal session override → that account
2. Otherwise, a learned/pinned recommendation for this repo → that account
3. Otherwise, your globally active account

If nothing resolves — outside a git repo, or no account with a GitHub username — the wrapper does nothing and plain `gh` runs as normal.

## If you turn it off

Pins don't disappear, but they stop applying. gitswitch says so plainly rather than pretending a local git config override still means what it used to:

```
work — Alice Smith <alice@company.com>
  (pinned to 'acme' — inactive, Session Isolation is off)
```

`gitswitch current`, the TUI, and your prompt all agree on this.

To turn it off: toggle it in the Utilities tab, or `gitswitch uninstall` to remove all shell integration at once. Either way the wrapper is removed immediately and bare `gh` behaves exactly as before.

## Troubleshooting

**`gh` still uses the wrong account**

Reload the shell. Terminals opened before you enabled it don't have the function.

**It looks like it's doing nothing**

`gh auth status` will still show your old global account — deliberately. Isolation overrides the *token* for individual calls and never touches gh's account list. Judge it by which account a `gh` command actually acted as.

**My pin isn't being used**

Check `gitswitch doctor`. A pin set while isolation was off is stored but inactive; re-run `gitswitch pin <name>` to switch isolation on and reapply it.

## Next

- **[HTTPS Push Routing](/docs/routing/https)** — the same resolution, for `git push`
- **[Identity Awareness](/docs/routing/identity-awareness)** — pins and learned habits
- **[Scopes](/docs/accounts/scopes)** — why a repo can outrank your global identity
