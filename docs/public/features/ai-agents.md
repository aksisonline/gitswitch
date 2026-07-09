---
title: AI Coding Agents (Claude Code)
description: Easy git and GitHub setup for AI coding agents — install the Claude Code skill and auto-fix commit attribution
---

Claude Code, Cursor, and other AI coding agents run `git commit` for you — but they don't know which of your git identities to use, and they won't run `gitswitch` first unless told to. This page covers the two commands built for that: `gitswitch claude` (setup) and `gitswitch reauthor` (cleanup).

## The problem

An agent working across a freelance repo and a work repo in the same session will happily commit as whichever identity is currently active in your global git config — usually the wrong one. You only notice after the commits already exist.

## `gitswitch claude` — install the skill

```bash
gitswitch claude
```

Installs the embedded gitswitch skill into Claude Code. The `SKILL.md` ships inside the binary — no network request, no separate download. Once installed, Claude Code recognizes phrases like "wrong account" or "pushing to the wrong repo" and calls `gitswitch` itself instead of hand-editing `git config`.

| Flag | Description |
|------|-------------|
| `--scope user` | Install to `~/.claude/skills/gitswitch/` — active for every project (default) |
| `--scope project` | Install to `.claude/skills/gitswitch/` — active for this project only |

Reload Claude Code (or start a new session) after installing.

## `gitswitch reauthor` — fix commits an agent already made

```bash
gitswitch reauthor <base> --to <nickname> [--from <old-email>] [--push] [--yes]
```

Rewrites the author and committer on every commit between `<base>` and `HEAD` to a stored profile's identity, in one call — no hand-scripted `git rebase -i` / `git commit --amend` loop.

- `<base>` — a commit-ish (`HEAD~3`, a SHA) or a bare number `N` meaning "the last N commits"
- `--to <nickname>` — required; the profile to attribute commits to
- `--from <old-email>` — only rewrite commits currently authored by this email, leaving the rest untouched
- `--push` — force-push (`--force-with-lease`) after rewriting
- `--yes` — skip confirmation prompts; built for scripts and agents running non-interactively

```bash
# an agent committed the last 3 commits as your personal account inside a work repo
gitswitch reauthor 3 --to work --from personal@gmail.com --push --yes
```

This rewrites history — if the branch is already pushed and you don't pass `--push`, force-push it yourself once you've checked the result.

## Next steps

- [Commit Identity Switching](/docs/features/commit-identity)
- [CLI Reference](/docs/cli/commands)
