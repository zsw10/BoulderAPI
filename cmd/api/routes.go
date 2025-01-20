package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.notAllowedResponse)

	router.HandlerFunc(http.MethodPost, "/account", app.createAccountHandler)
	router.Handler(http.MethodGet, "/metrics", expvar.Handler())

	return app.recoverPanic(app.rateLimit(router))
}
