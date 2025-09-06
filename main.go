package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/abdo-355/pokedex/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, *pokecache.Cache, string) error
}

type config struct {
	Next, Previous string
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	cache := pokecache.NewCache(7 * time.Second)

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
		"explore": {
			name:        "explore",
			description: "explore the pokemon available at a certain location",
			callback:    commandExplore,
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
				var err error
				if len(ci) < 2 {
					err = cmd.callback(&cfg, cache, "")
				} else {
					err = cmd.callback(&cfg, cache, strings.ToLower(ci[1]))
				}
				if err != nil {
					fmt.Println("Error:", err)
				}
			} else {
				fmt.Println("Unknown command")
			}

		}
	}
}

func commandExit(c *config, cache *pokecache.Cache, _ string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func commandHelp(c *config, cache *pokecache.Cache, _ string) error {
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

func commandMap(c *config, cache *pokecache.Cache, _ string) error {
	if c.Next == "" {
		c.Next = "https://pokeapi.co/api/v2/location-area/"
	}

	body, err := pokeRequest(c.Next, cache)
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

func commandMapb(c *config, cache *pokecache.Cache, _ string) error {
	if c.Previous == "" {
		c.Previous = "https://pokeapi.co/api/v2/location-area/"
	}

	body, err := pokeRequest(c.Previous, cache)
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

func pokeRequest(url string, cache *pokecache.Cache) ([]byte, error) {
	cd, cached := cache.Get(url)
	if cached {
		return cd, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	cache.Add(url, body)

	return body, nil
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

type EncounterResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	GameIndex int    `json:"game_index"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
		VersionDetails []struct {
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
			} `json:"version"`
			EncounterDetails []struct {
				Chance   int `json:"chance"`
				MinLevel int `json:"min_level"`
				MaxLevel int `json:"max_level"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func commandExplore(c *config, cache *pokecache.Cache, p string) error {
	body, err := pokeRequest("https://pokeapi.co/api/v2/location-area/"+p, cache)
	if err != nil {
		return err
	}

	r := EncounterResponse{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", p)
	fmt.Println("Found Pokemon:")

	for _, v := range r.PokemonEncounters {
		fmt.Printf("- %s\n", v.Pokemon.Name)
	}

	return nil
}
