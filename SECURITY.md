# Security Policy

gitswitch touches sensitive material by design: git identities, SSH/GPG
signing keys, OAuth tokens, and OS keychain entries. Reports involving any
of that get priority handling.

## Reporting a vulnerability

**Please do not open a public GitHub issue for security reports.**

Use GitHub's private reporting instead: go to the
[Security tab](https://github.com/aksisonline/gitswitch/security) →
**Report a vulnerability**. This opens a private advisory only you and the
maintainer can see, and lets us collaborate on a fix before anything is
public.

If you'd rather not use GitHub, email **abhiramkanna@edirq.com** with:

- What you found and why it's a security issue (not just a bug).
- Steps to reproduce, or a PoC if you have one.
- The gitswitch version (`gitswitch version`) and OS/shell you tested on.
- Whether you'd like credit under your name/handle when this is disclosed.

## Scope

In scope:

- Credential/token handling (`internal/secrets`, `internal/oauth`, the
  HTTPS credential helper, keychain storage).
- SSH/GPG key handling and git config writes (`internal/git`).
- Shell integration snippets (`internal/shell`) — anything that could turn
  installed shell hooks into arbitrary code execution.
- Anything that could leak, corrupt, or misattribute another profile's
  identity, keys, or tokens.

Out of scope:

- Issues that require the attacker to already have arbitrary code execution
  or write access to your local machine/config files — gitswitch is a local
  CLI with no server component or attack surface beyond your own shell.
- Social engineering, physical access, or third-party services gitswitch
  merely shells out to (`git`, `gh`) — report those upstream.

## Supported versions

| Version | Supported |
|---------|-----------|
| Latest stable release | ✅ |
| Latest canary/beta | ✅ (best-effort — betas move fast) |
| Anything older | ❌ — please upgrade (`gitswitch stable` / your package manager) and confirm the issue still reproduces |

## What to expect

- Acknowledgment within a few days.
- A fix ships as a patch release (or a beta first, if the fix needs testing)
  as soon as it's ready — there's no fixed SLA, but credential/key-handling
  issues jump the queue.
- We'll credit you in the release notes and the advisory, unless you ask us
  not to.
