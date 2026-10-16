package main

import (
	"errors"
	"log/slog"
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

	app.logger.Info("form", slog.Any("values", formValues))
	app.logger.Info("validator", slog.Any("errs", v.Errors))
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
}
