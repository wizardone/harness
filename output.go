package main

import (
	"github.com/charmbracelet/glamour/ansi"
	"github.com/charmbracelet/glamour/styles"
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
	// sentPromptStyle marks prompts already sent, so they stand out when scrolling back.
	sentPromptStyle = lipgloss.NewStyle().Reverse(true).Bold(true).Padding(0, 1)
)

// markdownStyles returns the monochrome style replies are rendered with: the plain notty style,
// with bold, italics, headings and inline code shown as text styles instead of Markdown symbols.
func markdownStyles() ansi.StyleConfig {
	style := styles.NoTTYStyleConfig
	bold, italic := true, true
	style.Strong = ansi.StylePrimitive{Bold: &bold}
	style.Emph = ansi.StylePrimitive{Italic: &italic}
	style.Heading.Bold = &bold
	style.H1.Prefix, style.H2.Prefix, style.H3.Prefix = "", "", ""
	style.Code = ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Bold: &bold}}
	return style
}

// promptMarker starts the prompt line and marks the highlighted option.
const promptMarker = "> "
