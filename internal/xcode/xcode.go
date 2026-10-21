// Package xcode integrates gitswitch with Xcode's Source Control accounts.
//
// Xcode keeps its GitHub accounts in the com.apple.dt.Xcode defaults domain
// under IDESourceControlHostAccounts_10: an array of dictionaries holding the
// username, host and a *reference* (TokenKeychainEntryName) to the token — the
// token itself lives in Xcode's private Keychain access group, which only
// Apple-signed Xcode processes can write. gitswitch therefore never touches
// tokens; it only manages which account Xcode prefers:
//
//   - Xcode matches a remote URL to an account first by a username embedded in
//     the URL (https://user@github.com/...), then by the owner path segment,
//     and otherwise falls back to the first account for the host. Promoting the
//     active profile's account to the front, and binding pinned repos with an
//     embedded username, makes Xcode follow gitswitch automatically.
package xcode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
)

const (
	// Domain is Xcode's defaults domain.
	Domain = "com.apple.dt.Xcode"
	// AccountsKey is the defaults key holding the source control accounts.
	AccountsKey = "IDESourceControlHostAccounts_10"
	// GitHubHost is the host URL Xcode records for github.com accounts.
	GitHubHost = "https://github.com"
)

// Account is one Xcode source control account. Raw keeps every key Xcode wrote
// so a rewrite preserves fields gitswitch doesn't understand.
type Account struct {
	Username string
	HostURL  string
	Name     string
	Raw      map[string]any
}

// Supported reports whether Xcode accounts can be read/written on this system.
func Supported() bool {
	if runtime.GOOS != "darwin" {
		return false
	}
	for _, bin := range []string{"defaults", "plutil"} {
		if _, err := exec.LookPath(bin); err != nil {
			return false
		}
	}
	return true
}

// IsRunning reports whether Xcode is currently open. Xcode only reads its
// accounts at launch, so a change made while it runs applies after a restart.
func IsRunning() bool {
	return exec.Command("pgrep", "-x", "Xcode").Run() == nil
}

// Store reads and writes Xcode's accounts through the defaults system (never by
// editing the plist file directly, which cfprefsd would overwrite). Domain is
// overridable so tests never touch the real Xcode preferences.
type Store struct {
	Domain string
}

// NewStore returns a Store for Xcode's real defaults domain.
func NewStore() *Store { return &Store{Domain: Domain} }

// Accounts returns the configured source control accounts in Xcode's order.
// A domain with no accounts yet yields an empty slice, not an error.
func (s *Store) Accounts() ([]Account, error) {
	raw, err := s.rawAccounts()
	if err != nil {
		return nil, err
	}
	out := make([]Account, 0, len(raw))
	for _, r := range raw {
		out = append(out, accountFromRaw(r))
	}
	return out, nil
}

func (s *Store) rawAccounts() ([]map[string]any, error) {
	exp := exec.Command("defaults", "export", s.Domain, "-")
	var plist bytes.Buffer
	exp.Stdout = &plist
	if err := exp.Run(); err != nil {
		return nil, nil // domain doesn't exist yet: no accounts
	}
	ext := exec.Command("plutil", "-extract", AccountsKey, "json", "-o", "-", "-")
	ext.Stdin = &plist
	var out, stderr bytes.Buffer
	ext.Stdout, ext.Stderr = &out, &stderr
	if err := ext.Run(); err != nil {
		if strings.Contains(stderr.String(), "No value at that key path") {
			return nil, nil
		}
		return nil, fmt.Errorf("reading Xcode accounts: %s", strings.TrimSpace(stderr.String()))
	}
	var raw []map[string]any
	if err := json.Unmarshal(out.Bytes(), &raw); err != nil {
		return nil, fmt.Errorf("parsing Xcode accounts: %w", err)
	}
	return raw, nil
}

// Save replaces the account list, preserving each account's raw fields.
func (s *Store) Save(accts []Account) error {
	raw := make([]map[string]any, len(accts))
	for i, a := range accts {
		raw[i] = a.Raw
	}
	js, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	conv := exec.Command("plutil", "-convert", "xml1", "-o", "-", "-")
	conv.Stdin = bytes.NewReader(js)
	var xml, stderr bytes.Buffer
	conv.Stdout, conv.Stderr = &xml, &stderr
	if err := conv.Run(); err != nil {
		return fmt.Errorf("encoding Xcode accounts: %s", strings.TrimSpace(stderr.String()))
	}
	if out, err := exec.Command("defaults", "write", s.Domain, AccountsKey, xml.String()).CombinedOutput(); err != nil {
		return fmt.Errorf("writing Xcode accounts: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

func accountFromRaw(r map[string]any) Account {
	str := func(k string) string { v, _ := r[k].(string); return v }
	return Account{Username: str("Username"), HostURL: str("HostURL"), Name: str("AccountName"), Raw: r}
}

// IsGitHub reports whether the account is a github.com account.
func (a Account) IsGitHub() bool {
	return strings.EqualFold(strings.TrimRight(a.HostURL, "/"), GitHubHost)
}

// FindGitHub returns the index of user's github.com account, or -1.
func FindGitHub(accts []Account, user string) int {
	for i, a := range accts {
		if a.IsGitHub() && strings.EqualFold(a.Username, user) {
			return i
		}
	}
	return -1
}

// Promote returns accts with user's github.com account moved to the front of
// the github.com accounts (other hosts keep their relative positions). changed
// is false when user has no Xcode account or is already first.
func Promote(accts []Account, user string) (out []Account, changed bool) {
	idx := FindGitHub(accts, user)
	if idx < 0 {
		return accts, false
	}
	first := -1
	for i, a := range accts {
		if a.IsGitHub() {
			first = i
			break
		}
	}
	if idx == first {
		return accts, false
	}
	out = make([]Account, 0, len(accts))
	out = append(out, accts[:first]...)
	out = append(out, accts[idx])
	for i := first; i < len(accts); i++ {
		if i != idx {
			out = append(out, accts[i])
		}
	}
	return out, true
}

// BindURL embeds user into an https github.com remote URL so Xcode matches the
// repo to that account regardless of account order. Non-GitHub, non-https
// (ssh) and already-bound URLs are returned unchanged with ok=false.
func BindURL(remote, user string) (string, bool) {
	u, err := url.Parse(remote)
	if err != nil || u.Scheme != "https" || !strings.EqualFold(u.Hostname(), "github.com") {
		return remote, false
	}
	if u.User != nil {
		if _, hasPw := u.User.Password(); hasPw || u.User.Username() == user {
			return remote, false // never rewrite a URL carrying a credential
		}
	}
	u.User = url.User(user)
	return u.String(), true
}

// UnbindURL strips an embedded username from an https github.com remote URL.
func UnbindURL(remote string) (string, bool) {
	u, err := url.Parse(remote)
	if err != nil || u.Scheme != "https" || !strings.EqualFold(u.Hostname(), "github.com") || u.User == nil {
		return remote, false
	}
	if _, hasPw := u.User.Password(); hasPw {
		return remote, false
	}
	u.User = nil
	return u.String(), true
}
