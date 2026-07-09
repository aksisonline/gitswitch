# TUI v0.2.0 Implementation Notes

Session date: 2026-06-23  
Branch: release/v0.2.0

---

## What was built

Full redesign of the gitswitch TUI. Every screen was either rewritten or newly created. The stack stayed pure Go + charmbracelet (bubbletea v2, lipgloss, bubbles). Harmonica added for spring animations.

---

## New files

| File | Purpose |
|---|---|
| `internal/tui/animation.go` | 6 harmonica spring constructors (tab slide, cursor, input focus, boot logo, pac-man X, arcade switch bounce) |
| `internal/tui/components.go` | Shared UI helpers: `truncate`, `renderToggle`, `renderItemBox`, `itemInnerW`, `padTo`, `titleWithRight` |
| `internal/tui/view_accounts.go` | Accounts tab view + `viewTabHeader()` shared tab nav strip |
| `internal/tui/view_onboarding.go` | All 6 onboarding screens: `viewWhatsNew`, `viewWizardWelcome`, `viewWizardDetect`, `viewWizardImport`, `viewWizardAddMore`, `viewWizardDone` |
| `internal/tui/view_utilities.go` | Utilities tab — 4-line item boxes for shell integration, pre-commit, credential helper |
| `internal/tui/view_settings.go` | Settings tab — config path box, theme picker box |

---

## Modified files

### internal/storage/store.go
- Added `SplashSeen020 bool` and `ShellEnabled bool` to `Prefs`

### internal/tui/model.go
- New `State` consts: `StateWhatsNew` (10) through `StateWizardDone` (15)
- New `Model` fields: `tabIndex`, `wizardStep`, `detectedProfiles`, `importSelected`, `splashSeen020`, `shellEnabled`, `utilityFocus`, `settingsFocus`
- `New()` first-run detection: no profiles → wizard, has profiles + `!SplashSeen020` → what's new splash

### internal/tui/view.go
- `View()` routes new states
- Default case dispatches on `tabIndex` (0→accounts, 1→utilities, 2→settings)
- Profile item email truncated to prevent Y-offset drift; nickname column capped at 22 chars

### internal/tui/update.go
- `'a'` key → `StateWizardAddMore` (OAuth/manual choice) instead of raw form
- `↑/↓` tab-aware: routes to profile cursor, utilityFocus, or settingsFocus per tab
- `enter` tab-aware: switches profile (tab 0), toggles shell integration (tab 1)
- `1/2/3`, `tab`, `shift+tab` for tab navigation
- `←/→` cycles theme in Settings tab
- `panelTopY()` estimates panel vertical position per state for accurate mouse hit detection
- `handleMouse()`: tab header clicks, profile row clicks, wizard button clicks, utility/settings item clicks

---

## Onboarding flow

```
First run (no profiles)
  StateWizardWelcome → StateWizardDetect → (auto-scan ~/.gitconfig, gh CLI)
    → StateWizardImport (if configs found) → StateWizardAddMore
    → StateWizardAddMore (if nothing found)
      → OAuth (coming v0.2.1, shows status message)
      → Manual → StateAdd (existing form)
      → Done → StateList

Existing user (profiles exist, SplashSeen020=false)
  StateWhatsNew → (any key/click) → StateList

Subsequent runs
  StateList directly
```

---

## 3-tab layout

```
[ Accounts ]   Utilities   Settings

Tab 0 — Accounts   switch profiles, same keyboard as before
Tab 1 — Utilities  shell integration toggle (live), pre-commit + credential (coming v0.2.1)
Tab 2 — Settings   config path display, theme picker (← → or c)
```

Mouse: click tab label to switch, click profile row to switch profile, click item box to toggle.

---

## Item box design

Fixed 4-line height for all Utilities/Settings items:

```
  ┌──────────────────────────────────────────────────────┐
  │ Shell Integration                          ● on      │
  │ Auto-switch identity when you cd into a repo.        │
  └──────────────────────────────────────────────────────┘
```

Toggle states:
- ON:  `  ● on  ` — green background, black text
- OFF: ` ○ off  ` — dim background

Disabled items show a `[v0.2.1]` chip instead of a toggle.

Mouse hit zone: `item = (relY - 5) / 5` where relY is Y relative to panel top border. Valid when `(relY-5) % 5 < 4`. This works because every item is exactly 4 lines + 1-line blank prefix (5 total).

---

## Mouse architecture

`panelTopY()` estimates the panel's absolute Y screen position per-state, accounting for lipgloss vertical centering (`lipgloss.Place`). Estimates are ~±2 lines accurate — good enough for 4-line item boxes.

`handleMouse()` hit zones:
- relY=3: tab header (click left/mid/right third → switch tab)
- relY=7+: profile rows in Accounts tab (one row per profile, constant height due to truncation)
- relY=5+ in Utilities/Settings: item boxes via `(relY-5)/5` formula

---

## Harmonica springs (animation.go)

| Spring | freq | damp | Use |
|---|---|---|---|
| Tab slide | 4.0 | 0.85 | Panel slide when switching tabs |
| Cursor row | 10.0 | 0.75 | Profile list cursor |
| Input focus | 6.0 | 0.80 | Form field expand |
| Boot logo | 3.0 | 0.90 | Intro screen slide-in |
| Pac-man X | 2.5 | 1.00 | Critically damped, steady march |
| Arcade switch | 8.0 | 0.55 | Bouncy feedback on profile switch |

Springs are defined but not yet wired to tick animations — that's a follow-up task.

---

## Known gaps / follow-ups

- Harmonica springs constructed but not yet driving actual animations (tick wiring needed)
- OAuth device flow not implemented (shows "coming v0.2.1" message)
- Shell integration toggle is UI-only — actual shell hook install/uninstall not wired
- Settings config-path edit not implemented
- Mouse Y hit zones use estimated panel height; could be made exact by measuring rendered output
- Layout visual polish pass still needed (spacing, alignment, compact mode)
- Arcade mode wizard screens need styling review
- `StateNoProfiles` is effectively superseded by `StateWizardWelcome` but still exists
