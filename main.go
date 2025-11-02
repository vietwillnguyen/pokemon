package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/internal/pokeapi"
	"strings"
)

var commandsMap map[string]cliCommand

type cliCommand struct {
	name        string
	description string
	callback    func(*replConfig, []string) error
}

type replConfig struct {
	pokeApiClient *pokeapi.Client
	Next          string `json:"next"`
	Previous      string `json:"previous"`
}

func commandMap(config *replConfig, args []string) error {
	locationAreasListResponse, err := config.pokeApiClient.GetLocationAreasList(config.Next, args)
	if err != nil {
		return err
	}
	for _, response := range locationAreasListResponse.Results {
		fmt.Println(response.Name)
	}
	// Update config with next/previous URLs for pagination
	config.Next = locationAreasListResponse.Next
	config.Previous = locationAreasListResponse.Previous
	return nil
}

func commandMapb(config *replConfig, args []string) error {
	if config.Previous == "" {
		fmt.Println("You're on the first page")
		config.Next = ""
		return nil
	}

	locationAreasListResponse, err := config.pokeApiClient.GetLocationAreasList(config.Previous, args)
	if err != nil {
		return err
	}
	for _, location := range locationAreasListResponse.Results {
		fmt.Println(location.Name)
	}
	// Update config with next/previous URLs for pagination
	config.Next = locationAreasListResponse.Next
	config.Previous = locationAreasListResponse.Previous
	return nil
}

func commandExplore(config *replConfig, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: explore <area_name>")
	}
	locationAreasDetailsResponse, err := config.pokeApiClient.GetLocationAreasDetail(args[0])
	if err != nil {
		return err
	}
	for _, pokemonEncounter := range locationAreasDetailsResponse.PokemonEncounters {
		fmt.Println(pokemonEncounter.Pokemon.Name)
	}
	return nil
}

func commandExit(config *replConfig, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *replConfig, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, cmd := range commandsMap {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func init() {
	commandsMap = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Display a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Display the map, subsequent calls will display the next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the map, subsequent calls will display the previous 20 locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore <area_name>",
			description: "list of all the Pokémon in <area_name>",
			callback:    commandExplore,
		},
	}
}

func cleanInput(input string) []string {
	lowered := strings.ToLower(input)
	words := strings.Fields(lowered)
	return words
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	config := &replConfig{
		pokeApiClient: pokeapi.NewClient(),
	}
	for {
		fmt.Print("Pokedex > ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		words := cleanInput(input)
		if len(words) == 0 {
			continue
		}
		command := words[0]
		args := words[1:]
		cmd, exists := commandsMap[command]
		if !exists {
			fmt.Printf("Unknown command\n")
			continue
		}
		err = cmd.callback(config, args)
		if err != nil {
			fmt.Println("Error executing command:", err)
		}
	}
}
