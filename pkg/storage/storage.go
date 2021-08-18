package storage

import "github.com/DanilLagunov/jokes-api/pkg/models"

type Storage interface {
	GetJokes() ([]models.Joke, error)
	AddJoke(title, body string) error
	GetJokeByText(text string) ([]models.Joke, error)
	GetJokeByID(id string) (models.Joke, error)
	GetRandomJokes() ([]models.Joke, error)
	GetFunniestJokes() ([]models.Joke, error)
}
