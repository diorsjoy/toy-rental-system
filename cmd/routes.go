package main

import (
	"github.com/julienschmidt/httprouter"
)

func routes() *httprouter.Router {
	// router instance
	router := httprouter.New()

	// convert our own helpers to http.Handler 404 code error using adapter SDP beyba xD

	// return instance
	return router
}
