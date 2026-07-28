# Changelog

All notable changes to gitswitch are documented here.
Format: `[version] — date — summary`

---

## [Unreleased]

### Fixed
- **`gitswitch install` now actually updates an existing setup** — it used to stop at "already installed" and leave your old shell hook in place, so the "shell integration updated" notice sent you to a command that did nothing. It now replaces the gitswitch block in your rc file (and only that block).
- **Repos that set the same email as your global config no longer suppress the nudge** — a repo with a local `user.email` identical to your global one overrides nothing, but gitswitch was treating it as pinned and silently skipping both the nudge and its usage learning.
- **Prompt footer no longer splits a key from its label** — the TUI footer measured itself against the panel's outer width instead of the usable width inside the border, so a long list could wrap between `1/2/3` and `tabs`.
- **Mouse works on every screen** — clicks are now hit-tested against the actual rendered output instead of hardcoded row offsets, so the Utilities and Settings boxes, wizard buttons, update prompt, tips, shell-confirm dialog, and arcade mode (whose score line shifted everything by one row) all respond correctly. Item-box clicks previously used a 5-row grid for 4-row boxes and often focused the wrong item or did nothing.
- **Channel-locked update checks** — canary (beta) builds no longer get update notifications for stable releases, and vice versa. Moving between channels is explicit via `gitswitch stable` / `gitswitch beta`. Canary builds also keep receiving betas of newer targets (previously only same-base betas were offered).
- **Arcade select flash no longer wraps** — the flash row was rendered 2 columns wider than the panel and broke the layout mid-animation; it now matches the list rows (including the email/username toggle).
- **Edit form survives delete-cancel** — pressing `n` on the delete confirmation rebuilds the edit form with the profile's values instead of a stale/blank seed.
- **Onboarding welcome** — removed the decorative "Import config file" button (feature doesn't exist); "Skip for now" is now clickable and actually skips.
- **Mouse now works on Add/Edit Profile and Delete confirm** — huh (the form library) has no mouse support of its own; clicking a field's title, description, or input now focuses it, matching the rest of the app. Delete confirm's "confirm delete" / "cancel" labels are now clickable too.

### Added
- **You can now see which identity is actually in charge** — the prompt and the TUI mark whether the current profile comes from your global config, from this repo, or from this terminal, so a pinned repo no longer looks identical to a global switch. Nothing changes for anyone without pins: no marker, no extra characters.
- **`p` pins a profile to the current repo from the TUI** — highlight a profile, press `p`, done; press `p` again on it to release the repo back to your global identity. `enter` still means "change my global identity", so the two never blur together.
- **Prompt shows the source at a glance** — `[work]` on your global identity, `[work●]` in a pinned repo, `[work◆]` in a terminal with its own identity. Run `gitswitch install` once to pick it up.
- **`gitswitch current` says where the identity came from** — e.g. `work — Alice W <alice@work.com>  (pinned to this repo)`, so you never have to read `git config` to find out why a commit was attributed the way it was.
- **Sign commits with an SSH key, no GPG needed** — pass your SSH key to `--sign-key` (`--sign-key ~/.ssh/id_ed25519.pub`) and gitswitch switches git to SSH signing for that profile. Previously it only ever set the key and assumed GPG, so an SSH key silently failed to sign (it now sets `gpg.format=ssh` alongside `user.signingkey`, clears it again for GPG key IDs, and adds git's required `key::` prefix for inline keys).
- **Pin a repo and it just stays right** — `gitswitch pin work` now writes the identity into that repo's own git config instead of only noting a preference, so commits there use it forever without switching your global identity or answering a prompt on every `cd`. Your gh account, HTTPS pushes, prompt, and `gitswitch current` all follow the repo while you're inside it; `gitswitch unpin` removes it again.
- **`gitswitch pin` with no name adopts what the repo already has** — for repos you configured by hand years ago (`git config --local user.email`), it matches the email to a stored profile and fills in the rest (signing key, SSH key, gh account). Repos like this are now recognised even without a pin, so gitswitch stops nudging you to switch in a repo that was already correct.
- **`gitswitch reauthor`** — rewrites author/committer identity on already-made commits to a stored profile in one command (`gitswitch reauthor <base> --to <nickname> [--from <old-email>] [--push]`), instead of hand-scripting `git rebase`/`git commit --amend`. `--from` scopes the rewrite to commits currently authored by a given email; `--push` force-pushes (`--force-with-lease`) after. Built so Claude/agents can fix pre-switch attribution in one call instead of a multi-step git dance.
- **Tabs in arcade mode** — `1/2/3`, `tab`, and mouse clicks now switch tabs in arcade mode too, unlocking the arcade skins of the Utilities and Settings tabs (SHELL HOOK, SAVE FILE LOCATION, …) that were previously unreachable.
- **Hover focus** — moving the mouse over profile rows and Utilities/Settings boxes moves the cursor/focus, like a real menu. Scroll wheel also moves focus on the Settings tab and in wizard lists.
- **Persistent arcade high score** — HIGH SCORE survives restarts (`arcade_hi_score` in config.json); switching identities scores 200, adding/editing a profile 500. Factory high score is a beatable 5000.
- **Wizard mouse support** — import list rows toggle on click, Import/Skip buttons work, detect screen advances on click. Detect and what's-new hints now say "or click" so the affordance isn't hidden.

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
