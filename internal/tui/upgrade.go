package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// notesViewportHeight calculates how many pre-wrapped visual lines fit inside
// the What's New / Upgrade modal without overflowing the terminal vertically.
func notesViewportHeight(termHeight int) int {
	if termHeight <= 0 {
		return 12
	}
	// Panel chrome (top/bottom borders, header, spacers, scroll hint, divider,
	// footer, and outer margin) takes ~12 lines. Cap at 16 lines so the modal
	// stays compact and scannable even on tall terminals.
	h := termHeight - 12
	if h > 16 {
		h = 16
	}
	if h < 4 {
		h = 4
	}
	return h
}

// cleanReleaseNotes strips noisy sections from raw GitHub release bodies
// (canary boilerplate banners, redundant headers, install scripts, and
// auto-generated PR/changelog footers) so the TUI splash stays concise.
func cleanReleaseNotes(raw string) string {
	raw = strings.ReplaceAll(raw, "\r\n", "\n")

	// When a release body includes both a cumulative [Unreleased] section and
	// GitHub's tag-to-tag "## What's Changed" section (e.g. canary releases),
	// keep only "## What's Changed" so only the most recent update is shown.
	for _, marker := range []string{"\n## What's Changed", "## What's Changed"} {
		if idx := strings.Index(raw, marker); idx >= 0 {
			raw = raw[idx:]
			break
		}
	}

	var kept []string
	skipSection := false
	inCodeBlock := false
	seenVersionHeader := false

	for _, line := range strings.Split(raw, "\n") {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}
		if inCodeBlock {
			if !skipSection {
				kept = append(kept, line)
			}
			continue
		}

		// Drop canary warning blockquotes at the top of beta releases.
		if strings.HasPrefix(trimmed, ">") && (strings.Contains(trimmed, "canary/beta release") || strings.Contains(trimmed, "gitswitch stable")) {
			continue
		}
		// Drop Full Changelog link line.
		if strings.HasPrefix(trimmed, "**Full Changelog**") || strings.HasPrefix(trimmed, "Full Changelog:") {
			continue
		}

		// Check headings for sections to skip or flatten.
		if strings.HasPrefix(trimmed, "#") {
			heading := strings.TrimSpace(strings.TrimLeft(trimmed, "#"))
			lower := strings.ToLower(heading)

			// Version heading like "## v0.4.4 — 2026-08-20" or "## [v0.4.4]":
			// skip the first one (redundant with modal title) and stop if a second
			// older version section appears.
			if strings.HasPrefix(lower, "v0.") || strings.HasPrefix(lower, "v1.") || strings.HasPrefix(lower, "[v") {
				if seenVersionHeader {
					break
				}
				seenVersionHeader = true
				skipSection = false
				continue
			}

			switch {
			case lower == "what's in this beta" || lower == "what's new" || lower == "what's changed":
				// Redundant with the modal title "What's New in vX.Y.Z"
				skipSection = false
				continue
			case strings.HasPrefix(lower, "install "):
				// User is already running the binary in the TUI; skip curl/ps1 blocks.
				skipSection = true
				continue
			default:
				skipSection = false
			}
		}

		if skipSection {
			continue
		}
		kept = append(kept, line)
	}
	return strings.TrimSpace(strings.Join(kept, "\n"))
}

// renderNotes cleans and pre-wraps markdown release notes into visual terminal
// lines so that 1 element == 1 screen row. This keeps the scroll viewport and
// panel height exact regardless of paragraph length.
func renderNotes(raw string, pw int) []string {
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(colorGreen)
	bulletStyle := lipgloss.NewStyle().Foreground(colorPurple).Bold(true)
	textStyle := lipgloss.NewStyle().Foreground(colorWhite)
	dimStyle := lipgloss.NewStyle().Foreground(colorDim)

	cleaned := cleanReleaseNotes(raw)
	cleaned = strings.ReplaceAll(cleaned, "**", "")

	// Usable inner width inside stylePanelBorder(pw) is pw - 2.
	// Leave 2 leading spaces + 1 trailing margin = pw - 5 for text.
	wrapW := pw - 5
	if wrapW < 20 {
		wrapW = 20
	}

	var lines []string
	prevBlank := true // prevent leading blank lines

	appendBlank := func() {
		if !prevBlank && len(lines) > 0 {
			lines = append(lines, "")
			prevBlank = true
		}
	}

	for _, line := range strings.Split(cleaned, "\n") {
		trimmed := strings.TrimSpace(line)
		switch {
		case trimmed == "":
			appendBlank()

		case strings.HasPrefix(trimmed, "### "), strings.HasPrefix(trimmed, "## "), strings.HasPrefix(trimmed, "# "):
			heading := strings.TrimSpace(strings.TrimLeft(trimmed, "#"))
			if heading == "" {
				continue
			}
			appendBlank()
			for _, wl := range wrapNoteLine(heading, wrapW) {
				lines = append(lines, "  "+headerStyle.Render(wl))
			}
			prevBlank = false

		case strings.HasPrefix(trimmed, "* "), strings.HasPrefix(trimmed, "- "):
			item := strings.TrimSpace(trimmed[2:])
			// Strip noisy PR author/URL suffixes from auto-generated GitHub bullets.
			if idx := strings.Index(item, " by @"); idx > 0 && strings.Contains(item[idx:], "http") {
				item = strings.TrimSpace(item[:idx])
			}
			for _, marker := range []string{" [minor]", " [major]", " [patch]"} {
				item = strings.TrimSuffix(item, marker)
			}
			wrapped := wrapNoteLine(item, wrapW-2)
			for i, wl := range wrapped {
				if i == 0 {
					lines = append(lines, "  "+bulletStyle.Render("•")+" "+textStyle.Render(wl))
				} else {
					lines = append(lines, "    "+textStyle.Render(wl))
				}
			}
			prevBlank = false

		default:
			for _, wl := range wrapNoteLine(trimmed, wrapW) {
				lines = append(lines, "  "+dimStyle.Render(wl))
			}
			prevBlank = false
		}
	}

	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

// wrapNoteLine word-wraps s to width columns (ANSI-aware, breaking long URLs/words if needed).
func wrapNoteLine(s string, width int) []string {
	if width <= 0 {
		return []string{s}
	}
	wrapped := ansi.Wordwrap(s, width, "")
	var out []string
	for _, l := range strings.Split(wrapped, "\n") {
		for lipgloss.Width(l) > width {
			part := ansi.Truncate(l, width, "")
			if part == "" {
				break
			}
			out = append(out, part)
			l = strings.TrimPrefix(l, part)
		}
		out = append(out, l)
	}
	return out
}
