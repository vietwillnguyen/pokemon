package cli

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/internal/models"
	"pokedexcli/internal/pokeapi"
	"strings"
)

// CleanInput normalizes user input by lowercasing and splitting into words
func CleanInput(input string) []string {
	lowered := strings.ToLower(input)
	words := strings.Fields(lowered)
	return words
}

// StartREPL initializes and starts the REPL loop
func StartREPL() {
	reader := bufio.NewReader(os.Stdin)
	config := &models.ReplConfig{
		Pokedex:       map[string]models.Pokemon{},
		PokeApiClient: pokeapi.NewClient(),
	}

	for {
		fmt.Print("Pokedex > ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		words := CleanInput(input)
		if len(words) == 0 {
			continue
		}

		command := words[0]
		args := words[1:]

		cmd, exists := CommandsMap[command]
		if !exists {
			fmt.Printf("Unknown command\n")
			continue
		}

		err = cmd.Callback(config, args)
		if err != nil {
			fmt.Println("Error executing command:", err)
		}
	}
}

