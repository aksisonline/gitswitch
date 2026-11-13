package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
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
	skill "github.com/aksisonline/gitswitch/skills/gitswitch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var version = "dev"

var store *storage.Store

// commandAllowsMissingGit is true when a command can run without Git: `doctor`,
// whose whole job is checking/installing it; bare invocation, which runs the
// same check-and-install step before the TUI launches; the self-contained
// hidden snake game; and everything that only talks about gitswitch itself —
// help, the version, updates, completion scripts.
func commandAllowsMissingGit(args []string) bool {
	if len(args) == 1 {
		return true
	}
	switch args[1] {
	case "doctor", "snake",
		"help", "--help", "-h",
		"version", "--version", "-v", "--v", "-V",
		"upgrade", "update", "beta", "stable",
		"completion", "__complete", "__completeNoDesc":
		return true
	}
	return false
}

func init() {
	if !git.IsGitInstalled() && !commandAllowsMissingGit(os.Args) {
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
	Short: "Switch git identities and GitHub accounts on one machine",
	Long: `Run without arguments to open the profile picker. First run on a machine with
no profiles yet checks for git/gh automatically and offers to install them,
then walks you through connecting a GitHub account — nothing else to run first.`,
	Args:              argRange(0, 1),
	SilenceErrors:     true,
	SilenceUsage:      true,
	ValidArgsFunction: completeNicknames,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureInitialized(); err != nil {
			return err
		}
		// Quick switch: gitswitch <nickname>
		if len(args) == 1 {
			if _, err := store.Get(args[0]); err != nil {
				profiles, _ := store.Load()
				return unknownTargetError(cmd.Root(), args[0], profiles)
			}
			return quickSwitch(args[0])
		}
		var whatsNewBody string
		if profiles, _ := store.Load(); len(profiles) == 0 {
			if err := ensurePrereqsBootstrap(); err != nil {
				return err
			}
			// First run: straight into setup. "What's new" means nothing to
			// someone who has never seen the old version, and fetching the
			// notes would hold a blank screen before the wizard appears.
			ver.MarkVersionSeen(store.ConfigDir(), version)
		} else if show, notes := ver.ShouldShowWhatsNew(store.ConfigDir(), version); show {
			whatsNewBody = notes
		}
		m, err := tui.New(store, version, whatsNewBody)
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
				return completeLogin("", "")
			}
		}
		return nil
	},
}

// completeLogin hands off straight to `gh auth login`, then builds or
// updates a profile from whichever account it left active on host.
// nickname, when empty, defaults to the GitHub username. Shared by the
// onboarding wizard's "Log in with GitHub" step (rootCmd) and the explicit
// `login` command — same flow, only the host/nickname inputs differ. Login
// failures are printed, not returned, matching both callers' original
// behavior of never surfacing a login error as a cobra usage error.
func completeLogin(host, nickname string) error {
	fmt.Println()
	fmt.Println("  ┌──────────────────────────────────────────┐")
	fmt.Println("  │  gitswitch · Log in with GitHub          │")
	fmt.Println("  └──────────────────────────────────────────┘")
	fmt.Println()

	if !git.IsGHInstalled() {
		fmt.Println("  ✗  gh CLI is required for gitswitch login — install it and try again")
		return nil
	}

	displayHost := host
	if displayHost == "" {
		displayHost = "github.com"
	}

	// Capture the active account(s) before login so we can tell which
	// account gh just authenticated afterwards — on a host that already
	// has an active account, gh may leave the old one active, and we must
	// not misattribute the new profile to it.
	before := git.ListGHUsers()

	args := []string{"auth", "login", "--scopes", "user:email"}
	if host != "" {
		args = append(args, "--hostname", host)
	}
	c := exec.Command("gh", args...)
	c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := c.Run(); err != nil {
		fmt.Println()
		fmt.Printf("  ✗  gh auth login failed: %v\n\n", err)
		return nil
	}

	acct, ok := pickLoggedInAccount(before, git.ListGHUsers(), displayHost)
	if !ok {
		fmt.Println()
		fmt.Println("  ⚠  Logged in, but couldn't detect the new account — run `gh auth status` to check")
		return nil
	}

	if nickname == "" {
		nickname = acct.Login
	}

	// verifiedEmail is "" when the lookup fails (missing scope, network error,
	// timeout) — kept separate from the noreply fallback below so a failed
	// lookup on re-login doesn't clobber an already-correct profile email
	// (mergeOAuthUpdate only preserves the existing value when this is blank).
	// IdentityFor fetches the gh token once and reuses it for both the name
	// and email lookups — a single `gh auth token` subprocess per account.
	name, verified := git.IdentityFor(acct.Login, displayHost)
	verifiedEmail := ""
	if len(verified) > 0 {
		verifiedEmail = verified[0]
	}
	email := verifiedEmail
	if email == "" {
		email = acct.Login + "@users.noreply.github.com"
	}
	if name == "" {
		name = acct.Login
	}

	if err := store.Add(nickname, name, email, "", "", acct.Login); err != nil {
		// Profile exists — merge fresh login fields, keep SSH/GPG keys.
		existing, _ := store.Get(nickname)
		_ = store.Update(nickname, mergeOAuthUpdate(existing, storage.Profile{
			Nickname: nickname,
			UserName: name,
			Email:    verifiedEmail,
			GHUser:   acct.Login,
		}))
	} else {
		_ = store.Update(nickname, storage.Profile{
			Nickname: nickname,
			UserName: name,
			Email:    email,
			GHUser:   acct.Login,
		})
	}

	// Make first profile active
	profiles, _ := store.Load()
	if len(profiles) == 1 {
		_ = store.SetActive(nickname)
	}

	fmt.Println()
	fmt.Printf("  ✓  Logged in as %s (%s)\n", acct.Login, displayHost)
	fmt.Printf("  ✓  Profile %q created\n", nickname)
	fmt.Println()
	fmt.Printf("  Next: run  gs switch %s  to activate\n\n", nickname)
	return nil
}

// activeAccountForHost picks the account in accounts that is active on
// host — git.GetGHUser() won't do here since it returns the first active
// account across ALL hosts, wrong when other hosts already have accounts
// active.

