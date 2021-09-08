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

func NewTemptale(path [7]string) Template {
	var t Template
	template, err := template.ParseFiles(
		path[0], path[1], path[2], path[3], path[4], path[5], path[6])
	if err != nil {
		fmt.Println("template parsing error: %w", err)
	}
	t.Template = template
	return t
}
