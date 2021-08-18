package api

import (
	"errors"
	"net/http"

	"github.com/DanilLagunov/jokes-api/pkg/models"
	memory_storage "github.com/DanilLagunov/jokes-api/pkg/storage/memory-storage"
)

func (h Handler) getJokes(w http.ResponseWriter, r *http.Request) {
	skip, seed, err := h.getPaginationParams(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	jokes, err := h.storage.GetJokes()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	page := models.NewPage(skip, seed, jokes)

	err = h.template.Template.ExecuteTemplate(w, h.template.GetJokesTemplate, page)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h Handler) addJoke(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	body := r.FormValue("body")

	h.storage.AddJoke(title, body)

	http.Redirect(w, r, "/jokes", 302)
}

func (h Handler) getJoke(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")
	id := r.URL.Query().Get("id")

	if text != "" {
		result, err := h.storage.GetJokeByText(text)
		if err != nil {
			if errors.Is(err, memory_storage.JokeNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = h.template.Template.ExecuteTemplate(w, h.template.GetJokesTemplate, result)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if id != "" {
		result, err := h.storage.GetJokeByID(id)
		if err != nil {
			if errors.Is(err, memory_storage.JokeNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = h.template.Template.ExecuteTemplate(w, h.template.GetJokesTemplate, result)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	err := h.template.Template.ExecuteTemplate(w, h.template.GetJokesTemplate, models.Joke{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	return

}

func (h Handler) getRandomJokes(w http.ResponseWriter, r *http.Request) {
	skip, seed, err := h.getPaginationParams(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	random, err := h.storage.GetRandomJokes()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	page := models.NewPage(skip, seed, random)

	err = h.template.Template.ExecuteTemplate(w, h.template.GetJokesTemplate, page)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h Handler) getFunniestJokes(w http.ResponseWriter, r *http.Request) {

	skip, seed, err := h.getPaginationParams(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	funniest, err := h.storage.GetFunniestJokes()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	page := models.NewPage(skip, seed, funniest)

	err = h.template.Template.ExecuteTemplate(w, h.template.GetJokesTemplate, page)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
