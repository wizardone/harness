package main

import (
	"errors"
	"strings"
	"testing"
)

func TestSendReportsReply(t *testing.T) {
	msg := send(fakeClient{send: echoReply}, "hello")()

	want := replyMsg{answer: "You entered: hello"}
	if msg != want {
		t.Fatalf("msg = %+v, want %+v", msg, want)
	}
}

func TestReplyTextShowsAnswerOrError(t *testing.T) {
	if got := replyText(replyMsg{answer: "hi"}); strings.TrimSpace(got) != "hi" {
		t.Fatalf("replyText = %q, want %q once trimmed", got, "hi")
	}
	if got := replyText(replyMsg{err: errors.New("api down")}); got != "error: api down" {
		t.Fatalf("replyText = %q, want %q", got, "error: api down")
	}
}

func TestReplyTextRendersMarkdown(t *testing.T) {
	got := replyText(replyMsg{answer: "- one\n- two"})

	if !strings.Contains(got, "• one") || !strings.Contains(got, "• two") {
		t.Fatalf("replyText = %q, want the list rendered with bullets", got)
	}
}
