// Package install provides the interactive gitswitch setup wizard.
// It runs when `gitswitch install` is called from a real terminal with no
// explicit flags. When piped/CI/--yes, it skips prompts and uses defaults.
package install

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aksisonline/gitswitch/internal/git"
	"github.com/aksisonline/gitswitch/internal/shell"
	"github.com/mattn/go-isatty"
)

// Options carries the resolved choices from a wizard run.
type Options struct {
	Shell            shell.Shell
	Framework        shell.Framework
	InstallShell     bool
	InstallHTTPS     bool
}

// Config controls wizard behaviour.
type Config struct {
	// ShellOverride is set when --shell flag is provided. Empty = auto-detect.
	ShellOverride string
	// Yes skips all prompts, accepting defaults (non-interactive).
	Yes bool
	// HTTPSDefault is the default for the HTTPS step (from --https flag).
	HTTPSDefault bool
}

// Run runs the wizard and returns the user's choices.
// When not interactive (piped, --yes, CI), it returns safe defaults immediately.
func Run(cfg Config, out io.Writer) (Options, error) {
	sh, fw := resolveShell(cfg.ShellOverride)
	opts := Options{Shell: sh, Framework: fw, InstallShell: true, InstallHTTPS: cfg.HTTPSDefault}

	interactive := isatty.IsTerminal(os.Stdin.Fd()) && !cfg.Yes && cfg.ShellOverride == ""

	if !interactive {
		return opts, nil
	}

	reader := bufio.NewReader(os.Stdin)

	printHeader(out)

	// ── Step 1: Shell integration ─────────────────────────────────────────────
	printStep(out, 1, 2, "Shell integration")
	printShellInfo(out, sh, fw)
	fmt.Fprintln(out)
	install, err := prompt(reader, out, "Install shell integration?", true)
	if err != nil {
		return opts, err
	}
	opts.InstallShell = install

	fmt.Fprintln(out)

	// ── Step 2: HTTPS credential routing ─────────────────────────────────────
	printStep(out, 2, 2, "HTTPS credential routing")
	printHTTPSInfo(out)
	fmt.Fprintln(out)
	httpsInstall, err := prompt(reader, out, "Register HTTPS credential helper?", true)
	if err != nil {
		return opts, err
	}
	opts.InstallHTTPS = httpsInstall

	fmt.Fprintln(out)

	return opts, nil
}

// ── helpers ──────────────────────────────────────────────────────────────────

func resolveShell(override string) (shell.Shell, shell.Framework) {
	var sh shell.Shell
	switch override {
	case "zsh":
		sh = shell.ShellZsh
	case "bash":
		sh = shell.ShellBash
	case "fish":
		sh = shell.ShellFish
	default:
		sh = shell.DetectShell()
	}
	return sh, shell.DetectFramework()
}

func printHeader(w io.Writer) {
	fmt.Fprintln(w)
	fmt.Fprintln(w, "  gitswitch setup")
	fmt.Fprintln(w, "  ───────────────────────────────────────────────────────")
	fmt.Fprintln(w)
}

func printStep(w io.Writer, n, total int, title string) {
	fmt.Fprintf(w, "  Step %d of %d — %s\n", n, total, title)
	fmt.Fprintln(w, "  ───────────────────────────────────────────────────────")
}

func shellName(sh shell.Shell) string {
	switch sh {
	case shell.ShellZsh:
		return "zsh"
	case shell.ShellBash:
		return "bash"
	case shell.ShellFish:
		return "fish"
	default:
		return "bash"
	}
}

func frameworkName(fw shell.Framework) string {
	switch fw {
	case shell.FrameworkOMZ:
		return "oh-my-zsh"
	case shell.FrameworkP10k:
		return "Powerlevel10k"
	case shell.FrameworkStarship:
		return "Starship"
	default:
		return ""
	}
}

