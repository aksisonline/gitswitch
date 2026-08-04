package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aksisonline/gitswitch/internal/git"
	"github.com/aksisonline/gitswitch/internal/history"
	"github.com/aksisonline/gitswitch/internal/shell"
	"github.com/aksisonline/gitswitch/internal/storage"
	ver "github.com/aksisonline/gitswitch/internal/version"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

type switchDoneMsg struct {
	profile  *storage.Profile
	warnings []string
	err      error
}

// credentialHelperDoneMsg reports the result of toggling the HTTPS credential
// helper. enabled reflects git.IsCredentialHelperInstalled() after the toggle —
// not just "we ran install" — since another tool's helper can still shadow us.
type credentialHelperDoneMsg struct {
	enabled bool
	err     error
}

// ghWrapperDoneMsg reports the result of toggling the `gh` CLI wrapper.
// installed reflects shell.IsGHWrapperInstalled() after the toggle, same
// pattern as shellDoneMsg — it touches an rc file so needs a shell reload.
type ghWrapperDoneMsg struct {
	installed bool
	err       error
}

// pinDoneMsg reports the result of pinning or releasing this repo's identity.
// A nil profile means the pin was removed.
type pinDoneMsg struct {
	profile          *storage.Profile
	warnings         []string
	isolationEnabled bool
	err              error
}

type upgradeDoneMsg struct {
	err error
}

type editorDoneMsg struct {
	err error
}

