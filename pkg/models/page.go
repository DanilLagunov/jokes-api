package models

type Page struct {
	Skip     int
	Seed     int
	CurrPage int
	MaxPage  int
	Content  []Joke
	Next     int
	Prev     int
}

func NewPage(skip, seed int, content []Joke) Page {
	if skip > len(content) || seed == 0 {
		return Page{skip, seed, 0, 0, []Joke{}, 0, 0}
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
		return Page{skip, seed, currPage, maxPage, content[skip:], next, prev}
	}

	return Page{skip, seed, currPage, maxPage, content[skip : skip+seed], next, prev}
}
