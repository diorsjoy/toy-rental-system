package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *Application) routes() *httprouter.Router {
	// router instance
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/toys", app.)

	// convert our own helpers to http.Handler 404 code error using adapter SDP beyba xD

	// return instance
	return router
}
