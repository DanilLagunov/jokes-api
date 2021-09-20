package models

import (
	"crypto/rand"
	"fmt"
)

type Joke struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
	Score int    `json:"score"`
}

func NewJoke(id, title, body string, score int) Joke {
	return Joke{id, title, body, score}
}

func GenerateID() (string, error) {
	b := make([]byte, 3)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("Byte reading error: %w", err)
	}
	return fmt.Sprintf("%x", b), nil
}
