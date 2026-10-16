package providers

import (
	"context"
	"fmt"
	"os"
)

type Client interface {
	Send(ctx context.Context, message string) (Reply, error)
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