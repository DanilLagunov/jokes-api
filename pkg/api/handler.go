package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/DanilLagunov/jokes-api/pkg/models"
	"github.com/DanilLagunov/jokes-api/pkg/storage"
	"github.com/gorilla/mux"
)

type Handler struct {
	Router   *mux.Router
	storage  storage.Storage
	template models.Template
}

func NewHandler(s storage.Storage) *Handler {
	h := &Handler{
		storage: s,
	}
	h.Router = h.initRoutes()
	h.template = models.NewTemptale()
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	start := time.Now()

	h.Router.ServeHTTP(w, req)

	// measure time
	fmt.Printf("request time is %v \n", time.Now().Sub(start))
}

func (h Handler) getPaginationParams(r *http.Request) (int, int, error) {
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