// pickLoggedInAccount returns the account on host that gh just
// authenticated via `gh auth login`. A newly-activated account (active
// afterwards but not before) wins, so logging in a second account on a
// host that already had one active doesn't misattribute the new profile
// to the previously-active account. When the same account was merely
// re-authenticated, it was active both before and after, so the
// currently-active one is returned as a fallback.
func pickLoggedInAccount(before, after []git.GHAccount, host string) (git.GHAccount, bool) {
	beforeActive := make(map[string]bool)
	for _, a := range before {
		if a.Host == host && a.Active {
			beforeActive[a.Login] = true
		}
	}
	for _, a := range after {
		if a.Host == host && a.Active && !beforeActive[a.Login] {
			return a, true
		}
	}
	for _, a := range after {
		if a.Host == host && a.Active {
			return a, true
		}
	}
	return git.GHAccount{}, false
}
func activeAccountForHost(accounts []git.GHAccount, host string) (git.GHAccount, bool) {
	for _, a := range accounts {
		if a.Host == host && a.Active {
			return a, true
		}
	}
	return git.GHAccount{}, false
}

// mergeOAuthUpdate builds the profile to save after a login/re-login: the fresh
// OAuth fields win, but SSH/GPG keys and any name/email already on the existing
// profile survive (OAuth never carries them, so a blank here means "unset").
func mergeOAuthUpdate(existing *storage.Profile, updated storage.Profile) storage.Profile {
	if existing == nil {
		return updated
	}
	updated.SignKey = existing.SignKey
	updated.SSHKey = existing.SSHKey
	if updated.UserName == "" {
		updated.UserName = existing.UserName
	}
	if updated.Email == "" {
		updated.Email = existing.Email
	}
	return updated
}

// applyGitConfig writes a profile's git identity (name, email, signing key,
// SSH key) into the given git config scope without touching the gh CLI.
func applyGitConfig(cfg *git.Config, p *storage.Profile) error {
	if err := cfg.SetUser(p.UserName, p.Email); err != nil {
		return err
	}
	if err := cfg.SetSignKey(p.SignKey); err != nil {
		return err
	}
	return cfg.SetSSHKey(p.SSHKey)
}

// applyProfile writes a profile's identity into one git config scope and points
// the gh CLI at the matching account, so git and GitHub never disagree about who
// you are. Returns a non-fatal warning (e.g. gh auth switch failed) instead of
// printing it directly — some callers (shell hooks) must stay silent.
func applyProfile(cfg *git.Config, p *storage.Profile) (warning string, err error) {
	if err := applyGitConfig(cfg, p); err != nil {
		return "", err
	}
	return git.SwitchGHUser(p.GHUser), nil
}

// effectiveProfile returns the identity commits in dir will actually use, and the
// scope it comes from: this terminal's session, the repo's own config, or the
// global active profile. Everything user-facing resolves through this, so the
// prompt and `gitswitch current` agree with what git will really write.
func effectiveProfile(dir string) (*storage.Profile, git.Scope, error) {
	if scope, email := git.ResolveIdentity(dir); scope == git.ScopeRepo || scope == git.ScopeSession {
		if scope == git.ScopeRepo && !shell.IsGHWrapperInstalled(shell.RCFile(shell.DetectShell())) {
			// Session Isolation is off — the repo-local override doesn't count as active.
			p, err := store.GetActive()
			return p, git.ScopeGlobal, err
		}
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

// quickSwitch switches the active global profile and, when invoked inside a
// repo that has a local identity or pin, also updates that repo's local git
// config so git, gh, and HTTPS pushes in the current repo immediately use the
// explicitly switched profile while preserving the remembered pin in history.json.
func quickSwitch(nickname string) error { return switchProfile(nickname, false) }

// switchProfile is quickSwitch with an option to print nothing on success
// (warnings and errors still surface), for scripts.
func switchProfile(nickname string, quiet bool) error {
	p, err := store.Get(nickname)
	if err != nil {
		return err
	}
	warn, err := applyProfile(git.New(true), p)
	if err != nil {
		return err
	}
	cwd, _ := os.Getwd()
	repoKey := history.GetRepoKeyForPath(cwd)
	pinnedNick := history.GetPinned(repoKey)
	if repoKey != "" && (pinnedNick != "" || git.LocalEmail(cwd) != "") {
		if err := applyGitConfig(git.New(false), p); err != nil {
			return err
		}
	}
	if err := store.SetActive(p.Nickname); err != nil {
		return err
	}
	if warn != "" {
		fmt.Fprintf(os.Stderr, "warning: %s\n", warn)
	}
	if quiet {
		xcodeAutoSync(p)
		return nil
	}
	if pinnedNick != "" && pinnedNick != p.Nickname {
		fmt.Printf("%s Switched to '%s' — %s <%s> (overriding repo pin '%s')\n", styleOK("✓"), p.Nickname, p.UserName, p.Email, pinnedNick)
	} else {
		fmt.Printf("%s Switched to '%s' — %s <%s>\n", styleOK("✓"), p.Nickname, p.UserName, p.Email)
	}
	if msg := xcodeAutoSync(p); msg != "" {
		fmt.Printf("%s %s\n", styleOK("✓"), msg)
	}
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

// ensurePrereqsBootstrap runs once, the first time gitswitch is invoked bare
// on a machine with zero profiles, before the TUI wizard launches. It's the
// same git/gh check-and-install `doctor` already does (reused, not
// duplicated) — this can't happen inside the bubbletea wizard itself because
// installing git/gh may need a sudo password prompt or (macOS git) pop the
// Xcode CLT GUI installer, which needs plain blocking terminal I/O.
func ensurePrereqsBootstrap() error {
	r := prereqs.Check()
	if r.Git.Installed && r.GH.Installed {
		return nil
	}
	fmt.Println()
	fmt.Println("  First run — checking git and gh...")
	prereqs.PrintWarnings(r)
	if !r.Git.Installed {
		if err := promptInstallTool("git", r.Git, false); err != nil {
			return err
		}
	}
	if !r.GH.Installed {
		if err := promptInstallTool("gh", r.GH, false); err != nil {
			return err
		}
	}
	return nil
}

var addCmd = &cobra.Command{
	Use:   "add <nickname> <user-name> <email>",
	Short: "Add a profile",
	Long: `Adds a profile: a nickname plus the name and email git should commit with.

To connect a GitHub account instead (name and email are filled in for you),
use 'gitswitch login'.`,
	Example: `  gitswitch add work "Ada Lovelace" ada@corp.example
  gitswitch add work "Ada Lovelace" ada@corp.example --gh-user ada-corp --ssh-key ~/.ssh/id_work`,
	Args: exactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		signKey, _ := cmd.Flags().GetString("sign-key")
		sshKey, _ := cmd.Flags().GetString("ssh-key")
		ghUser, _ := cmd.Flags().GetString("gh-user")
		if err := store.Add(args[0], args[1], args[2], signKey, sshKey, ghUser); err != nil {
			return err
		}
		fmt.Printf("%s Profile '%s' added\n", styleOK("✓"), args[0])
		return nil
	},
}

var switchCmd = &cobra.Command{
	Use:     "switch <nickname>",
	Aliases: []string{"use"},
	Short:   "Switch to a profile",
	Long: `Makes a profile the active identity: sets git's global name, email, signing
key and SSH key, and points the gh CLI at the matching account.

'gitswitch <nickname>' is the short form of this command.`,
	Example: `  gitswitch switch work
  gitswitch work              same thing, shorter
  gitswitch switch work -q    no output on success (for scripts)`,
	Args:              exactArgs(1),
	ValidArgsFunction: completeNicknames,
	RunE: func(cmd *cobra.Command, args []string) error {
		quiet, _ := cmd.Flags().GetBool("quiet")
		return switchProfile(args[0], quiet)
	},
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List profiles",
	Args:    noArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		profiles, err := store.Load()
		if err != nil {
			return err
		}
		jsonOut, _ := cmd.Flags().GetBool("json")
		if jsonOut {
			type profileJSON struct {
				Nickname string `json:"nickname"`
				UserName string `json:"user_name"`
				Email    string `json:"email"`
				GHUser   string `json:"gh_user,omitempty"`
				Active   bool   `json:"active"`
			}
			out := make([]profileJSON, len(profiles))
			for i, p := range profiles {
				out[i] = profileJSON{Nickname: p.Nickname, UserName: p.UserName, Email: p.Email, GHUser: p.GHUser, Active: p.Active}
			}
			b, err := json.MarshalIndent(out, "", "  ")
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", b)
			return nil
		}
		if len(profiles) == 0 {
			fmt.Println("No profiles")
			return nil
		}
		for _, p := range profiles {
			prefix := " "
			if p.Active {
				prefix = styleOK("✓")
			}
			fmt.Printf("%s  %-14s  %s <%s>\n", prefix, p.Nickname, p.UserName, p.Email)
		}
		return nil
	},
}

var removeCmd = &cobra.Command{
	Use:               "remove <nickname>",
	Aliases:           []string{"rm", "delete"},
	Short:             "Remove a profile",
	Args:              exactArgs(1),
	ValidArgsFunction: completeNicknames,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := store.Remove(args[0]); err != nil {
			return err
		}
		fmt.Printf("%s Profile '%s' removed\n", styleOK("✓"), args[0])
		return nil
	},
}