type versionCheckMsg struct {
	latest string
	notes  string
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = ws.Width
		m.height = ws.Height
		return m, nil
	}
	if vc, ok := msg.(versionCheckMsg); ok {
		m.latestVersion = vc.latest
		m.releaseNotes = vc.notes
		m.updateAvailable = ver.IsUpdateAvailable(m.currentVersion, vc.latest)
		if m.updateAvailable && m.state == StateList {
			m.state = StateUpdatePrompt
		}
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
	if ed, ok := msg.(editorDoneMsg); ok {
		if ed.err != nil {
			m.statusMsg = fmt.Sprintf("editor: %v", ed.err)
			m.statusIsErr = true
		} else if profiles, err := m.store.Load(); err != nil {
			m.statusMsg = fmt.Sprintf("reload config: %v", err)
			m.statusIsErr = true
		} else {
			m.profiles = profiles
			m.active = git.DetectActive(profiles)
			if m.cursor >= len(m.profiles) && m.cursor > 0 {
				m.cursor = len(m.profiles) - 1
			}
			m.statusMsg = "config reloaded"
			m.statusIsErr = false
		}
		return m, tea.EnableMouseAllMotion
	}
	if ud, ok := msg.(upgradeDoneMsg); ok {
		if ud.err != nil {
			m.statusMsg = fmt.Sprintf("upgrade failed: %v", ud.err)
			m.statusIsErr = true
		} else {
			m.statusMsg = fmt.Sprintf("upgraded to %s — restart to apply", m.latestVersion)
			m.statusIsErr = false
			m.updateAvailable = false
		}
		m.state = StateList
		return m, nil
	}
	if cd, ok := msg.(credentialHelperDoneMsg); ok {
		if cd.err != nil {
			m.statusMsg = fmt.Sprintf("credential helper: %v", cd.err)
			m.statusIsErr = true
		} else {
			m.credentialHelperEnabled = cd.enabled
			if cd.enabled {
				m.statusMsg = "HTTPS credential helper enabled"
			} else {
				m.statusMsg = "HTTPS credential helper removed"
			}
			m.statusIsErr = false
		}
		return m, nil
	}
	if gd, ok := msg.(ghWrapperDoneMsg); ok {
		if gd.err != nil {
			m.statusMsg = fmt.Sprintf("gh wrapper: %v", gd.err)
			m.statusIsErr = true
		} else {
			m.ghWrapperEnabled = gd.installed
			if gd.installed {
				rc := shell.RCFile(shell.DetectShell())
				m.statusMsg = fmt.Sprintf("Session Isolation enabled — run: source %s", rc)
				m.PendingReloadCmd = fmt.Sprintf("source %s", rc)
			} else {
				m.statusMsg = "Session Isolation removed — restart your shell to apply"
				m.PendingReloadCmd = ""
			}
			m.statusIsErr = false
		}
		return m, nil
	}
	if sd, ok := msg.(shellDoneMsg); ok {
		if sd.err != nil {
			m.statusMsg = fmt.Sprintf("shell integration: %v", sd.err)
			m.statusIsErr = true
		} else {
			m.shellEnabled = sd.installed
			_ = m.savePrefs()
			if sd.installed {
				rc := shell.RCFile(shell.DetectShell())
				m.statusMsg = fmt.Sprintf("installed — run: source %s", rc)
				m.PendingReloadCmd = fmt.Sprintf("source %s", rc)
			} else {
				m.statusMsg = "shell integration removed — restart your shell to apply"
				m.PendingReloadCmd = ""
			}
			m.statusIsErr = false
		}
		m.state = StateList
		m.tabIndex = m.shellReturnTab
		return m, nil
	}
	if mm, ok := msg.(tea.MouseMsg); ok {
		return m.handleMouse(mm)
	}
	if m.state == StateShellConfirm {
		return m.updateShellConfirm(msg)
	}
	if m.aliasEditing {
		if km, ok := msg.(tea.KeyMsg); ok {
			switch km.String() {
			case "enter":
				val := strings.TrimSpace(m.aliasInput.Value())
				if val != "" {
					m.shellAlias = val
					_ = m.savePrefs()
					if m.shellEnabled && !m.shellAliasDisabled {
						return m, m.reinstallShellCmd()
					}
					m.statusMsg = "alias renamed"
					m.statusIsErr = false
				}
				m.aliasEditing = false
				return m, nil
			case "esc", "ctrl+c":
				m.aliasEditing = false
				return m, nil
			default:
				var cmd tea.Cmd
				m.aliasInput, cmd = m.aliasInput.Update(msg)
				return m, cmd
			}
		}
		return m, nil
	}
	switch m.state {
	case StateWhatsNew:
		if km, ok := msg.(tea.KeyMsg); ok {
			switch km.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				if m.whatsNewScroll > 0 {
					m.whatsNewScroll--
				}
			case "down", "j", " ":
				m.whatsNewScroll++
			default:
				// any other key dismisses
				ver.MarkVersionSeen(m.store.ConfigDir(), m.currentVersion)
				m.splashSeen020 = true
				_ = m.savePrefs()
				m.state = StateList
			}
		}
		return m, nil
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
	case StateUpdatePrompt:
		return m.updateUpdatePrompt(msg)
	case StateList:
		return m.updateList(msg)
	case StateAdd, StateEdit:
		return m.updateForm(msg)
	case StateDeleteConfirm:
		return m.updateDelete(msg)
	case StateTips:
		return m.updateTips(msg)
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
				m.utilityFocus = (m.utilityFocus + 3) % 4
			case 2:
				if m.settingsFocus > 0 {
					m.settingsFocus--
				}
				m.aliasEditing = false
			}
		case "down", "j":
			m.statusMsg = ""
			switch m.tabIndex {
			case 0:
				if m.cursor < len(m.profiles)-1 {
					m.cursor++
				}
			case 1:
				m.utilityFocus = (m.utilityFocus + 1) % 4
			case 2:
				if m.settingsFocus < 2 {
					m.settingsFocus++
				}
				m.aliasEditing = false
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
					m.openShellConfirm(!m.shellEnabled)
					return m, nil
				}
				if m.utilityFocus == 1 {
					return m, m.toggleGHWrapperCmd()
				}
				if m.utilityFocus == 2 {
					return m, m.toggleCredentialHelperCmd()
				}
				if m.utilityFocus == 3 {
					m.autoPinDisabled = !m.autoPinDisabled
					_ = m.savePrefs()
					if m.autoPinDisabled {
						m.statusMsg = "auto-pin disabled"
					} else {
						m.statusMsg = "auto-pin enabled"
					}
					m.statusIsErr = false
					return m, nil
				}
			case 2: // Settings
				if m.settingsFocus == 0 {
					return m, m.openConfigEditor()
				}
				if m.settingsFocus == 2 {
					m.shellAliasDisabled = !m.shellAliasDisabled
					_ = m.savePrefs()
					if m.shellEnabled {
						return m, m.reinstallShellCmd()
					}
					if !m.shellAliasDisabled {
						m.statusMsg = "alias enabled — install shell integration to apply"
					} else {
						m.statusMsg = "alias disabled"
					}
					m.statusIsErr = false
					return m, nil
				}
			}
		case "a":
			if m.tabIndex != 0 {
				break
			}
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
			if m.tabIndex == 2 && m.settingsFocus == 2 {
				m.aliasInput.SetValue(m.shellAlias)
				m.aliasInput.Focus()
				m.aliasEditing = true
				return m, textinput.Blink
			}
			if m.tabIndex != 0 {
				break
			}
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
				m.statusMsg = ""
				// keep the seed so a delete-confirm "n" can rebuild the form
				m.formFields = [6]string{p.Nickname, p.UserName, p.Email, p.SignKey, p.SSHKey, p.GHUser}
				return m, m.openProfileForm(true, m.formFields)
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
				if err := m.savePrefs(); err != nil {
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
				if err := m.savePrefs(); err != nil {
					m.statusMsg = "could not save theme"
				}
			}
		case "right":
			if !m.arcadeMode && m.tabIndex == 2 {
				m.colorTheme = (m.colorTheme + 1) % 12
				if err := m.savePrefs(); err != nil {
					m.statusMsg = "could not save theme"
				}
			}
		case "v":
			if m.tabIndex == 0 {
				m.showUsername = !m.showUsername
				_ = m.savePrefs()
				if m.showUsername {
					m.statusMsg = "showing GitHub usernames"
				} else {
					m.statusMsg = "showing emails"
				}
				m.statusIsErr = false
			}
		case "p":
			if m.tabIndex != 0 || len(m.profiles) == 0 {
				break
			}
			if m.repoKey == "" {
				m.statusMsg = "not inside a git repo — nothing to pin to"
				m.statusIsErr = true
				return m, nil
			}
			p := m.profiles[m.cursor]
			// One key both ways: pinning the profile this repo is already pinned to
			// releases it, like the other toggles in the app.
			if m.scope == git.ScopeRepo && m.scopeProfile != nil && m.scopeProfile.Nickname == p.Nickname {
				return m, m.unpinProfileCmd()
			}
			return m, m.pinProfileCmd(p)
		case "1", "2", "3":
			m.tabIndex = int(msg.String()[0] - '1')
			m.statusMsg = ""
		case "tab":
			m.tabIndex = (m.tabIndex + 1) % 3
			m.statusMsg = ""
		case "shift+tab":
			m.tabIndex = (m.tabIndex + 2) % 3
			m.statusMsg = ""
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
			m.refreshRepoScope()
			if len(msg.warnings) > 0 {
				m.statusMsg = fmt.Sprintf("switched to %s (warning: %s)", msg.profile.Nickname, msg.warnings[0])
			} else {
				m.statusMsg = fmt.Sprintf("switched to %s", msg.profile.Nickname)
			}
			m.statusIsErr = false
		}
	case pinDoneMsg:
		if msg.err != nil {
			m.statusMsg = fmt.Sprintf("error: %v", msg.err)
			m.statusIsErr = true
			break
		}
		if msg.isolationEnabled {
			m.ghWrapperEnabled = true
		}
		m.refreshRepoScope()
		m.statusIsErr = false
		if msg.profile == nil {
			m.statusMsg = "unpinned — this repo now uses the global identity"
			break
		}
		suffix := ""
		if len(msg.warnings) > 0 {
			suffix = fmt.Sprintf(" (warning: %s)", msg.warnings[0])
		} else if msg.isolationEnabled {
			suffix = " (Session Isolation turned on)"
		}
		m.statusMsg = fmt.Sprintf("pinned %s to this repo%s", msg.profile.Nickname, suffix)
	}
	return m, nil
}

