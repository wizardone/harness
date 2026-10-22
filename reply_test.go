package main

import (
	"errors"
	"harness/providers"
	"strings"
	"testing"
)

func TestSendReportsReply(t *testing.T) {
	msg := send(fakeClient{send: echoReply}, []string{"hi", "hello there"}, "hello")()

	want := replyMsg{reply: providers.Reply{Text: "You entered: hello"}, prompt: "hello"}
	if msg != want {
		t.Fatalf("msg = %+v, want %+v", msg, want)
	}
}

func TestReplyTextShowsAnswerOrError(t *testing.T) {
	if got := replyText(replyMsg{reply: providers.Reply{Text: "hi"}}); !strings.Contains(got, "hi") {
		t.Fatalf("replyText = %q, want it to contain %q", got, "hi")
	}
	if got := replyText(replyMsg{err: errors.New("api down")}); got != "error: api down" {
		t.Fatalf("replyText = %q, want %q", got, "error: api down")
	}
}

func TestReplyTextRendersMarkdown(t *testing.T) {
	got := replyText(replyMsg{reply: providers.Reply{Text: "- one\n- two"}})

	if !strings.Contains(got, "• one") || !strings.Contains(got, "• two") {
		t.Fatalf("replyText = %q, want the list rendered with bullets", got)
	}
}

func TestReplyTextShowsTokenUsage(t *testing.T) {
	got := replyText(replyMsg{reply: providers.Reply{Text: "hi", InputTokens: 12, OutputTokens: 34}})

	if !strings.HasSuffix(got, "\n12 in · 34 out") {
		t.Fatalf("replyText = %q, want it to end with the token usage", got)
	}
}

func TestRenderMarkdownStylesInsteadOfSymbols(t *testing.T) {
	got := renderMarkdown("# One\n\n## Two\n\n### Three\n\n#### Four\n\n##### Five\n\n###### Six\n\n" +
		"**bold** and *italic* and `code` and ~~gone~~")

	for _, symbol := range []string{"#", "**", "*italic*", "`", "~~"} {
		if strings.Contains(got, symbol) {
			t.Fatalf("renderMarkdown = %q, want no %q", got, symbol)
		}
	}
	for _, styled := range []string{"\x1b[1mSix", "\x1b[1mbold", "\x1b[3mitalic", "\x1b[1mcode", "\x1b[9mgone"} {
		if !strings.Contains(got, styled) {
			t.Fatalf("renderMarkdown = %q, want it to contain %q", got, styled)
		}
	}
}
