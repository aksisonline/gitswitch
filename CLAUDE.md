# gitswitch — Claude Instructions

## Release Notes Standard

Every push to `main` that cuts a release MUST update `CHANGELOG.md` with release notes before committing. This is non-negotiable.

### Audience

gitswitch users range from students running their first git command to senior engineers with 10 repos. Release notes must make sense to both. Never assume the reader knows what a "credential helper" or "OAuth flow" is without one plain-English sentence of context.

### Format (required structure per release)

```md
## vX.Y.Z

### What's New          ← only for feat: / minor/major bumps
### Bug Fixes           ← for fix: commits
### Under the Hood      ← for chore:/refactor: — optional, keep brief
```

### Writing rules

- **Lead with the user benefit, not the implementation.** Bad: "Refactored credential store to use atomic writes." Good: "Profiles now save instantly — no more partial-write corruption on slow disks."
- **One sentence per item.** If you need two, the item is two items.
- **Plain English first, technical detail in parentheses.** "OAuth support — connect GitHub without managing tokens (uses device flow)."
- **Beginners need context, pros need precision.** Write the plain version, then add the precise qualifier in parens or after a dash.
- **Beta / canary callouts get their own block** — clear opt-in command, what's in it, nothing speculative.
- **No internal jargon** — no PR numbers, no internal ticket refs, no "as discussed".

### Commit convention (affects auto version bump)

| Prefix | Bump |
|--------|------|
| `feat:` | MINOR (0.X.0) — triggers What's New splash |
| `fix:` | PATCH |
| `chore:` / `docs:` / `refactor:` | PATCH |
| `feat!:` or `BREAKING CHANGE:` | MAJOR |

The CI reads `feat:` commits since the last tag to decide MINOR vs PATCH. If your change ships something users will notice — new command, new UI, new flow — use `feat:`. If it fixes something broken, use `fix:`.

### What's New splash trigger

The in-app splash screen fires when the installed version has a higher MINOR or MAJOR than the user's last-seen version. The splash content is the GitHub release body — which is `CHANGELOG.md` passed to GoReleaser via `--release-notes CHANGELOG.md`.

So: great CHANGELOG = great splash = users understand what changed without Googling.

## CI/CD

Workflow: `.github/workflows/release.yml`  
GoReleaser config: `.goreleaser.yaml`

Push to `main` → CI checks for code changes in `internal/` or `cmd/` → computes next version (MINOR if any `feat:` commit, PATCH otherwise) → pushes tag → GoReleaser builds + publishes GitHub release + updates Homebrew tap.

**Never push tags manually to remote.** CI owns remote tags.

Required secrets: `HOMEBREW_PAT`, `WEBSITE_DISPATCH_PAT` (GITHUB_TOKEN is auto-provided).
