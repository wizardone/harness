package main

import (
	"context"
	"harness/providers"

	"github.com/charmbracelet/glamour"
)

// markdownClient wraps a Client and renders its replies as Markdown for the terminal.
type markdownClient struct {
	client providers.Client
}

// Send forwards the prompt to the wrapped client and renders the reply.
//   - ctx: controls cancellation of the request.
//   - message: the prompt to send to the model.
func (markdown markdownClient) Send(ctx context.Context, message string) (string, error) {
	answer, err := markdown.client.Send(ctx, message)
	if err != nil {
		return "", err
	}
	return glamour.Render(answer, markdownStyle)
}

// Header returns the wrapped client's header unchanged.
func (markdown markdownClient) Header() providers.Header {
	return markdown.client.Header()
}
