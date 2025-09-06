package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	Next, Previous string
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
		"map": {
			name:        "map",
			description: "Displays the names of 20 location areas in the Pokemon world.",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the names of the previous 20 location areas in the Pokemon world.",
			callback:    commandMapb,
		},
	}

	cfg := config{
		"", "",
	}

	for {
		fmt.Print("Pokedex > ")

		r := scanner.Scan()
		if r {
			t := scanner.Text()
			ci := cleanInput(t)

			if cmd := cmds[strings.ToLower(ci[0])]; cmd.callback != nil {
				err := cmd.callback(&cfg)
				if err != nil {
					fmt.Println("Error:", err)
				}
			} else {
				fmt.Println("Unknown command")
			}

		}
	}
}

func commandExit(c *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func commandHelp(c *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println("help: Displays a help message")
	fmt.Println("map: displays the names of 20 location areas in the Pokemon world.")
	fmt.Println("mapb: displays the names of the previous 20 location areas in the Pokemon world.")
	fmt.Println("exit: Exit the Pokedex")

	return nil
}

type apiResponse struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func commandMap(c *config) error {
	if c.Next == "" {
		c.Next = "https://pokeapi.co/api/v2/location-area/"
	}

	res, err := http.Get(c.Next)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	r := apiResponse{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return err
	}

	c.Next = r.Next
	if r.Previous == "" {
		c.Previous = "https://pokeapi.co/api/v2/location-area/"
	} else {
		c.Previous = r.Previous
	}

	for _, v := range r.Results {
		fmt.Println(v.Name)
	}

	return nil
}

func commandMapb(c *config) error {
	if c.Previous == "" {
		c.Previous = "https://pokeapi.co/api/v2/location-area/"
	}

	res, err := http.Get(c.Previous)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	r := apiResponse{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return err
	}

	c.Previous = r.Previous
	c.Next = r.Next

	for _, v := range r.Results {
		fmt.Println(v.Name)
	}

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
