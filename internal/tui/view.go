package tui

import (
	"fmt"
	"strings"

	"github.com/aksisonline/gitswitch/internal/git"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const brandName = "Gitswitch"

// applyColors loads the active palette into the package-level styles.
func (m Model) applyColors() {
	if m.arcadeMode {
		applyTheme(arcadeTheme, true)
		return
	}
	theme := m.colorTheme
	if theme < 0 || theme >= len(normalThemes) {
		theme = 0
	}
	applyTheme(normalThemes[theme], false)
}

func (m Model) View() string {
	s := m.screen()
	top, left := m.origin(s)
	lines := s.lines
	if m.height > 0 {
		if top+len(lines) > m.height {
			top = max(0, m.height-len(lines))
		}
		if len(lines) > m.height {
			lines = lines[:m.height]
		}
	}
	pad := strings.Repeat(" ", left)
	var b strings.Builder
	b.WriteString(strings.Repeat("\n", top))
	for i, l := range lines {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(pad)
		b.WriteString(l)
	}
	return b.String()
}

// screenAt draws the current state at density d.
func (m Model) screenAt(d int) screen {
	pw := m.panelWidth()
	c := newCanvas(pw, d)
	var pre []string
	score := func() {
		if m.arcadeMode {
			pre = []string{clip(m.viewScoreLine(pw), pw+2)}
		}
	}
	switch m.state {
	case StateIntro:
		score()
		m.drawIntro(c)
	case StateTransition:
		score()
		m.drawTransition(c)
	case StateExitAnim:
		score()
		m.drawExitAnim(c)
	case StateAdd:
		m.drawForm(c, "Add Profile")
	case StateEdit:
		m.drawForm(c, "Edit  "+m.editingNick)
	case StateDeleteConfirm:
		m.drawDeleteConfirm(c)
	case StateTips:
		m.drawTips(c)
	case StateUpdatePrompt:
		m.drawUpdatePrompt(c)
	case StateWhatsNew:
		m.drawWhatsNew(c)
	case StateWizardWelcome:
		m.drawWizardWelcome(c)
	case StateWizardDetect:
		m.drawWizardDetect(c)
	case StateWizardImport:
		m.drawWizardImport(c)
	case StateWizardAddMore:
		m.drawWizardAddMore(c)
	case StateWizardDone:
		m.drawWizardDone(c)
	case StateShellConfirm:
		m.drawShellConfirm(c)
	default:
		score()
		switch {
		case m.tabIndex == 1 && m.state == StateList:
			m.drawUtilities(c)
		case m.tabIndex == 2 && m.state == StateList:
			m.drawSettings(c)
		default:
			m.drawAccounts(c, 0)
			// Still too tall at the tightest density: window the list around
			// the cursor instead of letting rows fall off the screen.
			if d == dTight && m.height > 0 && len(m.profiles) > 1 {
				if over := finishPanel(pw, c, pre...).height() - m.height; over > 0 {
					rows := max(1, len(m.profiles)-over-1)
					c = newCanvas(pw, d)
					m.drawAccounts(c, rows)
				}
			}
		}
	}
	return finishPanel(pw, c, pre...)
}

// viewScoreLine — arcade-only: top score panel like classic arcades.
func (m Model) viewScoreLine(pw int) string {
	left := styleScoreLabel.Render("1UP") + "  " + styleScore.Render(fmt.Sprintf("%06d", m.score))
	right := styleScoreLabel.Render("HIGH SCORE") + "  " + styleScore.Render(fmt.Sprintf("%06d", m.hiScore))
	gap := pw - lipgloss.Width(left) - lipgloss.Width(right) - 4
	if gap < 1 {
		gap = 1
	}
	return "  " + left + strings.Repeat(" ", gap) + right
}

// brandTitle renders the app name. For a moment after launch a bright band
// sweeps across it (m.shimmer is the band's position, 0 = off).
func (m Model) brandTitle() string {
	name := brandName
	if m.arcadeMode {
		name = strings.ToUpper(name)
	}
	if m.shimmer == 0 {
		return styleTitle.Render(name)
	}
	glint := lipgloss.NewStyle().Bold(true).Foreground(colorWhite)
	center := m.shimmer - 3
	var b strings.Builder
	for i, r := range name {
		if d := i - center; d >= -1 && d <= 1 {
			b.WriteString(glint.Render(string(r)))
		} else {
			b.WriteString(styleTitle.Render(string(r)))
		}
	}
	return b.String()
}

// header draws the brand line. subtitle names the current dialog ("" on the
// main tabs, where the tagline shows instead); right is an optional
// right-aligned extra such as the setup stepper or the version.
func (m Model) header(c *canvas, subtitle, right string) {
	var left string
	if m.arcadeMode {
		left = "  " + styleTitle.Render("ᗧ") + "  " + m.brandTitle() + " " + styleDivider.Render("· · ·")
		if subtitle != "" {
			left += styleItemDim.Render(" ► ") + styleFormTitle.Render(subtitle)
		}
	} else {
		left = "  " + styleTitle.Render("◆") + "  " + m.brandTitle()
		if subtitle != "" {
			left += styleItemDim.Render("  ›  ") + styleFormTitle.Render(subtitle)
		} else if tagline := "   " + styleBrand.Render("identity manager for git"); lipgloss.Width(left)+lipgloss.Width(tagline)+lipgloss.Width(right)+2 <= c.w {
			left += tagline
		}
	}
	if right != "" {
		right += " "
	}
	c.add(spread(left, right, c.w))
}

// versionTag is the dim version shown at the right of the main header.
func (m Model) versionTag() string {
	if m.currentVersion == "" {
		return ""
	}
	return styleItemDim.Render(m.currentVersion)
}

// currentLine reports the identity git will actually use here: a repo pin or a
// session identity overrides the global profile, and saying otherwise would be
// a lie.
func (m Model) currentLine(w int) string {
	cur, nameStyle := m.active, styleCheckmark
	if m.scopeProfile != nil {
		cur, nameStyle = m.scopeProfile, styleScopeMark
	}
	if cur == nil {
		if m.arcadeMode {
			return "  " + styleCurrent.Render("NO ACTIVE PLAYER")
		}
		return "  " + styleCurrent.Render("No active profile")
	}
	tags := ""
	switch m.scope {
	case git.ScopeRepo:
		tags += "  " + styleScopeMark.Render("pinned")
	case git.ScopeSession:
		tags += "  " + styleScopeMark.Render("session")
	}
	if m.pinnedInactive != "" {
		tags += "  " + styleItemDim.Render("pin: "+m.pinnedInactive+" (isolation off)")
	}
	if cur.SSHKey != "" {
		tags += "  " + styleItemDim.Render("ssh")
	}
	if cur.GHUser != "" {
		tags += "  " + styleItemDim.Render("gh:"+cur.GHUser)
	}
	label := "Current  "
	if m.arcadeMode {
		label = "PLAYER 1 ► "
	}
	nick := truncate(cur.Nickname, max(8, w/3))
	// budget for the secondary value = width − indent − label − nick − separator − tags
	secMax := w - 2 - lipgloss.Width(label) - lipgloss.Width(nick) - 5 - lipgloss.Width(tags)
	line := "  " + styleCurrent.Render(label) + nameStyle.Render(nick)
	if secMax >= 6 {
		line += styleCurrent.Render("  ·  ") + styleCurrentVal.Render(truncate(m.profileSecondary(cur.Email, cur.GHUser), secMax))
	}
	return line + tags
}

func (m Model) nickColumnWidth() int {
	nickColW := 12
	for _, p := range m.profiles {
		if nickW := lipgloss.Width(p.Nickname); nickW > nickColW {
			nickColW = nickW
		}
	}
	result := nickColW + 2
	if result > 22 {
		result = 22
	}
	return result
}

// profileSecondary returns the secondary column text for a profile, honoring
// the email/username toggle. Falls back to email when a username is missing.
func (m Model) profileSecondary(email, ghUser string) string {
	if m.showUsername && ghUser != "" {
		return "gh:" + ghUser
	}
	return email
}

// statusLine is the one-row slot under the content. It is always drawn — even
// empty — so a message appearing or clearing never changes the panel height.
// Priority: work in progress, then the last result, then a pending update.
func (m Model) statusLine() string {
	switch {
	case m.state == StateSelectFlash && m.selFlashFrame >= 1:
		return "  " + lipgloss.NewStyle().Foreground(colorYellow).Bold(true).Background(colorBgChip).Padding(0, 1).Render("+200  BONUS!")
	case m.busy != "":
		return "  " + lipgloss.NewStyle().Foreground(colorPurple).Bold(true).Render(m.spinner()) + " " +
			styleItemDim.Render(m.busy)
	case m.statusMsg != "" && m.statusIsErr:
		prefix := "✕ "
		if m.arcadeMode {
			prefix = "✗ "
		}
		return "  " + lipgloss.NewStyle().Foreground(colorRed).Bold(true).Render(prefix+m.statusMsg)
	case m.statusMsg != "":
		prefix := "● "
		if m.arcadeMode {
			prefix = "◎ "
		}
		return "  " + lipgloss.NewStyle().Foreground(colorGreen).Render(prefix+m.statusMsg)
	case m.updateAvailable && m.onTabs():
		if m.arcadeMode {
			return "  " + styleChipBonus.Render("★ BONUS STAGE") + "  " + styleScore.Render(m.latestVersion) +
				styleBrand.Render("  ·  press ") + styleFooterKey.Render("u")
		}
		return "  " + styleChipBox().Render("UPDATE") + "  " + styleCurrentVal.Render(m.latestVersion) +
			styleBrand.Render("  ·  press ") + styleFooterKey.Render("u") + styleBrand.Render(" to upgrade")
	}
	return ""
}

// formChrome is how many rows the add/edit screen spends outside the form
// itself: borders (2), header + blank (2), divider + keys (2).
const formChrome = 6

// formSize is the width/height handed to the huh form so it scrolls inside
// the panel instead of pushing the panel off-screen.
func (m Model) formSize() (w, h int) {
	w = m.panelWidth() - 6
	if m.height > 0 {
		h = max(4, m.height-formChrome)
	}
	return w, h
}

func (m Model) drawForm(c *canvas, subtitle string) {
	m.header(c, subtitle, "")
	c.blank()
	if m.form != nil {
		lines := strings.Split(strings.TrimRight(m.form.View(), "\n"), "\n")
		// huh sizes itself to the height it was given, but a wrapped
		// description or a validation message can still push it over. Window
		// the rows around the focused field (the one with the ┃ bar) so the
		// panel never outgrows the terminal.
		if _, h := m.formSize(); h > 0 && len(lines) > h {
			focus := 0
			for i, l := range lines {
				if strings.Contains(l, "┃") {
					focus = i
					break
				}
			}
			// A validation message is the form's last line — always keep it.
			var tail []string
			if last := lines[len(lines)-1]; strings.Contains(last, "*") {
				tail, lines, h = []string{last}, lines[:len(lines)-1], h-1
			}
			start := min(max(focus-1, 0), max(len(lines)-h, 0))
			lines = append(lines[start:min(start+h, len(lines))], tail...)
		}
		for _, l := range lines {
			c.add(l)
		}
	}
	pairs := [][2]string{{"↑/↓", "move"}, {"enter", "next / save"}, {"esc", "cancel"}}
	if m.state == StateEdit {
		pairs = append(pairs, [2]string{"ctrl+d", "delete"})
	}
	c.footer(m, pairs)
}

func (m Model) drawDeleteConfirm(c *canvas) {
	m.header(c, "Delete Profile", "")
	c.gap()
	if len(m.profiles) > 0 && m.cursor < len(m.profiles) {
		p := m.profiles[m.cursor]
		title, body := fmt.Sprintf("Delete '%s'?", p.Nickname), "This cannot be undone."
		if m.arcadeMode {
			title, body = fmt.Sprintf("GAME OVER?  '%s'", p.Nickname), "THIS ACTION CANNOT BE UNDONE."
		}
		chip := lipgloss.NewStyle().Background(colorRed).Foreground(lipgloss.Color("0")).Bold(true).Padding(0, 1).Render("WARNING")
		c.add("  " + chip + "  " + styleDeleteTitle.Render(title))
		c.add("  " + styleCurrent.Render(p.UserName+"  ·  "+p.Email))
		c.gap()
		c.add("  " + styleBrand.Render(body))
	}
	c.gap()
	c.footer(m, [][2]string{{"y", "confirm delete"}, {"n / esc", "cancel"}})
}

func (m Model) drawTips(c *canvas) {
	m.header(c, "CLI Quick Reference", "")
	tips := [][2]string{
		{"gitswitch <nickname>", "switch without opening this screen"},
		{"gitswitch current", "show the identity in use here"},
		{"gitswitch list", "list all saved profiles"},
		{"gitswitch pin <nickname>", "always use a profile in this repo"},
		{"gitswitch login", "connect another GitHub account"},
		{"gitswitch doctor", "check that everything is wired up"},
		{"gitswitch --help", "every command and flag"},
	}
	cmdStyle := lipgloss.NewStyle().Foreground(colorGreen).Bold(true)
	bullet := styleTitle.Render("›")
	label := "Skip the UI entirely"
	if m.arcadeMode {
		label = "CHEAT CODES"
	}
	c.gap()
	c.add("  " + lipgloss.NewStyle().Foreground(colorPurple).Bold(true).Render(label))
	if c.d != dTight {
		c.blank()
	}
	cmdW := 0
	for _, t := range tips {
		cmdW = max(cmdW, lipgloss.Width(t[0]))
	}
	for i, t := range tips {
		switch {
		case c.d == dRoomy && c.w < cmdW+40:
			// Narrow and roomy: command on one line, description beneath it.
			c.add("  " + bullet + "  " + cmdStyle.Render(t[0]))
			c.add("     " + styleItemDim.Render(t[1]))
			if i < len(tips)-1 {
				c.blank()
			}
		default:
			pad := strings.Repeat(" ", cmdW-lipgloss.Width(t[0]))
			c.add("  " + bullet + "  " + cmdStyle.Render(t[0]) + pad + "  " + styleItemDim.Render(t[1]))
		}
	}
	c.gap()
	c.footer(m, [][2]string{{"esc / ?", "back"}})
}

// ── Arcade-only animation views ───────────────────────────────────────────────

// introSlots is how many 2-cell track slots fit inside the panel.
func (m Model) introSlots() int {
	n := (m.panelWidth() - 4) / 2
	if n < 1 {
		n = 1
	}
	return n
}

func (m Model) drawIntro(c *canvas) {
	c.add("  " + styleTitle.Render("ᗧ") + "  " + styleTitle.Render("GITSWITCH") + " " + styleDivider.Render("· · ·"))
	c.add(lipgloss.NewStyle().Foreground(arcadeMazeBlue).Bold(true).Render("     ARCADE MODE"))
	c.blank()
	c.gap()
	n := m.introSlots()
	switch m.introPhase {
	case 0:
		c.add("  " + m.renderChaseTrack(n))
	case 1:
		c.add("  " + m.renderFrightenedTrack(n))
	default:
		c.add("  " + m.renderReadyTrack(n))
	}
	c.blank()
	// The READY row is always reserved so the panel doesn't grow in phase 2.
	ready := ""
	if m.introPhase == 2 {
		readyStyle := lipgloss.NewStyle().Foreground(colorYellow).Bold(true)
		if m.introReadyFrame%2 == 1 {
			readyStyle = lipgloss.NewStyle().Foreground(colorDim)
		}
		ready = "  " + readyStyle.Render("★   PLAYER 1   READY!   ★")
	}
	c.add(ready)
	c.gap()
	c.add(styleBrand.Render("  [any key to skip]"))
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
			gPos := g - retreat
			if gPos < 0 {
				gPos = 0
			} else if gPos >= numSlots {
				gPos = numSlots - 1
			}
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

// Tunnel-scroll transition: vertical pipes move horizontally with offsets.
func (m Model) drawTransition(c *canvas) {
	m.header(c, "", "")
	c.gap()
	cols := (c.w - 2) / 2
	if cols < 1 {
		cols = 1
	}
	pipeStyle := lipgloss.NewStyle().Foreground(arcadeMazeBlue).Bold(true)
	step := max(1, cols/6)
	pipeCol := max(0, cols-1-m.transFrame*step)
	rows := 5
	if c.d == dTight {
		rows = 3
	}
	for r := 0; r < rows; r++ {
		row := "  "
		for col := 0; col < cols; col++ {
			switch {
			case col == pipeCol:
				row += pipeStyle.Render("║") + " "
			case col == pipeCol+1:
				row += pipeStyle.Render("│") + " "
			case (r+col+m.transFrame)%4 == 0:
				row += styleDivider.Render("·") + " "
			default:
				row += "  "
			}
		}
		c.add(row)
	}
	c.gap()
	// Fixed-width dots so the label doesn't jitter as they cycle.
	dots := strings.Repeat(".", (m.transFrame%4)+1)
	c.add(lipgloss.NewStyle().Foreground(colorYellow).Bold(true).Render("  LOADING" + fmt.Sprintf("%-4s", dots)))
}

// Two-act exit animation: GAME OVER (frames 0-5) → INSERT COIN countdown (6-13).
// Both acts are the same height so the panel doesn't jump between them.
func (m Model) drawExitAnim(c *canvas) {
	c.gap()
	c.blank()
	if m.exitFrame < 6 {
		textStyle := lipgloss.NewStyle().Bold(true).Foreground(arcadeGhostRed)
		if m.exitFrame%2 == 1 {
			textStyle = lipgloss.NewStyle().Foreground(colorDim)
		}
		c.add("  " + textStyle.Render("★   G A M E   O V E R   ★"))
		c.blank()
		c.add("  " + lipgloss.NewStyle().Foreground(colorYellow).Bold(true).Render(fmt.Sprintf("FINAL SCORE  %06d", m.score)))
		c.blank()
		c.add(styleBrand.Render("    Thanks for playing"))
	} else {
		coinStyle := lipgloss.NewStyle().Foreground(colorYellow).Bold(true)
		if m.exitFrame%2 == 1 {
			coinStyle = lipgloss.NewStyle().Foreground(colorDim)
		}
		c.add("  " + coinStyle.Render("◉  INSERT COIN TO CONTINUE  ◉"))
		c.blank()
		c.add("             " + lipgloss.NewStyle().Foreground(arcadeGhostRed).Bold(true).Render(fmt.Sprintf("  %d  ", 8-(m.exitFrame-6))))
		c.blank()
		c.add(styleBrand.Render("    goodbye"))
	}
	c.blank()
	c.gap()
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
	// Budget against the panel's usable inner width, not its outer width: the
	// border adds 1 column of padding each side, so a line built to pw gets
	// re-wrapped by lipgloss — which splits a key from its own label.
	avail := pw - 2
	var lines []string
	currentLine := "  " + prefix
	currentW := 2 + lipgloss.Width(prefix)
	for i, p := range pairs {
		item := styleFooterKey.Render(p[0]) + styleFooter.Render(" "+p[1])
		itemW := lipgloss.Width(p[0]) + 1 + lipgloss.Width(p[1])
		if i > 0 {
			if currentW+sepW+itemW > avail {
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

func (m Model) drawWhatsNew(c *canvas) {
	pw := c.w + 2
	c.blank()
	c.add(lipgloss.NewStyle().Bold(true).Foreground(colorPurple).Render("  What's New in " + m.currentVersion))
	c.blank()

	noteLines := renderNotes(m.whatsNewBody, pw)
	visibleH := notesViewportHeight(m.height)
	total := len(noteLines)
	offset := min(max(m.whatsNewScroll, 0), max(total-visibleH, 0))
	end := min(offset+visibleH, total)
	for _, l := range noteLines[offset:end] {
		c.add(l)
	}

	keys := [][2]string{{"any key", "continue"}}
	if total > visibleH {
		c.blank()
		c.add(styleItemDim.Render(fmt.Sprintf("  (%d–%d of %d lines · ↑/↓ scroll)", offset+1, end, total)))
		keys = [][2]string{{"↑/↓", "scroll"}, {"any key", "continue"}}
	}
	c.blank()
	c.footer(m, keys)
}

// Ensure Model satisfies tea.Model at compile time.
var _ tea.Model = Model{}
