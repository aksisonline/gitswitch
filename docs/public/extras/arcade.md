---
title: Arcade Mode
description: The one where a git tool has a high score
---
You found it. Congratulations.

```bash
gitswitch pacman
```

It's a light switch, not a party trick — run it once and gitswitch launches in arcade mode from then on. Run it again to go back to being a serious professional tool.

## What changes

- An intro animation, every launch, because once is not enough
- Double borders, pellet-separated tabs, `▶ INSERT COIN` instead of `▶ Start Setup`
- Arcade skins for all three tabs — `SHELL HOOK`, `SAVE FILE LOCATION`, and friends
- A **HIGH SCORE** line

Everything still works exactly as normal. You can pin repos and route pushes with a score counter on screen. Nobody will know.

## Scoring

| Action | Points |
|---|---|
| Switch identity | 200 |
| Add or edit an account | 500 |

The factory high score is **5000** — deliberately beatable. Your high score persists across restarts (it's `arcade_hi_score` in `~/.config/gitswitch/config.json`), which means yes, it is technically possible to grind git identity management.

## The markers change too

The [scope markers](/docs/accounts/scopes) get arcade counterparts:

| Normal | Arcade | Meaning |
|---|---|---|
| `✓` | `★` | your global identity |
| `●` | `◆` | pinned to this repo |
| `◉` | `✦` | pinned *and* global |
| `◆` | `▲` | this terminal's session |

## Themes

Not arcade-specific, but while you're here — `c` in the TUI cycles 12 color themes:

Default · Ocean · Sunset · Forest · Mono · Rose · Arctic · Gold · Violet · Ember · Matrix · Steel

Your choice is saved. `Matrix` is exactly what you think it is.

## Next

- **[Quick Start](/docs/get-started/quick-start)** — back to work
- **[Configuration](/docs/cli/config)** — where the score and theme live on disk
