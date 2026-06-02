---
title: Quick Start
description: Add profiles and switch identities in under five minutes
---

This page covers adding your first profiles and switching between them.

## 1. First run

```bash
gitswitch
```

On first launch with no existing profiles, gitswitch reads your current `git config --global user.name` and `git config --global user.email` and saves them as a profile named `default`. Then it opens the TUI.

## 2. Add profiles

### Via CLI

```bash
# Git identity only
gitswitch add personal "Alice Smith" alice@gmail.com

# With SSH key
gitswitch add work "Alice Smith" alice@company.com \
  --ssh-key ~/.ssh/id_work

# Full setup: SSH key + GPG signing + GitHub CLI
gitswitch add corp "Alice Smith" alice@corp.com \
  --ssh-key ~/.ssh/id_corp \
  --sign-key ABCD1234EF567890 \
  --gh-user alice-corp
```

### Via TUI

Press `a` in the TUI to open the add-profile form. Navigate fields with `Tab`. Press `Enter` to save.

## 3. Switch profiles

### Quick switch (no UI)

```bash
gitswitch work
```

Output:

```
✓ Switched to 'work' — Alice Smith <alice@company.com>
```

### Interactive TUI

```bash
gitswitch
```

| Key | Action |
|-----|--------|
| `↑` / `↓` or `k` / `j` | Move cursor |
| `Enter` | Switch to selected profile |
| `a` | Add profile |
| `e` | Edit profile |
| `ctrl+d` (in edit form) | Delete profile |
| `c` | Cycle color theme |
| `?` | Show CLI quick reference |
| `q` / `ctrl+c` | Quit |

## 4. Check current identity

```bash
gitswitch current
```

```
work — Alice Smith <alice@company.com>
```

## 5. Set up shell integration

```bash
gitswitch install
```

This runs an interactive wizard that installs:

- **Prompt segment** — shows active identity when inside a git repo
- **Identity nudge** — on `cd` into a repo, prompts you to switch if a different identity is usually used there
- **Tab completion** — completes `gitswitch` commands and profile nicknames

After running, reload your shell:

```bash
source ~/.zshrc   # zsh
source ~/.bashrc  # bash
# or open a new terminal
```

## 6. Pin a repo to an identity

If you always want a specific profile in a repo, pin it:

```bash
cd ~/work/my-project
gitswitch pin work
```

```
Pinned 'work' to this repo
```

From now on, `gitswitch recommend` (used by the shell nudge hook) always returns `work` for that repo, regardless of usage history.

Remove the pin:

```bash
gitswitch unpin
```

```
Unpinned — identity recommendation now based on usage history
```

## List all profiles

```bash
gitswitch list
```

```
✓  personal        Alice Smith <alice@gmail.com>
   work            Alice Smith <alice@company.com>
   corp            Alice Smith <alice@corp.com>
```

The `✓` marks the currently active profile.

## Next steps

- [Commit Identity](/docs/features/commit-identity)
- [SSH Keys](/docs/features/ssh-keys)
- [Shell Integration](/docs/features/shell)
- [CLI Reference](/docs/cli/commands)
