package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) viewUpdatePrompt(pw int) string {
	header := m.viewHeader("") + m.viewTabHeader()
	iw := itemInnerW(pw)
	sep := "\n\n"

	var chip, title, body string
	if m.arcadeMode {
		chip = lipgloss.NewStyle().
			Background(colorBgChip).
			Foreground(colorYellow).
			Bold(true).
			Padding(0, 1).
			Render("★ BONUS STAGE")
		title = sep + "  " + chip + "  " + styleScore.Render("NEW VERSION DETECTED")
		body = sep +
			"  " + styleCurrentVal.Render(m.latestVersion) + styleBrand.Render("  is available") +
			sep +
			"  " + styleBrand.Render("  "+truncate("Update now to unlock new features and fixes.", iw-4))
	} else {
		chip = styleChipBox().Render("UPDATE")
		title = sep + "  " + chip + "  " + styleCurrentVal.Render("New version available")
		body = sep +
			"  " + styleBrand.Render("  "+m.latestVersion+" is ready to install.")

		// Release notes block (cleaned & capped so the prompt never overflows)
		if cleaned := cleanReleaseNotes(m.releaseNotes); cleaned != "" {
			body += sep + "  " + styleItemDim.Render("What's new:") + "\n"
			noteStyle := lipgloss.NewStyle().Foreground(colorDim)
			maxLines := 6
			shown := 0
			for _, line := range strings.Split(cleaned, "\n") {
				line = strings.TrimSpace(strings.TrimLeft(strings.TrimSpace(line), "#"))
				if line == "" {
					continue
				}
				line = strings.ReplaceAll(line, "**", "")
				if idx := strings.Index(line, " by @"); idx > 0 && strings.Contains(line[idx:], "http") {
					line = strings.TrimSpace(line[:idx])
				}
				if shown >= maxLines {
					body += "  " + noteStyle.Render("  …") + "\n"
					break
				}
				body += "  " + noteStyle.Render(truncate(line, iw-4)) + "\n"
				shown++
			}
		} else {
			body += sep +
				"  " + styleBrand.Render("  "+truncate("Your profiles are unchanged. Restart after upgrade to apply.", iw-4))
		}
	}

	yesLabel := "  Yes, update now  "
	noLabel := "  Skip  "
	if m.arcadeMode {
		yesLabel = "  INSERT COIN  "
		noLabel = "  CONTINUE  "
	}

	yesBtn := lipgloss.NewStyle().
		Background(colorGreen).
		Foreground(lipgloss.Color("0")).
		Bold(true).
		Padding(0, 1).
		Render(yesLabel)
	noBtn := lipgloss.NewStyle().
		Foreground(colorDim).
		Padding(0, 1).
		Render(noLabel)

	buttons := sep + "  " + yesBtn + "   " + noBtn

	footer := sep + divider(pw) + "\n" + m.footerKeys(pw, [][2]string{
		{"y / enter", "update"},
		{"n / esc", "skip"},
	})

	return stylePanelBorder(pw).Render(header + title + body + buttons + footer)
}
