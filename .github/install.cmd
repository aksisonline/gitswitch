@echo off
REM gitswitch Windows installer (Command Prompt / batch)
REM Usage:
REM   curl -fsSL https://raw.githubusercontent.com/aksisonline/gitswitch/main/.github/install.cmd -o install.cmd && install.cmd
REM   install.cmd v0.2.0   (pin a specific version)

setlocal enabledelayedexpansion

REM ── version ──────────────────────────────────────────────────────────────────
set VERSION=%1
if "%VERSION%"=="" (
    echo Fetching latest version...
    for /f "delims=" %%i in ('powershell -NoProfile -Command "(Invoke-RestMethod https://api.github.com/repos/aksisonline/gitswitch/releases/latest).tag_name"') do set VERSION=%%i
    if "!VERSION!"=="" (
        echo ERROR: Could not fetch latest version. Pass version as argument: install.cmd v0.2.0
        exit /b 1
    )
)
echo Installing gitswitch !VERSION!...

REM ── architecture ─────────────────────────────────────────────────────────────
if "%PROCESSOR_ARCHITECTURE%"=="ARM64" (
    set BINARY=gitswitch-windows-arm64.exe
) else if "%PROCESSOR_ARCHITECTURE%"=="AMD64" (
    set BINARY=gitswitch-windows-amd64.exe
) else (
    echo ERROR: Unsupported architecture: %PROCESSOR_ARCHITECTURE%
    exit /b 1
)

REM ── install dir ──────────────────────────────────────────────────────────────
if "%GS_INSTALL_DIR%"=="" (
    set INSTALL_DIR=%LOCALAPPDATA%\gitswitch
) else (
    set INSTALL_DIR=%GS_INSTALL_DIR%
)

if not exist "!INSTALL_DIR!" mkdir "!INSTALL_DIR!"

REM ── download ─────────────────────────────────────────────────────────────────
set RELEASE_URL=https://github.com/aksisonline/gitswitch/releases/download/!VERSION!/!BINARY!
set DEST=!INSTALL_DIR!\gitswitch.exe
set GS_DEST=!INSTALL_DIR!\gs.exe

echo Downloading !BINARY!...
curl -fsSL "!RELEASE_URL!" -o "!DEST!"
if errorlevel 1 (
    echo ERROR: Download failed. Check version or internet connection.
    exit /b 1
)

REM ── gs alias ─────────────────────────────────────────────────────────────────
copy /y "!DEST!" "!GS_DEST!" >nul

REM ── PATH ─────────────────────────────────────────────────────────────────────
echo Updating PATH...
powershell -NoProfile -Command ^
  "$p = [System.Environment]::GetEnvironmentVariable('PATH','User'); if ($p -notlike '*!INSTALL_DIR!*') { [System.Environment]::SetEnvironmentVariable('PATH', $p + ';!INSTALL_DIR!', 'User'); Write-Host 'PATH updated.' }"

REM ── done ─────────────────────────────────────────────────────────────────────
echo.
echo   OK  gitswitch.exe -^> !DEST!
echo   OK  gs.exe        -^> !GS_DEST!
echo.
echo   Restart your terminal, then run:  gs login
echo.
endlocal
