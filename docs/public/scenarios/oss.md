---
title: Open Source + Work
description: Separate git identities for day-job and personal open source contributions
---

Setup: 2 GitHub accounts — `alice-work` for the company and `alice` for personal/OSS. Different emails, different SSH keys, both active on the same machine.

## 1. Generate SSH keys

```bash
ssh-keygen -t ed25519 -f ~/.ssh/id_work -C "alice@company.com"
ssh-keygen -t ed25519 -f ~/.ssh/id_oss  -C "alice@github.com"
```

## 2. Add keys to GitHub

**Work account (alice-work):**

```bash
cat ~/.ssh/id_work.pub
# Paste into github.com/settings/keys as alice-work account
```

**Personal account (alice):**

```bash
cat ~/.ssh/id_oss.pub
# Paste into github.com/settings/keys as alice account
```

Test both:

```bash
ssh -i ~/.ssh/id_work -T git@github.com
# Hi alice-work! You've successfully authenticated...

ssh -i ~/.ssh/id_oss -T git@github.com
# Hi alice! You've successfully authenticated...
```

## 3. Create gitswitch profiles

```bash
gitswitch add work "Alice Smith" alice@company.com \
  --ssh-key ~/.ssh/id_work \
  --gh-user alice-work

gitswitch add oss "Alice" alice@github.com \
  --ssh-key ~/.ssh/id_oss \
  --gh-user alice
```

## 4. Pin repos

```bash
cd ~/work/internal-api
gitswitch pin work

cd ~/projects/my-library
gitswitch pin oss
```

## Daily workflows

### Working on company code

```bash
cd ~/work/internal-api
# nudge: "this repo usually uses work — switch? [y/N]"
y

git checkout -b fix/critical-bug
git commit -m "Fix critical bug"
# committed as alice@company.com, SSH auth as alice-work

git push origin fix/critical-bug
# gh creates PR from alice-work account if --gh-user is set
```

### Contributing to OSS

```bash
cd ~/projects/my-library
# nudge: "this repo usually uses oss — switch? [y/N]"
y

git checkout -b add-feature
git commit -m "Add feature"
# committed as alice@github.com, SSH auth as alice

git push origin add-feature
```

### Using gh CLI with both accounts

Because both profiles have `--gh-user` set, `gh` switches alongside the git identity:

```bash
# Working as OSS account
gitswitch oss
gh repo list                         # lists alice's repos
gh issue list --repo alice/my-lib    # issues on alice's repos

# Working as work account
gitswitch work
gh pr list --repo alice-work/api     # PRs on company repos
gh pr create --title "Fix" --body "" # creates PR as alice-work
```

## Verify correct attribution

Before pushing any commit:

```bash
gitswitch current
# oss — Alice <alice@github.com>

git log -1 --format="%an <%ae>"
# Alice <alice@github.com>

ssh -T git@github.com
# Hi alice! You've successfully authenticated...
```

## Fix a wrong-identity commit before pushing

```bash
# Committed as work identity but should be oss
gitswitch oss
git commit --amend --reset-author --no-edit
git log -1 --format="%an <%ae>"
# Alice <alice@github.com>
```

## Troubleshooting

**Commits are showing as alice-work on a personal repo**

The identity was wrong when the commit was made. If not yet pushed:

```bash
gitswitch oss
git commit --amend --reset-author --no-edit
```

If already pushed, the attribution is permanent. For future commits, use `gitswitch pin oss` in that repo.

**SSH key rejected**

```bash
gitswitch current
git config --global core.sshCommand
# Verify it points to the expected key

ssh -i ~/.ssh/id_oss -T git@github.com
# Should authenticate as alice
```

## Next steps

- [Multi-account GitHub](/docs/scenarios/multi-github)
- [Multi-client Freelancer](/docs/scenarios/freelancer)
- [GitHub Account Sync](/docs/features/github-sync)
