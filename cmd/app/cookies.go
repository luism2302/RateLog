package main

import (
	"net/http"
	"time"

	"github.com/luism2302/RateLog/internal/data"
)

const (
	authCookieName = "auth_token"
)

func setAuthCookie(responseWriter http.ResponseWriter, token data.Token) {
	cookie := http.Cookie{
		Path:     "/",
		Name:     authCookieName,
		Value:    token.Plaintext,
		Expires:  token.Expires,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(responseWriter, &cookie)
}

func removeAuthCookie(responseWriter http.ResponseWriter) {
	cookie := http.Cookie{
		Path:     "/",
		Name:     authCookieName,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(responseWriter, &cookie)
}
