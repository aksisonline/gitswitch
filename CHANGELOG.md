# Changelog

All notable changes to gitswitch are documented here.
Format: `[version] — date — summary`

---

## [Unreleased]

### Fixed
- **Mouse works on every screen** — clicks are now hit-tested against the actual rendered output instead of hardcoded row offsets, so the Utilities and Settings boxes, wizard buttons, update prompt, tips, shell-confirm dialog, and arcade mode (whose score line shifted everything by one row) all respond correctly. Item-box clicks previously used a 5-row grid for 4-row boxes and often focused the wrong item or did nothing.
- **Channel-locked update checks** — canary (beta) builds no longer get update notifications for stable releases, and vice versa. Moving between channels is explicit via `gitswitch stable` / `gitswitch beta`. Canary builds also keep receiving betas of newer targets (previously only same-base betas were offered).
- **Arcade select flash no longer wraps** — the flash row was rendered 2 columns wider than the panel and broke the layout mid-animation; it now matches the list rows (including the email/username toggle).
- **Edit form survives delete-cancel** — pressing `n` on the delete confirmation rebuilds the edit form with the profile's values instead of a stale/blank seed.
- **Onboarding welcome** — removed the decorative "Import config file" button (feature doesn't exist); "Skip for now" is now clickable and actually skips.

### Added
- **Tabs in arcade mode** — `1/2/3`, `tab`, and mouse clicks now switch tabs in arcade mode too, unlocking the arcade skins of the Utilities and Settings tabs (SHELL HOOK, SAVE FILE LOCATION, …) that were previously unreachable.
- **Hover focus** — moving the mouse over profile rows and Utilities/Settings boxes moves the cursor/focus, like a real menu. Scroll wheel also moves focus on the Settings tab and in wizard lists.
- **Persistent arcade high score** — HIGH SCORE survives restarts (`arcade_hi_score` in config.json); switching identities scores 200, adding/editing a profile 500. Factory high score is a beatable 5000.
- **Wizard mouse support** — import list rows toggle on click, Import/Skip buttons work, detect screen advances on click.

### Changed
- Active tab is now rendered as a highlighted chip; arcade tab strip uses pellet separators.

---

## [v0.2.0-beta.12] — 2026-06-25

### Fixed
- **Shell alias actually applied** — toggling or renaming the alias now auto-reinstalls the shell integration block in-place (uninstall + install) so the change takes effect on next shell reload. No manual reinstall required.
- `shell.Install` no longer overrides an empty alias with `"gs"` — empty now correctly skips the alias line.

### Added
- **Release notes in update screen** — when a newer version is available, the update prompt shows the GitHub release body (what's changed) fetched alongside the version check and cached to disk.

---

## [v0.2.0-beta.11] — 2026-06-25

### Added
- **Shell alias toggle** — Settings tab now shows a "Shell Alias" item with an on/off toggle (`enter` to toggle, `e` to rename). Default alias is `gs`.
- **Editable alias name** — Rename the short alias from within the TUI; value persisted to `config.json`.
- **Alias in shell snippets** — All shell integrations (zsh, bash, fish, p10k) now include `alias gs=gitswitch` (or the configured alias) in the installed block. Alias is omitted from the snippet when disabled.

### Fixed
- **Mouse clicks broken after editor** — Returning from an external editor (nano/vi) no longer kills mouse tracking; `tea.EnableMouseCellMotion` is re-sent on `editorDoneMsg`.

---

## [v0.2.0-beta.9] — 2026-06-25

### Added
- **Shell alias (initial)** — Editable alias field in Settings tab (`ShellAlias` in prefs, default `"gs"`).
- Shell snippets updated to include `alias gs=gitswitch`.

### Fixed
- **Mouse after editor** — Re-enable mouse cell motion on return from external config editor.

---

## [v0.2.0-beta.8] — 2026-06-24

### Fixed
- Derive panel top-Y from actual render height instead of hardcoded estimates.
- Stable status line height; disable utility nav with single item.
