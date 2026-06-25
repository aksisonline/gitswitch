# Changelog

All notable changes to gitswitch are documented here.
Format: `[version] — date — summary`

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
