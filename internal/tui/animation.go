package tui

import (
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
)

// Motion in the normal (non-arcade) UI is driven by one frame ticker that only
// runs while something is actually moving — a spinner, the tab indicator, the
// launch shimmer — and stops itself the moment everything has settled. An
// idle gitswitch redraws nothing.

const (
	frameInterval = 33 * time.Millisecond // ~30fps
	spinnerEvery  = 3                     // frames per spinner glyph
	scanStepEvery = 8                     // frames per setup-scan row (~260ms)
	scanHoldFor   = 10                    // frames to rest on the finished scan
	flashFrames   = 14                    // how long a switched row glows
)

// frameMsg advances every running animation by one frame.
type frameMsg time.Time

func frameCmd() tea.Cmd {
	return tea.Tick(frameInterval, func(t time.Time) tea.Msg { return frameMsg(t) })
}

// arcadeTickMsg drives arcade mode's own cut-scenes (intro, flash, exit).
type arcadeTickMsg time.Time

func arcadeTickCmd() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg { return arcadeTickMsg(t) })
}

var spinnerGlyphs = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// reduceMotion turns off decorative movement (shimmer, sliding tab indicator,
// staged setup scan). Spinners stay — they are status, not decoration.
var reduceMotion = os.Getenv("GITSWITCH_REDUCE_MOTION") != ""

// AnimState is one damped spring: a value easing toward a target.
type AnimState struct {
	Pos float64
	Vel float64
	spr harmonica.Spring
}

func newSpring(freq, damp float64) AnimState {
	return AnimState{spr: harmonica.NewSpring(harmonica.FPS(30), freq, damp)}
}

// Update advances the spring one frame toward target.
func (a *AnimState) Update(target float64) {
	a.Pos, a.Vel = a.spr.Update(a.Pos, a.Vel, target)
}

// Done reports whether the spring has settled on target.
func (a *AnimState) Done(target float64) bool {
	diff := a.Pos - target
	if diff < 0 {
		diff = -diff
	}
	return diff < 0.5 && a.Vel > -0.5 && a.Vel < 0.5
}

// spinner returns the current spinner glyph.
func (m Model) spinner() string {
	return spinnerGlyphs[(m.frame/spinnerEvery)%len(spinnerGlyphs)]
}

// needsFrames reports whether anything on screen is still moving.
func (m Model) needsFrames() bool {
	return m.busy != "" || m.busyItem != 0 || m.shimmer > 0 || m.flashLeft > 0 ||
		m.tabAnim || (m.state == StateWizardDetect && m.scanning)
}

// animate starts the frame ticker if something needs it and it isn't running.
func (m *Model) animate() tea.Cmd {
	if m.ticking || !m.needsFrames() {
		return nil
	}
	m.ticking = true
	return frameCmd()
}

// onFrame advances every running animation by one frame.
func (m Model) onFrame() (Model, tea.Cmd) {
	m.frame++
	if m.shimmer > 0 {
		m.shimmer++
		if m.shimmer > len(brandName)+6 {
			m.shimmer = 0
		}
	}
	if m.flashLeft > 0 {
		m.flashLeft--
	}
	if m.tabAnim {
		xs, ws := m.tabGeometry()
		tx, tw := float64(xs[m.tabIndex]), float64(ws[m.tabIndex])
		m.tabX.Update(tx)
		m.tabW.Update(tw)
		if m.tabX.Done(tx) && m.tabW.Done(tw) {
			m.tabAnim = false
		}
	}
	if m.state == StateWizardDetect && m.scanning {
		if m.scanStep < scanRows && m.frame%scanStepEvery == 0 {
			m.scanStep++
		}
		if m.scanStep >= scanRows && m.scanReady {
			m.scanHold++
			if m.scanHold >= scanHoldFor {
				m.finishScan()
			}
		}
	}
	if m.needsFrames() {
		return m, frameCmd()
	}
	m.ticking = false
	return m, nil
}

// setTab switches tab, sliding the indicator from the old tab to the new one.
func (m *Model) setTab(i int) {
	if i == m.tabIndex {
		return
	}
	xs, ws := m.tabGeometry()
	from := m.tabIndex
	m.tabIndex = i
	m.statusMsg = ""
	if reduceMotion || m.width == 0 {
		return
	}
	if !m.tabAnim {
		m.tabX, m.tabW = newSpring(16.0, 1.0), newSpring(16.0, 1.0) // critically damped: quick, no wobble
		m.tabX.Pos, m.tabW.Pos = float64(xs[from]), float64(ws[from])
	}
	m.tabAnim = true
}

// setBusy shows a spinner with label in the status slot until clearBusy.
func (m *Model) setBusy(label string) {
	m.busy = label
	m.statusMsg = ""
}

func (m *Model) clearBusy() {
	m.busy = ""
	m.busyItem = 0
}
