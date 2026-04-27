package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.arcadeMode {
		applyTheme(arcadeTheme, true)
	} else {
		applyTheme(normalThemes[m.colorTheme], false)
	}

	pw := m.panelWidth()
	var body string
	switch m.state {
	case StateIntro:
		body = m.viewIntro(pw)
	case StateSelectFlash:
		body = m.viewSelectFlash(pw)
	case StateTransition:
		body = m.viewTransition(pw)
	case StateExitAnim:
		body = m.viewExitAnim(pw)
	case StateAdd:
		body = m.viewForm("Add Profile", pw)
	case StateEdit:
		body = m.viewForm(fmt.Sprintf("Edit  %s", m.editingNick), pw)
	case StateDeleteConfirm:
		body = m.viewDeleteConfirm(pw)
	case StateTips:
		body = m.viewTips(pw)
	default:
		body = m.viewList(pw)
	}
	if m.width > 0 && m.height > 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, body)
	}
	return body
}

func (m Model) viewHeader(subtitle string) string {
	var icon string
	if m.arcadeMode {
		icon = styleTitle.Render("ᗧ")
	} else {
		icon = styleTitle.Render("✦")
	}
	title := styleTitle.Render("Git-Switcher")

	var heading string
	if subtitle != "" {
		if m.arcadeMode {
			pellets := styleDivider.Render("· · ·")
			heading = fmt.Sprintf("  %s  %s %s%s%s\n", icon, title, pellets, styleItemDim.Render(" ► "), styleFormTitle.Render(subtitle))
		} else {
			heading = fmt.Sprintf("  %s  %s%s%s\n", icon, title, styleItemDim.Render(" › "), styleFormTitle.Render(subtitle))
		}
	} else {
		if m.arcadeMode {
			pellets := styleDivider.Render("· · ·")
			heading = fmt.Sprintf("  %s  %s %s\n", icon, title, pellets)
		} else {
			heading = fmt.Sprintf("  %s  %s\n", icon, title)
		}
	}
	made := styleBrand.Render("  Made by ") + styleTitle.Render("AKS") +
		styleBrand.Render("  ·  ") + styleBrandLink.Render("abhiramkanna.com")
	return heading + made
}

func (m Model) viewList(pw int) string {
	header := m.viewHeader("")
	currentLine := m.viewCurrentLine()

	nickColW := m.nickColumnWidth()
	items := m.viewProfileItems(pw, nickColW)
	statusLine := m.viewStatusLine()

	var updateBanner string
	if m.updateAvailable {
		if m.arcadeMode {
			updateBanner = "\n\n  " +
				lipgloss.NewStyle().Foreground(colorYellow).Bold(true).Render("★ BONUS STAGE: "+m.latestVersion+" available") +
				"\n  " + styleBrand.Render("Press [u] to upgrade")
		} else {
			updateBanner = "\n\n  " +
				lipgloss.NewStyle().Foreground(colorGreen).Bold(true).Render("⬆ Update available: "+m.latestVersion) +
				"\n  " + styleBrand.Render("Press [u] to upgrade")
		}
	}

	footerPairs := [][2]string{
		{"↑/↓", "navigate"},
		{"enter", "switch"},
		{"a", "add"},
		{"e", "edit"},
		{"?", "cli tips"},
		{"q", "quit"},
	}
	if !m.arcadeMode {
		footerPairs = append(footerPairs, [2]string{"c", "theme"})
	}
	if m.updateAvailable {
		footerPairs = append(footerPairs, [2]string{"u", "upgrade"})
	}

	footer := "\n\n" + divider(pw) + "\n" + m.footerKeys(pw, footerPairs)
	return stylePanelBorder(pw).Render(header + currentLine + items + updateBanner + statusLine + footer)
}

func (m Model) viewCurrentLine() string {
	if m.active != nil {
		tags := ""
		if m.active.SSHKey != "" {
			tags += "  " + styleItemDim.Render("ssh")
		}
		if m.active.GHUser != "" {
			tags += "  " + styleItemDim.Render("gh:"+m.active.GHUser)
		}
		label := "Current  "
		if m.arcadeMode {
			label = "PLAYER 1 ► "
		}
		return "\n\n  " +
			styleCurrent.Render(label) +
			styleCheckmark.Render(m.active.Nickname) +
			styleCurrent.Render("  ·  ") +
			styleCurrentVal.Render(m.active.Email) +
			tags
	}
	if m.arcadeMode {
		return "\n\n  " + styleCurrent.Render("NO ACTIVE PLAYER")
	}
	return "\n\n  " + styleCurrent.Render("No active profile")
}

