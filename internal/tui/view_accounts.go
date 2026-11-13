package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// tabNames returns the three tab labels for the current mode.
func (m Model) tabNames() [3]string {
	if m.arcadeMode {
		return [3]string{"ACCOUNTS", "UTILITIES", "SETTINGS"}
	}
	return [3]string{"Accounts", "Utilities", "Settings"}
}

const tabGap = 2 // cells between tab labels

// tabGeometry returns each tab label's start column and width inside the
// panel. Every label is padded the same whether or not it is active, so
// changing tab never nudges the labels sideways. Shared by the renderer, the
// sliding indicator and the mouse handler.
func (m Model) tabGeometry() (xs, ws [3]int) {
	x := 2
	for i, name := range m.tabNames() {
		xs[i], ws[i] = x, lipgloss.Width(name)+2
		x += ws[i] + tabGap
	}
	return xs, ws
}

// drawTabs renders the tab strip and, under it, a rule whose bright segment
// marks the active tab. While a tab change is animating, the segment slides
// from the old tab to the new one.
func (m Model) drawTabs(c *canvas) {
	xs, ws := m.tabGeometry()
	active := lipgloss.NewStyle().Background(colorBgChip).Foreground(colorPurple).Bold(true).Padding(0, 1)
	if m.arcadeMode {
		active = active.Foreground(colorYellow)
	}
	idle := styleItemDim.Padding(0, 1)

	line := "  "
	for i, name := range m.tabNames() {
		if i > 0 {
			line += strings.Repeat(" ", tabGap)
		}
		if i == m.tabIndex {
			line += active.Render(name)
		} else {
			line += idle.Render(name)
		}
	}
	c.mark(zTabs, 0, line)

	ix, iw := xs[m.tabIndex], ws[m.tabIndex]
	if m.tabAnim {
		ix, iw = int(m.tabX.Pos+0.5), max(1, int(m.tabW.Pos+0.5))
	}
	ix = min(max(ix, 0), c.w)
	iw = min(iw, c.w-ix)
	on := lipgloss.NewStyle().Foreground(colorPurple).Bold(true)
	if m.arcadeMode {
		on = on.Foreground(colorYellow)
	}
	off := lipgloss.NewStyle().Foreground(colorBgHover)
	c.add(off.Render(strings.Repeat("─", ix)) +
		on.Render(strings.Repeat("━", iw)) +
		off.Render(strings.Repeat("─", max(0, c.w-ix-iw))))
}

// profileWindow returns the slice [start, end) of profiles to show when only
// rows of them fit, keeping the cursor roughly centered.
func (m Model) profileWindow(rows int) (start, end int) {
	n := len(m.profiles)
	if rows <= 0 || rows >= n {
		return 0, n
	}
	start = min(max(m.cursor-rows/2, 0), n-rows)
	return start, start + rows
}

// drawAccounts is the main account-switching tab. maxRows > 0 limits the list
// to that many rows (the terminal is too short to show every profile).
func (m Model) drawAccounts(c *canvas, maxRows int) {
	m.header(c, "", m.versionTag())
	if m.headGap() {
		c.blank()
	}
	m.drawTabs(c)
	c.gap()
	c.add(m.currentLine(c.w))
	c.gap()
	m.drawProfiles(c, maxRows)
	c.gap()
	c.add(m.statusLine())

	pairs := [][2]string{{"↑/↓", "navigate"}, {"enter", "switch"}}
	if m.repoKey != "" { // only offer the pin key when there is a repo to pin to
		pairs = append(pairs, [2]string{"p", "pin repo"})
	}
	pairs = append(pairs, [2]string{"a", "add"}, [2]string{"e", "edit"})
	if c.d != dTight {
		pairs = append(pairs, [2]string{"v", secondaryToggleLabel(m.showUsername)}, [2]string{"tab", "next tab"})
	}
	pairs = append(pairs, [2]string{"?", "help"}, [2]string{"q", "quit"})
	c.footer(m, pairs)
}

func (m Model) drawProfiles(c *canvas, maxRows int) {
	if len(m.profiles) == 0 {
		msg := "No profiles yet. Press a to add one."
		if m.arcadeMode {
			msg = "NO PLAYERS YET  ·  PRESS [A] TO INSERT COIN"
		}
		c.add("  " + styleItemDim.Render(msg))
		return
	}

	nickColW := m.nickColumnWidth()
	if limit := c.w / 3; nickColW > limit {
		nickColW = max(8, limit)
	}
	// Column budget: cursor(2) + state(2) + nick + gap(2) + secondary
	secW := max(6, c.w-nickColW-6)

	flashing := m.state == StateSelectFlash
	start, end := m.profileWindow(maxRows)
	for i := start; i < end; i++ {
		p := m.profiles[i]
		glyph := m.scopeGlyph(p)
		nick := padTo(p.Nickname, nickColW)
		sec := truncate(m.profileSecondary(p.Email, p.GHUser), secW)

		// A freshly switched-to profile glows for a moment.
		glow := m.flashLeft > 0 && p.Nickname == m.flashNick

		switch {
		case flashing && i == m.selFlashProfile: // arcade: the picked row strobes
			bg := []lipgloss.Color{colorYellow, arcadeGhostCyan, lipgloss.Color("255")}[m.selFlashFrame%3]
			bar := lipgloss.NewStyle().Background(bg).Foreground(lipgloss.Color("0")).Bold(true)
			c.mark(zProfile, i, bar.Render(padTo("ᗧ ★ "+nick+"  "+sec, c.w)))
		case i == m.cursor && !flashing:
			// Whole row is one highlighted bar with a leading arrow — unmistakable cursor.
			arrow := "❯"
			if m.arcadeMode {
				arrow = "ᗧ"
			}
			bar := styleItemActive(c.w)
			if glow && m.flashLeft%4 >= 2 {
				bar = bar.Foreground(colorGreen)
			}
			c.mark(zProfile, i, bar.Render(truncate(arrow+" "+glyph+" "+nick+"  "+sec, c.w)))
		default:
			stateStyle := styleItemDim
			switch {
			case m.scopeProfile != nil && p.Nickname == m.scopeProfile.Nickname:
				// The repo/session identity is the one git will actually use — call it
				// out louder than the global one it overrides.
				stateStyle = styleScopeMark
			case m.active != nil && p.Nickname == m.active.Nickname:
				stateStyle = styleCheckmark
			}
			nickStyle := styleItemInactive
			if glow {
				nickStyle = lipgloss.NewStyle().Foreground(colorGreen).Bold(true)
			}
			c.mark(zProfile, i, "  "+stateStyle.Render(glyph)+" "+nickStyle.Render(nick)+"  "+styleItemDim.Render(sec))
		}
	}
	if end-start < len(m.profiles) {
		c.add(styleItemDim.Render(fmt.Sprintf("  ↕ %d–%d of %d", start+1, end, len(m.profiles))))
	}
}

// secondaryToggleLabel returns the footer hint for the email/username toggle.
func secondaryToggleLabel(showUsername bool) string {
	if showUsername {
		return "emails"
	}
	return "usernames"
}
