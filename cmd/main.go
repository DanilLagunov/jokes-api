package main

import (
	"log"
	"net/http"
	"time"

	"github.com/DanilLagunov/jokes-api/pkg/api"
	"github.com/DanilLagunov/jokes-api/pkg/models"
	memory_storage "github.com/DanilLagunov/jokes-api/pkg/storage/memory-storage"
)

func main() {
	storage := memory_storage.NewMemoryStorage()
	template := models.NewTemptale()
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
