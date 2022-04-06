package api

import (
	"net/http"
	"time"

	"github.com/DanilLagunov/jokes-api/pkg/cache"
	"github.com/DanilLagunov/jokes-api/pkg/logger"
	"github.com/DanilLagunov/jokes-api/pkg/storage"
	"github.com/DanilLagunov/jokes-api/pkg/views"
	"github.com/gorilla/mux"
)

// Handler struct.
type Handler struct {
	Router   *mux.Router
	storage  storage.Storage
	template views.Template
	cache    cache.Cache
	logger   logger.Logger
}

// NewHandler creating a new Handler object.
func NewHandler(s storage.Storage, t views.Template, c cache.Cache, l logger.Logger) *Handler {
	h := &Handler{
		storage:  s,
		template: t,
		cache:    c,
		logger:   l,
	}
	h.Router = h.initRoutes()
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	start := time.Now()

	h.Router.ServeHTTP(w, req)

	// measure time
	h.logger.Log.Infof("request time is %v \n", time.Since(start))
}
