package main

import (
	"fmt"
	"strings"
)

// renderSentPrompt returns a sent prompt as a numbered ledger entry under a rule, so turns stand out when scrolling back.
//   - prompt: the text that was sent.
//   - turn: the number of this exchange, starting at 1.
//   - width: the terminal width the rule spans.
func renderSentPrompt(prompt string, turn, width int) string {
	return renderTurnRule(width) + "\n" + renderLedger(fmt.Sprintf("%02d", turn), promptStyle.Render(prompt))
}

// renderReply returns the reply, or its error, in the ledger under the prompt it answers,
// with no turn number of its own.
//   - reply: the outcome of sending a prompt.
func renderReply(reply replyMsg) string {
	return renderLedger("", replyText(reply))
}

// renderPromptArea returns the prompt with a line above and below it.
//   - prompt: the rendered prompt line.
//   - width: the terminal width the lines span; 0 until the size is known, when they span the prompt.
func renderPromptArea(prompt string, width int) string {
	return promptAreaStyle.Width(width).Render(prompt)
}

// renderOptions returns the options one per line, highlighting the one under the cursor.
//   - cursor: the index of the highlighted option.
func renderOptions(cursor int) string {
	lines := make([]string, len(options))
	for i, option := range options {
		lines[i] = hintStyle.Render("  " + option)
		if i == cursor {
			lines[i] = promptStyle.Render(promptMarker + option)
		}
	}
	return strings.Join(lines, "\n")
}

// ledgerMargin is how many characters the turn number column on the left takes.
const ledgerMargin = 4

// renderLedger returns the text with a faint margin and bar on the left of every line.
//   - label: shown in the margin of the first line only, such as the turn number; may be empty.
//   - text: the lines to put in the ledger.
func renderLedger(label, text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = hintStyle.Render(fmt.Sprintf("%*s │ ", ledgerMargin, label)) + line
		label = ""
	}
	return strings.Join(lines, "\n")
}

// renderTurnRule returns the faint line that separates one turn from the next, crossing the margin bar.
//   - width: the terminal width the line spans; at least 20 characters are drawn after the bar.
func renderTurnRule(width int) string {
	rest := max(width-ledgerMargin-2, 20)
	return hintStyle.Render(strings.Repeat("─", ledgerMargin+1) + "┼" + strings.Repeat("─", rest))
}
