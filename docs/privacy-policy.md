# Privacy Policy

**gitswitch** — Last updated: May 28, 2026

---

## The short version

gitswitch runs entirely on your machine. It collects nothing, sends nothing, and stores nothing outside of your local filesystem and OS keychain. There are no servers, no analytics, no telemetry, and no third parties involved.

---

## What gitswitch stores and where

All data gitswitch writes stays on your machine:

| What | Where |
|---|---|
| Profiles (name, email, signing key, SSH key paths, host config) | `~/.config/gitswitch/config.yaml` — readable only by your user account |
| PAT tokens and secrets | Your OS keychain (macOS Keychain, Windows Credential Manager, or Linux Secret Service) — never written to disk in plaintext |
| Repository usage history and pins | `~/.config/gitswitch/config.yaml` — local only |
| Shell hook version marker | `~/.config/gitswitch/hook-version` — local only |

gitswitch reads your existing `~/.gitconfig`, `~/.ssh/config`, and `~/.config/gh/hosts.yml` in order to update them. It does not read them for any other purpose.

---

## What gitswitch does not do

- Does not send any data to any server
- Does not connect to the internet (except when you explicitly run `gitswitch upgrade` to check for a new release from GitHub releases, or when `gitswitch import gh` reads your locally stored gh credentials to call the GitHub API on your behalf)
- Does not collect usage analytics or crash reports
- Does not track which commands you run
- Does not store your SSH private keys — it stores the path to them, which stays on your machine
- Does not have a backend, database, or user accounts

---

## GitHub API calls

Two commands make outbound network requests, both initiated explicitly by you:

- `gitswitch upgrade` — checks the GitHub Releases API (`api.github.com`) for the latest version and optionally downloads it. No identifying information is sent beyond what a standard HTTP request includes (your IP address, as with any GitHub download).
- `gitswitch import gh` — reads your locally stored `gh` CLI credentials and calls the GitHub API (`/user`, `/user/emails`) to fetch your profile information. This call goes to GitHub's servers using your own token. gitswitch does not receive or log this data.

---

## GDPR (European Union)

gitswitch does not process personal data on behalf of any organization. The data gitswitch stores (your name, email, and credentials) is stored locally on your own machine, under your own control. Under GDPR, you are the data controller for your own local data. gitswitch, as a locally executed tool with no backend, does not act as a data processor.

If you have questions or concerns, contact: hello@aksisonline.com

---

## CCPA (California)

gitswitch does not sell, share, or disclose personal information to any third party. No personal information is collected by the gitswitch project from users. The data stored by gitswitch lives on your machine and is controlled by you.

---

## Open source

gitswitch is open source. You can inspect exactly what it does with your data at any time:
https://github.com/aksisonline/gitswitch

---

## Changes to this policy

If this policy changes materially, the updated version will be posted at gitswitch.dev/privacy with an updated date. Since gitswitch collects no data, meaningful changes to this policy are unlikely.

---

## Contact

Questions: hello@aksisonline.com
