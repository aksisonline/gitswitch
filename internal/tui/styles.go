package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const minPanelWidth = 56
const maxPanelWidth = 120

// ThemeColors holds the four mutable colors per palette.
type ThemeColors struct {
	Primary   lipgloss.Color
	Accent    lipgloss.Color
	Highlight lipgloss.Color
	Error     lipgloss.Color
}

var normalThemes = [12]ThemeColors{
	{"141", "84", "226", "196"},   // 0: Default (original)
	{"39", "51", "123", "203"},    // 1: Ocean
	{"208", "220", "214", "196"},  // 2: Sunset
	{"34", "118", "154", "196"},   // 3: Forest
	{"255", "250", "245", "203"},  // 4: Mono
	{"213", "207", "219", "196"},  // 5: Rose
	{"51", "195", "159", "203"},   // 6: Arctic
	{"226", "214", "220", "196"},  // 7: Gold
	{"165", "177", "183", "196"},  // 8: Violet
	{"196", "202", "226", "196"},  // 9: Ember
	{"46", "40", "82", "196"},     // 10: Matrix
	{"111", "153", "195", "203"},  // 11: Steel
}

var themeNames = [12]string{
	"Default", "Ocean", "Sunset", "Forest",
	"Mono", "Rose", "Arctic", "Gold",
	"Violet", "Ember", "Matrix", "Steel",
}

var arcadeTheme = ThemeColors{
	Primary:   "226", // pac yellow
	Accent:    "214", // coin gold
	Highlight: "51",  // ghost cyan
	Error:     "196", // ghost red
}

// Arcade-specific fixed colors — always the same regardless of palette.
var (
	arcadeMazeBlue  = lipgloss.Color("27")
	arcadeGhostRed  = lipgloss.Color("196")
	arcadeGhostPink = lipgloss.Color("213")
	arcadeGhostCyan = lipgloss.Color("51")
	arcadePipe      = lipgloss.Color("46")
)

// Mutable color vars — reassigned by applyTheme. Default = original theme.
var (
	colorPurple  = lipgloss.Color("141")
	colorGreen   = lipgloss.Color("84")
	colorYellow  = lipgloss.Color("226")
	colorDim     = lipgloss.Color("241") // fixed
	colorRed     = lipgloss.Color("196")
	colorWhite   = lipgloss.Color("255") // fixed
	colorBgHover = lipgloss.Color("237") // fixed
	isArcadeMode = false
)

// Width-independent styles — rebuilt by applyTheme.
var (
	styleTitle = lipgloss.NewStyle().Bold(true).Foreground(colorPurple)
	styleBrand = lipgloss.NewStyle().Foreground(colorDim)

	styleBrandLink = lipgloss.NewStyle().
			Foreground(colorGreen).
			Underline(true)

	styleCurrent      = lipgloss.NewStyle().Foreground(colorDim)
	styleCurrentVal   = lipgloss.NewStyle().Foreground(colorWhite).Bold(true)
	styleCheckmark    = lipgloss.NewStyle().Foreground(colorGreen)
	styleItemDim      = lipgloss.NewStyle().Foreground(colorDim)
	styleItemInactive = lipgloss.NewStyle().Foreground(colorWhite)

	styleFooter    = lipgloss.NewStyle().Foreground(colorDim)
	styleFooterKey = lipgloss.NewStyle().Foreground(colorYellow)

	styleFormTitle        = lipgloss.NewStyle().Foreground(colorGreen)
	styleInputLabel       = lipgloss.NewStyle().Foreground(colorDim)
	styleInputLabelActive = lipgloss.NewStyle().Foreground(colorPurple).Bold(true)
	styleDeleteTitle      = lipgloss.NewStyle().Bold(true).Foreground(colorRed)
	styleSubtitle         = lipgloss.NewStyle().Foreground(colorDim).Italic(true)
	styleDivider          = lipgloss.NewStyle().Foreground(colorPurple)
)

// applyTheme reassigns color vars and rebuilds dependent style vars.
// Call at the start of View() with the current theme.
func applyTheme(tc ThemeColors, arcade bool) {
	isArcadeMode = arcade
	colorPurple = tc.Primary
	colorGreen = tc.Accent
	colorYellow = tc.Highlight
	colorRed = tc.Error

	styleTitle = lipgloss.NewStyle().Bold(true).Foreground(colorPurple)
	styleBrandLink = lipgloss.NewStyle().Foreground(colorGreen).Underline(true)
	styleCheckmark = lipgloss.NewStyle().Foreground(colorGreen)
	styleFooterKey = lipgloss.NewStyle().Foreground(colorYellow)
	styleFormTitle = lipgloss.NewStyle().Foreground(colorGreen)
	styleInputLabelActive = lipgloss.NewStyle().Foreground(colorPurple).Bold(true)
	styleDeleteTitle = lipgloss.NewStyle().Bold(true).Foreground(colorRed)
	styleDivider = lipgloss.NewStyle().Foreground(colorPurple)
}

// Width-dependent styles — read current vars at call time.

func stylePanelBorder(w int) lipgloss.Style {
	border := lipgloss.Border(lipgloss.RoundedBorder())
	if isArcadeMode {
		border = lipgloss.DoubleBorder()
	}
	return lipgloss.NewStyle().
		Border(border).
		BorderForeground(colorPurple).
		Width(w).
		PaddingLeft(1).
		PaddingRight(1)
}

func styleItemActive(w int) lipgloss.Style {
	return lipgloss.NewStyle().
		Background(colorBgHover).
		Foreground(colorYellow).
		Width(w)
}

func styleInputActive(w int) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(colorPurple).
		Width(w).
		PaddingLeft(1)
}

func styleInputInactive(w int) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(colorDim).
		Width(w).
		PaddingLeft(1)
}

func divider(w int) string {
	if isArcadeMode {
		dots := ""
		for i := 0; i < w; i += 2 {
			dots += "·"
			if i+1 < w {
				dots += " "
			}
		}
		return styleDivider.Render(dots)
	}
	return styleDivider.Render(strings.Repeat("─", w))
}
