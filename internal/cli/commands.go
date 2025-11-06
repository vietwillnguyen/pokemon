package cli

import (
	"fmt"
	"math/rand/v2"
	"os"
	"pokedexcli/internal/models"
	"strings"
	"time"
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
			Description: "Display the next 20 locations",
			Callback:    CommandMap,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Display the previous 20 locations",
			Callback:    CommandMapb,
		},
		"explore": {
			Name:        "explore <area_name>",
			Description: "List all Pokémon in a specific area",
			Callback:    CommandExplore,
		},
		"catch": {
			Name:        "catch <pokemon_name>",
			Description: "Attempt to catch a Pokémon",
			Callback:    CommandCatch,
		},
		"inspect": {
			Name:        "inspect <pokemon_name>",
			Description: "View details of a caught Pokémon",
			Callback:    CommandInspect,
		},
		"pokedex": {
			Name:        "pokedex",
			Description: "List all caught Pokémon",
			Callback:    CommandPokedex,
		},
	}
}

func CommandMap(config *models.ReplConfig, args []string) error {
	locationAreasListResponse, err := config.PokeApiClient.GetLocationAreasList(config.Next, args)
	if err != nil {
		return err
	}

	fmt.Printf("%s═══ Locations ═══%s\n", colorCyan, colorReset)
	for i, response := range locationAreasListResponse.Results {
		fmt.Printf("%s%2d.%s %s\n", colorGray, i+1, colorReset, response.Name)
	}

	// Update config with next/previous URLs for pagination
	config.Next = locationAreasListResponse.Next
	config.Previous = locationAreasListResponse.Previous

	if config.Next != "" {
		fmt.Printf("\n%sType 'map' for more locations%s\n", colorGray, colorReset)
	}
	return nil
}

func CommandMapb(config *models.ReplConfig, args []string) error {
	if config.Previous == "" {
		fmt.Printf("%s⚠ You're on the first page%s\n", colorYellow, colorReset)
		config.Next = ""
		return nil
	}

	locationAreasListResponse, err := config.PokeApiClient.GetLocationAreasList(config.Previous, args)
	if err != nil {
		return err
	}

	fmt.Printf("%s═══ Locations ═══%s\n", colorCyan, colorReset)
	for i, location := range locationAreasListResponse.Results {
		fmt.Printf("%s%2d.%s %s\n", colorGray, i+1, colorReset, location.Name)
	}

	// Update config with next/previous URLs for pagination
	config.Next = locationAreasListResponse.Next
	config.Previous = locationAreasListResponse.Previous

	if config.Previous != "" {
		fmt.Printf("\n%sType 'mapb' for previous locations%s\n", colorGray, colorReset)
	}
	return nil
}

func CommandExplore(config *models.ReplConfig, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: explore <area_name>")
	}

	areaName := args[0]
	fmt.Printf("%sExploring %s...%s\n", colorYellow, areaName, colorReset)

	locationAreasDetailsResponse, err := config.PokeApiClient.GetLocationAreasDetail(areaName)
	if err != nil {
		return err
	}

	if len(locationAreasDetailsResponse.PokemonEncounters) == 0 {
		fmt.Printf("%sNo Pokémon found in this area%s\n", colorGray, colorReset)
		return nil
	}

	fmt.Printf("\n%s═══ Pokémon Found in %s ═══%s\n", colorGreen, areaName, colorReset)
	for i, pokemonEncounter := range locationAreasDetailsResponse.PokemonEncounters {
		fmt.Printf("%s%2d.%s %s\n", colorGray, i+1, colorReset, pokemonEncounter.Pokemon.Name)
	}
	fmt.Printf("\n%sUse 'catch <pokemon_name>' to attempt a catch!%s\n", colorGray, colorReset)

	return nil
}

