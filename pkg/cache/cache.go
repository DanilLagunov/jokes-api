package cache

import (
	"errors"
	"time"

	"github.com/DanilLagunov/jokes-api/pkg/models"
)

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrItemExpired = errors.New("item expired")
)

type Cache interface {
	Get(key string) (models.Joke, error)
	Set(key string, value models.Joke, duration time.Duration)
}
