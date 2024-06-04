package main

import (
	"database/sql"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"toy-rental-system/internal/data"
	"toy-rental-system/service"
)

func routes() *httprouter.Router {
	// router instance
	router := httprouter.New()

	db, err := sql.Open("postgres", "your_connection_string_here")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	toyRepo := data.ToyModel{DB: db}
	toyService := service.NewToyService(toyRepo)
	// convert our own helpers to http.Handler 404 code error using adapter SDP beyba xD
	http.HandleFunc("/toys", toyService.ListToysHandler)
	http.HandleFunc("/toys/create", toyService.CreateToyHandler)
	http.HandleFunc("/toys/update", toyService.UpdateToyHandler)
	http.HandleFunc("/toys/delete", toyService.DeleteToyHandler)
	http.HandleFunc("/toys/show", toyService.ShowToyHandler)

	// return instance
	return router
}
