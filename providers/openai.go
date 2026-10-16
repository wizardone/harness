package providers

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
)

const (
	openAIModel    = "gpt-5.6-sol"
)

type OpenAIClient struct {
	client openai.Client
	model openai.ChatModel
}

func newOpenAICLient() OpenAIClient {
	return OpenAIClient{client: openai.NewClient(), model: openAIModel}
}

func (c OpenAIClient) Header() Header {
	return Header{Provider: "OpenAI", Model: c.model}
}

func (c OpenAIClient) Send(ctx context.Context, message string) (Reply, error) {
	resp, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openAIModel,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(message),
		},
	})
	if err != nil {
		return Reply{}, err
	}
	if len(resp.Choices) == 0 {
		return Reply{}, fmt.Errorf("openai returned no choices")
	}
	return Reply{Text: resp.Choices[0].Message.Content}, nil
}
