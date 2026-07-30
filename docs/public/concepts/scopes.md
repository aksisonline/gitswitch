---
title: Scopes — who you are, and why
description: Global, repo, and terminal identities — how git picks one, and how to tell at a glance
---

Three different things can decide who a commit is from. Knowing which one is in charge is the difference between "I think this is right" and "I know this is right".

## The three scopes

| Scope | Lives in | Applies to |
|---|---|---|
| **Global** | `~/.gitconfig` | every repo on your machine |
| **Repo** | that repo's `.git/config` | one repo, forever |
| **Session** | `GIT_CONFIG_*` environment variables | one terminal, and anything launched from it |

Narrower wins. A repo beats your global config; a terminal's session beats both. That's git's own rule — gitswitch doesn't invent it, it just makes it visible.

```
Session   ◆   ← beats everything
   ↑
 Repo     ●   ← beats global
   ↑
Global    ✓   ← the fallback
```

## How you set each one

```bash
gitswitch work          # global — "I'm working as `work` for now"
gitswitch pin work      # repo   — "this repo is always `work`"
```

There's no `gitswitch` command for the session scope yet — it exists because *other* tools create it (a terminal launched with `GIT_CONFIG_COUNT`/`GIT_CONFIG_KEY_0`/`GIT_CONFIG_VALUE_0` set, or `git -c user.email=…`). gitswitch detects and reports it correctly, which is what matters: if something in your setup is quietly overriding your identity, you'll see it instead of being baffled by it.

## Telling them apart

In the account list (`gitswitch` or `gitswitch list`):

| Marker | Meaning |
|---|---|
| `✓` | your global identity |
| `●` | pinned to this repo |
| `◉` | pinned to this repo **and** your global identity — everything agrees |
| `◆` | this terminal's session is overriding both |
| `·` | just an account you have, not currently in play |

In your shell prompt, the same marker rides next to the nickname:

```
~/personal/blog  [personal] ❯      ← global
~/work/api       [work●] ❯         ← pinned repo
~/somewhere      [work◆] ❯         ← session override
```

Nothing is added when you're on your plain global identity, so if you don't use pins your prompt looks exactly as it always did.

And from the CLI:

```bash
gitswitch current
```

```
work — Alice Smith <alice@company.com>  (pinned to this repo)
```

Machine-readable, for scripts and agents:

```bash
gitswitch current --json
```

```json
{
  "nickname": "work",
  "user_name": "Alice Smith",
  "email": "alice@company.com",
  "scope": "repo",
  "gh_user": "alice-corp",
  "credential_helper_active": true
}
```

## One thing that trips people up

A repo whose local `user.email` is *the same* as your global one is **not** treated as an override — because it isn't one. It changes nothing, so gitswitch reports scope `global`, shows no marker, and keeps its nudges switched on. Otherwise nearly every repo would look "pinned" for no reason.

## Pins need Session Isolation

A repo pin only actually takes effect while [Session Isolation](/docs/features/session-isolation) is on — that's the piece that keeps a repo's `gh` account and push token separate from the rest of your machine. It's on by default, and `gitswitch pin` turns it on for you the first time you pin something.

If you deliberately turned it off, gitswitch says so rather than pretending:

```
work — Alice Smith <alice@company.com>
  (pinned to 'acme' — inactive, Session Isolation is off)
```

## Next

- **[Identity Awareness](/docs/features/identity-awareness)** — pins, learned habits, and nudges
- **[Session Isolation](/docs/features/session-isolation)** — per-terminal `gh` accounts
- **[Commit Identity](/docs/features/commit-identity)** — what a switch actually writes
