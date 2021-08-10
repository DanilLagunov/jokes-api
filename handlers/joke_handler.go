package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/DanilLagunov/jokes-api/models"
	"github.com/DanilLagunov/jokes-api/utils"
)

type JokeHandler struct {
	jokes    []models.Joke
	template *template.Template
}

func NewJokeHandler(jokes []models.Joke, template *template.Template) JokeHandler {
	return JokeHandler{jokes, template}
}

func (h *JokeHandler) GetJokes(w http.ResponseWriter, r *http.Request) {
	skip, err := strconv.Atoi(r.URL.Query().Get("skip"))
	if err != nil {
		fmt.Println("Converting err, using default value")
		skip = 0
	}

	seed, err := strconv.Atoi(r.URL.Query().Get("seed"))
	if err != nil {
		fmt.Println("Converting err, using default value")
		seed = 10
	}

	page := models.NewPage(skip, seed, skip+seed, skip-seed, h.jokes[skip:skip+seed])

	err = h.template.ExecuteTemplate(w, "index", page)
	if err != nil {
		log.Fatal(err)
	}

}

func (h *JokeHandler) AddJoke(w http.ResponseWriter, r *http.Request) {
	// Generating ID and taking data from form
	title := r.FormValue("title")
	body := r.FormValue("body")
	score := 0

	// Check for uniqueness
CHECK:
	id := utils.GenerateId()
	for i := 0; i < len(h.jokes); i++ {
		if id == h.jokes[i].ID {
			goto CHECK
		}
	}

	// Creating new joke
	joke := models.NewJoke(id, title, body, score)
	h.jokes = append(h.jokes, joke)

	//Formatting and adding joke to the file
	rawDataOut, err := json.MarshalIndent(&h.jokes, "", "   ")
	if err != nil {
		log.Fatal("JSON marshaling failed: ", err)
	}
	err = ioutil.WriteFile("reddit_jokes.json", rawDataOut, 0)
	if err != nil {
		log.Fatal("Cannot write:", err)
	}

	//Redirecting to main page
	http.Redirect(w, r, "/jokes", 302)
}

func (h *JokeHandler) GetJoke(w http.ResponseWriter, r *http.Request) {
	//Taking search values from URL
	text := r.URL.Query().Get("text")
	id := r.URL.Query().Get("id")

	var result []models.Joke

	//Searching by text
	if text != "" {
		for _, item := range h.jokes {
			if strings.Contains(item.Title, text) || strings.Contains(item.Body, text) {
				result = append(result, item)
			}
		}
		err := h.template.ExecuteTemplate(w, "findjoke", result)
		if err != nil {
			log.Fatal(err)
		}
	}

	//Searching by ID
	if id != "" {
		for _, item := range h.jokes {
			if item.ID == id {
				result = append(result, item)
				err := h.template.ExecuteTemplate(w, "findjoke", result)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

func (h *JokeHandler) GetRandomJokes(w http.ResponseWriter, r *http.Request) {
	s := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(s)

	skip, err := strconv.Atoi(r.URL.Query().Get("skip"))
	if err != nil {
		fmt.Println("Converting err, using default value")
		skip = 0
	}

	seed, err := strconv.Atoi(r.URL.Query().Get("seed"))
	if err != nil {
		fmt.Println("Converting err, using default value")
		seed = 10
	}

	//Filling an array with random elements
	var rndjokes [300]models.Joke
	for i := 0; i < 300; i++ {
		rndjokes[i] = h.jokes[rnd.Intn(len(h.jokes))]
	}

	page := models.NewPage(skip, seed, skip+seed, skip-seed, h.jokes[skip:skip+seed])

	err = h.template.ExecuteTemplate(w, "random", page)
	if err != nil {
		log.Fatal(err)
	}
}

func (h *JokeHandler) GetFunniestJokes(w http.ResponseWriter, r *http.Request) {
	//Sorting an array by score
	skip, err := strconv.Atoi(r.URL.Query().Get("skip"))
	if err != nil {
		fmt.Println("Converting err, using default value")
		skip = 0
	}

	seed, err := strconv.Atoi(r.URL.Query().Get("seed"))
	if err != nil {
		fmt.Println("Converting err, using default value")
		seed = 10
	}

	var funniest []models.Joke

	funniest = append(funniest, h.jokes...)
	sort.Slice(funniest, func(i, j int) (less bool) {
		return funniest[i].Score > funniest[j].Score
	})

	page := models.NewPage(skip, seed, skip+seed, skip-seed, h.jokes[skip:skip+seed])

	err = h.template.ExecuteTemplate(w, "funniest", page)
	if err != nil {
		log.Fatal(err)
	}
}
