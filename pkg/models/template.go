package models

import (
	"fmt"
	"html/template"
)

type Template struct {
	Template                 *template.Template
	GetJokesTemplate         string
	GetJokeByParamTemplate   string
	GetRandomJokesTemplate   string
	GetFunniestJokesTemplate string
}

func NewTemptale() Template {
	var t Template
	template, err := template.ParseFiles(
		"../templates/index.html",
		"../templates/findjoke.html",
		"../templates/random.html",
		"../templates/funniest.html",
		"../templates/header.html",
		"../templates/footer.html")
	if err != nil {
		fmt.Println("template initializing error: %w", err)
	}
	t.Template = template
	t.GetJokesTemplate = "index"
	t.GetJokeByParamTemplate = "findjoke"
	t.GetRandomJokesTemplate = "random"
	t.GetFunniestJokesTemplate = "funniest"
	return t
}
