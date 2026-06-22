package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aksisonline/gitswitch/internal/credential"
	"github.com/aksisonline/gitswitch/internal/git"
	"github.com/aksisonline/gitswitch/internal/history"
	wizard "github.com/aksisonline/gitswitch/internal/install"
	"github.com/aksisonline/gitswitch/internal/prereqs"
	"github.com/aksisonline/gitswitch/internal/shell"
	"github.com/aksisonline/gitswitch/internal/storage"
	"github.com/aksisonline/gitswitch/internal/tui"
	ver "github.com/aksisonline/gitswitch/internal/version"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

//go:embed skill/SKILL.md
var skillMD []byte

var version = "dev"

var store *storage.Store

func init() {
	if !git.IsGitInstalled() {
		fmt.Fprintf(os.Stderr, "Error: git is not installed or not on PATH\n")
		os.Exit(1)
	}
	var err error
	store, err = storage.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "gitswitch [nickname]",
	Short: "Manage git profiles — run without args for interactive UI",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureInitialized(); err != nil {
			return err
		}
		// Quick switch: gitswitch <nickname>
		if len(args) == 1 {
			return quickSwitch(args[0])
		}
		m, err := tui.New(store, version)
		if err != nil {
			return err
		}
		_, err = tea.NewProgram(m, tea.WithAltScreen()).Run()
		return err
	},
}

func quickSwitch(nickname string) error {
	p, err := store.Get(nickname)
	if err != nil {
		return err
	}
	cfg := git.New(true)
	if err := cfg.SetUser(p.UserName, p.Email); err != nil {
		return err
	}
	if err := cfg.SetSignKey(p.SignKey); err != nil {
		return err
	}
	if err := cfg.SetSSHKey(p.SSHKey); err != nil {
		return err
	}
	if w := git.SwitchGHUser(p.GHUser); w != "" {
		fmt.Printf("warning: %s\n", w)
	}
	if err := store.SetActive(p.Nickname); err != nil {
		return err
	}
	fmt.Printf("✓ Switched to '%s' — %s <%s>\n", p.Nickname, p.UserName, p.Email)
	return nil
}

func ensureInitialized() error {
	profiles, err := store.Load()
	if err != nil {
		return err
	}
	if len(profiles) == 0 {
		fmt.Println("First startup: importing existing git config...")
		if err := store.ImportCurrent(); err != nil {
			fmt.Printf("Could not auto-import: %v\n", err)
			fmt.Println("Tip: gitswitch add <nickname> <user-name> <email>")
			return nil
		}
		fmt.Println("✓ Imported as 'default' profile")
	}
	return nil
}

var addCmd = &cobra.Command{
	Use:   "add <nickname> <user-name> <email>",
	Short: "Add new profile",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		signKey, _ := cmd.Flags().GetString("sign-key")
		sshKey, _ := cmd.Flags().GetString("ssh-key")
		ghUser, _ := cmd.Flags().GetString("gh-user")
		if err := store.Add(args[0], args[1], args[2], signKey, sshKey, ghUser); err != nil {
			return err
		}
		fmt.Printf("Profile '%s' added\n", args[0])
		return nil
	},
}

