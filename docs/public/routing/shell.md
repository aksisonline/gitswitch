---
title: Shell Integration
description: Prompt segment, nudges on cd, tab completion, and the gs alias
---
```bash
gitswitch shell
```

One wizard, three steps, each explaining itself and each skippable:

| Step | |
|---|---|
| **Shell integration** | Prompt segment, identity nudge on `cd`, tab completion, and the `gs` alias — this page |
| **HTTPS credential routing** | Right token per repo on `git push` — [HTTPS Push Routing](/docs/routing/https) |
| **Session Isolation** | Right `gh` account per repo, and working pins — [Session Isolation](/docs/routing/session-isolation) |

Skip the questions:

```bash
gitswitch shell --yes                # accept all defaults, for scripts and CI
gitswitch shell --shell zsh          # target a specific shell (also skips the wizard)
gitswitch shell --https=false        # everything except HTTPS routing
```

Then reload:

```bash
source ~/.zshrc    # zsh
source ~/.bashrc   # bash
# or open a new terminal
```

## It knows your prompt framework

| Framework | What gets written |
|---|---|
| Starship | `[custom.gitswitch]` block in `~/.config/starship.toml` |
| oh-my-zsh | A plugin at `~/.oh-my-zsh/custom/plugins/gitswitch/gitswitch.plugin.zsh` |
| Powerlevel10k | Segment function in your rc file, plus one manual step it prints for `~/.p10k.zsh` |
| Plain zsh / bash / fish | Prompt function, nudge hook, completion, and alias appended to your rc file |

Everything gitswitch writes lives between `# gitswitch` marker comments, and re-running `install` **replaces that block in place** — it never appends a second copy. That's also how you pick up an improved prompt after an upgrade, which is what the "shell integration updated — run `gitswitch shell`" notice is asking for.

## The prompt segment

Shows who you are, but only inside a git repo:

```
~/personal/blog  [personal] ❯
~/work/api       [work●] ❯
~/somewhere      [work◆] ❯
```

That trailing marker is the [scope](/docs/accounts/scopes) — nothing for your global identity, `●` for a pinned repo, `◆` for a terminal override. A glance tells you whether a commit here will be attributed the way you expect.

Under the hood:

```bash
gitswitch current --prompt
# work	141	●
```

Three tab-separated fields: nickname, the ANSI color for your current theme, and the marker. Starship gets `gitswitch current --short` instead (nickname + email, marker riding on the nickname) because it renders command output verbatim.

### Customizing it

**Starship** — edit the block in `~/.config/starship.toml`:

```toml
[custom.gitswitch]
command = "gitswitch current --short"
when = "git rev-parse --git-dir > /dev/null 2>&1"
symbol = " "
style = "bold cyan"
format = "[$symbol($output)]($style) "
```

`style` for color, `symbol` for an icon.

**oh-my-zsh / plain shells** — edit the `__gitswitch_prompt` function in the gitswitch block. Note that re-running `install` replaces that block, so keep changes you want to survive an upgrade somewhere else in the file.

## The nudge

`cd` into a repo where you usually use a different account and you get one line:

```
gitswitch: this repo usually uses work <alice@company.com> — switch? [y/N]
```

`y` switches, `n` or Enter carries on. It never blocks your prompt, and it defaults to no.

Pinned repos are never nudged — they already commit correctly. How the "usually" is worked out: [Identity Awareness](/docs/routing/identity-awareness).

## The `gs` alias

Because you'll type this a lot:

```bash
gs              # same as gitswitch
gs work
gs current
```

Installed with the shell integration. In the TUI's **Settings** tab you can rename it (`e`) or turn it off (`enter`) — either one reinstalls the shell block in place so it takes effect on your next shell reload.

## Tab completion

Commands and your own account nicknames:

```
$ gitswitch <Tab>
add  beta  claude  current  doctor  init  install  list  login  pacman  pin
reauthor  record  recommend  remove  setup  stable  switch  uninstall  unpin
upgrade  version

$ gitswitch w<Tab>
work
```

Needs zsh 5.0+, bash 4.0+, or fish 3.0+. The installed block sources it from `gitswitch completion <shell>`, which you can also wire up yourself if you manage your rc file some other way.

## Uninstall

```bash
gitswitch uninstall
```

Removes the marker block from your rc file (or the oh-my-zsh plugin), unregisters the HTTPS credential helper, and removes the `gh` wrapper — whichever of those were installed. Reload your shell to finish.

## Troubleshooting

**No prompt segment**

- Are you in a git repo? It hides outside one on purpose.
- Did you reload your shell?
- Does `gitswitch current --short` print anything?
- On Starship? Check `~/.config/starship.toml`, not your rc file.

**No nudges**

The threshold is ≥ 3 visits with ≥ 60% consistency before a nudge fires. `gitswitch recommend` shows what it would say right now.

**Completion not working**

Check your shell version against the minimums above, then re-run `gitswitch shell` and reload.

More in [Troubleshooting](/docs/troubleshooting).

## Next

- **[Identity Awareness](/docs/routing/identity-awareness)** — pins, learning, nudges
- **[HTTPS Push Routing](/docs/routing/https)** — the credential helper step
- **[Session Isolation](/docs/routing/session-isolation)** — the `gh` wrapper step
