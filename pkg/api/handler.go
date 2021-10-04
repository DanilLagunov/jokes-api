package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/DanilLagunov/jokes-api/pkg/storage"
	"github.com/DanilLagunov/jokes-api/pkg/views"
	"github.com/gorilla/mux"
)

// Handler struct.
type Handler struct {
	Router   *mux.Router
	storage  storage.Storage
	template views.Template
}

// NewHandler creating a new Handler object.
func NewHandler(s storage.Storage, t views.Template) *Handler {
	h := &Handler{
		storage:  s,
		template: t,
	}
	h.Router = h.initRoutes()
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	start := time.Now()

	h.Router.ServeHTTP(w, req)

	// measure time
	fmt.Printf("request time is %v \n", time.Since(start))
}
