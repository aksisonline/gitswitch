package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aksisonline/gitswitch/internal/storage"
	"github.com/aksisonline/gitswitch/internal/xcode"
	"github.com/spf13/cobra"
)

const (
	xcodeSettingsPath = "Xcode → Settings → Accounts"
	xcodePATURL       = "https://github.com/settings/tokens/new?scopes=repo,write:public_key,read:org,user&description=Xcode"
)

// xcodeMarker is the opt-in flag file: when present, `switch` and `pin` keep
// Xcode's account choice in step automatically.
func xcodeMarker() string { return filepath.Join(store.ConfigDir(), "xcode-sync") }

func xcodeAutoEnabled() bool {
	_, err := os.Stat(xcodeMarker())
	return err == nil
}

// xcodeAutoSync promotes p's Xcode GitHub account after a global switch when
// auto-sync is on. Best effort and silent on failure: a switch must never fail
// (or print noise from shell hooks) because of Xcode.
func xcodeAutoSync(p *storage.Profile) string {
	if p == nil || p.GHUser == "" || !xcodeAutoEnabled() || !xcode.Supported() {
		return ""
	}
	res, err := syncXcodeAccount(xcode.NewStore(), p.GHUser, false)
	if err != nil || !res.Changed {
		return ""
	}
	if xcode.IsRunning() {
		return fmt.Sprintf("Xcode now prefers %s (restart Xcode to apply)", p.GHUser)
	}
	return fmt.Sprintf("Xcode now prefers %s", p.GHUser)
}

// xcodeAutoBind pins a repo's https GitHub remotes to p's account for Xcode
// when auto-sync is on. Best effort, silent on failure.
func xcodeAutoBind(dir string, p *storage.Profile) {
	if p == nil || p.GHUser == "" || !xcodeAutoEnabled() || !xcode.Supported() {
		return
	}
	_, _ = bindRemotes(dir, p.GHUser)
}

// xcodeAutoUnbind undoes xcodeAutoBind on unpin.
func xcodeAutoUnbind(dir string) {
	if !xcodeAutoEnabled() {
		return
	}
	_, _ = unbindRemotes(dir)
}

type syncResult struct {
	User    string   `json:"user"`
	Found   bool     `json:"found"`
	Changed bool     `json:"changed"`
	Order   []string `json:"order"`
}

// syncXcodeAccount moves user's GitHub account to the front in Xcode. With
// dryRun it only reports what would change.
func syncXcodeAccount(s *xcode.Store, user string, dryRun bool) (syncResult, error) {
	res := syncResult{User: user}
	accts, err := s.Accounts()
	if err != nil {
		return res, err
	}
	res.Found = xcode.FindGitHub(accts, user) >= 0
	reordered, changed := xcode.Promote(accts, user)
	res.Changed = changed
	for _, a := range reordered {
		res.Order = append(res.Order, a.Username)
	}
	if changed && !dryRun {
		if err := s.Save(reordered); err != nil {
			return res, err
		}
	}
	return res, nil
}

func remoteURLs(dir string) (map[string]string, error) {
	out, err := exec.Command("git", "-C", dir, "remote").Output()
	if err != nil {
		return nil, fmt.Errorf("not inside a git repo")
	}
	urls := map[string]string{}
	for _, name := range strings.Fields(string(out)) {
		u, err := exec.Command("git", "-C", dir, "config", "--get", "remote."+name+".url").Output()
		if err == nil {
			urls[name] = strings.TrimSpace(string(u))
		}
	}
	return urls, nil
}

func rewriteRemotes(dir string, fn func(string) (string, bool)) ([]string, error) {
	urls, err := remoteURLs(dir)
	if err != nil {
		return nil, err
	}
	var changed []string
	for name, u := range urls {
		if nu, ok := fn(u); ok {
			if out, err := exec.Command("git", "-C", dir, "remote", "set-url", name, nu).CombinedOutput(); err != nil {
				return changed, fmt.Errorf("remote %s: %s", name, strings.TrimSpace(string(out)))
			}
			changed = append(changed, name)
		}
	}
	return changed, nil
}

func bindRemotes(dir, user string) ([]string, error) {
	return rewriteRemotes(dir, func(u string) (string, bool) { return xcode.BindURL(u, user) })
}

