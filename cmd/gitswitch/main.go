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
	gsoauth "github.com/aksisonline/gitswitch/internal/oauth"
	"github.com/aksisonline/gitswitch/internal/prereqs"
	secretsStore "github.com/aksisonline/gitswitch/internal/secrets"
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
	Short: "Switch between GitHub accounts on one machine",
	Long:  `Run without arguments to open the profile picker.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureInitialized(); err != nil {
			return err
		}
		// Quick switch: gitswitch <nickname>
		if len(args) == 1 {
			return quickSwitch(args[0])
		}
		var tuiOpts []tui.Option
		if show, notes := ver.ShouldShowWhatsNew(store.ConfigDir(), version); show {
			tuiOpts = append(tuiOpts, tui.WithWhatsNew(notes))
		}
		m, err := tui.New(store, version, tuiOpts...)
		if err != nil {
			return err
		}
		result, err := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseAllMotion()).Run()
		if err != nil {
			return err
		}
		if final, ok := result.(tui.Model); ok {
			if final.PendingReloadCmd != "" {
				fmt.Println()
				fmt.Println("  Shell integration installed. Reload your shell:")
				fmt.Println()
				fmt.Printf("    %s\n", final.PendingReloadCmd)
				fmt.Println()
			}
			if final.LaunchLogin {
				fmt.Println()
				fmt.Println("  Run  gs login  to connect your first GitHub account.")
				fmt.Println()
			}
			if final.LaunchOAuth {
				fmt.Println()
				fmt.Println("  ┌──────────────────────────────────────────┐")
				fmt.Println("  │  gitswitch · Log in with GitHub          │")
				fmt.Println("  └──────────────────────────────────────────┘")
				token, user, err := gsoauth.Login("", "")
				if err != nil {
					fmt.Printf("\n  ✗  %v\n\n", err)
					return nil
				}
				nickname := user.Login
				ref := fmt.Sprintf("gitswitch:%s:github.com", nickname)
				secrets := secretsStore.Default()
				if secrets.Available() {
					if err := secrets.Set(ref, token); err != nil {
						fmt.Printf("  ⚠  Could not store token in keychain: %v\n", err)
					}
				}
				name := user.Name
				if name == "" {
					name = user.Login
				}
				if err := store.Add(nickname, name, user.Email, "", "", user.Login); err != nil {
					_ = store.Update(nickname, storage.Profile{
						Nickname: nickname,
						UserName: name,
						Email:    user.Email,
						GHUser:   user.Login,
						TokenRef: ref,
					})
				} else {
					_ = store.Update(nickname, storage.Profile{
						Nickname: nickname,
						UserName: name,
						Email:    user.Email,
						GHUser:   user.Login,
						TokenRef: ref,
					})
				}
				profiles, _ := store.Load()
				if len(profiles) == 1 {
					_ = store.SetActive(nickname)
				}
				fmt.Printf("\n  ✓  Logged in as %s (github.com)\n", user.Login)
				fmt.Printf("  ✓  Profile %q created\n", nickname)
				if secrets.Available() {
					fmt.Println("  ✓  Token stored in keychain")
				}
				fmt.Printf("\n  Next: run  gs switch %s  to activate\n\n", nickname)
			}
		}
		return nil
	},
}

// applyProfile writes a profile's identity into one git config scope and points
// the gh CLI at the matching account, so git and GitHub never disagree about who
// you are.
func applyProfile(cfg *git.Config, p *storage.Profile) error {
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
	return nil
}

// effectiveProfile returns the identity commits in dir will actually use, and the
// scope it comes from: this terminal's session, the repo's own config, or the
// global active profile. Everything user-facing resolves through this, so the
// prompt and `gitswitch current` agree with what git will really write.
func effectiveProfile(dir string) (*storage.Profile, git.Scope, error) {
	if scope, email := git.ResolveIdentity(dir); scope == git.ScopeRepo || scope == git.ScopeSession {
		if p := store.GetByEmail(email); p != nil {
			return p, scope, nil
		}
	}
	p, err := store.GetActive()
	return p, git.ScopeGlobal, err
}

// scopeMarker is the glyph appended to the profile name in the shell prompt when
// the identity does not come from the global config. Empty for the global case —
// a user with no pins and no sessions sees nothing new.
func scopeMarker(s git.Scope) string {
	switch s {
	case git.ScopeRepo:
		return "●"
	case git.ScopeSession:
		return "◆"
	}
	return ""
}

func quickSwitch(nickname string) error {
	p, err := store.Get(nickname)
	if err != nil {
		return err
	}
	if err := applyProfile(git.New(true), p); err != nil {
		return err
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
		// Try to import current git config silently as a convenience.
		// If it fails, the TUI will show the welcome screen instead.
		_ = store.ImportCurrent()
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
		if err := applyProfile(git.New(true), p); err != nil {
			return err
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
		cwd, _ := os.Getwd()
		p, scope, err := effectiveProfile(cwd)
		if err != nil {
			return err
		}
		if p == nil {
			if !short {
				fmt.Println("No active profile")
			}
			return nil
		}
		marker := scopeMarker(scope)
		if short {
			// Starship renders this through a format string we don't control, so the
			// marker rides along on the nickname rather than as a new field.
			fmt.Printf("%s%s\t%s\n", p.Nickname, marker, p.Email)
			return nil
		}
		prompt, _ := cmd.Flags().GetBool("prompt")
		if prompt {
			prefs, _ := store.LoadPrefs()
			color := tui.ThemePromptColor(prefs.ColorTheme)
			// Third field: older installed hooks read only fields 1-2 and ignore it.
			fmt.Printf("%s\t%s\t%s\n", p.Nickname, color, marker)
			return nil
		}
		switch scope {
		case git.ScopeRepo:
			fmt.Printf("%s — %s <%s>  (pinned to this repo)\n", p.Nickname, p.UserName, p.Email)
		case git.ScopeSession:
			fmt.Printf("%s — %s <%s>  (this terminal's session)\n", p.Nickname, p.UserName, p.Email)
		default:
			fmt.Printf("%s — %s <%s>\n", p.Nickname, p.UserName, p.Email)
		}
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
		bin := filepath.Base(os.Args[0])
		fmt.Printf("%s %s\n", bin, version)
		latest := ver.CachedLatestVersion(store.ConfigDir(), version)
		if latest != "" && ver.IsUpdateAvailable(version, latest) {
			fmt.Printf("New version available: %s\n", latest)
			if isBrewInstall() {
				fmt.Printf("Run: brew upgrade gitswitch\n")
			} else {
				fmt.Printf("Run: %s upgrade\n", bin)
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
		_, err = tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseAllMotion()).Run()
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
	Short: "Check for updates and upgrade gitswitch",
	RunE: func(cmd *cobra.Command, args []string) error {
		if isBrewInstall() {
			fmt.Println("gitswitch was installed via Homebrew.")
			fmt.Println("Run: brew upgrade gitswitch")
			return nil
		}
		fmt.Println("Checking for updates...")
		latest, err := ver.FetchLatestVersionFreshFor(version)
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

var reauthorCmd = &cobra.Command{
	Use:   "reauthor <base>",
	Short: "Rewrite author/committer identity on already-made commits",
	Long: `Rewrites the author and committer of every commit between <base> and HEAD
