package prereqs

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

type ToolCheck struct {
	Installed   bool   `json:"installed"`
	Version     string `json:"version,omitempty"`
	OK          bool   `json:"ok"`
	InstallHint string `json:"install_hint,omitempty"`
	InstallURL  string `json:"install_url,omitempty"`
}

type CheckResult struct {
	Git *ToolCheck `json:"git"`
	GH  *ToolCheck `json:"gh"`
}

func Check() CheckResult {
	return CheckResult{
		Git: checkGit(),
		GH:  checkGH(),
	}
}

// AllOK returns true when all hard requirements are met.
func (r CheckResult) AllOK() bool {
	return r.Git.OK
}

func (r CheckResult) JSON() []byte {
	b, _ := json.MarshalIndent(r, "", "  ")
	return b
}

// PrintWarnings writes human-readable prereq status to stdout.
// Hard stops (git missing) print an error; soft warnings (gh missing) are advisory.
func PrintWarnings(r CheckResult) {
	if !r.Git.Installed {
		fmt.Println()
		fmt.Println("  ✗  git not found")
		fmt.Println()
		fmt.Println("     gitswitch requires git. Install it:")
		printInstallHint(r.Git)
		fmt.Println()
		return
	}
	if !r.Git.OK {
		fmt.Printf("  ⚠  git %s found — gitswitch requires git ≥ 2.28\n", r.Git.Version)
		printInstallHint(r.Git)
	}
	if !r.GH.Installed {
		fmt.Println()
		fmt.Println("  ⚠  gh CLI not found")
		fmt.Println()
		fmt.Println("     Without it, gitswitch can't sync your GitHub auth.")
		fmt.Println("     Install it (optional but recommended):")
		printInstallHint(r.GH)
		fmt.Println()
	}
}

func printInstallHint(t *ToolCheck) {
	switch runtime.GOOS {
	case "darwin":
		fmt.Printf("       macOS:   %s\n", t.InstallHint)
	case "linux":
		fmt.Printf("       Linux:   %s\n", t.InstallURL)
	case "windows":
		fmt.Printf("       Windows: %s\n", t.InstallHint)
	default:
		fmt.Printf("       %s\n", t.InstallURL)
	}
}

func checkGit() *ToolCheck {
	t := &ToolCheck{
		InstallHint: gitInstallHint(),
		InstallURL:  "https://git-scm.com/downloads",
	}
	_, err := exec.LookPath("git")
	if err != nil {
		return t
	}
	t.Installed = true
	out, err := exec.Command("git", "--version").Output()
	if err != nil {
		return t
	}
	raw := strings.TrimPrefix(strings.TrimSpace(string(out)), "git version ")
	// strip platform suffix like " (Apple Git-137)"
	if i := strings.Index(raw, " "); i > 0 {
		raw = raw[:i]
	}
	t.Version = raw
	t.OK = semverAtLeast(raw, 2, 28)
	return t
}

func checkGH() *ToolCheck {
	t := &ToolCheck{
		InstallHint: ghInstallHint(),
		InstallURL:  "https://github.com/cli/cli/releases/latest",
	}
	_, err := exec.LookPath("gh")
	if err != nil {
		return t
	}
	t.Installed = true
	out, err := exec.Command("gh", "--version").Output()
	if err != nil {
		t.OK = true // installed but can't read version — treat as OK
		return t
	}
	// "gh version 2.52.0 (2024-06-03)"
	line := strings.TrimSpace(strings.Split(string(out), "\n")[0])
	parts := strings.Fields(line)
	if len(parts) >= 3 {
		t.Version = parts[2]
		t.OK = semverAtLeast(t.Version, 2, 0)
	} else {
		t.OK = true
	}
	return t
}

// semverAtLeast returns true when vstr (e.g. "2.45.2") is ≥ major.minor.
func semverAtLeast(vstr string, major, minor int) bool {
	var maj, min, patch int
	_, err := fmt.Sscanf(vstr, "%d.%d.%d", &maj, &min, &patch)
	if err != nil {
		// Try "major.minor" without patch
		_, err = fmt.Sscanf(vstr, "%d.%d", &maj, &min)
		if err != nil {
			return false
		}
	}
	if maj != major {
		return maj > major
	}
	return min >= minor
}

func gitInstallHint() string {
	switch runtime.GOOS {
	case "darwin":
		return "xcode-select --install"
	case "windows":
		return "winget install Git.Git"
	default:
		return "sudo apt install git  OR  sudo dnf install git"
	}
}

func ghInstallHint() string {
	switch runtime.GOOS {
	case "darwin":
		return "brew install gh"
	case "windows":
		return "winget install GitHub.cli"
	default:
		return "sudo apt install gh  OR  sudo dnf install gh"
	}
}
