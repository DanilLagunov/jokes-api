package storage

import (
	"context"
	"errors"

	"github.com/DanilLagunov/jokes-api/pkg/models"
)

// ErrJokeNotFound describes the error when the joke is not found.
var ErrJokeNotFound = errors.New("joke not found")

// Storage interface.
type Storage interface {
	GetJokes(ctx context.Context, skip, seed int) ([]models.Joke, int, error)
	AddJoke(ctx context.Context, title, body string, score int) (models.Joke, error)
	GetJokesByText(ctx context.Context, skip, seed int, text string) ([]models.Joke, int, error)
	GetJokeByID(ctx context.Context, id string) (models.Joke, error)
	GetRandomJokes(ctx context.Context, seed int) ([]models.Joke, int, error)
	GetFunniestJokes(ctx context.Context, skip, seed int) ([]models.Joke, int, error)
}
