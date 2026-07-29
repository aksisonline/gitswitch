---
title: Identity Awareness
description: How gitswitch learns which identity to use per repository
---

This page covers how gitswitch tracks per-repo identity usage, when it nudges you to switch, and how pinning overrides the auto-recommender.

## How learning works

Every time you enter a git repo (via the shell nudge hook), gitswitch calls `gitswitch record` to increment a usage counter for the active profile under that repo's key.

The repo key is resolved from `git remote get-url origin`. If no remote exists, it falls back to `git rev-parse --show-toplevel` (the absolute repo root path). This means two clones of the same repo at different paths are tracked separately.

History is stored in `~/.config/gitswitch/history.json`:

```json
{
  "repos": {
    "git@github.com:company/api.git": {
      "identities": { "work": 14, "personal": 1 },
      "last_used": "work"
    }
  }
}
```

## Recommendation threshold

gitswitch recommends a profile when all three conditions are true:

1. The top identity has **≥ 3** recorded entries for that repo
2. The top identity holds **≥ 60%** of all entries for that repo
3. The top identity differs from the currently active profile

If a repo has a pin set, the pin always wins — the threshold is ignored.

## Shell nudge

With [shell integration](/docs/features/shell) installed, entering a git repo triggers:

```bash
gitswitch recommend
```

If a recommendation is available, the shell hook displays:

```
gitswitch: this repo usually uses work <alice@company.com> — switch? [y/N]
```

Press `y` to switch, `n` or `Enter` to skip. The nudge is non-blocking.

## Pin a repo

Pinning sets the identity for one repo permanently, instead of asking you to switch globally every time you enter it:

```bash
cd ~/work/api
gitswitch pin work
```

```
Pinned 'work' to this repo — Alice Smith <alice@company.com> (local git config, global identity untouched)
```

The pin is written to the repo's own `.git/config` (`user.name`, `user.email`, `user.signingkey`/`gpg.format`, `core.sshCommand`) — git's local scope always beats the global one, so commits in this repo use that identity even while another profile is active everywhere else. It is the same thing as running `git config --local user.email …` by hand, just with your whole profile applied at once.

Everything else follows the repo automatically, so git and GitHub never disagree about who you are:

- **gh CLI** — pinning switches `gh` to that profile's account, and entering the repo later switches it back (the shell hook keeps it in step).
- **HTTPS pushes** — the credential helper serves this repo's tokens from the pinned account.
- **Prompt and `gitswitch current`** — both report the repo's identity while you are inside it, not the global one, and mark it as pinned (see [Telling the scopes apart](#telling-the-scopes-apart)).

Your globally active profile is the one thing left alone — pinning one repo should not change what every other repo does.

Since the repo already commits correctly, the shell hook does not nudge you when you enter it.

A repo whose local `user.email` merely repeats your global one is not treated as pinned — it overrides nothing, so gitswitch leaves the nudges on and shows no marker.

### Pinning from the TUI

Inside the TUI, `p` pins the highlighted profile to the repo you launched it from, and pressing `p` again on the pinned profile releases it. `enter` always means the same thing it ever did — change your **global** identity — so the two never get confused.

The key only appears in the footer when you are actually inside a git repo.

### Telling the scopes apart

Three things can decide who you commit as. Both the TUI and the shell prompt mark which one is in play:

| Marker | Scope | Where it lives |
|---|---|---|
| `✓` | global identity | `~/.gitconfig` — the profile you last switched to |
| `●` | pinned to this repo | the repo's own `.git/config` |
| `◉` | pinned to this repo **and** your global identity | both agree |
| `◆` | this terminal's session | environment variables (see below) |

In the shell prompt the same markers ride next to the nickname — `[work]` for global, `[work●]` in a pinned repo. Nothing is added when you are on your global identity, so a setup with no pins looks exactly as it always did.

The session marker (`◆`) appears when git config environment variables override both your repo and your global config — for example a terminal launched with `GIT_CONFIG_COUNT`/`GIT_CONFIG_KEY_0`/`GIT_CONFIG_VALUE_0` set. Anything started from that terminal, including AI coding agents, inherits it. First-class commands for this are on the roadmap; gitswitch already reports it correctly today.

### Repos you already configured by hand

If a repo already has its own `user.email` (set with `git config --local` long before gitswitch), run `pin` with no nickname to adopt it:

```bash
gitswitch pin
```

```
Adopted this repo's existing identity <alice@company.com>
Pinned 'work' to this repo — Alice Smith <alice@company.com> (local git config; global identity unchanged)
```

gitswitch matches the email to a stored profile and fills in the rest (signing key, SSH key, gh account). Repos like this are recognised even without a pin: the prompt, credential helper, and gh account all follow the local identity, and no nudge appears. Passing a nickname over an existing local identity tells you what it is replacing before it does so.

Remove the pin:

```bash
gitswitch unpin
```

```
Unpinned — this repo now uses the global identity
```

This removes those keys from the repo's local config. Both commands require you to be inside a git repo.

## Manually check what would be recommended

```bash
gitswitch recommend
```

- Exits `0` and prints `nickname\tname\temail` if a recommendation exists
- Exits `1` silently if already on the right identity, no history, or threshold not met

```bash
gitswitch recommend --path ~/work/api
```

Use `--path` to check a different directory without `cd`ing into it.

## Record usage manually

```bash
gitswitch record
```

Records the currently active profile for the current directory. The shell hook calls this automatically — you rarely need to run it by hand.

## Reset history

To clear all learned history:

```bash
rm ~/.config/gitswitch/history.json
```

To remove the pin for the current repo without clearing all history:

```bash
gitswitch unpin
```

## Troubleshooting

**Nudges are not appearing when I `cd` into repos**

1. Confirm shell integration is installed: `gitswitch install`
2. Reload your shell: `source ~/.zshrc` or open a new terminal
3. Confirm you're inside a git repo: `git rev-parse --git-dir`
4. Check usage counts: `cat ~/.config/gitswitch/history.json`

The threshold requires ≥ 3 entries with ≥ 60% share. If you've visited fewer than 3 times with a consistent identity, no nudge appears yet.

**Wrong profile is being suggested**

The recommendation reflects the most-used identity based on recorded history. If the data is stale or incorrect, edit or delete `~/.config/gitswitch/history.json` to clear it.

## Next steps

- [Shell Integration](/docs/features/shell)
- [CLI Reference — recommend](/docs/cli/commands)
- [Real-world Scenarios](/docs/scenarios/freelancer)
