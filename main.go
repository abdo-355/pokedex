package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	cmds := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
	}

	for {
		fmt.Print("Pokedex > ")

		r := scanner.Scan()
		if r {
			t := scanner.Text()
			ci := cleanInput(t)

			if cmd := cmds[strings.ToLower(ci[0])]; cmd.callback != nil {
				err := cmd.callback()
				if err != nil {
					fmt.Println("Error:", err)
				}
			} else {
				fmt.Println("Unknown command")
			}

		}
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")

	return nil
}

func cleanInput(text string) []string {
	cleaned := strings.TrimSpace(text)

	split := strings.Split(cleaned, " ")

	var removedZero []string

	for _, v := range split {
		c := strings.TrimSpace(v)
		if c != "" {
			removedZero = append(removedZero, c)
		}
	}

	return removedZero
}
