---
title: Multi-client Freelancer
description: Managing separate git identities and SSH keys for multiple clients
---

Setup: 3 clients, 3 separate git identities, 3 separate SSH keys. Each client uses a different git host account.

## 1. Generate SSH keys

One key per client:

```bash
ssh-keygen -t ed25519 -f ~/.ssh/id_clienta -C "you@clienta.com"
ssh-keygen -t ed25519 -f ~/.ssh/id_clientb -C "you@clientb.com"
ssh-keygen -t ed25519 -f ~/.ssh/id_personal -C "you@personal.com"
```

## 2. Add public keys to each client's git service

For each client:

```bash
cat ~/.ssh/id_clienta.pub
# Paste into their GitHub/GitLab/Bitbucket SSH keys page

# Test
ssh -i ~/.ssh/id_clienta -T git@github.com
# Hi client-a-account! You've successfully authenticated...
```

## 3. Create gitswitch profiles

```bash
gitswitch add clienta "Your Name" you@clienta.com --ssh-key ~/.ssh/id_clienta
gitswitch add clientb "Your Name" you@clientb.com --ssh-key ~/.ssh/id_clientb
gitswitch add personal "Your Name" you@personal.com --ssh-key ~/.ssh/id_personal
```

Verify:

```bash
gitswitch list
```

```
   clienta         you@clienta.com
   clientb         you@clientb.com
   personal        you@personal.com
```

## 4. Pin each repo

Navigate to each client repo and pin it:

```bash
cd ~/clients/clienta/main-repo
gitswitch pin clienta

cd ~/clients/clientb/project-x
gitswitch pin clientb

cd ~/personal/side-project
gitswitch pin personal
```

With [shell integration](/docs/features/shell) installed, entering any of these repos will show a nudge:

```
gitswitch: this repo usually uses clienta — switch? [y/N]
```

## Daily workflow

### Starting work on a client project

```bash
cd ~/clients/clienta/main-repo
# nudge fires if identity is wrong
y  # switch immediately

gitswitch current
# clienta — Your Name <you@clienta.com>

git commit -m "Fix bug"
# committed as you@clienta.com using ~/.ssh/id_clienta
```

### Switching between clients

```bash
cd ~/clients/clientb/project-x
y  # accept nudge to switch to clientb

gitswitch current
# clientb — Your Name <you@clientb.com>
```

### Verify before pushing to a client repo

```bash
git log -1 --format="%an <%ae>"
# Your Name <you@clienta.com>

ssh -i ~/.ssh/id_clienta -T git@github.com
# Hi client-a-account! You've successfully authenticated...
```

## Fix a wrong-identity commit before pushing

If you committed with the wrong identity and haven't pushed yet:

```bash
gitswitch clienta       # switch to correct identity
git commit --amend --reset-author --no-edit
git log -1 --format="%an <%ae>"
# Your Name <you@clienta.com>
```

Once pushed, attribution is permanent.

## Troubleshooting

**`Permission denied (publickey)` when pushing**

```bash
# Check which identity is active
gitswitch current

# Check what SSH key is configured
git config --global core.sshCommand
# ssh -i /Users/you/.ssh/id_clienta -o IdentitiesOnly=yes

# Test the key directly
ssh -i ~/.ssh/id_clienta -T git@github.com
```

**Not sure which client owns a repo**

```bash
git remote -v
# origin  git@github.com:client-a-account/repo.git (fetch)

git log -3 --format="%an <%ae>"
# Shows which identities were used in recent commits
```

## Next steps

- [Multi-account GitHub](/docs/scenarios/multi-github)
- [Open Source + Work](/docs/scenarios/oss)
- [Identity Awareness](/docs/features/identity-awareness)
