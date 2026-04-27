# v0.2.0

## What's new

### Color themes

The interactive TUI now ships with 12 built-in color palettes. Press `c` while on the profile list to cycle through them. Your choice is remembered between sessions — no config edits needed.

**Available palettes:**

| # | Name | Character |
|---|------|-----------|
| 1 | Default | Original purple / green |
| 2 | Ocean | Deep blue / cyan |
| 3 | Sunset | Warm orange / gold |
| 4 | Forest | Dark green / lime |
| 5 | Mono | White / gray scale |
| 6 | Rose | Pink / magenta |
| 7 | Arctic | Cyan / ice blue |
| 8 | Gold | Yellow / amber |
| 9 | Violet | Magenta / lavender |
| 10 | Ember | Red / orange |
| 11 | Matrix | Bright green / terminal green |
| 12 | Steel | Slate blue / silver |

Theme preference is saved to `~/.config/gitswitch/config.json` and restored on next launch.

---

## Changes

- `c` key added to the profile list — cycles active color palette
- Theme selection persisted to `~/.config/gitswitch/config.json`
- Status bar confirms active theme name and index on each cycle

---

## No breaking changes

All existing profiles, CLI commands, and config files are unaffected.
