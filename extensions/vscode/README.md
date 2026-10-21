# gitswitch.dev for VS Code

Status bar widget for [gitswitch.dev](https://gitswitch.dev) ([GitHub](https://github.com/aksisonline/gitswitch)) — see and switch your active git identity without leaving the editor.

## Requires

The `gitswitch` CLI on your `PATH` (or set `gitswitch.binaryPath` in settings).

## Features

- Status bar shows the active profile nickname (`●` when pinned, `[session]` when overridden, and pre-commit guard status in the tooltip).
- Click the status bar item, or run **gitswitch.dev: Switch Profile**, to pick a different identity via Quick Pick.
- **gitswitch.dev: Pin Profile to This Repo** / **gitswitch.dev: Unpin This Repo** — pin an identity to the open workspace's local git config and `gitswitch` history.
- **gitswitch.dev: Toggle Pre-commit Guard** — enable or disable the global pre-commit identity guard hook (`gitswitch guard`).

## Development & Local Testing

```bash
bun install
bun run compile
```

Open `extensions/vscode` in VS Code and press **F5** to launch an Extension Development Host window with the extension loaded.