func (m Model) nickColumnWidth() int {
	nickColW := 12
	for _, p := range m.profiles {
		nickW := lipgloss.Width(p.Nickname)
		if nickW > nickColW {
			nickColW = nickW
		}
	}
	return nickColW + 2
}

func (m Model) viewProfileItems(pw, nickColW int) string {
	if len(m.profiles) == 0 {
		if m.arcadeMode {
			return "\n\n  " + styleItemDim.Render("NO PLAYERS YET  ·  PRESS [A] TO INSERT COIN")
		}
		return "\n\n  " + styleItemDim.Render("No profiles yet. Press [a] to add one.")
	}
	cursorChar := "  ❯"
	checkChar := func(active bool) string {
		if active {
			return " " + styleCheckmark.Render("✓") + " "
		}
		return styleItemDim.Render(" · ")
	}
	if m.arcadeMode {
		cursorChar = "  ᗧ"
		checkChar = func(active bool) string {
			if active {
				return " " + styleCheckmark.Render("★") + " "
			}
			return styleItemDim.Render(" · ")
		}
	}
	items := "\n"
	for i, p := range m.profiles {
		cursor := "   "
		if i == m.cursor {
			cursor = cursorChar
		}
		isActive := m.active != nil && p.Nickname == m.active.Nickname
		check := checkChar(isActive)
		nick := p.Nickname + strings.Repeat(" ", max(0, nickColW-lipgloss.Width(p.Nickname)))
		line := fmt.Sprintf("%s%s%s  %s", cursor, check, nick, p.Email)
		if i == m.cursor {
			line = styleItemActive(pw).Render(line)
		} else {
			line = "  " + styleItemInactive.Render(line)
		}
		items += "\n" + line
	}
	return items
}

func (m Model) viewStatusLine() string {
	if m.statusMsg == "" {
		return ""
	}
	if m.statusIsErr {
		prefix := ""
		if m.arcadeMode {
			prefix = "✗ "
		}
		return "\n\n  " + lipgloss.NewStyle().Foreground(colorRed).Render(prefix+m.statusMsg)
	}
	if m.arcadeMode {
		return "\n\n  " + styleBrand.Render("◎ "+m.statusMsg)
	}
	return "\n\n  " + styleBrand.Render(m.statusMsg)
}

func (m Model) viewForm(subtitle string, pw int) string {
	inputW := pw - 6
	header := m.viewHeader(subtitle)

	var fields string
	for i := range formLabels {
		fields += "\n\n"
		sub := styleBrand.Render("  " + formSubtitles[i])
		if i == m.formFocus {
			fields += "  " + styleInputLabelActive.Render(formLabels[i]) + sub + "\n"
			fields += "  " + styleInputActive(inputW).Render(m.formFields[i]+styleTitle.Render("█"))
		} else {
			fields += "  " + styleInputLabel.Render(formLabels[i]) + sub + "\n"
			v := m.formFields[i]
			if v == "" {
				v = " "
			}
			fields += "  " + styleInputInactive(inputW).Render(v)
		}
	}

	footerPairs := [][2]string{
		{"tab/↓", "next"},
		{"shift+tab/↑", "prev"},
		{"enter", "next / submit"},
		{"esc", "cancel"},
	}
	if m.state == StateEdit {
		footerPairs = append(footerPairs, [2]string{"ctrl+d", "delete"})
	}

	footer := "\n\n" + divider(pw) + "\n" + m.footerKeys(pw, footerPairs)
	return stylePanelBorder(pw).Render(header + fields + footer)
}

func (m Model) viewDeleteConfirm(pw int) string {
	header := m.viewHeader("Delete Profile")
	var body string
	if len(m.profiles) > 0 {
		p := m.profiles[m.cursor]
		var titleText, bodyText string
		if m.arcadeMode {
			titleText = fmt.Sprintf("GAME OVER?  '%s'", p.Nickname)
			bodyText = "THIS ACTION CANNOT BE UNDONE."
		} else {
			titleText = fmt.Sprintf("Delete '%s'?", p.Nickname)
			bodyText = "This cannot be undone."
		}
		body = "\n\n  " + styleDeleteTitle.Render(titleText) +
			"\n  " + styleCurrent.Render(fmt.Sprintf("%s  ·  %s", p.UserName, p.Email)) +
			"\n\n  " + styleBrand.Render(bodyText)
	}
	footer := "\n\n" + divider(pw) + "\n" + m.footerKeys(pw, [][2]string{
		{"y", "confirm delete"},
		{"n / esc", "cancel"},
	})
	return stylePanelBorder(pw).Render(header + body + footer)
}

