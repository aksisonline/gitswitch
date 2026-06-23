package tui

import "github.com/charmbracelet/harmonica"

// AnimState holds a single harmonica spring + its current position/velocity.
// Create one per animated value; call Update() on each tea.Tick.
type AnimState struct {
	Pos float64
	Vel float64
	spr harmonica.Spring
}

// newSpring creates an AnimState with the given frequency and damping ratio.
//   - freq: angular frequency (higher = faster response)
//   - damp: damping ratio (1.0 = critically damped, <1 = bouncy, >1 = overdamped)
func newSpring(freq, damp float64) AnimState {
	return AnimState{spr: harmonica.NewSpring(harmonica.FPS(60), freq, damp)}
}

// Update advances the spring one tick toward target. Call on every TickMsg.
func (a *AnimState) Update(target float64) {
	a.Pos, a.Vel = a.spr.Update(a.Pos, a.Vel, target)
}

// Done reports whether the spring has settled (within 0.5 units of target).
func (a *AnimState) Done(target float64) bool {
	diff := a.Pos - target
	if diff < 0 {
		diff = -diff
	}
	return diff < 0.5 && a.Vel > -0.1 && a.Vel < 0.1
}

// Predefined spring constructors for each animation type.
// Parameters derived from the v0.2.0 TUI design spec (Section 04).

// NewTabSpring returns a spring for sliding between tabs.
// freq=4.0, damp=0.85 — snappy with slight overshoot.
func NewTabSpring() AnimState { return newSpring(4.0, 0.85) }

// NewCursorSpring returns a spring for the cursor row highlight.
// freq=10.0, damp=0.75 — fast, smooth on large jumps.
func NewCursorSpring() AnimState { return newSpring(10.0, 0.75) }

// NewInputFocusSpring returns a spring for input field focus expand.
// freq=6.0, damp=0.80 — slight overshoot draws attention.
func NewInputFocusSpring() AnimState { return newSpring(6.0, 0.80) }

// NewLogoRevealSpring returns a spring for boot logo slide-in.
// freq=3.0, damp=0.90 — slow, heavy, no bounce.
func NewLogoRevealSpring() AnimState { return newSpring(3.0, 0.90) }

// NewPacmanSpring returns a spring for pac-man X position.
// freq=2.5, damp=1.0 — critically damped, steady march, no overshoot.
func NewPacmanSpring() AnimState { return newSpring(2.5, 1.0) }

// NewSwitchFlashSpring returns a spring for the arcade switch bounce.
// freq=8.0, damp=0.55 — underdamped, bouncy, satisfying.
func NewSwitchFlashSpring() AnimState { return newSpring(8.0, 0.55) }
