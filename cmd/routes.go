package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Update the routes() method to return a http.Handler instead of a *httprouter.Router.
func (app *application) routes() http.Handler {
	//router instance
	router := httprouter.New()

	//convert our own helpers to http.Handler 404 code error using adapter
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	toysHandler := *app.toyHandler
	// likewise, convert to 405 error, basically making custom which is supported by http.Handler
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodPost, "/subscribe", app.subscriptionHandler.Subscribe)
	router.HandlerFunc(http.MethodPost, "/toy", toysHandler.CreateToyHandler)
	router.HandlerFunc(http.MethodGet, "/toy/:id", toysHandler.ShowToyHandler)
	router.HandlerFunc(http.MethodGet, "/toys", toysHandler.ListToysHandler)
	router.HandlerFunc(http.MethodDelete, "/toy/:id", toysHandler.DeleteToyHandler)
	router.HandlerFunc(http.MethodPatch, "/toy/:id", toysHandler.UpdateToyHandler)

	return router

}
