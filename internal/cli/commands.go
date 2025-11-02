package cli

import (
	"fmt"
	"math/rand/v2"
	"os"
	"pokedexcli/internal/models"
)

// Command represents a CLI command
type Command struct {
	Name        string
	Description string
	Callback    func(*models.ReplConfig, []string) error
}

// CommandsMap holds all available commands
var CommandsMap map[string]Command

func init() {
	CommandsMap = map[string]Command{
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback:    CommandExit,
		},
		"help": {
			Name:        "help",
			Description: "Display a help message",
			Callback:    CommandHelp,
		},
		"map": {
			Name:        "map",
			Description: "Display the map, subsequent calls will display the next 20 locations",
			Callback:    CommandMap,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Display the map, subsequent calls will display the previous 20 locations",
			Callback:    CommandMapb,
		},
		"explore": {
			Name:        "explore <area_name>",
			Description: "list of all the Pokémon in <area_name>",
			Callback:    CommandExplore,
		},
		"catch": {
			Name:        "catch <pokemon_name>",
			Description: "attempt catch on <pokemon_name>",
			Callback:    CommandCatch,
		},
		"inspect": {
			Name:        "inspect <pokemon_name>",
			Description: "See more details on <pokemon_name>",
			Callback:    CommandInspect,
		},
		"pokedex": {
			Name:        "pokedex",
			Description: "list all caught pokemon",
			Callback:    CommandPokedex,
		},
	}
}

func CommandMap(config *models.ReplConfig, args []string) error {
	locationAreasListResponse, err := config.PokeApiClient.GetLocationAreasList(config.Next, args)
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

func CommandMapb(config *models.ReplConfig, args []string) error {
	if config.Previous == "" {
		fmt.Println("You're on the first page")
		config.Next = ""
		return nil
	}

	locationAreasListResponse, err := config.PokeApiClient.GetLocationAreasList(config.Previous, args)
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

func CommandExplore(config *models.ReplConfig, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: explore <area_name>")
	}
	locationAreasDetailsResponse, err := config.PokeApiClient.GetLocationAreasDetail(args[0])
	if err != nil {
		return err
	}
	for _, pokemonEncounter := range locationAreasDetailsResponse.PokemonEncounters {
		fmt.Println(pokemonEncounter.Pokemon.Name)
	}
	return nil
}

func CommandCatch(config *models.ReplConfig, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: catch <pokemon_name>")
	}
	pokemonResponse, err := config.PokeApiClient.GetPokemonInformation(args[0])
	if err != nil {
		return err
	}
	normalized := ((float32(pokemonResponse.BaseExperience) - 36.00) / (608.00 - 36.00))
	roll := rand.Float32()
	fmt.Printf("Throwing a Pokeball at %s...\n", args[0])
	if roll > normalized {
		fmt.Printf("%s was caught!\n", args[0])
		config.Pokedex[pokemonResponse.Name] = models.Pokemon{Name: pokemonResponse.Name}
	} else {
		fmt.Printf("%s escaped!\n", args[0])
	}
	return nil
}

func CommandInspect(config *models.ReplConfig, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: inspect <pokemon_name>")
	}
	_, pokemonExists := config.Pokedex[args[0]]
	if !pokemonExists {
		fmt.Printf("you have not caught that pokemon\n")
		return nil
	}
	pokemonResponse, err := config.PokeApiClient.GetPokemonInformation(args[0])
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

func CommandPokedex(config *models.ReplConfig, args []string) error {
	fmt.Println("Your Pokedex:")
	for _, pokemon := range config.Pokedex {
		fmt.Printf("- %s\n", pokemon.Name)
	}
	return nil
}

func CommandExit(config *models.ReplConfig, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func CommandHelp(config *models.ReplConfig, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")

	// Find longest command name for padding
	maxLen := 0
	for _, cmd := range CommandsMap {
		if len(cmd.Name) > maxLen {
			maxLen = len(cmd.Name)
		}
	}

	// Print aligned
	for _, cmd := range CommandsMap {
		fmt.Printf("  %-*s  %s\n", maxLen, cmd.Name, cmd.Description)
	}

	return nil
}
