package main

import (
	"context"
	"fmt"
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
	return glamour.Render(answer, "auto")
}

// Describe returns the wrapped client's description rendered as Markdown, or "Ups" if rendering fails.
func (markdown markdownClient) Describe() string {
	answer, err := glamour.Render(markdown.client.Describe(), "auto")
	if err != nil {
		return "Ups"
	}
	return answer
}

// printDefaults prints the provider details in a header box, with a hint on how to exit.
//   - markdown: the wrapped client whose rendered description is shown.
func printDefaults(markdown markdownClient){
	fmt.Println(headerStyle.Render(markdown.Describe()))
	fmt.Println(hintStyle.Render("Type quit to exit."))
}
