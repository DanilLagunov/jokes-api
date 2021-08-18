package models

import (
	"fmt"
	"html/template"
)

type Template struct {
	Template                 *template.Template
	GetJokesTemplate         string
	GetJokesByTextTemplate   string
	GetJokeByIDTemplate      string
	GetRandomJokesTemplate   string
	GetFunniestJokesTemplate string
}

func NewTemptale() Template {
	var t Template
	template, err := template.ParseFiles(
		"../templates/index.html",
		"../templates/get-joke-by-id.html",
		"../templates/get-jokes-by-text.html",
		"../templates/random.html",
		"../templates/funniest.html",
		"../templates/header.html",
		"../templates/footer.html")
	if err != nil {
		fmt.Println("template parsing error: %w", err)
	}
	t.Template = template
	t.GetJokesTemplate = "index"
	t.GetJokeByIDTemplate = "get-joke-by-id"
	t.GetJokesByTextTemplate = "get-jokes-by-text"
	t.GetRandomJokesTemplate = "random"
	t.GetFunniestJokesTemplate = "funniest"
	return t
}
