package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// Joke struct with constructor

type Joke struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
	Score int    `json:"score"`
}

func NewJoke(id, title, body string, score int) Joke {
	return Joke{id, title, body, score}
}

//GLOBAL VARIABLES

var jokes []Joke

var t *template.Template

// HANDLERS

func getJokes(w http.ResponseWriter, r *http.Request) {
	skip, err := strconv.Atoi(r.URL.Query().Get("skip"))
	if err != nil {
		log.Fatal(err)
	}

	seed, err := strconv.Atoi(r.URL.Query().Get("seed"))
	if err != nil {
		log.Fatal(err)
	}

	t.ExecuteTemplate(w, "index", jokes[skip:skip+seed])

}

func addJoke(w http.ResponseWriter, r *http.Request) {
	// Generating ID and taking data from form
	var id string
	title := r.FormValue("title")
	body := r.FormValue("body")
	score := 0

	// Check for uniqueness
CHECK:
	id = generateId()
	for i := 0; i < len(jokes); i++ {
		if id == jokes[i].ID {
			goto CHECK
		}
	}

	// Creating new joke
	joke := NewJoke(id, title, body, score)
	jokes = append(jokes, joke)

	//Formatting and adding joke to the file
	rawDataOut, err := json.MarshalIndent(&jokes, "", "   ")
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

func getJoke(w http.ResponseWriter, r *http.Request) {
	//Taking search values from URL
	text := r.URL.Query().Get("text")
	id := r.URL.Query().Get("id")

	var result []Joke

	//Searching by text
	if text != "" {
		for _, item := range jokes {
			if strings.Contains(item.Title, text) || strings.Contains(item.Body, text) {
				result = append(result, item)
			}
		}
		t.ExecuteTemplate(w, "findjoke", result)
	}

	//Searching by ID
	if id != "" {
		for _, item := range jokes {
			if item.ID == id {
				result = append(result, item)
				t.ExecuteTemplate(w, "findjoke", result)
			}
		}
	}
}

func getRandomJokes(w http.ResponseWriter, r *http.Request) {
	s := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(s)

	//Filling an array with random elements
	var rndjokes [100]Joke
	for i := 0; i < 100; i++ {
		rndjokes[i] = jokes[rnd.Intn(len(jokes))]
	}

	t.ExecuteTemplate(w, "random", rndjokes)
}

func getFunniestJokes(w http.ResponseWriter, r *http.Request) {
	//Sorting an array by score
	var funniest []Joke
	funniest = append(funniest, jokes...)
	sort.Slice(funniest, func(i, j int) (less bool) {
		return funniest[i].Score > funniest[j].Score
	})

	t.ExecuteTemplate(w, "funniest", funniest[0:99])
}

// ADDITIONAL FUNCS

func parseJSON(path string, list *[]Joke) {
	file, _ := os.Open(path)
	decoder := json.NewDecoder(file)

	err := decoder.Decode(&list)
	if err != nil {
		fmt.Println("Decode error")
	}
}

func generateId() string {
	b := make([]byte, 3)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// MAIN FUNC

func main() {
	parseJSON("reddit_jokes.json", &jokes)

	r := mux.NewRouter()

	t, _ = template.ParseFiles("templates/index.html", "templates/findjoke.html", "templates/random.html", "templates/funniest.html", "templates/header.html", "templates/footer.html")

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	r.HandleFunc("/jokes", getJokes).Methods("GET")
	r.HandleFunc("/jokes/add", addJoke).Methods("POST")
	r.HandleFunc("/jokes/search", getJoke).Methods("GET")
	r.HandleFunc("/jokes/random", getRandomJokes).Methods("GET")
	r.HandleFunc("/jokes/funniest", getFunniestJokes).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", r))
}
