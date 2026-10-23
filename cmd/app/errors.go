package main

import (
	"net/http"
)

func (app *application) logError(request *http.Request, err error) {
	app.logger.Error("system error", "method", request.Method, "URI", request.URL.Path, "err", err.Error())
}

func (app *application) serverErrorResponse(responseWriter http.ResponseWriter, request *http.Request, err error) {
	app.logError(request, err)
	app.renderError(responseWriter, request, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func (app *application) notFoundResponse(responseWriter http.ResponseWriter, request *http.Request) {
	app.renderError(responseWriter, request, http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

func (app *application) rateLimitExceededResponse(responseWriter http.ResponseWriter, request *http.Request) {
	app.renderError(responseWriter, request, http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests))
}

func (app *application) unprocessableEntityResponse(responseWriter http.ResponseWriter, request *http.Request) {
	app.renderError(responseWriter, request, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
}
