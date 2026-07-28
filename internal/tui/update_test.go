package tui

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aksisonline/gitswitch/internal/git"
	"github.com/aksisonline/gitswitch/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

// TestProfileFormFieldTitlesMatchForm guards against the mouse map (in
// update.go) drifting from the actual field titles (in form.go) — a rename
// there would silently break click-to-focus without this.
func TestProfileFormFieldTitlesMatchForm(t *testing.T) {
	d := &profileFormData{}
	form := newProfileForm(d, false, 60)
	form.Init()
	rendered := form.View()
	for _, f := range profileFormFieldTitles {
		if !strings.Contains(rendered, f.title) {
			t.Errorf("expected form to render title %q for key %q, got:\n%s", f.title, f.key, rendered)
		}
	}
}

// TestScopeGlyph pins the state-column glyphs that tell the user which identity
// git will actually use: the repo's, this terminal's, or the global one.
func TestScopeGlyph(t *testing.T) {
	profiles := []storage.Profile{{Nickname: "personal"}, {Nickname: "work"}}

	cases := []struct {
		name         string
		scope        git.Scope
		scopeProfile *storage.Profile
		active       *storage.Profile
		arcade       bool
		want         [2]string // glyph for personal, work
	}{
		{"global only", git.ScopeGlobal, nil, &profiles[0], false, [2]string{"✓", "·"}},
		{"pinned elsewhere than active", git.ScopeRepo, &profiles[1], &profiles[0], false, [2]string{"✓", "●"}},
		{"pinned to the active profile", git.ScopeRepo, &profiles[0], &profiles[0], false, [2]string{"◉", "·"}},
		{"session identity", git.ScopeSession, &profiles[1], &profiles[0], false, [2]string{"✓", "◆"}},
		{"arcade global", git.ScopeGlobal, nil, &profiles[0], true, [2]string{"★", "·"}},
		{"arcade pinned", git.ScopeRepo, &profiles[1], &profiles[0], true, [2]string{"★", "◆"}},
		{"arcade session", git.ScopeSession, &profiles[1], &profiles[0], true, [2]string{"★", "▲"}},
	}
	for _, c := range cases {
		m := Model{profiles: profiles, active: c.active, scope: c.scope, scopeProfile: c.scopeProfile, arcadeMode: c.arcade}
		for i, want := range c.want {
			if got := m.scopeGlyph(profiles[i]); got != want {
				t.Errorf("%s: glyph for %q = %q, want %q", c.name, profiles[i].Nickname, got, want)
			}
		}
	}
}

// The TUI must notice the repo it was launched from: without this wiring the
// glyphs and the Current line silently report the global identity while git uses
// the pinned one.
func TestNewDetectsRepoPin(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	gitconfig := filepath.Join(home, "gitconfig")
	if err := os.WriteFile(gitconfig, nil, 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GIT_CONFIG_GLOBAL", gitconfig)
	if err := exec.Command("git", "config", "--global", "user.email", "alice@personal.dev").Run(); err != nil {
		t.Fatal(err)
	}

	st, err := storage.New()
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Add("work", "Alice W", "alice@work.com", "", "", ""); err != nil {
		t.Fatal(err)
	}

	repo := t.TempDir()
	for _, args := range [][]string{
		{"-C", repo, "init"},
		{"-C", repo, "remote", "add", "origin", "git@github.com:acme/widget.git"},
		{"-C", repo, "config", "--local", "user.email", "alice@work.com"},
	} {
		if err := exec.Command("git", args...).Run(); err != nil {
			t.Fatalf("git %v: %v", args, err)
		}
	}
	t.Chdir(repo)

	m, err := New(st, "v0.0.0-test")
	if err != nil {
		t.Fatal(err)
	}
	if m.scope != git.ScopeRepo {
		t.Errorf("scope = %v, want ScopeRepo", m.scope)
	}
	if m.scopeProfile == nil || m.scopeProfile.Nickname != "work" {
		t.Errorf("scopeProfile = %v, want the 'work' profile", m.scopeProfile)
	}
	if m.repoKey == "" {
		t.Error("repoKey should be set inside a repo")
	}
	if got := m.scopeGlyph(storage.Profile{Nickname: "work"}); got != "●" {
		t.Errorf("glyph for the pinned profile = %q, want ●", got)
	}
}

// Pressing p outside a git repo must say so rather than silently doing nothing —
// there is no repo to pin to and no local config to write.
func TestPinKeyOutsideRepo(t *testing.T) {
	m := Model{
		store:    storage.NewAt(t.TempDir()),
		profiles: []storage.Profile{{Nickname: "personal", Email: "a@b.dev"}},
		state:    StateList,
		width:    80,
		height:   24,
	}
	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	got, ok := next.(Model)
	if !ok {
		t.Fatalf("Update returned %T, want Model", next)
	}
	if !got.statusIsErr || !strings.Contains(got.statusMsg, "not inside a git repo") {
		t.Errorf("statusMsg = %q (isErr=%v), want a 'not inside a git repo' error", got.statusMsg, got.statusIsErr)
	}
}

// TestFormFieldAt pins the add/edit form's click-to-field mapping: huh has no
// mouse support of its own, so a click anywhere in a field's block (title,
// description, input line) must resolve to that field's key.
func TestFormFieldAt(t *testing.T) {
	lines := []string{
		"  Nickname",
		"  A label for this identity",
		"  > personal",
		"",
		"  Name",
		"  git user.name",
		"  > Ada Lovelace",
		"",
		"  Signing key",
		"  optional",
		"  > ",
	}
	cases := []struct {
		relY int
		want string
		ok   bool
	}{
		{0, "nickname", true},
		{1, "nickname", true},
		{2, "nickname", true},
		{3, "nickname", true}, // blank spacer still belongs to the field above
		{4, "userName", true},
		{6, "userName", true},
		{8, "signKey", true},
		{10, "signKey", true},
	}
	for _, c := range cases {
		got, ok := formFieldAt(lines, c.relY)
		if ok != c.ok || got != c.want {
			t.Errorf("formFieldAt(relY=%d) = (%q, %v), want (%q, %v)", c.relY, got, ok, c.want, c.ok)
		}
	}

	if _, ok := formFieldAt(lines, -1); ok {
		t.Error("formFieldAt(relY=-1) should report not-found")
	}
	if _, ok := formFieldAt(nil, 0); ok {
		t.Error("formFieldAt(nil lines) should report not-found")
	}
}
