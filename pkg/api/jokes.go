package api

import (
	"errors"
	"log"
	"net/http"

	"github.com/DanilLagunov/jokes-api/pkg/models"
	file_storage "github.com/DanilLagunov/jokes-api/pkg/storage/file-storage"

	"github.com/DanilLagunov/jokes-api/pkg/views"
)

func (h Handler) getJokes(w http.ResponseWriter, r *http.Request) {
	skip, seed, err := views.GetPaginationParams(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	jokes, err := h.storage.GetJokes()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	paginationData := views.CreatePaginationData(skip, seed, jokes)

	err = h.template.Template.ExecuteTemplate(w, views.GetJokesTemplate, paginationData)
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
	skip, seed, err := views.GetPaginationParams(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	if text != "" {
		result, err := h.storage.GetJokeByText(text)
		if err != nil {
			if errors.Is(err, file_storage.JokeNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		paginationData := views.CreatePaginationData(skip, seed, result)
		err = h.template.Template.ExecuteTemplate(w, views.GetJokesByTextTemplate,
			views.SearchPaginationData{
				SearchRequest:  text,
				PaginationData: paginationData,
			})
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if id != "" {
		result, err := h.storage.GetJokeByID(id)
		if err != nil {
			if errors.Is(err, file_storage.JokeNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = h.template.Template.ExecuteTemplate(w, views.GetJokeByIDTemplate, result)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	err = h.template.Template.ExecuteTemplate(w, views.GetJokesByTextTemplate, []models.Joke{})
	if err != nil {
		log.Fatal(err)
	}
	return

}

func (h Handler) getRandomJokes(w http.ResponseWriter, r *http.Request) {
	skip, seed, err := views.GetPaginationParams(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	random, err := h.storage.GetRandomJokes()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	paginationData := views.CreatePaginationData(skip, seed, random)

	err = h.template.Template.ExecuteTemplate(w, views.GetRandomJokesTemplate, paginationData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h Handler) getFunniestJokes(w http.ResponseWriter, r *http.Request) {

	skip, seed, err := views.GetPaginationParams(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	funniest, err := h.storage.GetFunniestJokes()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	paginationData := views.CreatePaginationData(skip, seed, funniest)

	err = h.template.Template.ExecuteTemplate(w, views.GetFunniestJokesTemplate, paginationData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
