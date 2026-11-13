package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/aksisonline/gitswitch/internal/storage"
	"github.com/spf13/cobra"
)

// `--v` and `-V` are what people type when `--version` doesn't come to mind.
func TestNormalizeArgs(t *testing.T) {
	cases := map[string]string{"--v": "--version", "-V": "--version", "--version": "--version", "-v": "-v", "list": "list"}
	for in, want := range cases {
		args := []string{"gitswitch", in}
		normalizeArgs(args)
		if args[1] != want {
			t.Errorf("normalizeArgs(%q) = %q, want %q", in, args[1], want)
		}
	}
	// Only as the sole argument — never rewrite something inside a command.
	args := []string{"gitswitch", "add", "--v"}
	normalizeArgs(args)
	if args[2] != "--v" {
		t.Errorf("normalizeArgs rewrote an argument of a subcommand: %v", args)
	}
}

// Commands that only talk about gitswitch itself must work on a machine
// without git — including the version flag the Homebrew formula test runs.
func TestSelfCommandsAllowMissingGit(t *testing.T) {
	for _, arg := range []string{"--version", "-v", "--v", "-V", "version", "--help", "-h", "help", "upgrade", "beta", "stable", "completion"} {
		if !commandAllowsMissingGit([]string{"gitswitch", arg}) {
			t.Errorf("gitswitch %s should run without git installed", arg)
		}
	}
	for _, arg := range []string{"switch", "pin", "list", "login"} {
		if commandAllowsMissingGit([]string{"gitswitch", arg}) {
			t.Errorf("gitswitch %s should require git", arg)
		}
	}
}

func TestArgRangeMessages(t *testing.T) {
	cmd := &cobra.Command{Use: "add <nickname> <user-name> <email>"}
	root := &cobra.Command{Use: "gitswitch"}
	root.AddCommand(cmd)

	err := exactArgs(3)(cmd, []string{"work"})
	var ue *usageError
	if !errors.As(err, &ue) {
		t.Fatalf("want a usageError, got %v", err)
	}
	if want := "gitswitch add <nickname> <user-name> <email>"; !strings.Contains(err.Error(), want) {
		t.Errorf("error %q should spell out the expected arguments (%q)", err, want)
	}
	if err := exactArgs(3)(cmd, []string{"a", "b", "c"}); err != nil {
		t.Errorf("3 arguments rejected: %v", err)
	}
	if err := noArgs(cmd, []string{"x"}); err == nil || !strings.Contains(err.Error(), "takes no arguments") {
		t.Errorf("noArgs error = %v", err)
	}
}

func TestUnknownTargetSuggestsProfilesAndCommands(t *testing.T) {
	root := &cobra.Command{Use: "gitswitch"}
	root.SuggestionsMinimumDistance = 2
	root.AddCommand(&cobra.Command{Use: "list", Run: func(*cobra.Command, []string) {}})
	profiles := []storage.Profile{{Nickname: "work"}, {Nickname: "personal"}}

	if err := unknownTargetError(root, "wrk", profiles); !strings.Contains(err.Error(), "work") {
		t.Errorf("a near-miss profile name should be suggested: %v", err)
	}
	if err := unknownTargetError(root, "lst", profiles); !strings.Contains(err.Error(), "list") {
		t.Errorf("a near-miss command name should be suggested: %v", err)
	}
	if err := unknownTargetError(root, "zzzzzz", profiles); strings.Contains(err.Error(), "did you mean") {
		t.Errorf("nothing close — no suggestion expected: %v", err)
	}
}

func TestOneOf(t *testing.T) {
	cmd := &cobra.Command{Use: "shell"}
	if err := oneOf(cmd, "shell", "", "", "zsh", "bash", "fish"); err != nil {
		t.Errorf("empty (auto-detect) rejected: %v", err)
	}
	err := oneOf(cmd, "shell", "tcsh", "", "zsh", "bash", "fish")
	if err == nil || !strings.Contains(err.Error(), "zsh, bash, fish") {
		t.Errorf("want the allowed values listed, got %v", err)
	}
}

// Every user-facing command must show up in root help, the plumbing must not,
// and help must exit cleanly (the install script uses `--help` as its sanity
// check on a freshly downloaded binary).
func TestRootHelpListsCommands(t *testing.T) {
	root := &cobra.Command{Use: "gitswitch", Short: "Switch git identities"}
	root.AddGroup(&cobra.Group{ID: "identity", Title: "Identity"})
	root.AddCommand(
		&cobra.Command{Use: "switch <nickname>", Short: "Switch to a profile", GroupID: "identity", Run: func(*cobra.Command, []string) {}},
		&cobra.Command{Use: "record", Short: "plumbing", Hidden: true, Run: func(*cobra.Command, []string) {}},
	)
	installHelp(root)
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetArgs([]string{"--help"})
	if err := root.Execute(); err != nil {
		t.Fatalf("--help returned an error: %v", err)
	}
	help := out.String()
	for _, want := range []string{"Usage", "Identity", "switch", "Switch to a profile", "--version", "gitswitch <nickname>"} {
		if !strings.Contains(help, want) {
			t.Errorf("root help is missing %q:\n%s", want, help)
		}
	}
	if strings.Contains(help, "record") {
		t.Errorf("hidden plumbing command leaked into help:\n%s", help)
	}
}
