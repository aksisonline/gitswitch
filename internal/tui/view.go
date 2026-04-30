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
		theme := m.colorTheme
		if theme < 0 || theme >= len(normalThemes) {
			theme = 0
		}
		applyTheme(normalThemes[theme], false)
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

// viewScoreLine — arcade-only: top-right score panel like classic arcades.
func (m Model) viewScoreLine(pw int) string {
	left := styleScoreLabel.Render("1UP") + "  " + styleScore.Render(fmt.Sprintf("%06d", m.score))
	right := styleScoreLabel.Render("HIGH SCORE") + "  " + styleScore.Render(fmt.Sprintf("%06d", m.hiScore))
	gap := pw - lipgloss.Width(left) - lipgloss.Width(right) - 4
	if gap < 1 {
		gap = 1
	}
	return "  " + left + strings.Repeat(" ", gap) + right
}

func (m Model) viewHeader(subtitle string) string {
	if m.arcadeMode {
		return m.viewArcadeHeader(subtitle)
	}
	return m.viewNormalHeader(subtitle)
}

func (m Model) viewNormalHeader(subtitle string) string {
	icon := styleTitle.Render("◆")
	title := styleTitle.Render("Git-Switcher")
	tagline := styleBrand.Render("identity manager for git")

	var heading string
	if subtitle != "" {
		heading = fmt.Sprintf("  %s  %s%s%s\n",
			icon, title,
			styleItemDim.Render("  ›  "),
			styleFormTitle.Render(subtitle))
	} else {
		heading = fmt.Sprintf("  %s  %s   %s", icon, title, tagline)
	}
	return heading
}

func (m Model) viewArcadeHeader(subtitle string) string {
	icon := styleTitle.Render("ᗧ")
	title := styleTitle.Render("GIT-SWITCHER")
	pellets := styleDivider.Render("· · ·")

	var heading string
	if subtitle != "" {
		heading = fmt.Sprintf("  %s  %s %s%s%s\n",
			icon, title, pellets,
			styleItemDim.Render(" ► "),
			styleFormTitle.Render(subtitle))
	} else {
		heading = fmt.Sprintf("  %s  %s %s", icon, title, pellets)
	}
	return heading
}

func (m Model) viewList(pw int) string {
	// Full list height = 10 + len(profiles) lines; switch to compact 2 lines early.
	compact := m.height > 0 && m.height < 12+len(m.profiles)

	var top string
	if m.arcadeMode {
		top = m.viewScoreLine(pw) + "\n"
	}
	header := m.viewHeader("")
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
		{"?", "cli tips"},
		{"q", "quit"},
	}
	if !m.arcadeMode {
		footerPairs = append(footerPairs, [2]string{"c", "theme"})
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

func (m Model) viewCurrentLine(compact bool) string {
	sep := "\n\n"
	if compact {
		sep = "\n"
	}
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
		return sep + "  " +
			styleCurrent.Render(label) +
			styleCheckmark.Render(m.active.Nickname) +
			styleCurrent.Render("  ·  ") +
			styleCurrentVal.Render(m.active.Email) +
			tags
	}
	if m.arcadeMode {
		return sep + "  " + styleCurrent.Render("NO ACTIVE PLAYER")
	}
	return sep + "  " + styleCurrent.Render("No active profile")
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

	items := "\n"
	for i, p := range m.profiles {
		isCursor := i == m.cursor
		isActive := m.active != nil && p.Nickname == m.active.Nickname

		var ribbon, marker string
		if m.arcadeMode {
			if isCursor {
				ribbon = styleRibbonCursor.Render("ᗧ")
			} else {
				ribbon = " "
			}
			if isActive {
				marker = styleCheckmark.Render(" ★ ")
			} else {
				marker = styleItemDim.Render(" · ")
			}
		} else {
			if isCursor && isActive {
				ribbon = styleRibbonCursor.Render("▎")
			} else if isCursor {
				ribbon = styleRibbonCursor.Render("▎")
			} else if isActive {
				ribbon = styleRibbonActive.Render("▎")
			} else {
				ribbon = " "
			}
			if isActive {
				marker = styleCheckmark.Render(" ✓ ")
			} else {
				marker = styleItemDim.Render(" · ")
			}
		}

		nick := p.Nickname + strings.Repeat(" ", max(0, nickColW-lipgloss.Width(p.Nickname)))
		content := fmt.Sprintf("%s%s  %s", marker, nick, p.Email)

		var line string
		if isCursor {
			// rendered with bg highlight; prepend ribbon outside the bg block
			rowBody := styleItemActive(pw - 2).Render(" " + content)
			line = " " + ribbon + rowBody
		} else {
			line = " " + ribbon + " " + styleItemInactive.Render(content)
		}
		items += "\n" + line
	}
	return items
}

