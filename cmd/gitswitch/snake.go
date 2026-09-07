package main

import (
	"github.com/aksisonline/gitswitch/internal/snake"
	"github.com/aksisonline/gitswitch/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var snakeCmd = &cobra.Command{
	Use:    "snake",
	Hidden: true,
	Args:   cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		prefs, err := store.LoadPrefs()
		if err != nil {
			return err
		}
		theme := tui.ThemeColorsFor(prefs.ColorTheme, prefs.ArcadeMode)
		_, err = tea.NewProgram(snake.NewModel(theme), tea.WithAltScreen()).Run()
		return err
	},
}
