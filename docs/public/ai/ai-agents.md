---
title: AI Coding Agents
description: Stop your agent committing as the wrong person — and fix it when it already did
---
Claude Code, Cursor, and friends run `git commit` for you. They don't know which of your accounts belongs in this repo, and they won't run `gitswitch` first unless something tells them to. An agent working across a client repo and your day job in one session will happily attribute everything to whichever account happens to be active.

Two commands exist for exactly this.

## Prevention: pin your repos

Boring answer, best answer. A [pinned repo](/docs/routing/identity-awareness#pin-a-repo) writes the identity into the repo's own git config, so an agent doesn't have to know anything — git itself uses the right account:

```bash
gitswitch pin work
```

## `gitswitch skills` — teach your agent about this

```bash
gitswitch skills
```

Installs the gitswitch skill for whichever AI coding agent you're using. Tries [skills.sh](https://www.skills.sh) first (`npx skills add aksisonline/gitswitch`), which places it correctly for Claude Code, Cursor, Codex, and anything else that speaks the Agent Skills format — gitswitch doesn't need to know each tool's own layout. No network, or no `npx`? It falls back to installing directly for Claude Code (`~/.claude/skills/`) and the shared `.agents/skills/` convention (Codex and others).

After that, saying "I'm committing as the wrong account" gets you an actual fix instead of a hand-edited `git config`. The skill is embedded in the binary — no download, always matching your installed version.

| Flag | |
|---|---|
| `--scope user` | `~/.claude/skills/gitswitch/`, `~/.agents/skills/gitswitch/` — every project (default) |
| `--scope project` | `.claude/skills/gitswitch/`, `.agents/skills/gitswitch/` — this project only |
| `--offline` | Skip the skills.sh attempt, go straight to the offline installer |

Reload your agent (or start a new session) afterwards.

## `gitswitch reauthor` — fix commits that already exist

```bash
gitswitch reauthor <base> --to <nickname> [--from <old-email>] [--push] [--yes]
```

Rewrites the author *and* committer on every commit between `<base>` and `HEAD` to one of your accounts, in one command — no hand-rolled `git rebase -i` / `git commit --amend` loop.

| | |
|---|---|
| `<base>` | A commit-ish (`HEAD~3`, a SHA) or just a number: `3` means "the last 3 commits" |
| `--to <nickname>` | Required — the account to attribute them to |
| `--from <old-email>` | Only rewrite commits currently authored by this email; leave everything else alone |
| `--push` | Force-push afterwards, with `--force-with-lease` |
| `--yes` / `-y` | No confirmation prompts — for scripts and agents running unattended |

```bash
# an agent made the last 3 commits as your personal account, inside a work repo
gitswitch reauthor 3 --to work --from personal@gmail.com --push --yes
```

This rewrites history. If the branch is already pushed and you don't pass `--push`, force-push it yourself once you've looked at the result. If anyone else has pulled those commits, talk to them first.

## Machine-readable everything

Every state-reading command speaks JSON, so an agent never has to parse prose:

```bash
gitswitch current --json    # who am I here, and why (scope: global | repo | session)
gitswitch list --json       # all accounts
gitswitch doctor --json     # a one-shot manifest: version, profile count, git/gh state, is HTTPS routed, is isolation on
```

```json
{
  "nickname": "work",
  "user_name": "Alice Smith",
  "email": "alice@company.com",
  "scope": "repo",
  "gh_user": "alice-corp",
  "credential_helper_active": true
}
```

`"scope": "repo"` is the useful bit — it means the repo is pinned and the identity is guaranteed regardless of what's global.

## A workflow that holds up

```bash
gitswitch doctor --json        # 1. is the setup sane?
gitswitch current --json       # 2. who am I in this repo?
gitswitch pin work            # 3. wrong? claim the repo, don't just switch globally
# ...agent commits...
gitswitch reauthor 5 --to work --yes   # 4. if it got out ahead of you
```

Step 3 over a plain `gitswitch work` matters for agents specifically: a global switch is a machine-wide side effect that affects every other terminal and every other repo. A pin is local, permanent, and can't be undone by something else switching identities in the background.

## The session scope

If you launch an agent from a terminal with `GIT_CONFIG_COUNT`/`GIT_CONFIG_KEY_0`/`GIT_CONFIG_VALUE_0` set, everything it spawns inherits that identity, and it outranks both the repo and your global config. gitswitch detects and reports it as `"scope": "session"` (`◆` in the UI). There's no gitswitch command to *create* one yet — but if your harness makes one, gitswitch will tell you the truth about it rather than confidently reporting the wrong account.

## Next

- **[Pins & Identity Awareness](/docs/routing/identity-awareness)** — the prevention half
- **[Scopes](/docs/accounts/scopes)** — global, repo, session
- **[CLI Reference](/docs/cli/commands)** — all flags
