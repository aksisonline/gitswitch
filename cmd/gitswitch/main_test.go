package main

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/aksisonline/gitswitch/internal/git"
)

func TestActiveAccountForHost(t *testing.T) {
	accounts := []git.GHAccount{
		{Login: "alice", Host: "github.com", Active: false},
		{Login: "alice", Host: "ghe.acme.com", Active: true},
		{Login: "bob", Host: "github.com", Active: true},
	}
	got, ok := activeAccountForHost(accounts, "github.com")
	if !ok || got.Login != "bob" {
		t.Fatalf("got %+v, ok=%v; want bob/true", got, ok)
	}
	if _, ok := activeAccountForHost(accounts, "gitlab.example.com"); ok {
		t.Fatal("expected no match for unrelated host")
	}
}

// TestPickLoggedInAccount covers the multi-account-on-one-host case that
// the old activeAccountForHost logic got wrong: when a host already has an
// active account and you log in a *different* one, gh makes the new one
// active and the old one inactive — pickLoggedInAccount must return the
// newly-activated account, not the previously-active one.
func TestPickLoggedInAccount(t *testing.T) {
	// Before: abhiramkanna active, aksisonline present but inactive.
	before := []git.GHAccount{
		{Login: "abhiramkanna-edirq", Host: "github.com", Active: true},
		{Login: "aksisonline", Host: "github.com", Active: false},
	}
	// After logging in as aksisonline: it became active, abhiramkanna didn't.
	after := []git.GHAccount{
		{Login: "abhiramkanna-edirq", Host: "github.com", Active: false},
		{Login: "aksisonline", Host: "github.com", Active: true},
	}

	got, ok := pickLoggedInAccount(before, after, "github.com")
	if !ok || got.Login != "aksisonline" {
		t.Fatalf("got %+v, ok=%v; want aksisonline/true", got, ok)
	}

	// Re-auth of the already-active account: active login unchanged.
	same := []git.GHAccount{
		{Login: "abhiramkanna-edirq", Host: "github.com", Active: true},
		{Login: "aksisonline", Host: "github.com", Active: false},
	}
	if got, ok := pickLoggedInAccount(same, same, "github.com"); !ok || got.Login != "abhiramkanna-edirq" {
		t.Fatalf("re-auth: got %+v, ok=%v; want abhiramkanna-edirq/true", got, ok)
	}

	// Fresh host: only the new account exists/active afterwards.
	freshBefore := []git.GHAccount{{Login: "bob", Host: "github.com", Active: true}}
	freshAfter := []git.GHAccount{
		{Login: "bob", Host: "github.com", Active: true},
		{Login: "carol", Host: "ghe.acme.com", Active: true},
	}
	if got, ok := pickLoggedInAccount(freshBefore, freshAfter, "ghe.acme.com"); !ok || got.Login != "carol" {
		t.Fatalf("fresh host: got %+v, ok=%v; want carol/true", got, ok)
	}

	if _, ok := pickLoggedInAccount(before, after, "gitlab.example.com"); ok {
		t.Fatal("expected no match for unrelated host")
	}
}

// TestPickLoggedInAccount_Live switches the real active github.com account
// on this machine and verifies detection without misattributing, then
// restores the original active account. Skips if fewer than two
// github.com accounts are logged in.
func TestPickLoggedInAccount_Live(t *testing.T) {
	if !git.IsGHInstalled() {
		t.Skip("gh not installed")
	}
	all := git.ListGHUsers()
	var ghc []git.GHAccount
	for _, a := range all {
		if a.Host == "github.com" {
			ghc = append(ghc, a)
		}
	}
	if len(ghc) < 2 {
		t.Skip("need >=2 github.com accounts to exercise multi-account detection")
	}

	var active, other string
	for _, a := range ghc {
		if a.Active {
			active = a.Login
		} else {
			other = a.Login
		}
	}
	if active == "" || other == "" {
		t.Skip("no inactive github.com account to switch to")
	}

	// Make `other` active — simulating what `gh auth login` does for it.
	runGH := func(args ...string) error {
		cmd := exec.Command("gh", args...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Logf("gh %s: %v\n%s", strings.Join(args, " "), err, out)
		}
		return err
	}
	if err := runGH("auth", "switch", "--user", other); err != nil {
		t.Fatalf("could not switch active account to %s: %v", other, err)
	}
	// Always restore the original active account.
	defer runGH("auth", "switch", "--user", active)

	before := []git.GHAccount{{Login: active, Host: "github.com", Active: true}}
	after := git.ListGHUsers()

	got, ok := pickLoggedInAccount(before, after, "github.com")
	if !ok || got.Login != other {
		t.Fatalf("live: got %+v, ok=%v; want %s/true", got, ok, other)
	}
}

func TestCommandAllowsMissingGit(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want bool
	}{
		{"bare invocation", []string{"gitswitch"}, true},
		{"doctor", []string{"gitswitch", "doctor"}, true},
		{"snake", []string{"gitswitch", "snake"}, true},
		{"add", []string{"gitswitch", "add"}, false},
	}
	for _, c := range cases {
		if got := commandAllowsMissingGit(c.args); got != c.want {
			t.Errorf("%s: commandAllowsMissingGit(%v) = %v, want %v", c.name, c.args, got, c.want)
		}
	}
}
