package views

import (
	"fmt"
	"html/template"
	"path"
)

// GetJokesTemplate is a constant for calling the "index" template.
const GetJokesTemplate string = "index"

// GetJokesByTextTemplate is a constant for calling the "get-jokes-by-text" template.
const GetJokesByTextTemplate string = "get-jokes-by-text"

// GetJokeByIDTemplate is a constant for calling the "get-joke-by-id" template.
const GetJokeByIDTemplate string = "get-joke-by-id"

// GetRandomJokesTemplate is a constant for calling the "random" template.
const GetRandomJokesTemplate string = "random"

// GetFunniestJokesTemplate is a constant for calling the "funniest" template.
const GetFunniestJokesTemplate string = "funniest"

// Template struct.
type Template struct {
	Template *template.Template
}

// NewTemptale creating a new Template object.
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
