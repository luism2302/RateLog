package main

import (
	"errors"
	"net/http"

	"github.com/luism2302/RateLog/internal/data"
	"github.com/luism2302/RateLog/internal/validator"
)

func (app *application) authenticate(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		authCookie, err := request.Cookie(authCookieName)
		if err != nil {
			http.Redirect(responseWriter, request, "/login", http.StatusSeeOther)
			return
		}

		token := authCookie.Value
		v := validator.New()

		if data.ValidateToken(token, v); !v.Valid() {
			app.serverErrorResponse(responseWriter, request, errors.New("Token Validation Failed"))
			return
		}

		user, err := app.models.Users.GetUserForToken(token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.serverErrorResponse(responseWriter, request, err)
				return
			default:
				app.serverErrorResponse(responseWriter, request, errors.New("User Validation Failed"))
				return
			}
		}
		request = app.setAuthenticatedUser(request, user)
		next.ServeHTTP(responseWriter, request)
	})
}
