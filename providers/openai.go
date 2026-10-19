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

// Send sends the conversation to OpenAI and returns the first choice of the reply.
//   - ctx: cancels or times out the request.
//   - history: user prompts at even indexes and model replies at odd ones, ending with the new prompt.
func (c OpenAIClient) Send(ctx context.Context, history []string) (Reply, error) {
	messages := make([]openai.ChatCompletionMessageParamUnion, len(history))
	for i, text := range history {
		if i%2 == 0 {
			messages[i] = openai.UserMessage(text)
		} else {
			messages[i] = openai.AssistantMessage(text)
		}
	}
	resp, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    openAIModel,
		Messages: messages,
	})
	if err != nil {
		return Reply{}, err
	}
	if len(resp.Choices) == 0 {
		return Reply{}, fmt.Errorf("openai returned no choices")
	}
	return Reply{Text: resp.Choices[0].Message.Content}, nil
}
