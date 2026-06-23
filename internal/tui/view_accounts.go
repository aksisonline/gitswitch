package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// viewTabHeader renders the 3-tab navigation strip (Accounts / Utilities / Settings).
// Called at the top of every tab view, after the main brand header.
func (m Model) viewTabHeader() string {
	var tabs [3]string
	if m.arcadeMode {
		tabs = [3]string{"ACCOUNTS", "UTILITIES", "SETTINGS"}
	} else {
		tabs = [3]string{"Accounts", "Utilities", "Settings"}
	}
	var parts []string
	for i, t := range tabs {
		if i == m.tabIndex {
			parts = append(parts, styleTitle.Render("[ "+t+" ]"))
		} else {
			parts = append(parts, styleItemDim.Render(t))
		}
	}
	return "\n\n  " + strings.Join(parts, styleItemDim.Render("   "))
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
		{"a", "add"},
		{"e", "edit"},
		{"v", secondaryToggleLabel(m.showUsername)},
		{"?", "cli tips"},
		{"q", "quit"},
	}
	if !m.arcadeMode {
		footerPairs = append(footerPairs, [2]string{"1/2/3", "tabs"})
	}
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

// ensure fmt is used (for Sprintf in callers that might need it)
var _ = fmt.Sprintf
