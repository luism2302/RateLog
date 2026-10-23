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
		app.serverErrorResponse(responseWriter, request, err)
		return
	}

	game, err := app.models.Games.GetGameByID(gameID)
	if err != nil {
		app.serverErrorResponse(responseWriter, request, err)
		return
	}
	v := validator.New()
	app.renderTemplate(responseWriter, request, http.StatusOK, views.ReviewPage(user, &game, views.ReviewFormValues{}, v))
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

	v := validator.New()
	formValues, err := app.parseAndValidateReviewForm(request, v)
	if err != nil {
		app.serverErrorResponse(responseWriter, request, err)
		return
	}

	if !v.Valid() {
		app.renderTemplate(responseWriter, request, http.StatusUnprocessableEntity, views.ReviewPage(user, &game, formValues, v))
		return
	}

	review := data.NewReview(user.ID, gameID, formValues.Body, formValues.Rating)
	if err := app.models.Reviews.Insert(review); err != nil {
		app.serverErrorResponse(responseWriter, request, err)
		return
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
