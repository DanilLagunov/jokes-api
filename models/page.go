package models

type Page struct {
	Skip    int
	Seed    int
	Next    int
	Prev    int
	Content []Joke
}

func NewPage(skip, seed, next, prev int, content []Joke) Page {
	return Page{skip, seed, next, prev, content}
}
