package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// Density is how much vertical breathing room a screen is drawn with. Every
// screen is tried roomy first and re-drawn tighter until it fits the terminal,
// so a 24-row window gets the same layout as a 50-row one, just with less air.
const (
	dRoomy   = iota // blank separators, boxed items
	dCompact        // no optional blank lines, borderless items
	dTight          // one line per item, secondary copy dropped
)

// Below this the panel cannot hold its own tab strip — show a notice instead
// of a broken layout.
const (
	minTermW = 44
	minTermH = 12
)

// zoneKind names a clickable region. Screens record zones as they draw, and
// the mouse handler asks "what is on this row?" — so hit-testing can never
// drift from what is actually rendered.
type zoneKind int

const (
	zTabs    zoneKind = iota + 1 // the tab strip
	zProfile                     // a row in the Accounts list
	zItem                        // a row/box in Utilities or Settings
	zImport                      // a detected identity in the setup wizard
	zChoice                      // a button-like choice (wizard, dialogs)
)

type zone struct {
	kind   zoneKind
	idx    int
	y0, y1 int // inclusive rows
}

// canvas is the inside of the panel: a list of lines that are each guaranteed
// to fit the panel's inner width. One canvas line is exactly one screen row —
// nothing is left for lipgloss to re-wrap, which is what used to change panel
// heights (and shift the whole layout) when a long status message appeared.
type canvas struct {
	w     int // inner width in cells
	d     int // density
	lines []string
	zones []zone
}

func newCanvas(pw, d int) *canvas {
	w := pw - 2
	if w < 1 {
		w = 1
	}
	return &canvas{w: w, d: d}
}

// add appends s (which may span several lines), clipping each to the width.
func (c *canvas) add(s string) {
	for _, l := range strings.Split(s, "\n") {
		c.lines = append(c.lines, clip(l, c.w))
	}
}

// blank appends an empty line at every density.
func (c *canvas) blank() { c.lines = append(c.lines, "") }

// gap appends an empty line only when there is room to spare.
func (c *canvas) gap() {
	if c.d == dRoomy {
		c.blank()
	}
}

// mark appends s and records the rows it occupies as a clickable zone.
func (c *canvas) mark(kind zoneKind, idx int, s string) {
	y0 := len(c.lines)
	c.add(s)
	c.zones = append(c.zones, zone{kind: kind, idx: idx, y0: y0, y1: len(c.lines) - 1})
}

// footer appends the divider and key hints that close every screen.
func (c *canvas) footer(m Model, pairs [][2]string) {
	c.add(divider(c.w + 2))
	c.add(m.footerKeys(c.w+2, pairs))
}

// clip shortens a (possibly styled) line to w cells.
func clip(s string, w int) string {
	if lipgloss.Width(s) <= w {
		return s
	}
	return ansi.Truncate(s, w, "…")
}

// spread puts left and right on one line w cells wide. If both don't fit,
// right is dropped — it is always the optional half.
func spread(left, right string, w int) string {
	gap := w - lipgloss.Width(left) - lipgloss.Width(right)
	if right == "" || gap < 1 {
		return left
	}
	return left + strings.Repeat(" ", gap) + right
}

// screen is a fully rendered body: the bordered panel plus anything above it.
type screen struct {
	lines []string
	zones []zone
}

func (s screen) height() int { return len(s.lines) }

func (s screen) width() int {
	w := 0
	for _, l := range s.lines {
		if lw := lipgloss.Width(l); lw > w {
			w = lw
		}
	}
	return w
}

func (s screen) String() string { return strings.Join(s.lines, "\n") }

// zoneAt returns the zone covering row y, if any.
func (s screen) zoneAt(y int) (zone, bool) {
	for _, z := range s.zones {
		if y >= z.y0 && y <= z.y1 {
			return z, true
		}
	}
	return zone{}, false
}

// plain returns the ANSI-stripped text of row y ("" when out of range).
func (s screen) plain(y int) string {
	if y < 0 || y >= len(s.lines) {
		return ""
	}
	return ansi.Strip(s.lines[y])
}

// finishPanel wraps the canvas in the panel border. pre holds lines that sit
// above the border (arcade mode's score line).
func finishPanel(pw int, c *canvas, pre ...string) screen {
	panel := stylePanelBorder(pw).Render(strings.Join(c.lines, "\n"))
	lines := append(append([]string{}, pre...), strings.Split(panel, "\n")...)
	off := len(pre) + 1 // + top border
	zones := make([]zone, len(c.zones))
	for i, z := range c.zones {
		zones[i] = zone{kind: z.kind, idx: z.idx, y0: z.y0 + off, y1: z.y1 + off}
	}
	return screen{lines: lines, zones: zones}
}

func (m Model) tooSmall() bool {
	return (m.width > 0 && m.width < minTermW) || (m.height > 0 && m.height < minTermH)
}

// tooSmallScreen is the notice shown when the window cannot hold the panel.
func (m Model) tooSmallScreen() screen {
	msg := []string{
		styleTitle.Render("gitswitch"),
		"",
		styleItemDim.Render("Window too small."),
		styleItemDim.Render("Make it larger, or use"),
		styleFooterKey.Render("gitswitch <nickname>"),
	}
	var lines []string
	for _, l := range msg {
		if m.width > 0 {
			l = clip(l, m.width)
		}
		lines = append(lines, l)
	}
	if m.height > 0 && len(lines) > m.height {
		lines = lines[:m.height]
	}
	return screen{lines: lines}
}

// screen renders the current state at the roomiest density that fits.
func (m Model) screen() screen {
	m.applyColors()
	if m.tooSmall() {
		return m.tooSmallScreen()
	}
	var s screen
	for d := dRoomy; d <= dTight; d++ {
		s = m.screenAt(d)
		if m.height <= 0 || s.height() <= m.height {
			break
		}
	}
	return s
}

// origin is where the screen's top-left corner lands in the terminal.
//
// Horizontally everything is centered. Vertically, dialogs are centered too —
// but the three tabs share one anchor row, computed from the tallest of them.
// Tabs have different heights, and centering each one separately made the tab
// strip jump up and down under the pointer every time you changed tab.
func (m Model) origin(s screen) (top, left int) {
	if m.width > 0 {
		left = (m.width - s.width()) / 2
	}
	h := s.height()
	if m.onTabs() && !m.tooSmall() {
		for i := 0; i < 3; i++ {
			mm := m
			mm.state, mm.tabIndex = StateList, i
			if th := mm.screen().height(); th > h {
				h = th
			}
		}
	}
	if m.height > 0 {
		top = (m.height - h) / 2
	}
	if top < 0 {
		top = 0
	}
	if left < 0 {
		left = 0
	}
	return top, left
}

// headGap reports whether the tabs get a blank line under the header. It
// depends only on the terminal height — never on the density a particular tab
// ended up at — so the tab strip sits on the same row on all three tabs.
func (m Model) headGap() bool { return m.height == 0 || m.height >= 20 }

// onTabs reports whether the tabbed main screen is showing.
func (m Model) onTabs() bool {
	return m.state == StateList || m.state == StateSelectFlash
}

// bodyView returns the rendered body as one string.
func (m Model) bodyView() string { return m.screen().String() }
