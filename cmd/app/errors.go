package main

import (
	"net/http"

	"github.com/luism2302/RateLog/components"
)

func (app *application) logError(request *http.Request, err error) {
	app.logger.Error("system error", "method", request.Method, "URI", request.URL.RequestURI(), "err", err.Error())
}

func (app application) serverErrorResponse(responseWriter http.ResponseWriter, request *http.Request, err error) {
	app.logError(request, err)
	responseWriter.WriteHeader(http.StatusInternalServerError)
	app.renderTemplate(responseWriter, components.ErrTemplate("Internal Server Error"))
}

func (app application) notFoundResponse(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.WriteHeader(http.StatusNotFound)
	app.renderTemplate(responseWriter, components.ErrTemplate("Not Found"))
}
