package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Config struct {
	Global bool
}

// SwitchResult holds the outcome of a profile switch.
type SwitchResult struct {
	Warnings []string
}

func (r *SwitchResult) addWarning(msg string) {
	r.Warnings = append(r.Warnings, msg)
}

func New(global bool) *Config {
	return &Config{Global: global}
}

func (c *Config) scope() string {
	if c.Global {
		return "--global"
	}
	return "--local"
}

func (c *Config) SetUser(name, email string) error {
	if err := exec.Command("git", "config", c.scope(), "user.name", name).Run(); err != nil {
		return fmt.Errorf("set user.name: %w", err)
	}
	if err := exec.Command("git", "config", c.scope(), "user.email", email).Run(); err != nil {
		return fmt.Errorf("set user.email: %w", err)
	}
	return nil
}

// IsSSHSignKey reports whether a signing key value is an SSH key rather than a
// GPG key ID. SSH keys are either a literal public key ("ssh-ed25519 AAAA…",
// "key::…") or a path to one (~/.ssh/id_ed25519.pub); GPG keys are bare hex IDs
// or emails, which contain none of those markers.
func IsSSHSignKey(key string) bool {
	return strings.HasPrefix(key, "ssh-") ||
		strings.HasPrefix(key, "key::") ||
		strings.HasPrefix(key, "~") ||
		strings.ContainsAny(key, `/\`) ||
		strings.HasSuffix(key, ".pub")
}

// SetSignKey sets user.signingkey and the matching gpg.format, so an SSH key
// signs with SSH and a GPG key ID signs with OpenPGP. Both are unset together
// when key is empty.
func (c *Config) SetSignKey(key string) error {
	if key == "" {
		// Best-effort unset — ignore "key not set" (exit 5) but surface real failures.
		for _, k := range []string{"user.signingkey", "gpg.format"} {
			if err := exec.Command("git", "config", c.scope(), "--unset", k).Run(); err != nil && !isUnsetNothingErr(err) {
				return fmt.Errorf("unset %s: %w", k, err)
			}
		}
		return nil
	}
	if IsSSHSignKey(key) {
		switch {
		case strings.HasPrefix(key, "ssh-"):
			key = "key::" + key // git requires the key:: prefix for inline public keys
		case !strings.HasPrefix(key, "key::"):
			key = ExpandPath(key) // git does not resolve ~ in user.signingkey
		}
		if err := exec.Command("git", "config", c.scope(), "gpg.format", "ssh").Run(); err != nil {
			return fmt.Errorf("set gpg.format: %w", err)
		}
	} else if err := exec.Command("git", "config", c.scope(), "--unset", "gpg.format").Run(); err != nil && !isUnsetNothingErr(err) {
		return fmt.Errorf("unset gpg.format: %w", err)
	}
	if err := exec.Command("git", "config", c.scope(), "user.signingkey", key).Run(); err != nil {
		return fmt.Errorf("set signingkey: %w", err)
	}
	return nil
}

// ClearIdentity removes every key gitswitch writes in this config scope, so the
// scope falls back to whatever the next-wider one says. Used to undo a repo pin.
func (c *Config) ClearIdentity() error {
	for _, k := range []string{"user.name", "user.email", "user.signingkey", "gpg.format", "core.sshCommand"} {
		if err := exec.Command("git", "config", c.scope(), "--unset", k).Run(); err != nil && !isUnsetNothingErr(err) {
			return fmt.Errorf("unset %s: %w", k, err)
		}
	}
	return nil
}

// SetSSHKey sets core.sshCommand to force a specific SSH key for this config scope.
// Uses IdentitiesOnly=yes to prevent SSH agent fallback to other keys.
func (c *Config) SetSSHKey(keyPath string) error {
	if keyPath == "" {
		// Best-effort unset — ignore "key not set" (exit 5) but surface real failures.
		if err := exec.Command("git", "config", c.scope(), "--unset", "core.sshCommand").Run(); err != nil && !isUnsetNothingErr(err) {
			return fmt.Errorf("unset core.sshCommand: %w", err)
		}
		return nil
	}
	expanded := ExpandPath(keyPath)
	sshCmd := fmt.Sprintf("ssh -i %s -o IdentitiesOnly=yes", expanded)
	if err := exec.Command("git", "config", c.scope(), "core.sshCommand", sshCmd).Run(); err != nil {
		return fmt.Errorf("set core.sshCommand: %w", err)
	}
	return nil
}

// IsGitInstalled checks if git is available on PATH.
func IsGitInstalled() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// IsGHInstalled checks if the gh CLI is available on PATH.
func IsGHInstalled() bool {
	_, err := exec.LookPath("gh")
	return err == nil
}

// SwitchGHUser runs gh auth switch for the given username.
// Returns a warning string (not an error) if gh is unavailable or the switch fails —
// git config is the critical step; gh is optional.
func SwitchGHUser(ghUser string) string {
	if ghUser == "" {
		return ""
	}
	if !IsGHInstalled() {
		return "gh not installed (skipped gh auth switch)"
	}
	out, err := exec.Command("gh", "auth", "switch", "--user", ghUser).CombinedOutput()
	if err != nil {
		return fmt.Sprintf("gh auth switch failed: %s", strings.TrimSpace(string(out)))
	}
	return ""
}

// GetUser reads user.name and user.email from the given scope.
func (c *Config) GetUser() (name, email string, err error) {
	nameOut, err := exec.Command("git", "config", c.scope(), "user.name").Output()
	if err != nil {
		return "", "", nil // not set, not an error
	}
	emailOut, err := exec.Command("git", "config", c.scope(), "user.email").Output()
	if err != nil {
		return "", "", nil
	}
	return strings.TrimSpace(string(nameOut)), strings.TrimSpace(string(emailOut)), nil
}

// GetSSHKey parses the SSH key path out of core.sshCommand, e.g.
// "ssh -i ~/.ssh/id_work -o IdentitiesOnly=yes" → "~/.ssh/id_work"
func (c *Config) GetSSHKey() string {
	out, err := exec.Command("git", "config", c.scope(), "core.sshCommand").Output()
	if err != nil {
		return ""
	}
	parts := strings.Fields(strings.TrimSpace(string(out)))
	for i, p := range parts {
		if p == "-i" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

// Scope says where the identity git will actually use comes from. Ordered by
// git's own precedence: a narrower scope overrides every wider one.
type Scope int

const (
	ScopeNone    Scope = iota // no identity configured anywhere
	ScopeGlobal               // ~/.gitconfig — whichever profile is active
	ScopeRepo                 // this repo's .git/config — a gitswitch pin, or hand-set
	ScopeSession              // GIT_CONFIG_* env vars — this terminal and its children
)

// String names the scope for user-facing output.
func (s Scope) String() string {
	switch s {
	case ScopeRepo:
		return "repo"
	case ScopeSession:
		return "session"
	case ScopeGlobal:
		return "global"
	}
	return "none"
}

// ResolveIdentity reports the email git will actually author with in dir, and
// which scope supplies it. One `git config` call answers both: --show-scope
// labels every value with its source, and git lists them in precedence order,
// so the last line wins.
//
// A repo whose local user.email merely repeats the global one overrides nothing,
// so it reports ScopeGlobal — otherwise near every repo would look pinned.
func ResolveIdentity(dir string) (Scope, string) {
	out, err := exec.Command("git", "-C", dir, "config", "--show-scope", "--get-all", "user.email").Output()
	if err != nil {
		return ScopeNone, ""
	}
	var scope Scope
	var email, globalEmail string
	for _, line := range strings.Split(string(out), "\n") {
		label, value, ok := strings.Cut(strings.TrimSpace(line), "\t")
		if !ok || value == "" {
			continue
		}
		switch label {
		case "command": // GIT_CONFIG_* env vars, or git -c
			scope, email = ScopeSession, value
		case "local", "worktree":
			scope, email = ScopeRepo, value
		case "global", "system":
			scope, email = ScopeGlobal, value
			globalEmail = value
		}
	}
	if email == "" {
		return ScopeNone, ""
	}
	if scope != ScopeGlobal && strings.EqualFold(email, globalEmail) {
		return ScopeGlobal, email
	}
	return scope, email
}

// LocalEmail returns the repo-local user.email at dir, or "" if the repo has
// none. Prefer ResolveIdentity unless you specifically need the file value —
// this ignores session env vars and cannot tell a real override from a repo that
// merely repeats the global identity.
func LocalEmail(dir string) string {
	out, err := exec.Command("git", "-C", dir, "config", "--local", "user.email").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// HasLocalIdentity reports whether the repo at dir overrides the global identity
// — via a pin or a session — so there is nothing to recommend or learn for it.
func HasLocalIdentity(dir string) bool {
	scope, _ := ResolveIdentity(dir)
	return scope == ScopeRepo || scope == ScopeSession
}

// GetGHUser reads the currently active GitHub CLI account.
// Returns empty string if gh is not installed or no account is active.
func GetGHUser() string {
	out, err := exec.Command("gh", "auth", "status").CombinedOutput()
	if err != nil {
		return ""
	}
	lines := strings.Split(string(out), "\n")
	for i, line := range lines {
		if strings.Contains(line, "Active account: true") {
			// The username appears on the preceding line: "  ✓ account <username> ..."
			if i > 0 {
				prev := lines[i-1]
				fields := strings.Fields(prev)
				for j, f := range fields {
					if f == "account" && j+1 < len(fields) {
						return fields[j+1]
					}
				}
			}
		}
	}
	return ""
}

// credentialHelperValue is the entry gitswitch adds to a credential.helper
// chain. The leading '!' tells git to run it as a shell command (so it resolves
// gitswitch on PATH and appends the operation arg).
const credentialHelperValue = "!gitswitch credential"

// Two git rules govern this whole section, both verified against git 2.50:
//
//  1. For a given URL, git accumulates the values of `credential.helper` AND of
//     every matching `credential.<url>.helper` into ONE chain, in config-file
//     order — a URL-scoped key does not override the generic one, it joins it.
//  2. An empty helper value discards everything accumulated before it, and the
//     first helper that returns a credential wins.
//
// So registering gitswitch in the generic key alone is not enough: `gh auth
// setup-git` writes an empty reset plus its own helper under
// credential.https://github.com.helper, which either wipes our generic entry or
// answers ahead of it. gitswitch must be the first live entry in every chain.

// credentialHelperKeys returns every global config key that holds a credential
// helper chain — the generic one plus any URL-scoped keys another tool wrote —
// deduplicated, generic first.
func credentialHelperKeys() []string {
	keys := []string{"credential.helper"}
	seen := map[string]bool{"credential.helper": true}
	out, err := exec.Command("git", "config", "--global", "--name-only",
		"--get-regexp", `^credential\.([^=]*\.)?helper$`).Output()
	if err != nil {
		return keys // no keys set yet, or no global config
	}
	for _, line := range strings.Split(string(out), "\n") {
		if k := strings.TrimSpace(line); k != "" && !seen[k] {
			seen[k] = true
			keys = append(keys, k)
		}
	}
	return keys
}

// credentialHelperChain returns one key's values in order. Read NUL-separated so
// empty values survive — those are the resets that rule 2 above turns on, and
// silently dropping one would re-enable a helper the user's other tool disabled.
func credentialHelperChain(key string) []string {
	out, err := exec.Command("git", "config", "--global", "--get-all", "-z", key).Output()
	if err != nil || len(out) == 0 {
		return nil
	}
	return strings.Split(strings.TrimSuffix(string(out), "\x00"), "\x00")
}

// firstLiveHelper returns the index of the first value git would actually run:
// the one after the last reset, since everything before a reset is dead config.
func firstLiveHelper(values []string) int {
	live := 0
	for i, v := range values {
		if v == "" {
			live = i + 1
		}
	}
	return live
}

// helperChainIsCurrent reports whether gitswitch is already the first live entry.
func helperChainIsCurrent(values []string) bool {
	i := firstLiveHelper(values)
	return i < len(values) && values[i] == credentialHelperValue
}

// helperChainWithGitswitchFirst inserts gitswitch as the first live entry,
// keeping every other value — including resets — exactly where it was. Removing
// any existing gitswitch entry first makes this idempotent.
func helperChainWithGitswitchFirst(values []string) []string {
	stripped := make([]string, 0, len(values)+1)
	for _, v := range values {
		if v != credentialHelperValue {
			stripped = append(stripped, v)
		}
	}
	at := firstLiveHelper(stripped)
	out := make([]string, 0, len(stripped)+1)
	out = append(out, stripped[:at]...)
	out = append(out, credentialHelperValue)
	return append(out, stripped[at:]...)
}

// IsCredentialHelperInstalled reports whether gitswitch will actually be
// consulted first. False when another tool's URL-scoped helper shadows us, since
// a registration that never gets asked is not an installation — that answer is
// what makes `gitswitch install` offer to repair it.
func IsCredentialHelperInstalled() bool {
	for _, key := range credentialHelperKeys() {
		chain := credentialHelperChain(key)
		if key != "credential.helper" && len(chain) == 0 {
			continue // key vanished between the two reads
		}
		if !helperChainIsCurrent(chain) {
			return false
		}
	}
	return true
}

// HelperConflict names a chain where some other helper is asked before gitswitch,
// so gitswitch never gets to route that host's credentials.
type HelperConflict struct {
	Key    string // e.g. credential.https://github.com.helper
	Winner string // the helper git asks instead, "" when nothing is registered
}

// CredentialHelperConflicts lists the chains where gitswitch is not first, so
// `doctor` can name the tool that is shadowing it instead of just saying no.
func CredentialHelperConflicts() []HelperConflict {
	var out []HelperConflict
	for _, key := range credentialHelperKeys() {
		chain := credentialHelperChain(key)
		if helperChainIsCurrent(chain) {
			continue
		}
		if key != "credential.helper" && len(chain) == 0 {
			continue
		}
		winner := ""
		if i := firstLiveHelper(chain); i < len(chain) {
			winner = chain[i]
		}
		out = append(out, HelperConflict{Key: key, Winner: winner})
	}
	return out
}

// InstallCredentialHelper makes gitswitch the first helper git asks, for every
// host, while preserving every other helper (osxkeychain, gh's per-host entries)
// behind it. gitswitch stays silent for hosts and repos it cannot serve, so git
// falls through to those. Idempotent.
func InstallCredentialHelper() error {
	for _, key := range credentialHelperKeys() {
		existing := credentialHelperChain(key)
		if helperChainIsCurrent(existing) {
			continue
		}
		if err := rewriteHelperChain(key, existing, helperChainWithGitswitchFirst(existing)); err != nil {
			return err
		}
	}
	return nil
}

// rewriteHelperChain replaces one key's values wholesale, restoring the original
// list if any write fails — the user must never be left without their helpers.
func rewriteHelperChain(key string, existing, want []string) error {
	restore := func() {
		_ = exec.Command("git", "config", "--global", "--unset-all", key).Run()
		for _, v := range existing {
			_ = exec.Command("git", "config", "--global", "--add", key, v).Run()
		}
	}
	if err := exec.Command("git", "config", "--global", "--unset-all", key).Run(); err != nil {
		if !isUnsetNothingErr(err) {
			return fmt.Errorf("reset %s: %w", key, err)
		}
	}
	for _, v := range want {
		if err := exec.Command("git", "config", "--global", "--add", key, v).Run(); err != nil {
			restore()
			return fmt.Errorf("add %s %q: %w", key, v, err)
		}
	}
	return nil
}

// UninstallCredentialHelper removes only gitswitch's entries, from every chain,
// leaving each one otherwise as it was. Tolerates not being present.
func UninstallCredentialHelper() error {
	pattern := "^" + regexp.QuoteMeta(credentialHelperValue) + "$"
	for _, key := range credentialHelperKeys() {
		if err := exec.Command("git", "config", "--global", "--unset-all", key, pattern).Run(); err != nil {
			if !isUnsetNothingErr(err) {
				return fmt.Errorf("unset %s: %w", key, err)
			}
		}
	}
	return nil
}

// isUnsetNothingErr reports whether err is git's exit code 5, returned when an
// --unset/--unset-all targets a key/value that does not exist.
func isUnsetNothingErr(err error) bool {
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return ee.ExitCode() == 5
	}
	return false
}

// IsWorkingTreeClean reports whether there are no staged or unstaged changes.
// Reauthor requires this since it rebases.
func IsWorkingTreeClean() bool {
	out, err := exec.Command("git", "status", "--porcelain").Output()
	return err == nil && strings.TrimSpace(string(out)) == ""
}

// ResolveReauthorBase turns a rev-range argument into a rebase base.
// A bare integer N means "the last N commits" (HEAD~N); anything else
// (a ref, SHA, or HEAD~N already) is passed through untouched.
func ResolveReauthorBase(arg string) string {
	if n, err := strconv.Atoi(arg); err == nil && n > 0 {
		return fmt.Sprintf("HEAD~%d", n)
	}
	return arg
}

// Reauthor rewrites author (and committer) identity to name/email on every
// commit in (base, HEAD] whose current author email matches fromEmail. If
// fromEmail is empty, every commit in range is rewritten. Requires a clean
// working tree since it runs a non-interactive rebase.
func Reauthor(base, fromEmail, name, email string) error {
	amend := fmt.Sprintf("GIT_COMMITTER_NAME=%s GIT_COMMITTER_EMAIL=%s git commit --amend --author=%s --no-edit --no-verify",
		shellQuote(name), shellQuote(email), shellQuote(fmt.Sprintf("%s <%s>", name, email)))

	var script string
	if fromEmail != "" {
		script = fmt.Sprintf(`if [ "$(git log -1 --format=%%ae)" = %s ]; then %s; fi`, shellQuote(fromEmail), amend)
	} else {
		script = amend
	}

	cmd := exec.Command("git", "rebase", "--exec", script, base)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("rebase failed (repo left mid-rebase — run 'git rebase --abort' to undo): %w", err)
	}
	return nil
}

// PushForceWithLease force-pushes the current branch, refusing if the
// remote has moved since the local rebase started.
func PushForceWithLease() error {
	cmd := exec.Command("git", "push", "--force-with-lease")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// shellQuote wraps s in single quotes for safe use inside a sh -c script,
// escaping any embedded single quotes.
func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

// ExpandPath expands a leading ~/ to the user's home directory.
func ExpandPath(path string) string {
	if !strings.HasPrefix(path, "~/") {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	return filepath.Join(home, path[2:])
}
