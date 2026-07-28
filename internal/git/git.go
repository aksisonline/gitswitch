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

// LocalEmail returns the repo-local user.email at dir, or "" if the repo has
// none. A non-empty value means the repo commits with that identity no matter
// which profile is active globally — set by a gitswitch pin, or by hand.
func LocalEmail(dir string) string {
	out, err := exec.Command("git", "-C", dir, "config", "--local", "user.email").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// HasLocalIdentity reports whether the repo at dir carries its own identity, so
// there is nothing to recommend or learn for it.
func HasLocalIdentity(dir string) bool { return LocalEmail(dir) != "" }

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

// credentialHelperValue is the entry gitswitch adds to the global
// credential.helper chain. The leading '!' tells git to run it as a shell
// command (so it resolves gitswitch on PATH and appends the operation arg).
const credentialHelperValue = "!gitswitch credential"

// getGlobalCredentialHelpers returns the current global credential.helper
// chain as a slice of non-empty entries, in order.
func getGlobalCredentialHelpers() []string {
	out, err := exec.Command("git", "config", "--global", "--get-all", "credential.helper").Output()
	if err != nil {
		return nil
	}
	var helpers []string
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			helpers = append(helpers, line)
		}
	}
	return helpers
}

// IsCredentialHelperInstalled reports whether gitswitch is already registered
// in the global credential.helper chain.
func IsCredentialHelperInstalled() bool {
	for _, h := range getGlobalCredentialHelpers() {
		if h == credentialHelperValue {
			return true
		}
	}
	return false
}

// InstallCredentialHelper prepends gitswitch to the global credential.helper
// chain so it is consulted first, while preserving any existing helpers
// (osxkeychain, gh's per-host helper, etc.). Idempotent: a no-op if already
// installed. gitswitch stays silent for hosts/repos it cannot serve, so git
// falls through to the preserved helpers.
func InstallCredentialHelper() error {
	if IsCredentialHelperInstalled() {
		return nil
	}
	existing := getGlobalCredentialHelpers()

	// restore puts the original helper list back. Called on any failure after
	// --unset-all so the user is never left without their keychain helpers.
	restore := func() {
		_ = exec.Command("git", "config", "--global", "--unset-all", "credential.helper").Run()
		for _, e := range existing {
			_ = exec.Command("git", "config", "--global", "--add", "credential.helper", e).Run()
		}
	}

	// Wipe the list so we can re-add with gitswitch first.
	if err := exec.Command("git", "config", "--global", "--unset-all", "credential.helper").Run(); err != nil {
		if !isUnsetNothingErr(err) {
			return fmt.Errorf("reset credential.helper: %w", err)
		}
	}
	if err := exec.Command("git", "config", "--global", "--add", "credential.helper", credentialHelperValue).Run(); err != nil {
		restore()
		return fmt.Errorf("add credential.helper: %w", err)
	}
	for _, e := range existing {
		if err := exec.Command("git", "config", "--global", "--add", "credential.helper", e).Run(); err != nil {
			restore()
			return fmt.Errorf("restore credential.helper %q: %w", e, err)
		}
	}
	return nil
}

// UninstallCredentialHelper removes only gitswitch's entry from the global
// credential.helper chain, preserving all others. Tolerates the helper not
// being present.
func UninstallCredentialHelper() error {
	if err := exec.Command("git", "config", "--global", "--unset-all",
		"credential.helper", "^"+regexp.QuoteMeta(credentialHelperValue)+"$").Run(); err != nil {
		if isUnsetNothingErr(err) {
			return nil
		}
		return fmt.Errorf("unset credential.helper: %w", err)
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
