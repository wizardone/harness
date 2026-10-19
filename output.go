package main

import (
	"github.com/charmbracelet/lipgloss"
)

// Terminal styles; monochrome apart from errors, and plain text when output is not a terminal.
var (
	headerStyle = lipgloss.NewStyle().Bold(true)
	promptStyle = lipgloss.NewStyle().Bold(true)
	errorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#e05252"))
	hintStyle   = lipgloss.NewStyle().Faint(true)
	// promptAreaStyle draws a line above and below the prompt so it stands out from the replies.
	promptAreaStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true, false)
)

// markdownStyle is the glamour style replies are rendered with.
const markdownStyle = "notty"

// promptMarker starts the prompt line and marks the highlighted option.
const promptMarker = "> "
