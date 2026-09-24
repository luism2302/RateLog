package main

import (
	"fmt"
	"net/http"

	"github.com/luism2302/RateLog/components"
)

func (app *application) home(responseWriter http.ResponseWriter, request *http.Request) {
	if err := app.renderTemplate(responseWriter, components.Base()); err != nil {
		app.serverErrorResponse(responseWriter, request, err)
		return
	}
}

func (app *application) searchGameHandler(responseWriter http.ResponseWriter, request *http.Request) {
	title := app.readTitleParam(request)
	games, err := app.client.SearchGameByTitle(title)
	if err != nil {
		app.serverErrorResponse(responseWriter, request, err)
		return
	}

	fmt.Print(games)
}
