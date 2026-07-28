package tui

import (
	"strings"

	"github.com/aksisonline/gitswitch/internal/git"
	"github.com/charmbracelet/lipgloss"
)

// tabNames returns the three tab labels for the current mode. The mouse
// handler hit-tests against these same strings, so view and clicks never drift.
func (m Model) tabNames() [3]string {
	if m.arcadeMode {
		return [3]string{"ACCOUNTS", "UTILITIES", "SETTINGS"}
	}
	return [3]string{"Accounts", "Utilities", "Settings"}
}

// viewTabHeader renders the 3-tab navigation strip (Accounts / Utilities / Settings).
// Called at the top of every tab view, after the main brand header.
func (m Model) viewTabHeader() string {
	tabs := m.tabNames()
	activeStyle := lipgloss.NewStyle().
		Background(colorBgChip).
		Foreground(colorPurple).
		Bold(true).
		Padding(0, 1)
	if m.arcadeMode {
		activeStyle = activeStyle.Foreground(colorYellow)
	}
	sep := "   "
	if m.arcadeMode {
		sep = " · "
	}
	var parts []string
	for i, t := range tabs {
		if i == m.tabIndex {
			parts = append(parts, activeStyle.Render(t))
		} else {
			parts = append(parts, styleItemDim.Render(t))
		}
	}
	return "\n\n  " + strings.Join(parts, styleItemDim.Render(sep))
}

// viewAccountsTab is the main account-switching screen (tab 0).
// It's the original viewList content with the tab header prepended.
func (m Model) viewAccountsTab(pw int) string {
	compact := m.height > 0 && m.height < 12+len(m.profiles)

	var top string
	if m.arcadeMode {
		top = m.viewScoreLine(pw) + "\n"
	}

	header := m.viewHeader("") + m.viewTabHeader()
	currentLine := m.viewCurrentLine(compact)

	nickColW := m.nickColumnWidth()
	items := m.viewProfileItems(pw, nickColW)
	statusLine := m.viewStatusLine(compact)

	var updateBanner string
	if m.updateAvailable {
		bannerSep := "\n\n"
		if compact {
			bannerSep = "\n"
		}
		if m.arcadeMode {
			chip := lipgloss.NewStyle().
				Background(colorBgChip).
				Foreground(colorYellow).
				Bold(true).
				Padding(0, 1).
				Render("★ BONUS STAGE")
			updateBanner = bannerSep + "  " + chip + "  " +
				styleScore.Render(m.latestVersion) +
				styleBrand.Render("  available")
		} else {
			chip := styleChipBox().Render("UPDATE")
			updateBanner = bannerSep + "  " + chip + "  " +
				styleCurrentVal.Render(m.latestVersion) +
				styleBrand.Render("  ·  press [u] to upgrade")
		}
	}

	footerPairs := [][2]string{
		{"↑/↓", "navigate"},
		{"enter", "switch"},
	}
	// Only offer the pin key when there is a repo to pin to.
	if m.repoKey != "" {
		label := "pin to repo"
		if m.scope == git.ScopeRepo {
			label = "pin/unpin repo"
		}
		footerPairs = append(footerPairs, [2]string{"p", label})
	}
	footerPairs = append(footerPairs,
		[2]string{"a", "add"},
		[2]string{"e", "edit"},
		[2]string{"v", secondaryToggleLabel(m.showUsername)},
		[2]string{"?", "cli tips"},
		[2]string{"q", "quit"},
	)
	footerPairs = append(footerPairs, [2]string{"1/2/3", "tabs"})
	if m.updateAvailable {
		footerPairs = append(footerPairs, [2]string{"u", "upgrade"})
	}

	footerSep := "\n\n"
	if compact {
		footerSep = "\n"
	}
	footer := footerSep + divider(pw) + "\n" + m.footerKeys(pw, footerPairs)
	return top + stylePanelBorder(pw).Render(header+currentLine+items+updateBanner+statusLine+footer)
}

// secondaryToggleLabel returns the footer hint for the email/username toggle.
func secondaryToggleLabel(showUsername bool) string {
	if showUsername {
		return "emails"
	}
	return "usernames"
}
