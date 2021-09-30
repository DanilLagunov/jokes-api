package file_storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/DanilLagunov/jokes-api/pkg/models"
)

// ErrJokeNotFound describes the error when the joke is not found.
var ErrJokeNotFound = errors.New("joke not found")

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

// GetJokes method returns all jokes.
func (s *FileStorage) GetJokes() ([]models.Joke, error) {
	return s.Data, nil
}

// AddJoke method creating new joke.
func (s *FileStorage) AddJoke(title, body string) error {
	var id string
CHECK:
	id, err := models.GenerateID()
	if err != nil {
		return fmt.Errorf("ID generating error: %w", err)
	}

	for i := 0; i < len(s.Data); i++ {
		if id == s.Data[i].ID {
			goto CHECK
		}
	}

	joke := models.NewJoke(id, title, body, 0)
	s.Data = append(s.Data, joke)

	rawDataOut, err := json.MarshalIndent(&s.Data, "", "   ")
	if err != nil {
		return fmt.Errorf("marshalling error: %w", err)
	}

	err = ioutil.WriteFile(s.FilePath, rawDataOut, 0)
	if err != nil {
		return fmt.Errorf("cannot write: %w", err)
	}

	return nil
}

// GetJokeByText returns jokes which contain the desired text.
func (s *FileStorage) GetJokeByText(text string) ([]models.Joke, error) {
	var result []models.Joke

	for _, item := range s.Data {
		if strings.Contains(item.Title, text) || strings.Contains(item.Body, text) {
			result = append(result, item)
		}
	}
	if len(result) != 0 {
		return result, nil
	}
	return result, ErrJokeNotFound
}

// GetJokeByID returns joke that has the same id.
func (s *FileStorage) GetJokeByID(id string) (models.Joke, error) {
	for _, item := range s.Data {
		if item.ID == id {
			return item, nil
		}
	}
	return models.Joke{}, ErrJokeNotFound
}

// GetRandomJokes returns random jokes.
func (s *FileStorage) GetRandomJokes() ([]models.Joke, error) {
	random := make([]models.Joke, 0, len(s.Data))

	r := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(r)

	for i := 0; i < len(s.Data); i++ {
		random = append(random, s.Data[rnd.Intn(len(s.Data))])
	}

	return random, nil
}

// GetFunniestJokes returns jokes, sorted by score.
func (s *FileStorage) GetFunniestJokes() ([]models.Joke, error) {
	var funniest []models.Joke

	funniest = append(funniest, s.Data...)
	sort.Slice(funniest, func(i, j int) (less bool) {
		return funniest[i].Score > funniest[j].Score
	})

	return funniest, nil
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
