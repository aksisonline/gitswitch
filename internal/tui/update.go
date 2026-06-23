package tui

import (
	"fmt"
	"os/exec"
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
	if dcm, ok := msg.(detectConfigsMsg); ok {
		m.detectedProfiles = dcm.profiles
		m.importSelected = make([]bool, len(dcm.profiles))
		for i := range m.importSelected {
			m.importSelected[i] = true
		}
		if len(dcm.profiles) > 0 {
			m.state = StateWizardImport
		} else {
			m.state = StateWizardAddMore
		}
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
	if mm, ok := msg.(tea.MouseMsg); ok {
		return m.handleMouse(mm)
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
	case StateNoProfiles:
		return m.updateNoProfiles(msg)
	case StateWhatsNew:
		return m.updateWhatsNew(msg)
	case StateWizardWelcome, StateWizardDetect, StateWizardImport, StateWizardAddMore, StateWizardDone:
		return m.updateWizard(msg)
	case StateSelectFlash, StateTransition, StateExitAnim:
		if km, ok := msg.(tea.KeyMsg); ok && km.String() == "ctrl+c" {
			return m, tea.Quit
		}
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
			m.statusMsg = ""
			switch m.tabIndex {
			case 0:
				if m.cursor > 0 {
					m.cursor--
				}
			case 1:
				if m.utilityFocus > 0 {
					m.utilityFocus--
				}
			case 2:
				if m.settingsFocus > 0 {
					m.settingsFocus--
				}
			}
		case "down", "j":
			m.statusMsg = ""
			switch m.tabIndex {
			case 0:
				if m.cursor < len(m.profiles)-1 {
					m.cursor++
				}
			case 1:
				if m.utilityFocus < 2 {
					m.utilityFocus++
				}
			case 2:
				if m.settingsFocus < 1 {
					m.settingsFocus++
				}
			}
		case "enter":
			switch m.tabIndex {
			case 0: // Accounts — switch profile
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
			case 1: // Utilities — toggle focused item
				if m.utilityFocus == 0 {
					m.shellEnabled = !m.shellEnabled
					_ = m.store.SavePrefs(storage.Prefs{
						ColorTheme:    m.colorTheme,
						SplashSeen020: m.splashSeen020,
						ShellEnabled:  m.shellEnabled,
					})
					if m.shellEnabled {
						m.statusMsg = "shell integration enabled"
					} else {
						m.statusMsg = "shell integration disabled"
					}
				}
			case 2: // Settings — nothing interactive yet
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
			m.wizardStep = 0
			m.state = StateWizardAddMore
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
				if err := m.store.SavePrefs(storage.Prefs{ColorTheme: m.colorTheme, SplashSeen020: m.splashSeen020, ShellEnabled: m.shellEnabled}); err != nil {
					m.statusMsg = fmt.Sprintf("theme: %s — could not save: %v", themeNames[m.colorTheme], err)
					m.statusIsErr = true
				} else {
					m.statusMsg = fmt.Sprintf("theme: %s (%d/12)", themeNames[m.colorTheme], m.colorTheme+1)
					m.statusIsErr = false
				}
			}
		case "left":
			if !m.arcadeMode && m.tabIndex == 2 {
				m.colorTheme = (m.colorTheme + 11) % 12
				if err := m.store.SavePrefs(storage.Prefs{ColorTheme: m.colorTheme, SplashSeen020: m.splashSeen020, ShellEnabled: m.shellEnabled}); err != nil {
					m.statusMsg = "could not save theme"
				}
			}
		case "right":
			if !m.arcadeMode && m.tabIndex == 2 {
				m.colorTheme = (m.colorTheme + 1) % 12
				if err := m.store.SavePrefs(storage.Prefs{ColorTheme: m.colorTheme, SplashSeen020: m.splashSeen020, ShellEnabled: m.shellEnabled}); err != nil {
					m.statusMsg = "could not save theme"
				}
			}
		case "1", "2", "3":
			if !m.arcadeMode {
				m.tabIndex = int(msg.String()[0] - '1')
				m.statusMsg = ""
			}
		case "tab":
			if !m.arcadeMode {
				m.tabIndex = (m.tabIndex + 1) % 3
				m.statusMsg = ""
			}
		case "shift+tab":
			if !m.arcadeMode {
				m.tabIndex = (m.tabIndex + 2) % 3
				m.statusMsg = ""
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
			if msg.Paste {
				m.formFields[m.formFocus] += string(msg.Runes)
			} else if len(msg.String()) == 1 {
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
	if m.selFlashFrame >= 6 {
		m.state = StateList
		m.score += 200
		if len(m.profiles) > 0 && m.selFlashProfile < len(m.profiles) {
			return m, m.switchProfileCmd(m.profiles[m.selFlashProfile])
		}
		return m, nil
	}
	return m, arcadeTickCmd()
}

func (m Model) tickTransition() (tea.Model, tea.Cmd) {
	m.transFrame++
	if m.transFrame >= 6 {
		m.state = m.transTarget
		return m, nil
	}
	return m, arcadeTickCmd()
}

func (m Model) tickExitAnim() (tea.Model, tea.Cmd) {
	m.exitFrame++
	// 14 frames: GAME OVER flash (0-5) → INSERT COIN countdown (6-13)
	if m.exitFrame >= 14 {
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
		// PAC eats dots moving right
		m.introPos++
		m.introMouthOpen = !m.introMouthOpen
		if m.introPos >= numSlots {
			m.introPhase = 1
			m.introPos = 0
			m.introReadyFrame = 0
		}
		return m, arcadeTickCmd()
	case 1:
		// Power-pellet eaten — ghosts turn frightened, PAC chases left to right
		m.introPos++
		m.introMouthOpen = !m.introMouthOpen
		if m.introPos%3 == 0 && m.introGhostsEat < 4 {
			m.introGhostsEat++
		}
		if m.introPos >= numSlots {
			m.introPhase = 2
			m.introReadyFrame = 0
		}
		return m, arcadeTickCmd()
	case 2:
		// READY flash
		m.introReadyFrame++
		if m.introReadyFrame >= 10 {
			m.state = StateList
			return m, nil
		}
		return m, arcadeTickCmd()
	}
	return m, nil
}

// panelTopY returns the absolute screen Y of the panel's top border.
// Uses a state-aware height estimate since View() has a value receiver.
func (m Model) panelTopY() int {
	if m.height == 0 {
		return 0
	}
	var panelH int
	switch m.state {
	case StateWizardWelcome:
		panelH = 20
	case StateWizardDetect:
		panelH = 16
	case StateWizardImport:
		panelH = 14 + len(m.detectedProfiles)*3
	case StateWizardAddMore:
		panelH = 26 // header(1) + tabgap(2) + 3 items×5 + footer(3) + border(2) + margin
	case StateWizardDone:
		panelH = 20
	case StateWhatsNew:
		panelH = 24
	default:
		if m.tabIndex == 1 { // utilities: 3 items
			panelH = 22
		} else if m.tabIndex == 2 { // settings: 2 items
			panelH = 18
		} else {
			panelH = 13 + len(m.profiles)
			if panelH < 16 {
				panelH = 16
			}
		}
	}
	top := (m.height - panelH) / 2
	if top < 0 {
		top = 0
	}
	return top
}

func (m Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if msg.Action != tea.MouseActionPress {
		return m, nil
	}

	pw := m.panelWidth()
	panelLeft := (m.width - pw) / 2
	panelTop := m.panelTopY()
	// relY: Y relative to panel top border (0 = border, 1 = first content line)
	relY := msg.Y - panelTop
	// contentX: X within panel content (after left border char)
	contentX := msg.X - panelLeft - 1

	switch msg.Button {
	case tea.MouseButtonWheelUp:
		if m.state == StateList && m.cursor > 0 {
			m.cursor--
			m.statusMsg = ""
		}
		if m.state == StateWizardAddMore && m.wizardStep > 0 {
			m.wizardStep--
		}
	case tea.MouseButtonWheelDown:
		if m.state == StateList && m.cursor < len(m.profiles)-1 {
			m.cursor++
			m.statusMsg = ""
		}
		if m.state == StateWizardAddMore && m.wizardStep < 2 {
			m.wizardStep++
		}
	case tea.MouseButtonLeft:
		// What's new: any click dismisses
		if m.state == StateWhatsNew {
			m.state = StateList
			m.splashSeen020 = true
			_ = m.store.SavePrefs(storage.Prefs{ColorTheme: m.colorTheme, SplashSeen020: true})
			return m, nil
		}

		// Wizard welcome: click the CTA button area (relY 10-14)
		if m.state == StateWizardWelcome && relY >= 10 && relY <= 14 {
			m.state = StateWizardDetect
			return m, m.detectExistingConfigsCmd()
		}

		// Wizard add-more: click buttons
		// Layout: header(1) + blank+status(2) + blank+divider(2) + blank+question(2) + blank+oauth(3) + blank+manual(3) + blank+done(3)
		// relY: 0=border, 1=header, 2=blank, 3=status, 4=blank, 5=divider, 6=blank, 7=question
		//       8=blank, 9=oauth-top, 10=oauth-mid, 11=oauth-bot
		//       12=blank, 13=manual-top, 14=manual-mid, 15=manual-bot
		//       16=blank, 17=done-top, 18=done-mid, 19=done-bot
		if m.state == StateWizardAddMore {
			switch {
			case relY >= 9 && relY <= 11: // OAuth button
				m.wizardStep = 0
				m.statusMsg = "GitHub OAuth coming in v0.2.1 — use Add Manually for now"
			case relY >= 13 && relY <= 15: // Manual button
				m.wizardStep = 1
				m.state = StateAdd
				m.formFields = [6]string{}
				m.formFocus = 0
				m.statusMsg = ""
			case relY >= 17 && relY <= 19: // Done button
				m.wizardStep = 2
				m.state = StateList
				m.tabIndex = 0
			}
			return m, nil
		}

		if m.state == StateList && !m.arcadeMode {
			// Tab header click — relY==3 is the tab header line
			// Content: "  [ Accounts ]   Utilities   Settings"
			//            0123456789012345678901234567890123456
			// Accounts ends ~15, Utilities ends ~25, Settings rest
			if relY == 3 {
				switch {
				case contentX <= 15:
					m.tabIndex = 0
				case contentX <= 26:
					m.tabIndex = 1
				default:
					m.tabIndex = 2
				}
				m.statusMsg = ""
				return m, nil
			}
		}

		// Utilities tab item clicks
		// Item boxes start at relY=5, each item = 5 relY rows (4 lines + 1 blank prefix)
		if m.state == StateList && m.tabIndex == 1 {
			if relY >= 5 {
				itemOffset := relY - 5
				itemIdx := itemOffset / 5
				withinItem := itemOffset % 5
				// withinItem 0=blank prefix of item, 1=top border, 2=line1, 3=line2, 4=bottom border
				if withinItem > 0 && withinItem < 5 && itemIdx >= 0 && itemIdx <= 2 {
					m.utilityFocus = itemIdx
					if itemIdx == 0 {
						// clicking the shell integration toggle area — toggle it
						m.shellEnabled = !m.shellEnabled
						_ = m.store.SavePrefs(storage.Prefs{
							ColorTheme:    m.colorTheme,
							SplashSeen020: m.splashSeen020,
							ShellEnabled:  m.shellEnabled,
						})
						if m.shellEnabled {
							m.statusMsg = "shell integration enabled"
						} else {
							m.statusMsg = "shell integration disabled"
						}
					}
					return m, nil
				}
			}
		}

		// Settings tab item clicks
		if m.state == StateList && m.tabIndex == 2 {
			if relY >= 5 {
				itemOffset := relY - 5
				itemIdx := itemOffset / 5
				withinItem := itemOffset % 5
				if withinItem > 0 && withinItem < 5 && itemIdx >= 0 && itemIdx <= 1 {
					m.settingsFocus = itemIdx
					if itemIdx == 1 { // theme box
						// right side of line1 = cycle forward
						if contentX > pw/2 {
							m.colorTheme = (m.colorTheme + 1) % 12
						} else {
							m.colorTheme = (m.colorTheme + 11) % 12
						}
						_ = m.store.SavePrefs(storage.Prefs{
							ColorTheme:    m.colorTheme,
							SplashSeen020: m.splashSeen020,
							ShellEnabled:  m.shellEnabled,
						})
					}
					return m, nil
				}
			}
		}

		// Profile rows (accounts tab, relY 7+)
		// Layout: 0=border, 1=header, 2=blank, 3=tab-header, 4=blank, 5=current, 6=blank, 7+=items
		if m.state == StateList && m.tabIndex == 0 && len(m.profiles) > 0 {
			idx := relY - 7
			if idx >= 0 && idx < len(m.profiles) {
				m.cursor = idx
				p := m.profiles[m.cursor]
				if m.arcadeMode {
					m.selFlashFrame = 0
					m.selFlashProfile = m.cursor
					m.state = StateSelectFlash
					return m, arcadeTickCmd()
				}
				return m, m.switchProfileCmd(p)
			}
		}

		// Legacy StateNoProfiles handler
		if m.state == StateNoProfiles {
			switch relY - 8 {
			case 0:
				m.LaunchLogin = true
				return m, tea.Quit
			case 1:
				m.state = StateWizardAddMore
				m.wizardStep = 0
				m.formFields = [6]string{}
				m.statusMsg = ""
			}
		}
	}
	_ = contentX
	return m, nil
}

func (m Model) updateWhatsNew(msg tea.Msg) (tea.Model, tea.Cmd) {
	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		default:
			// any key dismisses
			m.state = StateList
			m.splashSeen020 = true
			_ = m.store.SavePrefs(storage.Prefs{
				ColorTheme:    m.colorTheme,
				SplashSeen020: true,
			})
		}
	}
	if mm, ok := msg.(tea.MouseMsg); ok && mm.Action == tea.MouseActionPress && mm.Button == tea.MouseButtonLeft {
		m.state = StateList
		m.splashSeen020 = true
		_ = m.store.SavePrefs(storage.Prefs{
			ColorTheme:    m.colorTheme,
			SplashSeen020: true,
		})
	}
	return m, nil
}

func (m Model) updateWizard(msg tea.Msg) (tea.Model, tea.Cmd) {
	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if m.state == StateWizardWelcome {
				return m, tea.Quit
			}
		case "esc":
			switch m.state {
			case StateWizardDetect:
				m.state = StateWizardWelcome
			case StateWizardImport:
				m.state = StateWizardDetect
			case StateWizardAddMore:
				// if profiles exist we got here via 'a' from main TUI; go back
				m.state = StateList
			case StateWizardDone:
				m.state = StateList
			}
		case "enter":
			switch m.state {
			case StateWizardWelcome:
				m.state = StateWizardDetect
				return m, m.detectExistingConfigsCmd()
			case StateWizardDetect:
				m.importSelected = make([]bool, len(m.detectedProfiles))
				for i := range m.importSelected {
					m.importSelected[i] = true
				}
				if len(m.detectedProfiles) > 0 {
					m.state = StateWizardImport
				} else {
					m.wizardStep = 0
					m.state = StateWizardAddMore
				}
			case StateWizardImport:
				for i, p := range m.detectedProfiles {
					if i < len(m.importSelected) && m.importSelected[i] {
						_ = m.store.Add(p.Nickname, p.UserName, p.Email, p.SignKey, p.SSHKey, p.GHUser)
					}
				}
				profiles, _ := m.store.Load()
				m.profiles = profiles
				m.wizardStep = 0
				m.state = StateWizardAddMore
			case StateWizardAddMore:
				switch m.wizardStep {
				case 0: // OAuth
					m.statusMsg = "GitHub OAuth coming in v0.2.1 — use Add Manually for now"
				case 1: // Manual
					m.state = StateAdd
					m.formFields = [6]string{}
					m.formFocus = 0
					m.statusMsg = ""
				default: // Done (wizardStep==2 or unset)
					m.state = StateList
					m.tabIndex = 0
				}
			case StateWizardDone:
				m.state = StateList
				m.tabIndex = 0
			}
		case " ":
			// space toggles checkbox in import list; advances elsewhere
			if m.state == StateWizardImport && m.wizardStep < len(m.importSelected) {
				m.importSelected[m.wizardStep] = !m.importSelected[m.wizardStep]
			}
		case "up", "k":
			switch m.state {
			case StateWizardImport:
				if m.wizardStep > 0 {
					m.wizardStep--
				}
			case StateWizardAddMore:
				if m.wizardStep > 0 {
					m.wizardStep--
				}
			}
		case "down", "j":
			switch m.state {
			case StateWizardImport:
				if m.wizardStep < len(m.detectedProfiles)-1 {
					m.wizardStep++
				}
			case StateWizardAddMore:
				if m.wizardStep < 2 {
					m.wizardStep++
				}
			}
		}
	}
	return m, nil
}

