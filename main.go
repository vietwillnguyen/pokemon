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
	callback    func(*replConfig) error
}

type replConfig struct {
	pokeApiClient *pokeapi.Client
	Next          string `json:"next"`
	Previous      string `json:"previous"`
}

func commandMap(config *replConfig) error {
	locationAreasResponse, err := config.pokeApiClient.GetLocationAreas(config.Next)
	if err != nil {
		return err
	}
	for _, location := range locationAreasResponse.Results {
		fmt.Println(location.Name)
	}
	// Update config with next/previous URLs for pagination
	config.Next = locationAreasResponse.Next
	config.Previous = locationAreasResponse.Previous
	return nil
}

/*
commandMapb is a command that displays the previous 20 locations on the map.

// Scenarios:
1. User is on the first page and tries to go back
*/
func commandMapb(config *replConfig) error {

	if config.Previous == "" {
		fmt.Println("You're on the first page")
		config.Next = ""
		return nil
	}

	locationAreasResponse, err := config.pokeApiClient.GetLocationAreas(config.Previous)
	if err != nil {
		return err
	}
	for _, location := range locationAreasResponse.Results {
		fmt.Println(location.Name)
	}
	// Update config with next/previous URLs for pagination
	config.Next = locationAreasResponse.Next
	config.Previous = locationAreasResponse.Previous
	return nil
}

func commandExit(config *replConfig) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *replConfig) error {
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
		// args := words[1:]
		cmd, exists := commandsMap[command]
		if !exists {
			fmt.Printf("Unknown command\n")
			continue
		}
		err = cmd.callback(config)
		if err != nil {
			fmt.Println("Error executing command:", err)
		}
	}
}
