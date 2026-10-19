package providers

import (
	"context"
	"fmt"
	"os"
)

type Client interface {
	// Send sends the conversation and returns the model's next reply.
	//   - ctx: cancels or times out the request.
	//   - history: the turns so far, user prompts at even indexes and model replies at odd ones,
	//     ending with the new prompt.
	Send(ctx context.Context, history []string) (Reply, error)
	Header() Header
}

// Header describes the provider for the start-up header.
type Header struct {
	Provider  string
	Model     string
	MaxTokens int // 0 when the provider sets no limit
}

// Reply describes the reply format an LLM returns
type Reply struct {
	Text string
	InputTokens int64
	OutputTokens int64
}

func isEnvSet(key string) bool {
 value, ok := os.LookupEnv(key)
 return ok && value != ""
}

func ChooseProvider() (Client, error) {
	if isEnvSet("ANTHROPIC_API_KEY"){
		return newAnthropicClient(), nil
	}

	if isEnvSet("OPENAI_API_KEY"){
		return newOpenAICLient(), nil
	}
	return nil, fmt.Errorf("no provider configured")
}