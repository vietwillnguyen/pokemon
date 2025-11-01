package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type LocationsAreasResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
}

func GetLocations(urlOrOffset string) (LocationsAreasResponse, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Build URL - if urlOrOffset is empty or a full URL, use it directly
	// Otherwise, treat it as an offset number
	var url string
	if urlOrOffset == "" {
		url = "https://pokeapi.co/api/v2/location-area?limit=20"
	} else if strings.HasPrefix(urlOrOffset, "http") {
		// It's a full URL from the API response
		url = urlOrOffset
	} else {
		// It's an offset number
		url = fmt.Sprintf("https://pokeapi.co/api/v2/location?offset=%s&limit=20", urlOrOffset)
	}

	// Get locations
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationsAreasResponse{}, err
	}

	// Send request
	res, err := client.Do(req)
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
	// Cache these values for later

	return locationsAreasResponse, nil
}
