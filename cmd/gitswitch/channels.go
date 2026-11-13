package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/aksisonline/gitswitch/internal/git"
	wizard "github.com/aksisonline/gitswitch/internal/install"
	"github.com/aksisonline/gitswitch/internal/shell"
	ver "github.com/aksisonline/gitswitch/internal/version"
	skill "github.com/aksisonline/gitswitch/skills/gitswitch"
	"github.com/spf13/cobra"
)

// Updates and release channels.
//
// One rule decides everything here: whoever installed the binary is the only
// one allowed to replace it. When Homebrew installed gitswitch, every update
// and every channel change is a `brew` command — gitswitch never downloads
// over a Homebrew-managed file (that leaves brew's records and the binary on
// disk disagreeing, and is not allowed in Homebrew formulae at all). Only
// installs made by the get.gitswitch.dev script update themselves.
//
// The two channels are two formulae in the same tap:
//
//	aksisonline/tap/gitswitch        stable
//	aksisonline/tap/gitswitch-beta   beta (canary) — conflicts with stable

var upgradeCmd = &cobra.Command{
	Use:     "upgrade",
	Aliases: []string{"update"},
	Short:   "Update gitswitch to the latest release on your channel",
	Long: `Updates gitswitch to the newest release on the channel you are on (stable
or beta).

How it updates depends on how gitswitch was installed:
  • Homebrew        runs 'brew upgrade' for you — Homebrew stays in charge
  • install script  downloads the new binary from get.gitswitch.dev

Your profiles and settings (~/.config/gitswitch/) are never touched.`,
	Example: `  gitswitch upgrade            update now
  gitswitch upgrade --check    only say whether an update exists
  gitswitch upgrade --dry-run  show the command that would run`,
	Args: noArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		check, _ := cmd.Flags().GetBool("check")
		dry, _ := cmd.Flags().GetBool("dry-run")
		inst := ver.DetectInstall()

		var latest string
		err := withSpinner("Checking for updates…", func() (err error) {
			latest, err = ver.FetchLatestVersionFreshFor(version)
			return err
		})
		if err != nil {
			return fmt.Errorf("could not check for updates: %w", err)
		}
		if !ver.IsUpdateAvailable(version, latest) {
			fmt.Printf("%s Already on the latest version (%s).\n", styleOK("✓"), version)
			return nil
		}

		fmt.Printf("Update available: %s → %s\n", version, styleOK(latest))
		if check {
			fmt.Printf("Run: %s\n", inst.UpgradeHint())
			return nil
		}
		upgrade, err := ver.UpgradeCommandFor(inst, latest)
		if err != nil {
			return err
		}
		fmt.Printf("Updating with %s:\n", inst.Label())
		printStep(os.Stdout, upgrade.Args...)
		if dry {
			fmt.Println(styleDim("Dry run — nothing was changed."))
			return nil
		}
		fmt.Println()

		if configDir, err := gitswitchConfigDir(); err == nil {
			ver.ScheduleWhatsNew(configDir, latest)
		}
		if err := upgrade.Run(); err != nil {
			return fmt.Errorf("upgrade failed: %w", err)
		}
		fmt.Printf("\n%s Updated to %s.\n", styleOK("✓"), latest)
		afterUpgrade()
		return nil
	},
}

// afterUpgrade runs the follow-ups that keep the rest of the setup in step
// with a new binary.
func afterUpgrade() {
	// Installed AI-agent skills are a copy of the file embedded in the binary.
	// This process is still the OLD binary, so ask the new one to refresh them.
	if outdated, upToDate := skill.IsOutdated(); len(outdated)+len(upToDate) > 0 {
		if newBin, err := exec.LookPath("gitswitch"); err == nil {
			_ = runCommand(newBin, "skills", "--update")
		}
	}

	// If shell integration is on but HTTPS routing never was, offer it now —
	// it is the step people on older installs are missing.
	sh := shell.DetectShell()
	if shell.IsInstalled(shell.RCFile(sh)) && !git.IsCredentialHelperInstalled() {
		fmt.Println("\nNew setup options are available in this version — opening setup…")
		opts, err := wizard.Run(wizard.Config{HTTPSDefault: true, GuardDefault: true})
		if err == nil && opts.InstallHTTPS {
			if herr := git.InstallCredentialHelper(); herr != nil {
				fmt.Printf("  %s could not register HTTPS credential helper: %v\n", styleWarn("warning:"), herr)
			} else {
				wizard.PrintSummary(os.Stdout, "", false, true, false, nil)
			}
		}
		return
	}
	fmt.Println("Open a new terminal (or restart gitswitch) to use it.")
}

