package tui

import (
	"fmt"
	"strings"

	"github.com/aksisonline/gitswitch/internal/git"
	"github.com/aksisonline/gitswitch/internal/storage"
	ver "github.com/aksisonline/gitswitch/internal/version"
	tea "github.com/charmbracelet/bubbletea"
)

type switchDoneMsg struct {
	profile  *storage.Profile
	warnings []string
	err      error
}

type upgradeDoneMsg struct {
	err error
}

type versionCheckMsg struct {
	latest string
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = ws.Width
		m.height = ws.Height
		return m, nil
	}
	if vc, ok := msg.(versionCheckMsg); ok {
		m.latestVersion = vc.latest
		m.updateAvailable = ver.IsUpdateAvailable(m.currentVersion, vc.latest)
		return m, nil
	}
	if _, ok := msg.(arcadeTickMsg); ok {
		switch m.state {
		case StateIntro:
			return m.tickIntro()
		case StateSelectFlash:
			return m.tickSelectFlash()
		case StateTransition:
			return m.tickTransition()
		case StateExitAnim:
			return m.tickExitAnim()
		}
		return m, nil
	}
	switch m.state {
	case StateIntro:
		if km, ok := msg.(tea.KeyMsg); ok {
			switch km.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			default:
				m.state = StateList
			}
		}
		return m, nil
	case StateList:
		return m.updateList(msg)
	case StateAdd, StateEdit:
		return m.updateForm(msg)
	case StateDeleteConfirm:
		return m.updateDelete(msg)
	case StateTips:
		return m.updateTips(msg)
	case StateSelectFlash, StateTransition, StateExitAnim:
		return m, nil
	}
	return m, nil
}

func (m Model) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			if m.arcadeMode {
				m.exitFrame = 0
				m.state = StateExitAnim
				return m, arcadeTickCmd()
			}
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
			if m.arcadeMode && len(m.profiles) > 0 {
				m.selFlashFrame = 0
				m.selFlashProfile = m.cursor
				m.state = StateSelectFlash
				return m, arcadeTickCmd()
			}
			if len(m.profiles) > 0 {
				p := m.profiles[m.cursor]
				return m, m.switchProfileCmd(p)
			}
		case "a":
			if m.arcadeMode {
				m.formFields = [6]string{}
				m.formFocus = 0
				m.statusMsg = ""
				m.transTarget = StateAdd
				m.transFrame = 0
				m.state = StateTransition
				return m, arcadeTickCmd()
			}
			m.state = StateAdd
			m.formFields = [6]string{}
			m.formFocus = 0
			m.statusMsg = ""
		case "e":
			if m.arcadeMode && len(m.profiles) > 0 {
				p := m.profiles[m.cursor]
				m.editingNick = p.Nickname
				m.formFields = [6]string{p.Nickname, p.UserName, p.Email, p.SignKey, p.SSHKey, p.GHUser}
				m.formFocus = 0
				m.statusMsg = ""
				m.transTarget = StateEdit
				m.transFrame = 0
				m.state = StateTransition
				return m, arcadeTickCmd()
			}
			if len(m.profiles) > 0 {
				p := m.profiles[m.cursor]
				m.state = StateEdit
				m.editingNick = p.Nickname
				m.formFields = [6]string{p.Nickname, p.UserName, p.Email, p.SignKey, p.SSHKey, p.GHUser}
				m.formFocus = 0
				m.statusMsg = ""
			}
		case "?":
			if m.arcadeMode {
				m.transTarget = StateTips
				m.transFrame = 0
				m.state = StateTransition
				return m, arcadeTickCmd()
			}
			m.state = StateTips
		case "c":
			if !m.arcadeMode {
				m.colorTheme = (m.colorTheme + 1) % 12
				m.statusMsg = fmt.Sprintf("theme: %s (%d/12)", themeNames[m.colorTheme], m.colorTheme+1)
				m.statusIsErr = false
				_ = m.store.SavePrefs(storage.Prefs{ColorTheme: m.colorTheme})
			}
		case "u":
			if m.updateAvailable {
				cmd, err := ver.UpgradeCommand(m.latestVersion)
				if err != nil {
					m.statusMsg = fmt.Sprintf("upgrade error: %v", err)
					m.statusIsErr = true
					return m, nil
				}
				return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
					return upgradeDoneMsg{err: err}
				})
			}
		}
	case upgradeDoneMsg:
		if msg.err != nil {
			m.statusMsg = fmt.Sprintf("upgrade failed: %v", msg.err)
			m.statusIsErr = true
		} else {
			m.statusMsg = fmt.Sprintf("upgraded to %s — restart to apply", m.latestVersion)
			m.statusIsErr = false
			m.updateAvailable = false
		}
	case switchDoneMsg:
		if msg.err != nil {
			m.statusMsg = fmt.Sprintf("error: %v", msg.err)
			m.statusIsErr = true
		} else {
			profiles, err := m.store.Load()
			if err != nil {
				m.statusMsg = fmt.Sprintf("error: %v", err)
				m.statusIsErr = true
				break
			}
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
			if m.arcadeMode {
				return m, m.startBackTransition()
			}
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
				cmd := m.submitForm()
				return m, cmd
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

func (m *Model) submitForm() tea.Cmd {
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
		return nil
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
		profiles, loadErr := m.store.Load()
		if loadErr != nil {
			m.statusMsg = fmt.Sprintf("error: %v", loadErr)
			m.statusIsErr = true
		} else {
			m.statusIsErr = false
			m.profiles = profiles
			m.active = git.DetectActive(profiles)
		}
	}
	if m.arcadeMode {
		return m.startBackTransition()
	}
	m.state = StateList
	return nil
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
					profiles, loadErr := m.store.Load()
					if loadErr != nil {
						m.statusMsg = fmt.Sprintf("error: %v", loadErr)
						m.statusIsErr = true
					} else {
						m.profiles = profiles
						m.active = git.DetectActive(profiles)
						if m.cursor >= len(m.profiles) && m.cursor > 0 {
							m.cursor--
						}
						m.statusMsg = fmt.Sprintf("deleted '%s'", nick)
						m.statusIsErr = false
					}
				}
			}
			if m.arcadeMode {
				m.transTarget = StateList
				m.transFrame = 0
				m.state = StateTransition
				return m, arcadeTickCmd()
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
			if m.arcadeMode {
				return m, m.startBackTransition()
			}
			m.state = StateList
		}
	}
	return m, nil
}

