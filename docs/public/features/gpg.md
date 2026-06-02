---
title: GPG Signing
description: How gitswitch switches GPG signing keys per profile
---

This page covers GPG key switching — what gitswitch sets, how to find your key ID, and how to verify it works.

## What it sets

When you switch to a profile with a GPG key, gitswitch runs:

```bash
git config --global user.signingkey ABCD1234EF567890
```

Git uses this key when signing commits (`git commit -S`) or when `commit.gpgsign true` is set.

## Add a profile with a GPG key

```bash
gitswitch add work "Alice Smith" alice@company.com \
  --sign-key ABCD1234EF567890
```

The value passed to `--sign-key` is a GPG key ID — the 16-character hex string after the `/` in `gpg --list-secret-keys`.

## Find your GPG key ID

```bash
gpg --list-secret-keys --keyid-format LONG
```

```
sec   rsa4096/ABCD1234EF567890 2023-01-15 [SC]
      ABCD1234EF5678901234567890ABCDEF12345678
uid                 [ultimate] Alice Smith <alice@company.com>
```

Use the 16-character ID: `ABCD1234EF567890`.

## Generate a new GPG key

```bash
gpg --gen-key
# Follow the prompts; use the same email as your git profile
```

After generating, find the ID with `gpg --list-secret-keys --keyid-format LONG`.

Add the public key to GitHub:

```bash
gpg --armor --export ABCD1234EF567890
# Paste the output into GitHub → Settings → SSH and GPG keys → New GPG key
```

## Sign commits

After switching to a profile with a signing key:

```bash
# Sign a single commit
git commit -S -m "Signed commit"

# Or always sign automatically
git config --global commit.gpgsign true
```

## Verify a signed commit

```bash
git verify-commit HEAD
```

```
gpg: Signature made Mon 15 Jan 2024
gpg:                using RSA key ABCD1234EF567890
gpg: Good signature from "Alice Smith <alice@company.com>"
```

Or view signatures in the log:

```bash
git log --show-signature -1
```

## Multiple profiles, multiple keys

```bash
gitswitch add personal "Alice" alice@gmail.com \
  --sign-key PERSONAL_KEY_ID

gitswitch add work "Alice Smith" alice@company.com \
  --sign-key WORK_KEY_ID
```

Switching profiles switches `user.signingkey` automatically.

## Troubleshooting

**`error: key "KEYID" does not contain a secret key`**

The key is not in your keyring:

```bash
gpg --list-secret-keys --keyid-format LONG
# Check if the key ID appears

gpg --import ~/path/to/backup.gpg
# Import if needed
```

**GPG hangs asking for passphrase repeatedly**

Configure `gpg-agent` to cache the passphrase:

```bash
# ~/.gnupg/gpg-agent.conf
default-cache-ttl 3600
max-cache-ttl 7200
```

Restart the agent:

```bash
gpgconf --kill gpg-agent
```

**Commits show "Unverified" on GitHub**

GitHub requires:
1. The public key is added to your GitHub account (Settings → SSH and GPG keys)
2. The email on the key matches your git `user.email`

Check both and ensure the email in your gitswitch profile matches the GPG key's UID email.

## Next steps

- [GitHub Account Sync](/docs/features/github-sync)
- [Commit Identity](/docs/features/commit-identity)
- [SSH Keys](/docs/features/ssh-keys)
