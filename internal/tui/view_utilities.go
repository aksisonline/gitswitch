package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Utility items (indices match utilityFocus and the 5-relY mouse grid):
//   0 = Shell Integration   (toggleable)
//   1 = Pre-commit Safety   (disabled, coming soon)
//   2 = Credential Helper   (disabled, coming soon)

func (m Model) viewUtilitiesTab(pw int) string {
	var top string
	if m.arcadeMode {
		top = m.viewScoreLine(pw) + "\n"
	}

	header := m.viewHeader("") + m.viewTabHeader()
	iw := itemInnerW(pw)

	// Build each item box
	items := m.utilItem(pw, iw, 0) + m.utilItem(pw, iw, 1) + m.utilItem(pw, iw, 2)

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

	case 1: // Pre-commit Safety Net
		title := "Pre-commit Safety Net"
		desc := "Warn before committing as the wrong identity in a pinned repo."
		chip := "v0.2.1"
		if m.arcadeMode {
			title = "PRE-COMMIT SAFETY NET"
			desc = "Stops wrong-player commits. Bonus stage."
			chip = "SOON"
		}
		chipStr := styleChipBox().Render(chip)
		line1 := titleWithRight(styleItemDim.Render(title), chipStr, iw)
		line2 := lipgloss.NewStyle().Foreground(colorDim).Render(truncate(desc, iw))
		line2 = padTo(line2, iw)
		return renderItemBox(pw, focused, true, line1, line2)

	case 2: // HTTPS Credential Helper
		title := "HTTPS Credential Helper"
		desc := "Manage GitHub HTTPS credentials automatically. No more 401s."
		chip := "v0.2.1"
		if m.arcadeMode {
			title = "CREDENTIAL HELPER"
			desc = "HTTPS auth. Automatic. No 401 game-overs."
			chip = "SOON"
		}
		chipStr := styleChipBox().Render(chip)
		line1 := titleWithRight(styleItemDim.Render(title), chipStr, iw)
		line2 := lipgloss.NewStyle().Foreground(colorDim).Render(truncate(desc, iw))
		line2 = padTo(line2, iw)
		return renderItemBox(pw, focused, true, line1, line2)
	}
	return ""
}