var switchCmd = &cobra.Command{
	Use:   "switch <nickname>",
	Short: "Switch to profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := store.Get(args[0])
		if err != nil {
			return err
		}
		cfg := git.New(true)
		if err := cfg.SetUser(p.UserName, p.Email); err != nil {
			return err
		}
		if err := cfg.SetSignKey(p.SignKey); err != nil {
			return err
		}
		if err := cfg.SetSSHKey(p.SSHKey); err != nil {
			return err
		}
		if w := git.SwitchGHUser(p.GHUser); w != "" {
			fmt.Printf("warning: %s\n", w)
		}
		if err := store.SetActive(p.Nickname); err != nil {
			return err
		}
		fmt.Printf("Switched to '%s' — %s <%s>\n", p.Nickname, p.UserName, p.Email)
		return nil
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		profiles, err := store.Load()
		if err != nil {
			return err
		}
		if len(profiles) == 0 {
			fmt.Println("No profiles")
			return nil
		}
		for _, p := range profiles {
			prefix := " "
			if p.Active {
				prefix = "✓"
			}
			fmt.Printf("%s  %-14s  %s <%s>\n", prefix, p.Nickname, p.UserName, p.Email)
		}
		return nil
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove <nickname>",
	Short: "Remove profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := store.Remove(args[0]); err != nil {
			return err
		}
		fmt.Printf("Profile '%s' removed\n", args[0])
		return nil
	},
}

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show current profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		short, _ := cmd.Flags().GetBool("short")
		p, err := store.GetActive()
		if err != nil {
			return err
		}
		if p == nil {
			if !short {
				fmt.Println("No active profile")
			}
			return nil
		}
		if short {
			fmt.Printf("%s\t%s\n", p.Nickname, p.Email)
			return nil
		}
		prompt, _ := cmd.Flags().GetBool("prompt")
		if prompt {
			prefs, _ := store.LoadPrefs()
			color := tui.ThemePromptColor(prefs.ColorTheme)
			fmt.Printf("%s\t%s\n", p.Nickname, color)
			return nil
		}
		fmt.Printf("%s — %s <%s>\n", p.Nickname, p.UserName, p.Email)
		if git.IsCredentialHelperInstalled() {
			fmt.Println("HTTPS credential helper: active")
		}
		return nil
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Import existing git config",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := store.ImportCurrent(); err != nil {
			return err
		}
		fmt.Println("✓ Imported current git config as 'default' profile")
		return nil
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show current version and check for updates",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("gitswitch %s\n", version)
		latest := ver.CachedLatestVersion(store.ConfigDir())
		if latest != "" && ver.IsUpdateAvailable(version, latest) {
			fmt.Printf("New version available: %s\n", latest)
			if isBrewInstall() {
				fmt.Println("Run: brew upgrade gitswitch")
			} else {
				fmt.Println("Run: gitswitch upgrade")
			}
		} else if latest != "" {
			fmt.Println("Already on latest version.")
		}
		return nil
	},
}

var pacmanCmd = &cobra.Command{
	Use:   "pacman",
	Short: "Launch Git-Switcher with arcade intro animation",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureInitialized(); err != nil {
			return err
		}
		m, err := tui.New(store, version, tui.WithArcadeMode())
		if err != nil {
			return err
		}
		_, err = tea.NewProgram(m, tea.WithAltScreen()).Run()
		return err
	},
}

// gitswitchConfigDir returns ~/.config/gitswitch, or an error when the home
// directory cannot be determined.
func gitswitchConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}
	return filepath.Join(home, ".config", "gitswitch"), nil
}

// isBrewInstall reports whether the running binary lives inside a Homebrew
// Cellar by resolving symlinks and checking the path.
func isBrewInstall() bool {
	exe, err := os.Executable()
	if err != nil {
		return false
	}
	resolved, err := filepath.EvalSymlinks(exe)
	if err != nil {
		resolved = exe
	}
	return strings.Contains(resolved, "/Cellar/") ||
		strings.Contains(resolved, "/linuxbrew/")
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade gitswitch to the latest version",
	RunE: func(cmd *cobra.Command, args []string) error {
		if isBrewInstall() {
			fmt.Println("gitswitch was installed via Homebrew.")
			fmt.Println("Run: brew upgrade gitswitch")
			return nil
		}
		fmt.Println("Checking for updates...")
		latest, err := ver.FetchLatestVersionFresh()
		if err != nil {
			return fmt.Errorf("could not fetch latest version: %w", err)
		}
		if !ver.IsUpdateAvailable(version, latest) {
			fmt.Printf("Already on latest version (%s).\n", version)
			return nil
		}
		fmt.Printf("Upgrading %s → %s...\n", version, latest)
		if err := ver.RunUpgrade(latest); err != nil {
			return fmt.Errorf("upgrade failed: %w", err)
		}
		fmt.Println("✓ Upgrade complete.")
		fmt.Println()

		// If shell is installed but credential helper is not, run the wizard
		// immediately so the user can activate new features without a separate step.
		sh := shell.DetectShell()
		if shell.IsInstalled(shell.RCFile(sh)) && !git.IsCredentialHelperInstalled() {
			fmt.Println("New features available in this version.")
			fmt.Println("Launching setup to activate them...")
			fmt.Println()
			opts, err := wizard.Run(wizard.Config{HTTPSDefault: true}, os.Stdout)
			if err == nil && opts.InstallHTTPS {
				if herr := git.InstallCredentialHelper(); herr != nil {
					fmt.Printf("  warning: could not register HTTPS credential helper: %v\n", herr)
				} else {
					wizard.PrintSummary(os.Stdout, "", false, true, nil)
				}
			}
		} else {
			fmt.Println("Restart your shell to use the new version.")
		}
		return nil
	},
}

