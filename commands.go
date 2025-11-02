package main

import (
	"fmt"
	"math/rand/v2"
	"os"
)

var commandsMap map[string]cliCommand

type cliCommand struct {
	name        string
	description string
	callback    func(*replConfig, []string) error
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

func commandCatch(config *replConfig, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: catch <pokemon_name>")
	}
	pokemonResponse, err := config.pokeApiClient.GetPokemonInformation(args[0])
	if err != nil {
		return err
	}
	normalized := ((float32(pokemonResponse.BaseExperience) - 36.00) / (608.00 - 36.00))
	roll := rand.Float32()
	fmt.Printf("Throwing a Pokeball at %s...\n", args[0])
	// fmt.Printf("baseExp: %d\n", pokemonResponse.BaseExperience)
	// fmt.Printf("your roll: %.2f, normalized: %.2f\n", roll, normalized)
	if roll > normalized {
		fmt.Printf("%s was caught!\n", args[0])
		config.pokedex[pokemonResponse.Name] = Pokemon{pokemonResponse.Name}
	} else {
		fmt.Printf("%s escaped!\n", args[0])
	}
	return nil
}

func commandInspect(config *replConfig, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: inspect <pokemon_name>")
	}
	_, pokemonExists := config.pokedex[args[0]]
	if !pokemonExists {
		fmt.Printf("you have not caught that pokemon\n")
		return nil
	}
	pokemonResponse, err := config.pokeApiClient.GetPokemonInformation(args[0])
	if err != nil {
		return fmt.Errorf("failed to get info for pokemon %s", args[0])
	}

	// print in required format
	fmt.Printf("Name: %s\n", pokemonResponse.Name)
	fmt.Printf("Height: %d\n", pokemonResponse.Height)
	fmt.Printf("Weight: %d\n", pokemonResponse.Weight)

	// stats
	fmt.Println("Stats:")
	for _, s := range pokemonResponse.Stats {
		// s.Stat.Name is like "special-attack"
		fmt.Printf("  -%s: %d\n", s.Stat.Name, s.BaseStat)
	}

	// types
	fmt.Println("Types:")
	for _, t := range pokemonResponse.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}

	return nil
}

func commandPokedex(config *replConfig, args []string) error {
	fmt.Println("Your Pokedex:")
	for _, pokemon := range config.pokedex {
		fmt.Printf("- %s\n", pokemon.Name)
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

	// Find longest command name for padding
	maxLen := 0
	for _, cmd := range commandsMap {
		if len(cmd.name) > maxLen {
			maxLen = len(cmd.name)
		}
	}

	// Print aligned
	for _, cmd := range commandsMap {
		fmt.Printf("  %-*s  %s\n", maxLen, cmd.name, cmd.description)
	}

	return nil
}
