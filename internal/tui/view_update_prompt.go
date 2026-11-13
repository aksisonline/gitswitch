package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Button labels on the update prompt — also the mouse targets.
const (
	updateYes       = "Update now"
	updateNo        = "Not now"
	updateYesArcade = "INSERT COIN"
	updateNoArcade  = "CONTINUE"
)

// drawUpdatePrompt confirms an upgrade and says exactly how it will happen —
// through Homebrew when Homebrew owns the binary, otherwise the install script.
func (m Model) drawUpdatePrompt(c *canvas) {
	m.header(c, "Update", "")
	c.gap()

	from := styleItemDim.Render(m.currentVersion)
	to := lipgloss.NewStyle().Foreground(colorGreen).Bold(true).Render(m.latestVersion)
	if m.arcadeMode {
		c.add("  " + styleChipBonus.Render("★ BONUS STAGE") + "  " + styleScore.Render("NEW VERSION DETECTED"))
	} else {
		c.add("  " + styleChipBox().Render("UPDATE") + "  " + styleCurrentVal.Render("New version available"))
	}
	c.add("  " + from + styleBrand.Render("  →  ") + to)
	c.gap()

	// How it updates — never a surprise.
	via := "Updates with " + m.install.Label() + ": "
	c.add("  " + styleBrand.Render(via) + styleFooterKey.Render(m.install.UpgradeHint()))
	if c.d != dTight {
		c.add("  " + styleBrand.Render("Your profiles and settings are not touched."))
	}

	// Release notes, capped to what the window can spare.
	if notes := m.updateNoteLines(c.w - 4); len(notes) > 0 && c.d != dTight {
		budget := 6
		if c.d == dCompact {
			budget = 3
		}
		c.gap()
		c.add("  " + styleItemDim.Render("What's new"))
		for i, line := range notes {
			if i >= budget {
				c.add("  " + styleItemDim.Render("  …"))
				break
			}
			c.add("  " + lipgloss.NewStyle().Foreground(colorPurple).Render("•") + " " + styleBrand.Render(line))
		}
	}

	c.gap()
	if m.arcadeMode {
		c.add(buttonRow(updateYesArcade, updateNoArcade))
	} else {
		c.add(buttonRow(updateYes, updateNo))
	}
	c.gap()
	c.footer(m, [][2]string{{"y / enter", "update"}, {"n / esc", "not now"}})
}

// updateNoteLines flattens the pending release's notes into short plain lines.
func (m Model) updateNoteLines(w int) []string {
	var out []string
	for _, line := range strings.Split(cleanReleaseNotes(m.releaseNotes), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimSpace(strings.TrimLeft(line, "*- "))
		line = strings.ReplaceAll(strings.ReplaceAll(line, "**", ""), "`", "")
		if idx := strings.Index(line, " by @"); idx > 0 && strings.Contains(line[idx:], "http") {
			line = strings.TrimSpace(line[:idx])
		}
		if line != "" {
			out = append(out, truncate(line, w))
		}
	}
	return out
}