to a stored profile's identity. <base> is a commit-ish (e.g. HEAD~3, a SHA) or
a bare number N meaning "the last N commits".

Use --from to only touch commits currently authored by a given email
(the pre-switch account), leaving other commits in range untouched.

This rewrites history. If the branch is already pushed, pass --push to
force-push (--force-with-lease) afterward, or do it yourself once you've
checked the result.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		to, _ := cmd.Flags().GetString("to")
		from, _ := cmd.Flags().GetString("from")
		push, _ := cmd.Flags().GetBool("push")
		yes, _ := cmd.Flags().GetBool("yes")
		if to == "" {
			return fmt.Errorf("--to <nickname> is required")
		}
		p, err := store.Get(to)
		if err != nil {
			return err
		}
		if !git.IsWorkingTreeClean() {
			return fmt.Errorf("working tree not clean — commit or stash changes before reauthoring")
		}
		base := git.ResolveReauthorBase(args[0])

		fmt.Printf("This rewrites commit history from %s to HEAD, setting author to %s <%s>", base, p.UserName, p.Email)
		if from != "" {
			fmt.Printf(" (commits currently authored by %s only)", from)
		}
		fmt.Println(".")
		if !yes && !confirm("Proceed?") {
			fmt.Println("Aborted.")
			return nil
		}

		if err := git.Reauthor(base, from, p.UserName, p.Email); err != nil {
			return err
		}
		fmt.Println("✓ History rewritten.")

		if push {
			if !yes && !confirm("Force-push (--force-with-lease) now?") {
				fmt.Println("Skipped push — history rewritten locally only. Push manually when ready.")
				return nil
			}
			if err := git.PushForceWithLease(); err != nil {
				return fmt.Errorf("push failed: %w", err)
			}
			fmt.Println("✓ Pushed.")
		} else {
			fmt.Println("Local only — pass --push to force-push, or push manually.")
		}
		return nil
	},
}

