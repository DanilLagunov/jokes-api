package models

import (
	"crypto/rand"
	"fmt"
)

type Joke struct {
	ID    string `json:"id" bson:"id"`
	Title string `json:"title" bson:"title"`
	Body  string `json:"body" bson:"body"`
	Score int    `json:"score" bson:"score"`
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
