---
title: Multi-account GitHub
description: Cloning, pushing, and switching between two GitHub accounts on one machine
---

Setup: 2 GitHub accounts on one machine — `alice` (personal) and `alice-corp` (work). Each needs its own SSH key so GitHub can identify which account is making the request.

## 1. Generate SSH keys

```bash
ssh-keygen -t ed25519 -f ~/.ssh/id_personal -C "alice@gmail.com"
ssh-keygen -t ed25519 -f ~/.ssh/id_work     -C "alice@company.com"
```

## 2. Add keys to GitHub

**Personal account (alice):**

```bash
cat ~/.ssh/id_personal.pub
# Add to: github.com/settings/keys (logged in as alice)
```

**Work account (alice-corp):**

```bash
cat ~/.ssh/id_work.pub
# Add to: github.com/settings/keys (logged in as alice-corp)
```

Test both:

```bash
ssh -i ~/.ssh/id_personal -T git@github.com
# Hi alice! You've successfully authenticated...

ssh -i ~/.ssh/id_work -T git@github.com
# Hi alice-corp! You've successfully authenticated...
```

## 3. Create gitswitch profiles

```bash
gitswitch add personal "Alice" alice@gmail.com \
  --ssh-key ~/.ssh/id_personal \
  --gh-user alice

gitswitch add work "Alice" alice@company.com \
  --ssh-key ~/.ssh/id_work \
  --gh-user alice-corp
```

## 4. Pin your repos

```bash
cd ~/projects/my-library
gitswitch pin personal

cd ~/work/internal-api
gitswitch pin work
```

## Cloning from each account

Always switch before cloning so `core.sshCommand` points to the right key:

```bash
# Clone a personal repo
gitswitch personal
git clone git@github.com:alice/my-project.git

# Clone a work repo
gitswitch work
git clone git@github.com:alice-corp/company-service.git
```

## Daily workflow

### Personal project

```bash
cd ~/projects/my-library
# nudge: "this repo usually uses personal — switch? [y/N]"
y

git commit -m "Add feature"
# Committed as alice@gmail.com using ~/.ssh/id_personal

git push  # Pushes as alice account
```

### Work project

```bash
cd ~/work/internal-api
# nudge: "this repo usually uses work — switch? [y/N]"
y

git commit -m "Fix bug"
# Committed as alice@company.com using ~/.ssh/id_work

git push  # Pushes as alice-corp account
```

### Check which account will push

```bash
gitswitch current
# work — Alice <alice@company.com>

ssh -T git@github.com
# Hi alice-corp! You've successfully authenticated...
```

## Managing notifications across both accounts

```bash
gitswitch personal
gh notification list   # personal account notifications

gitswitch work
gh notification list   # work account notifications
```

## Handling repos owned by one account that the other contributes to

If `alice-corp` wants to contribute to `alice`'s repo:

```bash
git clone git@github.com:alice/open-project.git
cd open-project
gitswitch work           # use work identity
gitswitch pin work       # pin so nudge reminds you

git commit -m "Add contribution"
git push                 # pushes as alice-corp
```

The contribution will be attributed to `alice-corp`. If you want it attributed to `alice`, switch to personal first.

## Fix a wrong-account push (before others pull)

```bash
# Pushed from wrong account
gitswitch personal
git commit --amend --reset-author --no-edit
git push --force-with-lease  # safe if no one else has pulled
```

After `--force-with-lease`, the commit is now attributed to the personal account.

## Troubleshooting

**`Permission denied (publickey)`**

```bash
gitswitch current
# Check which profile is active

git config --global core.sshCommand
# Verify it points to the expected key

ssh -i ~/.ssh/id_personal -T git@github.com
# Should show: Hi alice!
```

**Unsure which GitHub account owns a repo**

```bash
git remote -v
# origin  git@github.com:alice/my-project.git → personal account
# origin  git@github.com:alice-corp/service.git → work account
```

**Too many repos to manage individually**

Organize repos by account in a directory structure:

```
~/github/
├── alice/           # personal repos
│   └── my-project/
└── alice-corp/      # work repos
    └── service/
```

Pin once per directory as you set them up, and the identity awareness system tracks the rest.

## Next steps

- [Open Source + Work](/docs/scenarios/oss)
- [Multi-client Freelancer](/docs/scenarios/freelancer)
- [SSH Keys](/docs/features/ssh-keys)
