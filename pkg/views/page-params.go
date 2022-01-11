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
func CreatePageParams(skip, limit, amount int, content []models.Joke) JokesPageParams {
	if skip >= amount || limit == 0 {
		return JokesPageParams{skip, limit, 0, 0, []models.Joke{}, 0, 0}
	}

	currPage := skip/limit + 1
	next := skip + limit
	prev := skip - limit

	var maxPage int
	if amount%limit != 0 {
		maxPage = amount/limit + 1
	} else {
		maxPage = amount / limit
	}

	if skip+limit >= amount {
		return JokesPageParams{skip, limit, currPage, maxPage, content[:amount-skip], next, prev}
	}

	return JokesPageParams{skip, limit, currPage, maxPage, content, next, prev}
}

// SearchPageParams struct.
type SearchPageParams struct {
	SearchRequest string
	PageParams    JokesPageParams
}
