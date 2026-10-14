package main

import (
	"harness/providers"
	"testing"
)

func TestRenderHeader(t *testing.T) {
	t.Setenv("HOME", "/home/dev")

	tests := []struct {
		name   string
		header providers.Header
		dir    string
		want   string
	}{
		{
			"max tokens and home shortened",
			providers.Header{Provider: "Anthropic", Model: "claude-opus-5-5", MaxTokens: 16000},
			"/home/dev/projects/harness",
			"claude-opus-5-5\nAnthropic · max tokens 16000 · ~/projects/harness",
		},
		{
			"no max tokens",
			providers.Header{Provider: "OpenAI", Model: "gpt-5.6-sol"},
			"/home/dev",
			"gpt-5.6-sol\nOpenAI · ~",
		},
		{
			"directory outside home kept whole",
			providers.Header{Provider: "OpenAI", Model: "gpt-5.6-sol"},
			"/home/developer/harness",
			"gpt-5.6-sol\nOpenAI · /home/developer/harness",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := renderHeader(tt.header, tt.dir); got != tt.want {
				t.Fatalf("renderHeader() = %q, want %q", got, tt.want)
			}
		})
	}
}
