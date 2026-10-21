package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/luism2302/RateLog/internal/data"
	"github.com/luism2302/RateLog/internal/validator"
	"github.com/luism2302/RateLog/internal/views"
)

func (app *application) indexLoginHandler(responseWriter http.ResponseWriter, request *http.Request) {
	v := validator.New()
	err := app.renderTemplate(responseWriter, views.LoginPage(views.LoginFormValues{}, v))
	if err != nil {
		app.serverErrorResponse(responseWriter, request, err)
		return
	}
}

func (app *application) loginUserHandler(responseWriter http.ResponseWriter, request *http.Request) {
	v := validator.New()
	formValues, err := app.parseAndValidateLoginForm(request, v)
	if err != nil {
		app.serverErrorResponse(responseWriter, request, err)
		return
	}

	if !v.Valid() {
		app.renderTemplate(responseWriter, views.LoginPage(formValues, v))
		return
	}

	// authenticate user
	user, err := app.models.Users.GetUserByEmail(formValues.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			user.Password.Matches(formValues.Password)
			app.invalidResponse(responseWriter, request, err)
			return
		default:
			app.serverErrorResponse(responseWriter, request, err)
			return
		}
	}

	matches, err := user.Password.Matches(formValues.Password)
	if err != nil {
		app.serverErrorResponse(responseWriter, request, err)
		return
	}

	if !matches {
		v.AddError("password", "Incorrect Password")
		app.renderTemplate(responseWriter, views.LoginPage(formValues, v))
		return
	}

	// token
	token, err := app.models.Tokens.NewToken(user.ID, 24*time.Hour)
	if err != nil {
		app.serverErrorResponse(responseWriter, request, err)
		return
	}

	setAuthCookie(responseWriter, token)
	http.Redirect(responseWriter, request, "/search_game", http.StatusSeeOther)
	return
}

func (app *application) indexSignUpHandler(responseWriter http.ResponseWriter, request *http.Request) {
	v := validator.New()
	err := app.renderTemplate(responseWriter, views.SignupPage(views.SignupFormValues{}, v))
	if err != nil {
		app.serverErrorResponse(responseWriter, request, err)
		return
	}
}

func (app *application) createUserHandler(responseWriter http.ResponseWriter, request *http.Request) {
	v := validator.New()
	formValues, err := app.parseAndValidateSignupForm(request, v)
	if err != nil {
		app.serverErrorResponse(responseWriter, request, err)
		return
	}

	if !v.Valid() {
		app.renderTemplate(responseWriter, views.SignupPage(formValues, v))
		return
	}

	user, err := data.NewUser(formValues.Email, formValues.Name, formValues.Password)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrPasswordHashNil):
			panic("missing password hash for user")
		default:
			app.serverErrorResponse(responseWriter, request, err)
			return
		}
	}

	if err := app.models.Users.Insert(user); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "Email already in use")
			app.renderTemplate(responseWriter, views.SignupPage(formValues, v))
			return
		default:
			app.serverErrorResponse(responseWriter, request, err)
			return
		}
	}
	http.Redirect(responseWriter, request, "/login", http.StatusSeeOther)
	return
}

func (app *application) logoutUserHandler(responseWriter http.ResponseWriter, request *http.Request) {
	authCookie, err := request.Cookie(authCookieName)
	if err != nil {
		app.serverErrorResponse(responseWriter, request, err)
		return
	}

	token := authCookie.Value
	removeAuthCookie(responseWriter)

	err = app.models.Tokens.Delete(token)
	if err != nil {
		app.serverErrorResponse(responseWriter, request, err)
		return
	}
	http.Redirect(responseWriter, request, "/login", http.StatusSeeOther)
	return
}
