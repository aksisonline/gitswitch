package tui

import "github.com/charmbracelet/lipgloss"

// Utility items (indices match utilityFocus):
//   0 = Shell Integration       (toggleable)
//   1 = Session Isolation       (toggleable)
//   2 = HTTPS Credential Helper (toggleable)
//   3 = Auto-pin Repeat Repos   (toggleable)
//   4 = Pre-commit Guard        (toggleable)

func (m Model) drawUtilities(c *canvas) {
	m.header(c, "", m.versionTag())
	if m.headGap() {
		c.blank()
	}
	m.drawTabs(c)

	items := make([]itemView, len(utilItemSpecs))
	for i, spec := range utilItemSpecs {
		title, desc := spec.title, spec.desc
		if m.arcadeMode {
			title, desc = spec.arcadeTitle, spec.arcadeDesc
		}
		toggle := renderToggle(spec.enabled(m))
		if m.busyItem == i+1 {
			toggle = m.renderToggleBusy()
		}
		items[i] = itemView{
			title:  title,
			right:  toggle,
			detail: lipgloss.NewStyle().Foreground(colorDim).Render(desc),
		}
	}
	m.drawItems(c, m.utilityFocus, items)

	if c.d != dTight {
		c.blank()
	}
	c.add(m.statusLine())
	c.footer(m, [][2]string{
		{"↑/↓", "move"},
		{"enter", "toggle"},
		{"tab", "next tab"},
		{"q", "quit"},
	})
}

// utilItemSpec describes one utility toggle's copy (normal and arcade-mode)
// and how to read its current on/off state from the model.
type utilItemSpec struct {
	title, arcadeTitle string
	desc, arcadeDesc   string
	enabled            func(m Model) bool
}

var utilItemSpecs = [5]utilItemSpec{
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
	{
		title: "Pre-commit Guard", arcadeTitle: "COMMIT GUARD",
		desc:       "Refuse wrong-author commits before they enter git history.",
		arcadeDesc: "Blocks wrong-account commits before they hit git log.",
		enabled:    func(m Model) bool { return m.guardEnabled },
	},
}
