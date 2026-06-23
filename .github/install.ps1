# gitswitch Windows installer
# Usage (PowerShell):
#   irm https://raw.githubusercontent.com/aksisonline/gitswitch/main/.github/install.ps1 | iex
#   # or pin a version:
#   $env:GS_VERSION="v0.2.0"; irm .../install.ps1 | iex
#
# Usage (winget — coming soon):
#   winget install aksisonline.gitswitch

[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.SecurityProtocolType]::Tls12

$ErrorActionPreference = "Stop"

# ── version ────────────────────────────────────────────────────────────────────
$VERSION = $env:GS_VERSION
if (-not $VERSION) {
    try {
        $rel = Invoke-RestMethod -Uri "https://api.github.com/repos/aksisonline/gitswitch/releases/latest" `
            -Headers @{ "User-Agent" = "gitswitch-installer" }
        $VERSION = $rel.tag_name
    } catch {
        Write-Error "Could not fetch latest version from GitHub. Set `$env:GS_VERSION and retry."
        exit 1
    }
}
Write-Host "Installing gitswitch $VERSION for Windows..."

# ── architecture ───────────────────────────────────────────────────────────────
$arch = $env:PROCESSOR_ARCHITECTURE
if ($arch -eq "ARM64") {
    $BINARY = "gitswitch-windows-arm64.exe"
} elseif ($arch -eq "AMD64" -or $arch -eq "x86_64") {
    $BINARY = "gitswitch-windows-amd64.exe"
} else {
    Write-Error "Unsupported architecture: $arch"
    exit 1
}

# ── download ───────────────────────────────────────────────────────────────────
$RELEASE_URL = "https://github.com/aksisonline/gitswitch/releases/download/$VERSION/$BINARY"
$TEMP_FILE   = Join-Path $env:TEMP "gitswitch-install.exe"

Write-Host "Downloading $BINARY..."
try {
    Invoke-WebRequest -Uri $RELEASE_URL -OutFile $TEMP_FILE -UseBasicParsing
} catch {
    Write-Error "Download failed: $_"
    exit 1
}

# ── install dir ────────────────────────────────────────────────────────────────
# Prefer %LOCALAPPDATA%\gitswitch (no admin required).
# Set GS_INSTALL_DIR env var to override.
$INSTALL_DIR = $env:GS_INSTALL_DIR
if (-not $INSTALL_DIR) {
    $INSTALL_DIR = Join-Path $env:LOCALAPPDATA "gitswitch"
}
if (-not (Test-Path $INSTALL_DIR)) {
    New-Item -ItemType Directory -Path $INSTALL_DIR | Out-Null
}

$DEST = Join-Path $INSTALL_DIR "gitswitch.exe"
Copy-Item -Force $TEMP_FILE $DEST
Remove-Item $TEMP_FILE -ErrorAction SilentlyContinue

# ── gs alias ───────────────────────────────────────────────────────────────────
$GS_DEST = Join-Path $INSTALL_DIR "gs.exe"
Copy-Item -Force $DEST $GS_DEST

# ── PATH ───────────────────────────────────────────────────────────────────────
$USER_PATH = [System.Environment]::GetEnvironmentVariable("PATH", "User")
if ($USER_PATH -notlike "*$INSTALL_DIR*") {
    [System.Environment]::SetEnvironmentVariable(
        "PATH",
        "$USER_PATH;$INSTALL_DIR",
        "User"
    )
    Write-Host "Added $INSTALL_DIR to your PATH (restart your terminal to apply)."
}

# ── sanity check ───────────────────────────────────────────────────────────────
try {
    $ver = & $DEST version 2>&1
    Write-Host ""
    Write-Host "  v  gitswitch installed: $ver"
} catch {
    Write-Warning "Binary installed but failed version check. Check $DEST manually."
}

Write-Host ""
Write-Host "  v  gitswitch.exe -> $DEST"
Write-Host "  v  gs.exe        -> $GS_DEST"
Write-Host ""
Write-Host "  Restart your terminal, then run:  gs login"
Write-Host ""
