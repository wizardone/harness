package main

import (
	"context"
	"fmt"
	"harness/providers"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

// replyMsg carries the model's reply to a prompt, or the error from sending it.
type replyMsg struct {
	reply providers.Reply
	err    error
	prompt string
}

// send returns a command that sends the conversation with the new prompt to the model and reports the reply.
//   - client: sends the conversation to the model.
//   - history: the earlier prompts and replies, alternating; copied so the command has its own slice.
//   - prompt: the new text to send.
func send(client providers.Client, history []string, prompt string) tea.Cmd {
	return func() tea.Msg {
		reply, err := client.Send(context.Background(), append(slices.Clone(history), prompt))
		return replyMsg{reply: reply, err: err, prompt: prompt}
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