func (m Model) viewTips(pw int) string {
	header := m.viewHeader("CLI Quick Reference")

	type tip struct{ cmd, desc string }
	tips := []tip{
		{"gitswitch <nickname>", "quick switch without opening UI"},
		{"gitswitch current", "show the active profile"},
		{"gitswitch list", "list all saved profiles"},
		{"gitswitch add <nick> <user-name> <email>", "add a profile from the terminal"},
		{"gitswitch switch <nick>", "switch profile from the terminal"},
		{"gitswitch remove <nick>", "delete a profile"},
	}

	cmdStyle := lipgloss.NewStyle().Foreground(colorGreen)
	descStyle := lipgloss.NewStyle().Foreground(colorDim)
	labelStyle := lipgloss.NewStyle().Foreground(colorPurple).Bold(true)

	sectionLabel := "Skip the UI entirely:"
	if m.arcadeMode {
		sectionLabel = "CHEAT CODES:"
	}

	body := "\n\n  " + labelStyle.Render(sectionLabel) + "\n"
	for _, t := range tips {
		body += fmt.Sprintf("\n  %s\n  %s\n", cmdStyle.Render(t.cmd), descStyle.Render("  "+t.desc))
	}

	footer := "\n" + divider(pw) + "\n" + m.footerKeys(pw, [][2]string{{"esc / ?", "back"}})
	return stylePanelBorder(pw).Render(header + body + footer)
}

// ── Arcade-only animation views ───────────────────────────────────────────────

func (m Model) viewIntro(pw int) string {
	icon := styleTitle.Render("ᗧ")
	title := styleTitle.Render("GIT-SWITCHER")
	pellets := styleDivider.Render("· · ·")
	titleRow := fmt.Sprintf("  %s  %s %s", icon, title, pellets)
	subtitleRow := styleBrand.Render("     ARCADE MODE")

	track := "  " + m.renderIntroTrack(pw)

	var readyRow string
	if m.introPhase == 1 {
		readyStyle := lipgloss.NewStyle().Foreground(colorPurple).Bold(true)
		if m.introReadyFrame%2 == 1 {
			readyStyle = lipgloss.NewStyle().Foreground(colorDim)
		}
		readyRow = "\n\n  " + readyStyle.Render("★   PLAYER 1   READY!   ★")
	}

	skipHint := styleBrand.Render("  [any key to skip]")
	body := titleRow + "\n" + subtitleRow + "\n\n\n" + track + readyRow + "\n\n" + skipHint + "\n"
	return stylePanelBorder(pw).Render(body)
}

func (m Model) renderIntroTrack(pw int) string {
	numSlots := (pw - 2) / 2
	if numSlots < 1 {
		numSlots = 1
	}
	pacSlot := m.introPos
	if m.introPhase == 1 {
		pacSlot = numSlots - 1
	}
	ghost1 := pacSlot - 2
	ghost2 := pacSlot - 4
	ghost3 := pacSlot - 6

	g1 := lipgloss.NewStyle().Foreground(arcadeGhostRed)
	g2 := lipgloss.NewStyle().Foreground(arcadeGhostPink)
	g3 := lipgloss.NewStyle().Foreground(arcadeGhostCyan)

	track := ""
	for slot := 0; slot < numSlots; slot++ {
		var ch string
		switch {
		case slot == ghost3 && ghost3 >= 0:
			ch = g3.Render("ᗣ") + " "
		case slot == ghost2 && ghost2 >= 0:
			ch = g2.Render("ᗣ") + " "
		case slot == ghost1 && ghost1 >= 0:
			ch = g1.Render("ᗣ") + " "
		case slot == pacSlot:
			if m.introMouthOpen {
				ch = styleTitle.Render("ᗧ") + " "
			} else {
				ch = styleTitle.Render("●") + " "
			}
		case slot > pacSlot:
			ch = styleDivider.Render("·") + " "
		default:
			ch = "  "
		}
		track += ch
	}
	return track
}

