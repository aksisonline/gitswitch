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

All installations are idempotent — everything gitswitch adds lives between `# gitswitch shell integration` markers, and re-running `gitswitch install` replaces that block in place rather than appending a second copy. That is also how you pick up an improved prompt after upgrading, which is what the "shell integration updated — run `gitswitch install`" notice is asking you to do.

After installing, reload your shell:

```bash
source ~/.zshrc    # zsh
source ~/.bashrc   # bash
# or open a new terminal
```

## Prompt segment

The prompt segment shows the git identity in effect when you are inside a git repo. It calls:

```bash
gitswitch current --prompt
# work	141	●
```

The three fields are the profile nickname, the theme colour, and a **scope marker** saying where that identity comes from:

| Marker | Meaning |
|---|---|
| _(none)_ | your global identity — the profile you last switched to |
| `●` | this repo is [pinned](/docs/features/identity-awareness#pin-a-repo) and overrides the global identity |
| `◆` | this terminal's session overrides both |

So a glance at the prompt tells you whether a commit here will use the identity you think it will:

```
~/personal/blog  [personal] ❯
~/work/api       [work●] ❯
```

The segment is hidden when you're not inside a git repo. Starship uses `gitswitch current --short` instead, where the marker rides on the nickname (`work●`) because Starship renders the command output verbatim.

### Starship customization

The block added to `~/.config/starship.toml`:

```toml
[custom.gitswitch]
command = "gitswitch current --short"
when = "git rev-parse --git-dir > /dev/null 2>&1"
symbol = " "
style = "bold cyan"
format = "[$symbol($output)]($style) "
```

Edit `style` to change the color, or `symbol` to add an icon.

### oh-my-zsh customization

Edit the `__gitswitch_prompt` function in the gitswitch block of your `~/.zshrc` directly.

Re-running `gitswitch install` replaces that block with a fresh one, so keep customizations elsewhere in the file if you want them to survive an upgrade.

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

This removes the marker block from your rc file (or the oh-my-zsh plugin file), unregisters the HTTPS credential helper, and removes the `gh` CLI wrapper — all if they were installed. Reload your shell to complete removal.

## gh CLI isolation

`gh` (the GitHub CLI) only tracks one "active" account for your whole machine. If you switch identities in one terminal, every other terminal's bare `gh` commands (`gh pr create`, `gh issue list`, ...) silently start using that account too.

Toggle "gh CLI Isolation" on in the TUI's Utilities tab to fix this: it adds a `gh` shell function (its own marker block, independent of the rest of shell integration) that resolves the right account for your current repo before every `gh` call and passes it via `GH_TOKEN` for just that one command — the same resolution the HTTPS credential helper uses, so pushes and `gh` commands always agree. gh's global active-account file is never touched, so concurrent terminals never fight over it.

`gitswitch doctor` reports whether it's active. Reload your shell (or open a new terminal) after toggling.

## How the HTTPS credential helper is registered

Worth knowing if you also use `gh`, because git's rules here are surprising.

Git does not pick *one* credential helper. For a given URL it collects the values of `credential.helper` **and** every matching `credential.<host>.helper` into a single list, in the order they appear in your config files, then asks each one in turn until one returns a username and password. Two consequences:

- A host-specific entry does **not** override the generic one — it just joins the same list.
- A helper set to the **empty string** discards everything collected before it. `gh auth setup-git` writes exactly that, followed by its own helper, which is how it removes your keychain helper for github.com.

So gitswitch cannot simply add itself to `credential.helper`: `gh`'s empty entry would erase it. Instead `gitswitch install` makes gitswitch the first *live* entry in every helper list it finds, leaving everything else — including `gh`'s helper and your keychain — in place behind it:

```
[credential "https://github.com"]
	helper =                                     # gh's reset, preserved
	helper = !gitswitch credential               # gitswitch answers first
	helper = !/opt/homebrew/bin/gh auth git-credential   # fallback, untouched
```

gitswitch stays silent for any host or repo it cannot serve, so git falls straight through to the next helper. Nothing is removed, and `gitswitch uninstall` takes only gitswitch's own lines back out.

**Re-running `gh auth setup-git` (or an interactive `gh auth login`) undoes this.** gh rewrites that whole key with `--replace-all`, which drops gitswitch's line along with anything else there. `gitswitch doctor` detects it and names the helper answering ahead of gitswitch; `gitswitch install` puts things back.

Check the current state any time with:

```bash
gitswitch doctor
```

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
