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

func TestPromptSnippetsUsePromptFlag(t *testing.T) {
	cases := []struct{ name, snippet string }{
		{"zsh", nudgeSnippetZsh()},
		{"bash", nudgeSnippetBash()},
		{"fish", nudgeSnippetFish()},
		{"omz", omzPluginContent()},
		{"p10k", p10kSnippet()},
	}
	for _, c := range cases {
		if strings.Contains(c.snippet, "git config user.email") {
			t.Errorf("%s: prompt must not use git config user.email", c.name)
		}
		if !strings.Contains(c.snippet, "gitswitch current --prompt") {
			t.Errorf("%s: prompt must use 'gitswitch current --prompt'", c.name)
		}
		if strings.Contains(c.snippet, "gitswitch current --short") && c.name != "starship" {
			t.Errorf("%s: prompt must not use --short (reserved for Starship/scripts)", c.name)
		}
	}
}

func TestStarshipSnippetUsesShortFlag(t *testing.T) {
	s := starshipSnippet()
	if !strings.Contains(s, "gitswitch current --short") {
		t.Error("starship snippet must use --short (outputs nick+email for display)")
	}
	if strings.Contains(s, "gitswitch current --prompt") {
		t.Error("starship snippet must not use --prompt")
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

	snippet := nudgeSnippetZsh()
	if err := os.WriteFile(rc, []byte(snippet), 0644); err != nil {
		t.Fatal(err)
	}
	if !IsInstalled(rc) {
		t.Fatal("marker should be present after first write")
	}
	if !IsInstalled(rc) {
		t.Error("second IsInstalled check should still be true — idempotency guard works")
	}
}

// ── removeMarkerBlock ────────────────────────────────────────────────────────

func TestRemoveMarkerBlock_RemovesBlock(t *testing.T) {
	tmp := t.TempDir()
	rc := tmp + "/rc.sh"
	content := "before\n" + nudgeSnippetZsh() + "\nafter\n"
	if err := os.WriteFile(rc, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	if err := removeMarkerBlock(rc); err != nil {
		t.Fatalf("removeMarkerBlock: %v", err)
	}
	got, _ := os.ReadFile(rc)
	if strings.Contains(string(got), marker) {
		t.Error("marker still present after removal")
	}
	if !strings.Contains(string(got), "before") || !strings.Contains(string(got), "after") {
		t.Error("surrounding content was stripped")
	}
}

func TestRemoveMarkerBlock_Idempotent(t *testing.T) {
	tmp := t.TempDir()
	rc := tmp + "/rc.sh"
	content := "before\n" + nudgeSnippetZsh() + "\nafter\n"
	if err := os.WriteFile(rc, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	if err := removeMarkerBlock(rc); err != nil {
		t.Fatal(err)
	}
	// second call on file with no markers should be a no-op
	if err := removeMarkerBlock(rc); err != nil {
		t.Errorf("second removeMarkerBlock should be no-op, got: %v", err)
	}
}

func TestRemoveMarkerBlock_PreservesMode(t *testing.T) {
	tmp := t.TempDir()
	rc := tmp + "/rc.sh"
	content := nudgeSnippetZsh()
	if err := os.WriteFile(rc, []byte(content), 0755); err != nil {
		t.Fatal(err)
	}
	if err := removeMarkerBlock(rc); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(rc)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode() != 0755 {
		t.Errorf("file mode changed: got %v, want 0755", info.Mode())
	}
}

func TestRemoveMarkerBlock_UnbalancedBegin(t *testing.T) {
	tmp := t.TempDir()
	rc := tmp + "/rc.sh"
	content := marker + " begin\nsome content\n" // no end marker
	if err := os.WriteFile(rc, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	if err := removeMarkerBlock(rc); err == nil {
		t.Error("expected error for begin-without-end marker")
	}
}

func TestRemoveMarkerBlock_StarshipBlock(t *testing.T) {
	tmp := t.TempDir()
	toml := tmp + "/starship.toml"
	content := "[palettes]\n" + marker + " begin\n" + starshipSnippet() + "\n" + marker + " end\n"
	if err := os.WriteFile(toml, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	if err := removeMarkerBlock(toml); err != nil {
		t.Fatalf("removeMarkerBlock on starship.toml: %v", err)
	}
	got, _ := os.ReadFile(toml)
	if strings.Contains(string(got), "custom.gitswitch") {
		t.Error("starship block not removed")
	}
	if !strings.Contains(string(got), "[palettes]") {
		t.Error("surrounding toml content was stripped")
	}
}
