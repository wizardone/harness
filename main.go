package main

import (
	"bufio"
	"context"
	"fmt"
	"harness/providers"
	"io"
	"log"
	"os"
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

// run reads prompts from in until EOF or "quit", writing each reply to out.
func run(in io.Reader, out io.Writer, client providers.Client) error {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, promptStyle.Render("> "))

		if !scanner.Scan() {
			return scanner.Err()
		}

		prompt := scanner.Text()
		if prompt == "quit" {
			fmt.Fprintln(out, hintStyle.Render("Bye..."))
			return nil
		}

		answer, err := client.Send(context.Background(), prompt)
		if err != nil {
			fmt.Fprintln(out, errorStyle.Render(fmt.Sprintf("error: %v", err)))
			continue
		}

		fmt.Fprintln(out, answer)
	}
}
