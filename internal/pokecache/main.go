// Package pokecache provides a simple cache for pokemon data
package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entries map[string]CacheEntry
	mu      sync.Mutex
}

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(d time.Duration) *Cache {
	nc := &Cache{
		entries: make(map[string]CacheEntry),
	}

	nc.reapLoop(d)

	return nc
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = CacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if v, ok := c.entries[key]; ok {
		return v.val, true
	} else {
		return nil, false
	}
}

func (c *Cache) reapLoop(d time.Duration) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for {
		<-ticker.C
		c.deleteExpired(d)
	}
}

func (c *Cache) deleteExpired(i time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.entries {
		if now.Sub(entry.createdAt) > i {
			delete(c.entries, key)
		}
	}
}
