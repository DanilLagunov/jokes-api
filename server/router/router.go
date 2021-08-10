package router

import (
	"github.com/DanilLagunov/jokes-api/handlers"
	"github.com/gorilla/mux"
)

type Router struct {
	Route *mux.Router
}

func NewRouter() Router {
	return Router{Route: mux.NewRouter()}
}

func (r *Router) CreateRoutes(jokeHandler handlers.JokeHandler) {

	r.Route.HandleFunc("/jokes", jokeHandler.GetJokes).Methods("GET")
	r.Route.HandleFunc("/jokes/add", jokeHandler.AddJoke).Methods("POST")
	r.Route.HandleFunc("/jokes/search", jokeHandler.GetJoke).Methods("GET")
	r.Route.HandleFunc("/jokes/random", jokeHandler.GetRandomJokes).Methods("GET")
	r.Route.HandleFunc("/jokes/funniest", jokeHandler.GetFunniestJokes).Methods("GET")

}
