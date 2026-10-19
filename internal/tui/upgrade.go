package tui

import (
	"fmt"
	"strings"

	ver "github.com/aksisonline/gitswitch/internal/version"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
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

	case tea.MouseMsg:
		if m.state == upgradeAvailable || m.state == upgradeLatest {
			switch msg.Button {
			case tea.MouseButtonWheelUp:
				if m.scroll > 0 {
					m.scroll--
				}
			case tea.MouseButtonWheelDown:
				m.scroll++
			}
		}

	case tea.KeyMsg:
		switch m.state {
		case upgradeLoading:
			if msg.String() == "ctrl+c" {
				return m, tea.Quit
			}

		case upgradeBrew, upgradeErr, upgradeDone:
			return m, tea.Quit

		case upgradeLatest:
			switch msg.String() {
			case "up", "k":
				if m.scroll > 0 {
					m.scroll--
				}
			case "down", "j", " ":
				m.scroll++
			default:
				return m, tea.Quit
			}

		case upgradeAvailable:
			switch msg.String() {
			case "ctrl+c", "n", "N", "q", "esc":
				return m, tea.Quit
			case "y", "Y", "enter":
				// Schedule What's New for the new version on next launch.
				ver.ScheduleWhatsNew(m.configDir, m.latestVersion)
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

	visibleH := notesViewportHeight(m.height)
	total := len(noteLines)
	maxScroll := total - visibleH
	if maxScroll < 0 {
		maxScroll = 0
	}
	offset := m.scroll
	if offset > maxScroll {
		offset = maxScroll
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
		scrollHint = "\n  " + styleBrand.Render(fmt.Sprintf("(%d–%d of %d lines · ↑/↓ scroll)", offset+1, end, total))
	}

	notesBody := strings.Join(noteLines[offset:end], "\n")

	var footerKeys string
	if updateAvailable {
		if total > visibleH {
			footerKeys = m.footerLine(pw, "↑/↓  scroll  ·  y  upgrade  ·  n  cancel")
		} else {
			footerKeys = m.footerLine(pw, "y  upgrade  ·  n  cancel")
		}
	} else {
		if total > visibleH {
			footerKeys = m.footerLine(pw, "↑/↓  scroll  ·  any key  close")
		} else {
			footerKeys = m.footerLine(pw, "any key  close")
		}
	}

	body := "\n" + header + "\n\n" + notesBody + scrollHint + "\n"
	footer := "\n" + divider(pw) + "\n" + footerKeys
	return stylePanelBorder(pw).Render(body + footer)
}

// notesViewportHeight calculates how many pre-wrapped visual lines fit inside
// the What's New / Upgrade modal without overflowing the terminal vertically.
func notesViewportHeight(termHeight int) int {
	if termHeight <= 0 {
		return 12
	}
	// Panel chrome (top/bottom borders, header, spacers, scroll hint, divider,
	// footer, and outer margin) takes ~12 lines. Cap at 16 lines so the modal
	// stays compact and scannable even on tall terminals.
	h := termHeight - 12
	if h > 16 {
		h = 16
	}
	if h < 4 {
		h = 4
	}
	return h
}

// cleanReleaseNotes strips noisy sections from raw GitHub release bodies
// (canary boilerplate banners, redundant headers, install scripts, and
// auto-generated PR/changelog footers) so the TUI splash stays concise.
func cleanReleaseNotes(raw string) string {
	raw = strings.ReplaceAll(raw, "\r\n", "\n")

	// When a release body includes both a cumulative [Unreleased] section and
	// GitHub's tag-to-tag "## What's Changed" section (e.g. canary releases),
	// keep only "## What's Changed" so only the most recent update is shown.
	for _, marker := range []string{"\n## What's Changed", "## What's Changed"} {
		if idx := strings.Index(raw, marker); idx >= 0 {
			raw = raw[idx:]
			break
		}
	}

	var kept []string
	skipSection := false
	inCodeBlock := false
	seenVersionHeader := false

	for _, line := range strings.Split(raw, "\n") {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}
		if inCodeBlock {
			if !skipSection {
				kept = append(kept, line)
			}
			continue
		}

		// Drop canary warning blockquotes at the top of beta releases.
		if strings.HasPrefix(trimmed, ">") && (strings.Contains(trimmed, "canary/beta release") || strings.Contains(trimmed, "gitswitch stable")) {
			continue
		}
		// Drop Full Changelog link line.
		if strings.HasPrefix(trimmed, "**Full Changelog**") || strings.HasPrefix(trimmed, "Full Changelog:") {
			continue
		}

		// Check headings for sections to skip or flatten.
		if strings.HasPrefix(trimmed, "#") {
			heading := strings.TrimSpace(strings.TrimLeft(trimmed, "#"))
			lower := strings.ToLower(heading)

			// Version heading like "## v0.4.4 — 2026-08-20" or "## [v0.4.4]":
			// skip the first one (redundant with modal title) and stop if a second
			// older version section appears.
			if strings.HasPrefix(lower, "v0.") || strings.HasPrefix(lower, "v1.") || strings.HasPrefix(lower, "[v") {
				if seenVersionHeader {
					break
				}
				seenVersionHeader = true
				skipSection = false
				continue
			}

			switch {
			case lower == "what's in this beta" || lower == "what's new" || lower == "what's changed":
				// Redundant with the modal title "What's New in vX.Y.Z"
				skipSection = false
				continue
			case strings.HasPrefix(lower, "install "):
				// User is already running the binary in the TUI; skip curl/ps1 blocks.
				skipSection = true
				continue
			default:
				skipSection = false
			}
		}

		if skipSection {
			continue
		}
		kept = append(kept, line)
	}
	return strings.TrimSpace(strings.Join(kept, "\n"))
}

// renderNotes cleans and pre-wraps markdown release notes into visual terminal
// lines so that 1 element == 1 screen row. This keeps the scroll viewport and
// panel height exact regardless of paragraph length.
func renderNotes(raw string, pw int) []string {
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(colorGreen)
	bulletStyle := lipgloss.NewStyle().Foreground(colorPurple).Bold(true)
	textStyle := lipgloss.NewStyle().Foreground(colorWhite)
	dimStyle := lipgloss.NewStyle().Foreground(colorDim)

	cleaned := cleanReleaseNotes(raw)
	cleaned = strings.ReplaceAll(cleaned, "**", "")

	// Usable inner width inside stylePanelBorder(pw) is pw - 2.
	// Leave 2 leading spaces + 1 trailing margin = pw - 5 for text.
	wrapW := pw - 5
	if wrapW < 20 {
		wrapW = 20
	}

	var lines []string
	prevBlank := true // prevent leading blank lines

	appendBlank := func() {
		if !prevBlank && len(lines) > 0 {
			lines = append(lines, "")
			prevBlank = true
		}
	}

	for _, line := range strings.Split(cleaned, "\n") {
		trimmed := strings.TrimSpace(line)
		switch {
		case trimmed == "":
			appendBlank()

		case strings.HasPrefix(trimmed, "### "), strings.HasPrefix(trimmed, "## "), strings.HasPrefix(trimmed, "# "):
			heading := strings.TrimSpace(strings.TrimLeft(trimmed, "#"))
			if heading == "" {
				continue
			}
			appendBlank()
			for _, wl := range wrapNoteLine(heading, wrapW) {
				lines = append(lines, "  "+headerStyle.Render(wl))
			}
			prevBlank = false

		case strings.HasPrefix(trimmed, "* "), strings.HasPrefix(trimmed, "- "):
			item := strings.TrimSpace(trimmed[2:])
			// Strip noisy PR author/URL suffixes from auto-generated GitHub bullets.
			if idx := strings.Index(item, " by @"); idx > 0 && strings.Contains(item[idx:], "http") {
				item = strings.TrimSpace(item[:idx])
			}
			for _, marker := range []string{" [minor]", " [major]", " [patch]"} {
				item = strings.TrimSuffix(item, marker)
			}
			wrapped := wrapNoteLine(item, wrapW-2)
			for i, wl := range wrapped {
				if i == 0 {
					lines = append(lines, "  "+bulletStyle.Render("•")+" "+textStyle.Render(wl))
				} else {
					lines = append(lines, "    "+textStyle.Render(wl))
				}
			}
			prevBlank = false

		default:
			for _, wl := range wrapNoteLine(trimmed, wrapW) {
				lines = append(lines, "  "+dimStyle.Render(wl))
			}
			prevBlank = false
		}
	}

	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

// wrapNoteLine word-wraps s to width columns (ANSI-aware, breaking long URLs/words if needed).
func wrapNoteLine(s string, width int) []string {
	if width <= 0 {
		return []string{s}
	}
	wrapped := ansi.Wordwrap(s, width, "")
	var out []string
	for _, l := range strings.Split(wrapped, "\n") {
		for lipgloss.Width(l) > width {
			part := ansi.Truncate(l, width, "")
			if part == "" {
				break
			}
			out = append(out, part)
			l = strings.TrimPrefix(l, part)
		}
		out = append(out, l)
	}
	return out
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