var pinCmd = &cobra.Command{
	Use:   "pin <nickname>",
	Short: "Pin an identity to this repo — always recommended regardless of usage history",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := store.Get(args[0]); err != nil {
			return err
		}
		repoKey := history.GetRepoKey()
		if repoKey == "" {
			return fmt.Errorf("not inside a git repo")
		}
		if err := history.Pin(repoKey, args[0]); err != nil {
			return err
		}
		fmt.Printf("Pinned '%s' to this repo\n", args[0])
		return nil
	},
}

var unpinCmd = &cobra.Command{
	Use:   "unpin",
	Short: "Remove pinned identity for this repo, fall back to auto-recommendation",
	RunE: func(cmd *cobra.Command, args []string) error {
		repoKey := history.GetRepoKey()
		if repoKey == "" {
			return fmt.Errorf("not inside a git repo")
		}
		if err := history.Unpin(repoKey); err != nil {
			return err
		}
		fmt.Println("Unpinned — identity recommendation now based on usage history")
		return nil
	},
}

var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "Record current identity for this repo (called by shell hooks on repo entry)",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString("path")
		if path == "" {
			var err error
			path, err = os.Getwd()
			if err != nil {
				return err
			}
		}
		repoKey := history.GetRepoKeyForPath(path)
		if repoKey == "" {
			return nil
		}
		active, err := store.GetActive()
		if err != nil || active == nil {
			return nil
		}
		return history.Record(repoKey, active.Nickname)
	},
}

// errNoRecommendation signals a silent exit 1 from recommendCmd.
// SilenceErrors on the command prevents cobra from printing it.
var errNoRecommendation = fmt.Errorf("")

var recommendCmd = &cobra.Command{
	Use:          "recommend",
	Short:        "Print recommended profile for current repo (used by shell hooks)",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString("path")
		if path == "" {
			var err error
			path, err = os.Getwd()
			if err != nil {
				return errNoRecommendation
			}
		}

		repoKey := history.GetRepoKeyForPath(path)
		if repoKey == "" {
			return errNoRecommendation
		}

		active, _ := store.GetActive()
		currentNick := ""
		if active != nil {
			currentNick = active.Nickname
		}

		nick, ok := history.Recommend(repoKey, currentNick)
		if !ok {
			return errNoRecommendation
		}

		p, err := store.Get(nick)
		if err != nil {
			return errNoRecommendation
		}
		fmt.Printf("%s\t%s\t%s\n", p.Nickname, p.UserName, p.Email)
		return nil
	},
}

