package tui

import "github.com/charmbracelet/lipgloss"

const minPanelWidth = 56
const maxPanelWidth = 120

var (
	colorPurple  = lipgloss.Color("141")
	colorGreen   = lipgloss.Color("84")
	colorYellow  = lipgloss.Color("226")
	colorDim     = lipgloss.Color("241")
	colorRed     = lipgloss.Color("196")
	colorWhite   = lipgloss.Color("255")
	colorBgHover = lipgloss.Color("237")

	// Width-independent styles
	styleTitle = lipgloss.NewStyle().Bold(true).Foreground(colorPurple)
	styleBrand = lipgloss.NewStyle().Foreground(colorDim)

	styleBrandLink = lipgloss.NewStyle().
			Foreground(colorGreen).
			Underline(true)

	styleCurrent    = lipgloss.NewStyle().Foreground(colorDim)
	styleCurrentVal = lipgloss.NewStyle().Foreground(colorWhite).Bold(true)
	styleCheckmark  = lipgloss.NewStyle().Foreground(colorGreen)
	styleItemDim    = lipgloss.NewStyle().Foreground(colorDim)
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

// Width-dependent styles — call these with the computed panel width.

func stylePanelBorder(w int) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
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
	line := ""
	for i := 0; i < w; i++ {
		line += "─"
	}
	return styleDivider.Render(line)
}