type detectConfigsMsg struct {
	profiles []storage.Profile
}

func (m Model) detectExistingConfigsCmd() tea.Cmd {
	return func() tea.Msg {
		var found []storage.Profile
		// Read global gitconfig
		name, email := readGlobalGitConfig()
		if name != "" && email != "" {
			found = append(found, storage.Profile{
				Nickname: deriveNickname(email),
				UserName: name,
				Email:    email,
			})
		}
		return detectConfigsMsg{profiles: found}
	}
}

func readGlobalGitConfig() (name, email string) {
	cmdName := exec.Command("git", "config", "--global", "user.name")
	if out, err := cmdName.Output(); err == nil {
		name = strings.TrimSpace(string(out))
	}
	cmdEmail := exec.Command("git", "config", "--global", "user.email")
	if out, err := cmdEmail.Output(); err == nil {
		email = strings.TrimSpace(string(out))
	}
	return
}

func deriveNickname(email string) string {
	if at := strings.Index(email, "@"); at > 0 {
		return email[:at]
	}
	return email
}

func (m Model) updateNoProfiles(msg tea.Msg) (tea.Model, tea.Cmd) {
	km, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch strings.ToLower(km.String()) {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "a":
		m.state = StateAdd
		m.formFields = [6]string{}
		m.formFocus = 0
		m.statusMsg = ""
	case "l", "enter":
		m.LaunchLogin = true
		return m, tea.Quit
	}
	return m, nil
}
