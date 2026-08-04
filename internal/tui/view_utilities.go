package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Utility items (indices match utilityFocus and the 5-relY mouse grid):
//   0 = Shell Integration       (toggleable)
//   1 = Session Isolation       (toggleable)
//   2 = HTTPS Credential Helper (toggleable)
//   3 = Auto-pin Repeat Repos   (toggleable)

func (m Model) viewUtilitiesTab(pw int) string {
	var top string
	if m.arcadeMode {
		top = m.viewScoreLine(pw) + "\n"
	}

	header := m.viewHeader("") + m.viewTabHeader()
	iw := itemInnerW(pw)

	// Build each item box
	items := m.utilItem(pw, iw, 0) + m.utilItem(pw, iw, 1) + m.utilItem(pw, iw, 2) + m.utilItem(pw, iw, 3)

	sep := "\n\n"
	footer := sep + divider(pw) + "\n" + m.footerKeys(pw, [][2]string{
		{"enter", "toggle"},
		{"1/2/3", "tabs"},
		{"q", "quit"},
	})

	return top + stylePanelBorder(pw).Render(header+items+footer)
}

// utilItemSpec describes one utility toggle's copy (normal and arcade-mode)
// and how to read its current on/off state from the model.
type utilItemSpec struct {
	title, arcadeTitle string
	desc, arcadeDesc   string
	enabled            func(m Model) bool
}

var utilItemSpecs = [4]utilItemSpec{
	{
		title: "Shell Integration", arcadeTitle: "SHELL HOOK",
		desc:       "Auto-switch identity when you cd into a pinned repo.",
		arcadeDesc: "Auto-switch on cd. Works like a cheat code for repos.",
		enabled:    func(m Model) bool { return m.shellEnabled },
	},
	{
		title: "Session Isolation", arcadeTitle: "GH CO-OP MODE",
		desc:       "Per-repo git identity + gh account — needed for pins to work.",
		arcadeDesc: "No more fighting over one shared gh account. Every terminal plays its own.",
		enabled:    func(m Model) bool { return m.ghWrapperEnabled },
	},
	{
		title: "HTTPS Credential Helper", arcadeTitle: "CREDENTIAL HELPER",
		desc:       "Route HTTPS git pushes through the active profile's gh account.",
		arcadeDesc: "HTTPS auth. Automatic. No 401 game-overs.",
		enabled:    func(m Model) bool { return m.credentialHelperEnabled },
	},
	{
		title: "Auto-pin Repeat Repos", arcadeTitle: "AUTO-CLAIM",
		desc:       "Pin the account after 3 uses in the same repo.",
		arcadeDesc: "Play a repo 3 times with an account — it's yours.",
		enabled:    func(m Model) bool { return !m.autoPinDisabled },
	},
}

func (m Model) utilItem(pw, iw, idx int) string {
	if idx < 0 || idx >= len(utilItemSpecs) {
		return ""
	}
	spec := utilItemSpecs[idx]
	focused := m.utilityFocus == idx

	title, desc := spec.title, spec.desc
	if m.arcadeMode {
		title, desc = spec.arcadeTitle, spec.arcadeDesc
	}
	toggle := renderToggle(spec.enabled(m))
	line1 := titleWithRight(styleCurrentVal.Render(title), toggle, iw)
	line2 := lipgloss.NewStyle().Foreground(colorDim).Render(truncate(desc, iw))
	line2 = padTo(line2, iw)
	return renderItemBox(pw, focused, false, line1, line2)
}
