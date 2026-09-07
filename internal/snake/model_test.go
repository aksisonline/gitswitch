package snake

import (
	"testing"

	"github.com/aksisonline/gitswitch/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func TestInputDoesNotScheduleAdditionalTicks(t *testing.T) {
	m := NewModel(tui.ThemeColorsFor(0, false))

	next, cmd := m.Update(tea.KeyMsg{Type: tea.KeyDown})

	if cmd != nil {
		t.Fatal("input must not schedule a tick; the existing game tick owns the timer")
	}
	if next.(Model).game.Direction != Down {
		t.Fatalf("direction = %#v, want %#v", next.(Model).game.Direction, Down)
	}
}

func TestRestartSchedulesOneNewTick(t *testing.T) {
	m := NewModel(tui.ThemeColorsFor(0, false))
	m.game.Over = true

	next, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

	if cmd == nil {
		t.Fatal("restart must schedule a game tick")
	}
	if next.(Model).game.Over {
		t.Fatal("restart must begin a new game")
	}
}

func TestPausedTickStopsTimer(t *testing.T) {
	m := NewModel(tui.ThemeColorsFor(0, false))
	m.paused = true

	_, cmd := m.Update(tickMsg{})

	if cmd != nil {
		t.Fatal("a paused game must not schedule another tick")
	}
}
