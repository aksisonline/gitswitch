---
title: Troubleshooting
description: Start with gitswitch doctor, then find your symptom below
---

## Start here

```bash
gitswitch doctor
```

```
  ✓  git 2.50.1
  ✓  gh  2.95.0
  ✓  HTTPS pushes routed by gitswitch
  ✓  Session Isolation active (bare `gh` commands resolve per-repo)
```

Four checks: is git there, is `gh` there, do HTTPS pushes actually go through gitswitch, and is Session Isolation live. `--json` for scripts and agents.

Most problems below are one of those four saying no.

## "My commit went out as the wrong person"

```bash
gitswitch current
```

That tells you who you *are* here and **why** — global, pinned to this repo, or a terminal override. If it's not who you expected, [Scopes](/docs/concepts/scopes) explains who's winning.

Not pushed yet:

```bash
gitswitch work
git commit --amend --reset-author --no-edit
```

Already pushed, or more than one commit:

```bash
gitswitch reauthor 3 --to work --push
```

See [AI Coding Agents](/docs/features/ai-agents#gitswitch-reauthor-fix-commits-that-already-exist) for the full `reauthor` reference — it's the same command whether a human or an agent made the mess.

Then stop it recurring:

```bash
gitswitch pin work
```

## "My push went out as the wrong GitHub account"

For **HTTPS** remotes, run `gitswitch doctor`. If you see this:

```
  ✗  HTTPS pushes answered by another helper before gitswitch:
       credential.https://github.com.helper → !/opt/homebrew/bin/gh auth git-credential
```

then another credential helper is being asked before gitswitch and gitswitch never gets a say. Usually because `gh auth setup-git` ran (an interactive `gh auth login` does it too).

```bash
gitswitch install
```

puts gitswitch back in front without removing anything else. Details: [HTTPS Push Routing](/docs/features/https).

For **SSH** remotes, check the key actually being offered:

```bash
git config --global core.sshCommand
# must contain: -o IdentitiesOnly=yes
```

If it's missing or wrong, re-switch the profile (`gitswitch work`). See [SSH Keys](/docs/features/ssh-keys).

## "`gh pr create` used the wrong account"

`gh` tracks one active account for your whole machine, so another terminal can flip it under you. That's what [Session Isolation](/docs/features/session-isolation) is for:

```bash
gitswitch doctor    # is it active?
```

If it's off, turn it on in the TUI's **Utilities** tab (or re-run `gitswitch install`) and reload your shell — the wrapper is a shell function, so terminals opened before you enabled it won't have it:

```bash
source ~/.zshrc    # zsh
source ~/.bashrc   # bash
```

> `gh auth status` will still show your old global account. That's expected and not a bug — isolation deliberately never touches gh's global account file, it overrides the token for individual calls. Judge it by which account a `gh` command actually acted as.

## "My pin isn't doing anything"

Pins require Session Isolation. If it's off, the pin is stored but inactive, and gitswitch tells you so:

```
work — Alice Smith <alice@company.com>
  (pinned to 'acme' — inactive, Session Isolation is off)
```

```bash
gitswitch pin acme    # re-pinning turns Session Isolation back on
```

Also: a repo whose local `user.email` merely repeats your global one isn't an override, so it reports as `global` with no marker. That's correct — it overrides nothing.

## "The prompt doesn't show anything"

1. Are you inside a git repo? The segment hides outside one, on purpose. (`git rev-parse --git-dir`)
2. Did you reload your shell after `gitswitch install`?
3. Does the underlying command work? `gitswitch current --short`
4. Using Starship? The block goes in `~/.config/starship.toml`, not your rc file.

## "I'm not getting nudged"

The nudge only fires once there's a clear pattern: the top identity for that repo needs **≥ 3** recorded visits and **≥ 60%** share, and it has to differ from your current identity.

```bash
gitswitch recommend         # what would it say right now?
cat ~/.config/gitswitch/history.json
```

Exits `1` silently when there's nothing to recommend — that's normal, the shell hook relies on it.

Pinned repos are never nudged, since they already commit correctly.

## "Tab completion isn't working"

Needs zsh 5.0+, bash 4.0+, or fish 3.0+ (macOS ships bash 3.2 — `brew install bash` if you're on the system one). Then re-run `gitswitch install` and reload.

## "I told it to install, and it said 'already installed'"

Fixed in v0.3.0 — `gitswitch install` now replaces the gitswitch block in your rc file rather than stopping short. If you're seeing this, you're on an older build:

```bash
gitswitch upgrade
```

## "gitswitch can't read my profiles"

Your accounts live in `~/.config/gitswitch/config.yaml`. If you hand-edited it, check the YAML:

```bash
cat ~/.config/gitswitch/config.yaml
```

Upgrading from an older version migrated `profiles.json` → `config.yaml` and left `profiles.json.v1.bak` behind. If something went wrong in that migration, that backup is your safety net.

Corrupted history is handled automatically — gitswitch backs it up to `history.json.bak` and starts fresh.

Permissions, after restoring from a backup:

```bash
chmod 700 ~/.config/gitswitch
chmod 600 ~/.config/gitswitch/*
```

## Start completely over

```bash
gitswitch uninstall            # shell integration, HTTPS routing, gh wrapper
rm -rf ~/.config/gitswitch     # accounts, history, preferences
gitswitch install
```

Tokens in your OS keychain aren't removed by that — delete `gitswitch:*` entries in Keychain Access / your secret store if you want those gone too.

## Still stuck

- [Open an issue](https://github.com/aksisonline/gitswitch/issues) — include `gitswitch doctor --json` output
- [CLI Reference](/docs/cli/commands) — every command and flag
