# Changelog

All notable changes to gitswitch are documented here.
Format: `[version] — date — summary`

---

## [Unreleased]

### Under the Hood
- **The install command is shorter** — `curl -fsSL https://get.gitswitch.dev | bash` replaces the long `raw.githubusercontent.com/.../install.sh` URL everywhere it's documented. Behaves identically; `get.gitswitch.dev` proxies the same script with the right `Content-Type`, it's not a redirect.
- **Agent setup is now a single URL** — the install script itself carries the full setup playbook as comments, so handing an agent `https://get.gitswitch.dev` is enough; it reads its own instructions from the file it just downloaded instead of needing a separate copy-paste block.

## [v0.4.0] — 2026-08-03

### What's New
- **One command sets up everything now** — run `gitswitch` on a fresh machine and it checks whether git and the GitHub CLI (`gh`) are installed, offers to install whichever's missing, then walks you through connecting your GitHub account. No separate setup command to remember first.
- **`gitswitch install` is now `gitswitch shell`** — same shell-integration wizard (prompt, HTTPS push routing, Session Isolation), just a name that matches what it actually does.
- **`gitswitch claude` is now `gitswitch skills`, and it's no longer just for Claude Code** — it tries [skills.sh](https://www.skills.sh) first, which installs the gitswitch skill correctly for whichever AI agent you're using (Claude Code, Cursor, Codex, and more), and falls back to installing it directly if you're offline.
- **A copy-paste setup prompt for AI agents** — a new Agent Aided Setup guide gives you one block to hand your coding agent, and it installs gitswitch, checks git/gh, and connects your account for you.

### Bug Fixes
- **Windows install command was broken** — `irm .../install.ps1 | iex` downloaded a binary that was never actually published, so every Windows install failed. Windows builds are now included in every release.
- **Signing in through gitswitch's own GitHub login didn't always let Session Isolation or HTTPS push routing use the new account** — those two features get their tokens from the `gh` CLI, but a fresh sign-in never told `gh` about the account, so it silently didn't work until you separately ran `gh auth login`. Logging in now registers the account with `gh` too, so it works right away.

### Under the Hood
- **The product name is now spelled consistently as "Gitswitch" everywhere** — the TUI header, arcade mode, and onboarding screens previously showed "Git-Switcher" (with a hyphen); they now match the CLI and docs.
- **`gitswitch setup` is gone** — checking and installing git/gh now happens automatically the first time you run `gitswitch` on a machine with no accounts yet. Its machine-readable manifest for AI agents is now part of `gitswitch doctor --json`.
- **Docs got reorganized** — Get Started, Accounts, Automatic Routing, AI, Guides, and CLI Reference now match how the site actually groups them, instead of one flat folder for everything.

## [v0.3.0] — 2026-07-30

### Fixed
- **Arcade mode's PAC-MAN intro only ever played once** — turning it on with `gitswitch pacman` showed the intro that one time, but every launch after that jumped straight to the profile list, even though the exit animation kept working. It now replays on every launch while arcade mode is on, not just the moment you toggle it.
- **The "What's New" screen never showed up after switching to a canary release** — canary versions only ever change by their last number (like `v0.2.2-beta.1` → `v0.2.2-beta.2`), and the screen used to only appear for bigger version jumps. `gitswitch beta`/`gitswitch stable` now schedule it for every canary switch, and canary installs show it for any version change, not just big ones.
- **The "run: gitswitch install" reminder kept nagging you forever, even after you'd already run it** — it used to compare version numbers, so any update at all (even ones that changed nothing about your shell setup) brought the reminder back. It now tracks whether your shell setup itself is actually out of date, so the reminder only shows up when there's really something to update.
- **The "Pre-commit Safety Net" placeholder is gone from the Utilities tab** — it had been sitting there disabled and unimplemented; the Utilities tab is now Shell Integration, Session Isolation, HTTPS Credential Helper — three real toggles, no "coming soon." — finishing an update from the "new version available" screen now returns you to the list with a clear success or error message, instead of leaving the screen sitting there with no feedback (thanks @shitcodebykaushik).
- **Editing a profile could silently disconnect it from GitHub** — saving changes from the edit screen now keeps your existing OAuth login (`token_ref`) in place, so you don't have to log in again after a small edit like fixing a typo in your name (thanks @shitcodebykaushik).
- **Logging in again could wipe your SSH signing key or GPG key** — running `gitswitch login` a second time on an account that already had a signing key or SSH key configured now keeps those keys instead of clearing them (thanks @shitcodebykaushik).
- **Manually editing your config file didn't show up in the app** — after using Settings → edit config to open `config.yaml` in your editor, gitswitch now reloads your profiles right away, so you see your changes without restarting (thanks @shitcodebykaushik).
- **The setup wizard now actually looks for accounts you already have** — it scans your logged-in `gh` accounts and the private keys in `~/.ssh/` and offers them for import, instead of only reading your global git name and email (thanks @shitcodebykaushik).
- **HTTPS Credential Helper in the Utilities tab was still "coming soon" long after it shipped** — it's been a working `gitswitch install`/`uninstall` feature for a while; the TUI just never caught up. It's now a real on/off toggle like Shell Integration, reflecting whether gitswitch is actually routing your pushes (not just registered — see the ordering fix below).
- **HTTPS pushes could still use the wrong account** — if you had ever run `gh auth setup-git` (or answered yes to "Authenticate Git with your GitHub credentials?" during `gh auth login`), gh's per-host helper was asked before gitswitch and gitswitch never got a say, so pushes went out as whichever account gh had active. gitswitch now registers itself as the first helper git asks for every host, leaving gh and your keychain in place behind it as fallbacks. Run `gitswitch install` once to repair an existing setup.
- **`gitswitch doctor` now checks HTTPS routing** — it reports whether pushes actually go through gitswitch and names the helper answering ahead of it if not, instead of leaving you to read `git config` yourself.
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
- **Two terminals, two GitHub accounts, no more fighting** — `gh` (the GitHub CLI) only remembers one "active" account for your whole machine, so switching identities in one terminal used to silently flip which account `gh pr create`/`gh issue create`/etc. used in every other terminal too. Turn on "Session Isolation" in the Utilities tab (on by default for new installs) and every `gh` command figures out the right account for the repo you're in, on its own, without ever touching that shared global setting.
- **Pinning a repo and Session Isolation are now one feature, not two** — a pin only takes effect while Session Isolation is on, since that's the only thing that actually keeps a repo's identity separate from the rest of your machine. `gitswitch pin` turns Session Isolation on automatically the first time you use it, so this is invisible unless you'd deliberately turned it off — in that case, `gitswitch current`, the TUI, and the shell prompt now clearly show your global identity is active and that the pin is sitting there inactive, instead of quietly trusting a local git config override that no longer means what it used to.
- **`gitswitch install` now offers Session Isolation as a setup step**, on by default, so new installs don't need a separate trip to the Utilities tab to get per-repo `gh` account isolation and working pins.
- **`gitswitch pacman` is now a light switch, not a one-time party trick** — running it once turns arcade mode on for every future launch of gitswitch (not just that session); run it again to turn it back off.
- **Colorful `pin`/`unpin`/`current`/`doctor`/`list` and friends** — plain CLI commands now use the same green/yellow/red palette as the TUI and setup wizard, `gitswitch --help` groups commands by what they're for instead of one long alphabetical list, and `gitswitch current --json` / `gitswitch list --json` let scripts read your identity state without parsing prose (colors stay off automatically when output isn't going to a real terminal).
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
