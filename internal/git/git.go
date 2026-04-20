package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

func (c *Config) SetSignKey(key string) error {
	if key == "" {
		return nil
	}
	if err := exec.Command("git", "config", c.scope(), "user.signingkey", key).Run(); err != nil {
		return fmt.Errorf("set signingkey: %w", err)
	}
	return nil
}

// SetSSHKey sets core.sshCommand to force a specific SSH key for this config scope.
// Uses IdentitiesOnly=yes to prevent SSH agent fallback to other keys.
func (c *Config) SetSSHKey(keyPath string) error {
	if keyPath == "" {
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
