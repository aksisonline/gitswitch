package tui

import (
	"fmt"
	"strings"

	ver "github.com/aksisonline/gitswitch/internal/version"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type upgradeState int

const (
	upgradeLoading   upgradeState = iota
	upgradeLatest                 // already on latest — show current release notes
	upgradeAvailable              // update available — show new version notes + y/n
	upgradeRunning                // upgrade script executing (tea.ExecProcess)
	upgradeDone                   // done, waiting for keypress
	upgradeBrew                   // installed via Homebrew — can't auto-upgrade
	upgradeErr                    // network/fetch error
)

// UpgradeModel is a standalone bubbletea model for `gitswitch upgrade`.
type UpgradeModel struct {
	currentVersion string
	latestVersion  string
	notes          string
	state          upgradeState
	isBrew         bool
	configDir      string
	width          int
	height         int
	err            error
	upgradeErr     error
	scroll         int
}

type upgradeFetchedMsg struct {
	latest string
	notes  string
	err    error
}

type upgradeFinishedMsg struct{ err error }

func NewUpgradeModel(currentVersion string, isBrew bool, configDir string) UpgradeModel {
	return UpgradeModel{
		currentVersion: currentVersion,
		isBrew:         isBrew,
		configDir:      configDir,
		state:          upgradeLoading,
	}
}

func (m UpgradeModel) Init() tea.Cmd {
	if m.isBrew {
		return func() tea.Msg { return upgradeFetchedMsg{} }
	}
	return func() tea.Msg {
		latest, err := ver.FetchLatestVersionFresh()
		if err != nil {
			return upgradeFetchedMsg{err: err}
		}
		// Show notes for whichever version is most relevant.
		target := latest
		if !ver.IsUpdateAvailable(m.currentVersion, latest) {
			target = m.currentVersion
		}
		notes, _ := ver.FetchReleaseNotes(target)
		return upgradeFetchedMsg{latest: latest, notes: notes}
	}
}

func (m UpgradeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil

	case upgradeFetchedMsg:
		if m.isBrew {
			m.state = upgradeBrew
			return m, nil
		}
		if msg.err != nil {
			m.err = msg.err
			m.state = upgradeErr
			return m, nil
		}
		m.latestVersion = msg.latest
		m.notes = msg.notes
		if ver.IsUpdateAvailable(m.currentVersion, m.latestVersion) {
			m.state = upgradeAvailable
		} else {
			m.state = upgradeLatest
		}
		return m, nil

	case upgradeFinishedMsg:
		m.upgradeErr = msg.err
		m.state = upgradeDone
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case upgradeLoading:
			if msg.String() == "ctrl+c" {
				return m, tea.Quit
			}

		case upgradeLatest, upgradeBrew, upgradeErr, upgradeDone:
			return m, tea.Quit

		case upgradeAvailable:
			switch msg.String() {
			case "ctrl+c", "n", "N", "q", "esc":
				return m, tea.Quit
			case "y", "Y", "enter":
				// Mark current version so next launch shows What's New for new version.
				ver.MarkVersionSeen(m.configDir, m.currentVersion)
				m.state = upgradeRunning
				cmd, err := ver.UpgradeCommand(m.latestVersion)
				if err != nil {
					m.upgradeErr = err
					m.state = upgradeDone
					return m, nil
				}
				return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
					return upgradeFinishedMsg{err: err}
				})
			case "up", "k":
				if m.scroll > 0 {
					m.scroll--
				}
			case "down", "j", " ":
				m.scroll++
			}
		}
	}
	return m, nil
}

func (m UpgradeModel) View() string {
	applyTheme(normalThemes[0], false)
	pw := 72
	if m.width > 0 {
		if avail := m.width - 6; avail < pw {
			pw = avail
		}
		if pw < minPanelWidth {
			pw = minPanelWidth
		}
	}

	var body string
	switch m.state {
	case upgradeLoading:
		body = m.viewLoading(pw)
	case upgradeLatest:
		body = m.viewNotes(pw, false)
	case upgradeAvailable:
		body = m.viewNotes(pw, true)
	case upgradeRunning:
		body = stylePanelBorder(pw).Render("\n  " + styleTitle.Render("Upgrading…") + "\n\n")
	case upgradeDone:
		body = m.viewDone(pw)
	case upgradeBrew:
		body = m.viewBrew(pw)
	case upgradeErr:
		body = m.viewErr(pw)
	}

	if m.width > 0 && m.height > 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, body)
	}
	return body
}

func (m UpgradeModel) viewLoading(pw int) string {
	inner := "\n  " + styleBrand.Render("Checking for updates…") + "\n\n"
	return stylePanelBorder(pw).Render(inner)
}