var currentCmd = &cobra.Command{
	Use:     "current",
	Aliases: []string{"whoami"},
	Short:   "Show the identity in use here",
	Long: `Shows the identity git will commit with in the current directory, and where it
comes from: this repo's pin, this terminal's session, or the global profile.`,
	Args: noArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		short, _ := cmd.Flags().GetBool("short")
		cwd, _ := os.Getwd()
		p, scope, err := effectiveProfile(cwd)
		if err != nil {
			return err
		}
		jsonOut, _ := cmd.Flags().GetBool("json")
		if jsonOut {
			scopeStr := "global"
			switch scope {
			case git.ScopeRepo:
				scopeStr = "repo"
			case git.ScopeSession:
				scopeStr = "session"
			}
			out := struct {
				Nickname               string `json:"nickname,omitempty"`
				UserName               string `json:"user_name,omitempty"`
				Email                  string `json:"email,omitempty"`
				Scope                  string `json:"scope"`
				GHUser                 string `json:"gh_user,omitempty"`
				SSHKey                 string `json:"ssh_key,omitempty"`
				CredentialHelperActive bool   `json:"credential_helper_active"`
				GuardActive            bool   `json:"guard_active"`
			}{
				Scope:                  scopeStr,
				CredentialHelperActive: git.IsCredentialHelperInstalled(),
				GuardActive:            git.IsGuardHookInstalled(),
			}
			if p != nil {
				out.Nickname, out.UserName, out.Email, out.GHUser, out.SSHKey = p.Nickname, p.UserName, p.Email, p.GHUser, p.SSHKey
			}
			b, err := json.MarshalIndent(out, "", "  ")
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", b)
			return nil
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
		nick := styleOK(p.Nickname)
		pinned := history.GetPinned(history.GetRepoKeyForPath(cwd))
		switch scope {
		case git.ScopeRepo:
			if pinned != "" && pinned != p.Nickname {
				fmt.Printf("%s — %s <%s>  (overriding repo pin '%s')\n", nick, p.UserName, p.Email, pinned)
			} else {
				fmt.Printf("%s — %s <%s>  (pinned to this repo)\n", nick, p.UserName, p.Email)
			}
		case git.ScopeSession:
			fmt.Printf("%s — %s <%s>  (this terminal's session)\n", nick, p.UserName, p.Email)
		default:
			fmt.Printf("%s — %s <%s>\n", nick, p.UserName, p.Email)
		}
		if scope != git.ScopeRepo && pinned != "" {
			if !shell.IsGHWrapperInstalled(shell.RCFile(shell.DetectShell())) {
				fmt.Println(styleDim(fmt.Sprintf("  (pinned to '%s' — inactive, Session Isolation is off)", pinned)))
			} else if p.Nickname != pinned {
				fmt.Println(styleDim(fmt.Sprintf("  (pinned to '%s' — overridden by switch)", pinned)))
			}
		}
		if git.IsCredentialHelperInstalled() {
			fmt.Println("HTTPS credential helper: active")
		}
		if git.IsGuardHookInstalled() {
			fmt.Println("Pre-commit guard: active")
		}
		return nil
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Import your current git config as a profile",
	Args:  noArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := store.ImportCurrent(); err != nil {
			return err
		}
		fmt.Printf("%s Imported current git config as 'default' profile\n", styleOK("✓"))
		return nil
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version, channel and whether an update exists",
	Long: `Shows the installed version, the release channel, how gitswitch is managed
(Homebrew or the install script), and whether a newer release exists.

For just the version number, with no network access, use 'gitswitch --version'.`,
	Args: noArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonOut, _ := cmd.Flags().GetBool("json")
		inst := ver.DetectInstall()
		channel := "stable"
		if ver.IsBeta(version) || inst.Formula == ver.BetaFormula {
			channel = "beta"
		}
		latest := ver.CachedLatestVersion(store.ConfigDir(), version)
		updateAvailable := latest != "" && ver.IsUpdateAvailable(version, latest)

		if jsonOut {
			out := struct {
				Version         string `json:"version"`
				Channel         string `json:"channel"`
				Install         string `json:"install"`
				Formula         string `json:"formula,omitempty"`
				Latest          string `json:"latest,omitempty"`
				UpdateAvailable bool   `json:"update_available"`
				UpgradeCommand  string `json:"upgrade_command"`
			}{version, channel, "script", inst.Formula, latest, updateAvailable, inst.UpgradeHint()}
			if inst.IsBrew() {
				out.Install = "homebrew"
			}
			b, err := json.MarshalIndent(out, "", "  ")
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", b)
			return nil
		}

		fmt.Printf("gitswitch %s\n", version)
		fmt.Println(styleDim(fmt.Sprintf("%s channel · managed by %s", channel, inst.Label())))
		switch {
		case updateAvailable:
			fmt.Printf("Update available: %s — run: %s\n", styleOK(latest), inst.UpgradeHint())
		case latest != "":
			fmt.Println("You're on the latest version.")
		}
		return nil
	},
}

