package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"harness/providers"
)

func main() {
	client, err := providers.ChooseProvider()
	if err != nil {
		log.Fatalf("provider error: %v", err)
	}

	printDefaults(client)

	if err := run(os.Stdin, os.Stdout, markdownClient{client}); err != nil {
		log.Printf("input error: %v", err)
	}
}

// run reads prompts from in until EOF or "quit", writing each reply to out.
func run(in io.Reader, out io.Writer, client providers.Client) error {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, "> ")

		if !scanner.Scan() {
			return scanner.Err()
		}

		prompt := scanner.Text()
		if prompt == "quit" {
			fmt.Fprint(out, "Bye...")
			return nil
		}

		answer, err := client.Send(context.Background(), prompt)
		if err != nil {
			fmt.Fprintf(out, "error: %v\n", err)
			continue
		}

		fmt.Fprintln(out, answer)
	}
}
