# Pokedex CLI

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
Pokedex > 
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
Pokedex > map
Pokedex > explore pallet-town
Pokedex > catch pikachu
Pokedex > inspect pikachu
Pokedex > pokedex
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

