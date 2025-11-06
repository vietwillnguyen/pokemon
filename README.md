# Pokedex CLI

```
╔═══════════════════════════════════════╗
║                                       ║
║   ██████╗  ██████╗ ██╗  ██╗███████╗   ║ 
║   ██╔══██╗██╔═══██╗██║ ██╔╝██╔════╝   ║ 
║   ██████╔╝██║   ██║█████╔╝ █████╗     ║ 
║   ██╔═══╝ ██║   ██║██╔═██╗ ██╔══╝     ║ 
║   ██║     ╚██████╔╝██║  ██╗███████╗   ║ 
║   ╚═╝      ╚═════╝ ╚═╝  ╚═╝╚══════╝   ║ 
║                                       ║
║         COMMAND LINE INTERFACE        ║
║                                       ║
╚═══════════════════════════════════════╝

Type 'help' to see available commands    
```

A command-line interface for exploring the Pokemon world using the [PokeAPI](https://pokeapi.co/).

## Features

- Browse Pokemon location areas
- Explore areas to see available Pokemon
- Catch Pokemon with probability-based mechanics
- Inspect caught Pokemon details
- Manage your Pokedex collection

## Installation

```bash
# Clone the repository
git clone <repository-url>
cd pokedexcli

# Build the application
make build
# or
go build -o pokedexcli ./cmd/pokedexcli

# Run the application
./pokedexcli
# or
make run
```

## Usage

Once the application is running, you'll see a REPL prompt:

```
[0 caught] Pokedex > map
```

### Available Commands

- `help` - Display a help message with all available commands
- `map` - Display the next 20 location areas
- `mapb` - Display the previous 20 location areas
- `explore <area_name>` - List all Pokemon in a specific area
- `catch <pokemon_name>` - Attempt to catch a Pokemon
- `inspect <pokemon_name>` - View detailed information about a caught Pokemon
- `pokedex` - List all Pokemon in your collection
- `exit` - Exit the application

### Examples

```
[0 caught] Pokedex > map

═══ Locations ═══
 1. canalave-city-area        
 2. eterna-city-area
 3. pastoria-city-area        
 4. sunyshore-city-area       
 5. sinnoh-pokemon-league-area
 6. oreburgh-mine-1f
 7. oreburgh-mine-b1f
 8. valley-windworks-area     
 9. eterna-forest-area        
10. fuego-ironworks-area      
11. mt-coronet-1f-route-207   
12. mt-coronet-2f
13. mt-coronet-3f
14. mt-coronet-exterior-snowfall
15. mt-coronet-exterior-blizzard
16. mt-coronet-4f
17. mt-coronet-4f-small-room
18. mt-coronet-5f
19. mt-coronet-6f
20. mt-coronet-1f-from-exterior

Type 'map' for more locations

[0 caught] Pokedex > explore mt-coronet-exterior-snowfall

Exploring mt-coronet-exterior-snowfall...

═══ Pokémon Found in mt-coronet-exterior-snowfall ═══
 1. clefairy
 2. golbat
 3. machoke
 4. noctowl
 5. loudred
 6. nosepass
 7. medicham
 8. lunatone
 9. solrock
10. absol
11. chingling
12. bronzong
13. snover
14. abomasnow

Use 'catch <pokemon_name>' to attempt a catch!

[0 caught] Pokedex > catch abomasnow

Throwing a Pokéball at abomasnow...
Wobble... Wobble... Wobble... 
✗ Oh no! abomasnow broke free!
  Catch rate: 76.0% - Try again!

[0 caught] Pokedex > catch abomasnow

Throwing a Pokéball at abomasnow...
Wobble... Wobble... Wobble...
✓ Gotcha! abomasnow was caught!
  Base Experience: 173

[1 caught] Pokedex > inspect abomasnow


╔═══════════════════════════════════╗
║  ABOMASNOW                        ║
╚═══════════════════════════════════╝

Height: 22 decimetres
Weight: 1355 hectograms

Types:
  • grass
  • ice

Stats:
  hp:                 90 ███████░░░░░░░░░░░░░
  attack:             92 ███████░░░░░░░░░░░░░
  defense:            75 █████░░░░░░░░░░░░░░░
  special attack:     92 ███████░░░░░░░░░░░░░
  special defense:    85 ██████░░░░░░░░░░░░░░
  speed:              60 ████░░░░░░░░░░░░░░░░

[1 caught] Pokedex >
```

## Project Structure

```
pokedexcli/
├── cmd/
│   └── pokedexcli/          # Application entry point
│       └── main.go
├── internal/
│   ├── cli/                  # CLI command implementations
│   │   ├── commands.go
│   │   ├── repl.go
│   │   └── repl_test.go
│   ├── models/               # Domain models
│   │   └── models.go
│   ├── pokeapi/              # PokeAPI client
│   │   ├── client.go
│   │   ├── locations.go
│   │   └── pokemon.go
│   └── pokecache/            # HTTP response caching
│       ├── pokecache.go
│       └── pokecache_test.go
├── go.mod
├── Makefile
└── README.md
```

## Development

### Running Tests

```bash
make test
# or
go test ./...
```

### Building

```bash
make build
# or
go build -o pokedexcli ./cmd/pokedexcli
```

### Cleaning

```bash
make clean
```

## Architecture

- **CLI Layer** (`internal/cli/`): Handles user interaction and command routing
- **Models** (`internal/models/`): Domain models and application state
- **API Client** (`internal/pokeapi/`): PokeAPI integration with HTTP client
- **Cache** (`internal/pokecache/`): In-memory cache for API responses

The application uses a persistent HTTP client with caching to minimize API calls and improve performance.

## License

[Add your license here]

