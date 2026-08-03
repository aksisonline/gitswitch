---
title: Pins & Identity Awareness
description: Claim a repo for one account, or let gitswitch learn your habits and nudge you
---

Two ways to stop using the wrong account in a repo. Pinning is the deliberate one. Learning is the one that works even when you never got around to it.

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

- **`gh` commands** — resolve to that account inside this repo ([Session Isolation](/docs/features/session-isolation))
- **HTTPS pushes** — use that account's token ([HTTPS routing](/docs/features/https))
- **Prompt and `gitswitch current`** — report the repo's identity while you're in it, marked `●`

Your globally active account is the one thing left alone. Pinning one repo shouldn't change what every other repo does.

> **Pins need Session Isolation.** It's on by default, and `gitswitch pin` turns it on automatically the first time you pin something. If you've deliberately turned it off, the pin is stored but inactive — and gitswitch [says so](/docs/concepts/scopes#pins-need-session-isolation) rather than pretending.

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

1. `gitswitch install`, then reload your shell
2. Are you inside a git repo? (`git rev-parse --git-dir`)
3. Have you been there 3+ times with a consistent account? (`gitswitch recommend`)

**It's suggesting the wrong account**

The recommendation is only ever a reflection of recorded history. Fix it at the source — pin the repo (`gitswitch pin <name>`), which overrides the counting entirely.

## Next

- **[Scopes](/docs/concepts/scopes)** — global vs. repo vs. terminal
- **[Session Isolation](/docs/features/session-isolation)** — what makes pins take effect
- **[Guides](/docs/scenarios/multi-github)** — pinning in a real two-account setup
