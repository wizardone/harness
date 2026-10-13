package main

import (
	"errors"
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
	if got := replyText(replyMsg{answer: "hi"}); got != "hi" {
		t.Fatalf("replyText = %q, want %q", got, "hi")
	}
	if got := replyText(replyMsg{err: errors.New("api down")}); got != "error: api down" {
		t.Fatalf("replyText = %q, want %q", got, "error: api down")
	}
}
