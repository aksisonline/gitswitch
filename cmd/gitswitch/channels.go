package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	ver "github.com/aksisonline/gitswitch/internal/version"
	"github.com/spf13/cobra"
)

var betaCmd = &cobra.Command{
	Use:   "beta",
	Short: "Switch to the beta (canary) release channel",
	Long: `Downloads and installs the latest canary pre-release of gitswitch.

If gitswitch was installed via Homebrew, the formula will be removed first and
replaced with the canary binary. Your profiles and settings at
~/.config/gitswitch/ are untouched.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if ver.IsCanaryBuild(version) {
			fmt.Println("Already on the canary channel (" + version + ").")
			fmt.Println("Run 'gitswitch stable' to switch back to the stable channel.")
			return nil
		}
		return switchToCanary()
	},
}

var stableCmd = &cobra.Command{
	Use:   "stable",
	Short: "Switch back to the stable release channel",
	Long: `Reinstalls the latest stable release of gitswitch.

Use this to move off the canary channel. Your profiles and settings at
~/.config/gitswitch/ are preserved.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !ver.IsCanaryBuild(version) {
			fmt.Println("Already on the stable channel (" + version + ").")
			fmt.Println("Run 'gitswitch beta' to switch to the canary channel.")
			return nil
		}
		return switchToStable()
	},
}

// switchToCanary migrates the user from any install (brew or script) to the
// latest canary pre-release binary.
func switchToCanary() error {
	fmt.Println("Checking for the latest canary release...")
	canaryVersion, err := ver.FetchLatestCanaryVersion()
	if err != nil {
		return fmt.Errorf("could not find a canary release: %w", err)
	}

	fmt.Printf("\nLatest canary: %s (current: %s)\n", canaryVersion, version)
	fmt.Println()

	brewInstall := isBrewInstall()
	if brewInstall {
		fmt.Println("gitswitch is currently installed via Homebrew.")
		fmt.Println("Switching to the canary channel requires:")
		fmt.Println("  1. Uninstalling the Homebrew formula (brew uninstall gitswitch)")
		fmt.Println("  2. Installing the canary binary to /usr/local/bin")
		fmt.Println()
		fmt.Println("Your profiles (~/.config/gitswitch/) will NOT be affected.")
	} else {
		fmt.Println("This will replace the current binary with the canary build.")
		fmt.Println("Your profiles (~/.config/gitswitch/) will NOT be affected.")
	}

	fmt.Println()
	if !confirmPrompt("Proceed? [y/N] ") {
		fmt.Println("Cancelled.")
		return nil
	}

	if brewInstall {
		fmt.Println()
		fmt.Println("Removing Homebrew formula...")
		if err := runCommand("brew", "uninstall", "gitswitch"); err != nil {
			return fmt.Errorf("brew uninstall failed: %w\n\nRun manually: brew uninstall gitswitch", err)
		}
		fmt.Println("✓ Homebrew formula removed")
	}

	fmt.Println()
	fmt.Printf("Installing canary %s...\n", canaryVersion)
	if err := runInstallScript(canaryVersion); err != nil {
		return fmt.Errorf("canary install failed: %w", err)
	}

	fmt.Println()
	fmt.Println("✓ Switched to canary channel.")
	fmt.Println("  Run 'gitswitch stable' at any time to return to stable.")
	fmt.Println("  Restart your shell (or open a new terminal) to use the new binary.")
	return nil
}

// switchToStable migrates a canary install back to the latest stable release.
func switchToStable() error {
	fmt.Println("Checking for the latest stable release...")
	stableVersion, err := ver.FetchLatestVersionFresh()
	if err != nil {
		return fmt.Errorf("could not fetch stable version: %w", err)
	}

	fmt.Printf("\nLatest stable: %s (current: %s)\n", stableVersion, version)
	fmt.Println()
	fmt.Println("This will replace the canary binary with the latest stable release.")
	fmt.Println("Your profiles (~/.config/gitswitch/) will NOT be affected.")
	fmt.Println()
	fmt.Println("If you prefer Homebrew after this, run:")
	fmt.Println("  brew install aksisonline/tap/gitswitch")
	fmt.Println()

	if !confirmPrompt("Proceed? [y/N] ") {
		fmt.Println("Cancelled.")
		return nil
	}

	fmt.Println()
	fmt.Printf("Installing stable %s...\n", stableVersion)
	if err := runInstallScript(stableVersion); err != nil {
		return fmt.Errorf("stable install failed: %w", err)
	}

	fmt.Println()
	fmt.Println("✓ Switched to stable channel.")
	fmt.Println("  Restart your shell (or open a new terminal) to use the new binary.")
	return nil
}

// confirmPrompt prints the prompt and returns true only if the user types y or Y.
func confirmPrompt(prompt string) bool {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return false
	}
	answer := strings.TrimSpace(scanner.Text())
	return strings.EqualFold(answer, "y") || strings.EqualFold(answer, "yes")
}

// runCommand executes a command, streaming stdout/stderr to the terminal.
func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// runInstallScript downloads and runs the gitswitch install script for targetVersion.
func runInstallScript(targetVersion string) error {
	cmd, err := ver.UpgradeCommand(targetVersion)
	if err != nil {
		return err
	}
	return cmd.Run()
}
