---
title: Shell Integration
description: How to install the prompt segment, identity nudge, and tab completion
---

This page covers what `gitswitch install` sets up, which frameworks it supports, and how to customize the prompt segment.

## What gets installed

```bash
gitswitch install
```

Runs an interactive wizard that installs two things:

1. **Shell integration** — prompt segment, identity nudge on `cd`, tab completion
2. **HTTPS credential helper** — routes HTTPS git operations through the active profile's `gh` account (optional)

To skip the wizard and accept defaults:

```bash
gitswitch install --yes
```

To target a specific shell:

```bash
gitswitch install --shell zsh
gitswitch install --shell bash
gitswitch install --shell fish
```

## Framework detection

gitswitch detects your prompt framework and installs accordingly:

| Framework | What gets written |
|-----------|-------------------|
| Starship | `[custom.gitswitch]` block appended to `~/.config/starship.toml` |
| oh-my-zsh | Plugin created at `~/.oh-my-zsh/custom/plugins/gitswitch/gitswitch.plugin.zsh` |
| Powerlevel10k | Segment function dropped to rc file; manual step printed for `~/.p10k.zsh` |
| Raw zsh / bash / fish | Prompt function, nudge hook, and completion snippet appended to rc file |

All installations are idempotent — a `# gitswitch shell integration` marker is used to skip re-installation on subsequent runs.

After installing, reload your shell:

```bash
source ~/.zshrc    # zsh
source ~/.bashrc   # bash
# or open a new terminal
```

## Prompt segment

The prompt segment shows your active git identity when inside a git repo. It calls:

```bash
gitswitch current --short
# work	alice@company.com
```

Example prompt (with Starship):

```
~/work/api [work: alice@company.com] ❯
```

The segment is hidden when you're not inside a git repo.

### Starship customization

The block added to `~/.config/starship.toml`:

```toml
[custom.gitswitch]
command = "gitswitch current --short"
when = "git rev-parse --git-dir 2>/dev/null"
symbol = ""
format = "[$output]($style) "
style = "bold 141"
```

Edit `style` to change the color, or `symbol` to add an icon.

### oh-my-zsh customization

Edit `~/.oh-my-zsh/custom/plugins/gitswitch/gitswitch.plugin.zsh` directly.

## Identity nudge

When you `cd` into a git repo, the shell hook runs `gitswitch recommend`. If a recommendation exists (based on usage history or a pin), it prints:

```
gitswitch: this repo usually uses work <alice@company.com> — switch? [y/N]
```

Press `y` to switch immediately, or `n` / `Enter` to skip. The nudge does not block your prompt.

See [Identity Awareness](/docs/features/identity-awareness) for details on how recommendations are learned.

## Tab completion

After installation, tab completion is active for all `gitswitch` commands and your profile nicknames:

```
$ gitswitch <Tab>
add       claude    current   init      install   list      pin       record    recommend remove    switch    uninstall unpin     upgrade   version

$ gitswitch sw<Tab>
switch

$ gitswitch w<Tab>
work
```

## Uninstall

To remove all shell integration written by `gitswitch install`:

```bash
gitswitch uninstall
```

This removes the marker block from your rc file (or the oh-my-zsh plugin file) and unregisters the HTTPS credential helper if it was installed. Reload your shell to complete removal.

## Troubleshooting

**Prompt segment not appearing**

- Confirm `gitswitch install` completed without errors
- Reload your shell or open a new terminal
- Confirm you're inside a git repo: `git rev-parse --git-dir`
- Test manually: `gitswitch current --short`

**Nudges not appearing**

- Shell integration must be installed: `gitswitch install`
- The identity threshold needs ≥ 3 visits with ≥ 60% consistency before a nudge fires
- Check history: `cat ~/.config/gitswitch/history.json`

**Completion not working**

Verify the shell version:

```bash
zsh --version   # 5.0 or later
bash --version  # 4.0 or later
fish --version  # 3.0 or later
```

Reinstall if needed: `gitswitch install`

## Next steps

- [Identity Awareness](/docs/features/identity-awareness)
- [Quick Start](/docs/quick-start)
- [CLI Reference](/docs/cli/commands)
