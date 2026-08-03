---
title: Agent Aided Setup
description: One URL. Give it to your agent and it takes care of everything.
---

Give your coding agent this URL:

```
https://get.gitswitch.dev
```

That's the whole thing. Say something like "set up gitswitch for me: `https://get.gitswitch.dev`" to Claude Code, Codex, or whatever you're using, and it takes it from there.

## Why one URL is enough

That URL is both a shell script and its own instructions. `curl -fsSL https://get.gitswitch.dev | bash` installs gitswitch — and the same file is full of comments telling an agent exactly what to do once it has: check `gitswitch doctor --json`, run `gitswitch login`, hand you the browser step, confirm with `doctor --json` again. An agent that fetches the URL as text sees the entire playbook right there, without needing this page at all — this page exists for the curious human, not because the agent needs it.

Nothing to copy-paste, no separate manifest to keep in sync with the installer. One link, and the agent reads its own instructions from the thing it just downloaded.

## Why a human still shows up once

GitHub's OAuth device flow (`gitswitch login`) is deliberately a human-in-the-loop step: it prints a URL and a short code, and a person has to open that in a browser and approve it. An agent can't click through a browser on someone's behalf, and scripting around it (scraping the page, automating the browser) is exactly what device flow exists to prevent. The install script's own comments tell the agent this — surface the URL and code, ask the human to approve it, then wait; the command blocks until they do.

Once that's done, gitswitch also registers the new account with the `gh` CLI itself, so Session Isolation and HTTPS push routing have something to find immediately — no extra step for the agent or the human.

## Next

- **[Quick Start](/docs/get-started/quick-start)** — the human-readable walkthrough of the same flow
- **[AI Coding Agents](/docs/ai/ai-agents)** — pinning repos and fixing commits an agent already got wrong
- **[CLI Reference](/docs/cli/commands)** — every command and flag, including `doctor --json`'s full schema
