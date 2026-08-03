---
title: Multi-client Freelancer
description: Three clients, three identities, zero mental overhead
---

Three clients, each expecting commits from a different email, plus your own personal account. This is the setup that goes wrong most often — and the one pins solve most completely.

## Set up

```bash
gitswitch shell
```

Then one login per client:

```bash
gitswitch login --profile clienta
gitswitch login --profile clientb
gitswitch login --profile personal
```

```bash
gitswitch list
```

```
✓  clienta         Your Name <you@clienta.com>
   clientb         Your Name <you@clientb.com>
   personal        Your Name <you@personal.com>
```

> Client accounts should also be logged in to `gh` (`gh auth login`) — that's where gitswitch gets push tokens from. If a client only gave you an email and repo access (no GitHub org account), skip that; commits will still be attributed correctly.

## Claim every repo, once

```bash
cd ~/clients/clienta/main-repo && gitswitch pin clienta
cd ~/clients/clientb/project-x && gitswitch pin clientb
cd ~/personal/side-project     && gitswitch pin personal
```

Done. That's the last time you think about this.

## Every day

```bash
cd ~/clients/clienta/main-repo
git commit -m "Fix bug"
git push
# you@clienta.com, automatically
```

```bash
cd ~/clients/clientb/project-x
git commit -m "Ship feature"
git push
# you@clientb.com, no switching
```

Your prompt shows which client you're billing:

```
~/clients/clienta/main-repo  [clienta●] ❯
~/clients/clientb/project-x  [clientb●] ❯
```

## New client, new repo

```bash
gitswitch login --profile clientc
cd ~/clients/clientc/repo && gitswitch pin clientc
```

Two commands and the new client is handled forever.

## Client repo you didn't clone yourself?

If they set up the repo with a local `user.email` already:

```bash
gitswitch pin
```

No nickname needed — gitswitch reads the existing email, matches it to one of your accounts, and fills in the rest.

## Wrong identity on a client's repo

The one that actually matters, because clients read commit logs.

Not pushed yet:

```bash
gitswitch clienta
git commit --amend --reset-author --no-edit
```

Already pushed:

```bash
gitswitch reauthor 1 --to clienta --push
```

A whole batch of commits under the wrong email:

```bash
gitswitch reauthor 20 --to clienta --from you@personal.com --push
```

`--from` scopes it — only the commits currently authored by your personal email get rewritten, everything else is left alone. Coordinate with anyone else on the repo first; this rewrites history.

## Which client owns this repo again?

```bash
gitswitch current      # who am I here, and why
git remote -v          # whose repo is this
git log -3 --format="%an <%ae>"    # what have I been committing as
```

## Prefer your own SSH keys?

Optional — HTTPS needs nothing from you. But if a client handed you a deploy key:

```bash
gitswitch add clienta "Your Name" you@clienta.com --ssh-key ~/.ssh/id_clienta
```

## Next

- **[Two GitHub Accounts](/docs/guides/multi-github)** — the two-account version
- **[Pins & Identity Awareness](/docs/routing/identity-awareness)**
- **[Troubleshooting](/docs/troubleshooting)**
