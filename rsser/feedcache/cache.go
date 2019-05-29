// Package feedcache provides a quick and dirty in-memory cache for feeds.Feed
// objects.
package feedcache

import (
	"sync"
	"time"

	"github.com/gorilla/feeds"
)

// CacheValue represents a value in the cache. The feed and expiry is stored.
type CacheValue struct {
	feed      *feeds.Feed
	expiresAt time.Time
}

// Cache represents the in-memory cache.
type Cache struct {
	cache map[string]CacheValue
	mu    *sync.RWMutex
}

// New creates a new cache.
func New() *Cache {
	c := &Cache{
		cache: make(map[string]CacheValue),
		mu:    &sync.RWMutex{},
	}
	return c
}

// Set sets a value in the cache with the given ttl.
func (c *Cache) Set(key string, value *feeds.Feed, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v := CacheValue{
		feed:      value,
		expiresAt: time.Now().Add(ttl),
	}
	c.cache[key] = v
}

// Get sets a value in the cache. Returns nil if object is not found in the
// cache or if it has expired.
func (c *Cache) Get(key string) *feeds.Feed {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if v, ok := c.cache[key]; ok {
		if time.Now().Before(v.expiresAt) {
			return v.feed
		}
		delete(c.cache, key)
	}
	return nil
}
