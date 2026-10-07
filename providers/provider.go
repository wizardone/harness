package providers
import (
	"fmt"
	"os"
	"context"
)

type Client interface {
	Send(ctx context.Context, message string) (string, error)
	Describe() string
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