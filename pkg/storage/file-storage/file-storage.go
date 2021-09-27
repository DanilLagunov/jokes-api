package file_storage

import (
	"context"
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

var JokeNotFound = errors.New("joke not found")

type FileStorage struct {
	FilePath string
	Data     []models.Joke
}

func NewFileStorage(filePath string) *FileStorage {
	var storage FileStorage
	storage.FilePath = filePath
	err := parseJSON(storage.FilePath, &storage.Data)
	if err != nil {
		return &FileStorage{}
	}
	return &storage
}

func (s *FileStorage) GetJokes(ctx context.Context) ([]models.Joke, error) {
	defer ctx.Done()
	return s.Data, nil
}

func (s *FileStorage) AddJoke(ctx context.Context, title, body string) (models.Joke, error) {
	defer ctx.Done()
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

	joke := models.NewJoke(id, title, body, 0)
	s.Data = append(s.Data, joke)

	rawDataOut, err := json.MarshalIndent(&s.Data, "", "   ")
	if err != nil {
		return joke, fmt.Errorf("Marshalling error: %w", err)
	}
	err = ioutil.WriteFile(s.FilePath, rawDataOut, 0)
	if err != nil {
		return joke, fmt.Errorf("Cannot write: %w", err)
	}

	return joke, nil
}

func (s *FileStorage) GetJokeByText(ctx context.Context, text string) ([]models.Joke, error) {
	defer ctx.Done()
	var result []models.Joke
	for _, item := range s.Data {
		if strings.Contains(item.Title, text) || strings.Contains(item.Body, text) {
			result = append(result, item)
		}
	}
	if len(result) != 0 {
		return result, nil
	}
	return result, JokeNotFound
}

func (s *FileStorage) GetJokeByID(ctx context.Context, id string) (models.Joke, error) {
	defer ctx.Done()
	for _, item := range s.Data {
		if item.ID == id {
			return item, nil
		}
	}
	return models.Joke{}, JokeNotFound
}

func (s *FileStorage) GetRandomJokes(ctx context.Context) ([]models.Joke, error) {
	defer ctx.Done()
	r := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(r)
	var random []models.Joke

	for i := 0; i < 300; i++ {
		random = append(random, s.Data[rnd.Intn(len(s.Data))])
	}

	return random, nil
}

func (s *FileStorage) GetFunniestJokes(ctx context.Context) ([]models.Joke, error) {
	defer ctx.Done()
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
		return fmt.Errorf("Opening file error: %w", err)
	}
	decoder := json.NewDecoder(file)

	err = decoder.Decode(&list)
	if err != nil {
		return fmt.Errorf("Decode error: %w", err)
	}
	return nil
}
