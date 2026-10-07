package main

import (
	"context"
	"github.com/charmbracelet/glamour"
	"harness/providers"
)

// markdownClient wraps a Client and renders its replies as Markdown for the terminal.
type markdownClient struct {
	client providers.Client
}

// Send forwards the prompt to the wrapped client and renders the reply.
//   - ctx: controls cancellation of the request.
//   - message: the prompt to send to the model.
func (c markdownClient) Send(ctx context.Context, message string) (string, error) {
	answer, err := c.client.Send(ctx, message)
	if err != nil {
		return "", err
	}
	return glamour.Render(answer, "auto")
}

// Describe returns the wrapped client's description unchanged.
func (c markdownClient) Describe() string {
	return c.client.Describe()
}
