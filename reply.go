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
	reply providers.Reply
	err    error
}

// send returns a command that sends the prompt to the model and reports the reply.
//   - client: sends the prompt to the model.
//   - prompt: the text to send.
func send(client providers.Client, prompt string) tea.Cmd {
	return func() tea.Msg {
		reply, err := client.Send(context.Background(), prompt)
		return replyMsg{reply: reply, err: err}
	}
}

// replyText returns the reply rendered as Markdown with a faint token count underneath,
// or the styled error when sending failed.
// The reply is shown as plain text if it cannot be rendered.
//   - reply: the outcome of sending a prompt.
func replyText(reply replyMsg) string {
	if reply.err != nil {
		return errorStyle.Render(fmt.Sprintf("error: %v", reply.err))
	}
	text := reply.reply.Text
	if rendered, err := glamour.Render(text, markdownStyle); err == nil {
		text = rendered
	}
	usage := fmt.Sprintf("%d in · %d out", reply.reply.InputTokens, reply.reply.OutputTokens)
	return text + "\n" + hintStyle.Render(usage)
}
