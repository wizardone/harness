package main

import (
	"context"
	"fmt"
	"harness/providers"

	tea "github.com/charmbracelet/bubbletea"
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

// replyText returns the reply, or the styled error when sending failed.
//   - reply: the outcome of sending a prompt.
func replyText(reply replyMsg) string {
	if reply.err != nil {
		return errorStyle.Render(fmt.Sprintf("error: %v", reply.err))
	}
	return reply.answer
}
