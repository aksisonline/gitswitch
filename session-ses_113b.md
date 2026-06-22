# git-switcher v1 plan revamp: GitHub OAuth + gs

**Session ID:** ses_113b8da61ffeBHt797YVi7JCNg
**Created:** 6/22/2026, 3:52:26 AM
**Updated:** 6/22/2026, 4:29:00 AM

---

## User

 "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up" 

add new docs, so we are cancelling the rust migration plan cause it's too risky to give up on existing system. So keep those docs, but we will not do that. revamp plan, for v1 launch, we are not going rust, but we will add support for direct OAuth with GitHub, so the plan is that gitswitch (we will implement short command `gs` to invoke gitswitch other than `gitswitch`)

Revamp the v1 launch plan to this. and tell me what the plan looks like now.

---

## Assistant (Plan · Kimi K2.7 Code · 6.9s)

_Thinking:_

The user is in plan mode and wants me to revamp a v1 launch plan document. They've provided a file path and asked to add new docs about cancelling the Rust migration and instead adding direct OAuth with GitHub support and a `gs` short command alias for `gitswitch`.

Since this is plan mode, I must NOT edit files. I should read the existing docs and understand the current plan, then present a revamped plan to the user for approval.

Let me first read the file at the path provided.

I can’t edit files in plan mode, but I can read the current doc and propose the revamped plan for your approval.

Let me look at the existing file first.

**Tool: read**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up"
}
```

**Output:**
```
<path>/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up</path>
<type>directory</type>
<entries>
arcade-philosophy.md
competitor-analysis.md
diagnostic-design.md
explainer.md
gitswitch-v2-discussion.md
marketing.md
monetization.md
product-gaps-and-philosophy.md
project-memory.md
roadmap-full-timeline.md
roadmap-v0.2.0-isolation.md
roadmap-v0.2.0-tui.md
roadmap-v0.2.0.md
roadmap-v0.3.0.md
rust/
versioning.md
why-rust.md

(17 entries)
</entries>
```

---

## Assistant (Plan · Kimi K2.7 Code · 5.5s)

_Thinking:_

The path is a directory with many docs. The user mentioned "keep those docs, but we will not do that [Rust migration]". I should read the relevant docs to understand the current v1 launch plan. Key files likely include:
- roadmap-v0.3.0.md
- roadmap-full-timeline.md
- why-rust.md
- rust/ directory
- gitswitch-v2-discussion.md
- explainer.md

Let me read the most relevant ones in parallel.

**Tool: read**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/roadmap-full-timeline.md"
}
```

**Output:**
```
<path>/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/roadmap-full-timeline.md</path>
<type>file</type>
<content>
1: # gitswitch — Full Product Timeline
2: 
3: _Last updated: May 2026_
4: 
5: ---
6: 
7: ## v0.1.x — Current (Go, shipping now)
8: 
9: The foundation. Everything in this version works today.
10: 
11: **Identity management**
12: - Named profiles: name, email, SSH key, GPG/signing key, gh_user (single host: github.com)
13: - Switch global `~/.gitconfig` identity in one command
14: - `gitswitch current` — see who you are right now
15: 
16: **Shell integration**
17: - Shell hook fires on `cd`, logs repo + identity usage in background
18: - Learns per-repo patterns; suggests the right profile after 3+ uses
19: - Shell prompt segment: shows active profile (or warning if config doesn't match any profile)
20: - Supports zsh, bash, fish, Oh My Zsh
21: 
22: **TUI**
23: - Full-screen terminal UI launched by `gitswitch` with no arguments
24: - Three themes: arcade (loud, retro, animated), modern (clean), minimal (monochrome, SSH-safe)
25: - Pacman easter egg: hidden game using profile list as the level
26: 
27: **Per-repo pinning**
28: - `gitswitch pin <profile>` — this repo always uses this profile
29: - Automatic prompt when entering a pinned repo with wrong identity active
30: - `gitswitch unpin` to remove
31: 
32: **gh CLI sync**
33: - Writes the correct account to `~/.config/gh/hosts.yml` on every switch
34: - git and `gh` can never drift apart
35: 
36: **Distribution**
37: - Single static binary, no runtime dependencies
38: - `curl | sh` one-liner installer
39: - Homebrew tap
40: - Storage: `~/.gitswitch/profiles.json`
41: 
42: **Utility commands**
43: - `gitswitch doctor` — health check on all profiles
44: - `gitswitch upgrade` — self-update from GitHub releases
45: - `gitswitch hook-check` — verify shell hook is installed and current
46: - `gitswitch install` / `gitswitch uninstall`
47: - `gitswitch version`
48: 
49: ---
50: 
51: ## v0.2.0 — Next milestone (Go)
52: 
53: The enterprise and AI agent release. Two features ship in parallel: multi-host identity and session isolation.
54: 
55: ### Multi-host identity
56: 
57: One profile covers multiple GitHub hosts with different credentials per host.
58: 
59: **Profile model change**
60: ```yaml
61: nickname: corp
62: user_name: Abhiram Kanna
63: email: abhiram@corp.com
64: sign_key: 4A1B2C3D
65: hosts:
66:   github.com:
67:     ssh_key: ~/.ssh/id_personal
68:     gh_user: aksisonline
69:     token_ref: gitswitch:corp:github.com
70:   github.corp.com:
71:     ssh_key: ~/.ssh/id_corp
72:     gh_user: aksisonline-pp
73:     token_ref: gitswitch:corp:github.corp.com
74: ```
75: 
76: **What this enables**
77: - One `gitswitch switch corp` sets SSH config, gh auth, and git identity for all hosts simultaneously
78: - Works correctly with repos that have multiple remotes pointing to different hosts
79: - Multi-pin: `gitswitch pin corp --remote origin`, `gitswitch pin personal --remote fork`
80: 
81: **New storage**
82: - Migration from `profiles.json` → `~/.config/gitswitch/config.yaml` (single file: version + prefs + profiles + pins)
83: - Old file renamed `profiles.json.v1.bak` on first load
84: - PAT tokens stored in OS keychain (macOS Keychain, Windows Credential Manager, Linux Secret Service) — never in plaintext
85: 
86: **New commands**
87: - `gitswitch import gh` — discover all accounts from `gh auth status`, create profiles automatically
88: - `gitswitch import git` — import current git config as a new profile (old `init` behavior)
89: - `gitswitch host add <profile> <hostname>` — add a host entry to an existing profile
90: - `gitswitch host remove <profile> <hostname>`
91: - `gitswitch host list <profile>`
92: 
93: ### Session isolation
94: 
95: Each terminal window holds its own identity via environment variables — parallel sessions, independent identities.
96: 
97: **How it works**
98: - `gitswitch shell [profile]` — opens a subshell with identity locked in env vars
99: - Inside the session: git reads `GIT_AUTHOR_NAME`, `GIT_AUTHOR_EMAIL`, etc. in preference to `~/.gitconfig`
100: - `GH_TOKEN` set so `gh` uses the right account
101: - `GITSWITCH_SESSION_ID`, `GITSWITCH_SESSION_PROFILE` for introspection
102: - Exit the shell → env vars disappear → no cleanup needed
103: - Can still `gitswitch switch` within a session to change identity for that session only
104: 
105: **Why this matters for AI agents**
106: - Two Claude Code agents running in parallel, each inside their own `gitswitch shell`, cannot interfere with each other's identity
107: - Every commit from every agent is correctly attributed regardless of what the developer is doing globally
108: - `gitswitch claude` skill integration
109: 
110: ### Diagnostics and safety
111: 
112: **`gitswitch status`** (new) — per-repo narrative diagnostic
113: - What identity is active (global or session)
114: - Every remote: host, SSH config block present?, SSH key exists + permissions correct?, gh auth valid + non-expired?
115: - Pin state, recommendations
116: - Designed to answer "why is my push failing?"
117: 
118: **Pre-commit safety net** (auto-installed by `gitswitch install`)
119: - Checks if current identity matches the pinned profile for the repo
120: - Warns before the wrong-identity commit lands: "you're about to commit as personal@gmail.com — this repo is pinned to work"
121: - User can override with `git commit --no-verify`
122: - Can opt out via config
123: 
124: **Progressive CLI hints**
125: - After `gitswitch switch` → "tip: run `gitswitch pin` to make this automatic"
126: - After `gitswitch pin` → "run `gitswitch status` to see the full picture"
127: - After `gitswitch install` → "shell hook installed — gitswitch will suggest your usual identity when you enter a repo"
128: 
129: ### TUI refresh
130: 
131: - 5 screens: Profiles, Hosts, Repos, Session, Doctor
132: - Command palette (`:` key) for power users
133: - Wizards for profile creation and host addition
134: - New keybindings for session and host management
135: - Shell tab completion ships in v0.2.0
136: 
137: ---
138: 
139: ## Rust migration (parallel track across v0.2.0 → v1.0)
140: 
141: Happening alongside feature development, not blocking it.
142: 
143: - **During v0.2.0**: Go gets the features. Rust rewrite (`gitswitch-rs`) builds in parallel in a 9-crate workspace. Not user-facing yet.
144: - **During v0.3.0**: Rust binary reaches feature parity. Ships as an opt-in alternative binary. Parity tested via shared corpus.
145: - **v1.0**: Rust binary becomes the default. Go binary deprecated.
146: 
147: **Why Rust (locked decision, May 2026)**:
148: 1. Sub-millisecond shell hook — no GC, no runtime. Critical since hook fires on every `cd`.
149: 2. Smaller binary — ~1.5–3MB with LTO vs 4–15MB Go
150: 3. Ratatui + tachyonfx — best-in-class TUI ecosystem; needed for arcade effects (coalesce, dissolve, CRT scanline, ghost movement)
151: 4. Compiler-enforced secret handling — zeroize on drop, borrow checker prevents data races in SSH/GPG/PAT code
152: 
153: ---
154: 
155: ## v0.3.0 — Power user and team features
156: 
157: Items deliberately deferred from v0.2.0 planning. Exact sequencing TBD.
158: 
159: **In-repo `.gitswitch.yaml`**
160: - Per-repo config file optionally checked into the repo root
161: - Declares required profiles, identity policies (signing required, email suffix enforcement)
162: - Loaded by the shell hook automatically
163: - Precedence: session env > local pin > repo YAML > global switch
164: - Open questions: team-overridable vs locked fields, worktree interaction — needs design
165: 
166: **Team preset distribution**
167: - `gitswitch apply <file-or-url>` — bootstraps all profiles from a shared config in one command
168: - Signed with minisign so an attacker can't push a malicious preset
169: - Intended for engineering teams onboarding new developers
170: 
171: **Own OAuth flow**
172: - Direct device-flow OAuth for users who don't use `gh` CLI
173: - Required for GitHub Enterprise Server environments where `gh` isn't installed
174: - v0.2.0 imports from `gh auth` — this adds the fallback for the long tail
175: 
176: **Remote retargeting**
177: - `gitswitch remote retarget` — rewrites `origin` (and others) to a new host while preserving the path
178: - Useful for org migrations between GitHub hosting models
179: 
180: **Host auto-detection**
181: - Walk `git remote -v` on `cd` into a repo
182: - Auto-recommend the profile whose host map covers all remotes
183: - Reduces manual pinning for repos with clear host signals
184: 
185: **Remote-URL-aware routing**
186: - Pin by remote URL pattern instead of directory path
187: - `github.com/mycompany/* → work`, `github.com/aksisonline/* → personal`
188: - More robust than directory-based pinning for developers who don't organize repos conventionally
189: 
190: **git includeIf support (optional, not default)**
191: - Fires identity at git-command time via git's native config loading, not on `cd`
192: - No shell modification needed — alternative for users who want zero shell hooks
193: - gitswitch's shell hook remains the default; includeIf is an opt-in mode
194: 
195: **Other v0.3.0 items**
196: - Ephemeral one-shot run: `gitswitch run <profile> -- <cmd>` (run one command as a different identity without switching)
197: - Profile export/import (JSON backup and restore for new machine setup)
198: - Windows native support improvements (beyond WSL)
199: 
200: ---
201: 
202: ## v1.0 — Production milestone (Rust, stable)
203: 
204: The 1.0 release. Rust binary is the default. Stable CLI API with semver commitment. Enterprise-ready.
205: 
206: **HTTPS credential helper**
207: - `gitswitch` acts as a proper git `credential.helper`
208: - Set with: `git config --global credential.helper gitswitch`
209: - When git needs credentials for any HTTPS remote, it asks gitswitch — gitswitch returns the right PAT for the active profile
210: - Closes the gap for users who clone via HTTPS (the GitHub default) rather than SSH
211: - Means gitswitch handles the full auth stack: SSH (via `~/.ssh/config`), HTTPS (via credential helper), gh CLI (via `hosts.yml`)
212: 
213: **Rust TUI — full arcade build**
214: - Ratatui with crossterm backend
215: - tachyonfx effects: splash screen (coalesce → repeating evolve → parallel dissolve + delayed coalesce)
216: - `ScanlineOverlay` custom widget for CRT effect on arcade theme
217: - `gitswitch-arcade` crate: all easter egg games isolated, compiled in but gated behind unlock sequences
218: - Planned easter eggs beyond Pacman: roguelike dungeon (profile list as rooms), constellation map (repo usage visualized), ghost profiles (":konami" → matrix cascade)
219: 
220: **Stability guarantees**
221: - Semver commitment: no breaking CLI changes without a major version bump
222: - Stable config format: `~/.config/gitswitch/config.yaml` schema versioned
223: - Stable keychain key convention: `gitswitch:<nickname>:<hostname>`
224: 
225: **Distribution at v1.0**
226: - curl one-liner (primary)
227: - Homebrew (official tap)
228: - cargo install
229: - AUR (Arch Linux)
230: - Windows: native installer (NSIS or WiX)
231: - gitswitch.dev live with full docs, blog, and privacy policy
232: 
233: ---
234: 
235: ## What does NOT change across any version
236: 
237: - Every feature, including all enterprise toolsets, stays free forever
238: - No cloud component, no telemetry, no data leaves the machine
239: - No feature flags, no paywalls, no paid themes
240: - The arcade personality and easter eggs are never gatekept
241: - Open source, inspectable at any time

(End of file - total 241 lines)
</content>
```

**Tool: read**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/roadmap-v0.3.0.md"
}
```

**Output:**
```
<path>/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/roadmap-v0.3.0.md</path>
<type>file</type>
<content>
1: # gitswitch v0.3.0 — Roadmap
2: 
3: Items deferred from [v0.2.0 planning](roadmap-v0.2.0.md). v0.2.0 establishes multi-host identity, keychain-backed PATs, gh-driven profile import, multi-pin, and `gitswitch status`. v0.3.0 builds on that foundation to address team-scale distribution, in-repo config, and the long tail of authentication scenarios.
4: 
5: ---
6: 
7: ## 1. In-repo YAML config (`.gitswitch.yaml`)
8: 
9: A per-repo config file checked into the repo root. Declares which profile(s) should be used, which host(s) remotes resolve to, and identity policies the repo enforces.
10: 
11: ### Why
12: 
13: Today (and in 0.2.0), pins live in `~/.config/gitswitch/config.yaml` keyed by absolute repo path. That doesn't survive:
14: - A teammate cloning the repo to a different path.
15: - A CI runner.
16: - Switching machines.
17: 
18: In-repo config makes "this repo wants the `corp` profile" a property of the repo, not the local machine.
19: 
20: ### Open questions
21: 
22: - **Precedence**: in-repo config vs. local pin in `~/.config/gitswitch/config.yaml` — does in-repo win, or does local pin override (escape hatch)?
23: - **Policy enforcement**: should `.gitswitch.yaml` be able to *block* commits that don't match the policy (e.g. via a `prepare-commit-msg` hook)? Or only nudge?
24: - **Locked vs. overridable fields**: a team config probably wants to lock `email_suffix` but leave `ssh_key` path overridable per-dev.
25: - **Worktrees**: same `.gitswitch.yaml` applies to all worktrees of a repo? Or per-worktree opt-out?
26: - **Schema versioning**: how does an older `gitswitch` binary handle a `.gitswitch.yaml` from a newer schema version?
27: 
28: ### Sketch
29: 
30: ```yaml
31: # .gitswitch.yaml (in repo root)
32: version: 1
33: require:
34:   email_suffix: "@corp.com"
35:   signing: required
36: hosts:
37:   origin: github.corp.com
38:   fork:   github.com           # optional, validated against actual remote
39: suggested_profile: corp
40: ```
41: 
42: ---
43: 
44: ## 2. Team preset distribution
45: 
46: A team-shared bootstrap config — distributed via intranet URL, internal git repo, or shared drive — that onboards a new dev in one command:
47: 
48: ```
49: gitswitch apply https://wiki.corp.com/gitswitch-platform-team.yaml
50: gitswitch apply ./team-preset.yaml
51: ```
52: 
53: ### Components
54: 
55: - **Schema** for team presets (distinct from per-repo `.gitswitch.yaml` — this is a *template* that produces profiles, not policy applied to a repo).
56: - **Fetch + cache** for HTTPS sources. Local file path also supported.
57: - **Signing & verification** — so an attacker can't push a malicious preset to wiki / internal repo and steal credentials.
58:   - Candidate formats: [minisign](https://jedisct1.github.io/minisign/), [cosign](https://github.com/sigstore/cosign).
59:   - Trust roots: a team config declares its signing key fingerprint; first-apply prompts the user to trust it (TOFU); subsequent applies verify against the trusted key.
60:   - Revocation: list of revoked keys fetched from a known endpoint.
61: - **Apply semantics**: `apply` creates profile skeletons (host config, gh_user pattern, ssh_key naming hint) and prompts only for personal bits (the dev's actual SSH key path, their actual gh login). Pre-existing profiles aren't clobbered — only new ones are added.
62: 
63: ### Why this is hard enough to defer
64: 
65: The signing story is the load-bearing piece. Shipping `apply` without signing means corp security teams will reject it on day one. Building signing properly (key management, trust onboarding, revocation, UX for "this preset is no longer trusted") is a project on its own.
66: 
67: ---
68: 
69: ## 3. Own OAuth flow (device flow)
70: 
71: Direct OAuth device flow against GitHub, bypassing `gh` CLI.
72: 
73: ### Why
74: 
75: v0.2.0 uses `gh auth` as the import source — works for the majority. The long tail v0.3.0 must address:
76: - Users who don't have `gh` installed (rare but real).
77: - Enterprise GHES instances where `gh` won't authenticate (custom auth gateways, SAML enforcement).
78: - Headless setups where `gh auth login`'s browser flow is impractical.
79: 
80: ### Scope
81: 
82: - Register a gitswitch OAuth app on github.com (one-time; `client_id` baked into the binary, no client_secret since device flow doesn't require one).
83: - Allow user-supplied `client_id` for enterprise GHES (each org registers their own OAuth app).
84: - Scopes: `repo`, `read:user`, `user:email`, `gist`, `workflow` (matches `gh` defaults so the token is a drop-in for `gh auth`).
85: - Command: `gitswitch login [--host HOST] [--client-id ID]`.
86: 
87: ### Open questions
88: 
89: - Should the gitswitch OAuth app be owned by the gitswitch project or by each user's GitHub account?
90: - Refresh tokens: GitHub doesn't issue refresh tokens for device flow by default — UX for expiry?
91: 
92: ---
93: 
94: ## 4. Remote retargeting
95: 
96: `gitswitch remote retarget` — repo-host migration helper.
97: 
98: ### Why
99: 
100: When an org migrates between GitHub hosting models — An org's `github.corp.com` → `github.com/neworg` is the canonical case — every developer has hundreds of local clones with `origin` pointing at the old host. Manual `git remote set-url` per repo is painful.
101: 
102: ### Sketch
103: 
104: ```
105: gitswitch remote retarget --from github.corp.com --to github.com/neworg
106: ```
107: 
108: In each known repo (or the current repo):
109: 1. Walk `git remote -v`.
110: 2. For any remote matching `--from`, rewrite to the equivalent path under `--to`.
111: 3. Confirm before each rewrite (or `--yes` to skip).
112: 4. Optionally switch the active pin to the profile that owns the new host.
113: 
114: ### Open questions
115: 
116: - Should this walk *all* repos under a configured root (e.g. `~/code`) or only the current repo?
117: - Path translation: `github.corp.com:teamA/repo` → `github.com/neworg:teamA/repo`? Or do paths get rewritten too (`neworg-teamA/repo`)?
118: 
119: ---
120: 
121: ## 5. Host auto-detection
122: 
123: On `cd` into a repo (via the existing shell hook), walk `git remote -v` and auto-recommend the profile whose `Hosts` map covers all remotes.
124: 
125: ### Why
126: 
127: Reduces manual pinning. If a repo's only remote is `github.corp.com:team/x.git` and the user has exactly one profile with that host configured, gitswitch can confidently suggest that profile without the user having to `gitswitch pin` it.
128: 
129: ### Sketch
130: 
131: Extend [internal/history/history.go](../internal/history/history.go) `Recommend`:
132: 1. Existing path: usage frequency ≥ 60% → recommend.
133: 2. New path: if no usage history yet AND exactly one profile covers every host in `git remote -v` → recommend that profile.
134: 3. If multiple profiles cover all hosts → no recommendation (ambiguous; defer to user).
135: 
136: ### Open questions
137: 
138: - Should auto-detection auto-switch silently, or always nudge?
139: - Interaction with multi-pin: should auto-detection emit a multi-pin suggestion when different remotes are best served by different profiles?
140: 
141: ---
142: 
143: ## Likely sequencing for v0.3.0
144: 
145: 1. In-repo `.gitswitch.yaml` (foundation: in-repo config schema becomes the substrate for team presets).
146: 2. Host auto-detection (small, high-impact UX win; piggybacks on 0.2.0 shell hook).
147: 3. Team preset apply *without* signing (gated behind `--unsafe-trust` flag; collects feedback on schema before locking it in).
148: 4. Team preset signing (the big one).
149: 5. Own OAuth flow.
150: 6. Remote retargeting.
151: 
152: Order may shift based on which use cases surface from 0.2.0 adopters.

(End of file - total 152 lines)
</content>
```

**Tool: read**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/gitswitch-v2-discussion.md"
}
```

**Output:**
```
<path>/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/gitswitch-v2-discussion.md</path>
<type>file</type>
<content>
1: # gitswitch v2 — Strategy Discussion
2: _Date: May 28, 2026_
3: 
4: ---
5: 
6: ## The Conversation That Sparked This
7: 
8: A chat with Vaishnavi (PayPal engineer) surfaced a use case beyond the current gitswitch scope:
9: 
10: > "We have `github.paypal.com` and `github.com`. We need to push two different things sometimes — gh doesn't recognise properly. We have different PAT tokens."
11: 
12: PayPal is currently migrating away from this dual-instance setup (consolidating to `github.com` under `onepaypal`), but her point was clear: **other companies still have this problem, and many will for years.**
13: 
14: She also flagged **maintainability in the longer run** as a core reason developers need a proper solution — not just at setup time, but as an ongoing daily friction point.
15: 
16: ---
17: 
18: ## What This Reveals About the Market
19: 
20: ### Current gitswitch covers:
21: - Personal vs. work identity switching (name, email, SSH key, GPG key)
22: - Named profiles for individual developers
23: - Single-host GitHub environments
24: 
25: ### The gap Vaishnavi is describing:
26: - Two separate GitHub *hosts* (`github.com` + a self-hosted or enterprise instance)
27: - Different PAT tokens per host
28: - Different SSH/auth contexts per host
29: - All of this needing to coexist within a single developer's daily workflow
30: 
31: This is a **multi-tenancy problem at the infrastructure level**, not just the identity level. It's especially common in:
32: - Financial services (internal GitHub + public GitHub)
33: - Healthcare and defense (on-prem instances alongside cloud)
34: - Large consulting firms (client-specific GitHub orgs/instances)
35: - Any enterprise mid-migration between GitHub hosting models
36: 
37: ---
38: 
39: ## What gitswitch v2 Should Add
40: 
41: ### 1. GitHub Host Targeting
42: A profile should carry not just identity credentials but also a **remote host**:
43: - `github.com` (default)
44: - Custom enterprise host (e.g., `github.paypal.com`, `github.mycompany.com`)
45: 
46: gitswitch should know how to route `gh` auth and `git` remote operations to the right host when a profile is active.
47: 
48: ### 2. PAT Token Management Per Host
49: Different hosts require different PAT tokens. Rather than leaving developers to manage these via `gh auth switch` or environment variables, gitswitch could:
50: - Store PAT tokens per profile/host securely
51: - Swap the active token when switching profiles
52: - Keep `gh` and `git` auth in lockstep (already a roadmap theme)
53: 
54: ### 3. Team/Org Presets (Accelerated Priority)
55: Vaishnavi's comment about maintainability suggests teams want a **shared configuration** that can be distributed — not just individual developer profiles. This means:
56: - A shareable profiles config format (e.g., checked-in team config file)
57: - Onboarding new devs to a team's identity setup in one command
58: - Centralized management of which profiles exist for a team
59: 
60: ---
61: 
62: ## How This Fits the Existing Vision
63: 
64: The current positioning is:
65: > `git` handles version control → `gh` handles GitHub API/PR/auth → `gitswitch` handles **who you are** when you run either
66: 
67: The v2 extension is natural: add **"which GitHub you're talking to"** as part of identity. It's not a pivot — it's the same idea applied one layer deeper.
68: 
69: ```
70: gitswitch profile = {
71:   name,
72:   email,
73:   ssh_key,
74:   gpg_key,          // existing
75:   github_host,      // new
76:   pat_token,        // new
77:   gh_auth_context   // new (sync with gh CLI)
78: }
79: ```
80: 
81: ---
82: 
83: ## Strategic Notes
84: 
85: - You already have company devs using gitswitch — organic adoption at organizations is a strong signal.
86: - The PayPal case is the first **enterprise-grade use case** surfacing from real users.
87: - The problem is widespread but underserved — most teams solve this with a mix of shell aliases, manual `gh auth switch`, and tribal knowledge.
88: - gitswitch is positioned to be the clean, maintained, opinionated solution that none of the existing hacks are.
89: 
90: ---
91: 
92: ## Open Questions for Future Discussion
93: 
94: 1. How should host switching interact with repo-level config? (Should switching to a profile auto-update the remote URL in a repo's `.git/config`?)
95: 2. Should PAT tokens be stored in the OS keychain, or in gitswitch's own encrypted store?
96: 3. What does "team preset" distribution look like — a git repo, a URL, a config file in the project root?
97: 4. Is there a B2B angle here? Enterprises managing developer onboarding would pay for a managed profiles service.
98: 
99: ---
100: 
101: _Discussion between Abhiram Kanna and Claude, based on Vaishnavi (PayPal) feedback captured May 25, 2026._

(End of file - total 101 lines)
</content>
```

**Tool: read**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/why-rust.md"
}
```

**Output:**
```
<path>/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/why-rust.md</path>
<type>file</type>
<content>
1: # Why Rust — gitswitch Technical Rationale
2: 
3: _Decision locked: May 2026. The Go codebase remains the shipping v0.1.x tool; Rust is the target for v0.2.0+. See [docs/rust/roadmap-migration.md](rust/roadmap-migration.md) for the migration plan._
4: 
5: ---
6: 
7: ## The core constraint
8: 
9: gitswitch is not a web server or a batch processor. It is a **shell-integrated, terminal-resident tool** that:
10: 
11: - Runs on every `cd` via a shell hook (latency is user-visible)
12: - Renders a prompt segment on every command (latency is user-visible)
13: - Touches sensitive local state (SSH keys, GPG keys, PAT tokens, `~/.gitconfig`)
14: - Hosts a TUI that should feel snappy and visually rich
15: 
16: These constraints favor Rust specifically over Go. The decision is not "Rust is better than Go" in general — Go is excellent for backend services, CLIs that run once-and-exit infrequently, and network-heavy work. The decision is: **for this tool's exact execution profile, Rust's tradeoffs are exactly right**.
17: 
18: ---
19: 
20: ## Reason 1: Sub-millisecond hook execution
21: 
22: The shell hook fires on every directory change. Users feel this. A 30ms delay on every `cd` is perceptible; over a full workday it becomes friction.
23: 
24: **Go's issue**: Go ships with a runtime and a garbage collector. Even for a binary that does nothing, Go's runtime initialization adds overhead — typically 5–30ms depending on the system. The GC can also pause at inopportune moments.
25: 
26: **Rust's answer**: No runtime. No GC. Rust compiles to bare machine code. Cold-start for a Rust binary doing minimal work (read env var, print string, exit) is measurable in microseconds. For the shell hook path — which runs `gitswitch record --path .` on every `cd` — this difference is the difference between a transparent tool and an annoying one.
27: 
28: Reference: tools like [`starship`](https://starship.rs) (shell prompt, Rust), [`zoxide`](https://github.com/ajeetdsouza/zoxide) (directory jumper, Rust), and [`ripgrep`](https://github.com/BurntSushi/ripgrep) set the bar for what "fast enough to be invisible" means in shell-resident tools. All chose Rust for this reason.
29: 
30: ---
31: 
32: ## Reason 2: Binary size and distribution
33: 
34: gitswitch is distributed via Homebrew, a curl-to-shell install script, and direct binary download. Binary size and self-containment matter.
35: 
36: **Go's issue**: Go statically links its runtime. A simple CLI tool compiles to 4–15MB. The runtime is necessary and non-removable.
37: 
38: **Rust's answer**: With `lto = true`, `opt-level = "z"` (size-optimize), `strip = true`, and `codegen-units = 1`, a Rust binary for a tool of gitswitch's scope compiles to ~1.5–3MB. Fully self-contained. No dynamic runtime dependencies beyond the OS.
39: 
40: For the "minimalist developer" demographic that gitswitch targets — people who care about what's in their shell and dislike bloat — a small, auditable binary is a selling point.
41: 
42: ---
43: 
44: ## Reason 3: Terminal UI ecosystem
45: 
46: v0.2.0 ships a full-featured TUI (five screens, command palette, wizard forms, mouse support, three themes, arcade easter eggs). The Rust TUI ecosystem is currently the strongest in any language for building this kind of tool.
47: 
48: | Need | Go option | Rust option |
49: |---|---|---|
50: | TUI framework | Bubble Tea (Charm) | [ratatui](https://crates.io/crates/ratatui) — immediate-mode, widget gallery, modular backends |
51: | Animation / effects | Harmonica (springs only) | [tachyonfx](https://crates.io/crates/tachyonfx) — 50+ composable shader-like effects (coalesce, dissolve, evolve, slide) |
52: | Spinners | Manual / bubbles | [throbber-widgets-tui](https://crates.io/crates/throbber-widgets-tui) — customizable symbol sets |
53: | Big text / splash | Manual | [tui-big-text](https://crates.io/crates/tui-big-text) — block-letter rendering |
54: | Fuzzy palette | Manual | [nucleo-matcher](https://crates.io/crates/nucleo-matcher) — used by Helix, fastest available |
55: | Snapshot tests | vhs (external) | [insta](https://crates.io/crates/insta) — first-class inline snapshot testing |
56: | Panic safety | Manual terminal restore | [color-eyre](https://crates.io/crates/color-eyre) — panic hooks that restore raw mode |
57: 
58: The arcade layer specifically requires tachyonfx-class effects. No Go equivalent exists at that quality level. The hidden game roadmap (roguelike seeded from profile data, etc.) benefits from Rust's performance and ratatui's composable widget model.
59: 
60: ---
61: 
62: ## Reason 4: Security guarantees for sensitive paths
63: 
64: gitswitch touches:
65: - `~/.ssh/config` (SSH key paths, per-host blocks)
66: - `~/.gitconfig` (global identity)
67: - `~/.config/gh/hosts.yml` (GitHub PAT tokens)
68: - OS keychain (via the `keyring` crate)
69: - `~/.config/gitswitch/config.yaml` (profiles + pins)
70: 
71: **Go's concern**: Go is memory-safe for heap allocations but has no compile-time guarantees about resource lifetimes, data races without explicit mutexes, or secret-in-memory exposure. None of these are blocking issues in practice, but they require discipline.
72: 
73: **Rust's answer**: The borrow checker prevents data races at compile time. The `zeroize` crate integrates with secret-holding types to guarantee zero-on-drop — PAT tokens and passphrases for the encrypted-file fallback backend are wiped from memory on drop without the compiler optimizing it out. File permission checks and atomic writes (write-then-rename via `tempfile`) are idiomatic and the compiler enforces `Result` handling at every call site. Unhandled errors are a compile error, not a runtime surprise.
74: 
75: For a tool handling SSH keys and GitHub tokens, "the compiler made me handle the error" is a meaningful safety guarantee.
76: 
77: ---
78: 
79: ## Reason 5: The ecosystem is already there
80: 
81: The crates gitswitch needs are mature, production-used, and maintained:
82: 
83: | Crate | Purpose | Used by |
84: |---|---|---|
85: | [`clap`](https://crates.io/crates/clap) v4 | CLI argument parsing | Cargo itself, hundreds of major CLIs |
86: | [`serde`](https://crates.io/crates/serde) + [`yaml_serde`](https://crates.io/crates/yaml_serde) | Serialization | Ubiquitous |
87: | [`keyring`](https://crates.io/crates/keyring) v3 | OS keychain | Production keychain access |
88: | [`reqwest`](https://crates.io/crates/reqwest) | HTTP (GitHub API) | Ubiquitous |
89: | [`github-device-flow`](https://crates.io/crates/github-device-flow) | OAuth device flow | Purpose-built |
90: | [`ratatui`](https://crates.io/crates/ratatui) | TUI framework | Helix adjacent, massive community |
91: | [`directories`](https://crates.io/crates/directories) | XDG config paths | Cross-platform standard |
92: | [`assert_cmd`](https://crates.io/crates/assert_cmd) + [`insta`](https://crates.io/crates/insta) | Integration + snapshot tests | Standard Rust CLI testing stack |
93: 
94: No greenfield risk. All crates have existing production users.
95: 
96: ---
97: 
98: ## Where Go wins (acknowledged)
99: 
100: - **Faster to ship initially**: Go's simplicity meant gitswitch v0.1 was built and iterated quickly. That was the right call.
101: - **Cross-compilation**: Both are excellent; Go is slightly simpler (`GOOS=windows go build` vs `cross`).
102: - **Backend services**: If gitswitch ever grows a managed profiles API or team-preset hosting service, Go remains the better fit there.
103: 
104: The Go codebase is not a mistake — it got gitswitch to v0.1.20 with real users. Rust is the right choice for where the tool is going.
105: 
106: ---
107: 
108: ## The verdict
109: 
110: gitswitch lives in the terminal, fires on every directory change, handles secrets, and needs a polished animated TUI. Those constraints pick Rust. The migration plan is in [docs/rust/roadmap-migration.md](rust/roadmap-migration.md); the v0.2.0 feature plans are in [docs/rust/](rust/).
111: 
112: > "If gitswitch.dev were a cloud-hosted web app or API backend, Go would be the faster path to ship. But because gitswitch lives strictly in the user's terminal, interacts directly with local files, and intercepts shell navigation, Rust provides the performance, safety, and modern toolset required to make it a world-class developer utility."

(End of file - total 112 lines)
</content>
```

---

## Assistant (Plan · Kimi K2.7 Code · 23.0s)

_Thinking:_

Now I have a good picture. The user wants to:
1. Cancel the Rust migration plan (keep docs but mark as not happening / cancelled)
2. For v1 launch, stay in Go
3. Add direct OAuth with GitHub support
4. Implement `gs` as a short command alias for `gitswitch`

I need to present a revamped plan. Since this is plan mode, I should NOT edit files. I should present the plan to the user for approval, and ask clarifying questions where needed.

Let me think through what the new plan should look like:

**Key changes:**
- Rust migration is cancelled. v1.0 will ship the Go binary as default.
- Direct OAuth with GitHub becomes a v1.0 feature (was v0.3.0 deferred).
- `gs` command alias is added.
- v1.0 scope changes from "Rust stable" to "Go stable + OAuth + alias".

**Implications:**
- v0.2.0 remains Go with multi-host identity, session isolation, TUI refresh, etc.
- v0.3.0 may need reordering since OAuth moves to v1.0.
- v1.0 becomes "Production milestone (Go, stable)" instead of Rust.
- Need to update roadmap-full-timeline.md, roadmap-v0.3.0.md, why-rust.md, rust/roadmap-migration.md.
- Need to decide whether OAuth is device flow (as described in v0.3.0 doc) or web flow.
- Need to decide `gs` alias implementation: shell alias, symlink, or hardcoded command registration.

**Questions to ask the user:**
1. OAuth flow type: device flow (no browser redirect handling needed) or full web OAuth with localhost callback?
2. `gs` alias: should it be a shell alias installed by `gitswitch install`, a separate binary/symlink, or both?
3. v0.3.0 scope: since OAuth moves to v1.0, what stays in v0.3.0? Should we move anything up or keep it as is minus OAuth?
4. Rust docs: should we add a note/cancellation marker to existing docs, or create a new "rust-cancelled.md" doc explaining the decision?

Let me draft the plan clearly.

Here is the proposed revamp. No files have been edited yet — this is for your sign-off before I make changes.

---

## Proposed v1 Launch Plan Revamp

### What is changing

| Before | After |
|---|---|
| v1.0 ships a Rust rewrite as the default binary | v1.0 ships the Go binary as the stable default; Rust migration is cancelled |
| Direct GitHub OAuth is a v0.3.0 "long tail" feature | Direct GitHub OAuth is a **v1.0 launch requirement** |
| Only `gitswitch` command is supported | `gitswitch` **and** a short `gs` alias both invoke the tool |

### Why the change

- Rust migration is too risky relative to the existing shipping Go system.
- Direct OAuth removes the `gh` CLI dependency for authentication, which unblocks enterprise/headless users and makes first-time setup smoother.
- `gs` is a natural short alias that matches tools like `gh`, `glab`, etc.

---

## Revised timeline

### v0.1.x — Current (Go)
No change. Keep shipping.

### v0.2.0 — Next milestone (Go)
No change to scope:
- Multi-host identity
- Keychain-backed PATs
- Session isolation
- `gitswitch status`
- TUI refresh with 5 screens + command palette
- Pre-commit safety net
- Shell tab completion

### v0.3.0 — Power user / team features (Go)
Remove **Own OAuth flow** from this milestone. Remaining items:
- In-repo `.gitswitch.yaml`
- Team preset distribution (signed)
- Remote retargeting
- Host auto-detection
- Remote-URL-aware routing
- `git includeIf` support
- Ephemeral `gitswitch run`
- Profile export/import
- Windows native support improvements

Suggested new v0.3.0 sequencing:
1. In-repo `.gitswitch.yaml`
2. Host auto-detection
3. Team preset apply (with signing)
4. Remote retargeting
5. Remote-URL-aware routing
6. Everything else

### v1.0 — Production milestone (Go, stable)
This becomes the stable Go launch.

**New v1.0 requirements:**
1. **Direct GitHub OAuth flow**
   - Device flow against github.com (no client secret)
   - Command: `gitswitch login [--host HOST] [--client-id ID]`
   - Scopes: `repo`, `read:user`, `user:email`, `gist`, `workflow`
   - Tokens stored in OS keychain
   - Supports GHES with user-supplied `client_id`
   - Falls back cleanly if `gh` CLI auth exists

2. **`gs` short command**
   - `gs` is equivalent to `gitswitch`
   - Implementation options (need your pick):
     a. Shell alias injected by `gitswitch install`
     b. A tiny `gs` shim binary/symlink in the same release
     c. Both

3. **HTTPS credential helper**
   - `gitswitch` as a git `credential.helper`
   - Uses the active profile's stored PAT
   - Completes full auth stack: SSH + HTTPS + gh CLI

4. **Stability guarantees**
   - Semver commitment
   - Stable config format (`~/.config/gitswitch/config.yaml`)
   - Stable keychain key convention

### Rust migration
- Cancelled as an active engineering track.
- Existing Rust docs (`why-rust.md`, `rust/`) are **kept** but marked as cancelled/archived.

---

## Files to update

1. `docs/.never_gonna_give_you_up/roadmap-full-timeline.md`
   - Remove Rust migration track
   - Recast v1.0 as Go stable + OAuth + `gs`

2. `docs/.never_gonna_give_you_up/roadmap-v0.3.0.md`
   - Move OAuth section out
   - Update sequencing

3. `docs/.never_gonna_give_you_up/why-rust.md`
   - Add prominent "CANCELLED" header with rationale

4. `docs/.never_gonna_give_you_up/rust/roadmap-migration.md`
   - Mark cancelled/archived

5. New doc (optional): `docs/.never_gonna_give_you_up/rust-cancelled.md`
   - Brief decision record explaining why Rust was abandoned

---

## Questions before I edit

1. **OAuth flow**: Device flow only (simpler, no local HTTP server), or do you also want a localhost web-redirect flow?
2. **`gs` alias**: Shell alias only, a shim binary/symlink, or both?
3. **v0.3.0 scope**: Are you happy with OAuth leaving v0.3.0 entirely, or do you want a partial/fallback OAuth earlier?
4. **Rust docs**: Add a cancellation marker to existing docs, or create a separate decision record file?

Reply with your choices and I’ll make the updates.

---

## User

1. explain more
2. explain more
3. i didnt understand, explain me my own original plan and this new one
4. yes just add the cancel market on top for all docs and to the folder name

---

## Assistant (Plan · Kimi K2.7 Code · 33.5s)

_Thinking:_

The user wants me to explain more about my questions, and explain their original plan vs the new one. They're also confirming point 4: add cancel markers on top of all Rust docs and to the folder name.

However, I'm in plan mode and cannot edit files. I should explain things in detail and ask clarifying questions, but NOT make any changes yet.

Let me address each point clearly:

1. OAuth flow - explain device flow vs web-redirect flow, pros/cons
2. gs alias - explain shell alias vs shim binary/symlink vs both, pros/cons
3. Explain original plan vs new plan more clearly
4. Acknowledge they want cancellation markers on Rust docs and folder rename

I should be careful: the user said "yes just add the cancel market on top for all docs and to the folder name" - they want this done. But I'm in plan mode and cannot edit. I should explain that I'll do this when we exit plan mode, and present the plan.

Actually, looking at the system reminder again: "The user indicated that they do not want you to execute yet". So I should not execute. I should present the plan and ask if they're ready to proceed.

Wait, but the user is also asking for explanations. So my response should primarily be explanations, plus a note that I'll make the edits once they confirm.

Let me think through the explanations:

**OAuth device flow:**
- User runs `gitswitch login`
- gitswitch prints a one-time code and a URL (github.com/login/device)
- User opens browser, enters code, authorizes
- GitHub returns an access token
- No local server needed
- Good for SSH/headless environments
- GitHub device flow doesn't use refresh tokens by default

**OAuth web flow (localhost redirect):**
- User runs `gitswitch login`
- gitswitch opens a browser to GitHub OAuth authorize URL
- GitHub redirects to localhost callback
- gitswitch receives the code and exchanges for token
- Requires opening a browser and local HTTP server
- Can use refresh tokens
- Better UX for desktop users, worse for headless

**gs alias options:**
a) Shell alias: add `alias gs=gitswitch` to shell rc file during `gitswitch install`
   - Pros: trivial, no extra binary
   - Cons: doesn't work until shell reload, not available in subprocesses unless shell config loaded, doesn't show in `which gs`
   
b) Shim binary/symlink: ship a `gs` executable that calls `gitswitch`
   - Pros: works immediately after install, shows in PATH, can be used in scripts
   - Cons: slightly more packaging work, needs to handle being a symlink or small wrapper
   
c) Both: install shell alias AND ship `gs` binary
   - Pros: works everywhere
   - Cons: potential confusion if they diverge

**Original plan vs new plan:**
Original: v0.2 Go → parallel Rust rewrite → v0.3 Go and Rust opt-in → v1.0 Rust default. OAuth was deferred to v0.3 as fallback. Only `gitswitch` command.

New: v0.2 Go → v0.3 Go (no OAuth) → v1.0 Go stable with OAuth built-in and `gs` alias. Rust migration cancelled. OAuth promoted to v1.0 core feature.

Let me write a clear response.

Got it — here are the detailed explanations, then the next step.

---

## 1. OAuth flow: device flow vs. web-redirect flow

### Option A: Device flow (what was already sketched in `roadmap-v0.3.0.md`)

How it works:
1. User runs `gitswitch login`.
2. gitswitch asks GitHub for a device code.
3. Terminal prints:  
   `Open https://github.com/login/device and enter code: ABCD-EFGH`
4. User opens browser, pastes code, clicks Authorize.
5. gitswitch polls GitHub, gets the PAT, stores it in keychain.

**Pros**
- No local HTTP server needed.
- Works over SSH, in containers, in headless CI, on remote machines.
- Matches how `gh auth login` works, so users already know the pattern.
- Simpler to implement and secure (no open localhost port).

**Cons**
- GitHub device flow does not normally issue refresh tokens, so re-auth is needed when the token expires.
- Slightly more manual steps for the user than a one-click browser redirect.

### Option B: Web-redirect flow (localhost callback)

How it works:
1. User runs `gitswitch login`.
2. gitswitch starts a temporary local HTTP server on `localhost:8123`.
3. gitswitch opens the user's default browser to GitHub's OAuth authorize URL, with `redirect_uri=localhost:8123`.
4. User clicks Authorize in the browser.
5. Browser redirects back to `localhost:8123` with a code.
6. gitswitch exchanges the code for a PAT and stores it.

**Pros**
- One-click-ish experience on desktop.
- Can use refresh tokens for long-lived sessions.
- Feels more "modern app".

**Cons**
- Requires a browser and the ability to open it from the terminal.
- Does not work well over SSH/headless setups.
- Opening a localhost port, even briefly, can trigger firewalls or security tools.
- More code to write securely (CSRF state param, port collision handling, etc.).

**My recommendation for gitswitch**: **Device flow only for v1.0**. It fits the terminal-native, SSH-friendly identity tool you are building. You can always add web-redirect later if desktop users complain.

---

## 2. `gs` alias: shell alias vs. shim binary

### Option A: Shell alias only

During `gitswitch install`, append `alias gs=gitswitch` to `~/.zshrc`, `~/.bashrc`, `~/.config/fish/config.fish`, etc.

**Pros**
- Trivial to implement.
- No extra release artifact.
- Works instantly in new shell windows.

**Cons**
- Not available until the user reloads their shell config or opens a new terminal.
- Does not work inside scripts unless the script sources the shell config.
- `which gs` returns nothing, which confuses some users.
- Fish and other shells have different alias syntax.

### Option B: Shim binary / symlink

Ship a `gs` executable alongside `gitswitch`. It can be:
- A tiny compiled binary that `exec`s `gitswitch` with all args forwarded.
- A symlink to the `gitswitch` binary, where the binary detects `argv[0] == "gs"` and behaves the same.

**Pros**
- Works immediately after install.
- Shows up in `which gs`.
- Works in scripts and subshells.
- Consistent across bash/zsh/fish.

**Cons**
- Packaging slightly more complex (two binaries or one binary + symlink).
- If `gitswitch` is installed manually without the installer, the user must also create `gs`.

### Option C: Both

Ship `gs` as a real binary/symlink **and** also offer to add a shell alias for users who want tab-completion customization.

**Pros**
- Maximum compatibility.

**Cons**
- Two ways to do the same thing can cause confusion if they ever diverge.

**My recommendation**: **Option B — shim binary/symlink**. It is the most robust and matches how `gh`, `kubectl` (`k`), and other CLI tools ship short aliases. You can detect `argv[0]` inside the same binary so you do not even need a second executable file.

---

## 3. Original plan vs. new plan

### Your original plan

```
v0.1.x (Go, today)
    ↓
v0.2.0 (Go) ──→ start Rust rewrite in parallel
    ↓
v0.3.0 (Go + Rust opt-in binary)
    ↓
v1.0 (Rust becomes the default; Go deprecated)
```

- **Direct GitHub OAuth** was planned for v0.3.0 as a fallback for people who do not use `gh`.
- **Command name** was just `gitswitch`.

### The new plan

```
v0.1.x (Go, today)
    ↓
v0.2.0 (Go) ──→ NO Rust rewrite
    ↓
v0.3.0 (Go only, no OAuth here)
    ↓
v1.0 (Go becomes the stable production release)
```

- **Direct GitHub OAuth** is moved up and becomes a **core v1.0 launch feature**, not a fallback.
- **Rust migration is cancelled** — the Go codebase ships v1.0.
- **Command names** are `gitswitch` and `gs`.

### Why this matters

| Area | Original plan | New plan |
|---|---|---|
| Engineering risk | High — rewriting in Rust while shipping Go features | Low — keep improving the working Go code |
| Time to v1.0 | Longer, because Rust parity had to be reached first | Shorter, because Go is already production-shaped |
| Authentication story | v1.0 still leans on `gh` CLI | v1.0 owns auth end-to-end via OAuth |
| User onboarding | Install `gitswitch`, then maybe install `gh`, then auth | Install `gitswitch`, run `gs login`, done |
| Terminal ergonomics | Type `gitswitch` every time | Type `gs` for quick switches |

---

## 4. Rust doc cancellation markers

You want:
- A "CANCELLED" / "ARCHIVED" marker at the top of every Rust-related doc.
- The `rust/` folder renamed to something like `rust-cancelled/` or `rust-archived/`.

I will do this once we move out of plan mode. Specifically:

- `why-rust.md` → add cancellation banner at the top.
- `rust/roadmap-migration.md` → add cancellation banner at the top.
- Any other files inside `rust/` → add cancellation banner at the top.
- Rename folder `rust/` → `rust-cancelled/` (or `rust-archived/`, your pick).

---

## Final questions before I edit

1. **OAuth**: Confirm **device flow only** for v1.0?
2. **`gs`**: Confirm **shim binary / symlink** (same binary detects `argv[0]`)?
3. **Rust folder name**: Do you prefer `rust-cancelled/` or `rust-archived/`?
4. **Ready to proceed?** Say yes and I will apply the revamp to the docs.

---

## User

1. we will go device flow
2. option B it is
3. so basically we are on version 0.1.22 , idk where to go right now... theres so many other things like PAT Token management and session isolations for multi agent orchestration that i want, so I need to come up with a versioning plan properly
4. yup

so finish cancelling the plan, and then let's go and write new docs with new plans for these. 

---

## Assistant (Build · Kimi K2.7 Code · 8.9s)

_Thinking:_

Good, now I'm in build mode and can make changes. The user wants:
1. Device flow for OAuth
2. Shim binary/symlink for `gs`
3. Help come up with a proper versioning plan since they're at 0.1.22 and have many features they want
4. Rename rust folder and cancel the rust plan

They want me to:
1. Finish cancelling the rust plan (add markers, rename folder)
2. Write new docs with new plans

Let me start by doing the cancellation work. First, let me check the current state of files and the rust folder.

I need to:
1. Read all rust-related docs to add cancellation markers
2. Rename rust/ folder to rust-cancelled/
3. Update roadmap-full-timeline.md
4. Update roadmap-v0.3.0.md
5. Update why-rust.md
6. Create new plan docs

But first, I need to understand the current codebase and version. The user says they're on v0.1.22. Let me check the actual version in the code.

Let me start with a todo list since this is multi-step.

**Tool: todowrite**

**Input:**
```json
{
  "todos": [
    {
      "content": "Add cancellation markers to Rust docs and rename rust/ folder",
      "status": "in_progress",
      "priority": "high"
    },
    {
      "content": "Revamp roadmap-full-timeline.md with new v1 plan (Go stable, OAuth device flow, gs alias)",
      "status": "pending",
      "priority": "high"
    },
    {
      "content": "Revamp roadmap-v0.3.0.md to remove OAuth and reflect new sequencing",
      "status": "pending",
      "priority": "high"
    },
    {
      "content": "Check current project version and existing features to ground the new versioning plan",
      "status": "pending",
      "priority": "high"
    },
    {
      "content": "Draft new versioning/roadmap docs for v0.2.x, v0.3.x, v1.0 with PAT mgmt, session isolation, OAuth, gs",
      "status": "pending",
      "priority": "high"
    }
  ]
}
```

**Output:**
```
[
  {
    "content": "Add cancellation markers to Rust docs and rename rust/ folder",
    "status": "in_progress",
    "priority": "high"
  },
  {
    "content": "Revamp roadmap-full-timeline.md with new v1 plan (Go stable, OAuth device flow, gs alias)",
    "status": "pending",
    "priority": "high"
  },
  {
    "content": "Revamp roadmap-v0.3.0.md to remove OAuth and reflect new sequencing",
    "status": "pending",
    "priority": "high"
  },
  {
    "content": "Check current project version and existing features to ground the new versioning plan",
    "status": "pending",
    "priority": "high"
  },
  {
    "content": "Draft new versioning/roadmap docs for v0.2.x, v0.3.x, v1.0 with PAT mgmt, session isolation, OAuth, gs",
    "status": "pending",
    "priority": "high"
  }
]
```

---

## Assistant (Build · Kimi K2.7 Code · 3.2s)

_Thinking:_

Now let me read the rust folder contents and the main docs I need to edit. I'll read them in parallel.

**Tool: read**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust"
}
```

**Output:**
```
<path>/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust</path>
<type>directory</type>
<entries>
roadmap-migration.md
roadmap-v0.2.0-isolation.md
roadmap-v0.2.0-tui.md
roadmap-v0.2.0.md
roadmap-v0.3.0.md

(5 entries)
</entries>
```

**Tool: read**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust/roadmap-migration.md"
}
```

**Output:**
```
<path>/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust/roadmap-migration.md</path>
<type>file</type>
<content>
1: # gitswitch — Go → Rust Migration Plan
2: 
3: ## 1. Context & motivation
4: 
5: gitswitch today is ~4.5k LoC of Go (`cmd/gitswitch` + `internal/{git,history,shell,storage,tui,version}`), a single static binary distributed per OS/arch. The v0.2.0 roadmap (multi-host, session isolation, TUI refresh) roughly triples surface area: keychain backends, encrypted-file fallback, ratatui-grade animation, OAuth import, per-host SSH config, per-remote pin maps, session env emission, doctor probes. Pushing all of that through Go gets us there; pushing it through Rust gets us there with stronger ergonomics on every load-bearing axis.
6: 
7: **Why Rust now:**
8: 
9: - **TUI render loop**: [`ratatui`](https://docs.rs/ratatui) (forked from `tui-rs` in 2023, modularised in 0.30) with the [`ratatui-crossterm`](https://crates.io/crates/ratatui-crossterm) backend and [`tachyonfx`](https://crates.io/crates/tachyonfx) effects (50+ shader-like effects: materialize, dissolve, evolve, slide, coalesce) outclasses any Go TUI stack for the arcade aesthetic we want preserved. Bubble Tea is fine; ratatui + tachyonfx is in a different aesthetic league for splash + transitions.
10: - **Memory safety + correctness**: keychain backends, encrypted-file fallback, SSH config block surgery, and YAML in-place editing of `gh hosts.yml` all benefit from a borrow-checked, exhaustive-match language. The Go version already needs platform build tags for file locks ([`history/lock_unix.go`](../../internal/history/lock_unix.go) / [`lock_windows.go`](../../internal/history/lock_windows.go)); Rust's `cfg` story is at least as good.
11: - **Single static binary parity**: `cargo build --release` with `lto=true`, `strip=true`, `codegen-units=1`, MUSL on Linux yields a 2–4 MB statically linked binary — comparable to the Go output (currently ~7 MB).
12: - **OS-native integrations**: [`keyring`](https://docs.rs/keyring) v3 ships first-party `apple-native` (Security framework) and `windows-native` (WinCred) features, plus libsecret/secret-service on Linux. Materially more mature than [`zalando/go-keyring`](https://github.com/zalando/go-keyring) (which shells out to `security` on macOS).
13: - **CLI ergonomics**: [`clap`](https://docs.rs/clap) v4 derive macros (`Parser`, `Subcommand`, `#[clap(flatten)]`, `#[clap(global = true)]`) collapse the 550-line cobra wiring in [`cmd/gitswitch/main.go`](../../cmd/gitswitch/main.go) into a single typed enum tree.
14: 
15: This is not a "rewrite because Rust." It's "v0.2.0 doubles the surface; do it in the language whose ecosystem already shipped the answers."
16: 
17: ---
18: 
19: ## 2. Migration strategy — clean rewrite, parallel binary
20: 
21: **Decision: clean rewrite, not an FFI sandwich.**
22: 
23: Turborepo's [Rust-Go-Rust sandwich](https://turbo.build/blog/turbopack-benchmarks) (Rust entry point cgo'ing into Go for unported commands) is the right call for a 100k+ LoC build tool under active feature development. gitswitch is ~4.5k LoC, frozen-feature for the duration of the port (v0.2.0 ships in Rust, not Go). The sandwich's overhead — cgo build matrix, dual-runtime distribution, drift between two implementations — costs more than the rewrite saves.
24: 
25: Instead: **ship the Rust rewrite as `gitswitch-rs` alongside the Go `gitswitch`** until parity is certified. Same config dir, same shell hook contract, same YAML schema. Users opt-in by aliasing or by installing the rs build through a separate install flag. When the parity corpus passes 100%, the binaries swap names and the Go tree gets archived.
26: 
27: ### Parity discipline
28: 
29: The rewrite is gated by a behavior-locked **parity corpus** — a directory of shell scripts each running a scenario (`add → switch → list → status`, `import gh → pin --remote → push`, etc.) against both binaries with `HOME=$(mktemp -d)`, then diffing stdout/stderr/exit-code and the resulting `~/.config/gitswitch/` tree byte-for-byte. The corpus is the spec; passing it is the cutover gate.
30: 
31: ### Why not incremental piece-by-piece
32: 
33: The corrode.dev [Rust migration guide](https://corrode.dev/blog/migrating-from-go-to-rust/) recommends "piece-by-piece" porting through a long parallel-runtime period for codebases where the cost of feature-freezing the original is high. gitswitch's cost of feature-freezing is essentially zero — there are no contributors blocked on v0.1.x development. The piece-by-piece pattern's main payoff (de-risking) is achieved more cheaply here by the parallel-binary scheme above without the cgo overhead.
34: 
35: ---
36: 
37: ## 3. Phase plan (6–8 weeks)
38: 
39: ### Phase 0 — Skeleton (week 1)
40: 
41: - Cargo workspace (`crates/` layout below).
42: - `clap` v4 command surface mirroring the current Go cobra tree (every subcommand stubbed, returns `unimplemented!()`).
43: - CI matrix: `macos-latest`, `ubuntu-latest`, `windows-latest`, plus `aarch64-unknown-linux-musl` via [`cross`](https://github.com/cross-rs/cross).
44: - `cargo deny` + `cargo audit` wired into CI.
45: - Release scaffolding: [`cargo-dist`](https://opensource.axo.dev/cargo-dist/) for cross-built artifacts under the `gitswitch-rs` name.
46: - **Deliverable**: `gitswitch-rs --help` prints the v0.1 command surface; nothing else works.
47: 
48: ### Phase 1 — Core data + storage (weeks 2–3)
49: 
50: - `gitswitch-core`: `Profile`, `HostCreds`, `Pin`, `Prefs` types matching the [v0.2.0 YAML schema](../roadmap-v0.2.0.md#storage-file-format) verbatim.
51: - YAML I/O via [`yaml-rust2`](https://docs.rs/yaml-rust2) (preserves comments for the `init`-scaffolded config; `serde_yaml` is [deprecated since 2024-03](https://github.com/dtolnay/serde-yaml/releases/tag/0.9.34)).
52: - v1 JSON → v2 YAML migration (port [`internal/storage/store.go`](../../internal/storage/store.go) `legacyProfile` path 1:1; same `profiles.json.v1.bak` backup behavior).
53: - `gitswitch-secrets`: `keyring` v3 wrapper with the same `Store` trait shape as the v0.2.0 plan, encrypted-file fallback via [`chacha20poly1305`](https://docs.rs/chacha20poly1305) + [`scrypt`](https://docs.rs/scrypt).
54: - Snapshot tests with [`insta`](https://docs.rs/insta) on YAML output across a fixture corpus.
55: - **Deliverable**: `gitswitch list`, `add`, `remove`, `current` all work against `config.yaml` on all three OSes.
56: 
57: ### Phase 2 — Git + shell integration (weeks 3–4)
58: 
59: - `gitswitch-git`: shell-out wrapper around `git config` (port [`internal/git/git.go`](../../internal/git/git.go) verbatim — `std::process::Command`, no libgit2 dependency).
60: - SSH config block manager: idempotent `# gitswitch:<hostname>` block writer in `~/.ssh/config` (new for v0.2.0, no Go precursor).
61: - `gitswitch-shell`: zsh / bash / fish snippet installers (port [`internal/shell/shell.go`](../../internal/shell/shell.go)). **The snippets written to user shell rc files do not change** — the hook contract is part of the public API.
62: - `hook-check` parity (port [`internal/shell/shell.go`](../../internal/shell/shell.go) detection + stale-install warnings).
63: - **Deliverable**: `switch`, `install`, `uninstall`, `hook-check`, `record`, `recommend` reach v0.1 parity.
64: 
65: ### Phase 3 — TUI (weeks 4–5)
66: 
67: - `gitswitch-tui` (`ratatui` + `ratatui-crossterm`) implementing the [5-screen layout from the TUI roadmap](../roadmap-v0.2.0-tui.md): Profiles, Hosts, Repos, Session, Doctor.
68: - `gitswitch-arcade` housing the bespoke arcade pieces:
69:   - splash via [`tachyonfx`](https://docs.rs/tachyonfx) `coalesce` + `dissolve` effects (replaces the Bubble Tea splash in [`internal/tui/view.go`](../../internal/tui/view.go))
70:   - spinner via [`throbber-widgets-tui`](https://docs.rs/throbber-widgets-tui)
71:   - CRT-style border as a custom [`Widget`](https://docs.rs/ratatui/latest/ratatui/widgets/trait.Widget.html) — scanline overlay + chromatic-aberration corner blocks
72:   - 12-theme palette ported from [`internal/tui/styles.go`](../../internal/tui/styles.go)
73: - Golden-frame tests: `insta::assert_snapshot!` on `Buffer`-rendered output per screen + theme.
74: - **Deliverable**: TUI reaches v0.1 visual parity, plus the v0.2.0 screen IA from the TUI roadmap.
75: 
76: ### Phase 4 — Session isolation (weeks 5–6)
77: 
78: - `gitswitch-env`: session env builder shared by the `shell` / `env` commands and the TUI's exit-eval path (see [session isolation roadmap](../roadmap-v0.2.0-isolation.md)).
79: - Shell wrapper function (zsh/bash/fish) emitted by `install` — same surface as the Go plan, `eval`-style env emission via `--emit-env`.
80: - Scope-aware command dispatch: `switch`/`current`/`status`/`pin` read `GITSWITCH_SESSION` and route to session vs global state; `--global` / `--session` flags override.
81: - **Deliverable**: parallel terminals can hold distinct active identities; v0.2.0 isolation scenarios pass.
82: 
83: ### Phase 5 — gh import + OAuth + doctor (weeks 6–7)
84: 
85: - `gitswitch-gh`: parse `~/.config/gh/hosts.yml` directly + fallback to `gh auth status --show-token`; `SyncHost` for in-place `hosts.yml` edits.
86: - GitHub device flow via [`github-device-flow`](https://crates.io/crates/github-device-flow) (or [`oauth2`](https://docs.rs/oauth2) if we need the more general primitive) — opt-in path for users without `gh` installed.
87: - HTTP: [`ureq`](https://docs.rs/ureq) (sync, ~400 KB) for the `/user` + `/user/emails` enrichment calls and token-validity probes. No async runtime needed for the import flow; keeps startup fast.
88: - `gitswitch doctor`: walks every profile × host, probes keychain availability, token validity (HTTP 200 on `/user`), SSH key file presence + mode, `gh hosts.yml` sync state, shell hook version.
89: - `gitswitch status`: per-repo narrative diagnostic (port the [v0.2.0 status layout](../roadmap-v0.2.0.md#gitswitch-status--repo-narrative-diagnostic)).
90: - **Deliverable**: `import gh`, `import git`, `host add/remove/list`, `status`, `doctor` all functional.
91: 
92: ### Phase 6 — Parity certification + cutover (weeks 7–8)
93: 
94: - Run the parity corpus (see §8) against `gitswitch` (Go) and `gitswitch-rs` (Rust) on macOS + Linux + Windows CI runners.
95: - Diff every scenario's stdout/stderr/exit-code + `~/.config/gitswitch/` tree + `~/.ssh/config` + `~/.config/gh/hosts.yml`.
96: - Fix divergences until corpus is 100% green on all three platforms.
97: - Tag `gitswitch-rs v0.2.0-rc1`, dogfood for one week.
98: - Rename: archive `cmd/gitswitch/` (Go) to `legacy/` branch; rename Rust binary `gitswitch-rs` → `gitswitch`; ship as `v0.2.0`.
99: 
100: ---
101: 
102: ## 4. Workspace layout
103: 
104: ```
105: git-switcher/
106: ├── Cargo.toml                  # workspace manifest, [profile.release] lto/strip/codegen-units=1
107: ├── crates/
108: │   ├── gitswitch/              # bin crate — main.rs, clap dispatch only
109: │   ├── gitswitch-core/         # Profile/Host/Pin types, YAML load/save, v1 migration
110: │   ├── gitswitch-secrets/      # keyring wrapper + encrypted-file fallback
111: │   ├── gitswitch-git/          # git/gh shell-outs, SSH config block manager
112: │   ├── gitswitch-gh/           # GitHub API client, device flow, hosts.yml editor
113: │   ├── gitswitch-env/          # session env builder (shared shell/env/TUI exit path)
114: │   ├── gitswitch-tui/          # ratatui app, screens, theme, view models
115: │   ├── gitswitch-arcade/       # splash, tachyonfx effects, CRT widget
116: │   └── gitswitch-shell/        # install/uninstall/hook detection (zsh/bash/fish)
117: ├── tests/
118: │   ├── parity/                 # parity corpus shell scripts
119: │   ├── snapshots/              # insta snapshots (YAML + TUI frames)
120: │   └── integration/            # cargo test integration suites
121: ├── legacy/                     # archived Go tree, kept for reference until 0.2.1
122: └── docs/rust/                  # this folder
123: ```
124: 
125: ### Key dependency graph
126: 
127: ```
128: gitswitch
129:   ├── clap = "4"
130:   ├── gitswitch-core
131:   ├── gitswitch-secrets
132:   ├── gitswitch-git
133:   ├── gitswitch-gh        (feature = "gh")
134:   ├── gitswitch-env
135:   ├── gitswitch-tui       (feature = "tui")
136:   └── gitswitch-shell
137: 
138: gitswitch-tui
139:   ├── ratatui = "0.30"
140:   ├── ratatui-crossterm
141:   ├── gitswitch-arcade
142:   ├── throbber-widgets-tui
143:   └── (consumes all data crates read-only via &dyn traits)
144: 
145: gitswitch-arcade
146:   └── tachyonfx
147: ```
148: 
149: `tokio` is **not** in the workspace by default — the entire TUI loop and CLI dispatch is sync. `ureq` covers HTTP. If we later need concurrent doctor probes we add `tokio` behind a feature flag in `gitswitch-gh` only.
150: 
151: ### Core data types (target shape)
152: 
153: ```rust
154: // gitswitch-core/src/profile.rs
155: #[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
156: pub struct Profile {
157:     pub nickname: String,
158:     pub user_name: String,
159:     pub email: String,
160:     #[serde(default, skip_serializing_if = "String::is_empty")]
161:     pub sign_key: String,
162:     #[serde(default)]
163:     pub active: bool,
164:     #[serde(default, skip_serializing_if = "Option::is_none")]
165:     pub source: Option<ProfileSource>,
166:     #[serde(default)]
167:     pub hosts: BTreeMap<String, HostCreds>,
168: }
169: 
170: #[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
171: pub struct HostCreds {
172:     #[serde(default, skip_serializing_if = "String::is_empty")]
173:     pub ssh_key: String,
174:     #[serde(default, skip_serializing_if = "String::is_empty")]
175:     pub gh_user: String,
176:     #[serde(default, skip_serializing_if = "String::is_empty")]
177:     pub token_ref: String,
178:     #[serde(default, skip_serializing_if = "std::ops::Not::not")]
179:     pub gh_synced: bool,
180: }
181: 
182: #[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
183: #[serde(rename_all = "kebab-case")]
184: pub enum ProfileSource { Manual, GhImport, Imported }
185: 
186: #[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
187: pub struct Pin {
188:     #[serde(default, skip_serializing_if = "Option::is_none")]
189:     pub default: Option<String>,
190:     #[serde(default, skip_serializing_if = "BTreeMap::is_empty")]
191:     pub remotes: BTreeMap<String, String>,
192: }
193: 
194: #[derive(Debug, Clone, Serialize, Deserialize)]
195: pub struct Config {
196:     pub version: u32,
197:     #[serde(default)]
198:     pub prefs: Prefs,
199:     #[serde(default)]
200:     pub profiles: Vec<Profile>,
201:     #[serde(default)]
202:     pub pins: BTreeMap<String, Pin>,
203: }
204: ```
205: 
206: `BTreeMap` (not `HashMap`) for hosts/pins/remotes so YAML serialization is order-stable — required for the parity corpus byte-diff.
207: 
208: ### Secrets trait
209: 
210: ```rust
211: // gitswitch-secrets/src/lib.rs
212: pub trait SecretStore: Send + Sync {
213:     fn set(&self, reference: &str, token: &str) -> Result<()>;
214:     fn get(&self, reference: &str) -> Result<String>;
215:     fn delete(&self, reference: &str) -> Result<()>;
216:     fn available(&self) -> bool;
217: }
218: 
219: pub struct KeyringStore;          // keyring v3, OS-native
220: pub struct EncryptedFileStore {   // xchacha20poly1305 + scrypt
221:     path: PathBuf,
222:     passphrase: Zeroizing<String>,
223: }
224: ```
225: 
226: ### CLI surface
227: 
228: ```rust
229: // gitswitch/src/cli.rs
230: #[derive(Parser)]
231: #[command(name = "gitswitch", version, about)]
232: pub struct Cli {
233:     /// Run against session scope (default: auto-detect via GITSWITCH_SESSION)
234:     #[arg(long, global = true)]
235:     pub session: bool,
236:     /// Run against global scope, overriding session detection
237:     #[arg(long, global = true)]
238:     pub global: bool,
239:     /// Emit shell-eval env updates instead of taking effect directly
240:     #[arg(long, global = true, hide = true)]
241:     pub emit_env: bool,
242: 
243:     #[command(subcommand)]
244:     pub command: Option<Command>,
245: 
246:     /// Quick switch: gitswitch <nickname>
247:     pub nickname: Option<String>,
248: }
249: 
250: #[derive(Subcommand)]
251: pub enum Command {
252:     Add(AddArgs), Switch(SwitchArgs), List, Remove(RemoveArgs),
253:     Current, Init, Version, Upgrade,
254:     Pin(PinArgs), Unpin(UnpinArgs),
255:     Record, Recommend,
256:     Install, Uninstall, HookCheck,
257:     Claude, Pacman,
258:     // v0.2.0 additions
259:     Status, Doctor,
260:     #[command(subcommand)] Host(HostCommand),
261:     #[command(subcommand)] Import(ImportCommand),
262:     Shell, Env,
263: }
264: ```
265: 
266: Compare to the 550-line cobra wiring in [`cmd/gitswitch/main.go`](../../cmd/gitswitch/main.go) — the Rust version lands closer to 150 lines of dispatch + a per-command module each.
267: 
268: ---
269: 
270: ## 5. Crate selection table
271: 
272: | Concern              | Crate                                                                                              | One-line justification                                                                                                            |
273: | -------------------- | -------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------- |
274: | CLI parsing          | [`clap`](https://docs.rs/clap) v4 (derive)                                                         | Industry default; derive macro collapses 500+ lines of cobra wiring; `global = true` + `flatten` cover our scope/shared flags.    |
275: | TUI                  | [`ratatui`](https://docs.rs/ratatui) 0.30+ + [`ratatui-crossterm`](https://crates.io/crates/ratatui-crossterm) | Active fork of tui-rs, modular backends, immediate-mode model matches the v0.1 Bubble Tea approach.                               |
276: | TUI effects          | [`tachyonfx`](https://crates.io/crates/tachyonfx)                                                  | 50+ shader-like effects; what makes the arcade splash possible without hand-rolling animation primitives.                         |
277: | Spinners             | [`throbber-widgets-tui`](https://docs.rs/throbber-widgets-tui)                                     | Drop-in spinner widget, multiple styles; doctor + import flows.                                                                   |
278: | Keychain             | [`keyring`](https://docs.rs/keyring) v3                                                            | First-party `apple-native`, `windows-native` features; libsecret on Linux; mature error model.                                    |
279: | Encrypted fallback   | [`chacha20poly1305`](https://docs.rs/chacha20poly1305) + [`scrypt`](https://docs.rs/scrypt) + [`zeroize`](https://docs.rs/zeroize) | RustCrypto AEAD + standard KDF; `Zeroizing<String>` for passphrase handling.                                                      |
280: | YAML                 | [`yaml-rust2`](https://docs.rs/yaml-rust2)                                                         | Maintained fork; preserves comments (needed for `init`-scaffolded config). `serde_yaml` is deprecated.                            |
281: | Serialization        | [`serde`](https://docs.rs/serde) + [`serde_json`](https://docs.rs/serde_json)                      | Standard; JSON for legacy v1 migration + `--emit-env` IPC.                                                                        |
282: | HTTP                 | [`ureq`](https://docs.rs/ureq)                                                                     | Sync, ~400 KB, rustls-backed; no tokio dependency on a CLI that's otherwise fully sync.                                           |
283: | OAuth device flow    | [`github-device-flow`](https://crates.io/crates/github-device-flow) (fallback: [`oauth2`](https://docs.rs/oauth2)) | GitHub-specific; opt-in for users without `gh` installed.                                                                         |
284: | Terminal             | [`crossterm`](https://docs.rs/crossterm)                                                           | Cross-platform raw mode; bundled via `ratatui-crossterm`.                                                                         |
285: | Errors               | [`thiserror`](https://docs.rs/thiserror) (libs) + [`anyhow`](https://docs.rs/anyhow) (bin)         | Standard split; `thiserror` for typed library errors, `anyhow` for `main` chaining.                                               |
286: | Logging              | [`tracing`](https://docs.rs/tracing) + [`tracing-subscriber`](https://docs.rs/tracing-subscriber)  | `GITSWITCH_LOG=debug` for diagnostic flows without polluting normal output.                                                       |
287: | Tests                | [`insta`](https://docs.rs/insta) + [`assert_cmd`](https://docs.rs/assert_cmd) + [`tempfile`](https://docs.rs/tempfile) | Snapshot YAML + TUI frames; spawn binary in tmp HOME for integration tests.                                                       |
288: | Cross-compile        | [`cross`](https://github.com/cross-rs/cross)                                                       | Same cross-compile UX `cargo` ships natively, with Docker-isolated toolchains for MUSL Linux + Windows from any host.             |
289: | Releases             | [`cargo-dist`](https://opensource.axo.dev/cargo-dist/)                                             | Replaces the current `Makefile` / GitHub Actions ad-hoc release flow; auto-publishes per-OS tarballs + Homebrew formula update.   |
290: 
291: ---
292: 
293: ## 6. Data migration
294: 
295: **The v0.2.0 YAML schema is the boundary.** The Rust binary must read and write `~/.config/gitswitch/config.yaml` byte-identically to the Go binary's v0.2.0 output. Migration from v0.1's `profiles.json` is implemented behavior-for-behavior in `gitswitch-core`:
296: 
297: 1. On `Config::load()`, if `config.yaml` is absent and `profiles.json` exists:
298: 2. Read JSON via a `LegacyProfile` struct mirroring the v0.1 [`legacyProfile`](../../internal/storage/store.go#L21) (handles the `Name` → `UserName` rename and missing-`Nickname` fallback).
299: 3. Promote each v1 profile to the hosts-map shape; populate `Hosts["github.com"]` from the old `SSHKey` + `GHUser` fields.
300: 4. Write `config.yaml` (mode `0o600`).
301: 5. Rename `profiles.json` → `profiles.json.v1.bak`.
302: 
303: Test corpus: a folder of v1 `profiles.json` fixtures (empty, one-profile, multi-profile, missing-nickname, name-vs-username) — each must produce the same v2 `config.yaml` as the Go binary's migration, byte-for-byte.
304: 
305: `history.json` (per-repo pin storage, [`internal/history/history.go`](../../internal/history/history.go)) is similarly migrated into `config.yaml`'s `pins:` block as part of the same v1→v2 step. Locking semantics ([`flock` on Unix, no-op on Windows](../../internal/history/lock_unix.go)) port to [`fs2`](https://docs.rs/fs2) `FileExt::lock_exclusive` cross-platform.
306: 
307: ---
308: 
309: ## 7. Compatibility / risk register
310: 
311: | Risk                                         | Mitigation                                                                                                                                                                                                                                                            |
312: | -------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
313: | Shell hook script churn                      | The zsh / bash / fish snippets written by `install` are **frozen**. Only the binary they call changes. Snippet content asserted byte-identical to Go output by a parity test.                                                                                         |
314: | Cross-compile gaps                           | `cross` covers `aarch64-unknown-linux-musl`, `x86_64-unknown-linux-musl`, `x86_64-pc-windows-msvc`. macOS arm64 + x86_64 built natively on macOS runners. Same matrix the Go binary publishes today.                                                                  |
315: | Binary size regression                       | Release profile: `lto = "fat"`, `codegen-units = 1`, `strip = true`, `panic = "abort"`. Expected output 2–4 MB statically linked. CI gate: fail if release binary exceeds 1.2× Go binary size on any platform.                                                        |
316: | Async vs sync mismatch                       | TUI inner loop is fully sync (no tokio). HTTP via `ureq` (sync). OAuth + doctor concurrent probes opt into `tokio` only behind `gitswitch-gh`'s `async` feature flag — currently unused.                                                                              |
317: | Keychain on headless Linux                   | `SecretStore::available()` probes at runtime; falls back to `EncryptedFileStore` (xchacha20poly1305 + scrypt KDF, `~/.config/gitswitch/secrets.enc`, mode 0600). `GITSWITCH_SECRETS_BACKEND=file` env var forces fallback for CI.                                      |
318: | `gh hosts.yml` in-place edits                | `yaml-rust2` preserves comments + key order; we read → mutate the relevant `host:` entry → write back, never round-tripping unrelated user state. Snapshot tests cover preservation of unrelated hosts.                                                                |
319: | TUI snapshot brittleness                     | `insta` snapshots key off rendered `Buffer`s with normalized terminal width × height. Effect-frame snapshots taken at deterministic time steps (no real-time animation in test fixtures).                                                                             |
320: | Two-binary distribution during transition    | Releases tagged `0.2.0-rcN-rs` ship the Rust binary as `gitswitch-rs`; Go releases continue under `gitswitch` until parity is signed off. Both Homebrew taps and install scripts updated to install both side-by-side; users opt in via `--rs` flag on install script. |
321: | Skill / Claude integration drift             | The `claude` subcommand emits the same `SKILL.md` payload (embedded via [`include_str!`](https://doc.rust-lang.org/std/macro.include_str.html)) byte-identical to today's [`//go:embed skill/SKILL.md`](../../cmd/gitswitch/main.go#L20).                              |
322: | Windows `flock` semantics                    | Use [`fs2`](https://docs.rs/fs2) which wraps `LockFileEx` on Windows + `flock` on Unix — single API, real locking on both platforms (the current Go [`lock_windows.go`](../../internal/history/lock_windows.go) is a no-op; we get an upgrade for free).               |
323: 
324: ---
325: 
326: ## 8. Testing strategy
327: 
328: ### Unit (per crate)
329: - `gitswitch-core`: YAML round-trip, v1 migration on the fixture corpus, `Pin` resolution under default-vs-remote precedence.
330: - `gitswitch-secrets`: both backends pass the same `SecretStore` contract test suite. Encrypted-file backend tested for tamper detection (modified ciphertext → decryption error, not silent corruption).
331: - `gitswitch-git`: SSH config block management — add/update/remove idempotent across `~/.ssh/config` shapes (empty, multiple unrelated hosts, prior gitswitch blocks present).
332: - `gitswitch-gh`: `hosts.yml` parsing + in-place edit preserves unrelated host entries.
333: - `gitswitch-shell`: snippet emitter produces byte-identical zsh/bash/fish output to the Go binary's `install` (asserted via fixtures captured from `gitswitch` v0.1.x at port-start).
334: 
335: ### Integration (per-binary, in tmp HOME)
336: Use [`assert_cmd`](https://docs.rs/assert_cmd) + [`tempfile`](https://docs.rs/tempfile). Each test sets `HOME=tmpdir`, invokes the binary, asserts:
337: - stdout/stderr match expected
338: - exit code matches
339: - `~/.config/gitswitch/config.yaml` matches snapshot
340: - `~/.ssh/config` matches snapshot
341: - keychain entries (where mocked) match expectations
342: 
343: ### Parity corpus (cross-binary diff)
344: 
345: Shell script suite under `tests/parity/`:
346: ```
347: tests/parity/
348: ├── 01-add-list-switch.sh
349: ├── 02-multi-host-pin.sh
350: ├── 03-import-gh.sh
351: ├── 04-session-isolation.sh
352: ├── 05-v1-migration.sh
353: ├── 06-doctor-stale-token.sh
354: └── ...
355: ```
356: 
357: Each script:
358: 1. `export HOME=$(mktemp -d)`
359: 2. Run scenario against `GITSWITCH=$GO_BIN` → capture outputs + state tree → `golden/`
360: 3. Repeat with `GITSWITCH=$RS_BIN` → capture → `actual/`
361: 4. `diff -r golden/ actual/` — non-zero exit fails the run.
362: 
363: Runs on all three OS runners in CI on every PR to the Rust tree until cutover.
364: 
365: ### TUI snapshot tests
366: - `insta::assert_snapshot!` on `ratatui::buffer::Buffer` rendered output per screen × theme × terminal size.
367: - Effects pinned to deterministic frame indices (e.g. `splash @ frame=30`).
368: - Visual regression detected by snapshot diff; reviewer accepts via `cargo insta review`.
369: 
370: ---
371: 
372: ## 9. Cutover criteria
373: 
374: Rename `gitswitch-rs` → `gitswitch` and archive the Go tree only when **all** of the following hold:
375: 
376: 1. **Parity corpus 100% green** on macOS 13+, Ubuntu 22.04+, and Windows Server 2022 CI runners.
377: 2. **Binary size** within 1.2× of Go binary on every platform (`ls -l target/*/release/gitswitch` vs Go output).
378: 3. **Startup time** `time gitswitch list` ≤ Go binary + 50 ms on cold cache, measured across 10 runs on each platform.
379: 4. **TUI frame rate** sustained ≥ 60 fps idle (`tachyonfx` effects + base render) measured via `ratatui` instrumentation on a 120×40 terminal.
380: 5. **Doctor full sweep** ≤ Go binary's `doctor` time + 100 ms across a 5-profile-3-host fixture.
381: 6. **No critical issues** filed against `0.2.0-rcN-rs` builds in the one-week dogfood window.
382: 7. **Install scripts updated** — Homebrew formula, install.sh, all distribution channels swap to Rust binary atomically with the rename.
383: 8. **Migration smoke** — every fixture in `tests/parity/05-v1-migration*.sh` produces byte-identical `config.yaml` to the Go binary.
384: 
385: Cutover commit:
386: - Move `cmd/`, `internal/` → `legacy/`.
387: - Rename Rust binary artifact `gitswitch-rs` → `gitswitch` in `Cargo.toml` + `cargo-dist` config.
388: - Update `README.md`, install scripts, Homebrew formula.
389: - Tag `v0.2.0`.
390: 
391: Legacy Go tree kept on the `legacy/go-v0.1.x` branch indefinitely. Deleted from `main` 30 days post-cutover if no regressions surface.
392: 
393: ---
394: 
395: ## 10. Companion plans
396: 
397: The Rust workspace consumes the same feature decisions as the Go roadmap — only the implementation language changes. Per-feature Rust roadmaps live alongside this file:
398: 
399: - [`docs/rust/roadmap-v0.2.0.md`](roadmap-v0.2.0.md) — multi-host identity in Rust (Profile / HostCreds / Pin types, keychain + encrypted-file `SecretStore`, gh-driven import, SSH config block manager). Sister doc to [`docs/roadmap-v0.2.0.md`](../roadmap-v0.2.0.md).
400: - [`docs/rust/roadmap-v0.2.0-isolation.md`](roadmap-v0.2.0-isolation.md) — session isolation via env vars + shell wrapper, scope-aware command dispatch (`--global` / `--session` flags), `--emit-env` IPC. Sister doc to [`docs/roadmap-v0.2.0-isolation.md`](../roadmap-v0.2.0-isolation.md).
401: - [`docs/rust/roadmap-v0.2.0-tui.md`](roadmap-v0.2.0-tui.md) — ratatui + tachyonfx implementation of the 5-screen IA, arcade theme as a tachyonfx-driven layer, custom CRT widget. Sister doc to [`docs/roadmap-v0.2.0-tui.md`](../roadmap-v0.2.0-tui.md).
402: - [`docs/rust/roadmap-v0.3.0.md`](roadmap-v0.3.0.md) — deferred items (in-repo YAML config, team preset distribution + signing, own OAuth as primary path, host auto-detection). Sister doc to [`docs/roadmap-v0.3.0.md`](../roadmap-v0.3.0.md).
403: 
404: The Go roadmaps under [`docs/`](../) remain the canonical feature specs; the Rust roadmaps are pure implementation notes layered on top.

(End of file - total 404 lines)
</content>
```

---

## Assistant (Build · Kimi K2.7 Code · 10.7s)

_Thinking:_

Now I have all the info I need. Let me start by adding cancellation markers and renaming the folder. Since I'm in build mode, I can use bash commands. But the instructions say to use specialized tools for file operations. For renaming a folder, bash mv is appropriate. For editing files, I should use edit tool.

Actually, I need to be careful. The instructions say:
- "DO NOT use it for file operations (reading, writing, editing, searching, finding files) - use the specialized tools for this instead."

So for renaming the folder, I should use bash mv? The bash tool says "This tool is for terminal operations like git, npm, docker, etc. DO NOT use it for file operations". Hmm, but renaming a directory is a file operation. However, there's no dedicated "rename" tool. I think using bash mv for directory rename is acceptable since it's a terminal operation and there's no other tool. But let me think...

Actually, the bash tool examples include `mkdir`, which is a file operation. So directory operations with bash are probably fine.

For editing files, I must use the edit tool. I'll add cancellation banners to the top of each file.

Let me start with renaming the folder and editing files. First, let me check if there are relative links I need to update. The rust docs reference each other with relative links. If I rename the folder to `rust-cancelled/`, those links will break unless I update them. But since they're cancelled docs, maybe it's OK to leave them as historical references. Actually, better to update them to point to the new location or note they're archived.

Wait, the user said "add the cancel market on top for all docs and to the folder name". So I'll:
1. Rename `rust/` to `rust-cancelled/`
2. Add a cancellation banner at the top of each file in the folder
3. Add cancellation banner to `why-rust.md`
4. Update any internal relative links? The links would break after rename. But since these are archived docs, I could either update the links or leave them. Better to update links within the cancelled folder to reflect new path, or at least not break them. Actually, since they're all in the same folder, relative links between them (`roadmap-v0.2.0.md`) won't break. Only links going outside (`../../`) might still work since the parent changed from `rust/` to `rust-cancelled/` but the relative path depth is the same. Links to `../roadmap-v0.2.0.md` would still work. So no link updates needed.

Let me do the rename and edits.

**Tool: bash**

**Input:**
```json
{
  "command": "mv \"docs/.never_gonna_give_you_up/rust\" \"docs/.never_gonna_give_you_up/rust-cancelled\"",
  "description": "Rename rust folder to rust-cancelled"
}
```

**Output:**
```
(no output)
```

---

## Assistant (Build · Kimi K2.7 Code · 6.1s)

_Thinking:_

Now I need to add cancellation banners to all files in rust-cancelled/ and to why-rust.md. Let me read each file's first few lines to make precise edits.

Actually, I already have the contents. Let me use the edit tool to prepend a cancellation banner to each file. I need to read each file first? The edit tool says "You must use your Read tool at least once in the conversation before editing." I've already read roadmap-migration.md. But I haven't read the other rust files. I should read them before editing.

Let me read the remaining rust-cancelled files in parallel.

**Tool: read**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-v0.2.0.md"
}
```

**Output:**
```
<path>/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-v0.2.0.md</path>
<type>file</type>
<content>
1: # gitswitch v0.2.0 (Rust) — Multi-Host Identity Plan
2: 
3: ## Context
4: 
5: This is the Rust port of the v0.2.0 multi-host plan. Behavior parity with the Go v0.1 → v0.2 migration path is preserved: same on-disk schema, same command surface, same Vaishnavi multi-host scenario as the acceptance target. The Rust workspace structure follows [roadmap-migration.md](roadmap-migration.md).
6: 
7: Current gitswitch (v0.1.x, Go) treats identity as `(name, email, ssh_key, gpg_key, gh_user)` against a single implicit host (`github.com`). an enterprise multi-host use case — and other enterprise multi-host environments — surfaces a gap: developers operate against **two GitHub hosts simultaneously** with different PAT tokens, SSH keys, and `gh` contexts.
8: 
9: v0.2.0 closes the gap by making **host a first-class dimension of identity**, **leveraging `gh auth` as the import source** (rather than building our own OAuth), and adding per-repo diagnostics (`gitswitch status`) plus multi-profile pinning. Team-preset YAML distribution, our own OAuth flow, and in-repo config files are explicitly v0.3.0 scope — see [roadmap-v0.3.0.md](roadmap-v0.3.0.md).
10: 
11: ### Locked design decisions
12: 
13: 1. **Per-host credentials map inside one profile** (not one-profile-per-host).
14: 2. **PAT storage**: OS keychain primary, gh `hosts.yml` synced on switch, encrypted-file fallback for headless/CI.
15: 3. **OAuth strategy**: don't roll our own. **Import from `gh auth status` / `~/.config/gh/hosts.yml`** — gh CLI users already authenticated their accounts; we discover them, create gitswitch profiles, then layer our customizations (per-host SSH keys, per-repo pinning) on top.
16: 4. **No remote URL rewrite on switch** — host is per-remote; gitswitch makes credentials correct for whichever host git resolves at push time.
17: 5. **`gitswitch init` repurposed**: scaffolds local YAML config in `~/.config/gitswitch/` (not in-repo — in-repo YAML is v0.3.0) and imports current git config + gh auth state in one shot.
18: 6. **New `gitswitch status`** — per-repo narrative diagnostic.
19: 7. **Multi-pin** — a repo can pin different profiles per remote (origin → corp, fork → personal).
20: 
21: ---
22: 
23: ## Profile model (storage layer)
24: 
25: Crate: `gitswitch-core` (workspace member; see [roadmap-migration.md](roadmap-migration.md)).
26: 
27: ```rust
28: use serde::{Deserialize, Serialize};
29: use std::collections::HashMap;
30: 
31: #[derive(Debug, Clone, Serialize, Deserialize)]
32: pub struct Profile {
33:     pub nickname: String,
34:     pub user_name: String,
35:     pub email: String,
36:     #[serde(skip_serializing_if = "String::is_empty", default)]
37:     pub sign_key: String,
38:     #[serde(default)]
39:     pub active: bool,
40:     #[serde(skip_serializing_if = "Option::is_none")]
41:     pub source: Option<ProfileSource>,
42:     pub hosts: HashMap<String, HostCreds>,
43: }
44: 
45: #[derive(Debug, Clone, Serialize, Deserialize)]
46: pub struct HostCreds {
47:     #[serde(skip_serializing_if = "String::is_empty", default)]
48:     pub ssh_key: String,
49:     #[serde(skip_serializing_if = "String::is_empty", default)]
50:     pub gh_user: String,
51:     #[serde(skip_serializing_if = "String::is_empty", default)]
52:     pub token_ref: String,
53:     #[serde(default)]
54:     pub gh_synced: bool,
55: }
56: 
57: #[derive(Debug, Clone, Serialize, Deserialize)]
58: #[serde(rename_all = "kebab-case")]
59: pub enum ProfileSource {
60:     Manual,
61:     GhImport,
62:     Imported,
63: }
64: ```
65: 
66: The legacy top-level `ssh_key` / `gh_user` fields (v1) collapse into `hosts["github.com"]` during the v1→v2 upgrade. `Option<ProfileSource>` keeps the YAML clean — manual profiles without an explicit source serialize without the key.
67: 
68: Serialization uses [`serde`](https://crates.io/crates/serde) + [`yaml_serde`](https://crates.io/crates/yaml_serde) (NOT `serde_yaml`, which is deprecated upstream).
69: 
70: ### Storage file format
71: 
72: Single file at `~/.config/gitswitch/config.yaml` (XDG resolution via [`directories`](https://crates.io/crates/directories) — `ProjectDirs::from("dev", "aksisonline", "gitswitch")`).
73: 
74: ```rust
75: #[derive(Debug, Clone, Serialize, Deserialize)]
76: pub struct Config {
77:     pub version: u32,
78:     #[serde(default)]
79:     pub prefs: Prefs,
80:     #[serde(default)]
81:     pub profiles: Vec<Profile>,
82:     #[serde(default)]
83:     pub pins: HashMap<String, Pin>,
84: }
85: 
86: #[derive(Debug, Clone, Default, Serialize, Deserialize)]
87: pub struct Prefs {
88:     #[serde(default)]
89:     pub color_theme: u8,
90: }
91: ```
92: 
93: Wire format identical to the Go plan:
94: 
95: ```yaml
96: version: 2
97: prefs:
98:   color_theme: 0
99: profiles:
100:   - nickname: personal
101:     user_name: Abhiram Kanna
102:     email: abhiram@personal.dev
103:     sign_key: 4A...
104:     active: true
105:     source: manual
106:     hosts:
107:       github.com:
108:         ssh_key: ~/.ssh/id_personal
109:         gh_user: aksisonline
110:         token_ref: gitswitch:personal:github.com
111:         gh_synced: true
112:   - nickname: corp
113:     user_name: Abhiram Kanna
114:     email: abhiram@corp.com
115:     active: false
116:     source: gh-import
117:     hosts:
118:       github.com:
119:         ssh_key: ~/.ssh/id_personal
120:         gh_user: aksisonline
121:         token_ref: gitswitch:personal:github.com
122:       github.corp.com:
123:         ssh_key: ~/.ssh/id_corp
124:         gh_user: aksisonline-pp
125:         token_ref: gitswitch:corp:github.corp.com
126: pins:
127:   /Users/abhi/code/corp-platform:
128:     default: corp
129:     remotes:
130:       fork: personal
131: ```
132: 
133: Atomic writes use [`tempfile::NamedTempFile`](https://crates.io/crates/tempfile) in the same directory + `persist()` (rename is atomic on POSIX). File mode 0600 set via `std::os::unix::fs::PermissionsExt` before persist.
134: 
135: Concurrent-write protection: [`fs2::FileExt::try_lock_exclusive`](https://crates.io/crates/fs2) on a sibling `.lock` file for the duration of `Config::save()`.
136: 
137: ### Migration (v1 JSON → v2 YAML)
138: 
139: `gitswitch_core::migrate::v1_to_v2()`:
140: 
141: 1. If `config.yaml` exists, no-op.
142: 2. If `profiles.json` exists, parse with [`serde_json`](https://crates.io/crates/serde_json) using a `V1Profile` struct mirroring the Go v0.1 shape.
143: 3. For each `V1Profile`, build a v2 `Profile` with `hosts.insert("github.com".into(), HostCreds { ssh_key, gh_user, ..Default::default() })`.
144: 4. `Config::save()` writes `config.yaml` (mode 0600).
145: 5. Rename `profiles.json` → `profiles.json.v1.bak` via `std::fs::rename` (not delete — safety net).
146: 6. Idempotent: second run sees `config.yaml` exists, returns Ok early.
147: 
148: Errors propagate as `thiserror`-derived `MigrationError`; the binary wraps in `anyhow::Context` for user-readable messages.
149: 
150: ---
151: 
152: ## gh-driven import (replaces own-OAuth path)
153: 
154: New crate: `gitswitch-gh`. Module `gitswitch_gh::import`:
155: 
156: ```rust
157: use serde::{Deserialize, Serialize};
158: 
159: #[derive(Debug, Clone, Serialize, Deserialize)]
160: pub struct ImportedAccount {
161:     pub hostname: String,
162:     pub login: String,
163:     pub token: String,
164:     pub scopes: Vec<String>,
165:     pub source: ImportSource,
166: }
167: 
168: #[derive(Debug, Clone, Copy, Serialize, Deserialize)]
169: #[serde(rename_all = "kebab-case")]
170: pub enum ImportSource {
171:     HostsYml,
172:     GhAuthStatus,
173: }
174: 
175: pub fn discover() -> anyhow::Result<Vec<ImportedAccount>>;
176: pub fn enrich_account(acct: &ImportedAccount) -> anyhow::Result<gitswitch_core::Profile>;
177: ```
178: 
179: - `discover()` first reads `~/.config/gh/hosts.yml` directly via `yaml_serde::from_str` (preferred — no shell-out, faster, deterministic). Returns one entry per `(host, user)` pair.
180: - Fallback: if `hosts.yml` is missing or malformed, locate the `gh` binary via [`which::which`](https://crates.io/crates/which), shell out via `std::process::Command::new("gh").args(["auth", "status", "--show-token"])`, parse stderr line-by-line. (gh writes status to stderr.)
181: - `enrich_account()` issues `GET https://<api>/user` and `GET https://<api>/user/emails` with [`reqwest::blocking::Client`](https://crates.io/crates/reqwest) (`Authorization: Bearer <token>` + `User-Agent: gitswitch/0.2`). Picks the primary verified email or the `<id>+<login>@users.noreply.github.com` form. JSON responses parsed with `serde_json`.
182: - API host derivation: `github.com` → `api.github.com`; enterprise host `github.example.com` → `github.example.com/api/v3`.
183: 
184: ### `gitswitch import gh` flow
185: 
186: 1. Run `discover()`.
187: 2. For each account, prompt: `Create profile "<login>" for host "<hostname>"? [Y/n]`.
188: 3. On accept: build profile (`source: ProfileSource::GhImport`), call `secrets_store.set("gitswitch:<nickname>:<hostname>", &token)`, set `gh_synced = true`. Optionally prompt for SSH key path; if provided, upsert SSH config block.
189: 4. If two accounts share a login across hosts, offer to merge into one multi-host profile (single `Profile` with two `hosts` entries).
190: 
191: This subsumes the originally proposed OAuth device flow. Target audience already runs `gh auth login` for both `github.com` and any enterprise host — we ride on that.
192: 
193: ---
194: 
195: ## Keychain integration
196: 
197: New crate: `gitswitch-secrets`. Trait + two backends.
198: 
199: ```rust
200: pub trait SecretStore: Send + Sync {
201:     fn set(&self, refkey: &str, token: &str) -> anyhow::Result<()>;
202:     fn get(&self, refkey: &str) -> anyhow::Result<String>;
203:     fn delete(&self, refkey: &str) -> anyhow::Result<()>;
204:     fn available(&self) -> bool;
205: }
206: 
207: pub struct KeyringStore;          // wraps keyring::Entry
208: pub struct EncryptedFileStore {   // chacha20poly1305 + scrypt
209:     path: std::path::PathBuf,
210: }
211: 
212: pub fn pick() -> Box<dyn SecretStore> {
213:     if std::env::var("GITSWITCH_SECRETS_BACKEND").as_deref() == Ok("file") {
214:         return Box::new(EncryptedFileStore::default());
215:     }
216:     let kr = KeyringStore;
217:     if kr.available() { Box::new(kr) } else { Box::new(EncryptedFileStore::default()) }
218: }
219: ```
220: 
221: Dependencies:
222: - [`keyring`](https://crates.io/crates/keyring) v3 with `apple-native` (macOS Keychain) and `windows-native` (Credential Manager) features. Linux uses the default `linux-native-sync-persistent` backend (libsecret via D-Bus). Service name: `dev.aksisonline.gitswitch`; account: the `refkey`.
223: - [`chacha20poly1305`](https://crates.io/crates/chacha20poly1305) (XChaCha20-Poly1305 AEAD) for the encrypted-file backend.
224: - [`scrypt`](https://crates.io/crates/scrypt) for passphrase-derived keys (cost params: `log_n=15, r=8, p=1`).
225: - Random nonce per write via [`rand`](https://crates.io/crates/rand) `OsRng`.
226: 
227: `EncryptedFileStore` layout (`~/.config/gitswitch/secrets.enc`, mode 0600):
228: ```
229: | 16-byte salt | 24-byte nonce | ciphertext (serde_json map: refkey → token) |
230: ```
231: Passphrase prompted on first write, cached for the process via [`zeroize::Zeroizing<String>`](https://crates.io/crates/zeroize); `GITSWITCH_SECRETS_PASSPHRASE` env var honored for CI.
232: 
233: Token ref convention: `gitswitch:<nickname>:<hostname>`. Deterministic — `doctor` walks profiles and probes each.
234: 
235: ---
236: 
237: ## gh CLI sync on switch
238: 
239: Module `gitswitch_gh::sync`:
240: 
241: ```rust
242: pub fn sync_host(hostname: &str, gh_user: &str, token: &str) -> anyhow::Result<()>;
243: ```
244: 
245: Reads `~/.config/gh/hosts.yml` with `yaml_serde` into a `BTreeMap<String, serde_yaml_value::Value>`-style intermediate (preserve unknown keys for hosts we don't manage), upserts the entry for `hostname` with `user`, `oauth_token`, `git_protocol: ssh`, and writes atomically (tempfile + rename, mode 0600).
246: 
247: Called from `apply_profile()` per host in the active profile. Best-effort: if `which::which("gh").is_err()` or `hosts.yml` parent dir is missing, log + skip (don't abort the switch).
248: 
249: This keeps `gh` and `git push` in lockstep — the long-standing promise of gitswitch — across multiple hosts.
250: 
251: ---
252: 
253: ## SSH config management
254: 
255: New crate: `gitswitch-git`. Module `gitswitch_git::ssh_config`.
256: 
257: ```rust
258: pub fn upsert_block(hostname: &str, identity_file: &Path) -> anyhow::Result<()>;
259: pub fn remove_block(hostname: &str) -> anyhow::Result<()>;
260: pub fn has_block(hostname: &str) -> anyhow::Result<bool>;
261: ```
262: 
263: Idempotent per-host block in `~/.ssh/config`:
264: 
265: ```
266: # gitswitch:github.corp.com
267: Host github.corp.com
268:     IdentityFile ~/.ssh/id_corp
269:     IdentitiesOnly yes
270: ```
271: 
272: Implementation: byte-level line scan. No SSH config parser needed — we only own blocks marked with our magic comment. Algorithm:
273: 
274: 1. Read file (create if missing, mode 0600).
275: 2. Walk lines; identify our block by leading `# gitswitch:<hostname>` marker; the block ends at the next blank line or next `# gitswitch:` marker or next non-indented `Host ` directive.
276: 3. Rebuild file: original content minus the old block (if any) + new block appended.
277: 4. Atomic write via `tempfile::NamedTempFile::persist`.
278: 
279: Drops the global `core.sshCommand` approach in the Go v0.1 path. Per-host SSH config is correct and lets multi-remote pushes work without re-switching profiles between pushes.
280: 
281: ---
282: 
283: ## Switch flow
284: 
285: In `gitswitch-git`:
286: 
287: ```rust
288: pub fn apply_profile(p: &Profile, store: &dyn SecretStore) -> anyhow::Result<()>;
289: ```
290: 
291: 1. Run `git config --global user.name <p.user_name>`, `user.email <p.email>`, `user.signingkey <p.sign_key>` via `std::process::Command::new("git")`.
292: 2. For each `(hostname, creds)` in `p.hosts`:
293:    - SSH: `ssh_config::upsert_block(hostname, &creds.ssh_key)`.
294:    - Token: `store.get(&creds.token_ref)` → `gitswitch_gh::sync::sync_host(hostname, &creds.gh_user, &token)`.
295: 3. Update history + pins (delegated to `gitswitch-core::history`).
296: 
297: Legacy single-key SSH path retained only for migration tests; new flows don't call it.
298: 
299: ---
300: 
301: ## Per-repo pinning — extended to multi-pin
302: 
303: `gitswitch-core::history`:
304: 
305: ```rust
306: #[derive(Debug, Clone, Default, Serialize, Deserialize)]
307: pub struct Pin {
308:     #[serde(skip_serializing_if = "String::is_empty", default)]
309:     pub default: String,
310:     #[serde(skip_serializing_if = "HashMap::is_empty", default)]
311:     pub remotes: HashMap<String, String>,
312: }
313: 
314: pub fn resolve_pin(pin: &Pin, remote: Option<&str>) -> Option<&str>;
315: ```
316: 
317: A repo with `origin = github.corp.com:...` and `fork = github.com:me/...` can pin:
318: 
319: ```yaml
320: my-repo-key:
321:   default: corp
322:   remotes:
323:     fork: personal
324: ```
325: 
326: Backwards compatibility:
327: - `gitswitch recommend` and the `record` shell hook continue to work — if only `default` is set, behavior matches v0.1.
328: - If `remotes` is populated, the shell prompt segment shows the resolved pin for the current branch's upstream remote.
329: 
330: Command surface:
331: - `gitswitch pin <profile> [--remote NAME]` — without `--remote`, pins `default`; with `--remote`, pins that remote only.
332: - `gitswitch unpin [--remote NAME]` — mirrors.
333: 
334: ---
335: 
336: ## `gitswitch status` — repo narrative diagnostic
337: 
338: New command in `gitswitch-cli`. In a git repo, output a structured explanation:
339: 
340: ```
341: Repo: ~/code/corp-platform
342:   Branch: feature/x  →  origin/feature/x
343: 
344: Identity (global git config):
345:   user.name   = Abhiram Kanna
346:   user.email  = abhiram@corp.com
347:   signingkey  = 4A...                (gpg, present)
348:   Active profile: corp             (matches git config ✓)
349: 
350: Pins for this repo:
351:   default          → corp
352:   remote "fork"    → personal
353: 
354: Remotes:
355:   origin → git@github.corp.com:team/x.git
356:     Host: github.corp.com
357:     Profile that owns this host: corp
358:     SSH config block: present ✓        (~/.ssh/config: # gitswitch:github.corp.com)
359:     Identity file: ~/.ssh/id_corp    (exists, mode 600 ✓)
360:     gh auth: aksisonline-pp            (synced ✓, token valid as of 2026-05-28)
361:   fork → git@github.com:aksisonline/x.git
362:     Host: github.com
363:     Profile that owns this host: personal  (per fork-remote pin)
364:     SSH config block: present ✓
365:     Identity file: ~/.ssh/id_personal      (exists ✓)
366:     gh auth: aksisonline                   (synced ✓, token valid)
367: 
368: Recommendations: none.
369: ```
370: 
371: Terminal coloring via [`owo-colors`](https://crates.io/crates/owo-colors) (zero-dep, supports `NO_COLOR`). Git remote enumeration via shell-out to `git remote -v` parsed line-by-line (keeps parity with Go; avoids dragging `git2`/`libgit2` into the binary). If we later want native, swap to [`git2`](https://crates.io/crates/git2) behind a feature flag.
372: 
373: When a check fails, surface inline (`gh auth: stale token (401)`, `Identity file: missing`, `Pin references unknown profile "old-work"`). Composes `gitswitch_core::detect_active`, `gitswitch_core::recommend`, plus new `gitswitch_secrets` and `gitswitch_gh` probes.
374: 
375: `status` is per-repo; `doctor` (below) is global.
376: 
377: ---
378: 
379: ## Command surface
380: 
381: `gitswitch-cli` binary uses [`clap`](https://crates.io/crates/clap) v4 derive:
382: 
383: ```rust
384: use clap::{Parser, Subcommand};
385: 
386: #[derive(Parser)]
387: #[command(name = "gitswitch", version, about)]
388: struct Cli {
389:     #[command(subcommand)]
390:     cmd: Cmd,
391: }
392: 
393: #[derive(Subcommand)]
394: enum Cmd {
395:     /// Scaffold config.yaml + import current git/gh state.
396:     Init,
397: 
398:     /// Import accounts from gh or git config.
399:     Import {
400:         #[command(subcommand)]
401:         kind: ImportKind,
402:     },
403: 
404:     /// Manage hosts on an existing profile.
405:     Host {
406:         #[command(subcommand)]
407:         action: HostAction,
408:     },
409: 
410:     /// Per-repo narrative diagnostic.
411:     Status,
412: 
413:     /// Global health check across all profiles × hosts.
414:     Doctor {
415:         #[arg(long)]
416:         fix: bool,
417:     },
418: 
419:     /// Switch the active profile.
420:     Switch {
421:         nickname: String,
422:         #[arg(long)]
423:         global: bool,
424:         #[arg(long)]
425:         session: bool,
426:     },
427: 
428:     /// Add a new profile (repeatable --host for multi-host).
429:     Add {
430:         nickname: String,
431:         #[arg(long)]
432:         name: String,
433:         #[arg(long)]
434:         email: String,
435:         #[arg(long = "host", value_parser = parse_host_spec)]
436:         hosts: Vec<HostSpec>,
437:         #[arg(long)]
438:         sign_key: Option<String>,
439:     },
440: 
441:     Remove { nickname: String },
442:     List,
443:     Current,
444:     Record,
445:     Recommend,
446:     Pin {
447:         profile: String,
448:         #[arg(long)]
449:         remote: Option<String>,
450:     },
451:     Unpin {
452:         #[arg(long)]
453:         remote: Option<String>,
454:     },
455:     Install,
456:     Uninstall,
457:     Upgrade,
458:     HookCheck,
459:     Claude,
460:     Pacman,
461: }
462: 
463: #[derive(Subcommand)]
464: enum ImportKind {
465:     Gh,
466:     Git,
467: }
468: 
469: #[derive(Subcommand)]
470: enum HostAction {
471:     Add {
472:         profile: String,
473:         hostname: String,
474:         #[arg(long)]
475:         ssh_key: Option<String>,
476:         #[arg(long)]
477:         gh_user: Option<String>,
478:         #[arg(long)]
479:         token: Option<String>,
480:     },
481:     Remove {
482:         profile: String,
483:         hostname: String,
484:     },
485:     List {
486:         profile: String,
487:     },
488: }
489: ```
490: 
491: ### New
492: - `gitswitch import gh` — re-runnable gh-discovery flow.
493: - `gitswitch import git` — old `init` semantics (import current git config as new profile).
494: - `gitswitch host add <profile> <hostname> [--ssh-key] [--gh-user] [--token]`
495: - `gitswitch host remove <profile> <hostname>`
496: - `gitswitch host list <profile>`
497: - `gitswitch status` — per-repo narrative.
498: - `gitswitch doctor` — global health (every profile × every host: keychain reachable? token valid? SSH file present? gh hosts.yml in sync?).
499: 
500: ### Renamed / repurposed
501: - `gitswitch init` — ensure `~/.config/gitswitch/config.yaml` exists, import current git config as first profile, run `gh import` if `gh` is logged in. Old import-only behavior moves to `gitswitch import git`.
502: 
503: ### Extended
504: - `gitswitch add` — repeatable `--host hostname[=ssh_key_path]` flag for multi-host creation.
505: - `gitswitch pin` / `unpin` — gain `--remote NAME` flag for multi-pin.
506: 
507: ### Unchanged
508: `switch`, `list`, `remove`, `current`, `record`, `recommend`, `install`, `uninstall`, `version`, `upgrade`, `hook-check`, `claude`, `pacman`.
509: 
510: ---
511: 
512: ## Verification plan
513: 
514: ### Unit (per crate)
515: - `gitswitch-core::migrate` — v1 JSON → v2 YAML round-trip, idempotent on second run, backup file created.
516: - `gitswitch-secrets` — round-trip both backends (`KeyringStore` gated behind a feature flag in CI, `EncryptedFileStore` always tested).
517: - `gitswitch-git::ssh_config` — idempotent block management (add → add same → no-op; add → remove → absent; multiple hosts coexist).
518: - `gitswitch-gh::import` — parses fixture `hosts.yml` and fixture `gh auth status --show-token` stderr output.
519: 
520: ### Integration
521: Build the binary into a tmp HOME (via [`assert_cmd`](https://crates.io/crates/assert_cmd) + [`tempfile::TempDir`](https://crates.io/crates/tempfile)) and snapshot via [`insta`](https://crates.io/crates/insta):
522: - `gitswitch init` produces `config.yaml` + imports git config.
523: - `gitswitch import gh` (with a mocked `gh` binary on PATH) creates expected profiles.
524: - `gitswitch pin corp --remote origin` writes the expected pin structure.
525: - `gitswitch status` in a fake repo renders all sections.
526: 
527: ### End-to-end (Vaishnavi scenario reproduction)
528: Two profiles: `personal` (github.com) and `corp` (multi-host: github.com + github.corp.com).
529: 
530: In a repo with `origin = git@github.corp.com:team/x.git` and `fork = git@github.com:me/x.git`:
531: 1. `gitswitch pin corp`
532: 2. `gitswitch pin personal --remote fork`
533: 3. `git push origin` and `git push fork` both succeed via per-host SSH config.
534: 4. `gh repo view` works against both hosts after switching.
535: 5. `gitswitch status` reflects everything correctly.
536: 
537: ### Diagnostics
538: - Revoke a PAT externally → `gitswitch doctor` flags the expired token for the right host; `gitswitch status` shows it inline against the offending remote.
539: 
540: ### Migration safety
541: Smoke test: install Go v0.1.x, create two profiles, install Rust v0.2.0 binary, run `gitswitch list` → profiles intact, `profiles.json.v1.bak` present, `config.yaml` has hosts map populated with `github.com`.
542: 
543: ---
544: 
545: ## Critical files / crates
546: 
547: Workspace layout (full detail in [roadmap-migration.md](roadmap-migration.md)):
548: 
549: - **`gitswitch-core/`** — `Profile`, `HostCreds`, `Pin`, `Config`, YAML I/O, migration, history, `detect_active`, `recommend`.
550: - **`gitswitch-secrets/`** — `SecretStore` trait, `KeyringStore`, `EncryptedFileStore`.
551: - **`gitswitch-gh/`** — `import::discover`, `import::enrich_account`, `sync::sync_host`.
552: - **`gitswitch-git/`** — `apply_profile`, `ssh_config::{upsert_block, remove_block}`, git config shell-outs.
553: - **`gitswitch-cli/`** — `clap` command surface, command handlers, `status` rendering.
554: - **`gitswitch-tui/`** — (see TUI roadmap) presentation layer.
555: - [docs/diagnostic-design.md](../diagnostic-design.md) — extend with multi-host checks.
556: - [README.md](../../README.md) — refresh positioning to "who you are + which GitHub" on ship.
557: 
558: ## Companion plans
559: 
560: - [Session isolation](roadmap-v0.2.0-isolation.md) — env-scoped sessions so two terminals (or two AI agents) can hold two different active identities at the same time without touching global state. Shares the `gitswitch-secrets` crate with this plan.
561: - [TUI](roadmap-v0.2.0-tui.md) — Rust TUI on [`ratatui`](https://crates.io/crates/ratatui) (Bubble Tea analog). Full CLI feature parity (every command reachable via keybinds + `:` palette), arcade aesthetic preserved as a theme layer.
562: - [v0.3.0](roadmap-v0.3.0.md) — team preset distribution, own OAuth, in-repo config.
563: - [Migration](roadmap-migration.md) — workspace layout + Go → Rust port strategy.
564: 
565: ## Out of scope — see [v0.3.0](roadmap-v0.3.0.md)
566: 
567: In-repo YAML config, team preset distribution + signing, own OAuth flow, remote retargeting, host auto-detection.

(End of file - total 567 lines)
</content>
```

**Tool: read**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-v0.2.0-isolation.md"
}
```

**Output:**
```
<path>/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-v0.2.0-isolation.md</path>
<type>file</type>
<content>
1: # gitswitch v0.2.0 — Session Isolation (Rust Port)
2: 
3: > Rust translation of [roadmap-v0.2.0-isolation.md](../roadmap-v0.2.0-isolation.md). Mental model and all design decisions are identical; only the implementation idioms change. Read the Go plan first for full rationale.
4: 
5: ## Context
6: 
7: Same motivation as the Go plan:
8: 
9: - **AI agents** running two parallel coding sessions — one on a personal OSS repo (push as `aksisonline`), one on a client repo (push as `aksisonline-pp`). Shared `~/.gitconfig` → whichever profile switched last wins both.
10: - **Developers in multiple terminals / tmux panes** under different identities concurrently hit the same race.
11: - **CI / automation** wants parallel `git` ops under different identities without touching global state.
12: 
13: This plan ships **session isolation** in the Rust rewrite: a session is an isolated *scope* (a shell with its own active-identity env vars), not a locked profile. Inside a session you get **full gitswitch** — switch, pin, add hosts, status — all operating on session state. Built parallel to (not blocking on) the [multi-host work](roadmap-v0.2.0.md); shares the [`gitswitch-secrets`](roadmap-v0.2.0.md#crate-layout) keychain crate, otherwise orthogonal.
14: 
15: ## Mental model
16: 
17: ```
18: Global scope (single)              Session scope (many, parallel)
19: ───────────────────                ─────────────────────────────
20: ~/.gitconfig                       env vars in shell A: GITSWITCH_SESSION=A, identity=personal
21: ~/.config/gitswitch/config.yaml    env vars in shell B: GITSWITCH_SESSION=B, identity=corp
22:   ↑                                env vars in shell C: GITSWITCH_SESSION=C, identity=<none yet>
23:   └─ shared profile database ─────┘
24: ```
25: 
26: - **Profile database** (profiles + hosts + keychain refs): global, one source of truth.
27: - **Active identity**: scoped — global has one (in `~/.gitconfig`), each session has its own (in env vars).
28: - A session can be empty, single-profile, or freely switch between profiles during its lifetime. None of that touches global.
29: 
30: ## Locked design decisions
31: 
32: 1. **Mechanism = environment variables** + a shell function wrapper for in-session mutations. A session is a shell process with `GITSWITCH_SESSION` set; mutations happen via the wrapper, which `eval`s env updates into the current shell.
33: 2. **Session is a scope, not a profile.** `gitswitch shell` opens an empty session; switching profiles inside it is the normal flow.
34: 3. **All commands are scope-aware.** Inside a session, `switch` / `current` / `status` / `pin` read/write session state. Outside, they read/write global state. Flags `--global` / `--session` override.
35: 4. **Profile database is always global.** `add`, `remove`, `host add`, `host remove`, `import` mutate the shared YAML regardless of scope.
36: 5. **No daemon, no IPC, no shadow HOME.** Pure stateless: envs + the shell function.
37: 
38: ---
39: 
40: ## How it works
41: 
42: ### 1. `gitswitch install` adds a shell function
43: 
44: The shell-integration installer (the `gitswitch-shell` crate) gains a wrapper function emitted alongside the existing hook:
45: 
46: ```sh
47: # zsh / bash
48: gitswitch() {
49:   case "$1" in
50:     switch|use|shell)
51:       local out
52:       out="$(command gitswitch "$@" --emit-env 2>&1)" || { echo "$out" >&2; return 1; }
53:       eval "$out"
54:       ;;
55:     *)
56:       command gitswitch "$@"
57:       ;;
58:   esac
59: }
60: ```
61: 
62: ```fish
63: # fish
64: function gitswitch
65:     switch $argv[1]
66:         case switch use shell
67:             set -l out (command gitswitch $argv --emit-env 2>&1)
68:             or begin; echo $out >&2; return 1; end
69:             eval $out
70:         case '*'
71:             command gitswitch $argv
72:     end
73: end
74: ```
75: 
76: Without this wrapper, `gitswitch switch X` inside a session can't mutate the parent shell's envs. With it, every `switch` inside a session re-exports the right env block; outside a session, it falls back to global behavior.
77: 
78: Install snippet generation lives in `gitswitch-shell::install::wrapper_snippet(shell: Shell) -> String`.
79: 
80: ### 2. `gitswitch shell [profile]` opens a session
81: 
82: ```
83: $ gitswitch shell                           # empty session
84: [gs:session-a1b2] $ gitswitch switch corp
85: [gs:corp] $ git commit -m "..."           # committed as corp
86: [gs:corp] $ gitswitch switch personal
87: [gs:personal] $ git commit -m "..."         # committed as personal
88: [gs:personal] $ exit
89: $                                           # back to global scope
90: ```
91: 
92: Implementation in the `gitswitch` bin crate:
93: 
94: ```rust
95: use std::collections::HashMap;
96: use std::process::{Command, Stdio};
97: use anyhow::Result;
98: 
99: pub fn cmd_shell(profile: Option<String>, cmd: Option<String>, no_prompt: bool) -> Result<()> {
100:     let session_id = generate_session_id();              // e.g. "session-a1b2"
101:     let mut env: HashMap<String, String> = std::env::vars().collect();
102: 
103:     env.insert("GITSWITCH_SESSION".into(), session_id.clone());
104:     env.insert("GITSWITCH_SESSION_PID".into(), std::process::id().to_string());
105: 
106:     if let Some(p) = profile {
107:         let prof = storage::load_profile(&p)?;
108:         for ev in gitswitch_env::build(&prof)? {
109:             env.insert(ev.key, ev.value);
110:         }
111:         env.insert("GITSWITCH_SESSION_PROFILE".into(), p);
112:     }
113: 
114:     if !no_prompt {
115:         inject_prompt_marker(&mut env, &session_id);     // PS1/PROMPT/fish_prompt hook
116:     }
117: 
118:     let shell = detect_shell();                          // $SHELL or "/bin/sh"
119:     let mut builder = Command::new(&shell);
120:     builder
121:         .envs(&env)
122:         .stdin(Stdio::inherit())
123:         .stdout(Stdio::inherit())
124:         .stderr(Stdio::inherit());
125: 
126:     if let Some(c) = cmd {
127:         builder.args(["-c", &c]);
128:     }
129: 
130:     let status = builder.status()?;
131:     std::process::exit(status.code().unwrap_or(1));
132: }
133: 
134: fn generate_session_id() -> String {
135:     use rand::Rng;
136:     let suffix: String = rand::thread_rng()
137:         .sample_iter(&rand::distributions::Alphanumeric)
138:         .take(4)
139:         .map(char::from)
140:         .collect::<String>()
141:         .to_lowercase();
142:     format!("session-{suffix}")
143: }
144: ```
145: 
146: `Stdio::inherit()` keeps the subshell attached to the user's tty. On Unix we may swap to [`nix::unistd::execve`](https://docs.rs/nix/latest/nix/unistd/fn.execve.html) for exec-instead-of-spawn semantics (the subshell replaces the gitswitch process — cleaner `ps` output, no zombie parent). Spawn-and-wait is the default to keep the Windows code path identical.
147: 
148: Crates: [`rand`](https://crates.io/crates/rand), [`anyhow`](https://crates.io/crates/anyhow), [`nix`](https://crates.io/crates/nix) (Unix-only).
149: 
150: ### 3. Inside a session, commands are scope-aware
151: 
152: | Command | Outside session | Inside session |
153: |---|---|---|
154: | `gitswitch switch X` | Writes `~/.gitconfig`, syncs `gh hosts.yml`, writes `~/.ssh/config` blocks | `eval`s envs into the shell. Global config untouched. SSH config blocks still managed (per-host, not per-scope). |
155: | `gitswitch current` | Reads `~/.gitconfig` + active flag in YAML | Reads `GITSWITCH_SESSION_*` envs |
156: | `gitswitch status` | Identity section reads `git config --global` | Reads envs, labels source as "session" |
157: | `gitswitch pin X` | Global pin store keyed by repo path | Same — pins are per-repo, shared |
158: | `gitswitch list` | Reads profile DB | Reads profile DB (marks session-active profile) |
159: | `gitswitch add` / `remove` / `host add` / `host remove` / `import` | Profile DB | Profile DB (shared) |
160: | `gitswitch env <profile>` | Prints exports | Same |
161: | `gitswitch env --clear` | No-op | Prints unsets; ends session-scoped identity (session id stays) |
162: 
163: Override flags (clap v4 derive):
164: 
165: ```rust
166: #[derive(clap::Args)]
167: pub struct ScopeFlags {
168:     /// Force global behaviour even inside a session
169:     #[arg(long, conflicts_with = "session")]
170:     pub global: bool,
171: 
172:     /// Force session behaviour; errors if not inside a session
173:     #[arg(long)]
174:     pub session: bool,
175: }
176: ```
177: 
178: `gitswitch-env::scope` helpers (no allocation in the common path):
179: 
180: ```rust
181: pub fn in_session() -> bool {
182:     std::env::var_os("GITSWITCH_SESSION").is_some()
183: }
184: 
185: pub fn session_id() -> Option<String> {
186:     std::env::var("GITSWITCH_SESSION").ok()
187: }
188: 
189: pub fn active_profile_in_session() -> Option<String> {
190:     std::env::var("GITSWITCH_SESSION_PROFILE").ok()
191: }
192: ```
193: 
194: ### 4. Env vars that carry session identity
195: 
196: When a profile is active in a session (whether at entry or after an in-session `switch`), the wrapper exports:
197: 
198: ```sh
199: # Identity (git respects these per-process, overriding ~/.gitconfig)
200: export GIT_AUTHOR_NAME="Abhiram Kanna"
201: export GIT_AUTHOR_EMAIL="abhiram@corp.com"
202: export GIT_COMMITTER_NAME="Abhiram Kanna"
203: export GIT_COMMITTER_EMAIL="abhiram@corp.com"
204: 
205: # Signing — git's only env-based config injection mechanism
206: export GIT_CONFIG_COUNT=2
207: export GIT_CONFIG_KEY_0=user.signingkey
208: export GIT_CONFIG_VALUE_0=4A1B2C3D4E5F6789
209: export GIT_CONFIG_KEY_1=commit.gpgsign
210: export GIT_CONFIG_VALUE_1=true
211: 
212: # gh CLI per-host tokens (multi-host plan delivers these from keychain)
213: export GH_TOKEN=...
214: export GH_ENTERPRISE_TOKEN=...
215: 
216: # SSH — only set for single-host profiles; multi-host uses ~/.ssh/config blocks
217: export GIT_SSH_COMMAND="ssh -i ~/.ssh/id_corp -o IdentitiesOnly=yes"
218: 
219: # gitswitch session state
220: export GITSWITCH_SESSION="session-a1b2"
221: export GITSWITCH_SESSION_PROFILE="corp"
222: export GITSWITCH_SESSION_PID=12345
223: ```
224: 
225: Builder (`gitswitch-env::builder`):
226: 
227: ```rust
228: use anyhow::Result;
229: use serde::Serialize;
230: 
231: #[derive(Debug, Clone, Serialize)]
232: pub struct EnvVar {
233:     pub key: String,
234:     pub value: String,
235: }
236: 
237: pub fn build(profile: &Profile) -> Result<Vec<EnvVar>> {
238:     let mut out = vec![
239:         EnvVar { key: "GIT_AUTHOR_NAME".into(),     value: profile.user_name.clone() },
240:         EnvVar { key: "GIT_AUTHOR_EMAIL".into(),    value: profile.email.clone() },
241:         EnvVar { key: "GIT_COMMITTER_NAME".into(),  value: profile.user_name.clone() },
242:         EnvVar { key: "GIT_COMMITTER_EMAIL".into(), value: profile.email.clone() },
243:     ];
244: 
245:     // Signing + arbitrary git config via GIT_CONFIG_COUNT / KEY_n / VALUE_n
246:     let mut pairs: Vec<(&str, String)> = Vec::new();
247:     if let Some(key) = &profile.signing_key {
248:         pairs.push(("user.signingkey", key.clone()));
249:         pairs.push(("commit.gpgsign", "true".into()));
250:     }
251:     if !pairs.is_empty() {
252:         out.push(EnvVar { key: "GIT_CONFIG_COUNT".into(), value: pairs.len().to_string() });
253:         for (i, (k, v)) in pairs.into_iter().enumerate() {
254:             out.push(EnvVar { key: format!("GIT_CONFIG_KEY_{i}"),   value: k.into() });
255:             out.push(EnvVar { key: format!("GIT_CONFIG_VALUE_{i}"), value: v });
256:         }
257:     }
258: 
259:     // Tokens from keyring per host in profile.hosts
260:     for host in &profile.hosts {
261:         if let Some(tok) = gitswitch_secrets::get_token(&host.name, &profile.id)? {
262:             let var = match host.kind {
263:                 HostKind::GitHubCom => "GH_TOKEN",
264:                 HostKind::GitHubEnterprise => "GH_ENTERPRISE_TOKEN",
265:                 _ => continue,
266:             };
267:             out.push(EnvVar { key: var.into(), value: tok });
268:         }
269:     }
270: 
271:     // SSH — only if single-host
272:     if profile.hosts.len() == 1 {
273:         if let Some(key_path) = &profile.hosts[0].ssh_key {
274:             out.push(EnvVar {
275:                 key: "GIT_SSH_COMMAND".into(),
276:                 value: format!("ssh -i {} -o IdentitiesOnly=yes", key_path),
277:             });
278:         }
279:     }
280: 
281:     Ok(out)
282: }
283: ```
284: 
285: Formatter (`gitswitch-env::format`):
286: 
287: ```rust
288: use serde::Serialize;
289: 
290: pub enum EnvFormat { Posix, Fish, Json }
291: 
292: #[derive(Serialize)]
293: struct JsonOut<'a> {
294:     session: Option<String>,
295:     profile: Option<&'a str>,
296:     vars: std::collections::BTreeMap<&'a str, &'a str>,
297: }
298: 
299: pub fn emit(vars: &[EnvVar], fmt: EnvFormat) -> String {
300:     match fmt {
301:         EnvFormat::Posix => vars.iter()
302:             .map(|v| format!("export {}={}", v.key, shell_escape(&v.value)))
303:             .collect::<Vec<_>>()
304:             .join("\n"),
305:         EnvFormat::Fish => vars.iter()
306:             .map(|v| format!("set -gx {} {}", v.key, shell_escape(&v.value)))
307:             .collect::<Vec<_>>()
308:             .join("\n"),
309:         EnvFormat::Json => {
310:             let map = vars.iter().map(|v| (v.key.as_str(), v.value.as_str())).collect();
311:             serde_json::to_string_pretty(&JsonOut {
312:                 session: scope::session_id(),
313:                 profile: scope::active_profile_in_session().as_deref(),
314:                 vars: map,
315:             }).unwrap()
316:         }
317:     }
318: }
319: 
320: pub fn emit_clear(keys: &[&str], fmt: EnvFormat) -> String {
321:     match fmt {
322:         EnvFormat::Posix => keys.iter().map(|k| format!("unset {k}")).collect::<Vec<_>>().join("\n"),
323:         EnvFormat::Fish  => keys.iter().map(|k| format!("set -e {k}")).collect::<Vec<_>>().join("\n"),
324:         EnvFormat::Json  => serde_json::to_string(&serde_json::json!({ "unset": keys })).unwrap(),
325:     }
326: }
327: ```
328: 
329: When `gitswitch switch other` runs inside a session, the wrapper emits a fresh block that overwrites all `GIT_*`, `GH_*`, and `GITSWITCH_SESSION_PROFILE` but **preserves** `GITSWITCH_SESSION` (the scope id stays; only the identity changes).
330: 
331: Crates: [`serde`](https://crates.io/crates/serde), [`serde_json`](https://crates.io/crates/serde_json), [`shell-escape`](https://crates.io/crates/shell-escape).
332: 
333: ### 5. `gitswitch env <profile>` for non-spawn use
334: 
335: For agents / CI that can't spawn a subshell, `gitswitch env <profile>` prints the same block to stdout for `eval`ing in-place. `--format json` for programmatic consumption:
336: 
337: ```json
338: {
339:   "session": "session-a1b2",
340:   ""corp",
341:   "vars": {
342:     "GIT_AUTHOR_NAME": "Abhiram Kanna",
343:     "GIT_AUTHOR_EMAIL": "abhiram@corp.com",
344:     "...": "..."
345:   }
346: }
347: ```
348: 
349: `gitswitch env --clear` emits unsets — drops session-scoped identity without exiting the shell.
350: 
351: Python agent integration (same example as Go plan, unchanged because the binary surface is identical):
352: 
353: ```python
354: import subprocess, json, os
355: out = json.loads(subprocess.check_output(
356:     ["gitswitch", "env", "corp", "--format", "json"], text=True
357: ))
358: env = os.environ | out["vars"]
359: subprocess.run(["git", "commit", "-am", "..."], env=env)
360: ```
361: 
362: ---
363: 
364: ## AI-agent ergonomics
365: 
366: Two parallel coding sessions, both editing different repos, possibly switching identities mid-session.
367: 
368: **Pane A**:
369: ```
370: $ gitswitch shell personal
371: [gs:personal] $ cd ~/code/my-oss
372: [gs:personal] $ # work, commit as personal
373: [gs:personal] $ gitswitch switch contractor
374: [gs:contractor] $ # work on a different repo, commit as contractor
375: ```
376: 
377: **Pane B** (concurrent):
378: ```
379: $ gitswitch shell corp
380: [gs:corp] $ cd ~/code/client-proj
381: [gs:corp] $ # work, commit as corp
382: ```
383: 
384: Neither pane affects the other. Neither touches global.
385: 
386: ### One-shot for headless flows
387: 
388: ```
389: $ gitswitch shell corp --cmd "git commit -am 'fix: foo'"
390: ```
391: Spawns a session, runs the command, exits.
392: 
393: ### JSON env for AI tool integration
394: 
395: See the Python example above. Same surface as the Go plan because both binaries emit the same JSON schema.
396: 
397: ---
398: 
399: ## Edge cases
400: 
401: ### Switching profile inside session vs. opening a new session
402: Both work. `gitswitch shell` then multiple `switch` calls is the normal flow when an AI agent moves between tasks. `gitswitch shell <profile>` is the convenience case.
403: 
404: ### Global switch from outside a session
405: `gitswitch switch other` outside any session writes `~/.gitconfig` as today. Doesn't affect open sessions.
406: 
407: ### Global switch initiated from inside a session
408: `gitswitch switch other --global` from inside a session writes `~/.gitconfig` while leaving session envs alone.
409: 
410: ### Nested `gitswitch shell`
411: Inner shell gets a new `GITSWITCH_SESSION` id, shadowing the outer. On exit, the inner subshell process ends and the outer shell's process env is restored automatically. Prompt nesting: innermost only.
412: 
413: ### Token rotation mid-session
414: Keychain updated externally; session has stale `GH_TOKEN`. `gitswitch status` warns. Resolution: `gitswitch switch <current-profile>` re-reads keychain and re-emits envs.
415: 
416: ### Commit / push without realizing scope
417: `[gs:profile]` prompt + `gitswitch current` + `gitswitch status` make scope visible. Document "if unsure, run `gitswitch current`".
418: 
419: ### Session exit cleanup
420: The subshell process ends → envs disappear. No cleanup needed. No session registry persisted.
421: 
422: ### Repo-pin interaction
423: Pins are per-repo path, stored globally. A pin doesn't auto-activate in a session (session is explicit about its identity). `gitswitch status` shows both the session-active profile and the repo's pin, flags any mismatch.
424: 
425: ---
426: 
427: ## Interaction with the multi-host plan
428: 
429: Strictly orthogonal:
430: 
431: - **Multi-host** ([roadmap-v0.2.0.md](roadmap-v0.2.0.md)) extends the *data model*: profile = identity + host map.
432: - **Session isolation** extends the *scope*: per-process active identity instead of per-machine.
433: 
434: Both compose: any profile (single- or multi-host) can be the active identity inside a session. SSH per-host config (multi-host) handles SSH key selection by remote URL regardless of scope. Token routing (multi-host's keychain + gh `hosts.yml` sync) feeds the env emitter.
435: 
436: **Shared crates**:
437: - `gitswitch-secrets` — keychain read (also used by the multi-host plan).
438: - **New** `gitswitch-env` — `build(profile) -> Vec<EnvVar>` + format-agnostic `emit` (POSIX / fish / JSON). Used by both `shell` (spawn) and `env` (eval).
439: 
440: ---
441: 
442: ## Verification plan
443: 
444: ### Unit
445: - `gitswitch_env::builder::build` — single-host, multi-host, signing variants.
446: - `gitswitch_env::format::emit` — POSIX / fish / JSON output (snapshot tests via [`insta`](https://crates.io/crates/insta)).
447: - `gitswitch_env::scope::in_session` / `session_id` / `active_profile_in_session` — env-var presence cases.
448: 
449: ### Integration
450: - `gitswitch shell --cmd "git var GIT_AUTHOR_IDENT"` returns expected identity for the pre-switched profile.
451: - `gitswitch shell`, then in-session `gitswitch switch corp`; `git var GIT_AUTHOR_IDENT` returns corp — exercised via [`assert_cmd`](https://crates.io/crates/assert_cmd).
452: - Eval-into-existing-shell: spawn `bash -c 'eval "$(gitswitch env corp)"; env'`, assert exports present.
453: - Two concurrent `gitswitch shell --cmd "..."` invocations with different profiles run in parallel — assert outputs don't cross-contaminate (no shared mutable state between processes).
454: - `gitswitch switch corp --global` from inside a session updates `~/.gitconfig` without changing session envs.
455: 
456: ### Shell-wrapper portability
457: - Reinstall on a fresh `~/.zshrc` / `~/.bashrc` / fish config; assert wrapper function emitted; `gitswitch switch` inside a session updates envs without spawning a new shell.
458: 
459: ### Edge cases (manual)
460: - Nested `gitswitch shell` — exit unwinds correctly, parent prompt restored.
461: - `gitswitch env --clear` inside a session — identity envs gone, `GITSWITCH_SESSION` stays.
462: - `gitswitch shell --cmd "..."` — runs in fresh session, exits cleanly with child's status code.
463: 
464: ---
465: 
466: ## Command surface (additions)
467: 
468: | Command | Purpose |
469: |---|---|
470: | `gitswitch shell [profile] [--cmd "..."] [--no-prompt] [--shell <path>]` | Spawn isolated session subshell, optionally pre-switched |
471: | `gitswitch env <profile> [--format posix\|fish\|json]` | Print exports for eval / programmatic use |
472: | `gitswitch env --clear [--format ...]` | Print unsets (drops session identity without exiting) |
473: | `gitswitch switch <profile> [--global\|--session] --emit-env` | (Internal flag the wrapper uses) emit envs instead of mutating files |
474: 
475: ### Extended (scope-aware)
476: - `gitswitch switch` — session by default inside a session, global outside; explicit flags override.
477: - `gitswitch current` — labels source (session vs. global).
478: - `gitswitch status` — identity section reads envs when in a session.
479: 
480: ### Unchanged
481: - `gitswitch add` / `remove` / `host add` / `host remove` / `import gh` / `import git` / `list` / `pin` / `unpin` / `record` / `recommend` / `install` / `uninstall` / `version` / `upgrade` / `hook-check` / `claude` / `pacman` / `doctor`.
482: 
483: ## Critical files / crates
484: 
485: - **New crate** `gitswitch-env`:
486:   - `builder.rs` — `build(profile) -> Vec<EnvVar>`.
487:   - `format.rs` — `EnvFormat { Posix, Fish, Json }` + `emit` / `emit_clear`.
488:   - `scope.rs` — `in_session`, `session_id`, `active_profile_in_session`.
489: - `gitswitch-shell` (existing per [migration plan](roadmap-migration.md)) — `install::wrapper_snippet(shell)` emits the shell function wrapper alongside the hook.
490: - `gitswitch` (bin crate):
491:   - `cli/shell.rs` — `cmd_shell` (spawn subshell with env block).
492:   - `cli/env.rs` — `cmd_env` (emit eval-able block, `--clear`, `--format`).
493:   - `cli/switch.rs` — `--emit-env` flag; `--global` / `--session` scope flags.
494:   - `cli/current.rs`, `cli/status.rs` — scope-aware reads.
495: - `gitswitch-secrets` — shared with the [multi-host plan](roadmap-v0.2.0.md).
496: 
497: External crates: [`clap`](https://crates.io/crates/clap) v4 derive, [`rand`](https://crates.io/crates/rand), [`serde`](https://crates.io/crates/serde), [`serde_json`](https://crates.io/crates/serde_json), [`anyhow`](https://crates.io/crates/anyhow), [`shell-escape`](https://crates.io/crates/shell-escape), [`nix`](https://crates.io/crates/nix) (Unix), [`assert_cmd`](https://crates.io/crates/assert_cmd) + [`insta`](https://crates.io/crates/insta) (tests).
498: 
499: ## Out of scope (potential v0.3.0)
500: 
501: - **Per-session shadow `HOME` / `XDG_CONFIG_HOME`** — full FS-level isolation. Envs cover practical cases; FS isolation adds copy-lifecycle complexity.
502: - **Daemon-mediated sessions** — central inventory + token rotation push. Not worth the always-on process.
503: - **`direnv` integration** — auto-open a session on `cd` into a configured repo. Hooks into in-repo `.gitswitch.yaml` (v0.3.0).
504: 
505: ## Companion plans
506: 
507: - [Multi-host](roadmap-v0.2.0.md)
508: - [TUI](roadmap-v0.2.0-tui.md)
509: - [Migration](roadmap-migration.md)
510: - [v0.3.0](roadmap-v0.3.0.md)

(End of file - total 510 lines)
</content>
```

**Tool: read**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-v0.2.0-tui.md"
}
```

**Output:**
```
<path>/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-v0.2.0-tui.md</path>
<type>file</type>
<content>
1: # gitswitch v0.2.0 — TUI Refresh (Rust)
2: 
3: > Rust translation of the [Go TUI plan](../roadmap-v0.2.0-tui.md). Information architecture, keybinds, command palette, wizards, three themes — identical. This document locks the Rust crate stack and the Elm-architecture state shape, and shows how the arcade aesthetic is preserved through [tachyonfx](https://github.com/junkdog/tachyonfx) effects layered atop [ratatui](https://ratatui.rs/) widgets.
4: 
5: ## Context
6: 
7: The v0.1.x TUI is a single-screen list selector. v0.2.0's surface — multi-host profiles, per-remote pins, session scope, gh-driven import, doctor diagnostics — needs five screens, a command palette, and wizards. The arcade aesthetic (intro animation, theme colors, `pacman` mode) is what users like; usability around it has to grow up.
8: 
9: Goal: **every CLI command reachable through the TUI**, navigable by users who don't know the flag syntax, while keeping the arcade feel as a visual layer (boot splash, palette, transitions) rather than a gating mechanic.
10: 
11: The Rust ecosystem actually gives us a *richer* arcade than the Go plan can deliver: [`tachyonfx`](https://github.com/junkdog/tachyonfx) is a dedicated frame-by-frame effects library for ratatui (50+ effects: character evolution, materialization, directional sliding, color transforms — stateful, composable, applied after widget render). The Go equivalent was a hand-rolled `arcade/` subpackage plus harmonica springs. Here we get glyph-cycling, dissolve/coalesce, and CRT-style overlays out of the box.
12: 
13: Built parallel to the [multi-host](roadmap-v0.2.0.md) and [session isolation](roadmap-v0.2.0-isolation.md) plans; consumes both their data models. Shares no runtime state with them — pure presentation layer over `gitswitch-core` (storage), `gitswitch-secrets`, `gitswitch-env`, `gitswitch-gh`, `gitswitch-git`.
14: 
15: ## Stack decision
16: 
17: **[ratatui](https://crates.io/crates/ratatui) 0.30+ + [ratatui-crossterm](https://crates.io/crates/ratatui-crossterm) backend + [tachyonfx](https://crates.io/crates/tachyonfx) + [throbber-widgets-tui](https://crates.io/crates/throbber-widgets-tui).**
18: 
19: Why ratatui:
20: - Dominates Rust TUI since the fork from `tui-rs` in 2023; `0.30` modularization split the backends into `ratatui-crossterm` / `ratatui-termion` / `ratatui-termwiz` so we depend only on what we use.
21: - Immediate-mode draw model + atomic frame updates via crossterm's DEC synchronized output → smooth animations without flicker. Important for the arcade flavor.
22: - Huge third-party widget showcase (the [widget gallery](https://ratatui.rs/showcase/) covers tables, forms, charts, gauges, inputs — and crucially the animation/effects layer below).
23: 
24: Why not `cursive`:
25: - Callback-based, retained widget tree. Harder to compose with a frame-by-frame effects library like tachyonfx. Worse fit for the lazygit/k9s multi-panel, contextual-keybind feel that's the modern TUI baseline.
26: 
27: Why not [`iocraft`](https://github.com/ccbrown/iocraft):
28: - React-like component model with a JSX-flavored macro DSL — genuinely nice ergonomics, less boilerplate than ratatui's immediate-mode draw. Worth considering on its own merits.
29: - Rejected because the arcade aesthetic is a load-bearing requirement: tachyonfx + tui-big-text + throbber-widgets-tui + nucleo-matcher + tui-textarea are all ratatui-native and have no iocraft equivalents. Re-implementing the effect layer from scratch would be a project on its own.
30: - Also: ecosystem maturity gap. ratatui has a years-deep widget showcase, two reference TUI apps in production use (Helix uses ratatui internals; lazygit-style Rust tools standardize on it). iocraft is genuinely newer and smaller.
31: - Reconsider in v0.3.0 if iocraft's effects story catches up *and* a clean ratatui→iocraft migration path emerges.
32: 
33: Architecture pattern: **Elm-style** (`Model` / `Msg` / `update` / `view`) on top of ratatui's immediate-mode `Frame::render_widget`. ratatui doesn't impose an architecture; we impose one for testability — `update` is a pure function over state + message, `view` is pure render. Side effects routed through `Cmd` enums returned by `update`, executed by the runtime loop.
34: 
35: Supplementary crates:
36: - [`crossterm`](https://crates.io/crates/crossterm) — raw mode, key/mouse events, alternate screen, DEC synchronized output. Bundled via the `ratatui-crossterm` backend.
37: - [`tachyonfx`](https://crates.io/crates/tachyonfx) — effects + animation. Used for the arcade boot splash, focus transitions, palette overlay reveal, scanline overlay. See [examples](https://github.com/junkdog/tachyonfx/tree/main/examples).
38: - [`throbber-widgets-tui`](https://crates.io/crates/throbber-widgets-tui) — spinners. Custom symbol sets for the arcade theme (`▘▝▗▖` etc., 8-bit feel).
39: - [`tui-textarea`](https://crates.io/crates/tui-textarea) — multi-line text input for wizard pages with long fields.
40: - [`tui-input`](https://crates.io/crates/tui-input) — single-line text input.
41: - [`nucleo-matcher`](https://crates.io/crates/nucleo-matcher) — fuzzy matching for `:` command palette. Used by [Helix](https://github.com/helix-editor/helix); faster than `fuzzy-matcher`'s SkimMatcherV2 on large lists and provides deterministic ranking.
42: - [`tui-big-text`](https://crates.io/crates/tui-big-text) — block-letter rendering for the splash LOGO.
43: - [`ansi-to-tui`](https://crates.io/crates/ansi-to-tui) — render ANSI-colored subprocess output (e.g. doctor probe stdout) inside a `Paragraph`.
44: - [`color-eyre`](https://crates.io/crates/color-eyre) — panic handlers that restore the terminal before printing (critical: a panic in raw mode + alt screen will otherwise leave the user's terminal unusable).
45: - [`clap`](https://crates.io/crates/clap) — CLI dispatch into the TUI; keeps the `gitswitch` / `gitswitch pacman` / `gitswitch --no-tui` entry points.
46: - [`insta`](https://crates.io/crates/insta) — frame snapshot tests at fixed widths.
47: 
48: ## Information architecture
49: 
50: Five top-level screens reachable from a persistent left rail, plus modal flows (forms, confirmations, wizards). Inspired by lazygit's multi-panel layout, adapted to gitswitch's narrower scope.
51: 
52: ```
53: ┌─ gitswitch ────────────────────────────────────────────────────────────┐
54: │ ▌Profiles│   personal *active                                          │
55: │  Hosts   │   ┌─ profile detail ──────────────────────────────────────┐ │
56: │  Repos   │   │ nickname: personal                                    │ │
57: │  Session │   │ name:     Abhiram Kanna                               │ │
58: │  Doctor  │   │ email:    abhiram@personal.dev                        │ │
59: │          │   │ signing:  4A1B... (gpg, present ✓)                    │ │
60: │  ─────   │   │ source:   manual                                      │ │
61: │  Themes  │   │ hosts:                                                │ │
62: │  Help    │   │   github.com   ssh ~/.ssh/id_personal  gh aksisonline │ │
63: │          │   │                token ✓ (synced to gh hosts.yml)       │ │
64: │          │   │                                                       │ │
65: │          │   │ Actions: [s]witch  [e]dit  [d]uplicate  [x]delete     │ │
66: │          │   └───────────────────────────────────────────────────────┘ │
67: ├─────────────────────────────────────────────────────────────────────────┤
68: │ session: — (global)   active: personal   theme: arcade   ?: help        │
69: └─────────────────────────────────────────────────────────────────────────┘
70: ```
71: 
72: ### Screen 1: Profiles (default)
73: Left list of profiles; right detail pane shows the focused profile. Actions on the focused row: switch, edit, duplicate, delete. Multi-host profiles show one row per host expanded in the detail pane.
74: 
75: **CLI commands covered**: `switch`, `add`, `remove`, `current`, `import gh`, `import git`, `list`.
76: 
77: ### Screen 2: Hosts
78: View grouped by hostname (e.g. `github.com`, `github.corp.com`). For each host: which profiles cover it, current `gh hosts.yml` state, token validity (from `doctor` probes), SSH config block status. Per-row actions add/edit/remove host on a profile.
79: 
80: **CLI commands covered**: `host add`, `host remove`, `host list`.
81: 
82: ### Screen 3: Repos
83: Lists every repo gitswitch has seen via the shell hook (from `history.json`). Per-row: pin state (default + per-remote), last identity used, last activity timestamp. Detail pane = the `gitswitch status` output for that repo.
84: 
85: **CLI commands covered**: `pin`, `unpin`, `status`, `recommend`.
86: 
87: ### Screen 4: Session
88: - If launched outside a session: shows global state, with a `[n] New session` action that runs `gitswitch shell` (which exits the TUI cleanly, then spawns the session subshell).
89: - If launched inside a session: shows session id, session-active profile, env block. Actions: switch profile (session-scoped), clear identity, end session (exit).
90: 
91: **CLI commands covered**: `shell`, `env`, `current` (session view).
92: 
93: ### Screen 5: Doctor
94: Runs probes on demand and streams results in a scrollable `Paragraph` — keychain reachable, every profile's token valid against `/user`, SSH key files exist + permissions correct, `gh hosts.yml` in sync, shell hook version current. Per-finding action: `[f]ix` triggers the corresponding `--fix` flow (token re-auth, SSH config rewrite, hook reinstall).
95: 
96: **CLI commands covered**: `doctor`, `hook-check`, `install`, `uninstall`.
97: 
98: ### Persistent: Themes
99: Slide-out picker (same palette options the v0.1 TUI exposes), now including a named "arcade" theme that bundles palette + animated splash + retro borders + CRT scanline.
100: 
101: ### Persistent: Help
102: Auto-generated from the keymap. Custom `HelpView` widget renders the active screen's bindings + the always-on global bindings.
103: 
104: ---
105: 
106: ## Keybindings
107: 
108: Vim + arrow consistency throughout (matches lazygit / k9s muscle memory):
109: 
110: ### Global
111: | Key | Action |
112: |---|---|
113: | `1`–`5` | Jump to screen 1-5 |
114: | `?` | Toggle full-screen help overlay |
115: | `q` / `Ctrl-C` | Quit (with confirm if mid-edit) |
116: | `Tab` / `Shift-Tab` | Cycle focused panel within current screen |
117: | `:` | Command palette (fuzzy-find any CLI action) |
118: | `/` | Filter current list |
119: | `Esc` | Cancel modal / close overlay |
120: 
121: ### Profiles screen
122: | Key | Action |
123: |---|---|
124: | `j` / `k` / `↑` / `↓` | Move selection |
125: | `Enter` / `s` | Switch to selected profile (scope-aware) |
126: | `S` | Force global switch (when in a session) |
127: | `e` | Edit profile (modal form) |
128: | `a` | Add new profile (wizard) |
129: | `i` | Import flow (sub-menu) |
130: | `d` | Duplicate profile |
131: | `x` | Delete (with confirm) |
132: | `h` | Manage hosts on this profile (jumps to screen 2 filtered) |
133: 
134: ### Hosts screen
135: | Key | Action |
136: |---|---|
137: | `j` / `k` | Move |
138: | `a` | Add host to selected profile |
139: | `e` | Edit host creds (SSH key path, gh_user, token) |
140: | `x` | Remove host from profile |
141: | `t` | Re-fetch token from keychain → sync gh hosts.yml |
142: | `g` | Inspect gh hosts.yml diff |
143: 
144: ### Repos screen
145: | Key | Action |
146: |---|---|
147: | `j` / `k` | Move |
148: | `Enter` | Open status detail |
149: | `p` | Pin (default) — picks profile via sub-list |
150: | `P` | Pin a remote — picks remote, then profile |
151: | `u` | Unpin (with `?` sub-prompt for which) |
152: | `o` | `cd` into repo path on exit (writes hint file the shell hook reads) |
153: 
154: ### Session screen
155: | Key | Action |
156: |---|---|
157: | `n` | Spawn new session (exits TUI, runs `gitswitch shell`) |
158: | `s` | Switch profile within current session (session-scoped) |
159: | `c` | Clear session identity (keep session, drop profile) |
160: | `Ctrl-D` | End session (exits TUI + subshell) |
161: 
162: ### Doctor screen
163: | Key | Action |
164: |---|---|
165: | `r` | Re-run all probes |
166: | `j` / `k` | Move through findings |
167: | `f` | Apply fix for selected finding |
168: | `F` | Apply all fixes (with confirm) |
169: 
170: All keybinds remappable via `~/.config/gitswitch/config.yaml` (new `keymap:` section). Loaded once at startup into a `KeyMap` struct that the `update` functions dispatch from.
171: 
172: ---
173: 
174: ## Command palette (`:`)
175: 
176: Fuzzy-search every CLI action. Powered by [`nucleo-matcher`](https://crates.io/crates/nucleo-matcher) over a static `Vec<PaletteAction>` enumerating every CLI verb. Each result is a verb + brief help:
177: 
178: ```
179: :add profile
180: :switch personal
181: :host add corp github.corp.com
182: :import gh
183: :doctor --fix
184: :shell corp
185: :pin --remote origin
186: ```
187: 
188: This is the answer to "dumb users can do everything power users can via flags" — they don't need to know the CLI syntax; they fuzzy-find the verb and the palette walks them through the rest via the standard wizard.
189: 
190: ```rust
191: pub struct PaletteAction {
192:     pub id: &'static str,         // "profiles.switch"
193:     pub label: &'static str,      // "switch profile"
194:     pub cli_form: &'static str,   // "gitswitch switch <nickname>"
195:     pub kind: PaletteKind,        // Direct | Wizard | ExternalSpawn
196: }
197: ```
198: 
199: Nucleo's `Matcher::fuzzy_match` returns sorted scores; we keep top 20.
200: 
201: ---
202: 
203: ## Forms & wizards
204: 
205: [`tui-input`](https://crates.io/crates/tui-input) for single-line, [`tui-textarea`](https://crates.io/crates/tui-textarea) for multi-line. A small `gitswitch-tui::wizard` helper drives linear step pages with per-step validation, back/next/cancel.
206: 
207: ```rust
208: pub struct Wizard<S> {
209:     pub steps: Vec<Box<dyn WizardStep<S>>>,
210:     pub idx: usize,
211:     pub state: S,
212: }
213: 
214: pub trait WizardStep<S> {
215:     fn title(&self) -> &str;
216:     fn update(&mut self, state: &mut S, key: KeyEvent) -> StepOutcome;
217:     fn view(&self, state: &S, f: &mut Frame, area: Rect);
218:     fn validate(&self, state: &S) -> Result<(), String>;
219: }
220: 
221: pub enum StepOutcome { Stay, Next, Back, Cancel, Submit }
222: ```
223: 
224: ### Add profile wizard
225: 1. Nickname (text, validated unique).
226: 2. Identity source: [Import from gh] [Import from git config] [Enter manually].
227: 3. (manual path) Name + email + signing key.
228: 4. Hosts: add one or more `(hostname, ssh_key, gh_user, token?)` rows. Token input can paste a PAT or `[d]` to delegate to `gh auth login` (TUI suspends — see Side effects below).
229: 5. Confirm summary → save.
230: 
231: ### Import gh wizard
232: 1. Probe `gh auth status` → list discovered accounts.
233: 2. For each, checkbox: import as profile, merge into existing, skip.
234: 3. (merge path) pick target profile.
235: 4. Confirm → execute.
236: 
237: ### Pin wizard (P key)
238: 1. Pick remote from current repo's `git remote -v`.
239: 2. Pick profile.
240: 3. Apply.
241: 
242: ### Token entry
243: PAT inputs masked, displayed as `••••••••` with `[r]eveal` toggle. Tokens **never written to disk via TUI** — go straight into the keychain via `gitswitch-secrets`.
244: 
245: ---
246: 
247: ## Arcade preservation
248: 
249: The biggest section. The Go plan stops at "ASCII logo + sparkle particles + CRT vignette"; the Rust stack lets us go further because tachyonfx is purpose-built for this.
250: 
251: ### Boot splash
252: First-run + `theme=arcade` default: ~800ms reveal of big-text "GITSWITCH" via [`tui-big-text`](https://crates.io/crates/tui-big-text), composed with tachyonfx:
253: 
254: ```rust
255: use tachyonfx::{fx, Effect, Duration, Interpolation};
256: 
257: fn arcade_splash() -> Effect {
258:     fx::sequence(&[
259:         // 1. Random glyph soup → solidify into "GITSWITCH"
260:         fx::coalesce(Duration::from_millis(450))
261:             .with_cell_filter(big_text_cells()),
262:         // 2. Star particles flicker over the now-rendered logo
263:         fx::repeating(
264:             fx::evolve_into("·✦*✧⋆", Duration::from_millis(80)),
265:             6,
266:         ),
267:         // 3. Dissolve out, then coalesce the main view in
268:         fx::parallel(&[
269:             fx::dissolve(Duration::from_millis(220)),
270:             fx::delay(
271:                 Duration::from_millis(120),
272:                 fx::coalesce(Duration::from_millis(280)),
273:             ),
274:         ]),
275:     ])
276: }
277: ```
278: 
279: Skippable on any key (the runtime checks `effects.is_running()`; on keypress, calls `effects.skip_to_end()`). Subsequent runs default to "short" splash (just the coalesce step, ~200ms); fully off if `prefs.splash = "off"`.
280: 
281: ### Theme: "arcade"
282: - **Palette**: loud — lifted from the existing Go themes (`magenta`, `cyan`, `yellow`) — defined as `ratatui::style::Color` constants in `gitswitch-tui::theme`.
283: - **Borders**: doubled (`BorderType::Double`) on the focused panel.
284: - **CRT scanlines**: custom `ScanlineOverlay` widget — renders alternating dim lines (every other terminal row) at low opacity on top of the focused panel. Implemented as a post-render pass that mutates `Buffer` cell styles.
285: - **Active marker blink**: tachyonfx `repeating` over a `hsl_shift` of the focused row's background.
286: - **Spinner**: `throbber-widgets-tui` with the custom symbol set:
287:   ```rust
288:   Throbber::default()
289:       .throbber_set(throbber_widgets_tui::Set {
290:           symbols: &["▘","▝","▞","▖","▌","▟","▛","▜"],
291:           empty: " ",
292:       })
293:       .label("probing")
294:   ```
295: - **Optional CRT vignette**: background dim layer via a tachyonfx `hsl_shift` with negative lightness on the outer 2 rows/cols.
296: 
297: ### Theme: "modern" (new default)
298: Clean ratatui block styles. Quiet single-line borders. No continuous animation. Focus transitions use a single tachyonfx `slide_in` of the new panel content over 100ms with `Interpolation::QuadOut` — fast, calm, no afterglow.
299: 
300: ### Theme: "minimal"
301: Monochrome, no borders, no effects. For tmux+ssh, basic xterm, and partial screen-reader compatibility. (TUIs remain inherently hard for screen readers; minimal mode reduces noise.) `theme=minimal` short-circuits the `EffectsManager` to a no-op stub so we don't pay the tick cost.
302: 
303: ### Easter eggs
304: - **`gitswitch pacman`** hidden subcommand bypasses the main TUI entirely. Renders an arcade intro then a dot-eater that crawls the profile list. Ghost sprites move via tachyonfx `slide_in`/`slide_out` along the panel's column axis; the muncher is a one-cell widget with `evolve_into("◐◓◑◒")` cycling.
305: - **`:konami`** typed into the palette unlocks a neon palette swap + matrix-style background cascade. Implemented as a full-screen `evolve_from` effect over a character cell grid:
306:   ```rust
307:   fx::repeating(
308:       fx::evolve_from(
309:           "ｱｲｳｴｵｶｷｸｹｺ01",
310:           Duration::from_millis(60),
311:       ),
312:       u32::MAX,
313:   )
314:   ```
315: 
316: ---
317: 
318: ## State management
319: 
320: Elm-style. Single owned `App` struct; `update` returns optional `Cmd`s that the runtime executes; `view` is pure.
321: 
322: ```rust
323: pub struct App {
324:     pub screen: Screen,
325:     pub profiles: ProfilesState,
326:     pub hosts: HostsState,
327:     pub repos: ReposState,
328:     pub session: SessionState,
329:     pub doctor: DoctorState,
330:     pub palette: PaletteState,         // : palette overlay
331:     pub form: Option<FormState>,       // active modal
332:     pub help: HelpState,
333:     pub theme: Theme,
334:     pub scope: Scope,                  // global vs in-session
335:     pub store: Arc<Store>,
336:     pub secrets: Arc<dyn SecretStore>,
337:     pub gh: Arc<gh::Client>,
338:     pub status_bar: StatusBar,
339:     pub effects: EffectsManager,       // tachyonfx active-effects pool
340:     pub keymap: KeyMap,
341:     pub should_exit: bool,
342:     pub exit_payload: Option<String>,  // for session env eval handoff
343: }
344: 
345: pub enum Screen { Profiles, Hosts, Repos, Session, Doctor }
346: 
347: pub enum Msg {
348:     Key(KeyEvent),
349:     Resize(u16, u16),
350:     Tick(Duration),                    // ~16ms cadence for effects
351:     DataLoaded(DataKind),
352:     SubprocessDone(SubprocessResult),
353: }
354: 
355: pub enum Cmd {
356:     LoadData(DataKind),
357:     RunSubprocess(SubprocessSpec),
358:     ApplyProfile { profile: String, scope: Scope },
359:     SuspendForExternal(ExternalSpec),  // e.g. gh auth login
360:     ExitWithPayload(String),
361:     SpawnShell(String),
362:     None,
363: }
364: 
365: impl App {
366:     pub fn update(&mut self, msg: Msg) -> Option<Cmd> { /* dispatch */ }
367:     pub fn view(&self, f: &mut Frame) { /* render */ }
368: }
369: ```
370: 
371: ### Event loop
372: 
373: Synchronous; `crossterm::event::poll(Duration::from_millis(16))` for ~60fps tick. tachyonfx effects are driven by per-frame `elapsed` deltas:
374: 
375: ```rust
376: let mut last = Instant::now();
377: loop {
378:     if event::poll(Duration::from_millis(16))? {
379:         let msg = match event::read()? {
380:             Event::Key(k) => Msg::Key(k),
381:             Event::Resize(w, h) => Msg::Resize(w, h),
382:             _ => continue,
383:         };
384:         if let Some(cmd) = app.update(msg) { runtime.dispatch(cmd, &mut app); }
385:     }
386:     let now = Instant::now();
387:     let dt = now - last;
388:     last = now;
389:     app.update(Msg::Tick(dt));
390: 
391:     terminal.draw(|f| {
392:         app.view(f);
393:         app.effects.process(dt, f.buffer_mut(), f.area());
394:     })?;
395: 
396:     if app.should_exit { break; }
397: }
398: ```
399: 
400: ### Side effects
401: 
402: - **Switch profile outside session**: `Cmd::ApplyProfile { scope: Global }` → runtime runs `git`/`gh` synchronously, blocks a frame, dispatches `DataLoaded` to refresh.
403: - **Switch profile inside session**: `Cmd::ExitWithPayload(eval_string)` → app sets `should_exit + exit_payload`; main prints payload to stdout after terminal cleanup, so the shell wrapper's `eval "$(gitswitch …)"` picks it up. Same payload mechanism as `gitswitch switch --emit-env` in the [session plan](roadmap-v0.2.0-isolation.md).
404: - **Spawn shell**: `Cmd::SpawnShell(profile)` → app exits cleanly, parent `main` exec's `gitswitch shell <profile>` so the subshell replaces the process. The TUI never tries to host a subshell inside its viewport.
405: - **gh auth login**: `Cmd::SuspendForExternal(...)` → runtime calls `disable_raw_mode()` + `LeaveAlternateScreen`, runs `gh auth login --hostname X` in foreground, then re-enters alt screen + `enable_raw_mode()` + forces a full redraw. Effects pool is paused during suspend.
406: - **Doctor fixes**: run in-process on a worker thread; stream `String` lines back through an `mpsc::Receiver<Msg::SubprocessDone(...)>` consumed at the top of each loop iteration. Output rendered via [`ansi-to-tui`](https://crates.io/crates/ansi-to-tui) into a scrollable `Paragraph`.
407: 
408: ### Panic safety
409: 
410: `color-eyre`'s panic hook is installed alongside a custom hook that calls `disable_raw_mode()` + `LeaveAlternateScreen` *before* the panic message is printed. Otherwise a panic mid-frame strands the user in a cursor-hidden, raw-mode terminal.
411: 
412: ```rust
413: fn install_panic_hook() {
414:     let original = std::panic::take_hook();
415:     std::panic::set_hook(Box::new(move |info| {
416:         let _ = disable_raw_mode();
417:         let _ = execute!(io::stdout(), LeaveAlternateScreen, Show);
418:         original(info);
419:     }));
420: }
421: ```
422: 
423: ---
424: 
425: ## Migration & opt-out
426: 
427: - v0.1's TUI was the default no-arg behavior of `gitswitch`. v0.2 keeps that: bare `gitswitch` opens the new TUI.
428: - `gitswitch --no-tui` (or `GITSWITCH_NO_TUI=1`) keeps a CLI-only mode that prints the equivalent of screen 1.
429: - All TUI screens have a `:cli` palette action that prints "the CLI command that would do this" — discoverability bridge for users who want to learn the flag syntax.
430: 
431: ---
432: 
433: ## Verification plan
434: 
435: ### Unit
436: - Each state struct: `update` handles every `Msg` variant relevant to it; `view` renders without panic for empty / populated / error states.
437: - Wizard: per-step validation, back/next/cancel transitions.
438: - Palette: nucleo fuzzy ranking matches expected verbs for representative queries (`add prof` → `profiles.add` ranked first).
439: - KeyMap: YAML override loader merges with defaults; conflicting bindings rejected with named error.
440: 
441: ### Integration
442: - [`insta`](https://crates.io/crates/insta) frame snapshot tests. For each screen + state combination, render a `Buffer` at fixed widths (80, 100, 132, 180 cols) and snapshot the cell grid (`buffer_to_string`-style helper that flattens to ANSI). Re-runs detect layout regressions.
443: - A `TestHarness` that drives `Msg::Key` events through `app.update` without a real terminal — verifies command palette walk-throughs, wizard completion, switch flows.
444: 
445: ### Manual
446: - **Dumb-user walkthrough**: with no flags knowledge, run TUI, add a profile, import gh accounts, pin a repo, open doctor, fix a finding. Doable from keybinds + help bar alone.
447: - **Arcade feel preserved**: `theme=arcade`, splash renders, transitions feel snappy, no flicker. tachyonfx effects don't stutter at 60fps on a 132-col terminal.
448: - **Multi-host coverage**: edit a profile to add a second host, verify the wizard, switch to it, check hosts screen reflects synced state.
449: - **Session integration**: launch TUI inside `gitswitch shell`; switching a profile mutates session envs (TUI exits with eval payload, wrapper picks it up).
450: - **Resize + small terminal**: behaves sanely at 60 cols (panels collapse to single-column mode via `Layout` constraint switch).
451: - **Theme=minimal**: usable in tmux-over-ssh, no broken glyphs in basic xterm.
452: 
453: ---
454: 
455: ## Command surface
456: 
457: No new top-level CLI commands. Existing entry points unchanged:
458: - `gitswitch` (no args) → TUI default screen.
459: - `gitswitch pacman` → arcade intro + dot-eater (preserved).
460: - `gitswitch --no-tui` / `GITSWITCH_NO_TUI=1` → CLI-only.
461: 
462: New config in `~/.config/gitswitch/config.yaml`:
463: ```yaml
464: prefs:
465:   theme: arcade           # arcade | modern | minimal
466:   splash: short           # full | short | off
467:   keymap:                 # optional overrides
468:     profiles.switch: [enter, s, space]
469:     global.palette: [":"]
470: ```
471: 
472: ---
473: 
474: ## Critical files / crates
475: 
476: - **New crate** `crates/gitswitch-tui/`:
477:   - `app.rs` — root `App`, `Msg`, `Cmd`, update/view dispatch.
478:   - `screens/{profiles,hosts,repos,session,doctor}.rs` — one module per screen.
479:   - `forms/` — wizard helper, input wrappers, masked-token input.
480:   - `palette.rs` — `:` command palette overlay + nucleo matcher.
481:   - `keymap.rs` — central keymap + YAML override loader.
482:   - `theme.rs` — palette definitions, `Theme` enum, theme-aware style helpers.
483:   - `runtime.rs` — event loop, panic hook, suspend/resume around external spawns.
484: - **New crate** `crates/gitswitch-arcade/` — splash, `ScanlineOverlay` widget, theme=arcade effects bundle. Depends on `tachyonfx`, `tui-big-text`, `throbber-widgets-tui`. Cleanly separated so `theme=minimal` builds can `--no-default-features` it out.
485: - **Existing bin** `crates/gitswitch/` — wires the TUI; preserves `pacman` and `--no-tui`.
486: - **Existing lib** `crates/gitswitch-core/` — `Prefs` extended with `theme` / `splash` / `keymap` fields (sourced from the v0.2.0 multi-host plan).
487: 
488: ---
489: 
490: ## Mouse support — first-class in v0.2.0
491: 
492: Promoted from the Go plan's "out of scope" list. Crossterm exposes mouse events natively and ratatui's `Rect` widget bounds make hit-testing straightforward; the cost is small enough to ship now, and the UX win for non-keybind-fluent users (the explicit goal of this plan) is large.
493: 
494: ### What's enabled
495: 
496: - **Click to focus + select**: clicking a row in the Profiles / Hosts / Repos list moves selection. Double-click triggers the row's default action (`switch` on Profiles, `status` on Repos, etc.).
497: - **Click the left rail** to jump screens — same as `1`–`5` keys.
498: - **Click status-bar segments** (e.g. theme indicator) to open the corresponding picker.
499: - **Click `[s]witch` / `[e]dit` / `[d]uplicate` / `[x]delete` action chips** in the detail pane to invoke them.
500: - **Scroll-wheel** scrolls the focused list or viewport (doctor output, status detail).
501: - **Right-click on a profile / repo row** opens a contextual action menu (mirrors the keymap for that screen so users discover keybinds visually).
502: - **Drag-select disabled** — not useful here and complicates terminal text-selection behavior.
503: 
504: ### Hit-testing model
505: 
506: ratatui's layout pass already produces `Rect`s for every widget. We tag interactive regions in a per-frame `HitRegions` registry populated during `view`:
507: 
508: ```rust
509: pub struct HitRegions {
510:     regions: Vec<(Rect, HitTarget)>,
511: }
512: 
513: pub enum HitTarget {
514:     NavRail(Screen),
515:     ProfileRow(usize),
516:     HostRow(usize),
517:     RepoRow(usize),
518:     ActionChip(ChipId),
519:     StatusBarSegment(StatusSegment),
520:     PaletteResult(usize),
521: }
522: 
523: impl HitRegions {
524:     pub fn at(&self, col: u16, row: u16) -> Option<HitTarget> {
525:         self.regions.iter().rev()
526:             .find(|(rect, _)| rect.contains((col, row).into()))
527:             .map(|(_, target)| target.clone())
528:     }
529: }
530: ```
531: 
532: The `App` calls `view` to render *and* populate `HitRegions`. On `Event::Mouse`, `update` looks up the target and dispatches the same `Msg` a keybind would. Single code path for the action; mouse is just another input modality.
533: 
534: ### Crossterm wiring
535: 
536: In the runtime loop:
537: 
538: ```rust
539: use crossterm::event::{EnableMouseCapture, DisableMouseCapture};
540: 
541: terminal::enable_raw_mode()?;
542: execute!(stdout, EnterAlternateScreen, EnableMouseCapture)?;
543: 
544: // main loop polls event::read() — Event::Mouse delivered as MouseEvent { kind, column, row, modifiers }
545: 
546: // cleanup (also from panic hook):
547: execute!(stdout, LeaveAlternateScreen, DisableMouseCapture)?;
548: terminal::disable_raw_mode()?;
549: ```
550: 
551: `MouseEventKind::Down(MouseButton::Left)` → primary action. `Down(Right)` → context menu. `ScrollDown` / `ScrollUp` → forward to focused viewport.
552: 
553: ### Theme-aware affordances
554: 
555: - `theme = arcade` adds a subtle hover-glow effect (tachyonfx slow pulse on the rect under the cursor when mouse motion reported). Only when motion events enabled (`EnableMouseCapture` covers motion in capture mode; we throttle to ≤10Hz to keep render cost down).
556: - `theme = modern` shows underlined chip text under the cursor.
557: - `theme = minimal` no hover affordances — keyboard-first, mouse-tolerated.
558: 
559: ### Constraints we accept
560: 
561: - **Terminal multiplexer interaction**: some tmux/screen configs intercept mouse for their own pane management. Document the trade-off: users with `set -g mouse on` in tmux can hold `Shift` (or terminal-specific modifier) to bypass mux and pass mouse through to gitswitch. Not a fix we can make, just a documented FAQ.
562: - **SSH over slow links**: motion events can flood. Disable motion (keep click-only) when `GITSWITCH_NO_MOUSE_MOTION=1` or when running on `theme = minimal`.
563: - **Mouse off opt-out**: `GITSWITCH_NO_MOUSE=1` or `prefs.mouse: false` disables the whole capture — useful when users want native terminal text selection to remain primary.
564: 
565: ### Testing
566: 
567: - Unit tests for `HitRegions::at` with synthetic `Rect` layouts.
568: - Integration: feed scripted `Event::Mouse` sequences to `update`; assert state transitions match keybind equivalents.
569: - Manual matrix: macOS Terminal.app, iTerm2, Alacritty, Kitty, WezTerm, Ghostty, GNOME Terminal, Windows Terminal — verify click + scroll work; document any quirks.
570: 
571: ## Out of scope (v0.3.0)
572: 
573: - Screen-reader optimized "linear" mode beyond `theme=minimal`.
574: - Built-in editor for `~/.gitconfig` — not gitswitch's job, but `:cli` shows the right `git config` invocation.
575: - Plugin-style TUI extensions (custom user-supplied screens).
576: 
577: ---
578: 
579: ## Sources
580: 
581: - [ratatui (crates.io)](https://crates.io/crates/ratatui) · [ratatui.rs guide](https://ratatui.rs/) · [`docs.rs/ratatui`](https://docs.rs/ratatui)
582: - [ratatui 0.30 modularization changelog](https://ratatui.rs/highlights/v030/)
583: - [crossterm (crates.io)](https://crates.io/crates/crossterm) · [`docs.rs/crossterm`](https://docs.rs/crossterm)
584: - [tachyonfx (GitHub)](https://github.com/junkdog/tachyonfx) · [`crates.io/crates/tachyonfx`](https://crates.io/crates/tachyonfx) — effects + animations for ratatui.
585: - [throbber-widgets-tui](https://crates.io/crates/throbber-widgets-tui) — spinners with customizable symbol sets.
586: - [tui-textarea](https://crates.io/crates/tui-textarea) · [tui-input](https://crates.io/crates/tui-input).
587: - [tui-big-text](https://crates.io/crates/tui-big-text) — block-letter rendering for the splash logo.
588: - [nucleo-matcher](https://crates.io/crates/nucleo-matcher) — fuzzy matcher used by [Helix](https://github.com/helix-editor/helix).
589: - [ansi-to-tui](https://crates.io/crates/ansi-to-tui) — render ANSI-colored output into ratatui spans.
590: - [color-eyre](https://crates.io/crates/color-eyre) — panic handlers that don't trash the terminal.
591: - [insta](https://crates.io/crates/insta) — snapshot testing for rendered frames.
592: - [lazygit design](https://www.bwplotka.dev/2025/lazygit/) — multi-panel + contextual keybinds + vim consistency (referenced from the Go plan; aesthetic carries over).
593: - [k9s patterns](https://k9scli.io/) — table-centric resource browsing, `:` command palette.

(End of file - total 593 lines)
</content>
```

**Tool: read**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-v0.3.0.md"
}
```

**Output:**
```
<path>/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-v0.3.0.md</path>
<type>file</type>
<content>
1: # gitswitch v0.3.0 — Roadmap (Rust)
2: 
3: Items deferred from [v0.2.0 planning](roadmap-v0.2.0.md). v0.2.0 establishes multi-host identity, keychain-backed PATs, gh-driven profile import, multi-pin, session isolation, ratatui TUI with arcade theme. v0.3.0 builds on that foundation to address team-scale distribution, in-repo config, and the long tail of authentication scenarios.
4: 
5: ---
6: 
7: ## 1. In-repo YAML config (`.gitswitch.yaml`)
8: 
9: A per-repo config file checked into the repo root. Declares which profile(s) should be used, which host(s) remotes resolve to, and identity policies the repo enforces.
10: 
11: ### Why
12: 
13: Today (and in 0.2.0), pins live in `~/.config/gitswitch/config.yaml` keyed by absolute repo path. That doesn't survive:
14: - A teammate cloning the repo to a different path.
15: - A CI runner.
16: - Switching machines.
17: 
18: In-repo config makes "this repo wants the `corp` profile" a property of the repo, not the local machine.
19: 
20: ### Open questions
21: 
22: - **Precedence**: in-repo config vs. local pin in `~/.config/gitswitch/config.yaml` — does in-repo win, or does local pin override (escape hatch)?
23: - **Policy enforcement**: should `.gitswitch.yaml` be able to *block* commits that don't match the policy (e.g. via a `prepare-commit-msg` hook)? Or only nudge?
24: - **Locked vs. overridable fields**: a team config probably wants to lock `email_suffix` but leave `ssh_key` path overridable per-dev.
25: - **Worktrees**: same `.gitswitch.yaml` applies to all worktrees of a repo? Or per-worktree opt-out? Use [`git2::Repository::worktrees()`](https://docs.rs/git2) to enumerate.
26: - **Schema versioning**: how does an older `gitswitch` binary handle a `.gitswitch.yaml` from a newer schema version?
27: 
28: ### Sketch
29: 
30: ```yaml
31: # .gitswitch.yaml (in repo root)
32: version: 1
33: require:
34:   email_suffix: "@corp.com"
35:   signing: required
36: hosts:
37:   origin: github.corp.com
38:   fork:   github.com           # optional, validated against actual remote
39: suggested_profile: corp
40: ```
41: 
42: Serialized via [`yaml_serde`](https://crates.io/crates/yaml_serde) (NOT `serde_yaml`) + [`serde`](https://crates.io/crates/serde). Versioned schema struct, separate from the user-config schema. Repo root discovery via [`git2`](https://crates.io/crates/git2) `Repository::discover()`.
43: 
44: ```rust
45: use serde::{Deserialize, Serialize};
46: use std::collections::HashMap;
47: use std::path::Path;
48: 
49: #[derive(Serialize, Deserialize)]
50: pub struct RepoConfig {
51:     pub version: u32,
52:     pub require: Option<RepoRequire>,
53:     pub hosts: HashMap<String, String>,  // remoteName -> hostname
54:     pub suggested_profile: Option<String>,
55: }
56: 
57: #[derive(Serialize, Deserialize)]
58: pub struct RepoRequire {
59:     pub email_suffix: Option<String>,
60:     pub signing: Option<SigningRequirement>,  // None | Required
61: }
62: 
63: #[derive(Serialize, Deserialize)]
64: #[serde(rename_all = "lowercase")]
65: pub enum SigningRequirement {
66:     Required,
67: }
68: 
69: // Loader:
70: pub fn load(repo_root: &Path) -> anyhow::Result<Option<RepoConfig>>;
71: ```
72: 
73: ---
74: 
75: ## 2. Team preset distribution
76: 
77: A team-shared bootstrap config — distributed via intranet URL, internal git repo, or shared drive — that onboards a new dev in one command:
78: 
79: ```
80: gitswitch apply https://wiki.corp.com/gitswitch-platform-team.yaml
81: gitswitch apply ./team-preset.yaml
82: ```
83: 
84: ### Components
85: 
86: - **Schema** for team presets (distinct from per-repo `.gitswitch.yaml` — this is a *template* that produces profiles, not policy applied to a repo).
87: - **Fetch + cache** for HTTPS sources. Use [`reqwest`](https://crates.io/crates/reqwest) blocking client for HTTPS, [`cacache`](https://crates.io/crates/cacache) for content-addressable cache. Local file path also supported.
88: - **Signing & verification** — so an attacker can't push a malicious preset to wiki / internal repo and steal credentials.
89:   - Candidate libraries: [`minisign-verify`](https://crates.io/crates/minisign-verify) (pure Rust, no libsodium, audited — **primary recommendation**), [`cosign-rs`](https://crates.io/crates/cosign-rs) (less mature), [`ed25519-dalek`](https://crates.io/crates/ed25519-dalek) for raw signature crypto.
90:   - Trust roots: a team config declares its signing key fingerprint; first-apply prompts the user to trust it (TOFU); subsequent applies verify against the trusted key. Trust storage in `~/.config/gitswitch/trust.yaml` keyed by signing key fingerprint.
91:   - Revocation: a `revoked.yaml` fetched from a known endpoint per trust root.
92: - **Apply semantics**: `apply` creates profile skeletons (host config, gh_user pattern, ssh_key naming hint) and prompts only for personal bits (the dev's actual SSH key path, their actual gh login). Pre-existing profiles aren't clobbered — only new ones are added.
93: 
94: ### Sketch
95: 
96: ```rust
97: use minisign_verify::{PublicKey, Signature};
98: 
99: pub fn verify_preset(
100:     preset_bytes: &[u8],
101:     sig: &Signature,
102:     pubkey: &PublicKey,
103: ) -> anyhow::Result<()> {
104:     pubkey.verify(preset_bytes, sig, /*allow_legacy=*/ false)
105:         .map_err(|e| anyhow::anyhow!("preset signature verification failed: {e}"))
106: }
107: ```
108: 
109: ### Why this is hard enough to defer
110: 
111: The signing story is the load-bearing piece. Shipping `apply` without signing means corp security teams will reject it on day one. Building signing properly (key management, trust onboarding, revocation, UX for "this preset is no longer trusted") is a project on its own.
112: 
113: ---
114: 
115: ## 3. Own OAuth flow (device flow)
116: 
117: Direct OAuth device flow against GitHub, bypassing `gh` CLI.
118: 
119: ### Why
120: 
121: v0.2.0 uses `gh auth` as the import source — works for the majority. The long tail v0.3.0 must address:
122: - Users who don't have `gh` installed (rare but real).
123: - Enterprise GHES instances where `gh` won't authenticate (custom auth gateways, SAML enforcement).
124: - Headless setups where `gh auth login`'s browser flow is impractical.
125: 
126: ### Scope
127: 
128: - Register a gitswitch OAuth app on github.com (one-time; `client_id` baked into the binary, no client_secret since device flow doesn't require one).
129: - Allow user-supplied `client_id` for enterprise GHES (each org registers their own OAuth app).
130: - Scopes: `repo`, `read:user`, `user:email`, `gist`, `workflow` (matches `gh` defaults so the token is a drop-in for `gh auth`).
131: - Command: `gitswitch login [--host HOST] [--client-id ID]`.
132: - Crate choice: [`github-device-flow`](https://crates.io/crates/github-device-flow) for the github.com path (purpose-built, implements device flow + refresh); [`oauth2`](https://crates.io/crates/oauth2) for the user-supplied-client-id enterprise GHES path (more configurable).
133: - Tokens land in the `gitswitch-secrets` keychain entry (same backend as v0.2.0 PAT storage).
134: 
135: ### Sketch
136: 
137: ```rust
138: use github_device_flow::DeviceFlow;
139: 
140: pub async fn login(host: &str, client_id: &str) -> anyhow::Result<(String, Vec<String>)> {
141:     let flow = DeviceFlow::start(
142:         client_id,
143:         host,
144:         &["repo", "read:user", "user:email", "gist", "workflow"],
145:     ).await?;
146:     println!("Visit {} and enter code: {}", flow.verification_uri, flow.user_code);
147:     let token = flow.poll().await?;
148:     Ok((token.access_token, token.scopes))
149: }
150: ```
151: 
152: ### Open questions
153: 
154: - Should the gitswitch OAuth app be owned by the gitswitch project or by each user's GitHub account?
155: - Refresh tokens: GitHub doesn't issue refresh tokens for device flow by default — UX for expiry?
156: 
157: ---
158: 
159: ## 4. Remote retargeting
160: 
161: `gitswitch remote retarget` — repo-host migration helper.
162: 
163: ### Why
164: 
165: When an org migrates between GitHub hosting models — An org's `github.corp.com` → `github.com/neworg` is the canonical case — every developer has hundreds of local clones with `origin` pointing at the old host. Manual `git remote set-url` per repo is painful.
166: 
167: ### Sketch
168: 
169: ```
170: gitswitch remote retarget --from github.corp.com --to github.com/neworg
171: ```
172: 
173: In each known repo (or the current repo):
174: 1. Walk `git remote -v` (via [`git2`](https://crates.io/crates/git2) `Repository::remotes()` + `find_remote()`, or shell out to `git remote set-url` for parity with the Go shell-out approach).
175: 2. For any remote matching `--from`, rewrite to the equivalent path under `--to`.
176: 3. Confirm before each rewrite (or `--yes` to skip).
177: 4. Optionally switch the active pin to the profile that owns the new host.
178: 
179: Walking known repos: read `gitswitch-core`'s history-recorded repo paths.
180: 
181: ```rust
182: use std::path::Path;
183: 
184: pub struct RetargetChange {
185:     pub repo: std::path::PathBuf,
186:     pub remote: String,
187:     pub from_url: String,
188:     pub to_url: String,
189: }
190: 
191: pub fn retarget(
192:     repo: &Path,
193:     from_host: &str,
194:     to_host: &str,
195:     dry_run: bool,
196: ) -> anyhow::Result<Vec<RetargetChange>>;
197: ```
198: 
199: ### Open questions
200: 
201: - Should this walk *all* repos under a configured root (e.g. `~/code`) or only the current repo?
202: - Path translation: `github.corp.com:teamA/repo` → `github.com/neworg:teamA/repo`? Or do paths get rewritten too (`neworg-teamA/repo`)?
203: 
204: ---
205: 
206: ## 5. Host auto-detection
207: 
208: On `cd` into a repo (via the existing shell hook), walk `git remote -v` and auto-recommend the profile whose `hosts` map covers all remotes.
209: 
210: ### Why
211: 
212: Reduces manual pinning. If a repo's only remote is `github.corp.com:team/x.git` and the user has exactly one profile with that host configured, gitswitch can confidently suggest that profile without the user having to `gitswitch pin` it.
213: 
214: ### Sketch
215: 
216: Extend `gitswitch-core::history::Recommend()`:
217: 1. Existing path: usage frequency ≥ 60% → recommend.
218: 2. New path: if no usage history yet AND exactly one profile covers every host in `git remote -v` → recommend that profile.
219: 3. If multiple profiles cover all hosts → no recommendation (ambiguous; defer to user).
220: 
221: Shell hook calls `gitswitch recommend --path .` and returns a profile name or empty string.
222: 
223: ### Open questions
224: 
225: - Should auto-detection auto-switch silently, or always nudge?
226: - Interaction with multi-pin: should auto-detection emit a multi-pin suggestion when different remotes are best served by different profiles?
227: 
228: ---
229: 
230: ## Likely sequencing for v0.3.0
231: 
232: 1. In-repo `.gitswitch.yaml` (foundation: in-repo config schema becomes the substrate for team presets).
233: 2. Host auto-detection (small, high-impact UX win; piggybacks on 0.2.0 shell hook).
234: 3. Team preset apply *without* signing (gated behind `--unsafe-trust` flag; collects feedback on schema before locking it in).
235: 4. Team preset signing (the big one).
236: 5. Own OAuth flow.
237: 6. Remote retargeting.
238: 
239: Order may shift based on which use cases surface from 0.2.0 adopters.
240: 
241: ---
242: 
243: ## Companion plans
244: 
245: - [Migration](roadmap-migration.md)
246: - [Multi-host](roadmap-v0.2.0.md)
247: - [Session isolation](roadmap-v0.2.0-isolation.md)
248: - [TUI](roadmap-v0.2.0-tui.md)

(End of file - total 248 lines)
</content>
```

---

## Assistant (Build · Kimi K2.7 Code · 34.8s)

_Thinking:_

Now I have all the rust-cancelled file contents. I need to add cancellation banners to all of them and to why-rust.md. Let me create a standard cancellation banner and apply it.

The banner should be something like:

```
> **CANCELLED — June 2026**
>
> The Rust migration plan described in this document has been cancelled.
> gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.
> These documents are kept for historical reference only.
>
> See the updated roadmap in `roadmap-full-timeline.md` for the current plan.
```

I'll add this after the title/header of each file. For why-rust.md, I already read it. For the rust-cancelled files, I have their contents.

Let me do all the edits. Since each file has a different first line, I'll need unique oldString for each.

For why-rust.md:
oldString: "# Why Rust — gitswitch Technical Rationale\n\n_Decision locked: May 2026. The Go codebase remains the shipping v0.1.x tool; Rust is the target for v0.2.0+. See [docs/rust/roadmap-migration.md](rust/roadmap-migration.md) for the migration plan._"

For roadmap-migration.md:
oldString: "# gitswitch — Go → Rust Migration Plan\n\n## 1. Context & motivation"

For roadmap-v0.2.0.md:
oldString: "# gitswitch v0.2.0 (Rust) — Multi-Host Identity Plan\n\n## Context"

For roadmap-v0.2.0-isolation.md:
oldString: "# gitswitch v0.2.0 — Session Isolation (Rust Port)\n\n> Rust translation of [roadmap-v0.2.0-isolation.md](../roadmap-v0.2.0-isolation.md). Mental model and all design decisions are identical; only the implementation idioms change. Read the Go plan first for full rationale.\n\n## Context"

For roadmap-v0.2.0-tui.md:
oldString: "# gitswitch v0.2.0 — TUI Refresh (Rust)\n\n> Rust translation of the [Go TUI plan](../roadmap-v0.2.0-tui.md). Information architecture, keybinds, command palette, wizards, three themes — identical. This document locks the Rust crate stack and the Elm-architecture state shape, and shows how the arcade aesthetic is preserved through [tachyonfx](https://github.com/junkdog/tachyonfx) effects layered atop [ratatui](https://ratatui.rs/) widgets.\n\n## Context"

For roadmap-v0.3.0.md:
oldString: "# gitswitch v0.3.0 — Roadmap (Rust)\n\nItems deferred from [v0.2.0 planning](roadmap-v0.2.0.md). v0.2.0 establishes multi-host identity, keychain-backed PATs, gh-driven profile import, multi-pin, session isolation, ratatui TUI with arcade theme. v0.3.0 builds on that foundation to address team-scale distribution, in-repo config, and the long tail of authentication scenarios."

Let me prepare the new strings. I'll keep it consistent.

Actually, for the rust-cancelled/ files, I should make the banner reference the updated docs. Let me use:

```
> **CANCELLED — June 2026**
>
> The Rust migration plan described in this document has been cancelled.
> gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.
> This document is kept for historical reference only.
>
> See the updated roadmap in [`roadmap-full-timeline.md`](../roadmap-full-timeline.md) for the current plan.
```

For files inside rust-cancelled/, the relative path to roadmap-full-timeline.md is `../roadmap-full-timeline.md`.
For why-rust.md which is at the docs/.never_gonna_give_you_up/ level, the path is `roadmap-full-timeline.md`.

Wait, why-rust.md is in docs/.never_gonna_give_you_up/ and references roadmap-migration.md as `rust/roadmap-migration.md`. After rename, that link would break. I should update that too.

Let me also update the first line of why-rust.md to reflect cancellation.

For why-rust.md:
oldString: "# Why Rust — gitswitch Technical Rationale\n\n_Decision locked: May 2026. The Go codebase remains the shipping v0.1.x tool; Rust is the target for v0.2.0+. See [docs/rust/roadmap-migration.md](rust/roadmap-migration.md) for the migration plan._"

newString: "# Why Rust — gitswitch Technical Rationale\n\n> **CANCELLED — June 2026**\n>\n> The Rust migration described in this document has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.\n> This document is kept for historical reference only.\n>\n> See the updated roadmap in [`roadmap-full-timeline.md`](roadmap-full-timeline.md) for the current plan.\n\n_Decision locked: May 2026. The Go codebase remains the shipping v0.1.x tool; Rust was the target for v0.2.0+. See [docs/rust-cancelled/roadmap-migration.md](rust-cancelled/roadmap-migration.md) for the archived migration plan._"

Hmm, but that might be confusing. Let me think. The user wants cancellation markers. I should make it clear. Maybe replace the original subtitle entirely with the cancellation banner.

Actually, better:
oldString: "# Why Rust — gitswitch Technical Rationale\n\n_Decision locked: May 2026. The Go codebase remains the shipping v0.1.x tool; Rust is the target for v0.2.0+. See [docs/rust/roadmap-migration.md](rust/roadmap-migration.md) for the migration plan._"

newString: "# Why Rust — gitswitch Technical Rationale\n\n> **CANCELLED — June 2026**\n>\n> The Rust migration described in this document has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.\n> This document is kept for historical reference only.\n>\n> See the updated roadmap in [`roadmap-full-timeline.md`](roadmap-full-timeline.md) for the current plan.\n\n_Original decision locked: May 2026. The Go codebase was the shipping v0.1.x tool; Rust was the target for v0.2.0+. See [docs/rust-cancelled/roadmap-migration.md](rust-cancelled/roadmap-migration.md) for the archived migration plan._"

Good. Let me do all these edits.

**Tool: edit**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/why-rust.md",
  "oldString": "# Why Rust — gitswitch Technical Rationale\n\n_Decision locked: May 2026. The Go codebase remains the shipping v0.1.x tool; Rust is the target for v0.2.0+. See [docs/rust/roadmap-migration.md](rust/roadmap-migration.md) for the migration plan._",
  "newString": "# Why Rust — gitswitch Technical Rationale\n\n> **CANCELLED — June 2026**\n>\n> The Rust migration described in this document has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.\n> This document is kept for historical reference only.\n>\n> See the updated roadmap in [`roadmap-full-timeline.md`](roadmap-full-timeline.md) for the current plan.\n\n_Original decision locked: May 2026. The Go codebase was the shipping v0.1.x tool; Rust was the target for v0.2.0+. See [docs/rust-cancelled/roadmap-migration.md](rust-cancelled/roadmap-migration.md) for the archived migration plan._"
}
```

**Output:**
```
Edit applied successfully.
```

---

## Assistant (Build · Kimi K2.7 Code · 20.7s)

_Thinking:_

Good. Now let me edit the rust-cancelled files. I'll do them in parallel since they're independent edits.

**Tool: edit**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-migration.md",
  "oldString": "# gitswitch — Go → Rust Migration Plan\n\n## 1. Context & motivation",
  "newString": "# gitswitch — Go → Rust Migration Plan\n\n> **CANCELLED — June 2026**\n>\n> The Rust migration plan described in this document has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.\n> This document is kept for historical reference only.\n>\n> See the updated roadmap in [`roadmap-full-timeline.md`](../roadmap-full-timeline.md) for the current plan.\n\n## 1. Context & motivation"
}
```

**Output:**
```
Edit applied successfully.
```

**Tool: edit**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-v0.2.0.md",
  "oldString": "# gitswitch v0.2.0 (Rust) — Multi-Host Identity Plan\n\n## Context",
  "newString": "# gitswitch v0.2.0 (Rust) — Multi-Host Identity Plan\n\n> **CANCELLED — June 2026**\n>\n> The Rust migration plan described in this document has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.\n> This document is kept for historical reference only.\n>\n> See the updated roadmap in [`roadmap-full-timeline.md`](../roadmap-full-timeline.md) for the current plan.\n\n## Context"
}
```

**Output:**
```
Edit applied successfully.
```

**Tool: edit**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-v0.2.0-isolation.md",
  "oldString": "# gitswitch v0.2.0 — Session Isolation (Rust Port)\n\n> Rust translation of [roadmap-v0.2.0-isolation.md](../roadmap-v0.2.0-isolation.md). Mental model and all design decisions are identical; only the implementation idioms change. Read the Go plan first for full rationale.\n\n## Context",
  "newString": "# gitswitch v0.2.0 — Session Isolation (Rust Port)\n\n> **CANCELLED — June 2026**\n>\n> The Rust migration plan described in this document has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.\n> This document is kept for historical reference only.\n>\n> See the updated roadmap in [`roadmap-full-timeline.md`](../roadmap-full-timeline.md) for the current plan.\n\n> Rust translation of [roadmap-v0.2.0-isolation.md](../roadmap-v0.2.0-isolation.md). Mental model and all design decisions are identical; only the implementation idioms change. Read the Go plan first for full rationale.\n\n## Context"
}
```

**Output:**
```
Edit applied successfully.
```

**Tool: edit**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-v0.2.0-tui.md",
  "oldString": "# gitswitch v0.2.0 — TUI Refresh (Rust)\n\n> Rust translation of the [Go TUI plan](../roadmap-v0.2.0-tui.md). Information architecture, keybinds, command palette, wizards, three themes — identical. This document locks the Rust crate stack and the Elm-architecture state shape, and shows how the arcade aesthetic is preserved through [tachyonfx](https://github.com/junkdog/tachyonfx) effects layered atop [ratatui](https://ratatui.rs/) widgets.\n\n## Context",
  "newString": "# gitswitch v0.2.0 — TUI Refresh (Rust)\n\n> **CANCELLED — June 2026**\n>\n> The Rust migration plan described in this document has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.\n> This document is kept for historical reference only.\n>\n> See the updated roadmap in [`roadmap-full-timeline.md`](../roadmap-full-timeline.md) for the current plan.\n\n> Rust translation of the [Go TUI plan](../roadmap-v0.2.0-tui.md). Information architecture, keybinds, command palette, wizards, three themes — identical. This document locks the Rust crate stack and the Elm-architecture state shape, and shows how the arcade aesthetic is preserved through [tachyonfx](https://github.com/junkdog/tachyonfx) effects layered atop [ratatui](https://ratatui.rs/) widgets.\n\n## Context"
}
```

**Output:**
```
Edit applied successfully.
```

**Tool: edit**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-v0.3.0.md",
  "oldString": "# gitswitch v0.3.0 — Roadmap (Rust)\n\nItems deferred from [v0.2.0 planning](roadmap-v0.2.0.md). v0.2.0 establishes multi-host identity, keychain-backed PATs, gh-driven profile import, multi-pin, session isolation, ratatui TUI with arcade theme. v0.3.0 builds on that foundation to address team-scale distribution, in-repo config, and the long tail of authentication scenarios.",
  "newString": "# gitswitch v0.3.0 — Roadmap (Rust)\n\n> **CANCELLED — June 2026**\n>\n> The Rust migration plan described in this document has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.\n> This document is kept for historical reference only.\n>\n> See the updated roadmap in [`roadmap-full-timeline.md`](../roadmap-full-timeline.md) for the current plan.\n\nItems deferred from [v0.2.0 planning](roadmap-v0.2.0.md). v0.2.0 establishes multi-host identity, keychain-backed PATs, gh-driven profile import, multi-pin, session isolation, ratatui TUI with arcade theme. v0.3.0 builds on that foundation to address team-scale distribution, in-repo config, and the long tail of authentication scenarios."
}
```

**Output:**
```
Edit applied successfully.
```

---

## Assistant (Build · Kimi K2.7 Code · 4.0s)

_Thinking:_

Good. Now I need to check if there are any references to `docs/rust/` in other files that need updating to `docs/rust-cancelled/`. Let me search for that.

**Tool: grep**

**Input:**
```json
{
  "pattern": "docs/rust/|\\.\\./rust/|rust/roadmap",
  "path": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up"
}
```

**Output:**
```
Found 8 matches
/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-migration.md:
  Line 129: └── docs/rust/                  # this folder

  Line 406: - [`docs/rust/roadmap-v0.2.0.md`](roadmap-v0.2.0.md) — multi-host identity in Rust (Profile / HostCreds / Pin types, keychain + encrypted-file `SecretStore`, gh-driven import, SSH config block manager). Sister doc to [`docs/roadmap-v0.2.0.md`](../roadmap-v0.2.0.md).

  Line 407: - [`docs/rust/roadmap-v0.2.0-isolation.md`](roadmap-v0.2.0-isolation.md) — session isolation via env vars + shell wrapper, scope-aware command dispatch (`--global` / `--session` flags), `--emit-env` IPC. Sister doc to [`docs/roadmap-v0.2.0-isolation.md`](../roadmap-v0.2.0-isolation.md).

  Line 408: - [`docs/rust/roadmap-v0.2.0-tui.md`](roadmap-v0.2.0-tui.md) — ratatui + tachyonfx implementation of the 5-screen IA, arcade theme as a tachyonfx-driven layer, custom CRT widget. Sister doc to [`docs/roadmap-v0.2.0-tui.md`](../roadmap-v0.2.0-tui.md).

  Line 409: - [`docs/rust/roadmap-v0.3.0.md`](roadmap-v0.3.0.md) — deferred items (in-repo YAML config, team preset distribution + signing, own OAuth as primary path, host auto-detection). Sister doc to [`docs/roadmap-v0.3.0.md`](../roadmap-v0.3.0.md).


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/why-rust.md:
  Line 117: gitswitch lives in the terminal, fires on every directory change, handles secrets, and needs a polished animated TUI. Those constraints pick Rust. The migration plan is in [docs/rust/roadmap-migration.md](rust/roadmap-migration.md); the v0.2.0 feature plans are in [docs/rust/](rust/).


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/project-memory.md:
  Line 25: - `docs/rust/roadmap-migration.md` — 6-phase migration plan, 9-crate workspace, cutover criteria

  Line 26: - `docs/rust/` — all v0.2.0+ plans in Rust (multi-host, session isolation, TUI, v0.3.0)

```

---

## Assistant (Build · Kimi K2.7 Code · 4.6s)

_Thinking:_

I need to update these references:
1. rust-cancelled/roadmap-migration.md lines 129, 406-409: update `docs/rust/` to `docs/rust-cancelled/`
2. why-rust.md line 117: update `docs/rust/` to `docs/rust-cancelled/` and `[docs/rust/](rust/)` to `[docs/rust-cancelled/](rust-cancelled/)`
3. project-memory.md lines 25-26: update to rust-cancelled

But project-memory.md might need more updates since the plan has changed. Let me read it first.

**Tool: read**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/project-memory.md"
}
```

**Output:**
```
<path>/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/project-memory.md</path>
<type>file</type>
<content>
1: # gitswitch — Project Memory
2: 
3: _All locked decisions, design principles, and personal notes. Written May 2026._
4: 
5: ---
6: 
7: ## 1. Rust Decision (locked)
8: 
9: gitswitch v0.2.0+ will be a clean rewrite in Rust. Go remains the shipping v0.1.x tool. Decision locked May 2026 — not up for debate.
10: 
11: **Why Rust, specifically:**
12: 
13: 1. **Sub-millisecond shell hook** — the hook fires on every `cd`. Go has a GC and a runtime. Rust has neither. Hook latency is the most user-visible performance characteristic of the entire tool.
14: 2. **Smaller binary** — ~1.5–3MB with LTO vs 4–15MB for Go. Matters for a tool people install via curl.
15: 3. **Ratatui + tachyonfx** — Rust's TUI ecosystem is best-in-class. tachyonfx (50+ composable effects) has no Go equivalent. The arcade theme requires it.
16: 4. **Compiler-enforced secret safety** — zeroize on drop, borrow checker prevents data races on SSH/GPG/PAT material. The tool handles credentials; the compiler should enforce correctness.
17: 5. **Crate ecosystem** — clap v4 (derive), serde + yaml_serde, keyring v3, reqwest, ratatui, insta. All mature, all production-ready.
18: 
19: **Go is not wrong** — it got gitswitch to v0.1.20 with real users. Rust is right for where the tool is going.
20: 
21: **Key note on yaml_serde**: `serde_yaml` was deprecated March 2024. Use `yaml_serde` as drop-in replacement via cargo rename: `serde_yaml = { package = "yaml_serde", version = "0.10" }`.
22: 
23: **Reference docs:**
24: - `docs/why-rust.md` — full technical rationale
25: - `docs/rust/roadmap-migration.md` — 6-phase migration plan, 9-crate workspace, cutover criteria
26: - `docs/rust/` — all v0.2.0+ plans in Rust (multi-host, session isolation, TUI, v0.3.0)
27: 
28: ---
29: 
30: ## 2. Product Philosophy (non-negotiable)
31: 
32: gitswitch must be **beginner-friendly on first touch**, with advanced features that power users can discover themselves. Progressive disclosure, not gatekeeping.
33: 
34: The product should work for someone who has never heard of SSH keys AND for a senior engineer managing 5 GitHub hosts simultaneously.
35: 
36: **The layered model:**
37: 
38: ```
39: Layer 0 — Zero config (works the moment you install)
40:   gitswitch install    → shell hook added
41:   gitswitch            → shows who you are right now
42:   if gh is logged in   → auto-imports identity as first profile
43:   Done in 20 seconds.
44: 
45: Layer 1 — Basic switching (intentional, manual)
46:   gitswitch add personal / gitswitch add work
47:   gitswitch switch work
48:   gitswitch alone      → TUI picker
49: 
50: Layer 2 — Set and forget (automatic)
51:   gitswitch pin        → this repo always uses this profile
52:   pre-commit warning   → warns before wrong-identity commit
53:   shell hook suggests  → "you usually work here as work. Switch? [y/N]"
54: 
55: Layer 3 — Power user (explicit opt-in)
56:   gitswitch shell      → session isolation, per-terminal identity
57:   multi-host profiles  → one profile, multiple GitHub hosts
58:   gitswitch status     → full per-repo diagnostic
59:   gitswitch doctor     → all profiles health check
60:   HTTPS credential helper → future
61:   remote-URL routing   → future
62: ```
63: 
64: **Rule:** When designing any feature, ask which layer it belongs to. Layer 0–1 must have zero required config and no jargon. Layer 2–3 can be opt-in and complex. Never surface Layer 3 concepts to someone who hasn't reached Layer 2.
65: 
66: Each layer should be self-discoverable from the previous. CLI output hints at the next step. The TUI is the discovery surface. The man page is the reference.
67: 
68: ---
69: 
70: ## 3. Arcade Philosophy (non-negotiable)
71: 
72: The arcade layer — hidden games and easter eggs — is gitswitch's personality, not decoration.
73: 
74: **Abhiram's words (May 2026):** _"The arcade easter eggs are our personality, hidden beneath gitswitch, which is what makes the user wanna open the TUI from time to time for the fun of it."_
75: 
76: **Rules:**
77: - Never move arcade content behind paywalls, external services, or opt-in prompts
78: - Hidden, not gatekept — discoverable by anyone, available to everyone
79: - Never interrupts workflow — only accessible when you go looking
80: - One great hidden game beats ten mediocre spinners
81: - Quality over quantity — each easter egg should be worth finding
82: 
83: **Current easter eggs:**
84: - `gitswitch pacman` — hidden dot-eater game using profile list as level (exists in v0.1.x)
85: 
86: **Planned (locked in arcade roadmap):**
87: - Roguelike dungeon using repos as rooms
88: - Constellation map from commit graph
89: - Ghost profiles haunting the TUI after deletion
90: - Deep unlock `:konami` → matrix cascade effect
91: 
92: **Technical home in Rust:** `gitswitch-arcade` crate, separated from core logic.
93: 
94: **Reference doc:** `docs/arcade-philosophy.md`
95: 
96: ---
97: 
98: ## 4. Honest Gaps + Future Features
99: 
100: Personal notes on where gitswitch genuinely falls behind and what to do about it.
101: 
102: ### Confirmed gaps to fix
103: 
104: **HTTPS users are invisible to gitswitch** _(v0.3.0+ candidate)_
105: Most beginners clone via HTTPS — GitHub shows it by default. gitswitch handles SSH auth + gh CLI auth but does nothing for HTTPS credential routing. Acting as a `credential.helper` fixes this: git asks gitswitch "what's the PAT for github.com?" and gitswitch returns the right one for the active profile. git-ego (102 stars, Go) does this today.
106: 
107: **First-install has too many steps** _(fix in v0.2.0)_
108: Ideal: run `gitswitch`, detect `gh` logged in, offer to import as first profile, done in 20 seconds. Fix via `gitswitch import gh` flow.
109: 
110: **No pre-commit safety net** _(v0.2.0)_
111: Shell hook suggests but doesn't enforce. A pre-commit hook warning ("you're about to commit as personal@gmail.com in a repo pinned to work — sure?") means the mistake can never happen silently. No active competitor has this with a clean install story. Auto-install with `gitswitch install`, user can opt out.
112: 
113: **No progressive CLI hints** _(v0.2.0 UX polish)_
114: After `gitswitch switch work` → hint "run `gitswitch pin` to make this automatic." After first pin → hint "run `gitswitch status` to see the full picture." Zero new code. Teaches the product through normal use.
115: 
116: ### Intentional non-defaults
117: 
118: **git includeIf** — fires identity at git-command time via git's native config loading, no shell modification needed. Competitors use it (git-ego, karn, git-idm). Intentionally NOT the default for gitswitch — the shell hook is what enables suggestions, learning, and TUI integration. Keep as optional future mode, never the default.
119: 
120: **Remote-URL-aware routing** — pin by git remote URL pattern (`github.com/mycompany/*` → work) instead of directory path. More robust for scattered repos. Not a priority — directory pinning covers 90% of cases. v0.3.0+ list.
121: 
122: ### Future roadmap features
123: 
124: | Feature | Target | Notes |
125: |---|---|---|
126: | HTTPS credential helper | v0.3.0+ | `git config credential.helper gitswitch` |
127: | Remote-URL-aware pinning | v0.3.0+ | Pattern-based, not directory-based |
128: | Ephemeral one-shot run | v0.3.0+ | `gitswitch run <profile> -- <cmd>` (gpx does this) |
129: | Shell tab completion | v0.2.0 | Should ship with v0.2.0 |
130: | Profile export/import | v0.3.0+ | JSON backup for new machine onboarding |
131: | git includeIf (optional) | future | Zero-shell-modification alternative mode |
132: | Windows native support | low priority | WSL covers most cases |
133: 
134: ---
135: 
136: ## 5. What gitswitch has that nobody else does
137: 
138: 1. **TUI with arcade personality** — no active competitor has a real TUI; arcade layer is completely unique
139: 2. **Shell hook + learn + suggest** — fires on `cd`, builds per-repo history, suggests the right profile
140: 3. **gh CLI sync** — keeping git identity and `gh` auth in sync in one command; entirely unique
141: 4. **Single static binary, no runtime deps** — curl install; npm tools need Node, Python tools need Python, go install needs Go
142: 5. **Session isolation** (v0.2.0) — per-shell env-var-scoped identity for parallel AI agents; no competitor touches this
143: 6. **Multi-host profiles** (v0.2.0) — one profile covering multiple GitHub hosts with different SSH keys + PATs per host
144: 
145: ---
146: 
147: ## 6. Distribution strategy (from competitor research)
148: 
149: - **homebrew-core slot is open** — karn holds it with 71 installs/year (dead since 2021). No active competitor is in core. Every competitor is tap-only (invisible to `brew search`). Submit gitswitch after v0.2.0 ships — Homebrew bar is ~75 stars + stable release.
150: - **Set up tap now** — `brew tap aksisonline/gitswitch` as install path while core review is pending
151: - **crates.io is empty** — every Rust git identity crate is dead. gitswitch v0.2.0 owns the category by default on publish.
152: - **Do not bother with npm or pip** — both channels are dead for this problem space.
153: 
154: ---
155: 
156: ## 7. Monetization (privacy-first, locked)
157: 
158: No cloud. No telemetry. No feature flags. No paywalls. Everything free, forever — including enterprise features.
159: 
160: Revenue model: GitHub Sponsors + Open Collective + occasional consulting + open source grants.
161: 
162: **Sponsor tiers:** $5/mo (supporter) · $20/mo (power user) · $100/mo (company) · $500/mo (company + SLA) · $1000/mo (lead sponsor)
163: 
164: **Ceiling:** $3K–25K/month at maturity. One developer, fully sustained. That's the honest ceiling and it's enough.
165: 
166: **Never:** cloud sync, telemetry (even opt-in), open core / feature gating, ads, paid arcade content, selling data.
167: 
168: ---
169: 
170: ## 8. Key locked decisions summary
171: 
172: | Decision | Status | Notes |
173: |---|---|---|
174: | Rust for v0.2.0+ | Locked | Go stays for v0.1.x |
175: | yaml_serde (not serde_yaml) | Locked | serde_yaml deprecated Mar 2024 |
176: | PAT storage in OS keychain | Locked | Never plaintext |
177: | Arcade layer never paywalled | Locked | Core product personality |
178: | Layer 0 zero-config install | Locked | Non-negotiable UX constraint |
179: | No cloud features | Locked | Privacy-first always |
180: | git includeIf not default | Decided | Optional future mode only |
181: | Pre-commit safety net v0.2.0 | Planned | Auto-install with `gitswitch install` |
182: | homebrew-core submission | Post-v0.2.0 | After stable release + star count |

(End of file - total 182 lines)
</content>
```

---

## Assistant (Build · Kimi K2.7 Code · 22.4s)

_Thinking:_

I need to update project-memory.md to reflect the cancellation. The "Rust Decision (locked)" section needs to be changed. Also update references to `docs/rust/` to `docs/rust-cancelled/`. And maybe update the future roadmap features since OAuth is now v1.0 and HTTPS credential helper might move to v1.0 too.

But the user wants me to finish cancelling the plan and then write new docs with new plans. So I should at least update the Rust decision section and references. I can do a more comprehensive update later when writing the new plan docs.

Let me update project-memory.md:
1. Section 1: Change from "Rust Decision (locked)" to "Rust Decision (CANCELLED)" and explain Go remains.
2. Update references from docs/rust/ to docs/rust-cancelled/
3. Update key locked decisions summary table

I should also update the technical home in arcade section from `gitswitch-arcade` crate to something Go-appropriate, but maybe that's for later.

Let me do the edits to project-memory.md now.

**Tool: edit**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/project-memory.md",
  "oldString": "## 1. Rust Decision (locked)\n\ngitswitch v0.2.0+ will be a clean rewrite in Rust. Go remains the shipping v0.1.x tool. Decision locked May 2026 — not up for debate.\n\n**Why Rust, specifically:**\n\n1. **Sub-millisecond shell hook** — the hook fires on every `cd`. Go has a GC and a runtime. Rust has neither. Hook latency is the most user-visible performance characteristic of the entire tool.\n2. **Smaller binary** — ~1.5–3MB with LTO vs 4–15MB for Go. Matters for a tool people install via curl.\n3. **Ratatui + tachyonfx** — Rust's TUI ecosystem is best-in-class. tachyonfx (50+ composable effects) has no Go equivalent. The arcade theme requires it.\n4. **Compiler-enforced secret safety** — zeroize on drop, borrow checker prevents data races on SSH/GPG/PAT material. The tool handles credentials; the compiler should enforce correctness.\n5. **Crate ecosystem** — clap v4 (derive), serde + yaml_serde, keyring v3, reqwest, ratatui, insta. All mature, all production-ready.\n\n**Go is not wrong** — it got gitswitch to v0.1.20 with real users. Rust is right for where the tool is going.\n\n**Key note on yaml_serde**: `serde_yaml` was deprecated March 2024. Use `yaml_serde` as drop-in replacement via cargo rename: `serde_yaml = { package = \"yaml_serde\", version = \"0.10\" }`.\n\n**Reference docs:**\n- `docs/why-rust.md` — full technical rationale\n- `docs/rust/roadmap-migration.md` — 6-phase migration plan, 9-crate workspace, cutover criteria\n- `docs/rust/` — all v0.2.0+ plans in Rust (multi-host, session isolation, TUI, v0.3.0)",
  "newString": "## 1. Rust Decision (CANCELLED)\n\n> **CANCELLED — June 2026**\n>\n> The Rust rewrite plan is cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.\n> The cancellation reason: the migration risk outweighs the marginal gains for a tool that already works. Engineering effort is better spent on user-facing features (OAuth, session isolation, multi-host, HTTPS credential helper) than on a binary-compatible rewrite.\n> This section is kept for historical reference.\n\n~~gitswitch v0.2.0+ will be a clean rewrite in Rust. Go remains the shipping v0.1.x tool. Decision locked May 2026 — not up for debate.~~\n\n~~**Why Rust, specifically:**~~\n\n1. ~~**Sub-millisecond shell hook** — the hook fires on every `cd`. Go has a GC and a runtime. Rust has neither. Hook latency is the most user-visible performance characteristic of the entire tool.~~\n2. ~~**Smaller binary** — ~1.5–3MB with LTO vs 4–15MB for Go. Matters for a tool people install via curl.~~\n3. ~~**Ratatui + tachyonfx** — Rust's TUI ecosystem is best-in-class. tachyonfx (50+ composable effects) has no Go equivalent. The arcade theme requires it.~~\n4. ~~**Compiler-enforced secret safety** — zeroize on drop, borrow checker prevents data races on SSH/GPG/PAT material. The tool handles credentials; the compiler should enforce correctness.~~\n5. ~~**Crate ecosystem** — clap v4 (derive), serde + yaml_serde, keyring v3, reqwest, ratatui, insta. All mature, all production-ready.~~\n\n**Go is not wrong** — it got gitswitch to v0.1.x with real users. It remains the right choice for v1.0 and beyond.\n\n**Reference docs (archived):**\n- `docs/why-rust.md` — original technical rationale (now marked cancelled)\n- `docs/rust-cancelled/roadmap-migration.md` — archived 6-phase migration plan\n- `docs/rust-cancelled/` — all archived v0.2.0+ Rust plans"
}
```

**Output:**
```
Edit applied successfully.
```

---

## Assistant (Build · Kimi K2.7 Code · 4.9s)

_Thinking:_

Now I need to update the key locked decisions summary table at the bottom. Let me edit that too.

**Tool: edit**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/project-memory.md",
  "oldString": "| Rust for v0.2.0+ | Locked | Go stays for v0.1.x |",
  "newString": "| Rust for v0.2.0+ | CANCELLED | Go remains for all versions; see `docs/rust-cancelled/` |"
}
```

**Output:**
```
Edit applied successfully.
```

---

## Assistant (Build · Kimi K2.7 Code · 4.8s)

_Thinking:_

Now update the references in rust-cancelled/roadmap-migration.md. Lines 129 and 406-409. Let me read the relevant sections to get exact strings.

**Tool: read**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-migration.md",
  "offset": 120,
  "limit": 15
}
```

**Output:**
```
<path>/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-migration.md</path>
<type>file</type>
<content>
120: │   ├── gitswitch-env/          # session env builder (shared shell/env/TUI exit path)
121: │   ├── gitswitch-tui/          # ratatui app, screens, theme, view models
122: │   ├── gitswitch-arcade/       # splash, tachyonfx effects, CRT widget
123: │   └── gitswitch-shell/        # install/uninstall/hook detection (zsh/bash/fish)
124: ├── tests/
125: │   ├── parity/                 # parity corpus shell scripts
126: │   ├── snapshots/              # insta snapshots (YAML + TUI frames)
127: │   └── integration/            # cargo test integration suites
128: ├── legacy/                     # archived Go tree, kept for reference until 0.2.1
129: └── docs/rust/                  # this folder
130: ```
131: 
132: ### Key dependency graph
133: 
134: ```

(Showing lines 120-134 of 411. Use offset=135 to continue.)
</content>
```

**Tool: read**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-migration.md",
  "offset": 400,
  "limit": 15
}
```

**Output:**
```
<path>/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-migration.md</path>
<type>file</type>
<content>
400: ---
401: 
402: ## 10. Companion plans
403: 
404: The Rust workspace consumes the same feature decisions as the Go roadmap — only the implementation language changes. Per-feature Rust roadmaps live alongside this file:
405: 
406: - [`docs/rust/roadmap-v0.2.0.md`](roadmap-v0.2.0.md) — multi-host identity in Rust (Profile / HostCreds / Pin types, keychain + encrypted-file `SecretStore`, gh-driven import, SSH config block manager). Sister doc to [`docs/roadmap-v0.2.0.md`](../roadmap-v0.2.0.md).
407: - [`docs/rust/roadmap-v0.2.0-isolation.md`](roadmap-v0.2.0-isolation.md) — session isolation via env vars + shell wrapper, scope-aware command dispatch (`--global` / `--session` flags), `--emit-env` IPC. Sister doc to [`docs/roadmap-v0.2.0-isolation.md`](../roadmap-v0.2.0-isolation.md).
408: - [`docs/rust/roadmap-v0.2.0-tui.md`](roadmap-v0.2.0-tui.md) — ratatui + tachyonfx implementation of the 5-screen IA, arcade theme as a tachyonfx-driven layer, custom CRT widget. Sister doc to [`docs/roadmap-v0.2.0-tui.md`](../roadmap-v0.2.0-tui.md).
409: - [`docs/rust/roadmap-v0.3.0.md`](roadmap-v0.3.0.md) — deferred items (in-repo YAML config, team preset distribution + signing, own OAuth as primary path, host auto-detection). Sister doc to [`docs/roadmap-v0.3.0.md`](../roadmap-v0.3.0.md).
410: 
411: The Go roadmaps under [`docs/`](../) remain the canonical feature specs; the Rust roadmaps are pure implementation notes layered on top.

(End of file - total 411 lines)
</content>
```

---

## Assistant (Build · Kimi K2.7 Code · 13.4s)

_Thinking:_

I'll update those references. Since the file is archived, I could leave them, but better to keep links working. Let me edit them.

**Tool: edit**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-migration.md",
  "oldString": "└── docs/rust/                  # this folder",
  "newString": "└── docs/rust-cancelled/        # this folder (archived)"
}
```

**Output:**
```
Edit applied successfully.
```

**Tool: edit**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-migration.md",
  "oldString": "- [`docs/rust/roadmap-v0.2.0.md`](roadmap-v0.2.0.md) — multi-host identity in Rust (Profile / HostCreds / Pin types, keychain + encrypted-file `SecretStore`, gh-driven import, SSH config block manager). Sister doc to [`docs/roadmap-v0.2.0.md`](../roadmap-v0.2.0.md).\n- [`docs/rust/roadmap-v0.2.0-isolation.md`](roadmap-v0.2.0-isolation.md) — session isolation via env vars + shell wrapper, scope-aware command dispatch (`--global` / `--session` flags), `--emit-env` IPC. Sister doc to [`docs/roadmap-v0.2.0-isolation.md`](../roadmap-v0.2.0-isolation.md).\n- [`docs/rust/roadmap-v0.2.0-tui.md`](roadmap-v0.2.0-tui.md) — ratatui + tachyonfx implementation of the 5-screen IA, arcade theme as a tachyonfx-driven layer, custom CRT widget. Sister doc to [`docs/roadmap-v0.2.0-tui.md`](../roadmap-v0.2.0-tui.md).\n- [`docs/rust/roadmap-v0.3.0.md`](roadmap-v0.3.0.md) — deferred items (in-repo YAML config, team preset distribution + signing, own OAuth as primary path, host auto-detection). Sister doc to [`docs/roadmap-v0.3.0.md`](../roadmap-v0.3.0.md).",
  "newString": "- [`docs/rust-cancelled/roadmap-v0.2.0.md`](roadmap-v0.2.0.md) — multi-host identity in Rust (Profile / HostCreds / Pin types, keychain + encrypted-file `SecretStore`, gh-driven import, SSH config block manager). Sister doc to [`docs/roadmap-v0.2.0.md`](../roadmap-v0.2.0.md).\n- [`docs/rust-cancelled/roadmap-v0.2.0-isolation.md`](roadmap-v0.2.0-isolation.md) — session isolation via env vars + shell wrapper, scope-aware command dispatch (`--global` / `--session` flags), `--emit-env` IPC. Sister doc to [`docs/roadmap-v0.2.0-isolation.md`](../roadmap-v0.2.0-isolation.md).\n- [`docs/rust-cancelled/roadmap-v0.2.0-tui.md`](roadmap-v0.2.0-tui.md) — ratatui + tachyonfx implementation of the 5-screen IA, arcade theme as a tachyonfx-driven layer, custom CRT widget. Sister doc to [`docs/roadmap-v0.2.0-tui.md`](../roadmap-v0.2.0-tui.md).\n- [`docs/rust-cancelled/roadmap-v0.3.0.md`](roadmap-v0.3.0.md) — deferred items (in-repo YAML config, team preset distribution + signing, own OAuth as primary path, host auto-detection). Sister doc to [`docs/roadmap-v0.3.0.md`](../roadmap-v0.3.0.md)."
}
```

**Output:**
```
Edit applied successfully.
```

---

## Assistant (Build · Kimi K2.7 Code · 3.8s)

_Thinking:_

Now update why-rust.md line 117 reference. Let me read it to get exact string.

**Tool: read**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/why-rust.md",
  "offset": 108,
  "limit": 15
}
```

**Output:**
```
<path>/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/why-rust.md</path>
<type>file</type>
<content>
108: - **Cross-compilation**: Both are excellent; Go is slightly simpler (`GOOS=windows go build` vs `cross`).
109: - **Backend services**: If gitswitch ever grows a managed profiles API or team-preset hosting service, Go remains the better fit there.
110: 
111: The Go codebase is not a mistake — it got gitswitch to v0.1.20 with real users. Rust is the right choice for where the tool is going.
112: 
113: ---
114: 
115: ## The verdict
116: 
117: gitswitch lives in the terminal, fires on every directory change, handles secrets, and needs a polished animated TUI. Those constraints pick Rust. The migration plan is in [docs/rust/roadmap-migration.md](rust/roadmap-migration.md); the v0.2.0 feature plans are in [docs/rust/](rust/).
118: 
119: > "If gitswitch.dev were a cloud-hosted web app or API backend, Go would be the faster path to ship. But because gitswitch lives strictly in the user's terminal, interacts directly with local files, and intercepts shell navigation, Rust provides the performance, safety, and modern toolset required to make it a world-class developer utility."

(End of file - total 119 lines)
</content>
```

---

## Assistant (Build · Kimi K2.7 Code · 8.9s)

_Thinking:_

I need to update the verdict section to reflect cancellation. Let me edit it.

**Tool: edit**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/why-rust.md",
  "oldString": "## The verdict\n\ngitswitch lives in the terminal, fires on every directory change, handles secrets, and needs a polished animated TUI. Those constraints pick Rust. The migration plan is in [docs/rust/roadmap-migration.md](rust/roadmap-migration.md); the v0.2.0 feature plans are in [docs/rust/](rust/).\n\n> \"If gitswitch.dev were a cloud-hosted web app or API backend, Go would be the faster path to ship. But because gitswitch lives strictly in the user's terminal, interacts directly with local files, and intercepts shell navigation, Rust provides the performance, safety, and modern toolset required to make it a world-class developer utility.\"",
  "newString": "## The verdict\n\n> **CANCELLED — June 2026**\n>\n> The Rust migration plan has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.\n> The original rationale below is preserved for historical context only.\n\n~~gitswitch lives in the terminal, fires on every directory change, handles secrets, and needs a polished animated TUI. Those constraints pick Rust. The migration plan is in [docs/rust/roadmap-migration.md](rust/roadmap-migration.md); the v0.2.0 feature plans are in [docs/rust/](rust/).~~\n\nThe archived migration plan is in [`docs/rust-cancelled/roadmap-migration.md`](rust-cancelled/roadmap-migration.md); the archived v0.2.0+ Rust feature plans are in [`docs/rust-cancelled/`](rust-cancelled/).\n\n> \"If gitswitch.dev were a cloud-hosted web app or API backend, Go would be the faster path to ship. But because gitswitch lives strictly in the user's terminal, interacts directly with local files, and intercepts shell navigation, Rust provides the performance, safety, and modern toolset required to make it a world-class developer utility.\""
}
```

**Output:**
```
Edit applied successfully.
```

---

## Assistant (Build · Kimi K2.7 Code · 4.4s)

_Thinking:_

Also line 111 still says "Rust is the right choice for where the tool is going." I should strike that too since it's now cancelled. Let me update it.

**Tool: edit**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/why-rust.md",
  "oldString": "The Go codebase is not a mistake — it got gitswitch to v0.1.20 with real users. Rust is the right choice for where the tool is going.",
  "newString": "The Go codebase is not a mistake — it got gitswitch to v0.1.x with real users. ~~Rust is the right choice for where the tool is going.~~ Go remains the right choice for v1.0 and beyond."
}
```

**Output:**
```
Edit applied successfully.
```

---

## Assistant (Build · Kimi K2.7 Code · 6.5s)

_Thinking:_

Good. Now let me update the todo list and move on to revamping the roadmaps. But first, I should check the current project version and understand the codebase. The user mentioned they're on v0.1.22. Let me check the version file.

**Tool: todowrite**

**Input:**
```json
{
  "todos": [
    {
      "content": "Add cancellation markers to Rust docs and rename rust/ folder",
      "status": "completed",
      "priority": "high"
    },
    {
      "content": "Revamp roadmap-full-timeline.md with new v1 plan (Go stable, OAuth device flow, gs alias)",
      "status": "in_progress",
      "priority": "high"
    },
    {
      "content": "Revamp roadmap-v0.3.0.md to remove OAuth and reflect new sequencing",
      "status": "pending",
      "priority": "high"
    },
    {
      "content": "Check current project version and existing features to ground the new versioning plan",
      "status": "pending",
      "priority": "high"
    },
    {
      "content": "Draft new versioning/roadmap docs for v0.2.x, v0.3.x, v1.0 with PAT mgmt, session isolation, OAuth, gs",
      "status": "pending",
      "priority": "high"
    }
  ]
}
```

**Output:**
```
[
  {
    "content": "Add cancellation markers to Rust docs and rename rust/ folder",
    "status": "completed",
    "priority": "high"
  },
  {
    "content": "Revamp roadmap-full-timeline.md with new v1 plan (Go stable, OAuth device flow, gs alias)",
    "status": "in_progress",
    "priority": "high"
  },
  {
    "content": "Revamp roadmap-v0.3.0.md to remove OAuth and reflect new sequencing",
    "status": "pending",
    "priority": "high"
  },
  {
    "content": "Check current project version and existing features to ground the new versioning plan",
    "status": "pending",
    "priority": "high"
  },
  {
    "content": "Draft new versioning/roadmap docs for v0.2.x, v0.3.x, v1.0 with PAT mgmt, session isolation, OAuth, gs",
    "status": "pending",
    "priority": "high"
  }
]
```

---

## Assistant (Build · Kimi K2.7 Code · 5.2s)

_Thinking:_

Now let me check the current project version and existing code structure. I need to read relevant files. Let me search for version info and look at the Go code structure.

**Tool: bash**

**Input:**
```json
{
  "command": "git status --short && git log --oneline -5",
  "description": "Check git status and recent commits"
}
```

**Output:**
```
adc33e1 chore: fix issue template for CLI context and clarify Formula is a template
6822cb5 Documentation update
95ddb0f chore: update install script for v0.1.22
2cf27c6 Merge pull request #14 from aksisonline/fix/hook-check-version-noise
9f87782 fix: hook-check dev version noise + combined upgrade+HTTPS nudge

```

**Tool: glob**

**Input:**
```json
{
  "pattern": "**/version*"
}
```

**Output:**
```
No files found
```

**Tool: glob**

**Input:**
```json
{
  "pattern": "cmd/**/*.go"
}
```

**Output:**
```
/Users/aksisonline/Documents/Projects/Personal/git-switcher/cmd/gitswitch/main.go
```

---

## Assistant (Build · Kimi K2.7 Code · 4.6s)

_Thinking:_

No version file glob match. Let me look at main.go and also search for version constants.

**Tool: bash**

**Input:**
```json
{
  "command": "rg -n \"Version|version\" cmd/gitswitch/main.go internal/version/*.go 2>/dev/null || rg -n \"Version|version\" cmd/gitswitch/main.go",
  "description": "Find version constants in code"
}
```

**Output:**
```
cmd/gitswitch/main.go:17:	ver "github.com/aksisonline/gitswitch/internal/version"
cmd/gitswitch/main.go:25:var version = "dev"
cmd/gitswitch/main.go:54:		m, err := tui.New(store, version)
cmd/gitswitch/main.go:233:var versionCmd = &cobra.Command{
cmd/gitswitch/main.go:234:	Use:   "version",
cmd/gitswitch/main.go:235:	Short: "Show current version and check for updates",
cmd/gitswitch/main.go:237:		fmt.Printf("gitswitch %s\n", version)
cmd/gitswitch/main.go:238:		latest := ver.CachedLatestVersion(store.ConfigDir())
cmd/gitswitch/main.go:239:		if latest != "" && ver.IsUpdateAvailable(version, latest) {
cmd/gitswitch/main.go:240:			fmt.Printf("New version available: %s\n", latest)
cmd/gitswitch/main.go:247:			fmt.Println("Already on latest version.")
cmd/gitswitch/main.go:260:		m, err := tui.New(store, version, tui.WithArcadeMode())
cmd/gitswitch/main.go:296:	Short: "Upgrade gitswitch to the latest version",
cmd/gitswitch/main.go:304:		latest, err := ver.FetchLatestVersionFresh()
cmd/gitswitch/main.go:306:			return fmt.Errorf("could not fetch latest version: %w", err)
cmd/gitswitch/main.go:308:		if !ver.IsUpdateAvailable(version, latest) {
cmd/gitswitch/main.go:309:			fmt.Printf("Already on latest version (%s).\n", version)
cmd/gitswitch/main.go:312:		fmt.Printf("Upgrading %s → %s...\n", version, latest)
cmd/gitswitch/main.go:323:			fmt.Println("New features available in this version.")
cmd/gitswitch/main.go:335:			fmt.Println("Restart your shell to use the new version.")
cmd/gitswitch/main.go:508:			_ = shell.WriteHookVersion(configDir, version)
cmd/gitswitch/main.go:530:		if msg := shell.HookUpdateMessage(configDir, rcFile, version, git.IsCredentialHelperInstalled()); msg != "" {
cmd/gitswitch/main.go:609:	rootCmd.AddCommand(addCmd, switchCmd, listCmd, removeCmd, currentCmd, initCmd, versionCmd, upgradeCmd, pacmanCmd, pinCmd, unpinCmd, recordCmd, recommendCmd, installCmd, uninstallCmd, claudeCmd, hookCheckCmd, credentialCmd)
internal/version/check.go:1:package version
internal/version/check.go:25:	LatestVersion string    `json:"latest_version"`
internal/version/check.go:29:// CachedLatestVersion returns the latest version using a 24-hour disk cache.
internal/version/check.go:31:func CachedLatestVersion(configDir string) string {
internal/version/check.go:32:	cachePath := filepath.Join(configDir, "version-check.json")
internal/version/check.go:37:			return c.LatestVersion
internal/version/check.go:44:	latest, err := fetchLatestVersion(ctx)
internal/version/check.go:53:// FetchLatestVersionFresh always fetches from GitHub API, bypassing the cache.
internal/version/check.go:55:func FetchLatestVersionFresh() (string, error) {
internal/version/check.go:58:	return fetchLatestVersion(ctx)
internal/version/check.go:61:func fetchLatestVersion(ctx context.Context) (string, error) {
internal/version/check.go:88:func saveCache(cachePath, version string) error {
internal/version/check.go:89:	data, err := json.Marshal(cache{LatestVersion: version, CheckedAt: time.Now()})
internal/version/check.go:128:// UpgradeCommand returns a configured exec.Cmd for upgrading to targetVersion.
internal/version/check.go:129:// Returns an error if targetVersion fails semver validation (prevents shell injection).
internal/version/check.go:130:func UpgradeCommand(targetVersion string) (*exec.Cmd, error) {
internal/version/check.go:131:	if !semverRe.MatchString(targetVersion) {
internal/version/check.go:132:		return nil, fmt.Errorf("invalid version format: %q", targetVersion)
internal/version/check.go:134:	script := fmt.Sprintf(`curl -fsSL %s | bash -s -- %s`, installScriptURL, targetVersion)
internal/version/check.go:142:// RunUpgrade downloads and runs the install script to upgrade to the given version.
internal/version/check.go:143:func RunUpgrade(targetVersion string) error {
internal/version/check.go:144:	cmd, err := UpgradeCommand(targetVersion)

```

**Tool: bash**

**Input:**
```json
{
  "command": "ls internal/ && ls cmd/",
  "description": "List internal and cmd directories"
}
```

**Output:**
```
credential
git
history
install
shell
storage
tui
version
gitswitch

```

---

## Assistant (Build · Kimi K2.7 Code · 4.7s)

_Thinking:_

Good. The version is set via ldflags at build time, defaulting to "dev". The current release is v0.1.22. Now I need to understand the existing commands and features to ground the new versioning plan.

Let me read the main.go to see all commands, and maybe read the roadmap-v0.2.0.md since the rust version was based on it but the Go version should exist too.

**Tool: read**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/roadmap-v0.2.0.md"
}
```

**Output:**
```
<path>/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/roadmap-v0.2.0.md</path>
<type>file</type>
<content>
1: # gitswitch v0.2.0 — Multi-Host Identity Plan
2: 
3: ## Context
4: 
5: Current gitswitch (v0.1.x) treats identity as `(name, email, ssh_key, gpg_key, gh_user)` against a single implicit host (`github.com`). an enterprise multi-host use case — and other enterprise multi-host environments (finance, healthcare, defense, consulting) — surfaces a gap: developers operate against **two GitHub hosts simultaneously** with different PAT tokens, different SSH keys, different `gh` contexts, often in the same repo (mirror pushes) or in different repos in the same day.
6: 
7: v0.2.0 closes the gap by making **host a first-class dimension of identity**, **leveraging `gh auth` as the import source** (rather than building our own OAuth), and adding per-repo diagnostics (`gitswitch status`) plus multi-profile pinning. Team-preset YAML distribution, our own OAuth flow, and in-repo config files are explicitly v0.3.0 scope — see [roadmap-v0.3.0.md](roadmap-v0.3.0.md).
8: 
9: ### Locked design decisions
10: 
11: 1. **Per-host credentials map inside one profile** (not one-profile-per-host).
12: 2. **PAT storage**: OS keychain primary, gh `hosts.yml` synced on switch, encrypted-file fallback for headless/CI.
13: 3. **OAuth strategy**: don't roll our own. **Import from `gh auth status` / `~/.config/gh/hosts.yml`** — gh CLI users already authenticated their accounts; we discover them, create gitswitch profiles, then layer our customizations (per-host SSH keys, per-repo pinning) on top.
14: 4. **No remote URL rewrite on switch** — host is per-remote; gitswitch makes credentials correct for whichever host git resolves at push time.
15: 5. **`gitswitch init` repurposed**: scaffolds local YAML config in `~/.config/gitswitch/` (not in-repo — in-repo YAML is v0.3.0) and imports current git config + gh auth state in one shot.
16: 6. **New `gitswitch status`** — per-repo narrative diagnostic.
17: 7. **Multi-pin** — a repo can pin different profiles per remote (origin → corp, fork → personal).
18: 
19: ---
20: 
21: ## Profile model (storage layer)
22: 
23: Extend [internal/storage/store.go:10](../internal/storage/store.go) `Profile`:
24: 
25: ```go
26: type Profile struct {
27:     Nickname  string                `yaml:"nickname"`
28:     UserName  string                `yaml:"user_name"`
29:     Email     string                `yaml:"email"`
30:     SignKey   string                `yaml:"sign_key,omitempty"`
31:     Active    bool                  `yaml:"active"`
32:     Source    string                `yaml:"source,omitempty"`   // "manual" | "gh-import" | "imported"
33:     Hosts     map[string]HostCreds  `yaml:"hosts"`
34: }
35: 
36: type HostCreds struct {
37:     SSHKey   string `yaml:"ssh_key,omitempty"`
38:     GHUser   string `yaml:"gh_user,omitempty"`
39:     TokenRef string `yaml:"token_ref,omitempty"` // keychain key
40:     GHSynced bool   `yaml:"gh_synced,omitempty"` // token currently mirrored to gh hosts.yml
41: }
42: ```
43: 
44: Top-level `SSHKey` / `GHUser` removed from `Profile` — they migrate into `Hosts["github.com"]` during the v1→v2 upgrade.
45: 
46: ### Storage file format
47: 
48: Move from `~/.config/gitswitch/profiles.json` to `~/.config/gitswitch/config.yaml` (single file: version + prefs + profiles + pins).
49: 
50: ```yaml
51: version: 2
52: prefs:
53:   color_theme: 0
54: profiles:
55:   - nickname: personal
56:     user_name: Abhiram Kanna
57:     email: abhiram@personal.dev
58:     sign_key: 4A...
59:     active: true
60:     source: manual
61:     hosts:
62:       github.com:
63:         ssh_key: ~/.ssh/id_personal
64:         gh_user: aksisonline
65:         token_ref: gitswitch:personal:github.com
66:         gh_synced: true
67:   - nickname: corp
68:     user_name: Abhiram Kanna
69:     email: abhiram@corp.com
70:     active: false
71:     source: gh-import
72:     hosts:
73:       github.com:
74:         ssh_key: ~/.ssh/id_personal
75:         gh_user: aksisonline
76:         token_ref: gitswitch:personal:github.com
77:       github.corp.com:
78:         ssh_key: ~/.ssh/id_corp
79:         gh_user: aksisonline-pp
80:         token_ref: gitswitch:corp:github.corp.com
81: pins:
82:   /Users/abhi/code/corp-platform:
83:     default: corp
84:     remotes:
85:       fork: personal
86: ```
87: 
88: ### Migration (v1 JSON → v2 YAML)
89: 
90: On `Load()`, if `profiles.json` exists and `config.yaml` doesn't:
91: 1. Read JSON via the existing legacy path.
92: 2. Promote each v1 profile to the hosts-map shape, populating a single `github.com` entry from the old `SSHKey` / `GHUser` fields.
93: 3. Write `config.yaml` (mode 0600).
94: 4. Rename `profiles.json` → `profiles.json.v1.bak` (don't delete — safety net).
95: 
96: YAML library: `gopkg.in/yaml.v3` (supports comments, which `init` will use to emit a scaffolded config with annotations).
97: 
98: ---
99: 
100: ## gh-driven import (replaces own-OAuth path)
101: 
102: New package `internal/gh/import.go`:
103: 
104: ```go
105: type ImportedAccount struct {
106:     Hostname string
107:     Login    string
108:     Token    string
109:     Scopes   []string
110:     Source   string // "hosts.yml" | "gh-auth-status"
111: }
112: 
113: func Discover() ([]ImportedAccount, error)
114: func EnrichAccount(acct ImportedAccount) (storage.Profile, error)
115: ```
116: 
117: - `Discover()` reads `~/.config/gh/hosts.yml` directly (preferred, no shell-out) and falls back to parsing `gh auth status --show-token` output. Returns one entry per `(host, user)` pair authenticated.
118: - `EnrichAccount(acct)` calls `/user` and `/user/emails` via the token to fill `UserName` and `Email` (primary verified or noreply).
119: 
120: ### `gitswitch import gh` flow
121: 
122: 1. Run `Discover()`.
123: 2. For each account, prompt: `Create profile "<login>" for host "<hostname>"? [Y/n]`.
124: 3. On accept: write profile (`Source = "gh-import"`), store token in keychain at `gitswitch:<nickname>:<hostname>`, set `GHSynced = true`. Optionally prompt for SSH key path; if provided, write SSH config block.
125: 4. If two accounts share a login (same person on github.com and an enterprise host), offer to merge into one multi-host profile.
126: 
127: This subsumes the originally proposed OAuth device flow. Target audience already runs `gh auth login` for both `github.com` and any enterprise host — we ride on that.
128: 
129: ---
130: 
131: ## Keychain integration (new package: `internal/secrets/`)
132: 
133: Wrap [zalando/go-keyring](https://github.com/zalando/go-keyring) behind a thin interface so we can fall back when keyring isn't available:
134: 
135: ```go
136: type Store interface {
137:     Set(ref, token string) error
138:     Get(ref string) (string, error)
139:     Delete(ref string) error
140:     Available() bool
141: }
142: ```
143: 
144: Backends:
145: - `KeychainStore` (default) — macOS Keychain, libsecret on Linux, Credential Manager on Windows.
146: - `EncryptedFileStore` — `xchacha20poly1305` + `scrypt`-derived key from a passphrase, stored at `~/.config/gitswitch/secrets.enc` (mode 0600).
147: 
148: Selected at runtime via `Available()` probe. Override with `GITSWITCH_SECRETS_BACKEND=file`.
149: 
150: Token ref convention: `gitswitch:<nickname>:<hostname>`. Deterministic — `doctor` can walk profiles and probe each.
151: 
152: ---
153: 
154: ## gh CLI sync on switch
155: 
156: `internal/gh/sync.go`:
157: 
158: ```go
159: func SyncHost(hostname, ghUser, token string) error
160: ```
161: 
162: Writes / updates the host entry in `~/.config/gh/hosts.yml` using yaml.v3 in-place edits (preserves other host entries the user has). On every `applyProfile()`, walks the active profile's `Hosts` map and calls `SyncHost` per host. Best-effort: silent skip if `gh` isn't installed.
163: 
164: This keeps `gh` and `git push` in lockstep — the long-standing promise of gitswitch — across multiple hosts.
165: 
166: ---
167: 
168: ## SSH config management (new: `internal/ssh/config.go`)
169: 
170: Idempotent per-host block manager in `~/.ssh/config`:
171: 
172: ```
173: # gitswitch:github.corp.com
174: Host github.corp.com
175:     IdentityFile ~/.ssh/id_corp
176:     IdentitiesOnly yes
177: ```
178: 
179: Blocks keyed by `# gitswitch:<hostname>` marker comment. Add / update / remove via marker scan. Drops the global `core.sshCommand` approach in [internal/git/git.go:61](../internal/git/git.go) — per-host SSH config is correct and lets multi-remote pushes work without re-switching profiles between pushes.
180: 
181: ---
182: 
183: ## Switch flow (`cmd/gitswitch/main.go` `applyProfile`)
184: 
185: 1. `git config --global user.name / user.email / user.signingkey`.
186: 2. For each `(hostname, creds)` in active profile's `Hosts`:
187:    - SSH: upsert per-host block in `~/.ssh/config`.
188:    - Token: fetch from keychain → `gh.SyncHost(hostname, creds.GHUser, token)`.
189: 3. Update history + pins as today.
190: 
191: `SetSSHKey` in [internal/git/git.go](../internal/git/git.go) kept for migration path only; new flows don't call it.
192: 
193: ---
194: 
195: ## Per-repo pinning — extended to multi-pin
196: 
197: Today: `~/.config/gitswitch/history.json` stores `{repoKey: pinnedNickname}` (one profile per repo). Extend to:
198: 
199: ```go
200: type Pin struct {
201:     Default string            `yaml:"default,omitempty"` // fallback if no remote-specific pin
202:     Remotes map[string]string `yaml:"remotes,omitempty"` // remoteName -> nickname
203: }
204: ```
205: 
206: A repo with `origin = github.corp.com:...` and `fork = github.com:me/...` can pin:
207: 
208: ```yaml
209: my-repo-key:
210:   default: corp
211:   remotes:
212:     fork: personal
213: ```
214: 
215: Backwards compatibility:
216: - `gitswitch recommend` and the `record` shell hook continue to work — if only `default` is set, behavior matches v0.1.
217: - If `remotes` is populated, the shell prompt segment shows the resolved pin for the current branch's upstream remote.
218: 
219: Command surface:
220: - `gitswitch pin <profile> [--remote NAME]` — without `--remote`, pins `default`; with `--remote`, pins that remote only.
221: - `gitswitch unpin [--remote NAME]` — mirrors.
222: 
223: ---
224: 
225: ## `gitswitch status` — repo narrative diagnostic
226: 
227: New command. In a git repo, output a structured explanation:
228: 
229: ```
230: Repo: ~/code/corp-platform
231:   Branch: feature/x  →  origin/feature/x
232: 
233: Identity (global git config):
234:   user.name   = Abhiram Kanna
235:   user.email  = abhiram@corp.com
236:   signingkey  = 4A...                (gpg, present)
237:   Active profile: corp             (matches git config ✓)
238: 
239: Pins for this repo:
240:   default          → corp
241:   remote "fork"    → personal
242: 
243: Remotes:
244:   origin → git@github.corp.com:team/x.git
245:     Host: github.corp.com
246:     Profile that owns this host: corp
247:     SSH config block: present ✓        (~/.ssh/config: # gitswitch:github.corp.com)
248:     Identity file: ~/.ssh/id_corp    (exists, mode 600 ✓)
249:     gh auth: aksisonline-pp            (synced ✓, token valid as of 2026-05-28)
250:   fork → git@github.com:aksisonline/x.git
251:     Host: github.com
252:     Profile that owns this host: personal  (per fork-remote pin)
253:     SSH config block: present ✓
254:     Identity file: ~/.ssh/id_personal      (exists ✓)
255:     gh auth: aksisonline                   (synced ✓, token valid)
256: 
257: Recommendations: none.
258: ```
259: 
260: When any check fails, surface inline (`gh auth: stale token (401)`, `Identity file: missing`, `Pin references unknown profile "old-work"`). Composes existing helpers (`DetectActive`, `Recommend`) plus new keychain / gh probes.
261: 
262: `status` is per-repo; `doctor` (below) is global.
263: 
264: ---
265: 
266: ## Command surface
267: 
268: ### New
269: - `gitswitch import gh` — re-runnable gh-discovery flow.
270: - `gitswitch import git` — old `init` semantics (import current git config as new profile).
271: - `gitswitch host add <profile> <hostname> [--ssh-key] [--gh-user] [--token]`
272: - `gitswitch host remove <profile> <hostname>`
273: - `gitswitch host list <profile>`
274: - `gitswitch status` — per-repo narrative (above).
275: - `gitswitch doctor` — global health (every profile × every host: keychain reachable? token valid? SSH file present? gh hosts.yml in sync?).
276: 
277: ### Renamed / repurposed
278: - `gitswitch init` — ensure `~/.config/gitswitch/config.yaml` exists, import current git config as first profile, run `gh import` if `gh` is logged in. Old import-only behavior moves to `gitswitch import git`.
279: 
280: ### Extended
281: - `gitswitch add` — gains repeatable `--host hostname[=ssh_key_path]` flag for multi-host creation.
282: - `gitswitch pin` / `unpin` — gain `--remote NAME` flag for multi-pin.
283: 
284: ### Unchanged
285: `switch`, `list`, `remove`, `current`, `record`, `recommend`, `install`, `uninstall`, `version`, `upgrade`, `hook-check`, `claude`, `pacman`.
286: 
287: ---
288: 
289: ## Verification plan
290: 
291: ### Unit
292: - `internal/storage` migration (v1 JSON → v2 YAML, idempotent on second run, backup created).
293: - `internal/secrets` round-trip both backends.
294: - `internal/ssh/config` idempotent block management.
295: - `internal/gh/import` parses `gh auth status` and `hosts.yml`.
296: 
297: ### Integration
298: Build binary into tmp HOME and:
299: - `gitswitch init` produces `config.yaml` + imports git config.
300: - `gitswitch import gh` (with mocked `gh` state) creates expected profiles.
301: - `gitswitch pin corp --remote origin` writes the expected pin structure.
302: - `gitswitch status` in a fake repo renders all sections.
303: 
304: ### End-to-end (Vaishnavi scenario reproduction)
305: Two profiles: `personal` (github.com) and `corp` (multi-host: github.com + github.corp.com).
306: 
307: In a repo with `origin = git@github.corp.com:team/x.git` and `fork = git@github.com:me/x.git`:
308: 1. `gitswitch pin corp`
309: 2. `gitswitch pin personal --remote fork`
310: 3. `git push origin` and `git push fork` both succeed via per-host SSH config.
311: 4. `gh repo view` works against both hosts after switching.
312: 5. `gitswitch status` reflects everything correctly.
313: 
314: ### Diagnostics
315: - Revoke a PAT externally → `gitswitch doctor` flags the expired token for the right host; `gitswitch status` shows it inline against the offending remote.
316: 
317: ### Migration safety
318: Smoke test: install v0.1.x, create two profiles, upgrade binary, run `gitswitch list` → profiles intact, `profiles.json.v1.bak` present, `config.yaml` has hosts map populated with `github.com`.
319: 
320: ---
321: 
322: ## Critical files
323: 
324: - [internal/storage/store.go](../internal/storage/store.go) — Profile + Pin structs, YAML migration, load/save.
325: - [internal/git/git.go](../internal/git/git.go) — drop global `core.sshCommand` write path.
326: - [internal/history/history.go](../internal/history/history.go) — extend Pin to multi-remote shape; preserve existing API surface.
327: - [cmd/gitswitch/main.go](../cmd/gitswitch/main.go) — command wiring; rewrite `applyProfile`; new commands (`status`, `import gh`, `import git`, `host *`, `doctor`).
328: - **New**: `internal/secrets/` (keychain + encrypted file), `internal/gh/sync.go`, `internal/gh/import.go`, `internal/ssh/config.go`.
329: - [docs/diagnostic-design.md](diagnostic-design.md) — extend with multi-host checks.
330: - [README.md](../README.md) — refresh positioning to "who you are + which GitHub" on ship.
331: 
332: ## Companion plans
333: 
334: - [roadmap-v0.2.0-isolation.md](roadmap-v0.2.0-isolation.md) — env-scoped sessions so two terminals (or two AI agents) can hold two different active identities at the same time without touching global state. Shares the `internal/secrets/` keychain helper with this plan.
335: - [roadmap-v0.2.0-tui.md](roadmap-v0.2.0-tui.md) — TUI refresh on Bubble Tea v2 / Bubbles v2 / Lipgloss v2. Full CLI feature parity (every command reachable via keybinds + `:` palette), arcade aesthetic preserved as a theme layer. Pure presentation over the data models in this plan + the session plan.
336: 
337: ## Out of scope — see [roadmap-v0.3.0.md](roadmap-v0.3.0.md)
338: 
339: In-repo YAML config, team preset distribution + signing, own OAuth flow, remote retargeting, host auto-detection.

(End of file - total 339 lines)
</content>
```

---

## Assistant (Build · Kimi K2.7 Code · 52.3s)

_Thinking:_

Now I have a good understanding of the existing v0.2.0 Go plan. I need to:
1. Revamp roadmap-full-timeline.md to remove Rust, add OAuth to v1.0, add gs alias
2. Revamp roadmap-v0.3.0.md to remove OAuth
3. Create a new versioning/roadmap doc that fits v0.1.22 → v0.2.x → v0.3.x → v1.0

The user said "i didnt understand, explain me my own original plan and this new one" and then "so basically we are on version 0.1.22 , idk where to go right now... theres so many other things like PAT Token management and session isolations for multi agent orchestration that i want, so I need to come up with a versioning plan properly"

So I need to help create a proper versioning plan. PAT Token management and session isolation are already in v0.2.0 plan. The user wants these. I should organize the plan clearly.

Let me think about a sensible versioning plan from v0.1.22:

**v0.2.0** - The multi-host + session + TUI refresh release (big breaking-ish change with config migration)
- Multi-host identity
- PAT/keychain storage
- Session isolation for AI agents
- TUI refresh
- `gitswitch status`, `gitswitch doctor`
- `gitswitch import gh`, `gitswitch import git`
- `gitswitch host *`
- Multi-pin
- Pre-commit safety net
- Shell tab completion

**v0.3.0** - Team / repo-level config release
- In-repo `.gitswitch.yaml`
- Team preset distribution (with signing)
- Host auto-detection
- Remote retargeting
- Remote-URL-aware pinning
- Profile export/import
- `gitswitch run`
- git includeIf support

**v1.0** - Stable production release
- Direct GitHub OAuth device flow (`gitswitch login`)
- `gs` short command alias
- HTTPS credential helper
- Stability guarantees (semver, config schema)
- Official distribution channels

This makes sense. The user wants PAT mgmt and session isolation - these are v0.2.0. OAuth is v1.0.

Now let me write the updated docs. I'll overwrite roadmap-full-timeline.md with the new plan. Since I've already read it, I can use write. But wait, the write tool requires reading first if file exists. I've read it, so I can write.

Actually, I should use write to replace the whole file, or edit large sections. Given the extent of changes, write is cleaner.

Let me also write a new doc `roadmap-versioning-plan.md` or similar that explains the plan from v0.1.22 onwards. Or maybe update `roadmap-v0.3.0.md` and create a new `roadmap-v1.0.md`.

Actually, the user said "finish cancelling the plan, and then let's go and write new docs with new plans for these." So I should write new docs. Let me create:
1. Updated `roadmap-full-timeline.md` (revamped)
2. Updated `roadmap-v0.3.0.md` (OAuth removed, sequencing updated)
3. New `roadmap-v1.0.md` (detailed v1.0 plan: OAuth device flow, gs alias, HTTPS credential helper)
4. Maybe new `roadmap-versioning-plan.md` that ties it all together from v0.1.22

Actually, `roadmap-full-timeline.md` already serves as the tying-together doc. I can update it comprehensively and create `roadmap-v1.0.md` for details.

Let me start with roadmap-full-timeline.md.

**Tool: write**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/roadmap-full-timeline.md",
  "content": "# gitswitch — Full Product Timeline\n\n_Last updated: June 2026_  \n_Status: Rust migration **cancelled**. v1.0 will ship on the Go codebase._\n\n---\n\n## v0.1.x — Current (Go, shipping now)\n\nCurrent release: **v0.1.22**. The foundation. Everything in this version works today.\n\n**Identity management**\n- Named profiles: name, email, SSH key, GPG/signing key, gh_user (single host: github.com)\n- Switch global `~/.gitconfig` identity in one command\n- `gitswitch current` — see who you are right now\n\n**Shell integration**\n- Shell hook fires on `cd`, logs repo + identity usage in background\n- Learns per-repo patterns; suggests the right profile after 3+ uses\n- Shell prompt segment: shows active profile (or warning if config doesn't match any profile)\n- Supports zsh, bash, fish, Oh My Zsh\n\n**TUI**\n- Full-screen terminal UI launched by `gitswitch` with no arguments\n- Three themes: arcade (loud, retro, animated), modern (clean), minimal (monochrome, SSH-safe)\n- Pacman easter egg: hidden game using profile list as the level\n\n**Per-repo pinning**\n- `gitswitch pin <profile>` — this repo always uses this profile\n- Automatic prompt when entering a pinned repo with wrong identity active\n- `gitswitch unpin` to remove\n\n**gh CLI sync**\n- Writes the correct account to `~/.config/gh/hosts.yml` on every switch\n- git and `gh` can never drift apart\n\n**Distribution**\n- Single static binary, no runtime dependencies\n- `curl | sh` one-liner installer\n- Homebrew tap\n- Storage: `~/.gitswitch/profiles.json`\n\n**Utility commands**\n- `gitswitch doctor` — health check on all profiles\n- `gitswitch upgrade` — self-update from GitHub releases\n- `gitswitch hook-check` — verify shell hook is installed and current\n- `gitswitch install` / `gitswitch uninstall`\n- `gitswitch version`\n\n---\n\n## v0.2.0 — Identity + isolation milestone (Go)\n\nThe big infrastructure release. Three things ship together: multi-host identity, secure PAT storage, and session isolation. Config format migrates from JSON to YAML.\n\n### Multi-host identity\n\nOne profile covers multiple GitHub hosts with different credentials per host.\n\n**Profile model change**\n```yaml\nnickname: corp\nuser_name: Abhiram Kanna\nemail: abhiram@corp.com\nsign_key: 4A1B2C3D\nhosts:\n  github.com:\n    ssh_key: ~/.ssh/id_personal\n    gh_user: aksisonline\n    token_ref: gitswitch:corp:github.com\n  github.corp.com:\n    ssh_key: ~/.ssh/id_corp\n    gh_user: aksisonline-pp\n    token_ref: gitswitch:corp:github.corp.com\n```\n\n**What this enables**\n- One `gitswitch switch corp` sets SSH config, gh auth, and git identity for all hosts simultaneously\n- Works correctly with repos that have multiple remotes pointing to different hosts\n- Multi-pin: `gitswitch pin corp --remote origin`, `gitswitch pin personal --remote fork`\n\n**New storage**\n- Migration from `profiles.json` → `~/.config/gitswitch/config.yaml` (single file: version + prefs + profiles + pins)\n- Old file renamed `profiles.json.v1.bak` on first load\n- PAT tokens stored in OS keychain (macOS Keychain, Windows Credential Manager, Linux Secret Service) — never in plaintext\n- Encrypted-file fallback for headless/CI environments\n\n**New commands**\n- `gitswitch import gh` — discover all accounts from `gh auth status`, create profiles automatically\n- `gitswitch import git` — import current git config as a new profile (old `init` behavior)\n- `gitswitch host add <profile> <hostname>` — add a host entry to an existing profile\n- `gitswitch host remove <profile> <hostname>`\n- `gitswitch host list <profile>`\n\n### Session isolation\n\nEach terminal window holds its own identity via environment variables — parallel sessions, independent identities. Designed for AI agents and multi-repo terminal workflows.\n\n**How it works**\n- `gitswitch shell [profile]` — opens a subshell with identity locked in env vars\n- Inside the session: git reads `GIT_AUTHOR_NAME`, `GIT_AUTHOR_EMAIL`, etc. in preference to `~/.gitconfig`\n- `GH_TOKEN` / `GH_ENTERPRISE_TOKEN` set so `gh` uses the right account\n- `GITSWITCH_SESSION_ID`, `GITSWITCH_SESSION_PROFILE` for introspection\n- Exit the shell → env vars disappear → no cleanup needed\n- Can still `gitswitch switch` within a session to change identity for that session only\n\n**Why this matters for AI agents**\n- Two Claude Code agents running in parallel, each inside their own `gitswitch shell`, cannot interfere with each other's identity\n- Every commit from every agent is correctly attributed regardless of what the developer is doing globally\n- `gitswitch env <profile> --format json` for headless/agent consumption\n\n### Diagnostics and safety\n\n**`gitswitch status`** (new) — per-repo narrative diagnostic\n- What identity is active (global or session)\n- Every remote: host, SSH config block present?, SSH key exists + permissions correct?, gh auth valid + non-expired?\n- Pin state, recommendations\n- Designed to answer \"why is my push failing?\"\n\n**Pre-commit safety net** (auto-installed by `gitswitch install`)\n- Checks if current identity matches the pinned profile for the repo\n- Warns before the wrong-identity commit lands: \"you're about to commit as personal@gmail.com — this repo is pinned to work\"\n- User can override with `git commit --no-verify`\n- Can opt out via config\n\n**Progressive CLI hints**\n- After `gitswitch switch` → \"tip: run `gitswitch pin` to make this automatic\"\n- After `gitswitch pin` → \"run `gitswitch status` to see the full picture\"\n- After `gitswitch install` → \"shell hook installed — gitswitch will suggest your usual identity when you enter a repo\"\n\n### TUI refresh\n\n- 5 screens: Profiles, Hosts, Repos, Session, Doctor\n- Command palette (`:` key) for power users\n- Wizards for profile creation and host addition\n- New keybindings for session and host management\n- Shell tab completion ships in v0.2.0\n\n---\n\n## Rust migration — CANCELLED\n\n> **CANCELLED — June 2026**\n>\n> The Rust rewrite has been cancelled. The risk of binary-compatible feature parity + migration outweighed the marginal gains for a tool that already works.\n> gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.\n>\n> Archived docs are in [`rust-cancelled/`](rust-cancelled/).\n\n---\n\n## v0.3.0 — Repo + team features (Go)\n\nItems deliberately deferred from v0.2.0 planning. These make gitswitch work better across teams and machines.\n\n**In-repo `.gitswitch.yaml`**\n- Per-repo config file optionally checked into the repo root\n- Declares required profiles, identity policies (signing required, email suffix enforcement)\n- Loaded by the shell hook automatically\n- Precedence: session env > local pin > repo YAML > global switch\n- Open questions: team-overridable vs locked fields, worktree interaction — needs design\n\n**Team preset distribution**\n- `gitswitch apply <file-or-url>` — bootstraps all profiles from a shared config in one command\n- Signed with minisign so an attacker can't push a malicious preset\n- Intended for engineering teams onboarding new developers\n\n**Host auto-detection**\n- Walk `git remote -v` on `cd` into a repo\n- Auto-recommend the profile whose host map covers all remotes\n- Reduces manual pinning for repos with clear host signals\n\n**Remote retargeting**\n- `gitswitch remote retarget` — rewrites `origin` (and others) to a new host while preserving the path\n- Useful for org migrations between GitHub hosting models\n\n**Remote-URL-aware routing**\n- Pin by remote URL pattern instead of directory path\n- `github.com/mycompany/* → work`, `github.com/aksisonline/* → personal`\n- More robust than directory-based pinning for developers who don't organize repos conventionally\n\n**git includeIf support (optional, not default)**\n- Fires identity at git-command time via git's native config loading, not on `cd`\n- No shell modification needed — alternative for users who want zero shell hooks\n- gitswitch's shell hook remains the default; includeIf is an opt-in mode\n\n**Other v0.3.0 items**\n- Ephemeral one-shot run: `gitswitch run <profile> -- <cmd>` (run one command as a different identity without switching)\n- Profile export/import (JSON backup and restore for new machine setup)\n- Windows native support improvements (beyond WSL)\n\n**Likely sequencing**\n1. In-repo `.gitswitch.yaml`\n2. Host auto-detection\n3. Team preset apply *without* signing (gated behind `--unsafe-trust` flag)\n4. Team preset signing\n5. Remote retargeting\n6. Remote-URL-aware routing\n\nOrder may shift based on which use cases surface from 0.2.0 adopters.\n\n---\n\n## v1.0 — Production milestone (Go, stable)\n\nThe 1.0 release. Stable CLI API with semver commitment. Enterprise-ready. No language rewrite — just hardened, complete auth coverage.\n\n**Direct GitHub OAuth (device flow)**\n- `gitswitch login [--host HOST] [--client-id ID]`\n- Device flow against github.com; no local HTTP server needed\n- Scopes: `repo`, `read:user`, `user:email`, `gist`, `workflow`\n- `client_id` baked into the binary; user-supplied `client_id` supported for GHES\n- Tokens stored in OS keychain\n- Required for users without `gh` CLI and for GHES environments where `gh` won't authenticate\n\n**`gs` short command**\n- `gs` is equivalent to `gitswitch`\n- Implemented as a shim binary/symlink (or argv[0] detection in the same binary)\n- Works immediately after install, in scripts, and across shells\n\n**HTTPS credential helper**\n- `gitswitch` acts as a proper git `credential.helper`\n- Set with: `git config --global credential.helper gitswitch`\n- When git needs credentials for any HTTPS remote, it asks gitswitch — gitswitch returns the right PAT for the active profile\n- Closes the gap for users who clone via HTTPS (the GitHub default) rather than SSH\n- Means gitswitch handles the full auth stack: SSH (via `~/.ssh/config`), HTTPS (via credential helper), gh CLI (via `hosts.yml`)\n\n**Stability guarantees**\n- Semver commitment: no breaking CLI changes without a major version bump\n- Stable config format: `~/.config/gitswitch/config.yaml` schema versioned\n- Stable keychain key convention: `gitswitch:<nickname>:<hostname>`\n\n**Distribution at v1.0**\n- curl one-liner (primary)\n- Homebrew (official tap)\n- AUR (Arch Linux)\n- Windows: native installer (NSIS or WiX)\n- gitswitch.dev live with full docs, blog, and privacy policy\n\n---\n\n## What does NOT change across any version\n\n- Every feature, including all enterprise toolsets, stays free forever\n- No cloud component, no telemetry, no data leaves the machine\n- No feature flags, no paywalls, no paid themes\n- The arcade personality and easter eggs are never gatekept\n- Open source, inspectable at any time\n"
}
```

**Output:**
```
Wrote file successfully.
```

---

## Assistant (Build · Kimi K2.7 Code · 33.4s)

_Thinking:_

Now update roadmap-v0.3.0.md to remove OAuth and update sequencing. I'll use write to replace the whole file since I've read it.

**Tool: write**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/roadmap-v0.3.0.md",
  "content": "# gitswitch v0.3.0 — Roadmap\n\nItems deferred from [v0.2.0 planning](roadmap-v0.2.0.md). v0.2.0 establishes multi-host identity, keychain-backed PATs, gh-driven profile import, multi-pin, session isolation, and `gitswitch status`. v0.3.0 builds on that foundation to address team-scale distribution, in-repo config, and repo-host migrations.\n\n> **Note:** The original v0.3.0 plan included a direct GitHub OAuth device flow. That has been moved to the [v1.0 launch](roadmap-v1.0.md) because it is part of the \"gitswitch owns auth end-to-end\" production milestone, not a mid-beta fallback.\n\n---\n\n## 1. In-repo YAML config (`.gitswitch.yaml`)\n\nA per-repo config file checked into the repo root. Declares which profile(s) should be used, which host(s) remotes resolve to, and identity policies the repo enforces.\n\n### Why\n\nToday (and in 0.2.0), pins live in `~/.config/gitswitch/config.yaml` keyed by absolute repo path. That doesn't survive:\n- A teammate cloning the repo to a different path.\n- A CI runner.\n- Switching machines.\n\nIn-repo config makes \"this repo wants the `corp` profile\" a property of the repo, not the local machine.\n\n### Open questions\n\n- **Precedence**: in-repo config vs. local pin in `~/.config/gitswitch/config.yaml` — does in-repo win, or does local pin override (escape hatch)?\n- **Policy enforcement**: should `.gitswitch.yaml` be able to *block* commits that don't match the policy (e.g. via a `prepare-commit-msg` hook)? Or only nudge?\n- **Locked vs. overridable fields**: a team config probably wants to lock `email_suffix` but leave `ssh_key` path overridable per-dev.\n- **Worktrees**: same `.gitswitch.yaml` applies to all worktrees of a repo? Or per-worktree opt-out?\n- **Schema versioning**: how does an older `gitswitch` binary handle a `.gitswitch.yaml` from a newer schema version?\n\n### Sketch\n\n```yaml\n# .gitswitch.yaml (in repo root)\nversion: 1\nrequire:\n  email_suffix: \"@corp.com\"\n  signing: required\nhosts:\n  origin: github.corp.com\n  fork:   github.com           # optional, validated against actual remote\nsuggested_profile: corp\n```\n\n---\n\n## 2. Team preset distribution\n\nA team-shared bootstrap config — distributed via intranet URL, internal git repo, or shared drive — that onboards a new dev in one command:\n\n```\ngitswitch apply https://wiki.corp.com/gitswitch-platform-team.yaml\ngitswitch apply ./team-preset.yaml\n```\n\n### Components\n\n- **Schema** for team presets (distinct from per-repo `.gitswitch.yaml` — this is a *template* that produces profiles, not policy applied to a repo).\n- **Fetch + cache** for HTTPS sources. Local file path also supported.\n- **Signing & verification** — so an attacker can't push a malicious preset to wiki / internal repo and steal credentials.\n  - Candidate formats: [minisign](https://jedisct1.github.io/minisign/), [cosign](https://github.com/sigstore/cosign).\n  - Trust roots: a team config declares its signing key fingerprint; first-apply prompts the user to trust it (TOFU); subsequent applies verify against the trusted key.\n  - Revocation: list of revoked keys fetched from a known endpoint.\n- **Apply semantics**: `apply` creates profile skeletons (host config, gh_user pattern, ssh_key naming hint) and prompts only for personal bits (the dev's actual SSH key path, their actual gh login). Pre-existing profiles aren't clobbered — only new ones are added.\n\n### Why this is hard enough to defer\n\nThe signing story is the load-bearing piece. Shipping `apply` without signing means corp security teams will reject it on day one. Building signing properly (key management, trust onboarding, revocation, UX for \"this preset is no longer trusted\") is a project on its own.\n\n---\n\n## 3. Host auto-detection\n\nOn `cd` into a repo (via the existing shell hook), walk `git remote -v` and auto-recommend the profile whose `Hosts` map covers all remotes.\n\n### Why\n\nReduces manual pinning. If a repo's only remote is `github.corp.com:team/x.git` and the user has exactly one profile with that host configured, gitswitch can confidently suggest that profile without the user having to `gitswitch pin` it.\n\n### Sketch\n\nExtend [internal/history/history.go](../internal/history/history.go) `Recommend`:\n1. Existing path: usage frequency ≥ 60% → recommend.\n2. New path: if no usage history yet AND exactly one profile covers every host in `git remote -v` → recommend that profile.\n3. If multiple profiles cover all hosts → no recommendation (ambiguous; defer to user).\n\n### Open questions\n\n- Should auto-detection auto-switch silently, or always nudge?\n- Interaction with multi-pin: should auto-detection emit a multi-pin suggestion when different remotes are best served by different profiles?\n\n---\n\n## 4. Remote retargeting\n\n`gitswitch remote retarget` — repo-host migration helper.\n\n### Why\n\nWhen an org migrates between GitHub hosting models — An org's `github.corp.com` → `github.com/neworg` is the canonical case — every developer has hundreds of local clones with `origin` pointing at the old host. Manual `git remote set-url` per repo is painful.\n\n### Sketch\n\n```\ngitswitch remote retarget --from github.corp.com --to github.com/neworg\n```\n\nIn each known repo (or the current repo):\n1. Walk `git remote -v`.\n2. For any remote matching `--from`, rewrite to the equivalent path under `--to`.\n3. Confirm before each rewrite (or `--yes` to skip).\n4. Optionally switch the active pin to the profile that owns the new host.\n\n### Open questions\n\n- Should this walk *all* repos under a configured root (e.g. `~/code`) or only the current repo?\n- Path translation: `github.corp.com:teamA/repo` → `github.com/neworg:teamA/repo`? Or do paths get rewritten too (`neworg-teamA/repo`)?\n\n---\n\n## 5. Remote-URL-aware routing\n\nPin by remote URL pattern instead of directory path.\n\n```\ngithub.com/mycompany/* → work\ngithub.com/aksisonline/* → personal\n```\n\n### Why\n\nMore robust than directory-based pinning for developers who don't organize repos conventionally. Survives re-cloning to a different path.\n\n### Sketch\n\nAdd a `routes:` section to `~/.config/gitswitch/config.yaml`:\n\n```yaml\nroutes:\n  - pattern: \"github.com/mycompany/*\"\n    profile: work\n  - pattern: \"github.com/aksisonline/*\"\n    profile: personal\n```\n\nResolution precedence:\n1. Session env\n2. Local pin (`pins[repoPath]`)\n3. Remote-URL route\n4. In-repo `.gitswitch.yaml` suggested profile\n5. Global active profile\n\n---\n\n## Other v0.3.0 items\n\n- **Ephemeral one-shot run**: `gitswitch run <profile> -- <cmd>` — run one command as a different identity without switching.\n- **Profile export/import**: JSON backup and restore for new machine setup.\n- **git includeIf support (optional, not default)**: fires identity at git-command time via git's native config loading, no shell modification needed.\n- **Windows native support improvements** (beyond WSL).\n\n---\n\n## Likely sequencing for v0.3.0\n\n1. In-repo `.gitswitch.yaml` (foundation: in-repo config schema becomes the substrate for team presets).\n2. Host auto-detection (small, high-impact UX win; piggybacks on 0.2.0 shell hook).\n3. Team preset apply *without* signing (gated behind `--unsafe-trust` flag; collects feedback on schema before locking it in).\n4. Team preset signing (the big one).\n5. Remote retargeting.\n6. Remote-URL-aware routing.\n\nOrder may shift based on which use cases surface from 0.2.0 adopters.\n"
}
```

**Output:**
```
Wrote file successfully.
```

---

## Assistant (Build · Kimi K2.7 Code · 31.6s)

_Thinking:_

Now create the v1.0 roadmap doc with detailed OAuth and gs alias plans.

**Tool: write**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/roadmap-v1.0.md",
  "content": "# gitswitch v1.0 — Production Launch Roadmap\n\nThe 1.0 release marks the production-stable Go binary. There is no language rewrite. The focus is completing the authentication story (so gitswitch no longer depends on `gh` CLI for login), adding a short command alias, and making stability promises.\n\n---\n\n## 1. Direct GitHub OAuth (device flow)\n\n### Goal\n\nLet users authenticate with GitHub directly through gitswitch, without requiring the `gh` CLI. This unblocks:\n- Users who don't have `gh` installed.\n- GitHub Enterprise Server (GHES) environments where `gh` won't authenticate.\n- Headless setups where a browser-driven `gh auth login` is impractical.\n\n### Flow\n\nDevice flow is chosen because it is terminal-native, SSH-friendly, and requires no local HTTP server.\n\n```\n$ gitswitch login\nOpen https://github.com/login/device and enter code: ABCD-EFGH\nWaiting for authorization...\n✓ Logged in as aksisonline (github.com)\nCreated profile \"aksisonline\" with email abhiram@personal.dev\nToken stored in keychain.\n```\n\n### Command surface\n\n```\ngitswitch login [--host HOST] [--client-id ID] [--profile NAME]\n```\n\n| Flag | Default | Purpose |\n|---|---|---|\n| `--host` | `github.com` | GitHub host. Use `github.corp.com` for GHES. |\n| `--client-id` | built-in gitswitch OAuth app | Override for GHES orgs that register their own OAuth app. |\n| `--profile` | auto from login | Attach the new host to an existing profile instead of creating a new one. |\n\n### Scopes\n\n`repo`, `read:user`, `user:email`, `gist`, `workflow` — matches `gh auth` defaults so the token is a drop-in replacement.\n\n### Token storage\n\n- Tokens land in the OS keychain under the existing convention: `gitswitch:<nickname>:<hostname>`.\n- The profile record stores only the `token_ref`, never the token itself.\n- On `gitswitch switch <profile>`, the token is read from keychain and synced to `~/.config/gh/hosts.yml` (if `gh` is installed), just like imported tokens.\n\n### OAuth app ownership\n\nFor github.com: register a single gitswitch project OAuth app. The `client_id` is baked into the binary. Device flow does not require a `client_secret`.\n\nFor GHES: each organization registers their own OAuth app and passes `--client-id`.\n\n### Refresh tokens\n\nGitHub device flow does not issue refresh tokens by default. When a token expires, `gitswitch doctor` flags it and `gitswitch login --profile <name>` re-authenticates the same profile. We do not build a refresh-token workaround.\n\n### Implementation notes\n\n- HTTP client: reuse the existing Go HTTP stack (currently used by upgrade/version checks).\n- Poll interval and expiry follow GitHub's device flow spec (typically 5-second poll, 15-minute expiry).\n- No async runtime, no background daemon.\n\n---\n\n## 2. `gs` short command alias\n\n### Goal\n\nGive users a fast alias for `gitswitch`. `gs switch work`, `gs status`, `gs shell`.\n\n### Implementation\n\n**Option selected: shim binary / argv[0] detection.**\n\nThe release ships a `gs` executable that is either:\n- A second small binary that `exec`s `gitswitch` and forwards all args, or\n- A symlink to the same `gitswitch` binary, where the binary detects `os.Args[0] == \"gs\"` and behaves identically.\n\nThe symlink / argv[0] approach is preferred because it keeps the install footprint tiny and guarantees `gs` and `gitswitch` can never diverge.\n\n```go\nfunc main() {\n    if filepath.Base(os.Args[0]) == \"gs\" {\n        // same behavior, different invocation name\n    }\n    // existing gitswitch main\n}\n```\n\n### Install behavior\n\n- The `curl | sh` installer creates both `gitswitch` and `gs` in the install directory.\n- `gitswitch install` does **not** need to touch shell rc files for this; both binaries are on PATH.\n- No shell alias needed, so it works in scripts, subshells, and fresh terminals immediately.\n\n### Edge cases\n\n- If a user already has a `gs` binary on PATH (e.g. Ghostscript), the installer warns and offers `--no-gs`.\n- Package managers (Homebrew, AUR, etc.) ship both binaries.\n- `gs --version` prints `gitswitch v1.0.0 (gs)`.\n\n---\n\n## 3. HTTPS credential helper\n\n### Goal\n\nMake gitswitch work for HTTPS clones, not just SSH. Today HTTPS users are invisible to gitswitch — GitHub shows HTTPS by default, but gitswitch only manages SSH keys and `gh` auth.\n\n### How it works\n\n```sh\ngit config --global credential.helper gitswitch\n```\n\nWhen git needs credentials for an HTTPS remote, it invokes:\n\n```\ngitswitch credential fill\n```\n\ngitswitch inspects the URL host, looks up the active profile's host map, fetches the PAT from the keychain, and returns:\n\n```\nprotocol=https\nhost=github.com\nusername=aksisonline\npassword=<PAT>\n```\n\n### Commands\n\n- `gitswitch credential fill` — called by git.\n- `gitswitch credential approve` — no-op (we don't cache separately; keychain is the cache).\n- `gitswitch credential reject` — optionally delete the token from keychain and prompt re-login.\n\n### Integration with `gitswitch switch`\n\nNo extra work. The credential helper is stateless: it always reads the active profile + keychain at request time. Switching profiles changes which PAT is returned.\n\n### Security\n\n- The helper only responds over stdin/stdout with the git credential protocol.\n- It refuses to return a token unless the requesting host matches a host in the active profile.\n- It never logs tokens.\n\n---\n\n## 4. Stability guarantees\n\n### Semver\n\nv1.0.0 begins semver. Breaking CLI changes require a major version bump.\n\n### Config format\n\n- `~/.config/gitswitch/config.yaml` schema versioned.\n- Forward-compatible loaders: unknown fields preserved.\n- Migrations between schema versions are automatic and backed up.\n\n### Keychain convention\n\nToken references remain `gitswitch:<nickname>:<hostname>`.\n\n### Deprecation policy\n\n- Commands or flags may be deprecated in a minor release.\n- Deprecated items continue to work through the next major release.\n- Removal only at major version bumps.\n\n---\n\n## 5. Distribution\n\n- curl one-liner (primary)\n- Homebrew official tap\n- AUR (Arch Linux)\n- Windows native installer (NSIS or WiX)\n- gitswitch.dev live with docs, blog, privacy policy\n\n---\n\n## 6. Launch criteria\n\n- [ ] OAuth device flow works against github.com.\n- [ ] OAuth device flow works against a GHES instance with `--client-id`.\n- [ ] `gs` alias ships and works on macOS, Linux, Windows.\n- [ ] HTTPS credential helper works with HTTPS remotes.\n- [ ] Existing SSH + `gh` auth flows remain unchanged.\n- [ ] Config schema versioned and migration tested.\n- [ ] All v0.2.0 and v0.3.0 features stable and documented.\n"
}
```

**Output:**
```
Wrote file successfully.
```

---

## Assistant (Build · Kimi K2.7 Code · 8.4s)

_Thinking:_

Now I should update the todo list and also create a clear summary doc or maybe a `roadmap-versioning-plan.md` that explains the plan from v0.1.22. Actually, the user asked for versioning plan properly. I think the updated `roadmap-full-timeline.md` plus `roadmap-v1.0.md` covers it. But maybe I should also create a concise `VERSIONING.md` or update `project-memory.md` to reflect the new plan.

Actually, let me create a new doc `roadmap-versioning-plan.md` that is a concise decision record for the versioning plan from v0.1.22 onwards. This will help the user understand where they are and where they're going.

Let me also update project-memory.md's future roadmap table to reflect the new plan (OAuth moved to v1.0, HTTPS credential helper moved to v1.0, etc.).

First, let me update project-memory.md's future roadmap table and layer model since some things changed.

**Tool: todowrite**

**Input:**
```json
{
  "todos": [
    {
      "content": "Add cancellation markers to Rust docs and rename rust/ folder",
      "status": "completed",
      "priority": "high"
    },
    {
      "content": "Revamp roadmap-full-timeline.md with new v1 plan (Go stable, OAuth device flow, gs alias)",
      "status": "completed",
      "priority": "high"
    },
    {
      "content": "Revamp roadmap-v0.3.0.md to remove OAuth and reflect new sequencing",
      "status": "completed",
      "priority": "high"
    },
    {
      "content": "Check current project version and existing features to ground the new versioning plan",
      "status": "completed",
      "priority": "high"
    },
    {
      "content": "Draft new versioning/roadmap docs for v0.2.x, v0.3.x, v1.0 with PAT mgmt, session isolation, OAuth, gs",
      "status": "in_progress",
      "priority": "high"
    }
  ]
}
```

**Output:**
```
[
  {
    "content": "Add cancellation markers to Rust docs and rename rust/ folder",
    "status": "completed",
    "priority": "high"
  },
  {
    "content": "Revamp roadmap-full-timeline.md with new v1 plan (Go stable, OAuth device flow, gs alias)",
    "status": "completed",
    "priority": "high"
  },
  {
    "content": "Revamp roadmap-v0.3.0.md to remove OAuth and reflect new sequencing",
    "status": "completed",
    "priority": "high"
  },
  {
    "content": "Check current project version and existing features to ground the new versioning plan",
    "status": "completed",
    "priority": "high"
  },
  {
    "content": "Draft new versioning/roadmap docs for v0.2.x, v0.3.x, v1.0 with PAT mgmt, session isolation, OAuth, gs",
    "status": "in_progress",
    "priority": "high"
  }
]
```

---

## Assistant (Build · Kimi K2.7 Code · 10.5s)

_Thinking:_

Let me update project-memory.md's future roadmap table to reflect the new plan. I'll edit lines 124-132.

**Tool: edit**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/project-memory.md",
  "oldString": "| Feature | Target | Notes |\n|---|---|---|\n| HTTPS credential helper | v0.3.0+ | `git config credential.helper gitswitch` |\n| Remote-URL-aware pinning | v0.3.0+ | Pattern-based, not directory-based |\n| Ephemeral one-shot run | v0.3.0+ | `gitswitch run <profile> -- <cmd>` (gpx does this) |\n| Shell tab completion | v0.2.0 | Should ship with v0.2.0 |\n| Profile export/import | v0.3.0+ | JSON backup for new machine onboarding |\n| git includeIf (optional) | future | Zero-shell-modification alternative mode |\n| Windows native support | low priority | WSL covers most cases |",
  "newString": "| Feature | Target | Notes |\n|---|---|---|\n| Multi-host identity | v0.2.0 | One profile, multiple GitHub hosts |\n| PAT keychain storage | v0.2.0 | `gitswitch:<nickname>:<hostname>` |\n| Session isolation | v0.2.0 | `gitswitch shell`, env-var scoped |\n| `gitswitch status` / `doctor` | v0.2.0 | Per-repo and global diagnostics |\n| Shell tab completion | v0.2.0 | Ships with v0.2.0 |\n| In-repo `.gitswitch.yaml` | v0.3.0 | Repo-level identity policy |\n| Team preset distribution | v0.3.0 | `gitswitch apply <url>` with signing |\n| Host auto-detection | v0.3.0 | Recommend profile from `git remote -v` |\n| Remote retargeting | v0.3.0 | Rewrite remotes during host migration |\n| Remote-URL-aware pinning | v0.3.0 | Pattern-based, not directory-based |\n| Ephemeral one-shot run | v0.3.0 | `gitswitch run <profile> -- <cmd>` |\n| Profile export/import | v0.3.0 | JSON backup for new machine onboarding |\n| Direct GitHub OAuth | v1.0 | Device flow, no `gh` dependency |\n| `gs` short command | v1.0 | `gs` == `gitswitch` |\n| HTTPS credential helper | v1.0 | `git config credential.helper gitswitch` |\n| git includeIf (optional) | future | Zero-shell-modification alternative mode |\n| Windows native support | v0.3.0+ | Beyond WSL |"
}
```

**Output:**
```
Edit applied successfully.
```

---

## Assistant (Build · Kimi K2.7 Code · 5.5s)

_Thinking:_

Also update the "Technical home in Rust" line in arcade section since Rust is cancelled.

**Tool: edit**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/project-memory.md",
  "oldString": "**Technical home in Rust:** `gitswitch-arcade` crate, separated from core logic.",
  "newString": "**Technical home:** an `internal/arcade/` package in the Go tree, separated from core logic. (Originally planned as a `gitswitch-arcade` Rust crate; cancelled with the Rust migration.)"
}
```

**Output:**
```
Edit applied successfully.
```

---

## Assistant (Build · Kimi K2.7 Code · 14.1s)

_Thinking:_

Now let me create a concise versioning plan doc that ties everything together from v0.1.22.

**Tool: write**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/roadmap-versioning-plan.md",
  "content": "# gitswitch Versioning Plan — From v0.1.22 to v1.0\n\n_Last updated: June 2026_\n\n## Where we are\n\n- Current release: **v0.1.22**.\n- Language: **Go**.\n- Rust migration: **cancelled**.\n- Storage: `~/.gitswitch/profiles.json`.\n\n## Where we are going\n\n| Version | Theme | Big-ticket items | Why this grouping |\n|---|---|---|---|\n| **v0.2.0** | Identity + isolation | Multi-host profiles, keychain PATs, session isolation, TUI refresh, `status`/`doctor`, shell tab completion | These are the infrastructure pieces that change how profiles work. They ship together because they all touch the same storage model and shell integration. |\n| **v0.3.0** | Repo + team | In-repo `.gitswitch.yaml`, team presets, host auto-detection, remote retargeting, remote-URL routing, `gitswitch run` | These features assume v0.2.0's multi-host model is stable. They make gitswitch work across teams and machines. |\n| **v1.0** | Production auth | Direct GitHub OAuth device flow, `gs` alias, HTTPS credential helper, semver + schema stability | 1.0 is the \"gitswitch owns auth end-to-end\" release. No new core data model — just hardened, complete coverage. |\n\n## What this means for sequencing\n\n1. **Ship v0.2.0 before thinking about v0.3.0.** v0.2.0 is a breaking config migration (JSON → YAML, single-host → multi-host). Everything after it builds on the YAML schema.\n2. **Do not put OAuth in v0.2.0.** v0.2.0 already has enough load-bearing changes. OAuth is a v1.0 requirement because it completes the auth story and lets gitswitch stand alone without `gh`.\n3. **Keep the Go codebase.** No Rust rewrite. The time that would have gone to parity/cutover now goes to user-facing features.\n\n## What the user sees\n\n### v0.1.22 → v0.2.0\n```\ngitswitch switch corp          # now sets SSH + gh + git for ALL hosts in the profile\ngitswitch shell corp           # isolated terminal session\ngitswitch status               # \"why is my push failing?\"\ngitswitch import gh            # pull accounts from gh CLI\ngitswitch pin corp --remote origin\n```\n\n### v0.2.0 → v0.3.0\n```\ngitswitch apply ./team.yaml    # team onboarding\ngitswitch remote retarget ...  # org migration helper\ngitswitch run work -- ./script.sh\n```\n\n### v0.3.0 → v1.0\n```\ngitswitch login                # direct GitHub auth, no gh needed\ngs status                      # short alias\ngit config --global credential.helper gitswitch\n```\n\n## Versioning rules\n\n- **Minor versions (0.2.0, 0.3.0)** can introduce new commands and config fields, but must load old config automatically.\n- **Patch versions (0.2.1, 0.2.2)** are bug fixes and small UX improvements only.\n- **v1.0.0** starts semver. Breaking CLI changes require a major bump.\n- Features do not move backwards. If OAuth is v1.0, it doesn't slip into a late v0.3.x patch.\n\n## Open decisions\n\n- Should v0.2.0 be a single release or a short v0.2.x series? Recommendation: ship v0.2.0 when multi-host + session + TUI are done, then patch as needed.\n- Should v0.3.0 features be time-boxed? Recommendation: yes — if team presets drag, ship host auto-detection + retargeting as v0.3.0 and move presets to v0.4.0 rather than block.\n"
}
```

**Output:**
```
Wrote file successfully.
```

---

## Assistant (Build · Kimi K2.7 Code · 6.3s)

_Thinking:_

Now I should update the v0.2.0 roadmap to remove the "OAuth strategy: don't roll our own" locked decision since now we will roll our own in v1.0. Actually, for v0.2.0 the decision is still correct — v0.2.0 imports from gh. But I should add a note that direct OAuth is v1.0. Let me update roadmap-v0.2.0.md slightly.

Also, I should check if there are other docs that reference the Rust plan or outdated timelines. Let me search for references to rust/roadmap or "Rust becomes the default" etc.

**Tool: grep**

**Input:**
```json
{
  "pattern": "Rust becomes|Rust binary|Rust migration|docs/rust/|rust/roadmap|v0\\.2\\.0\\+.*Rust|Rust for v0\\.2\\.0",
  "path": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up"
}
```

**Output:**
```
Found 31 matches
/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/roadmap-versioning-plan.md:
  Line 9: - Rust migration: **cancelled**.


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/project-memory.md:
  Line 15: ~~gitswitch v0.2.0+ will be a clean rewrite in Rust. Go remains the shipping v0.1.x tool. Decision locked May 2026 — not up for debate.~~

  Line 30: - `docs/rust-cancelled/` — all archived v0.2.0+ Rust plans

  Line 96: **Technical home:** an `internal/arcade/` package in the Go tree, separated from core logic. (Originally planned as a `gitswitch-arcade` Rust crate; cancelled with the Rust migration.)

  Line 188: | Rust for v0.2.0+ | CANCELLED | Go remains for all versions; see `docs/rust-cancelled/` |


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/roadmap-full-timeline.md:
  Line 4: _Status: Rust migration **cancelled**. v1.0 will ship on the Go codebase._

  Line 141: ## Rust migration — CANCELLED


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/why-rust.md:
  Line 5: > The Rust migration described in this document has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.

  Line 33: **Rust's answer**: No runtime. No GC. Rust compiles to bare machine code. Cold-start for a Rust binary doing minimal work (read env var, print string, exit) is measurable in microseconds. For the shell hook path — which runs `gitswitch record --path .` on every `cd` — this difference is the difference between a transparent tool and an annoying one.

  Line 45: **Rust's answer**: With `lto = true`, `opt-level = "z"` (size-optimize), `strip = true`, and `codegen-units = 1`, a Rust binary for a tool of gitswitch's scope compiles to ~1.5–3MB. Fully self-contained. No dynamic runtime dependencies beyond the OS.

  Line 119: > The Rust migration plan has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.

  Line 122: ~~gitswitch lives in the terminal, fires on every directory change, handles secrets, and needs a polished animated TUI. Those constraints pick Rust. The migration plan is in [docs/rust/roadmap-migration.md](rust/roadmap-migration.md); the v0.2.0 feature plans are in [docs/rust/](rust/).~~

  Line 124: The archived migration plan is in [`docs/rust-cancelled/roadmap-migration.md`](rust-cancelled/roadmap-migration.md); the archived v0.2.0+ Rust feature plans are in [`docs/rust-cancelled/`](rust-cancelled/).


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-migration.md:
  Line 5: > The Rust migration plan described in this document has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.

  Line 40: The corrode.dev [Rust migration guide](https://corrode.dev/blog/migrating-from-go-to-rust/) recommends "piece-by-piece" porting through a long parallel-runtime period for codebases where the cost of feature-freezing the original is high. gitswitch's cost of feature-freezing is essentially zero — there are no contributors blocked on v0.1.x development. The piece-by-piece pattern's main payoff (de-risking) is achieved more cheaply here by the parallel-binary scheme above without the cgo overhead.

  Line 105: - Rename: archive `cmd/gitswitch/` (Go) to `legacy/` branch; rename Rust binary `gitswitch-rs` → `gitswitch`; ship as `v0.2.0`.

  Line 302: **The v0.2.0 YAML schema is the boundary.** The Rust binary must read and write `~/.config/gitswitch/config.yaml` byte-identically to the Go binary's v0.2.0 output. Migration from v0.1's `profiles.json` is implemented behavior-for-behavior in `gitswitch-core`:

  Line 327: | Two-binary distribution during transition    | Releases tagged `0.2.0-rcN-rs` ship the Rust binary as `gitswitch-rs`; Go releases continue under `gitswitch` until parity is signed off. Both Homebrew taps and install scripts updated to install both side-by-side; users opt in via `--rs` flag on install script. |

  Line 389: 7. **Install scripts updated** — Homebrew formula, install.sh, all distribution channels swap to Rust binary atomically with the rename.

  Line 394: - Rename Rust binary artifact `gitswitch-rs` → `gitswitch` in `Cargo.toml` + `cargo-dist` config.


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-v0.3.0.md:
  Line 5: > The Rust migration plan described in this document has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-v0.2.0-tui.md:
  Line 5: > The Rust migration plan described in this document has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-v0.2.0-isolation.md:
  Line 5: > The Rust migration plan described in this document has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-v0.2.0.md:
  Line 5: > The Rust migration plan described in this document has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/versioning.md:
  Line 88: 1. Rust binary reaches feature parity with Go (v0.2.0 feature set)

  Line 89: 2. Rust binary passes the full parity test corpus

  Line 96: - Rust binary is the default (Go binary deprecated, kept for one release cycle)

  Line 109: - [ ] Rust binary: all v0.1.x features ported

  Line 110: - [ ] Rust binary: v0.2.0 multi-host + session isolation ported

  Line 111: - [ ] Rust binary: ratatui TUI with arcade theme + tachyonfx effects


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/marketing.md:
  Line 36: **When**: v0.2.0 ships. That's the milestone worth announcing — multi-host, session isolation, full TUI, Rust binary.

```

---

## Assistant (Build · Kimi K2.7 Code · 8.2s)

_Thinking:_

I see there are more files to update:
1. `versioning.md` - needs updating to remove Rust references
2. `marketing.md` - needs updating to remove Rust binary mention

Let me read these and update them.

**Tool: read**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/versioning.md"
}
```

**Output:**
```
<path>/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/versioning.md</path>
<type>file</type>
<content>
1: # gitswitch — Versioning Strategy
2: 
3: _Written May 2026. Governs all releases from v0.1.21 onward._
4: 
5: ---
6: 
7: ## The version ladder
8: 
9: ```
10: v0.1.x  ──── current patch series (Go, in production)
11: v0.2.0  ──── multi-host + session isolation (Go)
12: v1.0.0  ──── stable release, Rust default, brand launch ← THE milestone
13: v1.x.x  ──── bug fixes and polish on stable
14: v2.0.0  ──── next major paradigm (defined closer to the time)
15: ```
16: 
17: ---
18: 
19: ## What each version number means
20: 
21: ```
22:   v  0  .  1  .  22
23:      │     │     └─ PATCH — bug fix, no new user-facing feature
24:      │     └─────── MINOR — new feature milestone, backwards compatible
25:      └─────────────  MAJOR — breaking CLI change or paradigm shift
26: ```
27: 
28: **Pre-1.0 rule**: minor bumps (0.x) are feature milestones. Patch bumps (0.1.x) are bug fixes and small additive features that don't require a plan milestone.
29: 
30: **Post-1.0 rule**: strict semver. No breaking CLI changes without bumping major. Deprecation notices before anything is removed.
31: 
32: ---
33: 
34: ## Current patch series — v0.1.x
35: 
36: Small additions and fixes to the Go codebase. No schema changes, no new commands that require planning.
37: 
38: | Version | Contents | Status |
39: |---|---|---|
40: | v0.1.20 | Latest stable (brew tap) | ✓ shipped |
41: | v0.1.21 | HTTPS credential helper (gh-delegated) + interactive install wizard | PR #13 — merge → tag |
42: | v0.1.22 | hook-check dev version noise fix + combined upgrade+HTTPS nudge | PR #14 — merge → tag |
43: | v0.1.23+ | Any further bug fixes before v0.2.0 | as needed |
44: 
45: **What belongs in v0.1.x:**
46: - Bug fixes to existing commands
47: - UX polish (messages, hints, wizard copy)
48: - Hook/install improvements
49: - Nothing that changes the profile schema or adds new required config
50: 
51: ---
52: 
53: ## v0.2.0 — The enterprise + AI agent release (Go)
54: 
55: Two features ship together. Both are planned in detail in `roadmap-v0.2.0.md` and `roadmap-v0.2.0-isolation.md`.
56: 
57: **Multi-host identity**
58: - Profile model gains `hosts` map (multiple GitHub hosts per profile)
59: - Storage migrates `profiles.json` → `~/.config/gitswitch/config.yaml`
60: - OS keychain for PAT tokens
61: - `gitswitch import gh` — auto-discover accounts from gh CLI
62: - `gitswitch host add/remove/list`
63: 
64: **Session isolation**
65: - `gitswitch shell [profile]` — per-terminal identity via env vars
66: - Parallel AI agents with independent identity scopes
67: - Full gitswitch works inside a session
68: 
69: **Diagnostics and safety**
70: - `gitswitch status` — per-repo narrative diagnostic
71: - Pre-commit safety net (auto-installed, warns before wrong-identity commit)
72: - Progressive CLI hints
73: 
74: **TUI refresh**
75: - 5 screens: Profiles, Hosts, Repos, Session, Doctor
76: - Command palette
77: - Shell tab completion
78: 
79: **Targeting**: once v0.1.x is stable and the multi-host plan is implementation-ready.
80: 
81: ---
82: 
83: ## v1.0.0 — The stable release
84: 
85: **This is the milestone that matters most.**
86: 
87: v1.0 ships when:
88: 1. Rust binary reaches feature parity with Go (v0.2.0 feature set)
89: 2. Rust binary passes the full parity test corpus
90: 3. gitswitch.dev is live with full docs, privacy policy, blog
91: 4. **Logo and brand are finalized and shipped**
92: 5. Homebrew-core submission ready (not just tap)
93: 6. Semver commitment in place — no breaking CLI changes without a major bump
94: 
95: **What v1.0 is:**
96: - Rust binary is the default (Go binary deprecated, kept for one release cycle)
97: - Stable config format: `~/.config/gitswitch/config.yaml` schema versioned, migration guaranteed
98: - Full arcade TUI with tachyonfx (splash, CRT, all easter eggs)
99: - HTTPS credential helper (full keychain-backed, not gh-delegated)
100: - AUR + cargo install + Windows installer added to distribution
101: - First public brand moment — logo, tagline, gitswitch.dev launch post
102: 
103: **What v1.0 is NOT:**
104: - Not a feature freeze. v1.0 ships everything from v0.2.0 + v0.3.0 that's ready.
105: - Not a rewrite of the Go code. The Rust rewrite happens in parallel and replaces the binary on v1.0.
106: - Items not ready go to v1.1.
107: 
108: **Blocking checklist for v1.0:**
109: - [ ] Rust binary: all v0.1.x features ported
110: - [ ] Rust binary: v0.2.0 multi-host + session isolation ported
111: - [ ] Rust binary: ratatui TUI with arcade theme + tachyonfx effects
112: - [ ] Parity test corpus passing
113: - [ ] `~/.config/gitswitch/config.yaml` stable schema with versioned migration
114: - [ ] gitswitch.dev live (docs, blog, privacy)
115: - [ ] Logo finalized + shipped in binary (`gitswitch version` shows it)
116: - [ ] homebrew-core PR submitted
117: 
118: ---
119: 
120: ## v1.x — Stable patch series
121: 
122: After v1.0, no new features in patch releases. Only:
123: - Bug fixes
124: - Security patches
125: - Performance improvements
126: - Compatibility updates (new shell versions, new gh CLI versions, new macOS)
127: 
128: New features go in v1.1, v1.2, etc.
129: 
130: **Examples of what goes where post-v1.0:**
131: 
132: | Change | Version |
133: |---|---|
134: | Fix: SSH key permissions check incorrect on Linux | v1.0.x |
135: | Fix: credential helper crashes on gh not installed | v1.0.x |
136: | Feat: remote-URL-aware routing | v1.1.0 |
137: | Feat: team preset distribution | v1.1.0 |
138: | Feat: in-repo `.gitswitch.yaml` | v1.2.0 |
139: | Breaking: rename `gitswitch current` to `gitswitch whoami` | v2.0.0 |
140: 
141: ---
142: 
143: ## v2.0.0 — Next major (future)
144: 
145: Not defined yet. Triggered when a breaking change or paradigm shift warrants it. Could be:
146: - Major config schema change
147: - New core model (e.g. workspace-level identity)
148: - Platform expansion that changes the CLI surface
149: 
150: Define closer to when v1.x matures.
151: 
152: ---
153: 
154: ## How to tag a release
155: 
156: ```bash
157: # After PR merges to main:
158: git tag -a v0.1.22 -m "fix: hook-check version noise"
159: git push origin v0.1.22
160: 
161: # Brew tap picks it up via the formula's url/sha update
162: # Update aksisonline/homebrew-tap Formula/gitswitch.rb → new version + sha
163: ```
164: 
165: Versioning in the binary is injected at link time by the Makefile/goreleaser:
166: ```
167: -ldflags="-X main.version=v0.1.22"
168: ```
169: 
170: Never commit a version bump to the source — `var version = "dev"` stays in source always.
171: 
172: ---
173: 
174: ## Commit → version bump rules
175: 
176: | Commit prefix | Version bump |
177: |---|---|
178: | `fix:` | PATCH (0.1.x) |
179: | `feat:` | PATCH while in v0.1.x series; MINOR for planned milestones (v0.2.0) |
180: | `refactor:` / `chore:` / `docs:` | PATCH if anything; often no release |
181: | `feat!:` / `BREAKING CHANGE:` | MAJOR (only post-v1.0) |
182: 
183: Pre-v1.0, no strict semver. Tag when the code is ready, not by formula.
184: 
185: ---
186: 
187: ## The path in one line
188: 
189: ```
190: now → v0.1.22 → v0.2.0 (multi-host+isolation) → Rust parity → v1.0.0 (brand+stable) → v1.x patches → v2.0.0
191: ```

(End of file - total 191 lines)
</content>
```

**Tool: read**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/marketing.md"
}
```

**Output:**
```
<path>/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/marketing.md</path>
<type>file</type>
<content>
1: # gitswitch — Marketing Strategy
2: 
3: _Privacy-first. No ads. No paid acquisition. Earned attention only._
4: 
5: _Date: May 2026_
6: 
7: ---
8: 
9: ## Core principle
10: 
11: Developers distrust marketing. They reward tools that solve a real problem, have good docs, and feel like they were built by someone who actually uses them. gitswitch's marketing is not campaigns — it's presence in the right places at the right moments, and content that earns its way in.
12: 
13: The arcade personality, the privacy stance, the Rust rewrite, the AI agent angle — all of these are marketing assets hiding as product decisions. Use them.
14: 
15: ---
16: 
17: ## The angles (what makes gitswitch worth talking about)
18: 
19: Every channel below needs a reason to post. These are gitswitch's genuine ones:
20: 
21: | Angle | Why it lands |
22: |---|---|
23: | "I kept committing as the wrong person" | Universal pain. Every multi-account dev has done this. Instantly relatable. |
24: | Privacy-first, fully local, no cloud | Rare positioning. Most dev tools are moving to cloud. This is a deliberate counter. |
25: | AI agents need identity isolation | Very 2026. Agents pushing code under the wrong identity is a new and growing problem. gitswitch solves it with session scoping. |
26: | Hidden Pacman mode / arcade easter eggs | Shareable. "This git tool has a hidden game" is a tweet that gets 500 likes from people who've never touched the tool. |
27: | Rust rewrite | The Rust community pays attention to "X rewritten in Rust" — and writes about it. |
28: | Shell hook that *learns* your identity per repo | The intelligence angle. Not just switching — automatic recommendation based on usage history. |
29: 
30: ---
31: 
32: ## Channel strategy
33: 
34: ### 1. Show HN — the launch moment
35: 
36: **When**: v0.2.0 ships. That's the milestone worth announcing — multi-host, session isolation, full TUI, Rust binary.
37: 
38: Hacker News Show HN is the highest-value single event for a developer CLI tool. Front-page Show HN gets 5K–15K visits, 3–8% convert to installs. Technical audience, direct feedback, no fluff.
39: 
40: **Title formula**: `Show HN: gitswitch – manage multiple GitHub identities and hosts from the terminal`
41: 
42: **Body**: technical story. Why the problem is harder than it looks (global git config, multi-host PATs, AI agents needing parallel identities). What the architecture does. Rust, no cloud, no telemetry. Link to the privacy stance directly.
43: 
44: **Do not**: make it sound like a product launch. HN readers punish marketing copy. Write like you're explaining it to a senior engineer.
45: 
46: **Timing**: Tuesday–Thursday, 9–11am US Eastern. Prepare for comments. Be present for 4 hours after posting.
47: 
48: ### 2. Reddit — sustained presence
49: 
50: Subreddits, in order of priority:
51: 
52: - **r/git** — the obvious home. Post genuinely. Answer questions about managing multiple accounts — gitswitch is often the right answer.
53: - **r/rust** — when v0.2.0 ships in Rust. "We rewrote our Go CLI tool in Rust — here's what changed." Technical writeup, not promotion.
54: - **r/devops** — enterprise angle. Multi-host GitHub, SSH key management at scale.
55: - **r/commandline** — terminal tool enthusiasts. The TUI and arcade angle lands here.
56: - **r/programming** — Show Reddit moment alongside Show HN.
57: - **r/MachineLearning / r/LocalLLaMA** — the AI agent identity isolation angle. "Your coding agent is committing as the wrong person."
58: 
59: **Rule**: never just post a link. Contribute to the conversation first. When gitswitch is genuinely the answer to someone's problem, mention it.
60: 
61: ### 3. Twitter / X
62: 
63: The devtools community is active here. Target audience: engineers who tweet about terminal setup, dotfiles, Rust, AI coding agents.
64: 
65: **Content types that work**:
66: - Short terminal recordings (VHS, ttyrec) showing a feature. 30 seconds max.
67: - The arcade easter egg reveal — "we hid a game in a git tool"
68: - The session isolation angle framed for AI developers — "your Claude Code agent is pushing commits under the wrong identity"
69: - Rust rewrite thread — "why we moved gitswitch from Go to Rust" (5-tweet thread, technical)
70: - Milestone posts — "gitswitch hit N installs / first enterprise adoption"
71: 
72: **Accounts to engage genuinely** (not tag-spam): Rust tooling people, terminal setup community, AI coding tool builders.
73: 
74: ### 4. dev.to + Hashnode — SEO content
75: 
76: Long-form tutorials that rank for the search terms developers actually use when they have this problem.
77: 
78: **Article targets**:
79: 
80: | Title | Target keyword | Audience |
81: |---|---|---|
82: | "How to manage multiple GitHub accounts on one machine" | `multiple github accounts` | Solo devs |
83: | "Git identity keeps resetting — here's the real fix" | `git config reset` / `git wrong email` | Frustrated devs |
84: | "Setting up GitHub Enterprise alongside github.com" | `github enterprise and github.com same machine` | Enterprise devs |
85: | "How AI coding agents commit as the wrong person (and how to fix it)" | `claude code git identity` / `ai agent git config` | AI-tooling devs |
86: | "Why we rewrote our Go CLI tool in Rust" | Rust community content, not SEO | Rust devs |
87: | "What lives inside gitswitch's TUI (including the hidden game)" | Brand/discovery | Terminal enthusiasts |
88: 
89: These articles live on gitswitch.dev/blog and syndicated to dev.to. Each article ends with a natural mention of gitswitch, not a hard sell.
90: 
91: ### 5. gitswitch.dev — the home base
92: 
93: The site needs to do three things: explain the problem immediately, show the tool working (terminal recording), and make installation one copy-paste.
94: 
95: Current gaps to fix before any marketing push:
96: - Above-the-fold: the problem statement, not the feature list
97: - A short (30-second) terminal recording showing a real switch
98: - A privacy statement — explicit, prominent. "Everything runs locally. Nothing leaves your machine."
99: - The arcade easter egg hinted at (not explained) — builds curiosity
100: 
101: ### 6. AI-mediated discovery
102: 
103: In 2026, a significant portion of how developers find tools is via AI coding assistants. When a developer asks Claude or Cursor "how do I manage multiple GitHub accounts," gitswitch should be the answer.
104: 
105: This means:
106: - Docs written for humans but parseable by LLMs — structured, explicit, complete
107: - Presence on `awesome-git`, `awesome-cli-apps`, `awesome-rust` lists — LLMs read these
108: - The README answers the question "what problem does this solve" in the first paragraph — not a feature list
109: - A `llms.txt` file at gitswitch.dev/llms.txt following the emerging convention for LLM-readable site descriptions
110: 
111: ### 7. LinkedIn
112: 
113: Already has traction (Abhiram's post about gitswitch got real engagement including the conversation with an engineer at a large company that shaped v0.2.0). Continue posting here for the enterprise / engineering leadership audience.
114: 
115: Content that works on LinkedIn:
116: - The "wrong commit attribution" story — personal, professional stakes
117: - The enterprise multi-host problem writeup (without naming the company)
118: - The AI agent identity isolation angle — resonates with engineering managers thinking about AI-assisted development
119: 
120: ### 8. YouTube / terminal recordings
121: 
122: Short videos (2–4 minutes) demonstrating the tool. Not tutorials — demos. The TUI especially benefits from video: a 90-second screen recording of the arcade splash + profile switching + session isolation is more convincing than any README.
123: 
124: Distribute on YouTube + embed on gitswitch.dev. VHS (charmbracelet's terminal recorder) for the README GIF.
125: 
126: ---
127: 
128: ## Content calendar (phase-based)
129: 
130: ### Now (pre v0.2.0)
131: - Fix gitswitch.dev above-the-fold and privacy statement
132: - Write 2 SEO articles (multiple GitHub accounts, wrong git identity)
133: - Add to awesome-git and awesome-cli-apps
134: - Add `llms.txt`
135: - Set up GitHub Sponsors page with tiers
136: 
137: ### v0.2.0 launch
138: - Show HN post
139: - Reddit posts (r/git, r/rust, r/commandline)
140: - Twitter launch thread (with TUI recording)
141: - dev.to article: "Why we rewrote gitswitch in Rust"
142: - LinkedIn post: multi-host enterprise angle
143: - Product Hunt (secondary — do 1 week after Show HN)
144: 
145: ### Post-launch (sustained)
146: - Monthly blog post on gitswitch.dev
147: - Reddit: answer questions about multiple GitHub accounts and mention gitswitch when genuinely relevant
148: - Twitter: short feature demos as they ship
149: - The arcade easter egg reveal — save this for when the TUI ships, it's a separate moment
150: 
151: ---
152: 
153: ## What success looks like
154: 
155: Not vanity metrics. Concrete signals:
156: 
157: - **GitHub stars**: 1K is credible, 5K is momentum, 10K is "this is a real tool"
158: - **Sponsor revenue**: first company sponsor ($100+/month) validates enterprise value
159: - **Organic mentions**: people recommending gitswitch in threads where it wasn't seeded
160: - **Consulting inbounds**: companies reaching out about custom integration
161: - **Contributions**: outside contributors opening PRs — signals community health
162: 
163: ---
164: 
165: ## What not to do
166: 
167: - No cold outreach to developers
168: - No "please share this" asks at the end of posts
169: - No posting the same content across platforms same-day without adapting tone
170: - No buying followers, stars, or reviews
171: - No exaggerating feature claims — the audience will test it immediately
172: 
173: ---
174: 
175: ## Sources
176: 
177: - [Developer marketing playbook (2026)](https://www.strategicnerds.com/blog/the-complete-developer-marketing-guide-2026)
178: - [Show HN vs Product Hunt for dev tools](https://medium.com/@baristaGeek/lessons-launching-a-developer-tool-on-hacker-news-vs-product-hunt-and-other-channels-27be8784338b)
179: - [Open source organic growth](https://business.daily.dev/resources/open-source-marketing-grow-developer-community-without-budget/)
180: - [AI-mediated discovery in 2026](https://www.strategicnerds.com/blog/the-complete-developer-marketing-guide-2026)
181: - [DevTools content marketing](https://www.saashero.net/content/devtools-content-marketing-strategies/)

(End of file - total 181 lines)
</content>
```

---

## Assistant (Build · Kimi K2.7 Code · 32.7s)

_Thinking:_

I need to update both files. For versioning.md, it's a significant rewrite. For marketing.md, smaller changes.

Let me rewrite versioning.md entirely since it has many Rust references and outdated v1.0 criteria.

**Tool: write**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/versioning.md",
  "content": "# gitswitch — Versioning Strategy\n\n_Written June 2026. Governs all releases from v0.1.22 onward._  \n_Status: Rust migration **cancelled**. v1.0 ships on Go._\n\n---\n\n## The version ladder\n\n```\nv0.1.x  ──── current patch series (Go, in production)\nv0.2.0  ──── multi-host + session isolation + TUI refresh (Go)\nv0.3.0  ──── repo + team features (Go)\nv1.0.0  ──── stable release, direct OAuth, `gs` alias, HTTPS credential helper ← THE milestone\nv1.x.x  ──── bug fixes and polish on stable\nv2.0.0  ──── next major paradigm (defined closer to the time)\n```\n\n---\n\n## What each version number means\n\n```\n  v  0  .  1  .  22\n     │     │     └─ PATCH — bug fix, no new user-facing feature\n     │     └─────── MINOR — new feature milestone, backwards compatible\n     └─────────────  MAJOR — breaking CLI change or paradigm shift\n```\n\n**Pre-1.0 rule**: minor bumps (0.x) are feature milestones. Patch bumps (0.1.x) are bug fixes and small additive features that don't require a plan milestone.\n\n**Post-1.0 rule**: strict semver. No breaking CLI changes without bumping major. Deprecation notices before anything is removed.\n\n---\n\n## Current patch series — v0.1.x\n\nSmall additions and fixes to the Go codebase. No schema changes, no new commands that require planning.\n\n| Version | Contents | Status |\n|---|---|---|\n| v0.1.22 | hook-check dev version noise fix + combined upgrade+HTTPS nudge | ✓ shipped |\n| v0.1.23+ | Any further bug fixes before v0.2.0 | as needed |\n\n**What belongs in v0.1.x:**\n- Bug fixes to existing commands\n- UX polish (messages, hints, wizard copy)\n- Hook/install improvements\n- Nothing that changes the profile schema or adds new required config\n\n---\n\n## v0.2.0 — The enterprise + AI agent release (Go)\n\nTwo features ship together. Both are planned in detail in `roadmap-v0.2.0.md` and `roadmap-v0.2.0-isolation.md`.\n\n**Multi-host identity**\n- Profile model gains `hosts` map (multiple GitHub hosts per profile)\n- Storage migrates `profiles.json` → `~/.config/gitswitch/config.yaml`\n- OS keychain for PAT tokens\n- `gitswitch import gh` — auto-discover accounts from gh CLI\n- `gitswitch host add/remove/list`\n\n**Session isolation**\n- `gitswitch shell [profile]` — per-terminal identity via env vars\n- Parallel AI agents with independent identity scopes\n- Full gitswitch works inside a session\n\n**Diagnostics and safety**\n- `gitswitch status` — per-repo narrative diagnostic\n- Pre-commit safety net (auto-installed, warns before wrong-identity commit)\n- Progressive CLI hints\n\n**TUI refresh**\n- 5 screens: Profiles, Hosts, Repos, Session, Doctor\n- Command palette\n- Shell tab completion\n\n**Targeting**: once v0.1.x is stable and the multi-host plan is implementation-ready.\n\n---\n\n## v0.3.0 — Repo + team features (Go)\n\nBuilds on v0.2.0's multi-host model. Detailed in `roadmap-v0.3.0.md`.\n\n- In-repo `.gitswitch.yaml` config\n- Team preset distribution with signing (`gitswitch apply`)\n- Host auto-detection from `git remote -v`\n- Remote retargeting for org migrations\n- Remote-URL-aware pinning\n- `gitswitch run <profile> -- <cmd>`\n- Profile export/import\n\n**Targeting**: after v0.2.0 is stable and the YAML schema has proven itself.\n\n---\n\n## v1.0.0 — The stable release\n\n**This is the milestone that matters most.**\n\nv1.0 ships when:\n1. v0.2.0 and v0.3.0 features are stable and documented\n2. Direct GitHub OAuth device flow is implemented\n3. `gs` short command alias ships\n4. HTTPS credential helper is implemented\n5. gitswitch.dev is live with full docs, privacy policy, blog\n6. **Logo and brand are finalized and shipped**\n7. Homebrew-core submission ready (not just tap)\n8. Semver commitment in place — no breaking CLI changes without a major bump\n\n**What v1.0 is:**\n- Stable Go binary\n- Stable config format: `~/.config/gitswitch/config.yaml` schema versioned, migration guaranteed\n- gitswitch owns auth end-to-end: SSH, HTTPS, and gh CLI sync\n- `gs` alias for quick invocation\n- AUR + Windows installer added to distribution\n- First public brand moment — logo, tagline, gitswitch.dev launch post\n\n**What v1.0 is NOT:**\n- Not a language rewrite. The Rust migration was cancelled.\n- Not a feature freeze. v1.0 ships everything from v0.2.0 + v0.3.0 that's ready.\n- Items not ready go to v1.1.\n\n**Blocking checklist for v1.0:**\n- [ ] All v0.1.x features stable in Go\n- [ ] v0.2.0 multi-host + session isolation shipped and stable\n- [ ] v0.3.0 repo/team features shipped (or deferred to v1.1)\n- [ ] Direct GitHub OAuth device flow works against github.com and GHES\n- [ ] `gs` alias ships on macOS, Linux, Windows\n- [ ] HTTPS credential helper works with HTTPS remotes\n- [ ] `~/.config/gitswitch/config.yaml` stable schema with versioned migration\n- [ ] gitswitch.dev live (docs, blog, privacy)\n- [ ] Logo finalized + shipped in binary (`gitswitch version` shows it)\n- [ ] homebrew-core PR submitted\n\n---\n\n## v1.x — Stable patch series\n\nAfter v1.0, no new features in patch releases. Only:\n- Bug fixes\n- Security patches\n- Performance improvements\n- Compatibility updates (new shell versions, new gh CLI versions, new macOS)\n\nNew features go in v1.1, v1.2, etc.\n\n**Examples of what goes where post-v1.0:**\n\n| Change | Version |\n|---|---|\n| Fix: SSH key permissions check incorrect on Linux | v1.0.x |\n| Fix: credential helper crashes on gh not installed | v1.0.x |\n| Feat: new easter egg game | v1.1.0 |\n| Feat: additional OAuth host support | v1.1.0 |\n| Breaking: rename `gitswitch current` to `gitswitch whoami` | v2.0.0 |\n\n---\n\n## v2.0.0 — Next major (future)\n\nNot defined yet. Triggered when a breaking change or paradigm shift warrants it. Could be:\n- Major config schema change\n- New core model (e.g. workspace-level identity)\n- Platform expansion that changes the CLI surface\n\nDefine closer to when v1.x matures.\n\n---\n\n## How to tag a release\n\n```bash\n# After PR merges to main:\ngit tag -a v0.1.22 -m \"fix: hook-check version noise\"\ngit push origin v0.1.22\n\n# Brew tap picks it up via the formula's url/sha update\n# Update aksisonline/homebrew-tap Formula/gitswitch.rb → new version + sha\n```\n\nVersioning in the binary is injected at link time by the Makefile/goreleaser:\n```\n-ldflags=\"-X main.version=v0.1.22\"\n```\n\nNever commit a version bump to the source — `var version = \"dev\"` stays in source always.\n\n---\n\n## Commit → version bump rules\n\n| Commit prefix | Version bump |\n|---|---|\n| `fix:` | PATCH (0.1.x) |\n| `feat:` | PATCH while in v0.1.x series; MINOR for planned milestones (v0.2.0) |\n| `refactor:` / `chore:` / `docs:` | PATCH if anything; often no release |\n| `feat!:` / `BREAKING CHANGE:` | MAJOR (only post-v1.0) |\n\nPre-v1.0, no strict semver. Tag when the code is ready, not by formula.\n\n---\n\n## The path in one line\n\n```\nnow → v0.1.22 → v0.2.0 (multi-host+isolation) → v0.3.0 (repo+team) → v1.0.0 (OAuth+gs+HTTPS+brand+stable) → v1.x patches → v2.0.0\n```\n"
}
```

**Output:**
```
Wrote file successfully.
```

---

## Assistant (Build · Kimi K2.7 Code · 13.5s)

_Thinking:_

Now update marketing.md to remove Rust references.

**Tool: edit**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/marketing.md",
  "oldString": "The arcade personality, the privacy stance, the Rust rewrite, the AI agent angle — all of these are marketing assets hiding as product decisions. Use them.",
  "newString": "The arcade personality, the privacy stance, the AI agent angle — all of these are marketing assets hiding as product decisions. Use them."
}
```

**Output:**
```
Edit applied successfully.
```

**Tool: edit**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/marketing.md",
  "oldString": "| Rust rewrite | The Rust community pays attention to \"X rewritten in Rust\" — and writes about it. |",
  "newString": "| Open source + fully local | Auditable, no cloud, no telemetry — a deliberate counter to cloud-everything tooling. |"
}
```

**Output:**
```
Edit applied successfully.
```

**Tool: edit**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/marketing.md",
  "oldString": "**When**: v0.2.0 ships. That's the milestone worth announcing — multi-host, session isolation, full TUI, Rust binary.",
  "newString": "**When**: v0.2.0 ships. That's the milestone worth announcing — multi-host, session isolation, full TUI refresh."
}
```

**Output:**
```
Edit applied successfully.
```

**Tool: edit**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/marketing.md",
  "oldString": "- **r/rust** — when v0.2.0 ships in Rust. \"We rewrote our Go CLI tool in Rust — here's what changed.\" Technical writeup, not promotion.",
  "newString": "- **r/rust** — not a primary channel since the Rust rewrite was cancelled. Reconsider only if there is a Rust-specific technical angle in the future."
}
```

**Output:**
```
Edit applied successfully.
```

**Tool: edit**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/marketing.md",
  "oldString": "- Rust rewrite thread — \"why we moved gitswitch from Go to Rust\" (5-tweet thread, technical)",
  "newString": "- Privacy-first thread — \"why gitswitch will never have a cloud component\" (5-tweet thread, technical)"
}
```

**Output:**
```
Edit applied successfully.
```

**Tool: edit**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/marketing.md",
  "oldString": "| \"Why we rewrote our Go CLI tool in Rust\" | Rust community content, not SEO | Rust devs |",
  "newString": "| \"Why gitswitch stays local-only\" | privacy-first tooling | security-minded devs |"
}
```

**Output:**
```
Edit applied successfully.
```

**Tool: edit**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/marketing.md",
  "oldString": "- Reddit posts (r/git, r/rust, r/commandline)\n- Twitter launch thread (with TUI recording)\n- dev.to article: \"Why we rewrote gitswitch in Rust\"",
  "newString": "- Reddit posts (r/git, r/commandline)\n- Twitter launch thread (with TUI recording)\n- dev.to article: \"Managing multiple GitHub identities without losing your mind\""
}
```

**Output:**
```
Edit applied successfully.
```

---

## Assistant (Build · Kimi K2.7 Code · 6.4s)

_Thinking:_

Now let me update roadmap-v0.2.0.md to add a note about OAuth being v1.0. The locked decision #3 says "OAuth strategy: don't roll our own" which is still true for v0.2.0, but we should note that direct OAuth is planned for v1.0.

**Tool: edit**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/roadmap-v0.2.0.md",
  "oldString": "3. **OAuth strategy**: don't roll our own. **Import from `gh auth status` / `~/.config/gh/hosts.yml`** — gh CLI users already authenticated their accounts; we discover them, create gitswitch profiles, then layer our customizations (per-host SSH keys, per-repo pinning) on top.",
  "newString": "3. **OAuth strategy for v0.2.0**: don't roll our own. **Import from `gh auth status` / `~/.config/gh/hosts.yml`** — gh CLI users already authenticated their accounts; we discover them, create gitswitch profiles, then layer our customizations (per-host SSH keys, per-repo pinning) on top. A direct GitHub OAuth flow is planned for the v1.0 launch; see [`roadmap-v1.0.md`](roadmap-v1.0.md)."
}
```

**Output:**
```
Edit applied successfully.
```

---

## Assistant (Build · Kimi K2.7 Code · 6.8s)

_Thinking:_

Now let me search for any remaining references to the old Rust plan or outdated OAuth references in the main docs.

**Tool: grep**

**Input:**
```json
{
  "pattern": "Rust.*v0\\.2\\.0|Rust.*v1\\.0|Rust binary|Rust rewrite|Rust migration|docs/rust/|rust/roadmap|own OAuth.*v0\\.3\\.0|OAuth.*v0\\.3\\.0",
  "path": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up",
  "include": "*.md"
}
```

**Output:**
```
Found 43 matches
/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/roadmap-v0.2.0.md:
  Line 7: v0.2.0 closes the gap by making **host a first-class dimension of identity**, **leveraging `gh auth` as the import source** (rather than building our own OAuth), and adding per-repo diagnostics (`gitswitch status`) plus multi-profile pinning. Team-preset YAML distribution, our own OAuth flow, and in-repo config files are explicitly v0.3.0 scope — see [roadmap-v0.3.0.md](roadmap-v0.3.0.md).


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/marketing.md:
  Line 53: - **r/rust** — not a primary channel since the Rust rewrite was cancelled. Reconsider only if there is a Rust-specific technical angle in the future.


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/versioning.md:
  Line 4: _Status: Rust migration **cancelled**. v1.0 ships on Go._

  Line 122: - Not a language rewrite. The Rust migration was cancelled.


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/roadmap-versioning-plan.md:
  Line 9: - Rust migration: **cancelled**.

  Line 24: 3. **Keep the Go codebase.** No Rust rewrite. The time that would have gone to parity/cutover now goes to user-facing features.


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/project-memory.md:
  Line 11: > The Rust rewrite plan is cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.

  Line 96: **Technical home:** an `internal/arcade/` package in the Go tree, separated from core logic. (Originally planned as a `gitswitch-arcade` Rust crate; cancelled with the Rust migration.)

  Line 165: - **crates.io is empty** — every Rust git identity crate is dead. gitswitch v0.2.0 owns the category by default on publish.

  Line 188: | Rust for v0.2.0+ | CANCELLED | Go remains for all versions; see `docs/rust-cancelled/` |


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/roadmap-full-timeline.md:
  Line 4: _Status: Rust migration **cancelled**. v1.0 will ship on the Go codebase._

  Line 141: ## Rust migration — CANCELLED

  Line 145: > The Rust rewrite has been cancelled. The risk of binary-compatible feature parity + migration outweighed the marginal gains for a tool that already works.


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/why-rust.md:
  Line 5: > The Rust migration described in this document has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.

  Line 10: _Original decision locked: May 2026. The Go codebase was the shipping v0.1.x tool; Rust was the target for v0.2.0+. See [docs/rust-cancelled/roadmap-migration.md](rust-cancelled/roadmap-migration.md) for the archived migration plan._

  Line 33: **Rust's answer**: No runtime. No GC. Rust compiles to bare machine code. Cold-start for a Rust binary doing minimal work (read env var, print string, exit) is measurable in microseconds. For the shell hook path — which runs `gitswitch record --path .` on every `cd` — this difference is the difference between a transparent tool and an annoying one.

  Line 45: **Rust's answer**: With `lto = true`, `opt-level = "z"` (size-optimize), `strip = true`, and `codegen-units = 1`, a Rust binary for a tool of gitswitch's scope compiles to ~1.5–3MB. Fully self-contained. No dynamic runtime dependencies beyond the OS.

  Line 111: The Go codebase is not a mistake — it got gitswitch to v0.1.x with real users. ~~Rust is the right choice for where the tool is going.~~ Go remains the right choice for v1.0 and beyond.

  Line 119: > The Rust migration plan has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.

  Line 122: ~~gitswitch lives in the terminal, fires on every directory change, handles secrets, and needs a polished animated TUI. Those constraints pick Rust. The migration plan is in [docs/rust/roadmap-migration.md](rust/roadmap-migration.md); the v0.2.0 feature plans are in [docs/rust/](rust/).~~


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-migration.md:
  Line 5: > The Rust migration plan described in this document has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.

  Line 22: This is not a "rewrite because Rust." It's "v0.2.0 doubles the surface; do it in the language whose ecosystem already shipped the answers."

  Line 30: Turborepo's [Rust-Go-Rust sandwich](https://turbo.build/blog/turbopack-benchmarks) (Rust entry point cgo'ing into Go for unported commands) is the right call for a 100k+ LoC build tool under active feature development. gitswitch is ~4.5k LoC, frozen-feature for the duration of the port (v0.2.0 ships in Rust, not Go). The sandwich's overhead — cgo build matrix, dual-runtime distribution, drift between two implementations — costs more than the rewrite saves.

  Line 32: Instead: **ship the Rust rewrite as `gitswitch-rs` alongside the Go `gitswitch`** until parity is certified. Same config dir, same shell hook contract, same YAML schema. Users opt-in by aliasing or by installing the rs build through a separate install flag. When the parity corpus passes 100%, the binaries swap names and the Go tree gets archived.

  Line 40: The corrode.dev [Rust migration guide](https://corrode.dev/blog/migrating-from-go-to-rust/) recommends "piece-by-piece" porting through a long parallel-runtime period for codebases where the cost of feature-freezing the original is high. gitswitch's cost of feature-freezing is essentially zero — there are no contributors blocked on v0.1.x development. The piece-by-piece pattern's main payoff (de-risking) is achieved more cheaply here by the parallel-binary scheme above without the cgo overhead.

  Line 105: - Rename: archive `cmd/gitswitch/` (Go) to `legacy/` branch; rename Rust binary `gitswitch-rs` → `gitswitch`; ship as `v0.2.0`.

  Line 302: **The v0.2.0 YAML schema is the boundary.** The Rust binary must read and write `~/.config/gitswitch/config.yaml` byte-identically to the Go binary's v0.2.0 output. Migration from v0.1's `profiles.json` is implemented behavior-for-behavior in `gitswitch-core`:

  Line 327: | Two-binary distribution during transition    | Releases tagged `0.2.0-rcN-rs` ship the Rust binary as `gitswitch-rs`; Go releases continue under `gitswitch` until parity is signed off. Both Homebrew taps and install scripts updated to install both side-by-side; users opt in via `--rs` flag on install script. |

  Line 389: 7. **Install scripts updated** — Homebrew formula, install.sh, all distribution channels swap to Rust binary atomically with the rename.

  Line 394: - Rename Rust binary artifact `gitswitch-rs` → `gitswitch` in `Cargo.toml` + `cargo-dist` config.

  Line 406: - [`docs/rust-cancelled/roadmap-v0.2.0.md`](roadmap-v0.2.0.md) — multi-host identity in Rust (Profile / HostCreds / Pin types, keychain + encrypted-file `SecretStore`, gh-driven import, SSH config block manager). Sister doc to [`docs/roadmap-v0.2.0.md`](../roadmap-v0.2.0.md).

  Line 409: - [`docs/rust-cancelled/roadmap-v0.3.0.md`](roadmap-v0.3.0.md) — deferred items (in-repo YAML config, team preset distribution + signing, own OAuth as primary path, host auto-detection). Sister doc to [`docs/roadmap-v0.3.0.md`](../roadmap-v0.3.0.md).


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-v0.3.0.md:
  Line 5: > The Rust migration plan described in this document has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-v0.2.0-tui.md:
  Line 5: > The Rust migration plan described in this document has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.

  Line 10: > Rust translation of the [Go TUI plan](../roadmap-v0.2.0-tui.md). Information architecture, keybinds, command palette, wizards, three themes — identical. This document locks the Rust crate stack and the Elm-architecture state shape, and shows how the arcade aesthetic is preserved through [tachyonfx](https://github.com/junkdog/tachyonfx) effects layered atop [ratatui](https://ratatui.rs/) widgets.


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-v0.2.0-isolation.md:
  Line 5: > The Rust migration plan described in this document has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.

  Line 10: > Rust translation of [roadmap-v0.2.0-isolation.md](../roadmap-v0.2.0-isolation.md). Mental model and all design decisions are identical; only the implementation idioms change. Read the Go plan first for full rationale.

  Line 20: This plan ships **session isolation** in the Rust rewrite: a session is an isolated *scope* (a shell with its own active-identity env vars), not a locked profile. Inside a session you get **full gitswitch** — switch, pin, add hosts, status — all operating on session state. Built parallel to (not blocking on) the [multi-host work](roadmap-v0.2.0.md); shares the [`gitswitch-secrets`](roadmap-v0.2.0.md#crate-layout) keychain crate, otherwise orthogonal.


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-v0.2.0.md:
  Line 5: > The Rust migration plan described in this document has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.

  Line 12: This is the Rust port of the v0.2.0 multi-host plan. Behavior parity with the Go v0.1 → v0.2 migration path is preserved: same on-disk schema, same command surface, same Vaishnavi multi-host scenario as the acceptance target. The Rust workspace structure follows [roadmap-migration.md](roadmap-migration.md).

  Line 16: v0.2.0 closes the gap by making **host a first-class dimension of identity**, **leveraging `gh auth` as the import source** (rather than building our own OAuth), and adding per-repo diagnostics (`gitswitch status`) plus multi-profile pinning. Team-preset YAML distribution, our own OAuth flow, and in-repo config files are explicitly v0.3.0 scope — see [roadmap-v0.3.0.md](roadmap-v0.3.0.md).

  Line 548: Smoke test: install Go v0.1.x, create two profiles, install Rust v0.2.0 binary, run `gitswitch list` → profiles intact, `profiles.json.v1.bak` present, `config.yaml` has hosts map populated with `github.com`.


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/competitor-analysis.md:
  Line 11: **Distribution picture**: no git identity tool has ever successfully landed in homebrew-core. Every active competitor is tap-only or worse. The homebrew-core slot is open and ownable. crates.io is a graveyard — every Rust git identity crate is dead. gitswitch v0.2.0 would ship to an empty shelf on both.

```

---

## Assistant (Build · Kimi K2.7 Code · 63.8s)

_Thinking:_

I see roadmap-v0.2.0.md line 7 still says "our own OAuth flow" is v0.3.0 scope. I need to update that to v1.0. Let me fix it.

Also project-memory.md line 165 mentions crates.io — I should update that since we're not shipping Rust. And competitor-analysis.md line 11 mentions crates.io — should update.

Let me fix roadmap-v0.2.0.md first.

**Tool: edit**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/roadmap-v0.2.0.md",
  "oldString": "v0.2.0 closes the gap by making **host a first-class dimension of identity**, **leveraging `gh auth` as the import source** (rather than building our own OAuth), and adding per-repo diagnostics (`gitswitch status`) plus multi-profile pinning. Team-preset YAML distribution, our own OAuth flow, and in-repo config files are explicitly v0.3.0 scope — see [roadmap-v0.3.0.md](roadmap-v0.3.0.md).",
  "newString": "v0.2.0 closes the gap by making **host a first-class dimension of identity**, **leveraging `gh auth` as the import source** (rather than building our own OAuth yet), and adding per-repo diagnostics (`gitswitch status`) plus multi-profile pinning. Team-preset YAML distribution and in-repo config files are v0.3.0 scope. A direct GitHub OAuth flow is v1.0 scope — see [roadmap-v1.0.md](roadmap-v1.0.md)."
}
```

**Output:**
```
Edit applied successfully.
```

---

## Assistant (Build · Kimi K2.7 Code · 5.5s)

_Thinking:_

Now update project-memory.md line 165 about crates.io.

**Tool: edit**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/project-memory.md",
  "oldString": "- **crates.io is empty** — every Rust git identity crate is dead. gitswitch v0.2.0 owns the category by default on publish.",
  "newString": "- **Go install is viable** — `go install github.com/aksisonline/gitswitch/cmd/gitswitch@latest` works for Go developers; no additional package manager needed for that audience."
}
```

**Output:**
```
Edit applied successfully.
```

---

## Assistant (Build · Kimi K2.7 Code · 4.0s)

_Thinking:_

Now update competitor-analysis.md line 11. Let me read a bit around it.

**Tool: read**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/competitor-analysis.md",
  "offset": 1,
  "limit": 25
}
```

**Output:**
```
<path>/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/competitor-analysis.md</path>
<type>file</type>
<content>
1: # gitswitch — Competitor Analysis
2: 
3: _Date: May 2026. Based on live GitHub data gathered via gh CLI, Homebrew formula API, crates.io, and npm._
4: 
5: ---
6: 
7: ## Summary verdict
8: 
9: gitswitch is the only tool in this space with a real TUI, shell-hook-driven auto-switching, and gh CLI integration. No competitor combines all three. The closest active threats are `bgreenwell/git-ego` (credential helper model, OS keychain, 102 stars and growing fast) and `letmutex/gitas` (Rust, GitHub OAuth, Homebrew, 138 stars). Neither has a TUI or auto-switching.
10: 
11: **Distribution picture**: no git identity tool has ever successfully landed in homebrew-core. Every active competitor is tap-only or worse. The homebrew-core slot is open and ownable. crates.io is a graveyard — every Rust git identity crate is dead. gitswitch v0.2.0 would ship to an empty shelf on both.
12: 
13: ---
14: 
15: ## Competitor table
16: 
17: | Repo | Stars | Language | Active? | TUI | SSH | Auto-switch | gh sync | Keychain | Install |
18: |---|---|---|---|---|---|---|---|---|---|
19: | **gitswitch** | — | Go→Rust | Yes | Yes (arcade) | Yes | Yes (hook) | Yes | Yes | curl 1-liner |
20: | geongeorge/Git-User-Switch | 648 | JavaScript | Dead 2022 | No | No | No | No | No | npm |
21: | prydonius/karn | 306 | Go | Dead 2022 | No | No | YAML map | No | No | manual |
22: | letmutex/gitas | 138 | Rust | Yes | Basic | No | No | No | Yes | curl/brew |
23: | frace/git-passport | 110 | Python | Dead 2015 | No | No | No | No | No | clone+symlink |
24: | samrocketman/git-identity-manager | 107 | Bash | Yes (2024) | No | Yes | includeIf | No | No | brew/script |
25: | bgreenwell/git-ego | 102 | Go | Yes (hot) | No | Yes | includeIf | No | Yes | go install |

(Showing lines 1-25 of 307. Use offset=26 to continue.)
</content>
```

---

## Assistant (Build · Kimi K2.7 Code · 8.0s)

_Thinking:_

I need to update competitor-analysis.md:
- Line 9: "letmutex/gitas (Rust, GitHub OAuth...)" — still accurate competitor info, no need to change.
- Line 11: remove crates.io reference
- Line 19: change Go→Rust to Go

Let me edit.

**Tool: edit**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/competitor-analysis.md",
  "oldString": "**Distribution picture**: no git identity tool has ever successfully landed in homebrew-core. Every active competitor is tap-only or worse. The homebrew-core slot is open and ownable. crates.io is a graveyard — every Rust git identity crate is dead. gitswitch v0.2.0 would ship to an empty shelf on both.",
  "newString": "**Distribution picture**: no git identity tool has ever successfully landed in homebrew-core. Every active competitor is tap-only or worse. The homebrew-core slot is open and ownable. gitswitch can aim for it once v1.0 is stable and the star threshold is met."
}
```

**Output:**
```
Edit applied successfully.
```

**Tool: edit**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/competitor-analysis.md",
  "oldString": "| **gitswitch** | — | Go→Rust | Yes | Yes (arcade) | Yes | Yes (hook) | Yes | Yes | curl 1-liner |",
  "newString": "| **gitswitch** | — | Go | Yes | Yes (arcade) | Yes | Yes (hook) | Yes | Yes | curl 1-liner |"
}
```

**Output:**
```
Edit applied successfully.
```

---

## Assistant (Build · Kimi K2.7 Code · 5.5s)

_Thinking:_

Let me search one more time for any remaining outdated references.

**Tool: grep**

**Input:**
```json
{
  "pattern": "own OAuth.*v0\\.3\\.0|OAuth.*v0\\.3\\.0|Rust.*v0\\.2\\.0\\+|Rust.*v1\\.0|Rust binary|Rust rewrite|Go→Rust|docs/rust/|rust/roadmap|crates\\.io.*gitswitch",
  "path": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up",
  "include": "*.md"
}
```

**Output:**
```
Found 34 matches
/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/project-memory.md:
  Line 11: > The Rust rewrite plan is cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.

  Line 188: | Rust for v0.2.0+ | CANCELLED | Go remains for all versions; see `docs/rust-cancelled/` |


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/roadmap-v0.2.0.md:
  Line 7: v0.2.0 closes the gap by making **host a first-class dimension of identity**, **leveraging `gh auth` as the import source** (rather than building our own OAuth yet), and adding per-repo diagnostics (`gitswitch status`) plus multi-profile pinning. Team-preset YAML distribution and in-repo config files are v0.3.0 scope. A direct GitHub OAuth flow is v1.0 scope — see [roadmap-v1.0.md](roadmap-v1.0.md).


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/marketing.md:
  Line 53: - **r/rust** — not a primary channel since the Rust rewrite was cancelled. Reconsider only if there is a Rust-specific technical angle in the future.


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/versioning.md:
  Line 4: _Status: Rust migration **cancelled**. v1.0 ships on Go._


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/roadmap-versioning-plan.md:
  Line 24: 3. **Keep the Go codebase.** No Rust rewrite. The time that would have gone to parity/cutover now goes to user-facing features.


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/roadmap-full-timeline.md:
  Line 4: _Status: Rust migration **cancelled**. v1.0 will ship on the Go codebase._

  Line 145: > The Rust rewrite has been cancelled. The risk of binary-compatible feature parity + migration outweighed the marginal gains for a tool that already works.


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/why-rust.md:
  Line 5: > The Rust migration described in this document has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.

  Line 10: _Original decision locked: May 2026. The Go codebase was the shipping v0.1.x tool; Rust was the target for v0.2.0+. See [docs/rust-cancelled/roadmap-migration.md](rust-cancelled/roadmap-migration.md) for the archived migration plan._

  Line 33: **Rust's answer**: No runtime. No GC. Rust compiles to bare machine code. Cold-start for a Rust binary doing minimal work (read env var, print string, exit) is measurable in microseconds. For the shell hook path — which runs `gitswitch record --path .` on every `cd` — this difference is the difference between a transparent tool and an annoying one.

  Line 45: **Rust's answer**: With `lto = true`, `opt-level = "z"` (size-optimize), `strip = true`, and `codegen-units = 1`, a Rust binary for a tool of gitswitch's scope compiles to ~1.5–3MB. Fully self-contained. No dynamic runtime dependencies beyond the OS.

  Line 111: The Go codebase is not a mistake — it got gitswitch to v0.1.x with real users. ~~Rust is the right choice for where the tool is going.~~ Go remains the right choice for v1.0 and beyond.

  Line 119: > The Rust migration plan has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.

  Line 122: ~~gitswitch lives in the terminal, fires on every directory change, handles secrets, and needs a polished animated TUI. Those constraints pick Rust. The migration plan is in [docs/rust/roadmap-migration.md](rust/roadmap-migration.md); the v0.2.0 feature plans are in [docs/rust/](rust/).~~


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-migration.md:
  Line 5: > The Rust migration plan described in this document has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.

  Line 32: Instead: **ship the Rust rewrite as `gitswitch-rs` alongside the Go `gitswitch`** until parity is certified. Same config dir, same shell hook contract, same YAML schema. Users opt-in by aliasing or by installing the rs build through a separate install flag. When the parity corpus passes 100%, the binaries swap names and the Go tree gets archived.

  Line 105: - Rename: archive `cmd/gitswitch/` (Go) to `legacy/` branch; rename Rust binary `gitswitch-rs` → `gitswitch`; ship as `v0.2.0`.

  Line 302: **The v0.2.0 YAML schema is the boundary.** The Rust binary must read and write `~/.config/gitswitch/config.yaml` byte-identically to the Go binary's v0.2.0 output. Migration from v0.1's `profiles.json` is implemented behavior-for-behavior in `gitswitch-core`:

  Line 327: | Two-binary distribution during transition    | Releases tagged `0.2.0-rcN-rs` ship the Rust binary as `gitswitch-rs`; Go releases continue under `gitswitch` until parity is signed off. Both Homebrew taps and install scripts updated to install both side-by-side; users opt in via `--rs` flag on install script. |

  Line 389: 7. **Install scripts updated** — Homebrew formula, install.sh, all distribution channels swap to Rust binary atomically with the rename.

  Line 394: - Rename Rust binary artifact `gitswitch-rs` → `gitswitch` in `Cargo.toml` + `cargo-dist` config.

  Line 409: - [`docs/rust-cancelled/roadmap-v0.3.0.md`](roadmap-v0.3.0.md) — deferred items (in-repo YAML config, team preset distribution + signing, own OAuth as primary path, host auto-detection). Sister doc to [`docs/roadmap-v0.3.0.md`](../roadmap-v0.3.0.md).


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-v0.3.0.md:
  Line 5: > The Rust migration plan described in this document has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-v0.2.0-tui.md:
  Line 5: > The Rust migration plan described in this document has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.

  Line 52: - [`clap`](https://crates.io/crates/clap) — CLI dispatch into the TUI; keeps the `gitswitch` / `gitswitch pacman` / `gitswitch --no-tui` entry points.

  Line 212: [`tui-input`](https://crates.io/crates/tui-input) for single-line, [`tui-textarea`](https://crates.io/crates/tui-textarea) for multi-line. A small `gitswitch-tui::wizard` helper drives linear step pages with per-step validation, back/next/cancel.


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-v0.2.0-isolation.md:
  Line 5: > The Rust migration plan described in this document has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.

  Line 20: This plan ships **session isolation** in the Rust rewrite: a session is an isolated *scope* (a shell with its own active-identity env vars), not a locked profile. Inside a session you get **full gitswitch** — switch, pin, add hosts, status — all operating on session state. Built parallel to (not blocking on) the [multi-host work](roadmap-v0.2.0.md); shares the [`gitswitch-secrets`](roadmap-v0.2.0.md#crate-layout) keychain crate, otherwise orthogonal.


/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/rust-cancelled/roadmap-v0.2.0.md:
  Line 5: > The Rust migration plan described in this document has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.

  Line 16: v0.2.0 closes the gap by making **host a first-class dimension of identity**, **leveraging `gh auth` as the import source** (rather than building our own OAuth), and adding per-repo diagnostics (`gitswitch status`) plus multi-profile pinning. Team-preset YAML distribution, our own OAuth flow, and in-repo config files are explicitly v0.3.0 scope — see [roadmap-v0.3.0.md](roadmap-v0.3.0.md).

  Line 79: Single file at `~/.config/gitswitch/config.yaml` (XDG resolution via [`directories`](https://crates.io/crates/directories) — `ProjectDirs::from("dev", "aksisonline", "gitswitch")`).

  Line 188: - `enrich_account()` issues `GET https://<api>/user` and `GET https://<api>/user/emails` with [`reqwest::blocking::Client`](https://crates.io/crates/reqwest) (`Authorization: Bearer <token>` + `User-Agent: gitswitch/0.2`). Picks the primary verified email or the `<id>+<login>@users.noreply.github.com` form. JSON responses parsed with `serde_json`.

  Line 229: - [`keyring`](https://crates.io/crates/keyring) v3 with `apple-native` (macOS Keychain) and `windows-native` (Credential Manager) features. Linux uses the default `linux-native-sync-persistent` backend (libsecret via D-Bus). Service name: `dev.aksisonline.gitswitch`; account: the `refkey`.

```

---

## Assistant (Build · Kimi K2.7 Code · 9.8s)

_Thinking:_

The remaining matches are mostly in rust-cancelled/ docs (which are archived and don't need further changes) and in why-rust.md (which is the rationale doc, kept historical). The main active docs look good.

However, I notice rust-cancelled/roadmap-v0.3.0.md line 116 and rust-cancelled/roadmap-v0.2.0.md line 16 still say OAuth is v0.3.0. Since these are archived docs, I could leave them, but to be consistent I should update them to note that OAuth moved to v1.0. Actually, the archived docs represent the plan at the time they were written. The cancellation banner at the top says they're historical. It's OK to leave them as-is, but it might be slightly confusing. I think it's fine to leave archived docs untouched except for the cancellation banner.

Now let me update the todo list to completed and do a final git status check.

**Tool: todowrite**

**Input:**
```json
{
  "todos": [
    {
      "content": "Add cancellation markers to Rust docs and rename rust/ folder",
      "status": "completed",
      "priority": "high"
    },
    {
      "content": "Revamp roadmap-full-timeline.md with new v1 plan (Go stable, OAuth device flow, gs alias)",
      "status": "completed",
      "priority": "high"
    },
    {
      "content": "Revamp roadmap-v0.3.0.md to remove OAuth and reflect new sequencing",
      "status": "completed",
      "priority": "high"
    },
    {
      "content": "Check current project version and existing features to ground the new versioning plan",
      "status": "completed",
      "priority": "high"
    },
    {
      "content": "Draft new versioning/roadmap docs for v0.2.x, v0.3.x, v1.0 with PAT mgmt, session isolation, OAuth, gs",
      "status": "completed",
      "priority": "high"
    }
  ]
}
```

**Output:**
```
[
  {
    "content": "Add cancellation markers to Rust docs and rename rust/ folder",
    "status": "completed",
    "priority": "high"
  },
  {
    "content": "Revamp roadmap-full-timeline.md with new v1 plan (Go stable, OAuth device flow, gs alias)",
    "status": "completed",
    "priority": "high"
  },
  {
    "content": "Revamp roadmap-v0.3.0.md to remove OAuth and reflect new sequencing",
    "status": "completed",
    "priority": "high"
  },
  {
    "content": "Check current project version and existing features to ground the new versioning plan",
    "status": "completed",
    "priority": "high"
  },
  {
    "content": "Draft new versioning/roadmap docs for v0.2.x, v0.3.x, v1.0 with PAT mgmt, session isolation, OAuth, gs",
    "status": "completed",
    "priority": "high"
  }
]
```

---

## Assistant (Build · Kimi K2.7 Code · 3.5s)

_Thinking:_

Now let me run git status to see all changes.

**Tool: bash**

**Input:**
```json
{
  "command": "git status --short",
  "description": "Show changed files"
}
```

**Output:**
```
(no output)
```

---

## Assistant (Build · Kimi K2.7 Code · 3.8s)

_Thinking:_

Hmm, git status shows no output? That's strange. Maybe the docs folder is gitignored or something. Let me check git status with more detail.

**Tool: bash**

**Input:**
```json
{
  "command": "git status && git diff --stat",
  "description": "Check full git status and diff"
}
```

**Output:**
```
On branch main
Your branch is up to date with 'origin/main'.

nothing to commit, working tree clean

```

---

## Assistant (Build · Kimi K2.7 Code · 5.9s)

_Thinking:_

Interesting. The working tree is clean. That means either:
1. The docs folder is gitignored
2. The files weren't actually changed (unlikely since edit tool said success)
3. The files are outside the git repo

Wait, the working directory is `/Users/aksisonline/Documents/Projects/Personal/git-switcher` and the files are in `/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/`. They should be in the repo.

But git status shows nothing. Maybe the `.never_gonna_give_you_up` folder is gitignored? The name starts with a dot, which might make it ignored by default. Or there's a gitignore entry.

Let me check gitignore and also verify the files exist with the changes.

**Tool: bash**

**Input:**
```json
{
  "command": "git check-ignore -v \"docs/.never_gonna_give_you_up\" \"docs/.never_gonna_give_you_up/roadmap-full-timeline.md\"",
  "description": "Check if docs are gitignored"
}
```

**Output:**
```
.gitignore:17:docs/.never_gonna_give_you_up/	docs/.never_gonna_give_you_up
.gitignore:17:docs/.never_gonna_give_you_up/	docs/.never_gonna_give_you_up/roadmap-full-timeline.md

```

**Tool: bash**

**Input:**
```json
{
  "command": "ls -la docs/ && ls -la docs/.never_gonna_give_you_up/ | head -20",
  "description": "List docs directory contents"
}
```

**Output:**
```
total 32
drwxr-xr-x@  6 aksisonline  staff   192 Jun  2 18:58 .
drwxr-xr-x@ 17 aksisonline  staff   544 Jun 22 03:51 ..
drwxr-xr-x@ 21 aksisonline  staff   672 Jun 22 04:04 .never_gonna_give_you_up
-rw-r--r--@  1 aksisonline  staff  8464 May 12 00:26 cli.md
-rw-r--r--@  1 aksisonline  staff  3413 May 29 14:46 privacy-policy.md
drwxr-xr-x@  9 aksisonline  staff   288 Jun  2 19:36 public
total 432
drwxr-xr-x@ 21 aksisonline  staff    672 Jun 22 04:04 .
drwxr-xr-x@  6 aksisonline  staff    192 Jun  2 18:58 ..
-rw-r--r--@  1 aksisonline  staff   4992 May 28 07:00 arcade-philosophy.md
-rw-r--r--@  1 aksisonline  staff  17938 Jun 22 04:07 competitor-analysis.md
-rw-r--r--@  1 aksisonline  staff   7245 May 16 16:33 diagnostic-design.md
-rw-r--r--@  1 aksisonline  staff  20440 May 30 01:37 explainer.md
-rw-r--r--@  1 aksisonline  staff   4393 May 28 05:56 gitswitch-v2-discussion.md
-rw-r--r--@  1 aksisonline  staff   9744 Jun 22 04:05 marketing.md
-rw-r--r--@  1 aksisonline  staff  11192 May 30 11:24 monetization.md
-rw-r--r--@  1 aksisonline  staff   8909 May 30 01:36 product-gaps-and-philosophy.md
-rw-r--r--@  1 aksisonline  staff  11050 Jun 22 04:07 project-memory.md
-rw-r--r--@  1 aksisonline  staff  10555 Jun 22 04:03 roadmap-full-timeline.md
-rw-r--r--@  1 aksisonline  staff  15662 May 30 01:38 roadmap-v0.2.0-isolation.md
-rw-r--r--@  1 aksisonline  staff  19607 May 30 01:38 roadmap-v0.2.0-tui.md
-rw-r--r--@  1 aksisonline  staff  15002 Jun 22 04:07 roadmap-v0.2.0.md
-rw-r--r--@  1 aksisonline  staff   7332 Jun 22 04:03 roadmap-v0.3.0.md
-rw-r--r--@  1 aksisonline  staff   6504 Jun 22 04:04 roadmap-v1.0.md
-rw-r--r--@  1 aksisonline  staff   3127 Jun 22 04:04 roadmap-versioning-plan.md
drwxr-xr-x@  7 aksisonline  staff    224 May 30 01:39 rust-cancelled

```

---

## Assistant (Build · Kimi K2.7 Code · 9.4s)

_Thinking:_

The docs/.never_gonna_give_you_up/ folder is gitignored. That's why git status shows nothing. The changes are saved to the filesystem but not tracked by git. That's fine — the user can decide whether to track them or not. I should mention this.

Actually, wait — if the folder is gitignored, then the original files were also not tracked. The user probably knows this. But I should mention that the changes are on disk and the folder is gitignored.

Let me verify the changes are actually on disk by reading a few key sections.

**Tool: read**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/roadmap-versioning-plan.md",
  "limit": 40
}
```

**Output:**
```
<path>/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/roadmap-versioning-plan.md</path>
<type>file</type>
<content>
1: # gitswitch Versioning Plan — From v0.1.22 to v1.0
2: 
3: _Last updated: June 2026_
4: 
5: ## Where we are
6: 
7: - Current release: **v0.1.22**.
8: - Language: **Go**.
9: - Rust migration: **cancelled**.
10: - Storage: `~/.gitswitch/profiles.json`.
11: 
12: ## Where we are going
13: 
14: | Version | Theme | Big-ticket items | Why this grouping |
15: |---|---|---|---|
16: | **v0.2.0** | Identity + isolation | Multi-host profiles, keychain PATs, session isolation, TUI refresh, `status`/`doctor`, shell tab completion | These are the infrastructure pieces that change how profiles work. They ship together because they all touch the same storage model and shell integration. |
17: | **v0.3.0** | Repo + team | In-repo `.gitswitch.yaml`, team presets, host auto-detection, remote retargeting, remote-URL routing, `gitswitch run` | These features assume v0.2.0's multi-host model is stable. They make gitswitch work across teams and machines. |
18: | **v1.0** | Production auth | Direct GitHub OAuth device flow, `gs` alias, HTTPS credential helper, semver + schema stability | 1.0 is the "gitswitch owns auth end-to-end" release. No new core data model — just hardened, complete coverage. |
19: 
20: ## What this means for sequencing
21: 
22: 1. **Ship v0.2.0 before thinking about v0.3.0.** v0.2.0 is a breaking config migration (JSON → YAML, single-host → multi-host). Everything after it builds on the YAML schema.
23: 2. **Do not put OAuth in v0.2.0.** v0.2.0 already has enough load-bearing changes. OAuth is a v1.0 requirement because it completes the auth story and lets gitswitch stand alone without `gh`.
24: 3. **Keep the Go codebase.** No Rust rewrite. The time that would have gone to parity/cutover now goes to user-facing features.
25: 
26: ## What the user sees
27: 
28: ### v0.1.22 → v0.2.0
29: ```
30: gitswitch switch corp          # now sets SSH + gh + git for ALL hosts in the profile
31: gitswitch shell corp           # isolated terminal session
32: gitswitch status               # "why is my push failing?"
33: gitswitch import gh            # pull accounts from gh CLI
34: gitswitch pin corp --remote origin
35: ```
36: 
37: ### v0.2.0 → v0.3.0
38: ```
39: gitswitch apply ./team.yaml    # team onboarding
40: gitswitch remote retarget ...  # org migration helper

(Showing lines 1-40 of 61. Use offset=41 to continue.)
</content>
```

**Tool: read**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/roadmap-v1.0.md",
  "limit": 30
}
```

**Output:**
```
<path>/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/roadmap-v1.0.md</path>
<type>file</type>
<content>
1: # gitswitch v1.0 — Production Launch Roadmap
2: 
3: The 1.0 release marks the production-stable Go binary. There is no language rewrite. The focus is completing the authentication story (so gitswitch no longer depends on `gh` CLI for login), adding a short command alias, and making stability promises.
4: 
5: ---
6: 
7: ## 1. Direct GitHub OAuth (device flow)
8: 
9: ### Goal
10: 
11: Let users authenticate with GitHub directly through gitswitch, without requiring the `gh` CLI. This unblocks:
12: - Users who don't have `gh` installed.
13: - GitHub Enterprise Server (GHES) environments where `gh` won't authenticate.
14: - Headless setups where a browser-driven `gh auth login` is impractical.
15: 
16: ### Flow
17: 
18: Device flow is chosen because it is terminal-native, SSH-friendly, and requires no local HTTP server.
19: 
20: ```
21: $ gitswitch login
22: Open https://github.com/login/device and enter code: ABCD-EFGH
23: Waiting for authorization...
24: ✓ Logged in as aksisonline (github.com)
25: Created profile "aksisonline" with email abhiram@personal.dev
26: Token stored in keychain.
27: ```
28: 
29: ### Command surface
30: 

(Showing lines 1-30 of 195. Use offset=31 to continue.)
</content>
```

---

## Assistant (Build · Kimi K2.7 Code · 6.8s)

_Thinking:_

Good, the changes are on disk. Now I should provide a clear summary to the user of what was done and what the new plan looks like. Since the user asked me to "finish cancelling the plan, and then let's go and write new docs with new plans for these", I should summarize both.

I should mention that the docs folder is gitignored so git status won't show changes.

Let me also verify the rust folder was renamed and check the cancellation banners.

**Tool: bash**

**Input:**
```json
{
  "command": "ls docs/.never_gonna_give_you_up/rust-cancelled/ && head -8 docs/.never_gonna_give_you_up/why-rust.md && head -8 docs/.never_gonna_give_you_up/rust-cancelled/roadmap-migration.md",
  "description": "Verify cancelled folder and banners"
}
```

**Output:**
```
roadmap-migration.md
roadmap-v0.2.0-isolation.md
roadmap-v0.2.0-tui.md
roadmap-v0.2.0.md
roadmap-v0.3.0.md
# Why Rust — gitswitch Technical Rationale

> **CANCELLED — June 2026**
>
> The Rust migration described in this document has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.
> This document is kept for historical reference only.
>
> See the updated roadmap in [`roadmap-full-timeline.md`](roadmap-full-timeline.md) for the current plan.
# gitswitch — Go → Rust Migration Plan

> **CANCELLED — June 2026**
>
> The Rust migration plan described in this document has been cancelled. gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.
> This document is kept for historical reference only.
>
> See the updated roadmap in [`roadmap-full-timeline.md`](../roadmap-full-timeline.md) for the current plan.

```

---

## Assistant (Build · Kimi K2.7 Code · 15.1s)

_Thinking:_

Everything looks good. Now I should present a final summary to the user. I need to:
1. Confirm the Rust cancellation is done (folder renamed, banners added)
2. Summarize the new plan from v0.1.22 → v1.0
3. List the new/updated docs
4. Mention gitignored status

I should keep it concise but complete.

Done. Rust plan is cancelled and the new versioning plan is written.

---

## What I changed

### Rust cancellation
- Renamed `docs/.never_gonna_give_you_up/rust/` → `rust-cancelled/`
- Added **CANCELLED — June 2026** banners to:
  - `why-rust.md`
  - `rust-cancelled/roadmap-migration.md`
  - `rust-cancelled/roadmap-v0.2.0.md`
  - `rust-cancelled/roadmap-v0.2.0-isolation.md`
  - `rust-cancelled/roadmap-v0.2.0-tui.md`
  - `rust-cancelled/roadmap-v0.3.0.md`
- Updated `project-memory.md`, `versioning.md`, `marketing.md`, and `competitor-analysis.md` to remove Rust-era assumptions.

### New / rewritten active docs
| File | What it is now |
|---|---|
| `roadmap-full-timeline.md` | Full product timeline from v0.1.22 → v1.0, Go-only |
| `roadmap-v0.3.0.md` | v0.3.0 with OAuth removed and sequencing updated |
| `roadmap-v1.0.md` | New detailed v1.0 plan: OAuth device flow, `gs` alias, HTTPS credential helper |
| `roadmap-versioning-plan.md` | Concise decision record from v0.1.22 to v1.0 |
| `versioning.md` | Updated version ladder and v1.0 launch criteria |

---

## The new plan, from where you are (v0.1.22)

```
v0.1.22 (now) → v0.2.0 → v0.3.0 → v1.0
```

| Version | Theme | What you get |
|---|---|---|
| **v0.2.0** | Identity + isolation | Multi-host profiles, PAT keychain storage, session isolation (`gitswitch shell`), `status`/`doctor`, TUI refresh, `import gh`, shell tab completion |
| **v0.3.0** | Repo + team | In-repo `.gitswitch.yaml`, team presets (`gitswitch apply`), host auto-detection, remote retargeting, remote-URL routing, `gitswitch run`, profile export/import |
| **v1.0** | Production auth | Direct GitHub OAuth device flow (`gitswitch login`), `gs` alias, HTTPS credential helper, semver + schema stability |

### Key decisions locked in
- **OAuth**: device flow, moved to **v1.0** (not v0.2.0 or v0.3.0)
- **`gs`**: shim binary / argv[0] detection, ships **v1.0**
- **Language**: Go for all versions; Rust rewrite cancelled
- **PAT management + session isolation**: both in **v0.2.0**

---

## Note

The `docs/.never_gonna_give_you_up/` folder is **gitignored**, so `git status` won't show these changes. They're saved on disk. If you want them tracked, you'll need to remove that gitignore line or move the docs out of that folder.

Want me to start implementing any of this, or adjust the plan first?

---

## User

let's switch things, we will add an easier multi account setup flow first paired with GitHub OAuth device flow and gs alias, HTTPS helper and stuff in v0.2.0 , after that we will do multi host profiles and the rest. Easy git setup first priority. so the plan is that gitswitch shall setup both git and gh cli for you with zero effort and easy switching.

---

## Assistant (Build · Kimi K2.7 Code · 79.9s)

_Thinking:_

The user wants to re-prioritize. They want v0.2.0 to focus on easy multi-account setup with GitHub OAuth device flow, `gs` alias, HTTPS credential helper — basically making gitswitch set up both git and gh CLI with zero effort. Then multi-host profiles and the rest come later.

So the new plan should be:
- v0.2.0: Easy multi-account setup + OAuth + `gs` + HTTPS credential helper + maybe improved TUI/flow
- v0.3.0: Multi-host profiles, session isolation, etc.
- v1.0: Stability / semver / production guarantees

Wait, but multi-host profiles and session isolation were the core of v0.2.0. If we move those to v0.3.0, then v0.2.0 becomes "easy setup and auth". That's a significant restructure.

Let me think about what's in each version:

**v0.2.0 — Easy multi-account + auth release**
- Direct GitHub OAuth device flow (`gitswitch login`)
- `gs` short command alias
- HTTPS credential helper
- Improved first-time setup / onboarding flow
- Easy account switching for git + gh CLI
- Maybe keep the TUI refresh? Or move to v0.3.0?
- Pre-commit safety net
- Progressive CLI hints

**v0.3.0 — Multi-host + isolation + TUI refresh**
- Multi-host profiles
- Session isolation
- Multi-pin
- `gitswitch status`/`doctor`
- TUI refresh
- Shell tab completion
- `import gh`, `host *`

**v1.0 — Stable production**
- Semver commitment
- Stable config format
- Full distribution (Homebrew core, AUR, Windows installer)
- gitswitch.dev launch

Hmm, but this changes a lot. The user said "Easy git setup first priority. so the plan is that gitswitch shall setup both git and gh cli for you with zero effort and easy switching." This means v0.2.0 should be about making setup trivial.

But wait — the existing v0.1.x already does git/gh switching. The user wants to add OAuth so it doesn't depend on `gh` CLI for auth, and `gs` alias, and HTTPS helper. So v0.2.0 becomes "auth independence + easy setup".

Multi-host profiles and session isolation are more advanced features. Moving them to v0.3.0 makes sense if the priority is easy setup.

However, there's a dependency issue: multi-host profiles require the YAML config migration and keychain storage. If we move multi-host to v0.3.0, do we still need the config migration and keychain in v0.2.0?

For OAuth and HTTPS credential helper, we need keychain storage (PATs). So yes, keychain storage moves to v0.2.0. But the config migration to YAML could be deferred if we store tokens in keychain and keep profiles.json for now. However, keychain token refs are cleaner in YAML. Hmm.

Actually, the cleanest approach:
- v0.2.0: OAuth + `gs` + HTTPS helper. This requires keychain token storage. We can keep profiles.json and add a keychain token store. Or do the YAML migration early.
- v0.3.0: Multi-host profiles, session isolation, TUI refresh, status/doctor, multi-pin, import gh, host commands.

But if v0.2.0 introduces OAuth tokens, and v0.3.0 changes the profile model, we might need two migrations. Better to do the YAML migration in v0.2.0 along with keychain, but keep the profile model single-host for now.

Actually, the user wants "easy multi account setup". This could mean:
1. User runs `gitswitch login`
2. gitswitch creates profiles automatically from GitHub accounts
3. User can switch with `gs switch work`
4. `gs` alias
5. HTTPS helper so HTTPS clones work

This doesn't necessarily require multi-host. A profile can still be single-host (github.com) for now.

So v0.2.0 plan:
- OAuth device flow (`gitswitch login`) creates profiles automatically
- `gs` alias
- HTTPS credential helper
- Improved onboarding / setup wizard
- Config migration to YAML (needed for stability)
- Keychain storage for PATs
- Pre-commit safety net (could be v0.2.0 or v0.3.0)
- Progressive hints

v0.3.0 plan:
- Multi-host profiles (add `hosts` map)
- Session isolation
- TUI refresh
- `status`/`doctor`
- Multi-pin
- `import gh`, `host *`
- Shell tab completion

Wait, but if v0.2.0 migrates to YAML with single-host profiles, then v0.3.0 changes to multi-host profiles, that's two schema changes. Could we keep profiles.json for v0.2.0 and migrate to YAML with multi-host in v0.3.0? That would mean v0.2.0 OAuth tokens need to be stored somewhere. Keychain refs could be stored in profiles.json.

Actually, maybe it's simpler: v0.2.0 does the YAML migration but keeps the profile model mostly the same as v0.1.x (single host). The token_ref field is added. Then v0.3.0 adds the `hosts` map.

Let me think about the profile model evolution:

v0.1.x:
```json
{
  "Name": "...",
  "Email": "...",
  "SSHKey": "...",
  "GHUser": "...",
  "Active": true
}
```

v0.2.0 proposal:
```yaml
version: 2
profiles:
  - nickname: personal
    user_name: Abhiram Kanna
    email: abhiram@personal.dev
    ssh_key: ~/.ssh/id_personal
    gh_user: aksisonline
    token_ref: gitswitch:personal:github.com
    active: true
```

v0.3.0 proposal:
```yaml
version: 3
profiles:
  - nickname: personal
    user_name: Abhiram Kanna
    email: abhiram@personal.dev
    sign_key: ...
    hosts:
      github.com:
        ssh_key: ~/.ssh/id_personal
        gh_user: aksisonline
        token_ref: gitswitch:personal:github.com
```

Two migrations is annoying. Maybe do one migration to v0.2.0 YAML with single-host + token_ref, then v0.3.0 adds hosts map.

Or, do the full multi-host migration in v0.2.0 but only support a single host per profile initially. The data model supports multi-host but the UI/commands don't expose it until v0.3.0. That way there's only one migration.

Actually, the cleanest: v0.2.0 does full multi-host YAML schema, but v0.3.0 adds the commands and UX for managing multiple hosts. v0.2.0 profiles have `hosts` with one entry (github.com). This is forward-compatible.

But the user said "after that we will do multi host profiles and the rest". So multi-host is explicitly v0.3.0. But the data model can still be forward-compatible.

OK, let me restructure the docs:

**v0.2.0 — Easy multi-account + auth (Go)**
- Direct GitHub OAuth device flow
- `gs` alias
- HTTPS credential helper
- `gitswitch login` creates profiles automatically
- Zero-effort first-time setup wizard
- Keychain PAT storage
- Config migration to YAML (forward-compatible with multi-host)
- Pre-commit safety net
- Progressive CLI hints
- Maybe shell tab completion?

**v0.3.0 — Multi-host + power user (Go)**
- Multi-host profiles (exposing the hosts map)
- Session isolation
- TUI refresh with 5 screens
- `status`/`doctor`
- Multi-pin
- `import gh`, `host *`
- Shell tab completion
- Remote retargeting, host auto-detection, etc.

**v1.0 — Stable production**
- Semver
- Team presets
- In-repo `.gitswitch.yaml`
- Remote-URL-aware routing
- Stable distribution

Wait, the user also mentioned "after that we will do multi host profiles and the rest". The "rest" includes session isolation, TUI refresh, etc. So v0.3.0 is multi-host + power features. v1.0 is stable + team features? Or maybe v1.0 is just stability.

I think the cleanest split:
- v0.2.0: Easy setup + OAuth + gs + HTTPS helper (the "anyone can install and use it" release)
- v0.3.0: Multi-host + session isolation + TUI refresh + status/doctor (the "power user / enterprise" release)
- v1.0: Stability + team features + production polish

Actually, team features (in-repo YAML, team presets) might fit better in v0.3.0 or v1.0. But since the user said "multi host profiles and the rest" after v0.2.0, I'll put multi-host, session isolation, TUI, status/doctor in v0.3.0. Then v1.0 is stability + team/enterprise features.

Hmm, but team presets and in-repo YAML depend on the multi-host model. They could go in v0.3.0 or v1.0. Let me put them in v1.0 since v0.3.0 is already big.

Actually, let me reconsider. The original plan had:
- v0.2.0: multi-host, session isolation, TUI, status/doctor
- v0.3.0: in-repo YAML, team presets, host auto-detection, remote retargeting, OAuth
- v1.0: Rust default, HTTPS helper, stability

New plan:
- v0.2.0: OAuth, gs, HTTPS helper, easy setup, keychain, YAML migration (single-host forward-compatible)
- v0.3.0: multi-host profiles, session isolation, TUI refresh, status/doctor, multi-pin, import gh, host commands, shell tab completion, host auto-detection, remote retargeting, remote-URL routing
- v1.0: semver, in-repo YAML, team presets, profile export/import, gitswitch.dev launch, distribution

This makes v0.3.0 big, but it's all the power-user/enterprise features together. v1.0 is "stable + team-scale".

Alternatively:
- v0.2.0: OAuth, gs, HTTPS helper, easy setup, keychain, YAML migration, TUI refresh (because TUI refresh improves the easy setup experience)
- v0.3.0: multi-host, session isolation, status/doctor, multi-pin, import gh, host commands, shell tab completion
- v1.0: stability, in-repo YAML, team presets, remote retargeting, remote-URL routing

Hmm, TUI refresh could go either way. If the priority is "easy setup", a refreshed TUI helps onboarding. But TUI refresh also exposes multi-host UI. Maybe keep TUI refresh minimal in v0.2.0 (just better onboarding screens) and full 5-screen refresh in v0.3.0.

I think the best approach is to write the new plan and present it to the user for confirmation before editing files again. But we're in build mode now, so I can make changes. However, given this is another big restructure, I should present the plan clearly and ask if that's what they want, then implement.

Actually, the user said "let's switch things" which is a directive. They want me to make the changes. But I should still confirm the specific split to avoid another rewrite. Let me present a concise proposed restructure and ask for confirmation before editing.

Wait, the system reminder says I'm in build mode now, but that doesn't mean I shouldn't confirm major decisions. Given the magnitude, I'll present the new structure and ask "Is this the split you want?" before editing files.

Actually, looking at the user's message: "let's switch things, we will add an easier multi account setup flow first paired with GitHub OAuth device flow and gs alias, HTTPS helper and stuff in v0.2.0 , after that we will do multi host profiles and the rest. Easy git setup first priority. so the plan is that gitswitch shall setup both git and gh cli for you with zero effort and easy switching."

This is clear enough. They want:
- v0.2.0: easy multi-account setup + OAuth + gs + HTTPS helper + "stuff"
- After that: multi-host profiles + the rest

I can implement this. But I want to be careful about what's in v0.3.0 vs v1.0. I'll make a reasonable split and present it.

Let me create the new docs. I'll need to update:
1. roadmap-full-timeline.md
2. roadmap-versioning-plan.md
3. versioning.md
4. Create a new v0.2.0 doc focused on easy setup + OAuth + gs + HTTPS
5. Move multi-host/session isolation back to v0.3.0
6. Update roadmap-v0.3.0.md accordingly
7. Update roadmap-v1.0.md
8. Update project-memory.md's roadmap table

Actually, I should create a new `roadmap-v0.2.0.md` that replaces the old one. The old one was about multi-host. I'll rewrite it.

Also, the old `roadmap-v0.2.0-isolation.md` and `roadmap-v0.2.0-tui.md` should be moved to v0.3.0 or kept but their filenames changed. Since these are planning docs, I can rewrite them.

Let me think about the cleanest doc structure:
- `roadmap-v0.2.0.md` — Easy multi-account setup + OAuth + gs + HTTPS helper
- `roadmap-v0.3.0.md` — Multi-host + session isolation + TUI refresh + status/doctor + power features
- `roadmap-v1.0.md` — Stability + team features
- `roadmap-versioning-plan.md` — Updated overview
- `roadmap-full-timeline.md` — Updated timeline
- `versioning.md` — Updated version ladder

I should also probably rename or archive the old `roadmap-v0.2.0-isolation.md` and `roadmap-v0.2.0-tui.md` since they're now v0.3.0 docs. Or I can leave them and update their content. Since they're not in the rust-cancelled folder, they're active docs. I should update them.

Actually, those docs are detailed implementation plans. I can:
1. Rename `roadmap-v0.2.0-isolation.md` to `roadmap-v0.3.0-isolation.md` and update internal references
2. Rename `roadmap-v0.2.0-tui.md` to `roadmap-v0.3.0-tui.md` and update references

Or just rewrite them in place with updated headers. The filenames matter less than the content, but consistency helps.

Given the user wants to "switch things", I'll do a clean restructure:
- Rewrite `roadmap-v0.2.0.md` for the new v0.2.0 focus
- Rename `roadmap-v0.2.0-isolation.md` → `roadmap-v0.3.0-isolation.md` (update content)
- Rename `roadmap-v0.2.0-tui.md` → `roadmap-v0.3.0-tui.md` (update content)
- Rewrite `roadmap-v0.3.0.md` to be the power-user release
- Update `roadmap-v1.0.md`
- Update `roadmap-full-timeline.md`, `roadmap-versioning-plan.md`, `versioning.md`
- Update `project-memory.md`

That's a lot. Let me start by presenting the proposed structure, then implement if the user confirms. Actually, since we're in build mode and the user gave a clear directive, I'll implement. But I'll present the structure first in my response and say "updating docs now".

Wait, I should use the question tool? No, the user gave a clear directive. I'll proceed with the restructure but keep the user informed.

Let me start by reading the current state of docs and then rewrite. I've already read most. I need to read roadmap-v0.2.0-isolation.md and roadmap-v0.2.0-tui.md (already read earlier).

Actually, I already read them in previous steps. Let me proceed.

I'll start with the timeline and versioning plan, then write the new v0.2.0 doc, then update v0.3.0, then rename the isolation/tui docs.

**Tool: todowrite**

**Input:**
```json
{
  "todos": [
    {
      "content": "Restructure roadmap: v0.2.0 = easy setup + OAuth + gs + HTTPS helper",
      "status": "in_progress",
      "priority": "high"
    },
    {
      "content": "Move multi-host/session/TUI/status to v0.3.0",
      "status": "pending",
      "priority": "high"
    },
    {
      "content": "Update v1.0 to stable + team features",
      "status": "pending",
      "priority": "high"
    },
    {
      "content": "Update cross-references in project-memory, versioning, full timeline",
      "status": "pending",
      "priority": "high"
    },
    {
      "content": "Rename isolation/TUI docs to v0.3.0 and update internal refs",
      "status": "pending",
      "priority": "medium"
    }
  ]
}
```

**Output:**
```
[
  {
    "content": "Restructure roadmap: v0.2.0 = easy setup + OAuth + gs + HTTPS helper",
    "status": "in_progress",
    "priority": "high"
  },
  {
    "content": "Move multi-host/session/TUI/status to v0.3.0",
    "status": "pending",
    "priority": "high"
  },
  {
    "content": "Update v1.0 to stable + team features",
    "status": "pending",
    "priority": "high"
  },
  {
    "content": "Update cross-references in project-memory, versioning, full timeline",
    "status": "pending",
    "priority": "high"
  },
  {
    "content": "Rename isolation/TUI docs to v0.3.0 and update internal refs",
    "status": "pending",
    "priority": "medium"
  }
]
```

---

## Assistant (Build · Kimi K2.7 Code · 39.6s)

_Thinking:_

Let me rewrite the core docs. I'll start with roadmap-full-timeline.md since it's the master doc.

**Tool: write**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/roadmap-full-timeline.md",
  "content": "# gitswitch — Full Product Timeline\n\n_Last updated: June 2026_  \n_Status: Rust migration **cancelled**. v1.0 will ship on the Go codebase._\n\n---\n\n## v0.1.x — Current (Go, shipping now)\n\nCurrent release: **v0.1.22**. The foundation. Everything in this version works today.\n\n**Identity management**\n- Named profiles: name, email, SSH key, GPG/signing key, gh_user (single host: github.com)\n- Switch global `~/.gitconfig` identity in one command\n- `gitswitch current` — see who you are right now\n\n**Shell integration**\n- Shell hook fires on `cd`, logs repo + identity usage in background\n- Learns per-repo patterns; suggests the right profile after 3+ uses\n- Shell prompt segment: shows active profile (or warning if config doesn't match any profile)\n- Supports zsh, bash, fish, Oh My Zsh\n\n**TUI**\n- Full-screen terminal UI launched by `gitswitch` with no arguments\n- Three themes: arcade (loud, retro, animated), modern (clean), minimal (monochrome, SSH-safe)\n- Pacman easter egg: hidden game using profile list as the level\n\n**Per-repo pinning**\n- `gitswitch pin <profile>` — this repo always uses this profile\n- Automatic prompt when entering a pinned repo with wrong identity active\n- `gitswitch unpin` to remove\n\n**gh CLI sync**\n- Writes the correct account to `~/.config/gh/hosts.yml` on every switch\n- git and `gh` can never drift apart\n\n**Distribution**\n- Single static binary, no runtime dependencies\n- `curl | sh` one-liner installer\n- Homebrew tap\n- Storage: `~/.gitswitch/profiles.json`\n\n**Utility commands**\n- `gitswitch doctor` — health check on all profiles\n- `gitswitch upgrade` — self-update from GitHub releases\n- `gitswitch hook-check` — verify shell hook is installed and current\n- `gitswitch install` / `gitswitch uninstall`\n- `gitswitch version`\n\n---\n\n## v0.2.0 — Easy multi-account + auth release (Go)\n\nThe priority for v0.2.0 is **zero-effort setup and switching**. A developer installs gitswitch, runs one command to log in with GitHub, and immediately has working git + gh identity. No manual profile creation, no `gh` CLI required, no SSH knowledge needed.\n\n### Direct GitHub OAuth (device flow)\n\n`gitswitch login` authenticates directly with GitHub and creates a profile automatically.\n\n**Flow**\n1. User runs `gitswitch login [--host HOST]`.\n2. Terminal prints a device code and `github.com/login/device` URL.\n3. User authorizes in the browser.\n4. gitswitch receives the PAT, stores it in the OS keychain, fetches `/user` + `/user/emails`, and creates a profile.\n\n**Why device flow**\n- No local HTTP server.\n- Works over SSH and in headless environments.\n- Matches the terminal-native feel of the tool.\n\n**Scopes**: `repo`, `read:user`, `user:email`, `gist`, `workflow` (matches `gh` defaults).\n\n### `gs` short command alias\n\n`gs` is equivalent to `gitswitch`. Implemented via a shim binary / symlink / argv[0] detection so it works in scripts and across shells immediately.\n\n```sh\ngs login\ngs switch work\ngs status\n```\n\n### HTTPS credential helper\n\n`gitswitch` acts as a git `credential.helper` so HTTPS clones and pushes work out of the box:\n\n```sh\ngit config --global credential.helper gitswitch\n```\n\nWhen git asks for credentials, gitswitch returns the PAT for the active profile and host. This closes the gap for beginners who clone via HTTPS (GitHub's default).\n\n### Zero-effort first-time setup\n\nRunning `gitswitch` with no profiles triggers an onboarding wizard:\n\n1. \"Log in with GitHub?\" → `gitswitch login`\n2. \"Add another account?\" → repeat login\n3. \"Which SSH key?\" → detect `~/.ssh/id_*` or skip\n4. Done — git + gh are configured.\n\n`gitswitch init` becomes the explicit version of this flow.\n\n### Secure PAT storage\n\n- Tokens stored in OS keychain: macOS Keychain, Windows Credential Manager, Linux libsecret.\n- Encrypted-file fallback for headless/CI environments.\n- Token ref convention: `gitswitch:<nickname>:<hostname>`.\n- Never written to disk in plaintext.\n\n### Config migration to YAML\n\n- Move from `~/.gitswitch/profiles.json` to `~/.config/gitswitch/config.yaml`.\n- Schema is forward-compatible with multi-host profiles (v0.3.0).\n- Old file backed up to `profiles.json.v1.bak`.\n\n### Pre-commit safety net\n\nAuto-installed by `gitswitch install`:\n- Warns before a wrong-identity commit: \"you're about to commit as personal@gmail.com — this repo is pinned to work\"\n- User can override with `git commit --no-verify`\n- Can opt out via config\n\n### Progressive CLI hints\n\n- After `gitswitch switch` → \"tip: run `gitswitch pin` to make this automatic\"\n- After `gitswitch pin` → \"run `gitswitch status` to see the full picture\" (status ships in v0.3.0)\n- After `gitswitch login` → \"run `gs switch <profile>` to change identity\"\n\n### What v0.2.0 does NOT include\n\nMulti-host profiles, session isolation, the 5-screen TUI refresh, `gitswitch status`, `gitswitch doctor`, and `import gh` are intentionally deferred to v0.3.0 so v0.2.0 can ship the easy-setup story quickly.\n\n---\n\n## Rust migration — CANCELLED\n\n> **CANCELLED — June 2026**\n>\n> The Rust rewrite has been cancelled. The risk of binary-compatible feature parity + migration outweighed the marginal gains for a tool that already works.\n> gitswitch will remain a Go project for the foreseeable future, including the v1.0 launch.\n>\n> Archived docs are in [`rust-cancelled/`](rust-cancelled/).\n\n---\n\n## v0.3.0 — Multi-host + power user release (Go)\n\nBuilds on v0.2.0's YAML config and keychain storage. Adds the enterprise and AI-agent features.\n\n### Multi-host identity\n\nOne profile covers multiple GitHub hosts with different credentials per host.\n\n```yaml\nnickname: corp\nuser_name: Abhiram Kanna\nemail: abhiram@corp.com\nsign_key: 4A1B2C3D\nhosts:\n  github.com:\n    ssh_key: ~/.ssh/id_personal\n    gh_user: aksisonline\n    token_ref: gitswitch:corp:github.com\n  github.corp.com:\n    ssh_key: ~/.ssh/id_corp\n    gh_user: aksisonline-pp\n    token_ref: gitswitch:corp:github.corp.com\n```\n\n- `gitswitch import gh` — discover accounts from `gh auth status` / `~/.config/gh/hosts.yml`\n- `gitswitch host add <profile> <hostname>` / `host remove` / `host list`\n- Multi-pin: `gitswitch pin corp --remote origin`, `gitswitch pin personal --remote fork`\n\n### Session isolation\n\n- `gitswitch shell [profile]` — per-terminal identity via env vars\n- Parallel AI agents with independent identity scopes\n- `gitswitch env <profile> --format json` for headless/agent use\n- Full gitswitch works inside a session\n\n### Diagnostics\n\n- `gitswitch status` — per-repo narrative diagnostic\n- `gitswitch doctor` — global health check across all profiles × hosts\n\n### TUI refresh\n\n- 5 screens: Profiles, Hosts, Repos, Session, Doctor\n- Command palette (`:` key)\n- Wizards for profile creation and host addition\n- Shell tab completion\n\n### Other v0.3.0 items\n\n- Host auto-detection from `git remote -v`\n- Remote retargeting for org migrations\n- Remote-URL-aware pinning\n- Ephemeral one-shot run: `gitswitch run <profile> -- <cmd>`\n- Profile export/import\n- Windows native support improvements\n\n### Likely sequencing\n\n1. Multi-host profile model + migration\n2. `import gh` + `host *` commands\n3. Session isolation\n4. `status` + `doctor`\n5. TUI refresh\n6. Host auto-detection + remote retargeting + remote-URL routing\n\n---\n\n## v1.0 — Production milestone (Go, stable)\n\nThe stable release. No new core data model. Focus is team-scale features and stability promises.\n\n**Team preset distribution**\n- `gitswitch apply <file-or-url>` — bootstraps all profiles from a shared config\n- Signed with minisign so an attacker can't push a malicious preset\n\n**In-repo `.gitswitch.yaml`**\n- Per-repo config file optionally checked into the repo root\n- Declares required profiles and identity policies\n- Precedence: session env > local pin > repo YAML > global switch\n\n**git includeIf support (optional, not default)**\n- Fires identity at git-command time via git's native config loading\n- Alternative for users who want zero shell hooks\n\n**Stability guarantees**\n- Semver commitment: no breaking CLI changes without a major version bump\n- Stable config format: `~/.config/gitswitch/config.yaml` schema versioned\n- Stable keychain key convention: `gitswitch:<nickname>:<hostname>`\n\n**Distribution at v1.0**\n- curl one-liner (primary)\n- Homebrew (official tap, homebrew-core submission ready)\n- AUR (Arch Linux)\n- Windows: native installer (NSIS or WiX)\n- gitswitch.dev live with full docs, blog, and privacy policy\n\n---\n\n## What does NOT change across any version\n\n- Every feature, including all enterprise toolsets, stays free forever\n- No cloud component, no telemetry, no data leaves the machine\n- No feature flags, no paywalls, no paid themes\n- The arcade personality and easter eggs are never gatekept\n- Open source, inspectable at any time\n"
}
```

**Output:**
```
Wrote file successfully.
```

---

## Assistant (Build · Kimi K2.7 Code · 34.7s)

_Thinking:_

Now write the new v0.2.0 doc.

**Tool: write**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/roadmap-v0.2.0.md",
  "content": "# gitswitch v0.2.0 — Easy Multi-Account + Auth Release\n\n## Goal\n\nMake gitswitch the fastest way to set up multiple GitHub accounts on a machine. A user should install, run `gs login`, and have working git + gh identity in under 60 seconds — no manual profile creation, no SSH expertise, no `gh` CLI required.\n\n---\n\n## 1. Direct GitHub OAuth (device flow)\n\n### Why\n\nv0.1.x requires `gh` CLI to already be authenticated, or the user must create profiles manually. This is friction for:\n- Beginners who don't have `gh` installed.\n- GHES users where `gh` won't authenticate.\n- Headless/SSH setups where browser-driven `gh auth login` is awkward.\n\n### Command\n\n```sh\ngitswitch login [--host HOST] [--client-id ID] [--profile NAME]\n```\n\n### Flow\n\n```\n$ gs login\nOpen https://github.com/login/device and enter code: ABCD-EFGH\nWaiting for authorization...\n✓ Logged in as aksisonline (github.com)\nCreated profile \"aksisonline\" with email abhiram@personal.dev\nToken stored in keychain.\n```\n\n### Implementation\n\n1. Start GitHub device flow: POST to `https://<host>/login/device/code`.\n2. Display `verification_uri` and `user_code`.\n3. Poll `https://<host>/login/oauth/access_token` until authorized or expired.\n4. Fetch `/user` and `/user/emails` with the token.\n5. Create profile:\n   - `nickname` = login\n   - `user_name` = name from `/user`\n   - `email` = primary verified email, or noreply fallback\n   - `gh_user` = login\n   - `token_ref` = `gitswitch:<login>:<hostname>`\n6. Store token in OS keychain at `token_ref`.\n7. Sync to `~/.config/gh/hosts.yml` if `gh` is installed.\n\n### Scopes\n\n`repo`, `read:user`, `user:email`, `gist`, `workflow`.\n\n### Token storage\n\n- OS keychain primary.\n- Encrypted-file fallback for headless/CI.\n- Ref convention: `gitswitch:<nickname>:<hostname>`.\n\n### Enterprise GHES\n\n- `--host github.corp.com` targets the enterprise host.\n- `--client-id` lets the org supply their own OAuth app.\n- API host derivation: `github.com` → `api.github.com`; `github.example.com` → `github.example.com/api/v3`.\n\n---\n\n## 2. `gs` short command alias\n\n### Goal\n\nA faster invocation name that matches `gh`, `glab`, etc.\n\n### Implementation\n\nShip `gs` as a shim/symlink to the same binary. The binary detects `filepath.Base(os.Args[0]) == \"gs\"` and behaves identically.\n\n### Install behavior\n\n- `curl | sh` installer creates both `gitswitch` and `gs`.\n- `gitswitch install` does not need to modify shell rc files for this.\n- Works in scripts, subshells, and fresh terminals immediately.\n\n### Conflict handling\n\n- If `gs` already exists on PATH (e.g. Ghostscript), warn and offer `--no-gs`.\n\n---\n\n## 3. HTTPS credential helper\n\n### Goal\n\nMake HTTPS remotes work with gitswitch-managed PATs. Most beginners clone via HTTPS; gitswitch currently only helps SSH users.\n\n### Setup\n\n```sh\ngit config --global credential.helper gitswitch\n```\n\n### Protocol\n\nWhen git needs credentials, it invokes:\n\n```sh\ngitswitch credential fill\n```\n\nstdin receives:\n\n```\nprotocol=https\nhost=github.com\n```\n\nstdout returns:\n\n```\nprotocol=https\nhost=github.com\nusername=aksisonline\npassword=<PAT>\n```\n\n### Security\n\n- Only responds for hosts present in the active profile.\n- Reads PAT from keychain at request time.\n- Never logs tokens.\n\n### Commands\n\n- `gitswitch credential fill` — return credentials for active profile + host\n- `gitswitch credential approve` — no-op (keychain is the cache)\n- `gitswitch credential reject` — optionally delete token and prompt re-login\n\n---\n\n## 4. Zero-effort first-time setup\n\n### Trigger\n\nRunning `gitswitch` or `gs` with no profiles opens an onboarding wizard:\n\n```\n$ gs\nNo profiles found. Let's set up gitswitch.\n\n[1] Log in with GitHub\n[2] Add a profile manually\n[3] Import from gh CLI\n\nChoice: 1\n... OAuth flow ...\n\nAdd another GitHub account? [y/N]: y\n... OAuth flow ...\n\nSet up SSH keys? (optional) [Y/n]: n\n\nDone. Try:\n  gs switch aksisonline\n  gs pin\n```\n\n### `gitswitch init`\n\nThe explicit version of the wizard, non-interactive flags supported:\n\n```sh\ngitswitch init --login --host github.com\n```\n\n### Profile defaults\n\n- First login becomes the active profile.\n- Each login creates a separate profile.\n- Nickname defaults to GitHub login; user can rename later.\n\n---\n\n## 5. Secure PAT storage (`internal/secrets/`)\n\n### Interface\n\n```go\ntype Store interface {\n    Set(ref, token string) error\n    Get(ref string) (string, error)\n    Delete(ref string) error\n    Available() bool\n}\n```\n\n### Backends\n\n- `KeychainStore` — macOS Keychain, libsecret on Linux, Credential Manager on Windows.\n- `EncryptedFileStore` — xchacha20poly1305 + scrypt-derived passphrase, stored at `~/.config/gitswitch/secrets.enc`.\n\n### Selection\n\nRuntime probe via `Available()`. Override with `GITSWITCH_SECRETS_BACKEND=file`.\n\n---\n\n## 6. Config migration to YAML\n\n### New format\n\n`~/.config/gitswitch/config.yaml`:\n\n```yaml\nversion: 2\nprefs:\n  theme: arcade\nprofiles:\n  - nickname: aksisonline\n    user_name: Abhiram Kanna\n    email: abhiram@personal.dev\n    ssh_key: ~/.ssh/id_personal\n    gh_user: aksisonline\n    token_ref: gitswitch:aksisonline:github.com\n    active: true\n  - nickname: work\n    user_name: Abhiram Kanna\n    email: abhiram@corp.com\n    ssh_key: ~/.ssh/id_corp\n    gh_user: aksisonline-pp\n    token_ref: gitswitch:work:github.com\n    active: false\npins:\n  /Users/abhi/code/corp-platform: work\n```\n\n### Migration\n\nOn first load after upgrade:\n1. Read `profiles.json`.\n2. Convert each profile to YAML schema.\n3. Write `config.yaml` (mode 0600).\n4. Rename `profiles.json` → `profiles.json.v1.bak`.\n\n### Forward compatibility\n\nThe v0.2.0 schema uses top-level `ssh_key` / `gh_user` per profile. v0.3.0 will fold these into a `hosts` map. The loader must accept both shapes or migrate cleanly.\n\n---\n\n## 7. Pre-commit safety net\n\n- Auto-installed by `gitswitch install`.\n- Warns before a wrong-identity commit in a pinned repo.\n- Override with `git commit --no-verify`.\n- Opt-out via config.\n\n---\n\n## 8. Progressive CLI hints\n\n- After `gs switch` → \"tip: run `gs pin` to make this automatic\"\n- After `gs pin` → \"run `gs status` to see the full picture\" (status ships in v0.3.0)\n- After `gs login` → \"run `gs switch <profile>` to change identity\"\n\n---\n\n## Command surface\n\n### New\n- `gitswitch login` — GitHub OAuth device flow\n- `gitswitch credential fill/approve/reject` — git credential helper protocol\n- `gs` — short alias for `gitswitch`\n\n### Repurposed\n- `gitswitch init` — interactive onboarding wizard\n\n### Unchanged\n- `switch`, `list`, `remove`, `current`, `pin`, `unpin`, `record`, `recommend`, `install`, `uninstall`, `version`, `upgrade`, `hook-check`, `claude`, `pacman`, `doctor`\n\n`doctor` continues its current v0.1.x health check; the expanded multi-host diagnostics ship in v0.3.0.\n\n---\n\n## Verification plan\n\n### Unit\n- OAuth device flow state machine against a mock GitHub device endpoint.\n- Keychain store round-trip and fallback selection.\n- Credential helper protocol parsing.\n- YAML migration from v1 JSON.\n\n### Integration\n- `gs login` creates a profile and stores a token in keychain.\n- `gitswitch credential fill` returns correct credentials for active profile.\n- HTTPS clone against a test repo uses gitswitch-provided PAT.\n- `gitswitch init` wizard creates expected profiles from mocked OAuth.\n\n### Manual\n- Install on a fresh machine with no `gh` CLI; authenticate via OAuth; push to HTTPS remote.\n- Verify `gs` works in bash, zsh, fish.\n- Verify token is not visible in any file.\n\n---\n\n## Out of scope — v0.3.0\n\nMulti-host profiles, session isolation, TUI refresh, `gitswitch status`, expanded `doctor`, `import gh`, `host *`, multi-pin, host auto-detection, remote retargeting.\n"
}
```

**Output:**
```
Wrote file successfully.
```

---

## Assistant (Build · Kimi K2.7 Code · 32.6s)

_Thinking:_

Now rewrite roadmap-v0.3.0.md to include multi-host, session isolation, TUI, etc.

**Tool: write**

**Input:**
```json
{
  "filePath": "/Users/aksisonline/Documents/Projects/Personal/git-switcher/docs/.never_gonna_give_you_up/roadmap-v0.3.0.md",
  "content": "# gitswitch v0.3.0 — Multi-Host + Power User Release\n\nItems deferred from [v0.2.0 planning](roadmap-v0.2.0.md). v0.2.0 establishes direct GitHub OAuth login, the `gs` alias, HTTPS credential helper, keychain-backed PAT storage, and the YAML config format. v0.3.0 builds on that foundation to add multi-host identity, session isolation, diagnostics, and the refreshed TUI.\n\n---\n\n## 1. Multi-host identity\n\nCurrent gitswitch treats identity as `(name, email, ssh_key, gpg_key, gh_user)` against a single implicit host (`github.com`). Enterprise multi-host environments — finance, healthcare, defense, consulting — need developers to operate against **two GitHub hosts simultaneously** with different PAT tokens, SSH keys, and `gh` contexts.\n\nv0.3.0 makes **host a first-class dimension of identity**.\n\n### Profile model change\n\n```yaml\nnickname: corp\nuser_name: Abhiram Kanna\nemail: abhiram@corp.com\nsign_key: 4A1B2C3D\nhosts:\n  github.com:\n    ssh_key: ~/.ssh/id_personal\n    gh_user: aksisonline\n    token_ref: gitswitch:corp:github.com\n  github.corp.com:\n    ssh_key: ~/.ssh/id_corp\n    gh_user: aksisonline-pp\n    token_ref: gitswitch:corp:github.corp.com\n```\n\nTop-level `ssh_key` / `gh_user` from v0.2.0 collapse into `hosts[\"github.com\"]` during the v2→v3 upgrade.\n\n### What this enables\n\n- One `gitswitch switch corp` sets SSH config, gh auth, and git identity for all hosts simultaneously.\n- Repos with multiple remotes pointing to different hosts work correctly.\n- Multi-pin: `gitswitch pin corp --remote origin`, `gitswitch pin personal --remote fork`.\n\n### New commands\n\n- `gitswitch import gh` — discover all accounts from `gh auth status` / `~/.config/gh/hosts.yml`, create profiles automatically.\n- `gitswitch import git` — import current git config as a new profile.\n- `gitswitch host add <profile> <hostname> [--ssh-key] [--gh-user] [--token]`\n- `gitswitch host remove <profile> <hostname>`\n- `gitswitch host list <profile>`\n\n### SSH config management\n\nIdempotent per-host blocks in `~/.ssh/config`:\n\n```\n# gitswitch:github.corp.com\nHost github.corp.com\n    IdentityFile ~/.ssh/id_corp\n    IdentitiesOnly yes\n```\n\nDrops the global `core.sshCommand` approach. Per-host config lets multi-remote pushes work without re-switching profiles.\n\n### gh CLI sync on switch\n\n`internal/gh/sync.go` updates `~/.config/gh/hosts.yml` per host in the active profile. Best-effort: silent skip if `gh` isn't installed.\n\n---\n\n## 2. Session isolation\n\nEach terminal window holds its own identity via environment variables — parallel sessions, independent identities.\n\n### How it works\n\n- `gitswitch shell [profile]` — opens a subshell with identity locked in env vars.\n- Inside the session: git reads `GIT_AUTHOR_NAME`, `GIT_AUTHOR_EMAIL`, etc. in preference to `~/.gitconfig`.\n- `GH_TOKEN` / `GH_ENTERPRISE_TOKEN` set so `gh` uses the right account.\n- `GITSWITCH_SESSION_ID`, `GITSWITCH_SESSION_PROFILE` for introspection.\n- Exit the shell → env vars disappear → no cleanup needed.\n- Can still `gitswitch switch` within a session to change identity for that session only.\n\n### Why this matters for AI agents\n\nTwo Claude Code agents running in parallel, each inside their own `gitswitch shell`, cannot interfere with each other's identity. Every commit is correctly attributed regardless of what the developer is doing globally.\n\n### Commands\n\n- `gitswitch shell [profile] [--cmd \"...\"] [--no-prompt]`\n- `gitswitch env <profile> [--format posix|fish|json]` — print exports for eval/programmatic use\n- `gitswitch env --clear` — drop session-scoped identity without exiting\n\n---\n\n## 3. Diagnostics\n\n### `gitswitch status` — repo narrative diagnostic\n\nIn a git repo, output a structured explanation:\n\n```\nRepo: ~/code/corp-platform\n  Branch: feature/x  →  origin/feature/x\n\nIdentity (global git config):\n  user.name   = Abhiram Kanna\n  user.email  = abhiram@corp.com\n  signingkey  = 4A...                (gpg, present)\n  Active profile: corp             (matches git config ✓)\n\nPins for this repo:\n  default          → corp\n  remote \"fork\"    → personal\n\nRemotes:\n  origin → git@github.corp.com:team/x.git\n    Host: github.corp.com\n    Profile that owns this host: corp\n    SSH config block: present ✓\n    Identity file: ~/.ssh/id_corp    (exists, mode 600 ✓)\n    gh auth: aksisonline-pp            (synced ✓, token valid)\n  fork → git@github.com:aksisonline/x.git\n    Host: github.com\n    Profile that owns this host: personal\n    SSH config block: present ✓\n    Identity file: ~/.ssh/id_personal      (exists ✓)\n    gh auth: aksisonline                   (synced ✓, token valid)\n\nRecommendations: none.\n```\n\n### `gitswitch doctor` — global health check\n\nWalks every profile × host:\n- keychain reachable?\n- token valid (HTTP 200 on `/user`)?\n- SSH key file exists + permissions correct?\n- `gh hosts.yml` in sync?\n- shell hook version current?\n\n---\n\n## 4. TUI refresh\n\n- 5 screens: Profiles, Hosts, Repos, Session, Doctor\n- Command palette (`:` key) for power users\n- Wizards for profile creation and host addition\n- New keybindings for session and host management\n- Shell tab completion ships in v0.3.0\n\n---\n\n## 5. Repo-level and routing features\n\n### Host auto-detection\n\nOn `cd` into a repo, walk `git remote -v` and auto-recommend the profile whose `Hosts` map covers all remotes.\n\n### Remote retargeting\n\n`gitswitch remote retarget --from github.corp.com --to github.com/neworg` rewrites remotes during org migrations.\n\n### Remote-URL-aware routing\n\nPin by remote URL pattern:\n\n```yaml\nroutes:\n  - pattern: \"github.com/mycompany/*\"\n    profile: work\n  - pattern: \"github.com/aksisonline/*\"\n    profile: personal\n```\n\n### Profile export/import\n\nJSON backup and restore for new machine setup.\n\n### Ephemeral one-shot run\n\n```sh\ngitswitch run <profile> -- <cmd>\n```\n\nRun one command as a different identity without switching.\n\n---\n\n## Likely sequencing for v0.3.0\n\n1. Multi-host profile model + v2→v3 migration\n2. SSH config block manager + gh sync per host\n3. `import gh` + `host *` commands\n4. Session isolation\n5. `status` + `doctor`\n6. TUI refresh\n7. Host auto-detection + remote retargeting + remote-URL routing\n\nOrder may shift based on which use cases surface from 0.2.0 adopters.\n\n---\n\n## Out of scope — v1.0\n\nIn-repo `.gitswitch.yaml`, team preset distribution + signing, semver stability, full distribution expansion.\n"
}
```

**Output:**
```
Wrote file successfully.
```

---

## Assistant (Build · Kimi K2.7 Code · 305.5s)

---

## User

one more thing I want to try and add is an agentic way to setup the tool, where the user can setup gitswitch, and then via our gitswitch, both their git and gh CLI all at the same time, using the agent, our tool can guide the agent, take inputs from the user and set it up for them, this should also be an available flow for super beginners. And then look into the doctor thing that I was trying to build, where gitswitch can act like a doctor to help diagnose broken local git setups, remote links, in-line-push-pull statuses, all of which are toggleable in the settings (new TUI for gitswitch with settings).

So we can be like, setup using Claude or agent, a simple link that the agent can parse, and help the user setup gitswitch through guided steps, the agent acts as the interface, if the user doesn't understand how to set it up.

---

## User

one more thing I want to try and add is an agentic way to setup the tool, where the user can setup gitswitch, and then via our gitswitch, both their git and gh CLI all at the same time, using the agent, our tool can guide the agent, take inputs from the user and set it up for them, this should also be an available flow for super beginners. And then look into the doctor thing that I was trying to build, where gitswitch can act like a doctor to help diagnose broken local git setups, remote links, in-line-push-pull statuses, all of which are toggleable in the settings (new TUI for gitswitch with settings).

So we can be like, setup using Claude or agent, a simple link that the agent can parse, and help the user setup gitswitch through guided steps, the agent acts as the interface, if the user doesn't understand how to set it up.

---

## Assistant (Build · Kimi K2.7 Code · 24.9s)

_Thinking:_

The user wants to add two more things:

1. **Agentic setup flow**: A way for users (especially beginners) to set up gitswitch via an AI agent (Claude or other). The agent parses a simple link or instruction and guides the user through setup, configuring git and gh CLI via gitswitch.

2. **Doctor diagnostics with toggleable settings in a new TUI settings screen**: The doctor should diagnose broken local git setups, remote links, in-line push/pull statuses. All diagnostic checks should be toggleable in settings.

These are significant additions. I need to think about where they fit in the versioning plan.

**Agentic setup flow:**
- This fits well with v0.2.0's "easy setup" theme. It could be a companion to the OAuth + wizard flow.
- It might require a web endpoint or a special `gitswitch setup --agent` mode that emits structured instructions.
- The "simple link that the agent can parse" suggests something like `gitswitch.dev/setup?token=...` or a local command that generates an LLM-readable setup manifest.
- This could be a v0.2.x addition or v0.2.0 feature.

**Doctor + settings TUI:**
- This is more advanced and fits v0.3.0 (which already has `status`/`doctor` and TUI refresh).
- The toggleable settings require a settings screen in the TUI.
- This could expand the doctor scope significantly.

Given the user's "one more thing" framing, I think they want these added to the plan. I should update the docs to reflect:
- Agentic setup as part of v0.2.0 (or maybe v0.2.x)
- Enhanced doctor + settings TUI as part of v0.3.0

Let me update the relevant docs:
1. roadmap-full-timeline.md
2. roadmap-v0.2.0.md (add agentic setup)
3. roadmap-v0.3.0.md (expand doctor + settings TUI)
4. roadmap-versioning-plan.md
5. versioning.md
6. project-memory.md

Also, I should consider whether the existing `claude` subcommand/skill is relevant. There was a `claude` command in v0.1.x and a `gitswitch claude` skill integration mentioned. The agentic setup could leverage this.

Let me also think about the implementation of the "simple link that the agent can parse". A few options:
- A URL like `https://gitswitch.dev/setup` that returns an `llms.txt` or structured setup guide.
- A local command `gitswitch setup --agent` that emits a JSON manifest of the current machine state and steps needed.
- A skill file (`SKILL.md`) that Claude can load, instructing it how to help set up gitswitch.

The user mentioned "a simple link that the agent can parse". This sounds like a web URL with setup instructions. Combined with the existing `//go:embed skill/SKILL.md` in main.go, perhaps the agentic setup uses a skill file.

Actually, since there's already a `claude` subcommand and skill integration, the agentic setup could be:
- `gitswitch setup --agent` prints a URL or a block of text the user can paste into Claude.
- The URL serves an LLM-readable setup guide with interactive steps.
- Claude asks the user questions, runs gitswitch commands, and reports back.

I should design this into the docs.

For the doctor + settings TUI:
- Settings screen in the TUI with toggles for each diagnostic check.
- Doctor checks: local git config health, remote URL validity, SSH config, gh auth, push/pull simulation, merge conflicts, uncommitted changes, etc.
- In-line push/pull status: show in the prompt or TUI whether the current branch is ahead/behind, has uncommitted changes, etc.

Let me update the docs. I'll add these sections.

Also, I still need to rename roadmap-v0.2.0-isolation.md and roadmap-v0.2.0-tui.md to v0.3.0. I should do that too.

Let me proceed with updates.

**Tool: todowrite**

**Input:**
```json
{
  "todos": [
    {
      "content": "Restructure roadmap: v0.2.0 = easy setup + OAuth + gs + HTTPS helper",
      "status": "completed",
      "priority": "high"
    },
    {
      "content": "Move multi-host/session/TUI/status to v0.3.0",
      "status": "completed",
      "priority": "high"
    },
    {
      "content": "Update v1.0 to stable + team features",
      "status": "completed",
      "priority": "high"
    },
    {
      "content": "Update cross-references in project-memory, versioning, full timeline",
      "status": "in_progress",
      "priority": "high"
    },
    {
      "content": "Rename isolation/TUI docs to v0.3.0 and update internal refs",
      "status": "pending",
      "priority": "medium"
    },
    {
      "content": "Add agentic setup flow to v0.2.0 plan",
      "status": "pending",
      "priority": "high"
    },
    {
      "content": "Add doctor + settings TUI details to v0.3.0 plan",
      "status": "pending",
      "priority": "high"
    }
  ]
}
```

**Output:**
```
[
  {
    "content": "Restructure roadmap: v0.2.0 = easy setup + OAuth + gs + HTTPS helper",
    "status": "completed",
    "priority": "high"
  },
  {
    "content": "Move multi-host/session/TUI/status to v0.3.0",
    "status": "completed",
    "priority": "high"
  },
  {
    "content": "Update v1.0 to stable + team features",
    "status": "completed",
    "priority": "high"
  },
  {
    "content": "Update cross-references in project-memory, versioning, full timeline",
    "status": "in_progress",
    "priority": "high"
  },
  {
    "content": "Rename isolation/TUI docs to v0.3.0 and update internal refs",
    "status": "pending",
    "priority": "medium"
  },
  {
    "content": "Add agentic setup flow to v0.2.0 plan",
    "status": "pending",
    "priority": "high"
  },
  {
    "content": "Add doctor + settings TUI details to v0.3.0 plan",
    "status": "pending",
    "priority": "high"
  }
]
```

---

## Assistant (Build · Kimi K2.7 Code · 0.7s)

---

