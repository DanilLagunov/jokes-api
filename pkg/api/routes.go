package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (h Handler) initRoutes() *mux.Router {
	h.Router = mux.NewRouter()
	h.Router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	h.Router.HandleFunc("/jokes", h.getJokes).Methods(http.MethodGet)
	h.Router.HandleFunc("/jokes/add", h.addJoke).Methods(http.MethodPost)
	h.Router.HandleFunc("/jokes/search?id={id}", h.getJokeByID).Methods(http.MethodGet)
	h.Router.HandleFunc("/jokes/search?text={text}", h.getJokesByText).Methods(http.MethodGet)
	h.Router.HandleFunc("/jokes/random", h.getRandomJokes).Methods(http.MethodGet)
	h.Router.HandleFunc("/jokes/funniest", h.getFunniestJokes).Methods(http.MethodGet)

	return h.Router
}
