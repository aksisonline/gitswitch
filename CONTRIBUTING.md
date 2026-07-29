# Contributing to gitswitch

Thanks for taking the time to contribute. gitswitch is a small Go CLI/TUI, and
the contribution loop is meant to be just as small.

## Before you start

- **Open an issue first for anything non-trivial.** Bug fixes and small
  polish are fine to send straight as a PR. New commands, new TUI screens, or
  anything that changes existing behavior should get a quick issue first so
  we don't cross wires on direction.
- **Target the `canary` branch, not `main`.** `main` only receives code that's
  already been tested on canary betas. Point your PR at `canary` — if you
  target `main` by mistake, we'll redirect it (or close it and port the fix
  ourselves with credit) rather than merge it out of sequence.

## Dev setup

```bash
git clone https://github.com/aksisonline/gitswitch
cd gitswitch
go build -o gitswitch ./cmd/gitswitch
go test ./...
go vet ./...
```

Go version: see `go.mod`. No other toolchain required.

## Making a change

- **Root-cause fixes over symptom patches.** If a bug shows up in one
  caller, check every other caller of the same function before patching —
  fix it once, where they all route through.
- **Match existing style.** No new abstractions for a single call site, no
  speculative config, no dependency for what a few lines of stdlib can do.
- **UI consistency.** If you change a command's behavior, a config field, or
  a status, grep the TUI views (`internal/tui/`) and CLI help text for stale
  references to it. This codebase has shipped backend features that the TUI
  kept describing as "coming soon" more than once — don't add another one.
- **Tests.** Non-trivial logic (a branch, a parser, a merge/dedup routine)
  should leave at least one test behind. Trivial one-liners don't need one.
- Run `go build ./...`, `go vet ./...`, and `go test ./...` before opening
  the PR — CI runs the same three.

## Commit messages

Conventional-commit-style prefixes (`feat:`, `fix:`, `chore:`, `refactor:`)
are appreciated for readability but not enforced on contributor PRs — we'll
tidy commit messages on merge if needed. You don't need to add `[minor]` /
`[major]` version markers or touch `CHANGELOG.md` yourself; that's handled
when the PR lands (see below).

## What happens after you open a PR

- CI builds, vets, and tests your branch.
- If it looks good, we'll add a `CHANGELOG.md` entry crediting you
  (`... (thanks @you)`) as part of merging — you don't need to write it.
- Fixes land on `canary` first, ship in a beta, get tested, and flow to
  `main` in a later release. This can take a little while; it's not a sign
  your PR was ignored.

## Reporting bugs / requesting features

Use the issue templates. For anything security-related (credential handling,
SSH/GPG keys, keychain storage, OAuth token storage), see
[SECURITY.md](SECURITY.md) instead of filing a public issue.

## License

By contributing, you agree your contributions are licensed under the
project's [Apache License 2.0](LICENSE).
