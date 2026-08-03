---
title: GitHub CLI Sync
description: Keeping the gh CLI's account in step with your git identity
---

Your commit identity and your `gh` account are two different things. gitswitch can keep them together.

## `--gh-user`

Give an account a GitHub username and switching to it also switches `gh`:

```bash
gitswitch add work "Alice Smith" alice@company.com --gh-user alice-corp
gitswitch work
# runs: gh auth switch --user alice-corp
```

`gitswitch login` sets this for you automatically — the username comes back from GitHub.

The username must already be logged in to `gh`:

```bash
gh auth status      # who's logged in?
gh auth login       # add another account
```

## When `gh` isn't there

Nothing breaks. The `gh` step is the last one and it's advisory only — you get a warning and the git config switch still happens:

```
warning: gh auth switch --user alice-corp: exec: "gh": executable file not found in $PATH
✓ Switched to 'work' — Alice Smith <alice@company.com>
```

Accounts with no `--gh-user` leave the active `gh` account alone entirely.

## `gh auth switch` is global — and that's a problem

This is worth understanding, because it's the reason [Session Isolation](/docs/routing/session-isolation) exists.

`gh auth switch` flips one machine-wide setting. Switch identities in one terminal and every *other* terminal's `gh pr create` silently starts using that account too, even in unrelated repos:

```bash
# terminal 1
gitswitch work        # gh is now alice-corp, everywhere

# terminal 2, in a personal repo, five minutes later
gh pr create          # ...opens the PR as alice-corp. Oops.
```

**Session Isolation fixes this properly.** With it on (the default), each `gh` call resolves the account for the repo you're actually in and passes that account's token for that one command — no global state, no cross-terminal interference:

```bash
cd ~/work/api      && gh pr create     # acts as alice-corp
cd ~/personal/blog && gh issue list    # acts as alice — same moment, other terminal
```

So think of `--gh-user` as "which account does this profile belong to", and Session Isolation as the thing that makes it apply *per repo* rather than globally.

## Verify

```bash
gh auth status
```

```
github.com
  ✓ Logged in to github.com account alice-corp (keyring)
  - Active account: true
```

> With Session Isolation on, `gh auth status` still reports your old global account — by design. Isolation never touches gh's global account file; it overrides the token per call. Judge it by which account a `gh` command actually acted as.

## Troubleshooting

**`gh auth switch` fails**

Run it by hand to see why:

```bash
gh auth switch --user alice-corp
```

Usually: that account isn't logged in to `gh` (`gh auth login`), or the username in the profile has a typo (`gh auth status` to check spelling).

**`gh` and my git identity disagree**

Expected for accounts without `--gh-user` — they're independent. Compare:

```bash
gh auth status       # the GitHub API account
gitswitch current    # the git commit identity
```

## Next

- **[Session Isolation](/docs/routing/session-isolation)** — per-repo, per-terminal `gh`
- **[HTTPS Push Routing](/docs/routing/https)** — the same idea, for `git push`
- **[Connecting Accounts](/docs/accounts/accounts)** — `gitswitch login`
