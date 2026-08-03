---
title: Open Source + Work
description: A company account and a public one, kept strictly apart
---

Two accounts: `alice-corp` for the day job, `alice` for your own open source. The stakes here are slightly different from the other setups — you really don't want a company email in a public commit log, or a personal account opening a PR on an internal repo.

## Set up

```bash
gitswitch install
```

```bash
gitswitch login --profile work    # log in as alice-corp
gitswitch login --profile oss     # run again, log in as alice
```

Both accounts should also be logged in to `gh` (`gh auth status`) — that's where push tokens come from.

## Claim your repos

```bash
cd ~/work/internal-api   && gitswitch pin work
cd ~/projects/my-library && gitswitch pin oss
```

The identity now lives in each repo's own git config, so it holds even when you switch globally for something else — and even when an AI agent commits without asking you first.

## Every day

```bash
cd ~/work/internal-api
git commit -m "Fix critical bug"
git push origin fix/critical-bug
# alice-corp, always
```

```bash
cd ~/projects/my-library
git commit -m "Add feature"
git push origin add-feature
# alice, always
```

`gh` comes along for the ride — with [Session Isolation](/docs/features/session-isolation) on (the default), `gh pr create` and `gh issue list` resolve to whichever account owns the repo you're in, so you can review a work PR and triage an OSS issue in two terminals at once.

## Contributing to a public repo you don't own

```bash
git clone https://github.com/someone/cool-project.git
cd cool-project
gitswitch pin oss
```

Pin *before* your first commit and there's nothing to clean up later.

## Signed commits for the public account

Many OSS projects want verified commits. If you have an SSH key, you already have a signing key:

```bash
gitswitch add oss "Alice" alice@users.noreply.github.com \
  --sign-key ~/.ssh/id_ed25519.pub --gh-user alice
git config --global commit.gpgsign true
```

Register it on GitHub as a **Signing Key** (a separate entry from an authentication key, even for the same key). Full detail: [Commit Signing](/docs/features/gpg).

> Using GitHub's `@users.noreply.github.com` address for public commits keeps your real email out of a permanent public log while still linking commits to your account.

## Company email in a public commit log

Fix it before anyone notices:

```bash
gitswitch reauthor 3 --to oss --from alice@company.com --push
```

Only the commits authored by the work email get rewritten. This rewrites history, so if the branch has contributors on it, say something first — and if it's already merged upstream, that ship has sailed; the fix is to pin the repo so it stops recurring.

## Verify before you push

```bash
gitswitch current
# oss — Alice <alice@users.noreply.github.com>  (pinned to this repo)

git log -1 --format="%an <%ae>"
```

## Prefer SSH over HTTPS?

Optional — HTTPS is the default and needs nothing from you:

```bash
gitswitch add work "Alice Smith" alice@company.com \
  --ssh-key ~/.ssh/id_work --gh-user alice-corp
```

## Next

- **[Two GitHub Accounts](/docs/scenarios/multi-github)**
- **[AI Coding Agents](/docs/features/ai-agents)** — agents are the main way this goes wrong now
- **[Commit Signing](/docs/features/gpg)**