func unbindRemotes(dir string) ([]string, error) { return rewriteRemotes(dir, xcode.UnbindURL) }

var xcodeCmd = &cobra.Command{
	Use:   "xcode",
	Short: "Keep Xcode's GitHub account in step with your gitswitch profile",
	Long: "Xcode keeps its own GitHub accounts (" + xcodeSettingsPath + ") and, when a repo's\n" +
		"owner doesn't match one of them, silently uses whichever account is listed first.\n" +
		"gitswitch can't write Xcode's tokens (they live in Apple's private keychain\n" +
		"group), but it can manage which account Xcode picks:\n\n" +
		"  xcode add     copy a profile's token so adding it to Xcode is one paste\n" +
		"  xcode sync    make the active profile's account Xcode's default\n" +
		"  xcode bind    pin this repo's remotes to its profile's account\n" +
		"  xcode enable  do sync/bind automatically on 'switch' and 'pin'\n\n" +
		"Commit authorship needs nothing extra — Xcode already reads git config.",
	RunE: func(cmd *cobra.Command, args []string) error { return runXcodeStatus(cmd) },
}

var xcodeStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Xcode's GitHub accounts and how they map to your profiles",
	RunE:  func(cmd *cobra.Command, args []string) error { return runXcodeStatus(cmd) },
}

func requireXcode() error {
	if !xcode.Supported() {
		return fmt.Errorf("Xcode integration is only available on macOS")
	}
	return nil
}

func runXcodeStatus(cmd *cobra.Command) error {
	if err := requireXcode(); err != nil {
		return err
	}
	accts, err := xcode.NewStore().Accounts()
	if err != nil {
		return err
	}
	profiles, err := store.Load()
	if err != nil {
		return err
	}
	jsonOut, _ := cmd.Flags().GetBool("json")

	type acctJSON struct {
		Username string `json:"username"`
		Host     string `json:"host"`
		Profile  string `json:"profile,omitempty"`
		Default  bool   `json:"default"`
	}
	type missingJSON struct {
		Profile string `json:"profile"`
		GHUser  string `json:"gh_user"`
	}
	var accounts []acctJSON
	seenGitHub := false
	for _, a := range accts {
		j := acctJSON{Username: a.Username, Host: a.HostURL}
		for _, p := range profiles {
			if p.GHUser != "" && a.IsGitHub() && strings.EqualFold(p.GHUser, a.Username) {
				j.Profile = p.Nickname
			}
		}
		if a.IsGitHub() && !seenGitHub {
			j.Default, seenGitHub = true, true
		}
		accounts = append(accounts, j)
	}
	var missing []missingJSON
	for _, p := range profiles {
		if p.GHUser != "" && xcode.FindGitHub(accts, p.GHUser) < 0 {
			missing = append(missing, missingJSON{p.Nickname, p.GHUser})
		}
	}
	active, _ := store.GetActive()
	activeNick := ""
	inSync := true
	if active != nil {
		activeNick = active.Nickname
		if active.GHUser != "" && xcode.FindGitHub(accts, active.GHUser) >= 0 {
			for _, a := range accounts {
				if a.Default && !strings.EqualFold(a.Username, active.GHUser) {
					inSync = false
				}
			}
		}
	}

	if jsonOut {
		b, err := json.MarshalIndent(map[string]any{
			"running":  xcode.IsRunning(),
			"auto":     xcodeAutoEnabled(),
			"active":   activeNick,
			"in_sync":  inSync,
			"accounts": accounts,
			"missing":  missing,
		}, "", "  ")
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", b)
		return nil
	}

	if len(accounts) == 0 {
		fmt.Printf("%s No GitHub accounts in Xcode yet (%s)\n", styleWarn("⚠"), xcodeSettingsPath)
	}
	for _, a := range accounts {
		mark, label := " ", ""
		if a.Default {
			mark, label = styleOK("★"), styleDim("  ← Xcode's default")
		}
		who := styleDim("(no gitswitch profile)")
		if a.Profile != "" {
			who = "→ " + a.Profile
		}
		fmt.Printf("  %s %-18s %s%s\n", mark, a.Username, who, label)
	}
	for _, m := range missing {
		fmt.Printf("  %s %-18s %s\n", styleWarn("+"), m.GHUser, styleDim("not in Xcode — run: gitswitch xcode add "+m.Profile))
	}
	fmt.Println()
	switch {
	case active != nil && !inSync:
		fmt.Printf("%s Xcode's default isn't your active profile '%s' — run: gitswitch xcode sync\n", styleWarn("⚠"), activeNick)
	case active != nil:
		fmt.Printf("%s Xcode's default matches your active profile\n", styleOK("✓"))
	}
	if xcodeAutoEnabled() {
		fmt.Printf("%s Auto-sync on (switch/pin keep Xcode in step)\n", styleOK("✓"))
	} else {
		fmt.Printf("  Auto-sync off — run 'gitswitch xcode enable' to follow switch/pin automatically\n")
	}
	if xcode.IsRunning() {
		fmt.Printf("%s Xcode is running — changes apply after you restart it\n", styleWarn("⚠"))
	}
	return nil
}

