package main

import (
	"database/sql"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"toy-rental-system/internal/data"
	"toy-rental-system/service"
)

// Update the routes() method to return a http.Handler instead of a *httprouter.Router.
func (app *application) routes2() http.Handler {
	//router instance
	router := httprouter.New()

	//convert our own helpers to http.Handler 404 code error using adapter
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	// likewise, convert to 405 error, basically making custom which is supported by http.Handler
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodPost, "/subscribe", app.subscriptionHandler.Subscribe)

}

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


	return router
}
