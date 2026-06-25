package tui

import (
	"fmt"

	"github.com/aksisonline/gitswitch/internal/shell"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// shellDoneMsg reports the result of an install/uninstall run.
type shellDoneMsg struct {
	installed bool
	result    string
	err       error
}

func shellDisplayName(sh shell.Shell) string {
	switch sh {
	case shell.ShellZsh:
		return "zsh"
	case shell.ShellBash:
		return "bash"
	case shell.ShellFish:
		return "fish"
	default:
		return "your shell"
	}
}

// openShellConfirm switches to the confirm dialog for installing/removing the
// shell hook. install reflects the action being requested.
func (m *Model) openShellConfirm(install bool) {
	m.pendingShellInstall = install
	m.shellReturnTab = m.tabIndex
	m.statusMsg = ""
	m.state = StateShellConfirm
}

// runShellAction performs the install or uninstall in a command.
func (m Model) runShellAction(install bool) tea.Cmd {
	return func() tea.Msg {
		sh := shell.DetectShell()
		fw := shell.DetectFramework()
		if install {
			alias := m.shellAlias
			if m.shellAliasDisabled {
				alias = ""
			}
			res, err := shell.Install(sh, fw, alias)
			return shellDoneMsg{installed: true, result: res, err: err}
		}
		res, err := shell.Uninstall(sh, fw)
		return shellDoneMsg{installed: false, result: res, err: err}
	}
}

func (m Model) updateShellConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y", "enter":
			return m, m.runShellAction(m.pendingShellInstall)
		case "n", "N", "esc", "q":
			m.state = StateList
			m.tabIndex = m.shellReturnTab
			return m, nil
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	// Mouse is intentionally not wired to confirm here — this is a confirmation
	// dialog, so we require an explicit y/n keypress to avoid accidental installs.
	return m, nil
}

func (m Model) viewShellConfirm(pw int) string {
	sh := shell.DetectShell()
	rc := shell.RCFile(sh)
	name := shellDisplayName(sh)

	header := m.viewHeader("Shell Integration")
	iw := pw - 2

	var title, body, yesLabel string
	if m.pendingShellInstall {
		title = "Enable shell integration?"
		body = fmt.Sprintf("Detected %s. gitswitch will add its auto-switch hook to:", name)
		yesLabel = "Yes, install it"
	} else {
		title = "Disable shell integration?"
		body = "gitswitch will safely remove its hook block from:"
		yesLabel = "Yes, remove it"
	}

	titleStyle := lipgloss.NewStyle().Foreground(colorPurple).Bold(true)
	pathStyle := lipgloss.NewStyle().Foreground(colorGreen)
	noteStyle := lipgloss.NewStyle().Foreground(colorDim)

	yesBtn := lipgloss.NewStyle().
		Background(colorGreen).
		Foreground(lipgloss.Color("0")).
		Bold(true).
		Padding(0, 2).
		Render(yesLabel)
	noBtn := lipgloss.NewStyle().
		Foreground(colorDim).
		Padding(0, 2).
		Render("Cancel")

	content := "\n\n  " + titleStyle.Render(title) +
		"\n\n  " + noteStyle.Render(truncate(body, iw-2)) +
		"\n  " + pathStyle.Render(truncate(rc, iw-2)) +
		"\n\n  " + noteStyle.Render(truncate("Only gitswitch's own block changes — the rest of your file is untouched.", iw-2)) +
		"\n\n  " + yesBtn + "   " + noBtn

	footer := "\n\n" + divider(pw) + "\n" + m.footerKeys(pw, [][2]string{
		{"y / enter", "confirm"},
		{"n / esc", "cancel"},
	})
	return stylePanelBorder(pw).Render(header + content + footer)
}
