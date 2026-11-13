package version

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Homebrew coordinates. Both formulae live in the same tap and conflict with
// each other (they install the same `gitswitch` binary), so exactly one is
// ever installed.
const (
	BrewTap       = "aksisonline/tap"
	StableFormula = "gitswitch"
	BetaFormula   = "gitswitch-beta"
)

// Method is how this binary got onto the machine — which decides who is
// allowed to replace it.
type Method int

const (
	// MethodScript is the get.gitswitch.dev install script (or a hand-placed
	// binary): gitswitch owns the file and may replace it itself.
	MethodScript Method = iota
	// MethodBrew means Homebrew owns the file. gitswitch must never overwrite
	// it — every update and channel change goes through `brew`, otherwise
	// Homebrew's records and the binary on disk disagree.
	MethodBrew
)

// Install describes how the running binary is managed.
type Install struct {
	Method  Method
	Formula string // Homebrew formula name; empty unless Method == MethodBrew
}

// DetectInstall inspects the running binary's real path.
func DetectInstall() Install {
	exe, err := os.Executable()
	if err != nil {
		return Install{}
	}
	if resolved, err := filepath.EvalSymlinks(exe); err == nil {
		exe = resolved
	}
	return installFromPath(exe)
}

// installFromPath reads the formula name out of a Homebrew Cellar path:
// <prefix>/Cellar/<formula>/<version>/bin/gitswitch.
func installFromPath(path string) Install {
	parts := strings.Split(filepath.ToSlash(path), "/")
	for i, part := range parts {
		if part == "Cellar" && i+1 < len(parts) && parts[i+1] != "" {
			return Install{Method: MethodBrew, Formula: parts[i+1]}
		}
	}
	if strings.Contains(filepath.ToSlash(path), "/linuxbrew/") {
		return Install{Method: MethodBrew, Formula: StableFormula}
	}
	return Install{}
}

// IsBrew reports whether Homebrew manages this install.
func (i Install) IsBrew() bool { return i.Method == MethodBrew }

// FormulaRef is the fully-qualified formula (tap/name), which works whether
// or not the user has other taps shadowing the short name.
func (i Install) FormulaRef() string { return BrewTap + "/" + i.Formula }

// UpgradeHint is the command a person would type to update this install.
func (i Install) UpgradeHint() string {
	if i.IsBrew() {
		return "brew upgrade " + i.Formula
	}
	return "gitswitch upgrade"
}

// Label names the update mechanism in a few words, for UI copy.
func (i Install) Label() string {
	if i.IsBrew() {
		return "Homebrew"
	}
	return "install script"
}

// UpgradeCommandFor returns the command that updates this install to
// targetVersion, wired to the terminal. Homebrew installs run `brew update`
// first: brew only refreshes taps about once a day on its own, so without it
// a release from the last few hours looks like "already installed".
func UpgradeCommandFor(i Install, targetVersion string) (*exec.Cmd, error) {
	if !i.IsBrew() {
		return UpgradeCommand(targetVersion)
	}
	if !semverRe.MatchString(targetVersion) {
		return nil, fmt.Errorf("invalid version format: %q", targetVersion)
	}
	if _, err := exec.LookPath("brew"); err != nil {
		return nil, fmt.Errorf("gitswitch was installed with Homebrew, but `brew` is not on PATH — run: %s", i.UpgradeHint())
	}
	cmd := exec.Command("sh", "-c", "brew update --quiet && brew upgrade "+i.FormulaRef())
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd, nil
}
