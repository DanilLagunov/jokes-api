package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/DanilLagunov/jokes-api/handlers"
	"github.com/DanilLagunov/jokes-api/models"
	"github.com/DanilLagunov/jokes-api/server/router"
	"github.com/DanilLagunov/jokes-api/utils"
)

type Server struct {
	routes router.Router
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start() {

	var jokes []models.Joke
	utils.ParseJSON("reddit_jokes.json", &jokes)

	template, err := template.ParseFiles("templates/index.html", "templates/findjoke.html", "templates/random.html", "templates/funniest.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Println("Template parsing error!")
	}

	jokeHandler := handlers.NewJokeHandler(jokes, template)

	s.routes = router.NewRouter()
	s.routes.Route.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("././assets/"))))
	s.routes.CreateRoutes(jokeHandler)

	log.Fatal(http.ListenAndServe(":8000", s.routes.Route))
}
