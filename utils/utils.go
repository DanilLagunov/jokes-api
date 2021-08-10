package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"

	"github.com/DanilLagunov/jokes-api/models"
)

func ParseJSON(path string, list *[]models.Joke) {
	file, _ := os.Open(path)
	decoder := json.NewDecoder(file)

	err := decoder.Decode(&list)
	if err != nil {
		fmt.Println("Decode error")
	}
}

func GenerateId() string {
	b := make([]byte, 3)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
