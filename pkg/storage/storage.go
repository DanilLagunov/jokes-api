package storage

import (
	"context"

	"github.com/DanilLagunov/jokes-api/pkg/models"
)

// Storage interface.
type Storage interface {
	GetJokes(ctx context.Context, skip, seed int) ([]models.Joke, int, error)
	AddJoke(ctx context.Context, title, body string) (models.Joke, error)
	GetJokesByText(ctx context.Context, skip, seed int, text string) ([]models.Joke, int, error)
	GetJokeByID(ctx context.Context, id string) (models.Joke, error)
	GetRandomJokes(ctx context.Context, seed int) ([]models.Joke, int, error)
	GetFunniestJokes(ctx context.Context, skip, seed int) ([]models.Joke, int, error)
}
