package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// truncate shortens s to maxRunes visible runes, adding "…" if trimmed.
func truncate(s string, maxRunes int) string {
	r := []rune(s)
	if len(r) <= maxRunes {
		return s
	}
	if maxRunes <= 1 {
		return "…"
	}
	return string(r[:maxRunes-1]) + "…"
}

// renderToggle renders a GUI-style toggle switch as a single inline string.
//
//	ON:  green filled  "  ● on  "
//	OFF: dim filled    " ○ off  "
func renderToggle(on bool) string {
	if on {
		return lipgloss.NewStyle().
			Background(colorGreen).
			Foreground(lipgloss.Color("#000000")).
			Bold(true).
			Render("  ● on  ")
	}
	return lipgloss.NewStyle().
		Background(colorBgHover).
		Foreground(colorDim).
		Render(" ○ off  ")
}

// renderItemBox renders a 4-line bordered item box used in the Utilities and
// Settings tabs. Each box occupies exactly 4 lines on screen (top border, 2
// content lines, bottom border) and is preceded by a single "\n" spacer that
// the caller should include.
//
// pw is the panel content width (from panelWidth()).
// The box is indented 2 spaces and expands to fill the panel width.
//
// line1 and line2 must each be exactly (pw-6) terminal columns wide.
// Use the renderItemLine helper to build them.
func renderItemBox(pw int, focused bool, disabled bool, line1, line2 string) string {
	dashW := pw - 6 // 2-space indent + ┌ + dashes + ┐ spans the usable inner width (pw-2)
	if dashW < 2 {
		dashW = 2
	}
	var b lipgloss.Style
	if focused && !disabled {
		b = lipgloss.NewStyle().Foreground(colorPurple).Bold(true)
	} else if disabled {
		b = lipgloss.NewStyle().Foreground(colorDim)
	} else {
		b = styleBrand
	}

	// Focused (and enabled) boxes get a yellow cursor arrow in place of the indent.
	indentTop, indentMid := "  ", "  "
	if focused && !disabled {
		arrow := lipgloss.NewStyle().Foreground(colorYellow).Bold(true).Render("❯")
		indentTop = arrow + " "
		indentMid = arrow + " "
	}

	top := indentTop + b.Render("┌"+strings.Repeat("─", dashW)+"┐")
	mid1 := indentMid + b.Render("│") + " " + line1 + " " + b.Render("│")
	mid2 := "  " + b.Render("│") + " " + line2 + " " + b.Render("│")
	bot := "  " + b.Render("└"+strings.Repeat("─", dashW)+"┘")

	return "\n" + top + "\n" + mid1 + "\n" + mid2 + "\n" + bot
}

// itemInnerW returns the writable character width inside a box built by
// renderItemBox for the given panel width.
// Formula: pw - 8 (2-space indent + left-│ + 1-space pad + content + 1-space pad
// + right-│, with the box outer width held to the usable inner width pw-2).
func itemInnerW(pw int) int {
	w := pw - 8
	if w < 0 {
		return 0
	}
	return w
}

// padTo pads or truncates s (by visual rune count) to exactly n terminal columns.
// Uses spaces for padding. ANSI codes in s are NOT counted — pass already-rendered
// strings only when measuring with lipgloss.Width separately.
func padTo(s string, n int) string {
	w := lipgloss.Width(s)
	if w >= n {
		return s
	}
	return s + strings.Repeat(" ", n-w)
}

// titleWithRight builds a line1 string: title left-aligned, right widget right-aligned,
// total width = innerW.
func titleWithRight(title, right string, innerW int) string {
	rw := lipgloss.Width(right)
	maxTitle := innerW - rw - 1
	if maxTitle < 0 {
		maxTitle = 0
	}
	t := truncate(title, maxTitle)
	pad := innerW - lipgloss.Width(t) - rw
	if pad < 0 {
		pad = 0
	}
	return t + strings.Repeat(" ", pad) + right
}
