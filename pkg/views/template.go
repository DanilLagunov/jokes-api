package views

import (
	"fmt"
	"html/template"
)

const GetJokesTemplate string = "index"
const GetJokesByTextTemplate string = "get-jokes-by-text"
const GetJokeByIDTemplate string = "get-joke-by-id"
const GetRandomJokesTemplate string = "random"
const GetFunniestJokesTemplate string = "funniest"

type Template struct {
	Template *template.Template
}

func NewTemptale() Template {
	var t Template
	template, err := template.ParseFiles(
		"./templates/index.html",
		"./templates/get-joke-by-id.html",
		"./templates/get-jokes-by-text.html",
		"./templates/random.html",
		"./templates/funniest.html",
		"./templates/header.html",
		"./templates/footer.html")
	if err != nil {
		fmt.Println("template parsing error: %w", err)
	}
	t.Template = template
	return t
}
