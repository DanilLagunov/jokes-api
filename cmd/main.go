package main

import (
	"log"
	"net/http"
	"time"

	"github.com/DanilLagunov/jokes-api/pkg/api"
	"github.com/DanilLagunov/jokes-api/pkg/storage/mongodb"
	"github.com/DanilLagunov/jokes-api/pkg/views"
)

func main() {
	storage, err := mongodb.NewDatabase("mongodb+srv://m001-student:m001-mongodb-basics@sandbox.evatv.mongodb.net/jokes-api?retryWrites=true&w=majority")
	if err != nil {
		log.Fatal(err)
	}

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
