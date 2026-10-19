package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aksisonline/gitswitch/internal/git"
	"github.com/aksisonline/gitswitch/internal/history"
	"github.com/aksisonline/gitswitch/internal/storage"
	"github.com/spf13/cobra"
)

// checkGuard verifies that a commit in dir will be authored by the right
// gitswitch profile. It returns nil when the commit should proceed, or an
// actionable error explaining why the commit was blocked and how to fix or
// bypass it.
//
// Enforcement rules:
//  1. Zero configured gitswitch profiles → allow (never brick git before onboarding).
//  2. No git user.email configured → block.
//  3. Effective commit email does not belong to any stored gitswitch profile → block.
//  4. Repo is pinned in history.json to profile P, but commit email != P.Email → block.
//  5. Repo is unpinned (global scope) and learned history recommends profile P
//     over the current profile → block.
func checkGuard(st *storage.Store, dir string) error {
	profiles, err := st.Load()
	if err != nil || len(profiles) == 0 {
		return nil
	}

	repoKey := history.GetRepoKeyForPath(dir)
	if repoKey == "" {
		return nil
	}

	scope, email := git.ResolveCommitEmail(dir)
	if email == "" {
		return fmt.Errorf(
			"gitswitch guard: no git user.email is configured for this commit\n" +
				"  Fix:    gitswitch switch <profile>   (or: gitswitch pin <profile>)\n" +
				"  Bypass: git commit --no-verify",
		)
	}

	committer := st.GetByEmail(email)
	if committer == nil {
		return fmt.Errorf(
			"gitswitch guard: commit author <%s> does not match any gitswitch profile\n"+
				"  Fix:    gitswitch switch <profile>   (or: gitswitch pin <profile>)\n"+
				"  Add it: gitswitch add <nickname> \"<Name>\" %s\n"+
				"  Bypass: git commit --no-verify",
			email, email,
		)
	}

	if pinnedNick := history.GetPinned(repoKey); pinnedNick != "" {
		if want, err := st.Get(pinnedNick); err == nil && want != nil {
			if !strings.EqualFold(email, want.Email) {
				return fmt.Errorf(
					"gitswitch guard: this repo is pinned to '%s' (%s <%s>), but git would commit as '%s' (<%s>)\n"+
						"  Fix:    gitswitch pin %s   (or: gitswitch switch %s)\n"+
						"  Bypass: git commit --no-verify",
					want.Nickname, want.UserName, want.Email,
					committer.Nickname, email,
					want.Nickname, want.Nickname,
				)
			}
		}
		return nil
	}

	if scope == git.ScopeGlobal {
		if recNick, ok := history.Recommend(repoKey, committer.Nickname); ok {
			if want, err := st.Get(recNick); err == nil && want != nil {
				return fmt.Errorf(
					"gitswitch guard: this repo usually commits as '%s' (%s <%s>), but active profile is '%s' (<%s>)\n"+
						"  Switch: gitswitch switch %s\n"+
						"  Pin it: gitswitch pin %s\n"+
						"  Allow:  gitswitch pin %s   (or bypass once: git commit --no-verify)",
					want.Nickname, want.UserName, want.Email,
					committer.Nickname, email,
					want.Nickname, want.Nickname, committer.Nickname,
				)
			}
		}
	}

	return nil
}

var guardCmd = &cobra.Command{
	Use:   "guard [enable|disable|check|status]",
	Short: "Enforce commit identity with a pre-commit guard hook",
	Long: `Blocks wrong-author commits before they enter git history.

When enabled, gitswitch installs a global pre-commit hook (in
~/.config/gitswitch/hooks, chaining into any existing repo .git/hooks/pre-commit
so Husky, lefthook, and pre-commit still run) that refuses a commit when:
  • the commit email doesn't match any gitswitch profile
  • the repo is pinned to a different profile
  • the repo's learned history strongly points to a different profile`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGuardStatus(cmd)
	},
}

var guardEnableCmd = &cobra.Command{
	Use:     "enable",
	Aliases: []string{"install", "on"},
	Short:   "Enable the global pre-commit identity guard hook",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := git.InstallGuardHook(); err != nil {
			return err
		}
		fmt.Printf("%s Pre-commit guard enabled — wrong-author commits will be refused before entering history\n", styleOK("✓"))
		return nil
	},
}

var guardDisableCmd = &cobra.Command{
	Use:     "disable",
	Aliases: []string{"uninstall", "off"},
	Short:   "Disable the global pre-commit identity guard hook",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := git.UninstallGuardHook(); err != nil {
			return err
		}
		fmt.Printf("%s Pre-commit guard disabled\n", styleOK("✓"))
		return nil
	},
}

var guardCheckCmd = &cobra.Command{
	Use:           "check",
	Short:         "Verify the current repo's commit identity (invoked by pre-commit hook)",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString("path")
		if path == "" {
			var err error
			path, err = os.Getwd()
			if err != nil {
				return err
			}
		}
		return checkGuard(store, path)
	},
}

var guardStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show pre-commit guard status and check the current repo",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGuardStatus(cmd)
	},
}

func runGuardStatus(cmd *cobra.Command) error {
	jsonOut, _ := cmd.Flags().GetBool("json")
	cwd, _ := os.Getwd()
	installed := git.IsGuardHookInstalled()
	checkErr := checkGuard(store, cwd)
	_, email := git.ResolveCommitEmail(cwd)
	var nick string
	if p := store.GetByEmail(email); p != nil {
		nick = p.Nickname
	}

	if jsonOut {
		out := struct {
			Installed bool   `json:"installed"`
			RepoOK    bool   `json:"repo_ok"`
			Nickname  string `json:"nickname,omitempty"`
			Email     string `json:"email,omitempty"`
			Error     string `json:"error,omitempty"`
		}{
			Installed: installed,
			RepoOK:    checkErr == nil,
			Nickname:  nick,
			Email:     email,
		}
		if checkErr != nil {
			out.Error = checkErr.Error()
		}
		b, err := json.MarshalIndent(out, "", "  ")
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", b)
		return nil
	}

	if installed {
		fmt.Printf("%s Pre-commit guard: active\n", styleOK("✓"))
	} else {
		fmt.Printf("%s Pre-commit guard: off (run 'gitswitch guard enable' to turn on)\n", styleWarn("⚠"))
	}

	if history.GetRepoKeyForPath(cwd) != "" {
		if checkErr != nil {
			fmt.Printf("%s Current repo check: blocked\n  %s\n", styleErr("✗"), strings.ReplaceAll(checkErr.Error(), "\n", "\n  "))
		} else if nick != "" {
			fmt.Printf("%s Current repo check: ok (%s <%s>)\n", styleOK("✓"), nick, email)
		}
	}
	return nil
}

func init() {
	guardCmd.Flags().Bool("json", false, "Output machine-readable JSON")
	guardStatusCmd.Flags().Bool("json", false, "Output machine-readable JSON")
	guardCheckCmd.Flags().String("path", "", "Directory to check (default: current working directory)")
	guardCmd.AddCommand(guardEnableCmd, guardDisableCmd, guardCheckCmd, guardStatusCmd)
}