func CommandCatch(config *models.ReplConfig, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: catch <pokemon_name>")
	}

	pokemonName := args[0]

	// Check if already caught
	if _, exists := config.Pokedex[pokemonName]; exists {
		fmt.Printf("%s✓ You've already caught %s!%s\n", colorYellow, pokemonName, colorReset)
		return nil
	}

	pokemonResponse, err := config.PokeApiClient.GetPokemonInformation(pokemonName)
	if err != nil {
		return err
	}

	// Higher base experience = harder to catch (inverted from before)
	// Normalize between 0 (easiest) and 1 (hardest)
	catchDifficulty := ((float32(pokemonResponse.BaseExperience) - 36.00) / (608.00 - 36.00))
	roll := rand.Float32()

	fmt.Printf("%sThrowing a Pokéball at %s...%s\n", colorYellow, pokemonName, colorReset)

	// Simulate 3 shakes with suspenseful pauses
	shakes := []string{"Wobble...", "Wobble...", "Wobble..."}

	for i, shake := range shakes {
		time.Sleep(800 * time.Millisecond)
		fmt.Print(shake)

		// For dramatic effect, check if Pokemon breaks free after each shake
		// Each shake has increasing chance to succeed
		shakeThreshold := catchDifficulty * float32(3-i) / 3.0

		if roll > shakeThreshold {
			// Pokemon breaks free
			if i < 2 { // Only break free on first two shakes for drama
				time.Sleep(500 * time.Millisecond)
				fmt.Print(" ")
			}
		} else {
			// Will succeed
			time.Sleep(500 * time.Millisecond)
			fmt.Print(" ")
		}
	}
	fmt.Println()

	// Final result check: easier Pokemon = higher success chance
	if roll > catchDifficulty {
		fmt.Printf("%s✓ Gotcha! %s was caught!%s\n", colorGreen, pokemonName, colorReset)
		fmt.Printf("  %sBase Experience: %d%s\n", colorGray, pokemonResponse.BaseExperience, colorReset)
		config.Pokedex[pokemonResponse.Name] = models.Pokemon{Name: pokemonResponse.Name}
	} else {
		fmt.Printf("%s✗ Oh no! %s broke free!%s\n", colorRed, pokemonName, colorReset)
		catchRate := (1.0 - catchDifficulty) * 100
		fmt.Printf("  %sCatch rate: %.1f%% - Try again!%s\n", colorGray, catchRate, colorReset)
	}
	return nil
}

func CommandInspect(config *models.ReplConfig, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: inspect <pokemon_name>")
	}

	pokemonName := args[0]
	_, pokemonExists := config.Pokedex[pokemonName]
	if !pokemonExists {
		fmt.Printf("%s✗ You haven't caught %s yet!%s\n", colorRed, pokemonName, colorReset)
		fmt.Printf("  %sUse 'catch %s' to attempt a catch%s\n", colorGray, pokemonName, colorReset)
		return nil
	}

	pokemonResponse, err := config.PokeApiClient.GetPokemonInformation(pokemonName)
	if err != nil {
		return fmt.Errorf("failed to get info for pokemon %s", pokemonName)
	}

	// Print header
	fmt.Printf("\n%s╔═══════════════════════════════════╗%s\n", colorCyan, colorReset)
	fmt.Printf("%s║  %-31s  ║%s\n", colorCyan, strings.ToUpper(pokemonResponse.Name), colorReset)
	fmt.Printf("%s╚═══════════════════════════════════╝%s\n\n", colorCyan, colorReset)

	// Basic info
	fmt.Printf("%sHeight:%s %d decimetres\n", colorBold, colorReset, pokemonResponse.Height)
	fmt.Printf("%sWeight:%s %d hectograms\n\n", colorBold, colorReset, pokemonResponse.Weight)

	// Types
	fmt.Printf("%sTypes:%s\n", colorBold, colorReset)
	for _, t := range pokemonResponse.Types {
		typeColor := getTypeColor(t.Type.Name)
		fmt.Printf("  • %s%s%s\n", typeColor, t.Type.Name, colorReset)
	}

	// Stats with visual bars
	fmt.Printf("\n%sStats:%s\n", colorBold, colorReset)
	for _, s := range pokemonResponse.Stats {
		statName := strings.ReplaceAll(s.Stat.Name, "-", " ")
		bar := generateStatBar(s.BaseStat)
		fmt.Printf("  %-18s %s%3d%s %s\n", statName+":", colorGray, s.BaseStat, colorReset, bar)
	}

	return nil
}

