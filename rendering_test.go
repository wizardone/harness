package main

import (
	"errors"
	"strings"
	"testing"
)

func TestRenderOptionsHighlightsCursor(t *testing.T) {
	want := "  " + options[0] + "\n" + promptMarker + options[1] + "\n" + "  " + options[2] + "\n" + "  " + options[3]

	if got := renderOptions(1); got != want {
		t.Fatalf("renderOptions(1) = %q, want %q", got, want)
	}
}

func TestRenderPromptAreaDrawsLinesAboveAndBelow(t *testing.T) {
	want := "──────\n> hi  \n──────"

	if got := renderPromptArea("> hi", 6); got != want {
		t.Fatalf("renderPromptArea = %q, want %q", got, want)
	}
}

func TestRenderSentPromptNumbersTheTurnUnderARule(t *testing.T) {
	want := renderTurnRule(30) + "\n  03 │ hi"

	if got := renderSentPrompt("hi", 3, 30); got != want {
		t.Fatalf("renderSentPrompt = %q, want %q", got, want)
	}
}

func TestRenderReplyPutsErrorsInTheLedger(t *testing.T) {
	want := "     │ error: api down"

	if got := renderReply(replyMsg{err: errors.New("api down")}); got != want {
		t.Fatalf("renderReply = %q, want %q", got, want)
	}
}

func TestRenderLedgerLabelsOnlyTheFirstLine(t *testing.T) {
	want := "  01 │ one\n     │ two"

	if got := renderLedger("01", "one\ntwo"); got != want {
		t.Fatalf("renderLedger = %q, want %q", got, want)
	}
}

func TestRenderTurnRuleCrossesTheMarginBar(t *testing.T) {
	tests := []struct {
		name  string
		width int
		want  string
	}{
		{"spans the width", 30, "─────┼" + strings.Repeat("─", 24)},
		{"keeps a minimum length", 0, "─────┼" + strings.Repeat("─", 20)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := renderTurnRule(tt.width); got != tt.want {
				t.Fatalf("renderTurnRule(%d) = %q, want %q", tt.width, got, tt.want)
			}
		})
	}
}
