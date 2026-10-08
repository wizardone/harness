package main

import (
	"context"
	"fmt"
	"harness/providers"
	"os"

	"github.com/charmbracelet/glamour"
	"golang.org/x/term"
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
	return glamour.Render(answer, "light")
}

// Header returns the wrapped client's header unchanged.
func (markdown markdownClient) Header() string {
	return markdown.client.Header()
}

// printDefaults prints the provider details in a header box spanning the terminal width, with a hint on how to exit.
//   - markdown: the wrapped client whose header is shown.
func printDefaults(markdown markdownClient){
	header := headerStyle
	if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		header = header.Width(width - header.GetHorizontalBorderSize())
	}
	fmt.Println(header.Render(markdown.Header()))
	fmt.Println(hintStyle.Render("Type quit to exit."))
}
