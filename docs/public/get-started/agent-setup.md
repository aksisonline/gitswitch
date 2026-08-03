---
title: Agent Aided Setup
description: Paste one block into your coding agent's chat and let it set up git, gh, and your GitHub account for you
---

You don't have to run any of the commands on this page yourself. Open Claude Code, Codex, or whatever agent you're already using, and paste this in:

```
Set up gitswitch on this machine, then connect my GitHub account:
1. Install it: curl -fsSL https://get.gitswitch.dev | bash
2. Run `gitswitch doctor --json` and resolve anything it flags (git/gh install, HTTPS routing).
3. Run `gitswitch login` — when it prints a browser URL/device code, tell me to open it and wait for my confirmation before continuing. Don't try to complete the OAuth step yourself.
4. Run `gitswitch doctor --json` once more to confirm everything is wired up, then tell me what profile got created.
```

That's the whole thing. The rest of this page explains *why* each step exists, in case you (or your agent) want to know before running it.

## Why `doctor --json` first

```bash
gitswitch doctor --json
```

```json
{
  "version": "0.4.0",
  "profiles": 0,
  "git": { "installed": true, "ok": true, "version": "2.50.1" },
  "gh": { "installed": false, "ok": false },
  "https": { "routed_by_gitswitch": false },
  "gh_wrapper_installed": false
}
```

Every field an agent needs to decide what's missing, in one call — no parsing prose, no guessing from exit codes. This is the same manifest that used to live behind `gitswitch setup --agent`; it's folded into `doctor --json` now so there's one machine-readable source of truth instead of two.

## Why gitswitch handles git/gh installation itself

The first time `gitswitch` runs bare on a machine with zero profiles, it checks for `git` and `gh` before doing anything else and offers to install whichever's missing, using the right package manager for the OS it's on (`apt`, `dnf`, `pacman`, `brew`, `winget`). An agent shouldn't hand-roll these — the install commands differ per platform, per package manager, and per tool, and gitswitch already gets it right. Just run `gitswitch` (or check `doctor --json` first) and let it do that work.

## Why step 3 needs a human

GitHub's OAuth device flow is deliberately a human-in-the-loop step: `gitswitch login` prints a URL and a short code, and a person has to open that URL in a browser and approve it. An agent can't click through a browser on someone's behalf, and trying to script around it (scraping the page, automating the browser) is exactly the kind of thing that breaks the security model device flow exists for. The right move for an agent here is to surface the URL and code to the user, ask them to complete it, and wait — `gitswitch login` blocks until the flow finishes or times out, so there's nothing to poll.

Once it's done, gitswitch also registers the new account with the `gh` CLI itself (`gh auth login --with-token`), so Session Isolation and HTTPS push routing — which both resolve identity through `gh auth token`/`gh auth switch` — actually have something to find. This happens automatically; nothing extra for the agent (or you) to do.

## Confirming it worked

```bash
gitswitch doctor --json
```

`"profiles"` should now be `1` or more, `"gh": { "installed": true }`, and — if the wizard's shell-integration step ran too — `"https": { "routed_by_gitswitch": true }` and `"gh_wrapper_installed": true`. If any of those are still off, [Troubleshooting](/docs/troubleshooting) covers the common causes.

## Next

- **[Quick Start](/docs/get-started/quick-start)** — the human-readable walkthrough of the same flow
- **[AI Coding Agents](/docs/ai/ai-agents)** — pinning repos and fixing commits an agent already got wrong
- **[CLI Reference](/docs/cli/commands)** — every command and flag, including `doctor --json`'s full schema
