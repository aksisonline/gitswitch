package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// truncate shortens s to maxWidth visible cells, adding "…" if trimmed.
// ANSI-aware so styled (Render()'d) titles don't get cut mid-escape-sequence.
func truncate(s string, maxWidth int) string {
	if lipgloss.Width(s) <= maxWidth {
		return s
	}
	if maxWidth <= 0 {
		return ""
	}
	if maxWidth == 1 {
		return "…"
	}
	return ansi.Truncate(s, maxWidth, "…")
}

// padTo pads or truncates s to exactly n terminal columns.
func padTo(s string, n int) string {
	if n <= 0 {
		return ""
	}
	return lipgloss.PlaceHorizontal(n, lipgloss.Left, truncate(s, n))
}

// toggleW is the fixed width of every toggle state, so flipping one never
// moves anything else on the row.
const toggleW = 8

// renderToggle renders a GUI-style toggle switch as a single inline string.
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

// renderToggleBusy is the toggle while its change is being applied.
func (m Model) renderToggleBusy() string {
	return lipgloss.NewStyle().
		Background(colorBgHover).
		Foreground(colorYellow).
		Bold(true).
		Render(padTo("   "+m.spinner(), toggleW))
}

// itemView is one row of the Utilities or Settings tab.
type itemView struct {
	title  string // plain text
	right  string // styled widget at the right edge (toggle, chip)
	detail string // styled second line (description or current value)
}

// drawItems renders a focusable list of items at the canvas density:
//
//	roomy   — a bordered box per item (title/widget, then detail)
//	compact — the same two lines without the box, a blank line between items
//	tight   — one line per item; only the focused item's detail is shown,
//	          on a single shared line under the list
//
// Each item is recorded as a zItem zone covering all of its rows.
func (m Model) drawItems(c *canvas, focus int, items []itemView) {
	cursor := lipgloss.NewStyle().Foreground(colorYellow).Bold(true).Render("❯")
	titleStyle := func(focused bool) lipgloss.Style {
		if focused {
			return lipgloss.NewStyle().Foreground(colorPurple).Bold(true)
		}
		return styleCurrentVal
	}

	if c.d == dRoomy {
		iw := c.w - 6 // indent(2) + │ + pad + content + pad + │
		if iw < 4 {
			iw = 4
		}
		for i, it := range items {
			focused := i == focus
			b := styleBrand
			indent := "  "
			if focused {
				b = lipgloss.NewStyle().Foreground(colorPurple).Bold(true)
				indent = cursor + " "
			}
			rule := strings.Repeat("─", iw+2)
			line1 := spreadFixed(titleStyle(focused).Render(it.title), it.right, iw)
			line2 := padTo(it.detail, iw)
			c.mark(zItem, i, strings.Join([]string{
				"  " + b.Render("┌"+rule+"┐"),
				indent + b.Render("│") + " " + line1 + " " + b.Render("│"),
				"  " + b.Render("│") + " " + line2 + " " + b.Render("│"),
				"  " + b.Render("└"+rule+"┘"),
			}, "\n"))
		}
		return
	}

	iw := c.w - 5 // cursor(2) + indent(2) … right margin(1)
	for i, it := range items {
		focused := i == focus
		lead := "  "
		if focused {
			lead = cursor + " "
		}
		line1 := lead + "  " + spreadFixed(titleStyle(focused).Render(it.title), it.right, iw)
		if c.d == dTight {
			c.mark(zItem, i, line1)
			continue
		}
		c.mark(zItem, i, line1+"\n"+"    "+truncate(it.detail, iw))
		if i < len(items)-1 {
			c.blank()
		}
	}
	if c.d == dTight && focus >= 0 && focus < len(items) {
		c.add("    " + truncate(items[focus].detail, iw))
	}
}

// spreadFixed lays out title (left) and right widget in exactly w cells,
// truncating the title if the two would collide.
func spreadFixed(title, right string, w int) string {
	rw := lipgloss.Width(right)
	t := truncate(title, max(0, w-rw-1))
	pad := max(0, w-lipgloss.Width(t)-rw)
	return t + strings.Repeat(" ", pad) + right
}

// choiceView is one button-like option (setup wizard, dialogs).
type choiceView struct {
	icon  string // styled, single cell
	label string // plain text
}

// drawChoices renders a vertical list of choices, each recorded as a zChoice
// zone. Roomy and compact draw bordered buttons; tight draws one highlighted
// row per choice.
func (m Model) drawChoices(c *canvas, focus int, choices []choiceView) {
	accent := colorPurple
	if m.arcadeMode {
		accent = colorYellow
	}
	if c.d == dTight {
		for i, ch := range choices {
			row := "  " + ch.icon + "  " + ch.label
			if i == focus {
				bar := lipgloss.NewStyle().Background(colorBgHover).Foreground(colorYellow).Bold(true)
				c.mark(zChoice, i, "  "+bar.Render(padTo("❯ "+ansi.Strip(ch.icon)+"  "+ch.label, min(44, c.w-4))))
			} else {
				c.mark(zChoice, i, "  "+styleItemInactive.Render(row))
			}
		}
		return
	}
	bw := min(44, c.w-4) // box outer width
	if bw < 12 {
		bw = 12
	}
	for i, ch := range choices {
		border, label := styleBrand, styleItemInactive
		lead := "  "
		if i == focus {
			border = lipgloss.NewStyle().Foreground(accent).Bold(true)
			label = lipgloss.NewStyle().Foreground(colorWhite).Bold(true)
			lead = lipgloss.NewStyle().Foreground(colorYellow).Bold(true).Render("❯") + " "
		}
		rule := strings.Repeat("─", bw-2)
		inner := padTo("  "+ch.icon+"  "+label.Render(ch.label), bw-2)
		c.mark(zChoice, i, strings.Join([]string{
			"  " + border.Render("╭"+rule+"╮"),
			lead + border.Render("│") + inner + border.Render("│"),
			"  " + border.Render("╰"+rule+"╯"),
		}, "\n"))
		if i < len(choices)-1 {
			c.gap()
		}
	}
}

// buttonRow renders a primary and a secondary inline button. Their labels are
// what the mouse handler hit-tests, so keep them unique on screen.
func buttonRow(primary, secondary string) string {
	yes := lipgloss.NewStyle().
		Background(colorGreen).
		Foreground(lipgloss.Color("0")).
		Bold(true).
		Padding(0, 2).
		Render(primary)
	no := lipgloss.NewStyle().
		Foreground(colorDim).
		Padding(0, 2).
		Render(secondary)
	return "  " + yes + "  " + no
}