func (m Model) updateForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Lazily build the form if we arrived here without one (e.g. arcade path).
	if m.form == nil {
		cmd := m.openProfileForm(m.state == StateEdit, m.formFields)
		return m, cmd
	}

	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "esc":
			return m.closeForm()
		case "ctrl+d":
			if m.state == StateEdit {
				m.form = nil
				m.formData = nil
				m.state = StateDeleteConfirm
				return m, nil
			}
		}
	}

	fm, cmd := m.form.Update(msg)
	if f, ok := fm.(*huh.Form); ok {
		m.form = f
	}

	switch m.form.State {
	case huh.StateCompleted:
		return m, m.submitForm()
	case huh.StateAborted:
		return m, tea.Quit
	}
	return m, cmd
}

// focusFormField moves the huh form's focus to the field with the given key
// by stepping NextField/PrevField exactly as far as needed — never past the
// target, since overshooting the last field submits the form.
func (m Model) focusFormField(key string) (tea.Model, tea.Cmd) {
	cur := m.form.GetFocusedField()
	if cur == nil {
		return m, nil
	}
	curIdx, tgtIdx := -1, -1
	for i, f := range profileFormFieldTitles {
		if f.key == cur.GetKey() {
			curIdx = i
		}
		if f.key == key {
			tgtIdx = i
		}
	}
	if curIdx < 0 || tgtIdx < 0 || curIdx == tgtIdx {
		return m, nil
	}
	var cmds []tea.Cmd
	if tgtIdx > curIdx {
		for i := curIdx; i < tgtIdx; i++ {
			cmds = append(cmds, m.form.NextField())
		}
	} else {
		for i := curIdx; i > tgtIdx; i-- {
			cmds = append(cmds, m.form.PrevField())
		}
	}
	return m, tea.Batch(cmds...)
}

