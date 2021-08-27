package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (h Handler) initRoutes() *mux.Router {
	router := mux.NewRouter()
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	router.HandleFunc("/jokes", h.getJokes).Methods(http.MethodGet)
	router.HandleFunc("/jokes/add", h.addJoke).Methods(http.MethodPost)
	router.HandleFunc("/jokes/search", h.getJoke).Methods(http.MethodGet)
	router.HandleFunc("/jokes/random", h.getRandomJokes).Methods(http.MethodGet)
	router.HandleFunc("/jokes/funniest", h.getFunniestJokes).Methods(http.MethodGet)

	h.Router = router
	return router
}
