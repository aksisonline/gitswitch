---
title: Multi-client Freelancer
description: Managing separate git identities for multiple clients, no manual setup
---

Setup: 3 clients, 3 separate identities. No SSH keys to generate, no keys to paste into anyone's settings page — gitswitch handles the whole thing.

## 1. One-time setup

```bash
gitswitch install   # shell integration, HTTPS routing, Session Isolation — all on by default
```

## 2. Connect each account

```bash
gitswitch login   # log in as client A's account in the browser
gitswitch login   # run again for client B
gitswitch login   # and again for your personal account
```

Each login creates a profile automatically. Check them:

```bash
gitswitch list
```

```
   clienta         you@clienta.com
   clientb         you@clientb.com
   personal        you@personal.com
```

## 3. Pin each repo

```bash
cd ~/clients/clienta/main-repo
gitswitch pin clienta

cd ~/clients/clientb/project-x
gitswitch pin clientb

cd ~/personal/side-project
gitswitch pin personal
```

Done. Every commit and push from here on uses the right identity automatically — no nudges, no remembering which key goes with which client.

## Daily workflow

```bash
cd ~/clients/clienta/main-repo
git commit -m "Fix bug"
git push
# committed and pushed as you@clienta.com, automatically
```

```bash
cd ~/clients/clientb/project-x
git commit -m "Ship feature"
git push
# committed and pushed as you@clientb.com — no switching required
```

Check which identity is active any time:

```bash
gitswitch current
# clienta — Your Name <you@clienta.com>  (pinned to this repo)
```

## Fixed a commit with the wrong identity?

Not pushed yet:

```bash
gitswitch clienta
git commit --amend --reset-author --no-edit
```

Already pushed:

```bash
gitswitch reauthor 1 --to clienta --push
# rewrites the last commit to the 'clienta' identity and force-pushes safely
```

## Not sure which client owns a repo?

```bash
git remote -v
git log -3 --format="%an <%ae>"
```

## Prefer SSH keys you already manage yourself?

Optional — HTTPS (the default) needs nothing from you. Point a profile at an existing key instead:

```bash
gitswitch add clienta "Your Name" you@clienta.com --ssh-key ~/.ssh/id_clienta
```

## Next steps

- [Multi-account GitHub](/docs/scenarios/multi-github)
- [Open Source + Work](/docs/scenarios/oss)
- [Identity Awareness](/docs/features/identity-awareness)
