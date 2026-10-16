package main

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/luism2302/RateLog/internal/validator"
	"github.com/luism2302/RateLog/internal/views"
)

func (app *application) renderTemplate(responseWriter http.ResponseWriter, template templ.Component) error {
	return template.Render(context.Background(), responseWriter)
}

func (app *application) readTitleQueryParam(request *http.Request) string {
	return request.URL.Query().Get("title")
}

func (app *application) parseAndValidateLoginForm(request *http.Request, v *validator.Validator) (views.LoginFormValues, error) {
	if err := request.ParseForm(); err != nil {
		return views.LoginFormValues{}, err
	}

	email := request.PostForm.Get("email")
	password := request.PostForm.Get("password")

	formValues := views.NewLoginFormValues(email, password)
	views.ValidateLoginForm(formValues, v)

	return formValues, nil
}

func (app *application) parseAndValidateSignupForm(request *http.Request, v *validator.Validator) (views.SignupFormValues, error) {
	if err := request.ParseForm(); err != nil {
		return views.SignupFormValues{}, err
	}

	email := request.PostForm.Get("email")
	name := request.PostForm.Get("name")
	password := request.PostForm.Get("password")
	confirmPassword := request.PostForm.Get("confirmPassword")

	formValues := views.NewSignupFormValues(email, name, password, confirmPassword)
	views.ValidateSignupForm(formValues, v)

	return formValues, nil
}

func (app *application) parseAndValidateAddGameForm(request *http.Request, v *validator.Validator) (views.AddGameFormValues, error) {
	if err := request.ParseForm(); err != nil {
		return views.AddGameFormValues{}, err
	}

	id, err := strconv.Atoi(request.PostForm.Get("id"))
	if err != nil {
		return views.AddGameFormValues{}, err
	}

	name := request.PostForm.Get("name")
	released := request.PostForm.Get("released")
	backgroundImage := request.PostForm.Get("backgroundImage")

	added, err := strconv.Atoi(request.PostForm.Get("added"))
	if err != nil {
		return views.AddGameFormValues{}, err
	}

	formValues := views.NewAddGameFormValues(id, added, name, released, backgroundImage)
	views.ValidateAddGameForm(formValues, v)

	parsedDate, err := time.Parse("2006-01-02 15:04:05 -0700 MST", formValues.Released)
	if err != nil {
		return views.AddGameFormValues{}, err
	}

	formValues.Released = parsedDate.Format(time.DateOnly)

	return formValues, nil
}
