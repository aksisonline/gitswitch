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

Pinning locks the recommendation for a repo, bypassing the usage threshold entirely:

```bash
cd ~/work/api
gitswitch pin work
```

```
Pinned 'work' to this repo
```

The pin is stored in `history.json` under a `"pinned"` field for that repo — no files are written to the repo itself.

Remove the pin:

```bash
gitswitch unpin
```

```
Unpinned — identity recommendation now based on usage history
```

Both commands require you to be inside a git repo.

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
