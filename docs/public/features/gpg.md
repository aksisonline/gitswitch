---
title: Commit Signing
description: A signing key per account — GPG, or the SSH key you already have
---

Signed commits get the green **Verified** badge on GitHub and prove a commit really came from you. Each of your accounts probably needs its own key — gitswitch swaps them along with everything else.

You have two options, and gitswitch supports both without a flag to choose between them: it looks at the key you gave it and figures it out.

## The easy one: sign with SSH

If you already have an SSH key, you already have a signing key. No GPG, no keyring, no passphrase agent.

```bash
gitswitch add work "Alice Smith" alice@company.com \
  --sign-key ~/.ssh/id_ed25519.pub
```

Switching to that account sets:

```bash
git config --global gpg.format   ssh
git config --global user.signingkey /Users/alice/.ssh/id_ed25519.pub
```

Then register it on GitHub as a **Signing Key** — Settings → SSH and GPG keys → New SSH key → key type **Signing Key**. This matters: an *authentication* key with the same contents will not make your commits show as Verified. It's a separate entry.

## The other one: GPG

```bash
gitswitch add work "Alice Smith" alice@company.com \
  --sign-key ABCD1234EF567890
```

Find your key ID (the 16 hex characters after the `/`):

```bash
gpg --list-secret-keys --keyid-format LONG
```

```
sec   rsa4096/ABCD1234EF567890 2023-01-15 [SC]
uid                 [ultimate] Alice Smith <alice@company.com>
```

Don't have one?

```bash
gpg --gen-key                              # use the same email as the account
gpg --armor --export ABCD1234EF567890      # paste into GitHub → New GPG key
```

## How gitswitch tells them apart

| What you pass to `--sign-key` | What gets set |
|---|---|
| `ABCD1234EF567890` — hex key ID | `user.signingkey`, and `gpg.format` **cleared** so git uses OpenPGP |
| `~/.ssh/id_ed25519.pub` — a path | `gpg.format=ssh` + `user.signingkey` with `~` expanded |
| `ssh-ed25519 AAAAC3Nz…` — inline key | `gpg.format=ssh` + `user.signingkey` with git's required `key::` prefix |
| nothing | both **cleared**, so a stale `gpg.format=ssh` can't break the next account |

Which means you can freely mix: a GPG account and an SSH-signing account side by side, each flipping git to the right backend on switch.

## Turn signing on

```bash
git commit -S -m "Signed commit"           # one commit
git config --global commit.gpgsign true    # or always
```

## Check it worked

```bash
git log --show-signature -1
git verify-commit HEAD
```

To verify *other people's* SSH-signed commits, git needs a list of who's allowed to sign what:

```bash
git config --global gpg.ssh.allowedSignersFile ~/.ssh/allowed_signers
echo "alice@company.com $(cat ~/.ssh/id_ed25519.pub)" >> ~/.ssh/allowed_signers
```

## Troubleshooting

**GitHub says "Unverified"**

Both of these must be true:

1. The public key is on your GitHub account — and for SSH keys, added as a **Signing Key**, not an authentication key.
2. The commit's email matches an email verified on that GitHub account, and the key's identity.

**`error: key "KEYID" does not contain a secret key`**

The key isn't in your keyring on this machine:

```bash
gpg --list-secret-keys --keyid-format LONG
gpg --import ~/path/to/backup.gpg
```

**GPG keeps asking for the passphrase**

```bash
# ~/.gnupg/gpg-agent.conf
default-cache-ttl 3600
max-cache-ttl 7200
```

```bash
gpgconf --kill gpg-agent    # then restart
```

**Signing broke after switching accounts**

Check both settings — the mismatch is almost always `gpg.format`:

```bash
git config --global gpg.format
git config --global user.signingkey
```

Re-switch the account to rewrite them as a pair.

## Next

- **[SSH Keys](/docs/features/ssh-keys)** — the authentication half
- **[Commit Identity](/docs/features/commit-identity)** — everything a switch writes