var betaCmd = &cobra.Command{
	Use:   "beta",
	Short: "Switch to the beta (canary) release channel",
	Long: `Moves this install to the beta channel: early builds of the next release.

With Homebrew this swaps the 'gitswitch' formula for 'gitswitch-beta' from the
same tap, so Homebrew keeps managing updates. With the install script it
replaces the binary with the latest beta build.

Your profiles and settings (~/.config/gitswitch/) are never touched.
Run 'gitswitch stable' to switch back.`,
	Example: `  gitswitch beta            switch (asks first)
  gitswitch beta --yes      switch without asking
  gitswitch beta --dry-run  show the commands that would run`,
	Args: noArgs,
	RunE: func(cmd *cobra.Command, args []string) error { return switchChannel(cmd, true) },
}

var stableCmd = &cobra.Command{
	Use:   "stable",
	Short: "Switch back to the stable release channel",
	Long: `Moves this install back to the stable channel.

With Homebrew this swaps 'gitswitch-beta' for the 'gitswitch' formula. With the
install script it replaces the binary with the latest stable release.

Your profiles and settings (~/.config/gitswitch/) are never touched.`,
	Example: `  gitswitch stable            switch (asks first)
  gitswitch stable --yes      switch without asking
  gitswitch stable --dry-run  show the commands that would run`,
	Args: noArgs,
	RunE: func(cmd *cobra.Command, args []string) error { return switchChannel(cmd, false) },
}

func channelName(beta bool) string {
	if beta {
		return "beta"
	}
	return "stable"
}

// switchChannel moves this install to the beta or stable channel.
func switchChannel(cmd *cobra.Command, toBeta bool) error {
	yes, _ := cmd.Flags().GetBool("yes")
	dry, _ := cmd.Flags().GetBool("dry-run")
	inst := ver.DetectInstall()

	onBeta := ver.IsBeta(version) || inst.Formula == ver.BetaFormula
	if onBeta == toBeta {
		fmt.Printf("Already on the %s channel (%s).\n", channelName(toBeta), version)
		fmt.Printf("Run 'gitswitch %s' to switch to the %s channel.\n", channelName(!toBeta), channelName(!toBeta))
		return nil
	}

	var target string
	err := withSpinner(fmt.Sprintf("Checking for the latest %s release…", channelName(toBeta)), func() (err error) {
		if toBeta {
			target, err = ver.FetchLatestCanaryVersion()
		} else {
			target, err = ver.FetchLatestVersionFresh()
		}
		return err
	})
	if err != nil {
		return fmt.Errorf("could not find a %s release: %w", channelName(toBeta), err)
	}
	if target == "" {
		return fmt.Errorf("no %s release is published yet", channelName(toBeta))
	}
	fmt.Printf("Latest %s: %s  (you have %s)\n\n", channelName(toBeta), styleOK(target), version)
	// Right after a stable release the newest beta is the one that became it.
	// "Switching to beta" would then be a downgrade — say so instead.
	if toBeta && !ver.IsUpdateAvailable(version, target) {
		fmt.Printf("There is no beta newer than %s yet — staying on stable.\n", version)
		fmt.Println(styleDim("Betas appear as work on the next release lands. Try again later."))
		return nil
	}

	if inst.IsBrew() {
		return brewChannelSwitch(inst, toBeta, target, yes, dry)
	}
	return scriptChannelSwitch(toBeta, target, yes, dry)
}

