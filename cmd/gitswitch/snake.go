package main

import (
	"github.com/aksisonline/gitswitch/internal/snake"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var snakeCmd = &cobra.Command{
	Use:    "snake",
	Hidden: true,
	Args:   cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := tea.NewProgram(snake.NewModel(), tea.WithAltScreen()).Run()
		return err
	},
}
