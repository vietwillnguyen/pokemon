// client.go
package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
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
		cache: pokecache.NewCache(10 * time.Minute),
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// fetchJSON is a private helper that handles the common HTTP + cache pattern
func (c *Client) fetchJSON(url string, target interface{}) error {
	// Check if cached val exists
	cachedVal, cachedValExists := c.cache.Get(url)
	if cachedValExists {
		// log.Printf("cached value exists, key: %s, val: %s", url, cachedVal)
		err := json.Unmarshal(cachedVal, &target)
		if err != nil {
			return fmt.Errorf("unexpected unmarshal error: %w", err)
		}
		return nil
	}
	// Cached val does not exist, must make request
	// Get locations
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Send request
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Check status code
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("request had unexpected status code: %d", res.StatusCode)
	}

	// Read response body, convert http response body to byte slice
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// Unmarshal the body into the target
	err = json.Unmarshal(body, target)
	if err != nil {
		return err
	}
	// Add value to cache for later
	c.cache.Add(url, body)
	return nil
}