var xcodeSyncCmd = &cobra.Command{
	Use:   "sync [nickname]",
	Short: "Make a profile's GitHub account Xcode's default (active profile if omitted)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireXcode(); err != nil {
			return err
		}
		var p *storage.Profile
		var err error
		if len(args) == 1 {
			p, err = store.Get(args[0])
		} else {
			p, err = store.GetActive()
		}
		if err != nil || p == nil {
			return fmt.Errorf("no profile to sync — pass a nickname")
		}
		if p.GHUser == "" {
			return fmt.Errorf("profile '%s' has no GitHub account — run: gitswitch login", p.Nickname)
		}
		dry, _ := cmd.Flags().GetBool("dry-run")
		res, err := syncXcodeAccount(xcode.NewStore(), p.GHUser, dry)
		if err != nil {
			return err
		}
		switch {
		case !res.Found:
			return fmt.Errorf("%s isn't added to Xcode yet — run: gitswitch xcode add %s", p.GHUser, p.Nickname)
		case !res.Changed:
			fmt.Printf("%s Xcode already prefers %s\n", styleOK("✓"), p.GHUser)
		case dry:
			fmt.Printf("Would make %s Xcode's default (order: %s)\n", p.GHUser, strings.Join(res.Order, ", "))
		default:
			fmt.Printf("%s Xcode now prefers %s (order: %s)\n", styleOK("✓"), p.GHUser, strings.Join(res.Order, ", "))
			if xcode.IsRunning() {
				fmt.Printf("%s Xcode is running — restart it to apply\n", styleWarn("⚠"))
			}
		}
		return nil
	},
}

var xcodeAddCmd = &cobra.Command{
	Use:   "add [nickname]",
	Short: "Add a profile's GitHub account to Xcode (copies its token for one-paste setup)",
	Long: "Xcode only accepts new accounts through " + xcodeSettingsPath + " — its token\n" +
		"storage is closed to other tools. This copies the profile's token from the gh CLI\n" +
		"to your clipboard so setup is: + → GitHub → username → paste.\n\n" +
		"Use --pat to open GitHub's token page with Xcode's scopes pre-selected instead.",
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireXcode(); err != nil {
			return err
		}
		var p *storage.Profile
		var err error
		if len(args) == 1 {
			p, err = store.Get(args[0])
		} else {
			p, err = store.GetActive()
		}
		if err != nil || p == nil {
			return fmt.Errorf("no profile selected — pass a nickname")
		}
		if p.GHUser == "" {
			return fmt.Errorf("profile '%s' has no GitHub account — run: gitswitch login", p.Nickname)
		}
		if accts, _ := xcode.NewStore().Accounts(); xcode.FindGitHub(accts, p.GHUser) >= 0 {
			fmt.Printf("%s %s is already in Xcode\n", styleOK("✓"), p.GHUser)
			return nil
		}
		usePAT, _ := cmd.Flags().GetBool("pat")
		if usePAT {
			_ = exec.Command("open", xcodePATURL).Run()
			fmt.Printf("Opened GitHub's token page (scopes pre-selected). Generate it as %s, then:\n", p.GHUser)
		} else {
			out, err := exec.Command("gh", "auth", "token", "--hostname", "github.com", "--user", p.GHUser).Output()
			tok := strings.TrimSpace(string(out))
			if err != nil || tok == "" {
				return fmt.Errorf("no gh token for %s — run 'gh auth login' or use --pat", p.GHUser)
			}
			pb := exec.Command("pbcopy")
			pb.Stdin = strings.NewReader(tok)
			if err := pb.Run(); err != nil {
				return fmt.Errorf("couldn't copy to clipboard: %w", err)
			}
			fmt.Printf("%s Token for %s copied to your clipboard. Now:\n", styleOK("✓"), p.GHUser)
		}
		fmt.Printf("  1. Open %s\n", xcodeSettingsPath)
		fmt.Printf("  2. Click + → GitHub\n")
		fmt.Printf("  3. Account: %s   Token: paste (⌘V)\n", p.GHUser)
		fmt.Printf("  4. Back here: gitswitch xcode sync %s\n", p.Nickname)
		return nil
	},
}

