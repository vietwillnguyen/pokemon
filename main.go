package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/internal/pokeapi"
	"strings"
)

type Pokemon struct {
	Name string
}

type replConfig struct {
	pokedex       map[string]Pokemon
	pokeApiClient *pokeapi.Client
	Next          string `json:"next"`
	Previous      string `json:"previous"`
}

func cleanInput(input string) []string {
	lowered := strings.ToLower(input)
	words := strings.Fields(lowered)
	return words
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
		"catch": {
			name:        "catch <pokemon_name>",
			description: "attempt catch on <pokemon_name>",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect <pokemon_name>",
			description: "See more details on <pokemon_name>",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "list all caught pokemon",
			callback:    commandPokedex,
		},
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	config := &replConfig{
		pokedex:       map[string]Pokemon{},
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
