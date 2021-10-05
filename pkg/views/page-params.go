package views

import (
	"github.com/DanilLagunov/jokes-api/pkg/models"
)

// JokesPageParams struct.
type JokesPageParams struct {
	Skip     int
	Seed     int
	CurrPage int
	MaxPage  int
	Content  []models.Joke
	Next     int
	Prev     int
}

// CreatePageParams creating a new JokesPageParams object.
func CreatePageParams(skip, seed, dataSize int, content []models.Joke) JokesPageParams {
	if skip >= dataSize || seed == 0 {
		return JokesPageParams{skip, seed, 0, 0, []models.Joke{}, 0, 0}
	}

	currPage := skip/seed + 1
	next := skip + seed
	prev := skip - seed

	var maxPage int
	if dataSize%seed != 0 {
		maxPage = dataSize/seed + 1
	} else {
		maxPage = dataSize / seed
	}

	if skip+seed >= dataSize {
		return JokesPageParams{skip, seed, currPage, maxPage, content[:dataSize-skip], next, prev}
	}

	return JokesPageParams{skip, seed, currPage, maxPage, content, next, prev}
}

// SearchPageParams struct.
type SearchPageParams struct {
	SearchRequest string
	PageParams    JokesPageParams
}
