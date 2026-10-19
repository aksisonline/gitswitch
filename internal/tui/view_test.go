package tui

import (
	"strings"
	"testing"

	"github.com/aksisonline/gitswitch/internal/git"
	"github.com/aksisonline/gitswitch/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// TestNoLineOverflowsPanel renders every screen at several terminal sizes and
// asserts no rendered line is wider than the panel's outer width. A wider line
// wraps in the terminal, which breaks both the layout and the mouse hit grid.
func TestNoLineOverflowsPanel(t *testing.T) {
	profiles := []storage.Profile{
		{Nickname: "personal", UserName: "Abhiram Kanna", Email: "a-very-long-address@example-domain-name.com", GHUser: "aksisonline", SSHKey: "~/.ssh/id_p"},
		{Nickname: "work-account-long-nickname", UserName: "Work", Email: "work@corp.example"},
	}

	states := []State{
		StateList, StateTips, StateDeleteConfirm, StateUpdatePrompt,
		StateWhatsNew, StateWizardWelcome, StateWizardDetect, StateWizardImport,
		StateWizardAddMore, StateWizardDone, StateShellConfirm,
		StateIntro, StateSelectFlash, StateTransition, StateExitAnim,
		StateAdd, StateEdit,
	}

	// Widest case: inside a repo pinned to a profile, so the extra "pinned" tag on
	// the Current line and the extra `p` footer key are both in play.
	for _, arcade := range []bool{false, true} {
		for _, size := range [][2]int{{66, 24}, {80, 24}, {90, 12}, {120, 40}} {
			for _, st := range states {
				for tab := 0; tab < 3; tab++ {
					m := Model{
						store:            storage.NewAt(t.TempDir()),
						profiles:         profiles,
						active:           &profiles[0],
						state:            st,
						width:            size[0],
						height:           size[1],
						arcadeMode:       arcade,
						tabIndex:         tab,
						latestVersion:    "v9.9.9",
						updateAvailable:  true,
						releaseNotes:     "## What's new\n- something quite long that should be truncated to fit the panel width without wrapping",
						whatsNewBody:     "## What's new\n- something quite long that should be wrapped to fit the panel width without overflowing",
						currentVersion:   "v0.5.0-beta.3",
						detectedProfiles: profiles,
						importSelected:   []bool{true, false},
						shellAlias:       "gs",
						repoKey:          "git@github.com:acme/widget.git",
						scope:            git.ScopeRepo,
						scopeProfile:     &profiles[1],
					}
					applyTheme(normalThemes[0], arcade)
					pw := m.panelWidth()
					body := m.bodyView(pw)
					for i, line := range strings.Split(body, "\n") {
						if w := lipgloss.Width(line); w > pw+2 {
							t.Errorf("arcade=%v size=%v state=%d tab=%d line %d: width %d > panel %d\n%q",
								arcade, size, st, tab, i, w, pw+2, line)
						}
					}
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
			store:                        storage.NewAt(t.TempDir()),
			state:                        StateWhatsNew,
			width:                        80,
			height:                       termH,
			currentVersion: "v0.5.0-beta.3",
			whatsNewBody:                 canaryBody,
		}
		applyTheme(normalThemes[0], false)
		pw := m.panelWidth()
		rendered := m.bodyView(pw)
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
		renderedBottom := mScrolled.bodyView(pw)
		if !strings.Contains(ansi.Strip(renderedBottom), "What's New in v0.5.0-beta.3") {
			t.Errorf("termH=%d: top header scrolled off screen after scrolling down", termH)
		}
	}
}

// TestMouseGridMatchesRender pins the mouse handler's row math to the actual
// rendered layout — the exact drift that broke clicks on the Utilities and
// Settings tabs.
func TestMouseGridMatchesRender(t *testing.T) {
	profiles := []storage.Profile{
		{Nickname: "personal", UserName: "A", Email: "a@example.com"},
		{Nickname: "work-account-long-nickname", UserName: "W", Email: "w@example.com"},
	}
	for _, arcade := range []bool{false, true} {
		for _, height := range []int{40, 13} { // roomy and compact
			m := Model{
				store:      storage.NewAt(t.TempDir()),
				profiles:   profiles,
				active:     &profiles[0],
				state:      StateList,
				width:      100,
				height:     height,
				arcadeMode: arcade,
				shellAlias: "gs",
			}
			applyTheme(normalThemes[0], arcade)
			off := 0
			if arcade {
				off = 1
			}

			// Accounts tab: second profile row must sit at profileRowStart+1.
			lines := strings.Split(m.bodyView(m.panelWidth()), "\n")
			found := -1
			for i, l := range lines {
				if strings.Contains(ansi.Strip(l), "work-account-long-nic") {
					found = i
					break
				}
			}
			if want := m.profileRowStart(off) + 1; found != want {
				t.Errorf("arcade=%v height=%d: profile row rendered at %d, mouse math expects %d", arcade, height, found, want)
			}

			// Utilities tab: first item box top border must sit where itemBoxAt starts.
			m.tabIndex = 1
			lines = strings.Split(m.bodyView(m.panelWidth()), "\n")
			found = -1
			for i, l := range lines {
				if strings.Contains(l, "┌") {
					found = i
					break
				}
			}
			if want := 4 + off; found != want {
				t.Errorf("arcade=%v height=%d: first item box at %d, mouse math expects %d", arcade, height, found, want)
			}
		}
	}
}
