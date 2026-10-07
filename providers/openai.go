package providers
import (
	"fmt"
	"context"
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

func (c OpenAIClient) Describe() string {
	return fmt.Sprintf("Provider: OpenAI | Model: %s", c.model)
}

func (c OpenAIClient) Send(ctx context.Context, message string) (string, error) {
	resp, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openAIModel,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(message),
		},
	})
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("openai returned no choices")
	}
	return resp.Choices[0].Message.Content, nil
}
