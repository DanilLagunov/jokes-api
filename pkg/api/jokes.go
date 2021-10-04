package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	file_storage "github.com/DanilLagunov/jokes-api/pkg/storage/file-storage"
	"github.com/DanilLagunov/jokes-api/pkg/views"
)

const requestTimeout time.Duration = time.Second * 2

func (h Handler) getJokes(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	skip, seed, err := getPaginationParams(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	jokes, err := h.storage.GetJokes(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("response writing error: %s", err)
		}
	}

	pageParams := views.CreatePageParams(skip, seed, jokes)

	err = h.template.Template.ExecuteTemplate(w, views.GetJokesTemplate, pageParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h Handler) addJoke(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	title := r.FormValue("title")
	body := r.FormValue("body")

	_, err := h.storage.AddJoke(ctx, title, body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("response writing error: %s", err)
		}

		return
	}

	http.Redirect(w, r, "/jokes", http.StatusFound)
}

func (h Handler) getJokesByText(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	text := r.URL.Query().Get("text")

	skip, seed, err := getPaginationParams(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	if text == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := h.storage.GetJokeByText(ctx, text)
	if errors.Is(err, file_storage.ErrJokeNotFound) {
		w.WriteHeader(http.StatusNotFound)

		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("response writing error: %s", err)
		}

		return
	}

	pageParams := views.CreatePageParams(skip, seed, result)

	err = h.template.Template.ExecuteTemplate(w, views.GetJokesByTextTemplate,
		views.SearchPageParams{
			SearchRequest: text,
			PageParams:    pageParams,
		})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h Handler) getJokeByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	id := r.URL.Query().Get("id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := h.storage.GetJokeByID(ctx, id)
	if errors.Is(err, file_storage.ErrJokeNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("response writing error: %s", err)
		}

		return
	}

	err = h.template.Template.ExecuteTemplate(w, views.GetJokeByIDTemplate, result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	return
}

func (h Handler) getRandomJokes(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	skip, seed, err := getPaginationParams(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	random, err := h.storage.GetRandomJokes(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("response writing error: %s", err)
		}
	}

	pageParams := views.CreatePageParams(skip, seed, random)

	err = h.template.Template.ExecuteTemplate(w, views.GetRandomJokesTemplate, pageParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h Handler) getFunniestJokes(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	skip, seed, err := getPaginationParams(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	funniest, err := h.storage.GetFunniestJokes(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("response writing error: %s", err)
		}
	}

	pageParams := views.CreatePageParams(skip, seed, funniest)

	err = h.template.Template.ExecuteTemplate(w, views.GetFunniestJokesTemplate, pageParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getPaginationParams(r *http.Request) (int, int, error) {
	var skip, seed int

	var err error

	skipStr := r.URL.Query().Get("skip")
	if skipStr == "" {
		fmt.Println("Skip is not specified, using default value")

		skip = 0
	} else {
		skip, err = strconv.Atoi(skipStr)
		if err != nil {
			return 0, 0, fmt.Errorf("skip is not valid: %w", err)
		}
		if skip < 0 {
			return 0, 0, fmt.Errorf("skip is negative: %w", err)
		}
	}

	seedStr := r.URL.Query().Get("seed")
	if seedStr == "" {
		fmt.Println("Seed is not specified, using default value")

		seed = 20
	} else {
		seed, err = strconv.Atoi(seedStr)
		if err != nil {
			return 0, 0, fmt.Errorf("seed is not valid: %w", err)
		}
		if seed < 0 {
			return 0, 0, fmt.Errorf("seed is negative: %w", err)
		}
	}
	return skip, seed, nil
}
