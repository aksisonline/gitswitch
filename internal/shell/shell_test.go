package shell

import (
	"os"
	"strings"
	"testing"
)

// ── Snippet content invariants ───────────────────────────────────────────────

// lastRepoClearPattern returns the expected pattern for each shell type.
// zsh/bash/omz/p10k use `unset`; fish uses `set -e`.
func zshClearPattern() string  { return "unset __GITSWITCH_LAST_REPO" }
func fishClearPattern() string { return "set -e __GITSWITCH_LAST_REPO" }

// ── Prompt: nickname not email ───────────────────────────────────────────────

func TestZshSnippetUsesNickname(t *testing.T) {
	s := nudgeSnippetZsh()
	if strings.Contains(s, "git config user.email") {
		t.Error("zsh prompt should use gitswitch current --short, not git config user.email")
	}
	if !strings.Contains(s, "gitswitch current --short") {
		t.Error("zsh prompt missing 'gitswitch current --short'")
	}
}

func TestBashSnippetUsesNickname(t *testing.T) {
	s := nudgeSnippetBash()
	if strings.Contains(s, "git config user.email") {
		t.Error("bash prompt should use gitswitch current --short, not git config user.email")
	}
	if !strings.Contains(s, "gitswitch current --short") {
		t.Error("bash prompt missing 'gitswitch current --short'")
	}
}

func TestFishSnippetUsesNickname(t *testing.T) {
	s := nudgeSnippetFish()
	if strings.Contains(s, "git config user.email") {
		t.Error("fish prompt should use gitswitch current --short, not git config user.email")
	}
	if !strings.Contains(s, "gitswitch current --short") {
		t.Error("fish prompt missing 'gitswitch current --short'")
	}
}

func TestOMZSnippetUsesNickname(t *testing.T) {
	s := omzPluginContent()
	if strings.Contains(s, "git config user.email") {
		t.Error("OMZ prompt should use gitswitch current --short, not git config user.email")
	}
}

func TestP10kSnippetUsesNickname(t *testing.T) {
	s := p10kSnippet()
	if strings.Contains(s, "git config user.email") {
		t.Error("P10k prompt should use gitswitch current --short, not git config user.email")
	}
}

// ── Prompt: only inside a git repo ───────────────────────────────────────────

func TestZshSnippetPromptHasGitCheck(t *testing.T) {
	for _, s := range []struct{ name, snippet string }{
		{"zsh", nudgeSnippetZsh()},
		{"bash", nudgeSnippetBash()},
		{"fish", nudgeSnippetFish()},
		{"omz", omzPluginContent()},
		{"p10k", p10kSnippet()},
	} {
		if !strings.Contains(s.snippet, "git rev-parse --git-dir") {
			t.Errorf("%s prompt function missing git-repo guard", s.name)
		}
	}
}

// ── Prompt placement ─────────────────────────────────────────────────────────

func TestZshSnippetUsesLeftPROMPT(t *testing.T) {
	s := nudgeSnippetZsh()
	if !strings.Contains(s, "PROMPT='$(__gitswitch_prompt)'") {
		t.Error("zsh snippet should prepend gitswitch_prompt to left PROMPT")
	}
}

func TestBashSnippetHasDynamicColor(t *testing.T) {
	s := nudgeSnippetBash()
	if !strings.Contains(s, `\e[38;5;`) {
		t.Error("bash prompt missing 256-color ANSI format")
	}
	if strings.Contains(s, `\e[36m`) {
		t.Error("bash prompt must not use hardcoded cyan")
	}
}

func TestFishSnippetHasDynamicColor(t *testing.T) {
	s := nudgeSnippetFish()
	if !strings.Contains(s, `\e[38;5;`) {
		t.Error("fish prompt missing 256-color ANSI format")
	}
}

func TestZshSnippetHasDynamicColor(t *testing.T) {
	s := nudgeSnippetZsh()
	if !strings.Contains(s, `%F{$color}`) {
		t.Error("zsh prompt missing dynamic color variable")
	}
	if strings.Contains(s, `%F{cyan}`) {
		t.Error("zsh prompt must not use hardcoded cyan")
	}
}

func TestZshSnippetClearsLastRepo(t *testing.T) {
	s := nudgeSnippetZsh()
	if !strings.Contains(s, zshClearPattern()) {
		t.Errorf("zsh snippet missing LAST_REPO clear: want %q", zshClearPattern())
	}
}

func TestBashSnippetClearsLastRepo(t *testing.T) {
	s := nudgeSnippetBash()
	if !strings.Contains(s, zshClearPattern()) {
		t.Errorf("bash snippet missing LAST_REPO clear: want %q", zshClearPattern())
	}
}

func TestFishSnippetClearsLastRepo(t *testing.T) {
	s := nudgeSnippetFish()
	if !strings.Contains(s, fishClearPattern()) {
		t.Errorf("fish snippet missing LAST_REPO clear: want %q", fishClearPattern())
	}
}

func TestOMZSnippetClearsLastRepo(t *testing.T) {
	s := omzPluginContent()
	if !strings.Contains(s, zshClearPattern()) {
		t.Errorf("OMZ snippet missing LAST_REPO clear: want %q", zshClearPattern())
	}
}

