package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestGuardHookInstallUninstallRoundtrip(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	isolatedGitConfig(t)

	if IsGuardHookInstalled() {
		t.Fatal("guard hook should not be installed in a fresh environment")
	}

	if err := InstallGuardHook(); err != nil {
		t.Fatalf("InstallGuardHook: %v", err)
	}
	if !IsGuardHookInstalled() {
		t.Fatal("IsGuardHookInstalled = false after install")
	}

	hooksDir, err := GuardHooksDir()
	if err != nil {
		t.Fatal(err)
	}
	if got := gitConfigValue(t, "core.hooksPath"); filepath.Clean(got) != filepath.Clean(hooksDir) {
		t.Errorf("core.hooksPath = %q, want %q", got, hooksDir)
	}

	// Idempotent second install
	if err := InstallGuardHook(); err != nil {
		t.Fatalf("second InstallGuardHook: %v", err)
	}
	if !IsGuardHookInstalled() {
		t.Fatal("IsGuardHookInstalled = false after second install")
	}

	if err := UninstallGuardHook(); err != nil {
		t.Fatalf("UninstallGuardHook: %v", err)
	}
	if IsGuardHookInstalled() {
		t.Error("IsGuardHookInstalled = true after uninstall")
	}
	if got := gitConfigValue(t, "core.hooksPath"); got != "" {
		t.Errorf("core.hooksPath should be unset after uninstall, got %q", got)
	}

	// Idempotent second uninstall
	if err := UninstallGuardHook(); err != nil {
		t.Errorf("second UninstallGuardHook: %v", err)
	}
}

func TestGuardHookPreservesAndRestoresPriorHooksPath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	isolatedGitConfig(t)

	priorHooks := filepath.Join(t.TempDir(), "custom-hooks")
	if err := os.MkdirAll(priorHooks, 0755); err != nil {
		t.Fatal(err)
	}
	if err := exec.Command("git", "config", "--global", "core.hooksPath", priorHooks).Run(); err != nil {
		t.Fatal(err)
	}

	if err := InstallGuardHook(); err != nil {
		t.Fatalf("InstallGuardHook: %v", err)
	}
	if !IsGuardHookInstalled() {
		t.Fatal("expected guard hook installed")
	}

	if err := UninstallGuardHook(); err != nil {
		t.Fatalf("UninstallGuardHook: %v", err)
	}
	if got := gitConfigValue(t, "core.hooksPath"); filepath.Clean(got) != filepath.Clean(priorHooks) {
		t.Errorf("core.hooksPath after uninstall = %q, want restored %q", got, priorHooks)
	}
}

func TestGuardPreCommitChainsToLocalHook(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	isolatedGitConfig(t)

	if err := InstallGuardHook(); err != nil {
		t.Fatal(err)
	}

	repo := t.TempDir()
	if err := exec.Command("git", "-C", repo, "init").Run(); err != nil {
		t.Fatal(err)
	}

	markerFile := filepath.Join(repo, "local-hook-ran")
	localHook := filepath.Join(repo, ".git", "hooks", "pre-commit")
	script := "#!/bin/sh\necho ok > " + shellQuote(markerFile) + "\n"
	if err := os.WriteFile(localHook, []byte(script), 0755); err != nil {
		t.Fatal(err)
	}
	_ = os.Chmod(localHook, 0755)

	hooksDir, _ := GuardHooksDir()
	// Run the pre-commit script with PATH excluding gitswitch so we isolate the
	// local-hook chaining behavior here (guard check itself is tested in cmd/gitswitch).
	cmd := exec.Command("/bin/sh", filepath.Join(hooksDir, "pre-commit"))
	cmd.Dir = repo
	cmd.Env = []string{"PATH=/usr/bin:/bin:/opt/homebrew/bin", "HOME=" + home}
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("pre-commit script failed: %v\n%s", err, out)
	}

	data, err := os.ReadFile(markerFile)
	if err != nil || strings.TrimSpace(string(data)) != "ok" {
		t.Errorf("expected local .git/hooks/pre-commit to be chained, marker err=%v data=%q", err, data)
	}
}

func TestResolveCommitEmail(t *testing.T) {
	isolatedGitConfig(t)
	repo := t.TempDir()
	if err := exec.Command("git", "-C", repo, "init").Run(); err != nil {
		t.Fatal(err)
	}
	if err := exec.Command("git", "config", "--global", "user.email", "alice@global.dev").Run(); err != nil {
		t.Fatal(err)
	}

	// When git exports GIT_AUTHOR_EMAIL equal to the configured email, keep the underlying scope.
	t.Setenv("GIT_AUTHOR_EMAIL", "alice@global.dev")
	if scope, email := ResolveCommitEmail(repo); scope != ScopeGlobal || email != "alice@global.dev" {
		t.Errorf("same GIT_AUTHOR_EMAIL: got (%v, %q), want (ScopeGlobal, alice@global.dev)", scope, email)
	}

	// When GIT_AUTHOR_EMAIL overrides the configured email, report ScopeSession and the override.
	t.Setenv("GIT_AUTHOR_EMAIL", "override@work.com")
	if scope, email := ResolveCommitEmail(repo); scope != ScopeSession || email != "override@work.com" {
		t.Errorf("override GIT_AUTHOR_EMAIL: got (%v, %q), want (ScopeSession, override@work.com)", scope, email)
	}
}