// closeForm tears down the form and returns to the list (or arcade transition).
func (m Model) closeForm() (tea.Model, tea.Cmd) {
	m.form = nil
	m.formData = nil
	if m.arcadeMode {
		return m, m.startBackTransition()
	}
	m.state = StateList
	return m, nil
}

func (m *Model) submitForm() tea.Cmd {
	d := m.formData
	m.form = nil
	m.formData = nil

	nickname := strings.TrimSpace(d.nickname)
	userName := strings.TrimSpace(d.userName)
	email := strings.TrimSpace(d.email)
	signKey := strings.TrimSpace(d.signKey)
	sshKey := strings.TrimSpace(d.sshKey)
	ghUser := strings.TrimSpace(d.ghUser)

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
			if m.arcadeMode {
				m.addScore(500)
			}
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
			return m.confirmDelete()
		case "n", "N", "esc":
			m.state = StateEdit
		}
	}
	return m, nil
}

// confirmDelete removes the profile under the cursor. Shared by the "y" key
// and the delete-confirm dialog's mouse button.
func (m Model) confirmDelete() (tea.Model, tea.Cmd) {
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

// toggleCredentialHelperCmd installs or removes the HTTPS credential helper,
// mirroring `gitswitch shell`/`uninstall --https`. Toggle direction is
// decided by current state so it behaves like the other Utilities toggles: one
// key/click, no confirm dialog (unlike shell integration, this never touches an
// rc file or requires a shell reload).
func (m Model) toggleCredentialHelperCmd() tea.Cmd {
	wasEnabled := m.credentialHelperEnabled
	return func() tea.Msg {
		var err error
		if wasEnabled {
			err = git.UninstallCredentialHelper()
		} else {
			err = git.InstallCredentialHelper()
		}
		if err != nil {
			return credentialHelperDoneMsg{enabled: wasEnabled, err: err}
		}
		return credentialHelperDoneMsg{enabled: git.IsCredentialHelperInstalled()}
	}
}

// toggleGHWrapperCmd installs or removes the `gh` CLI wrapper function,
// mirroring `gitswitch shell`/`uninstall`'s gh-wrapper handling. Unlike the
// credential helper, this writes to the shell rc file, so — like shell
// integration — it needs a reload to take effect.
func (m Model) toggleGHWrapperCmd() tea.Cmd {
	wasEnabled := m.ghWrapperEnabled
	sh := shell.DetectShell()
	return func() tea.Msg {
		var err error
		if wasEnabled {
			_, err = shell.UninstallGHWrapper(sh)
		} else {
			_, err = shell.InstallGHWrapper(sh)
		}
		if err != nil {
			return ghWrapperDoneMsg{installed: wasEnabled, err: err}
		}
		return ghWrapperDoneMsg{installed: shell.IsGHWrapperInstalled(shell.RCFile(sh))}
	}
}

// pinProfileCmd writes a profile into this repo's own git config and records the
// pin, leaving the global identity alone — the same two writes as `gitswitch pin`.
func (m Model) pinProfileCmd(p storage.Profile) tea.Cmd {
	wasEnabled := m.ghWrapperEnabled
	return func() tea.Msg {
		if !wasEnabled {
			if _, err := shell.InstallGHWrapper(shell.DetectShell()); err != nil {
				return pinDoneMsg{err: err}
			}
		}
		cfg := git.New(false)
		if err := cfg.SetUser(p.UserName, p.Email); err != nil {
			return pinDoneMsg{err: err}
		}
		if err := cfg.SetSignKey(p.SignKey); err != nil {
			return pinDoneMsg{err: err}
		}
		if err := cfg.SetSSHKey(p.SSHKey); err != nil {
			return pinDoneMsg{err: err}
		}
		var warnings []string
		if w := git.SwitchGHUser(p.GHUser); w != "" {
			warnings = append(warnings, w)
		}
		if err := history.Pin(m.repoKey, p.Nickname); err != nil {
			return pinDoneMsg{err: err}
		}
		return pinDoneMsg{profile: &p, warnings: warnings, isolationEnabled: !wasEnabled}
	}
}

// unpinProfileCmd clears this repo's local identity so it follows the global
// profile again.
func (m Model) unpinProfileCmd() tea.Cmd {
	return func() tea.Msg {
		if err := git.New(false).ClearIdentity(); err != nil {
			return pinDoneMsg{err: err}
		}
		if err := history.Unpin(m.repoKey); err != nil {
			return pinDoneMsg{err: err}
		}
		return pinDoneMsg{}
	}
}

// addScore bumps the arcade score and persists a beaten high score.
func (m *Model) addScore(points int) {
	m.score += points
	if m.score > m.hiScore {
		m.hiScore = m.score
		_ = m.savePrefs()
	}
}

func (m Model) tickSelectFlash() (tea.Model, tea.Cmd) {
	m.selFlashFrame++
	if m.selFlashFrame >= 6 {
		m.state = StateList
		m.addScore(200)
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
		if m.transTarget == StateAdd || m.transTarget == StateEdit {
			return m, m.openProfileForm(m.transTarget == StateEdit, m.formFields)
		}
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

// profileFormFieldTitles maps huh field keys to their rendered Title text, in
// the exact order newProfileForm builds them. Used to hit-test clicks on the
// add/edit form, since huh itself has no mouse support.
var profileFormFieldTitles = [6]struct{ key, title string }{
	{"nickname", "Nickname"},
	{"userName", "Name"},
	{"email", "Email"},
	{"ghUser", "GitHub username"},
	{"sshKey", "SSH key path"},
	{"signKey", "Signing key"},
}

// formFieldAt returns the key of the field block under relY: the last field
// title seen at or above that row. A field's block (title, description,
// input, spacing) all belong to the nearest title above them.
func formFieldAt(bodyLines []string, relY int) (string, bool) {
	key, found := "", false
	for i := 0; i <= relY && i < len(bodyLines); i++ {
		clean := ansi.Strip(bodyLines[i])
		for _, f := range profileFormFieldTitles {
			if strings.Contains(clean, f.title) {
				key, found = f.key, true
			}
		}
	}
	return key, found
}

// profileRowStart returns the relY of the first profile row on the Accounts
// tab, mirroring the layout built in viewAccountsTab.
func (m Model) profileRowStart(arcadeOff int) int {
	compact := m.height > 0 && m.height < 12+len(m.profiles)
	if compact {
		return 6 + arcadeOff
	}
	return 7 + arcadeOff
}

// itemBoxAt maps a relY to a Utilities/Settings item-box index. Each
// renderItemBox occupies exactly 4 rows, the first box starting right
// under the tab strip. maxIdx is the highest valid index for the tab being
// hit-tested (3 for Utilities' 4 boxes, 2 for Settings' 3).
func itemBoxAt(relY, arcadeOff, maxIdx int) (int, bool) {
	start := 4 + arcadeOff
	if relY < start {
		return 0, false
	}
	idx := (relY - start) / 4
	if idx > maxIdx {
		return 0, false
	}
	return idx, true
}

// handleMouse hit-tests against the rendered body (bodyView) instead of
// hardcoded row offsets, so clicks stay correct across compact layouts,
// arcade mode's score line, and variable-height screens. Hover (motion)
// moves the cursor/focus; a left press activates.
func (m Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	hover := msg.Action == tea.MouseActionMotion
	if !hover && msg.Action != tea.MouseActionPress {
		return m, nil
	}

	pw := m.panelWidth()
	body := m.bodyView(pw)
	bodyTop := (m.height - lipgloss.Height(body)) / 2
	if bodyTop < 0 {
		bodyTop = 0
	}
	bodyLeft := (m.width - lipgloss.Width(body)) / 2
	if bodyLeft < 0 {
		bodyLeft = 0
	}
	relY := msg.Y - bodyTop
	lineX := msg.X - bodyLeft

	bodyLines := strings.Split(body, "\n")
	var ln string // ANSI-stripped text of the row under the pointer
	if relY >= 0 && relY < len(bodyLines) {
		ln = ansi.Strip(bodyLines[relY])
	}
	// hitLabel: pointer is on this exact label. onRow: label anywhere on the row.
	hitLabel := func(label string) bool {
		i := strings.Index(ln, label)
		if i < 0 {
			return false
		}
		cellX := lipgloss.Width(ln[:i])
		return lineX >= cellX-1 && lineX <= cellX+lipgloss.Width(label)
	}
	onRow := func(labels ...string) bool {
		for _, l := range labels {
			if strings.Contains(ln, l) {
				return true
			}
		}
		return false
	}

	arcadeOff := 0
	if m.arcadeMode {
		arcadeOff = 1 // score line sits above the panel border
	}

	// Wheel scrolling
	if !hover && (msg.Button == tea.MouseButtonWheelUp || msg.Button == tea.MouseButtonWheelDown) {
		up := msg.Button == tea.MouseButtonWheelUp
		switch m.state {
		case StateList:
			switch m.tabIndex {
			case 0:
				if up && m.cursor > 0 {
					m.cursor--
				} else if !up && m.cursor < len(m.profiles)-1 {
					m.cursor++
				}
				m.statusMsg = ""
			case 2:
				if up && m.settingsFocus > 0 {
					m.settingsFocus--
				} else if !up && m.settingsFocus < 2 {
					m.settingsFocus++
				}
			}
		case StateWizardAddMore:
			if up && m.wizardStep > 0 {
				m.wizardStep--
			} else if !up && m.wizardStep < 2 {
				m.wizardStep++
			}
		case StateWizardImport:
			if up && m.wizardStep > 0 {
				m.wizardStep--
			} else if !up && m.wizardStep < len(m.detectedProfiles)-1 {
				m.wizardStep++
			}
		}
		return m, nil
	}

	// Hover only moves focus — never triggers actions.
	if hover {
		if m.state != StateList {
			return m, nil
		}
		switch m.tabIndex {
		case 0:
			if idx := relY - m.profileRowStart(arcadeOff); idx >= 0 && idx < len(m.profiles) {
				m.cursor = idx
			}
		case 1:
			if idx, ok := itemBoxAt(relY, arcadeOff, 3); ok {
				m.utilityFocus = idx
			}
		case 2:
			if idx, ok := itemBoxAt(relY, arcadeOff, 2); ok {
				m.settingsFocus = idx
			}
		}
		return m, nil
	}

	if msg.Button != tea.MouseButtonLeft {
		return m, nil
	}

	// A click anywhere ends alias editing (settings click below may re-enter).
	m.aliasEditing = false

	switch m.state {
	case StateIntro:
		m.state = StateList
		return m, nil

	case StateWhatsNew:
		m.state = StateList
		m.splashSeen020 = true
		_ = m.savePrefs()
		return m, nil

	case StateTips:
		if m.arcadeMode {
			return m, m.startBackTransition()
		}
		m.state = StateList
		return m, nil

	case StateUpdatePrompt:
		if hitLabel("Yes, update now") || hitLabel("INSERT COIN") {
			cmd, err := ver.UpgradeCommand(m.latestVersion)
			if err != nil {
				m.statusMsg = fmt.Sprintf("upgrade error: %v", err)
				m.statusIsErr = true
				m.state = StateList
				return m, nil
			}
			return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
				return upgradeDoneMsg{err: err}
			})
		}
		if hitLabel("Skip") || hitLabel("CONTINUE") {
			m.state = StateList
		}
		return m, nil

	case StateShellConfirm:
		// Precise click on a button only — same spirit as the y/n-only keys.
		if hitLabel("Yes, install it") || hitLabel("Yes, remove it") {
			return m, m.runShellAction(m.pendingShellInstall)
		}
		if hitLabel("Cancel") {
			m.state = StateList
			m.tabIndex = m.shellReturnTab
		}
		return m, nil

	case StateWizardWelcome:
		if onRow("Start Setup", "INSERT COIN") {
			m.state = StateWizardDetect
			return m, m.detectExistingConfigsCmd()
		}
		if onRow("Skip for now", "SKIP SETUP") {
			m.state = StateList
			m.tabIndex = 0
		}
		return m, nil

	case StateWizardDetect:
		// click continues, mirroring enter
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
		return m, nil

	case StateWizardImport:
		if hitLabel("Import selected") {
			m.importSelectedProfiles()
			return m, nil
		}
		if hitLabel("Skip") {
			m.wizardStep = 0
			m.state = StateWizardAddMore
			return m, nil
		}
		for i, p := range m.detectedProfiles {
			if strings.Contains(ln, p.Nickname) {
				m.wizardStep = i
				if i < len(m.importSelected) {
					m.importSelected[i] = !m.importSelected[i]
				}
				break
			}
		}
		return m, nil

	case StateWizardAddMore:
		switch {
		case onRow("Log in with GitHub", "LOG IN WITH GITHUB"):
			m.wizardStep = 0
			m.LaunchOAuth = true
			return m, tea.Quit
		case onRow("Add manually", "ADD MANUALLY"):
			m.wizardStep = 1
			m.state = StateAdd
			m.statusMsg = ""
			m.formFields = [6]string{}
			return m, m.openProfileForm(false, [6]string{})
		case onRow("Done, open", "DONE, OPEN"):
			m.wizardStep = 2
			m.state = StateList
			m.tabIndex = 0
		}
		return m, nil

	case StateWizardDone:
		m.state = StateList
		m.tabIndex = 0
		return m, nil

	case StateDeleteConfirm:
		if onRow("confirm delete") {
			return m.confirmDelete()
		}
		if onRow("cancel") {
			m.state = StateEdit
		}
		return m, nil

	case StateAdd, StateEdit:
		if m.form == nil {
			return m, nil
		}
		key, ok := formFieldAt(bodyLines, relY)
		if !ok {
			return m, nil
		}
		return m.focusFormField(key)

	case StateList:
		// Tab strip — the only row containing both later tab names.
		names := m.tabNames()
		if onRow(names[1]) && onRow(names[2]) {
			for i := 2; i >= 0; i-- {
				if pos := strings.Index(ln, names[i]); pos >= 0 && lineX >= pos-2 {
					m.tabIndex = i
					m.statusMsg = ""
					break
				}
			}
			return m, nil
		}

		switch m.tabIndex {
		case 0:
			if idx := relY - m.profileRowStart(arcadeOff); idx >= 0 && idx < len(m.profiles) {
				m.cursor = idx
				if m.arcadeMode {
					m.selFlashFrame = 0
					m.selFlashProfile = m.cursor
					m.state = StateSelectFlash
					return m, arcadeTickCmd()
				}
				return m, m.switchProfileCmd(m.profiles[idx])
			}
		case 1:
			if idx, ok := itemBoxAt(relY, arcadeOff, 3); ok {
				m.utilityFocus = idx
				switch idx {
				case 0:
					m.openShellConfirm(!m.shellEnabled)
				case 1:
					return m, m.toggleGHWrapperCmd()
				case 2:
					return m, m.toggleCredentialHelperCmd()
				case 3:
					m.autoPinDisabled = !m.autoPinDisabled
					_ = m.savePrefs()
					if m.autoPinDisabled {
						m.statusMsg = "auto-pin disabled"
					} else {
						m.statusMsg = "auto-pin enabled"
					}
					m.statusIsErr = false
				}
			}
		case 2:
			if idx, ok := itemBoxAt(relY, arcadeOff, 2); ok {
				m.settingsFocus = idx
				switch idx {
				case 0: // config location — only the edit chip launches $EDITOR
					if hitLabel("[✎ edit]") {
						return m, m.openConfigEditor()
					}
				case 1: // theme box — left half prev, right half next
					if !m.arcadeMode {
						if lineX > pw/2 {
							m.colorTheme = (m.colorTheme + 1) % 12
						} else {
							m.colorTheme = (m.colorTheme + 11) % 12
						}
						_ = m.savePrefs()
					}
				case 2: // alias box — rename chip edits, body toggles
					if hitLabel("[✎ rename]") {
						m.aliasInput.SetValue(m.shellAlias)
						m.aliasInput.Focus()
						m.aliasEditing = true
						return m, textinput.Blink
					}
					m.shellAliasDisabled = !m.shellAliasDisabled
					_ = m.savePrefs()
					if m.shellEnabled {
						return m, m.reinstallShellCmd()
					}
					if !m.shellAliasDisabled {
						m.statusMsg = "alias enabled — install shell integration to apply"
					} else {
						m.statusMsg = "alias disabled"
					}
					m.statusIsErr = false
				}
			}
		}
		return m, nil
	}
	return m, nil
}

// importSelectedProfiles adds the checked detected profiles and advances the
// wizard to the add-more step. Shared by the enter key and mouse click.
func (m *Model) importSelectedProfiles() {
	for i, p := range m.detectedProfiles {
		if i < len(m.importSelected) && m.importSelected[i] {
			_ = m.store.Add(p.Nickname, p.UserName, p.Email, p.SignKey, p.SSHKey, p.GHUser)
		}
	}
	profiles, _ := m.store.Load()
	m.profiles = profiles
	m.wizardStep = 0
	m.state = StateWizardAddMore
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
				m.importSelectedProfiles()
			case StateWizardAddMore:
				switch m.wizardStep {
				case 0: // OAuth
					m.LaunchOAuth = true
					return m, tea.Quit
				case 1: // Manual
					m.state = StateAdd
					m.statusMsg = ""
					return m, m.openProfileForm(false, [6]string{})
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
		return detectConfigsMsg{profiles: detectExistingProfiles()}
	}
}

// mergeProfileFields copies any field set in src into dst wherever dst is
// still empty. Shared by every merge tier in mergeIntoFound.
func mergeProfileFields(dst *storage.Profile, src storage.Profile) {
	if dst.UserName == "" {
		dst.UserName = src.UserName
	}
	if dst.Email == "" {
		dst.Email = src.Email
	}
	if dst.SignKey == "" {
		dst.SignKey = src.SignKey
	}
	if dst.SSHKey == "" {
		dst.SSHKey = src.SSHKey
	}
	if dst.GHUser == "" {
		dst.GHUser = src.GHUser
	}
}

// mergeIntoFound is detectExistingProfiles' candidate-merge logic, extracted
// so it's unit-testable without shelling out to git/gh. verifiedEmails is
// the GitHub-API cross-reference tier — pass nil for candidates that don't
// have one (e.g. the git-config candidate itself).
//
// Tiers, in order — each is additive, none can un-fire an earlier one:
//  1. exact Nickname collision (existing behavior)
//  2. a verifiedEmails entry matches an already-found candidate's Email —
//     this is what lets a gh account and a git-config profile that are
//     really the same person merge, instead of showing as two candidates
//     just because their nickname/email happened not to coincide
//  3. literal Email string coincidence via seenEmail (existing behavior)
func mergeIntoFound(found []storage.Profile, seenEmail map[string]bool, p storage.Profile, verifiedEmails []string) []storage.Profile {
	if p.Nickname == "" {
		return found
	}

	for i := range found {
		if found[i].Nickname == p.Nickname {
			mergeProfileFields(&found[i], p)
			return found
		}
	}

	for _, ve := range verifiedEmails {
		for i := range found {
			if found[i].Email != "" && strings.EqualFold(found[i].Email, ve) {
				mergeProfileFields(&found[i], p)
				return found
			}
		}
	}

	if p.Email != "" && seenEmail[p.Email] && p.GHUser == "" {
		return found
	}
	if p.Email != "" {
		seenEmail[p.Email] = true
	}
	return append(found, p)
}

// detectExistingProfiles scans git config, gh accounts, and ~/.ssh for
// identities the user can import. Used by the onboarding wizard.
func detectExistingProfiles() []storage.Profile {
	var found []storage.Profile
	seenEmail := map[string]bool{}

	add := func(p storage.Profile, verifiedEmails ...string) {
		found = mergeIntoFound(found, seenEmail, p, verifiedEmails)
	}

	cfg := git.New(true)
	name, email, _ := cfg.GetUser()
	if name != "" && email != "" {
		add(storage.Profile{
			Nickname: deriveNickname(email),
			UserName: name,
			Email:    email,
			SignKey:  cfg.GetSignKey(),
			SSHKey:   cfg.GetSSHKey(),
			GHUser:   git.GetGHUser(),
		})
	}

	// Best-effort cross-reference: for each gh-CLI account, fetch its
	// verified GitHub email(s) concurrently (one goroutine per account,
	// each writing only its own slice index — no mutex needed) so a match
	// against the git-config candidate above merges them into one profile
	// instead of surfacing as two. Silently empty on any failure (gh
	// missing, insufficient token scope, no network) — falls through to
	// today's nickname/email-string behavior with no regression.
	ghAccounts := git.ListGHUsers()
	verified := make([][]string, len(ghAccounts))
	var wg sync.WaitGroup
	for i, acct := range ghAccounts {
		wg.Add(1)
		go func(i int, acct git.GHAccount) {
			defer wg.Done()
			verified[i] = git.VerifiedEmailsFor(acct.Login, acct.Host)
		}(i, acct)
	}
	wg.Wait()

	for i, acct := range ghAccounts {
		nick := acct.Login
		if nick == "" {
			continue
		}
		add(storage.Profile{
			Nickname: nick,
			UserName: nick,
			Email:    nick + "@users.noreply.github.com",
			GHUser:   acct.Login,
		}, verified[i]...)
	}

	// Attach first SSH key found to the first profile if none set yet;
	// otherwise surface a profile-ready hint via a dedicated entry only
	// when we have no profiles at all.
	keys := git.ListSSHPrivateKeys()
	if len(keys) > 0 && len(found) > 0 {
		for i := range found {
			if found[i].SSHKey == "" {
				// Prefer tilde form for display/portability.
				k := keys[0]
				if home, err := os.UserHomeDir(); err == nil && strings.HasPrefix(k, home+string(os.PathSeparator)) {
					k = "~/" + strings.TrimPrefix(k, home+string(os.PathSeparator))
				}
				found[i].SSHKey = k
				break
			}
		}
	}

	return found
}

func deriveNickname(email string) string {
	if at := strings.Index(email, "@"); at > 0 {
		return email[:at]
	}
	return email
}

func (m Model) openConfigEditor() tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		editor = "vi"
	}
	configPath := filepath.Join(m.store.ConfigDir(), "config.yaml")
	cmd := exec.Command(editor, configPath)
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return editorDoneMsg{err: err}
	})
}

func (m Model) updateUpdatePrompt(msg tea.Msg) (tea.Model, tea.Cmd) {
	km, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch km.String() {
	case "y", "Y", "enter":
		cmd, err := ver.UpgradeCommand(m.latestVersion)
		if err != nil {
			m.statusMsg = fmt.Sprintf("upgrade error: %v", err)
			m.statusIsErr = true
			m.state = StateList
			return m, nil
		}
		return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
			return upgradeDoneMsg{err: err}
		})
	case "n", "N", "esc", "q":
		m.state = StateList
		return m, nil
	case "ctrl+c":
		return m, tea.Quit
	}
	return m, nil
}
