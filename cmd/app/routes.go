package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.ServeFiles("/static/*filepath", http.Dir("./static"))

	router.HandlerFunc(http.MethodGet, "/search_game", app.authenticate(http.HandlerFunc(app.searchGamesHandler)))
	router.HandlerFunc(http.MethodPost, "/add_game", app.authenticate(http.HandlerFunc(app.addGameHandler)))

	router.HandlerFunc(http.MethodGet, "/login", app.indexLoginHandler)
	router.HandlerFunc(http.MethodPost, "/login", app.loginUserHandler)
	router.HandlerFunc(http.MethodGet, "/signup", app.indexSignUpHandler)
	router.HandlerFunc(http.MethodPost, "/signup", app.createUserHandler)
	router.HandlerFunc(http.MethodPost, "/logout", app.authenticate(http.HandlerFunc(app.logoutUserHandler)))

	return router
}
