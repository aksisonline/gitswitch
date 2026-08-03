---
title: Two GitHub Accounts
description: Personal and work on one machine, start to finish
---

The most common setup, and the one gitswitch was built for. Two GitHub accounts — `alice` (personal) and `alice-corp` (work) — on one laptop.

No SSH keys to generate. No keys to paste into a settings page. No `~/.ssh/config` host aliases.

## Set up

```bash
gitswitch install
```

Say yes to all three steps (they're the defaults), then reload your shell.

```bash
gitswitch login --profile personal    # log in as alice
gitswitch login --profile work        # run it again, log in as alice-corp
```

```bash
gitswitch list
```

```
✓  personal        Alice Smith <alice@gmail.com>
   work            Alice Smith <alice@company.com>
```

> Also make sure both accounts are logged in to `gh` (`gh auth status`, `gh auth login`). gitswitch gets push tokens from `gh`, so that's what makes the automatic routing below actually route.

## Claim your repos

```bash
cd ~/projects/my-library  && gitswitch pin personal
cd ~/work/internal-api    && gitswitch pin work
```

That's the whole setup. Everything from here is just... working.

## Every day

```bash
cd ~/work/internal-api
git commit -m "Fix bug"
git push
# committed as alice@company.com, pushed as alice-corp
```

```bash
cd ~/projects/my-library
git commit -m "Add feature"
git push
# committed as alice@gmail.com, pushed as alice — no switching, no thinking
```

Your prompt tells you which is which:

```
~/projects/my-library  [personal●] ❯
~/work/internal-api    [work●] ❯
```

And if you want it spelled out:

```bash
gitswitch current
# work — Alice Smith <alice@company.com>  (pinned to this repo)
```

## Cloning

```bash
gitswitch personal
git clone git@github.com:alice/my-project.git
# or the https:// URL — gitswitch routes the right token either way

gitswitch work
git clone git@github.com:alice-corp/company-service.git
```

Then pin the new clone and forget about it.

## `gh` follows along

With [Session Isolation](/docs/features/session-isolation) on (the default), bare `gh` commands resolve to whichever account owns the repo you're standing in — even in two terminals at once:

```bash
cd ~/work/internal-api && gh pr create      # as alice-corp
cd ~/projects/my-library && gh issue list   # as alice, simultaneously
```

Notifications for each account, no switching:

```bash
cd ~/work/internal-api    && gh api /notifications
cd ~/projects/my-library  && gh api /notifications
```

## Contributing to someone else's repo

Pick the identity you want to be seen as, then pin it so the repo never asks again:

```bash
git clone git@github.com:someone/open-project.git
cd open-project
gitswitch pin personal
git commit -m "Add contribution" && git push
```

## Already pushed as the wrong account?

```bash
gitswitch reauthor 1 --to personal --push
# rewrites the last commit and force-pushes with --force-with-lease
```

Then pin the repo so it can't happen twice. If someone else has already pulled those commits, tell them before you rewrite.

## Lots of repos?

Organise by account and pin as you go — the [learning](/docs/features/identity-awareness) covers the ones you haven't got to yet:

```
~/github/
├── alice/          # personal
└── alice-corp/     # work
```

## Prefer your own SSH keys?

Entirely optional — HTTPS is the default and needs nothing from you. But if you already manage keys:

```bash
gitswitch add personal "Alice" alice@gmail.com \
  --ssh-key ~/.ssh/id_personal --gh-user alice
```

See [SSH Keys](/docs/features/ssh-keys).

## Next

- **[Open Source + Work](/docs/scenarios/oss)** — when one account is a company account
- **[Multi-client Freelancer](/docs/scenarios/freelancer)** — three or more accounts
- **[Pins & Identity Awareness](/docs/features/identity-awareness)**
