package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"toy-rental-system/internal/api/handler"
)

// Update the routes() method to return a http.Handler instead of a *httprouter.Router.
func (app *application) routes() http.Handler {
	//router instance
	router := httprouter.New()

	//convert our own helpers to http.Handler 404 code error using adapter
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	// likewise, convert to 405 error, basically making custom which is supported by http.Handler
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodPost, "/subscribe", handler.CreateSubscription)

	return router
}
