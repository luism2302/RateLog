package main

import (
	"context"
	"net/http"

	"github.com/luism2302/RateLog/internal/data"
)

type contextKey string

const authenticatedUserKey = contextKey("authenticatedUser")

func (app *application) setAuthenticatedUser(request *http.Request, user data.User) *http.Request {
	ctx := context.WithValue(request.Context(), authenticatedUserKey, user)
	return request.WithContext(ctx)
}

func (app *application) getAuthenticatedUser(request *http.Request) (*data.User, bool) {
	user, ok := request.Context().Value(authenticatedUserKey).(data.User)
	return &user, ok
}
