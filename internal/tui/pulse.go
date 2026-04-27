package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type arcadeTickMsg time.Time

func arcadeTickCmd() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg {
		return arcadeTickMsg(t)
	})
}
