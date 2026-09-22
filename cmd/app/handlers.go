package main

import (
	"net/http"

	"github.com/luism2302/RateLog/components"
)

func (app *application) home(responseWriter http.ResponseWriter, request *http.Request) {
	if err := renderTemplate(responseWriter, components.Base()); err != nil {
		app.serverErrorResponse(responseWriter, request, err)
		return
	}
}
