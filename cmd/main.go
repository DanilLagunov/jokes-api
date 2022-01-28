package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/DanilLagunov/jokes-api/pkg/api"
	"github.com/DanilLagunov/jokes-api/pkg/config"
	"github.com/DanilLagunov/jokes-api/pkg/storage/mongodb"
	"github.com/DanilLagunov/jokes-api/pkg/views"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	storage, err := mongodb.NewDatabase(cfg.DbURI, cfg.DbName, cfg.JokesCollection)
	if err != nil {
		log.Fatal(err)
	}

	template := views.NewTemptale("./templates/")

	server := http.Server{
		Addr:              ":" + strconv.Itoa(cfg.Port),
		Handler:           api.NewHandler(storage, template),
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
