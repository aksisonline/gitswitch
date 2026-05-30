package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// isolatedGitConfig points GIT_CONFIG_GLOBAL at a fresh temp file so that
// --global writes don't touch the developer's real ~/.gitconfig.
func isolatedGitConfig(t *testing.T) {
	t.Helper()
	cfg := filepath.Join(t.TempDir(), "gitconfig")
	if err := os.WriteFile(cfg, nil, 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GIT_CONFIG_GLOBAL", cfg)
}

func globalHelpers(t *testing.T) []string {
	t.Helper()
	out, err := exec.Command("git", "config", "--global", "--get-all", "credential.helper").Output()
	if err != nil {
		return nil // none set
	}
	var helpers []string
	for _, line := range strings.Split(string(out), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			helpers = append(helpers, line)
		}
	}
	return helpers
}

func TestInstallCredentialHelper_Idempotent(t *testing.T) {
	isolatedGitConfig(t)

	if IsCredentialHelperInstalled() {
		t.Fatal("should not be installed in a fresh config")
	}
	if err := InstallCredentialHelper(); err != nil {
		t.Fatalf("first install: %v", err)
	}
	if !IsCredentialHelperInstalled() {
		t.Fatal("should be installed after InstallCredentialHelper")
	}
	if err := InstallCredentialHelper(); err != nil {
		t.Fatalf("second install: %v", err)
	}

	count := 0
	for _, h := range globalHelpers(t) {
		if h == credentialHelperValue {
			count++
		}
	}
	if count != 1 {
		t.Errorf("gitswitch helper present %d times, want exactly 1", count)
	}
}

func TestInstallCredentialHelper_PreservesExistingAndOrdersFirst(t *testing.T) {
	isolatedGitConfig(t)

	// Pre-seed an existing helper (as if osxkeychain were configured).
	if err := exec.Command("git", "config", "--global", "--add", "credential.helper", "osxkeychain").Run(); err != nil {
		t.Fatal(err)
	}

	if err := InstallCredentialHelper(); err != nil {
		t.Fatalf("install: %v", err)
	}

	helpers := globalHelpers(t)
	if len(helpers) != 2 {
		t.Fatalf("want 2 helpers, got %d: %v", len(helpers), helpers)
	}
	if helpers[0] != credentialHelperValue {
		t.Errorf("gitswitch should be first, got order: %v", helpers)
	}
	found := false
	for _, h := range helpers {
		if h == "osxkeychain" {
			found = true
		}
	}
	if !found {
		t.Errorf("osxkeychain helper should be preserved, got: %v", helpers)
	}
}

func TestUninstallCredentialHelper_RemovesOnlyOurs(t *testing.T) {
	isolatedGitConfig(t)

	_ = exec.Command("git", "config", "--global", "--add", "credential.helper", "osxkeychain").Run()
	if err := InstallCredentialHelper(); err != nil {
		t.Fatal(err)
	}

	if err := UninstallCredentialHelper(); err != nil {
		t.Fatalf("uninstall: %v", err)
	}
	if IsCredentialHelperInstalled() {
		t.Error("gitswitch helper should be gone after uninstall")
	}

	helpers := globalHelpers(t)
	if len(helpers) != 1 || helpers[0] != "osxkeychain" {
		t.Errorf("osxkeychain should remain, got: %v", helpers)
	}

	// Second uninstall is a clean no-op (tolerates exit code 5).
	if err := UninstallCredentialHelper(); err != nil {
		t.Errorf("second uninstall should be a no-op, got: %v", err)
	}
}