// brewChannelSwitch swaps one Homebrew formula for the other. The two
// formulae conflict (same binary), so the old one is removed first — and put
// back if installing the new one fails, so a failed switch never leaves the
// machine without gitswitch.
func brewChannelSwitch(inst ver.Install, toBeta bool, target string, yes, dry bool) error {
	to := ver.StableFormula
	if toBeta {
		to = ver.BetaFormula
	}
	fromRef, toRef := inst.FormulaRef(), ver.BrewTap+"/"+to

	fmt.Println("gitswitch is managed by Homebrew, so the switch goes through brew:")
	printStep(os.Stdout, "brew", "update")
	printStep(os.Stdout, "brew", "uninstall", inst.Formula)
	printStep(os.Stdout, "brew", "install", toRef)
	fmt.Println(styleDim("Your profiles and settings (~/.config/gitswitch/) are not touched."))
	if dry {
		fmt.Println(styleDim("Dry run — nothing was changed."))
		return nil
	}
	if _, err := exec.LookPath("brew"); err != nil {
		return fmt.Errorf("`brew` is not on PATH — run the commands above from a shell where Homebrew is available")
	}
	fmt.Println()
	if !yes && !confirmPrompt("Proceed? [y/N] ") {
		fmt.Println("Cancelled — nothing was changed.")
		return nil
	}

	fmt.Println()
	if err := runCommand("brew", "update", "--quiet"); err != nil {
		return fmt.Errorf("brew update failed: %w", err)
	}
	// Make sure the destination formula exists BEFORE removing anything.
	if err := exec.Command("brew", "info", toRef).Run(); err != nil {
		return fmt.Errorf("Homebrew has no %s formula yet — nothing was changed", toRef)
	}
	if configDir, err := gitswitchConfigDir(); err == nil {
		ver.ScheduleWhatsNew(configDir, target)
	}
	if err := runCommand("brew", "uninstall", inst.Formula); err != nil {
		return fmt.Errorf("brew uninstall %s failed: %w — nothing else was changed", inst.Formula, err)
	}
	if err := runCommand("brew", "install", toRef); err != nil {
		fmt.Printf("\n%s Installing %s failed — putting %s back…\n", styleWarn("⚠"), to, inst.Formula)
		if rerr := runCommand("brew", "install", fromRef); rerr != nil {
			return fmt.Errorf("brew install %s failed (%v), and restoring %s failed too — run: brew install %s", toRef, err, inst.Formula, fromRef)
		}
		return fmt.Errorf("brew install %s failed: %w — %s was restored", toRef, err, inst.Formula)
	}

	fmt.Printf("\n%s Switched to the %s channel (%s).\n", styleOK("✓"), channelName(toBeta), target)
	fmt.Printf("  Updates now arrive with: brew upgrade %s\n", to)
	if toBeta {
		fmt.Println("  Run 'gitswitch stable' at any time to return to stable.")
	}
	return nil
}

// scriptChannelSwitch replaces a script-installed binary with the target
// channel's latest build.
func scriptChannelSwitch(toBeta bool, target string, yes, dry bool) error {
	install, err := ver.UpgradeCommand(target)
	if err != nil {
		return err
	}
	fmt.Printf("This replaces the current binary with the %s build:\n", channelName(toBeta))
	printStep(os.Stdout, install.Args...)
	fmt.Println(styleDim("Your profiles and settings (~/.config/gitswitch/) are not touched."))
	if dry {
		fmt.Println(styleDim("Dry run — nothing was changed."))
		return nil
	}
	fmt.Println()
	if !yes && !confirmPrompt("Proceed? [y/N] ") {
		fmt.Println("Cancelled — nothing was changed.")
		return nil
	}

	fmt.Println()
	if configDir, err := gitswitchConfigDir(); err == nil {
		ver.ScheduleWhatsNew(configDir, target)
	}
	if err := install.Run(); err != nil {
		return fmt.Errorf("%s install failed: %w", channelName(toBeta), err)
	}
	fmt.Printf("\n%s Switched to the %s channel (%s).\n", styleOK("✓"), channelName(toBeta), target)
	if toBeta {
		fmt.Println("  Run 'gitswitch stable' at any time to return to stable.")
	}
	fmt.Println("  Open a new terminal to use the new binary.")
	return nil
}

// printStep prints a command the way a person would type it. `sh -c "<x>"`
// (how compound commands are run) is shown as just <x>.
func printStep(w io.Writer, argv ...string) {
	if len(argv) == 3 && argv[0] == "sh" && argv[1] == "-c" {
		argv = argv[2:]
	}
	fmt.Fprintf(w, "  %s %s\n", styleDim("$"), strings.Join(argv, " "))
}

// confirmPrompt prints the prompt and returns true only if the user types y or yes.
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
