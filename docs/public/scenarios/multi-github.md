---
title: Multi-account GitHub
description: Cloning, pushing, and switching between two GitHub accounts on one machine, no manual setup
---

Setup: 2 GitHub accounts on one machine — `alice` (personal) and `alice-corp` (work). No SSH keys to generate, no keys to paste into GitHub, no manual git config.

## 1. One-time setup

```bash
gitswitch install   # shell integration, HTTPS routing, Session Isolation — all on by default
```

## 2. Connect both accounts

```bash
gitswitch login   # log in as alice in the browser
gitswitch login   # run again, log in as alice-corp
```

Each login creates a profile with the right name, email, and account automatically — nothing to type, nothing to test.

## 3. Pin your repos

```bash
cd ~/projects/my-library
gitswitch pin personal

cd ~/work/internal-api
gitswitch pin work
```

## Cloning from each account

```bash
gitswitch personal
git clone git@github.com:alice/my-project.git
# or the https:// URL — gitswitch routes the right token either way

gitswitch work
git clone git@github.com:alice-corp/company-service.git
```

## Daily workflow

```bash
cd ~/projects/my-library
git commit -m "Add feature"
git push
# committed and pushed as alice, automatically
```

```bash
cd ~/work/internal-api
git commit -m "Fix bug"
git push
# committed and pushed as alice-corp, automatically
```

```bash
gitswitch current
# work — Alice <alice@company.com>  (pinned to this repo)
```

## Managing notifications across both accounts

```bash
gitswitch personal
gh notification list   # personal account notifications

gitswitch work
gh notification list   # work account notifications
```

With [Session Isolation](/docs/features/shell#session-isolation) on (the default), bare `gh` commands already resolve to whichever account owns the repo you're in — no switching needed for `gh` calls made from inside a pinned repo.

## Contributing to someone else's repo from the "wrong" account

```bash
git clone git@github.com:alice/open-project.git
cd open-project
gitswitch work           # use the work identity for this contribution
gitswitch pin work        # pin it so the repo never asks again
git commit -m "Add contribution"
git push
```

## Fixed a push from the wrong account?

Not pulled by anyone else yet — safe to rewrite:

```bash
gitswitch reauthor 1 --to personal --push
# rewrites the last commit to the 'personal' identity and force-pushes safely
```

## Too many repos to manage individually?

Organize repos by account in a directory structure, then pin each one as you go — the identity awareness system tracks the rest:

```
~/github/
├── alice/           # personal repos
│   └── my-project/
└── alice-corp/      # work repos
    └── service/
```

## Prefer SSH keys you already manage yourself?

Optional — HTTPS (the default) needs nothing from you. Point a profile at an existing key instead:

```bash
gitswitch add personal "Alice" alice@gmail.com --ssh-key ~/.ssh/id_personal --gh-user alice
```

## Next steps

- [Open Source + Work](/docs/scenarios/oss)
- [Multi-client Freelancer](/docs/scenarios/freelancer)
- [SSH Keys](/docs/features/ssh-keys)
