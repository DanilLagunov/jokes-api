package fs

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/DanilLagunov/jokes-api/pkg/models"
	"github.com/DanilLagunov/jokes-api/pkg/storage"
)

// FileStorage struct.
type FileStorage struct {
	FilePath string
	Data     []models.Joke
}

// NewFileStorage creating a new FileStorage object.
func NewFileStorage(filePath string) *FileStorage {
	var storage FileStorage

	storage.FilePath = filePath

	err := parseJSON(storage.FilePath, &storage.Data)
	if err != nil {
		return &FileStorage{}
	}
	return &storage
}

// GetJokes method returns the number of jokes given by skip and limit parameters and total amount of jokes.
func (s *FileStorage) GetJokes(ctx context.Context, skip, seed int) ([]models.Joke, int, error) {
	if skip > len(s.Data) {
		return []models.Joke{}, 0, nil
	}
	if seed > len(s.Data) {
		return s.Data[skip:len(s.Data)], len(s.Data), nil
	}
	return s.Data[skip : skip+seed], len(s.Data), nil
}

// AddJoke method creating new joke.
func (s *FileStorage) AddJoke(ctx context.Context, title, body string, score int) (models.Joke, error) {
	var id string
CHECK:
	id, err := models.GenerateID()
	if err != nil {
		return models.Joke{}, fmt.Errorf("ID generating error: %w", err)
	}

	for i := 0; i < len(s.Data); i++ {
		if id == s.Data[i].ID {
			goto CHECK
		}
	}

	joke := models.NewJoke(id, title, body, score)
	s.Data = append(s.Data, joke)

	rawDataOut, err := json.MarshalIndent(&s.Data, "", "   ")
	if err != nil {
		return joke, fmt.Errorf("marshalling error: %w", err)
	}

	err = ioutil.WriteFile(s.FilePath, rawDataOut, 0)
	if err != nil {
		return joke, fmt.Errorf("cannot write: %w", err)
	}

	return joke, nil
}

// GetJokesByText returns the number jokes, which contain the desired text, given by skip and limit parameters and total amount of jokes.
func (s *FileStorage) GetJokesByText(ctx context.Context, skip, seed int, text string) ([]models.Joke, int, error) {
	var result []models.Joke

	for _, item := range s.Data {
		if strings.Contains(item.Title, text) || strings.Contains(item.Body, text) {
			result = append(result, item)
		}
	}
	if len(result) != 0 {
		if skip > len(result) {
			return []models.Joke{}, 0, nil
		}
		if seed > len(result) {
			return s.Data[skip:len(result)], len(result), nil
		}
		return result[skip : skip+seed], len(result), nil
	}
	return result, 0, storage.ErrJokeNotFound
}

// GetJokeByID returns joke that has the same id.
func (s *FileStorage) GetJokeByID(ctx context.Context, id string) (models.Joke, error) {
	for _, item := range s.Data {
		if item.ID == id {
			return item, nil
		}
	}
	return models.Joke{}, storage.ErrJokeNotFound
}

// GetRandomJokes returns the number of random jokes given by limit parameter and total amount of jokes.
func (s *FileStorage) GetRandomJokes(ctx context.Context, seed int) ([]models.Joke, int, error) {
	r := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(r)

	random := make([]models.Joke, seed)

	for i := 0; i < seed; i++ {
		random = append(random, s.Data[rnd.Intn(len(s.Data))])
	}
	if seed > len(s.Data) {
		return random[:len(s.Data)], len(s.Data), nil
	}
	return random, len(random), nil
}

// GetFunniestJokes returns the number of sorted jokes, given by skip and limit parameters and total amount of jokes.
func (s *FileStorage) GetFunniestJokes(ctx context.Context, skip, seed int) ([]models.Joke, int, error) {
	var funniest []models.Joke

	funniest = append(funniest, s.Data...)
	sort.Slice(funniest, func(i, j int) (less bool) {
		return funniest[i].Score > funniest[j].Score
	})

	if skip > len(funniest) {
		return []models.Joke{}, 0, nil
	}
	if seed > len(funniest) {
		return funniest[skip:], len(funniest), nil
	}
	return funniest[skip : skip+seed], len(funniest), nil
}

func parseJSON(path string, list *[]models.Joke) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("0pening file error: %w", err)
	}

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&list)
	if err != nil {
		return fmt.Errorf("decode error: %w", err)
	}
	return nil
}