var xcodeBindCmd = &cobra.Command{
	Use:   "bind [nickname]",
	Short: "Pin this repo's https remotes to a profile's account so Xcode always uses it",
	Long: "Embeds the profile's GitHub username in this repo's https remotes\n" +
		"(https://user@github.com/…). Xcode matches that username to its account, so the\n" +
		"repo uses the right one no matter which account is Xcode's default. Uses the\n" +
		"repo's pinned profile when no nickname is given. SSH remotes are left alone.",
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, _ := os.Getwd()
		var p *storage.Profile
		var err error
		if len(args) == 1 {
			p, err = store.Get(args[0])
		} else {
			p, _, err = effectiveProfile(cwd)
		}
		if err != nil || p == nil {
			return fmt.Errorf("no profile to bind — pass a nickname")
		}
		if p.GHUser == "" {
			return fmt.Errorf("profile '%s' has no GitHub account — run: gitswitch login", p.Nickname)
		}
		changed, err := bindRemotes(cwd, p.GHUser)
		if err != nil {
			return err
		}
		if len(changed) == 0 {
			fmt.Printf("Nothing to bind — no https github.com remotes need %s\n", p.GHUser)
			return nil
		}
		fmt.Printf("%s Bound %s to %s (Xcode will use that account here)\n", styleOK("✓"), strings.Join(changed, ", "), p.GHUser)
		return nil
	},
}

var xcodeUnbindCmd = &cobra.Command{
	Use:   "unbind",
	Short: "Remove the embedded username from this repo's https remotes",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, _ := os.Getwd()
		changed, err := unbindRemotes(cwd)
		if err != nil {
			return err
		}
		if len(changed) == 0 {
			fmt.Println("Nothing to unbind")
			return nil
		}
		fmt.Printf("%s Unbound %s\n", styleOK("✓"), strings.Join(changed, ", "))
		return nil
	},
}

var xcodeEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Keep Xcode in step automatically on 'switch' and 'pin'",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireXcode(); err != nil {
			return err
		}
		if err := os.WriteFile(xcodeMarker(), []byte("1\n"), 0o644); err != nil {
			return err
		}
		fmt.Printf("%s Xcode auto-sync on — 'switch' sets Xcode's default account, 'pin' binds the repo's remotes\n", styleOK("✓"))
		if p, _ := store.GetActive(); p != nil {
			if msg := xcodeAutoSync(p); msg != "" {
				fmt.Printf("%s %s\n", styleOK("✓"), msg)
			}
		}
		return nil
	},
}

var xcodeDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Stop managing Xcode automatically",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := os.Remove(xcodeMarker()); err != nil && !os.IsNotExist(err) {
			return err
		}
		fmt.Printf("%s Xcode auto-sync off (existing account order and bound remotes are left as-is)\n", styleOK("✓"))
		return nil
	},
}

func init() {
	xcodeCmd.Flags().Bool("json", false, "Output machine-readable JSON")
	xcodeStatusCmd.Flags().Bool("json", false, "Output machine-readable JSON")
	xcodeSyncCmd.Flags().Bool("dry-run", false, "Show what would change without writing")
	xcodeAddCmd.Flags().Bool("pat", false, "Open GitHub's token page instead of copying the gh token")
	xcodeCmd.AddCommand(xcodeStatusCmd, xcodeAddCmd, xcodeSyncCmd, xcodeBindCmd, xcodeUnbindCmd, xcodeEnableCmd, xcodeDisableCmd)
}
