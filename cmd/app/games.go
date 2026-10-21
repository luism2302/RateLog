package main

import (
	"errors"
	"net/http"

	"github.com/luism2302/RateLog/internal/data"
	"github.com/luism2302/RateLog/internal/validator"
	"github.com/luism2302/RateLog/internal/views"
)

func (app *application) searchGamesHandler(responseWriter http.ResponseWriter, request *http.Request) {
	title := app.readTitleQueryParam(request)
	user, ok := app.getAuthenticatedUser(request)
	if !ok {
		app.renderTemplate(responseWriter, views.SearchGamePage([]data.Game{}, nil))
		return
	}

	if title == "" {
		app.renderTemplate(responseWriter, views.SearchGamePage([]data.Game{}, user))
		return
	}

	games, err := app.client.SearchGameByTitle(title)
	if err != nil {
		app.serverErrorResponse(responseWriter, request, err)
		return
	}

	if err := app.renderTemplate(responseWriter, views.SearchGamePage(games, user)); err != nil {
		app.serverErrorResponse(responseWriter, request, err)
		return
	}
}

func (app *application) addGameHandler(responseWriter http.ResponseWriter, request *http.Request) {
	v := validator.New()
	formValues, err := app.parseAndValidateAddGameForm(request, v)
	if err != nil {
		app.serverErrorResponse(responseWriter, request, err)
		return
	}

	if !v.Valid() {
		app.serverErrorResponse(responseWriter, request, errors.New("Unable to parse add game form"))
		return
	}

	game, err := data.NewGame(formValues.ID, formValues.Added, formValues.Name, formValues.Released, formValues.BackgroundImage)
	if err != nil {
		app.serverErrorResponse(responseWriter, request, err)
		return
	}

	user, ok := app.getAuthenticatedUser(request)
	if !ok {
		http.Redirect(responseWriter, request, "/login", http.StatusSeeOther)
		return
	}
	if err := app.models.Games.Insert(game); err != nil {
		app.serverErrorResponse(responseWriter, request, err)
		return
	}

	if err := app.models.UserGames.Insert(user.ID, game.ID); err != nil {
		app.serverErrorResponse(responseWriter, request, err)
		return
	}

	http.Redirect(responseWriter, request, "/search_game", http.StatusSeeOther)
}

func (app *application) myGamesHandler(responseWriter http.ResponseWriter, request *http.Request) {
	user, ok := app.getAuthenticatedUser(request)
	if !ok {
		http.Redirect(responseWriter, request, "/login", http.StatusSeeOther)
		return
	}

	games, err := app.models.Games.GetGamesForUser(user.ID)
	if err != nil {
		app.serverErrorResponse(responseWriter, request, err)
		return
	}

	app.renderTemplate(responseWriter, views.MyGamesPage(games, user))
}
