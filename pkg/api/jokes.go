package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/DanilLagunov/jokes-api/pkg/storage"
	"github.com/DanilLagunov/jokes-api/pkg/views"
)

const requestTimeout time.Duration = time.Second * 2

func (h Handler) getJokes(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	skip, limit, err := getPaginationParams(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	jokes, amount, err := h.storage.GetJokes(ctx, skip, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("response writing error: %s", err)
		}
	}

	pageParams := views.CreatePageParams(skip, limit, amount, jokes)

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

	_, err := h.storage.AddJoke(ctx, title, body, 0)
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

	skip, limit, err := getPaginationParams(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	if text == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, amount, err := h.storage.GetJokesByText(ctx, skip, limit, text)
	if errors.Is(err, storage.ErrJokeNotFound) {
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

	pageParams := views.CreatePageParams(skip, limit, amount, result)

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

	result, err := h.cache.Get(id)
	if err == nil {
		err = h.template.Template.ExecuteTemplate(w, views.GetJokeByIDTemplate, result)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	fmt.Printf("cache error: %s", err)

	result, err = h.storage.GetJokeByID(ctx, id)
	if errors.Is(err, storage.ErrJokeNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("response writing error: %s", err)
		}

		return
	}

	h.cache.Set(id, result, 20*time.Second)

	err = h.template.Template.ExecuteTemplate(w, views.GetJokeByIDTemplate, result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h Handler) getRandomJokes(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	skip, limit, err := getPaginationParams(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	random, amount, err := h.storage.GetRandomJokes(ctx, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("response writing error: %s", err)
		}
	}

	pageParams := views.CreatePageParams(skip, limit, amount, random)

	err = h.template.Template.ExecuteTemplate(w, views.GetRandomJokesTemplate, pageParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h Handler) getFunniestJokes(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	skip, limit, err := getPaginationParams(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	funniest, amount, err := h.storage.GetFunniestJokes(ctx, skip, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("response writing error: %s", err)
		}
	}

	pageParams := views.CreatePageParams(skip, limit, amount, funniest)

	err = h.template.Template.ExecuteTemplate(w, views.GetFunniestJokesTemplate, pageParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getPaginationParams(r *http.Request) (int, int, error) {
	var skip, limit int

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

	limitStr := r.URL.Query().Get("seed")
	if limitStr == "" {
		fmt.Println("Seed is not specified, using default value")

		limit = 20
	} else {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return 0, 0, fmt.Errorf("seed is not valid: %w", err)
		}
		if limit < 0 {
			return 0, 0, fmt.Errorf("seed is negative: %w", err)
		}
	}
	return skip, limit, nil
}
