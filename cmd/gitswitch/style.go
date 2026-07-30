package main

import "github.com/charmbracelet/lipgloss"

// Palette mirrors internal/install/wizard.go and internal/tui/styles.go.
// lipgloss auto-detects whether stdout is a real terminal and strips ANSI
// codes on its own when it isn't, so piping these commands stays clean.
var (
	cPrimary   = lipgloss.Color("141") // purple
	cAccent    = lipgloss.Color("84")  // green
	cHighlight = lipgloss.Color("226") // yellow
	cDim       = lipgloss.Color("241") // gray
	cRed       = lipgloss.Color("196") // red
)

func styleOK(s string) string   { return lipgloss.NewStyle().Foreground(cAccent).Render(s) }
func styleWarn(s string) string { return lipgloss.NewStyle().Foreground(cHighlight).Render(s) }
func styleErr(s string) string  { return lipgloss.NewStyle().Foreground(cRed).Render(s) }
func styleDim(s string) string  { return lipgloss.NewStyle().Foreground(cDim).Render(s) }
