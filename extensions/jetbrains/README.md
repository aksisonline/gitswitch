# gitswitch for JetBrains IDEs

Status bar widget for [gitswitch](https://github.com/aksisonline/git-switcher) — see and switch your active git identity without leaving the IDE.

## Requires

The `gitswitch` CLI on your `PATH`. (Override its location with the `-Dgitswitch.binaryPath=...` JVM system property if needed.)

## Features

- Status bar shows the active profile nickname (`●` when pinned, `[session]` when overridden, and pre-commit guard status in the tooltip) — click it to switch via a popup list.
- **Tools → gitswitch → Switch Profile / Pin Profile to This Repo / Unpin This Repo / Toggle Pre-commit Guard**.

## Development & Local Testing

```bash
./gradlew runIde       # launch a sandbox IDE with the plugin loaded
./gradlew buildPlugin  # produces build/distributions/gitswitch-0.1.0.zip
```

To install into your regular JetBrains IDE: **Settings → Plugins → ⚙️ → Install Plugin from Disk…** and select `build/distributions/gitswitch-0.1.0.zip`.
