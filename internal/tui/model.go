package tui

import (
	"github.com/aksisonline/gitswitch/internal/git"
	"github.com/aksisonline/gitswitch/internal/shell"
	"github.com/aksisonline/gitswitch/internal/storage"
	ver "github.com/aksisonline/gitswitch/internal/version"
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
	StateNoProfiles
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
	active   *storage.Profile
	state    State
	width    int
	height   int

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

	// Utilities tab focus (0=shell, 1=precommit, 2=credential)
	utilityFocus int
	// Settings tab focus (0=config, 1=theme)
	settingsFocus int
	// Shell integration toggle
	shellEnabled bool
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

	LaunchLogin      bool
	LaunchOAuth      bool
	PendingReloadCmd string
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

type Option func(*Model)

func WithArcadeMode() Option {
	return func(m *Model) {
		m.arcadeMode = true
		m.state = StateIntro
		m.introMouthOpen = true
		m.hiScore = 99990
	}
}

func New(store *storage.Store, currentVersion string, opts ...Option) (*Model, error) {
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
	m := &Model{
		store:          store,
		profiles:       profiles,
		active:         active,
		state:          StateList,
		currentVersion: currentVersion,
		colorTheme:     prefs.ColorTheme,
		shellEnabled:   shell.IsInstalled(shell.RCFile(shell.DetectShell())),
		showUsername:   prefs.ShowUsername,
		splashSeen020:  prefs.SplashSeen020,
	}
	for _, opt := range opts {
		opt(m)
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
		latest := ver.CachedLatestVersion(configDir, cv)
		return versionCheckMsg{latest: latest}
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
		ColorTheme:    m.colorTheme,
		SplashSeen020: m.splashSeen020,
		ShellEnabled:  m.shellEnabled,
		ShowUsername:  m.showUsername,
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
