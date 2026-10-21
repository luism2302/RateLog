package main

import (
	"net/http"

	"github.com/luism2302/RateLog/internal/views"
)

func (app *application) logError(request *http.Request, err error) {
	app.logger.Error("system error", "method", request.Method, "URI", request.URL.RequestURI(), "err", err.Error())
}

func (app application) serverErrorResponse(responseWriter http.ResponseWriter, request *http.Request, err error) {
	app.logError(request, err)
	responseWriter.WriteHeader(http.StatusInternalServerError)
	app.renderTemplate(responseWriter, views.ErrTemplate("Internal Server Error"))
}

func (app application) invalidResponse(responseWriter http.ResponseWriter, request *http.Request, err error) {
	app.logError(request, err)
	responseWriter.WriteHeader(http.StatusUnprocessableEntity)
	app.renderTemplate(responseWriter, views.ErrTemplate("Invalid email or password"))
}

func (app application) notFoundResponse(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.WriteHeader(http.StatusNotFound)
	app.renderTemplate(responseWriter, views.ErrTemplate("Not Found"))
}
