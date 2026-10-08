package main

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
)

// fakeClient satisfies Client by delegating Send to a function.
type fakeClient struct {
	send func(ctx context.Context, prompt string) (string, error)
}

func (f fakeClient) Send(ctx context.Context, prompt string) (string, error) {
	return f.send(ctx, prompt)
}

func (f fakeClient) Header() string { return "fake" }

func echoReply(_ context.Context, prompt string) (string, error) {
	return "You entered: " + prompt, nil
}

func TestRunReadsInputUntilEOF(t *testing.T) {
	var output bytes.Buffer

	err := run(strings.NewReader("hello\nworld\n"), &output, fakeClient{send: echoReply})

	if err != nil {
		t.Fatalf("run returned an unexpected error: %v", err)
	}

	want := "> You entered: hello\n> You entered: world\n> "
	if output.String() != want {
		t.Fatalf("output = %q, want %q", output.String(), want)
	}
}

type failingReader struct{}

func (failingReader) Read([]byte) (int, error) {
	return 0, errors.New("read failed")
}

func TestRunReturnsInputErrors(t *testing.T) {
	var output bytes.Buffer

	err := run(failingReader{}, &output, fakeClient{send: echoReply})

	if err == nil {
		t.Fatal("run returned nil, want an input error")
	}
}

func TestRunPrintsReplyErrorsAndContinues(t *testing.T) {
	var output bytes.Buffer
	failing := func(context.Context, string) (string, error) {
		return "", errors.New("api down")
	}

	err := run(strings.NewReader("hello\n"), &output, fakeClient{send: failing})

	if err != nil {
		t.Fatalf("run returned an unexpected error: %v", err)
	}

	want := "> error: api down\n> "
	if output.String() != want {
		t.Fatalf("output = %q, want %q", output.String(), want)
	}
}

func TestRunStopsOnQuit(t *testing.T) {
	var output bytes.Buffer
	calls := 0
	counting := func(ctx context.Context, prompt string) (string, error) {
		calls++
		return echoReply(ctx, prompt)
	}

	err := run(strings.NewReader("hello\nquit\nworld\n"), &output, fakeClient{send: counting})

	if err != nil {
		t.Fatalf("run returned an unexpected error: %v", err)
	}

	want := "> You entered: hello\n> Bye...\n"
	if output.String() != want {
		t.Fatalf("output = %q, want %q", output.String(), want)
	}
	if calls != 1 {
		t.Fatalf("reply called %d times, want 1", calls)
	}
}
