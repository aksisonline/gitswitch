package tui

import (
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
)

// Settings items (indices match settingsFocus):
//   0 = Config Location  (display + edit button)
//   1 = Theme            (cycle with ← →)
//   2 = Shell Alias      (editable shortcut name, default "gs")

// Chip labels double as mouse targets — the handler looks for these strings.
const (
	chipEdit   = "[✎ edit]"
	chipRename = "[✎ rename]"
)

// themeSwatch previews a palette as four colored dots.
func themeSwatch(tc ThemeColors) string {
	out := ""
	for _, col := range []lipgloss.Color{tc.Primary, tc.Accent, tc.Highlight, tc.Error} {
		out += lipgloss.NewStyle().Foreground(col).Render("●")
	}
	return out
}

func (m Model) drawSettings(c *canvas) {
	m.header(c, "", m.versionTag())
	if m.headGap() {
		c.blank()
	}
	m.drawTabs(c)

	chip := lipgloss.NewStyle().Foreground(colorPurple)
	dim := lipgloss.NewStyle().Foreground(colorDim)

	configTitle, themeTitle, aliasTitle := "Config Location", "Theme", "Shell Alias"
	if m.arcadeMode {
		configTitle, themeTitle, aliasTitle = "SAVE FILE LOCATION", "VISUAL THEME", "COMMAND ALIAS"
	}

	theme := m.colorTheme
	if theme < 0 || theme >= len(normalThemes) {
		theme = 0
	}
	var themeDetail string
	if m.arcadeMode {
		themeDetail = lipgloss.NewStyle().Foreground(colorYellow).Render("ARCADE PALETTE") + dim.Render("  (locked in arcade mode)")
	} else {
		themeDetail = themeSwatch(normalThemes[theme]) + "  " +
			lipgloss.NewStyle().Foreground(colorYellow).Render(themeNames[theme]) +
			dim.Render(fmt.Sprintf("  (%d/%d)", theme+1, len(themeNames)))
	}

	var aliasDetail string
	switch {
	case m.aliasEditing && m.settingsFocus == 2:
		aliasDetail = lipgloss.NewStyle().Foreground(colorGreen).Render(m.aliasInput.View())
	case m.shellAliasDisabled:
		aliasDetail = dim.Render("○ disabled")
	default:
		aliasDetail = lipgloss.NewStyle().Foreground(colorGreen).Render("●") + " " + styleCurrentVal.Render(m.shellAlias)
	}

	m.drawItems(c, m.settingsFocus, []itemView{
		{title: configTitle, right: chip.Render(chipEdit), detail: dim.Render(filepath.Join(m.store.ConfigDir(), "config.yaml"))},
		{title: themeTitle, right: styleBrand.Render("← →"), detail: themeDetail},
		{title: aliasTitle, right: chip.Render(chipRename), detail: aliasDetail},
	})

	if c.d != dTight {
		c.blank()
	}
	c.add(m.statusLine())

	var hints [][2]string
	switch {
	case m.aliasEditing:
		hints = [][2]string{{"enter", "save"}, {"esc", "cancel"}}
	case m.settingsFocus == 0:
		hints = [][2]string{{"↑/↓", "move"}, {"enter", "open in editor"}, {"tab", "next tab"}, {"q", "quit"}}
	case m.settingsFocus == 2:
		hints = [][2]string{{"↑/↓", "move"}, {"enter", "on/off"}, {"e", "rename"}, {"tab", "next tab"}, {"q", "quit"}}
	default:
		hints = [][2]string{{"↑/↓", "move"}, {"← →", "change theme"}, {"tab", "next tab"}, {"q", "quit"}}
	}
	c.footer(m, hints)
}