func (m *Model) startBackTransition() tea.Cmd {
	m.transTarget = StateList
	m.transFrame = 0
	m.state = StateTransition
	return arcadeTickCmd()
}

func (m Model) switchProfileCmd(p storage.Profile) tea.Cmd {
	return func() tea.Msg {
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

func (m Model) tickSelectFlash() (tea.Model, tea.Cmd) {
	m.selFlashFrame++
	if m.selFlashFrame >= 4 {
		m.state = StateList
		if len(m.profiles) > 0 && m.selFlashProfile < len(m.profiles) {
			return m, m.switchProfileCmd(m.profiles[m.selFlashProfile])
		}
		return m, nil
	}
	return m, arcadeTickCmd()
}

func (m Model) tickTransition() (tea.Model, tea.Cmd) {
	m.transFrame++
	if m.transFrame >= 5 {
		m.state = m.transTarget
		return m, nil
	}
	return m, arcadeTickCmd()
}

func (m Model) tickExitAnim() (tea.Model, tea.Cmd) {
	m.exitFrame++
	if m.exitFrame >= 8 {
		return m, tea.Quit
	}
	return m, arcadeTickCmd()
}

func (m Model) tickIntro() (tea.Model, tea.Cmd) {
	if m.state != StateIntro {
		return m, nil
	}
	pw := m.panelWidth()
	numSlots := (pw - 2) / 2
	if numSlots < 1 {
		numSlots = 1
	}
	switch m.introPhase {
	case 0:
		m.introPos++
		m.introMouthOpen = !m.introMouthOpen
		if m.introPos >= numSlots {
			m.introPhase = 1
			m.introPos = numSlots - 1
		}
		return m, arcadeTickCmd()
	case 1:
		m.introReadyFrame++
		if m.introReadyFrame >= 10 {
			m.state = StateList
			return m, nil
		}
		return m, arcadeTickCmd()
	}
	return m, nil
}
