package tui

import (
	"fmt"
	"strings"

	"github.com/aksisonline/gitswitch/internal/git"
	"github.com/aksisonline/gitswitch/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type State int

const (
	StateList State = iota
	StateAdd
	StateEdit
	StateDeleteConfirm
	StateTips
)

type Model struct {
	store    *storage.Store
	profiles []storage.Profile
	cursor   int
	active   *storage.Profile
	state    State
	width    int
	height   int

	formFields  [6]string // nickname, user_name, email, sign_key, ssh_key, gh_user
	formFocus   int
	editingNick string

	statusMsg   string
	statusIsErr bool
}

var formLabels = [6]string{
	"Nickname",
	"User Name",
	"Email",
	"GPG Signing Key",
	"SSH Key Path",
	"GitHub Username",
}

var formSubtitles = [6]string{
	"label shown in this list — not written to git config",
	"git user.name — author name on commits",
	"git user.email — author email on commits",
	"git user.signingkey — optional, leave blank to skip",
	"sets core.sshCommand, e.g. ~/.ssh/id_work — optional",
	"for gh auth switch — optional, leave blank to skip",
}

type switchDoneMsg struct {
	profile  *storage.Profile
	warnings []string
	err      error
}

func New(store *storage.Store) (*Model, error) {
	profiles, err := store.Load()
	if err != nil {
		return nil, err
	}
	active := git.DetectActive(profiles)
	return &Model{
		store:    store,
		profiles: profiles,
		active:   active,
		state:    StateList,
	}, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) panelWidth() int {
	content := minPanelWidth
	for _, p := range m.profiles {
		needed := 3 + 3 + 14 + 2 + len(p.Email) + 6
		if needed > content {
			content = needed
		}
		if nickNeeded := 3 + 3 + len(p.Nickname) + 2 + len(p.Email) + 6; nickNeeded > content {
			content = nickNeeded
		}
	}
	for _, s := range formSubtitles {
		if needed := len(s) + 8; needed > content {
			content = needed
		}
	}
	if content > maxPanelWidth {
		content = maxPanelWidth
	}
	if m.width > 0 {
		if available := m.width - 6; content > available {
			content = available
		}
	}
	if content < minPanelWidth {
		content = minPanelWidth
	}
	return content
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = ws.Width
		m.height = ws.Height
		return m, nil
	}
	switch m.state {
	case StateList:
		return m.updateList(msg)
	case StateAdd, StateEdit:
		return m.updateForm(msg)
	case StateDeleteConfirm:
		return m.updateDelete(msg)
	case StateTips:
		return m.updateTips(msg)
	}
	return m, nil
}

func (m Model) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			m.statusMsg = ""
		case "down", "j":
			if m.cursor < len(m.profiles)-1 {
				m.cursor++
			}
			m.statusMsg = ""
		case "enter":
			if len(m.profiles) > 0 {
				p := m.profiles[m.cursor]
				return m, func() tea.Msg {
					cfg := git.New(true)
					if err := cfg.SetUser(p.UserName, p.Email); err != nil {
						return switchDoneMsg{err: err}
					}
					if err := cfg.SetSignKey(p.SignKey); err != nil {
						return switchDoneMsg{err: err}
					}
					if err := cfg.SetSSHKey(p.SSHKey); err != nil {
						return switchDoneMsg{err: err}
					}
					var warnings []string
					if w := git.SwitchGHUser(p.GHUser); w != "" {
						warnings = append(warnings, w)
					}
					if err := m.store.SetActive(p.Nickname); err != nil {
						return switchDoneMsg{err: err}
					}
					return switchDoneMsg{profile: &p, warnings: warnings}
				}
			}
		case "a":
			m.state = StateAdd
			m.formFields = [6]string{}
			m.formFocus = 0
			m.statusMsg = ""
		case "e":
			if len(m.profiles) > 0 {
				p := m.profiles[m.cursor]
				m.state = StateEdit
				m.editingNick = p.Nickname
				m.formFields = [6]string{p.Nickname, p.UserName, p.Email, p.SignKey, p.SSHKey, p.GHUser}
				m.formFocus = 0
				m.statusMsg = ""
			}
		case "?":
			m.state = StateTips
		}
	case switchDoneMsg:
		if msg.err != nil {
			m.statusMsg = fmt.Sprintf("error: %v", msg.err)
			m.statusIsErr = true
		} else {
			profiles, _ := m.store.Load()
			m.profiles = profiles
			m.active = git.DetectActive(profiles)
			if len(msg.warnings) > 0 {
				m.statusMsg = fmt.Sprintf("switched to %s (warning: %s)", msg.profile.Nickname, msg.warnings[0])
			} else {
				m.statusMsg = fmt.Sprintf("switched to %s", msg.profile.Nickname)
			}
			m.statusIsErr = false
		}
	}
	return m, nil
}

