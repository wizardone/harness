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

// Send sends the conversation to Claude and returns the joined text of the reply with its token usage.
//   - ctx: cancels or times out the request.
//   - history: user prompts at even indexes and model replies at odd ones, ending with the new prompt.
func (c AnthropicClient) Send(ctx context.Context, history []string) (Reply, error) {
	messages := make([]anthropic.MessageParam, len(history))
	for i, text := range history {
		if i%2 == 0 {
			messages[i] = anthropic.NewUserMessage(anthropic.NewTextBlock(text))
		} else {
			messages[i] = anthropic.NewAssistantMessage(anthropic.NewTextBlock(text))
		}
	}
	resp, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropicModel,
		MaxTokens: anthropicMaxTokens,
		Messages:  messages,
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
