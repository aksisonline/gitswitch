package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/aksisonline/gitswitch/internal/xcode"
)

func gitIn(t *testing.T, dir string, args ...string) string {
	t.Helper()
	out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func TestBindUnbindRemotes(t *testing.T) {
	dir := t.TempDir()
	gitIn(t, dir, "init", "-q")
	gitIn(t, dir, "remote", "add", "origin", "https://github.com/acme/app.git")
	gitIn(t, dir, "remote", "add", "ssh", "git@github.com:acme/app.git")
	gitIn(t, dir, "remote", "add", "lab", "https://gitlab.com/acme/app.git")

	changed, err := bindRemotes(dir, "bob")
	if err != nil {
		t.Fatal(err)
	}
	if len(changed) != 1 || changed[0] != "origin" {
		t.Fatalf("bound %v, want only origin (ssh and non-GitHub remotes untouched)", changed)
	}
	if got := gitIn(t, dir, "remote", "get-url", "origin"); got != "https://bob@github.com/acme/app.git" {
		t.Errorf("origin = %s", got)
	}
	if got := gitIn(t, dir, "remote", "get-url", "ssh"); got != "git@github.com:acme/app.git" {
		t.Errorf("ssh remote changed: %s", got)
	}

	if changed, _ := bindRemotes(dir, "bob"); len(changed) != 0 {
		t.Errorf("re-binding the same user must be a no-op, changed %v", changed)
	}
	// Switching the pinned account rebinds rather than stacking usernames.
	if _, err := bindRemotes(dir, "alice"); err != nil {
		t.Fatal(err)
	}
	if got := gitIn(t, dir, "remote", "get-url", "origin"); got != "https://alice@github.com/acme/app.git" {
		t.Errorf("rebind origin = %s", got)
	}

	if _, err := unbindRemotes(dir); err != nil {
		t.Fatal(err)
	}
	if got := gitIn(t, dir, "remote", "get-url", "origin"); got != "https://github.com/acme/app.git" {
		t.Errorf("unbind origin = %s", got)
	}
}

func TestBindRemotesOutsideRepo(t *testing.T) {
	if _, err := bindRemotes(t.TempDir(), "bob"); err == nil {
		t.Error("expected an error outside a git repo")
	}
}

func TestSyncXcodeAccount(t *testing.T) {
	if !xcode.Supported() {
		t.Skip("defaults/plutil unavailable")
	}
	domain := fmt.Sprintf("dev.gitswitch.test.cmdxcode.%d", os.Getpid())
	t.Cleanup(func() { _ = exec.Command("defaults", "delete", domain).Run() })
	s := &xcode.Store{Domain: domain}

	mk := func(u string) xcode.Account {
		return xcode.Account{Username: u, HostURL: xcode.GitHubHost, Raw: map[string]any{"Username": u, "HostURL": xcode.GitHubHost}}
	}
	if err := s.Save([]xcode.Account{mk("a"), mk("b")}); err != nil {
		t.Fatal(err)
	}

	res, err := syncXcodeAccount(s, "b", true)
	if err != nil || !res.Changed || !res.Found {
		t.Fatalf("dry run: %+v, %v", res, err)
	}
	if got, _ := s.Accounts(); got[0].Username != "a" {
		t.Fatal("dry run must not write")
	}

	if res, err = syncXcodeAccount(s, "b", false); err != nil || !res.Changed {
		t.Fatalf("sync: %+v, %v", res, err)
	}
	if got, _ := s.Accounts(); got[0].Username != "b" {
		t.Errorf("order after sync = %v", got)
	}
	if res, _ = syncXcodeAccount(s, "b", false); res.Changed {
		t.Error("second sync must be a no-op")
	}
	if res, _ = syncXcodeAccount(s, "ghost", false); res.Found || res.Changed {
		t.Errorf("unknown account: %+v", res)
	}
}