// generateStatBar creates a visual bar for stats
func generateStatBar(stat int) string {
	maxBarLength := 20
	filledLength := (stat * maxBarLength) / 255
	if filledLength > maxBarLength {
		filledLength = maxBarLength
	}

	bar := strings.Repeat("█", filledLength) + strings.Repeat("░", maxBarLength-filledLength)

	// Color based on stat value
	if stat >= 150 {
		return colorGreen + bar + colorReset
	} else if stat >= 100 {
		return colorYellow + bar + colorReset
	} else {
		return colorRed + bar + colorReset
	}
}

// getTypeColor returns appropriate color for Pokemon type
func getTypeColor(typeName string) string {
	typeColors := map[string]string{
		"fire":     colorRed,
		"water":    colorBlue,
		"grass":    colorGreen,
		"electric": colorYellow,
		"psychic":  colorPurple,
		"ice":      colorCyan,
		"dragon":   colorPurple,
		"dark":     colorGray,
		"fairy":    colorPurple,
	}

	if color, exists := typeColors[typeName]; exists {
		return color
	}
	return colorReset
}

func CommandPokedex(config *models.ReplConfig, args []string) error {
	if len(config.Pokedex) == 0 {
		fmt.Printf("%sYour Pokédex is empty!%s\n", colorYellow, colorReset)
		fmt.Printf("  %sUse 'explore' and 'catch' to start collecting Pokémon%s\n", colorGray, colorReset)
		return nil
	}

	fmt.Printf("%s╔═══════════════════════════════════╗%s\n", colorGreen, colorReset)
	fmt.Printf("%s║         YOUR POKÉDEX (%3d)        ║%s\n", colorGreen, len(config.Pokedex), colorReset)
	fmt.Printf("%s╚═══════════════════════════════════╝%s\n\n", colorGreen, colorReset)

	i := 1
	for _, pokemon := range config.Pokedex {
		fmt.Printf("%s%2d.%s %s\n", colorGray, i, colorReset, pokemon.Name)
		i++
	}

	fmt.Printf("\n%sUse 'inspect <pokemon_name>' to see details%s\n", colorGray, colorReset)
	return nil
}

func CommandExit(config *models.ReplConfig, args []string) error {
	fmt.Printf("\n%s╔═══════════════════════════════════╗%s\n", colorCyan, colorReset)
	fmt.Printf("%s║     Thanks for using Pokédex!    ║%s\n", colorCyan, colorReset)
	fmt.Printf("%s║      You caught %3d Pokémon       ║%s\n", colorCyan, len(config.Pokedex), colorReset)
	fmt.Printf("%s╚═══════════════════════════════════╝%s\n\n", colorCyan, colorReset)
	os.Exit(0)
	return nil
}

func CommandHelp(config *models.ReplConfig, args []string) error {
	fmt.Printf("%s╔═══════════════════════════════════════════════════════════╗%s\n", colorCyan, colorReset)
	fmt.Printf("%s║                  POKÉDEX COMMANDS                         ║%s\n", colorCyan, colorReset)
	fmt.Printf("%s╚═══════════════════════════════════════════════════════════╝%s\n\n", colorCyan, colorReset)

	// Group commands by category
	navigation := []string{"map", "mapb"}
	exploration := []string{"explore", "catch"}
	collection := []string{"pokedex", "inspect"}
	general := []string{"help", "exit"}

	printCommandGroup("Navigation", navigation)
	printCommandGroup("Exploration", exploration)
	printCommandGroup("Collection", collection)
	printCommandGroup("General", general)

	return nil
}

func printCommandGroup(title string, commands []string) {
	fmt.Printf("%s%s:%s\n", colorBold, title, colorReset)
	for _, cmdName := range commands {
		if cmd, exists := CommandsMap[cmdName]; exists {
			fmt.Printf("  %s%-25s%s %s\n",
				colorGreen, cmd.Name, colorReset, cmd.Description)
		}
	}
	fmt.Println()
}
