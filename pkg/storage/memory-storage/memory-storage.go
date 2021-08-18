package memory_storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/DanilLagunov/jokes-api/pkg/models"
)

var JokeNotFound = errors.New("joke not found")

type MemoryStorage []models.Joke

func NewMemoryStorage() MemoryStorage {
	var storage MemoryStorage
	parseJSON("../pkg/storage/memory-storage/reddit_jokes.json", &storage)
	return storage
}

func (s MemoryStorage) GetJokes() ([]models.Joke, error) {
	return s, nil
}

func (s MemoryStorage) AddJoke(title, body string) error {
	var id string
CHECK:
	id = models.GenerateID()
	for i := 0; i < len(s); i++ {
		if id == s[i].ID {
			goto CHECK
		}
	}

	joke := models.NewJoke(id, title, body, 0)
	s = append(s, joke)

	rawDataOut, err := json.MarshalIndent(&s, "", "   ")
	if err != nil {
		log.Fatal("JSON marshalling failed: ", err)
	}
	err = ioutil.WriteFile("../pkg/storage/memory-storage/reddit_jokes.json", rawDataOut, 0)
	if err != nil {
		log.Fatal("Cannot write:", err)
	}

	return nil
}

func (s MemoryStorage) GetJokeByText(text string) ([]models.Joke, error) {
	var result []models.Joke
	for _, item := range s {
		if strings.Contains(item.Title, text) || strings.Contains(item.Body, text) {
			result = append(result, item)
		}
	}
	if len(result) != 0 {
		return result, nil
	}
	return result, JokeNotFound
}

func (s MemoryStorage) GetJokeByID(id string) (models.Joke, error) {
	for _, item := range s {
		if item.ID == id {
			return item, nil
		}
	}
	return models.Joke{}, JokeNotFound
}

func (s MemoryStorage) GetRandomJokes() ([]models.Joke, error) {
	r := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(r)
	var random []models.Joke

	for i := 0; i < 300; i++ {
		random = append(random, s[rnd.Intn(len(s))])
	}

	return random, nil
}

func (s MemoryStorage) GetFunniestJokes() ([]models.Joke, error) {
	var funniest []models.Joke

	funniest = append(funniest, s...)
	sort.Slice(funniest, func(i, j int) (less bool) {
		return funniest[i].Score > funniest[j].Score
	})

	return funniest, nil
}

func parseJSON(path string, list *MemoryStorage) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Opening file error: %w", err)
	}
	decoder := json.NewDecoder(file)

	err = decoder.Decode(&list)
	if err != nil {
		fmt.Println("Decode error")
	}
}