func printShellInfo(w io.Writer, sh shell.Shell, fw shell.Framework) {
	name := shellName(sh)
	fwName := frameworkName(fw)
	if fwName != "" {
		fmt.Fprintf(w, "  Detected: %s  (with %s)\n", name, fwName)
	} else {
		fmt.Fprintf(w, "  Detected: %s\n", name)
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "  What it adds:")
	fmt.Fprintln(w, "    • Active profile shown in your prompt at all times")
	fmt.Fprintln(w, "    • A nudge when you cd into a repo usually worked on with a different identity")
	fmt.Fprintln(w, "    • Automatic usage tracking so gitswitch learns your patterns per repo")

	fmt.Fprintln(w)
	fmt.Fprintln(w, "  Before:")
	fmt.Fprintln(w, "    user@machine ~/code/work $")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "  After:")
	fmt.Fprintln(w, "    user@machine ~/code/work [work] $")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "    → Switching repos: \"You usually work here as 'personal'. Switch? [y/N]\"")
}

func printHTTPSInfo(w io.Writer) {
	fmt.Fprintln(w, "  When you push over HTTPS, git needs to know which account's token to use.")
	fmt.Fprintln(w, "  Without this, the wrong GitHub account can silently authenticate your push.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "  How it works:")
	fmt.Fprintln(w, "    • gitswitch registers as git's credential helper")
	fmt.Fprintln(w, "    • It reads the active/pinned profile for the repo you're pushing from")
	fmt.Fprintln(w, "    • It fetches that account's token from gh CLI (nothing stored by gitswitch)")
	fmt.Fprintln(w, "    • git authenticates as the right account automatically")
	fmt.Fprintln(w)

	if !git.IsGHInstalled() {
		fmt.Fprintln(w, "  ⚠  gh CLI not found — helper will install but stay inert until gh is set up")
	} else {
		fmt.Fprintln(w, "  Requires: gh CLI (found ✓)")
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "  Before:")
	fmt.Fprintln(w, "    $ git push  →  git prompts for username/password or uses wrong account")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "  After:")
	fmt.Fprintln(w, "    $ git push  →  gitswitch routes the right token for this repo, silently")
}

// prompt prints "[Y/n]" or "[y/N]" and reads a y/n response.
// Returns the default if the user just hits enter.
func prompt(r *bufio.Reader, w io.Writer, question string, defaultYes bool) (bool, error) {
	indicator := "[Y/n]"
	if !defaultYes {
		indicator = "[y/N]"
	}
	fmt.Fprintf(w, "  %s %s: ", question, indicator)

	line, err := r.ReadString('\n')
	if err != nil && err != io.EOF {
		return defaultYes, err
	}
	line = strings.TrimSpace(strings.ToLower(line))
	switch line {
	case "y", "yes":
		return true, nil
	case "n", "no":
		return false, nil
	default:
		return defaultYes, nil
	}
}

// PrintSummary prints the final "what was done" block.
func PrintSummary(w io.Writer, shellResult string, shellInstalled, httpsInstalled bool, httpsErr error) {
	fmt.Fprintln(w, "  ───────────────────────────────────────────────────────")
	fmt.Fprintln(w, "  Setup complete")
	fmt.Fprintln(w, "  ───────────────────────────────────────────────────────")
	fmt.Fprintln(w)

	if shellInstalled {
		fmt.Fprintf(w, "  ✓ Shell integration — %s\n", shellResult)
	} else {
		fmt.Fprintln(w, "  – Shell integration skipped")
	}

	if httpsErr != nil {
		fmt.Fprintf(w, "  ⚠ HTTPS credential helper — error: %v\n", httpsErr)
	} else if httpsInstalled {
		fmt.Fprintln(w, "  ✓ HTTPS credential helper registered")
	} else {
		fmt.Fprintln(w, "  – HTTPS credential helper skipped")
	}

	fmt.Fprintln(w)
	if shellInstalled {
		fmt.Fprintln(w, "  Reload your shell (or open a new terminal) to activate.")
	}
	fmt.Fprintln(w)

	if shellInstalled || httpsInstalled {
		fmt.Fprintln(w, "  Next:")
		if shellInstalled {
			fmt.Fprintln(w, "    gitswitch current   — see your active profile")
			fmt.Fprintln(w, "    gitswitch pin       — pin a repo to always use a specific profile")
		}
		fmt.Fprintln(w)
	}
}
