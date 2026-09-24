package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.ServeFiles("/static/*filepath", http.Dir("./static"))

	router.HandlerFunc(http.MethodGet, "/", app.home)

	router.HandlerFunc(http.MethodGet, "/games/:title", app.searchGameHandler)

	return router
}
