package providers

import (
	"context"
	"fmt"
	"os"
)

type Client interface {
	Send(ctx context.Context, message string) (string, error)
	Header() string
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