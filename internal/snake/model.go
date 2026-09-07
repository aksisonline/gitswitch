package snake

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg time.Time

// Model renders and controls a snake game in a terminal.
type Model struct {
	game   *Game
	paused bool
}

// NewModel starts a new game model.
func NewModel() Model {
	return Model{game: NewGame(time.Now().UnixNano())}
}

func (m Model) Init() tea.Cmd {
	return tick()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "up", "k", "w":
			m.game.Turn(Up)
		case "down", "j", "s":
			m.game.Turn(Down)
		case "left", "h", "a":
			m.game.Turn(Left)
		case "right", "l", "d":
			m.game.Turn(Right)
		case " ", "p":
			if !m.game.Over {
				m.paused = !m.paused
			}
		case "r":
			if m.game.Over {
				m = NewModel()
			}
		}
	case tickMsg:
		if !m.paused && !m.game.Over {
			m.game.Step()
		}
	}
	return m, tick()
}

func (m Model) View() string {
	var board strings.Builder
	board.WriteString("\n  ┌")
	board.WriteString(strings.Repeat("──", Width))
	board.WriteString("┐\n")
	for y := 0; y < Height; y++ {
		board.WriteString("  │")
		for x := 0; x < Width; x++ {
			point := Point{x, y}
			switch {
			case point == m.game.Food:
				board.WriteString(" ●")
			case point == m.game.Body[0]:
				board.WriteString(" █")
			case contains(m.game.Body[1:], point):
				board.WriteString(" ▓")
			default:
				board.WriteString("  ")
			}
		}
		board.WriteString("│\n")
	}
	board.WriteString("  └")
	board.WriteString(strings.Repeat("──", Width))
	board.WriteString("┘\n")
	board.WriteString(fmt.Sprintf("\n  gitswitch snake  •  Score: %d\n", m.game.Score))
	switch {
	case m.game.Won:
		board.WriteString("  You win! Press r to play again, or q to quit.\n")
	case m.game.Over:
		board.WriteString("  Game over. Press r to play again, or q to quit.\n")
	case m.paused:
		board.WriteString("  Paused. Press space or p to resume; q quits.\n")
	default:
		board.WriteString("  Arrow keys or WASD to move • Space/p pauses • q quits\n")
	}
	return board.String()
}

func tick() tea.Cmd {
	return tea.Tick(120*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func contains(points []Point, target Point) bool {
	for _, point := range points {
		if point == target {
			return true
		}
	}
	return false
}
