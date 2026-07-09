package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Settings sections (indices match settingsFocus and mouse grid):
//   0 = Config Location  (display + edit button)
//   1 = Theme            (cycle with ← →)
//   2 = Shell Alias      (editable shortcut name, default "gs")

func (m Model) viewSettingsTab(pw int) string {
	var top string
	if m.arcadeMode {
		top = m.viewScoreLine(pw) + "\n"
	}

	header := m.viewHeader("") + m.viewTabHeader()
	iw := itemInnerW(pw)

	// Config location box
	configPath := truncate(m.store.ConfigDir()+"/config.yaml", iw-6)
	editChip := lipgloss.NewStyle().Foreground(colorPurple).Render("[✎ edit]")
	configTitle := "Config Location"
	if m.arcadeMode {
		configTitle = "SAVE FILE LOCATION"
	}
	configLine1 := titleWithRight(styleItemDim.Render(configTitle), editChip, iw)
	configLine2 := lipgloss.NewStyle().Foreground(colorDim).Render(configPath)
	configLine2 = padTo(configLine2, iw)
	configBox := renderItemBox(pw, m.settingsFocus == 0, false, configLine1, configLine2)

	// Theme box
	themeTitle := "Theme"
	if m.arcadeMode {
		themeTitle = "VISUAL THEME"
	}
	cycleTip := styleBrand.Render("← →")
	themeLine1 := titleWithRight(styleCurrentVal.Render(themeTitle), cycleTip, iw)

	// Show current theme name + prev/next hint
	name := themeNames[m.colorTheme]
	var nameStr string
	if m.arcadeMode {
		nameStr = fmt.Sprintf("THEME %02d/%d", m.colorTheme+1, len(themeNames))
	} else {
		nameStr = fmt.Sprintf("%-20s  (%d/%d)", name, m.colorTheme+1, len(themeNames))
	}
	themeLine2 := lipgloss.NewStyle().Foreground(colorYellow).Render(truncate(nameStr, iw))
	themeLine2 = padTo(themeLine2, iw)
	themeBox := renderItemBox(pw, m.settingsFocus == 1, false, themeLine1, themeLine2)

	// Shell alias box
	aliasTitle := "Shell Alias"
	if m.arcadeMode {
		aliasTitle = "COMMAND ALIAS"
	}
	aliasEditChip := lipgloss.NewStyle().Foreground(colorPurple).Render("[✎ rename]")
	aliasLine1 := titleWithRight(styleItemDim.Render(aliasTitle), aliasEditChip, iw)
	var aliasLine2 string
	if m.aliasEditing && m.settingsFocus == 2 {
		aliasLine2 = lipgloss.NewStyle().Foreground(colorGreen).Render(m.aliasInput.View())
	} else if m.shellAliasDisabled {
		toggleDot := lipgloss.NewStyle().Foreground(colorDim).Render("○")
		aliasLine2 = padTo(toggleDot+" "+styleItemDim.Render("disabled"), iw)
	} else {
		toggleDot := lipgloss.NewStyle().Foreground(colorGreen).Render("●")
		aliasLine2 = padTo(toggleDot+" "+styleCurrentVal.Render(m.shellAlias), iw)
	}
	aliasBox := renderItemBox(pw, m.settingsFocus == 2, false, aliasLine1, aliasLine2)

	sep := "\n\n"
	var footerHints [][2]string
	if m.aliasEditing {
		footerHints = [][2]string{
			{"enter", "save"},
			{"esc", "cancel"},
		}
	} else if m.settingsFocus == 2 {
		footerHints = [][2]string{
			{"↑/↓", "move"},
			{"enter", "toggle on/off"},
			{"e", "rename alias"},
			{"q", "quit"},
		}
	} else {
		footerHints = [][2]string{
			{"↑/↓", "move"},
			{"← →", "cycle theme"},
			{"1/2/3", "tabs"},
			{"q", "quit"},
		}
	}
	footer := sep + divider(pw) + "\n" + m.footerKeys(pw, footerHints)

	return top + stylePanelBorder(pw).Render(header+configBox+themeBox+aliasBox+footer)
}
