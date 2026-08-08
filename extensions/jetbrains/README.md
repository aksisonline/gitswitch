# gitswitch for JetBrains IDEs

Status bar widget for [gitswitch](https://github.com/aksisonline/git-switcher) — see and switch your active git identity without leaving the IDE.

## Requires

The `gitswitch` CLI on your `PATH`. (Override its location with the `-Dgitswitch.binaryPath=...` JVM system property if needed.)

## Features

- Status bar shows the active profile nickname (and whether it's pinned/session/global) — click it to switch via a popup list.
- **Tools → gitswitch → Switch Profile / Pin Profile to This Repo / Unpin This Repo**.

## Development

```
./gradlew runIde       # launch a sandbox IDE with the plugin loaded
./gradlew buildPlugin  # produces build/distributions/gitswitch-*.zip
```
