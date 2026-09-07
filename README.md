<p align="center">
  <img src="https://gitswitch.dev/og-image.png" alt="gitswitch — Git, done right">
</p>

<h1 align="center">gitswitch</h1>

<p align="center">
  <a href="https://github.com/aksisonline/gitswitch/releases"><img src="https://shieldcn.dev/github/release/aksisonline/gitswitch.svg" alt="Latest release"></a>
  <a href="LICENSE"><img src="https://shieldcn.dev/github/license/aksisonline/gitswitch.svg" alt="License"></a>
</p>

Manage your git identities without touching config files. Switch the name, email, SSH key, and GitHub account used for your commits — instantly. Docs live at the [site](https://gitswitch.dev); this repo's `docs/` is their source.

**New to git entirely?** Run `gitswitch` — the first run checks git and `gh` are installed (offers to install them), then walks you through connecting a GitHub account and filling in name, email, and keys. You don't need a second account for this to be worth it.

---

## Why not just use `gh auth switch`?

`gh` (GitHub CLI) manages **API credentials** — the OAuth tokens that let `gh pr create`, `gh issue list`, etc. work. Switching with `gh auth switch` changes which account the `gh` CLI operates under.

It does **not** change your **git commit identity**.

Your commit identity — the name and email baked into every `git commit` — comes from:

```
git config --global user.name  "Your Name"
git config --global user.email "you@example.com"
```

These are independent. You can have `gh` authenticated as your work account while every commit shows your personal email. The two tools solve different problems:

| | `gh auth switch` | `gitswitch` |
|---|---|---|
| **Controls** | GitHub API tokens | `git config user.name/email` |
| **Affects** | `gh` CLI commands | Commit author identity |
| **SSH key switching** | No | Yes — sets `core.sshCommand` |
| **GPG signing key** | No | Yes — sets `user.signingkey` |
| **Works with GitLab/Bitbucket** | No | Yes — any git remote |
| **Talks to GitHub** | Yes (OAuth) | Optional — `gh auth switch` is best-effort |

**The typical problem:** You push a commit to your company repo, only to see it attributed to your personal email. `gh auth switch` would not have helped. `gitswitch` would.

---

## Install

**Homebrew** (recommended):
```bash
brew install aksisonline/tap/gitswitch
```

**Curl** (one-liner):
```bash
curl -fsSL https://get.gitswitch.dev | bash
```

Or **build from source**:
```bash
git clone https://github.com/aksisonline/gitswitch
cd gitswitch
make install
```

---

## Usage

### Interactive TUI

```bash
gitswitch
```

Opens a full terminal UI. First run auto-imports your existing `git config` as a `default` profile.

```
╭──────────────────────────────────────────────────────────────╮
│  ◆  Gitswitch   identity manager for git                     │
│                                                              │
│   Accounts    Utilities   Settings                           │
│                                                              │
│  Current  aks  ·  gh:aksisonline                             │
│                                                              │
│ ❯ ✓ aks             gh:aksisonline                           │
│   · work            gh:abhiramkanna-edirq                    │
│                                                              │
│ ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄   │
│  ↑/↓ navigate  ·  enter switch  ·  p pin to repo             │
│  a add  ·  e edit  ·  ? cli tips  ·  q quit  ·  1/2/3 tabs   │
╰──────────────────────────────────────────────────────────────╯
```

**Keys:**
| Key | Action |
|-----|--------|
| `↑`/`↓` (or `k`/`j`) | Navigate profiles |
| `enter` | Switch to selected profile |
| `p` | Pin the selected profile to the current repo |
| `a` | Add a profile |
| `e` | Edit the selected profile |
| `v` | Toggle emails ↔ usernames |
| `?` | CLI quick reference |
| `q` | Quit |
| `1`/`2`/`3` | Accounts / Utilities / Settings tabs |

---

### Quick switch (no UI)

```bash
gitswitch work
```

Switches immediately and exits. Useful in scripts or when you already know the profile.

### Other CLI commands

```bash
gitswitch current          # show active profile
gitswitch list             # list all profiles
gitswitch add work "Jane Doe" jane@company.com --ssh-key ~/.ssh/id_work
gitswitch switch work      # switch by name
gitswitch remove work      # remove a profile
gitswitch init             # re-import current git config
gitswitch version          # show version and check for updates
gitswitch upgrade          # upgrade to latest release
```

See **[docs/cli.md](docs/cli.md)** for full flag reference and examples.

---

## Shell integration

One command sets up your shell with three features at once:

```bash
gitswitch shell
```

- **Prompt segment** — shows the active git identity in your shell prompt whenever you're inside a repo
- **Identity nudge** — when you `cd` into a repo, suggests the identity you usually use there (one keypress to switch)
- **Tab completion** — completes `gitswitch` commands and profile nicknames

Detects your prompt framework automatically:

| Framework | What happens |
|-----------|-------------|
| Starship | Adds `[custom.gitswitch]` block to `~/.config/starship.toml` |
| oh-my-zsh | Creates a plugin at `~/.oh-my-zsh/custom/plugins/gitswitch/` |
| Powerlevel10k | Drops the segment function, prints manual step for `~/.p10k.zsh` |
| Raw zsh/bash/fish | Appends directly to your rc file |

Idempotent — safe to run multiple times. Reload your shell after running:

```bash
source ~/.zshrc   # or open a new terminal
```

---

## Identity awareness

gitswitch learns which identity you use in each repo and nudges you when something looks off.

### How it works

Each time you enter a git repo (new terminal or `cd`), gitswitch silently records which identity is currently active. Once a pattern is clear — the top identity has **≥3 entries** and **≥60% share** — it nudges:

```
gitswitch: this repo usually uses work <alice@company.com> — switch? [y/N]
```

One keypress. Defaults to N. Non-blocking.

Usage history is stored at `~/.config/gitswitch/history.json`.

### Pin a permanent identity to a repo

Requires Session Isolation (on by default via `gitswitch shell`; pinning turns it on automatically if it's off). For repos where you always want a specific identity, skip the learned-count logic entirely:

```bash
gitswitch pin work    # always recommend 'work' for this repo
gitswitch unpin       # remove the pin, fall back to auto-recommendation
```

The pin takes permanent priority over usage counts.

---

## AI agent skill

Install the gitswitch skill so your AI coding agent can detect and fix git identity problems automatically — Claude Code, Cursor, Codex, and anything else that speaks the Agent Skills format:

```bash
gitswitch skills                   # skills.sh first, falls back offline for Claude Code + .agents/skills
gitswitch skills --scope project   # this project only
gitswitch skills --offline         # skip skills.sh, go straight to the offline installer
```

The skill is embedded in the binary — the offline path needs no download, always matches your installed version. After installing, reload your agent or open a new session to activate.

---

## Profile fields

| Field | Git config key | Description |
|-------|---------------|-------------|
| Nickname | — | Label shown in the list. Not written to git config. |
| User Name | `user.name` | Author name on commits. |
| Email | `user.email` | Author email on commits. |
| GPG Signing Key | `user.signingkey` | Optional. For signed commits. |
| SSH Key Path | `core.sshCommand` | Optional. Path to SSH private key (e.g. `~/.ssh/id_work`). Sets `ssh -i <key> -o IdentitiesOnly=yes` to force that key and prevent SSH agent fallback. |
| GitHub Username | — | Optional. Runs `gh auth switch --user <username>` on switch. Fails gracefully if `gh` is not installed. |

Profiles are stored at `~/.config/gitswitch/profiles.json`. UI preferences (color theme) are stored at `~/.config/gitswitch/config.json`.

---

## How switching works

On switch, `gitswitch` runs (in order):

1. `git config --global user.name "<UserName>"`
2. `git config --global user.email "<Email>"`
3. `git config --global user.signingkey "<GPGKey>"` — if set
4. `git config --global core.sshCommand "ssh -i <SSHKey> -o IdentitiesOnly=yes"` — if set
5. `gh auth switch --user <GHUser>` — if set, **warning only** on failure (git config already applied)

Step 5 is best-effort. If `gh` is not installed or the account isn't logged in, git config still switches correctly.

---

## Common scenarios

**Contractor with multiple clients**
```bash
gitswitch add clienta "Your Name" you@clienta.com --ssh-key ~/.ssh/id_clienta
gitswitch add clientb "Your Name" you@clientb.com --ssh-key ~/.ssh/id_clientb
gitswitch clienta   # before working on client A's repo
```

**Open source contributor with a separate public identity**
```bash
gitswitch add oss  "Your Name" public@example.com --gh-user yourhandle-oss
gitswitch add day  "Your Name" you@company.com    --gh-user yourhandle-work
gitswitch oss   # before opening a PR on a public repo
```

**Multi-account GitHub setup**
```bash
gitswitch add personal "Alice" alice@gmail.com   --ssh-key ~/.ssh/id_personal --gh-user alice
gitswitch add work     "Alice" alice@company.com --ssh-key ~/.ssh/id_work     --gh-user alice-corp
```

**Always use the right identity automatically**
```bash
gitswitch shell      # set up shell integration once (enables Session Isolation by default)
gitswitch pin work     # pin 'work' to your work repo — never forget again
```

---

## Built with

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — terminal styling
- [Cobra](https://github.com/spf13/cobra) — CLI commands

---

## Credits

SSH key switching, `gh auth switch` integration, and profile state detection logic inspired by [dankozlowski/git-switcher](https://github.com/dankozlowski/git-switcher).
