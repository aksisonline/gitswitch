package snake

import (
	"fmt"
	"strings"
	"time"

	"github.com/aksisonline/gitswitch/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg time.Time

const stepInterval = 220 * time.Millisecond

// Model renders and controls a snake game in a terminal.
type Model struct {
	game   *Game
	paused bool
	width  int
	height int
	theme  tui.ThemeColors
}

// NewModel starts a new game model.
func NewModel(theme tui.ThemeColors) Model {
	return Model{
		game:   NewGame(time.Now().UnixNano()),
		width:  80,
		height: 24,
		theme:  theme,
	}
}

func (m Model) Init() tea.Cmd {
	return tick()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil
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
				m.game = NewGame(time.Now().UnixNano())
				m.paused = false
				return m, tick()
			}
		}
	case tickMsg:
		if m.paused || m.game.Over {
			return m, nil
		}
		m.game.Step()
		if !m.game.Over {
			return m, tick()
		}
		return m, nil
	}
	// Only tick messages schedule another tick. Scheduling from key handlers
	// creates concurrent timers, which makes the game speed up after a restart.
	return m, nil
}

func (m Model) View() string {
	var board strings.Builder
	border := lipgloss.NewStyle().Foreground(m.theme.Primary)
	head := lipgloss.NewStyle().Foreground(m.theme.Accent).Bold(true)
	body := lipgloss.NewStyle().Foreground(m.theme.Primary)
	food := lipgloss.NewStyle().Foreground(m.theme.Highlight).Bold(true)
	title := lipgloss.NewStyle().Foreground(m.theme.Primary).Bold(true)
	score := lipgloss.NewStyle().Foreground(m.theme.Accent).Bold(true)
	footer := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	alert := lipgloss.NewStyle().Foreground(m.theme.Error).Bold(true)

	board.WriteString(border.Render("┌" + strings.Repeat("──", Width) + "┐"))
	board.WriteByte('\n')
	for y := 0; y < Height; y++ {
		board.WriteString(border.Render("│"))
		for x := 0; x < Width; x++ {
			point := Point{x, y}
			switch {
			case point == m.game.Food:
				board.WriteString(food.Render(" ●"))
			case point == m.game.Body[0]:
				board.WriteString(head.Render(" █"))
			case contains(m.game.Body[1:], point):
				board.WriteString(body.Render(" ▓"))
			default:
				board.WriteString("  ")
			}
		}
		board.WriteString(border.Render("│"))
		board.WriteByte('\n')
	}
	board.WriteString(border.Render("└" + strings.Repeat("──", Width) + "┘"))

	contentWidth := Width * 2
	status := score.Render(fmt.Sprintf("Score: %d", m.game.Score))
	help := footer.Render("Arrow keys/WASD move • Space/p pauses • q quits")
	switch {
	case m.game.Won:
		help = alert.Render("You win! Press r to play again, or q to quit.")
	case m.game.Over:
		help = alert.Render("Game over. Press r to play again, or q to quit.")
	case m.paused:
		help = footer.Render("Paused. Press space or p to resume; q quits.")
	}
	content := strings.Join([]string{
		lipgloss.Place(contentWidth, 1, lipgloss.Center, lipgloss.Center, title.Render("gitswitch snake")),
		board.String(),
		lipgloss.Place(contentWidth, 1, lipgloss.Center, lipgloss.Center, status),
		lipgloss.Place(contentWidth, 1, lipgloss.Center, lipgloss.Center, help),
	}, "\n")
	panel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Primary).
		Padding(1, 1).
		Render(content)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, panel)
}

func tick() tea.Cmd {
	return tea.Tick(stepInterval, func(t time.Time) tea.Msg {
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
