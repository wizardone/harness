package providers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	anthropicoption "github.com/anthropics/anthropic-sdk-go/option"
	"github.com/openai/openai-go"
	openaioption "github.com/openai/openai-go/option"
)

// fakeAPI serves body with status and records the model named in the request.
func fakeAPI(t *testing.T, status int, body string, gotModel *string) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Model string `json:"model"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("decode request: %v", err)
		}
		*gotModel = req.Model
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)
	return server
}

func testAnthropicClient(url string) AnthropicClient {
	return AnthropicClient{
		client: anthropic.NewClient(
			anthropicoption.WithBaseURL(url),
			anthropicoption.WithAPIKey("test-key"),
			anthropicoption.WithMaxRetries(0),
		),
		model: anthropicModel,
	}
}

func testOpenAIClient(url string) OpenAIClient {
	return OpenAIClient{
		client: openai.NewClient(
			openaioption.WithBaseURL(url),
			openaioption.WithAPIKey("test-key"),
			openaioption.WithMaxRetries(0),
		),
		model: openAIModel,
	}
}

func TestIsEnvSet(t *testing.T) {
	t.Setenv("HARNESS_TEST_SET", "value")
	t.Setenv("HARNESS_TEST_EMPTY", "")

	if !isEnvSet("HARNESS_TEST_SET") {
		t.Fatal("isEnvSet returned false for a set variable")
	}
	if isEnvSet("HARNESS_TEST_EMPTY") {
		t.Fatal("isEnvSet returned true for an empty variable")
	}
	if isEnvSet("HARNESS_TEST_UNSET_VARIABLE") {
		t.Fatal("isEnvSet returned true for an unset variable")
	}
}

func TestChooseProviderPrefersAnthropic(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "anthropic-key")
	t.Setenv("OPENAI_API_KEY", "openai-key")

	client, err := ChooseProvider()

	if err != nil {
		t.Fatalf("chooseProvider returned an unexpected error: %v", err)
	}
	if _, ok := client.(AnthropicClient); !ok {
		t.Fatalf("client = %T, want AnthropicClient", client)
	}
}

func TestChooseProviderFallsBackToOpenAI(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "")
	t.Setenv("OPENAI_API_KEY", "openai-key")

	client, err := ChooseProvider()

	if err != nil {
		t.Fatalf("chooseProvider returned an unexpected error: %v", err)
	}
	if _, ok := client.(OpenAIClient); !ok {
		t.Fatalf("client = %T, want OpenAIClient", client)
	}
}

func TestChooseProviderReturnsErrorWithoutKeys(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "")
	t.Setenv("OPENAI_API_KEY", "")

	client, err := ChooseProvider()

	if err == nil {
		t.Fatalf("chooseProvider returned %T, want an error", client)
	}
}

func TestHeader(t *testing.T) {
	tests := []struct {
		name   string
		client Client
		want   string
	}{
		{"anthropic", newAnthropicClient(), "Provider: Anthropic | Model: claude-opus-5-5 | Max tokens: 16000"},
		{"openai", newOpenAICLient(), "Provider: OpenAI | Model: gpt-5.6-sol"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.client.Header(); got != tt.want {
				t.Fatalf("Header() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAnthropicSendJoinsTextBlocks(t *testing.T) {
	var gotModel string
	server := fakeAPI(t, http.StatusOK, `{
		"id": "msg_1", "type": "message", "role": "assistant", "model": "claude-opus-5-5",
		"content": [
			{"type": "thinking", "thinking": "", "signature": "sig"},
			{"type": "text", "text": "Hello"},
			{"type": "text", "text": " world"}
		],
		"stop_reason": "end_turn",
		"usage": {"input_tokens": 1, "output_tokens": 2}
	}`, &gotModel)

	answer, err := testAnthropicClient(server.URL).Send(context.Background(), "hi")

	if err != nil {
		t.Fatalf("Send returned an unexpected error: %v", err)
	}
	if answer != "Hello world" {
		t.Fatalf("answer = %q, want %q", answer, "Hello world")
	}
	if gotModel != anthropicModel {
		t.Fatalf("model = %q, want %q", gotModel, anthropicModel)
	}
}

func TestAnthropicSendReturnsRefusalError(t *testing.T) {
	var gotModel string
	server := fakeAPI(t, http.StatusOK, `{
		"id": "msg_1", "type": "message", "role": "assistant", "model": "claude-opus-5-5",
		"content": [],
		"stop_reason": "refusal",
		"stop_details": {"type": "refusal", "category": "cyber", "explanation": ""},
		"usage": {"input_tokens": 1, "output_tokens": 0}
	}`, &gotModel)

	_, err := testAnthropicClient(server.URL).Send(context.Background(), "hi")

	if err == nil || !strings.Contains(err.Error(), "cyber") {
		t.Fatalf("err = %v, want a refusal error naming the category", err)
	}
}

func TestAnthropicSendReturnsAPIErrors(t *testing.T) {
	var gotModel string
	server := fakeAPI(t, http.StatusBadRequest, `{
		"type": "error",
		"error": {"type": "invalid_request_error", "message": "bad request"}
	}`, &gotModel)

	_, err := testAnthropicClient(server.URL).Send(context.Background(), "hi")

	if err == nil {
		t.Fatal("Send returned nil, want an API error")
	}
}

func TestOpenAISendReturnsFirstChoice(t *testing.T) {
	var gotModel string
	server := fakeAPI(t, http.StatusOK, `{
		"id": "chatcmpl-1", "object": "chat.completion", "created": 0, "model": "gpt-5.6-sol",
		"choices": [
			{"index": 0, "finish_reason": "stop", "message": {"role": "assistant", "content": "Hello"}}
		]
	}`, &gotModel)

	answer, err := testOpenAIClient(server.URL).Send(context.Background(), "hi")

	if err != nil {
		t.Fatalf("Send returned an unexpected error: %v", err)
	}
	if answer != "Hello" {
		t.Fatalf("answer = %q, want %q", answer, "Hello")
	}
	if gotModel != openAIModel {
		t.Fatalf("model = %q, want %q", gotModel, openAIModel)
	}
}

func TestOpenAISendReturnsErrorWithoutChoices(t *testing.T) {
	var gotModel string
	server := fakeAPI(t, http.StatusOK, `{
		"id": "chatcmpl-1", "object": "chat.completion", "created": 0, "model": "gpt-5.6-sol",
		"choices": []
	}`, &gotModel)

	_, err := testOpenAIClient(server.URL).Send(context.Background(), "hi")

	if err == nil {
		t.Fatal("Send returned nil, want an error for empty choices")
	}
}

func TestOpenAISendReturnsAPIErrors(t *testing.T) {
	var gotModel string
	server := fakeAPI(t, http.StatusBadRequest, `{
		"error": {"message": "bad request", "type": "invalid_request_error"}
	}`, &gotModel)

	_, err := testOpenAIClient(server.URL).Send(context.Background(), "hi")

	if err == nil {
		t.Fatal("Send returned nil, want an API error")
	}
}
