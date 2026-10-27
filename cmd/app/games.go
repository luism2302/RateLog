package main

import (
	"errors"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/luism2302/RateLog/internal/client"
	"github.com/luism2302/RateLog/internal/data"
	"github.com/luism2302/RateLog/internal/validator"
	"github.com/luism2302/RateLog/internal/views"
)

func (app *application) searchGamesHandler(responseWriter http.ResponseWriter, request *http.Request) {
	v := validator.New()
	user, ok := app.getAuthenticatedUser(request)
	if !ok {
		app.renderTemplate(responseWriter, request, http.StatusOK, views.SearchGamePage([]data.Game{}, nil, v))
		return
	}

	title := app.readTitleQueryParam(request)
	title = strings.ToLower(strings.TrimSpace(title))

	if title == "" {
		app.renderTemplate(responseWriter, request, http.StatusOK, views.SearchGamePage([]data.Game{}, user, v))
		return
	}

	if utf8.RuneCountInString(title) > 100 {
		v.AddError("title", "Game title too long")
		app.renderTemplate(responseWriter, request, http.StatusUnprocessableEntity, views.SearchGamePage([]data.Game{}, user, v))
		return
	}

	games, err := app.client.SearchGameByTitle(user.ID, title, request.Context())
	if err != nil {
		switch {
		case errors.Is(err, client.ErrRateLimited):
			app.rateLimitExceededResponse(responseWriter, request)
			return
		case errors.Is(err, client.ErrUnavailable):
			app.serviceUnavailable(responseWriter, request)
			return
		default:
			app.serverErrorResponse(responseWriter, request, err)
			return
		}
	}

	app.renderTemplate(responseWriter, request, http.StatusOK, views.SearchGamePage(games, user, v))
}

func (app *application) addGameHandler(responseWriter http.ResponseWriter, request *http.Request) {
	user, ok := app.getAuthenticatedUser(request)
	if !ok {
		http.Redirect(responseWriter, request, "/login", http.StatusSeeOther)
		return
	}

	gameID, err := app.readGameIDParam(request)
	if err != nil {
		app.notFoundResponse(responseWriter, request)
		return
	}

	game, err := app.models.Games.GetGameByID(gameID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			game, err = app.client.SearchGameByID(user.ID, gameID, request.Context())
			if err != nil {
				switch {
				case errors.Is(err, client.ErrNotFound):
					app.notFoundResponse(responseWriter, request)
					return
				case errors.Is(err, client.ErrUnavailable) || errors.Is(err, client.ErrRateLimited):
					app.serviceUnavailable(responseWriter, request)
					return
				default:
					app.serverErrorResponse(responseWriter, request, err)
					return
				}
			}
			if err := app.models.Games.Insert(game); err != nil {
				app.serverErrorResponse(responseWriter, request, err)
				return
			}
		default:
			app.serverErrorResponse(responseWriter, request, err)
			return
		}
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

	app.renderTemplate(responseWriter, request, http.StatusOK, views.MyGamesPage(games, user))
}

func (app *application) deleteGameForUserHandler(responseWriter http.ResponseWriter, request *http.Request) {
	user, ok := app.getAuthenticatedUser(request)
	if !ok {
		http.Redirect(responseWriter, request, "/login", http.StatusSeeOther)
		return
	}

	gameID, err := app.readGameIDParam(request)
	if err != nil {
		app.notFoundResponse(responseWriter, request)
		return
	}

	if err := app.models.UserGames.Delete(user.ID, gameID); err != nil {
		switch {
		case errors.Is(err, data.ErrNoRowsAffected):
			app.notFoundResponse(responseWriter, request)
			return
		default:
			app.serverErrorResponse(responseWriter, request, err)
			return
		}
	}
	http.Redirect(responseWriter, request, "/my_games", http.StatusSeeOther)
}
