package tui

import (
	"fmt"
	"strings"

	"github.com/aksisonline/gitswitch/internal/git"
	"github.com/aksisonline/gitswitch/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type switchDoneMsg struct {
	profile  *storage.Profile
	warnings []string
	err      error
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
