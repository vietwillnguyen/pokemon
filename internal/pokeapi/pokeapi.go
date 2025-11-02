// pokeapi.go
package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"pokedexcli/internal/pokecache"
	"time"
)

type Client struct {
	cache  *pokecache.Cache
	client *http.Client
}

func NewClient() *Client {
	return &Client{
		cache: pokecache.NewCache(1 * time.Minute),
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

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

	// Use cached value if it exists
	cachedVal, cachedValExists := c.cache.Get(url)
	if cachedValExists {
		log.Printf("cached value exists, key: %s, val: %s", url, cachedVal)
		var locationAreasListResponse LocationAreasListResponse
		err := json.Unmarshal(cachedVal, &locationAreasListResponse)
		if err != nil {
			return LocationAreasListResponse{}, fmt.Errorf("unexpected error: %w", err)
		}
		return locationAreasListResponse, nil
	}

	// Get locations
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationAreasListResponse{}, err
	}

	// Send request
	res, err := c.client.Do(req)
	if err != nil {
		return LocationAreasListResponse{}, err
	}
	defer res.Body.Close()

	// Check status code
	if res.StatusCode != http.StatusOK {
		return LocationAreasListResponse{}, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Read response body, convert http response body to byte slice
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationAreasListResponse{}, err
	}

	// Unmarshal the body into the LocationAreasListResponse struct
	var locationAreasListResponse LocationAreasListResponse
	err = json.Unmarshal(body, &locationAreasListResponse)
	if err != nil {
		return LocationAreasListResponse{}, err
	}
	// Add value to cache for later
	c.cache.Add(url, body)
	return locationAreasListResponse, nil
}

func (c *Client) GetLocationAreasDetail(locationName string) (LocationAreasDetailsResponse, error) {
	if locationName == "" {
		return LocationAreasDetailsResponse{}, fmt.Errorf("Must supply a location name")
	}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", locationName)

	// Use cached value if it exists
	cachedVal, cachedValExists := c.cache.Get(url)
	if cachedValExists {
		log.Printf("cached value exists, key: %s, val: %s", url, cachedVal)
		var locationAreasDetailsResponse LocationAreasDetailsResponse
		err := json.Unmarshal(cachedVal, &locationAreasDetailsResponse)
		if err != nil {
			return LocationAreasDetailsResponse{}, fmt.Errorf("unexpected error: %w", err)
		}
		return locationAreasDetailsResponse, nil
	}

	// Get locations
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationAreasDetailsResponse{}, err
	}

	// Send request
	res, err := c.client.Do(req)
	if err != nil {
		return LocationAreasDetailsResponse{}, err
	}
	defer res.Body.Close()

	// Check status code
	if res.StatusCode != http.StatusOK {
		return LocationAreasDetailsResponse{}, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Read response body, convert http response body to byte slice
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationAreasDetailsResponse{}, err
	}

	// Unmarshal the body into the LocationAreasDetailsResponse struct
	var locationAreasDetailsResponse LocationAreasDetailsResponse
	err = json.Unmarshal(body, &locationAreasDetailsResponse)
	if err != nil {
		return LocationAreasDetailsResponse{}, err
	}
	// Add value to cache for later
	c.cache.Add(url, body)
	return locationAreasDetailsResponse, nil

}
