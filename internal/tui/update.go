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

// guardDoneMsg reports the result of toggling the global pre-commit guard hook.
type guardDoneMsg struct {
	enabled bool
	err     error
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

// Update is the bubbletea entry point. It handles the frame ticker, then
// delegates to update and starts the ticker again if the result is animating.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case frameMsg:
		return m.onFrame()
	case tea.KeyMsg:
		m.touched = true
	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress {
			m.touched = true
		}
	}
	next, cmd := m.update(msg)
	nm, ok := next.(Model)
	if !ok {
		return next, cmd
	}
	return nm, tea.Batch(cmd, nm.animate())
}

func (m Model) update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = ws.Width
		m.height = ws.Height
		if m.form != nil {
			// Keep the form inside the panel as the window changes.
			w, h := m.formSize()
			m.form = m.form.WithWidth(max(20, w))
			if h > 0 {
				m.form = m.form.WithHeight(h)
			}
		}
		return m, nil
	}
	if vc, ok := msg.(versionCheckMsg); ok {
		m.latestVersion = vc.latest
		m.releaseNotes = vc.notes
		m.updateAvailable = ver.IsUpdateAvailable(m.currentVersion, vc.latest)
		// Offer the update once, and only if the user hasn't started doing
		// something — the check lands a moment after launch, and swapping the
		// screen under a keypress turned "enter to switch" into "enter to
		// upgrade". After that the status line carries the reminder.
		if m.updateAvailable && m.state == StateList && !m.touched && m.skippedUpdate != vc.latest {
			m.state = StateUpdatePrompt
		}
		return m, nil
	}
	if dcm, ok := msg.(detectConfigsMsg); ok {
		m.detectedProfiles = dcm.profiles
		m.scanReady = true
		if !m.scanning {
			m.finishScan()
		}
		return m, nil
	}
	if sw, ok := msg.(switchDoneMsg); ok {
		return m.onSwitchDone(sw), nil
	}
	if pd, ok := msg.(pinDoneMsg); ok {
		return m.onPinDone(pd), nil
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
			m.statusMsg = fmt.Sprintf("upgraded to %s — restart gitswitch to use it", m.latestVersion)
			m.statusIsErr = false
			m.updateAvailable = false
		}
		m.state = StateList
		return m, nil
	}
	if cd, ok := msg.(credentialHelperDoneMsg); ok {
		m.clearBusy()
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
		m.clearBusy()
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
	if grd, ok := msg.(guardDoneMsg); ok {
		m.clearBusy()
		if grd.err != nil {
			m.statusMsg = fmt.Sprintf("pre-commit guard: %v", grd.err)
			m.statusIsErr = true
		} else {
			m.guardEnabled = grd.enabled
			if grd.enabled {
				m.statusMsg = "pre-commit guard enabled"
			} else {
				m.statusMsg = "pre-commit guard disabled"
			}
			m.statusIsErr = false
		}
		return m, nil
	}
	if sd, ok := msg.(shellDoneMsg); ok {
		m.clearBusy()
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
						m.aliasEditing = false
						m.setBusy("updating shell integration…")
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
			maxScroll := m.maxWhatsNewScroll()
			switch km.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				if m.whatsNewScroll > 0 {
					m.whatsNewScroll--
				}
			case "down", "j", " ":
				if m.whatsNewScroll < maxScroll {
					m.whatsNewScroll++
				}
			case "pgup":
				m.whatsNewScroll -= notesViewportHeight(m.height)
				if m.whatsNewScroll < 0 {
					m.whatsNewScroll = 0
				}
			case "pgdown":
				m.whatsNewScroll += notesViewportHeight(m.height)
				if m.whatsNewScroll > maxScroll {
					m.whatsNewScroll = maxScroll
				}
			default:
				// any other key dismisses
				m.dismissWhatsNew()
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

// finishScan leaves the setup wizard's Scan step once the scan has returned:
// to Import when something was found, otherwise straight to adding an account.
func (m *Model) finishScan() {
	m.scanning, m.scanStep, m.scanHold = false, 0, 0
	if m.state != StateWizardDetect || !m.scanReady {
		return
	}
	m.importSelected = make([]bool, len(m.detectedProfiles))
	for i := range m.importSelected {
		m.importSelected[i] = true
	}
	m.wizardStep = 0
	if len(m.detectedProfiles) > 0 {
		m.state = StateWizardImport
	} else {
		m.state = StateWizardAddMore
	}
}

func (m Model) onSwitchDone(msg switchDoneMsg) Model {
	m.clearBusy()
	if msg.err != nil {
		m.statusMsg = fmt.Sprintf("error: %v", msg.err)
		m.statusIsErr = true
		return m
	}
	profiles, err := m.store.Load()
	if err != nil {
		m.statusMsg = fmt.Sprintf("error: %v", err)
		m.statusIsErr = true
		return m
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
	m.flashNick, m.flashLeft = msg.profile.Nickname, flashFrames
	return m
}

func (m Model) onPinDone(msg pinDoneMsg) Model {
	m.clearBusy()
	if msg.err != nil {
		m.statusMsg = fmt.Sprintf("error: %v", msg.err)
		m.statusIsErr = true
		return m
	}
	if msg.isolationEnabled {
		m.ghWrapperEnabled = true
	}
	m.refreshRepoScope()
	m.statusIsErr = false
	if msg.profile == nil {
		m.statusMsg = "unpinned — this repo now uses the global identity"
		return m
	}
	suffix := ""
	if len(msg.warnings) > 0 {
		suffix = fmt.Sprintf(" (warning: %s)", msg.warnings[0])
	} else if msg.isolationEnabled {
		suffix = " (Session Isolation turned on)"
	}
	m.statusMsg = fmt.Sprintf("pinned %s to this repo%s", msg.profile.Nickname, suffix)
	return m
}

// switchTo starts switching to profile p, with a spinner until it lands.
func (m Model) switchTo(p storage.Profile) (tea.Model, tea.Cmd) {
	m.setBusy("switching to " + p.Nickname + "…")
	return m, m.switchProfileCmd(p)
}

// toggleUtility flips the Utilities item at idx. Shared by enter and click.
func (m Model) toggleUtility(idx int) (tea.Model, tea.Cmd) {
	m.utilityFocus = idx
	m.statusMsg = ""
	switch idx {
	case 0:
		m.openShellConfirm(!m.shellEnabled)
		return m, nil
	case 1:
		m.busyItem = idx + 1
		return m, m.toggleGHWrapperCmd()
	case 2:
		m.busyItem = idx + 1
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
	case 4:
		m.busyItem = idx + 1
		return m, m.toggleGuardCmd()
	}
	return m, nil
}

// cycleTheme moves the palette by delta and saves it.
func (m *Model) cycleTheme(delta int) {
	n := len(normalThemes)
	m.colorTheme = ((m.colorTheme+delta)%n + n) % n
	if err := m.savePrefs(); err != nil {
		m.statusMsg = fmt.Sprintf("theme: %s — could not save: %v", themeNames[m.colorTheme], err)
		m.statusIsErr = true
		return
	}
	m.statusMsg = fmt.Sprintf("theme: %s", themeNames[m.colorTheme])
	m.statusIsErr = false
}

// toggleAlias turns the shell alias on or off. Shared by enter and click.
func (m Model) toggleAlias() (tea.Model, tea.Cmd) {
	m.shellAliasDisabled = !m.shellAliasDisabled
	_ = m.savePrefs()
	if m.shellEnabled {
		m.setBusy("updating shell integration…")
		return m, m.reinstallShellCmd()
	}
	if !m.shellAliasDisabled {
		m.statusMsg = "alias enabled — turn on Shell Integration to apply"
	} else {
		m.statusMsg = "alias disabled"
	}
	m.statusIsErr = false
	return m, nil
}

// startAliasEdit opens the inline alias editor.
func (m Model) startAliasEdit() (tea.Model, tea.Cmd) {
	m.settingsFocus = 2
	m.aliasInput.SetValue(m.shellAlias)
	m.aliasInput.Focus()
	m.aliasEditing = true
	return m, textinput.Blink
}

// openAddAccount opens the "add an account" menu from the main list.
func (m *Model) openAddAccount() {
	m.wizardStep = 0
	m.state = StateWizardAddMore
	m.statusMsg = ""
}

func (m Model) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	km, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch km.String() {
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
			m.utilityFocus = (m.utilityFocus + len(utilItemSpecs) - 1) % len(utilItemSpecs)
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
			m.utilityFocus = (m.utilityFocus + 1) % len(utilItemSpecs)
		case 2:
			if m.settingsFocus < 2 {
				m.settingsFocus++
			}
		}
	case "home", "g":
		if m.tabIndex == 0 {
			m.cursor = 0
		}
	case "end", "G":
		if m.tabIndex == 0 && len(m.profiles) > 0 {
			m.cursor = len(m.profiles) - 1
		}
	case "enter", " ":
		switch m.tabIndex {
		case 0: // Accounts — switch profile
			if len(m.profiles) == 0 {
				if km.String() == "enter" {
					m.openAddAccount()
				}
				break
			}
			if m.arcadeMode {
				m.selFlashFrame = 0
				m.selFlashProfile = m.cursor
				m.state = StateSelectFlash
				return m, arcadeTickCmd()
			}
			return m.switchTo(m.profiles[m.cursor])
		case 1: // Utilities — toggle focused item
			return m.toggleUtility(m.utilityFocus)
		case 2: // Settings
			switch m.settingsFocus {
			case 0:
				return m, m.openConfigEditor()
			case 1:
				if !m.arcadeMode {
					m.cycleTheme(1)
				}
			case 2:
				return m.toggleAlias()
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
		m.openAddAccount()
	case "e":
		if m.tabIndex == 2 && m.settingsFocus == 2 {
			return m.startAliasEdit()
		}
		if m.tabIndex == 2 && m.settingsFocus == 0 {
			return m, m.openConfigEditor()
		}
		if m.tabIndex != 0 || len(m.profiles) == 0 {
			break
		}
		p := m.profiles[m.cursor]
		m.editingNick = p.Nickname
		m.statusMsg = ""
		// keep the seed so a delete-confirm "n" can rebuild the form
		m.formFields = [6]string{p.Nickname, p.UserName, p.Email, p.SignKey, p.SSHKey, p.GHUser}
		if m.arcadeMode {
			m.formFocus = 0
			m.transTarget = StateEdit
			m.transFrame = 0
			m.state = StateTransition
			return m, arcadeTickCmd()
		}
		m.state = StateEdit
		return m, m.openProfileForm(true, m.formFields)
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
			m.cycleTheme(1)
		}
	case "left", "h":
		if !m.arcadeMode && m.tabIndex == 2 && m.settingsFocus == 1 {
			m.cycleTheme(-1)
		}
	case "right", "l":
		if !m.arcadeMode && m.tabIndex == 2 && m.settingsFocus == 1 {
			m.cycleTheme(1)
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
			m.setBusy("unpinning…")
			return m, m.unpinProfileCmd()
		}
		m.setBusy("pinning " + p.Nickname + " to this repo…")
		return m, m.pinProfileCmd(p)
	case "1", "2", "3":
		m.setTab(int(km.String()[0] - '1'))
	case "tab":
		m.setTab((m.tabIndex + 1) % 3)
	case "shift+tab":
		m.setTab((m.tabIndex + 2) % 3)
	case "u":
		if m.updateAvailable {
			m.state = StateUpdatePrompt
		}
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
	m.state = m.afterFormState()
	return m, nil
}

// afterFormState is where the add/edit form returns to: the add-account menu
// while first-run setup is still in progress, otherwise the main list.
func (m *Model) afterFormState() State {
	if m.onboarding {
		m.wizardStep = 2 // "Done" — the likely next step after adding one
		return StateWizardAddMore
	}
	return StateList
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
	m.state = m.afterFormState()
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
			return m, m.openProfileForm(true, m.formFields)
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

// AfterSwitch, when set, runs after a successful global profile switch from the
// TUI. The CLI wires it to the Xcode account sync so picker switches behave like
// `gitswitch switch`. Best effort: it can't fail or block the switch.
var AfterSwitch func(p storage.Profile)

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
		if AfterSwitch != nil {
			AfterSwitch(p)
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

// toggleGuardCmd installs or removes the global pre-commit identity guard hook.
func (m Model) toggleGuardCmd() tea.Cmd {
	wasEnabled := m.guardEnabled
	return func() tea.Msg {
		var err error
		if wasEnabled {
			err = git.UninstallGuardHook()
		} else {
			err = git.InstallGuardHook()
		}
		if err != nil {
			return guardDoneMsg{enabled: wasEnabled, err: err}
		}
		return guardDoneMsg{enabled: git.IsGuardHookInstalled()}
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
			return m.switchTo(m.profiles[m.selFlashProfile])
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
	numSlots := m.introSlots()
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

// handleMouse hit-tests against the rendered screen: rows are resolved through
// the zones each view recorded while drawing, and buttons by their label text,
// so clicks stay correct at every density, in arcade mode, and at any size.
// Hover (motion) moves the cursor/focus; a left press activates.
func (m Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	hover := msg.Action == tea.MouseActionMotion
	if !hover && msg.Action != tea.MouseActionPress {
		return m, nil
	}

	scr := m.screen()
	top, left := m.origin(scr)
	relY := msg.Y - top
	lineX := msg.X - left
	ln := scr.plain(relY) // ANSI-stripped text of the row under the pointer
	z, onZone := scr.zoneAt(relY)
	if lineX < 0 || lineX >= scr.width() {
		onZone = false
	}

	// hitLabel: pointer is on this exact label. onRow: label anywhere on the row.
	hitLabel := func(label string) bool {
		i := strings.Index(ln, label)
		if i < 0 {
			return false
		}
		cellX := lipgloss.Width(ln[:i])
		return lineX >= cellX-2 && lineX < cellX+lipgloss.Width(label)+2
	}
	onRow := func(labels ...string) bool {
		for _, l := range labels {
			if strings.Contains(ln, l) {
				return true
			}
		}
		return false
	}

	// Wheel scrolling
	if !hover && (msg.Button == tea.MouseButtonWheelUp || msg.Button == tea.MouseButtonWheelDown) {
		step := 1
		if msg.Button == tea.MouseButtonWheelUp {
			step = -1
		}
		clamp := func(v, hi int) int { return min(max(v+step, 0), max(hi, 0)) }
		switch m.state {
		case StateList:
			switch m.tabIndex {
			case 0:
				m.cursor = clamp(m.cursor, len(m.profiles)-1)
				m.statusMsg = ""
			case 1:
				m.utilityFocus = clamp(m.utilityFocus, len(utilItemSpecs)-1)
			case 2:
				if !m.aliasEditing {
					m.settingsFocus = clamp(m.settingsFocus, 2)
				}
			}
		case StateWizardWelcome:
			m.wizardStep = clamp(m.wizardStep, 1)
		case StateWizardAddMore:
			m.wizardStep = clamp(m.wizardStep, 2)
		case StateWizardImport:
			m.wizardStep = clamp(m.wizardStep, len(m.detectedProfiles)-1)
		case StateWhatsNew:
			m.whatsNewScroll = clamp(m.whatsNewScroll, m.maxWhatsNewScroll())
		}
		return m, nil
	}

	// Hover only moves focus — never triggers actions.
	if hover {
		if !onZone || m.aliasEditing {
			return m, nil
		}
		switch {
		case m.state == StateList && z.kind == zProfile:
			m.cursor = z.idx
		case m.state == StateList && z.kind == zItem && m.tabIndex == 1:
			m.utilityFocus = z.idx
		case m.state == StateList && z.kind == zItem && m.tabIndex == 2:
			m.settingsFocus = z.idx
		case z.kind == zChoice && (m.state == StateWizardWelcome || m.state == StateWizardAddMore):
			m.wizardStep = z.idx
		case z.kind == zImport && m.state == StateWizardImport:
			m.wizardStep = z.idx
		}
		return m, nil
	}

	if msg.Button != tea.MouseButtonLeft {
		return m, nil
	}

	// A click anywhere ends alias editing (a settings click below may re-enter).
	m.aliasEditing = false

	switch m.state {
	case StateIntro:
		m.state = StateList
		return m, nil

	case StateWhatsNew:
		m.dismissWhatsNew()
		return m, nil

	case StateTips:
		if m.arcadeMode {
			return m, m.startBackTransition()
		}
		m.state = StateList
		return m, nil

	case StateUpdatePrompt:
		if hitLabel(updateYes) || hitLabel(updateYesArcade) {
			return m.runUpgrade()
		}
		if hitLabel(updateNo) || hitLabel(updateNoArcade) {
			m.skipUpdate()
		}
		return m, nil

	case StateShellConfirm:
		// Precise click on a button only — same spirit as the y/n-only keys.
		if hitLabel(shellYesInstall) || hitLabel(shellYesRemove) {
			m.setBusy("updating shell integration…")
			return m, m.runShellAction(m.pendingShellInstall)
		}
		if hitLabel(shellNo) {
			m.state = StateList
			m.tabIndex = m.shellReturnTab
		}
		return m, nil

	case StateWizardWelcome:
		if onZone && z.kind == zChoice {
			m.wizardStep = z.idx
			return m.wizardSelect()
		}
		return m, nil

	case StateWizardDetect:
		if m.scanReady {
			m.finishScan()
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
		if onZone && z.kind == zImport {
			m.wizardStep = z.idx
			if z.idx < len(m.importSelected) {
				m.importSelected[z.idx] = !m.importSelected[z.idx]
			}
		}
		return m, nil

	case StateWizardAddMore, StateWizardDone:
		if onZone && z.kind == zChoice {
			m.wizardStep = z.idx
			return m.wizardSelect()
		}
		return m, nil

	case StateDeleteConfirm:
		if onRow("confirm delete") && hitLabel("confirm delete") {
			return m.confirmDelete()
		}
		if hitLabel("cancel") {
			m.state = StateEdit
			return m, m.openProfileForm(true, m.formFields)
		}
		return m, nil

	case StateAdd, StateEdit:
		if m.form == nil {
			return m, nil
		}
		key, ok := formFieldAt(scr.lines, relY)
		if !ok {
			return m, nil
		}
		return m.focusFormField(key)

	case StateList:
		if !onZone {
			return m, nil
		}
		switch z.kind {
		case zTabs:
			innerX := lineX - 2 // panel border + padding
			xs, ws := m.tabGeometry()
			for i := range xs {
				if innerX >= xs[i]-1 && innerX <= xs[i]+ws[i] {
					m.setTab(i)
					break
				}
			}
		case zProfile:
			if z.idx >= len(m.profiles) {
				break
			}
			m.cursor = z.idx
			if m.arcadeMode {
				m.selFlashFrame = 0
				m.selFlashProfile = m.cursor
				m.state = StateSelectFlash
				return m, arcadeTickCmd()
			}
			return m.switchTo(m.profiles[z.idx])
		case zItem:
			if m.tabIndex == 1 {
				return m.toggleUtility(z.idx)
			}
			if m.tabIndex != 2 {
				break
			}
			m.settingsFocus = z.idx
			switch z.idx {
			case 0: // config location — only the edit chip launches $EDITOR
				if hitLabel(chipEdit) {
					return m, m.openConfigEditor()
				}
			case 1: // theme — left half previous, right half next
				if !m.arcadeMode {
					if lineX > scr.width()/2 {
						m.cycleTheme(1)
					} else {
						m.cycleTheme(-1)
					}
				}
			case 2: // alias — rename chip edits, anywhere else toggles
				if hitLabel(chipRename) {
					return m.startAliasEdit()
				}
				return m.toggleAlias()
			}
		}
		return m, nil
	}
	return m, nil
}

// dismissWhatsNew closes the release-notes splash and remembers it was seen.
func (m *Model) dismissWhatsNew() {
	ver.MarkVersionSeen(m.store.ConfigDir(), m.currentVersion)
	m.splashSeen020 = true
	_ = m.savePrefs()
	m.state = StateList
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
	m.active = git.DetectActive(profiles)
	// Everything was imported: nothing left to ask, so "Done" is the natural
	// next step. Otherwise start on "Log in with GitHub".
	m.wizardStep = 0
	if len(profiles) > 0 {
		m.wizardStep = 2
	}
	m.state = StateWizardAddMore
}

// startScan enters the setup wizard's Scan step and kicks off the real scan.
// The rows tick off on the frame ticker; with reduced motion (or before the
// first window size arrives) the wizard just advances when the scan returns.
func (m Model) startScan() (tea.Model, tea.Cmd) {
	m.state = StateWizardDetect
	m.scanReady, m.scanStep, m.scanHold = false, 0, 0
	m.scanning = !reduceMotion && m.width > 0
	return m, m.detectExistingConfigsCmd()
}

// leaveWizard returns to the main list and ends first-run setup.
func (m *Model) leaveWizard() {
	m.onboarding = false
	m.scanning = false
	m.state = StateList
	m.tabIndex = 0
}

// wizardSelect activates the focused choice on the current wizard step.
// Shared by enter and mouse click.
func (m Model) wizardSelect() (tea.Model, tea.Cmd) {
	switch m.state {
	case StateWizardWelcome:
		if m.wizardStep == 1 {
			m.leaveWizard()
			return m, nil
		}
		return m.startScan()
	case StateWizardDetect:
		if m.scanReady {
			m.finishScan()
		}
	case StateWizardImport:
		m.importSelectedProfiles()
	case StateWizardAddMore:
		switch m.wizardStep {
		case 0: // log in with GitHub — hands the terminal to `gh auth login`
			m.LaunchOAuth = true
			return m, tea.Quit
		case 1: // manual
			m.state = StateAdd
			m.statusMsg = ""
			m.formFields = [6]string{}
			return m, m.openProfileForm(false, [6]string{})
		default: // done
			if m.onboarding && len(m.profiles) > 0 {
				m.state = StateWizardDone
				return m, nil
			}
			m.leaveWizard()
		}
	case StateWizardDone:
		m.leaveWizard()
	}
	return m, nil
}

func (m Model) updateWizard(msg tea.Msg) (tea.Model, tea.Cmd) {
	km, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	// How many choices the focus can move across on this step.
	choices := 0
	switch m.state {
	case StateWizardWelcome:
		choices = 2
	case StateWizardImport:
		choices = len(m.detectedProfiles)
	case StateWizardAddMore:
		choices = 3
	}
	switch km.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "q":
		if m.state == StateWizardWelcome {
			return m, tea.Quit
		}
	case "esc":
		switch m.state {
		case StateWizardDetect, StateWizardImport:
			m.scanning = false
			m.wizardStep = 0
			m.state = StateWizardWelcome
		case StateWizardAddMore, StateWizardDone:
			m.leaveWizard()
		}
	case "s":
		switch m.state {
		case StateWizardWelcome:
			m.leaveWizard()
		case StateWizardImport:
			m.wizardStep = 0
			m.state = StateWizardAddMore
		}
	case "enter":
		return m.wizardSelect()
	case " ":
		// space toggles a checkbox in the import list
		if m.state == StateWizardImport && m.wizardStep < len(m.importSelected) {
			m.importSelected[m.wizardStep] = !m.importSelected[m.wizardStep]
		}
	case "up", "k":
		if m.wizardStep > 0 {
			m.wizardStep--
		}
	case "down", "j":
		if m.wizardStep < choices-1 {
			m.wizardStep++
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

	// Snapshot before recording, then record unconditionally — regardless of
	// which tier below ends up handling p, its own email must be marked seen
	// so a later, unrelated candidate sharing that literal email still gets
	// deduped by tier 3. Recording only inside tier 3's own branch (as
	// before) missed this whenever tier 1 or tier 2 merged and returned early.
	alreadySeenEmail := p.Email != "" && seenEmail[p.Email]
	if p.Email != "" {
		seenEmail[p.Email] = true
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

	if alreadySeenEmail && p.GHUser == "" {
		return found
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
	// display name and verified GitHub email(s) concurrently (one
	// goroutine per account, each writing only its own slice index — no
	// mutex needed) so a match against the git-config candidate above
	// merges them into one profile instead of surfacing as two, and the
	// candidate's user.name can be the real display name rather than the
	// bare handle — same as `gitswitch login`. IdentityFor makes one
	// `gh auth token` call per account and reuses it for both lookups.
	// Silently empty on any failure (gh missing, insufficient token
	// scope, no network) — falls through to today's nickname/email-string
	// behavior with no regression.
	ghAccounts := git.ListGHUsers()
	names := make([]string, len(ghAccounts))
	verified := make([][]string, len(ghAccounts))
	var wg sync.WaitGroup
	for i, acct := range ghAccounts {
		wg.Add(1)
		go func(i int, acct git.GHAccount) {
			defer wg.Done()
			names[i], verified[i] = git.IdentityFor(acct.Login, acct.Host)
		}(i, acct)
	}
	wg.Wait()

	for i, acct := range ghAccounts {
		nick := acct.Login
		if nick == "" {
			continue
		}
		userName := names[i]
		if userName == "" {
			userName = nick
		}
		add(storage.Profile{
			Nickname: nick,
			UserName: userName,
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

// runUpgrade hands the terminal to the upgrade command for this install —
// `brew upgrade` when Homebrew owns the binary, the install script otherwise.
func (m Model) runUpgrade() (tea.Model, tea.Cmd) {
	cmd, err := ver.UpgradeCommandFor(m.install, m.latestVersion)
	if err != nil {
		m.statusMsg = fmt.Sprintf("upgrade error: %v", err)
		m.statusIsErr = true
		m.state = StateList
		return m, nil
	}
	// Show the release notes on the next launch, whatever the bump size.
	ver.ScheduleWhatsNew(m.store.ConfigDir(), m.latestVersion)
	return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
		return upgradeDoneMsg{err: err}
	})
}

// skipUpdate closes the prompt and stops it reopening by itself for this
// release. The status-line reminder and the `u` key remain.
func (m *Model) skipUpdate() {
	m.skippedUpdate = m.latestVersion
	_ = m.savePrefs()
	m.state = StateList
}

func (m Model) updateUpdatePrompt(msg tea.Msg) (tea.Model, tea.Cmd) {
	km, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch km.String() {
	case "y", "Y", "enter":
		return m.runUpgrade()
	case "n", "N", "esc", "q":
		m.skipUpdate()
		return m, nil
	case "ctrl+c":
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) maxWhatsNewScroll() int {
	total := len(renderNotes(m.whatsNewBody, m.panelWidth()))
	vis := notesViewportHeight(m.height)
	if total > vis {
		return total - vis
	}
	return 0
}