func (m Model) updateForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.state == StateEdit {
		if km, ok := msg.(tea.KeyMsg); ok && km.String() == "ctrl+d" {
			m.state = StateDeleteConfirm
			return m, nil
		}
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.state = StateList
		case "tab", "down":
			if m.formFocus < 5 {
				m.formFocus++
			}
		case "shift+tab", "up":
			if m.formFocus > 0 {
				m.formFocus--
			}
		case "enter":
			if m.formFocus < 5 {
				m.formFocus++
			} else {
				m.submitForm()
			}
		case "backspace":
			f := &m.formFields[m.formFocus]
			if len(*f) > 0 {
				*f = (*f)[:len(*f)-1]
			}
		default:
			if len(msg.String()) == 1 {
				m.formFields[m.formFocus] += msg.String()
			}
		}
	}
	return m, nil
}

func (m *Model) submitForm() {
	nickname := strings.TrimSpace(m.formFields[0])
	userName := strings.TrimSpace(m.formFields[1])
	email := strings.TrimSpace(m.formFields[2])
	signKey := strings.TrimSpace(m.formFields[3])
	sshKey := strings.TrimSpace(m.formFields[4])
	ghUser := strings.TrimSpace(m.formFields[5])

	if nickname == "" || userName == "" || email == "" {
		m.statusMsg = "nickname, user name and email are required"
		m.statusIsErr = true
		m.state = StateList
		return
	}

	var err error
	if m.state == StateEdit {
		err = m.store.Update(m.editingNick, storage.Profile{
			Nickname: nickname,
			UserName: userName,
			Email:    email,
			SignKey:  signKey,
			SSHKey:   sshKey,
			GHUser:   ghUser,
		})
		if err == nil {
			m.statusMsg = fmt.Sprintf("updated '%s'", nickname)
		}
	} else {
		err = m.store.Add(nickname, userName, email, signKey, sshKey, ghUser)
		if err == nil {
			m.statusMsg = fmt.Sprintf("added '%s'", nickname)
		}
	}

	if err != nil {
		m.statusMsg = fmt.Sprintf("error: %v", err)
		m.statusIsErr = true
	} else {
		m.statusIsErr = false
		profiles, _ := m.store.Load()
		m.profiles = profiles
	}
	m.state = StateList
}

func (m Model) updateDelete(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			if len(m.profiles) > 0 {
				nick := m.profiles[m.cursor].Nickname
				if err := m.store.Remove(nick); err != nil {
					m.statusMsg = fmt.Sprintf("error: %v", err)
					m.statusIsErr = true
				} else {
					profiles, _ := m.store.Load()
					m.profiles = profiles
					m.active = git.DetectActive(profiles)
					if m.cursor >= len(m.profiles) && m.cursor > 0 {
						m.cursor--
					}
					m.statusMsg = fmt.Sprintf("deleted '%s'", nick)
					m.statusIsErr = false
				}
			}
			m.state = StateList
		case "n", "N", "esc":
			m.state = StateEdit
		}
	}
	return m, nil
}

func (m Model) updateTips(msg tea.Msg) (tea.Model, tea.Cmd) {
	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "q", "esc", "?", "ctrl+c":
			m.state = StateList
		}
	}
	return m, nil
}

// ── Views ─────────────────────────────────────────────────────────────────────

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

	var currentLine string
	if m.active != nil {
		tags := ""
		if m.active.SSHKey != "" {
			tags += "  " + styleItemDim.Render("ssh")
		}
		if m.active.GHUser != "" {
			tags += "  " + styleItemDim.Render("gh:"+m.active.GHUser)
		}
		currentLine = "\n\n  " +
			styleCurrent.Render("Current  ") +
			styleCheckmark.Render(m.active.Nickname) +
			styleCurrent.Render("  ·  ") +
			styleCurrentVal.Render(m.active.Email) +
			tags
	} else {
		currentLine = "\n\n  " + styleCurrent.Render("No active profile")
	}

	nickColW := 12
	for _, p := range m.profiles {
		if len(p.Nickname) > nickColW {
			nickColW = len(p.Nickname)
		}
	}
	nickColW += 2

	var items string
	if len(m.profiles) == 0 {
		items = "\n\n  " + styleItemDim.Render("No profiles yet. Press [a] to add one.")
	} else {
		items = "\n"
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
			nick := fmt.Sprintf("%-*s", nickColW, p.Nickname)
			line := fmt.Sprintf("%s%s%s  %s", cursor, check, nick, p.Email)
			if i == m.cursor {
				line = styleItemActive(pw).Render(line)
			} else {
				line = "  " + styleItemInactive.Render(line)
			}
			items += "\n" + line
		}
	}

	var statusLine string
	if m.statusMsg != "" {
		s := styleBrand.Render(m.statusMsg)
		if m.statusIsErr {
			s = lipgloss.NewStyle().Foreground(colorRed).Render(m.statusMsg)
		}
		statusLine = "\n\n  " + s
	}

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
		itemW := len(p[0]) + 1 + len(p[1])
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