func (m Model) viewSelectFlash(pw int) string {
	header := m.viewHeader("")
	currentLine := m.viewCurrentLine()
	nickColW := m.nickColumnWidth()

	items := "\n"
	for i, p := range m.profiles {
		if i == m.selFlashProfile {
			nick := p.Nickname + strings.Repeat(" ", max(0, nickColW-lipgloss.Width(p.Nickname)))
			line := fmt.Sprintf("  ᗧ · %s  %s", nick, p.Email)
			var flashStyle lipgloss.Style
			if m.selFlashFrame%2 == 0 {
				flashStyle = lipgloss.NewStyle().Background(colorYellow).Foreground(lipgloss.Color("0")).Width(pw).Bold(true)
			} else {
				flashStyle = lipgloss.NewStyle().Background(lipgloss.Color("255")).Foreground(lipgloss.Color("0")).Width(pw)
			}
			items += "\n" + flashStyle.Render(line)
		} else {
			isActive := m.active != nil && p.Nickname == m.active.Nickname
			check := styleItemDim.Render(" · ")
			if isActive {
				check = " " + styleCheckmark.Render("★") + " "
			}
			nick := p.Nickname + strings.Repeat(" ", max(0, nickColW-lipgloss.Width(p.Nickname)))
			line := fmt.Sprintf("   %s%s  %s", check, nick, p.Email)
			items += "\n  " + styleItemInactive.Render(line)
		}
	}

	footer := "\n\n" + divider(pw) + "\n" + m.footerKeys(pw, [][2]string{})
	return stylePanelBorder(pw).Render(header + currentLine + items + footer)
}

func (m Model) viewTransition(pw int) string {
	header := m.viewHeader("")

	numSlots := (pw - 2) / 2
	if numSlots < 1 {
		numSlots = 1
	}
	// Ghosts move right-to-left: start at far right, sweep left over 5 frames.
	step := numSlots / 5
	if step < 1 {
		step = 1
	}
	ghost1Slot := (numSlots - 1) - m.transFrame*step
	ghost2Slot := ghost1Slot + 2
	ghost3Slot := ghost1Slot + 4

	g1 := lipgloss.NewStyle().Foreground(arcadeGhostRed)
	g2 := lipgloss.NewStyle().Foreground(arcadeGhostPink)
	g3 := lipgloss.NewStyle().Foreground(arcadeGhostCyan)

	track := "  "
	for slot := 0; slot < numSlots; slot++ {
		var ch string
		switch slot {
		case ghost1Slot:
			ch = g1.Render("ᗣ") + " "
		case ghost2Slot:
			ch = g2.Render("ᗣ") + " "
		case ghost3Slot:
			ch = g3.Render("ᗣ") + " "
		default:
			ch = styleDivider.Render("·") + " "
		}
		track += ch
	}

	body := header + "\n\n\n" + track + "\n\n"
	return stylePanelBorder(pw).Render(body)
}

func (m Model) viewExitAnim(pw int) string {
	textStyle := lipgloss.NewStyle().Bold(true).Foreground(arcadeGhostRed)
	if m.exitFrame%2 == 1 {
		textStyle = lipgloss.NewStyle().Foreground(colorDim)
	}
	text := textStyle.Render("★   G A M E   O V E R   ★")
	sub := styleBrand.Render("  Thanks for playing")
	body := "\n\n\n\n  " + text + "\n\n  " + sub + "\n\n\n\n"
	return stylePanelBorder(pw).Render(body)
}

func (m Model) footerKeys(pw int, pairs [][2]string) string {
	var prefix string
	var sep string
	var sepW int
	if m.arcadeMode {
		prefix = lipgloss.NewStyle().Foreground(arcadeMazeBlue).Bold(true).Render("CONTROLS: ")
		sep = styleFooter.Render("  │  ")
		sepW = 5
	} else {
		prefix = "  "
		sep = styleFooter.Render("  ·  ")
		sepW = 5
	}
	var lines []string
	currentLine := "  " + prefix
	currentW := 2 + lipgloss.Width(prefix)
	for i, p := range pairs {
		item := styleFooterKey.Render(p[0]) + styleFooter.Render(" "+p[1])
		itemW := lipgloss.Width(p[0]) + 1 + lipgloss.Width(p[1])
		if i > 0 {
			if currentW+sepW+itemW > pw {
				lines = append(lines, currentLine)
				currentLine = "  " + item
				currentW = 2 + itemW
				continue
			}
			currentLine += sep + item
			currentW += sepW + itemW
		} else {
			currentLine += item
			currentW += itemW
		}
	}
	if currentLine != "  "+prefix && currentLine != "  " {
		lines = append(lines, currentLine)
	}
	return strings.Join(lines, "\n")
}

// Ensure Model satisfies tea.Model at compile time.
var _ tea.Model = Model{}
