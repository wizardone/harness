package main

import (
	"harness/providers"
	"io"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	client, err := providers.ChooseProvider()
	if err != nil {
		log.Fatalf("provider error: %v", err)
	}

	printDefaults(markdownClient{client: client})

	if err := run(os.Stdin, os.Stdout, markdownClient{client: client}); err != nil {
		log.Printf("input error: %v", err)
	}
}

// run reads prompts from in until Ctrl+C, Ctrl+D or "quit", writing each reply to out.
//   - in: where key presses are read from.
//   - out: where the prompt, options and replies are drawn.
//   - client: sends each submitted prompt to the model.
func run(in io.Reader, out io.Writer, client providers.Client) error {
	_, err := tea.NewProgram(newChatModel(client), tea.WithInput(in), tea.WithOutput(out)).Run()
	return err
}
