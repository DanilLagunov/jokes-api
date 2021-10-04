package storage

import (
	"context"

	"github.com/DanilLagunov/jokes-api/pkg/models"
)

// Storage interface.
type Storage interface {
	GetJokes(ctx context.Context) ([]models.Joke, error)
	AddJoke(ctx context.Context, title, body string) (models.Joke, error)
	GetJokeByText(ctx context.Context, text string) ([]models.Joke, error)
	GetJokeByID(ctx context.Context, id string) (models.Joke, error)
	GetRandomJokes(ctx context.Context) ([]models.Joke, error)
	GetFunniestJokes(ctx context.Context) ([]models.Joke, error)
}
