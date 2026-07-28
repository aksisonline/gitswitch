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

// TestResolveIdentity pins the scope rules the prompt marker, the TUI glyphs and
// the credential helper all read. The local==global case is the one that bites:
// plenty of repos carry a local user.email identical to the global one, and
// treating those as pinned silently kills the nudge hook for them.
func TestResolveIdentity(t *testing.T) {
	isolatedGitConfig(t)
	repo := t.TempDir()
	if err := exec.Command("git", "-C", repo, "init").Run(); err != nil {
		t.Fatal(err)
	}

	assert := func(what string, wantScope Scope, wantEmail string) {
		t.Helper()
		gotScope, gotEmail := ResolveIdentity(repo)
		if gotScope != wantScope || gotEmail != wantEmail {
			t.Errorf("%s: ResolveIdentity = (%v, %q), want (%v, %q)", what, gotScope, gotEmail, wantScope, wantEmail)
		}
		if want := wantScope == ScopeRepo || wantScope == ScopeSession; HasLocalIdentity(repo) != want {
			t.Errorf("%s: HasLocalIdentity = %v, want %v", what, !want, want)
		}
	}
	set := func(t *testing.T, args ...string) {
		t.Helper()
		if err := exec.Command("git", append([]string{"-C", repo, "config"}, args...)...).Run(); err != nil {
			t.Fatal(err)
		}
	}

	assert("nothing configured", ScopeNone, "")

	set(t, "--global", "user.email", "me@global.dev")
	assert("global only", ScopeGlobal, "me@global.dev")

	// A local value that just repeats the global one overrides nothing.
	set(t, "--local", "user.email", "ME@GLOBAL.DEV")
	assert("local same as global", ScopeGlobal, "ME@GLOBAL.DEV")

	set(t, "--local", "user.email", "me@work.com")
	assert("local overrides global", ScopeRepo, "me@work.com")

	// GIT_CONFIG_* env vars are how a session identity arrives; git reports them
	// as scope "command" and they beat the repo's own config.
	t.Setenv("GIT_CONFIG_COUNT", "1")
	t.Setenv("GIT_CONFIG_KEY_0", "user.email")
	t.Setenv("GIT_CONFIG_VALUE_0", "me@session.dev")
	assert("session overrides repo", ScopeSession, "me@session.dev")
}

func gitConfigValue(t *testing.T, key string) string {
	t.Helper()
	out, err := exec.Command("git", "config", "--global", key).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// TestSetSignKey_FormatFollowsKey pins the SSH-vs-GPG detection: gpg.format must
// track the kind of key given, or signing silently fails with the wrong backend.
func TestSetSignKey_FormatFollowsKey(t *testing.T) {
	isolatedGitConfig(t)
	c := New(true)

	cases := []struct {
		key, wantFormat, wantKeyPrefix string
	}{
		{"ABCD1234EF567890", "", "ABCD1234EF567890"},         // GPG key ID → openpgp (unset)
		{"~/.ssh/id_ed25519.pub", "ssh", ExpandPath("~/")},   // path → ssh, ~ expanded
		{"ssh-ed25519 AAAAC3Nz", "ssh", "key::ssh-ed25519"},  // inline key → ssh, key:: added
		{"key::ssh-ed25519 AAAAC3Nz", "ssh", "key::ssh-ed2"}, // already prefixed → untouched
	}
	for _, tc := range cases {
		if err := c.SetSignKey(tc.key); err != nil {
			t.Fatalf("SetSignKey(%q): %v", tc.key, err)
		}
		if got := gitConfigValue(t, "gpg.format"); got != tc.wantFormat {
			t.Errorf("SetSignKey(%q): gpg.format = %q, want %q", tc.key, got, tc.wantFormat)
		}
		if got := gitConfigValue(t, "user.signingkey"); !strings.HasPrefix(got, tc.wantKeyPrefix) {
			t.Errorf("SetSignKey(%q): user.signingkey = %q, want prefix %q", tc.key, got, tc.wantKeyPrefix)
		}
	}

	// Empty key clears both, so a GPG-less profile leaves no stale ssh format.
	if err := c.SetSignKey(""); err != nil {
		t.Fatalf("SetSignKey(\"\"): %v", err)
	}
	if got := gitConfigValue(t, "gpg.format"); got != "" {
		t.Errorf("gpg.format should be unset, got %q", got)
	}
	if got := gitConfigValue(t, "user.signingkey"); got != "" {
		t.Errorf("user.signingkey should be unset, got %q", got)
	}
}

// TestLocalIdentityRoundtrip covers the repo-pin path: a pin writes identity to
// the repo's own config (leaving global alone), unpin removes every key it wrote.
func TestLocalIdentityRoundtrip(t *testing.T) {
	isolatedGitConfig(t)
	repo := t.TempDir()
	if err := exec.Command("git", "-C", repo, "init").Run(); err != nil {
		t.Fatal(err)
	}
	t.Chdir(repo)

	if HasLocalIdentity(repo) {
		t.Fatal("fresh repo should have no local identity")
	}

	local := New(false)
	if err := local.SetUser("Ada", "ada@work.dev"); err != nil {
		t.Fatal(err)
	}
	if err := local.SetSignKey("~/.ssh/id_ed25519.pub"); err != nil {
		t.Fatal(err)
	}
	if !HasLocalIdentity(repo) {
		t.Error("pinned repo should report a local identity")
	}
	if got := gitConfigValue(t, "user.email"); got != "" {
		t.Errorf("global user.email should be untouched by a pin, got %q", got)
	}

	if err := local.ClearIdentity(); err != nil {
		t.Fatalf("ClearIdentity: %v", err)
	}
	if HasLocalIdentity(repo) {
		t.Error("unpinned repo should have no local identity")
	}
	for _, k := range []string{"user.name", "gpg.format", "user.signingkey"} {
		out, err := exec.Command("git", "-C", repo, "config", "--local", k).Output()
		if err == nil && strings.TrimSpace(string(out)) != "" {
			t.Errorf("local %s should be cleared, got %q", k, out)
		}
	}
	// Clearing an already-clean scope is a no-op, not an error (git exit 5).
	if err := local.ClearIdentity(); err != nil {
		t.Errorf("second ClearIdentity should be a no-op, got: %v", err)
	}
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
