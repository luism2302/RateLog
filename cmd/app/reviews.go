package main

import (
	"errors"
	"net/http"

	"github.com/luism2302/RateLog/internal/data"
	"github.com/luism2302/RateLog/internal/validator"
	"github.com/luism2302/RateLog/internal/views"
)

func (app *application) reviewGamePage(responseWriter http.ResponseWriter, request *http.Request) {
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
			app.notFoundResponse(responseWriter, request)
			return
		default:
			app.serverErrorResponse(responseWriter, request, err)
			return
		}
	}
	if err := app.models.UserGames.HasGame(user.ID, gameID); err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(responseWriter, request)
			return
		default:
			app.serverErrorResponse(responseWriter, request, err)
			return
		}
	}

	v := validator.New()

	review, err := app.models.Reviews.Get(user.ID, gameID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.renderTemplate(responseWriter, request, http.StatusOK, views.ReviewPage(user, &game, views.ReviewFormValues{}, false, v))
			return

		default:
			app.serverErrorResponse(responseWriter, request, err)
			return
		}
	}

	formValues := views.ReviewFormValues{Body: review.Body, Rating: review.Rating}
	app.renderTemplate(responseWriter, request, http.StatusOK, views.ReviewPage(user, &game, formValues, true, v))
}

func (app *application) createReviewHandler(responseWriter http.ResponseWriter, request *http.Request) {
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
			app.notFoundResponse(responseWriter, request)
			return
		default:
			app.serverErrorResponse(responseWriter, request, err)
			return
		}
	}

	if err := app.models.UserGames.HasGame(user.ID, gameID); err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(responseWriter, request)
			return
		default:
			app.serverErrorResponse(responseWriter, request, err)
			return
		}
	}

	v := validator.New()

	review, err := app.models.Reviews.Get(user.ID, gameID)
	isEdit := true

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			isEdit = false

		default:
			app.serverErrorResponse(responseWriter, request, err)
			return
		}
	}

	formValues, err := app.parseAndValidateReviewForm(responseWriter, request, v)
	if err != nil {
		switch {
		case errors.Is(err, ErrBodyTooLarge):
			app.renderTemplate(responseWriter, request, http.StatusUnprocessableEntity, views.ReviewPage(user, &game, views.ReviewFormValues{Body: review.Body, Rating: review.Rating}, isEdit, v))
			return
		default:
			app.badRequest(responseWriter, request)
			return
		}
	}

	if !v.Valid() {
		app.renderTemplate(responseWriter, request, http.StatusUnprocessableEntity, views.ReviewPage(user, &game, formValues, isEdit, v))
		return
	}

	review = data.NewReview(user.ID, gameID, formValues.Body, formValues.Rating)
	if err := app.models.Reviews.Insert(review); err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(responseWriter, request)
			return
		default:
			app.logError(request, err)
			v.AddError("form", "We couldn't save your review, please try again")
			app.renderTemplate(responseWriter, request, http.StatusInternalServerError, views.ReviewPage(user, &game, formValues, isEdit, v))
			return
		}
	}

	http.Redirect(responseWriter, request, "/my_reviews", http.StatusSeeOther)
}

func (app *application) myReviewsPageHandler(responseWriter http.ResponseWriter, request *http.Request) {
	user, ok := app.getAuthenticatedUser(request)
	if !ok {
		http.Redirect(responseWriter, request, "/login", http.StatusSeeOther)
		return
	}

	reviews, err := app.models.Reviews.GetReviewsForUser(user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(responseWriter, request)
			return
		default:
			app.serverErrorResponse(responseWriter, request, err)
			return
		}
	}

	app.renderTemplate(responseWriter, request, http.StatusOK, views.MyReviewsPage(user, reviews))
}
