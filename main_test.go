package main

import (
	"bytes"
	"context"
	"errors"
	"harness/providers"
	"strings"
	"testing"
)

// fakeClient satisfies Client by delegating Send to a function.
type fakeClient struct {
	send func(ctx context.Context, history []string) (providers.Reply, error)
}

// Send returns whatever the fake's send function returns.
//   - ctx: passed through to the send function.
//   - history: passed through to the send function.
func (f fakeClient) Send(ctx context.Context, history []string) (providers.Reply, error) {
	return f.send(ctx, history)
}

// Header returns a fixed header.
func (f fakeClient) Header() providers.Header {
	return providers.Header{Provider: "Fake", Model: "fake"}
}

// echoReply answers with the latest prompt it was given.
//   - history: the conversation, ending with the prompt to echo back.
func echoReply(_ context.Context, history []string) (providers.Reply, error) {
	return providers.Reply{Text: "You entered: " + history[len(history)-1]}, nil
}

func TestRunStopsOnQuit(t *testing.T) {
	var output bytes.Buffer
	calls := 0
	counting := func(ctx context.Context, history []string) (providers.Reply, error) {
		calls++
		return echoReply(ctx, history)
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
