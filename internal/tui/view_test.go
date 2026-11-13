package tui

import (
	"strings"
	"testing"

	"github.com/aksisonline/gitswitch/internal/git"
	"github.com/aksisonline/gitswitch/internal/storage"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// testProfiles is a deliberately awkward set: a long email, a long nickname.
var testProfiles = []storage.Profile{
	{Nickname: "personal", UserName: "Abhiram Kanna", Email: "a-very-long-address@example-domain-name.com", GHUser: "aksisonline", SSHKey: "~/.ssh/id_p"},
	{Nickname: "work-account-long-nickname", UserName: "Work", Email: "work@corp.example"},
}

var allStates = []State{
	StateList, StateTips, StateDeleteConfirm, StateUpdatePrompt,
	StateWhatsNew, StateWizardWelcome, StateWizardDetect, StateWizardImport,
	StateWizardAddMore, StateWizardDone, StateShellConfirm,
	StateIntro, StateSelectFlash, StateTransition, StateExitAnim,
	StateAdd, StateEdit,
}

// screenModel builds a model in the widest realistic configuration: inside a
// repo pinned to a profile (extra tag on the Current line, extra footer key),
// an update pending, a status message showing.
func screenModel(t *testing.T, st State, tab, w, h int, arcade bool) Model {
	t.Helper()
	m := Model{
		store:            storage.NewAt(t.TempDir()),
		profiles:         testProfiles,
		active:           &testProfiles[0],
		state:            st,
		width:            w,
		height:           h,
		arcadeMode:       arcade,
		tabIndex:         tab,
		latestVersion:    "v9.9.9",
		updateAvailable:  true,
		releaseNotes:     "## What's new\n- something quite long that should be truncated to fit the panel width without wrapping",
		whatsNewBody:     "## What's new\n- something quite long that should be wrapped to fit the panel width without overflowing",
		currentVersion:   "v0.5.0-beta.3",
		detectedProfiles: testProfiles,
		importSelected:   []bool{true, false},
		shellAlias:       "gs",
		repoKey:          "git@github.com:acme/widget.git",
		scope:            git.ScopeRepo,
		scopeProfile:     &testProfiles[1],
		onboarding:       st >= StateWizardWelcome && st <= StateWizardDone,
		aliasInput:       textinput.New(),
		formFields:       [6]string{"personal", "Abhiram Kanna", "a@example.com", "", "", "aksisonline"},
	}
	if st == StateAdd || st == StateEdit {
		m.openProfileForm(st == StateEdit, m.formFields)
	}
	return m
}

// TestEveryScreenFitsTheTerminal renders every screen, on every tab, in both
// modes, from a roomy window down to the smallest supported one, and asserts
// the output never exceeds the terminal in either direction. A line wider than
// the window wraps and a body taller than it scrolls the top off-screen —
// both break the layout and the mouse grid.
func TestEveryScreenFitsTheTerminal(t *testing.T) {
	sizes := [][2]int{
		{minTermW, minTermH}, {50, 16}, {60, 20}, {66, 24}, {80, 24}, {90, 12}, {100, 30}, {120, 40}, {200, 60},
	}
	for _, arcade := range []bool{false, true} {
		for _, size := range sizes {
			for _, st := range allStates {
				for tab := 0; tab < 3; tab++ {
					m := screenModel(t, st, tab, size[0], size[1], arcade)
					m.statusMsg = "Session Isolation enabled — run: source /Users/someone/with/a/long/home/.zshrc"
					view := m.View()
					lines := strings.Split(view, "\n")
					if len(lines) > size[1] {
						t.Errorf("arcade=%v size=%v state=%d tab=%d: %d rows > terminal height %d",
							arcade, size, st, tab, len(lines), size[1])
					}
					for i, line := range lines {
						if w := lipgloss.Width(line); w > size[0] {
							t.Errorf("arcade=%v size=%v state=%d tab=%d line %d: width %d > terminal width %d\n%q",
								arcade, size, st, tab, i, w, size[0], ansi.Strip(line))
						}
					}
				}
			}
		}
	}
}

// TestTooSmallWindowShowsNotice: below the minimum size the panel is replaced
// by a short notice rather than drawn broken.
func TestTooSmallWindowShowsNotice(t *testing.T) {
	for _, size := range [][2]int{{minTermW - 1, 30}, {100, minTermH - 1}, {20, 5}} {
		m := screenModel(t, StateList, 0, size[0], size[1], false)
		view := ansi.Strip(m.View())
		if !strings.Contains(view, "too small") {
			t.Errorf("size=%v: expected the too-small notice, got:\n%s", size, view)
		}
		for _, line := range strings.Split(view, "\n") {
			if lipgloss.Width(line) > size[0] {
				t.Errorf("size=%v: notice line wider than the window: %q", size, line)
			}
		}
	}
}

// TestLayoutDoesNotShift pins the things that used to move the panel around:
// a status message appearing, a toggle flipping, the update notice arriving,
// and changing tab must all leave the panel's size and position unchanged.
func TestLayoutDoesNotShift(t *testing.T) {
	for _, size := range [][2]int{{80, 24}, {100, 40}, {60, 18}} {
		base := screenModel(t, StateList, 0, size[0], size[1], false)
		base.updateAvailable = false

		// The tab strip must sit on the same terminal row on every tab.
		tabRow := func(m Model) int {
			for i, l := range strings.Split(ansi.Strip(m.View()), "\n") {
				if strings.Contains(l, "Accounts") && strings.Contains(l, "Settings") {
					return i
				}
			}
			return -1
		}
		want := tabRow(base)
		for tab := 1; tab < 3; tab++ {
			m := base
			m.tabIndex = tab
			if got := tabRow(m); got != want {
				t.Errorf("size=%v: tab strip on row %d for tab %d, row %d for tab 0", size, got, tab, want)
			}
		}

		for tab := 0; tab < 3; tab++ {
			m := base
			m.tabIndex = tab
			quiet := m.screen()

			variants := map[string]func(*Model){
				"status message":  func(m *Model) { m.statusMsg = strings.Repeat("a long status message ", 8) },
				"error message":   func(m *Model) { m.statusMsg, m.statusIsErr = "something failed", true },
				"busy spinner":    func(m *Model) { m.busy = "switching to work…" },
				"update notice":   func(m *Model) { m.updateAvailable = true },
				"toggles flipped": func(m *Model) { m.shellEnabled, m.guardEnabled, m.ghWrapperEnabled = true, true, true },
				"toggle applying": func(m *Model) { m.busyItem = 2 },
				"theme changed":   func(m *Model) { m.colorTheme = 7 },
				"alias disabled":  func(m *Model) { m.shellAliasDisabled = true },
				"row glowing":     func(m *Model) { m.flashNick, m.flashLeft = "personal", 5 },
			}
			for name, change := range variants {
				v := m
				change(&v)
				got := v.screen()
				if got.height() != quiet.height() || got.width() != quiet.width() {
					t.Errorf("size=%v tab=%d %s: panel %dx%d, was %dx%d", size, tab, name,
						got.width(), got.height(), quiet.width(), quiet.height())
				}
			}
		}
	}
}

// TestWhatsNewFitsScreenAndScrolls verifies that long release notes (including
// canary boilerplate, install code blocks, and multi-sentence bullet points)
// never overflow the terminal height, strip redundant boilerplate, and clamp
// scroll bounds properly.
func TestWhatsNewFitsScreenAndScrolls(t *testing.T) {
	canaryBody := `> **⚠ This is a canary/beta release** built from ` + "`canary`" + `. It may contain rough edges.

### What's in this beta

### What's New
- A hidden ` + "`gitswitch snake`" + ` game — take a quick terminal break without needing Git or a configured profile.
- ` + "`gitswitch login`" + ` now hands off straight to ` + "`gh`" + `'s own login flow instead of a separate branded one — you get whatever login method ` + "`gh`" + ` offers (browser, device code, or a token) instead of gitswitch's own device-code screen. The ` + "`--client-id`" + ` flag is gone since there's no longer a custom OAuth app; ` + "`--host`" + ` and ` + "`--profile`" + ` still work the same way. gitswitch no longer stores a token anywhere itself — ` + "`gh`" + ` already does.
- The setup wizard now recognizes when a ` + "`gh`" + ` account and your git config are the same person — it checks your verified GitHub email and merges them into one profile instead of showing you two things to sort out.

### Under the Hood
- The Windows install command is shorter too — ` + "`irm https://get.gitswitch.dev/install.ps1 | iex`" + ` replaces the long ` + "`raw.githubusercontent.com/.../install.ps1`" + ` URL, matching the shortened curl command from the previous release. Behaves identically; get.gitswitch.dev proxies the same script.
- Housekeeping pass — removed dead code, a one-call-site dependency, and a release-workflow step that built and uploaded binaries nothing ever downloaded. No behavior change.

### Install this beta

**macOS / Linux**
` + "```bash\ncurl -fsSL https://get.gitswitch.dev | bash -s -- v0.5.0-beta.3\n```" + `

**Windows — PowerShell**
` + "```powershell\n$env:GS_VERSION=\"v0.5.0-beta.3\"; irm https://get.gitswitch.dev/install.ps1 | iex\n```" + `

## What's Changed
* feat: pre-commit identity guard + VS Code/JetBrains extensions [minor] by @aksisonline in https://github.com/aksisonline/gitswitch/pull/44


**Full Changelog**: https://github.com/aksisonline/gitswitch/compare/v0.5.0-beta.2...v0.5.0-beta.3`

	for _, termH := range []int{0, 18, 24, 30, 40} {
		m := Model{
			store:          storage.NewAt(t.TempDir()),
			state:          StateWhatsNew,
			width:          80,
			height:         termH,
			currentVersion: "v0.5.0-beta.3",
			whatsNewBody:   canaryBody,
		}
		applyTheme(normalThemes[0], false)
		rendered := m.bodyView()
		lines := strings.Split(rendered, "\n")

		if termH >= 18 && len(lines) > termH {
			t.Errorf("termH=%d: rendered height %d exceeds terminal height %d", termH, len(lines), termH)
		}
		plain := ansi.Strip(rendered)
		if strings.Contains(plain, "Install this beta") || strings.Contains(plain, "Full Changelog") || strings.Contains(plain, "gitswitch login") {
			t.Errorf("termH=%d: expected boilerplate and older release sections to be stripped, got:\n%s", termH, plain)
		}
		if !strings.Contains(plain, "What's New in v0.5.0-beta.3") {

			t.Errorf("termH=%d: top header missing from rendered output:\n%s", termH, plain)
		}

		// Scroll down many times — should clamp at maxScroll and never panic or blank out.
		var model tea.Model = m
		for i := 0; i < 50; i++ {
			model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
		}
		mScrolled := model.(Model)
		maxS := mScrolled.maxWhatsNewScroll()
		if mScrolled.whatsNewScroll != maxS {
			t.Errorf("termH=%d: scroll=%d, want clamped maxScroll=%d", termH, mScrolled.whatsNewScroll, maxS)
		}
		renderedBottom := mScrolled.bodyView()
		if !strings.Contains(ansi.Strip(renderedBottom), "What's New in v0.5.0-beta.3") {
			t.Errorf("termH=%d: top header scrolled off screen after scrolling down", termH)
		}
	}
}

// click sends a left press at terminal cell (x, y).
func click(m Model, x, y int) Model {
	next, _ := m.Update(tea.MouseMsg{X: x, Y: y, Action: tea.MouseActionPress, Button: tea.MouseButtonLeft})
	return next.(Model)
}

// rowOf returns the terminal column and row where text last appears (the
// last, because a profile's name also shows on the Current line above the list).
func rowOf(t *testing.T, m Model, text string) (x, y int) {
	t.Helper()
	x, y = -1, -1
	for i, l := range strings.Split(ansi.Strip(m.View()), "\n") {
		if j := strings.Index(l, text); j >= 0 {
			x, y = lipgloss.Width(l[:j]), i
		}
	}
	if y < 0 {
		t.Fatalf("%q not on screen:\n%s", text, ansi.Strip(m.View()))
	}
	return x, y
}

// TestMouseHitsWhatIsRendered clicks on what is actually drawn — at roomy,
// compact and tight densities, in both modes — and checks the right thing
// reacts. This is the drift that used to break clicks on the Utilities and
// Settings tabs whenever the layout changed.
func TestMouseHitsWhatIsRendered(t *testing.T) {
	for _, arcade := range []bool{false, true} {
		for _, size := range [][2]int{{100, 40}, {80, 24}, {60, 14}} {
			m := screenModel(t, StateList, 0, size[0], size[1], arcade)
			m.updateAvailable = false

			// Tabs: clicking a label selects that tab.
			for i, name := range m.tabNames() {
				x, y := rowOf(t, m, name)
				if got := click(m, x+1, y); got.tabIndex != i {
					t.Errorf("arcade=%v size=%v: clicking tab %q selected tab %d", arcade, size, name, got.tabIndex)
				}
			}

			// Accounts: hovering the second profile moves the cursor to it.
			x, y := rowOf(t, m, "work-acc")
			next, _ := m.Update(tea.MouseMsg{X: x, Y: y, Action: tea.MouseActionMotion})
			if got := next.(Model).cursor; got != 1 {
				t.Errorf("arcade=%v size=%v: hovering profile 2 put the cursor on %d", arcade, size, got)
			}

			// Utilities: clicking each item focuses that item.
			m.tabIndex = 1
			for i, spec := range utilItemSpecs {
				title := spec.title
				if arcade {
					title = spec.arcadeTitle
				}
				x, y := rowOf(t, m, title)
				next, _ := m.Update(tea.MouseMsg{X: x, Y: y, Action: tea.MouseActionMotion})
				if got := next.(Model).utilityFocus; got != i {
					t.Errorf("arcade=%v size=%v: hovering utility %q focused item %d", arcade, size, title, got)
				}
			}

			// Settings: the rename chip opens the alias editor.
			m.tabIndex = 2
			x, y = rowOf(t, m, chipRename)
			if got := click(m, x+2, y); !got.aliasEditing || got.settingsFocus != 2 {
				t.Errorf("arcade=%v size=%v: clicking %s did not start alias editing", arcade, size, chipRename)
			}
		}
	}
}