var pacmanCmd = &cobra.Command{
	Use:   "pacman",
	Short: "Toggle arcade mode for the profile picker",
	Long:  "Toggles arcade mode — a retro skin for the profile picker — and opens it. Run it again to turn it back off.",
	Args:  noArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureInitialized(); err != nil {
			return err
		}
		prefs, err := store.LoadPrefs()
		if err != nil {
			prefs = storage.Prefs{}
		}
		prefs.ArcadeMode = !prefs.ArcadeMode
		if err := store.SavePrefs(prefs); err != nil {
			return err
		}
		// tui.New reads prefs.ArcadeMode straight off disk, including the intro
		// trigger, so no What's New body is needed here.
		m, err := tui.New(store, version, "")
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

var reauthorCmd = &cobra.Command{
	Use:   "reauthor <base>",
	Short: "Fix the author on commits already made",
	Long: `Rewrites the author and committer of every commit between <base> and HEAD
to a stored profile's identity. <base> is a commit-ish (e.g. HEAD~3, a SHA) or
a bare number N meaning "the last N commits".

Use --from to only touch commits currently authored by a given email
(the pre-switch account), leaving other commits in range untouched.

This rewrites history. If the branch is already pushed, pass --push to
force-push (--force-with-lease) afterward, or do it yourself once you've
checked the result.`,
	Example: `  gitswitch reauthor 3 --to work                    the last 3 commits
  gitswitch reauthor HEAD~5 --to work --from me@old.example
  gitswitch reauthor 3 --to work --push --yes       rewrite and force-push, no prompts`,
	Args: exactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		to, _ := cmd.Flags().GetString("to")
		from, _ := cmd.Flags().GetString("from")
		push, _ := cmd.Flags().GetBool("push")
		yes, _ := cmd.Flags().GetBool("yes")
		if to == "" {
			return usageErrorf(cmd, "--to <nickname> is required")
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
		if !yes && !confirmPrompt("Proceed? [y/N] ") {
			fmt.Println("Aborted.")
			return nil
		}

		if err := git.Reauthor(base, from, p.UserName, p.Email); err != nil {
			return err
		}
		fmt.Println("✓ History rewritten.")

		if push {
			if !yes && !confirmPrompt("Force-push (--force-with-lease) now? [y/N] ") {
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
// promptInstallTool shows the exact command needed to install a missing tool,
// asks for confirmation (skipped when yes is true), and runs it. Returns nil
// with no action taken if the tool has no InstallCommand for this platform
// (e.g. macOS git, which needs the Xcode CLT GUI dialog) or if the user declines.
func promptInstallTool(label string, t *prereqs.ToolCheck, yes bool) error {
	if len(t.InstallCommand) == 0 {
		return nil
	}
	fmt.Printf("  About to run: %s\n", strings.Join(t.InstallCommand, " "))
	if !yes && !confirmPrompt("  Proceed? [y/N] ") {
		fmt.Println("  Skipped.")
		return nil
	}
	if err := runCommand(t.InstallCommand[0], t.InstallCommand[1:]...); err != nil {
		return fmt.Errorf("%s install failed: %w", label, err)
	}
	fmt.Printf("  ✓ %s installed.\n", label)
	return nil
}

// A pin is a local identity: it writes the profile into this repo's
// .git/config, so commits here are correct no matter which profile is active
// globally. Nothing needs to be switched on entry.
var pinCmd = &cobra.Command{
	Use:   "pin [nickname]",
	Short: "Always use a profile in this repo",
	Long: "Pin an identity to this repo — writes it to the repo's local git config and\n" +
		"switches the gh CLI to the matching account.\n\n" +
		"With no nickname, adopts the identity the repo already has in its local git\n" +
		"config (useful for repos configured by hand before gitswitch).\n\n" +
		"Requires Session Isolation — turned on automatically if it's off.",
	Example: "  gitswitch pin work    always commit as \"work\" in this repo\n" +
		"  gitswitch pin         adopt the identity this repo already has\n" +
		"  gitswitch unpin       go back to the global identity",
	Args:              argRange(0, 1),
	ValidArgsFunction: completeNicknames,
	RunE: func(cmd *cobra.Command, args []string) error {
		repoKey := history.GetRepoKey()
		if repoKey == "" {
			return fmt.Errorf("not inside a git repo")
		}
		if rcFile := shell.RCFile(shell.DetectShell()); !shell.IsGHWrapperInstalled(rcFile) {
			if _, err := shell.InstallGHWrapper(shell.DetectShell()); err != nil {
				return err
			}
			fmt.Println("Session Isolation was off — turning it on (pins need it)")
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
				fmt.Printf("%s Replacing this repo's existing local identity <%s>\n", styleWarn("⚠"), existing)
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

		warn, err := applyProfile(git.New(false), p)
		if err != nil {
			return err
		}
		if err := history.Pin(repoKey, p.Nickname); err != nil {
			return err
		}
		if warn != "" {
			fmt.Printf("warning: %s\n", warn)
		}
		fmt.Printf("%s Pinned '%s' to this repo — %s <%s> (local git config; global identity unchanged)\n",
			styleOK("✓"), p.Nickname, p.UserName, p.Email)
		xcodeAutoBind(cwd, p)
		return nil
	},
}

var unpinCmd = &cobra.Command{
	Use:   "unpin",
	Short: "Remove this repo's pin and use the global identity again",
	Args:  noArgs,
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
		cwd, _ := os.Getwd()
		xcodeAutoUnbind(cwd)
		fmt.Printf("%s Unpinned — this repo now uses the global identity\n", styleOK("✓"))
		return nil
	},
}

var recordCmd = &cobra.Command{
	Use:    "record",
	Short:  "Record current identity for this repo (called by shell hooks on repo entry)",
	Hidden: true, // plumbing for the shell hook, not something to type
	Args:   noArgs,
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
		// A repo or session carrying its own identity (or an explicit pin) is
		// already decided: don't record the global profile against it (that would
		// teach the wrong one or overwrite an existing pin), but do point gh at
		// the matching account so pushes from here use the right user.
		scope, email := git.ResolveIdentity(path)
		if scope == git.ScopeRepo || scope == git.ScopeSession || history.GetPinned(repoKey) != "" || git.LocalEmail(path) != "" {
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

		// Record usage first — this is also the single locked read that gives us
		// the updated count, so we never read history.json twice per call.
		count, err := history.Record(repoKey, active.Nickname)
		if err != nil {
			return err
		}

		// Once the same account has been used in this repo enough times to call
		// it a pattern (not just one drive-by commit), auto-pin it so coming
		// back always uses the same account. Fires exactly once, the visit the
		// threshold is crossed. Silent — hooks must stay quiet either way.
		// Toggle: Utilities tab, or storage.Prefs.AutoPinDisabled.
		if count == history.AutoPinThreshold {
			if prefs, err := store.LoadPrefs(); err == nil && !prefs.AutoPinDisabled {
				if _, err := applyProfile(git.New(false), active); err == nil {
					_ = history.Pin(repoKey, active.Nickname)
				}
			}
		}

		return nil
	},
}

// errNoRecommendation signals a silent exit 1 from recommendCmd.
// SilenceErrors on the command prevents cobra from printing it.
var errNoRecommendation = fmt.Errorf("")

var recommendCmd = &cobra.Command{
	Use:           "recommend",
	Short:         "Print recommended profile for current repo (used by shell hooks)",
	Hidden:        true, // plumbing for the shell hook, not something to type
	Args:          noArgs,
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

		scope, email := git.ResolveIdentity(path)
		if scope == git.ScopeSession {
			return errNoRecommendation
		}

		// If the repo is pinned in history.json, check whether its effective
		// identity currently matches the pin. When the user overrode the pin
		// via `gitswitch switch <other>`, the next `cd` into the repo nudges
		// back to the pinned profile; otherwise it is already right.
		if pinned := history.GetPinned(repoKey); pinned != "" {
			if current := store.GetByEmail(email); current != nil && current.Nickname == pinned {
				return errNoRecommendation
			}
			p, err := store.Get(pinned)
			if err != nil {
				return errNoRecommendation
			}
			fmt.Printf("%s\t%s\t%s\n", p.Nickname, p.UserName, p.Email)
			return nil
		}

		// A hand-configured repo-local identity (with no history pin) already
		// commits as that identity — never nudge it to switch globally.
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

var skillsCmd = &cobra.Command{
	Use:   "skills",
	Short: "Teach your AI coding agent to use gitswitch",
	Long: `Installs or updates the gitswitch skill so an AI coding agent knows how to switch
identities, enforce commit guards, sync Xcode accounts, and fix wrong-account commits.

Prefers skills.sh (npx skills add), which installs into whichever
Agent-Skills-compatible tool you have locally — Claude Code, Cursor, Codex,
and more — without gitswitch needing to know each one's own layout. Falls
back to installing directly for Claude Code and the shared .agents/skills
convention (used by Codex and others) when npx isn't available or --offline
is passed.

Use --check to see if installed skills are outdated, or --update to bring all
installed skill files up to date with the current binary's embedded skill.`,
	Example: `  gitswitch skills                  install for your agent
  gitswitch skills --scope project  install into this repo only
  gitswitch skills --check          are installed copies current?
  gitswitch skills --update         refresh installed copies`,
	Args: noArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		checkOnly, _ := cmd.Flags().GetBool("check")
		updateOnly, _ := cmd.Flags().GetBool("update")
		scope, _ := cmd.Flags().GetString("scope")
		offline, _ := cmd.Flags().GetBool("offline")
		if err := oneOf(cmd, "scope", scope, "user", "project"); err != nil {
			return err
		}

		if checkOnly {
			outdated, upToDate := skill.IsOutdated()
			if len(outdated) == 0 && len(upToDate) == 0 {
				fmt.Println("No gitswitch skill installed yet. Run 'gitswitch skills' to install.")
				return nil
			}
			for _, p := range upToDate {
				fmt.Printf("%s %s (up to date)\n", styleOK("✓"), p)
			}
			if len(outdated) > 0 {
				for _, p := range outdated {
					fmt.Printf("%s %s (outdated)\n", styleWarn("!"), p)
				}
				fmt.Println("\nRun 'gitswitch skills --update' to bring them up to date.")
			}
			return nil
		}

		if updateOnly {
			updated, err := skill.UpdateAll()
			if err != nil {
				return fmt.Errorf("could not update skill: %w", err)
			}
			if len(updated) == 0 {
				outdated, upToDate := skill.IsOutdated()
				if len(upToDate) > 0 && len(outdated) == 0 {
					fmt.Printf("%s All installed gitswitch skills are already up to date.\n", styleOK("✓"))
				} else {
					fmt.Println("No installed gitswitch skill found. Run 'gitswitch skills' to install.")
				}
				return nil
			}
			for _, p := range updated {
				fmt.Printf("%s Updated %s\n", styleOK("✓"), p)
			}
			fmt.Println("  Reload your agent (or open a new session) to activate.")
			return nil
		}

		if !offline {
			if _, err := exec.LookPath("npx"); err == nil {
				fmt.Println("  Installing via skills.sh (works with Claude Code, Cursor, Codex, and more)...")
				c := exec.Command("npx", "-y", "skills", "add", "aksisonline/gitswitch")
				c.Stdin = os.Stdin
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				if err := c.Run(); err == nil {
					return nil
				}
				fmt.Println("  ⚠  skills.sh install failed (or no network) — falling back to the offline installer")
			}
		}

		var claudeBase, agentsBase string
		if scope == "project" {
			claudeBase, agentsBase = ".claude", ".agents"
		} else {
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			claudeBase = filepath.Join(home, ".claude")
			agentsBase = filepath.Join(home, ".agents")
		}

		var written []string
		for _, base := range []string{claudeBase, agentsBase} {
			dest := filepath.Join(base, "skills", "gitswitch")
			if err := os.MkdirAll(dest, 0755); err != nil {
				return fmt.Errorf("could not create skills directory: %w", err)
			}
			if err := os.WriteFile(filepath.Join(dest, "SKILL.md"), skill.Content, 0644); err != nil {
				return fmt.Errorf("could not write skill: %w", err)
			}
			written = append(written, dest)
		}

		for _, dest := range written {
			fmt.Printf("✓ Skill installed to %s\n", dest)
		}
		fmt.Println("  Reload your agent (or open a new session) to activate.")
		return nil
	},
}

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Set up shell integration, HTTPS routing and the commit guard",
	Long: `Interactive setup wizard. Detects your shell, shows what each step does,
and asks before making any changes. Use --yes to accept all defaults without
prompts (for scripts and CI).`,
	Example: `  gitswitch shell                  interactive setup
  gitswitch shell --yes            accept every default, no prompts
  gitswitch shell --shell fish     target a specific shell
  gitswitch shell --yes --guard=false --https=false   shell integration only`,
	Args: noArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		shellFlag, _ := cmd.Flags().GetString("shell")
		yes, _ := cmd.Flags().GetBool("yes")
		httpsDefault, _ := cmd.Flags().GetBool("https")
		guardDefault, _ := cmd.Flags().GetBool("guard")
		if err := oneOf(cmd, "shell", shellFlag, "", "zsh", "bash", "fish"); err != nil {
			return err
		}

		opts, err := wizard.Run(wizard.Config{
			ShellOverride: shellFlag,
			Yes:           yes,
			HTTPSDefault:  httpsDefault,
			GuardDefault:  guardDefault,
		})
		if err != nil {
			return fmt.Errorf("setup interrupted: %w", err)
		}

		var shellResult string
		if opts.InstallShell {
			// Reinstall, not Install: Install is a no-op when the marker block is
			// already there, so a user told "shell integration updated — run gitswitch
			// shell" would keep their old hook forever. Reinstall replaces it.
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

		var ghWrapperErr error
		if opts.InstallGHWrapper && !shell.IsGHWrapperInstalled(shell.RCFile(opts.Shell)) {
			_, ghWrapperErr = shell.InstallGHWrapper(opts.Shell)
		}

		var guardErr error
		if opts.InstallGuard {
			guardErr = git.InstallGuardHook()
		}

		wizard.PrintSummary(
			os.Stdout,
			shellResult,
			opts.InstallShell,
			opts.InstallHTTPS && httpsErr == nil,
			opts.InstallGHWrapper && ghWrapperErr == nil,
			httpsErr,
			opts.InstallGuard && guardErr == nil,
		)
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

// ghUserCmd backs the `gh` shell wrapper installed by the gh-wrapper toggle
// (see internal/shell InstallGHWrapper): it prints the gh login that applies
// to the calling directory, sharing resolution with the credential helper so
// bare `gh` commands and HTTPS pushes always agree on identity. Prints
// nothing (exit 0) when no profile applies, so the wrapper falls back to gh's
// own active account.
var ghUserCmd = &cobra.Command{
	Use:    "ghuser",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, _ := os.Getwd()
		if login := credential.ResolveGHUser(store, history.GetRepoKey(), cwd); login != "" {
			fmt.Println(login)
		}
		return nil
	},
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Connect a GitHub account",
	Long: `Connects a GitHub account and saves it as a profile. Sign-in is handed to the
GitHub CLI ('gh auth login'), so you get whatever method it offers — browser,
one-time code, or a token. gitswitch stores no token of its own.`,
	Example: `  gitswitch login
  gitswitch login --profile work                      choose the nickname
  gitswitch login --host github.mycompany.com         GitHub Enterprise`,
	Args: noArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		host, _ := cmd.Flags().GetString("host")
		profileName, _ := cmd.Flags().GetString("profile")
		return completeLogin(host, profileName)
	},
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove everything 'gitswitch shell' set up",
	Long: `Removes what 'gitswitch shell' added: the shell integration block, the HTTPS
credential helper, Session Isolation, and the pre-commit guard.

Your profiles stay, and so does the gitswitch program itself — remove that the
way you installed it (for Homebrew: brew uninstall gitswitch).`,
	Args: noArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		shellFlag, _ := cmd.Flags().GetString("shell")
		if err := oneOf(cmd, "shell", shellFlag, "", "zsh", "bash", "fish"); err != nil {
			return err
		}

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

		result, err := shell.Uninstall(sh)
		if err != nil {
			return fmt.Errorf("uninstall failed: %w", err)
		}
		ok, warn := styleOK("✓"), styleWarn("warning:")
		fmt.Printf("%s %s\n", ok, result)

		if git.IsCredentialHelperInstalled() {
			if err := git.UninstallCredentialHelper(); err != nil {
				fmt.Printf("  %s could not remove HTTPS credential helper: %v\n", warn, err)
			} else {
				fmt.Printf("%s HTTPS credential helper removed\n", ok)
			}
		}

		if shell.IsGHWrapperInstalled(shell.RCFile(sh)) {
			if _, err := shell.UninstallGHWrapper(sh); err != nil {
				fmt.Printf("  %s could not remove gh wrapper: %v\n", warn, err)
			} else {
				fmt.Printf("%s gh wrapper removed\n", ok)
			}
		}

		if git.IsGuardHookInstalled() {
			if err := git.UninstallGuardHook(); err != nil {
				fmt.Printf("  %s could not remove pre-commit guard: %v\n", warn, err)
			} else {
				fmt.Printf("%s Pre-commit guard removed\n", ok)
			}
		}

		fmt.Println("  Reload your shell (or open a new terminal) to complete removal.")
		if inst := ver.DetectInstall(); inst.IsBrew() {
			fmt.Println(styleDim("  To remove gitswitch itself: brew uninstall " + inst.Formula))
		}
		return nil
	},
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check your setup and offer to fix what's missing",
	Long: `Checks that git and the GitHub CLI are installed, and that HTTPS routing,
Session Isolation and the pre-commit guard are active. Offers to install git or
gh when they are missing.`,
	Example: `  gitswitch doctor
  gitswitch doctor --json    machine-readable report (for scripts and agents)
  gitswitch doctor --yes     install missing tools without asking`,
	Args: noArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonOut, _ := cmd.Flags().GetBool("json")
		yes, _ := cmd.Flags().GetBool("yes")
		r := prereqs.Check()
		conflicts := git.CredentialHelperConflicts()
		ghWrapperInstalled := shell.IsGHWrapperInstalled(shell.RCFile(shell.DetectShell()))
		guardInstalled := git.IsGuardHookInstalled()
		inst := ver.DetectInstall()
		copies := installsOnPath()
		if jsonOut {
			// Embedded so git/gh stay top-level, matching the pre-existing shape.
			// Version/Profiles absorb what `gitswitch setup --agent` used to report,
			// now that doctor --json is the one machine-readable manifest.
			report := struct {
				prereqs.CheckResult
				Version  string `json:"version"`
				Profiles int    `json:"profiles"`
				HTTPS    struct {
					RoutedByGitswitch bool                 `json:"routed_by_gitswitch"`
					Conflicts         []git.HelperConflict `json:"conflicts,omitempty"`
				} `json:"https"`
				GHWrapperInstalled bool `json:"gh_wrapper_installed"`
				GuardInstalled     bool `json:"guard_installed"`
				// How this binary is managed ("homebrew" or "script") and every
				// distinct gitswitch on PATH — more than one is a problem.
				Install       string   `json:"install"`
				InstallsFound []string `json:"installs_on_path"`
			}{CheckResult: r}
			report.Install = "script"
			if inst.IsBrew() {
				report.Install = "homebrew"
			}
			report.InstallsFound = copies
			report.Version = version
			profiles, _ := store.Load()
			report.Profiles = len(profiles)
			report.HTTPS.RoutedByGitswitch = len(conflicts) == 0
			report.HTTPS.Conflicts = conflicts
			report.GHWrapperInstalled = ghWrapperInstalled
			report.GuardInstalled = guardInstalled
			b, err := json.MarshalIndent(report, "", "  ")
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", b)
			return nil
		}
		ok, warn, bad := styleOK("✓"), styleWarn("⚠"), styleErr("✗")

		fmt.Println()
		if r.Git.Installed {
			fmt.Printf("  %s  git %s\n", ok, r.Git.Version)
		} else {
			fmt.Printf("  %s  git  not found\n", bad)
		}
		if r.GH.Installed {
			fmt.Printf("  %s  gh  %s\n", ok, r.GH.Version)
		} else {
			fmt.Printf("  %s  gh   not found (required for gitswitch login)\n", warn)
		}

		// HTTPS routing: being registered is not enough — another tool's per-host
		// helper can be asked first, and then gitswitch never sees the request.
		switch {
		case len(conflicts) == 0:
			fmt.Printf("  %s  HTTPS pushes routed by gitswitch\n", ok)
		case conflicts[0].Winner == "" && len(conflicts) == 1:
			fmt.Printf("  %s  HTTPS pushes not routed by gitswitch — run: gitswitch shell\n", warn)
		default:
			fmt.Printf("  %s  HTTPS pushes answered by another helper before gitswitch:\n", bad)
			for _, c := range conflicts {
				if c.Winner == "" {
					continue
				}
				fmt.Printf("       %s → %s\n", c.Key, c.Winner)
			}
			fmt.Println("       pushes may use the wrong account — run: gitswitch shell")
		}

		if ghWrapperInstalled {
			fmt.Printf("  %s  Session Isolation active (bare `gh` commands resolve per-repo)\n", ok)
		} else {
			fmt.Printf("  %s  Session Isolation off — bare `gh` commands use gh's single global account, and `gitswitch pin` has no effect\n", warn)
		}

		if guardInstalled {
			fmt.Printf("  %s  Pre-commit guard active (blocks wrong-author commits)\n", ok)
		} else {
			fmt.Printf("  %s  Pre-commit guard off — wrong-author commits are not blocked (run: gitswitch guard enable)\n", warn)
		}

		// Two installs shadowing each other is the classic "I upgraded and
		// nothing changed": the stale copy earlier on PATH keeps winning.
		if len(copies) > 1 {
			fmt.Printf("  %s  %d copies of gitswitch on PATH — the first one is the one that runs:\n", warn, len(copies))
			for _, c := range copies {
				fmt.Printf("       %s\n", c)
			}
			if inst.IsBrew() || strings.Contains(strings.Join(copies, " "), "/Cellar/") {
				fmt.Println("       keep the Homebrew one and delete the others, then: brew upgrade gitswitch")
			} else {
				fmt.Println("       keep one and delete the others")
			}
		} else {
			fmt.Printf("  %s  gitswitch %s, updated with %s (%s)\n", ok, version, inst.Label(), inst.UpgradeHint())
		}

		fmt.Println()
		prereqs.PrintWarnings(r)

		if !r.Git.Installed {
			if err := promptInstallTool("git", r.Git, yes); err != nil {
				return err
			}
		}
		if !r.GH.Installed {
			if err := promptInstallTool("gh", r.GH, yes); err != nil {
				return err
			}
		}
		return nil
	},
}

func main() {
	tui.AfterSwitch = func(p storage.Profile) { xcodeAutoSync(&p) }
	// Help lists commands in the order they are added here — the order a new
	// user needs them — rather than alphabetically.
	cobra.EnableCommandSorting = false
	rootCmd.SuggestionsMinimumDistance = 2
	rootCmd.AddCommand(
		loginCmd, addCmd, switchCmd, listCmd, currentCmd, removeCmd, initCmd,
		shellCmd, doctorCmd, pinCmd, unpinCmd, guardCmd, xcodeCmd, uninstallCmd,
		versionCmd, upgradeCmd, betaCmd, stableCmd,
		reauthorCmd, skillsCmd, pacmanCmd,
		recordCmd, recommendCmd, hookCheckCmd, credentialCmd, ghUserCmd, snakeCmd,
	)

	rootCmd.AddGroup(
		&cobra.Group{ID: "identity", Title: "Identity"},
		&cobra.Group{ID: "isolation", Title: "Setup & Session Isolation"},
		&cobra.Group{ID: "channels", Title: "Updates & Channels"},
		&cobra.Group{ID: "extras", Title: "Extras"},
	)
	for _, c := range []*cobra.Command{loginCmd, addCmd, switchCmd, listCmd, currentCmd, removeCmd, initCmd} {
		c.GroupID = "identity"
	}
	for _, c := range []*cobra.Command{shellCmd, doctorCmd, pinCmd, unpinCmd, guardCmd, xcodeCmd, uninstallCmd} {
		c.GroupID = "isolation"
	}
	for _, c := range []*cobra.Command{versionCmd, upgradeCmd, betaCmd, stableCmd} {
		c.GroupID = "channels"
	}
	for _, c := range []*cobra.Command{reauthorCmd, skillsCmd, pacmanCmd} {
		c.GroupID = "extras"
	}
	rootCmd.SetHelpCommandGroupID("extras")
	rootCmd.SetCompletionCommandGroupID("extras")

	// --version / -v (and --v, -V via normalizeArgs): offline and instant.
	// `gitswitch version` is the long form that also checks for updates.
	rootCmd.Version = version
	rootCmd.SetVersionTemplate("gitswitch {{.Version}}\n")
	rootCmd.Flags().BoolP("version", "v", false, "Show the installed version")
	installHelp(rootCmd)

	// Flags follow one vocabulary across commands:
	//   --json       machine-readable output (every command that reports state)
	//   -y, --yes    skip confirmation prompts (every command that asks)
	//   --dry-run    show what would happen, change nothing
	//   -q, --quiet  no output on success
	addCmd.Flags().String("sign-key", "", "Signing key: GPG key ID, or SSH key path for gpg.format=ssh")
	addCmd.Flags().String("ssh-key", "", "SSH private key path, e.g. ~/.ssh/id_work (sets core.sshCommand)")
	addCmd.Flags().String("gh-user", "", "GitHub username, so the gh CLI switches account with the profile")
	switchCmd.Flags().BoolP("quiet", "q", false, "Print nothing on success")
	currentCmd.Flags().Bool("short", false, "Output nickname and email tab-separated (for Starship and scripts)")
	currentCmd.Flags().Bool("prompt", false, "Output nickname and theme color tab-separated (for shell prompt functions)")
	currentCmd.Flags().Bool("json", false, "Output machine-readable JSON")
	listCmd.Flags().Bool("json", false, "Output machine-readable JSON")
	versionCmd.Flags().Bool("json", false, "Output machine-readable JSON")
	recordCmd.Flags().String("path", "", "Directory to record for (default: current working directory)")
	recommendCmd.Flags().String("path", "", "Directory to check (default: current working directory)")
	shellCmd.Flags().String("shell", "", "Shell to set up: zsh, bash or fish (default: auto-detect; skips the interactive wizard)")
	shellCmd.Flags().Bool("https", true, "Register the HTTPS credential helper")
	shellCmd.Flags().Bool("guard", true, "Enable the pre-commit identity guard")
	shellCmd.Flags().BoolP("yes", "y", false, "Accept all defaults without prompts (for scripts and CI)")
	uninstallCmd.Flags().String("shell", "", "Shell to clean up: zsh, bash or fish (default: auto-detect)")
	skillsCmd.Flags().String("scope", "user", "Where to install: user (~/.claude, ~/.agents) or project (this repo)")
	skillsCmd.Flags().Bool("offline", false, "Skip skills.sh (npx) and use the built-in installer")
	skillsCmd.Flags().Bool("check", false, "Report whether installed skills are up to date")
	skillsCmd.Flags().Bool("update", false, "Update installed skills to match this version")
	doctorCmd.Flags().Bool("json", false, "Output machine-readable JSON")
	doctorCmd.Flags().BoolP("yes", "y", false, "Install missing tools without asking")
	loginCmd.Flags().String("host", "", "GitHub host, for GitHub Enterprise (default: github.com)")
	loginCmd.Flags().String("profile", "", "Nickname for the new profile (default: your GitHub username)")
	reauthorCmd.Flags().String("to", "", "Profile to attribute the commits to (required)")
	reauthorCmd.Flags().String("from", "", "Only rewrite commits currently authored by this email")
	reauthorCmd.Flags().Bool("push", false, "Force-push (--force-with-lease) after rewriting")
	reauthorCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompts (for scripts and agents)")
	upgradeCmd.Flags().Bool("check", false, "Only report whether an update exists")
	upgradeCmd.Flags().Bool("dry-run", false, "Show the command that would run, change nothing")
	for _, c := range []*cobra.Command{betaCmd, stableCmd} {
		c.Flags().BoolP("yes", "y", false, "Switch without asking for confirmation")
		c.Flags().Bool("dry-run", false, "Show the commands that would run, change nothing")
	}

	// Shell completion for flag values.
	shells := fixedCompletions("zsh", "bash", "fish")
	_ = shellCmd.RegisterFlagCompletionFunc("shell", shells)
	_ = uninstallCmd.RegisterFlagCompletionFunc("shell", shells)
	_ = skillsCmd.RegisterFlagCompletionFunc("scope", fixedCompletions("user", "project"))
	_ = reauthorCmd.RegisterFlagCompletionFunc("to", func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		return nicknameCompletions(), cobra.ShellCompDirectiveNoFileComp
	})

	normalizeArgs(os.Args)
	if code := reportError(rootCmd.Execute()); code != 0 {
		os.Exit(code)
	}
}
