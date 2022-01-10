package cache

import (
	"errors"
	"sync"
	"time"

	"github.com/DanilLagunov/jokes-api/pkg/models"
)

// Cache struct.
type Cache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	items             map[string]Item
}

// Item struct.
type Item struct {
	Value      models.Joke
	Created    time.Time
	Expiration int64
}

// NewCache creating new Cache object.
func NewCache(defaultExpiration, cleanupInterval time.Duration) *Cache {
	items := make(map[string]Item)

	cache := Cache{
		items:             items,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}

	if cleanupInterval > 0 {
		cache.startCleaner()
	}

	return &cache
}

// Get return cache item by key.
func (c *Cache) Get(key string) (models.Joke, error) {
	c.RLock()

	defer c.RUnlock()

	item, found := c.items[key]
	if !found {
		return models.Joke{}, errors.New("key not found")
	}

	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return models.Joke{}, errors.New("item expired")
		}
	}

	return item.Value, nil
}

// Set puts new item into cache.
func (c *Cache) Set(key string, value models.Joke, duration time.Duration) {
	var expiration int64

	if duration == 0 {
		duration = c.defaultExpiration
	}

	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.Lock()

	defer c.Unlock()

	c.items[key] = Item{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}
}

func (c *Cache) startCleaner() {
	go c.cleaner()
}

func (c *Cache) cleaner() {
	for {
		<-time.After(c.cleanupInterval)

		if c.items == nil {
			return
		}

		if keys := c.findExpiredKeys(); len(keys) != 0 {
			c.clearItems(keys)
		}
	}
}

func (c *Cache) clearItems(keys []string) {
	c.Lock()

	defer c.Unlock()

	for _, k := range keys {
		delete(c.items, k)
	}
}

func (c *Cache) findExpiredKeys() (keys []string) {
	c.RLock()

	defer c.RUnlock()

	for k, i := range c.items {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}

	return
}
