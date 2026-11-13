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

// resolvedShellAlias is the alias to actually write: the configured one,
// unless the user has disabled it (blank tells shell.Install/Reinstall to
// skip writing an alias at all).
func (m Model) resolvedShellAlias() string {
	if m.shellAliasDisabled {
		return ""
	}
	return m.shellAlias
}

// runShellAction performs the install or uninstall in a command.
func (m Model) runShellAction(install bool) tea.Cmd {
	return func() tea.Msg {
		sh := shell.DetectShell()
		fw := shell.DetectFramework()
		if install {
			res, err := shell.Install(sh, fw, m.resolvedShellAlias())
			return shellDoneMsg{installed: true, result: res, err: err}
		}
		res, err := shell.Uninstall(sh)
		return shellDoneMsg{installed: false, result: res, err: err}
	}
}

// reinstallShellCmd removes the existing integration block and rewrites it
// with the current alias setting, then signals the shell to be reloaded.
func (m Model) reinstallShellCmd() tea.Cmd {
	return func() tea.Msg {
		sh := shell.DetectShell()
		fw := shell.DetectFramework()
		res, err := shell.Reinstall(sh, fw, m.resolvedShellAlias())
		return shellDoneMsg{installed: true, result: res, err: err}
	}
}

func (m Model) updateShellConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y", "enter":
			m.setBusy("updating shell integration…")
			return m, m.runShellAction(m.pendingShellInstall)
		case "n", "N", "esc", "q":
			m.state = StateList
			m.tabIndex = m.shellReturnTab
			return m, nil
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	// Mouse clicks are handled in handleMouse and require a precise hit on a
	// button — same spirit as the explicit y/n keys, no accidental installs.
	return m, nil
}

// Button labels on the shell-integration dialog — also the mouse targets.
const (
	shellYesInstall = "Yes, install it"
	shellYesRemove  = "Yes, remove it"
	shellNo         = "Cancel"
)

func (m Model) drawShellConfirm(c *canvas) {
	sh := shell.DetectShell()
	rc := shell.RCFile(sh)

	title, body, yes := "Disable shell integration?", "gitswitch will safely remove its hook block from:", shellYesRemove
	if m.pendingShellInstall {
		title = "Enable shell integration?"
		body = fmt.Sprintf("Detected %s. gitswitch will add its auto-switch hook to:", shellDisplayName(sh))
		yes = shellYesInstall
	}
	note := lipgloss.NewStyle().Foreground(colorDim)

	m.header(c, "Shell Integration", "")
	c.gap()
	c.add("  " + lipgloss.NewStyle().Foreground(colorPurple).Bold(true).Render(title))
	c.gap()
	c.add("  " + note.Render(body))
	c.add("  " + lipgloss.NewStyle().Foreground(colorGreen).Render(rc))
	if c.d != dTight {
		c.gap()
		c.add("  " + note.Render("Only gitswitch's own block changes; the rest stays as it is."))
	}
	c.gap()
	c.add(buttonRow(yes, shellNo))
	c.gap()
	c.footer(m, [][2]string{{"y / enter", "confirm"}, {"n / esc", "cancel"}})
}
