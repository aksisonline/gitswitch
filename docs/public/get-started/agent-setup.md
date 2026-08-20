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

`gitswitch login` hands off to `gh auth login`, which is deliberately a human-in-the-loop step: it prints a URL and a code (or opens a browser), and a person has to approve it. An agent can't click through a browser on someone's behalf, and scripting around it is exactly what this exists to prevent. The install script's own comments tell the agent this — surface whatever prompt `gh` shows, ask the human to approve it, then wait; the command blocks until they do.

Once that's done, the account is already registered with the `gh` CLI — `gh auth login` is what did the logging in — so Session Isolation and HTTPS push routing have something to find immediately, no extra step for the agent or the human.

## Next

- **[Quick Start](/docs/get-started/quick-start)** — the human-readable walkthrough of the same flow
- **[AI Coding Agents](/docs/ai/ai-agents)** — pinning repos and fixing commits an agent already got wrong
- **[CLI Reference](/docs/cli/commands)** — every command and flag, including `doctor --json`'s full schema
