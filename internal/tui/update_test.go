package tui

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aksisonline/gitswitch/internal/git"
	"github.com/aksisonline/gitswitch/internal/shell"
	"github.com/aksisonline/gitswitch/internal/storage"
	ver "github.com/aksisonline/gitswitch/internal/version"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
)

// runCmd executes cmd and returns its message. Commands that start an
// animation come back batched with the frame ticker; this unwraps the batch
// and returns the first message that isn't a frame tick.
func runCmd(cmd tea.Cmd) tea.Msg {
	if cmd == nil {
		return nil
	}
	msg := cmd()
	batch, ok := msg.(tea.BatchMsg)
	if !ok {
		return msg
	}
	for _, c := range batch {
		if m := runCmd(c); m != nil {
			if _, tick := m.(frameMsg); !tick {
				return m
			}
		}
	}
	return nil
}

// TestProfileFormFieldTitlesMatchForm guards against the mouse map (in
// update.go) drifting from the actual field titles (in form.go) — a rename
// there would silently break click-to-focus without this.
func TestProfileFormFieldTitlesMatchForm(t *testing.T) {
	d := &profileFormData{}
	form := newProfileForm(d, false, 60, 0)
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

// Pressing enter on the HTTPS Credential Helper item must actually flip git
// config, not just the toggle glyph — this is the toggle that used to be a
// permanently-disabled "coming soon" stub.
func TestCredentialHelperToggle(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("GIT_CONFIG_GLOBAL", filepath.Join(home, "gitconfig"))
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	if err := os.WriteFile(filepath.Join(home, "gitconfig"), nil, 0600); err != nil {
		t.Fatal(err)
	}

	m := Model{
		store: storage.NewAt(t.TempDir()), state: StateList, width: 84, height: 30,
		tabIndex: 1, utilityFocus: 2, credentialHelperEnabled: false,
	}
	if git.IsCredentialHelperInstalled() {
		t.Fatal("should not be installed in a fresh config")
	}

	next, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := next.(Model)
	if cmd == nil {
		t.Fatal("enter on the credential helper item should return a command")
	}
	msg := runCmd(cmd)
	got2, _ := got.Update(msg)
	final := got2.(Model)
	if !final.credentialHelperEnabled {
		t.Errorf("credentialHelperEnabled = false after enabling, statusMsg=%q", final.statusMsg)
	}
	if !git.IsCredentialHelperInstalled() {
		t.Error("git config should show gitswitch installed after the toggle")
	}

	// Toggle again: must remove it.
	next, cmd = final.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got = next.(Model)
	next2, _ := got.Update(runCmd(cmd))
	final = next2.(Model)
	if final.credentialHelperEnabled {
		t.Error("credentialHelperEnabled = true after disabling")
	}
	if git.IsCredentialHelperInstalled() {
		t.Error("git config should show gitswitch removed after the second toggle")
	}
}

func TestGuardToggle(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("GIT_CONFIG_GLOBAL", filepath.Join(home, "gitconfig"))
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	if err := os.WriteFile(filepath.Join(home, "gitconfig"), nil, 0600); err != nil {
		t.Fatal(err)
	}

	m := Model{
		store: storage.NewAt(t.TempDir()), state: StateList, width: 84, height: 30,
		tabIndex: 1, utilityFocus: 4, guardEnabled: false,
	}
	if git.IsGuardHookInstalled() {
		t.Fatal("guard hook should not be installed in a fresh config")
	}

	next, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := next.(Model)
	if cmd == nil {
		t.Fatal("enter on the guard item should return a command")
	}
	got2, _ := got.Update(runCmd(cmd))
	final := got2.(Model)
	if !final.guardEnabled {
		t.Errorf("guardEnabled = false after enabling, statusMsg=%q", final.statusMsg)
	}
	if !git.IsGuardHookInstalled() {
		t.Error("git config should show guard hook installed after toggle")
	}

	// Toggle again: must remove it.
	next, cmd = final.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got = next.(Model)
	next2, _ := got.Update(runCmd(cmd))
	final = next2.(Model)
	if final.guardEnabled {
		t.Error("guardEnabled = true after disabling")
	}
	if git.IsGuardHookInstalled() {
		t.Error("git config should show guard hook removed after second toggle")
	}
}

func TestAutoPinToggle(t *testing.T) {
	m := Model{
		store: storage.NewAt(t.TempDir()), state: StateList, width: 84, height: 30,
		tabIndex: 1, utilityFocus: 3, autoPinDisabled: false,
	}

	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	final := next.(Model)
	if !final.autoPinDisabled {
		t.Errorf("autoPinDisabled = false after toggling off, statusMsg=%q", final.statusMsg)
	}

	next, _ = final.Update(tea.KeyMsg{Type: tea.KeyEnter})
	final = next.(Model)
	if final.autoPinDisabled {
		t.Error("autoPinDisabled = true after toggling back on")
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
	// A repo pin only counts as active when Session Isolation is on.
	if _, err := shell.InstallGHWrapper(shell.DetectShell()); err != nil {
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

	m, err := New(st, "v0.0.0-test", "")
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

// Arcade mode persists across launches via `gitswitch pacman`, so the intro must
// trigger from the saved pref on every launch, not just the one-shot toggle.
func TestNewPersistedArcadeModeShowsIntro(t *testing.T) {
	st := storage.NewAt(t.TempDir())
	if err := st.SavePrefs(storage.Prefs{ArcadeMode: true}); err != nil {
		t.Fatal(err)
	}

	m, err := New(st, "v0.0.0-test", "")
	if err != nil {
		t.Fatal(err)
	}
	if m.state != StateIntro {
		t.Errorf("state = %v, want StateIntro", m.state)
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

// TestMergeIntoFoundNicknameCollision pins the existing behavior: two
// candidates sharing a nickname merge into one, missing fields filled in.
func TestMergeIntoFoundNicknameCollision(t *testing.T) {
	found := []storage.Profile{{Nickname: "alice", Email: "alice@gmail.com"}}
	seen := map[string]bool{"alice@gmail.com": true}

	found = mergeIntoFound(found, seen, storage.Profile{
		Nickname: "alice",
		GHUser:   "alice-corp",
	}, nil)

	if len(found) != 1 {
		t.Fatalf("expected merge, got %d profiles", len(found))
	}
	if found[0].GHUser != "alice-corp" {
		t.Errorf("expected GHUser filled in from nickname-collision merge, got %q", found[0].GHUser)
	}
}

// TestMergeIntoFoundVerifiedEmailMatch is the new tier: a gh account whose
// nickname and placeholder email don't coincide with the git-config
// candidate must still merge when its verified GitHub email matches.
func TestMergeIntoFoundVerifiedEmailMatch(t *testing.T) {
	found := []storage.Profile{{Nickname: "alice", Email: "alice@gmail.com"}}
	seen := map[string]bool{"alice@gmail.com": true}

	found = mergeIntoFound(found, seen, storage.Profile{
		Nickname: "alice-corp",
		Email:    "alice-corp@users.noreply.github.com",
		GHUser:   "alice-corp",
	}, []string{"someone-else@example.com", "Alice@Gmail.com"}) // case-insensitive match

	if len(found) != 1 {
		t.Fatalf("expected verified-email match to merge into one profile, got %d", len(found))
	}
	if found[0].GHUser != "alice-corp" {
		t.Errorf("expected GHUser set via verified-email merge, got %q", found[0].GHUser)
	}
	if found[0].Nickname != "alice" {
		t.Errorf("merge must keep the earlier candidate's nickname, got %q", found[0].Nickname)
	}
}

// TestMergeIntoFoundVerifiedEmailNoMatch confirms a non-matching verified
// email just appends a new candidate, same as if none had been fetched.
func TestMergeIntoFoundVerifiedEmailNoMatch(t *testing.T) {
	found := []storage.Profile{{Nickname: "alice", Email: "alice@gmail.com"}}
	seen := map[string]bool{"alice@gmail.com": true}

	found = mergeIntoFound(found, seen, storage.Profile{
		Nickname: "bob-corp",
		Email:    "bob-corp@users.noreply.github.com",
		GHUser:   "bob-corp",
	}, []string{"bob@work.com"})

	if len(found) != 2 {
		t.Fatalf("expected no match to append as a separate candidate, got %d profiles", len(found))
	}
}

// TestMergeIntoFoundLiteralEmailCoincidence pins the existing fallback tier:
// a second candidate with no GHUser and an already-seen literal email is
// dropped rather than appended.
func TestMergeIntoFoundLiteralEmailCoincidence(t *testing.T) {
	found := []storage.Profile{{Nickname: "alice", Email: "alice@gmail.com"}}
	seen := map[string]bool{"alice@gmail.com": true}

	found = mergeIntoFound(found, seen, storage.Profile{
		Nickname: "alice-dup",
		Email:    "alice@gmail.com",
	}, nil)

	if len(found) != 1 {
		t.Fatalf("expected literal-email coincidence with no GHUser to be dropped, got %d profiles", len(found))
	}
}

// TestMergeIntoFoundRecordsEmailEvenWhenMergedViaEarlierTier is a regression
// test for a real gap Copilot review caught: a candidate whose own Email is
// only recorded in seenEmail inside tier 3's branch means a merge via tier 1
// or tier 2 (which returns before reaching that code) never marks its email
// as seen — so a later, unrelated candidate sharing that same literal email
// would slip past tier 3's dedup instead of being caught.
func TestMergeIntoFoundRecordsEmailEvenWhenMergedViaEarlierTier(t *testing.T) {
	found := []storage.Profile{{Nickname: "alice", Email: "alice@gmail.com"}}
	seen := map[string]bool{"alice@gmail.com": true}

	// Merges via tier 1 (nickname collision) — found[0].Email is already
	// set, so mergeProfileFields won't overwrite it with this profile's own
	// email, but that email must still be recorded in seenEmail.
	found = mergeIntoFound(found, seen, storage.Profile{
		Nickname: "alice",
		Email:    "alice-work@gmail.com",
		GHUser:   "alice-corp",
	}, nil)
	if len(found) != 1 {
		t.Fatalf("expected nickname-collision merge, got %d profiles", len(found))
	}

	// A later, unrelated candidate sharing that same literal email (no
	// GHUser) must be deduped, not appended as a spurious second profile.
	found = mergeIntoFound(found, seen, storage.Profile{
		Nickname: "someone-else",
		Email:    "alice-work@gmail.com",
	}, nil)
	if len(found) != 1 {
		t.Fatalf("expected literal-email dedup to catch the merged-away candidate's email, got %d profiles", len(found))
	}
}

// The version check lands a moment after launch. It may open the update
// prompt only if the user hasn't touched anything yet — otherwise the screen
// swaps under their hands and "enter to switch" becomes "enter to upgrade".
func TestUpdatePromptNeverInterruptsTheUser(t *testing.T) {
	base := Model{store: storage.NewAt(t.TempDir()), state: StateList, width: 80, height: 24, currentVersion: "v0.5.0"}
	check := versionCheckMsg{latest: "v0.6.0"}

	next, _ := base.Update(check)
	if got := next.(Model); got.state != StateUpdatePrompt {
		t.Errorf("untouched: state = %d, want the update prompt", got.state)
	}

	touched, _ := base.Update(tea.KeyMsg{Type: tea.KeyDown})
	next, _ = touched.(Model).Update(check)
	got := next.(Model)
	if got.state != StateList {
		t.Errorf("after a keypress the prompt must not open by itself (state = %d)", got.state)
	}
	if !got.updateAvailable {
		t.Error("the update should still be flagged for the status line and the u key")
	}
	next, _ = got.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'u'}})
	if next.(Model).state != StateUpdatePrompt {
		t.Error("u should open the update prompt on request")
	}

	// "Not now" is remembered for that release.
	prompt, _ := base.Update(check)
	skipped, _ := prompt.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	if s := skipped.(Model); s.state != StateList || s.skippedUpdate != "v0.6.0" {
		t.Fatalf("n should close the prompt and remember the release (state=%d skipped=%q)", s.state, s.skippedUpdate)
	}
	again := base
	again.skippedUpdate = "v0.6.0"
	next, _ = again.Update(check)
	if next.(Model).state != StateList {
		t.Error("a skipped release must not reopen the prompt on the next launch")
	}
}

// The in-app upgrade must go through Homebrew when Homebrew owns the binary.
func TestUpdatePromptNamesTheUpdateMechanism(t *testing.T) {
	m := Model{store: storage.NewAt(t.TempDir()), state: StateUpdatePrompt, width: 80, height: 24,
		currentVersion: "v0.5.0", latestVersion: "v0.6.0",
		install: ver.Install{Method: ver.MethodBrew, Formula: ver.StableFormula}}
	if view := ansi.Strip(m.View()); !strings.Contains(view, "brew upgrade gitswitch") {
		t.Errorf("a Homebrew install should be told it updates with brew:\n%s", view)
	}
	m.install = ver.Install{}
	if view := ansi.Strip(m.View()); strings.Contains(view, "brew") {
		t.Errorf("a script install should not mention brew:\n%s", view)
	}
}

// First-run setup, driven by keys: welcome → scan → (nothing found) → add
// account → manual form. Skipping from the welcome screen must work from the
// keyboard too — it used to be reachable only with the mouse.
func TestSetupWizardKeyboardFlow(t *testing.T) {
	// The scan reads git config, gh accounts and ~/.ssh — point all of them
	// at an empty home so the test never touches (or depends on) this machine.
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("GH_CONFIG_DIR", filepath.Join(home, "gh"))
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("GIT_CONFIG_GLOBAL", filepath.Join(home, "gitconfig"))
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")

	m := Model{store: storage.NewAt(t.TempDir()), state: StateWizardWelcome, onboarding: true}
	key := func(m Model, k tea.KeyMsg) Model {
		next, cmd := m.Update(k)
		got := next.(Model)
		if msg := runCmd(cmd); msg != nil {
			if _, quit := msg.(tea.QuitMsg); !quit {
				next, _ = got.Update(msg)
				got = next.(Model)
			}
		}
		return got
	}
	enter, down := tea.KeyMsg{Type: tea.KeyEnter}, tea.KeyMsg{Type: tea.KeyDown}

	if got := key(key(m, down), enter); got.state != StateList || got.onboarding {
		t.Errorf("'Skip for now' should leave setup (state=%d onboarding=%v)", got.state, got.onboarding)
	}
	if got := key(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}); got.state != StateList {
		t.Errorf("s should skip setup, state=%d", got.state)
	}

	// Start setup: nothing to find in an empty home, so the scan goes
	// straight to adding an account.
	got := key(m, enter)
	if got.state != StateWizardAddMore {
		t.Fatalf("after an empty scan the wizard is on state %d, want Add account", got.state)
	}
	got.wizardStep = 1
	if got = key(got, enter); got.state != StateAdd || got.form == nil {
		t.Errorf("'Add manually' should open the form (state=%d form=%v)", got.state, got.form != nil)
	}
	// Esc from the form returns to the add-account menu while setup is running.
	if got = key(got, tea.KeyMsg{Type: tea.KeyEsc}); got.state != StateWizardAddMore {
		t.Errorf("esc from the form during setup went to state %d, want Add account", got.state)
	}
}
