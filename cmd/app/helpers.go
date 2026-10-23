package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/julienschmidt/httprouter"
	"github.com/luism2302/RateLog/internal/validator"
	"github.com/luism2302/RateLog/internal/views"
)

func (app *application) renderTemplate(responseWriter http.ResponseWriter, request *http.Request, status int, component templ.Component) {
	var buffer bytes.Buffer
	if err := component.Render(request.Context(), &buffer); err != nil {
		app.renderError(responseWriter, request, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		app.logError(request, err)
		return
	}

	responseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
	responseWriter.WriteHeader(status)

	_, err := io.Copy(responseWriter, &buffer)
	if err != nil {
		app.logError(request, err)
		return
	}
}

func (app *application) renderError(responseWriter http.ResponseWriter, request *http.Request, status int, msg string) {
	var buffer bytes.Buffer
	if err := views.ErrTemplate(msg).Render(request.Context(), &buffer); err != nil {
		http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		app.logError(request, err)
		return
	}

	responseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
	responseWriter.WriteHeader(status)

	_, err := io.Copy(responseWriter, &buffer)
	if err != nil {
		app.logError(request, err)
		return
	}
}

func (app *application) readTitleQueryParam(request *http.Request) string {
	return request.URL.Query().Get("title")
}

func (app *application) readGameIDParam(request *http.Request) (int, error) {
	params := httprouter.ParamsFromContext(request.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		return 0, err
	}

	if id <= 0 {
		return 0, errors.New("ID must be greater than zero")
	}
	return id, nil
}

func (app *application) parseAndValidateLoginForm(request *http.Request, v *validator.Validator) (views.LoginFormValues, error) {
	if err := request.ParseForm(); err != nil {
		return views.LoginFormValues{}, err
	}

	email := cleanEmail(request.PostForm.Get("email"))
	password := request.PostForm.Get("password")

	formValues := views.NewLoginFormValues(email, password)
	views.ValidateLoginForm(formValues, v)

	return formValues, nil
}

func (app *application) parseAndValidateSignupForm(request *http.Request, v *validator.Validator) (views.SignupFormValues, error) {
	if err := request.ParseForm(); err != nil {
		return views.SignupFormValues{}, err
	}

	email := cleanEmail(request.PostForm.Get("email"))
	name := request.PostForm.Get("name")
	password := request.PostForm.Get("password")
	confirmPassword := request.PostForm.Get("confirmPassword")

	formValues := views.NewSignupFormValues(email, name, password, confirmPassword)
	views.ValidateSignupForm(formValues, v)

	return formValues, nil
}

func (app *application) ValidateGameID(id int, v *validator.Validator) {
	if id <= 0 {
		v.AddError("id", "Must be greater than zero")
	}
}

func (app *application) parseAndValidateReviewForm(request *http.Request, v *validator.Validator) (views.ReviewFormValues, error) {
	if err := request.ParseForm(); err != nil {
		return views.ReviewFormValues{}, err
	}

	body := request.PostForm.Get("body")
	rating, err := strconv.Atoi(request.PostForm.Get("rating"))
	if err != nil {
		rating = -1
	}

	formValues := views.NewReviewFormValues(body, rating)
	views.ValidateReviewGameForm(formValues, v)

	return formValues, nil
}

func cleanEmail(email string) string {
	trimmed := strings.TrimSpace(email)
	return strings.ToLower(trimmed)
}
