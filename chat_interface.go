package main

import (
	"harness/providers"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// chatModel holds the prompt being typed and the state of the options list.
type chatModel struct {
	client   providers.Client
	input    textinput.Model
	menu     optionMenu
	quitting bool
	thinking bool
	width    int
	history []string
}

// newChatModel returns a model with a focused prompt.
//   - client: sends each submitted prompt to the model.
func newChatModel(client providers.Client) chatModel {
	input := textinput.New()
	input.Prompt = promptStyle.Render(promptMarker)
	input.Focus()
	return chatModel{client: client, input: input}
}

// Init starts the cursor blinking.
func (chat chatModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update reacts to key presses and replies.
//   - msg: the key press, reply or other message to handle.
func (chat chatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case replyMsg:
		chat.thinking = false
		if msg.err == nil {
			chat.history = append(chat.history, msg.prompt, msg.reply.Text)
		}
		return chat, tea.Println(renderReply(msg))
	case tea.KeyMsg:
		return chat.handleKey(msg)
	case tea.WindowSizeMsg:
		chat.width = msg.Width
	}
	return chat.updateInput(msg)
}

// handleKey quits, drives the options list, submits the prompt or types into it.
//   - key: the key pressed.
func (chat chatModel) handleKey(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch name := key.String(); {
	case name == "ctrl+c" || name == "ctrl+d":
		return chat, tea.Quit
	case chat.menu.open && name == "enter" && options[chat.menu.cursor] == "Quit":
      	chat.quitting = true
      	return chat, tea.Quit
	case chat.menu.open || name == ":" && chat.input.Value() == "":
		chat.menu = chat.menu.navigate(name)
		return chat, nil
	case name == "enter":
		return chat.submit()
	}
	return chat.updateInput(key)
}

// submit quits on "quit", otherwise echoes the prompt and sends it to the model.
func (chat chatModel) submit() (tea.Model, tea.Cmd) {
	prompt := chat.input.Value()
	chat.input.Reset()
	if prompt == "quit" {
		chat.quitting = true
		return chat, tea.Quit
	}
	chat.thinking = true
	turn := len(chat.history)/2 + 1
	return chat, tea.Sequence(tea.Println(renderSentPrompt(prompt, turn, chat.width)), send(chat.client, chat.history, prompt))
}

// updateInput passes the message to the prompt, for typing and cursor blinking.
//   - msg: the message for the prompt.
func (chat chatModel) updateInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	chat.input, cmd = chat.input.Update(msg)
	return chat, cmd
}

// View shows the prompt, with the options underneath when they are open.
func (chat chatModel) View() string {
	if chat.quitting {
		return hintStyle.Render("Bye...") + "\n"
	}
	if chat.thinking {
		return renderPromptArea(hintStyle.Render("Thinking..."), chat.width)
	}
	prompt := renderPromptArea(chat.input.View(), chat.width)
	if chat.menu.open {
		return prompt + "\n" + renderOptions(chat.menu.cursor)
	}
	return prompt
}

var options = []string{"Provider", "Model", "History", "Quit"}

// optionMenu holds whether the options list is shown and which option is highlighted.
type optionMenu struct {
	open   bool
	cursor int
}

// navigate moves the highlight, toggles the list on ":", or closes it.
//   - key: the name of the key pressed, such as "up" or "enter".
func (menu optionMenu) navigate(key string) optionMenu {
	switch key {
	case "up":
		menu.cursor = (menu.cursor - 1 + len(options)) % len(options)
	case "down":
		menu.cursor = (menu.cursor + 1) % len(options)
	case ":":
		menu.open = !menu.open
	case "enter", "esc":
		// TODO: on "enter", act on options[menu.cursor].
		menu.open = false
	}
	return menu
}
