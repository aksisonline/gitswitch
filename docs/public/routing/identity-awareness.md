---
title: Pins & Identity Awareness
description: Claim a repo for one account, let gitswitch learn your habits, and block wrong-author commits
---
Three ways to stop using the wrong account in a repo. Pinning is the deliberate one. Learning is the one that works even when you never got around to it. And the pre-commit guard catches mistakes before a commit ever enters history.

## Pin a repo

```bash
cd ~/work/api
gitswitch pin work
```

```
✓ Pinned 'work' to this repo — Alice Smith <alice@company.com> (local git config; global identity unchanged)
```

(If Session Isolation happened to be off, it says `Session Isolation was off — turning it on (pins need it)` first, then pins.)

Done forever. Every commit in that repo uses `work`, no nudge, no prompt, no remembering.

What that actually writes into the repo's own `.git/config`:

```ini
[user]
	name = Alice Smith
	email = alice@company.com
	signingkey = ~/.ssh/id_work.pub
[gpg]
	format = ssh
[core]
	sshCommand = ssh -i /Users/alice/.ssh/id_work -o IdentitiesOnly=yes
```

Which is just `git config --local` five times over — except gitswitch applies your whole account in one go, and git's local scope beats global, so this repo keeps its identity even while you're switched to something else everywhere else.

Everything else follows the repo too, so git and GitHub can't disagree about who you are:

- **`gh` commands** — resolve to that account inside this repo ([Session Isolation](/docs/routing/session-isolation))
- **HTTPS pushes** — use that account's token ([HTTPS routing](/docs/routing/https))
- **Prompt and `gitswitch current`** — report the repo's identity while you're in it, marked `●`

Your globally active account is the one thing left alone. Pinning one repo shouldn't change what every other repo does.

> **Pins need Session Isolation.** It's on by default, and `gitswitch pin` turns it on automatically the first time you pin something. If you've deliberately turned it off, the pin is stored but inactive — and gitswitch [says so](/docs/accounts/scopes#pins-need-session-isolation) rather than pretending.

### From the TUI

Highlight an account, press `p`. Press `p` again on the pinned one to release it. `enter` still means "change my **global** identity", so the two never blur.

The key only appears in the footer when you're actually inside a git repo.

### Repos you configured by hand years ago

Run `pin` with no name and gitswitch adopts what's already there:

```bash
gitswitch pin
```

```
Adopted this repo's existing identity <alice@company.com>
Pinned 'work' to this repo — Alice Smith <alice@company.com> (local git config; global identity unchanged)
```

It matches the existing `user.email` to one of your accounts and fills in the rest — signing key, SSH key, `gh` account. Repos like this are recognised even *without* a pin, so gitswitch won't nag you to switch in a repo that was already correct.

Passing a nickname over an existing local identity tells you what it's about to replace before it does.

### Unpin

```bash
gitswitch unpin
```

```
✓ Unpinned — this repo now uses the global identity
```

Removes those keys from the repo's local config. Both commands need to be run inside a git repo.

## Pre-commit Guard

Even before you pin a repo, the pre-commit guard blocks wrong-author commits before they enter git history. It is enabled by default during `gitswitch shell` and first-run onboarding, or you can toggle it anytime in the TUI (**Utilities → Pre-commit Guard**) or via the CLI:

```bash
gitswitch guard enable          # install global pre-commit hook
gitswitch guard status          # check hook + whether this repo passes (--json supported)
gitswitch guard disable         # remove hook and restore any previous core.hooksPath
```

It installs a global `core.hooksPath` wrapper (`~/.config/gitswitch/hooks/pre-commit`) that chains into any previous global hooks directory and each repo's own `.git/hooks/pre-commit`, so Husky, lefthook, and `pre-commit` keep running normally.

A commit is refused when:
- No `user.email` is configured in git
- The effective `user.email` doesn't match any stored gitswitch profile
- The repo is pinned to a different profile than the effective git identity
- An unpinned repo has a learned recommendation (≥ 3 visits, ≥ 60% consistency) for a different profile

## IDE & Xcode Integration (early)

Both integrations are early: they work today, but expect rough edges and changes between releases.

- **Xcode (macOS)**: `gitswitch xcode enable` keeps Xcode's default GitHub account and per-repo HTTPS remotes in sync automatically whenever you run `gitswitch switch` or `gitswitch pin`. Use `gitswitch xcode status` to see matched accounts or `gitswitch xcode add <nickname>` for one-paste token setup.
- **VS Code & JetBrains**: Native status-bar extensions (in `extensions/vscode` and `extensions/jetbrains`) show the active profile, scope (`●` pinned / `◆` session), and pre-commit guard status right in the editor status bar, with one-click switching, pinning, and guard toggling.

## Or let it learn

You won't pin everything. So gitswitch quietly keeps score.

Every time you enter a git repo, the shell hook records which account was active (`gitswitch record`). Once a pattern is obvious, it speaks up:

```
gitswitch: this repo usually uses work <alice@company.com> — switch? [y/N]
```

One keypress. Defaults to no. Never blocks your prompt.

### When it decides to say something

All three must be true:

1. The top account has **≥ 3** recorded visits for this repo
2. It holds **≥ 60%** of that repo's visits
3. It's different from your current account

A pin skips all of this — a pinned repo is never nudged, because it's already right.

### Where the counting lives

`~/.config/gitswitch/history.json`:

```json
{
  "repos": {
    "git@github.com:company/api.git": {
      "identities": { "work": 14, "personal": 1 },
      "last_used": "work",
      "pinned": "work"
    }
  }
}
```

The key is your `git remote get-url origin`, falling back to the absolute repo root path if there's no remote. Two clones of the same repo in different folders are tracked separately.

### Check or reset it

```bash
gitswitch recommend            # what would it suggest right now?
gitswitch recommend --path ~/work/api    # ...for somewhere else
```

Exits `0` and prints `nickname<TAB>name<TAB>email` when there's a recommendation; exits `1` silently when there isn't. That silence is deliberate — the shell hook depends on it.

```bash
gitswitch record               # record manually (the hook does this for you)
rm ~/.config/gitswitch/history.json      # forget everything
```

## Troubleshooting

**Nudges aren't appearing**

1. `gitswitch shell`, then reload your shell
2. Are you inside a git repo? (`git rev-parse --git-dir`)
3. Have you been there 3+ times with a consistent account? (`gitswitch recommend`)

**It's suggesting the wrong account**

The recommendation is only ever a reflection of recorded history. Fix it at the source — pin the repo (`gitswitch pin <name>`), which overrides the counting entirely.

## Next

- **[Scopes](/docs/accounts/scopes)** — global vs. repo vs. terminal
