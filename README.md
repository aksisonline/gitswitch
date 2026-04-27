# git-switcher

A terminal UI for managing multiple local git identities. Switch the name, email, SSH key, and GitHub account that are used for your commits — instantly, without touching config files manually.

![git-switcher TUI](https://img.shields.io/badge/built%20with-Go-00ADD8?style=flat-square&logo=go)
![license](https://img.shields.io/badge/license-MIT-green?style=flat-square)

> **Credits:** Switching logic (SSH key management, `gh auth switch` integration, profile state detection) is inspired by [dankozlowski/git-switcher](https://github.com/dankozlowski/git-switcher). The TUI, profile storage, and CLI design are original to this project.

**→ [Full CLI reference](docs/cli.md)**

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

| | `gh auth switch` | `git-switcher` |
|---|---|---|
| **Controls** | GitHub API tokens | `git config user.name/email` |
| **Affects** | `gh` CLI commands | Commit author identity |
| **SSH key switching** | No | Yes — sets `core.sshCommand` |
| **GPG signing key** | No | Yes — sets `user.signingkey` |
| **Works with GitLab/Bitbucket** | No | Yes — any git remote |
| **Talks to GitHub** | Yes (OAuth) | Optional — `gh auth switch` is best-effort |

**The typical problem:** You push a commit to your company repo, only to see it attributed to your personal email. `gh auth switch` would not have helped. `git-switcher` would.

---

## Install

**Homebrew** (recommended):
```bash
brew tap aksisonline/gitswitch
brew install gitswitch
```

**Curl** (one-liner):
```bash
curl -fsSL https://raw.githubusercontent.com/aksisonline/gitswitch/main/.github/install.sh | bash
```

Or **build from source**:
```bash
git clone https://github.com/aksisonline/gitswitch
cd git-switcher
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
╭──────────────────────────────────────────────────────╮
│  ✦  Git-Switcher                                     │
│     Made by AKS  ·  abhiramkanna.com                 │
│                                                      │
│  Current  username  ·  user@gmail.com                │
│                                                      │
│     ·  default       user@default.com                │
│  ❯  ✓  aksisonline   user@gmail.com                  │
│     ·  work          user@company.com                │
│                                                      │
│  ────────────────────────────────────────────────    │
│  ↑/↓ navigate  ·  enter switch  ·  a add             │
│  e edit  ·  ? cli tips  ·  c theme  ·  q quit        │
╰──────────────────────────────────────────────────────╯
```

**Keys:**
| Key | Action |
|-----|--------|
| `↑` / `↓` or `k` / `j` | Navigate profiles |
| `enter` | Switch to selected profile |
| `a` | Add new profile |
| `e` | Edit selected profile |
| `ctrl+d` (in edit) | Delete profile |
| `c` | Cycle color theme (12 palettes) |
| `?` | Show CLI quick reference |
| `q` | Quit |

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

## Profile fields

| Field | Git config key | Description |
|-------|---------------|-------------|
| Nickname | — | Label shown in the list. Not written to git config. |
| User Name | `user.name` | Author name on commits. |
| Email | `user.email` | Author email on commits. |
| GPG Signing Key | `user.signingkey` | Optional. For signed commits. |
| SSH Key Path | `core.sshCommand` | Optional. Path to SSH private key (e.g. `~/.ssh/id_work`). Sets `ssh -i <key> -o IdentitiesOnly=yes` to force that exact key and prevent SSH agent fallback. |
| GitHub Username | — | Optional. Runs `gh auth switch --user <username>` on switch. Fails gracefully if `gh` is not installed. |

Profiles are stored at `~/.config/gitswitch/profiles.json`. UI preferences (color theme) are stored at `~/.config/gitswitch/config.json`.

---

## How switching works

On switch, `git-switcher` runs (in order):

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

---

## Built with

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — terminal styling
- [Cobra](https://github.com/spf13/cobra) — CLI commands

---

## Credits

SSH key switching, `gh auth switch` integration, and profile state detection logic inspired by [dankozlowski/git-switcher](https://github.com/dankozlowski/git-switcher).

---

Made by [AKS](https://abhiramkanna.com)
