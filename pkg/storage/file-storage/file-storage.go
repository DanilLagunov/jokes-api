package file_storage

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

type FileStorage struct {
	FilePath string
	Data     []models.Joke
}

func NewFileStorage(filePath string) FileStorage {
	var storage FileStorage
	storage.FilePath = filePath
	parseJSON(storage.FilePath, &storage.Data)
	return storage
}

func (s FileStorage) GetJokes() ([]models.Joke, error) {
	return s.Data, nil
}

func (s FileStorage) AddJoke(title, body string) error {
	var id string
CHECK:
	id = models.GenerateID()
	for i := 0; i < len(s.Data); i++ {
		if id == s.Data[i].ID {
			goto CHECK
		}
	}

	joke := models.NewJoke(id, title, body, 0)
	s.Data = append(s.Data, joke)

	rawDataOut, err := json.MarshalIndent(&s, "", "   ")
	if err != nil {
		log.Fatal("JSON marshalling failed: ", err)
	}
	err = ioutil.WriteFile(s.FilePath, rawDataOut, 0)
	if err != nil {
		log.Fatal("Cannot write:", err)
	}

	return nil
}

func (s FileStorage) GetJokeByText(text string) ([]models.Joke, error) {
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

func (s FileStorage) GetJokeByID(id string) (models.Joke, error) {
	for _, item := range s.Data {
		if item.ID == id {
			return item, nil
		}
	}
	return models.Joke{}, JokeNotFound
}

func (s FileStorage) GetRandomJokes() ([]models.Joke, error) {
	r := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(r)
	var random []models.Joke

	for i := 0; i < 300; i++ {
		random = append(random, s.Data[rnd.Intn(len(s.Data))])
	}

	return random, nil
}

func (s FileStorage) GetFunniestJokes() ([]models.Joke, error) {
	var funniest []models.Joke

	funniest = append(funniest, s.Data...)
	sort.Slice(funniest, func(i, j int) (less bool) {
		return funniest[i].Score > funniest[j].Score
	})

	return funniest, nil
}

func parseJSON(path string, list *[]models.Joke) {
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
