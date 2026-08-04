---
title: SSH Keys
description: Force exactly the right key per account — no more agent roulette
---
By default, SSH offers your agent's keys one after another until something works. If two of your GitHub accounts each have a key loaded, "works" and "the one you meant" are not the same thing. That's how a push lands on the wrong account.

gitswitch pins it down.

## What it sets

Switching to an account with an SSH key writes:

```bash
git config --global core.sshCommand "ssh -i ~/.ssh/id_work -o IdentitiesOnly=yes"
```

`IdentitiesOnly=yes` is the whole trick: SSH will offer *only* the key named with `-i`, and nothing from the agent. Every `push`, `pull`, and `fetch` uses exactly that key.

Switching to an account **without** an SSH key unsets `core.sshCommand`, so the previous account's key never lingers.

## Point an account at a key

```bash
gitswitch add work "Alice Smith" alice@company.com --ssh-key ~/.ssh/id_work
```

Or add it to an existing account in the TUI (`e` to edit). The path is stored as you typed it and expanded when applied — `~` and `$HOME` both work.

## Don't have a key yet?

You may not need one. `gitswitch login` + [HTTPS routing](/docs/routing/https) covers pushing without any key management at all. If you'd still rather use SSH:

```bash
ssh-keygen -t ed25519 -f ~/.ssh/id_work -C "alice@company.com"
cat ~/.ssh/id_work.pub
# paste into GitHub → Settings → SSH and GPG keys → New SSH key
```

Check it:

```bash
ssh -i ~/.ssh/id_work -T git@github.com
# Hi alice-corp! You've successfully authenticated...
```

That greeting names the account GitHub thinks you are. Very useful when something's off.

## Verify what's active

```bash
gitswitch work
git config --global core.sshCommand
# ssh -i /Users/alice/.ssh/id_work -o IdentitiesOnly=yes
```

## Keys with passphrases

gitswitch doesn't cache passphrases — it only sets `core.sshCommand`. Put the key in your agent so git doesn't prompt:

```bash
ssh-add ~/.ssh/id_work
```

`IdentitiesOnly=yes` still works with an agent-held key: the agent does the crypto, but only the `-i` key is offered.

## One key per account

```bash
gitswitch add personal "Alice" alice@gmail.com  --ssh-key ~/.ssh/id_personal
gitswitch add work     "Alice" alice@company.com --ssh-key ~/.ssh/id_work
```

GitHub identifies you by which key authenticates, so with `IdentitiesOnly` on, the right account is guaranteed rather than probable.

## A different key for one repo only

```bash
cd ~/special-repo
gitswitch pin work    # applies work's key (and identity) to this repo alone
```

Or by hand, if the key doesn't belong to any account:

```bash
git config --local core.sshCommand "ssh -i ~/.ssh/id_special -o IdentitiesOnly=yes"
```

## Troubleshooting

**`Permission denied (publickey)`**

```bash
git config --global core.sshCommand      # is the path right?
ssh -i ~/.ssh/id_work -T git@github.com  # does the key work at all?
ssh-add -l                               # is a passphrase-protected key loaded?
```

**Authentication feels slow, or picks the wrong account**

That's SSH trying keys in turn, which means `IdentitiesOnly=yes` isn't in play:

```bash
git config --global core.sshCommand
# must contain: -o IdentitiesOnly=yes
```

Re-switch the account (`gitswitch work`) to rewrite it. If the repo is pinned, `gitswitch pin work` again.

## Next

- **[Commit Signing](/docs/accounts/gpg)** — you can sign with an SSH key too
- **[HTTPS Push Routing](/docs/routing/https)** — the no-keys-at-all alternative
- **[Troubleshooting](/docs/troubleshooting)**
