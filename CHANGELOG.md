# Changelog

## v0.2.0

### What's New

**gitswitch Beta is here — opt in to get the next generation of features early.**

The beta channel is live and already ships several big upgrades:

- **New UI** — rebuilt from the ground up, faster navigation and a cleaner look
- **Custom Alias** — give any profile a short name and switch to it in one word
- **OAuth Support** — connect your GitHub account without ever touching a token manually (uses GitHub's device flow)
- **Easier Setup Flow** — first-run wizard walks you through everything in under a minute, no config file editing required

To switch to beta:

```sh
gitswitch upgrade --channel beta
```

Or install fresh:

```sh
curl -fsSL https://gitswitch.dev/install.sh | bash -s -- --channel beta
```

Feedback and bug reports → https://github.com/aksisonline/gitswitch/issues

---

- Migrated CI/CD to GoReleaser — releases are now faster, Homebrew tap updates automatically, and checksums are always published
- Added What's New splash screen — shows release highlights the first time you launch after a significant update (minor or major version)

### Bug Fixes

- Fixed: forced upgrade when running a version that was pulled from GitHub
- Fixed: update banner no longer shows for a release that no longer exists (validates cached tag before displaying)
