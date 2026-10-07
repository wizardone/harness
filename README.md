# Go Harness

A minimal interactive Go program that sends each line you type to a large language model (Anthropic or OpenAI) and prints the reply.

## Requirements

- Go 1.23 or newer
- An API key for Anthropic or OpenAI

## Configuration

The provider is chosen from environment variables at start-up:

| Variable            | Provider  | Model             |
| ------------------- | --------- | ----------------- |
| `ANTHROPIC_API_KEY` | Anthropic | `claude-opus-5-5` |
| `OPENAI_API_KEY`    | OpenAI    | `gpt-5.6-sol`     |

If both are set, Anthropic is used. If neither is set, the program exits with `no provider configured`.

The models are set by the `anthropicModel` and `openAIModel` constants in `provider.go`.

## Start

From this directory, run:

```sh
export ANTHROPIC_API_KEY=your-key   # or OPENAI_API_KEY
go run .
```

The program prints the provider and model, then shows a `> ` prompt. Type a line and press Enter to send it to the model. The reply is printed, followed by a new prompt. If a request fails, the error is printed and the program waits for the next line.

To stop the program, type `quit` and press Enter. You can also press `Ctrl+D` on macOS/Linux, or `Ctrl+Z` then Enter on Windows.

## Test

```sh
go test ./...
```

The tests use local fake HTTP servers in place of the Anthropic and OpenAI APIs, so they need no API keys and make no real API calls.
