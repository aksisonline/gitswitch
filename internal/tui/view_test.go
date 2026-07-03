package tui

import (
	"strings"
	"testing"

	"github.com/aksisonline/gitswitch/internal/storage"
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

	for _, arcade := range []bool{false, true} {
		for _, size := range [][2]int{{66, 24}, {80, 24}, {90, 12}, {120, 40}} {
			for _, st := range states {
				for tab := 0; tab < 3; tab++ {
					m := Model{
						store:           storage.NewAt(t.TempDir()),
						profiles:        profiles,
						active:          &profiles[0],
						state:           st,
						width:           size[0],
						height:          size[1],
						arcadeMode:      arcade,
						tabIndex:        tab,
						latestVersion:   "v9.9.9",
						updateAvailable: true,
						releaseNotes:    "## What's new\n- something quite long that should be truncated to fit the panel width without wrapping",
						detectedProfiles: profiles,
						importSelected:   []bool{true, false},
						shellAlias:       "gs",
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
