package main

import (
	"log"
	"net/http"
	"time"

	"github.com/DanilLagunov/jokes-api/pkg/api"
	file_storage "github.com/DanilLagunov/jokes-api/pkg/storage/file-storage"
	"github.com/DanilLagunov/jokes-api/pkg/views"
)

func main() {
	storage := file_storage.NewFileStorage("./pkg/storage/file-storage/reddit_jokes.json")
	template := views.NewTemptale("./templates/")
	server := http.Server{
		Addr:              ":8000",
		Handler:           api.NewHandler(storage, template),
		ReadHeaderTimeout: time.Second * 30,
		ReadTimeout:       time.Second * 60,
		WriteTimeout:      time.Second * 60,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
