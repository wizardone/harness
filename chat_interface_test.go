package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// typed returns the key message for typing text.
//   - text: the characters typed.
func typed(text string) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(text)}
}

// press sends each key to a fresh echoing model and returns the resulting model and last command.
//   - keys: the key presses, in order.
func press(keys ...tea.KeyMsg) (chatModel, tea.Cmd) {
	chat := newChatModel(fakeClient{send: echoReply})
	var cmd tea.Cmd
	for _, key := range keys {
		var model tea.Model
		model, cmd = chat.Update(key)
		chat = model.(chatModel)
	}
	return chat, cmd
}

func TestEnterSubmitsPromptAndClearsInput(t *testing.T) {
	chat, cmd := press(typed("hello"), tea.KeyMsg{Type: tea.KeyEnter})

	if cmd == nil {
		t.Fatal("Enter returned no command, want the prompt to be sent")
	}
	if chat.input.Value() != "" {
		t.Fatalf("input = %q, want it cleared", chat.input.Value())
	}
}

func TestQuitShowsGoodbye(t *testing.T) {
	chat, cmd := press(typed("quit"), tea.KeyMsg{Type: tea.KeyEnter})

	if cmd == nil || cmd() != tea.Quit() {
		t.Fatal("quit did not return the quit command")
	}
	if chat.View() != "Bye...\n" {
		t.Fatalf("View() = %q, want %q", chat.View(), "Bye...\n")
	}
}

func TestCtrlCQuits(t *testing.T) {
	_, cmd := press(tea.KeyMsg{Type: tea.KeyCtrlC})

	if cmd == nil || cmd() != tea.Quit() {
		t.Fatal("Ctrl+C did not return the quit command")
	}
}

func TestColonOpensOptionsOnEmptyPrompt(t *testing.T) {
	chat, _ := press(typed(":"))

	if !chat.menu.open {
		t.Fatal("options closed, want them open")
	}
	if chat.input.Value() != "" {
		t.Fatalf("input = %q, want the colon not typed", chat.input.Value())
	}
}

func TestColonIsTypedAfterOtherText(t *testing.T) {
	chat, _ := press(typed("a"), typed(":"))

	if chat.menu.open {
		t.Fatal("options open, want them closed")
	}
	if chat.input.Value() != "a:" {
		t.Fatalf("input = %q, want %q", chat.input.Value(), "a:")
	}
}

func TestEnterClosesOptionsWithoutSending(t *testing.T) {
	chat, cmd := press(typed(":"), tea.KeyMsg{Type: tea.KeyEnter})

	if chat.menu.open {
		t.Fatal("options open, want them closed")
	}
	if cmd != nil {
		t.Fatal("Enter returned a command, want nothing sent")
	}
}

func TestRenderOptionsHighlightsCursor(t *testing.T) {
	want := "  " + options[0] + "\n" + promptMarker + options[1] + "\n" + "  " + options[2]

	if got := renderOptions(1); got != want {
		t.Fatalf("renderOptions(1) = %q, want %q", got, want)
	}
}

func TestNavigateMovesAndWraps(t *testing.T) {
	menu := optionMenu{open: true}.navigate("up")

	if menu.cursor != len(options)-1 {
		t.Fatalf("cursor = %d after up, want %d", menu.cursor, len(options)-1)
	}
	if menu = menu.navigate("down"); menu.cursor != 0 {
		t.Fatalf("cursor = %d after down, want 0", menu.cursor)
	}
}

func TestNavigateClosesOnEscOrColon(t *testing.T) {
	for _, key := range []string{"esc", ":"} {
		if menu := (optionMenu{open: true}).navigate(key); menu.open {
			t.Fatalf("options open after %q, want them closed", key)
		}
	}
}
