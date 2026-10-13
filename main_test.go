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

// Send returns whatever the fake's send function returns.
//   - ctx: passed through to the send function.
//   - prompt: passed through to the send function.
func (f fakeClient) Send(ctx context.Context, prompt string) (string, error) {
	return f.send(ctx, prompt)
}

// Header returns a fixed header.
func (f fakeClient) Header() string { return "fake" }

// echoReply answers with the prompt it was given.
//   - prompt: the text to echo back.
func echoReply(_ context.Context, prompt string) (string, error) {
	return "You entered: " + prompt, nil
}

func TestRunStopsOnQuit(t *testing.T) {
	var output bytes.Buffer
	calls := 0
	counting := func(ctx context.Context, prompt string) (string, error) {
		calls++
		return echoReply(ctx, prompt)
	}

	err := run(strings.NewReader("quit\r"), &output, fakeClient{send: counting})

	if err != nil {
		t.Fatalf("run returned an unexpected error: %v", err)
	}
	if !strings.Contains(output.String(), "Bye...") {
		t.Fatalf("output = %q, want it to contain %q", output.String(), "Bye...")
	}
	if calls != 0 {
		t.Fatalf("reply called %d times, want 0", calls)
	}
}

type failingReader struct{}

// Read always fails.
//   - []byte: unused buffer.
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
