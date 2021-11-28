package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DanilLagunov/jokes-api/pkg/api"
	"github.com/DanilLagunov/jokes-api/pkg/config"
	"github.com/DanilLagunov/jokes-api/pkg/storage/mongodb"
	"github.com/DanilLagunov/jokes-api/pkg/views"
)

func main() {
	cfgPath, err := config.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.NewConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	storage, err := mongodb.NewDatabase(cfg.Database.URI, cfg.Database.DBName, cfg.Database.JokesCollectionName)
	if err != nil {
		log.Fatal(err)
	}

	template := views.NewTemptale("./templates/")

	server := http.Server{
		Addr:              cfg.Server.Port,
		Handler:           api.NewHandler(storage, template),
		ReadHeaderTimeout: cfg.Server.Timeout.ReadHeader,
		ReadTimeout:       cfg.Server.Timeout.Read,
		WriteTimeout:      cfg.Server.Timeout.Write,
	}

	fmt.Println(server.ReadHeaderTimeout)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
