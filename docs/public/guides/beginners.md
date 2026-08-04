---
title: New to Git
description: One account, zero tutorials — gitswitch is still worth it on day one
---
You don't need a second GitHub account for gitswitch to be worth installing. If this is your first repo, gitswitch is the fastest way to get git and GitHub actually working — no tutorials on SSH keys, no hand-editing `.gitconfig`.

## What you'd normally have to do

Without gitswitch, a first-time setup usually looks like:

```bash
git config --global user.name  "Your Name"
git config --global user.email "you@example.com"
ssh-keygen -t ed25519 -C "you@example.com"
cat ~/.ssh/id_ed25519.pub   # paste this into GitHub → Settings → SSH Keys
```

Four commands, one of which asks you to copy a public key into a web form by hand — and none of it explains what an SSH key even is.

## What gitswitch does instead

```bash
gitswitch
```

The first time you run it on a machine with no profiles yet, it checks whether `git` and the GitHub CLI (`gh`) are on your machine before anything else. If either is missing, it shows you the exact install command and offers to run it for you:

```
✗  git not found

About to run: sudo apt install git
Proceed? [y/N]
```

Say yes, and it installs it. (On macOS, git itself comes from Apple's Xcode Command Line Tools — that one's a GUI prompt, so gitswitch shows you the command but doesn't drive the dialog for you. `gh` still installs itself the same way as everywhere else.)

Once git and `gh` are there:

```bash
gitswitch login
```

Opens GitHub in your browser (device flow — you type a short code, GitHub does the rest), then fills in your name, email, and GitHub account automatically:

```
✓  Logged in as yourname (github.com)
✓  Profile "yourname" created
✓  Token stored in keychain
```

That's it. `git commit` and `git push` now work, correctly attributed, with no SSH key you had to generate or paste anywhere.

## You're done — but you're not stuck

If you ever *do* pick up a second account (a job, a client, an open-source project you contribute to), everything here already applies — see [Multi-account GitHub](/docs/guides/multi-github). Nothing about starting with one account locks you out of adding a second later; there's no migration step.

## Next

- **[Quick Start](/docs/get-started/quick-start)** — the full command walkthrough
- **[Multi-account GitHub](/docs/guides/multi-github)** — when you add a second account
- **[Troubleshooting](/docs/troubleshooting)** — if gitswitch can't find a package manager it recognizes