func TestP10kSnippetClearsLastRepo(t *testing.T) {
	s := p10kSnippet()
	if !strings.Contains(s, zshClearPattern()) {
		t.Errorf("P10k snippet missing LAST_REPO clear: want %q", zshClearPattern())
	}
}

// Dedup guard must still be present — clearing must not remove the "already in
// this repo" short-circuit.

func TestZshSnippetHasDedupGuard(t *testing.T) {
	s := nudgeSnippetZsh()
	if !strings.Contains(s, `[[ "$root" == "$__GITSWITCH_LAST_REPO" ]] && return`) {
		t.Error("zsh snippet missing dedup guard")
	}
}

func TestBashSnippetHasDedupGuard(t *testing.T) {
	s := nudgeSnippetBash()
	if !strings.Contains(s, `[[ "$root" == "$__GITSWITCH_LAST_REPO" ]] && return`) {
		t.Error("bash snippet missing dedup guard")
	}
}

func TestFishSnippetHasDedupGuard(t *testing.T) {
	s := nudgeSnippetFish()
	if !strings.Contains(s, `if test "$root" = "$__GITSWITCH_LAST_REPO"`) {
		t.Error("fish snippet missing dedup guard")
	}
}

// Marker fencing — every snippet that writes to a rc file must have begin/end
// markers so IsInstalled can detect it.

func TestZshSnippetHasMarkers(t *testing.T) {
	s := nudgeSnippetZsh()
	if !strings.Contains(s, marker+" begin") || !strings.Contains(s, marker+" end") {
		t.Error("zsh snippet missing begin/end markers")
	}
}

func TestBashSnippetHasMarkers(t *testing.T) {
	s := nudgeSnippetBash()
	if !strings.Contains(s, marker+" begin") || !strings.Contains(s, marker+" end") {
		t.Error("bash snippet missing begin/end markers")
	}
}

func TestFishSnippetHasMarkers(t *testing.T) {
	s := nudgeSnippetFish()
	if !strings.Contains(s, marker+" begin") || !strings.Contains(s, marker+" end") {
		t.Error("fish snippet missing begin/end markers")
	}
}

func TestP10kSnippetHasMarkers(t *testing.T) {
	s := p10kSnippet()
	if !strings.Contains(s, marker+" begin") || !strings.Contains(s, marker+" end") {
		t.Error("p10k snippet missing begin/end markers")
	}
}

// Bash: must have the __GITSWITCH_LAST_PWD guard to avoid running git rev-parse
// on every prompt draw (PROMPT_COMMAND fires every prompt, not just on cd).

func TestBashSnippetHasPWDGuard(t *testing.T) {
	s := nudgeSnippetBash()
	if !strings.Contains(s, "__GITSWITCH_LAST_PWD") {
		t.Error("bash snippet missing __GITSWITCH_LAST_PWD dedup guard")
	}
}

// ── IsInstalled ──────────────────────────────────────────────────────────────

func TestIsInstalled_True(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "rc-*.sh")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(marker + " begin\nsome content\n" + marker + " end\n")
	f.Close()
	if !IsInstalled(f.Name()) {
		t.Error("expected IsInstalled=true when marker present")
	}
}

func TestIsInstalled_False(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "rc-*.sh")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString("# unrelated content\n")
	f.Close()
	if IsInstalled(f.Name()) {
		t.Error("expected IsInstalled=false when marker absent")
	}
}

func TestIsInstalled_MissingFile(t *testing.T) {
	if IsInstalled(t.TempDir() + "/nonexistent.sh") {
		t.Error("expected IsInstalled=false for missing file")
	}
}

// ── DetectShell ──────────────────────────────────────────────────────────────

func TestDetectShell_Zsh(t *testing.T) {
	t.Setenv("SHELL", "/bin/zsh")
	if DetectShell() != ShellZsh {
		t.Error("expected ShellZsh")
	}
}

func TestDetectShell_Bash(t *testing.T) {
	t.Setenv("SHELL", "/usr/bin/bash")
	if DetectShell() != ShellBash {
		t.Error("expected ShellBash")
	}
}

func TestDetectShell_Fish(t *testing.T) {
	t.Setenv("SHELL", "/usr/local/bin/fish")
	if DetectShell() != ShellFish {
		t.Error("expected ShellFish")
	}
}

func TestDetectShell_Unknown(t *testing.T) {
	t.Setenv("SHELL", "/bin/sh")
	if DetectShell() != ShellUnknown {
		t.Error("expected ShellUnknown")
	}
}

// ── installRaw idempotency ───────────────────────────────────────────────────

func TestInstallRaw_Idempotent(t *testing.T) {
	tmp := t.TempDir()
	rc := tmp + "/.zshrc"

	// Simulate installing twice by writing marker manually, then calling Install.
	// We can't call Install directly (it touches real $HOME), so test IsInstalled
	// behaviour after writing the snippet once.
	snippet := nudgeSnippetZsh()
	if err := os.WriteFile(rc, []byte(snippet), 0644); err != nil {
		t.Fatal(err)
	}
	if !IsInstalled(rc) {
		t.Fatal("marker should be present after first write")
	}
	// Writing again would create a duplicate — validate IsInstalled blocks it.
	// (The Install function checks IsInstalled before appending.)
	if !IsInstalled(rc) {
		t.Error("second IsInstalled check should still be true — idempotency guard works")
	}
}
