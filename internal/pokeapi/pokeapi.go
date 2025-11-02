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

type LocationsAreasResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
}

func (c *Client) GetLocationAreas(url string) (LocationsAreasResponse, error) {

	// Build URL - if url is empty or a full URL, use it directly
	// Otherwise, treat it as an offset number
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area?limit=20"
	}

	// Use cached value if it exists
	cachedVal, cachedValExists := c.cache.Get(url)
	if cachedValExists {
		log.Printf("cached value exists, key: %s, val: %s", url, cachedVal)
		var locationsAreasResponse LocationsAreasResponse
		err := json.Unmarshal(cachedVal, &locationsAreasResponse)
		if err != nil {
			return LocationsAreasResponse{}, fmt.Errorf("unexpected error: %w", err)
		}
		return locationsAreasResponse, nil
	}

	// Get locations
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationsAreasResponse{}, err
	}

	// Send request
	res, err := c.client.Do(req)
	if err != nil {
		return LocationsAreasResponse{}, err
	}
	defer res.Body.Close()

	// Check status code
	if res.StatusCode != http.StatusOK {
		return LocationsAreasResponse{}, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Read response body, convert http response body to byte slice
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationsAreasResponse{}, err
	}

	// Unmarshal the body into the LocationsAreasResponse struct
	var locationsAreasResponse LocationsAreasResponse
	err = json.Unmarshal(body, &locationsAreasResponse)
	if err != nil {
		return locationsAreasResponse, err
	}
	// Add value to cache for later
	c.cache.Add(url, body)
	return locationsAreasResponse, nil
}
