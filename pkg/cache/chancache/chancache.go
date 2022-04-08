package chancache

import (
	"context"
	"time"

	"github.com/DanilLagunov/jokes-api/pkg/cache"
	"github.com/DanilLagunov/jokes-api/pkg/models"
)

// ChanCache struct.
type ChanCache struct {
	getCh             chan GetRequest
	setCh             chan SetRequest
	cleanCh           chan struct{}
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	items             map[string]*cache.Item
}

// GetRequest describes the controller of Get func.
type GetRequest struct {
	key    string
	respCh chan *cache.Item
}

// SetRequest describes the controller of Set func.
type SetRequest struct {
	key  string
	item chan *cache.Item
}

// NewChannelCache creates new ChanCache object.
func NewChannelCache(ctx context.Context, defaultExpiration, cleanupInterval time.Duration) *ChanCache {
	items := make(map[string]*cache.Item)

	cache := ChanCache{
		items:             items,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
		getCh:             make(chan GetRequest),
		setCh:             make(chan SetRequest),
		cleanCh:           make(chan struct{}),
	}

	go cache.ChannelBasedCacheController(ctx)

	if cleanupInterval > 0 {
		go cache.cleaner()
	}

	return &cache
}

// ChannelBasedCacheController is a controller for ChanCache.
func (c *ChanCache) ChannelBasedCacheController(ctx context.Context) {
	for {
		select {
		case getReq := <-c.getCh:
			getReq.respCh <- c.items[getReq.key]
		case setReq := <-c.setCh:
			c.items[setReq.key] = <-setReq.item
		case <-c.cleanCh:
			<-c.cleanCh
		case <-ctx.Done():
			return
		}
	}
}

// Get return cache item by key.
func (c *ChanCache) Get(key string) (models.Joke, error) {
	responseCh := make(chan *cache.Item, 1)
	getReq := GetRequest{
		key:    key,
		respCh: responseCh,
	}
	c.getCh <- getReq

	item := <-getReq.respCh

	if item == nil {
		return models.Joke{}, cache.ErrKeyNotFound
	}

	currentTime := time.Now().UnixNano()
	if item.Expiration > 0 && currentTime > item.Expiration {
		return models.Joke{}, cache.ErrItemExpired
	}

	return item.Value, nil
}

// Set puts new item into cache.
func (c *ChanCache) Set(key string, value models.Joke, duration time.Duration) {
	var expiration int64

	item := make(chan *cache.Item, 1)

	if duration == 0 {
		duration = c.defaultExpiration
	}

	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	setReq := SetRequest{
		key:  key,
		item: item,
	}
	setReq.item <- &cache.Item{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}
	c.setCh <- setReq
}

func (c *ChanCache) cleaner() {
	for {
		<-time.After(c.cleanupInterval)

		if c.items == nil {
			return
		}

		c.clearExpiredItems()
	}
}

func (c *ChanCache) clearExpiredItems() {
	c.cleanCh <- struct{}{}

	currentTime := time.Now().UnixNano()
	for k, i := range c.items {
		if currentTime > i.Expiration && i.Expiration > 0 {
			delete(c.items, k)
		}
	}

	c.cleanCh <- struct{}{}
	return
}
