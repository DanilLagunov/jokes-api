package memcache

import (
	"sync"
	"time"

	"github.com/DanilLagunov/jokes-api/pkg/cache"
	"github.com/DanilLagunov/jokes-api/pkg/models"
)

// MemCache struct.
type MemCache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	items             map[string]*cache.Item
}

// Item struct.

// NewMemCache creates new MemCache object.
func NewMemCache(defaultExpiration, cleanupInterval time.Duration) *MemCache {
	items := make(map[string]*cache.Item)

	cache := MemCache{
		items:             items,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}

	if cleanupInterval > 0 {
		go cache.cleaner()
	}

	return &cache
}

// Get return cache item by key.
func (c *MemCache) Get(key string) (models.Joke, error) {
	c.RLock()

	defer c.RUnlock()

	item, found := c.items[key]
	if !found {
		return models.Joke{}, cache.ErrKeyNotFound
	}

	currentTime := time.Now().UnixNano()
	if item.Expiration > 0 {
		if currentTime > item.Expiration {
			return models.Joke{}, cache.ErrItemExpired
		}
	}

	return item.Value, nil
}

// Set puts new item into cache.
func (c *MemCache) Set(key string, value models.Joke, duration time.Duration) {
	var expiration int64

	if duration == 0 {
		duration = c.defaultExpiration
	}

	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.Lock()

	defer c.Unlock()

	c.items[key] = &cache.Item{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}
}

func (c *MemCache) cleaner() {
	for {
		<-time.After(c.cleanupInterval)

		if c.items == nil {
			return
		}

		c.clearExpiredItems()
	}
}

func (c *MemCache) clearExpiredItems() {
	c.Lock()

	defer c.Unlock()

	currentTime := time.Now().UnixNano()
	for k, i := range c.items {
		if currentTime > i.Expiration && i.Expiration > 0 {
			delete(c.items, k)
		}
	}

	return
}
