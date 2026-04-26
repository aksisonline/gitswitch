package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	pw := m.panelWidth()
	var body string
	switch m.state {
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
	icon := styleTitle.Render("✦")
	title := styleTitle.Render("Git-Switcher")
	var heading string
	if subtitle != "" {
		heading = fmt.Sprintf("  %s  %s%s%s\n", icon, title, styleItemDim.Render(" › "), styleFormTitle.Render(subtitle))
	} else {
		heading = fmt.Sprintf("  %s  %s\n", icon, title)
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

	footer := "\n\n" + divider(pw) + "\n" + m.footerKeys(pw, [][2]string{
		{"↑/↓", "navigate"},
		{"enter", "switch"},
		{"a", "add"},
		{"e", "edit"},
		{"?", "cli tips"},
		{"q", "quit"},
	})

	return stylePanelBorder(pw).Render(header + currentLine + items + statusLine + footer)
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
		return "\n\n  " +
			styleCurrent.Render("Current  ") +
			styleCheckmark.Render(m.active.Nickname) +
			styleCurrent.Render("  ·  ") +
			styleCurrentVal.Render(m.active.Email) +
			tags
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
		return "\n\n  " + styleItemDim.Render("No profiles yet. Press [a] to add one.")
	}
	items := "\n"
	for i, p := range m.profiles {
		cursor := "   "
		if i == m.cursor {
			cursor = "  ❯"
		}
		check := styleItemDim.Render(" · ")
		isActive := m.active != nil && p.Nickname == m.active.Nickname
		if isActive {
			check = " " + styleCheckmark.Render("✓") + " "
		}
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
	s := styleBrand.Render(m.statusMsg)
	if m.statusIsErr {
		s = lipgloss.NewStyle().Foreground(colorRed).Render(m.statusMsg)
	}
	return "\n\n  " + s
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
		body = "\n\n  " + styleDeleteTitle.Render(fmt.Sprintf("Delete '%s'?", p.Nickname)) +
			"\n  " + styleCurrent.Render(fmt.Sprintf("%s  ·  %s", p.UserName, p.Email)) +
			"\n\n  " + styleBrand.Render("This cannot be undone.")
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

	body := "\n\n  " + labelStyle.Render("Skip the UI entirely:") + "\n"
	for _, t := range tips {
		body += fmt.Sprintf("\n  %s\n  %s\n", cmdStyle.Render(t.cmd), descStyle.Render("  "+t.desc))
	}

	footer := "\n" + divider(pw) + "\n" + m.footerKeys(pw, [][2]string{{"esc / ?", "back"}})
	return stylePanelBorder(pw).Render(header + body + footer)
}

func (m Model) footerKeys(pw int, pairs [][2]string) string {
	sep := styleFooter.Render("  ·  ")
	sepW := 5
	var lines []string
	currentLine := "  "
	currentW := 2
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
	if currentLine != "  " {
		lines = append(lines, currentLine)
	}
	return strings.Join(lines, "\n")
}

// Ensure Model satisfies tea.Model at compile time.
var _ tea.Model = Model{}
