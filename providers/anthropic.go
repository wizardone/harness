package providers

import (
	"context"
	"fmt"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
)

const (
	anthropicModel = "claude-opus-5-5"
	anthropicMaxTokens = 16000
)

type AnthropicClient struct {
	client anthropic.Client
	model anthropic.Model
}

func newAnthropicClient() AnthropicClient {
	return AnthropicClient{client: anthropic.NewClient(), model: anthropicModel}
}

func (c AnthropicClient) Send(ctx context.Context, message string) (Reply, error) {
	resp, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropicModel,
		MaxTokens: anthropicMaxTokens,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(message)),
		},
	})
	if err != nil {
		return Reply{}, err
	}
	if resp.StopReason == anthropic.StopReasonRefusal {
		return Reply{}, fmt.Errorf("anthropic refused the request: %s", resp.StopDetails.Category)
	}

	var reply strings.Builder
	for _, block := range resp.Content {
		if text, ok := block.AsAny().(anthropic.TextBlock); ok {
			reply.WriteString(text.Text)
		}
	}
	return Reply{Text: reply.String(), InputTokens: resp.Usage.InputTokens, OutputTokens: resp.Usage.OutputTokens}, nil
}

func (c AnthropicClient) Header() Header {
	return Header{Provider: "Anthropic", Model: c.model, MaxTokens: anthropicMaxTokens}
}
