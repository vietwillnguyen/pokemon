// pokeapi.go
package pokeapi

import (
	"fmt"
)

// For listing/pagination (map/mapb commands)
type LocationAreasListResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationAreasDetailsResponse struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	Next                 string `json:"next"`
	Previous             string `json:"previous"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int           `json:"min_level"`
				MaxLevel        int           `json:"max_level"`
				ConditionValues []interface{} `json:"condition_values"`
				Chance          int           `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func (c *Client) GetLocationAreasList(url string, args []string) (LocationAreasListResponse, error) {
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	}
	if len(args) > 0 {
		url = fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", args[0])
	}
	locationAreasDetailsResponse := LocationAreasListResponse{}
	err := c.fetchJSON(url, &locationAreasDetailsResponse)
	if err != nil {
		return locationAreasDetailsResponse, nil
	}
	return locationAreasDetailsResponse, nil
}

func (c *Client) GetLocationAreasDetail(locationName string) (LocationAreasDetailsResponse, error) {
	if locationName == "" {
		return LocationAreasDetailsResponse{}, fmt.Errorf("Must supply a location name")
	}

	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", locationName)
	locationAreasDetailsResponse := LocationAreasDetailsResponse{}
	err := c.fetchJSON(url, &locationAreasDetailsResponse)
	if err != nil {
		return locationAreasDetailsResponse, nil
	}
	return locationAreasDetailsResponse, nil
}