// confirm prompts the user with a y/N question on stdin.
func confirm(question string) bool {
	fmt.Printf("%s [y/N] ", question)
	var resp string
	fmt.Scanln(&resp)
	resp = strings.ToLower(strings.TrimSpace(resp))
	return resp == "y" || resp == "yes"
}

// A pin is a local identity: it writes the profile into this repo's
// .git/config, so commits here are correct no matter which profile is active
// globally. Nothing needs to be switched on entry.
var pinCmd = &cobra.Command{
	Use:   "pin [nickname]",
	Short: "Pin an identity to this repo — writes it to the repo's local git config",
	Long: "Pin an identity to this repo — writes it to the repo's local git config and\n" +
		"switches the gh CLI to the matching account.\n\n" +
		"With no nickname, adopts the identity the repo already has in its local git\n" +
		"config (useful for repos configured by hand before gitswitch).",
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repoKey := history.GetRepoKey()
		if repoKey == "" {
			return fmt.Errorf("not inside a git repo")
		}
		cwd, _ := os.Getwd()
		existing := git.LocalEmail(cwd)

		var p *storage.Profile
		if len(args) == 1 {
			var err error
			if p, err = store.Get(args[0]); err != nil {
				return err
			}
			// Never silently replace a per-repo identity the user set themselves.
			if existing != "" && !strings.EqualFold(existing, p.Email) {
				fmt.Printf("Replacing this repo's existing local identity <%s>\n", existing)
			}
		} else {
			// Adopt what the repo already declares.
			if existing == "" {
				return fmt.Errorf("this repo has no local git identity to adopt — run 'gitswitch pin <nickname>'")
			}
			if p = store.GetByEmail(existing); p == nil {
				return fmt.Errorf("this repo commits as <%s>, which matches no stored profile — add it first, or run 'gitswitch pin <nickname>'", existing)
			}
			fmt.Printf("Adopted this repo's existing identity <%s>\n", existing)
		}

		if err := applyProfile(git.New(false), p); err != nil {
			return err
		}
		if err := history.Pin(repoKey, p.Nickname); err != nil {
			return err
		}
		fmt.Printf("Pinned '%s' to this repo — %s <%s> (local git config; global identity unchanged)\n",
			p.Nickname, p.UserName, p.Email)
		return nil
	},
}

