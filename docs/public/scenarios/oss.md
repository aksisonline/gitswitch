---
title: Open Source + Work
description: Separate git identities for day-job and personal open source contributions
---

Setup: 2 GitHub accounts — `alice-work` for the company and `alice` for personal/OSS.

No SSH keys to generate, no keys to paste into GitHub's settings page, no manual git config. gitswitch handles all of it.

## 1. One-time setup

```bash
gitswitch install   # shell integration, HTTPS routing, Session Isolation — all on by default
```

## 2. Connect both accounts

```bash
gitswitch login   # opens GitHub in your browser, log in as alice-work
gitswitch login   # run again, log in as alice this time
```

Each run creates a profile automatically — name, email, and a securely stored token, no typing required. Rename them if you want friendlier nicknames:

```bash
gitswitch list
gitswitch add work "Alice Smith" alice@company.com --gh-user alice-work   # only if you'd rather set it up by hand
```

## 3. Pin repos

```bash
cd ~/work/internal-api
gitswitch pin work

cd ~/projects/my-library
gitswitch pin oss
```

That's it. From here, every commit, push, and `gh` call in these repos just uses the right account — nothing left to remember.

## Daily workflow

```bash
cd ~/work/internal-api
# nudge: "this repo usually uses work — switch? [y/N]"  — only shows up if you didn't pin

git commit -m "Fix critical bug"
git push origin fix/critical-bug
# committed and pushed as alice-work, no thinking required
```

```bash
cd ~/projects/my-library
git commit -m "Add feature"
git push origin add-feature
# committed and pushed as alice, automatically
```

`gh` follows along too — with [Session Isolation](/docs/features/shell#session-isolation) on (the default), `gh pr create`, `gh issue list`, etc. always resolve to whichever account owns the repo you're in.

## Verify correct attribution

```bash
gitswitch current
# oss — Alice <alice@github.com>  (pinned to this repo)
```

## Fixed a commit with the wrong identity?

Not pushed yet:

```bash
gitswitch oss
git commit --amend --reset-author --no-edit
```

Already pushed — rewrite it in one command instead of hand-rolling a rebase:

```bash
gitswitch reauthor 3 --to oss --push
# rewrites the last 3 commits to the 'oss' identity and force-pushes safely
```

## Prefer SSH over HTTPS?

Optional — HTTPS (the default) needs nothing from you. If you'd rather use SSH keys you already manage yourself, point a profile at one:

```bash
gitswitch add work "Alice Smith" alice@company.com --ssh-key ~/.ssh/id_work --gh-user alice-work
```

## Next steps

- [Multi-account GitHub](/docs/scenarios/multi-github)
- [Multi-client Freelancer](/docs/scenarios/freelancer)
- [GitHub Account Sync](/docs/features/github-sync)