var claudeCmd = &cobra.Command{
	Use:   "claude",
	Short: "Install the gitswitch skill into Claude Code",
	RunE: func(cmd *cobra.Command, args []string) error {
		scope, _ := cmd.Flags().GetString("scope")

		var base string
		if scope == "project" {
			base = ".claude"
		} else {
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			base = filepath.Join(home, ".claude")
		}

		dest := filepath.Join(base, "skills", "gitswitch")
		if err := os.MkdirAll(dest, 0755); err != nil {
			return fmt.Errorf("could not create skills directory: %w", err)
		}

		skillPath := filepath.Join(dest, "SKILL.md")
		if err := os.WriteFile(skillPath, skillMD, 0644); err != nil {
			return fmt.Errorf("could not write skill: %w", err)
		}

		fmt.Printf("✓ Skill installed to %s\n", dest)
		fmt.Println("  Reload Claude Code (or open a new session) to activate.")
		return nil
	},
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Set up gitswitch — shell integration and HTTPS credential routing",
	Long: `Interactive setup wizard. Detects your shell, shows what each step does,
and asks before making any changes. Use --yes to accept all defaults without
prompts (for scripts and CI).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		shellFlag, _ := cmd.Flags().GetString("shell")
		yes, _ := cmd.Flags().GetBool("yes")
		httpsDefault, _ := cmd.Flags().GetBool("https")

		opts, err := wizard.Run(wizard.Config{
			ShellOverride: shellFlag,
			Yes:           yes,
			HTTPSDefault:  httpsDefault,
		}, os.Stdout)
		if err != nil {
			return fmt.Errorf("setup interrupted: %w", err)
		}

		var shellResult string
		if opts.InstallShell {
			shellResult, err = shell.Install(opts.Shell, opts.Framework)
			if err != nil {
				return fmt.Errorf("shell install failed: %w", err)
			}
			configDir, err := gitswitchConfigDir()
			if err != nil {
				return err
			}
			_ = shell.WriteHookVersion(configDir, version)
		}

		var httpsErr error
		if opts.InstallHTTPS {
			httpsErr = git.InstallCredentialHelper()
		}

		wizard.PrintSummary(os.Stdout, shellResult, opts.InstallShell, opts.InstallHTTPS && httpsErr == nil, httpsErr)
		return nil
	},
}

var hookCheckCmd = &cobra.Command{
	Use:    "hook-check",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		configDir, err := gitswitchConfigDir()
		if err != nil {
			return err
		}
		rcFile := shell.RCFile(shell.DetectShell())
		if msg := shell.HookUpdateMessage(configDir, rcFile, version, git.IsCredentialHelperInstalled()); msg != "" {
			fmt.Println(msg)
		}
		return nil
	},
}

// credentialCmd is a git credential helper (registered as
// `credential.helper = !gitswitch credential`). git invokes it with an
// operation arg (get/store/erase) and pipes the credential description on
// stdin. For `get`, gitswitch resolves the active/pinned profile for the
// current repo and delegates to gh to fetch that account's token — it stores
// nothing itself. store/erase are no-ops (gh owns its storage). On any case it
// cannot serve, it writes nothing and exits 0 so git falls through.
var credentialCmd = &cobra.Command{
	Use:           "credential [get|store|erase]",
	Hidden:        true,
	Args:          cobra.MaximumNArgs(1),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		op := ""
		if len(args) == 1 {
			op = args[0]
		}
		switch op {
		case "get", "fill":
			req, err := credential.ParseRequest(os.Stdin)
			if err != nil {
				return nil // graceful: exit 0, no output
			}
			return credential.Get(req, store, history.GetRepoKey(), os.Stdout)
		default:
			// store/approve/erase/reject/"" — gitswitch stores no tokens.
			return nil
		}
	},
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove shell integration written by 'gitswitch install'",
	RunE: func(cmd *cobra.Command, args []string) error {
		shellFlag, _ := cmd.Flags().GetString("shell")

		var sh shell.Shell
		switch shellFlag {
		case "zsh":
			sh = shell.ShellZsh
		case "bash":
			sh = shell.ShellBash
		case "fish":
			sh = shell.ShellFish
		default:
			sh = shell.DetectShell()
		}

		fw := shell.DetectFramework()

		result, err := shell.Uninstall(sh, fw)
		if err != nil {
			return fmt.Errorf("uninstall failed: %w", err)
		}
		fmt.Printf("✓ %s\n", result)

		if git.IsCredentialHelperInstalled() {
			if err := git.UninstallCredentialHelper(); err != nil {
				fmt.Printf("  warning: could not remove HTTPS credential helper: %v\n", err)
			} else {
				fmt.Println("✓ HTTPS credential helper removed")
			}
		}

		fmt.Println("  Reload your shell (or open a new terminal) to complete removal.")
		return nil
	},
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check that git and gh are installed and up to date",
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonOut, _ := cmd.Flags().GetBool("json")
		r := prereqs.Check()
		if jsonOut {
			fmt.Printf("%s\n", r.JSON())
			return nil
		}
		fmt.Println()
		if r.Git.Installed {
			fmt.Printf("  ✓  git %s\n", r.Git.Version)
		} else {
			fmt.Println("  ✗  git  not found")
		}
		if r.GH.Installed {
			fmt.Printf("  ✓  gh  %s\n", r.GH.Version)
		} else {
			fmt.Println("  ⚠  gh   not found (optional)")
		}
		fmt.Println()
		prereqs.PrintWarnings(r)
		return nil
	},
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Check requirements and show next steps",
	RunE: func(cmd *cobra.Command, args []string) error {
		agentMode, _ := cmd.Flags().GetBool("agent")
		r := prereqs.Check()

		if agentMode {
			profiles, _ := store.Load()
			manifest := map[string]interface{}{
				"tool":    "gitswitch",
				"version": version,
				"state": map[string]interface{}{
					"profiles":    len(profiles),
					"git":         r.Git,
					"gh":          r.GH,
				},
			}
			b, _ := json.MarshalIndent(manifest, "", "  ")
			fmt.Printf("%s\n", b)
			return nil
		}

		fmt.Println()
		fmt.Println("  Checking requirements...")
		fmt.Println()
		if r.Git.Installed {
			fmt.Printf("  ✓  git %s\n", r.Git.Version)
		} else {
			fmt.Println("  ✗  git  not found")
		}
		if r.GH.Installed {
			fmt.Printf("  ✓  gh  %s\n", r.GH.Version)
		} else {
			fmt.Println("  ⚠  gh   not found")
		}
		fmt.Println()
		prereqs.PrintWarnings(r)

		if r.AllOK() {
			profiles, _ := store.Load()
			if len(profiles) == 0 {
				fmt.Println("  No accounts yet.  Run  gs login  to get started.")
				fmt.Println()
			} else {
				fmt.Printf("  %d profile(s) configured.  Run  gs  to open the picker.\n", len(profiles))
				fmt.Println()
			}
		}
		return nil
	},
}

func main() {
	rootCmd.AddCommand(addCmd, switchCmd, listCmd, removeCmd, currentCmd, initCmd, versionCmd, upgradeCmd, pacmanCmd, pinCmd, unpinCmd, recordCmd, recommendCmd, installCmd, uninstallCmd, claudeCmd, hookCheckCmd, credentialCmd, doctorCmd, setupCmd)
	addCmd.Flags().String("sign-key", "", "GPG signing key (git user.signingkey)")
	addCmd.Flags().String("ssh-key", "", "SSH private key path, e.g. ~/.ssh/id_work (sets core.sshCommand)")
	addCmd.Flags().String("gh-user", "", "GitHub CLI username (for gh auth switch)")
	currentCmd.Flags().Bool("short", false, "Output nickname and email tab-separated (for Starship and scripts)")
	currentCmd.Flags().Bool("prompt", false, "Output nickname and theme color tab-separated (for shell prompt functions)")
	recordCmd.Flags().String("path", "", "Directory to record for (default: current working directory)")
	recommendCmd.Flags().String("path", "", "Directory to check (default: current working directory)")
	installCmd.Flags().String("shell", "", "Shell to target: zsh, bash, or fish (default: auto-detect; also skips interactive wizard)")
	installCmd.Flags().Bool("https", true, "Register HTTPS credential helper (default: true; prompted interactively when omitted)")
	installCmd.Flags().BoolP("yes", "y", false, "Accept all defaults without prompts (for scripts and CI)")
	uninstallCmd.Flags().String("shell", "", "Shell to uninstall for: zsh, bash, or fish (default: auto-detect)")
	claudeCmd.Flags().String("scope", "user", "Install scope: 'user' (~/.claude/skills) or 'project' (.claude/skills)")
	doctorCmd.Flags().Bool("json", false, "Output machine-readable JSON")
	setupCmd.Flags().Bool("agent", false, "Emit machine-readable setup manifest for AI agents")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
