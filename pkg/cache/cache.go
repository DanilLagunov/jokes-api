package cache

import (
	"errors"
	"time"

	"github.com/DanilLagunov/jokes-api/pkg/models"
)

var (
	// ErrKeyNotFound describes the error when the key is not found.
	ErrKeyNotFound = errors.New("key not found")
	// ErrItemExpired describes the error when the item is expired.
	ErrItemExpired = errors.New("item expired")
)

// Cache interface.
type Cache interface {
	Get(key string) (models.Joke, error)
	Set(key string, value models.Joke, duration time.Duration)
}
