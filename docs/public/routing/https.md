---
title: HTTPS Push Routing
description: Make git push over HTTPS use the right account's token, per repo
---

If your remotes are `https://github.com/...` URLs (they are, if you cloned by copy-pasting from the GitHub UI), then `git push` doesn't care about your commit identity at all — it asks a **credential helper** for a username and password. Your keychain usually answers with whatever it happened to save first. That's how you end up pushing to a work repo as your personal account.

gitswitch fixes this by answering that question itself, correctly, per repo.

## Turn it on

```bash
gitswitch shell
```

It's one of the wizard steps, on by default. Or toggle **HTTPS Credential Helper** in the TUI's **Utilities** tab.

```bash
gitswitch doctor    # confirm it's actually in charge
```

```
  ✓  git 2.50.1
  ✓  gh  2.95.0
  ✓  HTTPS pushes routed by gitswitch
  ✓  Session Isolation active (bare `gh` commands resolve per-repo)
```

> **Needs `gh`.** gitswitch stores no push tokens of its own for this — it asks the GitHub CLI for the right account's token (`gh auth token --user <account>`). Without `gh` installed and that account logged in, the helper installs but stays quiet and git falls through to your keychain as before.

## How it picks an account

Same resolution as [Session Isolation](/docs/routing/session-isolation), so pushes and `gh` commands can never disagree about who you are:

1. Is this repo [pinned](/docs/routing/identity-awareness#pin-a-repo), or does this terminal override the identity? Use that account.
2. Otherwise, does the repo have a learned/pinned recommendation in your history? Use that.
3. Otherwise, use your globally active account.

If none of that resolves to a GitHub account, gitswitch stays silent and git moves on to the next helper. It never breaks a push it can't help with.

## Why the ordering is fiddly (and why `gh` can steal it)

Worth knowing if you also use `gh`, because git's rules here surprise everyone.

Git doesn't pick *one* credential helper. For a given URL it collects `credential.helper` **and** every matching `credential.<host>.helper` into a single list, in config order, then asks each in turn until one answers. Two consequences:

- A host-specific entry does **not** replace the generic one — it just joins the list.
- A helper set to the **empty string** discards everything collected before it.

`gh auth setup-git` writes exactly that empty entry, followed by its own helper. That's how it removes your keychain helper for github.com — and it would erase gitswitch too if gitswitch just appended itself.

So `gitswitch shell` makes gitswitch the first *live* entry in every helper list it finds, leaving everything else in place behind it:

```ini
[credential "https://github.com"]
	helper =                                              # gh's reset, preserved
	helper = !gitswitch credential                        # gitswitch answers first
	helper = !/opt/homebrew/bin/gh auth git-credential    # fallback, untouched
```

Nothing is removed, and `gitswitch uninstall` takes only gitswitch's own line back out.

### Re-running `gh auth setup-git` undoes it

gh rewrites that whole key with `--replace-all`, dropping gitswitch's line along with anything else. An interactive `gh auth login` does it too if you answer yes to *"Authenticate Git with your GitHub credentials?"*.

`gitswitch doctor` catches it and names the culprit:

```
  ✗  HTTPS pushes answered by another helper before gitswitch:
       credential.https://github.com.helper → !/opt/homebrew/bin/gh auth git-credential
       pushes may use the wrong account — run: gitswitch shell
```

```bash
gitswitch shell    # puts gitswitch back in front, keeps gh's helper as fallback
```

## Prefer SSH?

Then you don't need any of this — SSH auth is decided by the key, and gitswitch handles that separately. See [SSH Keys](/docs/accounts/ssh-keys). Nothing stops you using both: HTTPS for some repos, SSH for others.

## Turn it off

```bash
gitswitch uninstall    # removes everything gitswitch shelled
```

Or toggle **HTTPS Credential Helper** off in the Utilities tab.

## Next

- **[Session Isolation](/docs/routing/session-isolation)** — the same routing, for `gh` commands
- **[Troubleshooting](/docs/troubleshooting)** — a push still went out as the wrong account
- **[Shell Integration](/docs/routing/shell)** — everything else `gitswitch shell` sets up
