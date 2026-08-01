package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Utility items (indices match utilityFocus and the 5-relY mouse grid):
//   0 = Shell Integration       (toggleable)
//   1 = Session Isolation       (toggleable)
//   2 = HTTPS Credential Helper (toggleable)
//   3 = Auto-pin New Repos      (toggleable)

func (m Model) viewUtilitiesTab(pw int) string {
	var top string
	if m.arcadeMode {
		top = m.viewScoreLine(pw) + "\n"
	}

	header := m.viewHeader("") + m.viewTabHeader()
	iw := itemInnerW(pw)

	// Build each item box
	items := m.utilItem(pw, iw, 0) + m.utilItem(pw, iw, 1) + m.utilItem(pw, iw, 2) + m.utilItem(pw, iw, 3)

	sep := "\n\n"
	footer := sep + divider(pw) + "\n" + m.footerKeys(pw, [][2]string{
		{"enter", "toggle"},
		{"1/2/3", "tabs"},
		{"q", "quit"},
	})

	return top + stylePanelBorder(pw).Render(header+items+footer)
}

func (m Model) utilItem(pw, iw, idx int) string {
	focused := m.utilityFocus == idx

	switch idx {
	case 0: // Shell Integration
		title := "Shell Integration"
		desc := "Auto-switch identity when you cd into a pinned repo."
		if m.arcadeMode {
			title = "SHELL HOOK"
			desc = "Auto-switch on cd. Works like a cheat code for repos."
		}
		toggle := renderToggle(m.shellEnabled)
		line1 := titleWithRight(styleCurrentVal.Render(title), toggle, iw)
		line2 := lipgloss.NewStyle().Foreground(colorDim).Render(truncate(desc, iw))
		line2 = padTo(line2, iw)
		return renderItemBox(pw, focused, false, line1, line2)

	case 1: // Session Isolation
		title := "Session Isolation"
		desc := "Per-repo git identity + gh account — needed for pins to work."
		if m.arcadeMode {
			title = "GH CO-OP MODE"
			desc = "No more fighting over one shared gh account. Every terminal plays its own."
		}
		toggle := renderToggle(m.ghWrapperEnabled)
		line1 := titleWithRight(styleCurrentVal.Render(title), toggle, iw)
		line2 := lipgloss.NewStyle().Foreground(colorDim).Render(truncate(desc, iw))
		line2 = padTo(line2, iw)
		return renderItemBox(pw, focused, false, line1, line2)

	case 2: // HTTPS Credential Helper
		title := "HTTPS Credential Helper"
		desc := "Route HTTPS git pushes through the active profile's gh account."
		if m.arcadeMode {
			title = "CREDENTIAL HELPER"
			desc = "HTTPS auth. Automatic. No 401 game-overs."
		}
		toggle := renderToggle(m.credentialHelperEnabled)
		line1 := titleWithRight(styleCurrentVal.Render(title), toggle, iw)
		line2 := lipgloss.NewStyle().Foreground(colorDim).Render(truncate(desc, iw))
		line2 = padTo(line2, iw)
		return renderItemBox(pw, focused, false, line1, line2)

	case 3: // Auto-pin New Repos
		title := "Auto-pin New Repos"
		desc := "First time in a repo, pin the active account to it."
		if m.arcadeMode {
			title = "AUTO-CLAIM"
			desc = "First visit to a repo claims it for the active player."
		}
		toggle := renderToggle(!m.autoPinDisabled)
		line1 := titleWithRight(styleCurrentVal.Render(title), toggle, iw)
		line2 := lipgloss.NewStyle().Foreground(colorDim).Render(truncate(desc, iw))
		line2 = padTo(line2, iw)
		return renderItemBox(pw, focused, false, line1, line2)
	}
	return ""
}
