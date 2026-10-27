package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/luism2302/RateLog/static"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.ServeFiles("/static/*filepath", http.FS(static.Css))

	router.Handler(http.MethodGet, "/", http.RedirectHandler("/search_game", http.StatusSeeOther))
	router.HandlerFunc(http.MethodGet, "/search_game", app.authenticate(http.HandlerFunc(app.searchGamesHandler)))
	router.HandlerFunc(http.MethodPost, "/my_games/:id", app.authenticate(http.HandlerFunc(app.addGameHandler)))
	router.HandlerFunc(http.MethodPost, "/my_games/:id/delete", app.authenticate(http.HandlerFunc(app.deleteGameForUserHandler)))
	router.HandlerFunc(http.MethodGet, "/my_games", app.authenticate(http.HandlerFunc(app.myGamesHandler)))
	router.HandlerFunc(http.MethodGet, "/my_games/:id/review", app.authenticate(http.HandlerFunc(app.reviewGamePage)))
	router.HandlerFunc(http.MethodPost, "/my_games/:id/review", app.authenticate(http.HandlerFunc(app.createReviewHandler)))
	router.HandlerFunc(http.MethodGet, "/my_reviews", app.authenticate(http.HandlerFunc(app.myReviewsPageHandler)))

	router.HandlerFunc(http.MethodGet, "/healthz", app.healthzHandler)
	router.HandlerFunc(http.MethodGet, "/login", app.indexLoginHandler)
	router.HandlerFunc(http.MethodPost, "/login", app.loginUserHandler)
	router.HandlerFunc(http.MethodGet, "/signup", app.indexSignUpHandler)
	router.HandlerFunc(http.MethodPost, "/signup", app.createUserHandler)
	router.HandlerFunc(http.MethodPost, "/logout", app.authenticate(http.HandlerFunc(app.logoutUserHandler)))

	antiCSRF := http.NewCrossOriginProtection()

	return app.recoverPanic(app.commonHeaders(antiCSRF.Handler(app.rateLimit(router))))
}
