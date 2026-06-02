---
title: SSH Key Management
description: How gitswitch forces a specific SSH key per profile
---

This page covers how SSH key switching works and how to set it up.

## What it sets

When you switch to a profile that has an SSH key, gitswitch runs:

```bash
git config --global core.sshCommand "ssh -i ~/.ssh/id_work -o IdentitiesOnly=yes"
```

`-o IdentitiesOnly=yes` tells SSH not to offer any other keys from your agent. Every `git push`, `git pull`, and `git fetch` uses exactly that key.

## Add a profile with an SSH key

```bash
gitswitch add work "Alice Smith" alice@company.com \
  --ssh-key ~/.ssh/id_work
```

The value passed to `--ssh-key` is stored as-is and expanded at switch time. Tilde (`~`) and `$HOME` both expand to your home directory.

## Generate a new SSH key

```bash
ssh-keygen -t ed25519 -f ~/.ssh/id_work -C "alice@company.com"
```

Add the public key to the service:

```bash
cat ~/.ssh/id_work.pub
# Paste this into GitHub → Settings → SSH and GPG keys → New SSH key
```

Test it:

```bash
ssh -i ~/.ssh/id_work -T git@github.com
# Hi alice-work! You've successfully authenticated, but GitHub does not provide shell access.
```

## Verify the active SSH key

After switching:

```bash
gitswitch work
git config --global core.sshCommand
# ssh -i /Users/alice/.ssh/id_work -o IdentitiesOnly=yes
```

## SSH keys with passphrases

gitswitch does not cache passphrases — it only sets `core.sshCommand`. If your key has a passphrase, add it to your SSH agent so git can use it without prompting:

```bash
ssh-add ~/.ssh/id_work
```

The `-o IdentitiesOnly=yes` flag is compatible with agent-provided keys: SSH will use the agent for decryption but still only offer the key specified with `-i`.

## Multiple keys for multiple GitHub accounts

```bash
gitswitch add personal "Alice" alice@gmail.com --ssh-key ~/.ssh/id_personal
gitswitch add work "Alice" alice@company.com --ssh-key ~/.ssh/id_work
```

GitHub identifies accounts by which SSH key authenticates. With `IdentitiesOnly=yes`, the correct key is always used — no cross-account authentication surprises.

## Per-repo SSH key override

If you need a different key for one repo without switching your global profile:

```bash
cd ~/special-repo
git config --local core.sshCommand "ssh -i ~/.ssh/id_special -o IdentitiesOnly=yes"
```

Local config takes priority over global.

## Troubleshooting

**`Permission denied (publickey)`**

```bash
# Verify the key path is correct
git config --global core.sshCommand

# Test the key directly
ssh -i ~/.ssh/id_work -T git@github.com

# If the key has a passphrase, ensure it's in the agent
ssh-add -l  # list loaded keys
ssh-add ~/.ssh/id_work
```

**Slow authentication**

If SSH tries multiple keys before succeeding, `IdentitiesOnly=yes` is not set:

```bash
git config --global core.sshCommand
# Must contain: -o IdentitiesOnly=yes
```

If it's missing, re-switch the profile:

```bash
gitswitch work
```

## Next steps

- [GPG Signing](/docs/features/gpg)
- [GitHub Account Sync](/docs/features/github-sync)
- [Identity Awareness](/docs/features/identity-awareness)
