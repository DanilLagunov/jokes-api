package views

import (
	"fmt"
	"html/template"
	"path"
)

const GetJokesTemplate string = "index"
const GetJokesByTextTemplate string = "get-jokes-by-text"
const GetJokeByIDTemplate string = "get-joke-by-id"
const GetRandomJokesTemplate string = "random"
const GetFunniestJokesTemplate string = "funniest"

type Template struct {
	Template *template.Template
}

func NewTemptale(folder string) Template {
	var t Template
	template, err := template.ParseFiles(
		path.Join(folder, "index.html"),
		path.Join(folder, "get-joke-by-id.html"),
		path.Join(folder, "get-jokes-by-text.html"),
		path.Join(folder, "random.html"),
		path.Join(folder, "funniest.html"),
		path.Join(folder, "header.html"),
		path.Join(folder, "footer.html"))
	if err != nil {
		fmt.Println("template parsing error: %w", err)
	}
	t.Template = template
	return t
}
