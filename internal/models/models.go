package models

import "pokedexcli/internal/pokeapi"

// Pokemon represents a caught Pokemon
type Pokemon struct {
	Name string
}

// ReplConfig holds the state of the REPL session
type ReplConfig struct {
	Pokedex       map[string]Pokemon
	PokeApiClient *pokeapi.Client
	Next          string
	Previous      string
}
