package models

import (
	"crypto/rand"
	"fmt"
)

// Joke struct.
type Joke struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
	Score int    `json:"score"`
}

// NewJoke creating a new Joke object.
func NewJoke(id, title, body string, score int) Joke {
	return Joke{id, title, body, score}
}

// GenerateID is a function to generate ID for the Joke.
func GenerateID() (string, error) {
	b := make([]byte, 3)

	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("byte reading error: %w", err)
	}

	return fmt.Sprintf("%x", b), nil
}
