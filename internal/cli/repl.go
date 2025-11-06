package cli

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/internal/models"
	"pokedexcli/internal/pokeapi"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
	colorBold   = "\033[1m"
)

// CleanInput normalizes user input by lowercasing and splitting into words
func CleanInput(input string) []string {
	lowered := strings.ToLower(input)
	words := strings.Fields(lowered)
	return words
}

// printBanner displays the Pokedex ASCII banner
func printBanner() {
	banner := `
╔═══════════════════════════════════════╗
║                                       ║
║   ██████╗  ██████╗ ██╗  ██╗███████╗  ║
║   ██╔══██╗██╔═══██╗██║ ██╔╝██╔════╝  ║
║   ██████╔╝██║   ██║█████╔╝ █████╗    ║
║   ██╔═══╝ ██║   ██║██╔═██╗ ██╔══╝    ║
║   ██║     ╚██████╔╝██║  ██╗███████╗  ║
║   ╚═╝      ╚═════╝ ╚═╝  ╚═╝╚══════╝  ║
║                                       ║
║         COMMAND LINE INTERFACE        ║
║                                       ║
╚═══════════════════════════════════════╝
`
	fmt.Print(colorCyan + banner + colorReset)
	fmt.Println("\nType 'help' to see available commands\n")
}

// printPrompt displays a styled prompt with status info
func printPrompt(config *models.ReplConfig) {
	caughtCount := len(config.Pokedex)
	fmt.Printf("%s[%d caught]%s %sPokedex >%s ",
		colorGray, caughtCount, colorReset,
		colorGreen, colorReset)
}

// printError displays an error message in red
func printError(err error) {
	fmt.Printf("%s✗ Error:%s %v\n", colorRed, colorReset, err)
}

// printSuccess displays a success message in green
func printSuccess(message string) {
	fmt.Printf("%s✓%s %s\n", colorGreen, colorReset, message)
}

// printWarning displays a warning message in yellow
func printWarning(message string) {
	fmt.Printf("%s⚠%s %s\n", colorYellow, colorReset, message)
}

// findSimilarCommand suggests similar commands for typos
func findSimilarCommand(input string) string {
	for cmd := range CommandsMap {
		if strings.HasPrefix(cmd, input[:min(len(input), len(cmd))]) {
			return cmd
		}
	}
	return ""
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// clearScreen clears the terminal (works on Unix-like systems)
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

// StartREPL initializes and starts the REPL loop
func StartREPL() {
	reader := bufio.NewReader(os.Stdin)
	config := &models.ReplConfig{
		Pokedex:       map[string]models.Pokemon{},
		PokeApiClient: pokeapi.NewClient(),
	}

	clearScreen()
	printBanner()

	for {
		printPrompt(config)
		input, err := reader.ReadString('\n')
		if err != nil {
			printError(fmt.Errorf("reading input: %w", err))
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
			similar := findSimilarCommand(command)
			if similar != "" {
				printWarning(fmt.Sprintf("Unknown command '%s'. Did you mean '%s'?", command, similar))
			} else {
				printWarning(fmt.Sprintf("Unknown command '%s'. Type 'help' for available commands.", command))
			}
			continue
		}

		fmt.Println() // Add spacing before command output
		err = cmd.Callback(config, args)
		if err != nil {
			printError(err)
		}
		fmt.Println() // Add spacing after command output
	}
}
