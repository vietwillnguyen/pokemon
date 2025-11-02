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
