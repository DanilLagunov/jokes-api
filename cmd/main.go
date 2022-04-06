package main

import (
	"net/http"
	"strconv"

	"github.com/DanilLagunov/jokes-api/pkg/api"
	"github.com/DanilLagunov/jokes-api/pkg/cache/memcache"
	"github.com/DanilLagunov/jokes-api/pkg/config"
	"github.com/DanilLagunov/jokes-api/pkg/logger"
	"github.com/DanilLagunov/jokes-api/pkg/storage/mongodb"
	"github.com/DanilLagunov/jokes-api/pkg/views"
)

func main() {
	logger := logger.InitLogger()

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Log.Fatalf("config creating error: %s", err)
	}

	storage, err := mongodb.NewDatabase(cfg.DbURI, cfg.DbName, cfg.JokesCollection)
	if err != nil {
		logger.Log.Fatalf("database creating error: %s", err)
	}

	template := views.NewTemptale("./templates/")

	cache := memcache.NewMemCache(cfg.CacheDefaultExpiration, cfg.CacheCleanupInterval)

	server := http.Server{
		Addr:              ":" + strconv.Itoa(cfg.Port),
		Handler:           api.NewHandler(storage, template, cache, *logger),
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
	}

	logger.Log.Infof("Server is listening on port %d", cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		logger.Log.Fatalf("server error: %s", err)
	}
}
