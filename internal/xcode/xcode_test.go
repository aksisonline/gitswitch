package xcode

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func acct(user, host string) Account {
	return accountFromRaw(map[string]any{"Username": user, "HostURL": host, "AccountName": "GitHub"})
}

func names(a []Account) []string {
	var out []string
	for _, x := range a {
		out = append(out, x.Username)
	}
	return out
}

func TestPromote(t *testing.T) {
	gh := GitHubHost
	accts := []Account{acct("a", "https://gitlab.com"), acct("b", gh), acct("c", gh), acct("d", gh)}

	out, changed := Promote(accts, "d")
	if !changed {
		t.Fatal("expected change")
	}
	if got := fmt.Sprint(names(out)); got != "[a d b c]" {
		t.Errorf("order = %s, want [a d b c] (gitlab keeps its slot)", got)
	}

	if _, changed := Promote(accts, "b"); changed {
		t.Error("already-first account must not change")
	}
	if _, changed := Promote(accts, "nobody"); changed {
		t.Error("unknown user must not change")
	}
	if _, changed := Promote(accts, "a"); changed {
		t.Error("account on another host must not be promoted among github accounts")
	}
	if out, _ := Promote(accts, "C"); fmt.Sprint(names(out)) != "[a c b d]" {
		t.Errorf("username match should be case-insensitive, got %v", names(out))
	}
}

func TestBindUnbindURL(t *testing.T) {
	cases := []struct {
		in, user, want string
		ok             bool
	}{
		{"https://github.com/acme/app.git", "bob", "https://bob@github.com/acme/app.git", true},
		{"https://bob@github.com/acme/app.git", "bob", "https://bob@github.com/acme/app.git", false},
		{"https://alice@github.com/acme/app.git", "bob", "https://bob@github.com/acme/app.git", true},
		{"git@github.com:acme/app.git", "bob", "git@github.com:acme/app.git", false},
		{"https://gitlab.com/acme/app.git", "bob", "https://gitlab.com/acme/app.git", false},
		{"https://x:tok@github.com/acme/app.git", "bob", "https://x:tok@github.com/acme/app.git", false},
	}
	for _, c := range cases {
		got, ok := BindURL(c.in, c.user)
		if got != c.want || ok != c.ok {
			t.Errorf("BindURL(%q,%q) = %q,%v want %q,%v", c.in, c.user, got, ok, c.want, c.ok)
		}
	}

	if got, ok := UnbindURL("https://bob@github.com/acme/app.git"); !ok || got != "https://github.com/acme/app.git" {
		t.Errorf("UnbindURL = %q,%v", got, ok)
	}
	if _, ok := UnbindURL("https://github.com/acme/app.git"); ok {
		t.Error("unbound URL must be unchanged")
	}
	if _, ok := UnbindURL("https://x:tok@github.com/acme/app.git"); ok {
		t.Error("URLs carrying a password are never rewritten")
	}
}

// TestStoreRoundTrip exercises the real defaults/plutil plumbing against a
// throwaway domain — never com.apple.dt.Xcode.
func TestStoreRoundTrip(t *testing.T) {
	if !Supported() {
		t.Skip("defaults/plutil unavailable")
	}
	domain := fmt.Sprintf("dev.gitswitch.test.xcode.%d", os.Getpid())
	t.Cleanup(func() { _ = exec.Command("defaults", "delete", domain).Run() })
	s := &Store{Domain: domain}

	if got, err := s.Accounts(); err != nil || len(got) != 0 {
		t.Fatalf("empty domain: %v, %v", got, err)
	}

	in := []Account{
		accountFromRaw(map[string]any{
			"Username": "b", "HostURL": GitHubHost, "AccountName": "GitHub", "CloneType": float64(0),
			"Authentication": map[string]any{"TokenKeychainEntryName": "Xcode-SCM-Token-1", "Type": "Bearer"},
		}),
		accountFromRaw(map[string]any{"Username": "a", "HostURL": GitHubHost, "AccountName": "GitHub"}),
	}
	if err := s.Save(in); err != nil {
		t.Fatal(err)
	}
	got, err := s.Accounts()
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprint(names(got)) != "[b a]" {
		t.Fatalf("round trip order = %v", names(got))
	}
	auth, _ := got[0].Raw["Authentication"].(map[string]any)
	if auth["TokenKeychainEntryName"] != "Xcode-SCM-Token-1" {
		t.Errorf("token reference lost on rewrite: %v", got[0].Raw)
	}

	reordered, _ := Promote(got, "a")
	if err := s.Save(reordered); err != nil {
		t.Fatal(err)
	}
	got, _ = s.Accounts()
	if fmt.Sprint(names(got)) != "[a b]" {
		t.Errorf("after promote = %v", names(got))
	}
}