func (m Model) viewStatusLine(compact bool) string {
	if m.statusMsg == "" {
		return ""
	}
	sep := "\n\n"
	if compact {
		sep = "\n"
	}
	if m.statusIsErr {
		prefix := "✕ "
		if m.arcadeMode {
			prefix = "✗ "
		}
		return sep + "  " + lipgloss.NewStyle().Foreground(colorRed).Bold(true).Render(prefix+m.statusMsg)
	}
	prefix := "● "
	if m.arcadeMode {
		prefix = "◎ "
	}
	return sep + "  " + lipgloss.NewStyle().Foreground(colorGreen).Render(prefix+m.statusMsg)
}

func (m Model) viewForm(subtitle string, pw int) string {
	inputW := pw - 6
	header := m.viewHeader(subtitle)

	// full=36 lines, compact=15 lines, mini=12 lines.
	// Switch 2 lines before the larger layout would clip.
	compact := m.height > 0 && m.height < 38
	mini := m.height > 0 && m.height < 17

	var fields string
	for i := range formLabels {
		counter := styleFieldCounter.Render(fmt.Sprintf("[%d/6]", i+1))
		isFocus := i == m.formFocus

		if mini {
			// All fields on a single line; active field uses an inline cursor.
			v := m.formFields[i]
			if isFocus {
				cursor := v + styleTitle.Render("█")
				row := "  " + counter + " " + styleInputLabelActive.Render(formLabels[i]) +
					styleItemDim.Render("  ") + styleCurrentVal.Render(cursor)
				fields += "\n" + row
			} else {
				display := v
				if display == "" {
					display = styleItemDim.Render("—")
				} else {
					display = styleCurrentVal.Render(display)
				}
				row := "  " + counter + " " + styleInputLabel.Render(formLabels[i]) +
					styleItemDim.Render("  ") + display
				fields += "\n" + row
			}
		} else if compact {
			// Active field: label + input box (no subtitle).
			// Inactive fields: single collapsed line, no input box.
			if isFocus {
				fields += "\n"
				fields += "  " + counter + " " + styleInputLabelActive.Render(formLabels[i]) + "\n"
				fields += "  " + styleInputActive(inputW).Render(m.formFields[i]+styleTitle.Render("█"))
			} else {
				v := m.formFields[i]
				display := styleItemDim.Render("—")
				if v != "" {
					display = styleCurrentVal.Render(v)
				}
				fields += "\n  " + counter + " " + styleItemDim.Render("· ") +
					styleInputLabel.Render(formLabels[i]) + styleItemDim.Render("  ") + display
			}
		} else {
			// Full layout: subtitles, input box for every field.
			fields += "\n\n"
			sub := styleBrand.Render("  " + formSubtitles[i])
			if isFocus {
				fields += "  " + counter + " " + styleInputLabelActive.Render(formLabels[i]) + sub + "\n"
				fields += "  " + styleInputActive(inputW).Render(m.formFields[i]+styleTitle.Render("█"))
			} else {
				fields += "  " + styleBrand.Render(fmt.Sprintf("[%d/6]", i+1)) + " " + styleInputLabel.Render(formLabels[i]) + sub + "\n"
				v := m.formFields[i]
				if v == "" {
					v = " "
				}
				fields += "  " + styleInputInactive(inputW).Render(v)
			}
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
	compact := m.height > 0 && m.height < 13 // full=11 lines
	header := m.viewHeader("Delete Profile")
	sep := "\n\n"
	if compact {
		sep = "\n"
	}
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
		warnChip := lipgloss.NewStyle().
			Background(colorRed).
			Foreground(lipgloss.Color("0")).
			Bold(true).
			Padding(0, 1).
			Render(" WARNING ")
		body = sep + "  " + warnChip + "  " + styleDeleteTitle.Render(titleText) +
			"\n  " + styleCurrent.Render(fmt.Sprintf("%s  ·  %s", p.UserName, p.Email)) +
			sep + "  " + styleBrand.Render(bodyText)
	}
	footer := sep + divider(pw) + "\n" + m.footerKeys(pw, [][2]string{
		{"y", "confirm delete"},
		{"n / esc", "cancel"},
	})
	return stylePanelBorder(pw).Render(header + body + footer)
}

func (m Model) viewTips(pw int) string {
	compact := m.height > 0 && m.height < 28 // full=26 lines
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

	cmdStyle := lipgloss.NewStyle().Foreground(colorGreen).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(colorDim)
	labelStyle := lipgloss.NewStyle().Foreground(colorPurple).Bold(true)
	bullet := styleTitle.Render("›")

	sectionLabel := "Skip the UI entirely"
	if m.arcadeMode {
		sectionLabel = "CHEAT CODES"
	}

	var body string
	if compact {
		body = "\n  " + labelStyle.Render(sectionLabel) + "\n"
		for _, t := range tips {
			body += fmt.Sprintf("\n  %s  %s  %s", bullet, cmdStyle.Render(t.cmd), descStyle.Render(t.desc))
		}
		body += "\n"
	} else {
		body = "\n\n  " + labelStyle.Render(sectionLabel) + "\n"
		for _, t := range tips {
			body += fmt.Sprintf("\n  %s  %s\n      %s\n", bullet, cmdStyle.Render(t.cmd), descStyle.Render(t.desc))
		}
	}

	footer := "\n" + divider(pw) + "\n" + m.footerKeys(pw, [][2]string{{"esc / ?", "back"}})
	return stylePanelBorder(pw).Render(header + body + footer)
}

// ── Arcade-only animation views ───────────────────────────────────────────────

func (m Model) viewIntro(pw int) string {
	score := m.viewScoreLine(pw)
	icon := styleTitle.Render("ᗧ")
	title := styleTitle.Render("GIT-SWITCHER")
	pellets := styleDivider.Render("· · ·")
	titleRow := fmt.Sprintf("  %s  %s %s", icon, title, pellets)
	subtitleRow := lipgloss.NewStyle().Foreground(arcadeMazeBlue).Bold(true).Render("     ARCADE MODE")

	track := "  " + m.renderIntroTrack(pw)

	var readyRow, sceneLabel string
	switch m.introPhase {
	case 0:
		sceneLabel = styleBrand.Render("  scene 1: chase")
	case 1:
		sceneLabel = lipgloss.NewStyle().Foreground(arcadeFrightened).Bold(true).Render("  POWER PELLET! ghosts frightened")
	case 2:
		readyStyle := lipgloss.NewStyle().Foreground(colorYellow).Bold(true)
		if m.introReadyFrame%2 == 1 {
			readyStyle = lipgloss.NewStyle().Foreground(colorDim)
		}
		readyRow = "\n\n  " + readyStyle.Render("★   PLAYER 1   READY!   ★")
	}

	skipHint := styleBrand.Render("  [any key to skip]")
	body := titleRow + "\n" + subtitleRow + "\n\n\n" + track + "\n" + sceneLabel + readyRow + "\n\n" + skipHint + "\n"
	return score + "\n" + stylePanelBorder(pw).Render(body)
}

func (m Model) renderIntroTrack(pw int) string {
	numSlots := (pw - 2) / 2
	if numSlots < 1 {
		numSlots = 1
	}

	switch m.introPhase {
	case 0:
		return m.renderChaseTrack(numSlots)
	case 1:
		return m.renderFrightenedTrack(numSlots)
	default:
		return m.renderReadyTrack(numSlots)
	}
}

// Phase 0: Pac chases dots from left, ghosts trail behind.
func (m Model) renderChaseTrack(numSlots int) string {
	pacSlot := m.introPos
	ghost1 := pacSlot - 2
	ghost2 := pacSlot - 4
	ghost3 := pacSlot - 6
	ghost4 := pacSlot - 8

	g1 := lipgloss.NewStyle().Foreground(arcadeGhostRed)
	g2 := lipgloss.NewStyle().Foreground(arcadeGhostPink)
	g3 := lipgloss.NewStyle().Foreground(arcadeGhostCyan)
	g4 := lipgloss.NewStyle().Foreground(arcadeGhostOrange)

	track := ""
	for slot := 0; slot < numSlots; slot++ {
		var ch string
		switch {
		case slot == ghost4 && ghost4 >= 0:
			ch = g4.Render("ᗣ") + " "
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
		case slot == numSlots-1:
			ch = lipgloss.NewStyle().Foreground(colorYellow).Bold(true).Render("●") + " "
		case slot > pacSlot:
			ch = styleDivider.Render("·") + " "
		default:
			ch = "  "
		}
		track += ch
	}
	return track
}

// Phase 1: Pac eats power pellet, ghosts are frightened (blue/white blink).
// Ghosts run RIGHT-to-LEFT, Pac chases LEFT-to-RIGHT and eats them on contact.
func (m Model) renderFrightenedTrack(numSlots int) string {
	pacSlot := m.introPos
	frightStyle := lipgloss.NewStyle().Foreground(arcadeFrightened).Bold(true)
	if m.introPos%2 == 1 {
		frightStyle = lipgloss.NewStyle().Foreground(arcadeFrightWhite).Bold(true)
	}

	// 4 ghosts running away (right-to-left): start at right edge, retreat as Pac advances.
	// Each ghost is "eaten" sequentially based on introGhostsEat.
	ghostBaseSlots := []int{numSlots - 2, numSlots - 4, numSlots - 6, numSlots - 8}
	// As pac moves, ghosts retreat — reduce slot by pac progress factor.
	retreat := pacSlot / 3

	track := ""
	for slot := 0; slot < numSlots; slot++ {
		ch := "  "
		// Pac
		if slot == pacSlot {
			if m.introMouthOpen {
				ch = styleTitle.Render("ᗧ") + " "
			} else {
				ch = styleTitle.Render("●") + " "
			}
			track += ch
			continue
		}
		// Ghosts (only those not yet eaten)
		ghostHere := false
		for gi, g := range ghostBaseSlots {
			gPos := g + retreat
			if gPos == slot && gi >= m.introGhostsEat {
				ch = frightStyle.Render("ᗣ") + " "
				ghostHere = true
				break
			}
			// Show "+200" briefly at slot just eaten
			if gi < m.introGhostsEat && gi == m.introGhostsEat-1 && slot == pacSlot-1 {
				ch = styleScoreLabel.Render("+") + styleScore.Render("200")
				ghostHere = true
				break
			}
		}
		if !ghostHere {
			if slot < pacSlot {
				ch = "  "
			} else {
				ch = styleDivider.Render("·") + " "
			}
		}
		track += ch
	}
	return track
}

// Phase 2 (READY): Pac + 4 ghosts lined up at end.
func (m Model) renderReadyTrack(numSlots int) string {
	pacSlot := numSlots - 1
	g1 := lipgloss.NewStyle().Foreground(arcadeGhostRed)
	g2 := lipgloss.NewStyle().Foreground(arcadeGhostPink)
	g3 := lipgloss.NewStyle().Foreground(arcadeGhostCyan)
	g4 := lipgloss.NewStyle().Foreground(arcadeGhostOrange)

	track := ""
	for slot := 0; slot < numSlots; slot++ {
		var ch string
		switch slot {
		case pacSlot:
			ch = styleTitle.Render("ᗧ") + " "
		case pacSlot - 2:
			ch = g1.Render("ᗣ") + " "
		case pacSlot - 4:
			ch = g2.Render("ᗣ") + " "
		case pacSlot - 6:
			ch = g3.Render("ᗣ") + " "
		case pacSlot - 8:
			ch = g4.Render("ᗣ") + " "
		default:
			ch = "  "
		}
		track += ch
	}
	return track
}

func (m Model) viewSelectFlash(pw int) string {
	score := m.viewScoreLine(pw)
	header := m.viewHeader("")
	currentLine := m.viewCurrentLine(m.height > 0 && m.height < 12+len(m.profiles))
	nickColW := m.nickColumnWidth()

	items := "\n"
	for i, p := range m.profiles {
		if i == m.selFlashProfile {
			nick := p.Nickname + strings.Repeat(" ", max(0, nickColW-lipgloss.Width(p.Nickname)))
			line := fmt.Sprintf("  ᗧ ★ %s  %s", nick, p.Email)
			var flashStyle lipgloss.Style
			switch m.selFlashFrame % 3 {
			case 0:
				flashStyle = lipgloss.NewStyle().Background(colorYellow).Foreground(lipgloss.Color("0")).Width(pw).Bold(true)
			case 1:
				flashStyle = lipgloss.NewStyle().Background(arcadeGhostCyan).Foreground(lipgloss.Color("0")).Width(pw).Bold(true)
			default:
				flashStyle = lipgloss.NewStyle().Background(lipgloss.Color("255")).Foreground(lipgloss.Color("0")).Width(pw).Bold(true)
			}
			items += "\n" + flashStyle.Render(line)
		} else {
			isActive := m.active != nil && p.Nickname == m.active.Nickname
			check := styleItemDim.Render(" · ")
			if isActive {
				check = styleCheckmark.Render(" ★ ")
			}
			nick := p.Nickname + strings.Repeat(" ", max(0, nickColW-lipgloss.Width(p.Nickname)))
			line := fmt.Sprintf("   %s%s  %s", check, nick, p.Email)
			items += "\n  " + styleItemInactive.Render(line)
		}
	}

	// score popup
	popup := ""
	if m.selFlashFrame >= 1 {
		popupStyle := lipgloss.NewStyle().
			Foreground(colorYellow).
			Bold(true).
			Background(colorBgChip).
			Padding(0, 1)
		popup = "\n\n  " + popupStyle.Render("+200  BONUS!")
	}

	footer := "\n\n" + divider(pw) + "\n" + m.footerKeys(pw, [][2]string{})
	return score + "\n" + stylePanelBorder(pw).Render(header+currentLine+items+popup+footer)
}

// Tunnel-scroll transition: vertical pipes move horizontally with offsets.
func (m Model) viewTransition(pw int) string {
	score := m.viewScoreLine(pw)
	header := m.viewHeader("")

	innerW := pw - 2
	if innerW < 1 {
		innerW = 1
	}

	pipeStyle := lipgloss.NewStyle().Foreground(arcadeMazeBlue).Bold(true)
	pelletStyle := styleDivider

	// 5 horizontal lines. Each frame, a vertical pipe column moves left.
	const numLines = 5
	step := innerW / 6
	if step < 1 {
		step = 1
	}
	pipeCol := innerW - 1 - m.transFrame*step

	body := ""
	for r := 0; r < numLines; r++ {
		row := "  "
		for c := 0; c < innerW/2; c++ {
			slot := c * 2
			switch {
			case slot == pipeCol:
				row += pipeStyle.Render("║") + " "
			case slot == pipeCol+2:
				row += pipeStyle.Render("│") + " "
			case (r+c+m.transFrame)%4 == 0:
				row += pelletStyle.Render("·") + " "
			default:
				row += "  "
			}
		}
		body += row + "\n"
	}

	loadingDots := strings.Repeat(".", (m.transFrame%4)+1)
	hint := lipgloss.NewStyle().Foreground(colorYellow).Bold(true).Render("  LOADING" + loadingDots)

	full := header + "\n\n" + body + "\n" + hint + "\n"
	return score + "\n" + stylePanelBorder(pw).Render(full)
}

// Two-act exit animation: GAME OVER (frames 0-5) → INSERT COIN countdown (6-13).
func (m Model) viewExitAnim(pw int) string {
	score := m.viewScoreLine(pw)

	var body string
	if m.exitFrame < 6 {
		textStyle := lipgloss.NewStyle().Bold(true).Foreground(arcadeGhostRed)
		if m.exitFrame%2 == 1 {
			textStyle = lipgloss.NewStyle().Foreground(colorDim)
		}
		text := textStyle.Render("★   G A M E   O V E R   ★")
		sub := styleBrand.Render("  Thanks for playing")
		finalScore := lipgloss.NewStyle().Foreground(colorYellow).Bold(true).
			Render(fmt.Sprintf("FINAL SCORE  %06d", m.score))
		body = "\n\n\n  " + text + "\n\n  " + finalScore + "\n\n  " + sub + "\n\n\n"
	} else {
		// INSERT COIN countdown
		secsLeft := 8 - (m.exitFrame - 6)
		coinStyle := lipgloss.NewStyle().Foreground(colorYellow).Bold(true)
		if (m.exitFrame % 2) == 1 {
			coinStyle = lipgloss.NewStyle().Foreground(colorDim)
		}
		coin := coinStyle.Render("◉  INSERT COIN TO CONTINUE  ◉")
		count := lipgloss.NewStyle().Foreground(arcadeGhostRed).Bold(true).
			Render(fmt.Sprintf("  %d  ", secsLeft))
		sub := styleBrand.Render("  goodbye")
		body = "\n\n\n  " + coin + "\n\n             " + count + "\n\n  " + sub + "\n\n\n"
	}
	return score + "\n" + stylePanelBorder(pw).Render(body)
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
		prefix = ""
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
