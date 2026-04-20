package tui

import (
	"github.com/aksisonline/gitswitch/internal/git"
	"github.com/aksisonline/gitswitch/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
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
