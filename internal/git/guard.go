package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	guardMarker     = "# gitswitch pre-commit guard"
	passthroughMark = "# gitswitch hook passthrough"
	prevHooksFile   = ".prev-hooks-path"
)

// passthroughHooks are standard client-side git hooks beside pre-commit. When
// core.hooksPath points at ~/.config/gitswitch/hooks, git stops looking in the
// repo's own .git/hooks/ directory for ALL hooks — so we install thin forwarder
// shims for these so a repo's local commit-msg or pre-push hook still runs.
var passthroughHooks = []string{
	"applypatch-msg",
	"pre-applypatch",
	"post-applypatch",
	"prepare-commit-msg",
	"commit-msg",
	"post-commit",
	"pre-rebase",
	"post-checkout",
	"post-merge",
	"pre-push",
	"post-rewrite",
}

// GuardHooksDir returns ~/.config/gitswitch/hooks.
func GuardHooksDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}
	return filepath.Join(home, ".config", "gitswitch", "hooks"), nil
}

func preCommitScript() string {
	return `#!/bin/sh
` + guardMarker + `
if command -v gitswitch >/dev/null 2>&1; then
  gitswitch guard check || exit $?
fi
prev_file="$(dirname "$0")/` + prevHooksFile + `"
if [ -f "$prev_file" ]; then
  prev_dir=$(cat "$prev_file")
  if [ -n "$prev_dir" ] && [ -x "$prev_dir/pre-commit" ]; then
    "$prev_dir/pre-commit" "$@" || exit $?
  fi
fi
git_dir=$(git rev-parse --git-dir 2>/dev/null)
if [ -n "$git_dir" ] && [ -x "$git_dir/hooks/pre-commit" ]; then
  exec "$git_dir/hooks/pre-commit" "$@"
fi
exit 0
`
}

func passthroughScript(hookName string) string {
	return `#!/bin/sh
` + passthroughMark + `
prev_file="$(dirname "$0")/` + prevHooksFile + `"
if [ -f "$prev_file" ]; then
  prev_dir=$(cat "$prev_file")
  if [ -n "$prev_dir" ] && [ -x "$prev_dir/` + hookName + `" ]; then
    "$prev_dir/` + hookName + `" "$@" || exit $?
  fi
fi
git_dir=$(git rev-parse --git-dir 2>/dev/null)
if [ -n "$git_dir" ] && [ -x "$git_dir/hooks/` + hookName + `" ]; then
  exec "$git_dir/hooks/` + hookName + `" "$@"
fi
exit 0
`
}

// IsGuardHookInstalled reports whether git's global core.hooksPath points to
// gitswitch's hooks directory and the pre-commit guard script is in place.
func IsGuardHookInstalled() bool {
	hooksDir, err := GuardHooksDir()
	if err != nil {
		return false
	}
	out, err := exec.Command("git", "config", "--global", "core.hooksPath").Output()
	if err != nil {
		return false
	}
	configured := ExpandPath(strings.TrimSpace(string(out)))
	if filepath.Clean(configured) != filepath.Clean(hooksDir) {
		return false
	}
	data, err := os.ReadFile(filepath.Join(hooksDir, "pre-commit"))
	if err != nil {
		return false
	}
	return strings.Contains(string(data), guardMarker)
}

// InstallGuardHook writes the pre-commit guard and passthrough shims into
// ~/.config/gitswitch/hooks and sets git's global core.hooksPath. Any prior
// non-gitswitch core.hooksPath is recorded and chained to before local hooks.
func InstallGuardHook() error {
	hooksDir, err := GuardHooksDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return fmt.Errorf("create hooks dir: %w", err)
	}

	// Preserve a pre-existing custom global core.hooksPath so we can chain to it
	// and restore it on uninstall.
	if out, err := exec.Command("git", "config", "--global", "core.hooksPath").Output(); err == nil {
		existing := strings.TrimSpace(string(out))
		expanded := ExpandPath(existing)
		if existing != "" && filepath.Clean(expanded) != filepath.Clean(hooksDir) {
			_ = os.WriteFile(filepath.Join(hooksDir, prevHooksFile), []byte(expanded), 0644)
		}
	}

	preCommitPath := filepath.Join(hooksDir, "pre-commit")
	if err := os.WriteFile(preCommitPath, []byte(preCommitScript()), 0755); err != nil {
		return fmt.Errorf("write pre-commit hook: %w", err)
	}
	if err := os.Chmod(preCommitPath, 0755); err != nil {
		return fmt.Errorf("chmod pre-commit hook: %w", err)
	}

	for _, name := range passthroughHooks {
		p := filepath.Join(hooksDir, name)
		if err := os.WriteFile(p, []byte(passthroughScript(name)), 0755); err != nil {
			return fmt.Errorf("write %s hook: %w", name, err)
		}
		_ = os.Chmod(p, 0755)
	}

	if err := exec.Command("git", "config", "--global", "core.hooksPath", hooksDir).Run(); err != nil {
		return fmt.Errorf("set core.hooksPath: %w", err)
	}
	return nil
}

// UninstallGuardHook removes gitswitch's global core.hooksPath (restoring any
// previous core.hooksPath that was saved on install) and deletes the managed
// hook scripts.
func UninstallGuardHook() error {
	hooksDir, err := GuardHooksDir()
	if err != nil {
		return err
	}

	prevFile := filepath.Join(hooksDir, prevHooksFile)
	prevBytes, prevErr := os.ReadFile(prevFile)
	prevPath := strings.TrimSpace(string(prevBytes))

	out, cfgErr := exec.Command("git", "config", "--global", "core.hooksPath").Output()
	if cfgErr == nil {
		configured := ExpandPath(strings.TrimSpace(string(out)))
		if filepath.Clean(configured) == filepath.Clean(hooksDir) {
			if prevErr == nil && prevPath != "" {
				if err := exec.Command("git", "config", "--global", "core.hooksPath", prevPath).Run(); err != nil {
					return fmt.Errorf("restore core.hooksPath: %w", err)
				}
			} else {
				if err := exec.Command("git", "config", "--global", "--unset", "core.hooksPath").Run(); err != nil && !isUnsetNothingErr(err) {
					return fmt.Errorf("unset core.hooksPath: %w", err)
				}
			}
		}
	}

	_ = os.Remove(prevFile)
	_ = os.Remove(filepath.Join(hooksDir, "pre-commit"))
	for _, name := range passthroughHooks {
		_ = os.Remove(filepath.Join(hooksDir, name))
	}
	return nil
}

// ResolveCommitEmail returns the effective scope and author email that a commit
// in dir will use. When GIT_AUTHOR_EMAIL is explicitly set to something other
// than the config-resolved email (e.g. `git commit --author=...` or an env
// override), that override wins and is reported as ScopeSession; when git itself
// exports GIT_AUTHOR_EMAIL matching the config email during hook execution, the
// underlying config scope (ScopeGlobal / ScopeRepo) is preserved.
func ResolveCommitEmail(dir string) (Scope, string) {
	scope, cfgEmail := ResolveIdentity(dir)
	if envEmail := strings.TrimSpace(os.Getenv("GIT_AUTHOR_EMAIL")); envEmail != "" {
		if !strings.EqualFold(envEmail, cfgEmail) {
			return ScopeSession, envEmail
		}
	}
	return scope, cfgEmail
}
