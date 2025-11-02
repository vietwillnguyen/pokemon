// pokecache.go
package pokecache

import (
	"log"
	"sync"
	"time"
)

type Cache struct {
	entries  map[string]CacheEntry
	mutex    sync.RWMutex
	inverval time.Duration
}

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(inverval time.Duration) *Cache {
	c := Cache{
		entries:  make(map[string]CacheEntry),
		inverval: inverval,
	}
	go c.reapLoop()
	log.Printf("NewCache: %v interval created\n", inverval)
	return &c
}

/*
*
- each time an inverval passes
- remove any entries older than inverval.
*/
func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.inverval)
	for {
		<-ticker.C
		log.Printf("reapLoop: %v interval passed.\n", c.inverval)
		for key, cacheEntry := range c.entries {
			// if expired, delete
			if time.Now().After(cacheEntry.createdAt) {
				log.Printf("reapLoop: deleted old entry at key: %v\n", key)
				delete(c.entries, key)
			}
		}
	}
}

func (c *Cache) Add(key string, val []byte) {
	log.Printf("Add: key %v\n", key)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.entries[key] = CacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	log.Printf("Get: key %v\n", key)
	// Lock the resource specifically to read only actions.
	// This allows other reads to happen
	// Once a mutex.Lock() is called, RLock() operations are blocked
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	cacheEntry, cacheEntryExists := c.entries[key]
	if !cacheEntryExists {
		return []byte{}, false
	}
	return cacheEntry.val, true
}
