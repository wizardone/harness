package main

import (
	"github.com/charmbracelet/lipgloss"
)

// Terminal styles; colours fall back to plain text when output is not a terminal.
var (
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")).Border(lipgloss.RoundedBorder()).Padding(0, 1)
	promptStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	errorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	hintStyle   = lipgloss.NewStyle().Faint(true)
)

// promptMarker starts the prompt line and marks the highlighted option.
const promptMarker = "> "
