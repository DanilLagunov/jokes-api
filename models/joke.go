package models

type Joke struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
	Score int    `json:"score"`
}

func NewJoke(id, title, body string, score int) Joke {
	return Joke{id, title, body, score}
}
