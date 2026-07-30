---
title: Session Isolation
description: Give every terminal its own GitHub CLI account and working repo pins, without fighting over gh's single global active account
---

> **Canary channel.** This feature ships in gitswitch's canary builds. Run `gitswitch beta` to switch your install to canary (`gitswitch stable` switches back at any time).

## The problem

`gh` (the GitHub CLI) only tracks one "active" account for your whole machine. Switching your gitswitch identity in one terminal silently flips which account every other open terminal's bare `gh` commands (`gh pr create`, `gh issue list`, `gh repo view`, ...) use too — even in a completely unrelated repo.

## What it does

Turning on Session Isolation installs a `gh` shell function that wraps the real `gh` binary. Before every `gh` invocation, the wrapper resolves the right GitHub account for the repo you're in — the same resolution the HTTPS credential helper uses for pushes — and passes that account's token via `GH_TOKEN` for just that one command. gh's own global "active account" file is never read or written, so concurrent terminals never fight over it.

Session Isolation also gates [repo pins](/docs/features/identity-awareness#pin-a-repo): a pin only takes effect while it's on, since that's what actually keeps a pinned repo's identity separate from the rest of your machine.

## Enable it

`gitswitch install` offers Session Isolation as a setup step, on by default for new installs. To turn it on later, open the TUI's **Utilities** tab:

```bash
gitswitch
```

Toggle **Session Isolation** on, then reload your shell (or open a new terminal) to pick up the wrapper.

Pinning a repo turns Session Isolation on automatically if it was off — `gitswitch pin work` just works, even before you've visited the Utilities tab.

## How account resolution works

For each `gh` call, the wrapper works out which profile applies to the current repo — your global identity, a [pinned repo](/docs/features/identity-awareness#pin-a-repo), or a session override — and fetches that profile's stored token. If no account resolves (no gh-linked profile, or you're outside a git repo), the wrapper falls through to plain `gh` using your existing global account — no change in behavior.

## If it's off

A pin set while Session Isolation is off doesn't disappear, but it doesn't apply either. `gitswitch current`, the TUI, and the shell prompt all show your global identity is active and note that the pin is present but inactive, instead of quietly trusting a local git config override that no longer means what it used to.

## Turn it off

Toggle Session Isolation off in the Utilities tab, or run `gitswitch uninstall` to remove all shell integration at once. Either way the wrapper function is removed immediately, bare `gh` behaves exactly like it did before, and any repo pins go inactive as described above.

## Troubleshooting

**`gh` commands still use the wrong account after enabling**

The wrapper is a shell function — a terminal opened before you turned the toggle on won't have it until its rc file is re-sourced:

```bash
source ~/.zshrc    # zsh
source ~/.bashrc   # bash
# or open a new terminal
```

**Isolation seems to do nothing**

The wrapper overrides the *token* used for a given call — it doesn't change what `gh auth status` reports, since that command reads gh's own global account list, which isolation intentionally never touches. Judge it by which account a `gh` command actually acted as, not by `gh auth status`.

**A repo pin isn't being used**

Check whether Session Isolation is on — a pin set while it's off is stored but inactive. Re-run `gitswitch pin <nickname>` to turn isolation on and re-apply the pin.

## Next steps

- [GitHub Account Sync](/docs/features/github-sync) — the `--gh-user` flag and `gh auth switch` on profile switch
- [Identity Awareness](/docs/features/identity-awareness) — pins and per-repo identity resolution
- [Shell Integration](/docs/features/shell) — everything else `gitswitch install` sets up
