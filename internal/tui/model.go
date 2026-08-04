package tui

import (
	"os"

	"github.com/aksisonline/gitswitch/internal/git"
	"github.com/aksisonline/gitswitch/internal/history"
	"github.com/aksisonline/gitswitch/internal/shell"
	"github.com/aksisonline/gitswitch/internal/storage"
	ver "github.com/aksisonline/gitswitch/internal/version"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// openProfileForm builds the huh add/edit form seeded from the given values
// and returns the form's Init command. edit=true wires the edit copy.
func (m *Model) openProfileForm(edit bool, seed [6]string) tea.Cmd {
	m.formData = &profileFormData{
		nickname: seed[0],
		userName: seed[1],
		email:    seed[2],
		signKey:  seed[3],
		sshKey:   seed[4],
		ghUser:   seed[5],
	}
	m.form = newProfileForm(m.formData, edit, m.panelWidth()-6)
	return m.form.Init()
}

type State int

const (
	StateList State = iota
	StateAdd
	StateEdit
	StateDeleteConfirm
	StateTips
	StateIntro
	StateSelectFlash
	StateTransition
	StateExitAnim
	StateWhatsNew      // one-time upgrade splash for v0.1.x users
	StateUpdatePrompt  // shown when a newer version is found on launch
	StateWizardWelcome // new-user onboarding step 0
	StateWizardDetect  // new-user step 1: scanning for existing configs
	StateWizardImport  // new-user step 2: import confirmation
	StateWizardAddMore // new-user step 3: add more accounts
	StateWizardDone    // new-user step 4: complete
	StateShellConfirm  // confirm install/uninstall of shell integration
)

type Model struct {
	store    *storage.Store
	profiles []storage.Profile
	cursor   int
	active   *storage.Profile // the *global* identity, from git.DetectActive
	state    State
	width    int
	height   int

	// Repo awareness: which identity the repo we were launched from actually uses,
	// and where it comes from. repoKey == "" means "not inside a git repo" and is
	// the guard for every pin action.
	repoDir        string
	repoKey        string
	scope          git.Scope
	scopeProfile   *storage.Profile // profile the repo/session resolves to, nil when global
	pinnedInactive string           // nickname pinned to this repo while Session Isolation is off

	formFields  [6]string // seed values when entering the add/edit form
	formFocus   int
	editingNick string

	// huh-powered add/edit form (nil unless in StateAdd/StateEdit)
	form     *huh.Form
	formData *profileFormData

	statusMsg   string
	statusIsErr bool

	currentVersion  string
	latestVersion   string
	releaseNotes    string
	updateAvailable bool

	colorTheme int // 0-11 normal palette index

	arcadeMode      bool
	introPos        int
	introMouthOpen  bool
	introPhase      int
	introReadyFrame int
	introGhostsEat  int // ghosts eaten in frightened phase (0..4)

	selFlashFrame   int
	selFlashProfile int

	transFrame  int
	transTarget State

	exitFrame int

	// pacman score state — purely cosmetic
	score   int
	hiScore int

	// Tab navigation (used when state == StateList)
	tabIndex int // 0=Accounts 1=Utilities 2=Settings

	// Utilities tab focus (0=shell, 1=gh wrapper/Session Isolation, 2=credential)
	utilityFocus int
	// Settings tab focus (0=config, 1=theme)
	settingsFocus int
	// Shell integration toggle
	shellEnabled bool
	// HTTPS credential helper toggle — true means gitswitch will actually be
	// asked first (git.IsCredentialHelperInstalled), not just "registered".
	credentialHelperEnabled bool
	// gh CLI wrapper toggle — true means bare `gh` commands resolve the
	// account per-repo instead of relying on gh's single global active
	// account (shell.IsGHWrapperInstalled).
	ghWrapperEnabled bool
	// autoPinDisabled turns off auto-pinning the active account to a repo once
	// it's been used there history.AutoPinThreshold times (recordCmd,
	// cmd/gitswitch/main.go). Zero value (false) means enabled — mirrors
	// storage.Prefs.AutoPinDisabled.
	autoPinDisabled bool
	// Accounts secondary column: false=email (default), true=GitHub username
	showUsername bool

	// Shell-integration confirm dialog state
	pendingShellInstall bool // true = about to install, false = about to remove
	shellReturnTab      int  // tab to return to after the dialog

	// New-user wizard
	wizardStep       int
	detectedProfiles []storage.Profile
	importSelected   []bool

	// Upgrade splash
	splashSeen020 bool

	// Shell alias (editable in Settings tab)
	shellAlias         string
	shellAliasDisabled bool
	aliasEditing       bool
	aliasInput         textinput.Model

	LaunchLogin      bool
	LaunchOAuth      bool
	PendingReloadCmd string

	whatsNewBody   string
	whatsNewScroll int
}

var formLabels = [6]string{
	"Nickname",
	"User Name",
	"Email",
	"Signing Key",
	"SSH Key Path",
	"GitHub Username",
}

var formSubtitles = [6]string{
	"label shown in this list — not written to git config",
	"git user.name — author name on commits",
	"git user.email — author email on commits",
	"GPG key ID or SSH key path — optional, leave blank to skip",
	"sets core.sshCommand, e.g. ~/.ssh/id_work — optional",
	"for gh auth switch — optional, leave blank to skip",
}

// refreshRepoScope re-reads the repo gitswitch was launched from: its key, and
// which identity git will actually author with there. Called at startup and after
// anything that writes git config, so the glyphs never go stale.
func (m *Model) refreshRepoScope() {
	m.repoDir, _ = os.Getwd()
	m.repoKey = history.GetRepoKeyForPath(m.repoDir)
	m.scope, m.scopeProfile, m.pinnedInactive = git.ScopeGlobal, nil, ""
	if m.repoKey == "" {
		return
	}
	scope, email := git.ResolveIdentity(m.repoDir)
	if scope == git.ScopeRepo && !m.ghWrapperEnabled {
		m.pinnedInactive = history.GetPinned(m.repoKey)
		return
	}
	if scope == git.ScopeRepo || scope == git.ScopeSession {
		if p := m.store.GetByEmail(email); p != nil {
			m.scope, m.scopeProfile = scope, p
		}
	}
}

// scopeGlyph returns the state-column glyph for a profile row: how this profile
// relates to the repo we're in and to the global identity. Single-width glyphs
// only — the column budget in panelWidth assumes one cell.
func (m Model) scopeGlyph(p storage.Profile) string {
	isGlobal := m.active != nil && p.Nickname == m.active.Nickname
	isScoped := m.scopeProfile != nil && p.Nickname == m.scopeProfile.Nickname

	switch {
	case isScoped && m.scope == git.ScopeSession:
		if m.arcadeMode {
			return "▲"
		}
		return "◆"
	case isScoped && isGlobal:
		if m.arcadeMode {
			return "✦"
		}
		return "◉"
	case isScoped:
		if m.arcadeMode {
			return "◆"
		}
		return "●"
	case isGlobal:
		if m.arcadeMode {
			return "★"
		}
		return "✓"
	}
	return "·"
}

// New builds the TUI model. whatsNewBody, when non-empty, opens straight to
// the What's New screen with that release-notes body instead of the list.
func New(store *storage.Store, currentVersion, whatsNewBody string) (*Model, error) {
	profiles, err := store.Load()
	if err != nil {
		return nil, err
	}
	active := git.DetectActive(profiles)
	prefs, err := store.LoadPrefs()
	if err != nil {
		prefs = storage.Prefs{}
	}
	if prefs.ColorTheme < 0 || prefs.ColorTheme >= len(normalThemes) {
		prefs.ColorTheme = 0
	}
	shellAlias := prefs.ShellAlias
	if shellAlias == "" {
		shellAlias = "gs"
	}
	aliasInput := textinput.New()
	aliasInput.CharLimit = 32
	aliasInput.Width = 20
	m := &Model{
		store:                   store,
		profiles:                profiles,
		active:                  active,
		state:                   StateList,
		currentVersion:          currentVersion,
		colorTheme:              prefs.ColorTheme,
		shellEnabled:            shell.IsInstalled(shell.RCFile(shell.DetectShell())),
		credentialHelperEnabled: git.IsCredentialHelperInstalled(),
		ghWrapperEnabled:        shell.IsGHWrapperInstalled(shell.RCFile(shell.DetectShell())),
		showUsername:            prefs.ShowUsername,
		splashSeen020:           prefs.SplashSeen020,
		hiScore:                 prefs.ArcadeHiScore,
		shellAlias:              shellAlias,
		shellAliasDisabled:      prefs.ShellAliasDisabled,
		aliasInput:              aliasInput,
		arcadeMode:              prefs.ArcadeMode,
		autoPinDisabled:         prefs.AutoPinDisabled,
	}
	m.refreshRepoScope()
	if whatsNewBody != "" {
		m.whatsNewBody = whatsNewBody
		m.state = StateWhatsNew
	}
	// Arcade mode persists across launches (via `gitswitch pacman`), so the intro
	// must trigger from the persisted flag, not just the one-shot toggle option —
	// otherwise every launch after the first skips straight to the list.
	if m.arcadeMode && m.state == StateList {
		m.state = StateIntro
		m.introMouthOpen = true
		// Beatable factory high score — real high scores persist via prefs.
		if m.hiScore < 5000 {
			m.hiScore = 5000
		}
	}
	if store.WasMigrated() {
		if store.BakCreated() {
			m.statusMsg = "Profiles migrated to config.yaml (backup: profiles.json.v1.bak)"
		} else {
			m.statusMsg = "Profiles migrated to config.yaml"
		}
	}
	if !m.arcadeMode {
		if len(profiles) == 0 {
			m.state = StateWizardWelcome
		} else if !prefs.SplashSeen020 {
			m.state = StateWhatsNew
		}
	}
	return m, nil
}

func (m Model) Init() tea.Cmd {
	configDir := m.store.ConfigDir()
	cv := m.currentVersion
	versionCmd := func() tea.Msg {
		rel := ver.CachedLatestRelease(configDir, cv)
		return versionCheckMsg{latest: rel.Version, notes: rel.Notes}
	}
	if m.arcadeMode {
		return tea.Batch(versionCmd, arcadeTickCmd())
	}
	return versionCmd
}

// savePrefs persists all current preference fields in one place so callers
// never accidentally clobber a field by omitting it from a struct literal.
func (m *Model) savePrefs() error {
	return m.store.SavePrefs(storage.Prefs{
		ColorTheme:         m.colorTheme,
		SplashSeen020:      m.splashSeen020,
		ShellEnabled:       m.shellEnabled,
		ShowUsername:       m.showUsername,
		ShellAlias:         m.shellAlias,
		ShellAliasDisabled: m.shellAliasDisabled,
		ArcadeHiScore:      m.hiScore,
		ArcadeMode:         m.arcadeMode,
		AutoPinDisabled:    m.autoPinDisabled,
	})
}

func (m Model) panelWidth() int {
	content := minPanelWidth
	for _, p := range m.profiles {
		needed := 3 + 3 + 14 + 2 + lipgloss.Width(p.Email) + 6
		if needed > content {
			content = needed
		}
		if nickNeeded := 3 + 3 + lipgloss.Width(p.Nickname) + 2 + lipgloss.Width(p.Email) + 6; nickNeeded > content {
			content = nickNeeded
		}
	}
	for _, s := range formSubtitles {
		if needed := lipgloss.Width(s) + 8; needed > content {
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