func (m UpgradeModel) viewBrew(pw int) string {
	title := styleTitle.Render("  Installed via Homebrew")
	msg := styleBrand.Render("  Run: ") + styleCurrentVal.Render("brew upgrade gitswitch")
	footer := "\n" + divider(pw) + "\n" + m.footerLine(pw, "any key  close")
	return stylePanelBorder(pw).Render("\n" + title + "\n\n  " + msg + "\n" + footer)
}

func (m UpgradeModel) viewErr(pw int) string {
	title := styleDeleteTitle.Render("  Could not check for updates")
	msg := styleBrand.Render(fmt.Sprintf("  %v", m.err))
	footer := "\n" + divider(pw) + "\n" + m.footerLine(pw, "any key  close")
	return stylePanelBorder(pw).Render("\n" + title + "\n\n  " + msg + "\n" + footer)
}

func (m UpgradeModel) viewDone(pw int) string {
	var title, msg string
	if m.upgradeErr != nil {
		title = styleDeleteTitle.Render("  Upgrade failed")
		msg = styleBrand.Render(fmt.Sprintf("  %v", m.upgradeErr))
	} else {
		title = styleCheckmark.Render("  ✓ Upgrade complete")
		msg = styleBrand.Render("  Restart your shell to use ") + styleCurrentVal.Render(m.latestVersion)
	}
	footer := "\n" + divider(pw) + "\n" + m.footerLine(pw, "any key  close")
	return stylePanelBorder(pw).Render("\n" + title + "\n\n  " + msg + "\n" + footer)
}

func (m UpgradeModel) viewNotes(pw int, updateAvailable bool) string {
	var header string
	if updateAvailable {
		from := styleCurrentVal.Render(m.currentVersion)
		to := styleCheckmark.Render(m.latestVersion)
		header = "  " + styleTitle.Render("Update available: ") + from + styleBrand.Render(" → ") + to
	} else {
		header = "  " + styleTitle.Render("You're on the latest — ") + styleCurrentVal.Render(m.currentVersion)
	}

	noteLines := renderNotes(m.notes, pw)

	visibleH := m.height - 12
	if visibleH < 4 {
		visibleH = 4
	}
	total := len(noteLines)
	offset := m.scroll
	if offset > total-visibleH {
		offset = total - visibleH
	}
	if offset < 0 {
		offset = 0
	}
	end := offset + visibleH
	if end > total {
		end = total
	}

	scrollHint := ""
	if total > visibleH {
		scrollHint = "\n  " + styleBrand.Render(fmt.Sprintf("(%d/%d lines — ↑/↓ scroll)", offset+visibleH, total))
	}

	notesBody := strings.Join(noteLines[offset:end], "\n")

	var footerKeys string
	if updateAvailable {
		footerKeys = m.footerLine(pw, "y  upgrade  ·  n  cancel")
	} else {
		footerKeys = m.footerLine(pw, "any key  close")
	}

	body := "\n" + header + "\n\n" + notesBody + scrollHint + "\n"
	footer := "\n" + divider(pw) + "\n" + footerKeys
	return stylePanelBorder(pw).Render(body + footer)
}

// renderNotes does minimal markdown cleanup — same logic as viewWhatsNew.
func renderNotes(raw string, _ int) []string {
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(colorGreen)
	bulletStyle := lipgloss.NewStyle().Foreground(colorDim)
	dimStyle := lipgloss.NewStyle().Foreground(colorDim)

	raw = strings.ReplaceAll(raw, "\r\n", "\n")
	raw = strings.ReplaceAll(raw, "**", "")
	var lines []string
	for _, line := range strings.Split(raw, "\n") {
		switch {
		case strings.HasPrefix(line, "## "):
			lines = append(lines, "\n  "+headerStyle.Render(strings.TrimPrefix(line, "## ")))
		case strings.HasPrefix(line, "# "):
			lines = append(lines, "\n  "+headerStyle.Render(strings.TrimPrefix(line, "# ")))
		case strings.HasPrefix(line, "* "):
			lines = append(lines, "  "+bulletStyle.Render("•")+" "+strings.TrimPrefix(line, "* "))
		case strings.HasPrefix(line, "- "):
			lines = append(lines, "  "+bulletStyle.Render("•")+" "+strings.TrimPrefix(line, "- "))
		case strings.TrimSpace(line) == "":
			lines = append(lines, "")
		default:
			lines = append(lines, "  "+dimStyle.Render(line))
		}
	}
	return lines
}

func (m UpgradeModel) footerLine(pw int, hint string) string {
	parts := strings.Split(hint, "  ·  ")
	var rendered []string
	for _, p := range parts {
		kv := strings.SplitN(p, "  ", 2)
		if len(kv) == 2 {
			rendered = append(rendered, styleFooterKey.Render(kv[0])+" "+styleFooter.Render(kv[1]))
		} else {
			rendered = append(rendered, styleFooter.Render(p))
		}
	}
	return "  " + strings.Join(rendered, styleFooter.Render("  ·  "))
}
