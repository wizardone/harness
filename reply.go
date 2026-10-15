package main

import (
	"context"
	"fmt"
	"harness/providers"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

// replyMsg carries the model's reply to a prompt, or the error from sending it.
type replyMsg struct {
	answer string
	err    error
}

// send returns a command that sends the prompt to the model and reports the reply.
//   - client: sends the prompt to the model.
//   - prompt: the text to send.
func send(client providers.Client, prompt string) tea.Cmd {
	return func() tea.Msg {
		answer, err := client.Send(context.Background(), prompt)
		return replyMsg{answer: answer, err: err}
	}
}

// replyText returns the reply rendered as Markdown, or the styled error when sending failed.
// The reply is shown as plain text if it cannot be rendered.
//   - reply: the outcome of sending a prompt.
func replyText(reply replyMsg) string {
	if reply.err != nil {
		return errorStyle.Render(fmt.Sprintf("error: %v", reply.err))
	}
	if rendered, err := glamour.Render(reply.answer, markdownStyle); err == nil {
		return rendered
	}
	return reply.answer
}