var unpinCmd = &cobra.Command{
	Use:   "unpin",
	Short: "Remove this repo's local identity, fall back to the global one",
	RunE: func(cmd *cobra.Command, args []string) error {
		repoKey := history.GetRepoKey()
		if repoKey == "" {
			return fmt.Errorf("not inside a git repo")
		}
		if err := git.New(false).ClearIdentity(); err != nil {
			return err
		}
		if err := history.Unpin(repoKey); err != nil {
			return err
		}
		fmt.Println("Unpinned — this repo now uses the global identity")
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
		// A repo or session carrying its own identity is already decided: don't
		// record the global profile against it (that would teach the wrong one), but
		// do point gh at the matching account so pushes from here use the right user.
		if scope, email := git.ResolveIdentity(path); scope == git.ScopeRepo || scope == git.ScopeSession {
			if p := store.GetByEmail(email); p != nil {
				// ponytail: unconditional `gh auth switch` — it is idempotent, and a
				// "is it already right?" check costs the same gh round-trip. Only
				// profiles with a gh_user pay it. Skip via a cached account if cd
				// latency ever shows up.
				_ = git.SwitchGHUser(p.GHUser) // best-effort; hooks must stay silent
			}
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
	Use:           "recommend",
	Short:         "Print recommended profile for current repo (used by shell hooks)",
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
		// A repo or terminal with its own identity (a pin, hand-set git config, or a
		// session) already commits correctly — never nudge it to switch globally.
		if git.HasLocalIdentity(path) {
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
			// Reinstall, not Install: Install is a no-op when the marker block is
			// already there, so a user told "shell integration updated — run gitswitch
			// install" would keep their old hook forever. Reinstall replaces it.
			shellResult, err = shell.Reinstall(opts.Shell, opts.Framework, "gs")
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
			cwd, _ := os.Getwd()
			return credential.Get(req, store, history.GetRepoKey(), cwd, os.Stdout)
		default:
			// store/approve/erase/reject/"" — gitswitch stores no tokens.
			return nil
		}
	},
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Connect a GitHub account",
	RunE: func(cmd *cobra.Command, args []string) error {
		host, _ := cmd.Flags().GetString("host")
		clientID, _ := cmd.Flags().GetString("client-id")
		profileName, _ := cmd.Flags().GetString("profile")

		fmt.Println()
		fmt.Println("  ┌──────────────────────────────────────────┐")
		fmt.Println("  │  gitswitch · Log in with GitHub          │")
		fmt.Println("  └──────────────────────────────────────────┘")

		token, user, err := gsoauth.Login(host, clientID)
		if err != nil {
			fmt.Println()
			fmt.Printf("  ✗  %v\n\n", err)
			return nil
		}

		nickname := profileName
		if nickname == "" {
			nickname = user.Login
		}

		// Store token in OS keychain
		ref := fmt.Sprintf("gitswitch:%s:%s", nickname, host)
		if host == "" {
			ref = fmt.Sprintf("gitswitch:%s:github.com", nickname)
		}
		secrets := secretsStore.Default()
		if secrets.Available() {
			if err := secrets.Set(ref, token); err != nil {
				fmt.Printf("  ⚠  Could not store token in keychain: %v\n", err)
			}
		}

		// Create profile
		name := user.Name
		if name == "" {
			name = user.Login
		}
		if err := store.Add(nickname, name, user.Email, "", "", user.Login); err != nil {
			// Profile exists — update TokenRef only
			_ = store.Update(nickname, storage.Profile{
				Nickname: nickname,
				UserName: name,
				Email:    user.Email,
				GHUser:   user.Login,
				TokenRef: ref,
			})
		} else {
			// Set TokenRef on the newly created profile
			_ = store.Update(nickname, storage.Profile{
				Nickname: nickname,
				UserName: name,
				Email:    user.Email,
				GHUser:   user.Login,
				TokenRef: ref,
			})
		}

		// Make first profile active
		profiles, _ := store.Load()
		if len(profiles) == 1 {
			_ = store.SetActive(nickname)
		}

		fmt.Println()
		fmt.Printf("  ✓  Logged in as %s (%s)\n", user.Login, func() string {
			if host == "" {
				return "github.com"
			}
			return host
		}())
		fmt.Printf("  ✓  Profile %q created\n", nickname)
		if secrets.Available() {
			fmt.Println("  ✓  Token stored in keychain")
		}
		fmt.Println()
		fmt.Printf("  Next: run  gs switch %s  to activate\n\n", nickname)
		return nil
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
		conflicts := git.CredentialHelperConflicts()
		if jsonOut {
			// Embedded so git/gh stay top-level, matching the pre-existing shape.
			report := struct {
				prereqs.CheckResult
				HTTPS struct {
					RoutedByGitswitch bool                 `json:"routed_by_gitswitch"`
					Conflicts         []git.HelperConflict `json:"conflicts,omitempty"`
				} `json:"https"`
			}{CheckResult: r}
			report.HTTPS.RoutedByGitswitch = len(conflicts) == 0
			report.HTTPS.Conflicts = conflicts
			b, err := json.MarshalIndent(report, "", "  ")
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", b)
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

		// HTTPS routing: being registered is not enough — another tool's per-host
		// helper can be asked first, and then gitswitch never sees the request.
		switch {
		case len(conflicts) == 0:
			fmt.Println("  ✓  HTTPS pushes routed by gitswitch")
		case conflicts[0].Winner == "" && len(conflicts) == 1:
			fmt.Println("  ⚠  HTTPS pushes not routed by gitswitch — run: gitswitch install")
		default:
			fmt.Println("  ✗  HTTPS pushes answered by another helper before gitswitch:")
			for _, c := range conflicts {
				if c.Winner == "" {
					continue
				}
				fmt.Printf("       %s → %s\n", c.Key, c.Winner)
			}
			fmt.Println("       pushes may use the wrong account — run: gitswitch install")
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
					"profiles": len(profiles),
					"git":      r.Git,
					"gh":       r.GH,
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
	rootCmd.AddCommand(addCmd, switchCmd, listCmd, removeCmd, currentCmd, initCmd, versionCmd, upgradeCmd, pacmanCmd, pinCmd, unpinCmd, recordCmd, recommendCmd, installCmd, uninstallCmd, claudeCmd, hookCheckCmd, credentialCmd, doctorCmd, setupCmd, loginCmd, betaCmd, stableCmd, reauthorCmd)
	addCmd.Flags().String("sign-key", "", "Signing key: GPG key ID, or SSH key path for gpg.format=ssh")
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
	loginCmd.Flags().String("host", "", "GitHub host (default: github.com)")
	loginCmd.Flags().String("client-id", "", "OAuth app client ID (overrides built-in)")
	loginCmd.Flags().String("profile", "", "Profile nickname to create (default: GitHub username)")
	reauthorCmd.Flags().String("to", "", "Profile nickname to attribute commits to (required)")
	reauthorCmd.Flags().String("from", "", "Only rewrite commits currently authored by this email")
	reauthorCmd.Flags().Bool("push", false, "Force-push (--force-with-lease) after rewriting")
	reauthorCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompts (for scripts and agents)")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
