package views

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/DanilLagunov/jokes-api/pkg/models"
)

type PaginationData struct {
	Skip     int
	Seed     int
	CurrPage int
	MaxPage  int
	Content  []models.Joke
	Next     int
	Prev     int
}

func CreatePaginationData(skip, seed int, content []models.Joke) PaginationData {
	if skip > len(content) || seed == 0 {
		return PaginationData{skip, seed, 0, 0, []models.Joke{}, 0, 0}
	}

	currPage := skip/seed + 1
	next := skip + seed
	prev := skip - seed

	var maxPage int
	if len(content)%seed != 0 {
		maxPage = len(content)/seed + 1
	} else {
		maxPage = len(content) / seed
	}

	if skip+seed >= len(content) {
		return PaginationData{skip, seed, currPage, maxPage, content[skip:], next, prev}
	}

	return PaginationData{skip, seed, currPage, maxPage, content[skip : skip+seed], next, prev}
}

func GetPaginationParams(r *http.Request) (int, int, error) {
	var skip, seed int
	var err error
	skipStr := r.URL.Query().Get("skip")
	if skipStr == "" {
		fmt.Println("Skip is not specified, using default value")
		skip = 0
	} else {
		skip, err = strconv.Atoi(skipStr)
		if err != nil {
			return 0, 0, fmt.Errorf("skip is not valid: %w", err)
		}
		if skip < 0 {
			return 0, 0, fmt.Errorf("skip is negative: %w", err)
		}
	}

	seedStr := r.URL.Query().Get("seed")
	if seedStr == "" {
		fmt.Println("Seed is not specified, using default value")
		seed = 20
	} else {
		seed, err = strconv.Atoi(seedStr)
		if err != nil {
			return 0, 0, fmt.Errorf("seed is not valid: %w", err)
		}
		if seed < 0 {
			return 0, 0, fmt.Errorf("seed is negative: %w", err)
		}
	}
	return skip, seed, nil
}
